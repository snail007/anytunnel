package main

import (
	_ "anytunnel/at-admin/app/routers"

	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}
