package controller

import (
	"xzlan/dao"
	"iris/mvc"
	"iris"
	"fmt"
	"encoding/json"
)

type RuleController struct {
	mvc.C
	MetricDao dao.Dao
}

type Rule struct {
	Id string
	dao.Rule
}

func (r *RuleController) GetBy(id string) mvc.View {
	v, err := r.MetricDao.Get(dao.RuleTable, id)
	if err != nil {
		fmt.Printf("rule/id error, %s", err)
	}
	if len(v) == 0 {
		return mvc.View{Name: "metric/rule.html", Layout: iris.NoLayout, Data: iris.Map{
			"Id":   id,
			"Rule": dao.Rule{},
		}}
	}
	var rule dao.Rule
	err = json.Unmarshal(v, &rule)
	if err != nil {
		fmt.Printf("rule/id error, %s", err)
	}
	fmt.Print(rule)
	return mvc.View{Name: "metric/rule.html", Layout: iris.NoLayout, Data: iris.Map{
		"Id":   id,
		"Rule": rule,
	}}
}

func (r *RuleController) PostAdd() iris.Map {
	rule := Rule{}
	err := r.Ctx.ReadJSON(&rule)
	if err != nil {
		return iris.Map{"code": iris.StatusInternalServerError, "msg": err.Error()}
	}
	var daoRule = dao.Rule{rule.Type, rule.Max, rule.Min, rule.Val, rule.Time,
		rule.Count, rule.Mails}
	err = r.MetricDao.PutByStruct(dao.RuleTable, rule.Id, daoRule)
	if err != nil {
		return iris.Map{"code": iris.StatusInternalServerError, "msg": err.Error()}
	}
	return iris.Map{"code": iris.StatusOK, "msg": "OK"}
}
