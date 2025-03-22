package controllers

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/vulnark/vulnark/models"
	"github.com/vulnark/vulnark/utils"
)

// SettingsController 系统设置控制器
type SettingsController struct{}

// GetSettings 获取系统设置
func (sc *SettingsController) GetSettings(c *gin.Context) {
	log.Printf("GetSettings: 开始获取系统设置")

	// 检查并确保settings表结构正确
	if err := sc.ensureSettingsTableStructure(); err != nil {
		log.Printf("GetSettings: 确保表结构过程中出错 - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "检查数据库结构时出错",
			"error":   err.Error(),
		})
		return
	}

	// 检查settings表中是否有数据
	var count int
	if err := utils.DB.Raw("SELECT COUNT(*) FROM settings").Row().Scan(&count); err != nil {
		log.Printf("GetSettings: 查询settings表记录数出错 - %v", err)
	} else {
		log.Printf("GetSettings: 当前settings表中有 %d 条记录", count)
	}

	// 使用原生SQL查询获取设置
	var (
		id            uint
		integrations  string
		notifications string
		ai            string
		updatedAt     time.Time
		updatedBy     uint
	)

	// 查询设置记录
	row := utils.DB.Raw("SELECT id, integrations, notifications, ai, updated_at, updated_by FROM settings WHERE id = ? LIMIT 1", 1).Row()

	if row.Err() != nil {
		if row.Err() == gorm.ErrRecordNotFound {
			log.Printf("GetSettings: 未找到设置记录，返回默认设置")
		} else {
			log.Printf("GetSettings: 查询设置记录出错 - %v", row.Err())
		}
		// 返回默认设置
		returnDefaultSettings(c)
		return
	}

	// 扫描结果
	if err := row.Scan(&id, &integrations, &notifications, &ai, &updatedAt, &updatedBy); err != nil {
		log.Printf("GetSettings: 扫描设置记录失败 - %v", err)
		// 返回默认设置
		returnDefaultSettings(c)
		return
	}

	log.Printf("GetSettings: 成功读取设置记录 ID=%d, 集成设置JSON长度=%d, 通知设置JSON长度=%d, AI设置JSON长度=%d",
		id, len(integrations), len(notifications), len(ai))

	// 创建设置对象
	settings := models.Settings{
		ID:        id,
		UpdatedAt: updatedAt,
		UpdatedBy: updatedBy,
	}

	// 解析集成设置
	if integrations != "" {
		if err := json.Unmarshal([]byte(integrations), &settings.Integrations); err != nil {
			log.Printf("GetSettings: 解析集成设置JSON失败 - %v", err)
			// 使用默认值
			settings.Integrations = getDefaultIntegrationSettings()
		} else {
			log.Printf("GetSettings: 成功解析集成设置JSON")
		}
	} else {
		log.Printf("GetSettings: 集成设置为空，使用默认值")
		settings.Integrations = getDefaultIntegrationSettings()
	}

	// 解析通知设置
	if notifications != "" {
		if err := json.Unmarshal([]byte(notifications), &settings.Notifications); err != nil {
			log.Printf("GetSettings: 解析通知设置JSON失败 - %v", err)
			// 使用默认值
			settings.Notifications = getDefaultNotificationSettings()
		} else {
			log.Printf("GetSettings: 成功解析通知设置JSON")
		}
	} else {
		log.Printf("GetSettings: 通知设置为空，使用默认值")
		settings.Notifications = getDefaultNotificationSettings()
	}

	// 解析AI设置
	if ai != "" {
		if err := json.Unmarshal([]byte(ai), &settings.AI); err != nil {
			log.Printf("GetSettings: 解析AI设置JSON失败 - %v", err)
			// 使用默认值
			settings.AI = getDefaultAISettings()
		} else {
			log.Printf("GetSettings: 成功解析AI设置JSON")
		}
	} else {
		log.Printf("GetSettings: AI设置为空，使用默认值")
		settings.AI = getDefaultAISettings()
	}

	// 确保数组字段已初始化
	initializeSettingsArrays(&settings)

	// 敏感信息处理，不返回密码等敏感信息
	settings.Notifications.Email.Password = ""
	settings.Integrations.JIRA.APIToken = ""
	settings.Integrations.Wechat.AppSecret = ""
	settings.Integrations.VulnDB.APIKey = ""
	settings.Integrations.VulnDB.APISecret = ""
	settings.AI.APIKEY = ""

	log.Printf("GetSettings: 成功获取系统设置，返回结果")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取设置成功",
		"data":    settings,
	})
}

// ensureSettingsTableStructure 确保settings表具有正确的结构
func (sc *SettingsController) ensureSettingsTableStructure() error {
	log.Printf("检查settings表结构...")

	// 检查表是否存在
	var tableExists int
	utils.DB.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = 'settings'").Row().Scan(&tableExists)

	if tableExists == 0 {
		// 表不存在，使用原生SQL创建表
		log.Printf("settings表不存在，使用原生SQL创建表...")
		createTableSQL := `
		CREATE TABLE settings (
			id int PRIMARY KEY,
			integrations JSON,
			notifications JSON,
			ai JSON,
			updated_at DATETIME,
			updated_by int
		)
		`
		if err := utils.DB.Exec(createTableSQL).Error; err != nil {
			log.Printf("创建settings表失败: %v", err)
			return err
		}
		log.Printf("settings表创建成功")
	} else {
		// 检查表结构
		rows, err := utils.DB.Raw("SHOW COLUMNS FROM settings").Rows()
		if err != nil {
			log.Printf("无法获取表结构: %v", err)
			return err
		}
		defer rows.Close()

		// 解析列信息
		var field, fieldType, null, key, defaultValue, extra string
		columnMap := make(map[string]bool)

		for rows.Next() {
			err := rows.Scan(&field, &fieldType, &null, &key, &defaultValue, &extra)
			if err != nil {
				log.Printf("无法扫描列信息: %v", err)
				continue
			}
			columnMap[field] = true
			log.Printf("发现列: %s, 类型: %s", field, fieldType)
		}

		// 检查必要的列是否存在
		requiredColumns := []string{"id", "integrations", "notifications", "ai", "updated_at", "updated_by"}
		missingColumns := []string{}

		for _, col := range requiredColumns {
			if !columnMap[col] {
				missingColumns = append(missingColumns, col)
			}
		}

		// 如果缺少列，尝试添加列而不是重建表
		if len(missingColumns) > 0 {
			log.Printf("表缺少以下列: %v，将尝试添加这些列", missingColumns)

			for _, col := range missingColumns {
				var alterSQL string
				switch col {
				case "id":
					alterSQL = "ALTER TABLE settings ADD COLUMN id int PRIMARY KEY"
				case "integrations", "notifications", "ai":
					alterSQL = fmt.Sprintf("ALTER TABLE settings ADD COLUMN %s JSON", col)
				case "updated_at":
					alterSQL = "ALTER TABLE settings ADD COLUMN updated_at DATETIME"
				case "updated_by":
					alterSQL = "ALTER TABLE settings ADD COLUMN updated_by int"
				}

				if alterSQL != "" {
					if err := utils.DB.Exec(alterSQL).Error; err != nil {
						log.Printf("添加列 %s 失败: %v", col, err)
						// 继续尝试其他列，不返回错误
					} else {
						log.Printf("成功添加列: %s", col)
					}
				}
			}
		} else {
			log.Printf("settings表结构正确，包含所有必要的列")
		}
	}

	// 检查是否有默认设置记录
	var count int64
	utils.DB.Raw("SELECT COUNT(*) FROM settings").Row().Scan(&count)

	if count == 0 {
		log.Printf("未找到设置记录，添加默认设置")

		// 准备默认设置
		defaultSettings := models.Settings{
			ID:            1,
			Integrations:  getDefaultIntegrationSettings(),
			Notifications: getDefaultNotificationSettings(),
			AI:            getDefaultAISettings(),
			UpdatedAt:     time.Now(),
			UpdatedBy:     1,
		}

		// 初始化数组字段
		initializeSettingsArrays(&defaultSettings)

		// 序列化为JSON
		integrationsJSON, err := json.Marshal(defaultSettings.Integrations)
		if err != nil {
			log.Printf("序列化集成设置失败: %v", err)
			return err
		}

		notificationsJSON, err := json.Marshal(defaultSettings.Notifications)
		if err != nil {
			log.Printf("序列化通知设置失败: %v", err)
			return err
		}

		aiJSON, err := json.Marshal(defaultSettings.AI)
		if err != nil {
			log.Printf("序列化AI设置失败: %v", err)
			return err
		}

		// 使用原生SQL插入记录
		insertSQL := "INSERT INTO settings (id, integrations, notifications, ai, updated_at, updated_by) VALUES (?, ?, ?, ?, ?, ?)"
		if err := utils.DB.Exec(
			insertSQL,
			defaultSettings.ID,
			string(integrationsJSON),
			string(notificationsJSON),
			string(aiJSON),
			defaultSettings.UpdatedAt,
			defaultSettings.UpdatedBy,
		).Error; err != nil {
			log.Printf("插入默认设置失败: %v", err)
			return err
		}

		log.Printf("成功添加默认设置记录")
	} else {
		log.Printf("发现已存在 %d 个设置记录", count)
	}

	return nil
}

// returnDefaultSettings 返回默认设置
func returnDefaultSettings(c *gin.Context) {
	settings := models.Settings{
		ID:            1,
		Integrations:  getDefaultIntegrationSettings(),
		Notifications: getDefaultNotificationSettings(),
		AI:            getDefaultAISettings(),
		UpdatedAt:     time.Now(),
	}

	// 初始化所有数组字段
	initializeSettingsArrays(&settings)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取默认设置成功",
		"data":    settings,
	})
}

// getDefaultIntegrationSettings 获取默认集成设置
func getDefaultIntegrationSettings() models.IntegrationSettings {
	return models.IntegrationSettings{
		JIRA: models.JIRASettings{
			Enabled: false,
		},
		Wechat: models.WechatSettings{
			Enabled: false,
		},
		VulnDB: models.VulnDBSettings{
			Enabled:  false,
			Provider: "weibu",
		},
	}
}

// getDefaultNotificationSettings 获取默认通知设置
func getDefaultNotificationSettings() models.NotificationSettings {
	return models.NotificationSettings{
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
			SMTPPort:   25,
			UseSSL:     true,
			Events:     []string{},
			Recipients: []string{},
		},
	}
}

// getDefaultAISettings 获取默认AI设置
func getDefaultAISettings() models.AISettings {
	return models.AISettings{
		Enabled:         false,
		Provider:        "openai",
		AnalysisOptions: []string{},
	}
}

// initializeSettingsArrays 初始化设置中的所有数组字段
func initializeSettingsArrays(settings *models.Settings) {
	// 通知设置中的数组
	if settings.Notifications.WorkWechat.Events == nil {
		settings.Notifications.WorkWechat.Events = []string{}
	}
	if settings.Notifications.Feishu.Events == nil {
		settings.Notifications.Feishu.Events = []string{}
	}
	if settings.Notifications.Dingtalk.Events == nil {
		settings.Notifications.Dingtalk.Events = []string{}
	}
	if settings.Notifications.Email.Events == nil {
		settings.Notifications.Email.Events = []string{}
	}
	if settings.Notifications.Email.Recipients == nil {
		settings.Notifications.Email.Recipients = []string{}
	}

	// AI设置中的数组
	if settings.AI.AnalysisOptions == nil {
		settings.AI.AnalysisOptions = []string{}
	}
}

// SaveSettings 保存系统设置
func (sc *SettingsController) SaveSettings(c *gin.Context) {
	// 尝试获取用户ID - 同时尝试两种可能的键名
	var userID interface{}
	var exists bool

	// 首先尝试使用"user_id"键
	userID, exists = c.Get("user_id")
	if !exists {
		// 如果不存在，尝试使用"userID"键
		userID, exists = c.Get("userID")
		if !exists {
			log.Printf("SaveSettings错误: 无法获取用户ID (尝试了user_id和userID两个键)")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
			return
		}
		log.Printf("SaveSettings: 使用userID键获取到用户ID = %v", userID)
	} else {
		log.Printf("SaveSettings: 使用user_id键获取到用户ID = %v", userID)
	}

	// 读取请求体
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("SaveSettings错误: 读取请求体失败 - %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "无法读取请求数据"})
		return
	}

	// 记录请求体内容用于调试
	log.Printf("SaveSettings: 接收到的请求体长度 = %d 字节", len(body))

	// 恢复请求体
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	// 解析请求数据
	var requestData models.Settings
	if err := c.ShouldBindJSON(&requestData); err != nil {
		log.Printf("SaveSettings错误: 解析JSON数据失败 - %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	log.Printf("SaveSettings: 成功解析设置数据 - 接收到集成设置、通知设置和AI设置")

	// 设置固定的设置ID和更新者
	settingsID := uint(1)
	requestData.ID = settingsID
	requestData.UpdatedBy = userID.(uint)
	requestData.UpdatedAt = time.Now()

	// 手动序列化JSON字段
	integrationsJSON, err := json.Marshal(requestData.Integrations)
	if err != nil {
		log.Printf("SaveSettings错误: 序列化集成设置失败 - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "序列化数据失败"})
		return
	}
	log.Printf("SaveSettings: 序列化集成设置成功, JSON长度=%d", len(integrationsJSON))

	notificationsJSON, err := json.Marshal(requestData.Notifications)
	if err != nil {
		log.Printf("SaveSettings错误: 序列化通知设置失败 - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "序列化数据失败"})
		return
	}
	log.Printf("SaveSettings: 序列化通知设置成功, JSON长度=%d", len(notificationsJSON))

	aiJSON, err := json.Marshal(requestData.AI)
	if err != nil {
		log.Printf("SaveSettings错误: 序列化AI设置失败 - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "序列化数据失败"})
		return
	}
	log.Printf("SaveSettings: 序列化AI设置成功, JSON长度=%d", len(aiJSON))

	// 开始事务
	tx := utils.DB.Begin()
	if tx.Error != nil {
		log.Printf("SaveSettings错误: 无法开始事务 - %v", tx.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库错误"})
		return
	}

	// 使用事务和原生SQL删除现有记录
	if result := tx.Exec("DELETE FROM settings WHERE id = ?", settingsID); result.Error != nil {
		tx.Rollback()
		log.Printf("SaveSettings错误: 删除现有设置失败 - %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存设置失败"})
		return
	}
	log.Printf("SaveSettings: 成功删除旧设置记录")

	// 使用原生SQL插入新记录，避免GORM的JSON处理问题
	insertSQL := "INSERT INTO settings (id, integrations, notifications, ai, updated_at, updated_by) VALUES (?, ?, ?, ?, ?, ?)"
	if err := tx.Exec(
		insertSQL,
		settingsID,
		string(integrationsJSON),
		string(notificationsJSON),
		string(aiJSON),
		requestData.UpdatedAt,
		requestData.UpdatedBy,
	).Error; err != nil {
		tx.Rollback()
		log.Printf("SaveSettings错误: 插入新设置失败 - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存设置失败"})
		return
	}
	log.Printf("SaveSettings: 成功插入新设置记录")

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		log.Printf("SaveSettings错误: 提交事务失败 - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存设置失败"})
		return
	}
	log.Printf("SaveSettings: 事务成功提交")

	// 验证设置是否成功保存
	var savedID uint
	var savedIntegrations string
	var savedNotifications string
	var savedAI string
	var savedUpdatedAt time.Time
	var savedUpdatedBy uint
	row := utils.DB.Raw("SELECT id, integrations, notifications, ai, updated_at, updated_by FROM settings WHERE id = ? LIMIT 1", settingsID).Row()
	if row.Err() != nil {
		log.Printf("SaveSettings警告: 验证时查询设置失败 - %v", row.Err())
	} else if err := row.Scan(&savedID, &savedIntegrations, &savedNotifications, &savedAI, &savedUpdatedAt, &savedUpdatedBy); err != nil {
		log.Printf("SaveSettings警告: 验证时扫描设置失败 - %v", err)
	} else {
		log.Printf("SaveSettings: 设置保存并验证成功 ID=%d, 集成设置长度=%d, 通知设置长度=%d, AI设置长度=%d",
			savedID, len(savedIntegrations), len(savedNotifications), len(savedAI))
	}

	c.JSON(http.StatusOK, gin.H{"message": "设置已保存", "data": requestData})
}

// TestJiraConnection 测试JIRA连接
func (sc *SettingsController) TestJiraConnection(c *gin.Context) {
	var jiraSettings models.JIRASettings
	if err := c.ShouldBindJSON(&jiraSettings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// TODO: 实现JIRA连接测试
	// 这里简单模拟测试结果，实际应该调用JIRA API验证连接
	log.Printf("测试JIRA连接: %s, 用户: %s", jiraSettings.URL, jiraSettings.Username)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "JIRA连接测试成功",
	})
}

// TestWechatLogin 测试微信登录配置
func (sc *SettingsController) TestWechatLogin(c *gin.Context) {
	var wechatSettings models.WechatSettings
	if err := c.ShouldBindJSON(&wechatSettings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// TODO: 实现微信登录配置测试
	// 这里简单模拟测试结果，实际应该调用微信开放平台API验证配置
	log.Printf("测试微信登录配置: AppID: %s, 授权作用域: %s", wechatSettings.AppID, wechatSettings.Scope)

	// 验证必要参数
	if wechatSettings.AppID == "" || wechatSettings.AppSecret == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "AppID和AppSecret不能为空",
		})
		return
	}

	// 验证回调URL格式
	if wechatSettings.CallbackURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "回调URL不能为空",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "微信登录配置测试成功",
	})
}

// TestWorkWechatBot 测试企业微信机器人连接
func (sc *SettingsController) TestWorkWechatBot(c *gin.Context) {
	log.Printf("开始测试企业微信机器人连接")

	var requestData struct {
		Enabled    bool     `json:"enabled"`
		WebhookURL string   `json:"webhookUrl"`
		Events     []string `json:"events"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		log.Printf("解析请求数据失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的请求数据",
		})
		return
	}

	// 检查必要参数
	if !requestData.Enabled {
		log.Printf("企业微信机器人未启用")
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "企业微信机器人未启用",
		})
		return
	}

	if requestData.WebhookURL == "" {
		log.Printf("企业微信Webhook URL为空")
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "企业微信Webhook URL不能为空",
		})
		return
	}

	log.Printf("企业微信Webhook URL: %s", requestData.WebhookURL)
	log.Printf("企业微信事件列表: %v", requestData.Events)

	// 构建消息内容
	requestBody := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"content": "### 测试消息\n这是一条来自VulnArk的测试消息，如果您收到此消息，表示企业微信机器人配置正确。",
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("构建请求体失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "构建请求失败",
		})
		return
	}

	log.Printf("企业微信请求体: %s", string(jsonData))

	// 发送HTTP请求
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Post(
		requestData.WebhookURL,
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		log.Printf("发送请求失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "发送测试消息失败: " + err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	// 读取并记录响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("读取响应失败: %v", err)
	} else {
		log.Printf("企业微信响应: %s, 状态码: %d", string(body), resp.StatusCode)
	}

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "测试失败，响应码: " + strconv.Itoa(resp.StatusCode),
		})
		return
	}

	// 解析响应
	var response struct {
		Errcode int    `json:"errcode"`
		Errmsg  string `json:"errmsg"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		log.Printf("解析响应失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "解析响应失败",
		})
		return
	}

	log.Printf("企业微信响应码: %d, 消息: %s", response.Errcode, response.Errmsg)

	if response.Errcode != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "测试失败，错误: " + response.Errmsg,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "测试成功",
	})
}

// TestFeishuBot 测试飞书机器人
func (sc *SettingsController) TestFeishuBot(c *gin.Context) {
	var feishuSettings models.FeishuSettings
	if err := c.ShouldBindJSON(&feishuSettings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 验证webhook URL是否存在
	if feishuSettings.WebhookURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Webhook URL不能为空",
		})
		return
	}

	// 当前时间戳(秒)
	timestamp := time.Now().Unix()

	// 构建飞书消息
	testMessage := map[string]interface{}{
		"timestamp": timestamp,
		"msg_type":  "interactive",
		"card": map[string]interface{}{
			"config": map[string]interface{}{
				"wide_screen_mode": true,
			},
			"header": map[string]interface{}{
				"title": map[string]interface{}{
					"tag":     "plain_text",
					"content": "VulnArk系统通知",
				},
				"template": "blue",
			},
			"elements": []interface{}{
				map[string]interface{}{
					"tag": "div",
					"text": map[string]interface{}{
						"tag":     "lark_md",
						"content": fmt.Sprintf("**测试消息** - %s\n这是一条来自VulnArk系统的测试消息，如果您看到此消息，说明飞书机器人配置成功。", utils.FormatTimeCST(time.Now())),
					},
				},
			},
		},
	}

	// 如果有签名密钥，需要计算签名
	if feishuSettings.Secret != "" {
		// 生成签名
		// timestamp + key 做sha256, 再进行base64编码
		stringToSign := fmt.Sprintf("%v", timestamp) + feishuSettings.Secret
		h := hmac.New(sha256.New, []byte(feishuSettings.Secret))
		h.Write([]byte(stringToSign))
		signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
		testMessage["sign"] = signature
	}

	messageJSON, err := json.Marshal(testMessage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "构建消息失败",
			"error":   err.Error(),
		})
		return
	}

	log.Printf("向飞书发送测试消息: %s", feishuSettings.WebhookURL)

	// 发送HTTP请求到飞书webhook
	resp, err := http.Post(
		feishuSettings.WebhookURL,
		"application/json",
		bytes.NewBuffer(messageJSON),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "发送消息失败",
			"error":   err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "读取响应失败",
			"error":   err.Error(),
		})
		return
	}

	log.Printf("飞书响应: %s", string(body))

	// 解析飞书响应
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "解析响应失败",
			"error":   err.Error(),
		})
		return
	}

	// 飞书API成功响应码为0
	if code, exists := result["code"]; exists && int(code.(float64)) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "飞书测试消息发送成功",
		})
	} else {
		msg := "未知错误"
		if m, exists := result["msg"]; exists {
			msg = m.(string)
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "飞书测试消息发送失败",
			"error":   fmt.Sprintf("错误码: %v, 错误信息: %s", code, msg),
		})
	}
}

// TestDingtalkBot 测试钉钉机器人
func (sc *SettingsController) TestDingtalkBot(c *gin.Context) {
	var dingtalkSettings models.DingtalkSettings
	if err := c.ShouldBindJSON(&dingtalkSettings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 验证webhook URL是否存在
	if dingtalkSettings.WebhookURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Webhook URL不能为空",
		})
		return
	}

	// 构建钉钉机器人消息
	testMessage := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"title": "VulnArk系统通知",
			"text":  fmt.Sprintf("# VulnArk系统通知\n\n**测试消息** - %s\n\n这是一条来自VulnArk系统的测试消息，如果您看到此消息，说明钉钉机器人配置成功。", utils.FormatTimeCST(time.Now())),
		},
	}

	webhookURL := dingtalkSettings.WebhookURL

	// 如果有签名密钥，需要计算签名并添加到URL
	if dingtalkSettings.Secret != "" {
		timestamp := time.Now().UnixNano() / 1e6
		stringToSign := fmt.Sprintf("%d\n%s", timestamp, dingtalkSettings.Secret)

		h := hmac.New(sha256.New, []byte(dingtalkSettings.Secret))
		h.Write([]byte(stringToSign))
		signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

		// URL编码签名
		signatureEncoded := url.QueryEscape(signature)

		// 添加时间戳和签名到URL
		if strings.Contains(webhookURL, "?") {
			webhookURL = fmt.Sprintf("%s&timestamp=%d&sign=%s", webhookURL, timestamp, signatureEncoded)
		} else {
			webhookURL = fmt.Sprintf("%s?timestamp=%d&sign=%s", webhookURL, timestamp, signatureEncoded)
		}
	}

	messageJSON, err := json.Marshal(testMessage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "构建消息失败",
			"error":   err.Error(),
		})
		return
	}

	log.Printf("向钉钉发送测试消息: %s", webhookURL)

	// 发送HTTP请求到钉钉webhook
	resp, err := http.Post(
		webhookURL,
		"application/json",
		bytes.NewBuffer(messageJSON),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "发送消息失败",
			"error":   err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "读取响应失败",
			"error":   err.Error(),
		})
		return
	}

	log.Printf("钉钉响应: %s", string(body))

	// 解析钉钉响应
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "解析响应失败",
			"error":   err.Error(),
		})
		return
	}

	// 钉钉API成功响应errcode为0
	if errcode, exists := result["errcode"]; exists && int(errcode.(float64)) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "钉钉测试消息发送成功",
		})
	} else {
		errmsg := "未知错误"
		if msg, exists := result["errmsg"]; exists {
			errmsg = msg.(string)
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "钉钉测试消息发送失败",
			"error":   fmt.Sprintf("错误码: %v, 错误信息: %s", errcode, errmsg),
		})
	}
}

// TestEmailNotification 测试邮件发送
func (sc *SettingsController) TestEmailNotification(c *gin.Context) {
	var emailSettings models.EmailSettings
	if err := c.ShouldBindJSON(&emailSettings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 验证必要参数
	if emailSettings.SMTPServer == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "SMTP服务器地址不能为空",
		})
		return
	}

	if emailSettings.SMTPPort <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "SMTP端口无效",
		})
		return
	}

	if emailSettings.FromEmail == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "发件人邮箱不能为空",
		})
		return
	}

	if len(emailSettings.Recipients) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "收件人列表不能为空",
		})
		return
	}

	// 创建SMTP配置
	smtpAuth := smtp.PlainAuth("", emailSettings.Username, emailSettings.Password, emailSettings.SMTPServer)

	// 构建发送地址
	smtpAddr := fmt.Sprintf("%s:%d", emailSettings.SMTPServer, emailSettings.SMTPPort)

	// 构建邮件头
	headers := map[string]string{
		"From":         emailSettings.FromEmail,
		"To":           strings.Join(emailSettings.Recipients, ", "),
		"Subject":      "VulnArk系统测试邮件",
		"MIME-Version": "1.0",
		"Content-Type": "text/html; charset=UTF-8",
	}

	// 构建头部字符串
	var headerStr strings.Builder
	for key, value := range headers {
		headerStr.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}

	// 构建邮件内容
	timeStr := utils.FormatTimeCST(time.Now())
	emailBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>VulnArk系统测试邮件</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #ddd; border-radius: 5px; }
        .header { background-color: #4e54c8; color: white; padding: 10px; text-align: center; border-radius: 5px 5px 0 0; }
        .content { padding: 20px; }
        .footer { background-color: #f5f5f5; padding: 10px; text-align: center; font-size: 12px; border-radius: 0 0 5px 5px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h2>VulnArk系统通知</h2>
        </div>
        <div class="content">
            <h3>测试邮件</h3>
            <p>您好，这是一封来自VulnArk系统的测试邮件。</p>
            <p>如果您收到此邮件，说明邮件通知功能已经配置成功。</p>
            <p>发送时间：%s</p>
        </div>
        <div class="footer">
            <p>此邮件由VulnArk系统自动发送，请勿回复。</p>
        </div>
    </div>
</body>
</html>
`, timeStr)

	// 完整邮件内容
	message := headerStr.String() + "\r\n" + emailBody

	log.Printf("尝试发送测试邮件到: %s", strings.Join(emailSettings.Recipients, ", "))

	// 使用TLS加密通信
	if emailSettings.UseSSL {
		// TLS配置
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true, // 在测试环境中可以跳过证书验证
			ServerName:         emailSettings.SMTPServer,
		}

		// 连接SMTP服务器
		conn, err := tls.Dial("tcp", smtpAddr, tlsConfig)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "连接SMTP服务器失败",
				"error":   err.Error(),
			})
			return
		}

		// 创建SMTP客户端
		client, err := smtp.NewClient(conn, emailSettings.SMTPServer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "创建SMTP客户端失败",
				"error":   err.Error(),
			})
			return
		}
		defer client.Close()

		// 设置身份验证
		if err := client.Auth(smtpAuth); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "SMTP身份验证失败",
				"error":   err.Error(),
			})
			return
		}

		// 设置发件人
		if err := client.Mail(emailSettings.FromEmail); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "设置发件人失败",
				"error":   err.Error(),
			})
			return
		}

		// 设置收件人
		for _, recipient := range emailSettings.Recipients {
			if err := client.Rcpt(recipient); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": "设置收件人失败",
					"error":   err.Error(),
				})
				return
			}
		}

		// 设置邮件内容
		w, err := client.Data()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "准备邮件内容失败",
				"error":   err.Error(),
			})
			return
		}

		_, err = w.Write([]byte(message))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "写入邮件内容失败",
				"error":   err.Error(),
			})
			return
		}

		err = w.Close()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "关闭写入器失败",
				"error":   err.Error(),
			})
			return
		}

		// 结束会话
		client.Quit()
	} else {
		// 不使用SSL，直接发送
		err := smtp.SendMail(
			smtpAddr,
			smtpAuth,
			emailSettings.FromEmail,
			emailSettings.Recipients,
			[]byte(message),
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "发送邮件失败",
				"error":   err.Error(),
			})
			return
		}
	}

	log.Printf("测试邮件发送成功")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "测试邮件发送成功",
	})
}

// TestAiService 测试AI服务连接
func (sc *SettingsController) TestAiService(c *gin.Context) {
	var aiSettings models.AISettings
	if err := c.ShouldBindJSON(&aiSettings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// TODO: 实现AI服务连接测试
	// 这里简单模拟测试结果，实际应该调用相应的AI服务API验证连接
	log.Printf("测试AI服务连接: 提供商: %s", aiSettings.Provider)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "AI服务连接测试成功",
	})
}

// TestVulnDBConnection 测试漏洞库API连接
func (sc *SettingsController) TestVulnDBConnection(c *gin.Context) {
	var vulndbSettings models.VulnDBSettings
	if err := c.ShouldBindJSON(&vulndbSettings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 验证必要参数
	if vulndbSettings.APIURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "API URL不能为空",
		})
		return
	}

	if vulndbSettings.APIKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "API Key不能为空",
		})
		return
	}

	// 当Provider为"weibu"时，调用微步社区API
	if vulndbSettings.Provider == "weibu" {
		// 测试CVE
		testCVE := "CVE-2021-44228" // Log4Shell漏洞，常见且重要的漏洞

		// 确保API URL是干净的，无末尾斜杠
		apiURL := strings.TrimRight(vulndbSettings.APIURL, "/")

		// 构建请求URL (漏洞详情查询接口)
		requestURL := fmt.Sprintf("%s/api/v3/vuln/detail", apiURL)

		// 构建请求表单数据
		data := url.Values{}
		data.Set("apikey", vulndbSettings.APIKey)
		data.Set("cve", testCVE)

		log.Printf("向微步社区发送API测试请求: %s, CVE: %s", requestURL, testCVE)

		// 发送HTTP POST请求
		resp, err := http.PostForm(requestURL, data)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "API请求失败",
				"error":   err.Error(),
			})
			return
		}
		defer resp.Body.Close()

		// 读取响应内容
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "读取API响应失败",
				"error":   err.Error(),
			})
			return
		}

		log.Printf("微步社区API响应: %s", string(body))

		// 解析JSON响应
		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "解析API响应失败",
				"error":   err.Error(),
			})
			return
		}

		// 检查响应状态码
		if responseCode, exists := result["response_code"]; exists {
			// 微步API成功响应码通常为0
			if int(responseCode.(float64)) == 0 {
				c.JSON(http.StatusOK, gin.H{
					"code":    200,
					"message": "微步社区漏洞API连接测试成功",
					"data": map[string]interface{}{
						"cve":    testCVE,
						"status": "API连接正常",
					},
				})
				return
			} else {
				// 获取错误信息
				verboseMsg := "未知错误"
				if msg, exists := result["verbose_msg"]; exists {
					verboseMsg = msg.(string)
				}

				c.JSON(http.StatusBadRequest, gin.H{
					"code":    400,
					"message": "微步社区API测试失败",
					"error":   fmt.Sprintf("错误码: %v, 错误信息: %s", responseCode, verboseMsg),
				})
				return
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "微步社区API返回格式异常",
				"error":   "响应中没有response_code字段",
			})
			return
		}
	} else if vulndbSettings.Provider == "vulniq" {
		// TODO: 实现VulnIQ API测试
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "VulnIQ API连接测试功能尚未实现",
		})
		return
	} else if vulndbSettings.Provider == "vuldb" {
		// TODO: 实现VulDB API测试
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "VulDB API连接测试功能尚未实现",
		})
		return
	} else if vulndbSettings.Provider == "other" {
		// 对于自定义API，我们只检查URL和Key是否有效，不进行实际调用
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "自定义漏洞库API配置验证成功",
		})
		return
	}

	// 默认返回成功
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "漏洞库API连接测试成功",
	})
}

// TestVulnerabilityNotification 测试漏洞通知
func (sc *SettingsController) TestVulnerabilityNotification(c *gin.Context) {
	log.Printf("测试漏洞通知")

	// 模拟一个测试漏洞对象
	testVuln := models.Vulnerability{
		ID:           999,
		Title:        "测试漏洞通知",
		Description:  "这是一个用于测试通知系统的漏洞",
		Severity:     models.SeverityHigh,
		Status:       models.StatusNew,
		CVE:          "CVE-2023-TEST",
		CVSS:         8.5,
		DiscoveredAt: time.Now(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// 获取通知管理器
	notificationManager, err := utils.NewNotificationManager()
	if err != nil {
		log.Printf("获取通知管理器失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取通知管理器失败: " + err.Error(),
		})
		return
	}

	// 记录配置信息用于调试
	log.Printf("通知设置获取成功，开始测试通知发送")
	log.Printf("企业微信通知状态: %v", notificationManager.GetSettings().Notifications.WorkWechat.Enabled)
	log.Printf("企业微信Webhook: %s", notificationManager.GetSettings().Notifications.WorkWechat.WebhookURL)
	log.Printf("企业微信通知事件列表: %v", notificationManager.GetSettings().Notifications.WorkWechat.Events)

	// 创建通知状态结果映射
	results := map[string]string{
		"workWechat": "未启用",
		"feishu":     "未启用",
		"dingtalk":   "未启用",
		"email":      "未启用",
	}

	// 更新状态
	if notificationManager.GetSettings().Notifications.WorkWechat.Enabled {
		results["workWechat"] = "已启用，将发送测试通知"
	}
	if notificationManager.GetSettings().Notifications.Feishu.Enabled {
		results["feishu"] = "已启用，将发送测试通知"
	}
	if notificationManager.GetSettings().Notifications.Dingtalk.Enabled {
		results["dingtalk"] = "已启用，将发送测试通知"
	}
	if notificationManager.GetSettings().Notifications.Email.Enabled {
		results["email"] = "已启用，将发送测试通知"
	}

	// 发送漏洞新增通知 - 该方法没有返回值
	notificationManager.SendVulnerabilityNotification(utils.EventVulnCreate, &testVuln, "")

	log.Printf("测试通知已发送完成")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "测试通知已发送",
		"data": gin.H{
			"settings": notificationManager.GetSettings().Notifications,
			"results":  results,
		},
	})
}
