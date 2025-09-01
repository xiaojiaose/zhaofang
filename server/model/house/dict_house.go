package house

import (
	"gorm.io/gorm"
	"time"
)

type DictHouse struct {
	ID uint `gorm:"primarykey" json:"ID"` // 主键ID

	EncryptHouseName string `json:"encryptHouseName"` // 包含 Unicode 转义序列
	UnitOpenId       string `json:"unitOpenId"`
	HouseOpenId      string `json:"houseOpenId"`

	CreatedAt time.Time      // 创建时间
	UpdatedAt time.Time      // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // 删除时间
}

func (DictHouse) TableName() string {
	return "dict_house"
}
