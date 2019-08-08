package main

//server api
import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/snail007/go-activerecord/mysql"
)

type Server struct{}

func NewServer() *Server {
	return &Server{}
}

//add a server
//method : POST
//params : name, token, user_id
func (server *Server) Add(responseWrite http.ResponseWriter, request *http.Request, params httprouter.Params) {

	name := request.PostFormValue("name")
	token := request.PostFormValue("token")
	userId := request.PostFormValue("user_id")

	if name == "" {
		jsonError(responseWrite, "名称不能为空", nil)
		return
	}
	if userId == "" {
		jsonError(responseWrite, "user_id错误", nil)
		return
	}
	if token == "" {
		jsonError(responseWrite, "token 不能为空", nil)
		return
	}

	serverValues := map[string]interface{}{
		"name":        name,
		"token":       token,
		"user_id":     userId,
		"create_time": time.Now().Unix(),
		"update_time": time.Now().Unix(),
	}

	db := G.DB()
	var rs *mysql.ResultSet
	rs, err := db.Exec(db.AR().Insert("server", serverValues))
	if err != nil {
		jsonError(responseWrite, err.Error(), nil)
		return
	}
	id := rs.LastInsertId

	jsonSuccess(responseWrite, "添加server成功", id)
}

//update server by server_id
//method : POST
//params : server_id, name, token, user_id
func (server *Server) Update(responseWrite http.ResponseWriter, request *http.Request, params httprouter.Params) {

	serverId := request.PostFormValue("server_id")
	name := request.PostFormValue("name")
	token := request.PostFormValue("token")

	if serverId == "" {
		jsonError(responseWrite, "server_id错误", nil)
		return
	}
	serverValues := map[string]interface{}{}
	if name != "" {
		serverValues["name"] = name
	}
	if token != "" {
		serverValues["token"] = token
	}

	db := G.DB()
	var rs *mysql.ResultSet
	rs, err := db.Query(db.AR().From("server").Where(map[string]interface{}{
		"server_id": serverId,
		"is_delete": 0,
	}).Limit(0, 1))
	if err != nil {
		jsonError(responseWrite, err.Error(), nil)
		return
	}
	if rs.Len() == 0 {
		jsonError(responseWrite, "server不存在", nil)
		return
	}

	serverValues["update_time"] = time.Now().Unix()
	rs, err = db.Exec(db.AR().Update("server", serverValues, map[string]interface{}{
		"server_id": serverId,
	}))
	if err != nil {
		jsonError(responseWrite, err.Error(), nil)
		return
	}
	id := rs.LastInsertId

	jsonSuccess(responseWrite, "修改server成功", id)
}

//delete server by server_id
//method : GET
//params : server_id
func (server *Server) Delete(responseWrite http.ResponseWriter, request *http.Request, params httprouter.Params) {

	serverId := request.FormValue("server_id")

	if serverId == "" {
		jsonError(responseWrite, "server_id 错误", nil)
		return
	}

	db := G.DB()
	var rs *mysql.ResultSet
	rs, err := db.Query(db.AR().From("server").Where(map[string]interface{}{
		"server_id": serverId,
		"is_delete": 0,
	}).Limit(0, 1))
	if err != nil {
		jsonError(responseWrite, err.Error(), nil)
		return
	}
	if rs.Len() == 0 {
		jsonError(responseWrite, "server 不存在", nil)
		return
	}

	serverValues := map[string]interface{}{
		"is_delete":   1,
		"update_time": time.Now().Unix(),
	}
	rs, err = db.Exec(db.AR().Update("server", serverValues, map[string]interface{}{
		"server_id": serverId,
	}))
	if err != nil {
		jsonError(responseWrite, err.Error(), nil)
		return
	}

	jsonSuccess(responseWrite, "删除 server 成功", nil)
}

//server list
//method : GET
//params : page(1), number(15), keyword("") user_id
func (server *Server) List(responseWrite http.ResponseWriter, request *http.Request, params httprouter.Params) {

	keyword := strings.Trim(request.FormValue("keyword"), "")
	page := request.FormValue("page")
	pageSize := request.FormValue("number")
	userId := request.FormValue("user_id")

	pageNumber := 1
	number := 15
	if page != "" {
		pageNumber, _ = strconv.Atoi(page)
	}
	if pageSize != "" {
		number, _ = strconv.Atoi(pageSize)
	}

	offset := (pageNumber - 1) * number

	db := G.DB()
	var rs *mysql.ResultSet
	var err error
	serverRows := []map[string]string{}
	if keyword != "" {
		sqlString := "SELECT * FROM at_server where is_delete=0"
		if userId != "" {
			sqlString += " AND user_id=" + userId
		}
		sqlString += " AND (" +
			"token LIKE '%" + keyword + "%' OR " +
			"name LIKE '%" + keyword + "%'" +
			")"
		sql := db.AR().Raw(sqlString)
		rs, err = db.Query(sql)
		if err != nil {
			jsonError(responseWrite, err.Error(), nil)
			return
		}
		serverRows = rs.Rows()
	} else {
		sql := db.AR().From("server")
		if userId != "" {
			sql = sql.Where(map[string]interface{}{"is_delete": 0, "user_id": userId})
		} else {
			sql = sql.Where(map[string]interface{}{"is_delete": 0})
		}
		sql = sql.Limit(offset, number)
		rs, err = db.Query(sql)
		if err != nil {
			jsonError(responseWrite, err.Error(), nil)
			return
		}
		serverRows = rs.Rows()
	}

	jsonSuccess(responseWrite, "ok", serverRows)
}

//get server by server_id
//method : GET
//params : server_id
func (server *Server) GetServerByServerId(responseWrite http.ResponseWriter, request *http.Request, params httprouter.Params) {

	serverId := request.FormValue("server_id")

	if serverId == "" {
		jsonError(responseWrite, "server_id 错误", nil)
		return
	}

	db := G.DB()
	var rs *mysql.ResultSet
	rs, err := db.Query(db.AR().From("server").Where(map[string]interface{}{
		"server_id": serverId,
		"is_delete": 0,
	}).Limit(0, 1))
	if err != nil {
		jsonError(responseWrite, err.Error(), nil)
		return
	}

	if rs.Len() == 0 {
		jsonError(responseWrite, "server 不存在", nil)
		return
	}

	jsonSuccess(responseWrite, "ok", rs.Row())
}
