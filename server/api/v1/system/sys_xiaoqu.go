package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type XiaoQuApi struct{}

// Edit
// @Tags     XiaoQu
// @Summary  编辑单个小区
// @Produce  application/json
// @Param    data  body      system.XiaoQu  true  "提交内容"
// @Success  200   {object}  response.Response{data=string}  "结果"
// @Router   /api/xiaoqu/edit [post]
func (receiver *XiaoQuApi) Edit(c *gin.Context) {
	var qu system.XiaoQu
	var xq system.XiaoQu
	err := c.ShouldBindJSON(&qu)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	//err = utils.Verify(qu, utils.ApiVerify)
	//if err != nil {
	//	response.FailWithMessage(err.Error(), c)
	//	return
	//}

	if qu.ID == 0 {
		xq, err = xiaoQuService.Create(qu)
	} else {
		xq, err = xiaoQuService.Edit(qu)
	}

	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(xq, c)
	return
}

// Show
// @Tags     XiaoQu
// @Summary  获取小区信息
// @Produce  application/json
// @Param    data  query    request.GetById  true  "id"
// @Success  200   {object}  response.Response{data=system.XiaoQu,msg=string}  "结果"
// @Router   /api/xiaoqu/show [get]
func (receiver *XiaoQuApi) Show(c *gin.Context) {
	var api request.GetById
	err := c.ShouldBindQuery(&api)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	qu, err2 := xiaoQuService.GetInfo(uint(api.ID))
	if err2 != nil {
		response.FailWithMessage(err2.Error(), c)
		return
	}
	response.OkWithDetailed(qu, "获取成功", c)
	return
}

// List
// @Tags     XiaoQu
// @Summary  获取小区列表
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.SearchXiaoqu    true  "页码, 每页大小"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "分页列表,返回包括列表,总数,页码,每页数量"
// @Router    /api/xiaoqu/list [post]
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
	list, total, err := xiaoQuService.GetList(pageInfo)
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
