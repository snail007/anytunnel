package main

import (
	utils "anytunnel/at-common"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

func initTrafficReporter() {
	//item map[string]uint64{}
	//url is empty , exit reporter
	url := cfg.GetString("url.traffic")
	if url == "" {
		return
	}
	go func() {
		var trafficReporterLastData = map[string]map[string]uint64{}
		interval := cfg.GetInt("url.traffic-interval")
		for {
			time.Sleep(time.Second * time.Duration(interval))
			var connCountMap = serverConns.GetTunnelConnCountMap()
			var trafficReporterNewData = map[string]map[string]uint64{}
			var reportData = map[string]map[string]interface{}{}
			reportData = map[string]map[string]interface{}{}
			trafficReporterNewData = trafficCounter.AllData()
			var isFirst = false
			if len(trafficReporterLastData) == 0 {
				isFirst = true
				trafficReporterLastData = trafficReporterNewData
			}
			//compare
			for k, vnew := range trafficReporterNewData {
				if vold, ok := trafficReporterLastData[k]; ok {
					oldPositiveCount := vold["positive"]
					oldNegativeCount := vold["negative"]
					newPositiveCount := vnew["positive"]
					newNegativeCount := vnew["negative"]
					positive := uint64(0)
					negative := uint64(0)
					if isFirst {
						positive = newPositiveCount
						negative = newNegativeCount
					} else {
						if newPositiveCount > oldPositiveCount {
							positive = newPositiveCount - oldPositiveCount
						}
						if newNegativeCount > oldNegativeCount {
							negative = newNegativeCount - oldNegativeCount
						}
					}
					tunnelID, _ := strconv.ParseUint(k, 10, 64)
					item, ok := connCountMap[tunnelID]
					if !ok {
						item := ConnItem{}
						item.Count = 0
						item.TunnelID = tunnelID
					}
					if positive > 0 || negative > 0 || item.Count > 0 {
						reportData[k] = map[string]interface{}{
							"serverToken": item.ServerToken,
							"tunnelID":    tunnelID,
							"connCount":   item.Count,
							"positive":    positive,
							"negative":    negative,
							"interval":    interval,
						}
					}
				}
			}
			//store last data
			trafficReporterLastData = trafficReporterNewData
			//report if needed
			if len(reportData) > 0 {
				//log.Warnf("ReportData : %v", reportData)
				var code int
				var err error
				var tryCount = 0
				for tryCount <= cfg.GetInt("url.fail-retry") {
					tryCount++
					d, _ := json.Marshal(reportData)
					_, code, err = HttpPostRaw(url, string(d), nil)
					if err == nil && code == cfg.GetInt("url.success-code") {
						break
					} else if err != nil {
						log.Warnf("report traffic fail to url %s, err: %s", url, err)
					} else {
						err = fmt.Errorf("token error")
						log.Warnf("report traffic fail to url %s, code: %d, except: %d", url, code, cfg.GetInt("url.success-code"))
					}
					if err != nil && tryCount <= cfg.GetInt("url.fail-retry") {
						time.Sleep(time.Second * time.Duration(cfg.GetInt("url.fail-wait")))
					}
				}
			}
		}
	}()
}
func HttpPost(URL string, data map[string]string, header map[string]string) (body []byte, code int, err error) {
	if cfg.GetBool("url.is-internal") {
		return utils.HttpPost(URL, data, header)
	} else {
		return utils.HttpPostNotInternal(URL, data, header)
	}
}
func HttpPostRaw(URL string, data string, header map[string]string) (body []byte, code int, err error) {
	if cfg.GetBool("url.is-internal") {
		return utils.HttpPostRaw(URL, data, header, true)
	} else {
		return utils.HttpPostRaw(URL, data, header, false)
	}
}
func HttpGet(URL string) (body []byte, code int, err error) {
	return utils.HttpGet(URL)
	if cfg.GetBool("url.is-internal") {
		return utils.HttpGet(URL)
	} else {
		return utils.HttpGetNotInternal(URL)
	}
}
