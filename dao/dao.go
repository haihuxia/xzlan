package dao

import (
	"github.com/boltdb/bolt"
	"time"
	"encoding/json"
	"strconv"
	"strings"
	"bytes"
)

// Dao 数据操作实现
type Dao struct {
	Db *bolt.DB
}

// NewDao 构造函数
func NewDao(dbPath string) (*Dao, error) {
	boltDb, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 2 * time.Second})
	if err != nil {
		return &Dao{}, err
	}
	err = boltDb.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(APITable))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(RuleTable))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(NoteTable))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(GlobalMailTable))
		if err != nil {
			return err
		}
		return nil
	})
	return &Dao{boltDb}, err
}

// CreateTable 创建表
func (dao *Dao) CreateTable(table string) error {
	return dao.Db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(table))
		return err
	})
}

// DeleteTable 删除表
func (dao *Dao) DeleteTable(table string) error {
	return dao.Db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(table))
	})
}

// PutByByte 存储数据
func (dao *Dao) PutByByte(table string, key string, value []byte) error {
	return dao.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(table))
		err := b.Put([]byte(key), value)
		return err
	})
}

// GetSeq 获取序列
func (dao *Dao) GetSeq(table string) (string, error) {
	var seq string
	err := dao.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(APITable))
		v, err := b.NextSequence()
		if err != nil {
			return err
		}
		seq = strconv.FormatUint(v, 10)
		return nil
	})
	return seq, err
}

// UpdateAPI 更新 api 数据
func (dao *Dao) UpdateAPI(api API) error {
	return dao.PutByStruct(APITable, api.ID, api)
}

// PutByStruct 存储数据
func (dao *Dao) PutByStruct(table string, key string, value interface{}) error {
	v, e := json.Marshal(value)
	if e != nil {
		return e
	}
	return dao.PutByByte(table, key, v)
}

// Delete 删除
func (dao *Dao) Delete(table string, key string) error {
	return dao.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(table))
		return b.Delete([]byte(key))
	})
}

// DeleteAPI 删除 api 数据
func (dao *Dao) DeleteAPI(id string) error {
	tx, err := dao.Db.Begin(true)
	if err != nil {
		return err
	}
	err = tx.Bucket([]byte(NoteTable)).Delete([]byte(id))
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Bucket([]byte(RuleTable)).Delete([]byte(id))
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Bucket([]byte(APITable)).Delete([]byte(id))
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

// Get 查询
func (dao *Dao) Get(table string, key string) ([]byte, error) {
	var v []byte
	err := dao.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(table))
		v = b.Get([]byte(key))
		return nil
	})
	return v, err
}

// GetAPIs 查询 api 数据
func (dao *Dao) GetAPIs(name string, method string) (apis []API, err error) {
	err = dao.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(APITable))
		var api API
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

// GetAPIsAll 查询所有 api 数据
func (dao *Dao) GetAPIsAll() (apis []API, err error) {
	err = dao.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(APITable))
		var api API
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

// GetNotesAll 查询所有记录
func (dao *Dao) GetNotesAll() (notes []Note, err error) {
	err = dao.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(NoteTable))
		var note Note
		return b.ForEach(func(k, v []byte) error {
			if e := json.Unmarshal(v, &note); e != nil {
				return e
			}
			notes = append(notes, note)
			return nil
		})
	})
	return
}

// GetGlobalMailsByPrefix 根据前缀查询全局邮箱
func (dao *Dao) GetGlobalMailsByPrefix(table, prefix []byte) (values []GlobalMail, err error) {
	err = dao.Db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(table).Cursor()
		for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
			mail := &GlobalMail{string(v)}
			values = append(values, *mail)
		}
		return nil
	})
	return
}

// GetGlobalMailsAll 查询所有全局邮箱
func (dao *Dao) GetGlobalMailsAll(table string) (values []GlobalMail, err error) {
	err = dao.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(table))
		return b.ForEach(func(k, v []byte) error {
			mail := &GlobalMail{string(v)}
			values = append(values, *mail)
			return nil
		})
	})
	return
}