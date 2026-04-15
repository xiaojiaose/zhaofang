package house

import (
	"errors"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/house"
	response2 "github.com/flipped-aurora/gin-vue-admin/server/model/house/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/search"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
)

type RewardService struct{}

func (service *RewardService) RecentContacts(userID uint) (list []response2.RewardRecentContact, err error) {
	var rows []struct {
		ResourceID uint
		MaxDate    time.Time
	}
	err = global.GVA_DB.Table("visit_daily").
		Select("resource_id, max(date) as max_date").
		Where("user_id = ? AND click > 0", userID).
		Group("resource_id").
		Order("max_date desc").
		Scan(&rows).Error
	if err != nil {
		return
	}
	for _, row := range rows {
		var resource house.Resource
		if e := global.GVA_DB.Where("id = ?", row.ResourceID).First(&resource).Error; e != nil {
			continue
		}
		var publisher system.SysUser
		_ = global.GVA_DB.Where("id = ?", resource.Owner).First(&publisher).Error
		list = append(list, response2.RewardRecentContact{
			ResourceID:          resource.ID,
			Xiaoqu:              resource.Xiaoqu,
			DoorNo:              resource.DoorNo,
			PublisherUserID:     publisher.ID,
			PublisherPhone:      publisher.Phone,
			PublisherWxNo:       publisher.WxNo,
			PublisherWxNickName: publisher.WxNickName,
			LastContactAt:       row.MaxDate.Format(time.DateTime),
		})
	}
	return
}

func (service *RewardService) Create(userID uint, req request.RewardApplicationCreate) error {
	var resource house.Resource
	if err := global.GVA_DB.Where("id = ?", req.ResourceID).First(&resource).Error; err != nil {
		return err
	}
	if resource.Owner == userID {
		return errors.New("不能对自己的房源申请出房有礼")
	}
	var applyUser system.SysUser
	var publisher system.SysUser
	if err := global.GVA_DB.Where("id = ?", userID).First(&applyUser).Error; err != nil {
		return err
	}
	if err := global.GVA_DB.Where("id = ?", resource.Owner).First(&publisher).Error; err != nil {
		return err
	}
	now := time.Now().UnixMilli()
	return global.GVA_DB.Create(&house.RewardApplication{
		ResourceID:              resource.ID,
		ApplyUserID:             userID,
		PublisherUserID:         publisher.ID,
		ApplyUserPhone:          applyUser.Phone,
		ApplyUserWxNo:           applyUser.WxNo,
		PublisherUserPhone:      publisher.Phone,
		PublisherUserWxNo:       publisher.WxNo,
		Remark:                  req.Remark,
		PublisherConfirmStatus:  house.RewardPublisherPending,
		AuditStatus:             house.RewardAuditNone,
		LastOperatedAtUnixMilli: now,
	}).Error
}

func (service *RewardService) GetPageForPublisher(userID uint, req request.RewardApplicationSearch) (list []response2.RewardApplicationResponse, total int64, err error) {
	service.AutoApproveExpired()
	db := global.GVA_DB.Model(&house.RewardApplication{}).Where("publisher_user_id = ?", userID)
	if req.PublisherConfirmStatus != "" {
		db = db.Where("publisher_confirm_status = ?", req.PublisherConfirmStatus)
	}
	if req.AuditStatus != "" {
		db = db.Where("audit_status = ?", req.AuditStatus)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	var apps []house.RewardApplication
	err = db.Order("id desc").Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&apps).Error
	if err != nil {
		return
	}
	list = service.attachResourceInfo(apps)
	return
}

func (service *RewardService) GetPageForAdmin(req request.RewardApplicationSearch) (list []response2.RewardApplicationResponse, total int64, err error) {
	service.AutoApproveExpired()
	db := global.GVA_DB.Model(&house.RewardApplication{}).Where("publisher_confirm_status = ?", house.RewardPublisherApproved)
	if req.AuditStatus != "" {
		db = db.Where("audit_status = ?", req.AuditStatus)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	var apps []house.RewardApplication
	err = db.Order("id desc").Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&apps).Error
	if err != nil {
		return
	}
	list = service.attachResourceInfo(apps)
	return
}

func (service *RewardService) PublisherAction(userID uint, req request.RewardApplicationAction) error {
	now := time.Now().UnixMilli()
	switch req.Action {
	case "approve":
		return global.GVA_DB.Model(&house.RewardApplication{}).
			Where("id = ? AND publisher_user_id = ?", req.ID, userID).
			Updates(map[string]interface{}{
				"publisher_confirm_status":    house.RewardPublisherApproved,
				"audit_status":                house.RewardAuditPending,
				"last_operated_at_unix_milli": now,
			}).Error
	case "reject":
		return global.GVA_DB.Model(&house.RewardApplication{}).
			Where("id = ? AND publisher_user_id = ?", req.ID, userID).
			Updates(map[string]interface{}{
				"publisher_confirm_status":    house.RewardPublisherRejected,
				"last_operated_at_unix_milli": now,
			}).Error
	default:
		return errors.New("不支持的操作")
	}
}

func (service *RewardService) AdminAction(req request.RewardApplicationAction) error {
	now := time.Now().UnixMilli()
	statusMap := map[string]string{
		"processing": house.RewardAuditProcessing,
		"approve":    house.RewardAuditApprovedPending,
		"paid":       house.RewardAuditPaid,
		"reject":     house.RewardAuditRejected,
	}
	status, ok := statusMap[req.Action]
	if !ok {
		return errors.New("不支持的操作")
	}
	return global.GVA_DB.Model(&house.RewardApplication{}).Where("id = ?", req.ID).Updates(map[string]interface{}{
		"audit_status":                status,
		"last_operated_at_unix_milli": now,
	}).Error
}

func (service *RewardService) AutoApproveExpired() {
	cutoff := time.Now().Add(-48 * time.Hour).UnixMilli()
	_ = global.GVA_DB.Model(&house.RewardApplication{}).
		Where("publisher_confirm_status = ? AND audit_status IN ?", house.RewardPublisherApproved, []string{house.RewardAuditPending, house.RewardAuditProcessing}).
		Where("last_operated_at_unix_milli < ?", cutoff).
		Update("audit_status", house.RewardAuditApprovedPending).Error
}

func (service *RewardService) CountByDate(start, end time.Time, phone string) (count int64, err error) {
	db := global.GVA_DB.Model(&house.RewardApplication{}).Where("created_at > ? AND created_at < ?", start, end)
	if phone != "" {
		db = db.Where("publisher_user_phone = ? OR apply_user_phone = ?", phone, phone)
	}
	err = db.Count(&count).Error
	return
}

func (service *RewardService) attachResourceInfo(apps []house.RewardApplication) []response2.RewardApplicationResponse {
	list := make([]response2.RewardApplicationResponse, 0, len(apps))
	for _, app := range apps {
		var resource house.Resource
		_ = global.GVA_DB.Where("id = ?", app.ResourceID).First(&resource).Error
		list = append(list, response2.RewardApplicationResponse{
			RewardApplication: app,
			Xiaoqu:            resource.Xiaoqu,
			DoorNo:            resource.DoorNo,
		})
	}
	return list
}

func (service *RewardService) BuildPhoneSummary(phone string) (summary search.StatisData, err error) {
	var resources []house.Resource
	err = global.GVA_DB.Where("phone = ?", phone).Find(&resources).Error
	if err != nil {
		return
	}
	summary.Add = len(resources)
	var resourceIDs []uint
	for _, resource := range resources {
		summary.View += resource.View
		summary.Follow += resource.Follow
		summary.Shared += resource.Shared
		summary.Click += resource.Click
		resourceIDs = append(resourceIDs, resource.ID)
	}
	if len(resourceIDs) == 0 {
		return
	}
	var rewardCount int64
	_ = global.GVA_DB.Model(&house.RewardApplication{}).Where("resource_id IN ?", resourceIDs).Count(&rewardCount).Error
	summary.RewardApply = int(rewardCount)
	return
}
