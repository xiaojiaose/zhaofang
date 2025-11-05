package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/house"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DictHouseApi struct{}

// CreateDictHouse 创建房屋字典
// @Tags      DictHouse
// @Summary   创建房屋字典
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      house.DictHouse                true  "房屋字典模型"
// @Success   200   {object}  response.Response{msg=string}  "创建成功"
// @Router    /dictHouse/createDictHouse [post]
func (d *DictHouseApi) CreateDictHouse(c *gin.Context) {
	var dictHouse house.DictHouse
	err := c.ShouldBindJSON(&dictHouse)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = global.GVA_DB.Create(&dictHouse).Error
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteDictHouse 删除房屋字典
// @Tags      DictHouse
// @Summary   删除房屋字典
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      house.DictHouse                true  "房屋字典模型"
// @Success   200   {object}  response.Response{msg=string}  "删除成功"
// @Router    /dictHouse/deleteDictHouse [delete]
func (d *DictHouseApi) DeleteDictHouse(c *gin.Context) {
	var dictHouse house.DictHouse
	err := c.ShouldBindJSON(&dictHouse)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = global.GVA_DB.Delete(&dictHouse).Error
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// UpdateDictHouse 更新房屋字典
// @Tags      DictHouse
// @Summary   更新房屋字典
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      house.DictHouse                true  "房屋字典模型"
// @Success   200   {object}  response.Response{msg=string}  "更新成功"
// @Router    /dictHouse/updateDictHouse [put]
func (d *DictHouseApi) UpdateDictHouse(c *gin.Context) {
	var dictHouse house.DictHouse
	err := c.ShouldBindJSON(&dictHouse)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = global.GVA_DB.Save(&dictHouse).Error
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindDictHouse 根据ID获取房屋字典
// @Tags      DictHouse
// @Summary   根据ID获取房屋字典
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      house.DictHouse                   true  "房屋字典模型"
// @Success   200   {object}  response.Response{data=house.DictHouse}  "查询成功"
// @Router    /dictHouse/findDictHouse [get]
func (d *DictHouseApi) FindDictHouse(c *gin.Context) {
	var dictHouse house.DictHouse
	err := c.ShouldBindJSON(&dictHouse)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = global.GVA_DB.Where("id = ?", dictHouse.ID).First(&dictHouse).Error
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithData(dictHouse, c)
}

// GetDictHouseList 分页获取房屋字典列表
// @Tags      DictHouse
// @Summary   分页获取房屋字典列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.PageInfo                        true  "分页参数"
// @Success   200   {object}  response.Response{data=response.PageResult}  "获取成功"
// @Router    /dictHouse/getDictHouseList [post]
func (d *DictHouseApi) GetDictHouseList(c *gin.Context) {
	var pageInfo request.PageInfo
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	
	limit := pageInfo.PageSize
	offset := pageInfo.PageSize * (pageInfo.Page - 1)
	
	var dictHouses []house.DictHouse
	var total int64
	
	db := global.GVA_DB.Model(&house.DictHouse{})
	err = db.Count(&total).Error
	if err != nil {
		global.GVA_LOG.Error("获取总数失败!", zap.Error(err))
		response.FailWithMessage("获取总数失败", c)
		return
	}
	
	err = db.Limit(limit).Offset(offset).Find(&dictHouses).Error
	if err != nil {
		global.GVA_LOG.Error("获取数据失败!", zap.Error(err))
		response.FailWithMessage("获取数据失败", c)
		return
	}
	
	response.OkWithDetailed(response.PageResult{
		List:     dictHouses,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}
