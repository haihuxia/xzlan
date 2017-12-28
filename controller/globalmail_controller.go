package controller

import (
	"github.com/kataras/iris/mvc"
	"xzlan/dao"
	"github.com/kataras/iris"
)

type GlobalMailController struct {
	mvc.C
	GlobalMailDao *dao.GlobalMailDao
}

// 查询
// get /globalmails
func (g *GlobalMailController) Get() iris.Map {
	var mails []dao.GlobalMail
	var err error
	mail := g.Ctx.URLParam("mail")
	if mail == "" {
		mails, err = g.GlobalMailDao.GetAll()
	} else {
		mails, err = g.GlobalMailDao.Get(mail)
	}
	if err != nil {
		return iris.Map{"code": 0, "msg": err.Error()}
	}
	return iris.Map{"code": 0, "msg": "", "count": len(mails), "data": mails}
}

// 新增
// get /globalmails/add
func (g *GlobalMailController) GetAdd() mvc.View {
	return mvc.View{Name: "globalmail/add.html", Layout: iris.NoLayout}
}

// 新增
// post /globalmails/add
func (g *GlobalMailController) PostAdd() iris.Map {
	m := &dao.GlobalMail{}
	err := g.Ctx.ReadJSON(&m)
	if err != nil {
		return iris.Map{"code": iris.StatusInternalServerError, "msg": err.Error()}
	}
	err = g.GlobalMailDao.Add(m.Mail)
	if err != nil {
		return iris.Map{"code": 0, "msg": err.Error()}
	}
	return iris.Map{"code": iris.StatusOK, "msg": "OK"}
}

// 删除
// delete /golbalmails/{mail}
func (g *GlobalMailController) DeleteBy(mail string) iris.Map {
	err := g.GlobalMailDao.Delete(mail)
	if err != nil {
		return iris.Map{"code": 0, "msg": err.Error()}
	}
	return iris.Map{"code": iris.StatusOK, "msg": "OK"}
}