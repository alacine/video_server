package dbops

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var (
	dbConn *sql.DB
	err    error
)

func init() {
	dbConn, err = sql.Open("mysql", "video:videoyes@tcp(172.50.0.2:3306)/video_server?charset=utf8mb4")
	if err != nil {
		panic(err.Error())
	}
}
