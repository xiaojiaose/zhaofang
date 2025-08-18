package initialize

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/sms"
	"github.com/gin-gonic/gin"
)

func InstallPlugin(PrivateGroup *gin.RouterGroup, PublicRouter *gin.RouterGroup, engine *gin.Engine) {
	if global.GVA_DB == nil {
		global.GVA_LOG.Info("项目暂未初始化，无法安装插件，初始化后重启项目即可完成插件安装")
		return
	}

	fmt.Println("短信==》", PublicRouter)
	PluginInit(PublicRouter, sms.CreateAliSmsPlug(global.GVA_CONFIG.AliyunSms.AccessKeyId, global.GVA_CONFIG.AliyunSms.AccessSecret, global.GVA_CONFIG.AliyunSms.SignName))

	bizPluginV1(PrivateGroup, PublicRouter)
	bizPluginV2(engine)
}
