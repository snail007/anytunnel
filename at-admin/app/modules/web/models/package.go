package models

import (
	"anytunnel/at-admin/app/utils"
	"strconv"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/snail007/go-activerecord/mysql"
)

type PackageModel struct {
}

func (p *PackageModel) GetPackageByPackageId(packageId string) (PackageModel map[string]string, err error) {
	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From("package").Where(map[string]interface{}{
		"package_id": packageId,
	}))
	if err != nil {
		return
	}
	PackageModel = rs.Row()
	return
}

func (p *PackageModel) Delete(packageId string) (err error) {
	db := DB
	_, err = db.Exec(db.AR().Update("package", map[string]interface{}{
		"is_delete": 1,
	}, map[string]interface{}{
		"package_id": packageId,
	}))
	if err != nil {
		return
	}
	return
}
func (p *PackageModel) Insert(_package map[string]interface{}) (id int64, err error) {
	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Exec(db.AR().Insert("package", _package))
	if err != nil {
		return
	}
	id = rs.LastInsertId
	return
}

func (p *PackageModel) Update(packageId string, _package map[string]interface{}) (id int64, err error) {
	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Exec(db.AR().Update("package", _package, map[string]interface{}{
		"package_id": packageId,
	}))
	if err != nil {
		return
	}
	id = rs.LastInsertId
	return
}

//根据搜索分页获取Package
func (p *PackageModel) GetPackagesByIDAndLimit(column, id string, limit int, number int) (packages []map[string]string, err error) {

	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From("package").Where(map[string]interface{}{
		column: id,
	}).OrderBy("end_time", "asc").Limit(limit, number))
	if err != nil {
		return
	}
	packages = rs.Rows()
	for k, v := range packages {
		_total, _ := strconv.ParseUint(v["bytes_total"], 10, 64)
		_totalLeft, _ := strconv.ParseUint(v["bytes_left"], 10, 64)
		_totalUse := _total - _totalLeft
		packages[k]["bytes_total_human"] = humanize.Bytes(_total)
		packages[k]["bytes_left_human"] = humanize.Bytes(_totalLeft)
		packages[k]["bytes_use_human"] = humanize.Bytes(_totalUse)
		packages[k]["status"] = `<span class="label label-success">有效</span>`
		endtime, _ := strconv.Atoi(v["end_time"])
		starttime, _ := strconv.Atoi(v["start_time"])
		if int64(endtime) <= time.Now().Unix() {
			packages[k]["status"] = `<span class="label label-danger">已过期</span>`
		}
		if int64(starttime) > time.Now().Unix() {
			packages[k]["status"] = `<span class="label label-info">未生效</span>`
		}
		if _totalLeft == 0 {
			packages[k]["status"] = `<span class="label label-warning">已用完</span>`
		}
	}
	return
}

//分页获取Package
func (p *PackageModel) GetPackagesByLimit(limit int, number int) (packages []map[string]string, err error) {

	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From("package").OrderBy("package_id", "desc").Limit(limit, number))
	if err != nil {
		return
	}
	packages = rs.Rows()
	for k, v := range packages {
		_total, _ := strconv.ParseUint(v["bytes_total"], 10, 64)
		_totalLeft, _ := strconv.ParseUint(v["bytes_left"], 10, 64)
		_totalUse := _total - _totalLeft
		packages[k]["bytes_total_human"] = humanize.Bytes(_total)
		packages[k]["bytes_left_human"] = humanize.Bytes(_totalLeft)
		packages[k]["bytes_use_human"] = humanize.Bytes(_totalUse)
		packages[k]["status"] = `<span class="label label-success">有效</span>`
		endtime, _ := strconv.Atoi(v["end_time"])
		starttime, _ := strconv.Atoi(v["start_time"])
		if int64(endtime) <= time.Now().Unix() {
			packages[k]["status"] = `<span class="label label-danger">已过期</span>`
		}
		if int64(starttime) > time.Now().Unix() {
			packages[k]["status"] = `<span class="label label-info">未生效</span>`
		}
		if _totalLeft == 0 {
			packages[k]["status"] = `<span class="label label-warning">已用完</span>`
		}
	}
	return
}

func (p *PackageModel) CountPackages() (count int, err error) {

	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().Select("count(*) as total").From("package"))
	if err != nil {
		return
	}
	count = utils.NewConvert().StringToInt(rs.Value("total"))
	return
}

func (p *PackageModel) CountPackagesByID(column, id string) (count int, err error) {

	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().
		Select("count(*) as total").
		From("package").
		Where(map[string]interface{}{
			column: id,
		}))
	if err != nil {
		return
	}
	count = utils.NewConvert().StringToInt(rs.Value("total"))
	return
}
