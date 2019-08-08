package controllers

import (
	"encoding/json"
	"strings"

	"github.com/astaxie/beego"
)

type BaseController struct {
	beego.Controller
}

func (this *BaseController) viewLayoutTitle(title, viewName, layout string) {
	this.Layout = "layout/" + layout + ".html"
	this.TplName = viewName + ".html"
	this.Data["title"] = title
	this.Render()
}
func (this *BaseController) viewLayout(viewName, layout string) {
	this.Layout = "layout/" + layout + ".html"
	this.TplName = viewName + ".html"
	this.Data["title"] = ""
	this.Render()
}
func (this *BaseController) view(viewName string) {
	this.Layout = "layout/default.html"
	this.TplName = viewName + ".html"
	this.Data["title"] = ""
	this.Render()
}
func (this *BaseController) viewError(errorMessage string, data ...interface{}) {
	this.Layout = "layout/page.html"
	errorType := "500"
	if len(data) > 0 {
		errorType = data[0].(string)
	}
	this.TplName = "error/" + errorType + ".html"
	this.Data["title"] = "system error"
	this.Data["errorMessage"] = errorMessage
	this.Render()
}
func (this *BaseController) viewTitle(title, viewName string) {
	this.Layout = "layout/default.html"
	this.TplName = viewName + ".html"
	this.Data["title"] = title
	this.Render()
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

//302跳转
func (this *BaseController) redirect(url string) {
	this.Redirect(url, 302)
	this.StopRun()
}
