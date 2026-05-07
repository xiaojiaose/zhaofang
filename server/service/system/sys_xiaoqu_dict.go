package system

import (
	"errors"
	"fmt"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	houseModel "github.com/flipped-aurora/gin-vue-admin/server/model/house"
	sysModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	sysReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	sysResp "github.com/flipped-aurora/gin-vue-admin/server/model/system/response"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (s *XiaoQuService) GetDictTree(xiaoquID uint) (result sysResp.XiaoQuDictTreeResponse, err error) {
	var xq sysModel.XiaoQu
	if err = global.GVA_DB.First(&xq, "id = ?", xiaoquID).Error; err != nil {
		return result, err
	}
	result.XiaoquID = xiaoquID
	result.CommunityID = xq.CommunityId
	if xq.CommunityId == 0 {
		result.Buildings = []sysResp.XiaoQuDictBuildingNode{}
		return result, nil
	}

	var buildings []houseModel.DictBuilding
	if err = global.GVA_DB.Where("community_id = ?", xq.CommunityId).Order("id asc").Find(&buildings).Error; err != nil {
		return result, err
	}

	result.Buildings = make([]sysResp.XiaoQuDictBuildingNode, 0, len(buildings))
	for _, building := range buildings {
		buildingNode := sysResp.XiaoQuDictBuildingNode{
			ID:             building.ID,
			BuildingOpenID: building.BuildingOpenID,
			Name:           building.EncryptBuildingName,
			Units:          []sysResp.XiaoQuDictUnitNode{},
		}

		var units []houseModel.DictUnit
		if err = global.GVA_DB.Where("building_open_id = ?", building.BuildingOpenID).Order("id asc").Find(&units).Error; err != nil {
			return result, err
		}

		buildingNode.Units = make([]sysResp.XiaoQuDictUnitNode, 0, len(units))
		for _, unit := range units {
			unitNode := sysResp.XiaoQuDictUnitNode{
				ID:         unit.ID,
				UnitOpenID: unit.UnitOpenID,
				Name:       unit.EncryptUnitName,
				Houses:     []sysResp.XiaoQuDictHouseNode{},
			}

			var houses []houseModel.DictHouse
			if err = global.GVA_DB.Where("unit_open_id = ?", unit.UnitOpenID).Order("id asc").Find(&houses).Error; err != nil {
				return result, err
			}

			unitNode.Houses = make([]sysResp.XiaoQuDictHouseNode, 0, len(houses))
			for _, h := range houses {
				unitNode.Houses = append(unitNode.Houses, sysResp.XiaoQuDictHouseNode{
					ID:          h.ID,
					HouseOpenID: h.HouseOpenId,
					Name:        h.EncryptHouseName,
				})
			}
			buildingNode.Units = append(buildingNode.Units, unitNode)
		}
		result.Buildings = append(result.Buildings, buildingNode)
	}
	return result, nil
}

func (s *XiaoQuService) UpsertDict(req sysReq.XiaoQuDictUpsertRequest) (result sysResp.XiaoQuDictTreeResponse, err error) {
	if req.XiaoquID == 0 {
		return result, errors.New("小区ID不能为空")
	}
	tx := global.GVA_DB.Begin()
	if tx.Error != nil {
		return result, tx.Error
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	communityID, err := s.ensureCommunityID(tx, req.XiaoquID)
	if err != nil {
		return result, err
	}
	for _, op := range req.BuildingOps {
		if err = s.upsertBuilding(tx, communityID, op); err != nil {
			return result, err
		}
	}
	for _, op := range req.UnitOps {
		if err = s.upsertUnit(tx, communityID, op); err != nil {
			return result, err
		}
	}
	for _, op := range req.HouseOps {
		if err = s.upsertHouse(tx, communityID, op); err != nil {
			return result, err
		}
	}
	if err = tx.Commit().Error; err != nil {
		return result, err
	}
	return s.GetDictTree(req.XiaoquID)
}

func (s *XiaoQuService) DeleteDict(req sysReq.XiaoQuDictDeleteRequest) (result sysResp.XiaoQuDictTreeResponse, err error) {
	if req.XiaoquID == 0 {
		return result, errors.New("小区ID不能为空")
	}
	tx := global.GVA_DB.Begin()
	if tx.Error != nil {
		return result, tx.Error
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	communityID, err := s.ensureCommunityID(tx, req.XiaoquID)
	if err != nil {
		return result, err
	}
	for _, id := range req.HouseIDs {
		if err = s.deleteHouse(tx, communityID, id); err != nil {
			return result, err
		}
	}
	for _, id := range req.UnitIDs {
		if err = s.deleteUnit(tx, communityID, id); err != nil {
			return result, err
		}
	}
	for _, id := range req.BuildingIDs {
		if err = s.deleteBuilding(tx, communityID, id); err != nil {
			return result, err
		}
	}
	if err = tx.Commit().Error; err != nil {
		return result, err
	}
	return s.GetDictTree(req.XiaoquID)
}

func (s *XiaoQuService) ensureCommunityID(tx *gorm.DB, xiaoquID uint) (int, error) {
	var xq sysModel.XiaoQu
	if err := tx.First(&xq, "id = ?", xiaoquID).Error; err != nil {
		return 0, err
	}
	if xq.CommunityId == 0 {
		return 0, errors.New("当前小区未配置community_id，无法维护楼盘字典")
	}
	return xq.CommunityId, nil
}

func (s *XiaoQuService) upsertBuilding(tx *gorm.DB, communityID int, op sysReq.XiaoQuDictBuildingOp) error {
	name := strings.TrimSpace(op.Name)
	if name == "" {
		return errors.New("楼栋名称不能为空")
	}
	if op.ID == 0 {
		var existing houseModel.DictBuilding
		if err := tx.Where("community_id = ? AND encrypt_building_name = ?", communityID, name).First(&existing).Error; err == nil {
			return errors.New("楼栋名称已存在")
		}
		record := houseModel.DictBuilding{
			CommunityID:         communityID,
			EncryptBuildingName: name,
			BuildingOpenID:      strings.TrimSpace(op.BuildingOpenID),
		}
		if record.BuildingOpenID == "" {
			record.BuildingOpenID = uuid.NewString()
		}
		return tx.Create(&record).Error
	}
	var record houseModel.DictBuilding
	if err := tx.Where("id = ? AND community_id = ?", op.ID, communityID).First(&record).Error; err != nil {
		return fmt.Errorf("楼栋ID不存在: %d", op.ID)
	}
	var existing houseModel.DictBuilding
	if err := tx.Where("community_id = ? AND encrypt_building_name = ? AND id != ?", communityID, name, op.ID).First(&existing).Error; err == nil {
		return errors.New("楼栋名称已存在")
	}
	record.EncryptBuildingName = name
	return tx.Save(&record).Error
}

func (s *XiaoQuService) upsertUnit(tx *gorm.DB, communityID int, op sysReq.XiaoQuDictUnitOp) error {
	name := strings.TrimSpace(op.Name)
	if name == "" {
		return errors.New("单元名称不能为空")
	}
	if op.ID == 0 {
		buildingOpenID := strings.TrimSpace(op.BuildingOpenID)
		if buildingOpenID == "" {
			return errors.New("新增单元必须传buildingOpenId")
		}
		var building houseModel.DictBuilding
		if err := tx.Where("building_open_id = ? AND community_id = ?", buildingOpenID, communityID).First(&building).Error; err != nil {
			return fmt.Errorf("楼栋OpenID不存在: %s", buildingOpenID)
		}
		var existing houseModel.DictUnit
		if err := tx.Where("building_open_id = ? AND encrypt_unit_name = ?", buildingOpenID, name).First(&existing).Error; err == nil {
			return errors.New("单元名称已存在")
		}
		record := houseModel.DictUnit{
			BuildingOpenID:  building.BuildingOpenID,
			EncryptUnitName: name,
			UnitOpenID:      strings.TrimSpace(op.UnitOpenID),
		}
		if record.UnitOpenID == "" {
			record.UnitOpenID = uuid.NewString()
		}
		return tx.Create(&record).Error
	}
	var record houseModel.DictUnit
	if err := tx.Where("id = ?", op.ID).First(&record).Error; err != nil {
		return fmt.Errorf("单元ID不存在: %d", op.ID)
	}
	var building houseModel.DictBuilding
	if err := tx.Where("building_open_id = ? AND community_id = ?", record.BuildingOpenID, communityID).First(&building).Error; err != nil {
		return fmt.Errorf("单元ID不存在: %d", op.ID)
	}
	var existing houseModel.DictUnit
	if err := tx.Where("building_open_id = ? AND encrypt_unit_name = ? AND id != ?", record.BuildingOpenID, name, op.ID).First(&existing).Error; err == nil {
		return errors.New("单元名称已存在")
	}
	record.EncryptUnitName = name
	return tx.Save(&record).Error
}

func (s *XiaoQuService) upsertHouse(tx *gorm.DB, communityID int, op sysReq.XiaoQuDictHouseOp) error {
	name := strings.TrimSpace(op.Name)
	if name == "" {
		return errors.New("房号名称不能为空")
	}
	if op.ID == 0 {
		unitOpenID := strings.TrimSpace(op.UnitOpenID)
		if unitOpenID == "" {
			return errors.New("新增房号必须传unitOpenId")
		}
		var unit houseModel.DictUnit
		if err := tx.Where("unit_open_id = ?", unitOpenID).First(&unit).Error; err != nil {
			return fmt.Errorf("单元OpenID不存在: %s", unitOpenID)
		}
		var building houseModel.DictBuilding
		if err := tx.Where("building_open_id = ? AND community_id = ?", unit.BuildingOpenID, communityID).First(&building).Error; err != nil {
			return fmt.Errorf("单元OpenID不存在: %s", unitOpenID)
		}
		var existing houseModel.DictHouse
		if err := tx.Where("unit_open_id = ? AND encrypt_house_name = ?", unitOpenID, name).First(&existing).Error; err == nil {
			return errors.New("房号名称已存在")
		}
		record := houseModel.DictHouse{
			UnitOpenId:       unit.UnitOpenID,
			EncryptHouseName: name,
			HouseOpenId:      strings.TrimSpace(op.HouseOpenID),
		}
		if record.HouseOpenId == "" {
			record.HouseOpenId = uuid.NewString()
		}
		return tx.Create(&record).Error
	}
	var record houseModel.DictHouse
	if err := tx.Where("id = ?", op.ID).First(&record).Error; err != nil {
		return fmt.Errorf("房号ID不存在: %d", op.ID)
	}
	var unit houseModel.DictUnit
	if err := tx.Where("unit_open_id = ?", record.UnitOpenId).First(&unit).Error; err != nil {
		return fmt.Errorf("房号ID不存在: %d", op.ID)
	}
	var building houseModel.DictBuilding
	if err := tx.Where("building_open_id = ? AND community_id = ?", unit.BuildingOpenID, communityID).First(&building).Error; err != nil {
		return fmt.Errorf("房号ID不存在: %d", op.ID)
	}
	var existing houseModel.DictHouse
	if err := tx.Where("unit_open_id = ? AND encrypt_house_name = ? AND id != ?", record.UnitOpenId, name, op.ID).First(&existing).Error; err == nil {
		return errors.New("房号名称已存在")
	}
	record.EncryptHouseName = name
	return tx.Save(&record).Error
}

func (s *XiaoQuService) deleteBuilding(tx *gorm.DB, communityID int, buildingID uint) error {
	if buildingID == 0 {
		return errors.New("楼栋删除必须传ID")
	}
	var record houseModel.DictBuilding
	if err := tx.Where("id = ? AND community_id = ?", buildingID, communityID).First(&record).Error; err != nil {
		return fmt.Errorf("楼栋ID不存在: %d", buildingID)
	}
	var unitOpenIDs []string
	if err := tx.Model(&houseModel.DictUnit{}).Where("building_open_id = ?", record.BuildingOpenID).Pluck("unit_open_id", &unitOpenIDs).Error; err != nil {
		return err
	}
	if err := tx.Where("building_open_id = ?", record.BuildingOpenID).Delete(&houseModel.DictUnit{}).Error; err != nil {
		return err
	}
	if len(unitOpenIDs) > 0 {
		if err := tx.Where("unit_open_id IN ?", unitOpenIDs).Delete(&houseModel.DictHouse{}).Error; err != nil {
			return err
		}
	}
	return tx.Delete(&record).Error
}

func (s *XiaoQuService) deleteUnit(tx *gorm.DB, communityID int, unitID uint) error {
	if unitID == 0 {
		return errors.New("单元删除必须传ID")
	}
	var record houseModel.DictUnit
	if err := tx.Where("id = ?", unitID).First(&record).Error; err != nil {
		return fmt.Errorf("单元ID不存在: %d", unitID)
	}
	var building houseModel.DictBuilding
	if err := tx.Where("building_open_id = ? AND community_id = ?", record.BuildingOpenID, communityID).First(&building).Error; err != nil {
		return fmt.Errorf("单元ID不存在: %d", unitID)
	}
	if err := tx.Where("unit_open_id = ?", record.UnitOpenID).Delete(&houseModel.DictHouse{}).Error; err != nil {
		return err
	}
	return tx.Delete(&record).Error
}

func (s *XiaoQuService) deleteHouse(tx *gorm.DB, communityID int, houseID uint) error {
	if houseID == 0 {
		return errors.New("房号删除必须传ID")
	}
	var record houseModel.DictHouse
	if err := tx.Where("id = ?", houseID).First(&record).Error; err != nil {
		return fmt.Errorf("房号ID不存在: %d", houseID)
	}
	var unit houseModel.DictUnit
	if err := tx.Where("unit_open_id = ?", record.UnitOpenId).First(&unit).Error; err != nil {
		return fmt.Errorf("房号ID不存在: %d", houseID)
	}
	var building houseModel.DictBuilding
	if err := tx.Where("building_open_id = ? AND community_id = ?", unit.BuildingOpenID, communityID).First(&building).Error; err != nil {
		return fmt.Errorf("房号ID不存在: %d", houseID)
	}
	return tx.Delete(&record).Error
}
