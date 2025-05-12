package api

import (
	"github.com/gorilla/sessions"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var store = sessions.NewCookieStore([]byte("super-secret-key"))

var tpl = template.Must(template.ParseFiles(
	"Main/front/index.html",
	"Main/Admin/front/admin.html",
	"Main/front/reg.html",
	"Main/front/auth.html",
	"Main/front/profile.html"))

func TopUpHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Показываем форму пополнения
		err := tpl.ExecuteTemplate(w, "profile.html", nil)
		if err != nil {
			http.Error(w, "Ошибка при отображении формы", http.StatusInternalServerError)
		}
		return

	case http.MethodPost:
		// === Проверяем сессию ===
		session, err := store.Get(r, "session")
		if err != nil {
			http.Error(w, "Ошибка сессии", http.StatusInternalServerError)
			return
		}
		username, ok := session.Values["username"].(string)
		if !ok || username == "" {
			http.Redirect(w, r, "/auth", http.StatusSeeOther)
			return
		}

		// === Читаем сумму ===
		amountStr := r.FormValue("amount")
		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil || amount <= 0 {
			http.Error(w, "Некорректная сумма", http.StatusBadRequest)
			return
		}

		// Описание транзакции — всегда одно и то же
		const description = "Пополнение баланса"

		db := InitDB()
		defer db.Close()

		// === Получаем user_id ===
		var userID int
		err = db.QueryRow("SELECT id FROM users WHERE name = ?", username).Scan(&userID)
		if err != nil {
			http.Error(w, "Пользователь не найден", http.StatusNotFound)
			log.Println("Ошибка при получении user_id:", err)
			return
		}

		// === Транзакция: UPDATE + INSERT ===
		tx, err := db.Begin()
		if err != nil {
			http.Error(w, "Ошибка транзакции", http.StatusInternalServerError)
			return
		}

		_, err = tx.Exec("UPDATE users SET money = money + ? WHERE id = ?", amount, userID)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Ошибка обновления баланса", http.StatusInternalServerError)
			return
		}

		_, err = tx.Exec(
			"INSERT INTO money_history (user_id, datetime, money, description) VALUES (?, NOW(), ?, ?)",
			userID, amount, description,
		)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Ошибка записи в историю", http.StatusInternalServerError)
			return
		}

		if err := tx.Commit(); err != nil {
			http.Error(w, "Ошибка коммита", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/profile", http.StatusSeeOther)

	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}
