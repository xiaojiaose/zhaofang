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
	WxNo       string `json:"wxNo"` // 微信号
	WxNickName string `json:"wxNickName"`
	HeaderImg  string `json:"headerImg"`
	Phone      string `json:"phone"`
}

type MyResourceResponse struct {
	house.Resource
	WxNo               string `json:"wxNo"`
	WxNickName         string `json:"wxNickName"`
	HeaderImg          string `json:"headerImg"`
	PublishQuotaTotal  int    `json:"publishQuotaTotal"`
	PublishQuotaUsed   int    `json:"publishQuotaUsed"`
	PublishQuotaRemain int    `json:"publishQuotaRemain"`
}

type DictBuildingResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
