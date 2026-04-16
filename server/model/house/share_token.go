package house

import (
	"time"

	"gorm.io/gorm"
)

type ShareToken struct {
	ID        uint `gorm:"primarykey" json:"ID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	UserID    uint           `json:"userId" gorm:"index;comment:用户ID"`
	// 分享页通过 token 访问，不直接暴露用户 id，方便后续单独失效和审计。
	Token        string `json:"token" gorm:"uniqueIndex;size:128;comment:分享token"`
	Status       string `json:"status" gorm:"default:active;comment:状态"`
	ExpireAtUnix int64  `json:"expireAtUnix" gorm:"index;comment:过期时间秒"`
}

func (ShareToken) TableName() string {
	return "house_share_tokens"
}
