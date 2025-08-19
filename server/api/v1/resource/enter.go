package resource

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	HouseResourceApi
	XiaoQuApi
}

var (
	XiaoQuService   = service.ServiceGroupApp.SystemServiceGroup.XiaoQuService
	ResourceService = service.ServiceGroupApp.HouseServiceGroup.ResourceService
	UserService     = service.ServiceGroupApp.SystemServiceGroup.UserService
)
