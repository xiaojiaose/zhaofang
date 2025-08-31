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
