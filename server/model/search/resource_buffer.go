package search

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/house"
	"strconv"
)

type ResourceBuffer struct {
	ID          string
	City        string `json:"city" `                // 所属城市
	Districts   string `json:"districts" `           // 所属商圈s
	DistrictIds string `json:"district_ids" `        // 所属商圈s
	XiaoquId    uint   `json:"xiaoqu_id"  `          // 所属小区id
	Xiaoqu      string `json:"xiaoqu"`               // 所属小区名字
	HouseType   string `json:"house_type,omitempty"` // 户型
	RentType    string `json:"rent_type,omitempty"`  // 出租类型
	Price       int    `json:"price"`                // 房源价格
	Feature     string `json:"feature" `             // 房源特色
	Status      string `json:"status" `              // 房源状态
}

func FromDeviceDB(entity *house.Resource) *ResourceBuffer {
	return &ResourceBuffer{
		ID:          strconv.Itoa(int(entity.ID)),
		City:        entity.City,
		Districts:   entity.Districts,
		DistrictIds: entity.DistrictIds,
		XiaoquId:    entity.XiaoquId,
		Xiaoqu:      entity.Xiaoqu,
		HouseType:   entity.HouseType,
		RentType:    entity.RentType,
		Price:       entity.Price,
		Feature:     entity.Feature,
		Status:      entity.Status,
	}
}

func FromDeviceES(data map[string]interface{}) *ResourceBuffer {
	result := &ResourceBuffer{}

	for key, value := range data {
		switch key {
		case "city":
			result.City = value.(string)
		case "districts":
			result.Districts = value.(string)
		case "district_ids":
			result.DistrictIds = value.(string)
		case "xiaoqu_id":
			result.XiaoquId = uint(value.(float64))
		case "xiaoqu":
			result.Xiaoqu = value.(string)
		case "house_type":
			result.HouseType = value.(string)
		case "rent_type":
			result.RentType = value.(string)
		case "price":
			result.Price = int(value.(int))
		case "feature":
			result.Feature = value.(string)
		case "status":
			result.Status = value.(string)
		}
	}

	return result
}

func (d *ResourceBuffer) ToData() map[string]interface{} {
	return map[string]interface{}{
		"city":         d.City,
		"districts":    d.Districts,
		"district_ids": d.DistrictIds,
		"xiaoqu_id":    d.XiaoquId,
		"xiaoqu":       d.Xiaoqu,
		"house_type":   d.HouseType,
		"rent_type":    d.RentType,
		"price":        d.Price,
		"feature":      d.Feature,
		"status":       d.Status,
	}
}
