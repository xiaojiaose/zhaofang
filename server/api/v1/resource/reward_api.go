package resource

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
)

// RewardAdminList
// @Tags     Admin
// @Summary  [新增] 后台成交有礼审核列表
// @Description [新增接口] 后台只展示“发布人已确认”的成交有礼申请单。
// @Accept   application/json
// @Produce  application/json
// @Param    data  body      request.RewardApplicationSearch  true  "查询参数"
// @Success  200   {object}  response.Response{data=response.PageResult,msg=string}  "审核列表"
// @Router   /api/house/reward/list [post]
func (h *HouseResourceApi) RewardAdminList(c *gin.Context) {
	// 后台列表只接“发布人已确认”的数据，具体过滤逻辑放在 service 层统一处理。
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
// @Summary  [新增] 后台操作成交有礼审核状态
// @Description [新增接口] 支持审核中、审通过待发放、已发放、未通过等后台状态流转。
// @Accept   application/json
// @Produce  application/json
// @Param    data  body      request.RewardApplicationAction  true  "操作参数"
// @Success  200   {object}  response.Response{data=string,msg=string}  "操作结果"
// @Router   /api/house/reward/action [post]
func (h *HouseResourceApi) RewardAdminAction(c *gin.Context) {
	// handler 只负责参数接收和响应，状态映射统一收口到 service，
	// 这样后台按钮文案和数据库状态之间的映射不会散落多处。
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
// @Summary  [新增] 增加经纪人联系方式查看次数
// @Description [新增接口] 按经纪人手机号给账号充值房东房源联系方式查看次数。
// @Accept   application/json
// @Produce  application/json
// @Param    data  body      request.ContactQuotaCreate  true  "增加次数参数"
// @Success  200   {object}  response.Response{data=string,msg=string}  "操作结果"
// @Router   /api/house/contactQuota/grant [post]
func (h *HouseResourceApi) ContactQuotaGrant(c *gin.Context) {
	// 次数充值按手机号查目标账号，便于后台运营按经纪人手机号直接操作。
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
// @Summary  [新增] 查看联系方式次数流水
// @Description [新增接口] 返回充值和消耗流水，便于后台核对剩余次数和使用情况。
// @Accept   application/json
// @Produce  application/json
// @Param    data  body      request.ContactQuotaSearch  true  "查询参数"
// @Success  200   {object}  response.Response{data=response.PageResult,msg=string}  "次数流水"
// @Router   /api/house/contactQuota/list [post]
func (h *HouseResourceApi) ContactQuotaList(c *gin.Context) {
	// 这里返回流水而不是汇总，便于后台同时查看充值记录和消耗记录。
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

// LandlordContactViewList
// @Tags     Admin
// @Summary  [新增] 房东房源联系方式查看记录列表
// @Description [新增接口] 房东房源被有权限用户查看联系方式后会自动生成记录，后台可在此分页查看和筛选。
// @Accept   application/json
// @Produce  application/json
// @Param    data  body      request.LandlordContactViewSearch  true  "查询参数"
// @Success  200   {object}  response.Response{data=response.PageResult{list=[]response.LandlordContactViewResponse},msg=string}  "记录列表"
// @Router   /api/house/landlordContactView/list [post]
func (h *HouseResourceApi) LandlordContactViewList(c *gin.Context) {
	var req request.LandlordContactViewSearch
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
	list, total, err := ContactQuotaService.GetLandlordContactViewPage(req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(response.PageResult{List: list, Total: total, Page: req.Page, PageSize: req.PageSize}, c)
}

// LandlordContactViewAction
// @Tags     Admin
// @Summary  [新增] 更新房东房源联系方式查看记录状态
// @Description [新增接口] 支持待审核记录流转到审核中、已通过未付款、已付款、未通过已拒绝。
// @Accept   application/json
// @Produce  application/json
// @Param    data  body      request.LandlordContactViewAction  true  "操作参数"
// @Success  200   {object}  response.Response{data=string,msg=string}  "操作结果"
// @Router   /api/house/landlordContactView/action [post]
func (h *HouseResourceApi) LandlordContactViewAction(c *gin.Context) {
	var req request.LandlordContactViewAction
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := ContactQuotaService.UpdateLandlordContactViewStatus(req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("操作成功", c)
}
