package main

//client api
import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/snail007/go-activerecord/mysql"
)

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

//add a client
//method : POST
//params : name, token, user_id
func (client *Client) Add(responseWrite http.ResponseWriter, request *http.Request, params httprouter.Params) {

	name := request.PostFormValue("name")
	token := request.PostFormValue("token")
	userId := request.PostFormValue("user_id")

	if name == "" {
		jsonError(responseWrite, "名称不能为空!", nil)
		return
	}
	if userId == "" {
		jsonError(responseWrite, "user_id 错误!", nil)
		return
	}
	if token == "" {
		jsonError(responseWrite, "token 不能为空!", nil)
		return
	}

	clientValues := map[string]interface{}{
		"name":        name,
		"token":       token,
		"user_id":     userId,
		"create_time": time.Now().Unix(),
		"update_time": time.Now().Unix(),
	}

	db := G.DB()
	var rs *mysql.ResultSet
	rs, err := db.Exec(db.AR().Insert("client", clientValues))
	if err != nil {
		jsonError(responseWrite, err.Error(), nil)
		return
	}
	id := rs.LastInsertId

	jsonSuccess(responseWrite, "添加client成功", id)
}

//update client by client_id
//method : POST
//params : client_id, name, token, user_id
func (client *Client) Update(responseWrite http.ResponseWriter, request *http.Request, params httprouter.Params) {

	clientId := request.PostFormValue("client_id")
	name := request.PostFormValue("name")
	token := request.PostFormValue("token")

	if clientId == "" {
		jsonError(responseWrite, "client_id错误!", nil)
		return
	}
	clientValues := map[string]interface{}{}
	if name != "" {
		clientValues["name"] = name
	}
	if token != "" {
		clientValues["token"] = token
	}

	db := G.DB()
	var rs *mysql.ResultSet
	rs, err := db.Query(db.AR().From("client").Where(map[string]interface{}{
		"client_id": clientId,
		"is_delete": 0,
	}).Limit(0, 1))
	if err != nil {
		jsonError(responseWrite, err.Error(), nil)
		return
	}
	if rs.Len() == 0 {
		jsonError(responseWrite, "client不存在", nil)
		return
	}
	clientValues["update_time"] = time.Now().Unix()
	rs, err = db.Exec(db.AR().Update("client", clientValues, map[string]interface{}{
		"client_id": clientId,
	}))
	if err != nil {
		jsonError(responseWrite, err.Error(), nil)
		return
	}
	id := rs.LastInsertId

	jsonSuccess(responseWrite, "修改client成功", id)
}

//delete client by client_id
//method : GET
//params : client_id
func (client *Client) Delete(responseWrite http.ResponseWriter, request *http.Request, params httprouter.Params) {

	clientId := request.FormValue("client_id")

	if clientId == "" {
		jsonError(responseWrite, "client_id错误", nil)
		return
	}

	db := G.DB()
	var rs *mysql.ResultSet
	rs, err := db.Query(db.AR().From("client").Where(map[string]interface{}{
		"client_id": clientId,
		"is_delete": 0,
	}).Limit(0, 1))
	if err != nil {
		jsonError(responseWrite, err.Error(), nil)
		return
	}
	if rs.Len() == 0 {
		jsonError(responseWrite, "client 不存在", nil)
		return
	}

	clientValues := map[string]interface{}{
		"is_delete":   1,
		"update_time": time.Now().Unix(),
	}
	rs, err = db.Exec(db.AR().Update("client", clientValues, map[string]interface{}{
		"client_id": clientId,
	}))
	if err != nil {
		jsonError(responseWrite, err.Error(), nil)
		return
	}

	jsonSuccess(responseWrite, "删除client成功", nil)
}

//client list
//method : GET
//params : page(1), number(15), keyword("") user_id
func (client *Client) List(responseWrite http.ResponseWriter, request *http.Request, params httprouter.Params) {

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
	clientRows := []map[string]string{}
	if keyword != "" {
		sqlString := "SELECT * FROM at_client where is_delete=0"
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
		clientRows = rs.Rows()
	} else {
		sql := db.AR().From("client")
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
		clientRows = rs.Rows()
	}

	jsonSuccess(responseWrite, "ok", clientRows)
}

//get client by client_id
//method : GET
//params : client_id
func (client *Client) GetClientByClientId(responseWrite http.ResponseWriter, request *http.Request, params httprouter.Params) {

	clientId := request.FormValue("client_id")

	if clientId == "" {
		jsonError(responseWrite, "client_id错误!", nil)
		return
	}

	db := G.DB()
	var rs *mysql.ResultSet
	rs, err := db.Query(db.AR().From("client").Where(map[string]interface{}{
		"client_id": clientId,
		"is_delete": 0,
	}).Limit(0, 1))
	if err != nil {
		jsonError(responseWrite, err.Error(), nil)
		return
	}

	if rs.Len() == 0 {
		jsonError(responseWrite, "client 不存在", nil)
		return
	}

	jsonSuccess(responseWrite, "ok", rs.Row())
}
