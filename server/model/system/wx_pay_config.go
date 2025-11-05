package system

// WxPayConfig 微信支付配置
type WxPayConfig struct {
	AppID     string `json:"appId"`     // 小程序AppID
	MchID     string `json:"mchId"`     // 商户号
	APIKey    string `json:"apiKey"`    // API密钥
	NotifyURL string `json:"notifyUrl"` // 支付回调地址
}
