package manage

import (
	"github.com/gin-gonic/gin"
	v1 "main.go/api/v1"
)

type ManageIndexConfigRouter struct {
}

func (r *ManageIndexConfigRouter) InitManageIndexConfigRouter(Router *gin.RouterGroup) {
	mallIndexConfigRouter := Router.Group("v1")
	var mallIndexConfigApi = v1.ApiGroupApp.ManageApiGroup.ManageIndexConfigApi
	{
		mallIndexConfigRouter.POST("indexConfigs", mallIndexConfigApi.CreateIndexConfig)        // 创建首页配置项
		mallIndexConfigRouter.POST("indexConfigs/delete", mallIndexConfigApi.DeleteIndexConfig) // 删除首页配置项
		mallIndexConfigRouter.PUT("indexConfigs", mallIndexConfigApi.UpdateIndexConfig)         // 更新首页配置项
		mallIndexConfigRouter.GET("indexConfigs/:id", mallIndexConfigApi.FindIndexConfig)       // 根据ID获取首页配置项
		mallIndexConfigRouter.GET("indexConfigs", mallIndexConfigApi.GetIndexConfigList)        // 获取首页配置项列表
	}
}
