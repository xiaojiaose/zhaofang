package human

import (
	"errors"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/human"
	request2 "github.com/flipped-aurora/gin-vue-admin/server/model/human/request"

	"gorm.io/gorm"
)

type HumanService struct{}

func (s *HumanService) UpdateHuman(human2 human.Human) error {

	if errors.Is(global.GVA_DB.Where("id = ?", human2.ID).First(&human.Human{}).Updates(human2).Error, gorm.ErrRecordNotFound) {
		return global.GVA_DB.Create(&human2).Error
	}
	return nil
}

func (s *HumanService) UpdateWxUser(human2 *human.WxUser) error {
	if errors.Is(global.GVA_DB.Where("id = ?", human2.ID).First(&human.WxUser{}).Updates(human2).Error, gorm.ErrRecordNotFound) {
		return global.GVA_DB.Create(&human2).Error
	}
	return nil
}

func (s *HumanService) FindHumanByMobile(mobile string) (h *human.Human) {

	var hh human.Human
	if errors.Is(global.GVA_DB.Where("mobile = ?", mobile).First(&hh).Error, gorm.ErrRecordNotFound) {
		return
	}
	return &hh
}

func (s *HumanService) FindUserByMobile(mobile string) (h *human.WxUser) {
	var hh human.WxUser
	if errors.Is(global.GVA_DB.Where("mobile = ?", mobile).First(&hh).Error, gorm.ErrRecordNotFound) {
		return
	}
	return &hh
}

func (s *HumanService) FindUserByOpenid(openid string) (h *human.WxUser) {
	var hh human.WxUser
	if errors.Is(global.GVA_DB.Where("openid = ?", openid).First(&hh).Error, gorm.ErrRecordNotFound) {
		return
	}
	return &hh
}

func (s *HumanService) FindHumanById(id uint) (h *human.Human, err error) {

	var hh human.Human

	err = global.GVA_DB.Where("id = ?", id).First(&hh).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		global.GVA_LOG.Warn(fmt.Sprintf("not found id: %s", id))
		return
	}
	return &hh, nil
}

func (s *HumanService) CreateWxUser(user *human.WxUser) error {

	var xq human.WxUser
	tx := global.GVA_DB.Where("openid = ?", user.Openid).First(&xq)

	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return global.GVA_DB.Save(user).Error
	}
	tx.Updates(user)

	return tx.Error

}

func (s *HumanService) GetPage(api request2.SearchHumanParams, info request.PageInfo, order string, desc bool) (list interface{}, total int64, err error) {

	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&human.Human{})
	var apiList []human.Human

	if api.Name != "" {
		db = db.Where("name = ?", api.Name)
	}

	if api.Mobile != "" {
		db = db.Where("mobile = ?", api.Mobile)
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
