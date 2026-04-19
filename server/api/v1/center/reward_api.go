package center

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	response2 "github.com/flipped-aurora/gin-vue-admin/server/model/house/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

// RewardRecent
// @Tags     Center
// @Summary  [新增] 获取最近联系过的房源发布人记录
// @Description [新增接口] 基于 /center/house/mobile 的联系方式点击记录，返回当前登录用户最近联系过的房源发布人。
// @Produce  application/json
// @Success  200  {object}  response.Response{data=[]map[string]interface{},msg=string}  "最近联系记录"
// @Router   /center/reward/recent [get]
func (h *HouseResourceApi) RewardRecent(c *gin.Context) {
	// 小程序首页“出房有礼”入口展示的是最近联系过的发布人，
	// 数据来源于联系方式点击记录，而不是用户主动维护的收藏/通讯录。
	list, err := RewardService.RecentContacts(utils.GetUserID(c))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(list, c)
}

// RewardApply
// @Tags     Center
// @Summary  [新增] 发起出房有礼申请
// @Description [新增接口] 申请人提交房源和备注后，后端会自动补齐申请人/发布人的手机号和微信号快照。
// @Accept   application/json
// @Produce  application/json
// @Param    data  body      request.RewardApplicationCreate  true  "申请参数"
// @Success  200   {object}  response.Response{data=string,msg=string}  "申请结果"
// @Router   /center/reward/apply [post]
func (h *HouseResourceApi) RewardApply(c *gin.Context) {
	// 前端这里只提交房源和备注；
	// 发布人微信号、申请人微信号等展示数据由后端从用户资料中补快照。
	var req request.RewardApplicationCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := RewardService.Create(utils.GetUserID(c), req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("申请成功", c)
}

// RewardPublisherList
// @Tags     Center
// @Summary  [新增] 发布人查看出房有礼审核列表
// @Description [新增接口] 房源发布人查看自己的申请单列表，不等同于后台总审核池。
// @Accept   application/json
// @Produce  application/json
// @Param    data  body      request.RewardApplicationSearch  true  "查询参数"
// @Success  200   {object}  response.Response{data=response.PageResult,msg=string}  "审核列表"
// @Router   /center/reward/publisher/list [post]
func (h *HouseResourceApi) RewardPublisherList(c *gin.Context) {
	// 这是“房源发布人视角”的审核列表，不是后台总审核池。
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
	list, total, err := RewardService.GetPageForPublisher(utils.GetUserID(c), req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(response.PageResult{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, c)
}

// RewardPublisherAction
// @Tags     Center
// @Summary  [新增] 发布人确认或拒绝出房有礼申请
// @Description [新增接口] 发布人确认后，申请单才会进入后台审核状态。
// @Accept   application/json
// @Produce  application/json
// @Param    data  body      request.RewardApplicationAction  true  "操作参数"
// @Success  200   {object}  response.Response{data=string,msg=string}  "操作结果"
// @Router   /center/reward/publisher/action [post]
func (h *HouseResourceApi) RewardPublisherAction(c *gin.Context) {
	// 发布人只允许做确认/拒绝，真正的后台审核动作不走这个接口。
	var req request.RewardApplicationAction
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := RewardService.PublisherAction(utils.GetUserID(c), req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("操作成功", c)
}

// CreateShare
// @Tags     Center
// @Summary  [新增] 创建我的房源地图分享链接
// @Description [新增接口] 为当前登录用户生成 7 天有效的房源地图分享 token。
// @Accept   application/json
// @Produce  application/json
// @Param    data  body      request.ResourceShareCreate  false  "分享参数"
// @Success  200   {object}  response.Response{data=map[string]interface{},msg=string}  "分享结果"
// @Router   /center/house/share [post]
func (h *HouseResourceApi) CreateShare(c *gin.Context) {
	// 分享链接允许前端不传有效期，后端统一兜底为 7 天。
	var req request.ResourceShareCreate
	_ = c.ShouldBindJSON(&req)
	entity, err := ShareService.Create(utils.GetUserID(c), req.ExpireDays)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(entity, c)
}

// SharedMap
// @Tags     Center
// @Summary  [新增] 通过分享token获取地图点位列表
// @Description [新增接口] 未登录也可访问；先返回按小区聚合后的地图点位，前端点击点位后再调用 /center/house/share/list 查看该小区下的房源。
// @Produce  application/json
// @Param    token  query     string  true  "分享token"
// @Success  200    {object}  response.Response{data=response2.SharedMapResponse,msg=string}  "分享房源地图点位"
// @Router   /center/house/share [get]
func (h *HouseResourceApi) SharedMap(c *gin.Context) {
	// 分享页未登录可访问，所以这里只依赖 token，不读登录态。
	token := c.Query("token")
	entity, err := ShareService.GetByToken(token)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, err := ResourceService.SharedMapAgg(entity.UserID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	user, _ := userService.FindUserById(int(entity.UserID))
	resp := response2.SharedMapResponse{
		List:     list,
		ExpireAt: entity.ExpireAtUnix,
	}
	if user != nil {
		resp.UserInfo = *user
	}
	response.OkWithData(resp, c)
}

// SharedMapList
// @Tags     Center
// @Summary  [新增] 通过分享token和小区ID获取点位下房源列表
// @Description [新增接口] 地图点位点击后调用，仍然通过 ZincSearch 按 owner + xiaoqu_id 过滤，再回表补全房源详情。
// @Produce  application/json
// @Param    token    query     string  true  "分享token"
// @Param    xiaoquId  query     int     true  "小区ID"
// @Param    page      query     int     false "页码"
// @Param    pageSize  query     int     false "每页数量"
// @Success  200    {object}  response.Response{data=response.PageResult{list=[]house.Resource},msg=string}  "分享小区房源列表"
// @Router   /center/house/share/list [get]
func (h *HouseResourceApi) SharedMapList(c *gin.Context) {
	token := c.Query("token")
	entity, err := ShareService.GetByToken(token)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	var req request.PageInfo
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	xiaoquID := c.Query("xiaoquId")
	xID, _ := strconv.Atoi(xiaoquID)
	list, total, err := ResourceService.SharedMapList(entity.UserID, uint(xID), req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(response.PageResult{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, c)
}
