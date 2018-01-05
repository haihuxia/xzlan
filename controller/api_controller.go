package controller

import (
	"xzlan/dao"
	"fmt"
	"xzlan/alert"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
)

type ApiController struct {
	mvc.C
	ApiDao   *dao.ApiDao
	RuleDao  *dao.RuleDao
	ApiAlert *alert.Alert
}

// 查询
// get /apis?name=&method=
func (a *ApiController) Get() iris.Map {
	name := a.Ctx.URLParam("name")
	method := a.Ctx.URLParam("method")
	apis, err := a.ApiDao.GetBy(name, method)
	if err != nil {
		return iris.Map{"code": 0, "msg": err.Error()}
	}
	return iris.Map{"code": 0, "msg": "", "count": len(apis), "data": apis}
}

// 查询
// get /apis/{id}
func (a *ApiController) GetBy(id string) mvc.View {
	v, err := a.ApiDao.Get(id)
	if err != nil {
		fmt.Printf("apis/id error, %s", err)
	}
	return mvc.View{Name: "api/edit.html", Layout: iris.NoLayout, Data: v}
}

// 新增
// get /apis/add
func (a *ApiController) GetAdd() mvc.View {
	return mvc.View{Name: "api/add.html", Layout: iris.NoLayout}
}

// 新增
// post /apis/add
func (a *ApiController) PostAdd() iris.Map {
	api := dao.Api{}
	err := a.Ctx.ReadJSON(&api)
	if err != nil {
		return iris.Map{"code": iris.StatusInternalServerError, "msg": err.Error()}
	}
	api.Status = "stop"
	err = a.ApiDao.Add(api)
	if err != nil {
		return iris.Map{"code": iris.StatusInternalServerError, "msg": err.Error()}
	}
	return iris.Map{"code": iris.StatusOK, "msg": "OK"}
}

// 修改
// post /apis/edit
func (a *ApiController) PostEdit() iris.Map {
	api := dao.Api{}
	err := a.Ctx.ReadJSON(&api)
	if err != nil {
		return iris.Map{"code": iris.StatusInternalServerError, "msg": err.Error()}
	}
	err = a.ApiDao.Update(api)
	if err != nil {
		return iris.Map{"code": iris.StatusInternalServerError, "msg": err.Error()}
	}
	return iris.Map{"code": iris.StatusOK, "msg": "OK"}
}

// 删除
// delete /apis/{id}
func (a *ApiController) DeleteBy(id string) iris.Map {
	err := a.ApiDao.Delete(id)
	if err != nil {
		return iris.Map{"code": iris.StatusInternalServerError, "msg": err.Error()}
	}
	return iris.Map{"code": iris.StatusOK, "msg": "OK"}
}

// 允许
// post /apis/run/{id}
func (a *ApiController) PostRunBy(id string) iris.Map {
	api, err := a.ApiDao.Get(id)
	if err != nil {
		return iris.Map{"code": iris.StatusInternalServerError, "msg": err.Error()}
	}
	//if api.Status == "running" {
	//	return iris.Map{"code": iris.StatusBadRequest, "msg": "该任务已启动无需处理"}
	//}

	rule, err := a.RuleDao.Get(api.Id)
	if err != nil {
		return iris.Map{"code": iris.StatusInternalServerError, "msg": err.Error()}
	}
	if rule.Type == "" {
		return iris.Map{"code": iris.StatusInternalServerError, "msg": "告警规则未配置"}
	}
	go a.ApiAlert.RunJob(api, rule)
	api.Status = "running"
	err = a.ApiDao.Update(api)
	if err != nil {
		return iris.Map{"code": iris.StatusInternalServerError, "msg": err.Error()}
	}
	return iris.Map{"code": iris.StatusOK, "msg": "OK"}
}

// 停止
// post /apis/stop/{id}
func (a *ApiController) PostStopBy(id string) iris.Map {
	a.ApiAlert.Stop(id)
	// 更新状态
	api, err := a.ApiDao.Get(id)
	if err != nil {
		return iris.Map{"code": iris.StatusInternalServerError, "msg": err.Error()}
	}
	if api.Status == "stop" {
		return iris.Map{"code": iris.StatusBadRequest, "msg": "该任务已停止无需处理"}
	}
	api.Status = "stop"
	err = a.ApiDao.Update(api)
	if err != nil {
		return iris.Map{"code": iris.StatusInternalServerError, "msg": err.Error()}
	}
	return iris.Map{"code": iris.StatusOK, "msg": "OK"}
}
