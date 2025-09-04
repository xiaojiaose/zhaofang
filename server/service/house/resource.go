package house

import (
	"context"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/house"
	"github.com/flipped-aurora/gin-vue-admin/server/model/search"
	"gorm.io/gorm"
	"time"
)

type ResourceService struct{}

func (service *ResourceService) FilterOptions() (list map[string]map[string]string, err error) {
	list = make(map[string]map[string]string)
	list["rentType"] = map[string]string{"1": "整租", "2": "分整租", "3": "合租"}
	list["houseType"] = map[string]string{"1": "1居", "2": "2居", "3": "3居", "4": "4居+", "5": "开间"}
	list["feature"] = map[string]string{"1": "可短租", "2": "包物业", "3": "南北通透", "4": "全南项"}
	list["price"] = map[string]string{"1": "500以下", "2": "500-1000元", "3": "1000-1500元", "4": "1500-2000元", "5": "2000-2500元", "6": "2500-3000元", "7": "3000元以上"}
	return
}

var priceRanges = map[string][]int{
	"1": {0, 500},
	"2": {500, 1000},
	"3": {1000, 1500},
	"4": {1500, 2000},
	"5": {2000, 2500},
	"6": {2500, 3000},
	"7": {3000, 100000},
}

func (service *ResourceService) GetPriceByOption(key string) []int {
	return priceRanges[key]
}

func (service *ResourceService) CreateOrUpdate(resource *house.Resource) (err error) {
	var doorNo string
	if resource.BuildingId != "" {
		var building house.DictBuilding
		err = global.GVA_DB.Model(&house.DictBuilding{}).Where("building_open_id = ? ", resource.BuildingId).First(&building).Error
		if err != nil {
			global.GVA_LOG.Error(err.Error())
			err = nil
		} else {
			doorNo = building.EncryptBuildingName + "号楼 "
		}
	}

	if resource.UnitId != "" {
		var unit house.DictUnit
		err = global.GVA_DB.Model(&house.DictUnit{}).Where("unit_open_id = ? ", resource.UnitId).First(&unit).Error
		if err != nil {
			global.GVA_LOG.Error(err.Error())
			err = nil
		} else {
			doorNo = doorNo + unit.EncryptUnitName + "单元 "
		}

	}

	if resource.HouseId != "" {
		var house1 house.DictHouse
		err = global.GVA_DB.Model(&house.DictHouse{}).Where("house_open_id = ? ", resource.HouseId).First(&house1).Error
		if err != nil {
			global.GVA_LOG.Error(err.Error())
			err = nil
		} else {
			doorNo = doorNo + house1.EncryptHouseName + "室"
		}
	}

	if doorNo != "" {
		resource.DoorNo = doorNo
	}

	resource.UpdatedLastAt = time.Now()
	err = global.GVA_DB.Where("id = ?", resource.ID).First(&house.Resource{}).Updates(&resource).Error
	if err != nil && err.Error() == "record not found" {
		err = global.GVA_DB.Create(resource).Error
	}
	if err == nil {
		err = global.Gva_ResourceSearch.Add(context.Background(), *search.FromDeviceDB(resource))
		if err != nil {
			return err
		}
	}

	return
}

func (service *ResourceService) GetInfo(id uint) (resource *house.Resource, err error) {
	err = global.GVA_DB.Model(&house.Resource{}).Where("id = ? ", id).First(&resource).Error

	return
}

func (service *ResourceService) DelByUser(id, owner uint) (err error) {
	return global.GVA_DB.Where("id = ? and owner = ?", id, owner).Delete(&house.Resource{}).Error
}

func (service *ResourceService) FollowViewClickAdd(id uint, field string) (err error) {
	err = global.GVA_DB.Model(&house.Resource{}).Where("id = ? ", id).UpdateColumn(field, gorm.Expr(fmt.Sprintf("%s + ?", field), 1)).Error
	return
}

func (service *ResourceService) FollowViewClickSub(id uint, field string) (err error) {
	err = global.GVA_DB.Model(&house.Resource{}).Where("id = ? ", id).UpdateColumn(field, gorm.Expr(fmt.Sprintf("%s - ?", field), 1)).Error
	return
}

func (service *ResourceService) SetState(ids []uint, value string) (err error) {
	err = global.GVA_DB.Model(&house.Resource{}).Where("id in ? ", ids).Update("status", value).Error
	return
}
func (service *ResourceService) SetApprovalStatus(ids []uint, value string) (err error) {
	err = global.GVA_DB.Model(&house.Resource{}).Where("id in ? ", ids).Updates(map[string]interface{}{"approval_status": value, "status": "待出租"}).Error
	return
}

func (service *ResourceService) GetListByIds(ids []uint) (resources []*house.Resource, err error) {
	err = global.GVA_DB.Model(&house.Resource{}).Where("id in ? ", ids).Find(&resources).Error
	return
}

func (service *ResourceService) GetPage(xiaoquId, userId uint, appStatus string, status string, info request.PageInfo, order string, desc bool) (list interface{}, total int64, err error) {

	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&house.Resource{})
	var apiList []house.Resource

	if xiaoquId != 0 {
		db = db.Where("xiaoqu_id = ?", xiaoquId)
	}

	if userId != 0 {
		db = db.Where("owner = ?", userId)
	}

	if len(appStatus) > 0 {
		db = db.Where("approval_status = ?", appStatus)
	}

	if len(status) > 0 {
		db = db.Where("status = ?", status)
	}

	if len(info.Keyword) > 0 {
		db = db.Where("door_no like ?", "%"+info.Keyword+"%")
	}

	err = db.Count(&total).Error

	if err != nil {
		return apiList, total, err
	} else {
		db = db.Limit(limit).Offset(offset)
		if order != "" {
			OrderStr := order
			if desc {
				OrderStr = order + " desc"
			}

			err = db.Order(OrderStr).Find(&apiList).Error
		} else {
			err = db.Order("id desc").Find(&apiList).Error
		}
	}
	return apiList, total, err
}

func (service *ResourceService) GetApprovalPage(xiaoquId, userId uint, appStatus string, info request.PageInfo, order string, desc bool) (list interface{}, total int64, err error) {

	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&house.Resource{})
	var apiList []house.Resource

	if xiaoquId != 0 {
		db = db.Where("xiaoqu_id = ?", xiaoquId)
	}

	if userId != 0 {
		db = db.Where("user_id = ?", userId)
	}

	if len(appStatus) > 0 {
		db = db.Where("approval_status = ?", appStatus)
	}

	err = db.Count(&total).Error

	if err != nil {
		return apiList, total, err
	} else {
		db = db.Limit(limit).Offset(offset)
		if order != "" {
			OrderStr := order
			if desc {
				OrderStr = order + " desc"
			}

			err = db.Order(OrderStr).Find(&apiList).Error
		} else {
			err = db.Order("id desc").Find(&apiList).Error
		}
	}
	return apiList, total, err
}
