package search

import "time"

type VisitRecord struct {
	ID        uint      `gorm:"primarykey" json:"ID"`   // 主键ID
	UserId    uint      `gorm:"user_id" json:"user_id"` // 用户ID
	Date      time.Time `gorm:"date" json:"date"`       // 访问日期
	CreatedAt time.Time `json:"-"`                      // 创建时间
}

func (v *VisitRecord) TableName() string {
	return "visit_records"
}
