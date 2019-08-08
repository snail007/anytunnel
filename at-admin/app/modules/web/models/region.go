package models

import (
	"github.com/snail007/go-activerecord/mysql"
)

type Region struct {
}

func (p *Region) GetRegionByRegionId(regionId string) (region map[string]string, err error) {
	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From("region").Where(map[string]interface{}{
		"region_id": regionId,
		"is_delete": 0,
	}))
	if err != nil {
		return
	}
	region = rs.Row()
	return
}
func (p *Region) HasRole(regionId string) (has bool, err error) {
	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From("role_region").Where(map[string]interface{}{
		"region_id": regionId,
	}).Limit(0, 1))
	if err != nil {
		return
	}
	if rs.Len() > 0 {
		has = true
	}
	return
}
func (p *Region) HasSubRegion(regionId string) (has bool, err error) {
	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From("region").Where(map[string]interface{}{
		"parent_id": regionId,
		"is_delete": 0,
	}).Limit(0, 1))
	if err != nil {
		return
	}
	if rs.Len() > 0 {
		has = true
	}
	return
}
func (p *Region) HasCluster(regionId string) (has bool, err error) {
	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From("region_cluster").Where(map[string]interface{}{
		"region_id": regionId,
	}).Limit(0, 1))
	if err != nil {
		return
	}
	if rs.Len() > 0 {
		has = true
	}
	return
}
func (p *Region) Delete(regionId string) (err error) {
	db := DB
	_, err = db.Exec(db.AR().Update("region", map[string]interface{}{
		"is_delete": 1,
	}, map[string]interface{}{
		"region_id": regionId,
	}))
	if err != nil {
		return
	}
	return
}
func (p *Region) Insert(region map[string]interface{}) (id int64, err error) {
	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Exec(db.AR().Insert("region", region))
	if err != nil {
		return
	}
	id = rs.LastInsertId
	return
}

func (p *Region) Update(regionId string, region map[string]interface{}) (id int64, err error) {
	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Exec(db.AR().Update("region", region, map[string]interface{}{
		"region_id": regionId,
	}))
	if err != nil {
		return
	}
	id = rs.LastInsertId
	return
}

//获取所有的Region
func (this *Region) GetAllRegions() (regions []map[string]string, err error) {
	db := DB
	res, err := db.Query(db.AR().From("region").Where(map[string]interface{}{
		"is_delete": 0,
	}).OrderBy("parent_id", "ASC").OrderBy("region_id", "ASC"))
	if err != nil {
		return
	}
	regions = res.Rows()
	return
}

//获取所有顶级Region
func (this *Region) GetTopRegions() (regions map[string]map[string]string, err error) {
	db := DB
	res, err := db.Query(db.AR().From("region").Where(map[string]interface{}{
		"is_delete": 0,
		"parent_id": 0,
	}).OrderBy("region_id", "ASC"))
	if err != nil {
		return
	}
	regions = res.MapRows("region_id")
	return
}

//获取所有二级Region
func (this *Region) GetSubRegions() (regions []map[string]string, err error) {
	db := DB
	res, err := db.Query(db.AR().From("region").Where(map[string]interface{}{
		"is_delete": 0,
	}).OrderBy("parent_id", "ASC"))
	if err != nil {
		return
	}
	regionsAll := res.Rows()
	for _, r := range regionsAll {
		if r["parent_id"] != "0" {
			r["parent"] = ""
			for _, r1 := range regionsAll {
				if r["parent_id"] == r1["region_id"] {
					r["parent"] = r1["name"]
					break
				}
			}
			regions = append(regions, r)
		}
	}
	return
}
