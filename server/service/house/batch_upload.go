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
		xiaoquName := values["楼盘名称"]
		if xiaoquName == "" {
			xiaoquName = values["小区名称"]
		}
		doorNo := values["楼栋房间号"]
		if doorNo == "" {
			doorNo = values["户室号"]
		}
		roomCode := values["房间号"]
		price := atoi(values["价格"])
		commission := atoi(values["返佣金额"])
		remarks := values["备注"]
		feature := values["标签"]
		rentType := values["出租类型"]
		houseType := values["房屋类型"]
		if houseType == "" {
			houseType = values["房间类型"]
		}

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
		findErr := global.GVA_DB.Where("owner = ? AND xiaoqu_id = ? AND door_no = ? AND room_code = ?", userID, xq.ID, doorNo, roomCode).First(&entity).Error
		if findErr != nil {
			// 首次导入的新房源优先尝试复用同账号、同地址/房号的历史图片，
			// 避免每次批量同步都要求重新传图。
			entity = house.Resource{
				Owner:       userID,
				XiaoquId:    xq.ID,
				Xiaoqu:      xq.Name,
				Districts:   xq.Districts,
				DistrictIds: xq.DistrictIds,
				City:        xq.City,
				Region:      xq.Area,
				DoorNo:      doorNo,
				RoomCode:    roomCode,
				Attachments: findHistoricalAttachments(userID, xq.ID, doorNo, roomCode),
			}
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

func findHistoricalAttachments(userID, xiaoquID uint, doorNo, roomCode string) common.AttachmentMap {
	// 图片复用只在“同账号 + 同小区 + 同户室 + 同房间号”范围内查找，
	// 避免把其他账号的历史图片错误带入当前导入结果。
	var resource house.Resource
	if err := global.GVA_DB.Where("owner = ? AND xiaoqu_id = ? AND door_no = ? AND room_code = ?", userID, xiaoquID, doorNo, roomCode).Order("id desc").First(&resource).Error; err == nil {
		return resource.Attachments
	}
	return common.AttachmentMap{}
}

func validateDoorNoByDict(xq system.XiaoQu, doorNo string) error {
	trimmed := strings.TrimSpace(doorNo)
	if trimmed == "" {
		return fmt.Errorf("为空")
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
	for _, b := range buildings {
		name := strings.TrimSpace(b.EncryptBuildingName)
		if name == "" {
			continue
		}
		if strings.Contains(trimmed, name+"号楼") || strings.Contains(trimmed, name) {
			return nil
		}
	}
	return fmt.Errorf("未匹配到该小区楼栋字典")
}
