package main

import (
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"tytyber.ru/API"
)

var store = sessions.NewCookieStore([]byte("super-secret-key"))

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var creds struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Ошибка декодирования", http.StatusBadRequest)
		return
	}

	var storedPassword string
	var userID int
	var rules int

	err := db.QueryRow("SELECT id, password, rules FROM users WHERE name = ?", creds.Name).
		Scan(&userID, &storedPassword, &rules)
	if err != nil {
		http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(creds.Password)); err != nil {
		http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
		return
	}

	session, _ := store.Get(r, "session")
	session.Values["authenticated"] = true
	session.Values["userID"] = userID
	session.Values["rules"] = rules
	session.Save(r, w)

	w.WriteHeader(http.StatusOK)
}
