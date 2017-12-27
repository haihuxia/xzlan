package alert

import (
	"xzlan/dao"
	"xzlan/mail"
	"gopkg.in/olivere/elastic.v5"
	"fmt"
	"context"
	"time"
	"strconv"
	"strings"
	"encoding/json"
)

var chanMap = make(map[string]chan bool)

type Alert struct {
	ApiDao  *dao.ApiDao
	RuleDao *dao.RuleDao
	NoteDao *dao.NoteDao
	Mail    *mail.Mail
	EsUrl   string
}

type Message struct {
	Message string `json:"message"`
}

func NewAlert(apiDao *dao.ApiDao, ruleDao *dao.RuleDao, noteDao *dao.NoteDao, mail *mail.Mail, esUrl string) *Alert {
	return &Alert{apiDao, ruleDao, noteDao, mail, esUrl}
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
		a.RunJob(apis[i], rule)
	}
	return nil
}

func (a *Alert) Stop(id string) error {
	if v, ok := chanMap[id]; ok {
		v <- true
	}
	return nil
}

func (a *Alert) RunJob(api dao.Api, rule dao.Rule) {
	tick := time.Tick(60e9)
	stop := make(chan bool)
	chanMap[api.Id] = stop
	var flag = false
	for {
		select {
		case <-tick:
			a.job(api, rule)
		case <-stop:
			flag = true
		}
		if flag {
			break
		}
	}
	fmt.Println("Stop !")
}

func (a *Alert) job(api dao.Api, rule dao.Rule) {
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
	result, err := client.Search("logstash-" + t).Query(query).Do(context.Background())
	if err != nil {
		fmt.Printf("elk query error %s \n", err)
		return
	}
	c, err := strconv.ParseInt(rule.Count, 10, 64)
	if err != nil {
		fmt.Printf("parseInt error %s \n", err)
		return
	}
	if result.Hits.TotalHits >= c {
		// 添加告警记录
		note := "接口：" + api.Name + "." + api.Method + "，耗时检查：实际 " +
			strconv.FormatInt(result.Hits.TotalHits, 10) + " 次 >= 限制 " + rule.Count + " 次"
		notifyTime, err := a.NoteDao.Add(note, api.Id)
		if err != nil {
			fmt.Printf("add note error %s \n", err)
		}
		fmt.Printf("%s %s %d hit: %d 【命中】 \n", api.Name, api.Method, c, result.Hits.TotalHits)

		// 判断是否需要发送邮件
		if rule.Delay != "" && api.NotifyTime != "" {
			notifyTime, err := time.Parse("2006-01-02 15:04:05", api.NotifyTime)
			if err != nil {
				fmt.Printf("time.Parse error %s \n", err)
			}
			delayTime, err := dao.DelayToTime(rule.Delay, notifyTime)
			if err != nil {
				fmt.Printf("time.Parse error %s \n", err)
			}
			if delayTime.After(time.Now()) {
				a.NoteDao.Add(note + ", 下次通知时间：" + delayTime.Format("2006-01-02 15:04:05"), api.Id)
				return
			}
		}

		// 发送邮件
		body := "接口：<b>" + api.Name + "</b><br/>方法：<b>" + api.Method + "</b><br/>耗时匹配次数：<b>" +
			strconv.FormatInt(result.Hits.TotalHits, 10) + "</b>（告警规则：大于 " + rule.Min + "ms, " +
			rule.Count + "次）<br/><br/>原始日志："
		for i := 0; i < len(result.Hits.Hits); i++ {
			b, _ := result.Hits.Hits[i].Source.MarshalJSON()
			var m Message
			json.Unmarshal(b, &m)
			body = body + "<br/>" + m.Message
		}
		tos := strings.Split(rule.Mails, ";")
		for i := 0; i < len(tos); i++ {
			if tos[i] == "" {
				continue
			}
			err = a.Mail.Send(tos[i], body)
			if err == nil {
				fmt.Printf("邮件发送成功 %s \n", rule.Mails)
			} else {
				fmt.Printf("error: %s \n", err)
			}
		}
		// 修改状态，告警已通知
		a.NoteDao.Update(note, api.Id, "0")
		api.NotifyTime = notifyTime
		a.ApiDao.Update(api)
	} else {
		fmt.Printf("%s %s %d hit: %d 【未命中】 \n", api.Name, api.Method, c, result.Hits.TotalHits)
	}
}
