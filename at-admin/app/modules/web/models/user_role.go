package models

import (
	"time"

	"github.com/snail007/go-activerecord/mysql"
)

type UserRole struct {
	UserRoleID uint   `orm:"pk"` // user_role_id
	UserID     uint64 // user_id
	RoleID     uint   // role_id
	CreateTime uint   // create_time
	UpdateTime uint   // update_time
}

func (this *UserRole) GetUserRolesByUserId(userId string) (userRoles []map[string]string, err error) {

	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From("user_role").Where(map[string]interface{}{
		"user_id": userId,
	}))
	if err != nil {
		return
	}
	userRoles = rs.Rows()

	return
}

// 插入 user_id 和 role_id 对应关系
func (this *UserRole) Insert(userId string, roleIds []string) (res bool, err error) {

	res = false
	db := DB
	//先删除
	_, err = db.Exec(db.AR().Delete("user_role", map[string]interface{}{
		"user_id": userId,
	}))
	if err != nil {
		return
	}
	//添加
	userRoles := []map[string]interface{}{}
	for _, roleId := range roleIds {
		userRole := map[string]interface{}{
			"role_id":     roleId,
			"user_id":     userId,
			"create_time": time.Now().Unix(),
			"update_time": time.Now().Unix(),
		}
		userRoles = append(userRoles, userRole)
	}

	_, err = db.Exec(db.AR().InsertBatch("user_role", userRoles))
	if err != nil {
		return
	}
	res = true
	return
}
