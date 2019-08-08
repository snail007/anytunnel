package controllers

import (
	"anytunnel/at-web/app/utils"
	"strings"
	"anytunnel/at-web/app/business"
)

type ClientController struct {
	BaseUserController
}

func (this *ClientController) List() {

	userId := this.getUserId()
	keyword := this.GetString("keyword")
	data := map[string]string {
		"user_id":  userId,
		"keyword":  keyword,
	}
	clientValues, err:= business.NewBase().GetRequest("client_list_uri", data)
	if(err != nil) {
		this.viewError("request error")
	}
	this.Data["clientValues"] = clientValues
	this.Data["keyword"] = keyword
	this.viewLayoutTitle("client列表", "client/list", "page")
}

//add client
func (this *ClientController) Add() {
	this.Data["clientValue"] = map[string]string{
		"client_id": "0",
		"name": "",
		"token": "",
	}
	this.viewLayoutTitle("添加client", "client/form", "page")
}

//save client
func (this *ClientController) Save() {

	name := strings.Trim(this.GetString("name"), "");

	data := map[string]string {
		"name":     name,
		"token":    utils.NewMisc().RandString(32),
		"user_id":  this.getUserId(),
	}
	_, err := business.NewBase().PostRequest("client_create_uri", data, nil)
	if(err != nil) {
		this.jsonError(err.Error())
	}

	this.jsonSuccess("添加Client成功", nil, "/client/list");
}

//edit client
func (this *ClientController) Edit() {

	clientId := this.GetString("client_id")
	if(clientId == "") {
		this.viewError("client_id error")
	}

	data := map[string]string{
		"client_id":clientId,
	}
	clientValues, err:= business.NewBase().GetRequest("client_info_uri", data)
	if(err != nil) {
		this.viewError(err.Error())
	}

	this.Data["clientValue"] = clientValues
	this.viewLayoutTitle("修改client", "client/form", "page")
}

//save client
func (this *ClientController) Modify() {

	name := this.GetString("name");
	clientId := this.GetString("client_id");

	data := map[string]string{
		"name": name,
		"client_id":clientId,
	}
	_, err := business.NewBase().PostRequest("client_update_uri", data, nil)
	if(err != nil) {
		this.jsonError(err.Error())
	}
	this.jsonSuccess("修改Client成功", nil, "/client/list");
}

//delete client
func (this *ClientController) Delete() {

	clientId := this.GetString("client_id")
	data := map[string]string{
		"client_id": clientId,
	}
	_, err := business.NewBase().GetRequest("client_delete_uri", data)
	if(err != nil) {
		this.jsonError(err.Error())
	}
	this.jsonSuccess("删除client成功", nil, "client/list")
}

//reset token
func (this *ClientController) ResetToken() {

	clientId := this.GetString("client_id");

	data := map[string]string{
		"token":        utils.NewMisc().RandString(32),
		"client_id":    clientId,
	}
	_, err := business.NewBase().PostRequest("client_update_uri", data, nil)
	if(err != nil) {
		this.jsonError(err.Error())
	}

	this.jsonSuccess("重置Token成功", nil, "/client/list");
}