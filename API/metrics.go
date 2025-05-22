package api

import (
	"log"
	"time"
)

type Dataset struct {
	Label   string  `json:"label"`
	Data    []int   `json:"data"`
	Fill    bool    `json:"fill"`
	Tension float64 `json:"tension"`
}

// ChartData — структура для передачи на фронт
type ChartData struct {
	Labels   []string  `json:"labels"`
	Datasets []Dataset `json:"datasets"`
}

func AllUsersMetrics() int {
	db := InitDB()
	defer db.Close()

	var users int
	err := db.QueryRow("SELECT count(*) FROM users").Scan(&users)
	if err != nil {
		log.Printf("не смогли посчитать пользователей: %v", err)
		return 0
	}

	return users
}

func AllMoneyaddMetrics() float64 {
	db := InitDB()
	defer db.Close()

	var money float64
	err := db.QueryRow("SELECT SUM(money) FROM users").Scan(&money)
	if err != nil {
		log.Printf("не смогли посчитать пользователей: %v", err)
		return 0
	}

	return money
}

func VisitorsMakeMetrics(older bool) {
	db := InitDB() // InitDB возвращает *sql.DB
	defer db.Close()

	query := `INSERT INTO visitors (dateTime, older)VALUES (?, ?)`
	if _, err := db.Exec(query, time.Now(), older); err != nil {
		log.Printf("Ошибка при вставке метрики посетителя: %v", err)
	}
}

func GetVisitorsMetrics() *ChartData {
	db := InitDB()
	defer db.Close()

	query := `
    SELECT
        DATE_FORMAT(dateTime, '%b') AS month,
        SUM(CASE WHEN older = 0 THEN 1 ELSE 0 END) AS new_visitors,
        SUM(CASE WHEN older = 1 THEN 1 ELSE 0 END) AS old_visitors
    FROM visitors
    GROUP BY month
    ORDER BY MIN(dateTime);
`

	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		log.Fatalln(err)
	}
	var labels []string
	var newData, oldData []int

	for rows.Next() {
		var month string
		var newCount, oldCount int
		err := rows.Scan(&month, &newCount, &oldCount)
		if err != nil {
			log.Println("Ошибка сканирования строки:", err)
			continue
		}

		labels = append(labels, month)
		newData = append(newData, newCount)
		oldData = append(oldData, oldCount)
	}

	chart := &ChartData{
		Labels: labels,
		Datasets: []Dataset{
			{
				Label:   "New Visitor",
				Data:    newData,
				Fill:    true,
				Tension: 0.4,
			},
			{
				Label:   "Old Visitor",
				Data:    oldData,
				Fill:    true,
				Tension: 0.4,
			},
		},
	}

	return chart
}
