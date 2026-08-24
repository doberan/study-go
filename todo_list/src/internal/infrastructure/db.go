package infrastructure

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func InitDB() *sql.DB {
	db, err := sql.Open("mysql", "todo:todo@tcp(127.0.0.1:3306)/todo")
	if err != nil {
		panic(err)
	}

	return db
}
