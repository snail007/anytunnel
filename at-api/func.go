package main

import (
	utils "anytunnel/at-common"
	"anytunnel/at-common/qqwry"
	"fmt"
	"strings"
	"time"
)

func authCheckCS(csType, token string) (userID string, err error) {
	rs, err := db.Query(db.AR().From(csType).Where(map[string]interface{}{
		"is_delete": 0,
		"token":     token,
	}).Limit(0, 1))
	if err != nil {
		return
	}
	if rs.Len() == 0 {
		return "", fmt.Errorf("token error")
	}
	userID = rs.Value("user_id")
	return
}

func authCheckUser(userID string) (err error) {
	rs, err := db.Query(db.AR().From("user").Where(map[string]interface{}{
		"user_id": userID,
		//"is_active": 1,
	}).Limit(0, 1))
	if err != nil {
		return
	}
	fmt.Println(userID)
	if rs.Len() == 0 {
		return fmt.Errorf("empty user")
	}
	if rs.Value("is_forbidden") == "1" {
		reason := rs.Value("forbidden_reason")
		if reason == "" {
			reason = "your account was forbidden"
		}
		return fmt.Errorf(reason)
	}
	return
}
func authCheckUserTraffic(userID string) (err error) {
	now := time.Now().Unix()
	rs, err := db.Query(db.AR().From("package").Where(map[string]interface{}{
		"user_id":       userID,
		"start_time <=": now,
		"end_time >=":   now,
		"bytes_left >":  0,
	}).Limit(0, 1))
	if err != nil {
		return
	}
	if rs.Len() == 0 {
		return fmt.Errorf("your traffic has no more")
	}
	return
}
func authCheckIPWhiteList(ip, csType string) (ok bool, err error) {
	rs, err := db.Query(db.AR().From("ip_list").Where(map[string]interface{}{
		"ip":           ip,
		"cs_type":      csType,
		"is_forbidden": 0,
	}).Limit(1))
	if err != nil {
		return
	}
	ok = rs.Len() == 1
	return
}
func authCheckIPBlackList(ip, csType string) (ok bool, err error) {
	rs, err := db.Query(db.AR().From("ip_list").Where(map[string]interface{}{
		"ip":           ip,
		"cs_type":      csType,
		"is_forbidden": 1,
	}).Limit(1))
	if err != nil {
		return
	}
	ok = rs.Len() == 0
	return
}
func authCheckAreaWhiteList(country, csType string) (ok bool, err error) {
	rs, err := db.Query(db.AR().From("area").Where(map[string]interface{}{
		"cs_type":      csType,
		"is_forbidden": 0,
	}))
	if err != nil {
		return
	}
	for _, v := range rs.Rows() {
		if strings.HasPrefix(country, v["name"]) {
			ok = true
			return
		}
	}
	return
}
func authCheckAreaBlackList(country, csType string) (ok bool, err error) {
	rs, err := db.Query(db.AR().From("area").Where(map[string]interface{}{
		"cs_type":      csType,
		"is_forbidden": 1,
	}))
	if err != nil {
		return
	}
	ok = true
	for _, v := range rs.Rows() {
		if strings.HasPrefix(country, v["name"]) {
			ok = false
			return
		}
	}
	return
}
func authCheckRoleArea(userID string, q qqwry.ResultQQwry, csType string) (err error) {
	rs, err := db.Query(db.AR().From("user_role").Where(map[string]interface{}{
		"user_id": userID,
	}))
	if err != nil {
		return
	}
	if rs.Len() == 0 {
		return fmt.Errorf("empty user role")
	}
	roleIDs := rs.Values("role_id")
	rs, err = db.Query(db.AR().From("role").Where(map[string]interface{}{
		"role_id": roleIDs,
	}))
	if err != nil {
		return
	}
	if rs.Len() == 0 {
		return fmt.Errorf("empty role")
	}
	serverHasAll := false
	clientHasAll := false
	serverArea := []string{}
	clientArea := []string{}
	for _, role := range rs.Rows() {
		if role["server_area"] == "all" {
			serverHasAll = true
		}
		if role["client_area"] == "all" {
			clientHasAll = true
		}
		serverArea = append(serverArea, role["server_area"])
		clientArea = append(clientArea, role["client_area"])
	}
	//检查server或者client的ip是否在允许的范围内
	if csType == "server" {
		if !serverHasAll {

			serverIPArea := ""
			if q.Country != "" {
				if utils.IsChina(q.Country) {
					serverIPArea = "china"
				} else {
					serverIPArea = "foreign"
				}
			}
			found := false
			for _, v := range serverArea {
				if v == serverIPArea {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("Your IP %s was forbidden,%s", q.IP, q.Country)
			}
		}
	}
	if csType == "client" {
		if !clientHasAll {
			clientIPArea := ""
			if q.Country != "" {
				if utils.IsChina(q.Country) {
					clientIPArea = "china"
				} else {
					clientIPArea = "foreign"
				}
			}
			found := false
			for _, v := range clientArea {
				if v == clientIPArea {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("Your IP %s was forbidden,%s", q.IP, q.Country)
			}
		}
	}
	return
}
