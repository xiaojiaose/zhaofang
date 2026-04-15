package house

import (
	"errors"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/house"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
)

type ContactQuotaService struct{}

func (service *ContactQuotaService) GrantByPhone(operatorID uint, req request.ContactQuotaCreate) error {
	var user system.SysUser
	if err := global.GVA_DB.Where("phone = ?", req.UserPhone).First(&user).Error; err != nil {
		return err
	}
	user.ContactViewQuotaTotal += req.Amount
	if err := global.GVA_DB.Model(&system.SysUser{}).Where("id = ?", user.ID).Update("contact_view_quota_total", user.ContactViewQuotaTotal).Error; err != nil {
		return err
	}
	return global.GVA_DB.Create(&house.ContactQuotaLog{
		UserID:          user.ID,
		OperatorID:      operatorID,
		UserPhone:       user.Phone,
		Action:          house.ContactQuotaActionGrant,
		ChangeAmount:    req.Amount,
		RemainingAmount: user.ContactViewQuotaTotal,
		Remark:          req.Remark,
	}).Error
}

func (service *ContactQuotaService) Consume(userID, resourceID uint, remark string) error {
	var user system.SysUser
	if err := global.GVA_DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return err
	}
	if user.ContactViewQuotaTotal <= 0 {
		return errors.New("查看次数不足，请联系客服购买")
	}
	user.ContactViewQuotaTotal--
	if err := global.GVA_DB.Model(&system.SysUser{}).Where("id = ?", user.ID).Update("contact_view_quota_total", user.ContactViewQuotaTotal).Error; err != nil {
		return err
	}
	return global.GVA_DB.Create(&house.ContactQuotaLog{
		UserID:            user.ID,
		OperatorID:        userID,
		UserPhone:         user.Phone,
		Action:            house.ContactQuotaActionUse,
		ChangeAmount:      -1,
		UsedAmount:        1,
		RemainingAmount:   user.ContactViewQuotaTotal,
		Remark:            remark,
		RelatedResourceID: resourceID,
	}).Error
}

func (service *ContactQuotaService) GetPage(req request.ContactQuotaSearch) (list []house.ContactQuotaLog, total int64, err error) {
	db := global.GVA_DB.Model(&house.ContactQuotaLog{})
	if req.UserPhone != "" {
		db = db.Where("user_phone = ?", req.UserPhone)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Order("id desc").Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&list).Error
	return
}
