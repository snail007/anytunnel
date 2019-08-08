package main

import (
	utils "anytunnel/at-common"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/astaxie/beego"
)

func userIsNoTraffic(userID string) bool {
	now := time.Now().Unix()
	rs, err := db.Query(db.AR().Select("user_id").From("package").Where(map[string]interface{}{
		"user_id":      userID,
		"start_time <": now,
		"end_time >":   now,
		"bytes_left >": 0,
	}).Limit(1))
	if err != nil {
		log.Printf("query user package ERR:%s", err)
		return false
	}
	return rs.Len() == 0
}
func killUser(userID string, forbidden bool, reason string) {
	//log.Printf("kill user : %s", userID)
	//1.kill all user's online server
	rs, err := db.Query(db.AR().Select("cs_id,cluster_id").From("online").Where(map[string]interface{}{
		"user_id": userID,
		"cs_type": "server",
	}))
	if err != nil {
		return
	}
	allOnlineServerIDs := rs.Rows()
	for _, online := range allOnlineServerIDs {
		killCS(online["cluster_id"], online["cs_id"], true)
	}
	//2.kill all user's online client
	rs, err = db.Query(db.AR().Select("cs_id,cluster_id").From("online").Where(map[string]interface{}{
		"user_id": userID,
		"cs_type": "client",
	}))
	if err != nil {
		return
	}
	allOnlineClientIDs := rs.Rows()
	for _, online := range allOnlineClientIDs {
		killCS(online["cluster_id"], online["cs_id"], false)
	}
	//3.kill all user's online tunnel
	killUserTunnel(userID)
	//4.if user forbidden or not
	if forbidden {
		_, err := db.Exec(db.AR().Update("user", map[string]interface{}{"is_forbidden": 1, "forbidden_reason": reason}, map[string]interface{}{"user_id": userID}))
		if err != nil {
			return
		}
	}
}

func killCS(clusterID, csID string, isServer bool) (err error) {
	typ := "server"
	if !isServer {
		typ = "client"
	}
	log.Printf("kill %s : %s", typ, csID)
	rs, err := db.Query(db.AR().From("cluster").Where(map[string]interface{}{"is_delete": 0}))
	if err != nil {
		return
	}
	if rs.Len() == 0 {
		err = fmt.Errorf("empty cluster for cluster_id:%s", clusterID)
		return
	}
	cluster := rs.Row()

	rs, err = db.Query(db.AR().From(typ).Where(map[string]interface{}{
		typ + "_id": csID,
		"is_delete": 0,
	}))
	if err != nil {
		return
	}
	if rs.Len() == 0 {
		err = fmt.Errorf("empty %s for %s_id:%s", typ, typ, csID)
		return
	}
	item := rs.Row()
	url := fmt.Sprintf("https://%s:%s/%s/offline/%s", cluster["ip"], beego.AppConfig.String("cluster.api.port"), typ, item["token"])
	body, code, err := utils.HttpGet(url)
	if err != nil {
		log.Printf("kill server fail,code: %d,body: %s ,ERR: %s ", code, string(body), err)
	}
	return
}
func killUserTunnel(userID string) {

	rs, err := db.Query(db.AR().From("conn").Where(map[string]interface{}{"user_id": userID}))
	if err != nil {
		return
	}
	if rs.Len() == 0 {
		//err = fmt.Errorf("killUserTunnel , empty tunnel for user_id:%s", userID)
		return
	}
	for _, conn := range rs.Rows() {
		log.Printf("kill tunnel:%s , userID: %s", conn["tunnel_id"], userID)
		rs, err = db.Query(db.AR().From("cluster").Where(map[string]interface{}{
			"cluster_id": conn["cluster_id"],
			"is_delete":  0,
		}))
		if err != nil {
			return
		}
		if rs.Len() == 0 {
			//err = fmt.Errorf("killUserTunnel , empty cluster for cluster_id:%s", conn["cluster_id"])
			return
		}
		cluster := rs.Row()
		url := fmt.Sprintf("https://%s:%s/port/close/%s", cluster["ip"], beego.AppConfig.String("cluster.api.port"), conn["tunnel_id"])
		body, code, err := utils.HttpGet(url)
		if err != nil {
			log.Printf("kill tunnel fail,code: %d,body: %s ,ERR: %s ", code, string(body), err)
			return
		}
	}
	db.Exec(db.AR().Delete("conn", map[string]interface{}{
		"user_id": userID,
	}))
}

func csIsOnline(csID, csType string) bool {
	rs, err := db.Query(db.AR().From("online").Where(map[string]interface{}{
		"cs_id":   csID,
		"cs_type": csType,
	}))
	if err != nil {
		return false
	}
	return rs.Len() == 1
}
func killTunnel(tunnelID, clusterID string) (err error) {
	rs, err := db.Query(db.AR().From("cluster").Where(map[string]interface{}{
		"cluster_id": clusterID,
		"is_delete":  0,
	}))
	if err != nil {
		log.Printf("kill tunnel fail,ERR:%s", err)
		return
	}
	if rs.Len() == 0 {
		return
	}
	cluster := rs.Row()
	url := fmt.Sprintf("https://%s:%s/port/close/%s", cluster["ip"], beego.AppConfig.String("cluster.api.port"), tunnelID)
	body, code, err := utils.HttpGet(url)
	if err != nil {
		log.Printf("kill tunnel fail,code: %d,body: %s ,ERR: %s ", code, string(body), err)
		return
	}
	return
}
func GetCS(csID, csType string) (row map[string]string, err error) {
	rs, err := db.Query(db.AR().From(csType).Where(map[string]interface{}{
		csType + "_id": csID,
		"is_delete":    0,
	}))
	if err != nil {
		log.Printf("query %s fail,ERR:%s", csType, err)
		return
	}
	if rs.Len() == 0 {
		err = fmt.Errorf("%s not exists", csType)
		return
	}
	row = rs.Row()
	return
}
func GetCluster(clusterID string) (row map[string]string, err error) {
	rs, err := db.Query(db.AR().From("cluster").Where(map[string]interface{}{
		"cluster_id": clusterID,
		"is_delete":  0,
	}))
	if err != nil {
		log.Printf("query clusterfail,ERR:%s", err)
		return
	}
	if rs.Len() == 0 {
		err = fmt.Errorf("cluster not exists ,clusterID:%s", clusterID)
		return
	}
	row = rs.Row()
	return
}
func csIsOnCluster(clusterID, csID, csType string) (exists bool) {
	cluster, err := GetCluster(clusterID)
	if err != nil {
		return
	}
	cs, err := GetCS(csID, csType)
	if err != nil {
		return
	}
	url := fmt.Sprintf("https://%s:%s/%s/status/%s", cluster["ip"], beego.AppConfig.String("cluster.api.port"), csType, cs["token"])
	body, code, err := utils.HttpGet(url)
	if err != nil {

		log.Printf("get %s status fail,code: %d,body: %s ,ERR: %s ", csType, code, string(body), err)
		return
	}
	res := map[string]interface{}{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Printf("parse csIsOnCluster result ERR:%s,body:%s", err, string(body))
	}
	respCode, err := strconv.ParseUint(fmt.Sprintf("%.0f", res["code"]), 10, 64)
	if err != nil {
		return
	}
	return respCode == 1
}
