package email

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"qckim-backend/internal/logger"
)

// EmailSender - структура для отправки писем
type EmailSender struct {
	SMTPHost  string
	SMTPPort  string
	Username  string
	Password  string
	FromEmail string
}

// NewEmailSender - создаёт новый экземпляр EmailSender
func NewEmailSender() *EmailSender {
	return &EmailSender{
		SMTPHost:  "smtp.yandex.ru",
		SMTPPort:  "587",
		Username:  "support@qck.im",
		Password:  "ljisglawpariypou",
		FromEmail: "support@qck.im",
	}
}

// SendEmail - отправляет письмо (синхронно)
func (e *EmailSender) SendEmail(toEmail, subject, body string) error {
	msg := MailStrBuilder(e.FromEmail, toEmail, subject, body)

	server := e.SMTPHost + ":" + e.SMTPPort

	client, err := smtp.Dial(server)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP: %w", err)
	}
	defer client.Close()

	if err = client.Hello("localhost"); err != nil {
		return fmt.Errorf("failed to send EHLO: %w", err)
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         e.SMTPHost,
	}
	if err = client.StartTLS(tlsConfig); err != nil {
		return fmt.Errorf("failed to start TLS: %w", err)
	}

	auth := smtp.PlainAuth("", e.Username, e.Password, e.SMTPHost)
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	if err = client.Mail(e.FromEmail); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	if err = client.Rcpt(toEmail); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	wc, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to send data: %w", err)
	}
	defer wc.Close()

	_, err = wc.Write([]byte(msg))
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	if err = client.Quit(); err != nil {
		logger.Debug("SMTP connection closed with warning", "error", err)
	}

	return nil
}

// EmailTask - структура задачи отправки письма
type EmailTask struct {
	To      string
	Subject string
	Body    string
}

// приватный канал для задач отправки писем
var emailQueue = make(chan EmailTask, 100)

// StartEmailWorker - запускает обработчик писем в отдельной горутине
func StartEmailWorker() {
	go func() {
		sender := NewEmailSender()
		for task := range emailQueue {
			err := sender.SendEmail(task.To, task.Subject, task.Body)
			if err != nil {
				logger.Error("Failed to send email", "error", err)
			} else {
				logger.Debug("Email sent", "to", task.To)
			}
		}
	}()
}

// SendEmailAsync - публичная функция для добавления задачи в очередь
func SendEmailAsync(to, subject, body string) {
	emailQueue <- EmailTask{To: to, Subject: subject, Body: body}
}
