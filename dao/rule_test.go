package dao

import (
	"testing"
	"fmt"
)

func TestRuleDao_Get(t *testing.T) {
	dao, err := NewDao("/Users/tiger/project/logs/go/xzlan.db")
	defer dao.Db.Close()

	if err != nil {
		t.Fatal(err)
	}
	ruleDao := NewRuleDao(dao)
	rule, err := ruleDao.Get("1")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(rule)
}
