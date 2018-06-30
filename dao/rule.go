package dao

import (
	"encoding/json"
	"time"
	"strconv"
)

// RuleTable 表名
const RuleTable = "rule"

// Rule 告警规则
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

// RuleDao 数据操作
type RuleDao struct {
	dao *Dao
}

// NewRuleDao 构造函数
func NewRuleDao(dao *Dao) *RuleDao {
	return &RuleDao{dao}
}

// Get 查询
func (r *RuleDao) Get(id string) (rule Rule, err error) {
	v, err := r.dao.Get(RuleTable, id)
	if err != nil {
		return rule, err
	}
	err = json.Unmarshal(v, &rule)
	return
}

// Add 新增
func (r *RuleDao) Add(id string, rule Rule) error {
	return r.dao.PutByStruct(RuleTable, id, rule)
}

// DelayToTime 计算延时
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