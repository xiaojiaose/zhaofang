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
	// 出房有礼的“最近联系过的人”并不是单独维护一张表，
	// 而是直接复用 /center/house/mobile 产生的 click 行为记录。
	// 这里先按当前用户聚合最近点击过的房源，再反查房源发布人信息返回给前端。
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
	// 申请单在创建时就把申请人和发布人的手机号/微信号快照落库，
	// 这样即使后续用户资料被修改，历史审核记录仍然能保持当时的展示数据。
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
	// 进入列表前先做一次兜底自动流转，保证页面看到的状态尽量接近最终业务状态。
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
	// 后台只处理“发布人已经确认”的申请，
	// 所以这里默认过滤掉还没到后台审核阶段的数据。
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
		// 发布人确认后，申请单才正式进入后台审核池。
		return global.GVA_DB.Model(&house.RewardApplication{}).
			Where("id = ? AND publisher_user_id = ?", req.ID, userID).
			Updates(map[string]interface{}{
				"publisher_confirm_status":    house.RewardPublisherApproved,
				"audit_status":                house.RewardAuditPending,
				"last_operated_at_unix_milli": now,
			}).Error
	case "reject":
		// 拒绝后不进入后台审核，所以这里只改发布人确认状态。
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
	// 前端操作文案和数据库里的审核状态并不完全一致，
	// 这里集中做一次映射，避免在多个 handler 中散落同样的判断逻辑。
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
	// 业务要求：发布人确认后，如果后台 48 小时内没有处理，
	// 系统默认流转为“审通过待发放”。
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
	// 审核列表需要展示小区和户室号，这些字段来源仍然是房源表，
	// 所以这里在返回前做一次轻量拼装。
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
	// 数据中心的手机号维度统计，是先用房源手机号反查房源，
	// 再把这些房源的浏览、分享、点击和奖励申请汇总出来。
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
