package task

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/house"
	"github.com/flipped-aurora/gin-vue-admin/server/model/search"
	sysModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"log"
	"time"
)

func StatisticVisits(db *gorm.DB) error {
	fmt.Println("定时统计访问量 start")

	now := time.Now().Local()
	yesterdayStart := time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, time.Local)
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

func StatisticSalerVisit(db *gorm.DB) error {

	now := time.Now().Local()
	yesterdayStart := time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, time.Local)
	// 更精确地表示为23:59:59.999999999
	yesterdayEndExact := yesterdayStart.Add(24*time.Hour - 1*time.Nanosecond)

	//err := db.Raw("select user_id from sys_operation_records where created_at > ? and created_at < ? order by created_at desc ", yesterdayStart, yesterdayEndExact).Error
	//if err != nil {
	//	global.GVA_LOG.Error("统计sys_operation_records失败!", zap.Error(err))
	//}

	var lastID uint = 0
	batchSize := 1000
	var allUserIDs []uint

	visitRecordMap := make(map[uint]search.VisitRecord)
	for {
		var records []struct {
			CreatedAt time.Time `gorm:"column:created_at"`
			UserID    uint      `gorm:"column:user_id"`
			ID        uint      `gorm:"column:id"`
		}

		query := db.Table("sys_operation_records").
			Select("id, user_id, created_at").
			Where("created_at > ? AND created_at < ?", yesterdayStart, yesterdayEndExact).
			Order("id DESC").
			Limit(batchSize)

		if lastID > 0 {
			query = query.Where("id < ?", lastID)
		}

		err := query.Find(&records).Error
		if err != nil {
			log.Println("查询失败:", err)
			return err
		}

		if len(records) == 0 {
			break
		}

		// 提取用户ID
		for _, record := range records {
			if _, ok := visitRecordMap[record.UserID]; !ok {
				visitRecordMap[record.UserID] = search.VisitRecord{
					UserId: record.UserID,
					Date:   record.CreatedAt,
				}
			}

			lastID = record.ID // 更新最后一条记录的ID
		}

		fmt.Printf("已处理批次，获取 %d 条记录，总计: %d\n", len(records), len(allUserIDs))

		if len(records) < batchSize {
			break // 最后一页
		}

		time.Sleep(5 * time.Millisecond)
	}

	fmt.Printf("总共获取 %d 个用户\n", len(visitRecordMap))

	if len(visitRecordMap) > 0 {
		for _, record := range visitRecordMap {
			record.CreatedAt = time.Now()
			err := db.Model(&search.VisitRecord{}).Create(&record).Error
			if err == nil {
				log.Println("visitRecord失败:", err)
			} else {
				log.Println(fmt.Sprintf("visitRecord info: %v :", record))
			}
		}
	}
	return nil
}
