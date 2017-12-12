package store

// 告警接口
type Api struct {
	Id string `toml:"id"`
	Name string `toml:"name"`
	Method string `toml:"method"`
}

// 告警规则
type Rule struct {
	Id string `toml:"id"`
	ValueRule string `toml:"valueRule"`
	TimeRule string `toml:"timeRule"`
	Mails []string `toml:"mails"`
}

// 全局通知邮件列表
type GlobalMail struct {
	Mails []string `toml:"mails"`
}

type TomlData struct {
	Api []Api `toml:"api"`
	Rule []Rule `toml:"rule"`
	GlobalMail GlobalMail `toml:"globalMail"`
}
