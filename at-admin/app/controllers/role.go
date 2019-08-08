package controllers

import (
	"anytunnel/at-admin/app/models"
	"regexp"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

type RoleController struct {
	BaseAdminController
}

func (this *RoleController) Delete() {
	roleModel := models.Role{}
	roleId := this.GetString("role_id")
	hasUser, err := roleModel.HasUser(roleId)
	if err != nil {
		this.jsonError(err)
	}
	if hasUser {
		this.jsonError("角色用户非空,不能删除")
	}
	err = roleModel.Delete(roleId)
	if err != nil {
		this.jsonError(err)
	}
	this.jsonSuccess("")

}
func (this *RoleController) Add() {
	roleModel := models.Role{}
	if this.Ctx.Input.IsPost() {
		_, data := this.getRoleFromPost(false)
		_, err := roleModel.Insert(data)
		if err != nil {
			this.jsonError(err)
		}
		this.jsonSuccess("")
	} else {

		this.Data["action"] = "add"
		this.viewLayout("role/form", "form")
	}

}
func (this *RoleController) Edit() {
	roleModel := models.Role{}
	roleId := this.GetString("role_id")
	if this.Ctx.Input.IsPost() {
		_, data := this.getRoleFromPost(true)
		_, err := roleModel.Update(roleId, data)
		if err != nil {
			this.jsonError(err)
		}
		this.jsonSuccess("")
	} else {
		role, err := roleModel.GetRoleByRoleId(roleId)
		if err != nil {
			this.jsonError(err, "")
		}
		if len(role) == 0 {
			this.jsonError("角色不存在")
		}
		this.Data["role"] = role
		this.Data["action"] = "edit"
		this.viewLayout("role/form", "form")
	}

}

func (this *RoleController) List() {
	roleModel := models.Role{}

	roles, err := roleModel.GetAllRoles()
	if err != nil {
		this.viewError(err.Error())
	}

	this.Data["roles"] = roles
	this.view("role/list")
}

func (this *RoleController) getRoleFromPost(isUpdate bool) (roleId string, role map[string]interface{}) {
	role = map[string]interface{}{
		"name":      this.GetString("name"),
		"is_delete": 0,
	}
	err := validation.Validate(role["name"],
		validation.Required.Error("名称不能为空"),
		validation.Match(regexp.MustCompile("^.{1,15}$")).Error("名称长度必须是1-15字符"))
	if err != nil {
		this.jsonError(err.Error())
	}
	if isUpdate {
		role["update_time"] = time.Now().Unix()
	} else {
		role["create_time"] = time.Now().Unix()
	}
	return
}
