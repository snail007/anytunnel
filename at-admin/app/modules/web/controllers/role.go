package controllers

import (
	"anytunnel/at-admin/app/modules/web/models"
	"regexp"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

type RoleController struct {
	BaseController
}

func (this *RoleController) Delete() {
	roleModel := models.Role{}
	roleId := this.GetString("role_id")
	// if roleId == "0" {
	// 	this.JsonError("系统角色不能删除")
	// }
	if roleId == "1" {
		this.JsonError("默认角色不能删除")
	}
	hasUser, err := roleModel.HasUser(roleId)
	if err != nil {
		this.JsonError(err)
	}
	if hasUser {
		this.JsonError("角色用户非空,不能删除")
	}
	err = roleModel.SetRegions(roleId, []string{})
	if err != nil {
		this.JsonError(err)
	}
	err = roleModel.Delete(roleId)
	if err != nil {
		this.JsonError(err)
	}
	this.JsonSuccess("")

}
func (this *RoleController) Add() {
	roleModel := models.Role{}
	if this.Ctx.Input.IsPost() {
		_, data := this.getRoleFromPost(false)
		_, err := roleModel.Insert(data)
		if err != nil {
			this.JsonError(err)
		}
		this.JsonSuccess("")
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
			this.JsonError(err)
		}
		this.JsonSuccess("")
	} else {
		role, err := roleModel.GetRoleByRoleId(roleId)
		if err != nil {
			this.JsonError(err, "")
		}
		if len(role) == 0 {
			this.ViewError("角色不存在")
		}
		this.Data["tunnel_mode_arr"] = strings.Split(role["tunnel_mode"], ",")
		this.Data["role"] = role
		this.Data["action"] = "edit"
		this.viewLayout("role/form", "form")
	}

}

func (this *RoleController) List() {
	roleModel := models.Role{}

	roles, err := roleModel.GetAllRoles()
	if err != nil {
		this.ViewError(err.Error())
	}

	this.Data["roles"] = roles
	this.view("role/list")
}

func (this *RoleController) Regions() {
	roleModel := models.Role{}
	roleID := this.GetString("role_id")
	if this.Ctx.Input.IsPost() {
		regionIds := this.GetStrings("region-ids")
		err := roleModel.SetRegions(roleID, regionIds)
		if err != nil {
			this.JsonError(err)
		}
		this.JsonSuccess("")
	} else {
		regionModel := models.Region{}
		regions, err := regionModel.GetAllRegions()
		if err != nil {
			this.ViewError(err.Error())
		}
		regionIds, err := roleModel.GetRoleRegionIds(roleID)
		if err != nil {
			this.ViewError(err.Error())
		}
		this.Data["regionIds"] = regionIds
		this.Data["roleId"] = this.GetString("role_id")
		this.Data["regions"] = regions
		this.view("role/region")
	}
}
func (this *RoleController) getRoleFromPost(isUpdate bool) (roleId string, role map[string]interface{}) {
	role = map[string]interface{}{
		"name":        this.GetString("name"),
		"server_area": this.GetString("server_area"),
		"client_area": this.GetString("client_area"),
		"bandwidth":   this.GetString("bandwidth"),
		"is_delete":   0,
	}
	errs := validation.Errors{
		"名称": validation.Validate(role["name"],
			validation.Required.Error("不能为空"),
			validation.Match(regexp.MustCompile("^.{1,15}$")).Error("长度必须是1-15字符")),
		"Server区域": validation.Validate(role["server_area"],
			validation.Required.Error("不能为空"),
			validation.Match(regexp.MustCompile("^(china|foreign|all)$")).Error("错误")),
		"Client区域": validation.Validate(role["client_area"],
			validation.Required.Error("不能为空"),
			validation.Match(regexp.MustCompile("^(china|foreign|all)$")).Error("错误")),
		"带宽": validation.Validate(role["bandwidth"],
			validation.Required.Error("不能为空"),
			validation.Match(regexp.MustCompile("^(0|([1-9][0-9]*))$")).Error("必须是大于等于0的整数")),
	}
	err := errs.Filter()
	if err != nil {
		this.JsonError(err)
	}
	role["tunnel_mode"] = strings.Join(this.GetStrings("tunnel_mode"), ",")
	if isUpdate {
		role["update_time"] = time.Now().Unix()
	} else {
		role["create_time"] = time.Now().Unix()
	}
	return
}
