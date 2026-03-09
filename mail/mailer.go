package mail

import (
	"fmt"
	"net/smtp"
)

type Mailer struct {
	host string
	port int
	from string
}

func NewMailer(host string, port int, from string) *Mailer {
	return &Mailer{host: host, port: port, from: from}
}

func (m *Mailer) SendMail(to, subject, body string) error {
	addr := fmt.Sprintf("%s:%d", m.host, m.port)

	msg := []byte("To: " + to + "\r\nSubject: " + subject + "\r\n\r\n" + body)

	return smtp.SendMail(addr, nil, m.from, []string{to}, msg)
}

func SendVerificationEmail(email, code string) error {
	mailer := NewMailer("mailhog", 1025, "noreply@cityviewpoint.ru")

	body := fmt.Sprintf(
		"Your verification code is %s\n",
		code,
	)

	return mailer.SendMail(email, "Email verification", body)
}
