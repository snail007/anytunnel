package controllers

import (
	"anytunnel/at-web/app/utils"
	"strings"
	"anytunnel/at-web/app/business"
)

type ServerController struct {
	BaseUserController
}

func (this *ServerController) List() {
	
	userId := this.getUserId()
	keyword := this.GetString("keyword")
	data := map[string]string {
		"user_id":  userId,
		"keyword":  keyword,
	}
	serverValues, err:= business.NewBase().GetRequest("server_list_uri", data)
	if(err != nil) {
		this.viewError("request error")
	}
	this.Data["serverValues"] = serverValues
	this.Data["keyword"] = keyword
	this.viewLayoutTitle("server列表", "server/list", "page")
}

//add server
func (this *ServerController) Add() {
	this.Data["serverValue"] = map[string]string{
		"server_id": "0",
		"name": "",
		"token": "",
	}
	this.viewLayoutTitle("添加server", "server/form", "page")
}

//save server
func (this *ServerController) Save() {
	
	name := strings.Trim(this.GetString("name"), "");
	
	data := map[string]string {
		"name":     name,
		"token":    utils.NewMisc().RandString(32),
		"user_id":  this.getUserId(),
	}
	_, err := business.NewBase().PostRequest("server_create_uri", data, nil)
	if(err != nil) {
		this.jsonError(err.Error())
	}
	
	this.jsonSuccess("添加Server成功", nil, "/server/list");
}

//edit server
func (this *ServerController) Edit() {
	
	serverId := this.GetString("server_id")
	if(serverId == "") {
		this.viewError("server_id error")
	}
	
	data := map[string]string{
		"server_id":serverId,
	}
	serverValues, err:= business.NewBase().GetRequest("server_info_uri", data)
	if(err != nil) {
		this.viewError(err.Error())
	}
	
	this.Data["serverValue"] = serverValues
	this.viewLayoutTitle("修改server", "server/form", "page")
}

//save server
func (this *ServerController) Modify() {
	
	name := this.GetString("name");
	serverId := this.GetString("server_id");
	
	data := map[string]string{
		"name": name,
		"server_id":serverId,
	}
	_, err := business.NewBase().PostRequest("server_update_uri", data, nil)
	if(err != nil) {
		this.jsonError(err.Error())
	}
	this.jsonSuccess("修改Server成功", nil, "/server/list");
}

//delete server
func (this *ServerController) Delete() {
	
	serverId := this.GetString("server_id")
	data := map[string]string{
		"server_id": serverId,
	}
	_, err := business.NewBase().GetRequest("server_delete_uri", data)
	if(err != nil) {
		this.jsonError(err.Error())
	}
	this.jsonSuccess("删除server成功", nil, "server/list")
}

//reset token
func (this *ServerController) ResetToken() {
	
	serverId := this.GetString("server_id");
	
	data := map[string]string{
		"token":        utils.NewMisc().RandString(32),
		"server_id":    serverId,
	}
	_, err := business.NewBase().PostRequest("server_update_uri", data, nil)
	if(err != nil) {
		this.jsonError(err.Error())
	}
	
	this.jsonSuccess("重置Token成功", nil, "/server/list");
}