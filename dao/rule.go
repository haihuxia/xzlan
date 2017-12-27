package dao

import (
	"encoding/json"
	"time"
	"strconv"
)

const RuleTable = "rule"

// 告警规则
// Api.Id 与 rule 表的 key 一致
// 以 Api.Id 为 key
type Rule struct {
	Type string `json:"type"`
	Max string `json:"max"`
	Min string `json:"min"`
	Val string `json:"val"`
	Time string `json:"time"`
	Count string `json:"count"`
	Delay string `json:"delay"`
	Mails string `json:"mails"`
}

type RuleDao struct {
	dao *Dao
}

func NewRuleDao(dao *Dao) *RuleDao {
	return &RuleDao{dao}
}

func (r *RuleDao) Get(id string) (rule Rule, err error) {
	v, err := r.dao.Get(RuleTable, id)
	if err != nil {
		return rule, err
	}
	err = json.Unmarshal(v, &rule)
	return
}

func (r *RuleDao) Add(id string, rule Rule) error {
	return r.dao.PutByStruct(RuleTable, id, rule)
}

func DelayToTime(delay string, t time.Time) (time.Time, error) {
	if delay == "" {
		return t, nil
	}
	i := len(delay) - 1
	last := string(delay[i])
	num, err := strconv.ParseInt(string(delay[:i]), 10, 64)
	if err != nil {
		return t, err
	}
	switch last {
	case "d":
		// 天
		t = t.AddDate(0, 0, int(num))
	case "m":
		// 分钟
		t = t.Add(time.Duration(num) * time.Minute)
	case "h":
		// 小时
		t = t.Add(time.Duration(num) * time.Hour)
	}
	return t, nil
}