package request

// SystemOAuth2Error OAuth2 错误返回
type SystemOAuth2Error struct {
	Code int    `json:"code" gorm:"comment:分类"`   // 错误代码
	Info string `json:"info" gorm:"comment:错误详情"` // 错误详情
}

// SystemOAuth2Request OAuth2 集成配置
type SystemOAuth2Request struct {
	Classify       uint   `json:"classify" gorm:"comment:分类"`              // 分类
	Status         bool   `json:"status" gorm:"comment:状态"`                // 状态
	ServerURL      string `json:"server_url" gorm:"comment:服务器地址"`         // OAuth2 服务器地址
	AuthorizeURL   string `json:"authorize_url" gorm:"comment:申请认证的URL"`   // 申请认证的URL
	TokenURL       string `json:"token_url" gorm:"comment:获取Token的URL"`    // 获取Token的URL
	UserinfoURL    string `json:"userinfo_url" gorm:"comment:获取用户信息URL"`   // 获取用户信息的URL
	LogoutURL      string `json:"logout_url" gorm:"comment:退出登录回调URL"`     // 退出登录回调URL
	AppID          string `json:"app_id" gorm:"comment:Client ID"`         // Client ID
	AppSecret      string `json:"app_secret" gorm:"comment:Client Secret"` // Client Secret
	UserNameField  string `json:"user_name_field" gorm:"comment:用户名字段"`    // 用户名字段
	UserEmailField string `json:"user_email_field" gorm:"comment:邮箱字段"`    // 邮箱字段
	UserIDField    string `json:"user_id_field" gorm:"comment:用户唯一标识字段"`   // 用户唯一标识字段
	Test           bool   `json:"test" gorm:"default:0;comment:是否测试链接联通性"` // 是否测试链接联通性
	Code           string `json:"code" gorm:"default:0;comment:code代码"`    // code代码
}
