package gaia

import (
	"context"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/gaia"
	"github.com/gin-gonic/gin"
)

type SystemApi struct{}

// GetDingTalk 获取钉钉系统配置
// @Tags System
// @Summary 获取钉钉系统配置
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query gaia.Tenants true "用id查询tenants表"
// @Success 200 {object} response.Response{data=gaia.Tenants,msg=string} "查询成功"
// @Router /gaia/system/dingtalk [get]
func (systemApi *SystemApi) GetDingTalk(c *gin.Context) {
	var host string
	var config = make(map[string]interface{})
	if host, _ = global.GVA_Dify_REDIS.Get(context.Background(), "api_host").Result(); len(host) == 0 {
		host = global.GVA_CONFIG.Gaia.Url
	}
	config["host"] = host
	config["config"] = systemIntegratedService.GetIntegratedConfig(gaia.SystemIntegrationDingTalk)
	response.OkWithData(config, c)
}

// SetDingTalk 设置钉钉系统配置
// @Tags System
// @Summary 设置钉钉系统配置
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query gaia.Tenants true "用id查询tenants表"
// @Success 200 {object} response.Response{data=gaia.Tenants,msg=string} "查询成功"
// @Router /gaia/system/dingtalk [post]
func (systemApi *SystemApi) SetDingTalk(c *gin.Context) {
	var err error
	var req gaia.SystemIntegration
	if err = c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// update
	req.Classify = gaia.SystemIntegrationDingTalk
	if err = systemIntegratedService.SetIntegratedConfig(
		req, "", req.Test); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData("ok", c)
}
