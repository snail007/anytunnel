package models

type RolePrivilege struct {
	RolePrivilegeID uint `orm:"pk"` // role_privilege_id
	RoleID          uint // role_id
	PrivilegeID     uint // privilege_id
	IsDelete        int8 // is_delete
	CreateTime      uint // create_time
	UpdateTime      uint // update_time

}

//根据 role_id 获取权限
func (rolePrivilege *RolePrivilege) GetRolePrivilegesByRoleId(roleId int) (rolePrivileges []map[string]string, err error) {

	db := G.DB()
	res, err := db.Query(db.AR().From("role_privilege").Where(map[string]interface{}{
		"role_id": roleId,
	}))
	if(err != nil) {
		return
	}

	rolePrivileges = res.Rows()
	return
}

//角色授权
func (rolePrivilege *RolePrivilege) GrantRolePrivileges(roleId int, privilegeIds []string) (res bool, err error) {

	res = false
	db := G.DB()
	//先删除
	_, err = db.Exec(db.AR().Delete("role_privilege", map[string]interface{}{
		"role_id": roleId,
	}))
	if(err != nil) {
		return
	}

	rolePrivileges := []map[string]interface{}{}
	for _, privilegeId := range privilegeIds {
		rolePrivilege := map[string]interface{}{
			"role_id": roleId,
			"privilege_id": privilegeId,
		}
		rolePrivileges = append(rolePrivileges, rolePrivilege)
	}
	//批量插入
	_, err = db.Exec(db.AR().InsertBatch("role_privilege", rolePrivileges))
	if err != nil {
		return
	}
	res = true
	return
}