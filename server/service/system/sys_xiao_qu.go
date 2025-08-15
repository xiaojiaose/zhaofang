package system

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"gorm.io/gorm"
)

type XiaoQuService struct{}

func (s *XiaoQuService) Edit(u system.XiaoQu) (xiaoQu system.XiaoQu, err error) {
	var xq system.XiaoQu
	tx := global.GVA_DB.Where("name = ?", u.Name).First(&xq)
	if tx.Error != nil && !errors.Is(tx.Error, gorm.ErrRecordNotFound) { // 判断名是否被占用
		if xq.ID != u.ID {
			return xq, errors.New("小区名已占用")
		}
		return xq, err
	}

	tx = global.GVA_DB.Where("id = ?", u.ID).First(&xq)
	err = tx.Updates(&u).Error

	return u, err

}

func (s *XiaoQuService) EditNum(id uint, houseNum int) (err error) {
	tx := global.GVA_DB.Model(&system.XiaoQu{}).Where("ID = ?", id).Update("house_num", houseNum)

	return tx.Error

}

func (s *XiaoQuService) Create(u system.XiaoQu) (xiaoQu system.XiaoQu, err error) {
	var xq system.XiaoQu
	tx := global.GVA_DB.Where("name = ?", u.Name).First(&xq)
	if !errors.Is(tx.Error, gorm.ErrRecordNotFound) { // 判断名是否被占用
		if xq.ID != u.ID {
			return xq, errors.New("小区名已占用")
		}
	}

	err = global.GVA_DB.Create(&u).Error
	return u, err
}

func (s *XiaoQuService) GetInfo(id uint) (user system.XiaoQu, err error) {
	var qu system.XiaoQu
	err = global.GVA_DB.First(&qu, "id = ?", id).Error
	if err != nil {
		return qu, err
	}
	return qu, err
}

func (s *XiaoQuService) GetDistance(lat float64, lng float64) (xqList []system.XiaoQu, err error) {
	tx := global.GVA_DB.Raw("SELECT id,name, longitude, latitude, (6371 * ACOS(COS(RADIANS(@lat)) * COS(RADIANS(latitude)) * COS(RADIANS(longitude) - RADIANS(@lng)) + SIN(RADIANS(@lat)) * SIN(RADIANS(latitude)))) AS distance FROM xiao_qu HAVING distance < 1 ORDER BY distance;",
		sql.Named("lat", lat), sql.Named("lng", lng)).
		Find(&xqList)
	return xqList, tx.Error
}

func (s *XiaoQuService) GetList(info request.SearchXiaoqu) (list []system.XiaoQu, total int64, err error) {
	limit := info.PageSize
	if limit == 0 {
		limit = 5000
	}
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&system.XiaoQu{})
	if info.Keyword != "" {
		db = db.Where(" name LIKE ?", "%"+info.Keyword+"%")
	}
	if info.CityId != "" {
		db = db.Where(" city = ?", info.CityId)
	}
	db2 := global.GVA_DB.Model(&system.XiaoQu{})
	if len(info.Districts) > 0 {
		db2 = db2.Where("district_ids LIKE ?", "%|"+fmt.Sprintf("%v", info.Districts[0])+"|%")
		for _, v := range info.Districts[1:] {
			db2 = db2.Or("district_ids LIKE ?", "%|"+fmt.Sprintf("%v", v)+"|%")
		}
		db.Where(db2)
	}
	var xqList []system.XiaoQu
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Find(&xqList).Error
	return xqList, total, err
}

func (s *XiaoQuService) GetAreaList(info request.SearchArea) (list []system.Area, total int64, err error) {
	db := global.GVA_DB.Model(&system.Area{})
	if info.CityId != "" {
		db = db.Where(" city = ?", info.CityId)
	}
	var xqList []system.Area
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Find(&xqList).Error
	return xqList, total, err
}

func (s *XiaoQuService) GetDistrictList(info request.SearchDistrict) (list []system.District, total int64, err error) {
	db := global.GVA_DB.Model(&system.District{})
	if info.AreaId != 0 {
		db = db.Where(" area_id = ?", info.AreaId)
	}
	var xqList []system.District
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Find(&xqList).Error
	return xqList, total, err
}
