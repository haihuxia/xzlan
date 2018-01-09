package alert

import (
	"xzlan/dao"
	"xzlan/mail"
	"gopkg.in/olivere/elastic.v5"
	"context"
	"time"
	"strconv"
	"strings"
	"encoding/json"
	"log"
	"errors"
)

var chanMap = make(map[string]chan bool)

type Alert struct {
	ApiDao        *dao.ApiDao
	RuleDao       *dao.RuleDao
	NoteDao       *dao.NoteDao
	GlobalMailDao *dao.GlobalMailDao
	Mail          *mail.Mail
	EsUrl         string
	EsIndex       string
}

type Message struct {
	Message string `json:"message"`
}

func NewAlert(apiDao *dao.ApiDao, ruleDao *dao.RuleDao, noteDao *dao.NoteDao, globalMailDao *dao.GlobalMailDao,
	mail *mail.Mail, esUrl string, esIndex string) *Alert {
	return &Alert{apiDao, ruleDao, noteDao, globalMailDao, mail,
		esUrl, esIndex}
}

func (a *Alert) Start() error {
	apis, err := a.ApiDao.GetAll()
	if err != nil {
		return err
	}
	for i := range apis {
		rule, err := a.RuleDao.Get(apis[i].Id)
		if err != nil {
			return err
		}
		if rule.Type == "" {
			return errors.New("no config alert rule")
		}
		a.RunJob(apis[i].Id)
	}
	return nil
}

func (a *Alert) Stop(id string) error {
	if v, ok := chanMap[id]; ok {
		v <- true
	}
	return nil
}

func (a *Alert) RunJob(id string) {
	tick := time.Tick(60e9)
	stop := make(chan bool)
	chanMap[id] = stop
	var flag = false
	for {
		select {
		case <-tick:
			// 不能传对象，否则无法加载最新修改的值
			a.job(id)
		case <-stop:
			flag = true
		}
		if flag {
			break
		}
	}
}

func (a *Alert) job(id string) {
	api, err := a.ApiDao.Get(id)
	if err != nil {
		log.Printf("job ApiDao.Get error %s \n", err)
		return
	}
	rule, err := a.RuleDao.Get(id)
	if err != nil {
		log.Printf("job RuleDao.Get error %s \n", err)
		return
	}
	if rule.Type != "min" {
		return
	}
	client, err := elastic.NewClient(elastic.SetURL(a.EsUrl))
	query := elastic.NewBoolQuery()
	query = query.Must(elastic.NewMatchPhraseQuery("interface", api.Name))
	query = query.Must(elastic.NewMatchPhraseQuery("method", api.Method))
	if rule.Type == "min" {
		query = query.Filter(elastic.NewRangeQuery("elapsed").Gte(rule.Min))
		query = query.Filter(elastic.NewRangeQuery("@timestamp").Gte("now-" + rule.Time + "m").Lt("now"))
	}
	t := time.Now().Format("2006.01.02")
	result, err := client.Search(a.EsIndex + t).Query(query).Do(context.Background())
	if err != nil {
		log.Printf("elk query error %s \n", err)
		return
	}
	c, err := strconv.ParseInt(rule.Count, 10, 64)
	if err != nil {
		log.Printf("parseInt error %s \n", err)
		return
	}
	if result.Hits.TotalHits >= c {
		// 添加告警记录
		note := "接口：" + api.Name + "." + api.Method + "，耗时检查：实际 " +
			strconv.FormatInt(result.Hits.TotalHits, 10) + " 次 >= 限制 " + rule.Count + " 次"
		notifyTime, err := a.NoteDao.Add(note, api.Id)
		if err != nil {
			log.Printf("add note error %s \n", err)
		}
		log.Printf("%s %s %d hit: %d 【命中】 \n", api.Name, api.Method, c, result.Hits.TotalHits)
		// 判断是否需要发送邮件
		if rule.Delay != "" && api.NotifyTime != "" {
			notifyTime, err := time.Parse("2006-01-02 15:04:05", api.NotifyTime)
			if err != nil {
				log.Printf("time.Parse error %s \n", err)
			}
			delayTime, err := dao.DelayToTime(rule.Delay, notifyTime)
			if err != nil {
				log.Printf("time.Parse error %s \n", err)
			}
			if delayTime.After(time.Now()) {
				a.NoteDao.Add(note+", 下次通知时间：" + delayTime.Format("2006-01-02 15:04:05"), api.Id)
				return
			}
		}

		// 发送邮件
		body := "接口：<b>" + api.Name + "</b><br/>方法：<b>" + api.Method + "</b><br/>耗时匹配次数：<b>" +
			strconv.FormatInt(result.Hits.TotalHits, 10) + "</b>（告警规则：大于 " + rule.Min + "ms, " +
			rule.Count + "次）<br/><br/>截取日志片段："
		// 最多截取 5 条原始日志
		l := len(result.Hits.Hits)
		if l > 5 {
			l = 5
		}
		for i := 0; i < l; i++ {
			b, _ := result.Hits.Hits[i].Source.MarshalJSON()
			var m Message
			json.Unmarshal(b, &m)
			body = body + "<br/>" + m.Message + "<br/>"
		}
		tos := strings.Split(rule.Mails, ";")
		for i := 0; i < len(tos); i++ {
			if tos[i] == "" {
				continue
			}
			err = a.Mail.Send(tos[i], body)
			if err == nil {
				log.Printf("email send [suc] %s \n", tos[i])
			} else {
				log.Printf("email send [falil] %s err: %s \n", tos[i], err)
			}
		}
		globalMails, err := a.GlobalMailDao.GetAll()
		if err != nil {
			log.Printf("alert GlobalMailDao.GetAll error %s \n", err)
		}
		for i := 0; i < len(globalMails); i++ {
			err = a.Mail.Send(globalMails[i].Mail, body)
			if err == nil {
				log.Printf("email send [suc] %s \n", globalMails[i].Mail)
			} else {
				log.Printf("email send [falil] %s err: %s \n", globalMails[i].Mail, err)
			}
		}
		// 修改状态，告警已通知
		a.NoteDao.Update(note, api.Id, "0")
		api.NotifyTime = notifyTime
		a.ApiDao.Update(api)
	} else {
		log.Printf("%s %s %d hit: %d 【未命中】 \n", api.Name, api.Method, c, result.Hits.TotalHits)
	}
}
