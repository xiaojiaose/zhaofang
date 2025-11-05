package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/house"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DictUnitApi struct{}

// CreateDictUnit 创建单元字典
// @Tags      DictUnit
// @Summary   创建单元字典
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      house.DictUnit                 true  "单元字典模型"
// @Success   200   {object}  response.Response{msg=string}  "创建成功"
// @Router    /dictUnit/createDictUnit [post]
func (d *DictUnitApi) CreateDictUnit(c *gin.Context) {
	var dictUnit house.DictUnit
	err := c.ShouldBindJSON(&dictUnit)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = global.GVA_DB.Create(&dictUnit).Error
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteDictUnit 删除单元字典
// @Tags      DictUnit
// @Summary   删除单元字典
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      house.DictUnit                 true  "单元字典模型"
// @Success   200   {object}  response.Response{msg=string}  "删除成功"
// @Router    /dictUnit/deleteDictUnit [delete]
func (d *DictUnitApi) DeleteDictUnit(c *gin.Context) {
	var dictUnit house.DictUnit
	err := c.ShouldBindJSON(&dictUnit)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = global.GVA_DB.Delete(&dictUnit).Error
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// UpdateDictUnit 更新单元字典
// @Tags      DictUnit
// @Summary   更新单元字典
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      house.DictUnit                 true  "单元字典模型"
// @Success   200   {object}  response.Response{msg=string}  "更新成功"
// @Router    /dictUnit/updateDictUnit [put]
func (d *DictUnitApi) UpdateDictUnit(c *gin.Context) {
	var dictUnit house.DictUnit
	err := c.ShouldBindJSON(&dictUnit)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = global.GVA_DB.Save(&dictUnit).Error
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindDictUnit 根据ID获取单元字典
// @Tags      DictUnit
// @Summary   根据ID获取单元字典
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      house.DictUnit                    true  "单元字典模型"
// @Success   200   {object}  response.Response{data=house.DictUnit}  "查询成功"
// @Router    /dictUnit/findDictUnit [get]
func (d *DictUnitApi) FindDictUnit(c *gin.Context) {
	var dictUnit house.DictUnit
	err := c.ShouldBindJSON(&dictUnit)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = global.GVA_DB.Where("id = ?", dictUnit.ID).First(&dictUnit).Error
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithData(dictUnit, c)
}

// GetDictUnitList 分页获取单元字典列表
// @Tags      DictUnit
// @Summary   分页获取单元字典列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.PageInfo                        true  "分页参数"
// @Success   200   {object}  response.Response{data=response.PageResult}  "获取成功"
// @Router    /dictUnit/getDictUnitList [post]
func (d *DictUnitApi) GetDictUnitList(c *gin.Context) {
	var pageInfo request.PageInfo
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	
	limit := pageInfo.PageSize
	offset := pageInfo.PageSize * (pageInfo.Page - 1)
	
	var dictUnits []house.DictUnit
	var total int64
	
	db := global.GVA_DB.Model(&house.DictUnit{})
	err = db.Count(&total).Error
	if err != nil {
		global.GVA_LOG.Error("获取总数失败!", zap.Error(err))
		response.FailWithMessage("获取总数失败", c)
		return
	}
	
	err = db.Limit(limit).Offset(offset).Find(&dictUnits).Error
	if err != nil {
		global.GVA_LOG.Error("获取数据失败!", zap.Error(err))
		response.FailWithMessage("获取数据失败", c)
		return
	}
	
	response.OkWithDetailed(response.PageResult{
		List:     dictUnits,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}
