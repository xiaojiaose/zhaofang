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

func (service *ResourceService) BatchUpload(userID uint, header *multipart.FileHeader) (record *house.BatchUploadRecord, err error) {
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
	}
	if err = global.GVA_DB.Create(record).Error; err != nil {
		return nil, err
	}

	_ = global.GVA_DB.Model(&house.Resource{}).Where("owner = ? AND status = ?", userID, "待出租").Update("status", "已下架").Error

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

		var entity house.Resource
		findErr := global.GVA_DB.Where("owner = ? AND xiaoqu_id = ? AND door_no = ? AND room_code = ?", userID, xq.ID, doorNo, roomCode).First(&entity).Error
		if findErr != nil {
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
	result := make([]string, 0, len(headers))
	for _, header := range headers {
		result = append(result, strings.TrimSpace(header))
	}
	return result
}

func mapRow(headers []string, row []string) map[string]string {
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
	var resource house.Resource
	if err := global.GVA_DB.Where("owner = ? AND xiaoqu_id = ? AND door_no = ? AND room_code = ?", userID, xiaoquID, doorNo, roomCode).Order("id desc").First(&resource).Error; err == nil {
		return resource.Attachments
	}
	return common.AttachmentMap{}
}
