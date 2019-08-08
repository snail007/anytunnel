package controllers

import (
	"anytunnel/at-admin/app/modules/web/models"
	"anytunnel/at-admin/app/utils"
	"fmt"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type IpListController struct {
	BaseController
}

func (this *IpListController) Delete() {
	ip_listModel := models.IpList{}
	ip_listId := this.GetString("ip_list_id")
	err := ip_listModel.Delete(ip_listId)
	if err != nil {
		this.JsonError(err)
	}
	this.JsonSuccess("")

}
func (this *IpListController) Add() {
	ip_listModel := models.IpList{}
	if this.Ctx.Input.IsPost() {
		_, data := this.getIpListFromPost(false)
		_, err := ip_listModel.Insert(data)
		if err != nil {
			this.JsonError(err)
		}
		this.JsonSuccess("")
	} else {
		regionModel := models.Region{}
		regions, err := regionModel.GetSubRegions()
		if err != nil {
			this.ViewError(err.Error())
		}
		this.Data["regions"] = regions
		this.Data["action"] = "add"
		this.viewLayout("ip_list/form", "form")
	}

}
func (this *IpListController) Edit() {
	ip_listModel := models.IpList{}
	ip_listId := this.GetString("ip_list_id")
	if this.Ctx.Input.IsPost() {
		_, data := this.getIpListFromPost(true)
		_, err := ip_listModel.Update(ip_listId, data)
		if err != nil {
			this.JsonError(err)
		}
		this.JsonSuccess("")
	} else {
		ip_list, err := ip_listModel.GetIpListByIpListId(ip_listId)
		if err != nil {
			this.JsonError(err, "")
		}
		if len(ip_list) == 0 {
			this.ViewError("IpList不存在")
		}
		regionModel := models.Region{}
		regions, err := regionModel.GetSubRegions()
		if err != nil {
			this.ViewError(err.Error())
		}
		this.Data["regions"] = regions
		this.Data["ip_list"] = ip_list
		this.Data["action"] = "edit"
		this.viewLayout("ip_list/form", "form")
	}

}
func (this *IpListController) List() {
	column := strings.Trim(this.GetString("type", ""), " ")
	keyword := strings.Trim(this.GetString("id", ""), " ")
	page, err := this.GetInt("page", 1)
	if err != nil {
		page = 1
	}
	//每页的条数
	pageSize := 10
	limit := (page - 1) * pageSize

	ip_listModel := models.IpList{}
	var ip_lists = []map[string]string{}
	var ip_listCount = 0

	if keyword == "" {
		ip_listCount, err = ip_listModel.CountIpLists()
		ip_lists, err = ip_listModel.GetIpListsByLimit(limit, pageSize)
	} else {
		ip_listCount, err = ip_listModel.CountIpListsByID(column, keyword)
		ip_lists, err = ip_listModel.GetIpListsByIDAndLimit(column, keyword, limit, pageSize)
	}
	if err != nil {
		this.ViewError(err.Error())
	}
	this.Data["ip_lists"] = ip_lists
	this.Data["page"] = utils.NewMisc().Page(ip_listCount, page, pageSize, fmt.Sprintf("/web/ip_list/list?page={page}&type=%s&id=%s", column, keyword))
	this.viewLayoutTitle("用户列表", "ip_list/list", "form")
}
func (this *IpListController) getIpListFromPost(isUpdate bool) (ip_listId string, ip_list map[string]interface{}) {
	ip_list = map[string]interface{}{
		"ip":           this.GetString("ip"),
		"cs_type":      this.GetString("cs_type"),
		"is_forbidden": this.GetString("is_forbidden"),
	}
	errs := validation.Errors{
		"IP": validation.Validate(ip_list["ip"],
			validation.Required.Error("不能为空"),
			is.IPv4.Error("格式错误")),
		"类型": validation.Validate(ip_list["cs_type"],
			validation.Required.Error("不能为空"),
			validation.In("server", "client").Error("错误"),
		),
		"访问控制": validation.Validate(ip_list["is_forbidden"],
			validation.Required.Error("不能为空"),
			validation.In("0", "1").Error("错误"),
		),
	}
	err := errs.Filter()
	if err != nil {
		this.JsonError(err)
	}
	if isUpdate {
		ip_list["update_time"] = time.Now().Unix()
	} else {
		ip_list["create_time"] = time.Now().Unix()
	}
	return
}
