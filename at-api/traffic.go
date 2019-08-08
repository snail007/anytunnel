package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/snail007/go-activerecord/mysql"
)

type UserTraffic struct {
	ClusterID   uint64
	UserID      string
	ServerToken string
	TunnelID    uint64
	ConnCount   uint64
	Upload      uint64
	Download    uint64
	Bytes       uint64
}

var trafficChn = make(chan UserTraffic, 50000)

func init() {
	go func() {
		db := G.DB()
		for {
			traffic := <-trafficChn
			if traffic.Bytes > 0 {
				//update user traffic
				now := time.Now().Unix()
				sql := db.AR().Update("package", map[string]interface{}{
					"bytes_left -": traffic.Bytes,
					"update_time":  now,
				}, map[string]interface{}{
					"user_id":       traffic.UserID,
					"start_time <=": now,
					"end_time >=":   now,
					"bytes_left >":  0,
				}).OrderBy("end_time", "ASC").Limit(1)
				_, err := db.Exec(sql)
				if err != nil {
					log.Warnf("traffic counter ERR:%s", err)
				}
			}
			if traffic.ConnCount > 0 {
				//update tunnel conn count
				rs, err := db.Query(db.AR().From("server").Where(map[string]interface{}{
					"token": traffic.ServerToken,
				}).Limit(0, 1))
				if err != nil {
					log.Warnf("traffic counter ERR:%s", err)
				}
				if rs.Len() == 0 {
					log.Warnf("traffic counter ERR:empty sever for token:%s", traffic.ServerToken)
				}

				_, err = db.Exec(db.AR().Replace("conn", map[string]interface{}{
					"user_id":     traffic.UserID,
					"cluster_id":  traffic.ClusterID,
					"server_id":   rs.Value("server_id"),
					"tunnel_id":   traffic.TunnelID,
					"count":       traffic.ConnCount,
					"upload":      traffic.Upload,
					"download":    traffic.Download,
					"update_time": time.Now().Unix(),
				}))
				if err != nil {
					log.Warnf("traffic counter ERR:%s", err)
				}
			}
		}
	}()
}

type Traffic struct{}

func NewTraffic() *Traffic {
	return &Traffic{}
}

//Traffic update a tunnel traffic bytes used,对应cluster的:url.traffic
//method : POST
//params : json {"1":{"positive":3231,"negative":5324,"connCount":32."serverToken":"0840d2i30ofs"}}
func (tunnel *Traffic) Traffic(responseWrite http.ResponseWriter, request *http.Request, params httprouter.Params) {
	defer request.Body.Close()
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		jsonError(responseWrite, err, nil)
		return
	}
	jsonStr := string(body)
	dataMap := map[string]map[string]interface{}{}
	err = json.Unmarshal([]byte(jsonStr), &dataMap)
	if err != nil {
		log.Warnf("update traffic ERR:%s", err)
		fmt.Fprint(responseWrite, err.Error())
		return
	}
	if len(dataMap) == 0 {
		log.Warnf("update traffic ERR:empty data")
		fmt.Fprint(responseWrite, err.Error())
		return
	}
	tunnelIDs := []string{}
	for k := range dataMap {
		tunnelIDs = append(tunnelIDs, k)
	}
	db := G.DB()
	var rs *mysql.ResultSet
	addr := request.RemoteAddr
	clusterIP := addr[0:strings.Index(addr, ":")]
	rs, err = db.Query(db.AR().
		Select("cluster_id").
		From("cluster").
		Where(map[string]interface{}{
			"ip":        clusterIP,
			"is_delete": 0,
		}).Limit(1))
	if err != nil {
		log.Warnf("update traffic ERR:%s", err)
		fmt.Fprint(responseWrite, err.Error())
		return
	}
	if rs.Len() == 0 {
		err = fmt.Errorf("update traffic ERR:empty cluster for ip :%s", clusterIP)
		log.Warnf("%s", err)
		fmt.Fprint(responseWrite, err.Error())
		return
	}
	_clusterID := rs.Value("cluster_id")
	clusterID, _ := strconv.ParseUint(_clusterID, 10, 64)
	rs, err = db.Query(db.AR().
		Select("user_id,tunnel_id").
		From("tunnel").
		Where(map[string]interface{}{
			"tunnel_id": tunnelIDs,
			"is_delete": 0,
		}))
	if err != nil {
		log.Warnf("update traffic ERR:%s", err)
		fmt.Fprint(responseWrite, err.Error())
		return
	}
	userIds := rs.MapValues("tunnel_id", "user_id")
	_data := []UserTraffic{}
	data := []UserTraffic{}
	for k, v := range dataMap {
		userID, ok := userIds[k]
		if !ok {
			continue
		}
		nVal, err1 := strconv.ParseUint(fmt.Sprintf("%.0f", v["negative"].(float64)), 10, 64)
		pVal, err2 := strconv.ParseUint(fmt.Sprintf("%.0f", v["positive"].(float64)), 10, 64)
		connCount, err3 := strconv.ParseUint(fmt.Sprintf("%.0f", v["connCount"].(float64)), 10, 64)
		interval, err4 := strconv.ParseUint(fmt.Sprintf("%.0f", v["interval"].(float64)), 10, 64)
		if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
			log.Warnf("update traffic ERR:%s,%s,%s", err1, err2, err3)
			fmt.Fprint(responseWrite, "data format incorrect")
			return
		}
		tunnelID, _ := strconv.ParseUint(k, 10, 64)
		_data = append(_data, UserTraffic{
			UserID:      userID,
			Bytes:       uint64(nVal + pVal),
			ConnCount:   connCount,
			ServerToken: v["serverToken"].(string),
			TunnelID:    tunnelID,
			ClusterID:   clusterID,
			Upload:      pVal / interval,
			Download:    nVal / interval,
		})
	}

	for _, v := range _data {
		found := false
		for k1, v1 := range data {
			if v.UserID == v1.UserID {
				found = true
				v1.Bytes = v.Bytes + v1.Bytes
				data[k1] = v1
				break
			}
		}
		if !found {
			data = append(data, v)
		}
	}
	if len(data) > 0 {
		for _, v := range data {
			select {
			case trafficChn <- v:
			default:
				log.Warnf("traffic counter queue is full")
			}
		}
	}
	responseWrite.WriteHeader(http.StatusNoContent)
}
