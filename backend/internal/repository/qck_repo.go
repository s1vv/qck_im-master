package repository

import (
	"database/sql"
	"qckim-backend/internal/logger"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type QckRepo struct {
	db *sql.DB
}

func NewQckRepo(db *sql.DB) *QckRepo {

	// Настройка пула соединений
	db.SetMaxOpenConns(100)                // Максимальное количество открытых соединений
	db.SetMaxIdleConns(50)                 // Количество соединений, которые можно держать открытыми в ожидании
	db.SetConnMaxLifetime(5 * time.Minute) // Максимальное время жизни соединения

	return &QckRepo{db: db}
}

// Close закрывает соединение с базой данных
func (r *QckRepo) Close() {
	if err := r.db.Close(); err != nil {
		logger.Fatal("Ошибка закрытия БД: %v", err)
	}
}

func (r *QckRepo) GetDB() *sql.DB {
	return r.db
}
