package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type XiaoQuRouter struct{}

func (s *XiaoQuRouter) InitXiaoQuRouter(Router *gin.RouterGroup) {
	xiaoquRouter := Router.Group("xiaoqu").Use(middleware.OperationRecord())
	xiaoquRouterWithoutRecord := Router.Group("xiaoqu")
	{
		xiaoquRouter.POST("edit", xiaoquApi.Edit)
		xiaoquRouter.POST("dict/upsert", xiaoquApi.UpsertDict)
		xiaoquRouter.POST("dict/delete", xiaoquApi.DeleteDict)
	}
	{
		xiaoquRouterWithoutRecord.GET("show", xiaoquApi.Show)
		xiaoquRouterWithoutRecord.POST("list", xiaoquApi.List)
		xiaoquRouterWithoutRecord.GET("dict/tree", xiaoquApi.GetDictTree)
	}
}
