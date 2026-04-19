package house

import (
	"time"

	"gorm.io/gorm"
)

type ShareToken struct {
	ID        uint           `gorm:"primarykey" json:"ID"` // 主键ID
	CreatedAt time.Time      // 创建时间
	UpdatedAt time.Time      // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`                   // 删除时间
	UserID    uint           `json:"userId" gorm:"index;comment:用户ID"` // 分享用户ID
	// 分享页通过 token 访问，不直接暴露用户 id，方便后续单独失效和审计。
	Token        string `json:"token" gorm:"uniqueIndex;size:128;comment:分享token"` // 分享token
	Status       string `json:"status" gorm:"default:active;comment:状态"`           // 状态
	ExpireAtUnix int64  `json:"expireAtUnix" gorm:"index;comment:过期时间秒"`           // 过期时间秒级时间戳
}

func (ShareToken) TableName() string {
	return "house_share_tokens"
}
