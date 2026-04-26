package main

import (
	"github.com/flipped-aurora/gin-vue-admin/server/core"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/initialize"
	service "github.com/flipped-aurora/gin-vue-admin/server/service/house"
	"go.uber.org/zap"
)

func Init() {
	global.GVA_VP = core.Viper() // 初始化Viper
	initialize.OtherInit()
	global.GVA_LOG = core.Zap() // 初始化zap日志库
	zap.ReplaceGlobals(global.GVA_LOG)
	global.GVA_DB = initialize.Gorm() // gorm连接数据库
	global.InitZincSearch(
		global.GVA_CONFIG.ZincSearch.Url,
		global.GVA_CONFIG.ZincSearch.Username,
		global.GVA_CONFIG.ZincSearch.Password,
		global.GVA_CONFIG.ZincSearch.ResourceIndex,
	)
}
func main() {
	Init()
	houseService := service.ResourceService{}
	total, err := houseService.RebuildAllResourceIndex(1000)
	if err != nil {
		global.GVA_LOG.Error("重建房源索引失败", zap.Error(err))
		return
	}
	global.GVA_LOG.Info("重建房源索引完成", zap.Int("count", total))
}
