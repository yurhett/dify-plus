package gaia

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/faabiosr/cachego/file"
	"github.com/fastwego/dingding"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/gaia"
	"github.com/flipped-aurora/gin-vue-admin/server/model/gaia/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type SystemIntegratedService struct{}

// GetIntegratedConfig
// @Tags System Integrated
// @Summary 获取系统集成配置
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
func (e *SystemIntegratedService) GetIntegratedConfig(classID uint) (integrate gaia.SystemIntegration) {
	// classID是否在
	var err error
	if err = global.GVA_DB.Where("classify = ?", classID).First(&integrate).Error; err != nil {
		integrate = gaia.SystemIntegration{
			Classify: classID,
			Status:   false,
		}
		// 创建相关集成信息
		global.GVA_DB.Create(&integrate)
	}
	// 隐藏部分加密信息
	var secret string
	if secret, err = utils.DecryptBlowfish(integrate.AppSecret, global.GVA_CONFIG.JWT.SigningKey); err == nil {
		integrate.AppSecret = utils.AddAsteriskToString(secret)
	}
	integrate.CorpID = utils.AddAsteriskToString(integrate.CorpID)
	return integrate
}

// SetIntegratedConfig
// @Tags System Integrated
// @Summary 设置系统集成配置
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @param: integrate gaia.SystemIntegration, code string, test bool
// @return: err error
func (e *SystemIntegratedService) SetIntegratedConfig(
	integrate gaia.SystemIntegration, code string, test bool) (err error) {
	// classID是否在
	var log gaia.SystemIntegration
	if err = global.GVA_DB.Where("classify = ?", integrate.Classify).First(&log).Error; err != nil {
		return err
	}
	// AppSecret
	var secret string
	if secret, err = utils.DecryptBlowfish(log.AppSecret, global.GVA_CONFIG.JWT.SigningKey); err == nil {
		encodeSecret := utils.AddAsteriskToString(secret)
		if encodeSecret != integrate.AppSecret {
			if secret, err = utils.EncryptBlowfish(
				[]byte(integrate.AppSecret), global.GVA_CONFIG.JWT.SigningKey); err != nil {
				return errors.New("AppSecret加密失败")
			}
			// save
			log.AppSecret = secret
		} else {
			// 为什么不用 integrate.AppSecret, 被加*了
			if secret, err = utils.DecryptBlowfish(log.AppSecret, global.GVA_CONFIG.JWT.SigningKey); err != nil {
				return errors.New("AppSecret解析失败")
			}
			integrate.AppSecret = secret
		}
	}
	// CorpID
	if utils.AddAsteriskToString(log.CorpID) != integrate.CorpID {
		log.CorpID = integrate.CorpID
	}
	// AppID 不加密，直接赋值
	log.AppID = integrate.AppID
	// 关闭不需要请求
	if integrate.Status || test {
		// 测试连接
		if err = e.TestConnection(integrate, code); err != nil {
			return errors.New("连接失败:" + err.Error())
		}
	}
	// Test completed
	if test {
		return err
	}
	// save
	if err = global.GVA_DB.Model(&gaia.SystemIntegration{}).Where(
		"id=?", log.Id).Updates(&map[string]interface{}{
		"config":     integrate.Config,
		"status":     integrate.Status,
		"agent_id":   integrate.AgentID,
		"app_key":    integrate.AppKey,
		"app_secret": log.AppSecret,
		"corp_id":    log.CorpID,
		"app_id":     log.AppID,
	}).Error; err != nil {
		return err
	}
	return nil
}

// DingTalkConfigAvailable
// @Tags System Integrated
// @Summary 测试钉钉配置是否可用
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @param: req gaia.SystemIntegration
// @return: *dingding.Client, error
func (e *SystemIntegratedService) DingTalkConfigAvailable(req gaia.SystemIntegration) (*dingding.Client, error) {
	var err error
	var reqs *http.Request
	dingding.ServerUrl = "https://api.dingtalk.com"
	// 特殊需要，检查可用性就不设置缓存了
	return dingding.NewClient(&dingding.DefaultAccessTokenManager{
		Id:    uuid.New().String(),
		Cache: file.New(os.TempDir()),
		Name:  "x-acs-dingtalk-access-token",
		GetRefreshRequestFunc: func() *http.Request {
			params := url.Values{}
			params.Add("appkey", req.AppKey)
			params.Add("appsecret", req.AppSecret)
			reqs, err = http.NewRequest(http.MethodGet, "https://oapi.dingtalk.com/gettoken?"+params.Encode(), nil)
			return reqs
		},
	}), err
}

// TestConnection 测试连接
// @Tags System Integrated
// @Summary 测试系统集成连接
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @param: integrate gaia.SystemIntegration, code string
// @return: error
func (e *SystemIntegratedService) TestConnection(integrate gaia.SystemIntegration, code string) error {
	switch integrate.Classify {
	case gaia.SystemIntegrationDingTalk:
		// 测试钉钉连接
		if _, err := e.DingTalkConfigAvailable(integrate); err != nil {
			return errors.New("钉钉链接失败: " + err.Error())
		}
		return nil
	case gaia.SystemIntegrationOAuth2:
		// 测试OAuth2连接
		return e.TestOAuth2Connection(integrate, code)
	default:
		return errors.New("不支持的集成类型")
	}
}

// TestOAuth2Connection 测试OAuth2连接
// @Tags System Integrated
// @Summary 测试OAuth2连接
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @param: integrate gaia.SystemIntegration, code string
// @return: error
func (e *SystemIntegratedService) TestOAuth2Connection(integrate gaia.SystemIntegration, code string) (err error) {
	// 解析Config字段
	var configMap request.SystemOAuth2Request
	if err = json.Unmarshal([]byte(integrate.Config), &configMap); err != nil {
		global.GVA_LOG.Error("解析OAuth2配置失败!", zap.Error(err))
		return err
	}
	// 没有code的（保存操作）
	if len(code) == 0 {
		return nil
	}
	// 检查必要字段
	if configMap.ServerURL == "" || configMap.TokenURL == "" || integrate.AppID == "" || integrate.AppSecret == "" {
		return errors.New("请填写完整的 OAuth2 配置信息")
	}

	// 验证URL格式
	if !strings.HasPrefix(configMap.ServerURL, "http://") && !strings.HasPrefix(configMap.ServerURL, "https://") {
		return errors.New("ServerURL 必须以 http:// 或 https:// 开头")
	}

	// 合成请求byte
	formData := url.Values{}
	formData.Set("grant_type", "authorization_code")
	formData.Set("code", code)
	// redirect_uri 必须与授权时一致
	redirectUri := strings.TrimSpace(configMap.RedirectUri)
	if redirectUri == "" {
		// 当配置中没有指定redirect_uri时，使用与前端相同的默认值
		var host string
		if host, _ = global.GVA_Dify_REDIS.Get(context.Background(), "api_host").Result(); len(host) == 0 {
			host = global.GVA_CONFIG.Gaia.Url
		}
		redirectUri = fmt.Sprintf("%s/admin/api/base/auth2/callback", host)
	}
	formData.Set("redirect_uri", redirectUri)
	// 支持basic与post两种认证方式
	// 默认使用client_secret_post，除非配置为basic
	tokenAuthMethod := strings.ToLower(strings.TrimSpace(configMap.TokenAuthMethod))
	useBasic := tokenAuthMethod == "client_secret_basic"
	if !useBasic {
		formData.Set("client_secret", integrate.AppSecret)
		formData.Set("client_id", integrate.AppID)
	}

	// 构建token URL，确保正确的URL格式
	baseURL := strings.TrimSuffix(configMap.ServerURL, "/")
	tokenPath := strings.TrimPrefix(configMap.TokenURL, "/")
	tokenURL := fmt.Sprintf("%s/%s", baseURL, tokenPath)

	// 记录构建的URL用于调试
	global.GVA_LOG.Info("构建OAuth2 token URL", 
		zap.String("serverURL", configMap.ServerURL),
		zap.String("tokenURL", configMap.TokenURL), 
		zap.String("constructedURL", tokenURL))

	// 发送请求
	var req *http.Request
	client := &http.Client{}
	req, err = http.NewRequest("POST", tokenURL, strings.NewReader(formData.Encode()))
	if err != nil {
		global.GVA_LOG.Error("创建测试请求失败", zap.Error(err), zap.String("url", tokenURL))
		return errors.New(fmt.Sprintf("创建测试请求失败: %s", err.Error()))
	}

	// 设置认证与Content-Type
	if useBasic {
		req.SetBasicAuth(integrate.AppID, integrate.AppSecret)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		global.GVA_LOG.Error("测试 OAuth2 连接失败", zap.Error(err), zap.String("url", tokenURL))
		return errors.New(fmt.Sprintf("连接 OAuth2 服务器失败: %s", err.Error()))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		global.GVA_LOG.Error("测试 OAuth2 连接失败", 
			zap.Int("status", resp.StatusCode), 
			zap.String("url", tokenURL),
			zap.String("method", "POST"))
		return errors.New(fmt.Sprintf("OAuth2 服务器返回错误状态码: %d (URL: %s)", resp.StatusCode, tokenURL))
	}
	var bodyByte []byte
	if bodyByte, err = io.ReadAll(resp.Body); err != nil {
		return fmt.Errorf("OAuth2 request io.ReadAll: %s", resp.Status)
	}

	var tokenMap request.SystemOAuth2Error
	if err = json.Unmarshal(bodyByte, &tokenMap); err == nil && tokenMap.Code != 0 {
		return fmt.Errorf("OAuth2 Eroor: %s", tokenMap.Info)
	}

	return nil
}
