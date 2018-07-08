package dao

import (
	"testing"
)

func TestDao_DeleteTable(t *testing.T) {
	dao, err := NewDao("/Users/tiger/project/logs/go/xzlan.db")
	defer dao.Db.Close()

	if err != nil {
		t.Fatal(err)
	}
	err = dao.DeleteTable(NoteTable)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDao_PutByStruct(t *testing.T) {
	dao, err := NewDao("/Users/tiger/project/logs/go/xzlan.db")
	defer dao.Db.Close()

	if err != nil {
		t.Fatal(err)
	}
	api := API{"1", "user", "get", "查询接口", "stop", ""}
	err = dao.PutByStruct(APITable, "1", api)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDao_Delete(t *testing.T) {
	dao, err := NewDao("/Users/tiger/project/logs/go/xzlan.db")
	defer dao.Db.Close()

	if err != nil {
		t.Fatal(err)
	}
	err = dao.Delete(RuleTable, "1")
	if err != nil {
		t.Fatal(err)
	}
}
