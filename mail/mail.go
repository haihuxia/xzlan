package mail

import (
	"net/smtp"
	"strings"
	"strconv"
	"bytes"
	"html/template"
	"fmt"
)

type Mail struct {
	user string
	password string
	host string
	port int64
	htmlTplUrl string
}

func NewMail(user, password, host, htmlTplUrl string) *Mail {
	return &Mail{user, password, host, 25, htmlTplUrl}
}

func (m *Mail) Send(to string, content string) error {
	if m.htmlTplUrl != "" {
		return m.html(to, content)
	}
	return m.text(to, content)
}

func (m *Mail) text(to string, content string) error {
	var body bytes.Buffer
	body.Write([]byte("To: " + to + "\r\nFrom: 告警<" + m.user +
		">\r\nSubject: 【告警】接口耗时超限\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n" + content))
	sendTo := strings.Split(to, ";")
	auth := smtp.PlainAuth("", m.user, m.password, m.host)
	return smtp.SendMail(m.host + ":" + strconv.FormatInt(m.port, 10), auth, m.user, sendTo, body.Bytes())
}

func (m *Mail) html(to string, content string) error {
	t, err := template.ParseFiles(m.htmlTplUrl)
	if err != nil {
		fmt.Println(err)
	}
	var body bytes.Buffer
	body.Write([]byte("To: " + to + "\r\nFrom: 告警<" + m.user +
		">\r\nSubject: 【告警】接口耗时超限\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n"))
	err = t.ExecuteTemplate(&body, t.Name(), template.HTML(content))
	if err != nil {
		fmt.Println(err)
	}
	sendTo := strings.Split(to, ";")
	auth := smtp.PlainAuth("", m.user, m.password, m.host)
	return smtp.SendMail(m.host + ":" + strconv.FormatInt(m.port, 10), auth, m.user, sendTo, body.Bytes())
}