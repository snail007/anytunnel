package controllers

import (
	"anytunnel/at-admin/app/modules/web/models"
	"anytunnel/at-admin/app/utils"
	"regexp"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

type ClientController struct {
	BaseController
}

func (this *ClientController) Reset() {
	clientModel := models.Client{}
	clientId := this.GetString("client_id")
	err := clientModel.Offline(clientId)
	if err != nil {
		this.JsonError(err)
	}
	err = clientModel.Reset(clientId)
	if err != nil {
		this.JsonError(err)
	}
	this.JsonSuccess("")
}
func (this *ClientController) Delete() {
	clientModel := models.Client{}
	clientId := this.GetString("client_id")
	hasTunnel, err := clientModel.HasTunnelRef(clientId)
	if err != nil {
		this.JsonError(err)
	}
	if hasTunnel {
		this.JsonError("Tunnel引用,不能删除")
	}
	err = clientModel.Delete(clientId)
	if err != nil {
		this.JsonError(err)
	}
	this.JsonSuccess("")

}
func (this *ClientController) Add() {
	clientModel := models.Client{}
	if this.Ctx.Input.IsPost() {
		_, data := this.getClientFromPost(false)
		_, err := clientModel.Insert(data)
		if err != nil {
			this.JsonError(err)
		}
		this.JsonSuccess("")
	} else {
		this.Data["action"] = "add"
		this.viewLayout("client/form", "form")
	}

}
func (this *ClientController) Edit() {
	clientModel := models.Client{}
	clientId := this.GetString("client_id")
	if this.Ctx.Input.IsPost() {
		_, data := this.getClientFromPost(true)
		_, err := clientModel.Update(clientId, data)
		if err != nil {
			this.JsonError(err)
		}
		this.JsonSuccess("")
	} else {
		client, err := clientModel.GetClientByClientId(clientId)
		if err != nil {
			this.JsonError(err, "")
		}
		if len(client) == 0 {
			this.ViewError("Client不存在")
		}
		this.Data["client"] = client
		this.Data["action"] = "edit"
		this.viewLayout("client/form", "form")
	}

}

func (this *ClientController) List() {

	userID := strings.Trim(this.GetString("user_id"), " ")
	page, err := this.GetInt("page", 1)
	if err != nil {
		page = 1
	}
	//每页的条数
	pageSize := 10
	limit := (page - 1) * pageSize

	clientModel := models.Client{}
	var clients = []map[string]string{}
	var clientCount = 0

	if userID == "" {
		clientCount, err = clientModel.CountClients()
		clients, err = clientModel.GetClientsByLimit(limit, pageSize)
	} else {
		clientCount, err = clientModel.CountClientsByUserID(userID)
		clients, err = clientModel.GetClientsByUserIDAndLimit(userID, limit, pageSize)
	}
	if err != nil {
		this.ViewError(err.Error())
	}
	this.Data["userID"] = userID
	this.Data["clients"] = clients
	this.Data["page"] = utils.NewMisc().Page(clientCount, page, pageSize, "/web/client/list?page={page}&user_id="+userID)
	this.viewLayoutTitle("Client列表", "client/list", "form")
}
func (this *ClientController) getClientFromPost(isUpdate bool) (clientId string, client map[string]interface{}) {
	port, _ := this.GetInt("local_port")
	client = map[string]interface{}{
		"name":       this.GetString("name"),
		"local_host": this.GetString("local_host"),
		"local_port": port,
		"user_id":    0,
		"is_delete":  0,
	}
	errs := validation.Errors{
		"名称": validation.Validate(client["name"],
			validation.Required.Error("不能为空"),
			validation.Match(regexp.MustCompile("^.{1,15}$")).Error("名称长度必须是1-15字符")),
		"本地网络Host": validation.Validate(client["local_host"],
			validation.Required.Error("不能为空")),
		"本地网络端口": validation.Validate(client["local_port"],
			validation.Required.Error("不能为空"),
			validation.Min(1).Error("最小值1"),
			validation.Max(65535).Error("最大值65535")),
	}
	err := errs.Filter()
	if err != nil {
		this.JsonError(err)
	}
	if isUpdate {
		client["update_time"] = time.Now().Unix()
	} else {
		client["token"] = utils.NewMisc().RandString(32)
		client["create_time"] = time.Now().Unix()
	}
	return
}
