package app

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"gopkg.in/gorp.v1"
)

var DB *gorp.DbMap

func InitDatabase() *gorp.DbMap {
	// username:password@protocol(address)/dbname?params
	dsn := os.Getenv("DB_USERNAME") + ":" + os.Getenv("DB_PASSWORD") + "@tcp(" + os.Getenv("DB_HOST") + ":3306)/" + os.Getenv("DB_DATABASE") + "?parseTime=true"
	db, err := sql.Open(os.Getenv("DB_CONNECTION"), dsn)
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
