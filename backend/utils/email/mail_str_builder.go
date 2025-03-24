package email

import (
	"encoding/base64"
	"fmt"
	"strings"
)

func MailStrBuilder(from, to, subject, body string) string {
	nBytes := (len(to) + len(from) + len(subject) + len(body) + 200)

	var b strings.Builder
	b.Grow(nBytes)

	b.WriteString("From: ")
	b.WriteString(from)
	b.WriteString("\r\nTo: ")
	b.WriteString(to)
	b.WriteString("\r\nSubject: ")
	b.WriteString(subject)
	b.WriteString("\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=\"UTF-8\"\r\n\r\n")
	b.WriteString(body)
	b.WriteString("\nЕсли вы получили это письмо случайно - игнорируйте его")

	return b.String()
}

// EncodeSubject кодирует тему письма в Base64 (MIME)
func EncodeSubject(subject string) string {
	return fmt.Sprintf("=?utf-8?B?%s?=", base64.StdEncoding.EncodeToString([]byte(subject)))
}
