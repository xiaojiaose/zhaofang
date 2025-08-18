package model

type AliModel struct {
	Phones        []string `json:"phones"`        // 手机号 多个
	TemplateCode  string   `json:"templateCode"`  // 模板code
	TemplateParam string   `json:"templateParam"` // 参数验证码
}

type Ali struct {
	Phones []string `json:"phones"` // 手机号 多个
}
