package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
)

type CommonAPi struct{}

var CITYS = []response.ReturnValue{
	{
		Name:  "杭州",
		Value: "1",
	},
	{
		Name:  "绍兴",
		Value: "2",
	},
}

var DESIGNE = []response.ReturnValue{
	{
		Name:  "侯文静",
		Value: "1",
	},
	{
		Name:  "赵峰",
		Value: "2",
	},
}

// CityList
// @Tags     common
// @Summary  获取城市列表
// @Produce  application/json
// @Success  200   {object}  response.Response{data=response.ReturnList{list=[]response.ReturnValue}}  "结果"
// @Router   /api/common/cities [get]
func (receiver *CommonAPi) CityList(c *gin.Context) {
	response.OkWithData(response.ReturnList{List: CITYS}, c)
	return
}

// DesignerList
// @Tags     common
// @Summary  获取设计师列表
// @Produce  application/json
// @Success  200   {object}  response.Response{data=response.ReturnList{list=[]response.ReturnValue}}  "结果"
// @Router   /api/common/designers [get]
func (receiver *CommonAPi) DesignerList(c *gin.Context) {

	response.OkWithData(response.ReturnList{List: DESIGNE}, c)
	return
}
