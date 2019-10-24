package dao

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
)

var dsn = "root:123456@tcp(mysql:3306)/digimon?charset=utf8&loc=Local"

var db *sqlx.DB

func Init() {
	var err error
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatalln(err)
	}
}
