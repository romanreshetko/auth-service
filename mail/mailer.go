package mail

import (
	"fmt"
	"net/smtp"
	"os"
)

type Mailer struct {
	host string
	port string
	from string
}

func NewMailer(from string) *Mailer {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	return &Mailer{host: host, port: port, from: from}
}

func (m *Mailer) SendMail(to, subject, body string) error {
	user := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASS")
	addr := fmt.Sprintf("%s:%s", m.host, m.port)
	auth := smtp.PlainAuth("", user, password, m.host)
	msg := []byte("To: " + to + "\r\n" +
		"From: " + m.from + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"\r\n" +
		body)

	return smtp.SendMail(addr, auth, m.from, []string{to}, msg)
}

func SendVerificationEmail(mailer *Mailer, email, code string) error {
	body := fmt.Sprintf(
		"Код подтверждения адреса почты: %s\n",
		code,
	)

	return mailer.SendMail(email, "Подтверждение почты", body)
}

func SendTemporaryPasswordEmail(mailer *Mailer, email, password string) error {
	body := fmt.Sprintf(
		"Ваш временный пароль: %s\n Смените его после входа",
		password,
	)

	return mailer.SendMail(email, "Временный пароль", body)
}
