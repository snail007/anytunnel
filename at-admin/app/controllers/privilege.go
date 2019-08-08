package controllers

import (
	"anytunnel/at-admin/app/models"
	"regexp"
	"time"

	"github.com/go-ozzo/ozzo-validation"
)

type PrivilegeController struct {
	BaseAdminController
}

func (this *PrivilegeController) Delete() {
	privilegeModel := models.Privilege{}
	privilegeId := this.GetString("privilege_id")
	privilege, err := privilegeModel.GetPrivilegeByPrivilegeId(privilegeId)
	if err != nil {
		this.jsonError(err)
	}
	hasSub, err := privilegeModel.HasSub(privilegeId)
	if err != nil {
		this.jsonError(err)
	}
	if privilege["type"] != "controller" && hasSub {
		this.jsonError("子菜单非空,不能删除")
	}
	err = privilegeModel.Delete(privilegeId)
	if err != nil {
		this.jsonError(err)
	}
	this.jsonSuccess("")

}
func (this *PrivilegeController) Add() {
	privilegeModel := models.Privilege{}
	if this.Ctx.Input.IsPost() {
		_, data := this.getPrivilegeFromPost(false)
		_, err := privilegeModel.Insert(data)
		if err != nil {
			this.jsonError(err)
		}
		this.jsonSuccess("")
	} else {
		navigators, menus, _, err := privilegeModel.GetTypedPrivileges("1", "-1")
		if err != nil {
			this.jsonError(err, "")
		}
		this.Data["navigators"] = navigators
		this.Data["menus"] = menus
		this.Data["action"] = "add"
		this.viewLayout("privilege/form", "form")
	}

}
func (this *PrivilegeController) Edit() {
	privilegeModel := models.Privilege{}
	privilegeId := this.GetString("privilege_id")
	if this.Ctx.Input.IsPost() {
		_, data := this.getPrivilegeFromPost(true)
		_, err := privilegeModel.Update(privilegeId, data)
		if err != nil {
			this.jsonError(err)
		}
		this.jsonSuccess("")
	} else {
		navigators, menus, _, err := privilegeModel.GetTypedPrivileges("1", "-1")
		if err != nil {
			this.jsonError(err, "")
		}
		privilege, err := privilegeModel.GetPrivilegeByPrivilegeId(privilegeId)
		if err != nil {
			this.jsonError(err, "")
		}
		if len(privilege) == 0 {
			this.jsonError("权限不存在")
		}
		this.Data["privilege"] = privilege
		this.Data["navigators"] = navigators
		this.Data["menus"] = menus
		this.Data["action"] = "edit"
		this.viewLayout("privilege/form", "form")
	}

}
func (this *PrivilegeController) getPrivilegeFromPost(isUpdate bool) (privilegeId string, privilege map[string]interface{}) {
	parentId := 0
	if this.GetString("type") == "menu" {
		parentId, _ = this.GetInt("parent_n")
	} else if this.GetString("type") == "controller" {
		parentId, _ = this.GetInt("parent_m")
	}
	privilege = map[string]interface{}{
		"name":       this.GetString("name"),
		"parent_id":  parentId,
		"type":       this.GetString("type"),
		"controller": this.GetString("controller"),
		"action":     this.GetString("action"),
		"icon":       this.GetString("icon"),
		"is_display": this.GetString("is_display"),
		"sequence":   this.GetString("sequence"),
		"target":     this.GetString("target"),
	}
	err := validation.Errors{
		"名称": validation.Validate(privilege["name"],
			validation.Required.Error("不能为空"),
			validation.Match(regexp.MustCompile("^.{1,15}$")).Error("长度必须是1-15字符")),
		"类型": validation.Validate(privilege["type"],
			validation.Required.Error("不能为空"),
			validation.In("navigator", "menu", "controller").Error("错误")),
		"排序": validation.Validate(privilege["sequence"],
			validation.Required.Error("不能为空"),
			validation.Match(regexp.MustCompile("^[1-9]{1}[0-9]*$")).Error("必须是大于0的数字")),
	}.Filter()
	if err != nil {
		this.jsonError(err)
	}
	if isUpdate {
		privilege["update_time"] = time.Now().Unix()
	} else {
		privilege["create_time"] = time.Now().Unix()
	}
	if this.GetString("type") != "controller" {
		privilege["is_display"] = 1
	}
	return
}

// 权限列表
func (this *PrivilegeController) List() {
	privilegeModel := new(models.Privilege)

	navigators, menus, controllers, err := privilegeModel.GetTypedPrivileges("1", "-1")
	if err != nil {
		this.viewError(err.Error())
	}

	this.Data["navigators"] = navigators
	this.Data["menus"] = menus
	this.Data["controllers"] = controllers

	this.view("privilege/list")
}
