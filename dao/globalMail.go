package dao

// GlobalMailTable 表名
const GlobalMailTable = "globalMail"

// GlobalMail 全局通知邮件
type GlobalMail struct {
	Mail string `json:"mail"`
}

// GlobalMailDao 数据操作
type GlobalMailDao struct {
	dao *Dao
}

// NewGlobalMailDao 构造函数
func NewGlobalMailDao(dao *Dao) *GlobalMailDao {
	return &GlobalMailDao{dao}
}

// Add 新增
func (g *GlobalMailDao) Add(mail string) error {
	return g.dao.PutByByte(GlobalMailTable, mail, []byte(mail))
}

// Delete 删除
func (g *GlobalMailDao) Delete(mail string) error {
	return g.dao.Delete(GlobalMailTable, mail)
}

// Get 查询
func (g *GlobalMailDao) Get(mail string) ([]GlobalMail, error) {
	return g.dao.GetGlobalMailsByPrefix([]byte(GlobalMailTable), []byte(mail))
}

// GetAll 查询所有
func (g *GlobalMailDao) GetAll() ([]GlobalMail, error) {
	return g.dao.GetGlobalMailsAll(GlobalMailTable)
}