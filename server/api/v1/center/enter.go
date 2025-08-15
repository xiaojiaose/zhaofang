package center

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiCenter struct {
	WxUserApi
	XiaoQuApi
	HouseResourceApi
}

var (
	XiaoQuService   = service.ServiceGroupApp.SystemServiceGroup.XiaoQuService
	ResourceService = service.ServiceGroupApp.HouseServiceGroup.ResourceService
)
