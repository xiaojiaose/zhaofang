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

	list["feature"] = map[string]string{"1": "可短租", "2": "包物业", "3": "南北通透", "4": "全南项"}
	list["price"] = map[string]string{"1": "500以下", "2": "500-1000元", "3": "1000-1500元", "4": "1500-2000元", "5": "2000-2500元", "6": "2500-3000元", "7": "3000元以上"}

	return
}

func (service *ResourceService) FilterOptions1() (list []request.RentType, err error) {
	list = append(list,
		request.RentType{
			Name:      "整租",
			HouseType: []string{"1居", "2居", "3居", "4居+", "开间"}, // 可短租，有电梯，可注册办公，密码看房、包物业
			Feature:   []string{"可短租", "包物业", "有电梯", "密码看房", "可办公注册"},
		},
		request.RentType{
			Name:      "分整租",
			HouseType: []string{"1居", "2居", "3居", "4居+"}, // 可短租，有电梯、有原卫、有阳台，有燃气，朝南
			Feature:   []string{"带阳台", "有电梯", "有原卫", "朝南", "有燃气", "可短租"},
		},
		request.RentType{
			Name:      "合租",
			HouseType: []string{"2居", "3居", "4居+"}, // 可短租、有电梯，有独卫，有阳台，可做饭、纯女生
			Feature:   []string{"带阳台", "可短租", "有电梯", "有独卫", "可做饭", "纯女生"},
		},
		request.RentType{
			Name:      "房东房源",
			HouseType: []string{"1居", "2居", "3居", "4居+", "开间"}, // 可短租，有电梯，可注册办公，密码看房、包物业
			Feature:   []string{"可短租", "包物业", "有电梯", "密码看房", "可办公注册"},
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
	// 房源的"团队房源标识"和默认联系方式都依赖发布人资料，
	// 每次保存前都重新补齐，避免前端漏传或传了旧值。
	if err = service.fillUserRelatedFields(resource); err != nil {
		return err
	}
	var doorNo string
	if resource.BuildingId != "" {
		// 如果接的是字典式楼栋/单元/房号，就在保存前拼出可展示的 DoorNo；
		// 这样前台和搜索索引都能直接复用统一的门牌字段。
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
	// 只有新建房源（ID==0）且要上架时才校验额度，编辑已有房源不校验。
	if resource.ID == 0 && resource.Status == "待出租" {
		if err = service.ensurePublishQuota(resource.Owner, resource.ID); err != nil {
			return err
		}
	}
	err = global.GVA_DB.Where("id = ?", resource.ID).First(&house.Resource{}).Updates(&resource).Error
	if err != nil && err.Error() == "record not found" {
		err = global.GVA_DB.Create(resource).Error
	}
	if err == nil {
		// 房源库和地图检索都依赖 ES/Zinc 索引，所以每次保存成功后都立即同步搜索索引。
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

func (service *ResourceService) SyncIndexByIDs(ids []uint) error {
	// 任何"先改 DB，再查 Zinc"的链路都可能遇到索引延迟。
	// 这里提供统一的按房源 ID 回写入口，确保状态/团队标识等关键字段能及时落到 Zinc。
	if len(ids) == 0 {
		return nil
	}
	var resources []house.Resource
	if err := global.GVA_DB.Where("id IN ?", ids).Find(&resources).Error; err != nil {
		return err
	}
	for _, item := range resources {
		if err := global.Gva_ResourceSearch.Add(context.Background(), *search.FromDeviceDB(&item)); err != nil {
			return err
		}
	}
	return nil
}

func (service *ResourceService) SyncIndexByOwner(owner uint) error {
	// 用户维度变更（例如找房超市标识）会影响其名下全部房源在地图中的可见性。
	// 该方法用于"按 owner 一次性全量回写索引"。
	if owner == 0 {
		return nil
	}
	var resources []house.Resource
	if err := global.GVA_DB.Where("owner = ?", owner).Find(&resources).Error; err != nil {
		return err
	}
	for _, item := range resources {
		if err := global.Gva_ResourceSearch.Add(context.Background(), *search.FromDeviceDB(&item)); err != nil {
			return err
		}
	}
	return nil
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
	err = global.GVA_DB.Model(&house.Resource{}).Where("id = ? AND "+field+" > 0", id).UpdateColumn(field, gorm.Expr(fmt.Sprintf("%s - ?", field), 1)).Error
	return
}

func (service *ResourceService) SetState(ids []uint, value string) (err error) {
	if value == "待出租" {
		// 批量上架时逐条校验额度，
		// 这样可以复用和单条保存一致的限制逻辑。
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
func (service *ResourceService) SetApprovalStatus(ids []uint, approvalStatus string, status string) (err error) {
	updates := make(map[string]interface{})
	updates["approval_status"] = approvalStatus
	if status != "" {
		updates["status"] = status
	}
	err = global.GVA_DB.Model(&house.Resource{}).Where("id in ?", ids).Updates(updates).Error
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

func (service *ResourceService) GetListByIdsSafe(ids []uint, status string, allowTeam bool) (resources []*house.Resource, err error) {
	// 通过搜索引擎拿到 id 后，数据库层再做一次权限和状态兜底过滤，
	// 防止索引延迟导致"已下架/无权限团队房源"被错误返回到小程序。
	if len(ids) == 0 {
		return []*house.Resource{}, nil
	}
	db := global.GVA_DB.Model(&house.Resource{}).Where("id in ?", ids)
	if status != "" {
		db = db.Where("status = ?", status)
	}
	if !allowTeam {
		// 兼容历史数据：老数据可能没有回填 is_team_house，按空值视作非团队房源。
		db = db.Where("(is_team_house = ? OR is_team_house IS NULL)", false)
	}

	var rows []*house.Resource
	if err = db.Find(&rows).Error; err != nil {
		return
	}

	// 回表后按 Zinc 返回的 id 顺序还原，避免数据库 in 查询打乱展示顺序。
	rowMap := make(map[uint]*house.Resource, len(rows))
	for _, row := range rows {
		rowMap[row.ID] = row
	}
	resources = make([]*house.Resource, 0, len(rows))
	for _, id := range ids {
		if row, ok := rowMap[id]; ok {
			resources = append(resources, row)
		}
	}
	return
}

func (service *ResourceService) GetPage(xiaoquId, userId uint, appStatus string, status string, info request.PageInfo, order string, desc bool, Other request.SearchOther) (list interface{}, total int64, err error) {
	// 这里是后台和若干业务列表共用的数据库查询入口；
	// ES 地图检索走另一条链路，但"我的房源""后台审核列表"等仍然依赖这里的条件拼装。
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
	// 团队房源字段虽然来源于用户标识，但已经冗余到房源表，
	// 所以这里可以直接按房源维度过滤，不需要联表查用户。
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
		db = db.Where("updated_last_at >= ? and updated_last_at < ?", Other.UpdatedAtStart, Other.UpdatedAtLast)
	}

	if Other.Status != "" {
		db = db.Where("status = ?", Other.Status)
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
			// 追加 id 作为次级排序键，避免主排序值重复时分页顺序不稳定导致翻页结果重复/漏数据。
			err = db.Order(OrderStr).Order("id desc").Find(&apiList).Error
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
	// 团队房源虽然来源于用户的找房超市标识，
	// 但最终在房源表上做了冗余存储，方便地图、ES 和后台列表直接查询。
	if err := global.GVA_DB.Model(&house.Resource{}).Where("owner = ?", userID).Update("is_team_house", isTeam).Error; err != nil {
		return err
	}
	// 关键修复：用户标识变更后，立即同步其名下全部房源到 Zinc，
	// 避免"DB 已改、索引未改"导致无权限用户仍看到团队房源。
	return service.SyncIndexByOwner(userID)
}

func (service *ResourceService) fillUserRelatedFields(resource *house.Resource) error {
	if resource.Owner == 0 {
		return nil
	}
	var user system.SysUser
	if err := global.GVA_DB.Where("id = ?", resource.Owner).First(&user).Error; err != nil {
		return err
	}
	// 团队房源标识始终以用户资料为准，不信任前端直接传的值。
	resource.IsTeamHouse = user.IsFindHouseSupermarket
	// 联系手机号默认跟随发布人账号，除非业务侧明确覆盖。
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
	// 编辑已上架房源时需要把自己排除掉，否则会把"保存已有房源"误判成超额。
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
			// 追加 id 作为次级排序键，避免主排序值重复时分页顺序不稳定导致翻页结果重复/漏数据。
			err = db.Order(OrderStr).Order("id desc").Find(&apiList).Error
		} else {
			err = db.Order("id desc").Find(&apiList).Error
		}
	}
	return apiList, total, err
}

func (service *ResourceService) RebuildAllResourceIndex(pageSize int) (count int, err error) {
	// 一次性全量重建：用于清理历史脏索引。
	// 按主键分页扫描 DB 全量房源，再逐条 Add 到 Zinc 覆盖旧文档。
	if pageSize <= 0 {
		pageSize = 1000
	}
	page := 1
	for {
		var rows []house.Resource
		e := global.GVA_DB.
			Model(&house.Resource{}).
			Order("id asc").
			Limit(pageSize).
			Offset((page - 1) * pageSize).
			Find(&rows).Error
		if e != nil {
			return count, e
		}
		if len(rows) == 0 {
			break
		}
		for _, item := range rows {
			if e = global.Gva_ResourceSearch.Add(context.Background(), *search.FromDeviceDB(&item)); e != nil {
				return count, e
			}
			count++
		}
		if len(rows) < pageSize {
			break
		}
		page++
	}
	return count, nil
}
