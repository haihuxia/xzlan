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
	Mail    *mail.Mail
	EsUrl   string
}

type Message struct {
	Message string `json:"message"`
}

func NewAlert(apiDao *dao.ApiDao, ruleDao *dao.RuleDao, mail *mail.Mail, esUrl string) *Alert {
	return &Alert{apiDao, ruleDao, mail, esUrl}
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
		fmt.Printf("error %s \n", err)
	}
	c, err := strconv.ParseInt(rule.Count, 10, 64)
	if err != nil {
		fmt.Printf("error %s \n", err)
	}
	if result.Hits.TotalHits >= c {
		fmt.Printf("%s %s %d hit: %d 【命中】 \n", api.Name, api.Method, c, result.Hits.TotalHits)
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
		fmt.Printf("tos: %s", tos)
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
	} else {
		fmt.Printf("%s %s %d hit: %d 【未命中】 \n", api.Name, api.Method, c, result.Hits.TotalHits)
	}
}
