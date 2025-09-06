package search

import (
	"time"
)

type StatisData struct {
	ID        uint      `gorm:"primarykey" json:"ID"`       // 主键ID
	CreatedAt time.Time `json:"-"`                          // 创建时间
	UpdatedAt time.Time `json:"-"`                          // 更新时间
	AddSaler  int       `json:"add_saler" form:"add_saler"` // 新增经纪人
	UseSaler  int       `json:"use_saler" form:"use_saler"` // 使用的经纪人
	Add       int       `json:"add" form:"add"`             // 新增帖子
	View      int       `json:"view" form:"view"`           // 帖子浏览数
	Follow    int       `json:"follow" form:"follow"`       // 帖子关注数
	Shared    int       `json:"shared" form:"shared"`       // 帖子分享数
	Click     int       `json:"click" form:"click"`         // 联系方式被点击数
	Date      time.Time `json:"-"`
}

func (StatisData) TableName() string {
	return "statis_data"
}
