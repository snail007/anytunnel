package controllers

import (
	system "anytunnel/at-admin/app/controllers"
)

const moduleName = "web"

type BaseController struct {
	system.BaseAdminController
}

func (this *BaseController) viewLayoutTitle(title, viewName, layout string) {
	this.ViewLayoutTitle(moduleName, title, viewName, layout)
}
func (this *BaseController) viewLayout(viewName, layout string) {
	this.ViewLayout(moduleName, viewName, layout)
}
func (this *BaseController) view(viewName string) {
	this.View(moduleName, viewName)
}
func (this *BaseController) viewTitle(title, viewName string) {
	this.ViewTitle(moduleName, title, viewName)
}
