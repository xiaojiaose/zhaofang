package center

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
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
// @Success   200  {object}  response.Response{data=map[string]interface{},msg=string}  "个人中心接口"
// @Router    /center/profile [get]
func (wx *WxUserApi) WxProfile(c *gin.Context) {
	user, err := userService.FindUserById(int(utils.GetUserID(c)))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(user, c)
	return
}

// SetWxProfile
// @Tags     Center
// @Summary  小程序设置微信资料
// @accept    application/json
// @Produce   application/json
// @Param    data  body      systemReq.WxProfileSync  true  "微信昵称、头像、微信号"
// @Success   200  {object}  response.Response{data=string,msg=string}  "设置微信资料"
// @Router    /center/profile [post]
func (wx *WxUserApi) SetWxProfile(c *gin.Context) {
	var req systemReq.WxProfileSync
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := userService.SetWxProfile(utils.GetUserID(c), req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	user, err := userService.FindUserById(int(utils.GetUserID(c)))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(user, c)
}
