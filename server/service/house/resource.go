package house

import (
	"context"
	"errors"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/house"
	"github.com/flipped-aurora/gin-vue-admin/server/model/search"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type ResourceService struct{}

func (service *ResourceService) FilterOptions() (list map[string]map[string]string, err error) {
	list = make(map[string]map[string]string)
	list["rentType"] = map[string]string{"1": "整租", "2": "分整租", "3": "合租", "4": "房东房源"}
	list["houseType"] = map[string]string{"1": "1居", "2": "2居", "3": "3居", "4": "4居+", "5": "开间", "6": "主卧", "7": "次卧", "8": "暗间", "9": "房东房源"}
	list["feature"] = map[string]string{"1": "可短租", "2": "包物业", "3": "南北通透", "4": "全南项", "5": "协助对接房东", "6": "可带看分佣"}
	list["price"] = map[string]string{"1": "500以下", "2": "500-1000元", "3": "1000-1500元", "4": "1500-2000元", "5": "2000-2500元", "6": "2500-3000元", "7": "3000元以上"}

	return
}

func (service *ResourceService) FilterOptions1() (list []request.RentType, err error) {
	list = append(list,
		request.RentType{
			Name:      "整租",
			HouseType: []string{"1居", "2居", "3居", "4居+", "开间"}, // 可短租，有电梯，可注册办公，密码看房、包物业
			Feature:   []string{"可短租", "包物业", "有电梯", "密码看房", "可办公注册", "协助对接房东", "可带看分佣"},
		},
		request.RentType{
			Name:      "分整租",
			HouseType: []string{"1居", "2居", "3居", "4居+"}, // 可短租，有电梯、有原卫、有阳台，有燃气，朝南
			Feature:   []string{"带阳台", "有电梯", "有原卫", "朝南", "有燃气", "可短租", "协助对接房东", "可带看分佣"},
		},
		request.RentType{
			Name:      "合租",
			HouseType: []string{"2居", "3居", "4居+"}, // 可短租、有电梯，有独卫，有阳台，可做饭、纯女生
			Feature:   []string{"带阳台", "可短租", "有电梯", "有独卫", "可做饭", "纯女生", "协助对接房东", "可带看分佣"},
		},
		request.RentType{
			Name:      "房东房源",
			HouseType: []string{"房东房源"},
			Feature:   []string{"协助对接房东", "可带看分佣"},
		},
	)

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
	if err = service.fillUserRelatedFields(resource); err != nil {
		return err
	}
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

	if imgs, ok := resource.Attachments["house"]; ok && len(imgs) > 0 {
		resource.HasPic = true
	}

	resource.UpdatedLastAt = time.Now()
	if resource.Status == "" {
		resource.Status = "待出租"
	}
	if resource.Status == "待出租" {
		if err = service.ensurePublishQuota(resource.Owner, resource.ID); err != nil {
			return err
		}
	}
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
	err = global.GVA_DB.Where("id = ? and owner = ?", id, owner).Delete(&house.Resource{}).Error
	if err == nil {
		err = global.Gva_ResourceSearch.Del(context.Background(), strconv.Itoa(int(id)))
		if err != nil {
			return err
		}
	}

	return
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
	if value == "待出租" {
		for _, id := range ids {
			info, e := service.GetInfo(id)
			if e != nil {
				return e
			}
			if err = service.ensurePublishQuota(info.Owner, info.ID); err != nil {
				return err
			}
		}
	}
	err = global.GVA_DB.Model(&house.Resource{}).Where("id in ? ", ids).Update("status", value).Error
	if err == nil {
		for _, id := range ids {
			info, e := service.GetInfo(id)
			if e == nil {
				err = global.Gva_ResourceSearch.Add(context.Background(), *search.FromDeviceDB(info))
			}
		}
	}
	return
}
func (service *ResourceService) SetApprovalStatus(ids []uint, value string) (err error) {
	err = global.GVA_DB.Model(&house.Resource{}).Where("id in ? ", ids).Updates(map[string]interface{}{"approval_status": value, "status": "待出租"}).Error
	if err == nil {
		for _, id := range ids {
			info, e := service.GetInfo(id)
			if e == nil {
				err = global.Gva_ResourceSearch.Add(context.Background(), *search.FromDeviceDB(info))
			}
		}
	}
	return
}

func (service *ResourceService) GetListByIds(ids []uint) (resources []*house.Resource, err error) {
	err = global.GVA_DB.Model(&house.Resource{}).Where("id in ?", ids).Order("updated_last_at desc").Find(&resources).Error
	return
}

func (service *ResourceService) GetPage(xiaoquId, userId uint, appStatus string, status string, info request.PageInfo, order string, desc bool, Other request.SearchOther) (list interface{}, total int64, err error) {

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

	if Other.Phone != "" {
		db = db.Where("phone = ?", Other.Phone)
	}
	if Other.IsTeamHouse == "true" {
		db = db.Where("is_team_house = ?", true)
	}
	if Other.IsTeamHouse == "false" {
		db = db.Where("is_team_house = ?", false)
	}
	if Other.HouseType != "" {
		db = db.Where("house_type = ?", Other.HouseType)
	}
	if Other.HasCommission == "true" {
		db = db.Where("commission_price > 0")
	}
	switch Other.HasPic {
	case "true":
		db = db.Where("has_pic = ?", true)
	case "false":
		db = db.Where("has_pic = ?", false)
	}

	if Other.RentType != "" {
		db = db.Where("rent_type = ?", Other.RentType)
	}

	if Other.UpdatedAtLast.IsZero() == false && Other.UpdatedAtStart.IsZero() == false {
		db = db.Where("updated_last_at > ? and updated_last_at < ?", Other.UpdatedAtStart, Other.UpdatedAtLast)
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

func (service *ResourceService) CountOnShelfByUser(userID uint, excludeIDs ...uint) (count int64, err error) {
	db := global.GVA_DB.Model(&house.Resource{}).Where("owner = ? AND status = ?", userID, "待出租")
	if len(excludeIDs) > 0 {
		db = db.Where("id NOT IN ?", excludeIDs)
	}
	err = db.Count(&count).Error
	return
}

func (service *ResourceService) RefreshUserTeamHouses(userID uint, isTeam bool) error {
	return global.GVA_DB.Model(&house.Resource{}).Where("owner = ?", userID).Update("is_team_house", isTeam).Error
}

func (service *ResourceService) fillUserRelatedFields(resource *house.Resource) error {
	if resource.Owner == 0 {
		return nil
	}
	var user system.SysUser
	if err := global.GVA_DB.Where("id = ?", resource.Owner).First(&user).Error; err != nil {
		return err
	}
	resource.IsTeamHouse = user.IsFindHouseSupermarket
	if resource.Phone == "" {
		resource.Phone = user.Phone
	}
	return nil
}

func (service *ResourceService) ensurePublishQuota(userID uint, resourceID uint) error {
	if userID == 0 {
		return nil
	}
	var user system.SysUser
	if err := global.GVA_DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return err
	}
	if user.PublishQuotaTotal <= 0 {
		return errors.New("当前账号未开通上架权限")
	}
	count, err := service.CountOnShelfByUser(userID, resourceID)
	if err != nil {
		return err
	}
	if int(count) >= user.PublishQuotaTotal {
		return fmt.Errorf("已达到最大上架数量：%d", user.PublishQuotaTotal)
	}
	return nil
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
