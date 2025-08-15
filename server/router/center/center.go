package center

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/gin-gonic/gin"
)

type CenterRouter struct{}

func (s *CenterRouter) InitCenterRouter(Router *gin.RouterGroup) {
	//houseRouter := Router.Group("center")
	//wxUserApi := v1.ApiGroupApp.ApiCenterGroup.WxUserApi
	//{ // /api/center/verification
	//	//houseRouter.GET("login", wxUserApi.Login) // 创建
	//
	//}

}

func (s *CenterRouter) InitCenterAuthRouter(Router *gin.RouterGroup) {
	houseRouter := Router.Group("center")
	centerApi := v1.ApiGroupApp.ApiCenterGroup.WxUserApi
	resourceApi := v1.ApiGroupApp.ApiCenterGroup.HouseResourceApi
	xiaoQuApi := v1.ApiGroupApp.ApiCenterGroup.XiaoQuApi
	fileUploadApi := v1.ApiGroupApp.ExampleApiGroup.FileUploadAndDownloadApi

	{
		houseRouter.GET("index", centerApi.WxProfile)   //
		houseRouter.POST("xiaoqu", xiaoQuApi.List)      //
		houseRouter.GET("distance", xiaoQuApi.Distance) //
		houseRouter.GET("test", resourceApi.Test)       //
		houseRouter.GET("area", resourceApi.FilterArea)
		houseRouter.GET("options", resourceApi.FilterOptions)
		houseRouter.POST("house/create", resourceApi.Create)
		houseRouter.GET("house/view", resourceApi.View)
		houseRouter.POST("house/edit", resourceApi.Edit)
		houseRouter.POST("upload", fileUploadApi.UploadFile1)

		//houseRouter.GET("incomeInfo", wxUserApi.GetIncomeInfo) //
	}

}
