package API

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func initDB() {
	var err error
	dsn := "root:@tcp(localhost:3306)/tytyber_api?parseTime=true"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("БД недоступна:", err)
	}
	fmt.Println("Подключение к БД установлено")
}
