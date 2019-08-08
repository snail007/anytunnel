package controllers

import (
	"anytunnel/at-admin/app/modules/web/models"
	"anytunnel/at-admin/app/utils"
	"fmt"
	"strings"
)

type TunnelController struct {
	BaseController
}

func (this *TunnelController) List() {
	column := strings.Trim(this.GetString("type", ""), " ")
	keyword := strings.Trim(this.GetString("id", ""), " ")
	page, err := this.GetInt("page", 1)
	if err != nil {
		page = 1
	}
	//每页的条数
	pageSize := 10
	limit := (page - 1) * pageSize

	tunnelModel := models.Tunnel{}
	var tunnels = []map[string]string{}
	var tunnelCount = 0

	if keyword == "" {
		tunnelCount, err = tunnelModel.CountTunnels()
		tunnels, err = tunnelModel.GetTunnelsByLimit(limit, pageSize)
	} else {
		tunnelCount, err = tunnelModel.CountTunnelsByID(column, keyword)
		tunnels, err = tunnelModel.GetTunnelsByIDAndLimit(column, keyword, limit, pageSize)
	}
	if err != nil {
		this.ViewError(err.Error())
	}
	this.Data["tunnels"] = tunnels
	this.Data["page"] = utils.NewMisc().Page(tunnelCount, page, pageSize, fmt.Sprintf("/web/tunnel/list?page={page}&type=%s&id=%s", column, keyword))
	this.viewLayoutTitle("用户列表", "tunnel/list", "form")
}
