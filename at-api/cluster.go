package main

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/snail007/go-activerecord/mysql"
)

type Cluster struct{}

func NewCluster() *Cluster {
	return &Cluster{}
}

//report a cluster status
//method : POST
//params : sys_conn_number,bandwidth,tunnel_conn_number
func (this *Cluster) Report(responseWrite http.ResponseWriter, request *http.Request, params httprouter.Params) {

	_sysConn := request.PostFormValue("sys_conn_number")
	sysConn, err := strconv.Atoi(_sysConn)
	if err != nil {
		jsonError(responseWrite, "sys_conn_number is error!", nil)
		return
	}
	_bandwidth := request.PostFormValue("bandwidth")
	bandwidth, err := strconv.Atoi(_bandwidth)
	if err != nil {
		jsonError(responseWrite, "bandwidth is error!", nil)
		return
	}
	_tunnelConn := request.PostFormValue("tunnel_conn_number")
	tunnelConn, err := strconv.Atoi(_tunnelConn)
	if err != nil {
		jsonError(responseWrite, "tunnel_conn_number is error!", nil)
		return
	}
	addr := request.RemoteAddr
	ip := addr[0:strings.Index(addr, ":")]
	clusterValue := map[string]interface{}{
		"ip":                 ip,
		"system_conn_number": sysConn,
		"tunnel_conn_number": tunnelConn,
		"bandwidth":          bandwidth,
	}
	where := map[string]interface{}{
		"ip":        ip,
		"is_delete": 0,
	}
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From("cluster").Where(where))
	if err != nil {
		jsonError(responseWrite, err.Error(), nil)
		log.Warnf("Cluster Report ERR:%s", err)
		return
	}
	if rs.Len() > 0 {
		clusterValue["update_time"] = time.Now().Unix()
		_, err = db.Exec(db.AR().Update("cluster", clusterValue, where))
	} else {
		clusterValue["create_time"] = time.Now().Unix()
		_, err = db.Exec(db.AR().Insert("cluster", clusterValue))
	}
	if err != nil {
		jsonError(responseWrite, err.Error(), nil)
		log.Warnf("Cluster Report ERR:%s", err)
		return
	}
	responseWrite.WriteHeader(http.StatusNoContent)
}

//cluster list
//method : POST
//params : sys_conn_number,bandwidth,tunnel_conn_number
func (this *Cluster) List(responseWrite http.ResponseWriter, request *http.Request, params httprouter.Params) {

	_sysConn := request.PostFormValue("sys_conn_number")
	sysConn, err := strconv.Atoi(_sysConn)
	if err != nil {
		jsonError(responseWrite, "sys_conn_number is error!", nil)
		return
	}
	_bandwidth := request.PostFormValue("bandwidth")
	bandwidth, err := strconv.Atoi(_bandwidth)
	if err != nil {
		jsonError(responseWrite, "bandwidth is error!", nil)
		return
	}
	_tunnelConn := request.PostFormValue("tunnel_conn_number")
	tunnelConn, err := strconv.Atoi(_tunnelConn)
	if err != nil {
		jsonError(responseWrite, "tunnel_conn_number is error!", nil)
		return
	}
	addr := request.RemoteAddr
	ip := addr[0:strings.Index(addr, ":")]
	clusterValue := map[string]interface{}{
		"ip":                 ip,
		"system_conn_number": sysConn,
		"tunnel_conn_number": tunnelConn,
		"bandwidth":          bandwidth,
	}
	where := map[string]interface{}{
		"ip":        ip,
		"is_delete": 0,
	}
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From("cluster").Where(where))
	if err != nil {
		jsonError(responseWrite, err.Error(), nil)
		log.Warnf("Cluster Report ERR:%s", err)
		return
	}
	if rs.Len() > 0 {
		clusterValue["update_time"] = time.Now().Unix()
		_, err = db.Exec(db.AR().Update("cluster", clusterValue, where))
	} else {
		clusterValue["create_time"] = time.Now().Unix()
		_, err = db.Exec(db.AR().Insert("cluster", clusterValue))
	}
	if err != nil {
		jsonError(responseWrite, err.Error(), nil)
		log.Warnf("Cluster Report ERR:%s", err)
		return
	}
	responseWrite.WriteHeader(http.StatusNoContent)
}