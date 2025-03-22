package models

import (
	"time"
)

// JIRASettings JIRA集成设置
type JIRASettings struct {
	Enabled        bool   `json:"enabled"`
	URL            string `json:"url"`
	APIToken       string `json:"apiToken"`
	Username       string `json:"username"`
	DefaultProject string `json:"defaultProject"`
}

// WechatSettings 微信扫码登录设置
type WechatSettings struct {
	Enabled        bool   `json:"enabled"`
	AppID          string `json:"appId"`
	AppSecret      string `json:"appSecret"`
	CallbackURL    string `json:"callbackUrl"`
	Scope          string `json:"scope"`          // 应用授权作用域
	QrCodeState    string `json:"qrCodeState"`    // 二维码状态参数
	RedirectDomain string `json:"redirectDomain"` // 重定向域名
}

// VulnDBSettings 漏洞库API设置
type VulnDBSettings struct {
	Enabled    bool   `json:"enabled"`
	Provider   string `json:"provider"`   // 漏洞库提供商: weibu, vulniq, vuldb等
	APIURL     string `json:"apiUrl"`     // API地址
	APIKey     string `json:"apiKey"`     // API密钥
	APISecret  string `json:"apiSecret"`  // API密钥对应的Secret (可选)
	Parameters string `json:"parameters"` // 其他参数 (JSON格式)
}

// WorkWechatSettings 企业微信机器人设置
type WorkWechatSettings struct {
	Enabled    bool     `json:"enabled"`
	WebhookURL string   `json:"webhookUrl"`
	Events     []string `json:"events"`
}

// FeishuSettings 飞书机器人设置
type FeishuSettings struct {
	Enabled    bool     `json:"enabled"`
	WebhookURL string   `json:"webhookUrl"`
	Secret     string   `json:"secret"`
	Events     []string `json:"events"`
}

// DingtalkSettings 钉钉机器人设置
type DingtalkSettings struct {
	Enabled    bool     `json:"enabled"`
	WebhookURL string   `json:"webhookUrl"`
	Secret     string   `json:"secret"`
	Events     []string `json:"events"`
}

// EmailSettings 邮件通知设置
type EmailSettings struct {
	Enabled    bool     `json:"enabled"`
	SMTPServer string   `json:"smtpServer"`
	SMTPPort   int      `json:"smtpPort"`
	FromEmail  string   `json:"fromEmail"`
	Username   string   `json:"username"`
	Password   string   `json:"password"`
	UseSSL     bool     `json:"useSsl"`
	Events     []string `json:"events"`
	Recipients []string `json:"recipients"`
}

// AISettings AI功能设置
type AISettings struct {
	Enabled         bool     `json:"enabled"`
	Provider        string   `json:"provider"`
	APIKEY          string   `json:"apiKey"`
	APIURL          string   `json:"apiUrl"`
	AnalysisOptions []string `json:"analysisOptions"`
}

// IntegrationSettings 集成设置
type IntegrationSettings struct {
	JIRA   JIRASettings   `json:"jira"`
	Wechat WechatSettings `json:"wechat"`
	VulnDB VulnDBSettings `json:"vulnDb"`
}

// NotificationSettings 通知设置
type NotificationSettings struct {
	WorkWechat WorkWechatSettings `json:"workWechat"`
	Feishu     FeishuSettings     `json:"feishu"`
	Dingtalk   DingtalkSettings   `json:"dingtalk"`
	Email      EmailSettings      `json:"email"`
}

// Settings 系统设置
type Settings struct {
	ID            uint                 `json:"id" gorm:"primary_key"`
	Integrations  IntegrationSettings  `json:"integrations" gorm:"type:json"`
	Notifications NotificationSettings `json:"notifications" gorm:"type:json"`
	AI            AISettings           `json:"ai" gorm:"type:json"`
	UpdatedAt     time.Time            `json:"updated_at"`
	UpdatedBy     uint                 `json:"updated_by"`
}

// TableName 表名
func (Settings) TableName() string {
	return "settings"
}

// Setting 单个系统设置项，用于记录设置项的变更
type Setting struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	Name      string    `json:"name" gorm:"type:varchar(255);not null;index"`
	Category  string    `json:"category" gorm:"type:varchar(100);not null"`
	Value     string    `json:"value" gorm:"type:text"`
	UpdatedBy uint      `json:"updated_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 表名
func (Setting) TableName() string {
	return "setting_logs"
}
