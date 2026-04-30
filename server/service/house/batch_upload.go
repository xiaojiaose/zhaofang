package house

import (
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common"
	"github.com/flipped-aurora/gin-vue-admin/server/model/house"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

var validRentTypes = map[string]struct{}{
	"整租":   {},
	"分整租":  {},
	"合租":   {},
	"房东房源": {},
}

var validHouseTypes = map[string]struct{}{
	"1居":  {},
	"2居":  {},
	"3居":  {},
	"4居+": {},
	"开间":  {},
	"主卧":  {},
	"次卧":  {},
	"暗间":  {},
	//"房东房源": {},
}

func (service *ResourceService) BatchUpload(userID uint, header *multipart.FileHeader, remark string) (record *house.BatchUploadRecord, err error) {
	// 批量上传按“同步当前账号房源”的语义实现：
	// 上传前先把当前账号已上架房源统一下架，
	// 再把 Excel 里解析出来的房源重新更新/创建为上架状态。
	file, err := header.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, err
	}
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("excel 没有工作表")
	}
	rows, err := f.GetRows(sheets[0])
	if err != nil {
		return nil, err
	}
	if len(rows) < 2 {
		return nil, fmt.Errorf("excel 没有可导入数据")
	}

	headers := normalizeHeaders(rows[0])
	batchNo := fmt.Sprintf("batch-%d", time.Now().UnixNano())
	record = &house.BatchUploadRecord{
		UserID:   userID,
		FileName: header.Filename,
		BatchNo:  batchNo,
		Remark:   strings.TrimSpace(remark),
	}
	if err = global.GVA_DB.Create(record).Error; err != nil {
		return nil, err
	}

	// 这里先统一下架，最终以上传文件中的内容为准。
	// 关键修复：下架前先拿到受影响 id，下架后立即按 id 回写 Zinc，
	// 避免“数据库已下架但索引还是待出租”造成地图误展示。
	var onShelfIDs []uint
	_ = global.GVA_DB.Model(&house.Resource{}).Where("owner = ? AND status = ?", userID, "待出租").Pluck("id", &onShelfIDs).Error
	if len(onShelfIDs) > 0 {
		if err = global.GVA_DB.Model(&house.Resource{}).Where("id IN ?", onShelfIDs).Update("status", "已下架").Error; err != nil {
			return nil, err
		}
		if err = service.SyncIndexByIDs(onShelfIDs); err != nil {
			return nil, err
		}
	}

	var failures []string
	for idx, row := range rows[1:] {
		if isEmptyRow(row) {
			continue
		}
		record.TotalCount++

		values := mapRow(headers, row)
		xiaoquName := values["小区名称"]
		buildingPart := firstNonEmpty(values["楼栋"], values["楼栋号"])
		unitPart := firstNonEmpty(values["单元"], values["单元号"])
		housePart := firstNonEmpty(values["门牌号"], values["户号"])
		doorNo := composeDoorNoFromParts(buildingPart, unitPart, housePart)
		roomCode := values["房间号"]
		price := atoi(values["价格"])
		commission := atoi(values["返佣金额"])
		remarks := values["备注"]
		feature := values["标签"]
		rentType := values["出租类型"]
		houseType := values["房屋类型"]

		if xiaoquName == "" || doorNo == "" {
			record.FailedCount++
			failures = append(failures, fmt.Sprintf("第%d行缺少楼盘名称或户室号", idx+2))
			continue
		}

		var xq system.XiaoQu
		if err = global.GVA_DB.Where("name = ?", xiaoquName).First(&xq).Error; err != nil {
			record.FailedCount++
			failures = append(failures, fmt.Sprintf("第%d行小区不存在:%s", idx+2, xiaoquName))
			continue
		}
		if err = validateDoorNoByDict(xq, doorNo); err != nil {
			record.FailedCount++
			failures = append(failures, fmt.Sprintf("第%d行楼栋房间号不合法:%s", idx+2, err.Error()))
			continue
		}
		buildingID, unitID, houseID, err := resolveDoorDictIDs(xq, doorNo)
		if err != nil {
			record.FailedCount++
			failures = append(failures, fmt.Sprintf("第%d行楼栋房间号字典解析失败:%s", idx+2, err.Error()))
			continue
		}

		if rentType != "" {
			if _, ok := validRentTypes[rentType]; !ok {
				record.FailedCount++
				failures = append(failures, fmt.Sprintf("第%d行出租类型不合法:%s", idx+2, rentType))
				continue
			}
		}
		if houseType != "" {
			if _, ok := validHouseTypes[houseType]; !ok {
				record.FailedCount++
				failures = append(failures, fmt.Sprintf("第%d行房屋类型不合法:%s", idx+2, houseType))
				continue
			}
		}

		var entity house.Resource
		// 批量导入的“同一条房源”判定：同账号 + 小区 + 楼栋 + 单元 + 门牌号。
		// 这里使用楼盘字典 openId 匹配，避免“1号楼/1”等展示文本差异导致误新增。
		findErr := global.GVA_DB.Where(
			"owner = ? AND xiaoqu_id = ? AND building_id = ? AND unit_id = ? AND house_id = ?",
			userID, xq.ID, buildingID, unitID, houseID,
		).First(&entity).Error
		if findErr != nil {
			// 首次导入的新房源，尽量复用历史模板补齐更多字段（楼层、面积、字典ID、联系方式、图片等），
			// 减少批量导入后需要手工二次编辑的成本。
			if template, ok := findHistoricalTemplateResource(userID, xq.ID, buildingID, unitID, houseID); ok {
				entity = template
				// 新建时需要清掉主键和时间戳，避免把模板记录覆盖掉。
				entity.ID = 0
				entity.CreatedAt = time.Time{}
				entity.UpdatedAt = time.Time{}
				entity.DeletedAt = gorm.DeletedAt{}
				entity.Follow = 0
				entity.View = 0
				entity.Shared = 0
				entity.Click = 0
			} else {
				entity = house.Resource{
					Attachments: findHistoricalAttachments(userID, xq.ID, buildingID, unitID, houseID),
				}
			}
			// 新建记录时，定位字段以本次 Excel + 小区/楼盘字典数据为准。
			entity.Owner = userID
			entity.XiaoquId = xq.ID
			entity.Xiaoqu = xq.Name
			entity.Districts = xq.Districts
			entity.DistrictIds = xq.DistrictIds
			entity.City = xq.City
			entity.Region = xq.Area
			entity.DoorNo = doorNo
			entity.RoomCode = roomCode
			entity.BuildingId = buildingID
			entity.UnitId = unitID
			entity.HouseId = houseID
		}

		entity.RentType = defaultString(rentType, entity.RentType)
		entity.HouseType = defaultString(houseType, entity.HouseType)
		entity.Price = price
		entity.CommissionPrice = commission
		entity.Feature = feature
		entity.Remarks = remarks
		entity.Status = "待出租"
		entity.UpdatedLastAt = time.Now()

		// 继续复用统一的 CreateOrUpdate，确保团队房源、联系方式、上架额度等规则保持一致。
		if err = service.CreateOrUpdate(&entity); err != nil {
			record.FailedCount++
			failures = append(failures, fmt.Sprintf("第%d行导入失败:%s", idx+2, err.Error()))
			continue
		}
		record.SuccessCount++
	}
	record.ResultSummary = strings.Join(failures, "\n")
	err = global.GVA_DB.Model(&house.BatchUploadRecord{}).Where("id = ?", record.ID).Updates(record).Error
	return record, err
}

func normalizeHeaders(headers []string) []string {
	// Excel 表头允许存在多余空格，先统一 trim，减少模板格式带来的干扰。
	result := make([]string, 0, len(headers))
	for _, header := range headers {
		result = append(result, strings.TrimSpace(header))
	}
	return result
}

func mapRow(headers []string, row []string) map[string]string {
	// 通过表头名转 map，后续读取字段时不依赖固定列序，Excel 模板更容易兼容。
	result := make(map[string]string, len(headers))
	for idx, header := range headers {
		if idx >= len(row) {
			result[header] = ""
			continue
		}
		result[header] = strings.TrimSpace(row[idx])
	}
	return result
}

func isEmptyRow(row []string) bool {
	for _, item := range row {
		if strings.TrimSpace(item) != "" {
			return false
		}
	}
	return true
}

func atoi(v string) int {
	var n int
	fmt.Sscanf(strings.TrimSpace(v), "%d", &n)
	return n
}

func defaultString(v, fallback string) string {
	if strings.TrimSpace(v) == "" {
		return fallback
	}
	return strings.TrimSpace(v)
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

func composeDoorNoFromParts(building, unit, houseNo string) string {
	b := strings.TrimSpace(building)
	u := strings.TrimSpace(unit)
	h := strings.TrimSpace(houseNo)
	if b == "" || u == "" || h == "" {
		return ""
	}
	// 存储统一格式：楼栋 单元 门牌号（空格分隔），便于字典校验和复用匹配。
	return strings.Join([]string{b, u, h}, " ")
}

func findHistoricalAttachments(userID, xiaoquID uint, buildingID, unitID, houseID string) common.AttachmentMap {
	// 图片复用在“同账号 + 小区 + 楼栋 + 单元 + 门牌号”范围内查找，
	// 与导入主匹配键保持一致。
	var resource house.Resource
	if err := global.GVA_DB.Where(
		"owner = ? AND xiaoqu_id = ? AND building_id = ? AND unit_id = ? AND house_id = ?",
		userID, xiaoquID, buildingID, unitID, houseID,
	).Order("id desc").First(&resource).Error; err == nil {
		return resource.Attachments
	}
	return common.AttachmentMap{}
}

func findHistoricalTemplateResource(userID, xiaoquID uint, buildingID, unitID, houseID string) (resource house.Resource, ok bool) {
	// 模板优先级：
	// 1) 同账号+小区+楼栋+单元+门牌号
	// 2) 同账号+同小区最近一条，用于复用该账号在该小区的历史录入信息和图片。
	err := global.GVA_DB.Where(
		"owner = ? AND xiaoqu_id = ? AND building_id = ? AND unit_id = ? AND house_id = ?",
		userID, xiaoquID, buildingID, unitID, houseID,
	).Order("updated_last_at desc, id desc").First(&resource).Error
	if err == nil {
		return resource, true
	}

	err = global.GVA_DB.Where("owner = ? AND xiaoqu_id = ?", userID, xiaoquID).Order("updated_last_at desc, id desc").First(&resource).Error
	if err == nil {
		return resource, true
	}

	return house.Resource{}, false
}

func validateDoorNoByDict(xq system.XiaoQu, doorNo string) error {
	trimmed := strings.TrimSpace(doorNo)
	if trimmed == "" {
		return fmt.Errorf("为空")
	}
	parts := strings.Fields(trimmed)
	if len(parts) != 3 {
		return fmt.Errorf("格式不正确，需为“楼栋 单元 门牌号”并用空格分隔")
	}
	buildingPart := normalizeDoorPart(parts[0], "号楼")
	unitPart := normalizeDoorPart(parts[1], "单元")
	housePart := normalizeDoorPart(parts[2], "室")
	if buildingPart == "" || unitPart == "" || housePart == "" {
		return fmt.Errorf("格式不正确，楼栋/单元/门牌号不能为空")
	}

	// 小区未关联社区字典时，只做非空校验，避免历史小区数据被误拦截。
	if xq.CommunityId == 0 {
		return nil
	}
	var buildings []house.DictBuilding
	if err := global.GVA_DB.Where("community_id = ?", xq.CommunityId).Find(&buildings).Error; err != nil {
		return err
	}
	if len(buildings) == 0 {
		// 字典无楼栋数据时不阻断导入。
		return nil
	}
	var matchedBuilding *house.DictBuilding
	for _, b := range buildings {
		if normalizeDoorPart(b.EncryptBuildingName, "号楼") == buildingPart {
			tmp := b
			matchedBuilding = &tmp
			break
		}
	}
	if matchedBuilding == nil {
		return fmt.Errorf("未匹配到该小区楼栋字典")
	}

	var units []house.DictUnit
	if err := global.GVA_DB.Where("building_open_id = ?", matchedBuilding.BuildingOpenID).Find(&units).Error; err != nil {
		return err
	}
	if len(units) == 0 {
		return fmt.Errorf("该楼栋未维护单元字典")
	}
	var matchedUnit *house.DictUnit
	for _, u := range units {
		if normalizeDoorPart(u.EncryptUnitName, "单元") == unitPart {
			tmp := u
			matchedUnit = &tmp
			break
		}
	}
	if matchedUnit == nil {
		return fmt.Errorf("未匹配到该楼栋单元字典")
	}

	var houses []house.DictHouse
	if err := global.GVA_DB.Where("unit_open_id = ?", matchedUnit.UnitOpenID).Find(&houses).Error; err != nil {
		return err
	}
	if len(houses) == 0 {
		return fmt.Errorf("该单元未维护房号字典")
	}
	for _, h := range houses {
		if normalizeDoorPart(h.EncryptHouseName, "室") == housePart {
			return nil
		}
	}
	return fmt.Errorf("未匹配到该单元房号字典")
}

func normalizeDoorPart(raw, suffix string) string {
	value := strings.TrimSpace(raw)
	value = strings.TrimSuffix(value, suffix)
	value = strings.TrimSpace(value)
	return value
}

func resolveDoorDictIDs(xq system.XiaoQu, doorNo string) (buildingOpenID, unitOpenID, houseOpenID string, err error) {
	parts := strings.Fields(strings.TrimSpace(doorNo))
	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("格式不正确，需为“楼栋 单元 门牌号”并用空格分隔")
	}
	buildingPart := normalizeDoorPart(parts[0], "号楼")
	unitPart := normalizeDoorPart(parts[1], "单元")
	housePart := normalizeDoorPart(parts[2], "室")
	if buildingPart == "" || unitPart == "" || housePart == "" {
		return "", "", "", fmt.Errorf("格式不正确，楼栋/单元/门牌号不能为空")
	}
	if xq.CommunityId == 0 {
		return "", "", "", fmt.Errorf("小区未维护 community_id")
	}

	var buildings []house.DictBuilding
	if err = global.GVA_DB.Where("community_id = ?", xq.CommunityId).Find(&buildings).Error; err != nil {
		return "", "", "", err
	}
	var matchedBuilding *house.DictBuilding
	for _, b := range buildings {
		if normalizeDoorPart(b.EncryptBuildingName, "号楼") == buildingPart {
			tmp := b
			matchedBuilding = &tmp
			break
		}
	}
	if matchedBuilding == nil {
		return "", "", "", fmt.Errorf("未匹配到该小区楼栋字典")
	}

	var units []house.DictUnit
	if err = global.GVA_DB.Where("building_open_id = ?", matchedBuilding.BuildingOpenID).Find(&units).Error; err != nil {
		return "", "", "", err
	}
	var matchedUnit *house.DictUnit
	for _, u := range units {
		if normalizeDoorPart(u.EncryptUnitName, "单元") == unitPart {
			tmp := u
			matchedUnit = &tmp
			break
		}
	}
	if matchedUnit == nil {
		return "", "", "", fmt.Errorf("未匹配到该楼栋单元字典")
	}

	var houses []house.DictHouse
	if err = global.GVA_DB.Where("unit_open_id = ?", matchedUnit.UnitOpenID).Find(&houses).Error; err != nil {
		return "", "", "", err
	}
	for _, h := range houses {
		if normalizeDoorPart(h.EncryptHouseName, "室") == housePart {
			return matchedBuilding.BuildingOpenID, matchedUnit.UnitOpenID, h.HouseOpenId, nil
		}
	}
	return "", "", "", fmt.Errorf("未匹配到该单元房号字典")
}
