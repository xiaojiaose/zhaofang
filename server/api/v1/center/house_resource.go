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

var houseType = map[string]string{
	//"房东房源": "房东房源",
	"1居":  "1居",
	"2居":  "2居",
	"3居":  "3居",
	"4居+": "4居+",
	"开间": "开间",
	"主卧": "主卧",
	"次卧": "次卧",
	"暗间": "暗间",
}

// View
// @Tags     Center
// @Summary  [变更] 查看房源详情（房东房源隐藏门牌号）
// @Description [变更接口] 房东房源详情页不再返回门牌号，避免地址暴露过细。
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
	if info.HouseType == "房东房源" {
		// 房东房源详情页不展示门牌号，避免把地址暴露得过细。
		r.DoorNo = ""
	}

	flist, err := FavoriteService.GetByUserIdRIds(utils.GetUserID(c), []uint{uint(req.ID)})
	if err != nil {
		global.GVA_LOG.Error("获取失败 favorite.GetByUserIdRIds!", zap.Error(err))
	}
	for _, favorite := range flist {
		if favorite.ResourceId == uint(req.ID) {
			r.Follow = true
		}
	}

	err = ResourceService.FollowViewClickAdd(uint(req.ID), "view")
	if err != nil {
		global.GVA_LOG.Error("view add failed !", zap.Error(err))
	}

	StatisService.InsertRecord(uint(req.ID), "view", utils.GetUserID(c))

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
	err = ResourceService.SetState(req.Ids, status)
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
// @Summary  [变更] 获取房源手机号（房东房源扣查看次数）
// @Description [变更接口] 房东房源查看联系方式前，会先校验并扣减当前登录用户的可用查看次数。
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
	if len(info.Phone) == 0 {
		u, err := userService.FindUserById(int(info.Owner))
		if err != nil {
			return
		}
		info.Phone = u.Phone
	}
	if isLandlordResource(*info) {
		// 房东房源联系方式需要消耗当前登录用户的可用查看次数，
		// 扣减成功后会自动生成一条“联系方式查看记录”（默认待审核）供后台展示。
		err = ContactQuotaService.Consume(utils.GetUserID(c), uint(req.ID), "查看房东房源联系方式")
		if err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
	}

	err = ResourceService.FollowViewClickAdd(uint(req.ID), "click")
	if err != nil {
		global.GVA_LOG.Error("view add failed !", zap.Error(err))
	}
	StatisService.InsertRecord(uint(req.ID), "click", utils.GetUserID(c))
	response.OkWithDetailed(map[string]string{"mobile": info.Phone}, "获取成功", c)
}

func isLandlordResource(resource house.Resource) bool {
	// 只保留  rent_type 就行了。可以先这样 没啥影响
	// 历史数据里“房东房源”可能落在 house_type 或 rent_type，
	// 且可能出现前后缀/空格，统一用包含匹配做兼容。
	houseTyp := strings.TrimSpace(resource.HouseType)
	rentType := strings.TrimSpace(resource.RentType)
	return strings.Contains(houseTyp, "房东房源") || strings.Contains(rentType, "房东房源")
}
func applyHouseSourceCondition(condition *searchx.Condition, houseSource string, allowTeam bool) {
	for _, item := range strings.Split(houseSource, ",") {
		switch strings.TrimSpace(item) {
		case "commission", "有返佣":
			condition.Ranges = append(condition.Ranges, searchx.Range{Field: "commission_price", GreatEqual: "1"})
		case "landlord", "房东房源":
			condition.Terms = append(condition.Terms, searchx.Term{Field: "rent_type", Value: "房东房源"})
		case "team", "团队房源":
			if allowTeam {
				condition.Terms = append(condition.Terms, searchx.Term{Field: "is_team_house", Value: "1"})
			}
		}
	}
}

// @Tags      Center
// @Summary   [变更] 地图聚合查询房源小区列表
// @Description [变更接口] 支持返佣、房东房源、团队房源筛选；普通用户后端默认过滤团队房源。
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
	global.GVA_LOG.Info("xiaoquAgg 参数", zap.String("req", fmt.Sprintf("%+v", req)))

	condition := searchx.Condition{
		Terms: []searchx.Term{
			{Field: "status", Value: "待出租"},
		},
		Ors: []searchx.Condition{
			{
				Terms: []searchx.Term{},
			},
		},
		Nots: []searchx.Condition{
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
		for _, v := range houseType {
			if strings.Contains(req.HouseType, v) {
				continue
			}
			condition.Nots[0].Terms = append(condition.Nots[0].Terms, searchx.Term{Field: "house_type", Value: v})
		}
	}

	if len(req.RentType) > 0 {
		condition.Terms = append(condition.Terms, searchx.Term{Field: "rent_type", Value: req.RentType + "*"})
	}
	if !userHasTeamPermission(utils.GetUserID(c)) {
		// 普通用户即使前端手工构造参数，也不能看到团队房源，这里做后端兜底过滤。
		condition.Terms = append(condition.Terms, searchx.Term{Field: "is_team_house", Value: "0"})
	}

	if len(req.HouseSource) > 0 {
		applyHouseSourceCondition(&condition, req.HouseSource, userHasTeamPermission(utils.GetUserID(c)))
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
// @Summary   [变更] 地图查询指定条件房源列表
// @Description [变更接口] 支持返佣、房东房源、团队房源筛选；普通用户后端默认过滤团队房源。
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.ResourceSearch   true  "查询条件"
// @Success   200   {object}  response.Response{data=response.PageResult{list=[]response2.ResourceResponse},msg=string}  "指定查询条件  返回指定小区房源列表"
// @Router    /center/house/xiaoquAggList [post]
func (h *HouseResourceApi) ListByXiaoquAggList(c *gin.Context) {
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
		Nots: []searchx.Condition{
			{
				Terms: []searchx.Term{},
			},
		},
	}

	if len(req.Feature) > 0 {
		for _, f := range strings.Split(req.Feature, ",") {
			condition.Terms = append(condition.Terms, searchx.Term{Field: "feature", Value: "*" + f + "*"})
		}
	}

	if len(req.HouseType) > 0 {
		for _, v := range houseType {
			if strings.Contains(req.HouseType, v) {
				continue
			}
			condition.Nots[0].Terms = append(condition.Nots[0].Terms, searchx.Term{Field: "house_type", Value: v})
		}
	}

	if len(req.RentType) > 0 {
		condition.Terms = append(condition.Terms, searchx.Term{Field: "rent_type", Value: req.RentType + "*"})
	}
	if !userHasTeamPermission(utils.GetUserID(c)) {
		// 聚合接口和列表接口都需要做同样的权限兜底，避免两边数据口径不一致。
		condition.Terms = append(condition.Terms, searchx.Term{Field: "is_team_house", Value: "0"})
	}
	if len(req.HouseSource) > 0 {
		applyHouseSourceCondition(&condition, req.HouseSource, userHasTeamPermission(utils.GetUserID(c)))
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

	if req.PageSize == 0 {
		req.PageSize = 50
	}
	if req.Page > 0 {
		req.Page = req.Page - 1
	}
	houseList, total, err := global.Gva_ResourceSearch.Search(context.Background(), condition, searchx.QueryParams{
		Fields: []string{"house_id", "xiaoqu", "xiaoqu_id", "_id"},
		Size:   req.PageSize,
		Page:   req.Page,
		Asc:    false,
		By:     "last_update",
	})

	var houseIds []uint
	for _, buffer := range houseList {
		x, _ := strconv.ParseUint(buffer.HouseId, 10, 0)
		houseIds = append(houseIds, uint(x))
	}

	var resources []*house.Resource
	if len(houseIds) > 0 {
		resources, err = ResourceService.GetListByIdsSafe(houseIds, "待出租", userHasTeamPermission(utils.GetUserID(c)))
		if err != nil {
			return
		}
	}

	response.OkWithDetailed(response.PageResult{
		List:     resources,
		Total:    int64(total),
		Page:     req.Page,
		PageSize: req.PageSize,
	}, "获取成功", c)
}

// @Tags      Center
// @Summary   [变更] 分页获取指定小区房源列表
// @Description [变更接口] 普通用户查询时默认排除团队房源，保持与地图筛选口径一致。
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

	other := request.SearchOther{}
	if !userHasTeamPermission(utils.GetUserID(c)) {
		// 走数据库分页的接口同样排除团队房源，保持与 ES 检索结果一致。
		other.IsTeamHouse = "false"
	}
	list, total, err := ResourceService.GetPage(pageInfo.XiaoquId, 0, "", "待出租", pageInfo.PageInfo, "updated_last_at", true, other)
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
// @Summary   [变更] 我的房源列表
// @Description [变更接口] 返回微信资料、返佣金额、已上架数量和剩余可上架数量。
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
	list, total, err := ResourceService.GetPage(0, userId, "", "", pageInfo.PageInfo, pageInfo.OrderKey, pageInfo.Desc, request.SearchOther{})
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	user, _ := userService.FindUserById(int(userId))
	var result []response2.MyResourceResponse
	publishUsed := 0
	if user != nil {
		count, _ := ResourceService.CountOnShelfByUser(userId)
		publishUsed = int(count)
	}
	for _, item := range list.([]house.Resource) {
		result = append(result, response2.MyResourceResponse{
			Resource:           item,
			WxNo:               user.WxNo,
			WxNickName:         user.WxNickName,
			HeaderImg:          user.HeaderImg,
			PublishQuotaTotal:  user.PublishQuotaTotal,
			PublishQuotaUsed:   publishUsed,
			PublishQuotaRemain: maxInt(user.PublishQuotaTotal-publishUsed, 0),
		})
	}
	response.OkWithDetailed(response.PageResult{
		List:     result,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
	return
}

// Create
// @Tags     Center
// @Summary  [变更] 创建房源
// @Description [变更接口] 支持房东房源类型、返佣金额，并自动联动团队房源标识和发布人手机号。
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

	loginUser := utils.GetUserInfo(c)
	req.Owner = loginUser.BaseClaims.ID // 获取登陆用户
	req.Phone = loginUser.BaseClaims.Mobile

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
	req.Region = xiaoqu.Area
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

// @Tags      Center
// @Summary   删除房源
// @accept    application/json
// @Produce   application/json
// @Param    data  query    string  true  "id"
// @Success   200   {object}  response.Response{data=string}  "结果"
// @Router    /center/house/del [post]
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

// Edit
// @Tags     Center
// @Summary  [变更] 编辑房源
// @Description [变更接口] 支持修改房东房源类型、返佣金额和房源补充字段。
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
	origin.CommissionPrice = req.CommissionPrice
	origin.Attachments = req.Attachments
	origin.HouseType = req.HouseType
	origin.RentType = req.RentType
	origin.Remarks = req.Remarks
	origin.RoomCode = req.RoomCode

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
// @Summary  [变更] 获取房源筛选选项
// @Description [变更接口] 增加房东房源、返佣和团队房源相关筛选项。
// @Produce  application/json
// @Success  200   {object}  response.Response{data=map[string]map[string]string}  "结果"
// @Router   /center/options [get]
func (h *HouseResourceApi) FilterOptions(c *gin.Context) {

	options, err := ResourceService.FilterOptions()
	if err != nil {
		return
	}

	response.OkWithDetailed(map[string]interface{}{
		"houseType": options,
		"price":     map[string]string{"1": "500以下", "2": "500-1000元", "3": "1000-1500元", "4": "1500-2000元", "5": "2000-2500元", "6": "2500-3000元", "7": "3000元以上"},
		//"houseSource":      map[string]string{"2": "有返佣", "3": "房东房源", "4": "团队房源"},
		"houseSource":      map[string]string{"2": "有返佣", "4": "团队房源"},
		"canViewTeamHouse": userHasTeamPermission(utils.GetUserID(c)),
	}, "获取成功", c)

}

// FilterTypeOptions
// @Tags     Center
// @Summary  [变更] 获取新版房型筛选选项
// @Description [变更接口] 增加房东房源类型，以及协助对接房东、可带看分佣等亮点选项。
// @Produce  application/json
// @Success  200   {object}  response.Response{data=map[string]interface{}}  "结果"
// @Router   /center/type/options [get]
func (h *HouseResourceApi) FilterTypeOptions(c *gin.Context) {

	options, err := ResourceService.FilterOptions1()
	if err != nil {
		return
	}

	re := map[string]interface{}{
		"houseType": options,
		"price":     map[string]string{"1": "500以下", "2": "500-1000元", "3": "1000-1500元", "4": "1500-2000元", "5": "2000-2500元", "6": "2500-3000元", "7": "3000元以上"},
		//"houseSource": map[string]string{"2": "有返佣", "3": "房东房源"},
		"houseSource": map[string]string{"2": "有返佣"},
	}
	if userHasTeamPermission(utils.GetUserID(c)) {
		//re["houseSource"] = map[string]string{"2": "有返佣", "3": "房东房源", "4": "团队房源"}
		re["houseSource"] = map[string]string{"2": "有返佣", "4": "团队房源"}
	}
	response.OkWithDetailed(re, "获取成功", c)

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
	StatisService.InsertRecord(uint(req.ID), "follow", utils.GetUserID(c))

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
	StatisService.InsertRecord(uint(req.ID), "follow", utils.GetUserID(c), -1)
	response.Ok(c)
	return
}

func userHasTeamPermission(userID uint) bool {
	if userID == 0 {
		return false
	}
	user, err := userService.FindUserById(int(userID))
	if err != nil || user == nil {
		return false
	}
	return user.IsFindHouseSupermarket
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
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

// @Tags     Center
// @Summary  房源分享回调接口
// @Produce  application/json
// @Param    data  query    string  true  "id"
// @Success  200   {object}  response.Response{}  "结果 ok"
// @Router   /center/house/shared [get]
func (h *HouseResourceApi) Shared(c *gin.Context) {
	var req request.GetById
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = ResourceService.FollowViewClickSub(uint(req.ID), "shared")
	if err != nil {
		global.GVA_LOG.Error("shared failed!", zap.Error(err))
	}

	StatisService.InsertRecord(uint(req.ID), "shared", utils.GetUserID(c))

	response.Ok(c)
}
