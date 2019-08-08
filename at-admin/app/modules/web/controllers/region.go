package controllers

import (
	"anytunnel/at-admin/app/modules/web/models"
	"regexp"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

type RegionController struct {
	BaseController
}

func (this *RegionController) Delete() {
	regionModel := models.Region{}
	regionId := this.GetString("region_id")
	hasRole, err := regionModel.HasRole(regionId)
	if err != nil {
		this.JsonError(err)
	}
	if hasRole {
		this.JsonError("角色引用非空,不能删除")
	}
	HasCluster, err := regionModel.HasCluster(regionId)
	if err != nil {
		this.JsonError(err)
	}
	if HasCluster {
		this.JsonError("Cluster引用非空,不能删除")
	}
	HasSubRegion, err := regionModel.HasSubRegion(regionId)
	if err != nil {
		this.JsonError(err)
	}
	if HasSubRegion {
		this.JsonError("子区域非空,不能删除")
	}
	err = regionModel.Delete(regionId)
	if err != nil {
		this.JsonError(err)
	}
	this.JsonSuccess("")

}
func (this *RegionController) Add() {
	regionModel := models.Region{}
	if this.Ctx.Input.IsPost() {
		_, data := this.getRegionFromPost(false)
		_, err := regionModel.Insert(data)
		if err != nil {
			this.JsonError(err)
		}
		this.JsonSuccess("")
	} else {
		regions, err := regionModel.GetTopRegions()
		if err != nil {
			this.ViewError(err.Error())
		}
		this.Data["regions"] = regions
		this.Data["action"] = "add"
		this.viewLayout("region/form", "form")
	}

}
func (this *RegionController) Edit() {
	regionModel := models.Region{}
	regionId := this.GetString("region_id")
	if this.Ctx.Input.IsPost() {
		_, data := this.getRegionFromPost(true)
		_, err := regionModel.Update(regionId, data)
		if err != nil {
			this.JsonError(err)
		}
		this.JsonSuccess("")
	} else {
		region, err := regionModel.GetRegionByRegionId(regionId)
		if err != nil {
			this.JsonError(err, "")
		}
		if len(region) == 0 {
			this.ViewError("区域不存在")
		}
		regions, err := regionModel.GetTopRegions()
		if err != nil {
			this.ViewError(err.Error())
		}
		this.Data["regions"] = regions
		this.Data["region"] = region
		this.Data["action"] = "edit"
		this.viewLayout("region/form", "form")
	}

}

func (this *RegionController) List() {
	regionModel := models.Region{}
	regions, err := regionModel.GetAllRegions()
	if err != nil {
		this.ViewError(err.Error())
	}
	this.Data["regions"] = regions
	this.view("region/list")
}

func (this *RegionController) getRegionFromPost(isUpdate bool) (regionId string, region map[string]interface{}) {
	region = map[string]interface{}{
		"name":      this.GetString("name"),
		"parent_id": this.GetString("parent_id"),
		"is_delete": 0,
	}
	err := validation.Validate(region["name"],
		validation.Required.Error("名称不能为空"),
		validation.Match(regexp.MustCompile("^.{1,15}$")).Error("名称长度必须是1-15字符"))
	if err != nil {
		this.JsonError(err.Error())
	}
	if isUpdate {
		region["update_time"] = time.Now().Unix()
	} else {
		region["create_time"] = time.Now().Unix()
	}
	return
}
