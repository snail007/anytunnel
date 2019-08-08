package models

import (
	"strconv"

	humanize "github.com/dustin/go-humanize"
	"github.com/snail007/go-activerecord/mysql"
)

type Cluster struct {
}

func (p *Cluster) GetClusterByClusterId(clusterId string) (cluster map[string]string, err error) {
	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From("cluster").Where(map[string]interface{}{
		"cluster_id": clusterId,
		"is_delete":  0,
	}))
	if err != nil {
		return
	}
	cluster = rs.Row()
	return
}
func (p *Cluster) HasTunnel(clusterId string) (has bool, err error) {
	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From("tunnel").Where(map[string]interface{}{
		"cluster_id": clusterId,
	}).Limit(0, 1))
	if err != nil {
		return
	}
	if rs.Len() > 0 {
		has = true
	}
	return
}
func (p *Cluster) Delete(clusterId string) (err error) {
	db := DB
	_, err = db.Exec(db.AR().Update("cluster", map[string]interface{}{
		"is_delete": 1,
	}, map[string]interface{}{
		"cluster_id": clusterId,
	}))
	if err != nil {
		return
	}
	return
}
func (p *Cluster) Insert(cluster map[string]interface{}) (id int64, err error) {
	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Exec(db.AR().Insert("cluster", cluster))
	if err != nil {
		return
	}
	id = rs.LastInsertId
	return
}

func (p *Cluster) Update(clusterId string, cluster map[string]interface{}) (id int64, err error) {
	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Exec(db.AR().Update("cluster", cluster, map[string]interface{}{
		"cluster_id": clusterId,
	}))
	if err != nil {
		return
	}
	id = rs.LastInsertId
	return
}

//禁用
func (p *Cluster) Forbidden(clusterId string) (err error) {
	db := DB
	_, err = db.Exec(db.AR().Update("cluster", map[string]interface{}{
		"is_disable": 1,
	}, map[string]interface{}{
		"cluster_id": clusterId,
	}))
	if err != nil {
		return
	}
	return
}

//恢复
func (p *Cluster) Review(clusterId string) (err error) {
	db := DB
	_, err = db.Exec(db.AR().Update("cluster", map[string]interface{}{
		"is_disable": 0,
	}, map[string]interface{}{
		"cluster_id": clusterId,
	}))
	if err != nil {
		return
	}
	return
}

//获取所有的
func (this *Cluster) GetAllClusters() (clusters []map[string]string, err error) {
	db := DB
	res, err := db.Query(db.AR().From("cluster").Where(map[string]interface{}{
		"is_delete": 0,
	}).OrderBy("name", "ASC"))
	if err != nil {
		return
	}
	clusters = res.Rows()
	for k, v := range clusters {
		_bandwidth, _ := strconv.ParseUint(v["bandwidth"], 10, 64)
		clusters[k]["bandwidth_human"] = humanize.Bytes(_bandwidth)
	}
	return
}
