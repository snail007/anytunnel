package controllers

import (
	utils "anytunnel/at-common"
	"encoding/json"
	"strings"

	"github.com/astaxie/beego"
)

type UserController struct {
	BaseUserController
}

//user center
func (this *UserController) Index() {
	this.Data["userValue"] = this.loginUser
	this.viewLayoutTitle("用户中心", "user/index", "user")
}

//user welcome
func (this *UserController) Welcome() {
	this.viewLayoutTitle("默认首页", "user/welcome", "page")
}

//user profile
func (this *UserController) Profile() {
	this.Data["userValue"] = this.loginUser
	this.viewLayoutTitle("默认首页", "user/profile", "page")
}

//user password
func (this *UserController) Password() {
	this.Data["userValue"] = this.loginUser
	this.viewLayoutTitle("默认首页", "user/password", "page")
}

//user save
func (this *UserController) Save() {

	userId := strings.Trim(this.GetString("user_id"), "")
	nickname := strings.Trim(this.GetString("nickname"), "")

	if userId == "" {
		this.jsonError("user_id error")
	}

	updateUri := beego.AppConfig.String("user_update_uri")
	if updateUri == "" {
		this.jsonError("修改资料失败")
	}

	data := map[string]string{
		"nickname": nickname,
		"user_id":  userId,
	}
	body, code, err := utils.HttpPost(updateUri, data, nil)
	if code != 200 || err != nil {
		this.jsonError("修改资料失败")
	}
	var results map[string]interface{}
	json.Unmarshal(body, &results)
	if results["code"].(float64) != 1 {
		this.jsonError(results["message"].(string))
	}

	//更新session和cookie
	this.refreshSession(userId)

	this.jsonSuccess("修改资料成功", nil, "/user/profile")
}

//user repass
func (this *UserController) Repass() {

	userId := strings.Trim(this.GetString("user_id"), "")
	oldPass := strings.Trim(this.GetString("old_pass"), "")
	newPass := strings.Trim(this.GetString("new_pass"), "")
	confirmPass := strings.Trim(this.GetString("confirm_pass"), "")

	if userId == "" {
		this.jsonError("user_id error")
	}
	if oldPass == "" {
		this.jsonError("旧密码错误")
	}
	if utils.NewEncrypt().Md5Encode(oldPass) != this.loginUser["password"].(string) {
		this.jsonError("旧密码错误")
	}
	if newPass == "" {
		this.jsonError("新密码不能为空")
	}
	if confirmPass == "" {
		this.jsonError("确认密码不能为空")
	}
	if confirmPass != newPass {
		this.jsonError("确认密码和新密码不一致")
	}

	updateUri := beego.AppConfig.String("user_update_uri")
	if updateUri == "" {
		this.jsonError("修改密码失败")
	}

	data := map[string]string{
		"password": newPass,
		"user_id":  userId,
	}
	body, code, err := utils.HttpPost(updateUri, data, nil)
	if code != 200 || err != nil {
		this.jsonError("修改密码失败")
	}
	var results map[string]interface{}
	json.Unmarshal(body, &results)
	if results["code"].(float64) != 1 {
		this.jsonError(results["message"].(string))
	}

	//更新session和cookie
	this.refreshSession(userId)

	this.jsonSuccess("修改密码成功", nil, "/user/password")
}
