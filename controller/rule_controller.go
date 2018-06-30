package controller

import (
	"xzlan/dao"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris"
	"log"
)

// RuleController 规则
type RuleController struct {
	mvc.C
	RuleDao *dao.RuleDao
	APIDao  *dao.APIDao
}

// Rule 规则
type Rule struct {
	ID string
	dao.Rule
}

// GetBy get /rule/{id}
func (r *RuleController) GetBy(id string) mvc.View {
	v, err := r.RuleDao.Get(id)
	if err != nil {
		log.Printf("rule/id error, %s", err)
	}
	return mvc.View{Name: "rule/rule.html", Layout: iris.NoLayout, Data: iris.Map{
		"ID":   id,
		"Rule": v,
	}}
}

// PostAdd post /rule/add
func (r *RuleController) PostAdd() iris.Map {
	rule := Rule{}
	err := r.Ctx.ReadJSON(&rule)
	if err != nil {
		return iris.Map{"code": iris.StatusInternalServerError, "msg": err.Error()}
	}
	var daoRule = dao.Rule{rule.Type, rule.Max, rule.Min, rule.Val, rule.Time,
		rule.Count, rule.Delay, rule.Mails}
	err = r.RuleDao.Add(rule.ID, daoRule)
	if err != nil {
		return iris.Map{"code": iris.StatusInternalServerError, "msg": err.Error()}
	}
	api, err := r.APIDao.Get(rule.ID)
	if err != nil {
		return iris.Map{"code": iris.StatusInternalServerError, "msg": err.Error()}
	}
	api.NotifyTime = ""
	err = r.APIDao.Update(api)
	if err != nil {
		return iris.Map{"code": iris.StatusInternalServerError, "msg": err.Error()}
	}
	return iris.Map{"code": iris.StatusOK, "msg": "OK"}
}
