package task

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/house"
	"github.com/flipped-aurora/gin-vue-admin/server/model/search"
	sysModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

func StatisticVisits(db *gorm.DB) error {
	fmt.Println("定时统计访问量 start")

	now := time.Now().UTC()
	yesterdayStart := time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, time.UTC)
	// 更精确地表示为23:59:59.999999999
	yesterdayEndExact := yesterdayStart.Add(24*time.Hour - 1*time.Nanosecond)

	var totalViews, totalShared, totalFollow, totalClick int64

	// 方法3.1: 使用 Raw 和 Scan
	err := db.Raw("SELECT COALESCE(SUM(view), 0), COALESCE(SUM(shared), 0), COALESCE(SUM(follow), 0), COALESCE(SUM(click), 0) FROM visit_daily where date > ? and date < ?", yesterdayStart, yesterdayEndExact).Row().
		Scan(&totalViews, &totalShared, &totalFollow, &totalClick)
	if err != nil {
		global.GVA_LOG.Error("统计访问量失败!", zap.Error(err))
		return err
	}

	var AddSaler int64
	err = db.Model(&sysModel.SysUser{}).Where("authority_id = ? and created_at > ? and created_at < ?", 555, yesterdayStart, yesterdayEndExact).Count(&AddSaler).Error
	if err != nil {
		global.GVA_LOG.Error("统计sys_user失败!", zap.Error(err))
	}

	var AddHouse int64
	err = db.Model(&house.Resource{}).Where("created_at > ? and created_at < ?", yesterdayStart, yesterdayEndExact).Count(&AddHouse).Error
	if err != nil {
		global.GVA_LOG.Error("统计house_resources失败!", zap.Error(err))
	}

	var userIDs []uint
	err = db.Raw("select user_id from sys_operation_records where created_at > ? and created_at < ? group by user_id", yesterdayStart, yesterdayEndExact).Pluck("user_id", &userIDs).Error
	if err != nil {
		global.GVA_LOG.Error("统计sys_operation_records失败!", zap.Error(err))
	}

	db.Model(&search.StatisData{}).Where("date = ?", yesterdayStart).Save(&search.StatisData{
		View:     int(totalViews),
		Shared:   int(totalShared),
		Follow:   int(totalFollow),
		Click:    int(totalClick),
		Date:     yesterdayStart,
		AddSaler: int(AddSaler),
		UseSaler: len(userIDs),
		Add:      int(AddHouse),
	})
	fmt.Println("定时统计访问量 end")

	return err
}
