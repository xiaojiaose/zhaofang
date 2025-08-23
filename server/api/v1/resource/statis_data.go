package resource

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/search"
	"github.com/gin-gonic/gin"
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
