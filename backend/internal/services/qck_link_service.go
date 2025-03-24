package services

import (
	"database/sql"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// QckLinkRepo - репозиторий для работы с qck_links
type QckLinkRepo struct {
	db *sql.DB
}

// NewQckLink - Создает новый репозиторий
func NewQckLink(db *sql.DB) *QckLinkRepo {
	return &QckLinkRepo{db: db}
}

func (r *QckLinkRepo) ActivateLink(userID int64, link, password string) error {
	var storedHash string

	// Получаем хеш пароля из БД
	err := r.db.QueryRow("SELECT password_hash FROM qck_links WHERE qck_link = ?", link).Scan(&storedHash)
	if err != nil {
		return fmt.Errorf("ошибка при получении пароля: %w", err)
	}

	// Сравниваем введенный пароль с хешем из БД
	if err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password)); err != nil {
		return fmt.Errorf("неверный пароль")
	}

	// Выполняем обновление данных
	res, err := r.db.Exec("UPDATE qck_links SET is_active=1, user_id = ? WHERE qck_link = ?", userID, link)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении данных: %w", err)
	}

	// Проверяем, были ли изменены записи
	affectedRows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка при проверке обновления: %w", err)
	}

	if affectedRows == 0 {
		return fmt.Errorf("ни одна запись не обновлена. Проверь данные запроса")
	}

	fmt.Printf("Обновлено %d записей\n", affectedRows)
	return nil
}

func (r *QckLinkRepo) UpdateDataLink(name, description, link string) error {
	if name == "" && description == "" {
		_, err := r.db.Exec("UPDATE qck_links SET description = ? WHERE qck_link = ?", description, link)
		return err
	}
	_, err := r.db.Exec("UPDATE qck_links SET name = ?, description = ? WHERE qck_link = ?", name, description, link)
	return err
}

func (r *QckLinkRepo) GetQckLinkID(link, password string) (int64, error) {
	var storedHash string
	err := r.db.QueryRow("SELECT password_hash FROM qck_links WHERE qck_link = ?", link).Scan(&storedHash)
	if err != nil {
		return 0, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password)); err != nil {
		return 0, err
	}

	var linkID int64
	err = r.db.QueryRow("SELECT id FROM qck_links WHERE qck_link = ?", link).Scan(&linkID)
	if err != nil {
		return 0, err
	}

	return linkID, nil
}

func (r *QckLinkRepo) GetQckLinkDescription(link string) (string, string, error) {
	var name, description string
	err := r.db.QueryRow("SELECT name, description FROM qck_links WHERE qck_link = ?", link).Scan(&name, &description)
	if err != nil {
		return "", "", err
	}
	return name, description, nil
}

func (r *QckLinkRepo) GetAllUserLinks(userID int64) ([]map[string]string, error) {
	rows, err := r.db.Query("SELECT qck_link, description, name FROM qck_links WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dataLinks []map[string]string
	for rows.Next() {
		var qckLink, description, name string
		if err := rows.Scan(&qckLink, &description, &name); err != nil {
			return nil, err
		}
		dataLinks = append(dataLinks, map[string]string{
			"qck_link":    qckLink,
			"description": description,
			"name":        name,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return dataLinks, nil
}
