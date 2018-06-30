package controller

import (
	"xzlan/dao"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris"
)

// NoteController 通知记录
type NoteController struct {
	mvc.C
	NoteDao *dao.NoteDao
}

// Get get /notes/
func (n *NoteController) Get() iris.Map {
	notes, err := n.NoteDao.GetAll()
	if err != nil {
		return iris.Map{"code": 0, "msg": err.Error()}
	}
	return iris.Map{"code": 0, "msg": "", "count": len(notes), "data": notes}
}