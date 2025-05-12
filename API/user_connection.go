package api

import (
	"database/sql"
	"fmt"
	"time"
)

type User struct {
	Nickname     string
	Email        string
	Rules        string
	DateRegistry time.Time
	UserKey      string
	UserMoney    float64
	WalletID     string
}

type WalletOperation struct {
	Date        time.Time
	Money       float64
	Description string
}

func GetUserProfile(db *sql.DB, username string) (*User, error) {
	u := &User{}
	var dateStr string

	query := `
        SELECT 
            name,
            email,
            rules,
            date,
            user_key,
            money,
            wallet_id
        FROM users
        WHERE name = ?
    `
	// Сканим date_registry в dateStr, а не сразу в time.Time
	err := db.QueryRow(query, username).Scan(
		&u.Nickname,
		&u.Email,
		&u.Rules,
		&dateStr, // <- сюда
		&u.UserKey,
		&u.UserMoney,
		&u.WalletID,
	)
	if err != nil {
		return nil, err
	}
	// Парсим строку в time.Time
	t, err := time.Parse("2006-01-02 15:04:05", dateStr)
	if err != nil {
		return nil, fmt.Errorf("не удалось распарсить дату '%s': %w", dateStr, err)
	}
	u.DateRegistry = t

	return u, nil
}

func GetWalletOperations(db *sql.DB, username string) ([]WalletOperation, error) {
	var operations []WalletOperation
	var userid int

	// Получаем ID пользователя по имени
	err := db.QueryRow(`SELECT id FROM users WHERE name = ?`, username).Scan(&userid)
	if err != nil {
		return nil, fmt.Errorf("не удалось найти пользователя: %w", err)
	}

	// Запрашиваем все записи из money_history по user_id, отсортированные по дате
	rows, err := db.Query(`
		SELECT datetime, money, description 
		FROM money_history 
		WHERE user_id = ? 
		ORDER BY datetime DESC
	`, userid)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить операции: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var w WalletOperation
		var dateStr string

		// Читаем строку
		err := rows.Scan(&dateStr, &w.Money, &w.Description)
		if err != nil {
			return nil, fmt.Errorf("ошибка при чтении операции: %w", err)
		}

		// Парсим дату (предполагаем формат ISO-8601)
		w.Date, err = time.Parse("2006-01-02 15:04:05", dateStr)
		if err != nil {
			return nil, fmt.Errorf("ошибка парсинга даты '%s': %w", dateStr, err)
		}

		operations = append(operations, w)
	}

	// Проверка ошибок после итерации
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при обходе результатов: %w", err)
	}

	return operations, nil
}
