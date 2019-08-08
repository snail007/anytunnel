package models

import (
	"strconv"
	"strings"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/snail007/go-activerecord/mysql"
)

type Role struct {
}

func (p *Role) GetRoleByRoleId(roleId string) (role map[string]string, err error) {
	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From("role").Where(map[string]interface{}{
		"role_id":   roleId,
		"is_delete": 0,
	}))
	if err != nil {
		return
	}
	role = rs.Row()
	return
}
func (p *Role) HasUser(roleId string) (has bool, err error) {
	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From("user_role").Where(map[string]interface{}{
		"role_id": roleId,
	}).Limit(0, 1))
	if err != nil {
		return
	}
	if rs.Len() > 0 {
		has = true
	}
	return
}
func (p *Role) Delete(roleId string) (err error) {
	db := DB
	_, err = db.Exec(db.AR().Update("role", map[string]interface{}{
		"is_delete": 1,
	}, map[string]interface{}{
		"role_id": roleId,
	}))
	if err != nil {
		return
	}
	return
}
func (p *Role) Insert(role map[string]interface{}) (id int64, err error) {
	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Exec(db.AR().Insert("role", role))
	if err != nil {
		return
	}
	id = rs.LastInsertId
	return
}

func (p *Role) Update(roleId string, role map[string]interface{}) (id int64, err error) {
	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Exec(db.AR().Update("role", role, map[string]interface{}{
		"role_id": roleId,
	}))
	if err != nil {
		return
	}
	id = rs.LastInsertId
	return
}

//获取所有的角色
func (this *Role) GetAllRoles() (roles []map[string]string, err error) {
	db := DB
	res, err := db.Query(db.AR().From("role").Where(map[string]interface{}{
		"is_delete": 0,
	}).OrderBy("role_id", "ASC"))
	if err != nil {
		return
	}
	roles = res.Rows()
	for k, v := range roles {
		c, _ := strconv.ParseUint(v["bandwidth"], 10, 64)
		v["bandwidth_human"] = humanize.Bytes(c)
		m := strings.Split(v["tunnel_mode"], ",")
		mmap := map[string]string{"0": "普通", "1": "高级", "2": "特殊"}
		_m := []string{}
		for _, val := range m {
			_m = append(_m, mmap[val])
		}
		v["tunnel_mode_human"] = strings.Join(_m, ",")
		roles[k] = v
	}
	return
}

// 根据用户名查找所有的角色
func (this *Role) GetRolesByUserId(userId string) (roles []map[string]string, err error) {

	userRoleModel := UserRole{}
	userRoles, err := userRoleModel.GetUserRolesByUserId(userId)
	if err != nil {
		return
	}

	for _, userRole := range userRoles {
		role, _ := this.GetRoleByRoleId(userRole["role_id"])
		roles = append(roles, role)
	}

	return
}
func (this *Role) SetRegions(roleId string, region_ids []string) (err error) {
	db := DB
	_, err = db.Exec(db.AR().Delete("role_region", map[string]interface{}{
		"role_id": roleId,
	}))
	if err != nil {
		return
	}
	data := []map[string]interface{}{}
	for _, v := range region_ids {
		data = append(data, map[string]interface{}{
			"role_id":     roleId,
			"region_id":   v,
			"create_time": time.Now().Unix(),
		})
	}
	if len(data) > 0 {
		_, err = db.Exec(db.AR().InsertBatch("role_region", data))
	}
	return
}
func (this *Role) GetRoleRegionIds(roleId string) (ids []string, err error) {
	db := DB
	rs, err := db.Query(db.AR().From("role_region").Where(map[string]interface{}{
		"role_id": roleId,
	}))
	if err != nil {
		return
	}
	ids = rs.Values("region_id")
	return
}
