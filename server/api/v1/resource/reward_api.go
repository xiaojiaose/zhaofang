package resource

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
)

// RewardAdminList
// @Tags     Admin
// @Summary  后台成交有礼审核列表
// @Accept   application/json
// @Produce  application/json
// @Param    data  body      request.RewardApplicationSearch  true  "查询参数"
// @Success  200   {object}  response.Response{data=response.PageResult,msg=string}  "审核列表"
// @Router   /api/house/reward/list [post]
func (h *HouseResourceApi) RewardAdminList(c *gin.Context) {
	var req request.RewardApplicationSearch
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}
	list, total, err := RewardService.GetPageForAdmin(req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(response.PageResult{List: list, Total: total, Page: req.Page, PageSize: req.PageSize}, c)
}

// RewardAdminAction
// @Tags     Admin
// @Summary  后台操作成交有礼审核状态
// @Accept   application/json
// @Produce  application/json
// @Param    data  body      request.RewardApplicationAction  true  "操作参数"
// @Success  200   {object}  response.Response{data=string,msg=string}  "操作结果"
// @Router   /api/house/reward/action [post]
func (h *HouseResourceApi) RewardAdminAction(c *gin.Context) {
	var req request.RewardApplicationAction
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := RewardService.AdminAction(req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("操作成功", c)
}

// ContactQuotaGrant
// @Tags     Admin
// @Summary  增加经纪人联系方式查看次数
// @Accept   application/json
// @Produce  application/json
// @Param    data  body      request.ContactQuotaCreate  true  "增加次数参数"
// @Success  200   {object}  response.Response{data=string,msg=string}  "操作结果"
// @Router   /api/house/contactQuota/grant [post]
func (h *HouseResourceApi) ContactQuotaGrant(c *gin.Context) {
	var req request.ContactQuotaCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := ContactQuotaService.GrantByPhone(utils.GetUserID(c), req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("添加成功", c)
}

// ContactQuotaList
// @Tags     Admin
// @Summary  查看联系方式次数流水
// @Accept   application/json
// @Produce  application/json
// @Param    data  body      request.ContactQuotaSearch  true  "查询参数"
// @Success  200   {object}  response.Response{data=response.PageResult,msg=string}  "次数流水"
// @Router   /api/house/contactQuota/list [post]
func (h *HouseResourceApi) ContactQuotaList(c *gin.Context) {
	var req request.ContactQuotaSearch
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}
	list, total, err := ContactQuotaService.GetPage(req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(response.PageResult{List: list, Total: total, Page: req.Page, PageSize: req.PageSize}, c)
}
