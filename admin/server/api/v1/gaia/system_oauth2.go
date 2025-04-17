package gaia

import (
	"context"
	"encoding/json"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/gaia"
	"github.com/flipped-aurora/gin-vue-admin/server/model/gaia/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SystemOAuth2Api struct{}

// GetOAuth2Config
// @Tags System
// @Summary 获取OAuth2集成配置
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /gaia/system/oauth2 [get]
func (s *SystemOAuth2Api) GetOAuth2Config(c *gin.Context) {
	var configMap request.SystemOAuth2Request
	var config = make(map[string]interface{})
	// 直接使用service层获取request.SystemOAuth2Request结构
	integrated := systemIntegratedService.GetIntegratedConfig(gaia.SystemIntegrationOAuth2)
	_ = json.Unmarshal([]byte(integrated.Config), &configMap)
	// setting
	configMap.AppID = integrated.AppID
	configMap.Status = integrated.Status
	configMap.Classify = integrated.Classify
	configMap.AppSecret = integrated.AppSecret
	var host string
	if host, _ = global.GVA_Dify_REDIS.Get(context.Background(), "api_host").Result(); len(host) == 0 {
		host = global.GVA_CONFIG.Gaia.Url
	}
	config["host"] = host
	config["config"] = configMap
	response.OkWithData(config, c)
}

// SetOAuth2Config
// @Tags System
// @Summary 修改OAuth2集成配置
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.SystemOAuth2Request true "修改数据"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"设置成功"}"
// @Router /gaia/system/oauth2 [post]
func (s *SystemOAuth2Api) SetOAuth2Config(c *gin.Context) {
	var req request.SystemOAuth2Request
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 序列化为JSON
	configBytes, err := json.Marshal(&map[string]string{
		"server_url":       req.ServerURL,
		"authorize_url":    req.AuthorizeURL,
		"token_url":        req.TokenURL,
		"userinfo_url":     req.UserinfoURL,
		"logout_url":       req.LogoutURL,
		"user_name_field":  req.UserNameField,
		"user_email_field": req.UserEmailField,
		"user_id_field":    req.UserIDField,
	})
	if err != nil {
		global.GVA_LOG.Error("序列化OAuth2配置失败!", zap.Error(err))
		response.FailWithMessage("配置序列化失败", c)
		return
	}

	// 更新配置
	if err = systemIntegratedService.SetIntegratedConfig(gaia.SystemIntegration{
		Classify:  gaia.SystemIntegrationOAuth2,
		Config:    string(configBytes),
		AppSecret: req.AppSecret,
		Status:    req.Status,
		AppID:     req.AppID,
	}, req.Code, req.Test); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData("设置成功", c)
}
