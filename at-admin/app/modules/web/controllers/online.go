package controllers

import (
	"anytunnel/at-admin/app/modules/web/models"
	"anytunnel/at-admin/app/utils"
	"fmt"
	"strings"
)

type OnlineController struct {
	BaseController
}

func (this *OnlineController) List() {
	cs := strings.Trim(this.GetString("cs", ""), " ")
	column := strings.Trim(this.GetString("type", ""), " ")
	keyword := strings.Trim(this.GetString("keyword", ""), " ")
	page, err := this.GetInt("page", 1)
	if err != nil {
		page = 1
	}
	//每页的条数
	pageSize := 10
	limit := (page - 1) * pageSize

	onlineModel := models.Online{}
	var onlines = []map[string]string{}
	var onlineCount = 0

	if keyword == "" {
		onlineCount, err = onlineModel.CountOnlines(cs)
		onlines, err = onlineModel.GetOnlinesByLimit(cs, limit, pageSize)
	} else {
		onlineCount, err = onlineModel.CountOnlinesByID(cs, column, keyword)
		onlines, err = onlineModel.GetOnlinesByIDAndLimit(cs, column, keyword, limit, pageSize)
	}
	if err != nil {
		this.ViewError(err.Error())
	}
	this.Data["cs"] = cs
	this.Data["onlines"] = onlines
	this.Data["page"] = utils.NewMisc().Page(onlineCount, page, pageSize, fmt.Sprintf("/web/online/list?page={page}&type=%s&keyword=%s&cs=%s", column, keyword, cs))
	this.viewLayoutTitle(cs+"列表", "online/list", "form")
}
