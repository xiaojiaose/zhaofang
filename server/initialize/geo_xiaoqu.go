package initialize

import (
	"database/sql"
	"errors"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	test12 "github.com/flipped-aurora/gin-vue-admin/server/utils/test"
	"github.com/paulmach/orb"
	"go.uber.org/zap"
	"strconv"
	"time"
)

func GeoXiaoqu() {
	loadGeoXiaoqu()
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			loadGeoXiaoqu()
		}
	}()
}

func loadGeoXiaoqu() {
	global.GVA_LOG.Info("正在加载小区地理数据...")
	var xqList []system.XiaoQu
	tx := global.GVA_DB.Raw("SELECT id,name, longitude, latitude FROM xiao_qu;").Find(&xqList)
	if tx.Error != nil && !errors.Is(tx.Error, sql.ErrNoRows) {
		global.GVA_LOG.Error("小区地理数据加载失败!", zap.Error(tx.Error))
		return
	}

	// 先构建完整索引，再一次性替换全局引用，避免刷新期间暴露不完整数据
	geoService := test12.NewGeoServiceWithoutGlobal()

	for _, v := range xqList {
		lat, _ := strconv.ParseFloat(v.Latitude, 64)
		lon, _ := strconv.ParseFloat(v.Longitude, 64)
		geoService.AddCommunity(&test12.Community{
			ID:       strconv.Itoa(int(v.ID)),
			Name:     v.Name,
			Location: orb.Point{lat, lon},
		})
	}

	// 索引构建完成后替换全局引用
	test12.GeoSearch = geoService

	global.GVA_LOG.Info("小区地理数据加载完成!", zap.Int("count", len(xqList)))
}
