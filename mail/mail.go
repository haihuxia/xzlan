package mail

import (
	"net/smtp"
	"strings"
	"strconv"
	"bytes"
	"html/template"
)

type Mail struct {
	user string
	password string
	host string
	port int64
}

func NewMail(user, password, host string) Mail {
	return Mail{user, password, host, 25}
}

func (m *Mail) Send(to string) error {
	t, _ := template.ParseFiles("../static/template.html")
	var body bytes.Buffer
	body.Write([]byte("To: " + to + "\r\nFrom: 告警<" + m.user +
		">\r\nSubject: 告警\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n"))
	t.ExecuteTemplate(&body, "template.html", "aaaaa")

	sendTo := strings.Split(to, ";")
	auth := smtp.PlainAuth("", m.user, m.password, m.host)
	return smtp.SendMail(m.host + ":" + strconv.FormatInt(m.port, 10), auth, m.user, sendTo, body.Bytes())
}
