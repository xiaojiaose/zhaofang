package resource

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	response2 "github.com/flipped-aurora/gin-vue-admin/server/model/house/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

type XiaoQuApi struct{}

// Show
// @Tags     Center
// @Summary  获取小区信息
// @Produce  application/json
// @Param    data  query    request.GetById  true  "id"
// @Success  200   {object}  response.Response{data=system.XiaoQu,msg=string}  "结果"
// @Router   /center/xiaoqu/show [get]
func (receiver *XiaoQuApi) Show(c *gin.Context) {
	var api request.GetById
	err := c.ShouldBindQuery(&api)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	qu, err2 := XiaoQuService.GetInfo(uint(api.ID))
	if err2 != nil {
		response.FailWithMessage(err2.Error(), c)
		return
	}
	response.OkWithDetailed(qu, "获取成功", c)
	return
}

// Distance
// @Tags     Center
// @Summary  根据坐标获取一公里小区列表
// @accept    application/json
// @Produce   application/json
// @Param   lat    query    int     true  "lat"
// @Param   lng  query    int  false "lng"
// @Success   200   {object}  response.Response{data=[]system.XiaoQu,msg=string}
// @Router    /center/distance [get]
func (receiver *XiaoQuApi) Distance(c *gin.Context) {
	lat, _ := strconv.ParseFloat(c.Query("lat"), 64)
	lng, _ := strconv.ParseFloat(c.Query("lng"), 64)
	qus, err := XiaoQuService.GetDistance(lat, lng)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(qus, "获取成功", c)
}

// List
// @Tags     Center
// @Summary  获取小区列表
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.SearchXiaoqu    true  "页码, 每页大小"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "分页列表,返回包括列表,总数,页码,每页数量"
// @Router    /center/xiaoqu/list [post]
func (receiver *XiaoQuApi) List(c *gin.Context) {
	var pageInfo request.SearchXiaoqu
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	//err = utils.Verify(pageInfo, utils.PageInfoVerify)
	//if err != nil {
	//	response.FailWithMessage(err.Error(), c)
	//	return
	//}
	list, total, err := XiaoQuService.GetList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)

}

// Show
// @Tags     Center
// @Summary  指定小区id获取楼栋列表
// @Produce  application/json
// @Param    data  query    request.GetById  true  "小区 id"
// @Success  200   {object}  response.Response{data=house.DictBuilding,msg=string}  "结果"
// @Router   /base/building [get]
func (receiver *XiaoQuApi) GetBuilding(c *gin.Context) {
	var api request.GetById
	err := c.ShouldBindQuery(&api)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	xq, _ := XiaoQuService.GetInfo(uint(api.ID))

	qu, err2 := DictService.GetBuilding(strconv.Itoa(xq.CommunityId))
	if err2 != nil {
		response.FailWithMessage(err2.Error(), c)
		return
	}
	var ll []response2.DictBuildingResponse
	for _, building := range qu {
		ll = append(ll, response2.DictBuildingResponse{
			Id:   building.BuildingOpenID,
			Name: building.EncryptBuildingName,
		})
	}
	response.OkWithDetailed(ll, "获取成功", c)
	return
}

// Show
// @Tags     Center
// @Summary  指定楼栋id获取单元列表
// @Produce  application/json
// @Param    data  query    request.GetByIdStr  true  "楼栋 id"
// @Success  200   {object}  response.Response{data=house.DictBuilding,msg=string}  "结果"
// @Router   /base/unit [get]
func (receiver *XiaoQuApi) GetUnit(c *gin.Context) {
	var api request.GetByIdStr
	err := c.ShouldBindQuery(&api)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	qu, err2 := DictService.GetUnit(api.ID)
	if err2 != nil {
		response.FailWithMessage(err2.Error(), c)
		return
	}
	var ll []response2.DictBuildingResponse
	for _, building := range qu {
		ll = append(ll, response2.DictBuildingResponse{
			Id:   building.UnitOpenID,
			Name: building.EncryptUnitName,
		})
	}
	response.OkWithDetailed(ll, "获取成功", c)
	return
}

// Show
// @Tags     Center
// @Summary  指定单元id获取门牌列表
// @Produce  application/json
// @Param    data  query    request.GetByIdStr  true  "单元 id"
// @Success  200   {object}  response.Response{data=house.DictBuilding,msg=string}  "结果"
// @Router   /base/house [get]
func (receiver *XiaoQuApi) GetHouse(c *gin.Context) {
	var api request.GetByIdStr
	err := c.ShouldBindQuery(&api)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	qu, err2 := DictService.GetHouse(api.ID)
	if err2 != nil {
		response.FailWithMessage(err2.Error(), c)
		return
	}
	var ll []response2.DictBuildingResponse
	for _, building := range qu {
		ll = append(ll, response2.DictBuildingResponse{
			Id:   building.HouseOpenId,
			Name: building.EncryptHouseName,
		})
	}
	response.OkWithDetailed(ll, "获取成功", c)
	return
}
