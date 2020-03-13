package controller

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"log"
	"xzlan/alert"
	"xzlan/dao"
)

// APIController 接口对象
type APIController struct {
	Ctx iris.Context
	APIDao   *dao.APIDao
	RuleDao  *dao.RuleDao
	APIAlert *alert.Alert
}

// Get get /apis?name=&method=
func (a *APIController) Get() iris.Map {
	name := a.Ctx.URLParam("name")
	method := a.Ctx.URLParam("method")
	apis, err := a.APIDao.GetBy(name, method)
	if err != nil {
		return iris.Map{"code": 0, "msg": err.Error()}
	}
	return iris.Map{"code": 0, "msg": "", "count": len(apis), "data": apis}
}

// GetBy get /apis/{id}
func (a *APIController) GetBy(id string) mvc.View {
	v, err := a.APIDao.Get(id)
	if err != nil {
		log.Printf("apis/id error, %s", err)
	}
	return mvc.View{Name: "api/edit.html", Layout: iris.NoLayout, Data: v}
}

// GetAdd get /apis/add
func (a *APIController) GetAdd() mvc.View {
	return mvc.View{Name: "api/add.html", Layout: iris.NoLayout}
}

// PostAdd post /apis/add
func (a *APIController) PostAdd() iris.Map {
	api := dao.API{}
	err := a.Ctx.ReadJSON(&api)
	if err != nil {
		return iris.Map{"code": iris.StatusInternalServerError, "msg": err.Error()}
	}
	api.Status = "stop"
	err = a.APIDao.Add(api)
	if err != nil {
		return iris.Map{"code": iris.StatusInternalServerError, "msg": err.Error()}
	}
	return iris.Map{"code": iris.StatusOK, "msg": "OK"}
}

// PostEdit post /apis/edit
func (a *APIController) PostEdit() iris.Map {
	api := dao.API{}
	err := a.Ctx.ReadJSON(&api)
	if err != nil {
		return iris.Map{"code": iris.StatusInternalServerError, "msg": err.Error()}
	}
	err = a.APIDao.Update(api)
	if err != nil {
		return iris.Map{"code": iris.StatusInternalServerError, "msg": err.Error()}
	}
	return iris.Map{"code": iris.StatusOK, "msg": "OK"}
}

// DeleteBy delete /apis/{id}
func (a *APIController) DeleteBy(id string) iris.Map {
	a.APIAlert.Stop(id)
	err := a.APIDao.Delete(id)
	if err != nil {
		return iris.Map{"code": iris.StatusInternalServerError, "msg": err.Error()}
	}
	return iris.Map{"code": iris.StatusOK, "msg": "OK"}
}

// PostRunBy post /apis/run/{id}
func (a *APIController) PostRunBy(id string) iris.Map {
	api, err := a.APIDao.Get(id)
	if err != nil {
		return iris.Map{"code": iris.StatusInternalServerError, "msg": err.Error()}
	}
	//if api.Status == "running" {
	//	return iris.Map{"code": iris.StatusBadRequest, "msg": "该任务已启动无需处理"}
	//}
	rule, err := a.RuleDao.Get(api.ID)
	if err != nil {
		return iris.Map{"code": iris.StatusInternalServerError, "msg": err.Error()}
	}
	if rule.Type == "" {
		return iris.Map{"code": iris.StatusInternalServerError, "msg": "告警规则未配置"}
	}
	go a.APIAlert.RunJob(id)
	api.Status = "running"
	err = a.APIDao.Update(api)
	if err != nil {
		return iris.Map{"code": iris.StatusInternalServerError, "msg": err.Error()}
	}
	return iris.Map{"code": iris.StatusOK, "msg": "OK"}
}

// PostStopBy post /apis/stop/{id}
func (a *APIController) PostStopBy(id string) iris.Map {
	a.APIAlert.Stop(id)
	// 更新状态
	api, err := a.APIDao.Get(id)
	if err != nil {
		return iris.Map{"code": iris.StatusInternalServerError, "msg": err.Error()}
	}
	if api.Status == "stop" {
		return iris.Map{"code": iris.StatusBadRequest, "msg": "该任务已停止无需处理"}
	}
	api.Status = "stop"
	err = a.APIDao.Update(api)
	if err != nil {
		return iris.Map{"code": iris.StatusInternalServerError, "msg": err.Error()}
	}
	return iris.Map{"code": iris.StatusOK, "msg": "OK"}
}
