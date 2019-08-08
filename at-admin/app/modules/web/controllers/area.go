package controllers

import (
	"anytunnel/at-admin/app/modules/web/models"
	"anytunnel/at-admin/app/utils"
	"fmt"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

type AreaController struct {
	BaseController
}

func (this *AreaController) Delete() {
	areaModel := models.Area{}
	areaId := this.GetString("area_id")
	err := areaModel.Delete(areaId)
	if err != nil {
		this.JsonError(err)
	}
	this.JsonSuccess("")

}
func (this *AreaController) Add() {
	areaModel := models.Area{}
	if this.Ctx.Input.IsPost() {
		_, data := this.getAreaFromPost(false)
		_, err := areaModel.Insert(data)
		if err != nil {
			this.JsonError(err)
		}
		this.JsonSuccess("")
	} else {
		regionModel := models.Region{}
		regions, err := regionModel.GetSubRegions()
		if err != nil {
			this.ViewError(err.Error())
		}
		this.Data["regions"] = regions
		this.Data["action"] = "add"
		this.viewLayout("area/form", "form")
	}

}
func (this *AreaController) Edit() {
	areaModel := models.Area{}
	areaId := this.GetString("area_id")
	if this.Ctx.Input.IsPost() {
		_, data := this.getAreaFromPost(true)
		_, err := areaModel.Update(areaId, data)
		if err != nil {
			this.JsonError(err)
		}
		this.JsonSuccess("")
	} else {
		area, err := areaModel.GetAreaByAreaId(areaId)
		if err != nil {
			this.JsonError(err, "")
		}
		if len(area) == 0 {
			this.ViewError("Area不存在")
		}
		regionModel := models.Region{}
		regions, err := regionModel.GetSubRegions()
		if err != nil {
			this.ViewError(err.Error())
		}
		this.Data["regions"] = regions
		this.Data["area"] = area
		this.Data["action"] = "edit"
		this.viewLayout("area/form", "form")
	}

}
func (this *AreaController) List() {
	column := strings.Trim(this.GetString("type", ""), " ")
	keyword := strings.Trim(this.GetString("id", ""), " ")
	page, err := this.GetInt("page", 1)
	if err != nil {
		page = 1
	}
	//每页的条数
	pageSize := 10
	limit := (page - 1) * pageSize

	areaModel := models.Area{}
	var areas = []map[string]string{}
	var areaCount = 0

	if keyword == "" {
		areaCount, err = areaModel.CountAreas()
		areas, err = areaModel.GetAreasByLimit(limit, pageSize)
	} else {
		areaCount, err = areaModel.CountAreasByID(column, keyword)
		areas, err = areaModel.GetAreasByIDAndLimit(column, keyword, limit, pageSize)
	}
	if err != nil {
		this.ViewError(err.Error())
	}
	this.Data["areas"] = areas
	this.Data["page"] = utils.NewMisc().Page(areaCount, page, pageSize, fmt.Sprintf("/web/area/list?page={page}&type=%s&id=%s", column, keyword))
	this.viewLayoutTitle("用户列表", "area/list", "form")
}
func (this *AreaController) getAreaFromPost(isUpdate bool) (areaId string, area map[string]interface{}) {
	area = map[string]interface{}{
		"name":         this.GetString("name"),
		"cs_type":      this.GetString("cs_type"),
		"is_forbidden": this.GetString("is_forbidden"),
	}
	errs := validation.Errors{
		"区域名称": validation.Validate(area["name"],
			validation.Required.Error("不能为空")),
		"类型": validation.Validate(area["cs_type"],
			validation.Required.Error("不能为空"),
			validation.In("server", "client").Error("错误"),
		),
		"访问控制": validation.Validate(area["is_forbidden"],
			validation.Required.Error("不能为空"),
			validation.In("0", "1").Error("错误"),
		),
	}
	err := errs.Filter()
	if err != nil {
		this.JsonError(err)
	}
	if isUpdate {
		area["update_time"] = time.Now().Unix()
	} else {
		area["create_time"] = time.Now().Unix()
	}
	return
}
