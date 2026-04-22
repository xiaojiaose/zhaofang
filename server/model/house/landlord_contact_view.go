package house

import (
	"time"

	"gorm.io/gorm"
)

const (
	// 后台展示状态：点击“联系方式”后先进入待审核，再由后台流转。
	LandlordViewStatusPending    = "待审核"
	LandlordViewStatusProcessing = "审核中"
	LandlordViewStatusApproved   = "已通过，未付款"
	LandlordViewStatusPaid       = "已付款"
	LandlordViewStatusRejected   = "未通过，已拒绝"
)

type LandlordContactView struct {
	ID                      uint           `gorm:"primarykey" json:"ID"` // 主键ID
	CreatedAt               time.Time      // 创建时间
	UpdatedAt               time.Time      // 更新时间
	DeletedAt               gorm.DeletedAt `gorm:"index" json:"-"`                                     // 删除时间
	ResourceID              uint           `json:"resourceId" gorm:"index;comment:房源ID"`               // 房源ID
	Xiaoqu                  string         `json:"xiaoqu" gorm:"comment:小区"`                           // 小区
	DoorNo                  string         `json:"doorNo" gorm:"comment:门牌号"`                          // 门牌号
	ResourceCreatedAt       time.Time      `json:"resourceCreatedAt" gorm:"comment:房源录入时间"`            // 房源录入时间
	PublisherUserID         uint           `json:"publisherUserId" gorm:"index;comment:发布人ID"`         // 发布人ID
	PublisherName           string         `json:"publisherName" gorm:"comment:发布人名称"`                 // 发布人名称
	PublisherPhone          string         `json:"publisherPhone" gorm:"index;comment:发布人手机号"`         // 发布人手机号
	ViewerUserID            uint           `json:"viewerUserId" gorm:"index;comment:查看人ID"`            // 查看人ID
	ViewerName              string         `json:"viewerName" gorm:"comment:查看人名称"`                    // 查看人名称
	ViewerPhone             string         `json:"viewerPhone" gorm:"index;comment:查看人手机号"`            // 查看人手机号
	Status                  string         `json:"status" gorm:"default:待审核;comment:当前状态"`             // 当前状态
	LastOperatedAtUnixMilli int64          `json:"lastOperatedAtUnixMilli" gorm:"comment:最后操作时间毫秒时间戳"` // 最后操作时间毫秒时间戳
}

func (LandlordContactView) TableName() string {
	return "house_landlord_contact_views"
}
