package house

import (
	"time"

	"gorm.io/gorm"
)

const (
	// grant 表示后台充值，use 表示前台查看联系方式时扣减。
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
	ChangeAmount      int            `json:"changeAmount" gorm:"comment:本次变更次数，充值为正，扣减为负"`
	UsedAmount        int            `json:"usedAmount" gorm:"comment:本次消耗次数"`
	RemainingAmount   int            `json:"remainingAmount" gorm:"comment:该次操作后的剩余次数"`
	Remark            string         `json:"remark" gorm:"type:text;comment:备注"`
	RelatedResourceID uint           `json:"relatedResourceId" gorm:"comment:关联房源ID"`
}

func (ContactQuotaLog) TableName() string {
	return "house_contact_quota_logs"
}
