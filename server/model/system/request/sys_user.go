package request

import (
	common "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
)

// Register User register structure
type Register struct {
	Username     string `json:"userName" example:"用户名"`
	Password     string `json:"passWord" example:"密码"`
	NickName     string `json:"nickName" example:"昵称"`
	HeaderImg    string `json:"headerImg" example:"头像链接"`
	AuthorityId  uint   `json:"authorityId" swaggertype:"string" example:"int 角色id"`
	Enable       int    `json:"enable" swaggertype:"string" example:"int 是否启用"`
	AuthorityIds []uint `json:"authorityIds" swaggertype:"string" example:"[]uint 角色id"`
	Phone        string `json:"phone" example:"电话号码"`
	Email        string `json:"email" example:"电子邮箱"`
}

// Login User login structure
type Login struct {
	Username  string `json:"username"`  // 用户名
	Password  string `json:"password"`  // 密码
	Captcha   string `json:"captcha"`   // 验证码
	CaptchaId string `json:"captchaId"` // 验证码ID

	Mobile string `json:"mobile"` // 手机号
	Sms    string `json:"sms"`    // 短信验证码
}

type MobileLogin struct {
	Mobile string `json:"mobile"` // 手机号
	Sms    string `json:"sms"`    // 短信验证码
}

// ChangePasswordReq Modify password structure
type ChangePasswordReq struct {
	ID          uint   `json:"-"`           // 从 JWT 中提取 user id，避免越权
	Password    string `json:"password"`    // 密码
	NewPassword string `json:"newPassword"` // 新密码
}

type ResetPassword struct {
	ID       uint   `json:"ID" form:"ID"`
	Password string `json:"password" form:"password" gorm:"comment:用户登录密码"` // 用户登录密码
}

// SetUserAuth Modify user's auth structure
type SetUserAuth struct {
	AuthorityId uint `json:"authorityId"` // 角色ID
}

// SetUserAuthorities Modify user's auth structure
type SetUserAuthorities struct {
	ID           uint
	AuthorityIds []uint `json:"authorityIds"` // 角色ID
}

type ChangeUserInfo struct {
	ID           uint                  `gorm:"primarykey"`                                                                               // 主键ID
	NickName     string                `json:"nickName" gorm:"default:系统用户;comment:用户昵称"`                                        // 用户昵称
	Phone        string                `json:"phone"  gorm:"comment:用户手机号"`                                                         // 用户手机号
	AuthorityIds []uint                `json:"authorityIds" gorm:"-"`                                                                    // 角色ID
	Email        string                `json:"email"  gorm:"comment:用户邮箱"`                                                           // 用户邮箱
	HeaderImg    string                `json:"headerImg" gorm:"default:https://qmplusimg.henrongyi.top/gva_header.jpg;comment:用户头像"` // 用户头像
	Enable       int                   `json:"enable" gorm:"comment:冻结用户"`                                                           //冻结用户
	Authorities  []system.SysAuthority `json:"-" gorm:"many2many:sys_user_authority;"`
}

type GetUserList struct {
	common.PageInfo
	Username    string `json:"username" form:"username"`
	NickName    string `json:"nickName" form:"nickName"`
	Phone       string `json:"phone" form:"phone"` // 手机号搜索
	Email       string `json:"email" form:"email"`
	AuthorityId int    `json:"authorityId" form:"authorityId"`
	Bind        int    `json:"bind" form:"bind"` // 绑定的状态 1 已绑定 2 未绑定
}

type WxLogin struct {
	Mobile string `json:"mobile"` // 出租人手机号
	Smn    string `json:"smn"`    // 短息验证码
	WxLoginRequest
}

type WxLoginRequest struct {
	Code          string `json:"code"`
	EncryptedData string `json:"encryptedData"`
	Iv            string `json:"iv"`
	RawData       string `json:"rawData"`
	Signature     string `json:"signature"`
}

type UserInfo struct {
	NickName  string `json:"nickName"`
	AvatarURL string `json:"avatarUrl"`
	Gender    int    `json:"gender"`
	City      string `json:"city"`
	Country   string `json:"country"`
	Province  string `json:"province"`
}
