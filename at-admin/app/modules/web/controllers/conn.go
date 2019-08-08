package controllers

import (
	"anytunnel/at-admin/app/modules/web/models"
	"anytunnel/at-admin/app/utils"
	"fmt"
	"strings"
)

type ConnController struct {
	BaseController
}

func (this *ConnController) List() {
	orderBy := strings.Trim(this.GetString("orderby"), " ")
	if orderBy == "" {
		orderBy = "count"
	}
	col := strings.Trim(this.GetString("col"), " ")
	keyword := strings.Trim(this.GetString("keyword"), " ")
	page, err := this.GetInt("page", 1)
	if err != nil {
		page = 1
	}
	//每页的条数
	pageSize := 10
	limit := (page - 1) * pageSize

	connModel := models.Conn{}
	var conns = []map[string]string{}
	var connCount = 0
	if keyword == "" {
		connCount, err = connModel.CountConns()
		conns, err = connModel.GetConnsByLimit(limit, pageSize, orderBy)
	} else {
		connCount, err = connModel.CountConnsByUserID(col, keyword)
		conns, err = connModel.GetConnsByUserIDAndLimit(col, keyword, orderBy, limit, pageSize)
	}
	if err != nil {
		this.ViewError(err.Error())
	}

	this.Data["conns"] = conns
	this.Data["col"] = col
	this.Data["keyword"] = keyword
	this.Data["orderby"] = orderBy
	this.Data["p"] = page
	this.Data["page"] = utils.NewMisc().Page(connCount, page, pageSize, fmt.Sprintf("/web/conn/list?page={page}&col=%s&keyword=%s&orderby=%s", col, keyword, orderBy))
	this.viewLayoutTitle("隧道连接数列表", "conn/list", "form")
}
