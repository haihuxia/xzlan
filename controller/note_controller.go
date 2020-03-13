package controller

import (
	"github.com/kataras/iris/v12"
	"xzlan/dao"
)

// NoteController 通知记录
type NoteController struct {
	Ctx iris.Context
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