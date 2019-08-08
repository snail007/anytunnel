package models

import (
	"anytunnel/at-admin/app/utils"
	common "anytunnel/at-common"
	"fmt"

	"github.com/astaxie/beego"

	"github.com/snail007/go-activerecord/mysql"
)

type Server struct {
}

func (p *Server) GetServerByServerId(serverId string) (server map[string]string, err error) {
	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From("server").Where(map[string]interface{}{
		"server_id": serverId,
	}))
	if err != nil {
		return
	}
	server = rs.Row()
	return
}
func (p *Server) HasTunnelRef(serverId string) (has bool, err error) {
	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From("tunnel").Where(map[string]interface{}{
		"server_id": serverId,
	}).Limit(0, 1))
	if err != nil {
		return
	}
	if rs.Len() > 0 {
		has = true
	}
	return
}
func (p *Server) Reset(serverId string) (err error) {
	db := DB
	_, err = db.Exec(db.AR().Update("server", map[string]interface{}{
		"token": utils.NewMisc().RandString(32),
	}, map[string]interface{}{
		"server_id": serverId,
	}))
	if err != nil {
		return
	}
	return
}
func (p *Server) Offline(serverId string) (err error) {
	db := DB
	rs, err := db.Query(db.AR().From("server").Where(map[string]interface{}{
		"server_id": serverId,
		"is_delete": 0,
	}))
	if err != nil {
		return
	}
	if rs.Len() == 0 {
		return
	}
	token := rs.Value("token")
	rs, err = db.Query(db.AR().From("online").Where(map[string]interface{}{
		"token": token,
		"type":  "server",
	}))
	if err != nil {
		return
	}
	if rs.Len() == 0 {
		return
	}
	clusterIP := rs.Value("ip")
	url := fmt.Sprintf("https://%s:%s/server/offline/%s", clusterIP, beego.AppConfig.String("cluster.api.port"), token)
	body, code, err := common.HttpGet(url)
	if err != nil {
		return
	}
	if code != 200 {
		err = fmt.Errorf("access server offline url %s fail,code:%d ,body:%s", url, code, body)
	}
	return
}
func (p *Server) Delete(serverId string) (err error) {
	db := DB
	_, err = db.Exec(db.AR().Update("server", map[string]interface{}{
		"is_delete": 1,
	}, map[string]interface{}{
		"server_id": serverId,
	}))
	if err != nil {
		return
	}
	return
}
func (p *Server) Insert(server map[string]interface{}) (id int64, err error) {
	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Exec(db.AR().Insert("server", server))
	if err != nil {
		return
	}
	id = rs.LastInsertId
	return
}

func (p *Server) Update(serverId string, server map[string]interface{}) (id int64, err error) {
	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Exec(db.AR().Update("server", server, map[string]interface{}{
		"server_id": serverId,
	}))
	if err != nil {
		return
	}
	id = rs.LastInsertId
	return
}

//根据user_id分页获取Server
func (server *Server) GetServersByUserIDAndLimit(userID string, limit int, number int) (servers []map[string]string, err error) {
	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From("server").Where(map[string]interface{}{
		"user_id":   userID,
		"is_delete": 0,
	}).OrderBy("server_id", "desc").Limit(limit, number))
	if err != nil {
		return
	}
	servers = rs.Rows()
	return
}

//分页获取Server
func (server *Server) GetServersByLimit(limit int, number int) (servers []map[string]string, err error) {

	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From("server").Where(map[string]interface{}{
		"is_delete":  0,
		"user_id <>": 0,
	}).
		OrderBy("server_id", "desc").
		Limit(limit, number))
	if err != nil {
		return
	}
	servers = rs.Rows()

	return
}

func (server *Server) CountServers() (count int, err error) {

	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().Select("count(*) as total").Where(map[string]interface{}{
		"is_delete":  0,
		"user_id <>": 0,
	}).From("server"))
	if err != nil {
		return
	}
	count = utils.NewConvert().StringToInt(rs.Value("total"))
	return
}

func (server *Server) CountServersByUserID(userID string) (count int, err error) {

	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().
		Select("count(*) as total").
		From("server").
		Where(map[string]interface{}{
			"user_id":   userID,
			"is_delete": 0,
		}))
	if err != nil {
		return
	}
	count = utils.NewConvert().StringToInt(rs.Value("total"))
	return
}
