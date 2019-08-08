package controllers

import (
	"anytunnel/at-admin/app/utils"
	"encoding/json"
	"strings"

	"github.com/astaxie/beego"
)

type BaseController struct {
	beego.Controller
}

//验证登录
func (this *BaseController) isLogin() bool {
	passport := beego.AppConfig.String("author.passport")
	//fmt.Println(passport)
	cookie := this.Ctx.GetCookie(passport)
	//cookie 失效
	if cookie == "" {
		//fmt.Println("cookie 失效")
		return false
	}
	user := this.GetSession("author")
	//fmt.Println(user)
	//session 失效
	if user == nil {
		//fmt.Println("session 失效")
		return false
	}
	encrypt := new(utils.Encrypt)
	cookieValue, _ := encrypt.Base64Decode(cookie)
	//fmt.Println("get cookie " + cookie)
	identifyList := strings.Split(cookieValue, "@")
	//fmt.Println(cookieValue, identifyList)
	if cookieValue == "" || len(identifyList) != 2 {
		//fmt.Println("cookieValue 无效")
		return false
	}
	name := identifyList[0]
	identify := identifyList[1]
	userValue := user.(map[string]string)

	//对比cookie 和 session name
	if name != userValue["username"] {
		//fmt.Println("对比cookie 和 session name 无效")
		//fmt.Println(userValue)
		return false
	}
	//对比客户端UAG and IP
	if identify != utils.NewEncrypt().Md5Encode(this.Ctx.Request.UserAgent()+this.getClientIp()+userValue["password"]) {
		//fmt.Println("对比客户端UAG and IP 无效")
		return false
	}
	//success
	return true
}
func (this *BaseController) viewLayoutTitle(title, viewName, layout string) {
	this.ViewLayoutTitle("", title, viewName, layout)
}
func (this *BaseController) viewLayout(viewName, layout string) {
	this.ViewLayout("", viewName, layout)
}
func (this *BaseController) view(viewName string) {
	this.View("", viewName)
}
func (this *BaseController) viewTitle(title, viewName string) {
	this.ViewTitle("", title, viewName)
}
func (this *BaseController) ViewLayoutTitle(module, title, viewName, layout string) {
	if module != "" {
		this.Layout = "layout/modules/" + module + "/" + layout + ".html"
		this.TplName = "modules/" + module + "/" + viewName + ".html"
	} else {
		this.Layout = "layout/" + layout + ".html"
		this.TplName = viewName + ".html"
	}
	this.Data["title"] = title
	this.Render()
}
func (this *BaseController) ViewLayout(module, viewName, layout string) {
	if module != "" {
		this.Layout = "layout/modules/" + module + "/" + layout + ".html"
		this.TplName = "modules/" + module + "/" + viewName + ".html"
	} else {
		this.Layout = "layout/" + layout + ".html"
		this.TplName = viewName + ".html"
	}
	this.Data["title"] = ""
	this.Render()
}
func (this *BaseController) View(module, viewName string) {
	if module != "" {
		this.Layout = "layout/modules/" + module + "/default.html"
		this.TplName = "modules/" + module + "/" + viewName + ".html"
	} else {
		this.Layout = "layout/default.html"
		this.TplName = viewName + ".html"
	}
	this.Data["title"] = ""
	this.Render()
}
func (this *BaseController) ViewTitle(module, title, viewName string) {
	if module != "" {
		this.Layout = "layout/modules/" + module + "/default.html"
		this.TplName = "modules/" + module + "/" + viewName + ".html"
	} else {
		this.Layout = "layout/default.html"
		this.TplName = viewName + ".html"
	}
	this.Data["title"] = title
	this.Render()
}
func (this *BaseController) ViewError(errorMessage string, data ...interface{}) {
	this.viewError(errorMessage, data...)
}
func (this *BaseController) viewError(errorMessage string, data ...interface{}) {
	this.Layout = "layout/default.html"
	errorType := "500"
	if len(data) > 0 {
		errorType = data[0].(string)
	}
	this.TplName = "error/" + errorType + ".html"
	this.Data["title"] = "system error"
	this.Data["errorMessage"] = errorMessage
	this.Render()
}
func (this *BaseController) JsonSuccess(message interface{}, data ...interface{}) {
	this.jsonSuccess(message, data...)
}

func (this *BaseController) jsonSuccess(message interface{}, data ...interface{}) {
	url := ""
	sleep := 500
	var _data interface{}
	if len(data) > 0 {
		_data = data[0]
	}
	if len(data) > 1 {
		url = data[1].(string)
	}
	if len(data) > 2 {
		sleep = data[2].(int)
	}

	this.Data["json"] = JSONResponse{
		Code:    1,
		Message: message,
		Data:    _data,
		Redirect: map[string]interface{}{
			"url":   url,
			"sleep": sleep,
		},
	}
	//this.ServeJSON()
	j, err := json.MarshalIndent(this.Data["json"], "", "\t")
	if err != nil {
		this.Abort(err.Error())
	} else {
		this.Abort(string(j))
	}

}
func (this *BaseController) JsonError(message interface{}, data ...interface{}) {
	this.jsonError(message, data...)
}
func (this *BaseController) jsonError(message interface{}, data ...interface{}) {
	url := ""
	sleep := 500
	var _data interface{}
	if len(data) > 0 {
		_data = data[0]
	}
	if len(data) > 1 {
		url = data[1].(string)
	}
	if len(data) > 2 {
		sleep = data[2].(int)
	}
	this.Data["json"] = JSONResponse{
		Code:    0,
		Message: message,
		Data:    _data,
		Redirect: map[string]interface{}{
			"url":   url,
			"sleep": sleep,
		},
	}
	j, err := json.MarshalIndent(this.Data["json"], "", " \t")
	if err != nil {
		this.Abort(err.Error())
	} else {
		this.Abort(string(j))
	}
}

//获取用户IP地址
func (this *BaseController) getClientIp() string {
	s := strings.Split(this.Ctx.Request.RemoteAddr, ":")
	return s[0]
}
