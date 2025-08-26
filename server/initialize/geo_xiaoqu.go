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
)

func GeoXiaoqu() {
	global.GVA_LOG.Info("正在初始化小区数据...")
	global.GVA_DB.Find(&[]system.XiaoQu{})
	var xqList []system.XiaoQu
	tx := global.GVA_DB.Raw("SELECT id,name, longitude, latitude FROM xiao_qu;").Find(&xqList)
	if tx.Error != nil && !errors.Is(tx.Error, sql.ErrNoRows) {
		global.GVA_LOG.Error("小区数据初始化失败!", zap.Error(tx.Error))
	}

	geoService := test12.NewGeoService()

	for _, v := range xqList {
		lat, _ := strconv.ParseFloat(v.Latitude, 64)
		lon, _ := strconv.ParseFloat(v.Longitude, 64)
		geoService.AddCommunity(&test12.Community{
			ID:       strconv.Itoa(int(v.ID)),
			Name:     v.Name,
			Location: orb.Point{lat, lon}, // {116.3984, 39.9093},
		})
	}
	global.GVA_LOG.Info("小区数据初始化完成!")
}
