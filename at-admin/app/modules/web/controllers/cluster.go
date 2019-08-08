package controllers

import (
	"anytunnel/at-admin/app/modules/web/models"
	common "anytunnel/at-common"
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/astaxie/beego"
	validation "github.com/go-ozzo/ozzo-validation"
	is "github.com/go-ozzo/ozzo-validation/is"
)

type ClusterController struct {
	BaseController
}

func (this *ClusterController) Statistic() {
	clusterID := this.GetString("cluster_id")
	if this.Ctx.Input.IsGet() {
		this.Data["clusterID"] = clusterID
		this.viewLayout("cluster/statistic", "default")
	} else {
		clusterModel := models.Cluster{}
		cluster, err := clusterModel.GetClusterByClusterId(clusterID)
		if err != nil {
			this.JsonError(err)
		}
		clusterIP := cluster["ip"]
		url := fmt.Sprintf("https://%s:%s/traffic/count", clusterIP, beego.AppConfig.String("cluster.api.port"))
		body, _, err := common.HttpGet(url)
		if err != nil {
			this.JsonError(err)
		}
		status := map[string]interface{}{}
		err = json.Unmarshal(body, &status)
		if err != nil {
			this.JsonError(err)
		}
		_data, ok := status["data"]
		if !ok {
			this.JsonError("no data")
		}
		data := _data.(map[string]interface{})
		rs, err := models.DB.Query(models.DB.AR().From("cluster").Where(map[string]interface{}{
			"cluster_id": clusterID,
		}))
		if err != nil {
			this.JsonError(err)
		}
		if rs.Len() == 0 {
			this.JsonError("no cluster")
		}
		data["cluster"] = rs.Row()
		status["data"] = data
		bytes, _ := json.Marshal(status)
		this.Ctx.WriteString(string(bytes))
	}

}
func (this *ClusterController) Forbidden() {
	clusterModel := models.Cluster{}
	clusterId := this.GetString("cluster_id")
	err := clusterModel.Forbidden(clusterId)
	if err != nil {
		this.JsonError(err)
	}
	this.JsonSuccess("")

}

func (this *ClusterController) Review() {
	clusterModel := models.Cluster{}
	clusterId := this.GetString("cluster_id")
	err := clusterModel.Review(clusterId)
	if err != nil {
		this.JsonError(err)
	}
	this.JsonSuccess("")

}

func (this *ClusterController) Delete() {
	clusterModel := models.Cluster{}
	clusterId := this.GetString("cluster_id")
	HasTunnel, err := clusterModel.HasTunnel(clusterId)
	if err != nil {
		this.JsonError(err)
	}
	if HasTunnel {
		this.JsonError("Tunnel引用非空,不能删除")
	}
	err = clusterModel.Delete(clusterId)
	if err != nil {
		this.JsonError(err)
	}
	this.JsonSuccess("")

}
func (this *ClusterController) Add() {
	clusterModel := models.Cluster{}
	if this.Ctx.Input.IsPost() {
		_, data := this.getClusterFromPost(false)
		_, err := clusterModel.Insert(data)
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
		this.viewLayout("cluster/form", "form")
	}

}
func (this *ClusterController) Edit() {
	clusterModel := models.Cluster{}
	clusterId := this.GetString("cluster_id")
	if this.Ctx.Input.IsPost() {
		_, data := this.getClusterFromPost(true)
		_, err := clusterModel.Update(clusterId, data)
		if err != nil {
			this.JsonError(err)
		}
		this.JsonSuccess("")
	} else {
		cluster, err := clusterModel.GetClusterByClusterId(clusterId)
		if err != nil {
			this.JsonError(err, "")
		}
		if len(cluster) == 0 {
			this.ViewError("Cluster不存在")
		}
		regionModel := models.Region{}
		regions, err := regionModel.GetSubRegions()
		if err != nil {
			this.ViewError(err.Error())
		}
		this.Data["regions"] = regions
		this.Data["cluster"] = cluster
		this.Data["action"] = "edit"
		this.viewLayout("cluster/form", "form")
	}

}

func (this *ClusterController) List() {
	clusterModel := models.Cluster{}

	clusters, err := clusterModel.GetAllClusters()
	if err != nil {
		this.ViewError(err.Error())
	}
	regionModel := models.Region{}
	regions, err := regionModel.GetSubRegions()
	if err != nil {
		this.ViewError(err.Error())
	}
	r := map[string]string{
		"region_id": "0",
	}
	regions = append(regions, r)
	this.Data["regions"] = regions
	this.Data["clusters"] = clusters
	this.view("cluster/list")
}

func (this *ClusterController) getClusterFromPost(isUpdate bool) (clusterId string, cluster map[string]interface{}) {
	cluster = map[string]interface{}{
		"name":      this.GetString("name"),
		"ip":        this.GetString("ip"),
		"region_id": this.GetString("region_id"),
	}
	errs := validation.Errors{
		"名称": validation.Validate(cluster["name"],
			validation.Required.Error("不能为空"),
			validation.Match(regexp.MustCompile("^.{1,15}$")).Error("长度必须是1-15字符")),
		"IP": validation.Validate(cluster["ip"],
			validation.Required.Error("不能为空"),
			is.IPv4),
		"区域": validation.Validate(cluster["region_id"],
			validation.Match(regexp.MustCompile("^[1-9][0-9]*$")).Error("格式错误")),
	}
	err := errs.Filter()
	if err != nil {
		this.JsonError(err)
	}
	if isUpdate {
		cluster["update_time"] = time.Now().Unix()
	} else {
		cluster["is_disable"] = 1
		cluster["create_time"] = time.Now().Unix()
	}
	return
}
