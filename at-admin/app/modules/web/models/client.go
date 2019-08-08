package models

import (
	"anytunnel/at-admin/app/utils"
	common "anytunnel/at-common"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/snail007/go-activerecord/mysql"
)

type Client struct {
}

func (p *Client) GetClientByClientId(clientId string) (client map[string]string, err error) {
	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From("client").Where(map[string]interface{}{
		"client_id": clientId,
	}))
	if err != nil {
		return
	}
	client = rs.Row()
	return
}
func (p *Client) HasTunnelRef(clientId string) (has bool, err error) {
	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From("tunnel").Where(map[string]interface{}{
		"client_id": clientId,
	}).Limit(0, 1))
	if err != nil {
		return
	}
	if rs.Len() > 0 {
		has = true
	}
	return
}
func (p *Client) Reset(clientId string) (err error) {
	db := DB
	_, err = db.Exec(db.AR().Update("client", map[string]interface{}{
		"token": utils.NewMisc().RandString(32),
	}, map[string]interface{}{
		"client_id": clientId,
	}))
	if err != nil {
		return
	}
	return
}
func (p *Client) Delete(clientId string) (err error) {
	db := DB
	_, err = db.Exec(db.AR().Update("client", map[string]interface{}{
		"is_delete": 1,
	}, map[string]interface{}{
		"client_id": clientId,
	}))
	if err != nil {
		return
	}
	return
}
func (p *Client) Insert(client map[string]interface{}) (id int64, err error) {
	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Exec(db.AR().Insert("client", client))
	if err != nil {
		return
	}
	id = rs.LastInsertId
	return
}

func (p *Client) Update(clientId string, client map[string]interface{}) (id int64, err error) {
	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Exec(db.AR().Update("client", client, map[string]interface{}{
		"client_id": clientId,
	}))
	if err != nil {
		return
	}
	id = rs.LastInsertId
	return
}

//获取所有的Client
func (this *Client) GetAllClients() (clients []map[string]string, err error) {
	db := DB
	res, err := db.Query(db.AR().From("client").Where(map[string]interface{}{
		"is_delete": 0,
	}).OrderBy("client_id", "ASC"))
	if err != nil {
		return
	}
	clients = res.Rows()
	return
}

//根据user_id分页获取Client
func (client *Client) GetClientsByUserIDAndLimit(userID string, limit int, number int) (clients []map[string]string, err error) {

	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From("client").Where(map[string]interface{}{
		"user_id":   userID,
		"is_delete": 0,
	}).OrderBy("client_id", "desc").Limit(limit, number))
	if err != nil {
		return
	}
	clients = rs.Rows()

	return
}

//分页获取Client
func (client *Client) GetClientsByLimit(limit int, number int) (clients []map[string]string, err error) {

	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From("client").Where(map[string]interface{}{
		"is_delete":  0,
		"user_id <>": 0,
	}).OrderBy("client_id", "desc").Limit(limit, number))
	if err != nil {
		return
	}
	clients = rs.Rows()

	return
}

func (client *Client) CountClients() (count int, err error) {

	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().Select("count(*) as total").Where(map[string]interface{}{
		"is_delete":  0,
		"user_id <>": 0,
	}).From("client"))
	if err != nil {
		return
	}
	count = utils.NewConvert().StringToInt(rs.Value("total"))
	return
}

func (client *Client) CountClientsByUserID(userID string) (count int, err error) {

	db := DB
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().
		Select("count(*) as total").
		From("client").
		Where(map[string]interface{}{
			"user_id":   userID,
			"is_delete": 0,
		}))
	if err != nil {
		return
	}
	count = utils.NewConvert().StringToInt(rs.Value("total"))
	return
}
func (p *Client) Offline(serverId string) (err error) {
	db := DB
	rs, err := db.Query(db.AR().From("client").Where(map[string]interface{}{
		"client_id": serverId,
		"is_delete": 0,
	}))
	if err != nil {
		return
	}
	if rs.Len() == 0 {
		return
	}
	token := rs.Value("token")
	rs, err = db.Query(db.AR().From("online").Where(map[string]interface{}{
		"token": token,
		"type":  "client",
	}))
	if err != nil {
		return
	}
	if rs.Len() == 0 {
		return
	}
	clusterIP := rs.Value("ip")
	url := fmt.Sprintf("https://%s:%s/client/offline/%s", clusterIP, beego.AppConfig.String("cluster.api.port"), token)
	body, code, err := common.HttpGet(url)
	if err != nil {
		return
	}
	if code != 200 {
		err = fmt.Errorf("access client offline url %s fail,code:%d ,body:%s", url, code, body)
	}
	return
}
