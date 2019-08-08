package controllers

import (
	"anytunnel/at-admin/app/models"
	"fmt"
	"strconv"
	"time"

	"github.com/dustin/go-humanize"

	"github.com/astaxie/beego"
)

type MainController struct {
	BaseAdminController
}

func (this *MainController) Default() {
	if this.Ctx.Input.IsGet() {
		this.view("main/default")
	} else {
		data := map[string]interface{}{}
		db := models.G.DB("base")
		rs, err := db.Query(db.AR().Select("count(*) as total,cs_type").From("online").GroupBy("cs_type"))
		if err != nil {
			this.jsonError(err.Error())
		}
		d := rs.MapValues("cs_type", "total")
		clients := "0"
		servers := "0"
		if v, ok := d["client"]; ok {
			clients = v
		}
		if v, ok := d["server"]; ok {
			servers = v
		}
		rs, err = db.Query(db.AR().Select("count(*) as total").From("cluster").Where(map[string]interface{}{
			"update_time >": time.Now().Unix() - 60,
			"is_disable":    0,
			"is_delete":     0,
		}))
		if err != nil {
			this.jsonError(err.Error())
		}
		clusters := rs.Value("total")

		now := time.Now().Unix()

		rs, err = db.Query(db.AR().Select("sum(bytes_total) as total").From("package").Where(map[string]interface{}{
			"start_time <=": now,
			"end_time >=":   now,
		}))
		if err != nil {
			this.jsonError(err.Error())
		}
		traffic := rs.Value("total")
		if traffic == "" {
			traffic = "0"
		}
		rs, err = db.Query(db.AR().Select("sum(bytes_left) as total").From("package").Where(map[string]interface{}{
			"start_time <=": now,
			"end_time >=":   now,
		}))
		if err != nil {
			this.jsonError(err.Error())
		}
		trafficLeft := rs.Value("total")
		if trafficLeft == "" {
			trafficLeft = "0"
		}
		data["online"] = map[string]interface{}{
			"client":  clients,
			"server":  servers,
			"cluster": clusters,
		}
		_traffic, _ := strconv.Atoi(traffic)
		_trafficLeft, _ := strconv.Atoi(trafficLeft)
		_trafficUse := _traffic - _trafficLeft
		data["traffic"] = map[string]interface{}{
			"total":            traffic,
			"total_human":      humanize.Bytes(uint64(_traffic)),
			"total_left":       _trafficLeft,
			"total_left_human": humanize.Bytes(uint64(_trafficLeft)),
			"total_use":        _trafficUse,
			"total_use_human":  humanize.Bytes(uint64(_trafficUse)),
		}
		this.jsonSuccess("", data)
	}
}
func (this *MainController) Index() {
	this.Layout = "layout/main.html"
	this.TplName = "main/index.html"
	this.Data["title"] = fmt.Sprintf("%s - %s", beego.AppConfig.String("sys.name"), beego.AppConfig.String("sys.fullname"))
	privilege := models.Privilege{}
	navigators, menus, controllers, err := privilege.GetTypedPrivileges(this.loginUser["user_id"], "1")
	if err != nil {
		this.jsonError(err, "")
	}
	this.Data["navigators"] = navigators
	this.Data["menus"] = menus
	this.Data["controllers"] = controllers
	this.Render()
}
