package main

import (
	"github.com/flipped-aurora/gin-vue-admin/server/core"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/initialize"
	_ "go.uber.org/automaxprocs"
	"go.uber.org/zap"
	"google.golang.org/protobuf/runtime/protoimpl"
)

//go:generate go env -w GO111MODULE=on
//go:generate go env -w GOPROXY=https://goproxy.cn,direct
//go:generate go mod tidy
//go:generate go mod download

// 这部分 @Tag 设置用于排序, 需要排序的接口请按照下面的格式添加
// swag init 对 @Tag 只会从入口文件解析, 默认 main.go
// 也可通过 --generalInfo flag 指定其他文件
// @Tag.Name        Base
// @Tag.Name        SysUser
// @Tag.Description 用户

// @title                       Gin-Vue-Admin Swagger API接口文档
// @version                     v2.8.4
// @description                 使用gin+vue进行极速开发的全栈开发基础平台
// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        x-token
// @BasePath                    /

type BatteryData struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Action        string                 `protobuf:"bytes,1,opt,name=action,proto3" json:"action,omitempty"` // 采集动作 "POWER_IN", "POWER_OUT", "START_UP", "SHUT_DOWN", "INITIATIVE_UPLOAD"
	Info          *Info                  `protobuf:"bytes,2,opt,name=info,proto3" json:"info,omitempty"`
	Stamp         int64                  `protobuf:"varint,3,opt,name=stamp,proto3" json:"stamp,omitempty"` // 采集时间
	Type          int32                  `protobuf:"varint,4,opt,name=type,proto3" json:"type,omitempty"`   // 采集种类 1 主动，2 被动
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

type Info struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Cap           float64                `protobuf:"fixed64,1,opt,name=cap,proto3" json:"cap,omitempty"`   // 总容量
	Ccr           int32                  `protobuf:"varint,2,opt,name=ccr,proto3" json:"ccr,omitempty"`    // 充电状态
	Bhs           string                 `protobuf:"bytes,3,opt,name=bhs,proto3" json:"bhs,omitempty"`     // 电池健康程度 未知 良好 过热 没电 过电压 温度过低 未知错误
	Soc           int32                  `protobuf:"varint,4,opt,name=soc,proto3" json:"soc,omitempty"`    // 当前电量
	Chg           string                 `protobuf:"bytes,5,opt,name=chg,proto3" json:"chg,omitempty"`     // 充电类型 充电器 USB 无线充电 未充电
	St            string                 `protobuf:"bytes,6,opt,name=st,proto3" json:"st,omitempty"`       // 电池状态 未知 充电中 放电中 未充电 电池满
	Tech          string                 `protobuf:"bytes,7,opt,name=tech,proto3" json:"tech,omitempty"`   // 电池技术
	Temp          float64                `protobuf:"fixed64,8,opt,name=temp,proto3" json:"temp,omitempty"` // 温度
	V             int32                  `protobuf:"varint,9,opt,name=v,proto3" json:"v,omitempty"`        // 电压
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func main() {
	// 初始化系统
	initializeSystem()
	// 运行服务器
	core.RunServer()
}

// initializeSystem 初始化系统所有组件
// 提取为单独函数以便于系统重载时调用
func initializeSystem() {
	global.GVA_VP = core.Viper() // 初始化Viper
	initialize.OtherInit()
	global.GVA_LOG = core.Zap() // 初始化zap日志库
	zap.ReplaceGlobals(global.GVA_LOG)
	global.GVA_DB = initialize.Gorm() // gorm连接数据库
	initialize.Timer()
	initialize.DBList()
	initialize.SetupHandlers() // 注册全局函数
	if global.GVA_DB != nil {
		initialize.RegisterTables() // 初始化表
	}
	global.InitZincSearch(global.GVA_CONFIG.ZincSearch.Url, global.GVA_CONFIG.ZincSearch.Username, global.GVA_CONFIG.ZincSearch.Password)
	initialize.GeoXiaoqu()
}
