package house

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/house"
)

type FavoriteService struct {
}

func (service *FavoriteService) CreateOrUpdate(favorite *house.Favorite) (err error) {
	err = global.GVA_DB.Where("id = ?", favorite.ID).First(&house.Favorite{}).Updates(&favorite).Error
	if err != nil && err.Error() == "record not found" {
		err = global.GVA_DB.Create(favorite).Error
	}

	return
}

func (service *FavoriteService) Delete(userId, resource uint) (err error) {
	err = global.GVA_DB.Where("user_id = ? and resource_id = ?", userId, resource).Delete(&house.Favorite{}).Error
	return
}

func (service *FavoriteService) GetPage(userId uint, info request.PageInfo, order string, desc bool) (list interface{}, total int64, err error) {

	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&house.Favorite{})
	var apiList []house.Favorite

	db = db.Where("user_id = ?", userId)

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
