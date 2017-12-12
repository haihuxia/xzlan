package store

import (
	"os"
	"github.com/BurntSushi/toml"
	"fmt"
	"path/filepath"
	"bytes"
	"io/ioutil"
)

func Load(filePath string) (*TomlData, error) {
	d := &TomlData {}
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return d, err
	}
	tomlAbsPath, err := filepath.Abs(filePath)
	if err != nil {
		return d, err
	}
	if _, err := toml.DecodeFile(tomlAbsPath, &d); err != nil {
		return d, err
	}
	fmt.Println("")
	fmt.Printf("d: %q", d)
	return d, nil
}

func Save(filePath string, data TomlData) error {
	var tomlBuffer bytes.Buffer
	e := toml.NewEncoder(&tomlBuffer)
	err := e.Encode(data)
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(filePath, tomlBuffer.Bytes(), os.FileMode(0666)); err != nil {
		return err
	}
	return nil
}
