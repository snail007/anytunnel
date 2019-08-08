package models

import (
	"anytunnel/at-admin/app/models"

	"github.com/snail007/go-activerecord/mysql"
)

var G *mysql.DBGroup
var DB *mysql.DB

func init() {
	G = models.G
	DB = G.DB("base")
}
