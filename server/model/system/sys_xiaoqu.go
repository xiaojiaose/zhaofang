package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

type XiaoQu struct {
	global.GVA_MODEL
	Name        string `json:"name"`         // name
	City        string `json:"city"`         // city
	Area        string `json:"area"`         // area
	AreaId      int    `json:"area_id"`      // area_id
	Position    string `json:"position"`     // postion
	Districts   string `json:"districts"`    // districts
	DistrictIds string `json:"district_ids"` // district_ids
	Address     string `json:"address"`      // address
	Latitude    string `json:"latitude"`     // latitude
	Longitude   string `json:"longitude"`    // longitude
	Able        bool   `json:"able"`         // 是否启用
	HouseNum    int    `json:"house_num"`    // 房源数量
}

func (s *XiaoQu) TableName() string {
	return "xiao_qu"
}
