package resource

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/house"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type HouseResourceApi struct {
}

// View
// @Tags     Admin
// @Summary   房源详情
// @Produce  application/json
// @Param    data  query    string  true  "id"
// @Success  200   {object}  response.Response{data=house.Resource}  "结果"
// @Router   /api/house/view [get]
func (h *HouseResourceApi) View(c *gin.Context) {
	var req request.GetById
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	info, err := ResourceService.GetInfo(uint(req.ID))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithDetailed(info, "获取成功", c)
}

// @Tags      Admin
// @Summary   房源上下架 可批量
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.HouseStateReq   true  "分页获取API列表"
// @Success   200
// @Router    /house/state [post]
func (h *HouseResourceApi) States(c *gin.Context) {
	var req request.HouseStateReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	status := "待出租"
	if req.State == 2 {
		status = "已下架"
	}
	err = ResourceService.SetState(req.Ids, status)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.Ok(c)
	return
}

// @Tags      Admin
// @Summary   批量审核房源
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.HouseStateReq   true  "分页获取API列表"
// @Success   200
// @Router    /house/approvalState [post]
func (h *HouseResourceApi) ApprovalStatus(c *gin.Context) {
	var req request.HouseStateReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	approvalStatus := "未通过"
	if req.State == 1 {
		approvalStatus = "通过"
	}
	err = ResourceService.SetApprovalStatus(req.Ids, approvalStatus)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.Ok(c)
	return
}

// @Tags      Admin
// @Summary   房源审核列表
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.SearchNameResource   true  "分页获取API列表"
// @Success   200   {object}  response.Response{data=response.PageResult{list=[]house.Resource},msg=string}  "分页获取API列表,返回包括列表,总数,页码,每页数量"
// @Router    /api/house/list [post]
func (h *HouseResourceApi) List(c *gin.Context) {
	var pageInfo request.SearchNameResource
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if pageInfo.PageInfo.Page == 0 {
		pageInfo.PageInfo.Page = 1
	}

	if pageInfo.PageInfo.PageSize == 0 {
		pageInfo.PageInfo.PageSize = 50
	}
	var (
		list  interface{}
		total int64
		uId   uint
	)

	//if len(pageInfo.Phone) > 0 {
	//	u := UserService.FindUserByMobile(pageInfo.Phone)
	//	if u != nil {
	//		uId = u.ID
	//	} else {
	//		response.FailWithMessage("经纪人手机号不存在", c)
	//		return
	//	}
	//}

	list, total, err = ResourceService.GetPage(pageInfo.XiaoquId, uId, pageInfo.ApprovalStatus, "", pageInfo.PageInfo, pageInfo.OrderKey, pageInfo.Desc, request.SearchOther{Phone: pageInfo.Phone})
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

	return
}

// @Tags      Admin
// @Summary   删除房源
// @accept    application/json
// @Produce   application/json
// @Param    data  query    string  true  "id"
// @Success   200   {object}  response.Response{data=string}  "结果"
// @Router    /api/house/del [post]
func (h *HouseResourceApi) DeleteByUserId(c *gin.Context) {
	var req request.GetById
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	userId := utils.GetUserID(c) // 获取登陆用户

	err = ResourceService.DelByUser(uint(req.ID), userId)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.Ok(c)
	return
}

// @Tags      Admin
// @Summary   我发的房源列表
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.MySearchResource   true  "分页获取API列表"
// @Success   200   {object}  response.Response{data=response.PageResult{list=[]house.Resource},msg=string}  "分页获取API列表,返回包括列表,总数,页码,每页数量"
// @Router    /api/house/my [post]
func (h *HouseResourceApi) ListByUserId(c *gin.Context) {
	var pageInfo request.MySearchResource
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if pageInfo.PageInfo.Page == 0 {
		pageInfo.PageInfo.Page = 1
	}

	if pageInfo.PageInfo.PageSize == 0 {
		pageInfo.PageInfo.PageSize = 50
	}
	userId := utils.GetUserID(c) // 获取登陆用户
	pageInfo.PageInfo.Keyword = pageInfo.DoorNo
	list, total, err := ResourceService.GetPage(pageInfo.XiaoquId, userId, "", pageInfo.Status, pageInfo.PageInfo, pageInfo.OrderKey, pageInfo.Desc, request.SearchOther{})
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
	return
}

// Create
// @Tags     Admin
// @Summary  创建房源
// @Produce  application/json
// @Param    data  body      house.Resource  true  "初始化内容"
// @Success  200   {object}  response.Response{data=string}  "结果"
// @Router   /api/house/create [post]
func (h *HouseResourceApi) Create(c *gin.Context) {
	var req house.Resource
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	req.Owner = utils.GetUserID(c) // 获取登陆用户
	if req.XiaoquId == 0 {
		response.FailWithMessage("小区id不能为空", c)
		return
	}

	xiaoqu, err := XiaoQuService.GetInfo(req.XiaoquId)
	if err != nil {
		response.FailWithMessage("小区不存在", c)
		return
	}

	req.Xiaoqu = xiaoqu.Name
	req.Districts = xiaoqu.Districts
	req.DistrictIds = xiaoqu.DistrictIds
	req.City = xiaoqu.City
	req.Area = xiaoqu.Area
	req.Status = "待出租"
	err = ResourceService.CreateOrUpdate(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	} else {
		global.GVA_Lock.Lock()
		defer global.GVA_Lock.Unlock()
		xx, _ := XiaoQuService.GetInfo(req.XiaoquId)
		xx.HouseNum++
		err = XiaoQuService.EditNum(xx.ID, xx.HouseNum)
		if err != nil {
			global.GVA_LOG.Error(err.Error())
			return
		}
	}
	response.Ok(c)
}

// Edit
// @Tags     Admin
// @Summary  编辑 房源
// @Produce  application/json
// @Param    data  body      house.Resource  true  "初始化内容"
// @Success  200   {object}  response.Response{data=string}  "结果"
// @Router   /api/house/edit [post]
func (h *HouseResourceApi) Edit(c *gin.Context) {
	var req house.Resource
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	origin, err := ResourceService.GetInfo(req.ID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	origin.Feature = req.Feature
	origin.Price = req.Price
	origin.Attachments = req.Attachments

	err = ResourceService.CreateOrUpdate(origin)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.Ok(c)
}

// FilterArea
// @Tags     Admin
// @Summary  获取区域列表（丰台区、朝阳区）
// @Produce  application/json
// @Param    data  query    string  true  "cityId"
// @Success  200   {object}  response.Response{data=map[string]response.Area}  "结果"
// @Router   /api/house/area [get]
func (h *HouseResourceApi) FilterArea(c *gin.Context) {
	cityId := c.Query("cityId")

	areaList, _, err := XiaoQuService.GetAreaList(request.SearchArea{CityId: cityId})
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list := make(map[string]*response.Area)

	for _, area := range areaList {
		list[area.Name] = &response.Area{
			Name: area.Name,
			Id:   int(area.ID),
			Sort: area.Sort,
		}

		districtList, _, err := XiaoQuService.GetDistrictList(request.SearchDistrict{AreaId: area.ID})
		if err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
		for _, district := range districtList {
			list[area.Name].Districts = append(list[area.Name].Districts, response.Districts{
				Name: district.Name,
				Id:   district.ID,
			})
		}

	}
	response.OkWithDetailed(list, "获取成功", c)

	return
}

// FilterOptions
// @Tags     Admin
// @Summary  筛选用到的选择项
// @Produce  application/json
// @Success  200   {object}  response.Response{data=map[string]map[string]string}  "结果"
// @Router   /api/house/options [get]
func (h *HouseResourceApi) FilterOptions(c *gin.Context) {

	options, err := ResourceService.FilterOptions()
	if err != nil {
		return
	}

	response.OkWithDetailed(options, "获取成功", c)

}

// FilterTypeOptions
// @Tags     Admin
// @Summary  new 房型筛选用到的选择项
// @Produce  application/json
// @Success  200   {object}  response.Response{data=map[string]interface{}}  "结果"
// @Router   /api/house/type/options [get]
func (h *HouseResourceApi) FilterTypeOptions(c *gin.Context) {
	options, err := ResourceService.FilterOptions1()
	if err != nil {
		return
	}

	response.OkWithDetailed(map[string]interface{}{
		"houseType": options,
		"price":     map[string]string{"1": "500以下", "2": "500-1000元", "3": "1000-1500元", "4": "1500-2000元", "5": "2000-2500元", "6": "2500-3000元", "7": "3000元以上"},
	}, "获取成功", c)

}
