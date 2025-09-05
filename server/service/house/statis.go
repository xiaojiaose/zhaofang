package house

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/search"
	"time"
)

type StatisService struct {
}

func (s StatisService) ByDate(start, end time.Time) (list []search.StatisData, err error) {
	db := global.GVA_DB.Model(&search.StatisData{})
	db.Where("start > ? and end < ?", start, end)
	err = db.Find(&list).Error
	return
}

func (s StatisService) VisitRecord(userId uint) (list []search.VisitRecord, err error) {
	db := global.GVA_DB.Model(&search.VisitRecord{})
	if userId > 0 {
		db.Where("user_id > ? ", userId)
	}
	err = db.Find(&list).Error
	return
}
