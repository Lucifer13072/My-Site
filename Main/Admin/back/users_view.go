package mainadmin

import (
	"bytes"
	_ "github.com/go-sql-driver/mysql" // MySQL драйвер
	"html/template"
	api "tytyber.ru/API"
)

// User — структура, соответствующая записи в таблице users
type User struct {
	ID    int
	Name  string
	Email string
	Money string
	Rules int // avatar_url
}

// fetchUsers достаёт всех пользователей из MySQL
func FetchUsers() ([]User, error) {
	db := api.InitDB()

	rows, err := db.Query("SELECT id, name, email, money, rules FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Money, &u.Rules); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

// HTML-шаблон карточки
var cardTpl = template.Must(template.New("card").Parse(`
<div class="user-card">
	<div style="margin-top: 10px; margin-bottom: 10px" class="container-fluid">
		<div class="row">
			<div class="col-8">
				<img src="/assets/userAvatars/{{.Name}}.jpg" class="avatar">
				<h2 class="card-title">{{.Name}}</h2>
				<p class="card-text">Почта: {{.Email}}</p>
				<p class="card-text">Деньги: {{.Money}}₽</p>
				<p class="card-text">Права: {{.Rules}} ID: {{.ID}}</p>
			</div>
			<div class="col-4">
				<form action="/admin/deleteusers" method="POST" onsubmit="return confirm('Точно удалить этого юзера?');">
					<input type="hidden" name="user_id" value="{{.ID}}">
					<button class="card-button" type="submit">Delete user</button>
			  	</form>
				<form action="/admin/rulesreplace?id={{.ID}}" method="POST" onsubmit="return confirm('Точно поменять права у этого юзера?');">
				  <input class="input-form" type="number" name="amount" placeholder="rules" required min="0" max="3" step="1">
				  <button class="card-button" type="submit">Назначить</button>
				</form>
			</div>
		</div>
	</div>
</div>
`))

func BuildUserCards(users []User) (template.HTML, error) {
	var buf bytes.Buffer
	for _, u := range users {
		if err := cardTpl.Execute(&buf, u); err != nil {
			return "", err
		}
	}
	// Превращаем строку в template.HTML, чтобы не экранировать div’
	return template.HTML(buf.String()), nil
}
