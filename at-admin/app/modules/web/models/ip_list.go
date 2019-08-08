package models

import (
	"anytunnel/at-admin/app/utils"

	"github.com/snail007/go-activerecord/mysql"
)

type IpList struct {
}

func (p *IpList) GetIpListByIpListId(ip_listId string) (ip_list map[string]string, err error) {
	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From("ip_list").Where(map[string]interface{}{
		"ip_list_id": ip_listId,
	}))
	if err != nil {
		return
	}
	ip_list = rs.Row()
	return
}

func (p *IpList) Delete(ip_listId string) (err error) {
	db := DB
	_, err = db.Exec(db.AR().Delete("ip_list", map[string]interface{}{
		"ip_list_id": ip_listId,
	}))
	if err != nil {
		return
	}
	return
}
func (p *IpList) Insert(ip_list map[string]interface{}) (id int64, err error) {
	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Exec(db.AR().Insert("ip_list", ip_list))
	if err != nil {
		return
	}
	id = rs.LastInsertId
	return
}

func (p *IpList) Update(ip_listId string, ip_list map[string]interface{}) (id int64, err error) {
	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Exec(db.AR().Update("ip_list", ip_list, map[string]interface{}{
		"ip_list_id": ip_listId,
	}))
	if err != nil {
		return
	}
	id = rs.LastInsertId
	return
}

//根据搜索分页获取IpList
func (ip_list *IpList) GetIpListsByIDAndLimit(column, id string, limit int, number int) (ip_lists []map[string]string, err error) {

	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From("ip_list").Where(map[string]interface{}{
		column: id,
	}).OrderBy("ip_list_id", "desc").Limit(limit, number))
	if err != nil {
		return
	}
	ip_lists = rs.Rows()
	return
}

//分页获取IpList
func (ip_list *IpList) GetIpListsByLimit(limit int, number int) (ip_lists []map[string]string, err error) {

	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From("ip_list").
		OrderBy("ip_list_id", "desc").
		Limit(limit, number))
	if err != nil {
		return
	}
	ip_lists = rs.Rows()
	return
}

func (ip_list *IpList) CountIpLists() (count int, err error) {

	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().Select("count(*) as total").From("ip_list"))
	if err != nil {
		return
	}
	count = utils.NewConvert().StringToInt(rs.Value("total"))
	return
}

func (ip_list *IpList) CountIpListsByID(column, id string) (count int, err error) {

	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().
		Select("count(*) as total").
		From("ip_list").
		Where(map[string]interface{}{
			column: id,
		}))
	if err != nil {
		return
	}
	count = utils.NewConvert().StringToInt(rs.Value("total"))
	return
}
