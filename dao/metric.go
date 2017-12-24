package dao

// 告警接口
// Api 与 Rule 是一对多的关系
type Api struct {
	Id string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Method string `json:"method,omitempty"`
	Remark string `json:"remark,omitempty"`
	Status string `json:"status,omitempty"`
}

// 告警规则
// Api 与 Rule 是一对多的关系
type Rule struct {
	Type string `json:"type,omitempty"`
	Max string `json:"max,omitempty"`
	Min string `json:"min,omitempty"`
	Val string `json:"val,omitempty"`
	Time string `json:"time,omitempty"`
	Count string `json:"count,omitempty"`
	Mails string `json:"mails,omitempty"`
}

// 全局通知邮件列表
type GlobalMail struct {
	Mails string `json:"mails,omitempty"`
}

type Message struct {
	Message string `json:"message,omitempty"`
}
