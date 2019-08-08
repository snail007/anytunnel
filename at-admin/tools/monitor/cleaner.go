package main

import (
	"log"
	"time"

	"github.com/astaxie/beego"
)

//定时锁定流量用完的用户
func InitUserTrafficLocker() {
	go func() {
		log.Printf("user traffic monitor is running")
		sleep := beego.AppConfig.DefaultInt("user.scan.round.sleep.seconds", 60)
		lastID := "0"
		pagesize := beego.AppConfig.DefaultInt("user.scan.pagesize", 100)
		for {
			sql := db.AR().Select("user_id").From("user").Where(map[string]interface{}{
				"user_id >":    lastID,
				"is_active":    1,
				"is_forbidden": 0,
			}).Limit(pagesize).OrderBy("user_id", "asc")
			rs, err := db.Query(sql)
			if err != nil {
				log.Printf("query user ERR:%s", err)
				time.Sleep(time.Second * 30)
				continue
			}
			for _, userID := range rs.Values("user_id") {
				if userIsNoTraffic(userID) {
					killUser(userID, false, "")
				}
				lastID = userID
			}
			if rs.Len() == 0 {
				lastID = "0"
				time.Sleep(time.Second * time.Duration(sleep))
			}
		}
	}()
}

//定时清理无效的conn信息
func InitConnCleaner() {
	go func() {
		log.Printf("conn cleaner is running")
		for {
			db.Exec(db.AR().Delete("conn", map[string]interface{}{
				"update_time <": time.Now().Unix() - 300,
			}))
			time.Sleep(time.Second * time.Duration(beego.AppConfig.DefaultInt("conn.scan.round.sleep.seconds", 300)))
		}
	}()
}

//定时清理异常的隧道
func InitTunnelCleaner() {
	go func() {
		log.Printf("tunnel cleaner is running")
		sleep := beego.AppConfig.DefaultInt("tunnel.scan.round.sleep.seconds", 300)
		lastID := "0"
		pagesize := beego.AppConfig.DefaultInt("tunnel.scan.pagesize", 300)
		for {
			updateData := []map[string]interface{}{}
			sql := db.AR().Select("server_id,client_id,cluster_id,tunnel_id").From("tunnel").Where(map[string]interface{}{
				"tunnel_id >": lastID,
				"is_open":     1,
				"is_delete":   0,
			}).Limit(pagesize).OrderBy("tunnel_id", "asc")
			rs, err := db.Query(sql)
			if err != nil {
				log.Printf("query tunnel ERR:%s", err)
				time.Sleep(time.Second * 30)
				continue
			}
			if rs.Len() == 0 {
				lastID = "0"
				time.Sleep(time.Second * time.Duration(sleep))
			}
			for _, row := range rs.Rows() {
				if !csIsOnline(row["server_id"], "server") || !csIsOnline(row["client_id"], "client") {
					err := killTunnel(row["tunnel_id"], row["cluster_id"])
					if err == nil {
						updateData = append(updateData, map[string]interface{}{
							"tunnel_id": row["tunnel_id"],
							"is_open":   0,
							"status":    0,
						})
					}
				}
				if len(updateData) > 0 {
					sql := db.AR().UpdateBatch("tunnel", updateData, []string{"tunnel_id"})
					rs, err = db.Exec(sql)
					if err != nil {
						log.Printf("update tunnel ERR:%s", err)
						time.Sleep(time.Second * 30)
						continue
					}
				}
				lastID = row["tunnel_id"]
			}
		}
	}()
}

//定时清理异常的online
func InitOnlineCleaner() {
	go func() {
		log.Printf("online cleaner is running")
		sleep := beego.AppConfig.DefaultInt("online.scan.round.sleep.seconds", 300)
		lastID := "0"
		pagesize := beego.AppConfig.DefaultInt("online.scan.pagesize", 300)
		for {
			deleteData := []string{}
			sql := db.AR().Select("cluster_id,cs_id,cs_type,online_id").From("online").Where(map[string]interface{}{
				"online_id >": lastID,
			}).Limit(pagesize).OrderBy("online_id", "asc")
			rs, err := db.Query(sql)
			if err != nil {
				log.Printf("query online ERR:%s", err)
				time.Sleep(time.Second * 30)
				continue
			}
			if rs.Len() == 0 {
				lastID = "0"
				time.Sleep(time.Second * time.Duration(sleep))
			}
			for _, row := range rs.Rows() {
				if !csIsOnCluster(row["cluster_id"], row["cs_id"], row["cs_type"]) {
					deleteData = append(deleteData, row["online_id"])
				}
				if len(deleteData) > 0 {
					sql := db.AR().Delete("online", map[string]interface{}{
						"online_id": deleteData,
					})
					rs, err = db.Exec(sql)
					if err != nil {
						log.Printf("delete online ERR:%s", err)
						time.Sleep(time.Second * 30)
						continue
					}
					log.Printf("clean online %v", deleteData)
				}
				lastID = row["online_id"]
			}
		}
	}()
}
