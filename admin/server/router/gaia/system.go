package gaia

import (
	"github.com/gin-gonic/gin"
)

type SystemRouter struct{}

// InitSystemRouter 初始化系统路由
func (s *SystemRouter) InitSystemRouter(Router *gin.RouterGroup) {
	systemRouter := Router.Group("gaia/system")
	{
		systemRouter.GET("dingtalk", systemApi.GetDingTalk)          // 获取钉钉系统配置
		systemRouter.POST("dingtalk", systemApi.SetDingTalk)         // 设置钉钉系统配置
		systemRouter.GET("oauth2", systemOAuth2Api.GetOAuth2Config)  // 获取OAuth2配置
		systemRouter.POST("oauth2", systemOAuth2Api.SetOAuth2Config) // 设置OAuth2配置
	}
}
