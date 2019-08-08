package models

import (
	"fmt"
	"os"

	"github.com/astaxie/beego"
	"github.com/snail007/go-activerecord/mysql"
)

var G *mysql.DBGroup

func init() {
	host := beego.AppConfig.String("db_host")
	port, _ := beego.AppConfig.Int("db_port")
	user := beego.AppConfig.String("db_user")
	pass := beego.AppConfig.String("db_pass")
	dbname := beego.AppConfig.String("db_name")
	dbTablePrefix := beego.AppConfig.String("db_table_prefix")
	maxIdle, _ := beego.AppConfig.Int("db_conn_max_idle")
	maxConn, _ := beego.AppConfig.Int("db_conn_max_connection")
	G = mysql.NewDBGroup("default")
	cfg := mysql.NewDBConfigWith(host, port, dbname, user, pass)
	cfg.SetMaxIdleConns = maxIdle
	cfg.SetMaxOpenConns = maxConn
	cfg.TablePrefix = dbTablePrefix
	cfg.TablePrefixSqlIdentifier = "__PREFIX__"

	//add database anytunnel_system
	err := G.Regist("default", cfg)
	if err != nil {
		beego.Error(fmt.Errorf("regist db error:%s,with config : %v", err, cfg))
		os.Exit(100)
	}

	//add database anytunnel_base
	cfg.Database = beego.AppConfig.String("db.base.name")
	cfg.TablePrefix = beego.AppConfig.String("db.base.table_prefix")
	err = G.Regist("base", cfg)
	if err != nil {
		beego.Error(fmt.Errorf("regist db error:%s,with config : %v", err, cfg))
		os.Exit(100)
	}
}
