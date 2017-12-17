package alert

import (
	"xzlan/dao"
	"xzlan/mail"
)

type Alert struct {
	MetricDao dao.Dao
	Mail mail.Mail
	EsUrl string
}

