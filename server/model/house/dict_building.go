package house

import (
	"gorm.io/gorm"
	"time"
)

type DictBuilding struct {
	ID uint `gorm:"primarykey" json:"ID"` // 主键ID

	BuildingOpenID      string `json:"buildingOpenId"`
	CommunityID         int    `json:"communityId"`
	EncryptBuildingName string `json:"buildingName"` // 注意：这里包含 Unicode 转义序列

	CreatedAt time.Time      // 创建时间
	UpdatedAt time.Time      // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // 删除时间
}

func (DictBuilding) TableName() string {
	return "dict_building"
}
