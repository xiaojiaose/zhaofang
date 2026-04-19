package search

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/house"
	"strconv"
	"time"
)

type ResourceBuffer struct {
	ID              string
	HouseId         string `json:"house_id"`
	Owner           string `json:"owner"`
	City            string `json:"city" `                // 所属城市
	Districts       string `json:"districts" `           // 所属商圈s
	DistrictIds     string `json:"district_ids" `        // 所属商圈s
	XiaoquId        uint   `json:"xiaoqu_id"  `          // 所属小区id
	Xiaoqu          string `json:"xiaoqu"`               // 所属小区名字
	HouseType       string `json:"house_type,omitempty"` // 户型
	RentType        string `json:"rent_type,omitempty"`  // 出租类型
	Price           int    `json:"price"`                // 房源价格
	CommissionPrice int    `json:"commission_price"`     // 房源返佣
	Feature         string `json:"feature" `             // 房源特色
	IsTeamHouse     string `json:"is_team_house"`        // 是否团队房源
	Status          string `json:"status" `              // 房源状态
	LastUpdate      string `json:"last_update"`          // 最后更新
}

func FromDeviceDB(entity *house.Resource) *ResourceBuffer {
	return &ResourceBuffer{
		ID:              strconv.Itoa(int(entity.ID)),
		HouseId:         strconv.Itoa(int(entity.ID)),
		Owner:           strconv.Itoa(int(entity.Owner)),
		City:            entity.City,
		Districts:       entity.Districts,
		DistrictIds:     entity.DistrictIds,
		XiaoquId:        entity.XiaoquId,
		Xiaoqu:          entity.Xiaoqu,
		HouseType:       entity.HouseType,
		RentType:        entity.RentType,
		Price:           entity.Price,
		CommissionPrice: entity.CommissionPrice,
		Feature:         entity.Feature,
		IsTeamHouse:     boolToSearchValue(entity.IsTeamHouse),
		Status:          entity.Status,
		LastUpdate:      entity.UpdatedLastAt.Format(time.RFC3339),
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
			switch v := value.(type) {
			case int:
				result.Price = v
			case float64:
				result.Price = int(v)
			}
		case "commission_price":
			switch v := value.(type) {
			case int:
				result.CommissionPrice = v
			case float64:
				result.CommissionPrice = int(v)
			}
		case "feature":
			result.Feature = value.(string)
		case "is_team_house":
			result.IsTeamHouse = value.(string)
		case "status":
			result.Status = value.(string)
		case "house_id":
			result.HouseId = value.(string)
		case "owner":
			switch v := value.(type) {
			case string:
				result.Owner = v
			case float64:
				result.Owner = strconv.Itoa(int(v))
			case int:
				result.Owner = strconv.Itoa(v)
			}
		}
	}

	return result
}

func (d *ResourceBuffer) ToData() map[string]interface{} {
	return map[string]interface{}{
		"house_id":         d.HouseId,
		"owner":            d.Owner,
		"city":             d.City,
		"districts":        d.Districts,
		"district_ids":     d.DistrictIds,
		"xiaoqu_id":        d.XiaoquId,
		"xiaoqu":           d.Xiaoqu,
		"house_type":       d.HouseType,
		"rent_type":        d.RentType,
		"price":            d.Price,
		"commission_price": d.CommissionPrice,
		"feature":          d.Feature,
		"is_team_house":    d.IsTeamHouse,
		"status":           d.Status,
		"last_update":      d.LastUpdate,
	}
}

func boolToSearchValue(v bool) string {
	if v {
		return "1"
	}
	return "0"
}
