package models

import (
	"anytunnel/at-admin/app/utils"

	"github.com/snail007/go-activerecord/mysql"
)

type Area struct {
}

func (p *Area) GetAreaByAreaId(areaId string) (area map[string]string, err error) {
	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From("area").Where(map[string]interface{}{
		"area_id": areaId,
	}))
	if err != nil {
		return
	}
	area = rs.Row()
	return
}

func (p *Area) Delete(areaId string) (err error) {
	db := DB
	_, err = db.Exec(db.AR().Delete("area", map[string]interface{}{
		"area_id": areaId,
	}))
	if err != nil {
		return
	}
	return
}
func (p *Area) Insert(area map[string]interface{}) (id int64, err error) {
	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Exec(db.AR().Insert("area", area))
	if err != nil {
		return
	}
	id = rs.LastInsertId
	return
}

func (p *Area) Update(areaId string, area map[string]interface{}) (id int64, err error) {
	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Exec(db.AR().Update("area", area, map[string]interface{}{
		"area_id": areaId,
	}))
	if err != nil {
		return
	}
	id = rs.LastInsertId
	return
}

//根据搜索分页获取Area
func (area *Area) GetAreasByIDAndLimit(column, id string, limit int, number int) (areas []map[string]string, err error) {

	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From("area").Where(map[string]interface{}{
		column: id,
	}).OrderBy("area_id", "desc").Limit(limit, number))
	if err != nil {
		return
	}
	areas = rs.Rows()
	return
}

//分页获取Area
func (area *Area) GetAreasByLimit(limit int, number int) (areas []map[string]string, err error) {

	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From("area").
		OrderBy("area_id", "desc").
		Limit(limit, number))
	if err != nil {
		return
	}
	areas = rs.Rows()
	return
}

func (area *Area) CountAreas() (count int, err error) {

	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().Select("count(*) as total").From("area"))
	if err != nil {
		return
	}
	count = utils.NewConvert().StringToInt(rs.Value("total"))
	return
}

func (area *Area) CountAreasByID(column, id string) (count int, err error) {

	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().
		Select("count(*) as total").
		From("area").
		Where(map[string]interface{}{
			column: id,
		}))
	if err != nil {
		return
	}
	count = utils.NewConvert().StringToInt(rs.Value("total"))
	return
}
