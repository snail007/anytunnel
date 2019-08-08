package controllers

import (
	"anytunnel/at-admin/app/modules/web/models"
	"anytunnel/at-admin/app/utils"
	"regexp"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

type ServerController struct {
	BaseController
}

func (this *ServerController) Reset() {
	serverModel := models.Server{}
	serverId := this.GetString("server_id")
	err := serverModel.Offline(serverId)
	if err != nil {
		this.JsonError(err)
	}
	err = serverModel.Reset(serverId)
	if err != nil {
		this.JsonError(err)
	}
	this.JsonSuccess("")

}
func (this *ServerController) Delete() {
	serverModel := models.Server{}
	serverId := this.GetString("server_id")
	hasTunnel, err := serverModel.HasTunnelRef(serverId)
	if err != nil {
		this.JsonError(err)
	}
	if hasTunnel {
		this.JsonError("Tunnel引用,不能删除")
	}
	err = serverModel.Delete(serverId)
	if err != nil {
		this.JsonError(err)
	}
	this.JsonSuccess("")

}
func (this *ServerController) Add() {
	serverModel := models.Server{}
	if this.Ctx.Input.IsPost() {
		_, data := this.getServerFromPost(false)
		_, err := serverModel.Insert(data)
		if err != nil {
			this.JsonError(err)
		}
		this.JsonSuccess("")
	} else {
		this.Data["action"] = "add"
		this.viewLayout("server/form", "form")
	}
}
func (this *ServerController) Edit() {
	serverModel := models.Server{}
	serverId := this.GetString("server_id")
	if this.Ctx.Input.IsPost() {
		_, data := this.getServerFromPost(true)
		_, err := serverModel.Update(serverId, data)
		if err != nil {
			this.JsonError(err)
		}
		this.JsonSuccess("")
	} else {
		server, err := serverModel.GetServerByServerId(serverId)
		if err != nil {
			this.JsonError(err, "")
		}
		if len(server) == 0 {
			this.ViewError("Server不存在")
		}
		this.Data["server"] = server
		this.Data["action"] = "edit"
		this.viewLayout("server/form", "form")
	}

}
func (this *ServerController) List() {

	userID := strings.Trim(this.GetString("user_id"), " ")
	page, err := this.GetInt("page", 1)
	if err != nil {
		page = 1
	}
	//每页的条数
	pageSize := 10
	limit := (page - 1) * pageSize

	serverModel := models.Server{}
	var servers = []map[string]string{}
	var serverCount = 0
	if userID == "" {
		serverCount, err = serverModel.CountServers()
		servers, err = serverModel.GetServersByLimit(limit, pageSize)
	} else {
		serverCount, err = serverModel.CountServersByUserID(userID)
		servers, err = serverModel.GetServersByUserIDAndLimit(userID, limit, pageSize)
	}
	if err != nil {
		this.ViewError(err.Error())
	}

	this.Data["userID"] = userID
	this.Data["servers"] = servers
	this.Data["page"] = utils.NewMisc().Page(serverCount, page, pageSize, "/web/server/list?page={page}&user_id="+userID)
	this.viewLayoutTitle("Server列表", "server/list", "form")
}

func (this *ServerController) getServerFromPost(isUpdate bool) (serverId string, server map[string]interface{}) {
	server = map[string]interface{}{
		"name": this.GetString("name"),
		// "user_id":   this.GetString("user_id"),
		"user_id":   0,
		"is_delete": 0,
	}
	err := validation.Validate(server["name"],
		validation.Required.Error("名称不能为空"),
		validation.Match(regexp.MustCompile("^.{1,15}$")).Error("名称长度必须是1-15字符"))
	if err != nil {
		this.JsonError(err.Error())
	}
	if isUpdate {
		server["update_time"] = time.Now().Unix()
	} else {
		server["token"] = utils.NewMisc().RandString(32)
		server["create_time"] = time.Now().Unix()
	}
	return
}
