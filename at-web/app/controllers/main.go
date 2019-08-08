package controllers

type MainController struct {
	BaseController
}

func (this *MainController) Index() {
	isLogin := "1"
	if(this.GetSession("user") == nil) {
		isLogin = "0"
	}
	this.Data["isLogin"] = isLogin
	this.viewLayoutTitle("AnyTunnelCloud", "web/index", "index")
}
