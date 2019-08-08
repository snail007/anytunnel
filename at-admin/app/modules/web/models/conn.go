package models

import (
	"anytunnel/at-admin/app/utils"
	common "anytunnel/at-common"
	"fmt"

	"github.com/astaxie/beego"

	"github.com/snail007/go-activerecord/mysql"
)

type Conn struct {
}

func (p *Conn) GetConnByConnId(connId string) (conn map[string]string, err error) {
	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From("conn").Where(map[string]interface{}{
		"conn_id": connId,
	}))
	if err != nil {
		return
	}
	conn = rs.Row()
	return
}
func (p *Conn) HasTunnelRef(connId string) (has bool, err error) {
	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From("tunnel").Where(map[string]interface{}{
		"conn_id": connId,
	}).Limit(0, 1))
	if err != nil {
		return
	}
	if rs.Len() > 0 {
		has = true
	}
	return
}
func (p *Conn) Reset(connId string) (err error) {
	db := DB
	_, err = db.Exec(db.AR().Update("conn", map[string]interface{}{
		"token": utils.NewMisc().RandString(32),
	}, map[string]interface{}{
		"conn_id": connId,
	}))
	if err != nil {
		return
	}
	return
}
func (p *Conn) Offline(connId string) (err error) {
	db := DB
	rs, err := db.Query(db.AR().From("conn").Where(map[string]interface{}{
		"conn_id":   connId,
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
		"type":  "conn",
	}))
	if err != nil {
		return
	}
	if rs.Len() == 0 {
		return
	}
	clusterIP := rs.Value("ip")
	url := fmt.Sprintf("https://%s:%s/conn/offline/%s", clusterIP, beego.AppConfig.String("cluster.api.port"), token)
	body, code, err := common.HttpGet(url)
	if err != nil {
		return
	}
	if code != 200 {
		err = fmt.Errorf("access conn offline url %s fail,code:%d ,body:%s", url, code, body)
	}
	return
}
func (p *Conn) Delete(connId string) (err error) {
	db := DB
	_, err = db.Exec(db.AR().Update("conn", map[string]interface{}{
		"is_delete": 1,
	}, map[string]interface{}{
		"conn_id": connId,
	}))
	if err != nil {
		return
	}
	return
}
func (p *Conn) Insert(conn map[string]interface{}) (id int64, err error) {
	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Exec(db.AR().Insert("conn", conn))
	if err != nil {
		return
	}
	id = rs.LastInsertId
	return
}

func (p *Conn) Update(connId string, conn map[string]interface{}) (id int64, err error) {
	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Exec(db.AR().Update("conn", conn, map[string]interface{}{
		"conn_id": connId,
	}))
	if err != nil {
		return
	}
	id = rs.LastInsertId
	return
}

//根据user_id分页获取Conn
func (conn *Conn) GetConnsByUserIDAndLimit(col, keyword, orderby string, limit int, number int) (conns []map[string]string, err error) {
	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From("conn").Where(map[string]interface{}{
		col: keyword,
	}).OrderBy(orderby, "desc").Limit(limit, number))
	if err != nil {
		return
	}
	conns = rs.Rows()
	return
}

//分页获取Conn
func (conn *Conn) GetConnsByLimit(limit int, number int, orderby string) (conns []map[string]string, err error) {

	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From("conn").
		OrderBy(orderby, "desc").
		Limit(limit, number))
	if err != nil {
		return
	}
	conns = rs.Rows()

	return
}

func (conn *Conn) CountConns() (count int, err error) {

	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().Select("count(*) as total").From("conn"))
	if err != nil {
		return
	}
	count = utils.NewConvert().StringToInt(rs.Value("total"))
	return
}

func (conn *Conn) CountConnsByUserID(col, keyword string) (count int, err error) {

	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().
		Select("count(*) as total").
		From("conn").
		Where(map[string]interface{}{
			col: keyword,
		}))
	if err != nil {
		return
	}
	count = utils.NewConvert().StringToInt(rs.Value("total"))
	return
}
