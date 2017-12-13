package controller

import (
	"xzlan/dao"
	"iris/mvc"
	"fmt"
	"iris"
	"encoding/json"
)

type MetricController struct {
	mvc.C
	MetricDao dao.Dao
}

func (m *MetricController) Get() iris.Map {
	name := m.Ctx.URLParam("name")
	method := m.Ctx.URLParam("method")
	a, err := m.MetricDao.GetApis(name, method)
	if err != nil {
		fmt.Printf("apis error, %s", err)
		return iris.Map{"code": 0, "msg": err.Error()}
	}
	return iris.Map{"code": 0, "msg": "", "count": len(a), "data": a}
}

func (m *MetricController) GetBy(id string) mvc.View {
	v, err := m.MetricDao.Get(dao.ApiTable, id)
	if err != nil {
		fmt.Printf("apis/id error, %s", err)
	}
	var a dao.Api
	err = json.Unmarshal(v, &a)
	if err != nil {
		fmt.Printf("apis/id error, %s", err)
	}
	return mvc.View{Name: "metric/editApi.html", Layout: iris.NoLayout, Data: a}
}

func (m *MetricController) GetAdd() mvc.View {
	return mvc.View{Name: "metric/addApi.html", Layout: iris.NoLayout}
}

func (m *MetricController) PostAdd() iris.Map {
	api := dao.Api{}
	err := m.Ctx.ReadJSON(&api)
	if err != nil {
		return iris.Map{"code": iris.StatusInternalServerError, "msg": err.Error()}
	}
	err = m.MetricDao.PutApi(api)
	if err != nil {
		return iris.Map{"code": iris.StatusInternalServerError, "msg": err.Error()}
	}
	return iris.Map{"code": iris.StatusOK, "msg": "OK"}
}

func (m *MetricController) PostEdit() iris.Map {
	api := dao.Api{}
	err := m.Ctx.ReadJSON(&api)
	if err != nil {
		return iris.Map{"code": iris.StatusInternalServerError, "msg": err.Error()}
	}
	err = m.MetricDao.Update(dao.ApiTable, api.Id, api)
	if err != nil {
		return iris.Map{"code": iris.StatusInternalServerError, "msg": err.Error()}
	}
	return iris.Map{"code": iris.StatusOK, "msg": "OK"}
}

func (m *MetricController) DeleteBy(id string) iris.Map {
	err := m.MetricDao.Delete(dao.ApiTable, id)
	if err != nil {
		return iris.Map{"code": iris.StatusInternalServerError, "msg": err.Error()}
	}
	return iris.Map{"code": iris.StatusOK, "msg": "OK"}
}
