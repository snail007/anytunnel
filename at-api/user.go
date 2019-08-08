package main

import (
	"anytunnel/at-web/app/utils"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/snail007/go-activerecord/mysql"
)

const DEFAULT_USER_ROLE = 1

const MODE_BASE  = "0"
const MODE_SENIOR  = "1"
const MODE_SPECIAL  = "2"

type User struct{}

func NewUser() *User {
	return &User{}
}

//get user by user_id
//method : GET
//params : user_id
func (this *User) GetUserByUserId(responseWrite http.ResponseWriter, request *http.Request, params httprouter.Params) {

	userId := request.FormValue("user_id")
	if userId == "" {
		jsonError(responseWrite, "user_id is not empty!", nil)
		return
	}
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err := db.Query(db.AR().From("user").Where(map[string]interface{}{
		"user_id": userId,
	}).Limit(0, 1))
	if err != nil {
		jsonError(responseWrite, err.Error(), nil)
		return
	}
	jsonSuccess(responseWrite, "ok", rs.Row())
}

//get user by name
//method : GET
//params : username
func (this *User) GetUserByName(responseWrite http.ResponseWriter, request *http.Request, params httprouter.Params) {

	username := request.FormValue("username")
	if username == "" {
		jsonError(responseWrite, "username is not empty!", nil)
		return
	}
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err := db.Query(db.AR().From("user").Where(map[string]interface{}{
		"username": username,
	}).Limit(0, 1))
	if err != nil {
		jsonError(responseWrite, err.Error(), nil)
		return
	}
	jsonSuccess(responseWrite, "ok", rs.Row())

}

//get user by email
//method : GET
//params : email
func (this *User) GetUserByEmail(responseWrite http.ResponseWriter, request *http.Request, params httprouter.Params) {

	email := request.FormValue("email")
	if email == "" {
		jsonError(responseWrite, "email is not empty!", nil)
		return
	}
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err := db.Query(db.AR().From("user").Where(map[string]interface{}{
		"email": email,
	}).Limit(0, 1))
	if err != nil {
		jsonError(responseWrite, err.Error(), nil)
		return
	}
	jsonSuccess(responseWrite, "ok", rs.Row())
}

//create user
//method : POST
//params : username, password, email
func (this *User) Create(responseWrite http.ResponseWriter, request *http.Request, params httprouter.Params) {

	username := request.PostFormValue("username")
	password := request.PostFormValue("password")
	email := request.PostFormValue("email")
	nickname := request.PostFormValue("nickname")

	if username == "" {
		jsonError(responseWrite, "用户名不能为空!", nil)
		return
	}
	if password == "" {
		jsonError(responseWrite, "密码不能为空!", nil)
	}
	if email == "" {
		jsonError(responseWrite, "邮箱不能为空!", nil)
	}

	//验证用户名是否存在
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err := db.Query(db.AR().From("user").Where(map[string]interface{}{
		"username": username,
	}).Limit(0, 1))
	if err != nil {
		jsonError(responseWrite, err.Error(), nil)
		return
	}
	if len(rs.Row()) > 0 {
		jsonError(responseWrite, "抱歉, 用户名已经存在", nil)
		return
	}

	//验证邮箱是否存在
	rs, err = db.Query(db.AR().From("user").Where(map[string]interface{}{
		"email": email,
	}).Limit(0, 1))
	if err != nil {
		jsonError(responseWrite, err.Error(), nil)
		return
	}
	if len(rs.Row()) > 0 {
		jsonError(responseWrite, "抱歉，邮箱已经被注册", nil)
		return
	}
	//添加用户
	userValues := map[string]interface{}{
		"username":     username,
		"password":     utils.NewEncrypt().Md5Encode(password),
		"email":        email,
		"nickname":     nickname,
		"is_active":    0,
		"is_forbidden": 0,
		"create_time":  time.Now().Unix(),
		"update_time":  time.Now().Unix(),
	}
	rs, err = db.Exec(db.AR().Insert("user", userValues))
	if err != nil {
		jsonError(responseWrite, err.Error(), nil)
		return
	}
	userId := rs.LastInsertId
	userValues["user_id"] = userId
	//创建用户和角色对应关系
	userRole := map[string]interface{}{
		"role_id":     DEFAULT_USER_ROLE,
		"user_id":     userId,
		"create_time": time.Now().Unix(),
		"update_time": time.Now().Unix(),
	}
	rs, err = db.Exec(db.AR().Insert("user_role", userRole))
	if err != nil {
		jsonError(responseWrite, err.Error(), nil)
		return
	}

	jsonSuccess(responseWrite, "创建成功", userValues)
}

//update user
//method : POST
//params : nickname, password, email
func (this *User) Update(responseWrite http.ResponseWriter, request *http.Request, params httprouter.Params) {

	userId := request.PostFormValue("user_id")
	password := request.PostFormValue("password")
	nickname := request.PostFormValue("nickname")

	if userId == "" {
		jsonError(responseWrite, "user_id error!", nil)
		return
	}
	userValue := map[string]interface{}{}
	if password != "" {
		userValue["password"] = utils.NewEncrypt().Md5Encode(password)
	}
	if nickname != "" {
		userValue["nickname"] = nickname
	}
	userValue["update_time"] = time.Now().Unix()
	db := G.DB()
	rs, err := db.Exec(db.AR().Update("user", userValue, map[string]interface{}{
		"user_id": userId,
	}))
	if err != nil {
		jsonError(responseWrite, err.Error(), nil)
		return
	}
	id := rs.LastInsertId
	jsonSuccess(responseWrite, "修改成功", id)
}

//get roles By user_id
//method : GET
//params : user_id
//return : roles
func (this *User) GetRolesByUserId(responseWrite http.ResponseWriter, request *http.Request, params httprouter.Params) {

	userId := request.FormValue("user_id")
	if userId == "" {
		jsonError(responseWrite, "user_id is not empty!", nil)
		return
	}
	db := G.DB()
	var rs *mysql.ResultSet
	//查找用户角色对应关系
	rs, err := db.Query(db.AR().From("user_role").Where(map[string]interface{}{
		"user_id": userId,
	}))
	if err != nil {
		jsonError(responseWrite, err.Error(), nil)
		return
	}
	userRoles := rs.Rows()
	roleIds := []string{}
	for _, userRole := range userRoles {
		roleIds = append(roleIds, userRole["role_id"])
	}

	//根据roleIds查找角色
	rs, err = db.Query(db.AR().From("role").Where(map[string]interface{}{
		"role_id": roleIds,
	}))
	if err != nil {
		jsonError(responseWrite, err.Error(), nil)
		return
	}
	roles := rs.Rows()

	jsonSuccess(responseWrite, "ok", roles)
}

//user cluster list
//method : GET
//params : user_id
//return : cluster region
func (this *User) Cluster(responseWrite http.ResponseWriter, request *http.Request, params httprouter.Params) {

	userId := request.FormValue("user_id")
	mode := request.FormValue("mode")
	if userId == "" {
		jsonError(responseWrite, "user_id 错误!", nil)
		return
	}
	if(mode == "") {
		mode = "0"
	}

	//查找用户角色对应关系
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err := db.Query(db.AR().From("user_role").Where(map[string]interface{}{
		"user_id": userId,
	}))
	if err != nil {
		jsonError(responseWrite, err.Error(), nil)
		return
	}
	if rs.Len() == 0 {
		jsonError(responseWrite, "user_id no role", nil)
		return
	}
	userRoles := rs.Rows()
	roleIds := []string{}
	for _, userRole := range userRoles {
		roleIds = append(roleIds, userRole["role_id"])
	}

	//查找角色地区对应关系
	rs, err = db.Query(db.AR().From("role_region").Where(map[string]interface{}{
		"role_id": roleIds,
	}))
	if err != nil {
		jsonError(responseWrite, err.Error(), nil)
		return
	}
	if rs.Len() == 0 {
		jsonError(responseWrite, "该角色没有可选地区", nil)
		return
	}
	roleRegions := rs.Rows()
	regionIds := []string{}
	for _, roleRegion := range roleRegions {
		regionIds = append(regionIds, roleRegion["region_id"])
	}

	//查找所有的地区
	rs, err = db.Query(db.AR().From("region").Where(map[string]interface{}{
		"region_id": regionIds,
		"is_delete": 0,
	}))
	if err != nil {
		jsonError(responseWrite, err.Error(), nil)
		return
	}
	if rs.Len() == 0 {
		jsonError(responseWrite, "no region", nil)
		return
	}
	regions := rs.Rows()

	onlineClusterIds := []string{}
	//基础模式，需要查找系统server在线的cluster
	//特殊模式，需要查找系统的client在线的cluster
	if(mode == MODE_BASE || mode == MODE_SPECIAL) {
		cs_type := "client"
		if(mode == MODE_BASE) {
			cs_type = "server"
		}
		rs, err = db.Query(db.AR().From("online").Where(map[string]interface{}{
			"user_id": 0,
			"cs_type": cs_type,
		}))
		if err != nil {
			jsonError(responseWrite, err.Error(), nil)
			return
		}
		if(len(rs.Rows()) > 0) {
			for _, online := range rs.Rows() {
				onlineClusterIds = append(onlineClusterIds, online["cluster_id"])
			}
		}
	}

	//组合数据
	regionClusters := map[string][]map[string]interface{}{}
	for _, region := range regions {
		//查找父级的地区
		rs, err = db.Query(db.AR().From("region").Where(map[string]interface{}{
			"region_id": region["parent_id"],
			"is_delete": 0,
		}))
		if err != nil {
			jsonError(responseWrite, err.Error(), nil)
			return
		}
		if(rs.Len() == 0) {
			continue
		}
		firstRegionName := rs.Row()["name"]

		regionCluster := map[string]interface{}{
			"region_id" :   region["region_id"],
			"name" :        region["name"],
			"clusters":     []map[string]string{},
		}

		whereMap := map[string]interface{}{
			"region_id": region["region_id"],
			"is_disable": 0,
			"is_delete": 0,
		}
		//基础模式和特殊模式需要同时满足在线的cluster条件
		if(mode == MODE_SPECIAL || mode == MODE_BASE) {
			whereMap["cluster_id"] = onlineClusterIds
		}
		//查找满足地区的cluster
		rs, err = db.Query(db.AR().From("cluster").Where(whereMap))
		if(err != nil) {
			jsonError(responseWrite, err.Error(), nil)
			return
		}
		if(rs.Len() == 0) {
			continue
		}
		clusters := rs.Rows()
		regionCluster["clusters"] = clusters

		regionClusters[firstRegionName] = append(regionClusters[firstRegionName], regionCluster)
	}

	jsonSuccess(responseWrite, "", regionClusters)

}