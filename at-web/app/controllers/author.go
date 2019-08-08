package controllers

import (
	"strings"
	"github.com/golang-collections/lib.go/validation/validator"
	"github.com/astaxie/beego"
	"anytunnel/at-web/app/utils"
	"anytunnel/at-web/app/business"
)

type AuthorController struct {
	BaseUserController
}

// login
func (this *AuthorController) Login() {
	this.viewLayoutTitle("AnyTunnelCloud", "author/login", "login")
}

// register
func (this *AuthorController) Register() {
	this.viewLayoutTitle("AnyTunnelCloud", "author/register", "login")
}

//register save
func (this *AuthorController) Signin() {

	username := strings.Trim(this.GetString("username"), "")
	password := strings.Trim(this.GetString("password"), "")
	email := strings.Trim(this.GetString("email"), "")

	if(username == "") {
		this.jsonError("注册失败：用户名不能为空")
	}
	if(password == "") {
		this.jsonError("注册失败：密码不能为空")
	}
	if(email == "") {
		this.jsonError("注册失败：邮箱不能为空")
	}
	if(!validator.IsEmail(email)) {
		this.jsonError("注册失败：邮箱格式错误")
	}
	data := map[string]string{
		"username":    username,
		"password":    password,
		"email":       email,
	}
	user, err := business.NewBase().PostRequest("create_user", data, nil);
	if(err != nil) {
		this.jsonError("注册失败")
	}
	userMap := user.(map[string]interface{})
	userId := utils.NewConvert().FloatToString(userMap["user_id"].(float64), 'f', 0, 64)
	//保存session和cookie
	this.refreshSession(userId)
	this.jsonSuccess("恭喜，注册成功", nil, "/user/index")
}

//login save
func (this *AuthorController) Signup() {

	username := strings.Trim(this.GetString("username"), "")
	password := strings.Trim(this.GetString("password"), "")

	if(username == "") {
		this.jsonError("登录失败：用户名不能为空")
	}
	if(password == "") {
		this.jsonError("登录失败：密码不能为空")
	}

	data := map[string]string{
		"username":username,
	}
	user, err := business.NewBase().GetRequest("get_user_by_name", data)
	if(err != nil) {
		this.jsonError("登录失败：server error")
	}
	userMap := user.(map[string]interface{})
	if(len(userMap) == 0) {
		this.jsonError("登录失败：用户名不存在或密码错误")
	}
	if(userMap["is_forbidden"].(string) == "1") {
		this.jsonError("登录失败：该用户已被屏蔽")
	}
	if(userMap["password"].(string) != utils.NewEncrypt().Md5Encode(password)) {
		this.jsonError("登录失败：用户名或密码错误")
	}
	//保存session和cookie
	this.refreshSession(userMap["user_id"].(string))

	this.jsonSuccess("恭喜，登录成功", nil, "/user/index")
}

func (this *AuthorController) Logout() {
	passport := beego.AppConfig.String("login.passport")
	this.Ctx.SetCookie(passport, "")
	this.SetSession("user", nil)
	this.redirect("/");
}