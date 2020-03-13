package dao

import (
	"fmt"
	"testing"
)

func TestNoteDao_Add(t *testing.T) {
	dao, err := NewDao("/Users/tiger/project/logs/go/xzlan.db")
	defer dao.Db.Close()

	if err != nil {
		t.Fatal(err)
	}
	noteDao := NewNoteDao(dao)
	_, err = noteDao.Add("接口：user.add，耗时检查：实际 5 次 >= 限制 5 次", "2")
	if err != nil {
		t.Fatal(err)
	}
}

func TestNoteDao_GetAll(t *testing.T) {
	dao, err := NewDao("/Users/tiger/project/logs/go/xzlan.db")
	defer dao.Db.Close()

	if err != nil {
		t.Fatal(err)
	}
	noteDao := NewNoteDao(dao)
	notes, err := noteDao.GetAll()
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < len(notes); i++ {
		fmt.Printf("value[%d]: %s \n", i, notes[i])
	}
}
