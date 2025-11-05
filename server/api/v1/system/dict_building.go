package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/house"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DictBuildingApi struct{}

// CreateDictBuilding 创建楼栋字典
// @Tags      DictBuilding
// @Summary   创建楼栋字典
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      house.DictBuilding             true  "楼栋字典模型"
// @Success   200   {object}  response.Response{msg=string}  "创建成功"
// @Router    /dictBuilding/createDictBuilding [post]
func (d *DictBuildingApi) CreateDictBuilding(c *gin.Context) {
	var dictBuilding house.DictBuilding
	err := c.ShouldBindJSON(&dictBuilding)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = global.GVA_DB.Create(&dictBuilding).Error
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteDictBuilding 删除楼栋字典
// @Tags      DictBuilding
// @Summary   删除楼栋字典
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      house.DictBuilding             true  "楼栋字典模型"
// @Success   200   {object}  response.Response{msg=string}  "删除成功"
// @Router    /dictBuilding/deleteDictBuilding [delete]
func (d *DictBuildingApi) DeleteDictBuilding(c *gin.Context) {
	var dictBuilding house.DictBuilding
	err := c.ShouldBindJSON(&dictBuilding)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = global.GVA_DB.Delete(&dictBuilding).Error
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// UpdateDictBuilding 更新楼栋字典
// @Tags      DictBuilding
// @Summary   更新楼栋字典
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      house.DictBuilding             true  "楼栋字典模型"
// @Success   200   {object}  response.Response{msg=string}  "更新成功"
// @Router    /dictBuilding/updateDictBuilding [put]
func (d *DictBuildingApi) UpdateDictBuilding(c *gin.Context) {
	var dictBuilding house.DictBuilding
	err := c.ShouldBindJSON(&dictBuilding)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = global.GVA_DB.Save(&dictBuilding).Error
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindDictBuilding 根据ID获取楼栋字典
// @Tags      DictBuilding
// @Summary   根据ID获取楼栋字典
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      house.DictBuilding                      true  "楼栋字典模型"
// @Success   200   {object}  response.Response{data=house.DictBuilding}  "查询成功"
// @Router    /dictBuilding/findDictBuilding [get]
func (d *DictBuildingApi) FindDictBuilding(c *gin.Context) {
	var dictBuilding house.DictBuilding
	err := c.ShouldBindJSON(&dictBuilding)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = global.GVA_DB.Where("id = ?", dictBuilding.ID).First(&dictBuilding).Error
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithData(dictBuilding, c)
}

// GetDictBuildingList 分页获取楼栋字典列表
// @Tags      DictBuilding
// @Summary   分页获取楼栋字典列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.PageInfo                        true  "分页参数"
// @Success   200   {object}  response.Response{data=response.PageResult}  "获取成功"
// @Router    /dictBuilding/getDictBuildingList [post]
func (d *DictBuildingApi) GetDictBuildingList(c *gin.Context) {
	var pageInfo request.PageInfo
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	
	limit := pageInfo.PageSize
	offset := pageInfo.PageSize * (pageInfo.Page - 1)
	
	var dictBuildings []house.DictBuilding
	var total int64
	
	db := global.GVA_DB.Model(&house.DictBuilding{})
	err = db.Count(&total).Error
	if err != nil {
		global.GVA_LOG.Error("获取总数失败!", zap.Error(err))
		response.FailWithMessage("获取总数失败", c)
		return
	}
	
	err = db.Limit(limit).Offset(offset).Find(&dictBuildings).Error
	if err != nil {
		global.GVA_LOG.Error("获取数据失败!", zap.Error(err))
		response.FailWithMessage("获取数据失败", c)
		return
	}
	
	response.OkWithDetailed(response.PageResult{
		List:     dictBuildings,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}
