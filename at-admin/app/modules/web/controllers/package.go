package controllers

import (
	"anytunnel/at-admin/app/modules/web/models"
	"anytunnel/at-admin/app/utils"
	"fmt"
	"regexp"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

type PackageController struct {
	BaseController
}

func (this *PackageController) List() {
	column := "user_id"
	keyword := strings.Trim(this.GetString("keyword", ""), " ")
	page, err := this.GetInt("page", 1)
	if err != nil {
		page = 1
	}
	//每页的条数
	pageSize := 10
	limit := (page - 1) * pageSize

	packageModel := models.PackageModel{}
	var packages = []map[string]string{}
	var packageCount = 0

	if keyword == "" {
		packageCount, err = packageModel.CountPackages()
		packages, err = packageModel.GetPackagesByLimit(limit, pageSize)
	} else {
		packageCount, err = packageModel.CountPackagesByID(column, keyword)
		packages, err = packageModel.GetPackagesByIDAndLimit(column, keyword, limit, pageSize)
	}
	if err != nil {
		this.ViewError(err.Error())
	}
	this.Data["packages"] = packages
	this.Data["page"] = utils.NewMisc().Page(packageCount, page, pageSize, fmt.Sprintf("/web/package/list?page={page}&type=%s&keyword=%s", column, keyword))
	this.viewLayoutTitle("流量列表", "package/list", "form")
}
func (this *PackageController) Add() {
	packageModel := models.PackageModel{}
	if this.Ctx.Input.IsPost() {
		_, data := this.getPackageFromPost(false)
		_, err := packageModel.Insert(data)
		if err != nil {
			this.JsonError(err)
		}
		this.JsonSuccess("")
	} else {
		this.Data["action"] = "add"
		this.viewLayout("package/form", "form")
	}
}
func (this *PackageController) Edit() {
	packageModel := models.PackageModel{}
	packageId := this.GetString("package_id")
	if this.Ctx.Input.IsPost() {
		_, data := this.getPackageFromPost(true)
		_, err := packageModel.Update(packageId, data)
		if err != nil {
			this.JsonError(err)
		}
		this.JsonSuccess("")
	} else {
		_package, err := packageModel.GetPackageByPackageId(packageId)
		if err != nil {
			this.JsonError(err, "")
		}
		if len(_package) == 0 {
			this.ViewError("Package不存在")
		}
		this.Data["package"] = _package
		this.Data["action"] = "edit"
		this.viewLayout("package/form", "form")
	}

}
func (this *PackageController) getPackageFromPost(isUpdate bool) (packageId string, _package map[string]interface{}) {
	_package = map[string]interface{}{
		"comment":    this.GetString("comment"),
		"bytes_left": this.GetString("bytes_left"),
		"user_id":    this.GetString("user_id"),
		"start_time": this.GetString("start_time"),
		"end_time":   this.GetString("end_time"),
	}
	errs := validation.Errors{
		"来源": validation.Validate(_package["comment"],
			validation.Required.Error("不能为空"),
			validation.Match(regexp.MustCompile("^.{1,10}$")).Error("长度必须是1-10字符")),
		"字节数": validation.Validate(_package["bytes_left"],
			validation.Required.Error("不能为空"),
			validation.Match(regexp.MustCompile("^[0-9]+$")).Error("必须是数字")),
		"用户ID": validation.Validate(_package["user_id"],
			validation.Required.Error("不能为空"),
			validation.Match(regexp.MustCompile("^[1-9][0-9]*$")).Error("必须是数字")),
		"生效时间": validation.Validate(_package["start_time"],
			validation.Required.Error("不能为空"),
			validation.Match(regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)).Error("格式错误")),
		"过期时间": validation.Validate(_package["end_time"],
			validation.Required.Error("不能为空"),
			validation.Match(regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)).Error("格式错误")),
	}

	err := errs.Filter()

	if err != nil {
		this.JsonError(err)
	}
	userModel := models.User{}
	user, err := userModel.GetUserByUserId(_package["user_id"].(string))
	if err != nil {
		this.JsonError(err)
	}
	if len(user) == 0 {
		this.JsonError("用户不存在")
	}
	_start, _ := time.ParseInLocation("2006-01-02 15:04:05", _package["start_time"].(string)+" 00:00:00", time.Local)
	_end, _ := time.ParseInLocation("2006-01-02 15:04:05", _package["end_time"].(string)+" 23:59:59", time.Local)
	_package["start_time"] = _start.Unix()
	_package["end_time"] = _end.Unix()
	if isUpdate {
		_package["update_time"] = time.Now().Unix()
	} else {
		_package["bytes_total"] = _package["bytes_left"]
		_package["create_time"] = time.Now().Unix()
	}
	return
}
