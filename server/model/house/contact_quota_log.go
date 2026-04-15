package house

import (
	"time"

	"gorm.io/gorm"
)

const (
	ContactQuotaActionGrant = "grant"
	ContactQuotaActionUse   = "use"
)

type ContactQuotaLog struct {
	ID                uint `gorm:"primarykey" json:"ID"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
	UserID            uint           `json:"userId" gorm:"index;comment:用户ID"`
	OperatorID        uint           `json:"operatorId" gorm:"index;comment:操作人ID"`
	UserPhone         string         `json:"userPhone" gorm:"index;comment:用户手机号"`
	Action            string         `json:"action" gorm:"comment:操作类型"`
	ChangeAmount      int            `json:"changeAmount" gorm:"comment:变更次数"`
	UsedAmount        int            `json:"usedAmount" gorm:"comment:已使用次数"`
	RemainingAmount   int            `json:"remainingAmount" gorm:"comment:剩余次数"`
	Remark            string         `json:"remark" gorm:"type:text;comment:备注"`
	RelatedResourceID uint           `json:"relatedResourceId" gorm:"comment:关联房源ID"`
}

func (ContactQuotaLog) TableName() string {
	return "house_contact_quota_logs"
}
