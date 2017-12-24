package controller

import (
	"xzlan/dao"
	"iris/mvc"
	"fmt"
	"iris"
	"encoding/json"
	"xzlan/alert"
)

type ApiController struct {
	mvc.C
	MetricDao dao.Dao
	ApiAlert alert.Alert
}

func (a *ApiController) Get() iris.Map {
	name := a.Ctx.URLParam("name")
	method := a.Ctx.URLParam("method")
	apis, err := a.MetricDao.GetApis(name, method)
	if err != nil {
		fmt.Printf("apis error, %s", err)
		return iris.Map{"code": 0, "msg": err.Error()}
	}
	return iris.Map{"code": 0, "msg": "", "count": len(apis), "data": apis}
}

func (a *ApiController) GetBy(id string) mvc.View {
	v, err := a.MetricDao.Get(dao.ApiTable, id)
	if err != nil {
		fmt.Printf("apis/id error, %s", err)
	}
	var api dao.Api
	err = json.Unmarshal(v, &api)
	if err != nil {
		fmt.Printf("apis/id error, %s", err)
	}
	return mvc.View{Name: "metric/editApi.html", Layout: iris.NoLayout, Data: api}
}

func (a *ApiController) GetAdd() mvc.View {
	return mvc.View{Name: "metric/addApi.html", Layout: iris.NoLayout}
}

func (a *ApiController) PostAdd() iris.Map {
	api := dao.Api{}
	err := a.Ctx.ReadJSON(&api)
	if err != nil {
		return iris.Map{"code": iris.StatusInternalServerError, "msg": err.Error()}
	}
	api.Status = "stop"
	err = a.MetricDao.PutApi(api)
	if err != nil {
		return iris.Map{"code": iris.StatusInternalServerError, "msg": err.Error()}
	}
	return iris.Map{"code": iris.StatusOK, "msg": "OK"}
}

func (a *ApiController) PostEdit() iris.Map {
	api := dao.Api{}
	err := a.Ctx.ReadJSON(&api)
	if err != nil {
		return iris.Map{"code": iris.StatusInternalServerError, "msg": err.Error()}
	}
	err = a.MetricDao.PutByStruct(dao.ApiTable, api.Id, api)
	if err != nil {
		return iris.Map{"code": iris.StatusInternalServerError, "msg": err.Error()}
	}
	return iris.Map{"code": iris.StatusOK, "msg": "OK"}
}

func (a *ApiController) DeleteBy(id string) iris.Map {
	err := a.MetricDao.Delete(dao.ApiTable, id)
	if err != nil {
		return iris.Map{"code": iris.StatusInternalServerError, "msg": err.Error()}
	}
	return iris.Map{"code": iris.StatusOK, "msg": "OK"}
}

func (a *ApiController) PostAlertBy(id string) iris.Map {
	api, err := a.MetricDao.GetApiBy(id)
	if err != nil {
		return iris.Map{"code": iris.StatusInternalServerError, "msg": err.Error()}
	}
	rule, err := a.MetricDao.GetRuleBy(api.Id)
	if err != nil {
		return iris.Map{"code": iris.StatusInternalServerError, "msg": err.Error()}
	}
	go a.ApiAlert.RunJob(api, rule)
	api.Status = "running"
	err = a.MetricDao.UpdateApi(api)
	if err != nil {
		return iris.Map{"code": iris.StatusInternalServerError, "msg": err.Error()}
	}
	return iris.Map{"code": iris.StatusOK, "msg": "OK"}
}
