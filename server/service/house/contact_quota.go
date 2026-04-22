package house

import (
	"errors"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/house"
	response2 "github.com/flipped-aurora/gin-vue-admin/server/model/house/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"gorm.io/gorm"
	"time"
)

type ContactQuotaService struct{}

func (service *ContactQuotaService) GrantByPhone(operatorID uint, req request.ContactQuotaCreate) error {
	// 联系方式查看次数挂在用户账号上，而不是挂在某一套房源上。
	// 后台按经纪人手机号加次数，本质上是在给该账号做额度充值。
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
	// 房东房源查看联系方式时按“查看一次扣一次”处理，
	// 同时记录消耗流水，方便后台追踪剩余次数和使用次数。
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var user system.SysUser
		if err := tx.Where("id = ?", userID).First(&user).Error; err != nil {
			return err
		}
		if user.ContactViewQuotaTotal <= 0 {
			return errors.New("查看次数不足，请联系客服购买")
		}
		user.ContactViewQuotaTotal--
		if err := tx.Model(&system.SysUser{}).Where("id = ?", user.ID).Update("contact_view_quota_total", user.ContactViewQuotaTotal).Error; err != nil {
			return err
		}
		if err := tx.Create(&house.ContactQuotaLog{
			UserID:            user.ID,
			OperatorID:        userID,
			UserPhone:         user.Phone,
			Action:            house.ContactQuotaActionUse,
			ChangeAmount:      -1,
			UsedAmount:        1,
			RemainingAmount:   user.ContactViewQuotaTotal,
			Remark:            remark,
			RelatedResourceID: resourceID,
		}).Error; err != nil {
			return err
		}
		return service.createLandlordContactViewRecord(tx, userID, resourceID)
	})
}

func (service *ContactQuotaService) GetPage(req request.ContactQuotaSearch) (list []house.ContactQuotaLog, total int64, err error) {
	// 这里返回的是流水列表，不是用户汇总表。
	// 这样后台能同时看到充值记录和使用记录。
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

func (service *ContactQuotaService) GetLandlordContactViewPage(req request.LandlordContactViewSearch) (list []response2.LandlordContactViewResponse, total int64, err error) {
	db := global.GVA_DB.Model(&house.LandlordContactView{})
	if req.Status != "" {
		db = db.Where("status = ?", req.Status)
	}
	if req.PublisherPhone != "" {
		db = db.Where("publisher_phone = ?", req.PublisherPhone)
	}
	if req.ViewerPhone != "" {
		db = db.Where("viewer_phone = ?", req.ViewerPhone)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	var records []house.LandlordContactView
	err = db.Order("id desc").Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&records).Error
	if err != nil {
		return
	}
	for _, item := range records {
		list = append(list, response2.LandlordContactViewResponse{
			ID:              item.ID,
			ResourceID:      item.ResourceID,
			Xiaoqu:          item.Xiaoqu,
			DoorNo:          item.DoorNo,
			PublisherUserID: item.PublisherUserID,
			PublisherName:   item.PublisherName,
			PublisherPhone:  item.PublisherPhone,
			ResourceAt:      item.ResourceCreatedAt.Format(time.DateTime),
			ViewerUserID:    item.ViewerUserID,
			ViewerName:      item.ViewerName,
			ViewerPhone:     item.ViewerPhone,
			ViewAt:          item.CreatedAt.Format(time.DateTime),
			Status:          item.Status,
			LastOperateAt:   time.UnixMilli(item.LastOperatedAtUnixMilli).Format(time.DateTime),
		})
	}
	return
}

func (service *ContactQuotaService) UpdateLandlordContactViewStatus(req request.LandlordContactViewAction) error {
	statusMap := map[string]string{
		"processing": house.LandlordViewStatusProcessing,
		"approve":    house.LandlordViewStatusApproved,
		"paid":       house.LandlordViewStatusPaid,
		"reject":     house.LandlordViewStatusRejected,
	}
	status, ok := statusMap[req.Action]
	if !ok {
		return errors.New("不支持的操作")
	}
	return global.GVA_DB.Model(&house.LandlordContactView{}).Where("id = ?", req.ID).Updates(map[string]interface{}{
		"status":                      status,
		"last_operated_at_unix_milli": time.Now().UnixMilli(),
	}).Error
}

func (service *ContactQuotaService) createLandlordContactViewRecord(tx *gorm.DB, viewerID, resourceID uint) error {
	var resource house.Resource
	if err := tx.Where("id = ?", resourceID).First(&resource).Error; err != nil {
		return err
	}
	if resource.HouseType != "房东房源" {
		// 只有房东房源联系方式查看才会进入该记录表。
		return nil
	}
	var publisher system.SysUser
	if err := tx.Where("id = ?", resource.Owner).First(&publisher).Error; err != nil {
		return err
	}
	var viewer system.SysUser
	if err := tx.Where("id = ?", viewerID).First(&viewer).Error; err != nil {
		return err
	}
	now := time.Now().UnixMilli()
	return tx.Create(&house.LandlordContactView{
		ResourceID:              resource.ID,
		Xiaoqu:                  resource.Xiaoqu,
		DoorNo:                  resource.DoorNo,
		ResourceCreatedAt:       resource.CreatedAt,
		PublisherUserID:         publisher.ID,
		PublisherName:           preferredUserName(publisher),
		PublisherPhone:          publisher.Phone,
		ViewerUserID:            viewer.ID,
		ViewerName:              preferredUserName(viewer),
		ViewerPhone:             viewer.Phone,
		Status:                  house.LandlordViewStatusPending,
		LastOperatedAtUnixMilli: now,
	}).Error
}

func preferredUserName(user system.SysUser) string {
	if user.WxNickName != "" {
		return user.WxNickName
	}
	if user.NickName != "" {
		return user.NickName
	}
	return user.Username
}
