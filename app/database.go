package app

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"gopkg.in/gorp.v1"
)

var DB *gorp.DbMap

func InitDatabase() *gorp.DbMap {
	// username:password@protocol(address)/dbname?params
	dsn := "root:@tcp(localhost:3306)/fa_diel?parseTime=true"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic("failed to connect database")
	}

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

	DB = dbmap

	return DB
}

func GetDB() *gorp.DbMap {
	return DB
}
