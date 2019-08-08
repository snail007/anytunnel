package routers

import (
	"anytunnel/at-web/app/controllers"
	"anytunnel/at-web/app/utils"
	"html/template"
	"math/rand"
	"net/http"

	"github.com/astaxie/beego"
)

func init() {
	beego.AppConfig.Set("sys.name", "ATC")
	beego.AppConfig.Set("sys.fullname", "AnyTunnel Cloud")

	beego.BConfig.ServerName = beego.AppConfig.String("sys.name")
	beego.SetStaticPath("/static/", "static")
	beego.BConfig.WebConfig.AutoRender = false
	beego.BConfig.WebConfig.Session.SessionName = "ssidw"
	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.BConfig.RouterCaseSensitive = false
	beego.AutoRouter(&controllers.MainController{})
	beego.AutoRouter(&controllers.AuthorController{})
	beego.AutoRouter(&controllers.UserController{})
	beego.AutoRouter(&controllers.TunnelController{})
	beego.AutoRouter(&controllers.ClientController{})
	beego.AutoRouter(&controllers.ServerController{})
	beego.Router("/", &controllers.MainController{}, "*:Index")
	beego.AddFuncMap("randInt", randInt)
	beego.AddFuncMap("dateFormat", utils.NewDate().Format)
	beego.ErrorHandler("404", page_not_found)
	beego.ErrorHandler("500", page_not_found)
	beego.BConfig.WebConfig.ViewsPath = "app/views"
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
