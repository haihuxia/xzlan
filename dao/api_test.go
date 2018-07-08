package dao

import (
	"testing"
	"fmt"
)

func TestAPIDao_Add(t *testing.T) {
	dao, err := NewDao("/Users/tiger/project/logs/go/xzlan.db")
	defer dao.Db.Close()

	if err != nil {
		t.Fatal(err)
	}
	apiDao := NewAPIDao(dao)
	api := API{"", "role", "add", "新增角色", "stop", ""}
	err = apiDao.Add(api)
	if err != nil {
		t.Fatal(err)
	}
}

func TestApiDao_GetBy(t *testing.T) {
	dao, err := NewDao("/Users/tiger/project/logs/go/xzlan.db")
	defer dao.Db.Close()

	if err != nil {
		t.Fatal(err)
	}
	apiDao := NewAPIDao(dao)

	var name = "user"
	var method string
	apis, err := apiDao.GetBy(name, method)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < len(apis); i++ {
		fmt.Printf("value[%d]: %s \n", i, apis[i])
	}
}

func TestApiDao_GetAll(t *testing.T) {
	dao, err := NewDao("/Users/tiger/project/logs/go/xzlan.db")
	defer dao.Db.Close()

	if err != nil {
		t.Fatal(err)
	}
	apiDao := NewAPIDao(dao)
	apis, err := apiDao.GetAll()
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < len(apis); i++ {
		fmt.Printf("value[%d]: %s \n", i, apis[i])
	}
}



