package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/human"
)

// 租户业主信息
type ReqHuman struct {
	human.Human
	//Suites []property.Suite
}
