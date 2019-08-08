package business

import (
	"github.com/astaxie/beego"
	"anytunnel/at-common"
	"encoding/json"
	"fmt"
	"anytunnel/at-web/app/utils"
)

type BusinessBase struct {

}

func NewBase() *BusinessBase  {
	return &BusinessBase{}
}

//Get Request Api
func (this *BusinessBase) GetRequest(confKey string, urlQuerys map[string]string) (data interface{}, err error) {
	uri := beego.AppConfig.String(confKey)
	if(uri == "") {
		return data, fmt.Errorf("%s", "uri conf error")
	}
	query := utils.NewUrls().HttpQueryBuild(urlQuerys)
	body, code, err := at_common.HttpGet(uri + "?" + query)
	if(code != 200) {
		return data, fmt.Errorf("%s", "get httpcode error")
	}
	if(err != nil) {
		return data, err
	}
	var results map[string]interface{}
	json.Unmarshal(body, &results)
	if(results["code"].(float64) != 1) {
		return data, fmt.Errorf("%s", results["message"].(string))
	}
	return results["data"], nil
}

//Post Request Api
func (this *BusinessBase) PostRequest(confKey string, urlQuerys map[string]string, header map[string]string) (data interface{}, err error) {
	uri := beego.AppConfig.String(confKey)
	if(uri == "") {
		return data, fmt.Errorf("%s", "uri conf error")
	}
	body, code, err := at_common.HttpPost(uri, urlQuerys, header)
	if(code != 200) {
		return data, fmt.Errorf("%s", "httpcode error")
	}
	if(err != nil) {
		return data, err
	}
	var results map[string]interface{}
	json.Unmarshal(body, &results)
	if(results["code"].(float64) != 1) {
		return data, fmt.Errorf("%s", results["message"].(string))
	}
	return results["data"], nil
}
