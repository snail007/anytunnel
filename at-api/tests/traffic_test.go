package test

import (
	api "anytunnel/at-api"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/snail007/go-activerecord/mysql"
)

func TestTraffic(t *testing.T) {
	raw := `{"1":{"positive":3231,"negative":5324},"2":{"negative":3997696,"positive":0}}`
	dataMap := map[string]map[string]uint{}
	err := json.Unmarshal([]byte(raw), &dataMap)
	if err != nil {
		fmt.Println("ERR:", err)
		return
	} else {
		db := api.G.DB()
		var rs *mysql.ResultSet
		rs, err := db.Query(db.AR().
			Select("user_id,tunnel_id").
			From("tunnel").
			Where(map[string]interface{}{
				"tunnel_id": []string{"1", "2"},
				"is_delete": 0,
			}))
		if err != nil {
			fmt.Println("ERR:", err)
			return
		}
		userIds := rs.MapValues("tunnel_id", "user_id")
		data := []map[string]interface{}{}
		for k, v := range dataMap {
			userID, ok := userIds[k]
			if !ok {
				continue
			}
			data = append(data, map[string]interface{}{
				"bytes +":     v["negative"] + v["positive"],
				"user_id":     userID,
				"month":       time.Now().Format("200601"),
				"update_time": time.Now().Unix(),
			})
		}
		//fmt.Println(data)
		rs, err = db.Exec(db.AR().
			UpdateBatch("traffic", data, []string{"user_id", "month"}))
		if err != nil {
			fmt.Println("ERR:", err)
			return
		}
		fmt.Println("Affected Rows : ", rs.RowsAffected)
	}
}
