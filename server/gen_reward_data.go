//go:build ignore

package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/model/house"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("用法: go run gen_reward_data.go <数量>")
		fmt.Println("示例: go run gen_reward_data.go 1000")
		os.Exit(1)
	}

	var count int
	fmt.Sscanf(os.Args[1], "%d", &count)
	if count <= 0 {
		fmt.Println("数量必须大于0")
		os.Exit(1)
	}

	// 初始化配置
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("读取配置失败: %v\n", err)
		os.Exit(1)
	}

	var c struct {
		DB struct {
			Path     string `mapstructure:"path"`
			Port     string `mapstructure:"port"`
			Username string `mapstructure:"username"`
			Password string `mapstructure:"password"`
			Dbname   string `mapstructure:"db-name"`
		} `mapstructure:"mysql"`
	}
	if err := viper.UnmarshalKey("mysql", &c.DB); err != nil {
		fmt.Printf("解析配置失败: %v\n", err)
		os.Exit(1)
	}

	// 连接数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.DB.Username, c.DB.Password, c.DB.Path, c.DB.Port, c.DB.Dbname)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("连接数据库失败: %v\n", err)
		os.Exit(1)
	}

	// 先查一下有多少资源可用
	var resourceCount int64
	db.Model(&house.Resource{}).Count(&resourceCount)
	fmt.Printf("当前有 %d 条房源\n", resourceCount)

	if resourceCount == 0 {
		fmt.Println("没有房源数据，请先导入房源")
		os.Exit(1)
	}

	// 查一下有多少用户
	var userCount int64
	db.Model(&system.SysUser{}).Count(&userCount)
	fmt.Printf("当前有 %d 个用户\n", userCount)

	if userCount < 2 {
		fmt.Println("用户数据不足，至少需要2个用户")
		os.Exit(1)
	}

	// 获取所有资源ID
	var resourceIDs []uint
	db.Model(&house.Resource{}).Pluck("id", &resourceIDs)

	// 获取所有用户ID
	var userIDs []uint
	db.Model(&system.SysUser{}).Where("id > 0").Pluck("id", &userIDs)

	// 生成数据
	publisherStatuses := []string{"待确认", "已拒绝", "已确认"}
	auditStatuses := []string{"未进入审核", "待审核", "审核中", "审核通过待发放", "已发放", "未通过"}

	var applications []house.RewardApplication
	inserted := 0

	for i := 0; i < count; i++ {
		resourceID := resourceIDs[rand.Intn(len(resourceIDs))]
		applyUserID := userIDs[rand.Intn(len(userIDs))]
		publisherUserID := userIDs[rand.Intn(len(userIDs))]

		// 避免申请人和发布人是同一个人
		for applyUserID == publisherUserID && len(userIDs) > 1 {
			publisherUserID = userIDs[rand.Intn(len(userIDs))]
		}

		// 随机生成时间在最近30天内
		daysAgo := rand.Intn(30)
		createdAt := time.Now().AddDate(0, 0, -daysAgo)
		lastOperated := createdAt.Add(time.Duration(rand.Intn(24*60)) * time.Minute)

		app := house.RewardApplication{
			ResourceID:              resourceID,
			ApplyUserID:             applyUserID,
			PublisherUserID:         publisherUserID,
			ApplyUserPhone:          fmt.Sprintf("138%08d", rand.Intn(100000000)),
			ApplyUserWxNo:           fmt.Sprintf("wx_%d", rand.Intn(100000)),
			PublisherUserPhone:      fmt.Sprintf("139%08d", rand.Intn(100000000)),
			PublisherUserWxNo:       fmt.Sprintf("wx_%d", rand.Intn(100000)),
			Remark:                  fmt.Sprintf("测试申请备注 %d", i+1),
			PublisherConfirmStatus:  publisherStatuses[rand.Intn(len(publisherStatuses))],
			AuditStatus:             auditStatuses[rand.Intn(len(auditStatuses))],
			LastOperatedAtUnixMilli: lastOperated.UnixMilli(),
		}
		applications = append(applications, app)

		// 每100条批量插入一次
		if len(applications) >= 100 {
			if err := db.Create(&applications).Error; err != nil {
				fmt.Printf("插入失败: %v\n", err)
				os.Exit(1)
			}
			inserted += len(applications)
			fmt.Printf("已插入 %d/%d 条\n", inserted, count)
			applications = applications[:0]
		}
	}

	// 插入剩余的数据
	if len(applications) > 0 {
		if err := db.Create(&applications).Error; err != nil {
			fmt.Printf("插入失败: %v\n", err)
			os.Exit(1)
		}
		inserted += len(applications)
		fmt.Printf("已插入 %d/%d 条\n", inserted, count)
	}

	// 验证
	var totalCount int64
	db.Model(&house.RewardApplication{}).Count(&totalCount)
	fmt.Printf("\n完成！共插入 %d 条数据，当前表共有 %d 条记录\n", inserted, totalCount)
}
