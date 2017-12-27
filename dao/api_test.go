package dao

import (
	"testing"
	"fmt"
	"time"
)

func TestApiDao_Add(t *testing.T) {
	dao, err := NewDao("/Users/tiger/project/logs/go/xzlan.db")
	defer dao.Db.Close()

	if err != nil {
		t.Fatal(err)
	}
	apiDao := NewApiDao(dao)
	api := Api{"", "role", "add", "新增角色", "stop", ""}
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
	apiDao := NewApiDao(dao)

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
	apiDao := NewApiDao(dao)
	apis, err := apiDao.GetAll()
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < len(apis); i++ {
		fmt.Printf("value[%d]: %s \n", i, apis[i])
	}

	tt, _ := time.Parse("2006-01-02 15:04:05", "2017-12-27 10:28:28")
	tt = tt.Add(2 * time.Hour)
	fmt.Println(tt)
	fmt.Printf("在之后吗？ %t \n", tt.After(time.Now()))

	var aa = "14d"
	fmt.Println(string(aa[len(aa)-1]))
	fmt.Println(string(aa[:len(aa)-1]))
}



