package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

//1.Cluster认证CS调用的接口
//2.server和client登录获取cluster的接口
//3.Cluster上报server和client上线下线的接口
type ClusterAPI struct{}

func NewClusterAPI() ClusterAPI {
	return ClusterAPI{}
}

//Auth cs by cs_token,对应cluster的:url.auth
//method : GET
//params : token type
func (cs *ClusterAPI) Auth(responseWrite http.ResponseWriter, request *http.Request, params httprouter.Params) {
	//1.参数校验
	ip := request.FormValue("ip")
	if ip == "" {
		responseWrite.Write([]byte("ip required"))
		return
	}
	token := request.FormValue("token")
	if token == "" {
		responseWrite.Write([]byte("token required"))
		return
	}
	csType := request.FormValue("type")
	if csType == "" {
		responseWrite.Write([]byte("type required"))
		return
	}
	if csType != "server" && csType != "client" {
		responseWrite.Write([]byte("type error"))
		return
	}
	//2.检查server或者client是否存在
	userID, err := authCheckCS(csType, token)
	if err != nil {
		responseWrite.Write([]byte(err.Error()))
		return
	}
	//3.当server或者client的属主不是系统用户时,进行权限检查
	if userID != "0" {
		// //4.检查对应的user是否被禁止
		// err = authCheckUser(userID)
		// if err != nil {
		// 	responseWrite.Write([]byte(err.Error()))
		// 	return
		// }
		// //5.检查对应的user流量是否用完
		// err = authCheckUserTraffic(userID)
		// if err != nil {
		// 	responseWrite.Write([]byte(err.Error()))
		// 	return
		// }
		// //6.0 IP白名单检查
		// ok, err := authCheckIPWhiteList(ip, csType)
		// if err != nil {
		// 	responseWrite.Write([]byte(err.Error()))
		// 	return
		// }
		// q := qqwry.Find(ip)
		// q.IP = ip
		// //IP白名单检查没有通过,继续进行区域白名单检查
		// if !ok {
		// 	ok, err := authCheckAreaWhiteList(q.Country, csType)
		// 	if err != nil {
		// 		responseWrite.Write([]byte(err.Error()))
		// 		return
		// 	}
		// 	//区域白名单检查没有通过,继续进行IP黑名单检查
		// 	if !ok {
		// 		ok, err = authCheckIPBlackList(ip, csType)
		// 		if err != nil {
		// 			responseWrite.Write([]byte(err.Error()))
		// 			return
		// 		}
		// 		//IP黑名单检查通过,继续进行区域黑名单检查
		// 		if ok {
		// 			ok, err = authCheckAreaBlackList(q.Country, csType)
		// 			if err != nil {
		// 				responseWrite.Write([]byte(err.Error()))
		// 				return
		// 			}
		// 			//区域黑名单检查通过,继续进行角色区域检查
		// 			if ok {
		// 				//取出用户所有角色区域对应登录设置
		// 				err = authCheckRoleArea(userID, q, csType)
		// 				if err != nil {
		// 					responseWrite.Write([]byte(err.Error()))
		// 					return
		// 				}
		// 			} else {
		// 				//区域黑名单检查没有通过
		// 				responseWrite.Write([]byte(fmt.Sprintf("Your IP Area %s was bocked", q.Country)))
		// 				return
		// 			}
		// 		} else {
		// 			//IP黑名单检查没有通过
		// 			responseWrite.Write([]byte(fmt.Sprintf("Your IP %s was blocked", ip)))
		// 			return
		// 		}
		// 	}
		// }
	}
	responseWrite.WriteHeader(http.StatusNoContent)
}

//Status update online/offline Status cs by cs_token,对应cluster的:url.status
//method : GET
//params : token  type  action
func (cs *ClusterAPI) Status(responseWrite http.ResponseWriter, request *http.Request, params httprouter.Params) {

	token := request.FormValue("token")
	if token == "" {
		jsonError(responseWrite, "token required!", nil)
		return
	}
	csIP := request.FormValue("ip")
	if token == "" {
		jsonError(responseWrite, "ip required!", nil)
		return
	}
	csType := request.FormValue("type")
	if csType == "" {
		jsonError(responseWrite, "type required!", nil)
		return
	}
	action := request.FormValue("action")
	if action == "" {
		jsonError(responseWrite, "action required!", nil)
		return
	}
	if csType != "server" && csType != "client" {
		jsonError(responseWrite, "type error!", nil)
		return
	}
	if action != "online" && action != "offline" {
		jsonError(responseWrite, "action error!", nil)
		return
	}
	db := G.DB()
	id := ""
	rs, err := db.Query(db.AR().From(csType).Where(map[string]interface{}{
		"token": token,
	}).Limit(0, 1))
	if err != nil {
		jsonError(responseWrite, err.Error(), nil)
		return
	}
	if rs.Len() == 0 {
		err = fmt.Errorf("empty "+csType+" for token:%s", token)
		jsonError(responseWrite, err.Error(), nil)
		return
	}
	id = rs.Value(csType + "_id")
	userID := rs.Value("user_id")

	clusterID := ""
	clusterIP := request.RemoteAddr[0:strings.Index(request.RemoteAddr, ":")]
	rs, err = db.Query(db.AR().From("cluster").Where(map[string]interface{}{
		"ip":        clusterIP,
		"is_delete": 0,
	}).Limit(0, 1))
	if err != nil {
		jsonError(responseWrite, err.Error(), nil)
		return
	}
	if rs.Len() == 0 {
		err = fmt.Errorf("empty cluster")
		jsonError(responseWrite, err.Error(), nil)
		return
	}
	clusterID = rs.Value("cluster_id")

	where := map[string]interface{}{}
	typ := ""
	where["cs_id"] = id
	if csType == "server" {
		typ = "server"
	}
	if csType == "client" {
		typ = "client"
	}

	where["cs_type"] = typ
	data := map[string]interface{}{
		"cluster_id":  clusterID,
		"cs_ip":       csIP,
		"cs_id":       id,
		"cs_type":     typ,
		"user_id":     userID,
		"create_time": time.Now().Unix(),
	}

	if action == "online" {
		_, err := db.Exec(db.AR().Replace("online", data))
		if err != nil {
			jsonError(responseWrite, err.Error(), nil)
			return
		}
	} else {
		_, err := db.Exec(db.AR().Delete("online", where))
		if err != nil {
			jsonError(responseWrite, err.Error(), nil)
			return
		}
		if csType == "server" {
			db.Exec(db.AR().Delete("conn", map[string]interface{}{
				"server_id": id,
			}))
		}
	}
	responseWrite.WriteHeader(http.StatusNoContent)
}

//Cluster get a Cluster ip for c/s
//method : GET
//params : token  type
func (cs *ClusterAPI) Cluster(responseWrite http.ResponseWriter, request *http.Request, params httprouter.Params) {

	token := request.FormValue("token")
	if token == "" {
		responseWrite.Write([]byte("token required"))
		return
	}
	csType := request.FormValue("type")
	if csType == "" {
		responseWrite.Write([]byte("type required"))
		return
	}
	if csType != "server" && csType != "client" {
		responseWrite.Write([]byte("type error"))
		return
	}
	ip, err := getClusterIPByToken(csType, token)
	if err != nil {
		responseWrite.WriteHeader(http.StatusInternalServerError)
		responseWrite.Write([]byte(err.Error()))
		return
	}
	responseWrite.Write([]byte(ip))
}
func getClusterIPByToken(typ, token string) (clusterIP string, err error) {
	from := "server"
	col := "server_id"
	if typ == "client" {
		from = "client"
		col = "client_id"
	}
	db := G.DB()

	rs, err := db.Query(db.AR().From(from).Where(map[string]interface{}{
		"is_delete": 0,
		"token":     token,
	}).Limit(0, 1))
	if err != nil {
		return
	}
	if rs.Len() == 0 {
		err = fmt.Errorf("token error")
		return
	}
	colID := rs.Value(col)

	rs, err = db.Query(db.AR().From("tunnel").Where(map[string]interface{}{
		"is_delete": 0,
		col:         colID,
	}).Limit(0, 1))
	if err != nil {
		return
	}
	if rs.Len() == 0 {
		err = fmt.Errorf("empty tunnel")
		return
	}
	tunnelClusterID := rs.Value("cluster_id")

	rs, err = db.Query(db.AR().From("cluster").Where(map[string]interface{}{
		"is_delete":  0,
		"cluster_id": tunnelClusterID,
		"is_disable": 0,
	}).Limit(0, 1))
	if err != nil {
		return
	}
	if rs.Len() == 0 {
		err = fmt.Errorf("empty cluster")
		return
	}
	clusterIP = rs.Value("ip")
	return
}
