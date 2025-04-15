package main

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
	"tytyber.ru/API"
)

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var creds struct {
		Name     string `json:"name"`
		Password string `json:"password"`
		Rules    int    `json:"rules"`
	}

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Ошибка декодирования", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Ошибка хеширования пароля", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("INSERT INTO users (name, password, date, rules) VALUES (?, ?, ?, ?)",
		creds.Name, string(hashedPassword), time.Now(), creds.Rules)
	if err != nil {
		http.Error(w, "Пользователь уже существует", http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
