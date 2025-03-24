package services

import (
	"database/sql"
	"fmt"
	"qckim-backend/internal/models"
	"qckim-backend/utils/cryptPass"
	"qckim-backend/utils/jwt"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

// UserRepo - Репозиторий для работы с пользователями
type UserRepo struct {
	db *sql.DB
}

// NewUserRepo - Создает новый репозиторий для работы с пользователями
func NewUser(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

// ChangePassword - Изменение пароля пользователя
func (r *UserRepo) ChangePassword(userID int64, hashedPassword string) error {
	_, err := r.db.Exec("UPDATE users SET password_hash = ? WHERE id = ?", hashedPassword, userID)
	return err
}

// GetUserIDByRefreshToken - Получение ID пользователя по refresh токену
func (r *UserRepo) GetUserIDByRefreshToken(token string) (int64, error) {
	var userID int64
	err := r.db.QueryRow("SELECT user_id FROM refresh_tokens WHERE token = ? AND expires_at > NOW()", token).Scan(&userID)
	return userID, err
}

// SaveUser - сохранение пользователя и активация ссылки
func (r *UserRepo) SaveUser(login, email, password string, linkID int64) (int64, error) {
	hashedPassword, err := cryptPass.HashPassword(password)
	if err != nil {
		return 0, err
	}

	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	// Вставка пользователя
	result, err := tx.Exec("INSERT INTO users (login, email, password_hash) VALUES (?, ?, ?)", login, email, hashedPassword)
	if err != nil {
		// Проверяем, является ли ошибка дублированием
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			if msg := mysqlErr.Message; msg != "" {
				// Определяем, по какому полю ошибка
				if strings.Contains(msg, "for key 'users.login'") {
					return 0, fmt.Errorf("%w: %s", models.ErrDuplicateLogin, login)
				} else if strings.Contains(msg, "for key 'users.email'") {
					return 0, fmt.Errorf("%w: %s", models.ErrDuplicateEmail, email)
				}
			}
		}
		return 0, fmt.Errorf("ошибка сохранения пользователя: %w", err)
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// Активация ссылки
	_, err = tx.Exec("UPDATE qck_links SET user_id = ?, is_active = 1 WHERE id = ?", userID, linkID)
	if err != nil {
		return 0, err
	}

	// Фиксация транзакции
	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return userID, nil
}

// CheckUser - Проверка пользователя по email и паролю
func (r *UserRepo) CheckUser(login, password string) (int64, bool, error) {
	var userID int64
	var storedHash string
	var isActive []byte
	err := r.db.QueryRow("SELECT id, password_hash, is_active FROM users WHERE login = ?", login).Scan(&userID, &storedHash, &isActive)
	if err != nil {
		return 0, false, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password)); err != nil {
		return 0, false, err
	}

	return userID, isActive[0] == 1, nil
}

// ActivateUser - Активация пользователя
func (r *UserRepo) ActivateUser(userID int64) error {
	_, err := r.db.Exec("UPDATE users SET is_active = 1 WHERE id = ?", userID)
	return err
}

// CreateResetToken - Создание токена для сброса пароля
func (r *UserRepo) CreateResetToken(email string) (string, error) {
	println("func (r *UserRepo) CreateResetToken(email string)")
	var userID int64
	token, err := jwt.GenerateToken32()
	if err != nil {
		return "", err
	}

	err = r.db.QueryRow("SELECT id FROM users WHERE email = ?", email).Scan(&userID)
	println("юзер, err", userID, err)
	if err != nil {
		return "", err
	}

	_, err = r.db.Exec("INSERT INTO password_reset_tokens (user_id, token, expires_at) VALUES (?, ?, ?)", userID, token, time.Now().Add(30*time.Minute))
	if err != nil {
		return "", err
	}

	return token, nil
}
