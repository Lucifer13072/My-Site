package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

var tpl = template.Must(template.ParseFiles("Main/front/index.html", "Main/Admin/front/loginAdmin.html", "Main/front/reg.html", "Main/front/auth.html"))

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

func authHandle(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "auth.html", nil)
	if err != nil {
		http.Error(w, "Ошибка при рендере страницы", http.StatusInternalServerError)
	}
}
func regHandle(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "reg.html", nil)
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

	// Роуты
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/admin/login", adminHandler)
	mux.HandleFunc("/auth", authHandle)
	mux.HandleFunc("/register", regHandle)

	log.Println("Сервер запущен на порту", port)
	log.Print(http.ListenAndServe(":"+port, mux))
}
