package store

import (
	"testing"
	"fmt"
)

var filePath = "/Users/tiger/abc.tml"

func TestLoad(t *testing.T) {
	d, err := Load(filePath)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%p", d)
}

func TestSave(t *testing.T) {
	var data = TomlData{
		[]Api{{"1", "1", "1"}},
		[]Rule{{"2", "2", "2", []string{"2", "2"}}},
		GlobalMail{[]string{"3"}},
	}
	if err := Save(filePath, data); err != nil {
		t.Fatal(err)
	}
}