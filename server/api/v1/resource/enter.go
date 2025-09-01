package resource

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	HouseResourceApi
	XiaoQuApi
	StatisDataApi
}

var (
	XiaoQuService   = service.ServiceGroupApp.SystemServiceGroup.XiaoQuService
	DictService     = service.ServiceGroupApp.SystemServiceGroup.DictService
	ResourceService = service.ServiceGroupApp.HouseServiceGroup.ResourceService
	UserService     = service.ServiceGroupApp.SystemServiceGroup.UserService
	StatisService   = service.ServiceGroupApp.HouseServiceGroup.StatisService
)
