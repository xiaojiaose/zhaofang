package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

type Area struct {
	global.GVA_MODEL
	Name string `json:"name"` // name
	City string `json:"city"` // city
	Sort int    `json:"sort"` // sort
	Able bool   `json:"able"` // 是否启用
}

func (s *Area) TableName() string {
	return "areas"
}
