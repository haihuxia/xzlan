package dao

import (
	"testing"
	"fmt"
	"encoding/json"
)

func TestDao_CreateTable(t *testing.T) {
	dao, err := NewDao("/Users/tiger/project/logs/go/xzlan.db")
	defer dao.Db.Close()

	if err != nil {
		t.Fatal(err)
	}
	err = dao.CreateTable(ApiTable)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDao_DeleteTable(t *testing.T) {
	dao, err := NewDao("/Users/tiger/project/logs/go/xzlan.db")
	defer dao.Db.Close()

	if err != nil {
		t.Fatal(err)
	}
	err = dao.DeleteTable(ApiTable)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDao_Put(t *testing.T) {
	dao, err := NewDao("/Users/tiger/project/logs/go/xzlan.db")
	defer dao.Db.Close()

	if err != nil {
		t.Fatal(err)
	}
	api := Api{"1", "user", "get", "查询接口"}
	apiJson, e := json.Marshal(api)
	if e != nil {
		t.Fatal(e)
	}
	err = dao.Put(ApiTable, "1", apiJson)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDao_PutApi(t *testing.T) {
	dao, err := NewDao("/Users/tiger/project/logs/go/xzlan.db")
	defer dao.Db.Close()

	if err != nil {
		t.Fatal(err)
	}
	api := Api{"", "user", "get", "查询接口"}
	err = dao.PutApi(api)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDao_Get(t *testing.T) {
	dao, err := NewDao("/Users/tiger/project/logs/go/xzlan.db")
	defer dao.Db.Close()

	if err != nil {
		t.Fatal(err)
	}
	api, e := dao.Get(ApiTable, "1")
	if api == nil || e != nil {
		t.Fatal(e)
	}
	var a Api
	e = json.Unmarshal(api, &a)
	if e != nil {
		t.Fatal(e)
	}
	fmt.Printf("value: %s", a)
}

func TestDao_GetApisAll(t *testing.T) {
	dao, err := NewDao("/Users/tiger/project/logs/go/xzlan.db")
	defer dao.Db.Close()

	if err != nil {
		t.Fatal(err)
	}
	v, e := dao.GetApisAll()
	if e != nil {
		t.Fatal(e)
	}
	for i := 0; i < len(v); i++ {
		fmt.Printf("value[%d]: %s \n", i, v[i])
	}
}

func TestDao_GetApis(t *testing.T) {
	dao, err := NewDao("/Users/tiger/project/logs/go/xzlan.db")
	defer dao.Db.Close()

	if err != nil {
		t.Fatal(err)
	}
	var name = "user"
	var method string
	v, e := dao.GetApis(name, method)
	if e != nil {
		t.Fatal(e)
	}
	for i := 0; i < len(v); i++ {
		fmt.Printf("\n value[%d]: %s", i, v[i])
	}
}

func TestDao_Delete(t *testing.T) {
	dao, err := NewDao("/Users/tiger/project/logs/go/xzlan.db")
	defer dao.Db.Close()

	if err != nil {
		t.Fatal(err)
	}
	err = dao.Delete(ApiTable, "1")
	if err != nil {
		t.Fatal(err)
	}
}
