package dao

import (
	"encoding/json"
)

const ApiTable = "api"

// 告警接口
// Api.Id 与 rule 表的 key 一致
// 以 Api.Id 为 key
type Api struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Method string `json:"method"`
	Remark string `json:"remark"`
	Status string `json:"status"`
	NotifyTime string `json:"notifyTime"`
}

type ApiDao struct {
	dao *Dao
}

func NewApiDao(dao *Dao) *ApiDao {
	return &ApiDao{dao}
}

func (a *ApiDao) Add(api Api) error {
	seq, err := a.dao.GetSeq(ApiTable)
	if err != nil {
		return err
	}
	api.Id = seq
	return a.dao.PutByStruct(ApiTable, seq, api)
}

func (a *ApiDao) Get(id string) (api Api, err error) {
	v, err := a.dao.Get(ApiTable, id)
	if err != nil {
		return api, err
	}
	err = json.Unmarshal(v, &api)
	return
}

func (a *ApiDao) GetBy(name string, method string) ([]Api, error) {
	return a.dao.GetApis(name, method)
}

func (a *ApiDao) GetAll() ([]Api, error) {
	return a.dao.GetApisAll()
}

func (a *ApiDao) Update(api Api) error {
	return a.dao.PutByStruct(ApiTable, api.Id, api)
}

func (a *ApiDao) Delete(id string) error {
	return a.dao.Delete(ApiTable, id)
}