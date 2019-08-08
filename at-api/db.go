package main

import (
	"fmt"
	"os"

	"github.com/snail007/go-activerecord/mysql"
)

var G *mysql.DBGroup
var db *mysql.DB

func initDB() {
	G = mysql.NewDBGroup("default")

	host := cfg.GetString("database.base.host")
	port := cfg.GetInt("database.base.port")
	dbName := cfg.GetString("database.base.dbName")
	user := cfg.GetString("database.base.user")
	pass := cfg.GetString("database.base.pass")
	maxIdle := cfg.GetInt("database.base.maxIdle")
	maxConn := cfg.GetInt("database.base.maxConn")
	dbTablePrefix := cfg.GetString("database.base.dbTablePrefix")

	mysqlConfig := mysql.NewDBConfigWith(host, port, dbName, user, pass)
	mysqlConfig.SetMaxIdleConns = maxIdle
	mysqlConfig.SetMaxOpenConns = maxConn
	mysqlConfig.TablePrefix = dbTablePrefix
	mysqlConfig.TablePrefixSqlIdentifier = "__PREFIX__"

	err := G.Regist("default", mysqlConfig)
	if err != nil {
		fmt.Println("regist db error!" + err.Error())
		os.Exit(100)
	}
	db = G.DB()
}
