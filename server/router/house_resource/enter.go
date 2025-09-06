package house_resource

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type ApiRouter struct{}

func (s *ApiRouter) InitApiRouter(Router *gin.RouterGroup) {
	//houseRouter := Router.Group("center")
	//wxUserApi := v1.ApiGroupApp.ApiCenterGroup.WxUserApi
	//{ // /api/center/verification
	//	//houseRouter.GET("login", wxUserApi.Login) // 创建
	//
	//}

}

func (s *ApiRouter) InitApiAuthRouter(Router *gin.RouterGroup) {
	houseRouter := Router.Group("house")
	houseRecordRouter := Router.Group("house").Use(middleware.OperationRecord())
	resourceApi := v1.ApiGroupApp.House
	statisApi := v1.ApiGroupApp.Statis
	xiaoQuApi := v1.ApiGroupApp.XiaoQu
	fileUploadApi := v1.ApiGroupApp.ExampleApiGroup.FileUploadAndDownloadApi

	{
		houseRouter.POST("xiaoqu", xiaoQuApi.List)
		houseRouter.GET("area", resourceApi.FilterArea)
		houseRouter.GET("options", resourceApi.FilterOptions)
		houseRouter.GET("type/options", resourceApi.FilterTypeOptions)
		houseRecordRouter.POST("create", resourceApi.Create)
		houseRouter.GET("view", resourceApi.View)
		houseRecordRouter.POST("list", resourceApi.List)
		houseRecordRouter.POST("del", resourceApi.DeleteByUserId)
		houseRecordRouter.POST("my", resourceApi.ListByUserId)
		houseRecordRouter.POST("edit", resourceApi.Edit)
		houseRouter.POST("upload", fileUploadApi.UploadFile1)
		houseRecordRouter.POST("approvalState", resourceApi.ApprovalStatus)
		houseRecordRouter.POST("state", resourceApi.States)
	}

	{
		houseRecordRouter.GET("/statis/view", statisApi.View)
		houseRecordRouter.POST("/statis/visit", statisApi.VisitRecord)
		houseRecordRouter.POST("/statis/list", statisApi.VisitHouse)
		houseRecordRouter.POST("/statis/all", statisApi.VisitHouse)
	}

}
