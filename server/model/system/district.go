package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

type District struct {
	global.GVA_MODEL
	Name      string `json:"name"`    // name
	CityId    string `json:"city_id"` // city
	AreaId    string `json:"area_id"`
	Latitude  string `json:"latitude"`  // latitude
	Longitude string `json:"longitude"` // longitude
	Able      bool   `json:"able"`      // 是否启用
}

func (s *District) TableName() string {
	return "districts"
}
