package dao

const GlobalMailTable = "globalMail"

// 全局通知邮件
type GlobalMail struct {
	Mail string `json:"mail"`
}

type GlobalMailDao struct {
	dao *Dao
}

func NewGlobalMailDao(dao *Dao) *GlobalMailDao {
	return &GlobalMailDao{dao}
}

func (g *GlobalMailDao) Add(mail string) error {
	return g.dao.PutByByte(GlobalMailTable, mail, []byte(mail))
}

func (g *GlobalMailDao) Delete(mail string) error {
	return g.dao.Delete(GlobalMailTable, mail)
}

func (g *GlobalMailDao) Get(mail string) ([]GlobalMail, error) {
	return g.dao.GetGlobalMailsByPrefix([]byte(GlobalMailTable), []byte(mail))
}

func (g *GlobalMailDao) GetAll() ([]GlobalMail, error) {
	return g.dao.GetGlobalMailsAll(GlobalMailTable)
}