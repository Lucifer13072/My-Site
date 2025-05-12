package main

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"tytyber.ru/API"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid password")
	ErrInvalidUserRole = errors.New("invalid user role")
)

func AuthenticateUser(name, password string) (int, error) {
	if name == "" || password == "" {
		return -1, ErrInvalidInput
	}

	db := api.InitDB()
	defer db.Close()

	var hashedPassword string
	var rules int
	err := db.QueryRow("SELECT password, rules FROM users WHERE name = ?", name).Scan(&hashedPassword, &rules)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1, ErrUserNotFound
		}
		return -1, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return -1, ErrInvalidPassword
	}

	// Проверка допустимых ролей (0-3)
	if rules < 0 || rules > 3 {
		return -1, ErrInvalidUserRole
	}

	return rules, nil
}
