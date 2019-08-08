package routers

import (
	"anytunnel/at-admin/app/controllers"
	web_controllers "anytunnel/at-admin/app/modules/web/controllers"
	"anytunnel/at-admin/app/utils"
	"html/template"
	"math/rand"
	"net/http"

	"github.com/astaxie/beego"
)

func init() {
	beego.AppConfig.Set("author.passport", "usmpassport")
	beego.AppConfig.Set("sys.name", "ATCAS")
	beego.AppConfig.Set("sys.fullname", "AnyTunnel Cloud Admin System")

	beego.BConfig.ServerName = beego.AppConfig.String("sys.name")
	beego.SetStaticPath("/static/", "static")
	beego.BConfig.WebConfig.AutoRender = false
	beego.BConfig.WebConfig.Session.SessionName = "ssid"
	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.BConfig.RouterCaseSensitive = false
	beego.Router("/", &controllers.MainController{}, "*:Index")
	beego.AutoRouter(&controllers.MainController{})
	beego.AutoRouter(&controllers.SystemController{})
	beego.AutoRouter(&controllers.LoginController{})
	beego.AutoRouter(&controllers.UserController{})
	beego.AutoRouter(&controllers.PrivilegeController{})
	beego.AutoRouter(&controllers.Role_PrivilegeController{})
	beego.AutoRouter(&controllers.RoleController{})
	beego.AddFuncMap("randInt", randInt)
	beego.AddFuncMap("dateFormat", utils.NewDate().Format)
	beego.ErrorHandler("404", page_not_found)
	beego.ErrorHandler("500", page_not_found)
	beego.BConfig.WebConfig.ViewsPath = "app/views"
	userNS := beego.NewNamespace("/web",
		beego.NSAutoRouter(&web_controllers.UserController{}),
		beego.NSAutoRouter(&web_controllers.RoleController{}),
		beego.NSAutoRouter(&web_controllers.ServerController{}),
		beego.NSAutoRouter(&web_controllers.ClientController{}),
		beego.NSAutoRouter(&web_controllers.RegionController{}),
		beego.NSAutoRouter(&web_controllers.ClusterController{}),
		beego.NSAutoRouter(&web_controllers.TunnelController{}),
		beego.NSAutoRouter(&web_controllers.OnlineController{}),
		beego.NSAutoRouter(&web_controllers.PackageController{}),
		beego.NSAutoRouter(&web_controllers.ConnController{}),
		beego.NSAutoRouter(&web_controllers.IpListController{}),
		beego.NSAutoRouter(&web_controllers.AreaController{}),
	)
	beego.AddNamespace(userNS)
}
func page_not_found(rw http.ResponseWriter, r *http.Request) {
	t, _ := template.New("500-full.html").ParseFiles(beego.BConfig.WebConfig.ViewsPath + "/error/500-full.html")
	data := make(map[string]interface{})
	data["content"] = ""
	t.Execute(rw, data)
}

func randInt(start, end int) int {
	return rand.Intn(end) + start
}
