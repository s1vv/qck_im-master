package email

import (
	"os"
	"path/filepath"
)

func saveEmailToFile(toEmail, msg string) error {
	dir := "test_emails"

	// Создаем директорию, если она не существует
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	// Записываем сообщение в файл
	filePath := filepath.Join(dir, toEmail+".txt")
	return os.WriteFile(filePath, []byte(msg), 0644)
}
