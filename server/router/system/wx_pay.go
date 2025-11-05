package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/gin-gonic/gin"
)

type WxPayRouter struct{}

func (w *WxPayRouter) InitWxPayRouter(Router *gin.RouterGroup) {
	wxPayRouter := Router.Group("wxpay")
	wxPayApi := v1.ApiGroupApp.SystemApiGroup.WxPayApi
	{
		wxPayRouter.POST("createOrder", wxPayApi.CreatePayOrder)   // 创建支付订单
		wxPayRouter.POST("notify", wxPayApi.PayNotify)            // 支付回调
		wxPayRouter.GET("queryOrder", wxPayApi.QueryOrder)        // 查询本地订单状态
		wxPayRouter.GET("queryWxOrder", wxPayApi.QueryWxOrder)    // 查询微信订单状态
	}
}
