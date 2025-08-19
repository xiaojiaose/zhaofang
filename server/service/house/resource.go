package house

import (
	"context"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/house"
	"github.com/flipped-aurora/gin-vue-admin/server/model/search"
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
		db = db.Where("user_id = ?", userId)
	}

	if len(appStatus) > 0 {
		db = db.Where("approval_status = ?", appStatus)
	}

	if len(status) > 0 {
		db = db.Where("status = ?", status)
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
