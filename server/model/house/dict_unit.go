package house

import (
	"gorm.io/gorm"
	"time"
)

type DictUnit struct {
	ID uint `gorm:"primarykey" json:"ID"` // 主键ID

	BuildingOpenID  string `json:"buildingOpenId"`
	EncryptUnitName string `json:"encryptUnitName"` // 包含 Unicode 转义序列
	UnitOpenID      string `json:"unitOpenId"`

	CreatedAt time.Time      // 创建时间
	UpdatedAt time.Time      // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // 删除时间
}

func (DictUnit) TableName() string {
	return "dict_unit"
}
