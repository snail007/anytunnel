package controllers

import (
	"anytunnel/at-admin/app/models"
	"strings"
	"time"
)

type JSONResponse struct {
	Code     int                    `json:"code"`
	Message  interface{}            `json:"message"`
	Data     interface{}            `json:"data"`
	Redirect map[string]interface{} `json:"redirect"`
}
type BaseAdminController struct {
	BaseController
	loginUser map[string]string
}

func (this *BaseAdminController) Prepare() {
	if !this.isLogin() {
		this.Redirect("/login/index.html", 302)
		this.StopRun()
	}
	user := this.GetSession("author").(map[string]string)
	this.Data["loginUser"] = user
	this.loginUser = user
	this.checkAccess()
	this.Data["TimeNowYear"] = time.Now().Format("2006")
	this.Layout = "layout/default.html"
}
func (this *BaseAdminController) checkAccess() {
	controllerName, actionName := this.GetControllerAndAction()
	controllerName = strings.ToLower(controllerName[0 : len(controllerName)-10])
	actionName = strings.ToLower(actionName)
	if (controllerName == "main" && actionName == "index") || controllerName == "main" && actionName == "default" {
		return
	}
	//检查权限
	if "1" != this.loginUser["user_id"] {
		privilege := models.Privilege{}
		_, _, controllers, err := privilege.GetTypedPrivileges(this.loginUser["user_id"], "-1")
		if err != nil {
			this.jsonError(err.Error(), "")
		}
		found := false
		for _, c := range controllers {
			action := strings.ToLower(c["action"])
			if strings.Contains(action, ".") {
				action = action[:strings.LastIndex(action, ".")]
			}
			if controllerName == strings.ToLower(c["controller"]) && actionName == action {
				found = true
				break
			}
		}
		if !found {
			if this.IsAjax() {
				this.jsonError("您无权限进行此操作")
			} else {
				this.viewError("您无权限进行此操作")
			}

		}
	}
}
