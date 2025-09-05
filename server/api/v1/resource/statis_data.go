package resource

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/house"
	response2 "github.com/flipped-aurora/gin-vue-admin/server/model/house/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/search"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type StatisDataApi struct {
}

// View
// @Tags     Admin
// @Summary   数据中心
// @Produce  application/json
// @Param    data  query   request.GetStatis  true   "start, end"
// @Success  200   {object}  response.Response{data=search.StatisData}  "结果"
// @Router   /api/statis/view [get]
func (s *StatisDataApi) View(c *gin.Context) {
	var req request.GetStatis
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, err := StatisService.ByDate(req.Start, req.End)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	var re search.StatisData
	for _, data := range list {
		re.View += data.View
		re.Add += data.Add
		re.AddSaler += data.AddSaler
		re.UseSaler += data.UseSaler
		re.Click += data.Click
		re.Follow += data.Follow
		re.Shared += data.Shared
	}
	response.OkWithData(re, c)
	return
}

// View
// @Tags     Admin
// @Summary   访问数据
// @Produce  application/json
// @Param    data  query   request.VisitReq  true   "start, end"
// @Success  200   {object}  response.Response{data=request.VisitResponse}  "结果"
// @Router   /api/statis/visit [get]
func (s *StatisDataApi) VisitRecord(c *gin.Context) {
	var req request.VisitReq
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	var (
		u      *system.SysUser
		userId uint
	)

	if req.Phone != "" {
		u = UserService.FindUserByMobile(req.Phone)
		if u != nil {
			userId = u.ID
		}

	} else if req.WxNo != "" {
		u = UserService.FindUserByWxNo(req.WxNo)
		if u != nil {
			userId = u.ID
		}
	}

	list, err := StatisService.VisitRecord(userId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	var userIds []uint
	for _, data := range list {
		userIds = append(userIds, data.UserId)
	}
	userMap := make(map[uint]system.SysUser)
	if len(userIds) > 0 {
		users, _ := UserService.GetUsersByIds(userIds)
		for i, user := range users {
			userMap[userIds[i]] = user
		}
	}

	var re []request.VisitResponse
	for _, data := range list {
		r := request.VisitResponse{
			Date: data.Date,
		}
		if u, ok := userMap[data.UserId]; ok {
			r.WxNo = u.WxNo
			r.Phone = u.Phone
		}
		re = append(re, r)
	}
	response.OkWithData(re, c)
	return
}

// View
// @Tags     Admin
// @Summary   帖子数据
// @Produce  application/json
// @Param     data  body      request.SearchNameResource   true  "分页获取API列表"
// @Success   200   {object}  response.Response{data=response.PageResult{list=[]response2.ResourceVisitResponse},msg=string}  "分页获取API列表,返回包括列表,总数,页码,每页数量"
// @Router   /api/statis/visit/house [post]
func (s *StatisDataApi) VisitHouse(c *gin.Context) {
	var req request.SearchHouseResource
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if req.PageInfo.Page == 0 {
		req.PageInfo.Page = 1
	}

	if req.PageInfo.PageSize == 0 {
		req.PageInfo.PageSize = 50
	}
	var userId uint
	if len(req.WxNo) > 0 {
		u := UserService.FindUserByWxNo(req.WxNo)
		if u != nil {
			userId = u.ID
		}
	}
	list, total, err := ResourceService.GetPage(req.XiaoquId, userId, req.ApprovalStatus, "", req.PageInfo, req.OrderKey, req.Desc,
		request.SearchOther{
			Phone:          req.Phone,
			UpdatedAtLast:  req.UpdatedAtLast,
			UpdatedAtStart: req.UpdatedAtStart,
			HasPic:         req.HasPic,
		},
	)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	var userIds []uint
	for _, data := range list.([]house.Resource) {
		userIds = append(userIds, data.Owner)
	}
	userMap := make(map[uint]system.SysUser)
	if len(userIds) > 0 {
		users, _ := UserService.GetUsersByIds(userIds)
		for i, user := range users {
			userMap[userIds[i]] = user
		}
	}

	var ll []response2.ResourceVisitResponse
	for _, data := range list.([]house.Resource) {
		r := response2.ResourceVisitResponse{
			Resource: data,
		}
		if u, ok := userMap[data.Owner]; ok {
			r.WxNo = u.WxNo
			r.Phone = u.Phone
		}
		ll = append(ll, r)
	}
	response.OkWithDetailed(response.PageResult{
		List:     ll,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, "获取成功", c)
	return
}
