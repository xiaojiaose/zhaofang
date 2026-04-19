package response

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/house"
)

type ResourceResponse struct {
	house.Resource
	Follow    bool   `json:"follow"`    // 关注是否
	Latitude  string `json:"latitude"`  // latitude
	Longitude string `json:"longitude"` // longitude
}

type ResourceVisitResponse struct {
	house.Resource
	WxNo       string `json:"wxNo"`       // 微信号
	WxNickName string `json:"wxNickName"` // 微信昵称
	HeaderImg  string `json:"headerImg"`  // 头像
	Phone      string `json:"phone"`      // 手机号
}

type MyResourceResponse struct {
	house.Resource
	WxNo               string `json:"wxNo"`               // 微信号
	WxNickName         string `json:"wxNickName"`         // 微信昵称
	HeaderImg          string `json:"headerImg"`          // 头像
	PublishQuotaTotal  int    `json:"publishQuotaTotal"`  // 可上架总数
	PublishQuotaUsed   int    `json:"publishQuotaUsed"`   // 已上架数量
	PublishQuotaRemain int    `json:"publishQuotaRemain"` // 剩余可上架数量
}

type DictBuildingResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
