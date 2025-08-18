package house

import (
	"gorm.io/gorm"
	"time"
)

type Favorite struct {
	ID         uint           `gorm:"primarykey" json:"ID"` // 主键ID
	CreatedAt  time.Time      // 创建时间
	UpdatedAt  time.Time      // 更新时间
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"` // 删除时间
	UserId     uint           `json:"user_id"`
	ResourceId uint           `json:"resource_id"`
}

func (Favorite) TableName() string {
	return "favorite"
}
