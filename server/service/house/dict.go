package house

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/house"
)

type DictService struct {
}

func (d DictService) GetBuilding(communityId string) (resources []house.DictBuilding, err error) {
	err = global.GVA_DB.Model(&house.DictBuilding{}).Where("community_id = ? ", communityId).Find(&resources).Error
	return
}

func (d DictService) GetUnit(buildingId string) (resources []house.DictUnit, err error) {
	err = global.GVA_DB.Model(&house.DictUnit{}).Where("building_open_id = ? ", buildingId).Find(&resources).Error
	return
}

func (d DictService) GetHouse(unitId string) (resources []house.DictHouse, err error) {
	err = global.GVA_DB.Model(&house.DictHouse{}).Where("unit_open_id = ? ", unitId).Find(&resources).Error
	return
}
