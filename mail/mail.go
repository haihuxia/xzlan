package mail

import (
	"net/smtp"
	"strings"
	"strconv"
	"bytes"
	"html/template"
	"log"
)

// Mail 邮箱
type Mail struct {
	user string
	password string
	host string
	port int64
	htmlTplURL string
}

// NewMail 构造函数
func NewMail(user, password, host, htmlTplURL string) *Mail {
	return &Mail{user, password, host, 25, htmlTplURL}
}

// Send 发送邮件
func (m *Mail) Send(to string, content string) error {
	if m.htmlTplURL != "" {
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
	t, err := template.ParseFiles(m.htmlTplURL)
	if err != nil {
		log.Printf("template.ParseFiles error %s \n", err)
	}
	var body bytes.Buffer
	body.Write([]byte("To: " + to + "\r\nFrom: 告警<" + m.user +
		">\r\nSubject: 【告警】接口耗时超限\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n"))
	err = t.ExecuteTemplate(&body, t.Name(), template.HTML(content))
	if err != nil {
		log.Printf("ExecuteTemplate error %s \n", err)
	}
	sendTo := strings.Split(to, ";")
	auth := smtp.PlainAuth("", m.user, m.password, m.host)
	return smtp.SendMail(m.host + ":" + strconv.FormatInt(m.port, 10), auth, m.user, sendTo, body.Bytes())
}