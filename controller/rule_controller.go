package controller

import (
	"xzlan/dao"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris"
	"log"
)

type RuleController struct {
	mvc.C
	RuleDao *dao.RuleDao
	ApiDao  *dao.ApiDao
}

type Rule struct {
	Id string
	dao.Rule
}

// 查询
// get /rule/{id}
func (r *RuleController) GetBy(id string) mvc.View {
	v, err := r.RuleDao.Get(id)
	if err != nil {
		log.Printf("rule/id error, %s", err)
	}
	return mvc.View{Name: "rule/rule.html", Layout: iris.NoLayout, Data: iris.Map{
		"Id":   id,
		"Rule": v,
	}}
}

// 新增
// post /rule/add
func (r *RuleController) PostAdd() iris.Map {
	rule := Rule{}
	err := r.Ctx.ReadJSON(&rule)
	if err != nil {
		return iris.Map{"code": iris.StatusInternalServerError, "msg": err.Error()}
	}
	var daoRule = dao.Rule{rule.Type, rule.Max, rule.Min, rule.Val, rule.Time,
		rule.Count, rule.Delay, rule.Mails}
	err = r.RuleDao.Add(rule.Id, daoRule)
	if err != nil {
		return iris.Map{"code": iris.StatusInternalServerError, "msg": err.Error()}
	}
	api, err := r.ApiDao.Get(rule.Id)
	if err != nil {
		return iris.Map{"code": iris.StatusInternalServerError, "msg": err.Error()}
	}
	api.NotifyTime = ""
	err = r.ApiDao.Update(api)
	if err != nil {
		return iris.Map{"code": iris.StatusInternalServerError, "msg": err.Error()}
	}
	return iris.Map{"code": iris.StatusOK, "msg": "OK"}
}
