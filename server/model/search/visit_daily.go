package search

import "time"

type VisitDaily struct {
	Date       time.Time `json:"date"`       // 日期
	ResourceId uint      `json:"resourceId"` // 房源id
	Follow     int       `json:"follow"`     // 关注次数
	View       int       `json:"view"`       // 浏览次数
	Click      int       `json:"click"`      // 电话获取次数
	Shared     int       `json:"shared"`     // 分享次数
}

func (VisitDaily) TableName() string {
	return "visit_daily"
}
