package controller

import (
	"iris/mvc"
	"xzlan/dao"
	"iris"
)

type NoteController struct {
	mvc.C
	NoteDao *dao.NoteDao
}

// 查询
// get /notes/
func (n *NoteController) Get() iris.Map {
	notes, err := n.NoteDao.GetAll()
	if err != nil {
		return iris.Map{"code": 0, "msg": err.Error()}
	}
	return iris.Map{"code": 0, "msg": "", "count": len(notes), "data": notes}
}