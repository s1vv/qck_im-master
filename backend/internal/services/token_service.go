package services

import (
	"database/sql"
	"time"
)

type TokenRepo struct {
	db *sql.DB
}

func NewToken(db *sql.DB) *TokenRepo {
	return &TokenRepo{db: db}
}

func (r *TokenRepo) SaveToken(userID int64, ip, userAgent, token string) (err error) {
	defer func() {
		if err != nil {
			token = ""
		}
	}()

	expiresAt := time.Now().Add(1 * time.Hour)

	query := "SELECT token FROM refresh_tokens WHERE user_id = ?"
	row := r.db.QueryRow(query, userID)

	var existingToken string
	err = row.Scan(&existingToken)
	if err == sql.ErrNoRows {
		query = "INSERT INTO refresh_tokens (user_id, token, expires_at, ip_address, user_agent) VALUES (?, ?, ?, ?, ?)"
		_, err = r.db.Exec(query, userID, token, expiresAt, ip, userAgent)
		return err
	} else if err == nil {
		query = "UPDATE refresh_tokens SET token = ?, expires_at = ?, ip_address = ?, user_agent = ? WHERE user_id = ?"
		_, err = r.db.Exec(query, token, expiresAt, ip, userAgent, userID)
		return err
	} else {
		return err
	}
}

func (r *TokenRepo) InvalidateUserTokens(userID int64) error {
	query := "DELETE FROM refresh_tokens WHERE user_id = ?"
	_, err := r.db.Exec(query, userID)
	return err
}

func (r *TokenRepo) CheckResetToken(token string) (userID int64, expiresAt time.Time, err error) {
	var expiresAtStr string
	query := "SELECT user_id, expires_at FROM password_reset_tokens WHERE token = ?"
	err = r.db.QueryRow(query, token).Scan(&userID, &expiresAtStr)
	if err != nil {
		return
	}
	expiresAt, err = time.Parse("2006-01-02 15:04:05", expiresAtStr)
	return
}
