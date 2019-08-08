package models

import (
	"anytunnel/at-admin/app/utils"

	"github.com/snail007/go-activerecord/mysql"
)

type Online struct {
}

func (p *Online) GetOnlineByOnlineId(onlineId string) (online map[string]string, err error) {
	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From("online").Where(map[string]interface{}{
		"online_id": onlineId,
	}))
	if err != nil {
		return
	}
	online = rs.Row()
	return
}

func (p *Online) Delete(onlineId string) (err error) {
	db := DB
	_, err = db.Exec(db.AR().Update("online", map[string]interface{}{
		"is_delete": 1,
	}, map[string]interface{}{
		"online_id": onlineId,
	}))
	if err != nil {
		return
	}
	return
}
func (p *Online) Insert(online map[string]interface{}) (id int64, err error) {
	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Exec(db.AR().Insert("online", online))
	if err != nil {
		return
	}
	id = rs.LastInsertId
	return
}

func (p *Online) Update(onlineId string, online map[string]interface{}) (id int64, err error) {
	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Exec(db.AR().Update("online", online, map[string]interface{}{
		"online_id": onlineId,
	}))
	if err != nil {
		return
	}
	id = rs.LastInsertId
	return
}

//根据搜索分页获取Online
func (online *Online) GetOnlinesByIDAndLimit(cs, column, id string, limit int, number int) (onlines []map[string]string, err error) {

	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From("online").Where(map[string]interface{}{
		column:    id,
		"cs_type": cs,
	}).OrderBy("online_id", "desc").Limit(limit, number))
	if err != nil {
		return
	}
	onlines = rs.Rows()

	return
}

//分页获取Online
func (online *Online) GetOnlinesByLimit(cs string, limit int, number int) (onlines []map[string]string, err error) {

	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From("online").Where(map[string]interface{}{
		"cs_type": cs,
	}).
		OrderBy("online_id", "desc").
		Limit(limit, number))
	if err != nil {
		return
	}
	onlines = rs.Rows()

	return
}

func (online *Online) CountOnlines(cs string) (count int, err error) {

	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().Select("count(*) as total").Where(map[string]interface{}{
		"cs_type": cs,
	}).From("online"))
	if err != nil {
		return
	}
	count = utils.NewConvert().StringToInt(rs.Value("total"))
	return
}

func (online *Online) CountOnlinesByID(cs, column, id string) (count int, err error) {

	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().
		Select("count(*) as total").
		From("online").
		Where(map[string]interface{}{
			column:    id,
			"cs_type": cs,
		}))
	if err != nil {
		return
	}
	count = utils.NewConvert().StringToInt(rs.Value("total"))
	return
}
