package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

var tpl = template.Must(template.ParseFiles("Main/front/index.html", "Main/Admin/front/loginAdmin.html"))

func indexHandler(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, "Ошибка при рендере страницы", http.StatusInternalServerError)
	}
}

// Обработчик страницы админки (простой вариант без аутентификации)
func adminHandler(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "loginAdmin.html", nil)
	if err != nil {
		http.Error(w, "Ошибка при рендере страницы", http.StatusInternalServerError)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	mux := http.NewServeMux()

	// Раздача статических файлов
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("Main/front/assets"))))
	mux.Handle("/blog/images/", http.StripPrefix("/blog/images/", http.FileServer(http.Dir("API/blog/images/"))))
	mux.Handle("/admin/assets/", http.StripPrefix("/admin/assets/", http.FileServer(http.Dir("Main/front/admin/assets/"))))
	mux.HandleFunc("/auth/register", registerHandler)
	mux.HandleFunc("/auth/login", loginHandler)
	mux.HandleFunc("/auth/google", googleLoginHandler)
	mux.HandleFunc("/auth/google/callback", googleCallbackHandler)

	// Роуты
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/admin/login", adminHandler)

	log.Println("Сервер запущен на порту", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
