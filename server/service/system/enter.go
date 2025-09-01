package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/sms/service"
	"github.com/flipped-aurora/gin-vue-admin/server/service/house"
)

type ServiceGroup struct {
	JwtService
	ApiService
	MenuService
	UserService
	CasbinService
	InitDBService
	AutoCodeService
	BaseMenuService
	AuthorityService
	DictionaryService
	SystemConfigService
	OperationRecordService
	DictionaryDetailService
	AuthorityBtnService
	SysExportTemplateService
	SysParamsService
	SysVersionService
	XiaoQuService
	house.DictService
	AutoCodePlugin   autoCodePlugin
	AutoCodePackage  autoCodePackage
	AutoCodeHistory  autoCodeHistory
	AutoCodeTemplate autoCodeTemplate
	service.AliSmsService
}
