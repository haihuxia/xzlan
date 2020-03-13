package dao

import (
	"encoding/json"
	"time"
)

// NoteTable 告警记录
const NoteTable = "note"

// Note 告警记录
// 以 Api.Id 为 key
type Note struct {
	Time string `json:"time"`
	Status string `json:"status"`
	Remark string `json:"remark"`
}

// NoteDao 数据操作
type NoteDao struct {
	dao *Dao
}

// NewNoteDao 构造函数
func NewNoteDao(dao *Dao) *NoteDao {
	return &NoteDao{dao}
}

// Add 新增
func (n *NoteDao) Add(remark string, id string) (string, error) {
	// 1: 未通知, 2: 不通知，0: 已通知
	note := &Note{now(), "1", remark}
	return note.Time, n.dao.PutByStruct(NoteTable, id, note)
}

// Update 更新
func (n *NoteDao) Update(remark string, id string, status string) error {
	note := &Note{now(), status, remark}
	return n.dao.PutByStruct(NoteTable, id, note)
}

// Get 查询
func (n *NoteDao) Get(id string) (note Note, err error) {
	v, err := n.dao.Get(NoteTable, id)
	if err != nil {
		return note, err
	}
	err = json.Unmarshal(v, &note)
	return
}

// GetAll 查询所有
func (n *NoteDao) GetAll() ([]Note, error) {
	return n.dao.GetNotesAll()
}

func now() string {
	return time.Now().Format("2006-01-02 15:04:05")
}