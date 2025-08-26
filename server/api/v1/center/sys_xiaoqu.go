package center

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	test12 "github.com/flipped-aurora/gin-vue-admin/server/utils/test"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"

	"github.com/paulmach/orb"
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

func (receiver *XiaoQuApi) DistanceTree(c *gin.Context) {
	lat, _ := strconv.ParseFloat(c.Query("lat"), 64)
	lng, _ := strconv.ParseFloat(c.Query("lng"), 64)
	qus, err := XiaoQuService.GetDistance(lat, lng)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(qus, "获取成功", c)
}

// Distance
// @Tags     Center
// @Summary  根据坐标获取一公里小区列表
// @accept    application/json
// @Produce   application/json
// @Param   lat    query    int     true  "lat"
// @Param   lng  query    int  true "lng"
// @Success   200   {object}  response.Response{data=[]system.XiaoQu,msg=string}
// @Router    /center/distance [get]
func (receiver *XiaoQuApi) Distance(c *gin.Context) {
	lat, _ := strconv.ParseFloat(c.Query("lat"), 64)
	lng, _ := strconv.ParseFloat(c.Query("lng"), 64)
	// 1. 定义目标点
	point := orb.Point{lat, lng} // orb使用[Lng, Lat]顺序
	//point := orb.Point{116.3974, 39.9093} // 北京
	radius := 1000.0 // 1公里

	results, err := test12.GeoSearch.FindNearbyCommunities(point, radius)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var xqList []system.XiaoQu
	for _, result := range results {
		x := system.XiaoQu{Name: result.Name, Latitude: fmt.Sprintf("%f", result.Point().Lon()), Longitude: fmt.Sprintf("%f", result.Point().Lat())}
		id, _ := strconv.ParseFloat(result.ID, 64)
		x.ID = uint(id)
		xqList = append(xqList, x)
	}
	fmt.Printf("找到 %d 个社区在 %.0f 米范围内\n", len(results), radius)
	response.OkWithDetailed(xqList, "获取成功", c)
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
