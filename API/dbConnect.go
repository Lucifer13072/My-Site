package api

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func InitDB() *sql.DB {
	// Задаём параметры подключения: юзер, пароль, хост, порт, имя базы
	dsn := "root:@tcp(192.168.1.9:3306)/tytyber_api"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Ошибка при открытии соединения с MySQL:", err)
	}

	// Проверим, что база реально доступна
	if err := db.Ping(); err != nil {
		log.Fatal("MySQL не отвечает:", err)
	}

	return db
}
