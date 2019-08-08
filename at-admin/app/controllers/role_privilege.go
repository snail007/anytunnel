package controllers

import "anytunnel/at-admin/app/models"

type Role_PrivilegeController struct {
	BaseAdminController
}

//角色授权
func (this *Role_PrivilegeController) Add() {

	roleId, err := this.GetInt("role_id", 0)
	if roleId == 0 {
		this.viewError("role_id is error! ")
	}

	rolePrivilegeModel := models.RolePrivilege{}
	rolePrivileges, err := rolePrivilegeModel.GetRolePrivilegesByRoleId(roleId)
	if err != nil {
		this.viewError(err.Error())
	}

	privilegeModel := models.Privilege{}
	navigators, menus, controllers, err := privilegeModel.GetTypedPrivileges("1", "-1")
	if err != nil {
		this.viewError(err.Error())
	}

	this.Data["navigators"] = navigators
	this.Data["menus"] = menus
	this.Data["controllers"] = controllers
	this.Data["rolePrivileges"] = rolePrivileges
	this.Data["privileges"] = rolePrivileges
	this.Data["role_id"] = roleId

	this.viewLayout("role/privilege", "form")
}

func (this *Role_PrivilegeController) Save() {

	privilegeIds := this.GetStrings("privilege_id", []string{})
	roleId, err := this.GetInt("role_id", 0)

	if err != nil {
		this.jsonError("角色授权失败：role_id error " + err.Error())
	}
	if roleId == 0 {
		this.jsonError("角色授权失败：role_id is error!")
	}
	if len(privilegeIds) == 0 {
		this.jsonError("角色授权失败：no select privilege!")
	}

	rolePrivilegeModel := models.RolePrivilege{}

	res, err := rolePrivilegeModel.GrantRolePrivileges(roleId, privilegeIds)
	if !res {
		this.jsonError("角色授权失败：" + err.Error())
	}

	this.jsonSuccess("角色授权成功")
}
