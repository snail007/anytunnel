package controllers

import (
	"anytunnel/at-admin/app/modules/web/models"
	"anytunnel/at-admin/app/utils"
	"regexp"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type UserController struct {
	BaseController
}

func (this *UserController) Forbidden() {
	userModel := models.User{}
	userId := this.GetString("user_id")
	// if userId == "0" {
	// 	this.JsonError("系统用户不能禁用")
	// }
	err := userModel.Forbidden(userId)
	if err != nil {
		this.JsonError(err)
	}
	this.JsonSuccess("")

}

func (this *UserController) Review() {
	userModel := models.User{}
	userId := this.GetString("user_id")
	// if userId == "0" {
	// 	this.JsonError("系统用户不能操作")
	// }
	err := userModel.Review(userId)
	if err != nil {
		this.JsonError(err)
	}
	this.JsonSuccess("")

}

// func (this *UserController) ChangePassword() {
// 	userModel := models.User{}
// 	if this.Ctx.Input.IsPost() {
// 		newpassword := this.GetString("password")
// 		oldpassword := this.GetString("password_old")
// 		errs := validation.Errors{
// 			"旧密码": validation.Validate(oldpassword,
// 				validation.Required.Error("不能为空")),
// 			"新密码": validation.Validate(newpassword,
// 				validation.Required.Error("不能为空"),
// 				validation.Match(regexp.MustCompile("^([0-9]+[a-zA-Z]+[_]*){1,16}$")).Error("必须同时包含数字和字母,且1-15字符")),
// 		}
// 		err := errs.Filter()
// 		if err != nil {
// 			this.JsonError(err)
// 		}
// 		err = userModel.ChangePassword(this.loginUser["user_id"], newpassword, oldpassword)
// 		if err != nil {
// 			this.JsonError("修改密码失败:" + err.Error())
// 		}
// 		this.JsonSuccess("")
// 	} else {
// 		this.viewLayout("user/changepassword", "form")
// 	}
// }
func (this *UserController) Add() {
	return
	userModel := models.User{}
	roleModel := models.Role{}
	userRoleModel := models.UserRole{}

	if this.Ctx.Input.IsPost() {
		_, data := this.getUserFromPost(false)
		username := this.GetString("username")
		roleIds := this.GetStrings("role_ids", []string{})
		if len(roleIds) == 0 {
			this.JsonError("没有选择角色")
		}
		HasUsername, err := userModel.HasUsername(username)
		if err != nil {
			this.JsonError(err)
		}
		if HasUsername {
			this.JsonError("用户名已经存在")
		}
		userId, err := userModel.Insert(data)
		if err != nil {
			this.JsonError("添加用户失败：" + err.Error())
		}

		//添加用户与角色对应关系
		_, err = userRoleModel.Insert(utils.NewConvert().IntToString(userId, 10), roleIds)
		if err != nil {
			this.JsonError("添加用户角色失败：" + err.Error())
		}

		this.JsonSuccess("")
	} else {

		roles := []map[string]string{}
		allRoles, _ := roleModel.GetAllRoles()
		for _, allRole := range allRoles {
			role := allRole
			role["is_default"] = "0"
			roles = append(roles, role)
		}
		this.Data["action"] = "add"
		this.Data["roles"] = roles
		this.viewLayout("user/form", "form")
	}

}
func (this *UserController) Edit() {
	userModel := models.User{}
	roleModel := models.Role{}
	userRoleModel := models.UserRole{}

	userId := this.GetString("user_id")
	if this.Ctx.Input.IsPost() {
		// _, data := this.getUserFromPost(true)
		// username := this.GetString("username")
		roleIds := this.GetStrings("role_ids", []string{})
		if len(roleIds) == 0 {
			this.JsonError("没有选择角色")
		}
		// HasSameUsername, err := userModel.HasSameUsername(userId, username)
		// if err != nil {
		// 	this.JsonError(err)
		// }
		// if HasSameUsername {
		// 	this.JsonError("用户名已经存在")
		// }
		// _, err = userModel.Update(userId, data)
		// if err != nil {
		// 	this.JsonError("修改用户失败：" + err.Error())
		// }

		//添加用户与角色对应关系
		_, err := userRoleModel.Insert(userId, roleIds)
		if err != nil {
			this.JsonError("修改用户角色失败：" + err.Error())
		}

		this.JsonSuccess("")
	} else {
		roles := []map[string]string{}
		user, err := userModel.GetUserByUserId(userId)
		allRoles, _ := roleModel.GetAllRoles()
		userRoles, _ := userRoleModel.GetUserRolesByUserId(userId)
		for _, allRole := range allRoles {
			role := allRole
			if len(userRoles) == 0 {
				role["is_default"] = "0"
			} else {
				for _, userRoles := range userRoles {
					if allRole["role_id"] == userRoles["role_id"] {
						role["is_default"] = "1"
						break
					}
					role["is_default"] = "0"
				}
			}
			roles = append(roles, role)
		}
		if err != nil {
			this.JsonError(err, "")
		}
		if len(user) == 0 {
			this.JsonError("用户不存在")
		}
		this.Data["user"] = user
		this.Data["roles"] = roles
		this.Data["action"] = "edit"
		this.viewLayout("user/form", "form")
	}

}
func (this *UserController) getUserFromPost(isUpdate bool) (userId string, user map[string]interface{}) {
	userModel := models.User{}
	user = map[string]interface{}{
		"username": this.GetString("username"),
		"nickname": this.GetString("nickname"),
		"email":    this.GetString("email"),
	}
	errs := validation.Errors{
		"邮箱": validation.Validate(user["email"],
			validation.Required.Error("不能为空"),
			is.Email.Error("格式错误")),
		"昵称": validation.Validate(user["nickname"],
			validation.Required.Error("不能为空"),
			validation.Match(regexp.MustCompile("^.{1,15}$")).Error("长度必须是1-15字符")),

		//"角色": validation.Validate(user["role_ids"]),
		//	validation.Required.Error("没有选择角色"),
	}
	if !isUpdate {
		errs["用户名"] = validation.Validate(user["username"],
			validation.Required.Error("不能为空"),
			validation.Match(regexp.MustCompile("^[0-9_a-zA-Z]{1,15}$")).Error("只能包含数字字母和下划线,且1-15字符"))
		errs["密码"] = validation.Validate(this.GetString("password"),
			validation.Required.Error("不能为空"),
			validation.Match(regexp.MustCompile("^([0-9]+[a-zA-Z]+[_]*){1,16}$")).Error("必须同时包含数字和字母,且1-15字符"))
	}
	err := errs.Filter()
	if err != nil {
		this.JsonError(err)
	}
	if isUpdate {
		if this.GetString("password") != "" {
			user["password"] = userModel.EncodePassword(this.GetString("password"))
		}
		user["update_time"] = time.Now().Unix()
		delete(user, "username")
	} else {
		user["is_forbidden"] = 0
		user["password"] = userModel.EncodePassword(this.GetString("password"))
		user["create_time"] = time.Now().Unix()
	}
	return
}

func (this *UserController) List() {

	keyword := strings.Trim(this.GetString("keyword", ""), " ")
	page, err := this.GetInt("page", 1)
	if err != nil {
		page = 1
	}
	//每页的条数
	pageSize := 10
	limit := (page - 1) * pageSize

	userModel := models.User{}
	var users = []map[string]string{}
	var userCount = 0

	if keyword == "" {
		userCount, err = userModel.CountUsers()
		users, err = userModel.GetUsersByLimit(limit, pageSize)
	} else {
		userCount, err = userModel.CountUsersByKeyword(keyword)
		users, err = userModel.GetUsersByKeywordAndLimit(keyword, limit, pageSize)
	}
	if err != nil {
		this.ViewError(err.Error())
	}

	roleModel := models.Role{}
	userRoles := map[string]string{}

	//用户角色
	for _, user := range users {
		userId := user["user_id"]
		var names = ""
		roles, err := roleModel.GetRolesByUserId(userId)
		if err != nil {
			this.ViewError(err.Error())
		}
		for _, role := range roles {
			names += "," + role["name"]
		}
		userRoles[user["user_id"]] = strings.Replace(names, ",", "", 1)
	}

	this.Data["users"] = users
	this.Data["userRoles"] = userRoles
	this.Data["keyword"] = keyword
	this.Data["page"] = utils.NewMisc().Page(userCount, page, pageSize, "/user/list?page={page}")
	this.viewLayoutTitle("用户列表", "user/list", "form")
}
