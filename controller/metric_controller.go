package controller

import (
	"xzlan/dao"
	"iris/mvc"
	"fmt"
	"iris"
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

func (m *MetricController) DeleteBy(id string) iris.Map {
	err := m.MetricDao.Delete("api", id)
	if err != nil {
		return iris.Map{"code": iris.StatusInternalServerError, "msg": err.Error()}
	}
	return iris.Map{"code": iris.StatusOK, "msg": "OK"}
}
