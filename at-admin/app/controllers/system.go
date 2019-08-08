package controllers

type SystemController struct {
	BaseAdminController
}

func (this *SystemController) Base() {
	this.Layout = "layout/page.html"
	this.TplName = "system/base.html"
	this.Render()
}
