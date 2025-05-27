package main

import ( // для хеширования паролей
	"fmt"
	"github.com/gorilla/sessions"
	"html/template"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	api "tytyber.ru/API"
	admin "tytyber.ru/Main/Admin/back"
)

var tpl = template.Must(template.ParseFiles(
	"Main/front/index.html",
	"Main/Admin/front/admin.html",
	"Main/Admin/front/users.html",
	"Main/front/reg.html",
	"Main/front/auth.html",
	"Main/front/profile.html",
	"Main/front/blog.html",
	"Main/front/404.html",
	"Main/front/403.html"))

var store = sessions.NewCookieStore([]byte("super-secret-key"))

// statusRecorder нужен, чтобы перехватить код ответа
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

func init() {
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   31536000, // кука живёт 100 час
		HttpOnly: true,     // нельзя читать куку из JS
		Secure:   false,    // ОБЯЗАТЕЛЬНО false для HTTP
	}
}

// existsFile проверяет, существует ли файл (или директория) по заданному пути
func existsFile(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")

	isLoggedIn := session.Values["username"] != nil

	data := map[string]interface{}{
		"isLoggedIn": isLoggedIn,                 // Если имя пользователя есть в сессии
		"username":   session.Values["username"], // Имя пользователя
		"rules":      session.Values["rules"],    // Роли
	}

	if isLoggedIn {
		api.VisitorsMakeMetrics(true)
	} else {
		api.VisitorsMakeMetrics(false)
	}

	err := tpl.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		http.Error(w, "Ошибка при рендере страницы", http.StatusInternalServerError)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	session.Options.MaxAge = -1
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Обработчик страницы админки (простой вариант без аутентификации)
func adminHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, "Ошибка сессии", http.StatusInternalServerError)
		return
	}

	rules, ok := session.Values["rules"].(int)
	if !ok || rules > 2 {
		http.Error(w, "Доступ запрещён", http.StatusForbidden)
		return
	}
	// Проверяем, залогинен ли пользователь
	username, ok := session.Values["username"].(string)
	if !ok || username == "" {
		// Если нет — редиректим на страницу авторизации
		http.Redirect(w, r, "/auth", http.StatusSeeOther)
		return
	}

	db := api.InitDB()
	defer db.Close()

	profile, err := api.GetUserProfile(db, username)
	if err != nil {
		log.Printf("Не удалось получить профиль %s: %v", username, err)
		http.Error(w, "Профиль не найден", http.StatusNotFound)
		return
	}
	var avatar string
	if existsFile("Main/front/assets/userAvatars/" + username + ".jpg") {
		avatar = "assets/userAvatars/" + username + ".jpg"
	} else {
		avatar = "assets/images/icon-profile.png"
	}

	usersMetrics := api.AllUsersMetrics()
	money := api.AllMoneyaddMetrics()
	visitorMetrics := api.GetVisitorsMetrics()

	datas := map[string]interface{}{
		"isLoggedIn": true,
		"name":       profile.Nickname,
		"email":      profile.Email,
		"rules":      profile.Rules,
		"date":       profile.DateRegistry.Format("02.01.2006 15:04"),
		"user_key":   profile.UserKey,
		"avatar":     avatar,
		"usersmeric": usersMetrics,
		"money":      money,
		"visitors":   visitorMetrics,
	}

	if err = tpl.ExecuteTemplate(w, "admin.html", datas); err != nil {
		http.Error(w, "Ошибка при рендере страницы", http.StatusInternalServerError)
	}
}

func authHandle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		err := tpl.ExecuteTemplate(w, "auth.html", nil)
		if err != nil {
			http.Error(w, "Ошибка при рендере страницы", http.StatusInternalServerError)
		}
	case http.MethodPost:
		name := r.FormValue("username")
		password := r.FormValue("password")

		rules, err := AuthenticateUser(name, password)
		if err != nil {
			log.Println("Ошибка авторизации:", err)
			http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
			return
		}

		// Сохраняем сессию
		session, _ := store.Get(r, "session")
		session.Values["username"] = name
		session.Values["rules"] = rules
		err = session.Save(r, w)
		if err != nil {
			log.Println("Ошибка сохранения сессии:", err)
			http.Error(w, "Ошибка сессии", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func regHandle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		err := tpl.ExecuteTemplate(w, "reg.html", nil)
		if err != nil {
			http.Error(w, "Ошибка при рендере страницы", http.StatusInternalServerError)
		}
	case http.MethodPost:
		name := r.FormValue("username")
		mail := r.FormValue("email")
		password := r.FormValue("password")

		err := RegisterUser(name, password, mail)
		if err != nil {
			log.Println("Ошибка регистрации:", err)
			http.Error(w, "Ошибка регистрации: "+err.Error(), http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, "/auth", http.StatusSeeOther)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем сессию
	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, "Ошибка сессии", http.StatusInternalServerError)
		return
	}

	// Проверяем, залогинен ли пользователь
	username, ok := session.Values["username"].(string)
	if !ok || username == "" {
		// Если нет — редиректим на страницу авторизации
		http.Redirect(w, r, "/auth", http.StatusSeeOther)
		return
	}

	db := api.InitDB()
	defer db.Close()

	profile, err := api.GetUserProfile(db, username)
	if err != nil {
		log.Printf("Не удалось получить профиль %s: %v", username, err)
		http.Error(w, "Профиль не найден", http.StatusNotFound)
		return
	}
	walletop, err := api.GetWalletOperations(db, username)
	if err != nil {
		log.Printf("Не удалось получить операции кошелька: %v", err)
		// Можно либо вернуть ошибку, либо просто оставить пустой список
		walletop = []api.WalletOperation{}
	}
	var avatar string
	if existsFile("Main/front/assets/userAvatars/" + username + ".jpg") {
		avatar = "assets/userAvatars/" + username + ".jpg"
	} else {
		avatar = "assets/images/icon-profile.png"
	}

	datas := map[string]interface{}{
		"isLoggedIn":       true,
		"name":             profile.Nickname,
		"email":            profile.Email,
		"rules":            profile.Rules,
		"date":             profile.DateRegistry.Format("02.01.2006 15:04"),
		"user_key":         profile.UserKey,
		"money":            fmt.Sprintf("%.2f", profile.UserMoney),
		"walletid":         profile.WalletID,
		"walletOperations": walletop,
		"avatar":           avatar,
	}

	// Рендерим профиль
	if err := tpl.ExecuteTemplate(w, "profile.html", datas); err != nil {
		http.Error(w, "Ошибка при рендере страницы", http.StatusInternalServerError)
	}
}

func avatarHandler(w http.ResponseWriter, r *http.Request) {
	// Максимальный размер загружаемого файла — 10MB
	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, "Ошибка сессии", http.StatusInternalServerError)
		return
	}

	// Проверяем, залогинен ли пользователь
	username, ok := session.Values["username"].(string)
	if !ok || username == "" {
		// Если нет — редиректим на страницу авторизации
		http.Redirect(w, r, "/auth", http.StatusSeeOther)
		return
	}

	r.ParseMultipartForm(10 << 20)

	file, _, err := r.FormFile("avatar")
	if err != nil {
		http.Error(w, "Ошибка загрузки файла", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Декодим картинку
	img, format, err := image.Decode(file)
	if err != nil {
		http.Error(w, "Не удалось распознать изображение", http.StatusUnsupportedMediaType)
		return
	}

	// Проверка размеров
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	if width > 500 || height > 500 {
		http.Error(w, "Размер изображения превышает 500x500", http.StatusBadRequest)
		return
	}

	// Создание папки, если нет
	outputDir := "Main/front/assets/userAvatars"
	os.MkdirAll(outputDir, os.ModePerm)

	var ext string
	switch strings.ToLower(format) {
	case "jpeg", "jpg":
		ext = ".jpg"
	case "png":
		ext = ".png"
	default:
		http.Error(w, "Формат не поддерживается", http.StatusUnsupportedMediaType)
		return
	}

	// Создание файла
	savePath := filepath.Join(outputDir, username+ext)

	outFile, err := os.Create(savePath)
	if err != nil {
		http.Error(w, "Не удалось сохранить файл", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	// Сохраняем в нужном формате
	switch strings.ToLower(format) {
	case "jpeg", "jpg":
		err = jpeg.Encode(outFile, img, nil)
	case "png":
		err = png.Encode(outFile, img)
	default:
		http.Error(w, "Формат изображения не поддерживается", http.StatusUnsupportedMediaType)
		return
	}

	if err != nil {
		http.Error(w, "Ошибка сохранения изображения", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func blogHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")

	isLoggedIn := session.Values["username"] != nil

	posts, err := api.GetPosts()
	if err != nil {
		http.Error(w, "Не удалось получить список пользователей: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Генерируем HTML карточек
	cardsHTML, err := api.BuildPostsCards(posts)
	if err != nil {
		http.Error(w, "Ошибка рендеринга карточек: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Передаём в основной шаблон
	data := map[string]interface{}{
		"posts":      cardsHTML,
		"isLoggedIn": isLoggedIn,                 // Если имя пользователя есть в сессии
		"username":   session.Values["username"], // Имя пользователя
		"rules":      session.Values["rules"],
	}

	err = tpl.ExecuteTemplate(w, "403.html", data)
	if err != nil {
		http.Error(w, "Ошибка при рендере страницы", http.StatusInternalServerError)
	}
}

func adminUsers(w http.ResponseWriter, r *http.Request) {

	users, err := admin.FetchUsers()
	if err != nil {
		http.Error(w, "Не удалось получить список пользователей: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Генерируем HTML карточек
	cardsHTML, err := admin.BuildUserCards(users)
	if err != nil {
		http.Error(w, "Ошибка рендеринга карточек: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Передаём в основной шаблон
	data := map[string]interface{}{
		"userscards": cardsHTML,
	}

	err = tpl.ExecuteTemplate(w, "users.html", data)
	if err != nil {
		http.Error(w, "Ошибка при рендере страницы", http.StatusInternalServerError)
	}
}

func proxyAPIHandler(w http.ResponseWriter, r *http.Request) {
	api.UserIndifecation()
}

func replaceRulesHandle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Если нужен GET – возвращаем список, но, наверно, GET вам не нужен
		http.Error(w, "GET не поддерживается", http.StatusMethodNotAllowed)
		return

	case http.MethodPost:
		// Парсим form и берём id из query
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Ошибка разбора формы", http.StatusBadRequest)
			return
		}
		idStr := r.URL.Query().Get("id")
		amount := r.FormValue("amount")
		if idStr == "" || amount == "" {
			http.Error(w, "Нет id или amount", http.StatusBadRequest)
			return
		}

		db := api.InitDB()
		defer db.Close()

		// Обновляем правило по user_id
		query := `UPDATE users SET rules = $1 WHERE id = $2;`
		if _, err := db.Exec(query, amount, idStr); err != nil {
			log.Printf("Ошибка обновления rules для user_id=%s: %v", idStr, err)
			http.Error(w, "Не удалось обновить права", http.StatusInternalServerError)
			return
		}

		// Редирект обратно на страницу со всеми юзерами
		http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
		return

	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func deleteUserHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Разбираем форму
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Ошибка разбора формы", http.StatusBadRequest)
		return
	}
	userID := r.FormValue("user_id")
	if userID == "" {
		http.Error(w, "Не передан user_id", http.StatusBadRequest)
		return
	}

	// Подключаемся к БД
	db := api.InitDB()
	defer db.Close()

	// Удаляем юзера
	res, err := db.Exec(`DELETE FROM users WHERE id = $1;`, userID)
	if err != nil {
		log.Printf("Ошибка удаления user_id=%s: %v", userID, err)
		http.Error(w, "Не удалось удалить пользователя", http.StatusInternalServerError)
		return
	}
	// Опционально можно проверить, затронулась ли строка
	rows, _ := res.RowsAffected()
	if rows == 0 {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	// Редирект обратно на список юзеров
	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}

func notFoundPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)

	// опционально: передаём данные в шаблон
	data := struct {
		URL string
	}{
		URL: r.URL.Path,
	}

	// рендерим именно 404.html
	err := tpl.ExecuteTemplate(w, "404.html", data)
	if err != nil {
		// если шаблон упал — возвращаем простой текст
		http.Error(w, "Ошибка при рендере страницы 404", http.StatusInternalServerError)
	}
}

func forbiddenHandler(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "403.html", nil)
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
	mux.Handle("/admin/assets/", http.StripPrefix("/admin/assets/", http.FileServer(http.Dir("Main/Admin/front/assets/"))))

	// Роуты
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/admin", adminHandler)
	mux.HandleFunc("/auth", authHandle)
	mux.HandleFunc("/register", regHandle)
	mux.HandleFunc("/profile", profileHandler)
	mux.HandleFunc("/logout", logoutHandler)
	mux.HandleFunc("/addmoney", api.TopUpHandler)
	mux.HandleFunc("/upload-avatar", avatarHandler)
	mux.HandleFunc("/blog", blogHandler)
	mux.HandleFunc("/proxyapi", proxyAPIHandler)
	mux.HandleFunc("/admin/users", adminUsers)
	mux.HandleFunc("/admin/rulesreplace", replaceRulesHandle)
	mux.HandleFunc("/admin/deleteusers", deleteUserHandle)

	// Обернём mux в statusRecorder, чтобы ловить 404
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h, pattern := mux.Handler(r)

		switch {
		case pattern == "":
			// Ни один маршрут не подошел
			notFoundPage(w, r)
			return
		case pattern == "/" && r.URL.Path != "/":
			// Если бы мы ловили "/" — а это не /
			notFoundPage(w, r)
			return
		default:
			// Всё ок — передаём исполнению оригинальный хендлер
			h.ServeHTTP(w, r)
		}
	})

	log.Println("🚀 Сервер запущен на порту", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
