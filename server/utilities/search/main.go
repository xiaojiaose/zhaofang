package main

import (
	"context"
	"github.com/flipped-aurora/gin-vue-admin/server/core"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/initialize"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/house"
	"github.com/flipped-aurora/gin-vue-admin/server/model/search"
	service "github.com/flipped-aurora/gin-vue-admin/server/service/house"
	"go.uber.org/zap"
)

func Init() {
	global.GVA_VP = core.Viper() // 初始化Viper
	initialize.OtherInit()
	global.GVA_LOG = core.Zap() // 初始化zap日志库
	zap.ReplaceGlobals(global.GVA_LOG)
	global.GVA_DB = initialize.Gorm() // gorm连接数据库
	global.InitZincSearch(global.GVA_CONFIG.ZincSearch.Url, global.GVA_CONFIG.ZincSearch.Username, global.GVA_CONFIG.ZincSearch.Password)
}
func main() {
	Init()
	houseService := service.ResourceService{}
	list, _, err := houseService.GetPage(0, 0, "", "", request.PageInfo{Page: 1, PageSize: 1000}, "", false, request.SearchOther{})
	if err != nil {
		return
	}

	for _, h := range list.([]house.Resource) {
		func(hh house.Resource) {
			err = global.Gva_ResourceSearch.Add(context.Background(), *search.FromDeviceDB(&hh))
		}(h)
	}

}
