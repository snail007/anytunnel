package main

//online api
import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/snail007/go-activerecord/mysql"
)

type Online struct{}

func NewOnline() *Online {
	return &Online{}
}

//get online client
//method : GET
//params : cluster_id
func (online *Online) GetClientByClusterId(responseWrite http.ResponseWriter, request *http.Request, params httprouter.Params) {

	clusterId := request.FormValue("cluster_id")

	if clusterId == "" {
		jsonError(responseWrite, "cluster_id错误!", nil)
		return
	}

	db := G.DB()
	var rs *mysql.ResultSet
	rs, err := db.Query(db.AR().From("online").Where(map[string]interface{}{
		"cluster_id": clusterId,
		"cs_type": "client",
	}))
	if err != nil {
		jsonError(responseWrite, err.Error(), nil)
		return
	}

	if rs.Len() == 0 {
		jsonError(responseWrite, "client不存在", nil)
		return
	}

	jsonSuccess(responseWrite, "ok", rs.Rows())
}

//get online server
//method : GET
//params : cluster_id
func (online *Online) GetServerByClusterId(responseWrite http.ResponseWriter, request *http.Request, params httprouter.Params) {

	clusterId := request.FormValue("cluster_id")

	if clusterId == "" {
		jsonError(responseWrite, "cluster_id错误!", nil)
		return
	}

	db := G.DB()
	var rs *mysql.ResultSet
	rs, err := db.Query(db.AR().From("online").Where(map[string]interface{}{
		"cluster_id": clusterId,
		"cs_type": "server",
	}))
	if err != nil {
		jsonError(responseWrite, err.Error(), nil)
		return
	}

	if rs.Len() == 0 {
		jsonError(responseWrite, "server 不存在", nil)
		return
	}

	jsonSuccess(responseWrite, "ok", rs.Rows())
}
