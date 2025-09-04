package center

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
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

func (s *CenterRouter) InitCenterAuthRouter(Router *gin.RouterGroup, publicRouter *gin.RouterGroup) {
	houseRouterRecord := Router.Group("center").Use(middleware.OperationRecord())
	houseRouter := Router.Group("center")
	public := publicRouter.Group("center")
	centerApi := v1.ApiGroupApp.ApiCenterGroup.WxUserApi
	resourceApi := v1.ApiGroupApp.ApiCenterGroup.HouseResourceApi
	xiaoQuApi := v1.ApiGroupApp.ApiCenterGroup.XiaoQuApi
	fileUploadApi := v1.ApiGroupApp.ExampleApiGroup.FileUploadAndDownloadApi

	{
		houseRouterRecord.GET("index", centerApi.WxProfile)   //
		houseRouterRecord.POST("xiaoqu/list", xiaoQuApi.List) //
		houseRouterRecord.GET("xiaoqu/show", xiaoQuApi.Show)
		houseRouter.GET("distance", xiaoQuApi.Distance)         //
		houseRouter.GET("distanceTree", xiaoQuApi.DistanceTree) //

		houseRouter.GET("test", resourceApi.Test) //
		houseRouter.GET("area", resourceApi.FilterArea)
		houseRouter.GET("options", resourceApi.FilterOptions)
		houseRouterRecord.POST("house/create", resourceApi.Create)
		houseRouterRecord.POST("/house/del", resourceApi.DeleteByUserId)
		public.GET("/house/view", resourceApi.View)
		houseRouterRecord.GET("house/mobile", resourceApi.GetMobile)
		houseRouter.POST("house/xiaoquAgg", resourceApi.ListByXiaoquAgg)
		houseRouter.POST("house/xiaoquAggList", resourceApi.ListByXiaoquAggList)
		houseRouterRecord.POST("house/listByXiaoqu", resourceApi.ListByXiaoquId)
		houseRouterRecord.POST("house/my", resourceApi.ListByUserId)
		houseRouterRecord.POST("house/edit", resourceApi.Edit)
		houseRouterRecord.POST("upload", fileUploadApi.UploadFile1)
		houseRouterRecord.GET("favorite/add", resourceApi.FavoriteAdd)
		houseRouterRecord.GET("favorite/del", resourceApi.FavoriteDel)
		houseRouterRecord.POST("favorite/list", resourceApi.FavoriteList)
		houseRouterRecord.POST("state", resourceApi.States)

		//houseRouter.GET("incomeInfo", wxUserApi.GetIncomeInfo) //
	}

}
