package system

import "github.com/flipped-aurora/gin-vue-admin/server/plugin/sms/service"

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
	AutoCodePlugin   autoCodePlugin
	AutoCodePackage  autoCodePackage
	AutoCodeHistory  autoCodeHistory
	AutoCodeTemplate autoCodeTemplate
	service.AliSmsService
}
