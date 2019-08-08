package controllers

import (
	"strings"
	"github.com/astaxie/beego"
	"anytunnel/at-web/app/utils"
	"anytunnel/at-web/app/business"
	"fmt"
)

type JSONResponse struct {
	Code     int                    `json:"code"`
	Message  interface{}            `json:"message"`
	Data     interface{}            `json:"data"`
	Redirect map[string]interface{} `json:"redirect"`
}
type BaseUserController struct {
	BaseController
	loginUser map[string]interface{}
	controllerName string
	methodName string
}

//验证登录
func (this *BaseUserController) isLogin() bool {
	controllerName, actionName := this.GetControllerAndAction()
	this.controllerName = strings.ToLower(controllerName[0 : len(controllerName)-10])
	this.methodName = actionName
	//忽略 /author /error
	if(this.controllerName == "author" || this.controllerName == "error") {
		return true;
	}
	passport := beego.AppConfig.String("login.passport")
	cookie := this.Ctx.GetCookie(passport)
	//cookie 失效
	if(cookie == "") {
		return false
	}
	user := this.GetSession("user")
	//session 失效
	if(user == nil) {
		return false
	}
	cookieValue, _ := utils.NewEncrypt().Base64Decode(cookie)
	if(cookieValue == "") {
		return false
	}
	identifyList := strings.Split(cookieValue, "@")
	username := identifyList[0]
	identify := identifyList[1]
	userValue := user.(map[string]interface{})
	//对比cookie 和 session username
	if(username != userValue["username"].(string)) {
		return false
	}
	//对比客户端UAG and IP
	if(identify != utils.NewEncrypt().Md5Encode(this.Ctx.Request.UserAgent() + this.getClientIp() + userValue["password"].(string))) {
		return false
	}

	this.loginUser = userValue;
	//this.refreshSession()
	//success
	return true
}

func (this *BaseUserController) Prepare() {

	if !this.isLogin() {
		this.Redirect("/author/login", 302)
		this.StopRun()
	}

	this.Layout = "layout/default.html"
}
func (this *BaseUserController) checkAccess() {

}

//重置 session 和 cookie
func (this *BaseUserController) refreshSession(userId string)  {

	data := map[string]string {
		"user_id" : userId,
	}
	user, err := business.NewBase().GetRequest("get_user_by_id", data)
	if(err != nil) {
		fmt.Println("重置session错误")
	}
	if(len(user.(map[string]interface{})) > 0) {
		username := user.(map[string]interface{})["username"].(string)
		password := user.(map[string]interface{})["password"].(string)
		//重置 session
		this.SetSession("user", user)
		//重置 cookie
		identify := utils.NewEncrypt().Md5Encode(this.Ctx.Request.UserAgent() + this.getClientIp() + password)
		passportValue := utils.NewEncrypt().Base64Encode(username + "@" + identify)
		passport := beego.AppConfig.String("login.passport")
		cookieTime := beego.AppConfig.String("login.cookie.time")
		this.Ctx.SetCookie(passport, passportValue, cookieTime)
	}
	this.loginUser = user.(map[string]interface{});
}

//get user_id
func (this *BaseUserController) getUserId() string {
	return this.loginUser["user_id"].(string)
}

//get username
func (this *BaseUserController) getUsername() string {
	return this.loginUser["username"].(string)
}

//get email
func (this *BaseUserController) getEmail() string {
	return this.loginUser["email"].(string)
}
