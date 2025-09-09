package initialize

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/task"

	"github.com/robfig/cron/v3"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

func Timer() {
	go func() {
		var option []cron.Option
		option = append(option, cron.WithSeconds())
		// 清理DB定时任务
		_, err := global.GVA_Timer.AddTaskByFunc("ClearDB", "@daily", func() {
			err := task.ClearTable(global.GVA_DB) // 定时任务方法定在task文件包中
			if err != nil {
				fmt.Println("timer error:", err)
			}
		}, "定时清理数据库【日志，黑名单】内容", option...)
		if err != nil {
			fmt.Println("add timer error:", err)
		}

		// 其他定时任务定在这里 参考上方使用方法

		_, err = global.GVA_Timer.AddTaskByFunc("StatistVisits", "0 0 2 * * *", func() {
			err = task.StatisticVisits(global.GVA_DB) // 定时任务方法定在task文件包中
			if err != nil {
				fmt.Println("timer StatistVisits error:", err)
			}
		}, "", option...)
		if err != nil {
			fmt.Println("add StatistVisits error:", err)
		}

		_, err = global.GVA_Timer.AddTaskByFunc("StatisticSalerVisit", "0 0 1 * * *", func() {
			err = task.StatisticSalerVisit(global.GVA_DB) // 定时任务方法定在task文件包中
			if err != nil {
				fmt.Println("timer StatisticSalerVisit error:", err)
			}
		}, "", option...)
		if err != nil {
			fmt.Println("add StatisticSalerVisit error:", err)
		}

	}()
}
