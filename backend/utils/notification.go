package utils

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/smtp"
	"net/url"
	"strings"
	"time"

	"github.com/vulnark/vulnark/models"
)

// 事件类型常量
const (
	EventAssetCreate      = "资产新增"
	EventAssetUpdate      = "资产更新"
	EventAssetDelete      = "资产删除"
	EventVulnCreate       = "漏洞新增"
	EventVulnStatusChange = "漏洞状态变更"
	EventVulnUpdate       = "漏洞更新"
	EventVulnDelete       = "漏洞删除"
)

// NotificationManager 通知管理器
type NotificationManager struct {
	settings *models.Settings
}

// NewNotificationManager 创建通知管理器
func NewNotificationManager() (*NotificationManager, error) {
	log.Printf("初始化通知管理器...")

	// 直接使用原生SQL查询获取设置
	var (
		id                                          uint
		integrationsJSON, notificationsJSON, aiJSON []byte
		updatedAt                                   time.Time
		updatedBy                                   uint
	)

	// 使用原生SQL查询避免GORM的自动JSON转换
	row := DB.Raw("SELECT id, integrations, notifications, ai, updated_at, updated_by FROM settings WHERE id = ? LIMIT 1", 1).Row()
	if err := row.Scan(&id, &integrationsJSON, &notificationsJSON, &aiJSON, &updatedAt, &updatedBy); err != nil {
		log.Printf("设置数据查询失败: %v", err)

		// 如果查询失败，尝试使用默认设置
		log.Printf("尝试返回默认设置...")
		return &NotificationManager{
			settings: getDefaultSettings(),
		}, nil
	}

	// 创建设置对象
	settings := &models.Settings{
		ID:        id,
		UpdatedAt: updatedAt,
		UpdatedBy: updatedBy,
	}

	// 记录JSON原始数据（仅用于调试）
	log.Printf("集成设置JSON: %s", string(integrationsJSON))
	log.Printf("通知设置JSON: %s", string(notificationsJSON))
	log.Printf("AI设置JSON: %s", string(aiJSON))

	// 手动解析JSON字段
	if err := json.Unmarshal(integrationsJSON, &settings.Integrations); err != nil {
		log.Printf("解析集成设置JSON失败: %v，将使用默认集成设置", err)
		settings.Integrations = getDefaultSettings().Integrations
	}

	if err := json.Unmarshal(notificationsJSON, &settings.Notifications); err != nil {
		log.Printf("解析通知设置JSON失败: %v，将使用默认通知设置", err)
		settings.Notifications = getDefaultSettings().Notifications
	}

	if err := json.Unmarshal(aiJSON, &settings.AI); err != nil {
		log.Printf("解析AI设置JSON失败: %v，将使用默认AI设置", err)
		settings.AI = getDefaultSettings().AI
	}

	log.Printf("通知管理器初始化成功")

	return &NotificationManager{
		settings: settings,
	}, nil
}

// getDefaultSettings 返回默认的设置对象
func getDefaultSettings() *models.Settings {
	return &models.Settings{
		ID: 1,
		Integrations: models.IntegrationSettings{
			JIRA: models.JIRASettings{
				Enabled: false,
			},
			Wechat: models.WechatSettings{
				Enabled: false,
			},
			VulnDB: models.VulnDBSettings{
				Enabled: false,
			},
		},
		Notifications: models.NotificationSettings{
			WorkWechat: models.WorkWechatSettings{
				Enabled: false,
				Events:  []string{},
			},
			Feishu: models.FeishuSettings{
				Enabled: false,
				Events:  []string{},
			},
			Dingtalk: models.DingtalkSettings{
				Enabled: false,
				Events:  []string{},
			},
			Email: models.EmailSettings{
				Enabled:    false,
				Events:     []string{},
				Recipients: []string{},
			},
		},
		AI: models.AISettings{
			Enabled:         false,
			AnalysisOptions: []string{},
		},
		UpdatedAt: time.Now(),
		UpdatedBy: 1,
	}
}

// GetSettings 获取通知设置
func (m *NotificationManager) GetSettings() *models.Settings {
	return m.settings
}

// 检查事件是否需要通知
func containsEvent(events []string, event string) bool {
	if len(events) == 0 {
		log.Printf("事件列表为空，事件: %s 不匹配", event)
		return false
	}

	log.Printf("检查事件 %s 是否在列表中: %v", event, events)

	for _, e := range events {
		// 添加更多日志以确认比较过程
		log.Printf("比较事件: %s vs %s", e, event)

		// 尝试多种匹配方式
		if e == event || strings.Contains(e, event) || strings.Contains(event, e) {
			log.Printf("事件匹配成功: %s", event)
			return true
		}
	}

	log.Printf("事件 %s 不在列表中", event)
	return false
}

// SendAssetNotification 发送资产相关通知
func (m *NotificationManager) SendAssetNotification(event string, asset *models.Asset) {
	// 添加详细日志
	log.Printf("开始处理资产通知, 事件: %s, 资产ID: %d, 名称: %s", event, asset.ID, asset.Name)

	// 记录设置状态
	log.Printf("通知设置状态: 企业微信=%v, 飞书=%v, 钉钉=%v, 邮件=%v",
		m.settings.Notifications.WorkWechat.Enabled,
		m.settings.Notifications.Feishu.Enabled,
		m.settings.Notifications.Dingtalk.Enabled,
		m.settings.Notifications.Email.Enabled)

	// 记录事件列表
	log.Printf("企业微信事件列表: %v", m.settings.Notifications.WorkWechat.Events)
	log.Printf("飞书事件列表: %v", m.settings.Notifications.Feishu.Events)
	log.Printf("钉钉事件列表: %v", m.settings.Notifications.Dingtalk.Events)
	log.Printf("邮件事件列表: %v", m.settings.Notifications.Email.Events)

	// 构建消息内容
	assetType := string(asset.Type)
	assetStatus := string(asset.Status)

	var title string
	var content string

	switch event {
	case EventAssetCreate:
		title = "【资产新增】" + asset.Name
		content = fmt.Sprintf("资产名称: %s\n资产类型: %s\nIP地址: %s\n状态: %s\n部门: %s\n负责人: %s",
			asset.Name, assetType, asset.IPAddress, assetStatus, asset.Department, asset.Owner)
	case EventAssetUpdate:
		title = "【资产更新】" + asset.Name
		content = fmt.Sprintf("资产名称: %s\n资产类型: %s\nIP地址: %s\n状态: %s\n部门: %s\n负责人: %s",
			asset.Name, assetType, asset.IPAddress, assetStatus, asset.Department, asset.Owner)
	case EventAssetDelete:
		title = "【资产删除】" + asset.Name
		content = fmt.Sprintf("资产名称: %s\n资产类型: %s\nIP地址: %s\n部门: %s",
			asset.Name, assetType, asset.IPAddress, asset.Department)
	}

	// 发送各种类型的通知
	m.sendWorkWechatNotification(event, title, content)
	m.sendFeishuNotification(event, title, content)
	m.sendDingtalkNotification(event, title, content)
	m.sendEmailNotification(event, title, content)
}

// SendVulnerabilityNotification 发送漏洞相关通知
func (m *NotificationManager) SendVulnerabilityNotification(event string, vuln *models.Vulnerability, oldStatus string) {
	// 添加详细日志
	log.Printf("开始处理漏洞通知, 事件: %s, 漏洞ID: %d, 标题: %s", event, vuln.ID, vuln.Title)

	// 记录设置状态
	log.Printf("通知设置状态: 企业微信=%v, 飞书=%v, 钉钉=%v, 邮件=%v",
		m.settings.Notifications.WorkWechat.Enabled,
		m.settings.Notifications.Feishu.Enabled,
		m.settings.Notifications.Dingtalk.Enabled,
		m.settings.Notifications.Email.Enabled)

	// 记录事件列表
	log.Printf("企业微信事件列表: %v", m.settings.Notifications.WorkWechat.Events)
	log.Printf("飞书事件列表: %v", m.settings.Notifications.Feishu.Events)
	log.Printf("钉钉事件列表: %v", m.settings.Notifications.Dingtalk.Events)
	log.Printf("邮件事件列表: %v", m.settings.Notifications.Email.Events)

	// 构建消息内容
	var title string
	var content string

	severityText := ""
	switch vuln.Severity {
	case models.SeverityCritical:
		severityText = "严重"
	case models.SeverityHigh:
		severityText = "高危"
	case models.SeverityMedium:
		severityText = "中危"
	case models.SeverityLow:
		severityText = "低危"
	case models.SeverityInfo:
		severityText = "信息"
	default:
		severityText = string(vuln.Severity)
	}

	statusText := ""
	switch vuln.Status {
	case models.StatusNew:
		statusText = "新发现"
	case models.StatusVerified:
		statusText = "已验证"
	case models.StatusInProgress:
		statusText = "处理中"
	case models.StatusFixed:
		statusText = "已修复"
	case models.StatusClosed:
		statusText = "已关闭"
	case models.StatusFalsePositive:
		statusText = "误报"
	default:
		statusText = string(vuln.Status)
	}

	switch event {
	case EventVulnCreate:
		title = fmt.Sprintf("【新增漏洞】%s (%s)", vuln.Title, severityText)
		content = fmt.Sprintf("漏洞名称: %s\n严重程度: %s\n状态: %s\nCVE: %s\n",
			vuln.Title, severityText, statusText, vuln.CVE)
	case EventVulnStatusChange:
		oldStatusText := ""
		switch models.VulnStatus(oldStatus) {
		case models.StatusNew:
			oldStatusText = "新发现"
		case models.StatusVerified:
			oldStatusText = "已验证"
		case models.StatusInProgress:
			oldStatusText = "处理中"
		case models.StatusFixed:
			oldStatusText = "已修复"
		case models.StatusClosed:
			oldStatusText = "已关闭"
		case models.StatusFalsePositive:
			oldStatusText = "误报"
		default:
			oldStatusText = oldStatus
		}

		title = fmt.Sprintf("【漏洞状态变更】%s (%s)", vuln.Title, severityText)
		content = fmt.Sprintf("漏洞名称: %s\n严重程度: %s\n状态变更: %s → %s\nCVE: %s\n",
			vuln.Title, severityText, oldStatusText, statusText, vuln.CVE)
	case EventVulnUpdate:
		title = fmt.Sprintf("【漏洞更新】%s (%s)", vuln.Title, severityText)
		content = fmt.Sprintf("漏洞名称: %s\n严重程度: %s\n状态: %s\nCVE: %s\n",
			vuln.Title, severityText, statusText, vuln.CVE)
	case EventVulnDelete:
		title = fmt.Sprintf("【漏洞删除】%s (%s)", vuln.Title, severityText)
		content = fmt.Sprintf("漏洞名称: %s\n严重程度: %s\nCVE: %s\n",
			vuln.Title, severityText, vuln.CVE)
	}

	// 发送各种类型的通知
	m.sendWorkWechatNotification(event, title, content)
	m.sendFeishuNotification(event, title, content)
	m.sendDingtalkNotification(event, title, content)
	m.sendEmailNotification(event, title, content)
}

// 企业微信通知
func (m *NotificationManager) sendWorkWechatNotification(event, title, content string) {
	// 检查是否启用企业微信通知
	if !m.settings.Notifications.WorkWechat.Enabled {
		log.Printf("企业微信通知未启用，跳过发送。事件: %s", event)
		return
	}

	// 检查事件是否在通知列表中
	if !containsEvent(m.settings.Notifications.WorkWechat.Events, event) {
		log.Printf("事件 %s 不在企业微信通知事件列表中，跳过发送", event)
		return
	}

	// 检查WebhookURL是否配置
	webhookURL := m.settings.Notifications.WorkWechat.WebhookURL
	if webhookURL == "" {
		log.Println("企业微信WebhookURL未配置")
		return
	}

	log.Printf("准备发送企业微信通知，事件: %s，WebhookURL: %s", event, webhookURL)

	// 构建请求体
	requestBody := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"content": fmt.Sprintf("### %s\n%s", title, content),
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("JSON序列化企业微信通知失败: %v", err)
		return
	}

	// 发送请求
	resp, err := http.Post(
		webhookURL,
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		log.Printf("发送企业微信通知失败: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("企业微信通知响应异常, 状态码: %d", resp.StatusCode)
		return
	}

	log.Printf("企业微信通知发送成功: %s", title)
}

// 飞书通知
func (m *NotificationManager) sendFeishuNotification(event, title, content string) {
	// 检查是否启用飞书通知
	if !m.settings.Notifications.Feishu.Enabled {
		log.Printf("飞书通知未启用，跳过发送。事件: %s", event)
		return
	}

	// 检查事件是否在通知列表中
	if !containsEvent(m.settings.Notifications.Feishu.Events, event) {
		log.Printf("事件 %s 不在飞书通知事件列表中，跳过发送", event)
		return
	}

	// 检查WebhookURL是否配置
	webhookURL := m.settings.Notifications.Feishu.WebhookURL
	if webhookURL == "" {
		log.Println("飞书WebhookURL未配置")
		return
	}

	log.Printf("准备发送飞书通知，事件: %s，WebhookURL: %s", event, webhookURL)

	// 构建请求体
	requestBody := map[string]interface{}{
		"msg_type": "interactive",
		"card": map[string]interface{}{
			"header": map[string]interface{}{
				"title": map[string]interface{}{
					"tag":     "plain_text",
					"content": title,
				},
				"template": "blue",
			},
			"elements": []map[string]interface{}{
				{
					"tag": "div",
					"text": map[string]interface{}{
						"tag":     "lark_md",
						"content": content,
					},
				},
				{
					"tag": "note",
					"elements": []map[string]interface{}{
						{
							"tag":     "plain_text",
							"content": fmt.Sprintf("发送时间: %s", FormatTimeCST(NowCST())),
						},
					},
				},
			},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("JSON序列化飞书通知失败: %v", err)
		return
	}

	// 如果配置了签名密钥，需要计算签名
	var requestBytes []byte
	if secret := m.settings.Notifications.Feishu.Secret; secret != "" {
		// 计算签名
		timestamp := time.Now().Unix()
		stringToSign := fmt.Sprintf("%d\n%s", timestamp, string(jsonData))

		mac := hmac.New(sha256.New, []byte(secret))
		mac.Write([]byte(stringToSign))
		signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

		// 添加签名到请求
		signedRequest := map[string]interface{}{
			"timestamp": timestamp,
			"sign":      signature,
		}
		for k, v := range requestBody {
			signedRequest[k] = v
		}

		// 重新序列化
		requestBytes, err = json.Marshal(signedRequest)
		if err != nil {
			log.Printf("JSON序列化飞书带签名通知失败: %v", err)
			return
		}
	} else {
		requestBytes = jsonData
	}

	// 发送请求
	resp, err := http.Post(
		webhookURL,
		"application/json",
		bytes.NewBuffer(requestBytes),
	)

	if err != nil {
		log.Printf("发送飞书通知失败: %v", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应内容
	responseBody, _ := io.ReadAll(resp.Body)
	log.Printf("飞书响应: %s, 状态码: %d", string(responseBody), resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		log.Printf("飞书通知响应异常, 状态码: %d", resp.StatusCode)
		return
	}

	log.Printf("飞书通知发送成功: %s", title)
}

// 钉钉通知
func (m *NotificationManager) sendDingtalkNotification(event, title, content string) {
	// 检查是否启用钉钉通知
	if !m.settings.Notifications.Dingtalk.Enabled {
		log.Printf("钉钉通知未启用，跳过发送。事件: %s", event)
		return
	}

	// 检查事件是否在通知列表中
	if !containsEvent(m.settings.Notifications.Dingtalk.Events, event) {
		log.Printf("事件 %s 不在钉钉通知事件列表中，跳过发送", event)
		return
	}

	// 检查WebhookURL是否配置
	webhookURL := m.settings.Notifications.Dingtalk.WebhookURL
	if webhookURL == "" {
		log.Println("钉钉WebhookURL未配置")
		return
	}

	log.Printf("准备发送钉钉通知，事件: %s，WebhookURL: %s", event, webhookURL)

	// 构建请求体
	requestBody := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"title": title,
			"text":  fmt.Sprintf("### %s\n%s\n\n###### 发送时间: %s", title, content, FormatTimeCST(NowCST())),
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("JSON序列化钉钉通知失败: %v", err)
		return
	}

	// 处理钉钉安全设置
	finalURL := webhookURL
	if secret := m.settings.Notifications.Dingtalk.Secret; secret != "" {
		timestamp := time.Now().UnixNano() / 1e6
		stringToSign := fmt.Sprintf("%d\n%s", timestamp, secret)

		mac := hmac.New(sha256.New, []byte(secret))
		mac.Write([]byte(stringToSign))
		signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))
		signatureEncoded := url.QueryEscape(signature)

		// 添加签名参数到URL
		if strings.Contains(finalURL, "?") {
			finalURL = fmt.Sprintf("%s&timestamp=%d&sign=%s", finalURL, timestamp, signatureEncoded)
		} else {
			finalURL = fmt.Sprintf("%s?timestamp=%d&sign=%s", finalURL, timestamp, signatureEncoded)
		}
	}

	// 发送请求
	resp, err := http.Post(
		finalURL,
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		log.Printf("发送钉钉通知失败: %v", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应内容
	responseBody, _ := io.ReadAll(resp.Body)
	log.Printf("钉钉响应: %s, 状态码: %d", string(responseBody), resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		log.Printf("钉钉通知响应异常, 状态码: %d", resp.StatusCode)
		return
	}

	log.Printf("钉钉通知发送成功: %s", title)
}

// 邮件通知
func (m *NotificationManager) sendEmailNotification(event, title, content string) {
	// 检查是否启用邮件通知
	if !m.settings.Notifications.Email.Enabled {
		log.Printf("邮件通知未启用，跳过发送。事件: %s", event)
		return
	}

	// 检查事件是否在通知列表中
	if !containsEvent(m.settings.Notifications.Email.Events, event) {
		log.Printf("事件 %s 不在邮件通知事件列表中，跳过发送", event)
		return
	}

	// 检查Recipients是否配置
	recipients := m.settings.Notifications.Email.Recipients
	if len(recipients) == 0 {
		log.Println("邮件Recipients未配置")
		return
	}

	// 检查SMTP配置
	smtpServer := m.settings.Notifications.Email.SMTPServer
	smtpPort := m.settings.Notifications.Email.SMTPPort
	fromEmail := m.settings.Notifications.Email.FromEmail
	username := m.settings.Notifications.Email.Username
	password := m.settings.Notifications.Email.Password
	useSSL := m.settings.Notifications.Email.UseSSL

	if smtpServer == "" || smtpPort == 0 || fromEmail == "" {
		log.Println("邮件SMTP配置不完整")
		return
	}

	log.Printf("准备发送邮件通知，事件: %s，Recipients: %v", event, recipients)

	// 构建SMTP地址
	smtpAddr := fmt.Sprintf("%s:%d", smtpServer, smtpPort)

	// 构建认证
	var auth smtp.Auth
	if username != "" && password != "" {
		auth = smtp.PlainAuth("", username, password, smtpServer)
	}

	// 构建邮件内容
	headers := map[string]string{
		"From":         fromEmail,
		"To":           strings.Join(recipients, ", "),
		"Subject":      title,
		"MIME-Version": "1.0",
		"Content-Type": "text/html; charset=UTF-8",
	}

	// 构建头部字符串
	var headerStr strings.Builder
	for key, value := range headers {
		headerStr.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}

	// 构建HTML内容
	htmlContent := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>%s</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #ddd; border-radius: 5px; }
        .header { background-color: #4e54c8; color: white; padding: 10px; text-align: center; border-radius: 5px 5px 0 0; }
        .content { padding: 20px; white-space: pre-line; }
        .footer { background-color: #f5f5f5; padding: 10px; text-align: center; font-size: 12px; border-radius: 0 0 5px 5px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h2>%s</h2>
        </div>
        <div class="content">
            %s
        </div>
        <div class="footer">
            <p>发送时间: %s</p>
            <p>此邮件由VulnArk系统自动发送，请勿回复。</p>
        </div>
    </div>
</body>
</html>
`, title, title, content, FormatTimeCST(NowCST()))

	// 完整邮件内容
	message := headerStr.String() + "\r\n" + htmlContent

	// 发送邮件
	var err error
	if useSSL {
		// 使用TLS加密通信
		// TLS配置
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true, // 在测试环境中可以跳过证书验证
			ServerName:         smtpServer,
		}

		// 连接SMTP服务器
		conn, err := tls.Dial("tcp", smtpAddr, tlsConfig)
		if err != nil {
			log.Printf("连接SMTP服务器失败: %v", err)
			return
		}

		// 创建SMTP客户端
		client, err := smtp.NewClient(conn, smtpServer)
		if err != nil {
			log.Printf("创建SMTP客户端失败: %v", err)
			return
		}
		defer client.Close()

		// 设置身份验证
		if auth != nil {
			if err := client.Auth(auth); err != nil {
				log.Printf("SMTP身份验证失败: %v", err)
				return
			}
		}

		// 设置发件人
		if err := client.Mail(fromEmail); err != nil {
			log.Printf("设置发件人失败: %v", err)
			return
		}

		// 设置收件人
		for _, recipient := range recipients {
			if err := client.Rcpt(recipient); err != nil {
				log.Printf("设置收件人 %s 失败: %v", recipient, err)
				continue
			}
		}

		// 设置邮件内容
		w, err := client.Data()
		if err != nil {
			log.Printf("准备邮件内容失败: %v", err)
			return
		}

		_, err = w.Write([]byte(message))
		if err != nil {
			log.Printf("写入邮件内容失败: %v", err)
			return
		}

		err = w.Close()
		if err != nil {
			log.Printf("关闭写入器失败: %v", err)
			return
		}

		// 结束会话
		client.Quit()
	} else {
		// 不使用SSL
		err = smtp.SendMail(
			smtpAddr,
			auth,
			fromEmail,
			recipients,
			[]byte(message),
		)
	}

	if err != nil {
		log.Printf("发送邮件通知失败: %v", err)
		return
	}

	log.Printf("邮件通知发送成功: %s", title)
}
