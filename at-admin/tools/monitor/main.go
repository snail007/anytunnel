package main

import (
	"anytunnel/at-admin/app/models"
)

var db = models.G.DB("base")

func main() {
	//定时清理无效的conn信息
	InitConnCleaner()
	//定时关闭异常的隧道
	InitTunnelCleaner()
	//定时清理异常的online
	InitOnlineCleaner()
	//定时锁定流量用完的用户
	InitUserTrafficLocker()
	select {}
}
