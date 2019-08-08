package main

import (
	_ "anytunnel/at-web/app/routers"

	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}
