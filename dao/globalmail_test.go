package dao

import (
	"testing"
	"fmt"
)

func TestGlobalMailDao_GetAll(t *testing.T) {
	dao, err := NewDao("/Users/tiger/project/logs/go/xzlan.db")
	defer dao.Db.Close()

	if err != nil {
		t.Fatal(err)
	}
	globalmailDao := NewGlobalMailDao(dao)
	apis, err := globalmailDao.GetAll()
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < len(apis); i++ {
		fmt.Printf("value[%d]: %s \n", i, apis[i])
	}
}
