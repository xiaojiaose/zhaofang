package response

import "github.com/flipped-aurora/gin-vue-admin/server/model/system"

type SharedMapXiaoqu struct {
	XiaoquId  uint   `json:"xiaoquId"`  // 小区ID
	Name      string `json:"name"`      // 小区名称
	Latitude  string `json:"latitude"`  // 坐标纬度
	Longitude string `json:"longitude"` // 坐标经度
	Num       int    `json:"num"`       // 房源数量
}

type SharedMapResponse struct {
	UserInfo system.SysUser    `json:"userInfo"` // 发布人资料
	ExpireAt int64             `json:"expireAt"` // 分享过期时间
	List     []SharedMapXiaoqu `json:"list"`     // 小区聚合点位
}
