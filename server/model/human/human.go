package human

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common"
)

// 租户业主信息
type Human struct {
	global.GVA_MODEL
	Name          string               `json:"name"`                           // 共有人
	IdType        string               `json:"id_type"`                        // 证件类型
	IdNumber      string               `json:"id_number"`                      // 证件号
	Mobile        string               `json:"mobile"`                         // 手机号
	Attachments   common.AttachmentMap `json:"attachments" gorm:"TYPE:json"`   // 证件照片
	IdentityTypes common.IdentityTypes `json:"identityTypes" gorm:"TYPE:json"` // 身份标识
}

func (land *Human) TableName() string {
	return "humans"
}
