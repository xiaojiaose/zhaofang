package center

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/house"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
)

type HouseResourceApi struct {
}

// View
// @Tags     Center
// @Summary  创建|编辑 房源
// @Produce  application/json
// @Param    data  query    string  true  "id"
// @Success  200   {object}  response.Response{data=house.Resource}  "结果"
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

	response.OkWithDetailed(info, "获取成功", c)
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

	claims := utils.GetUserInfo(c) // 获取登陆用户
	req.Owner = claims.BaseClaims.ID

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
// @Summary  获取区域列表（丰台区、朝阳区）
// @Produce  application/json
// @Param    data  query    string  true  "cityId"
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
				Name: district.Name,
				Id:   district.ID,
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
