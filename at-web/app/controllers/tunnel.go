package controllers

import (
	"anytunnel/at-web/app/business"
	"strings"
	"anytunnel/at-web/app/utils"
	"fmt"
)

const MODE_BASE  = "0"
const MODE_SENIOR  = "1"
const MODE_SPECIAL  = "2"

type TunnelController struct {
	BaseUserController
}

//tunnel list
func (this *TunnelController) List() {

	userId := this.getUserId()
	keyword := this.GetString("keyword")
	data := map[string]string {
		"user_id" : userId,
		"keyword" : keyword,
	}
	tunnelValues, err := business.NewBase().GetRequest("tunnel_list_uri", data)
	if(err != nil) {
		this.viewError("request error")
	}
	this.Data["tunnelValues"] = tunnelValues
	this.Data["keyword"] = keyword
	this.viewLayoutTitle("隧道列表", "tunnel/list", "page")
}

//tunnel add 1. tunnel mode
func (this *TunnelController) Mode() {

	userId := this.getUserId()
	data := map[string]string{
		"user_id": userId,
	}
	roles, err := business.NewBase().GetRequest("get_user_role", data)
	if(err != nil) {
		this.viewError("request error " + err.Error())
	}
	modeBase := "0"
	modeSenior := "0"
	modeSpecial := "0"
	for _, role := range roles.([]interface{}) {
		role := role.(map[string]interface{})
		tunnelMode := strings.Split(role["tunnel_mode"].(string), ",")
		for _, tunnel := range tunnelMode {
			if(tunnel == MODE_BASE) {
				modeBase = "1"
			}
			if(tunnel == MODE_SENIOR) {
				modeSenior = "1"
			}
			if(tunnel == MODE_SPECIAL) {
				modeSpecial = "1"
			}
		}
	}
	this.Data["modeBase"] = modeBase
	this.Data["modeSenior"] = modeSenior
	this.Data["modeSpecial"] = modeSpecial
	//获取所有的 server
	this.viewLayoutTitle("选择模式", "tunnel/mode", "page")
}

//tunnel add 2. tunnel cluster
func (this *TunnelController) Cluster() {

	clusterId := this.GetString("cluster_id", "");
	mode := this.GetString("mode", "0");

	data := map[string]string{
		"user_id": this.getUserId(),
		"mode": mode,
	}
	regionClusters, err := business.NewBase().GetRequest("user_cluster_list",data)
	if(err != nil) {
		this.viewError("request error " + err.Error())
	}
	this.Data["regionClusters"] = regionClusters
	this.Data["clusterId"] = clusterId
	this.Data["mode"] = mode
	this.viewLayoutTitle("选择节点", "tunnel/cluster", "page")
}

//tunnel add 2. tunnel add
func (this *TunnelController) Add() {

	if(this.Ctx.Request.Method != "POST") {
		this.viewError("request error")
	}
	clusterId := this.GetString("cluster_id", "");
	mode := this.GetString("mode", "0");
	if(clusterId == "") {
		this.viewError("没有选择节点")
	}

	isHaveClient := "0"
	isHaveServer := "0"
	//基础模式，只能部署client
	if(mode == MODE_BASE) {
		isHaveClient = "1"
	}
	//高级模式，只能部署client和server
	if(mode == MODE_SENIOR) {
		isHaveClient = "1"
		isHaveServer = "1"
	}
	//特殊模式，只能部署server
	if(mode == MODE_SPECIAL) {
		isHaveServer = "1"
	}

	userId := this.getUserId()
	data := map[string]string{
		"user_id": userId,
	}
	clientValues := []interface{}{}
	serverValues := []interface{}{}
	systemClients := []interface{}{}
	systemServers := []interface{}{}
	systemServer := map[string]interface{}{}
	systemClient := map[string]interface{}{}

	//获取系统在线client
	res, _ := business.NewBase().GetRequest("online_client_by_clusterId", map[string]string{
		"cluster_id": clusterId,
	})
	if(res != nil) {
		systemClients = res.([]interface{})
	}

	//获取系统在线server
	res, _ = business.NewBase().GetRequest("online_server_by_clusterId", map[string]string{
		"cluster_id": clusterId,
	})
	if(res != nil) {
		systemClients = res.([]interface{})
	}

	if(isHaveClient == "1") {
		res, _ := business.NewBase().GetRequest("client_list_uri", data)
		clientValues = res.([]interface{})
		//存在，取交集
		if(len(systemClients) > 0) {
			for index, clientValue := range clientValues {
				for _, systemClient := range systemClients {
					if(clientValue.(map[string]interface{})["client_id"].(string) == systemClient.(map[string]interface{})["cs_id"].(string)) {
						clientValues = append(clientValues[:index], clientValues[index+1:]...)
						break
					}
				}
			}
		}
	}else {
		//随机获取一个系统 client
		systemClient = utils.NewMisc().RandSlice(systemClients).(map[string]interface{})
	}
	if(isHaveServer == "1") {
		res, _ := business.NewBase().GetRequest("server_list_uri", data)
		serverValues = res.([]interface{})
		//存在，取交集
		if(len(systemServers) > 0) {
			for index, serverValue := range serverValues {
				for _, systemServer := range systemServers {
					if(serverValue.(map[string]string)["client_id"] == systemServer.(map[string]string)["cs_id"]) {
						serverValues = append(serverValues[:index], serverValues[index+1:]...)
						break
					}
				}
			}
		}
	}else {
		//随机获取一个系统 server
		systemServer = utils.NewMisc().RandSlice(systemServers).(map[string]interface{})
	}

	this.Data["isHaveClient"] = isHaveClient
	this.Data["isHaveServer"] = isHaveServer
	this.Data["clientValues"] = clientValues
	this.Data["serverValues"] = serverValues
	this.Data["systemClient"] = systemClient
	this.Data["systemServer"] = systemServer
	this.Data["clusterId"] = clusterId
	this.Data["mode"] = mode
	//获取所有的 server
	this.viewLayoutTitle("添加隧道", "tunnel/form", "page")
}

func (this *TunnelController) Save() {

	if(this.Ctx.Request.Method != "POST") {
		this.viewError("request error")
	}

	name := this.GetString("name", "");
	serverId := this.GetString("server_id", "");
	serverListenIp := this.GetString("server_listen_ip", "");
	serverListenPort := this.GetString("server_listen_port", "");
	clientId := this.GetString("client_id", "");
	clientLocalHost := this.GetString("client_local_host", "");
	clientLocalPort := this.GetString("client_local_port", "");
	mode := this.GetString("mode", "0");
	clusterId := this.GetString("cluster_id", "0");
	userId := this.getUserId()

	data := map[string]string{
		"name": strings.Trim(name, ""),
		"server_id": strings.Trim(serverId, ""),
		"server_listen_ip": strings.Trim(serverListenIp, ""),
		"server_listen_port": strings.Trim(serverListenPort, ""),
		"client_id": strings.Trim(clientId, ""),
		"client_local_host": strings.Trim(clientLocalHost, ""),
		"client_local_port": strings.Trim(clientLocalPort, ""),
		"mode": mode,
		"user_id": userId,
		"protocol": "1",
		"cluster_id": clusterId,
	}

	res, err := business.NewBase().PostRequest("tunnel_create_uri", data, nil)
	if(err != nil) {
		this.jsonError(err.Error())
	}

	fmt.Println(res)
	this.jsonSuccess("添加隧道成功", nil, "/tunnel/list")
}

//启动隧道
func (this *TunnelController) Open()  {

	tunnelId := this.GetString("tunnel_id", "0");
	if(tunnelId == "0") {
		this.jsonError("tunnel_id error")
	}

	data := map[string]string{
		"tunnel_id": tunnelId,
	}

	_, err := business.NewBase().GetRequest("tunnel_open_uri", data)
	if(err != nil) {
		this.jsonError(err.Error())
	}

	this.jsonSuccess("开启隧道成功", nil, "/tunnel/list")
}

//关闭隧道
func (this *TunnelController) Close()  {

	tunnelId := this.GetString("tunnel_id", "0");
	if(tunnelId == "0") {
		this.jsonError("tunnel_id error")
	}

	data := map[string]string{
		"tunnel_id": tunnelId,
	}

	_, err := business.NewBase().GetRequest("tunnel_close_uri", data)
	if(err != nil) {
		this.jsonError(err.Error())
	}

	this.jsonSuccess("关闭隧道成功")
}

//删除隧道
func (this *TunnelController) Delete() {

	tunnelId := this.GetString("tunnel_id", "0");
	if(tunnelId == "0") {
		this.jsonError("tunnel_id error")
	}

	data := map[string]string{
		"tunnel_id": tunnelId,
	}

	_, err := business.NewBase().GetRequest("tunnel_delete_uri", data)
	if(err != nil) {
		this.jsonError(err.Error())
	}

	this.jsonSuccess("删除隧道成功")
}

//重启隧道
func (this *TunnelController) Refresh()  {

	tunnelId := this.GetString("tunnel_id", "0");
	if(tunnelId == "0") {
		this.jsonError("tunnel_id error")
	}

	data := map[string]string{
		"tunnel_id": tunnelId,
	}

	//关闭
	_, err := business.NewBase().GetRequest("tunnel_close_uri", data)
	if(err != nil) {
		this.jsonError(err.Error())
	}

	//打开
	_, err = business.NewBase().GetRequest("tunnel_open_uri", data)
	if(err != nil) {
		this.jsonError(err.Error())
	}

	this.jsonSuccess("重启隧道成功", nil, "/tunnel/list")
}