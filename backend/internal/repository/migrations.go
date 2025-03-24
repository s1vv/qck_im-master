package repository

import (
	"database/sql"
	"fmt"
)

// RunMigrations создаёт необходимые таблицы в базе
func RunMigrations(db *sql.DB) error {
	queries := []string{

		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			login VARCHAR(10) UNIQUE NOT NULL,
			email VARCHAR(60) UNIQUE NOT NULL,
			password_hash VARCHAR(64) NOT NULL,
			created_at DATETIME DEFAULT NOW(),
			notifications_enabled BIT DEFAULT 1,
			is_active BIT DEFAULT 0
		);`,

		`CREATE TABLE IF NOT EXISTS qck_links (
			id SERIAL PRIMARY KEY,
			qck_link CHAR(8) UNIQUE NOT NULL,
			password_hash VARCHAR(64) NOT NULL,
			user_id BIGINT UNSIGNED,
			name VARCHAR(20) DEFAULT "qck.im",
			description VARCHAR(1000) DEFAULT "qck.im",
			is_active BIT DEFAULT 0,
			FOREIGN KEY (user_id) REFERENCES users(id) 
		);`,

		`CREATE TABLE IF NOT EXISTS refresh_tokens (
			user_id BIGINT UNSIGNED UNIQUE,
			token VARCHAR(512) NOT NULL,
			expires_at DATETIME NOT NULL,
			ip_address VARCHAR(40),
			user_agent VARCHAR(255),
			PRIMARY KEY (user_id),
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);`,

		`CREATE TABLE IF NOT EXISTS password_reset_tokens (
			user_id BIGINT UNSIGNED UNIQUE,
			token VARCHAR(64) NOT NULL,
			expires_at DATETIME NOT NULL,
			PRIMARY KEY (user_id),
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);`,

		// `CREATE UNIQUE INDEX idx_qck_link ON qck_links(qck_link);`

	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("ошибка при миграции: %w", err)
		}
	}

	return nil
}
