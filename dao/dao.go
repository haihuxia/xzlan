package dao

import (
	"github.com/boltdb/bolt"
	"time"
	"encoding/json"
	"strconv"
	"strings"
)

const (
	ApiTable = "api"
	RuleTable = "rule"
)

type Dao struct {
	Db *bolt.DB
}

func NewDao(dbPath string) (Dao, error) {
	boltDb, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 2 * time.Second})
	if err != nil {
		return Dao{}, err
	}
	return Dao{boltDb}, nil
}

func (dao *Dao) CreateTable(table string) error {
	return dao.Db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(table))
		return err
	})
}

func (dao *Dao) DeleteTable(table string) error {
	return dao.Db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(table))
	})
}

func (dao *Dao) PutByByte(table string, key string, value []byte) error {
	return dao.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(table))
		err := b.Put([]byte(key), value)
		return err
	})
}

func (dao *Dao) PutApi(api Api) error {
	err := dao.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ApiTable))
		v, err := b.NextSequence()
		if err != nil {
			return err
		}
		seq := strconv.FormatUint(v, 10)
		api.Id = string(seq)
		return nil
	})
	if err != nil {
		return err
	}
	return dao.PutByStruct(ApiTable, api.Id, api)
}

func (dao *Dao) PutByStruct(table string, key string, value interface{}) error {
	v, e := json.Marshal(value)
	if e != nil {
		return e
	}
	return dao.PutByByte(table, key, v)
}

func (dao *Dao) Delete(table string, key string) error {
	return dao.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(table))
		return b.Delete([]byte(key))
	})
}

func (dao *Dao) Get(table string, key string) ([]byte, error) {
	var v []byte
	err := dao.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(table))
		v = b.Get([]byte(key))
		return nil
	})
	return v, err
}

func (dao *Dao) GetApis(name string, method string) (apis []Api, err error) {
	err = dao.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ApiTable))
		var api Api
		return b.ForEach(func(k, v []byte) error {
			if e := json.Unmarshal(v, &api); e != nil {
				return e
			}
			if name != "" && method != "" {
				if strings.Index(api.Name, name) != -1 && strings.Index(api.Method, method) != -1 {
					apis = append(apis, api)
				}
			}
			if name != "" && method == "" && strings.Index(api.Name, name) != -1 {
				apis = append(apis, api)
			}
			if name == "" && method != "" && strings.Index(api.Method, method) != -1 {
				apis = append(apis, api)
			}
			if name == "" && method == "" {
				apis = append(apis, api)
			}
			return nil
		})
	})
	return
}

func (dao *Dao) GetApisAll() (apis []Api, err error) {
	err = dao.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ApiTable))
		var api Api
		return b.ForEach(func(k, v []byte) error {
			if e := json.Unmarshal(v, &api); e != nil {
				return e
			}
			apis = append(apis, api)
			return nil
		})
	})
	return
}
