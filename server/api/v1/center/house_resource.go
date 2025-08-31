package center

import (
	"context"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/house"
	response2 "github.com/flipped-aurora/gin-vue-admin/server/model/house/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/linxdeep/linxdeep-framework/pkg/searchx"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

type HouseResourceApi struct {
}

// View
// @Tags     Center
// @Summary  查看 房源
// @Produce  application/json
// @Param    data  query    string  true  "id"
// @Success  200   {object}  response.Response{data=response2.ResourceResponse}  "结果"
// @Router   /center/house/view [get]
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
	xq, err := XiaoQuService.GetInfo(info.XiaoquId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	r := response2.ResourceResponse{
		Resource:  *info,
		Latitude:  xq.Latitude,
		Longitude: xq.Longitude,
	}

	err = ResourceService.FollowViewClickAdd(uint(req.ID), "view")
	if err != nil {
		global.GVA_LOG.Error("view add failed !", zap.Error(err))
	}

	response.OkWithDetailed(r, "获取成功", c)
}

// @Tags      Center
// @Summary   房源上下架 可批量
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.HouseStateReq   true  "参数"
// @Success   200
// @Router    /center/house/state [post]
func (h *HouseResourceApi) States(c *gin.Context) {
	var req request.HouseStateReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	status := "待出租"
	if req.State == 1 {
		status = "已下架"
	}
	err = ResourceService.SetApprovalStatus(req.Ids, status)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.Ok(c)
	return
}

// View
// @Tags     Center
// @Summary  获取 房源手机号
// @Produce  application/json
// @Param    data  query    string  true  "id"
// @Success  200   {object}  response.Response{data=map[string]string}  "结果 {'mobile': '13222222222'}"
// @Router   /center/house/mobile [get]
func (h *HouseResourceApi) GetMobile(c *gin.Context) {
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

	u, err := userService.FindUserById(int(info.Owner))
	if err != nil {
		return
	}

	err = ResourceService.FollowViewClickAdd(uint(req.ID), "click")
	if err != nil {
		global.GVA_LOG.Error("view add failed !", zap.Error(err))
	}

	response.OkWithDetailed(map[string]string{"mobile": u.Phone}, "获取成功", c)
}

// @Tags      Center
// @Summary   指定查询条件  返回小区列表 包含每个小区的房源数量（聚合）
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.ResourceSearch   true  "查询条件"
// @Success   200   {object}  response.Response{data=[]response.XiaoquRsp}  "返回小区列表 包含每个小区的房源数量（聚合）"
// @Router    /center/house/xiaoquAgg [post]
func (h *HouseResourceApi) ListByXiaoquAgg(c *gin.Context) {
	var req request.ResourceSearch
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	condition := searchx.Condition{
		Terms: []searchx.Term{
			{Field: "status", Value: "待出租"},
		},
		Ors: []searchx.Condition{
			{
				Terms: []searchx.Term{},
			},
		},
		Aggs: []searchx.Agg{
			{Field: "xiaoqu_id"},
		},
	}

	if len(req.Feature) > 0 {
		for _, f := range strings.Split(req.Feature, ",") {
			condition.Terms = append(condition.Terms, searchx.Term{Field: "feature", Value: "*" + f + "*"})
		}
	}
	if len(req.HouseType) > 0 {
		condition.Terms = append(condition.Terms, searchx.Term{Field: "house_type", Value: req.HouseType + "*"})
	}

	if len(req.RentType) > 0 {
		condition.Terms = append(condition.Terms, searchx.Term{Field: "rent_type", Value: req.RentType + "*"})
	}
	if req.Price > 0 {
		priceOption := ResourceService.GetPriceByOption(strconv.Itoa(req.Price))
		g := priceOption[0]
		l := priceOption[1]
		condition.Ranges = append(condition.Ranges, searchx.Range{Field: "price", GreatEqual: fmt.Sprintf("%d", g), LessEqual: fmt.Sprintf("%d", l)})
	}

	for _, i := range req.XiaoquId {
		condition.Ors[0].Terms = append(condition.Ors[0].Terms, searchx.Term{Field: "xiaoqu_id", Value: strconv.Itoa(i)})
	}
	_, _, agg, err := global.Gva_ResourceSearch.SearchAgg(context.Background(), condition, searchx.QueryParams{
		Fields: []string{"id", "xiaoqu", "xiaoqu_id"},
		Size:   0,
	})

	var list []response.XiaoquRsp
	for id, num := range agg["xiaoqu_id"] {
		xId, _ := strconv.Atoi(id)
		var name string
		xiaoq, e := XiaoQuService.GetInfo(uint(xId))
		if e == nil {
			name = xiaoq.Name
		}

		list = append(list, response.XiaoquRsp{
			Id:   xId,
			Name: name,
			Num:  num,
		})
	}
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(list, "获取成功", c)
}

// @Tags      Center
// @Summary   指定小区id 分页获取房源列表
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.SearchResource   true  "分页获取API列表"
// @Success   200   {object}  response.Response{data=response.PageResult{list=[]response2.ResourceResponse},msg=string}  "分页获取API列表,返回包括列表,总数,页码,每页数量"
// @Router    /center/house/listByXiaoqu [post]
func (h *HouseResourceApi) ListByXiaoquId(c *gin.Context) {
	var pageInfo request.SearchResource
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

	list, total, err := ResourceService.GetPage(pageInfo.XiaoquId, 0, "", "待出租", pageInfo.PageInfo, pageInfo.OrderKey, pageInfo.Desc)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	var ids []uint
	for _, i2 := range list.([]house.Resource) {
		ids = append(ids, i2.ID)
	}

	flist, err := FavoriteService.GetByUserIdRIds(utils.GetUserID(c), ids)
	if err != nil {
		global.GVA_LOG.Error("获取失败 favorite.GetByUserIdRIds!", zap.Error(err))
	}
	fmap := make(map[uint]bool)
	for _, favorite := range flist {
		fmap[favorite.ResourceId] = true
	}
	var res []response2.ResourceResponse
	for _, i2 := range list.([]house.Resource) {
		r := response2.ResourceResponse{
			Resource: i2,
		}
		if _, ok := fmap[i2.ID]; ok {
			r.Follow = true
		}
		res = append(res, r)
	}
	response.OkWithDetailed(response.PageResult{
		List:     res,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
	return
}

// @Tags      Center
// @Summary   我发的房源列表
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.FavoriteSearch   true  "分页获取API列表"
// @Success   200   {object}  response.Response{data=response.PageResult{list=[]house.Resource},msg=string}  "分页获取API列表,返回包括列表,总数,页码,每页数量"
// @Router    /center/house/my [post]
func (h *HouseResourceApi) ListByUserId(c *gin.Context) {
	var pageInfo request.FavoriteSearch
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
	list, total, err := ResourceService.GetPage(0, userId, "", "", pageInfo.PageInfo, pageInfo.OrderKey, pageInfo.Desc)
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
// @Tags     Center
// @Summary  创建|编辑 房源
// @Produce  application/json
// @Param    data  body      house.Resource  true  "初始化内容"
// @Success  200   {object}  response.Response{data=string}  "结果"
// @Router   /center/house/create [post]
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
// @Tags     Center
// @Summary  创建|编辑 房源
// @Produce  application/json
// @Param    data  body      house.Resource  true  "初始化内容"
// @Success  200   {object}  response.Response{data=string}  "结果"
// @Router   /center/house/edit [post]
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

func (h *HouseResourceApi) Test(c *gin.Context) {
	//list, _, err := XiaoQuService.GetAreaList(request.SearchArea{CityId: "1"})
	//if err != nil {
	//	response.FailWithMessage(err.Error(), c)
	//	return
	//}
	//
	//areaMap := make(map[string]uint)
	//for _, v := range list {
	//	areaMap[v.Name] = v.ID
	//}
	//
	//listq, _, err := XiaoQuService.GetDistrictList(request.SearchDistrict{})
	//if err != nil {
	//	response.FailWithMessage(err.Error(), c)
	//	return
	//}
	//
	//districtMap := make(map[string]uint)
	//for _, v := range listq {
	//	districtMap[v.Name] = v.ID
	//}
	//
	//// 裕华区	紫荆家园	嘉华路11号	南货场	石家庄市裕华区嘉华路11号 紫荆家园	114.50723	37.976481
	//for _, xx := range strings.Split(payload, "\n") {
	//	fields := strings.Split(xx, "\t")
	//
	//	if len(fields) < 7 {
	//		continue
	//	}
	//	areaName := fields[0]
	//	name := fields[1]
	//	postion := fields[2]
	//	districts := strings.Replace(fields[3], " ", "", -1)
	//
	//	address := fields[4]
	//	latitude := fields[5]
	//	longitude := fields[6]
	//
	//	if areaMap[areaName] == 0 {
	//		fmt.Println(areaName, "不存在")
	//	}
	//	var districtsIds string
	//	for _, s := range strings.Split(districts, ",") {
	//		districtsIds += fmt.Sprintf("|%d|,", districtMap[s])
	//	}
	//	districtsIds = districtsIds[:len(districtsIds)-1]
	//
	//	fmt.Printf("insert into xiao_qu (name, area_id, area, district_ids, districts, position, address, latitude, longitude) "+
	//		"values ('%s', %d, '%s', '%s', '%s', '%s', '%s', '%v', '%v'); \n ", name, areaMap[areaName], areaName, districtsIds, districts, postion, address, latitude, longitude)
	//
	//}
	response.Ok(c)
}

// FilterArea
// @Tags     Center
// @Summary  指定城市（默认石家庄 1 ） 获取区域列表（丰台区、朝阳区）
// @Produce  application/json
// @Param cityId query string true "城市id"  default(1)
// @Success  200   {object}  response.Response{data=map[string]response.Area}  "结果"
// @Router   /center/area [get]
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
				Name:      district.Name,
				Id:        district.ID,
				Latitude:  district.Latitude,
				Longitude: district.Longitude,
			})
		}

	}
	response.OkWithDetailed(list, "获取成功", c)

	return
}

// FilterOptions
// @Tags     Center
// @Summary  筛选用到的选择项
// @Produce  application/json
// @Success  200   {object}  response.Response{data=map[string]map[string]string}  "结果"
// @Router   /center/options [get]
func (h *HouseResourceApi) FilterOptions(c *gin.Context) {

	options, err := ResourceService.FilterOptions()
	if err != nil {
		return
	}

	response.OkWithDetailed(options, "获取成功", c)

}

// @Tags      Center
// @Summary   favorite 添加房源
// @accept    application/json
// @Produce   application/json
// @Param id query string true "房源id"
// @Success   200   {object}  response.Response{msg=string}
// @Router    /center/favorite/add [get]
func (h *HouseResourceApi) FavoriteAdd(c *gin.Context) {
	var req request.GetById
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	origin, err := ResourceService.GetInfo(uint(req.ID))
	if origin == nil && err != nil {
		return
	}

	userId := utils.GetUserID(c) // 获取登陆用户
	err = FavoriteService.CreateOrUpdate(&house.Favorite{ResourceId: uint(req.ID), UserId: userId})
	if err == nil {
		err = ResourceService.FollowViewClickAdd(uint(req.ID), "follow")
		if err != nil {
			global.GVA_LOG.Error("follow add failed !", zap.Error(err))
		}
	}
	response.Ok(c)
	return
}

// @Tags      Center
// @Summary   favorite 取消房源
// @accept    application/json
// @Produce   application/json
// @Param id query string true "房源id"
// @Success   200   {object}  response.Response{msg=string}
// @Router    /center/favorite/del [get]
func (h *HouseResourceApi) FavoriteDel(c *gin.Context) {
	var req request.GetById
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	origin, err := ResourceService.GetInfo(uint(req.ID))
	if origin == nil && err != nil {
		return
	}

	userId := utils.GetUserID(c) // 获取登陆用户
	err = FavoriteService.Delete(userId, uint(req.ID))
	if err == nil {
		err = ResourceService.FollowViewClickSub(uint(req.ID), "follow")
		if err != nil {
			global.GVA_LOG.Error("follow sub failed !", zap.Error(err))
		}
	}
	response.Ok(c)
	return
}

// @Tags      Center
// @Summary   favorite 列表
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.FavoriteSearch   true  "分页获取API列表"
// @Success   200   {object}  response.Response{data=response.PageResult{list=[]house.Resource},msg=string}  "分页获取API列表,返回包括列表,总数,页码,每页数量"
// @Router    /center/favorite/list [post]
func (h *HouseResourceApi) FavoriteList(c *gin.Context) {
	var pageInfo request.FavoriteSearch
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
	list, total, err := FavoriteService.GetPage(userId, pageInfo.PageInfo, pageInfo.OrderKey, pageInfo.Desc)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	var houseIds []uint
	for _, f := range list.([]house.Favorite) {
		houseIds = append(houseIds, f.ResourceId)
	}
	var resources []*house.Resource
	if len(houseIds) > 0 {
		resources, err = ResourceService.GetListByIds(houseIds)
		if err != nil {
			return
		}
	}

	response.OkWithDetailed(response.PageResult{
		List:     resources,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
	return
}
