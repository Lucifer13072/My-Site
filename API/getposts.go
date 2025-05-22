package api

import (
	"bytes"
	_ "github.com/go-sql-driver/mysql" // MySQL драйвер
	"html/template"
)

type Post struct {
	Title    string
	Desc     string
	SortDesc string
	Image    string
	ID       string
	Likes    string
	Dislikes string
}

func GetPosts() ([]Post, error) {
	db := InitDB()

	rows, err := db.Query("SELECT id, title, description, short_desc, image, likes, dislikes FROM blog_posts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []Post
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.ID, &p.Title, &p.Desc, &p.SortDesc, &p.Image, &p.Likes, &p.Dislikes); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, rows.Err()
}

var cardTpl = template.Must(template.New("card").Parse(`
<div class="post-card">
	<div class="container">
		<img src="/blog/images/{{.Image}}" alt="{{.Title}}" class="post-card__image"/>
		<div class="post-card__content">
			<h2 class="post-card__title">{{.Title}}</h2>
			<p class="post-card__desc">{{.SortDesc}}</p>
		</div>
		<hr>
		<div class="post-card__footer">
			<form method="POST" action="/like/{{.ID}}" style="display:inline">
				<button type="submit" class="post-card__like-button">👍 {{.Likes}}</button>
			</form>
			<form method="POST" action="/dislike/{{.ID}}" style="display:inline">
				<button type="submit" class="post-card__dislike-button">👎 {{.Dislikes}}</button>
			</form>
			<a href="/posts/{{.ID}}" class="post-card__link">Читать далее</a>
		</div>
	</div>
</div>
`))

func BuildPostsCards(posts []Post) (template.HTML, error) {
	var buf bytes.Buffer
	for _, u := range posts {
		if err := cardTpl.Execute(&buf, u); err != nil {
			return "", err
		}
	}
	// Превращаем строку в template.HTML, чтобы не экранировать div’
	return template.HTML(buf.String()), nil
}
