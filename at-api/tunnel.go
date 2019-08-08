package main

import (
	"anytunnel/at-common"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/snail007/go-activerecord/mysql"
	"fmt"
)

type Tunnel struct{}

func NewTunnel() *Tunnel {
	return &Tunnel{}
}

//add a tunnel
//method : POST
//params :  mode, name, user_id, cluster_id, server_id,client_id, protocol,
//          server_listen_port, server_listen_ip, client_local_port,client_local_host
func (tunnel *Tunnel) Add(responseWrite http.ResponseWriter, request *http.Request, params httprouter.Params) {

	mode := request.PostFormValue("mode")
	name := request.PostFormValue("name")
	userId := request.PostFormValue("user_id")
	clusterId := request.PostFormValue("cluster_id")
	serverId := request.PostFormValue("server_id")
	clientId := request.PostFormValue("client_id")
	protocol := request.PostFormValue("protocol")
	serverListenPort := request.PostFormValue("server_listen_port")
	serverListenIp := request.PostFormValue("server_listen_ip")
	clientLocalPort := request.PostFormValue("client_local_port")
	clientLocalHost := request.PostFormValue("client_local_host")

	if mode == "" {
		jsonError(responseWrite, "必须选择一种模式!", nil)
		return
	}
	if(mode != "0" && mode != "1" && mode != "2") {
		jsonError(responseWrite, "模式选择错误!", nil)
		return
	}
	if name == "" {
		jsonError(responseWrite, "隧道名称不能为空!", nil)
		return
	}
	if userId == "" {
		jsonError(responseWrite, "user_id 不能为空!", nil)
		return
	}
	if clusterId == "" {
		jsonError(responseWrite, "没有选择节点!", nil)
		return
	}
	if serverId == "" {
		jsonError(responseWrite, "没有选择server!", nil)
		return
	}
	if clientId == "" {
		jsonError(responseWrite, "没有选择client!", nil)
		return
	}
	if protocol == "" {
		jsonError(responseWrite, "没有选择协议!", nil)
		return
	}
	if serverListenPort == "" {
		jsonError(responseWrite, "server监听端口不能为空!", nil)
		return
	}
	if serverListenIp == "" {
		jsonError(responseWrite, "server绑定Ip不能为空!", nil)
		return
	}
	if clientLocalPort == "" {
		jsonError(responseWrite, "client监听端口不能为空!", nil)
		return
	}
	if clientLocalHost == "" {
		jsonError(responseWrite, "client绑定Ip不能为空!", nil)
		return
	}

	db := G.DB()
	var rs *mysql.ResultSet

	//cluster存在
	rs, err := db.Query(db.AR().From("cluster").Where(map[string]interface{}{
		"cluster_id": clusterId,
		"is_delete": 0,
	}).Limit(0, 1))
	if err != nil {
		jsonError(responseWrite, err.Error(), nil)
		return
	}
	if(rs.Len() == 0) {
		jsonError(responseWrite, "cluster不存在!", nil)
	}

	//基础模式
	//  1.server 必须为系统server
	//  2.server 必须在所选的cluster在线
	//  3.client 存在
	if(mode == "0") {
		//必须为系统server
		rs, err := db.Query(db.AR().From("server").Where(map[string]interface{}{
			"server_id": serverId,
			"user_id": 0,
			"is_delete": 0,
		}).Limit(0, 1))
		if err != nil {
			jsonError(responseWrite, err.Error(), nil)
			return
		}
		if(rs.Len() == 0) {
			jsonError(responseWrite, "基础模式下必须为系统server", nil)
		}
		//在cluster上且在线
		rs, err = db.Query(db.AR().From("online").Where(map[string]interface{}{
			"cs_id": serverId,
			"cs_type": "server",
			"cluster_id": clusterId,
		}).Limit(0, 1))
		if err != nil {
			jsonError(responseWrite, err.Error(), nil)
			return
		}
		if(rs.Len() == 0) {
			jsonError(responseWrite, "基础模式下server必须在线", nil)
		}
		//client存在
		rs, err = db.Query(db.AR().From("client").Where(map[string]interface{}{
			"client_id": clientId,
			"user_id": userId,
			"is_delete": 0,
		}).Limit(0, 1))
		if err != nil {
			jsonError(responseWrite, err.Error(), nil)
			return
		}
		if(rs.Len() == 0) {
			jsonError(responseWrite, "client不存在", nil)
		}
	}
	//高级模式
	//  1.client存在
	//  2.server存在
	if(mode == "1") {
		//server存在
		rs, err := db.Query(db.AR().From("server").Where(map[string]interface{}{
			"server_id": serverId,
			"user_id": userId,
			"is_delete": 0,
		}).Limit(0, 1))
		if err != nil {
			jsonError(responseWrite, err.Error(), nil)
			return
		}
		if(rs.Len() == 0) {
			jsonError(responseWrite, "server不存在!", nil)
		}
		//client存在
		rs, err = db.Query(db.AR().From("client").Where(map[string]interface{}{
			"client_id": clientId,
			"user_id": userId,
			"is_delete": 0,
		}).Limit(0, 1))
		if err != nil {
			jsonError(responseWrite, err.Error(), nil)
			return
		}
		if(rs.Len() == 0) {
			jsonError(responseWrite, "client不存在", nil)
		}
	}
	//特殊模式
	//  1.client 必须为系统server
	//  2.client 必须在所选的cluster在线
	//  3.server 存在
	if(mode == "2") {
		//必须为系统client
		rs, err := db.Query(db.AR().From("server").Where(map[string]interface{}{
			"client_id": clientId,
			"user_id": 0,
			"is_delete": 0,
		}).Limit(0, 1))
		if err != nil {
			jsonError(responseWrite, err.Error(), nil)
			return
		}
		if(rs.Len() == 0) {
			jsonError(responseWrite, "特殊模式下必须为系统client", nil)
		}
		//在cluster上且在线
		rs, err = db.Query(db.AR().From("online").Where(map[string]interface{}{
			"cs_id": clientId,
			"cs_type": "client",
			"cluster_id": clusterId,
		}).Limit(0, 1))
		if err != nil {
			jsonError(responseWrite, err.Error(), nil)
			return
		}
		if(rs.Len() == 0) {
			jsonError(responseWrite, "特殊模式下client必须在线", nil)
		}
		//server存在
		rs, err = db.Query(db.AR().From("server").Where(map[string]interface{}{
			"server_id": clientId,
			"user_id": userId,
			"is_delete": 0,
		}).Limit(0, 1))
		if err != nil {
			jsonError(responseWrite, err.Error(), nil)
			return
		}
		if(rs.Len() == 0) {
			jsonError(responseWrite, "server不存在", nil)
		}
	}

	//server_ip 和 cluster_id

	//add tunnel
	tunnelValues := map[string]interface{}{
		"mode":                 mode,
		"name":                 name,
		"user_id":              userId,
		"cluster_id":           clusterId,
		"server_id":            serverId,
		"client_id":            clientId,
		"protocol":             protocol,
		"server_listen_port":   serverListenPort,
		"server_listen_ip":     serverListenIp,
		"client_local_port":    clientLocalPort,
		"client_local_host":    clientLocalHost,
		"create_time":          time.Now().Unix(),
		"update_time":          time.Now().Unix(),
	}

	rs, err = db.Exec(db.AR().Insert("tunnel", tunnelValues))
	if err != nil {
		jsonError(responseWrite, err.Error(), nil)
		return
	}
	id := rs.LastInsertId

	jsonSuccess(responseWrite, "添加隧道成功", id)
}

//update by tunnel
//method : POST
//params :  mode, name, user_id, cluster_id, server_id,client_id, protocol,
//          server_listen_port, server_listen_ip, client_local_port,client_local_host
func (tunnel *Tunnel) Update(responseWrite http.ResponseWriter, request *http.Request, params httprouter.Params) {

	tunnelId := request.PostFormValue("tunnel_id")
	mode := request.PostFormValue("mode")
	name := request.PostFormValue("name")
	userId := request.PostFormValue("user_id")
	clusterId := request.PostFormValue("cluster_id")
	serverId := request.PostFormValue("server_id")
	clientId := request.PostFormValue("client_id")
	protocol := request.PostFormValue("protocol")
	serverListenPort := request.PostFormValue("server_listen_port")
	serverListenIp := request.PostFormValue("server_listen_ip")
	clientLocalPort := request.PostFormValue("client_local_port")
	clientLocalHost := request.PostFormValue("client_local_host")

	if tunnelId == "" {
		jsonError(responseWrite, "tunnel_id is not empty!", nil)
		return
	}
	if mode == "" {
		jsonError(responseWrite, "mode is not empty!", nil)
		return
	}
	if(mode != "1" && mode != "0") {
		jsonError(responseWrite, "mode is error!", nil)
		return
	}
	if name == "" {
		jsonError(responseWrite, "name is not empty!", nil)
		return
	}
	if userId == "" {
		jsonError(responseWrite, "user_id is not empty!", nil)
		return
	}
	if clusterId == "" {
		jsonError(responseWrite, "cluster_id is not empty!", nil)
		return
	}
	if serverId == "" {
		jsonError(responseWrite, "server_id is not empty!", nil)
		return
	}
	if clientId == "" {
		jsonError(responseWrite, "client_id is not empty!", nil)
		return
	}
	if protocol == "" {
		jsonError(responseWrite, "protocol is not empty!", nil)
		return
	}
	if serverListenPort == "" {
		jsonError(responseWrite, "server_listen_port is not empty!", nil)
		return
	}
	if serverListenIp == "" {
		jsonError(responseWrite, "server_listen_ip is not empty!", nil)
		return
	}
	if clientLocalPort == "" {
		jsonError(responseWrite, "client_local_port is not empty!", nil)
		return
	}
	if clientLocalHost == "" {
		jsonError(responseWrite, "client_local_host is not empty!", nil)
		return
	}

	db := G.DB()
	var rs *mysql.ResultSet
	rs, err := db.Query(db.AR().From("tunnel").Where(map[string]interface{}{
		"tunnel_id": tunnelId,
		"is_delete": 0,
	}).Limit(0, 1))
	if err != nil {
		jsonError(responseWrite, err.Error(), nil)
		return
	}
	if rs.Len() == 0 {
		jsonError(responseWrite, "tunnel is not exist", nil)
		return
	}

	//tunnel values
	tunnelValues := map[string]interface{}{
		"mode":                 mode,
		"name":                 name,
		"user_id":              userId,
		"cluster_id":           clusterId,
		"server_id":            serverId,
		"client_id":            clientId,
		"protocol":             protocol,
		"server_listen_port":   serverListenPort,
		"server_listen_ip":     serverListenIp,
		"client_local_port":    clientLocalPort,
		"client_localhost":     clientLocalHost,
		"update_time":          time.Now().Unix(),
	}

	rs, err = db.Exec(db.AR().Update("tunnel", tunnelValues, map[string]interface{}{
		"tunnel_id": tunnelId,
	}))
	if err != nil {
		jsonError(responseWrite, err.Error(), nil)
		return
	}
	id := rs.LastInsertId

	jsonSuccess(responseWrite, "update tunnel success", id)
}

//delete a tunnel by tunnel_id
//method : GET
//params : tunnel_id
func (tunnel *Tunnel) Delete(responseWrite http.ResponseWriter, request *http.Request, params httprouter.Params) {

	tunnelId := request.FormValue("tunnel_id")
	if tunnelId == "" {
		jsonError(responseWrite, "tunnel_id error!", nil)
		return
	}

	db := G.DB()
	var rs *mysql.ResultSet
	rs, err := db.Query(db.AR().From("tunnel").Where(map[string]interface{}{
		"tunnel_id": tunnelId,
		"is_delete": 0,
	}).Limit(0, 1))
	if err != nil {
		jsonError(responseWrite, err.Error(), nil)
		return
	}
	if rs.Len() == 0 {
		jsonError(responseWrite, "隧道不存在", nil)
		return
	}

	//判断是否打开
	tunnelRow := rs.Row()
	if(tunnelRow["is_open"] == "1") {
		jsonError(responseWrite, "请先关闭隧道再删除", nil)
		return
	}

	tunnelValues := map[string]interface{}{
		"is_delete":   1,
		"update_time": time.Now().Unix(),
	}
	rs, err = db.Exec(db.AR().Update("tunnel", tunnelValues, map[string]interface{}{
		"tunnel_id": tunnelId,
	}))
	if err != nil {
		jsonError(responseWrite, err.Error(), nil)
		return
	}

	jsonSuccess(responseWrite, "delete tunnel success", nil)
}

//tunnel list
//method : GET
//params: keyword("") page(1) number(15) user_id
func (tunnel *Tunnel) List(responseWrite http.ResponseWriter, request *http.Request, params httprouter.Params) {

	keyword := strings.Trim(request.FormValue("keyword"), "")
	page := request.FormValue("page")
	pageSize := request.FormValue("number")
	userId := request.FormValue("user_id")

	pageNumber := 1
	number := 15
	if page != "" {
		pageNumber, _ = strconv.Atoi(page)
	}
	if pageSize != "" {
		number, _ = strconv.Atoi(pageSize)
	}

	offset := (pageNumber - 1) * number

	db := G.DB()
	var rs *mysql.ResultSet
	var err error
	tunnelRows := []map[string]string{}
	if keyword != "" {
		sqlString := "SELECT * FROM at_tunnel where is_delete=0"
		if userId != "" {
			sqlString += " AND user_id=" + userId
		}
		sqlString += " AND (" +
			"server_local_ip='" + keyword + "' OR " +
			"client_local_host='" + keyword + "'" +
			")"
		sql := db.AR().Raw(sqlString)
		rs, err = db.Query(sql)
		if err != nil {
			jsonError(responseWrite, err.Error(), nil)
			return
		}
		tunnelRows = rs.Rows()
	} else {
		sql := db.AR().From("tunnel")
		if userId != "" {
			sql = sql.Where(map[string]interface{}{"is_delete": 0, "user_id": userId})
		} else {
			sql = sql.Where(map[string]interface{}{"is_delete": 0})
		}
		sql = sql.Limit(offset, number)
		rs, err = db.Query(sql)
		if err != nil {
			jsonError(responseWrite, err.Error(), nil)
			return
		}
		tunnelRows = rs.Rows()
	}

	jsonSuccess(responseWrite, "ok", tunnelRows)
}

//get tunnel by tunnel_id
//method : GET
//params : tunnel_id
func (tunnel *Tunnel) GetTunnelByTunnelId(responseWrite http.ResponseWriter, request *http.Request, params httprouter.Params) {

	tunnelId := request.FormValue("tunnel_id")
	if tunnelId == "" {
		jsonError(responseWrite, "tunnel_id is error!", nil)
		return
	}

	db := G.DB()
	var rs *mysql.ResultSet
	rs, err := db.Query(db.AR().From("tunnel").Where(map[string]interface{}{
		"tunnel_id": tunnelId,
		"is_delete": 0,
	}).Limit(0, 1))
	if err != nil {
		jsonError(responseWrite, err.Error(), nil)
		return
	}
	if rs.Len() == 0 {
		jsonError(responseWrite, "tunnel is not exists", nil)
		return
	}

	jsonSuccess(responseWrite, "ok", rs.Row())
}

//get tunnel by client_id
//method : GET
//params : client_id
func (tunnel *Tunnel) GetTunnelsByClientId(responseWrite http.ResponseWriter, request *http.Request, params httprouter.Params) {

	clientId := request.FormValue("client_id")
	if clientId == "" {
		jsonError(responseWrite, "client_id is error!", nil)
		return
	}
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err := db.Query(db.AR().From("tunnel").Where(map[string]interface{}{
		"client_id":     clientId,
		"is_delete": 0,
	}))
	if err != nil {
		jsonError(responseWrite, err.Error(), nil)
		return
	}

	jsonSuccess(responseWrite, "ok", rs.Rows())
}

//get tunnel by server_id
//method : GET
//params : server_id
func (tunnel *Tunnel) GetTunnelsByServerId(responseWrite http.ResponseWriter, request *http.Request, params httprouter.Params) {

	serverId := request.FormValue("server_id")
	if serverId == "" {
		jsonError(responseWrite, "server_id is error!", nil)
		return
	}
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err := db.Query(db.AR().From("tunnel").Where(map[string]interface{}{
		"client_id":     serverId,
		"is_delete": 0,
	}))
	if err != nil {
		jsonError(responseWrite, err.Error(), nil)
		return
	}

	jsonSuccess(responseWrite, "ok", rs.Rows())
}

//open tunnel
//method : GET
//params : tunnel_id
func (tunnel *Tunnel) TunnelOpen(responseWrite http.ResponseWriter, request *http.Request, params httprouter.Params) {

	tunnelId := request.FormValue("tunnel_id")
	if tunnelId == "" {
		jsonError(responseWrite, "tunnel_id is error!", nil)
		return
	}
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err := db.Query(db.AR().From("tunnel").Where(map[string]interface{}{
		"tunnel_id": tunnelId,
		"is_delete": 0,
	}).Limit(0, 1))
	if err != nil {
		jsonError(responseWrite, err.Error(), nil)
		return
	}
	if rs.Len() == 0 {
		jsonError(responseWrite, "隧道不存在", nil)
	}
	tunnelValue := rs.Row()
	clientId := tunnelValue["client_id"]
	serverId := tunnelValue["server_id"]
	//clusterId := tunnelValue["cluster_id"]

	//client_id 存在
	rs, err = db.Query(db.AR().From("client").Where(map[string]interface{}{
		"client_id": clientId,
		"is_delete": 0,
	}).Limit(0, 1))
	if err != nil {
		jsonError(responseWrite, err.Error(), nil)
		return
	}
	if rs.Len() == 0 {
		jsonError(responseWrite, "client 不存在", nil)
		return
	}
	clientValue := rs.Row()

	//server_id 存在
	rs, err = db.Query(db.AR().From("server").Where(map[string]interface{}{
		"server_id": serverId,
		"is_delete": 0,
	}).Limit(0, 1))
	if err != nil {
		jsonError(responseWrite, err.Error(), nil)
		return
	}
	if rs.Len() == 0 {
		jsonError(responseWrite, "server 不存在", nil)
		return
	}
	serverValue := rs.Row()

	//client 部署在线
	rs, err = db.Query(db.AR().From("online").Where(map[string]interface{}{
		"cs_id": clientId,
		"cs_type": "client",
	}).Limit(0, 1))
	if err != nil {
		jsonError(responseWrite, err.Error(), nil)
		return
	}
	if rs.Len() == 0 {
		jsonError(responseWrite, "client 未部署", nil)
		return
	}

	//server 部署在线
	rs, err = db.Query(db.AR().From("online").Where(map[string]interface{}{
		"cs_id": serverId,
		"cs_type": "server",
	}).Limit(0, 1))
	if err != nil {
		jsonError(responseWrite, err.Error(), nil)
		return
	}
	if rs.Len() == 0 {
		jsonError(responseWrite, "server 未部署", nil)
		return
	}

	//查找 cluster ip
	clusterId := rs.Row()["cluster_id"]
	rs, err = db.Query(db.AR().From("cluster").Where(map[string]interface{}{
		"cluster_id": clusterId,
	}).Limit(0, 1))
	if err != nil {
		jsonError(responseWrite, err.Error(), nil)
		return
	}
	if(rs.Len() == 0) {
		jsonError(responseWrite, "cluster 不存在", nil)
		return
	}
	clusterIp := rs.Row()["ip"]

	//tunnel 开启请求地址
	tunnelOpenUri := strings.Replace("{host}", cfg.GetString("uri.tunnel_open"), clusterIp, 1)
	//:TunnelID/:ServerToken/:ServerBindIP/:ServerListenPort/:ClientToken/:ClientLocalHost/:ClientLocalPort/:Protocol
	url := tunnelOpenUri + "/" +
		tunnelValue["tunnel_id"] + "/" +
		serverValue["server_token"] + "/" +
		tunnelValue["server_listen_ip"] + "/" +
		tunnelValue["server_listen_port"] + "/" +
		clientValue["client_token"] + "/" +
		tunnelValue["client_local_host"] + "/" +
		tunnelValue["client_local_port"] + "/" +
		tunnelValue["protocol"]
	fmt.Println(url)
	body, _, err := at_common.HttpGet(url)
	if err != nil {
		jsonError(responseWrite, err.Error(), nil)
		return
	}

	type Response struct {
		Code    int
		Message string
		Data    interface{}
	}
	var res Response
	json.Unmarshal(body, &res)
	if res.Data == 0 {
		jsonError(responseWrite, res.Message, nil)
		return
	}

	//修改tunnel状态
	tunnelUpdateValue := map[string]interface{}{
		"status":   1,
		"is_open":   1,
		"update_time": time.Now().Unix(),
	}
	rs, err = db.Exec(db.AR().Update("tunnel", tunnelUpdateValue, map[string]interface{}{
		"tunnel_id": tunnelId,
	}))
	if err != nil {
		jsonError(responseWrite, err.Error(), nil)
		return
	}

	jsonSuccess(responseWrite, "开启隧道成功", nil)
}

//close tunnel
//method : GET
//params : tunnel_id
func (tunnel *Tunnel) TunnelClose(responseWrite http.ResponseWriter, request *http.Request, params httprouter.Params) {

	tunnelId := request.FormValue("tunnel_id")
	if tunnelId == "" {
		jsonError(responseWrite, "tunnel_id is error!", nil)
		return
	}
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err := db.Query(db.AR().From("tunnel").Where(map[string]interface{}{
		"tunnel_id": tunnelId,
		"is_delete": 0,
	}).Limit(0, 1))
	if err != nil {
		jsonError(responseWrite, err.Error(), nil)
		return
	}
	if rs.Len() == 0 {
		jsonError(responseWrite, "隧道不存在", nil)
		return
	}

	//查找 cluster ip
	clusterId := rs.Row()["cluster_id"]
	rs, err = db.Query(db.AR().From("cluster").Where(map[string]interface{}{
		"cluster_id": clusterId,
	}).Limit(0, 1))
	if err != nil {
		jsonError(responseWrite, err.Error(), nil)
		return
	}
	if(rs.Len() == 0) {
		jsonError(responseWrite, "cluster 不存在", nil)
		return
	}
	clusterIp := rs.Row()["ip"]
	//tunnel 关闭请求地址
	tunnelCloseUri := strings.Replace("{host}", cfg.GetString("uri.tunnel_close"), clusterIp, 1)

	type Response struct {
		Code    int
		Message string
		Data    interface{}
	}
	var res Response
	url := tunnelCloseUri + "/" + tunnelId
	body, _, err := at_common.HttpGet(url)
	if err != nil {
		jsonError(responseWrite, err.Error(), nil)
		return
	}
	json.Unmarshal(body, &res)
	if res.Data == 0 {
		jsonError(responseWrite, res.Message, nil)
		return
	}

	jsonSuccess(responseWrite, "隧道关闭成功", nil)
}

//tunnel status
//method : GET
//params : tunnel_id
func (tunnel *Tunnel) TunnelStatus(responseWrite http.ResponseWriter, request *http.Request, params httprouter.Params) {

	tunnelId := request.FormValue("tunnel_id")
	if tunnelId == "" {
		jsonError(responseWrite, "tunnel_id is error!", nil)
		return
	}
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err := db.Query(db.AR().From("tunnel").Where(map[string]interface{}{
		"tunnel_id": tunnelId,
		"is_delete": 0,
	}).Limit(0, 1))
	if err != nil {
		jsonError(responseWrite, err.Error(), nil)
		return
	}
	if rs.Len() == 0 {
		jsonError(responseWrite, "隧道不存在", nil)
		return
	}

	//查找 cluster ip
	clusterId := rs.Row()["cluster_id"]
	rs, err = db.Query(db.AR().From("cluster").Where(map[string]interface{}{
		"cluster_id": clusterId,
	}).Limit(0, 1))
	if err != nil {
		jsonError(responseWrite, err.Error(), nil)
		return
	}
	if(rs.Len() == 0) {
		jsonError(responseWrite, "cluster 不存在", nil)
		return
	}
	clusterIp := rs.Row()["ip"]
	//tunnel 关闭请求地址
	tunnelStatusUri := strings.Replace("{host}", cfg.GetString("uri.tunnel_status"), clusterIp, 1)

	type Response struct {
		Code    int
		Message string
		Data    interface{}
	}
	var res Response
	url := tunnelStatusUri + "/" + tunnelId
	body, _, err := at_common.HttpGet(url)
	if err != nil {
		jsonError(responseWrite, err.Error(), nil)
		return
	}
	json.Unmarshal(body, &res)
	if res.Data == 0 {
		jsonError(responseWrite, res.Message, nil)
		return
	}

	jsonSuccess(responseWrite, "ok", res.Message)
}
