package center

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
)

type WxUserApi struct {
}

// Index
// @Tags     Center
// @Summary   首页
// @accept    application/json
// @Produce   application/json
// @Success   200  {object}  response.Response{data=string,msg=string}  "首页接口"
// @Router    /center/index [get]
func (wx *WxUserApi) Index(c *gin.Context) {
	//uid := utils.GetUserID(c)
	//mobile := utils.GetMobile(c)
	//
	//var resHeader center.ResCenterHeader
	//
	//human := humanService.FindHumanByMobile(mobile)
	//resHeader.OwnerName = human.Name
	//resHeader.Mobile = human.Mobile
	//
	//landList := landContractService.GetModelByMobile(mobile)
	//for _, contract := range landList {
	//	resHeader.Suites = append(resHeader.Suites, contract.Suite)
	//}

	response.OkWithData("", c)
	return

}

// WxProfile
// @Tags     Center
// @Summary   个人中心接口
// @accept    application/json
// @Produce   application/json
// @Success   200  {object}  response.Response{data=string,msg=string}  "个人中心接口"
// @Router    /center/profile [get]
func (wx *WxUserApi) WxProfile(c *gin.Context) {
	response.OkWithData("", c)
	return
}
