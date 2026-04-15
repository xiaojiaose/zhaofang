package house

import (
	"time"

	"gorm.io/gorm"
)

type ShareToken struct {
	ID           uint `gorm:"primarykey" json:"ID"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	UserID       uint           `json:"userId" gorm:"index;comment:用户ID"`
	Token        string         `json:"token" gorm:"uniqueIndex;size:128;comment:分享token"`
	Status       string         `json:"status" gorm:"default:active;comment:状态"`
	ExpireAtUnix int64          `json:"expireAtUnix" gorm:"index;comment:过期时间秒"`
}

func (ShareToken) TableName() string {
	return "house_share_tokens"
}
