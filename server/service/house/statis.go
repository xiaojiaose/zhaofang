package house

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/search"
	"time"
)

type StatisService struct {
}

func (s StatisService) ByDate(start, end time.Time) (list []search.StatisData, err error) {
	db := global.GVA_DB.Model(&search.StatisData{})
	db.Where("date > ? and date < ?", start, end)
	err = db.Find(&list).Error
	return
}

func (s StatisService) VisitRecord(userId uint, info request.PageInfo, order string, desc bool) (list interface{}, total int64, err error) {

	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&search.VisitRecord{})
	var apiList []search.VisitRecord
	if userId != 0 {
		db = db.Where("user_id = ?", userId)
	}

	err = db.Count(&total).Error

	if err != nil {
		return apiList, total, err
	} else {
		db = db.Limit(limit).Offset(offset)
		if order != "" {
			OrderStr := order
			if desc {
				OrderStr = order + " desc"
			}

			err = db.Order(OrderStr).Find(&apiList).Error
		} else {
			err = db.Order("id desc").Find(&apiList).Error
		}
	}
	return apiList, total, err
}

func (s StatisService) InsertRecord(houseId uint, field string, nums ...int) (err error) {
	v := 1
	if len(nums) > 0 {
		v = nums[0]
	}
	return global.GVA_DB.Exec(fmt.Sprintf("insert into visit_daily (resource_id, %s, date) value (?, ?, ?)", field), houseId, v, time.Now()).Error
}
