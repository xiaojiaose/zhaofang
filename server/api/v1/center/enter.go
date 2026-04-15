package center

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiCenter struct {
	WxUserApi
	XiaoQuApi
	HouseResourceApi
}

var (
	XiaoQuService       = service.ServiceGroupApp.SystemServiceGroup.XiaoQuService
	ResourceService     = service.ServiceGroupApp.HouseServiceGroup.ResourceService
	FavoriteService     = service.ServiceGroupApp.HouseServiceGroup.FavoriteService
	RewardService       = service.ServiceGroupApp.HouseServiceGroup.RewardService
	ShareService        = service.ServiceGroupApp.HouseServiceGroup.ShareService
	ContactQuotaService = service.ServiceGroupApp.HouseServiceGroup.ContactQuotaService
	userService         = service.ServiceGroupApp.SystemServiceGroup.UserService
	StatisService       = service.ServiceGroupApp.HouseServiceGroup.StatisService
)
