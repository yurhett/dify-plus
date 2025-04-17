import service from '@/utils/request'

// @Tags systrm
// @Summary 获取钉钉集成配置
// @Security ApiKeyAuth
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"返回成功"}"
// @Router /gaia/system/dingtalk [get]
export const getSystemDingTalk = () => {
    return service({
        url: '/gaia/system/dingtalk',
        method: 'get'
    })
}

// @Tags systrm
// @Summary 修改钉钉集成配置
// @Security ApiKeyAuth
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"返回成功"}"
// @Router /gaia/system/dingtalk [post]
export const setSystemDingTalk = (data) => {
    return service({
        url: '/gaia/system/dingtalk',
        method: 'post',
        data,
    })
}

// @Tags systrm
// @Summary 获取OAuth2集成配置
// @Security ApiKeyAuth
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"返回成功"}"
// @Router /gaia/system/oauth2 [get]
export const getSystemOAuth2 = () => {
    return service({
        url: '/gaia/system/oauth2',
        method: 'get'
    })
}

// @Tags systrm
// @Summary 修改OAuth2集成配置
// @Security ApiKeyAuth
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"返回成功"}"
// @Router /gaia/system/oauth2 [post]
export const setSystemOAuth2 = (data) => {
    return service({
        url: '/gaia/system/oauth2',
        method: 'post',
        data,
    })
}
