package house

import (
	"time"

	"gorm.io/gorm"
)

const (
	// 发布人确认状态，描述“房源发布人是否认可这笔申请”。
	RewardPublisherPending  = "待确认"
	RewardPublisherRejected = "已拒绝"
	RewardPublisherApproved = "已确认"

	// 后台审核状态，描述“平台侧处理到哪一步”。
	// 这里和发布人确认状态拆开存，避免一个字段同时承担两段流程导致语义混乱。
	RewardAuditNone            = "未进入审核"
	RewardAuditPending         = "待审核"
	RewardAuditProcessing      = "审核中"
	RewardAuditApprovedPending = "审通过待发放"
	RewardAuditPaid            = "已发放"
	RewardAuditRejected        = "未通过"
)

type RewardApplication struct {
	ID              uint `gorm:"primarykey" json:"ID"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	ResourceID      uint           `json:"resourceId" gorm:"index;comment:房源ID"`
	ApplyUserID     uint           `json:"applyUserId" gorm:"index;comment:申请人ID"`
	PublisherUserID uint           `json:"publisherUserId" gorm:"index;comment:发布人ID"`
	// 下面这些联系人信息是申请创建时的快照，
	// 避免用户后续修改手机号/微信号后，历史审核单展示被“穿透更新”。
	ApplyUserPhone          string `json:"applyUserPhone" gorm:"comment:申请人手机号"`
	ApplyUserWxNo           string `json:"applyUserWxNo" gorm:"comment:申请人微信号"`
	PublisherUserPhone      string `json:"publisherUserPhone" gorm:"comment:发布人手机号"`
	PublisherUserWxNo       string `json:"publisherUserWxNo" gorm:"comment:发布人微信号"`
	Remark                  string `json:"remark" gorm:"type:text;comment:申请备注"`
	PublisherConfirmStatus  string `json:"publisherConfirmStatus" gorm:"default:待确认;comment:发布人确认状态"`
	AuditStatus             string `json:"auditStatus" gorm:"default:未进入审核;comment:后台审核状态"`
	LastOperatedAtUnixMilli int64  `json:"lastOperatedAtUnixMilli" gorm:"comment:最后操作时间毫秒"`
}

func (RewardApplication) TableName() string {
	return "house_reward_applications"
}
