package dao

import (
	"encoding/json"
)

// APITable 表名
const APITable = "api"

// API 告警接口
// Api.Id 与 rule 表的 key 一致
// 以 Api.Id 为 key
type API struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Method string `json:"method"`
	Remark string `json:"remark"`
	Status string `json:"status"`
	NotifyTime string `json:"notifyTime"`
}

// APIDao 数据操作
type APIDao struct {
	dao *Dao
}

// NewAPIDao 构造函数
func NewAPIDao(dao *Dao) *APIDao {
	return &APIDao{dao}
}

// Add 新增
func (a *APIDao) Add(api API) error {
	seq, err := a.dao.GetSeq(APITable)
	if err != nil {
		return err
	}
	api.ID = seq
	return a.dao.PutByStruct(APITable, seq, api)
}

// Get 查询
func (a *APIDao) Get(id string) (api API, err error) {
	v, err := a.dao.Get(APITable, id)
	if err != nil {
		return api, err
	}
	err = json.Unmarshal(v, &api)
	return
}

// GetBy 查询
func (a *APIDao) GetBy(name string, method string) ([]API, error) {
	return a.dao.GetAPIs(name, method)
}

// GetAll 查询所有
func (a *APIDao) GetAll() ([]API, error) {
	return a.dao.GetAPIsAll()
}

// Update 更新
func (a *APIDao) Update(api API) error {
	return a.dao.PutByStruct(APITable, api.ID, api)
}

// Delete 删除
func (a *APIDao) Delete(id string) error {
	return a.dao.DeleteAPI(id)
}