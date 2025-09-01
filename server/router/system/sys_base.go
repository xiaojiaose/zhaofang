package system

import (
	"github.com/gin-gonic/gin"
)

type BaseRouter struct{}

func (s *BaseRouter) InitBaseRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	baseRouter := Router.Group("base")
	{
		baseRouter.POST("login", baseApi.Login)
		baseRouter.POST("captcha", baseApi.Captcha)
		baseRouter.POST("mobileLogin", baseApi.MobileLogin)
		baseRouter.POST("sendSms", baseApi.SendSms)
		baseRouter.GET("building", xiaoquApi.GetBuilding)
		baseRouter.GET("unit", xiaoquApi.GetUnit)
		baseRouter.GET("house", xiaoquApi.GetHouse)
	}
	return baseRouter
}

func (s *BaseRouter) InitWxBaseRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	baseRouter := Router.Group("wx")
	{
		baseRouter.POST("login", wxApi.WxLogin)
		baseRouter.POST("getMobile", wxApi.GetWxMobile)
	}

	return baseRouter
}
