package models

import (
	"anytunnel/at-admin/app/utils"

	"github.com/snail007/go-activerecord/mysql"
)

type Tunnel struct {
}

func (p *Tunnel) GetTunnelByTunnelId(tunnelId string) (tunnel map[string]string, err error) {
	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From("tunnel").Where(map[string]interface{}{
		"tunnel_id": tunnelId,
	}))
	if err != nil {
		return
	}
	tunnel = rs.Row()
	return
}

func (p *Tunnel) Delete(tunnelId string) (err error) {
	db := DB
	_, err = db.Exec(db.AR().Update("tunnel", map[string]interface{}{
		"is_delete": 1,
	}, map[string]interface{}{
		"tunnel_id": tunnelId,
	}))
	if err != nil {
		return
	}
	return
}
func (p *Tunnel) Insert(tunnel map[string]interface{}) (id int64, err error) {
	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Exec(db.AR().Insert("tunnel", tunnel))
	if err != nil {
		return
	}
	id = rs.LastInsertId
	return
}

func (p *Tunnel) Update(tunnelId string, tunnel map[string]interface{}) (id int64, err error) {
	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Exec(db.AR().Update("tunnel", tunnel, map[string]interface{}{
		"tunnel_id": tunnelId,
	}))
	if err != nil {
		return
	}
	id = rs.LastInsertId
	return
}

//根据搜索分页获取Tunnel
func (tunnel *Tunnel) GetTunnelsByIDAndLimit(column, id string, limit int, number int) (tunnels []map[string]string, err error) {

	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From("tunnel").Where(map[string]interface{}{
		column:      id,
		"is_delete": 0,
	}).OrderBy("tunnel_id", "desc").Limit(limit, number))
	if err != nil {
		return
	}
	tunnels = rs.Rows()
	return
}

//分页获取Tunnel
func (tunnel *Tunnel) GetTunnelsByLimit(limit int, number int) (tunnels []map[string]string, err error) {

	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From("tunnel").Where(map[string]interface{}{
		"is_delete": 0,
	}).
		OrderBy("tunnel_id", "desc").
		Limit(limit, number))
	if err != nil {
		return
	}
	tunnels = rs.Rows()
	return
}

func (tunnel *Tunnel) CountTunnels() (count int, err error) {

	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().Select("count(*) as total").Where(map[string]interface{}{
		"is_delete": 0,
	}).From("tunnel"))
	if err != nil {
		return
	}
	count = utils.NewConvert().StringToInt(rs.Value("total"))
	return
}

func (tunnel *Tunnel) CountTunnelsByID(column, id string) (count int, err error) {

	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().
		Select("count(*) as total").
		From("tunnel").
		Where(map[string]interface{}{
			column:      id,
			"is_delete": 0,
		}))
	if err != nil {
		return
	}
	count = utils.NewConvert().StringToInt(rs.Value("total"))
	return
}
