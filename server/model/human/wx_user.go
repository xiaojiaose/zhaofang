package human

import "github.com/flipped-aurora/gin-vue-admin/server/global"

// 微信登录用户信息
type WxUser struct {
	global.GVA_MODEL
	Mobile     string `json:"mobile"`     // 手机好
	Headimgurl string `json:"headimgurl"` // 头像
	Nickname   string `json:"nickname"`   // 昵称
	Openid     string `json:"openid"`     // openid
	Sex        int    `json:"sex"`        // 性别
}

func (land WxUser) TableName() string {
	return "wx_users"
}
