package dao

// 告警接口
// Api.Id 与 Rule.Id 是一对多的关系
type Api struct {
	Id string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Method string `json:"method,omitempty"`
	Remark string `json:"remark"`
}

// 告警规则
// Api.Id 与 Rule.Id 是一对多的关系
type Rule struct {
	Id string `json:"id,omitempty"`
	ValueRule string `json:"valueRule,omitempty"`
	TimeRule string `json:"timeRule,omitempty"`
	Mails []string `json:"mails,omitempty"`
}

// 全局通知邮件列表
type GlobalMail struct {
	Mails []string `json:"mails,omitempty"`
}
