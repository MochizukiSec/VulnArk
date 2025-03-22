package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vulnark/vulnark/models"
	"github.com/vulnark/vulnark/utils"
)

// IntegrationController 处理与CI/CD集成相关的接口
type IntegrationController struct{}

// ReceiveScanResult 接收CI/CD管道中的扫描结果
func (i *IntegrationController) ReceiveScanResult(c *gin.Context) {
	// 获取集成类型
	integrationType := c.Param("type")

	// 获取API密钥
	apiKey := c.GetHeader("X-API-Key")
	if apiKey == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未提供API密钥",
		})
		return
	}

	// 验证API密钥
	var integration models.CIIntegration
	if err := utils.DB.Where("api_key = ? AND type = ? AND enabled = ?", apiKey, integrationType, true).First(&integration).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "无效的API密钥或集成类型",
		})
		return
	}

	// 读取请求体
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "读取请求体失败: " + err.Error(),
		})
		return
	}

	// 根据集成类型处理不同格式的数据
	var vulnerabilities []models.Vulnerability
	var processErr error

	switch integrationType {
	case "jenkins":
		vulnerabilities, processErr = processJenkinsResult(body, integration)
	case "gitlab":
		vulnerabilities, processErr = processGitlabResult(body, integration)
	case "github":
		vulnerabilities, processErr = processGithubResult(body, integration)
	case "custom":
		vulnerabilities, processErr = processCustomResult(body, integration)
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "不支持的集成类型: " + integrationType,
		})
		return
	}

	if processErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "处理扫描结果失败: " + processErr.Error(),
		})
		return
	}

	// 记录扫描结果
	successCount := 0
	errorCount := 0

	for _, vuln := range vulnerabilities {
		// 检查是否已存在相同漏洞
		var existingVuln models.Vulnerability
		if err := utils.DB.Where("cve = ? OR title = ?", vuln.CVE, vuln.Title).First(&existingVuln).Error; err == nil {
			// 更新现有漏洞
			existingVuln.UpdatedAt = time.Now()
			existingVuln.Status = vuln.Status
			existingVuln.Description = vuln.Description
			existingVuln.Severity = vuln.Severity
			existingVuln.References = vuln.References

			if err := utils.DB.Save(&existingVuln).Error; err != nil {
				log.Printf("更新漏洞失败: %v", err)
				errorCount++
				continue
			}

			successCount++
		} else {
			// 创建新漏洞
			vuln.CreatedAt = time.Now()
			vuln.UpdatedAt = time.Now()

			if err := utils.DB.Create(&vuln).Error; err != nil {
				log.Printf("创建漏洞失败: %v", err)
				errorCount++
				continue
			}

			successCount++
		}
	}

	// 记录集成历史
	history := models.IntegrationHistory{
		IntegrationID:   integration.ID,
		IntegrationType: integrationType,
		Status:          "success",
		TotalRecords:    len(vulnerabilities),
		SuccessCount:    successCount,
		ErrorCount:      errorCount,
		ExecutedAt:      time.Now(),
	}

	if err := utils.DB.Create(&history).Error; err != nil {
		log.Printf("记录集成历史失败: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "扫描结果处理成功",
		"data": gin.H{
			"total":         len(vulnerabilities),
			"success_count": successCount,
			"error_count":   errorCount,
		},
	})
}

// 处理来自Jenkins的扫描结果
func processJenkinsResult(data []byte, integration models.CIIntegration) ([]models.Vulnerability, error) {
	var result struct {
		Findings []struct {
			Title       string `json:"title"`
			Severity    string `json:"severity"`
			Description string `json:"description"`
			CVE         string `json:"cve_id"`
			References  string `json:"references"`
		} `json:"findings"`
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	var vulnerabilities []models.Vulnerability

	for _, finding := range result.Findings {
		vuln := models.Vulnerability{
			Title:       finding.Title,
			Severity:    models.Severity(finding.Severity),
			Description: finding.Description,
			CVE:         finding.CVE,
			References:  finding.References,
			Status:      models.StatusNew,
			Source:      "jenkins-ci",
		}
		vulnerabilities = append(vulnerabilities, vuln)
	}

	return vulnerabilities, nil
}

// 处理来自GitLab CI的扫描结果
func processGitlabResult(data []byte, integration models.CIIntegration) ([]models.Vulnerability, error) {
	var result struct {
		Vulnerabilities []struct {
			Name        string `json:"name"`
			Severity    string `json:"severity"`
			Description string `json:"description"`
			CVE         string `json:"cve"`
			Solution    string `json:"solution"`
			Location    struct {
				File  string `json:"file"`
				Start int    `json:"start_line"`
				End   int    `json:"end_line"`
			} `json:"location"`
		} `json:"vulnerabilities"`
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	var vulnerabilities []models.Vulnerability

	for _, vuln := range result.Vulnerabilities {
		locationInfo := ""
		if vuln.Location.File != "" {
			locationInfo = "文件: " + vuln.Location.File
			if vuln.Location.Start > 0 {
				locationInfo += ", 行: " + string(vuln.Location.Start)
			}
		}

		v := models.Vulnerability{
			Title:       vuln.Name,
			Severity:    models.Severity(vuln.Severity),
			Description: vuln.Description,
			CVE:         vuln.CVE,
			Solution:    vuln.Solution,
			Status:      models.StatusNew,
			Source:      "gitlab-ci",
			References:  fmt.Sprintf("File: %s, Lines: %d-%d", vuln.Location.File, vuln.Location.Start, vuln.Location.End),
		}
		vulnerabilities = append(vulnerabilities, v)
	}

	return vulnerabilities, nil
}

// 处理来自GitHub Actions的扫描结果
func processGithubResult(data []byte, integration models.CIIntegration) ([]models.Vulnerability, error) {
	var result struct {
		Results []struct {
			RuleID      string `json:"rule_id"`
			RuleName    string `json:"rule_name"`
			Severity    string `json:"severity"`
			Description string `json:"description"`
			Path        string `json:"path"`
			StartLine   int    `json:"start_line"`
			EndLine     int    `json:"end_line"`
		} `json:"results"`
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	var vulnerabilities []models.Vulnerability

	for _, finding := range result.Results {
		locationInfo := ""
		if finding.Path != "" {
			locationInfo = "文件: " + finding.Path
			if finding.StartLine > 0 {
				locationInfo += ", 行: " + string(finding.StartLine)
			}
		}

		vuln := models.Vulnerability{
			Title:       finding.RuleName,
			Severity:    models.Severity(finding.Severity),
			Description: finding.Description,
			CVE:         finding.RuleID, // 使用规则ID作为CVE字段
			Status:      models.StatusNew,
			Source:      "github-action",
			References:  locationInfo,
		}
		vulnerabilities = append(vulnerabilities, vuln)
	}

	return vulnerabilities, nil
}

// 处理自定义格式的扫描结果
func processCustomResult(data []byte, integration models.CIIntegration) ([]models.Vulnerability, error) {
	var result struct {
		Vulnerabilities []struct {
			Title       string `json:"title"`
			Severity    string `json:"severity"`
			Description string `json:"description"`
			Identifier  string `json:"identifier"`
			Reference   string `json:"reference"`
			Status      string `json:"status"`
			Source      string `json:"source"`
		} `json:"vulnerabilities"`
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	var vulnerabilities []models.Vulnerability

	for _, v := range result.Vulnerabilities {
		status := models.StatusNew
		if v.Status != "" {
			status = models.VulnStatus(v.Status)
		}

		vuln := models.Vulnerability{
			Title:       v.Title,
			Severity:    models.Severity(v.Severity),
			Description: v.Description,
			CVE:         v.Identifier,
			Status:      status,
			Source:      "custom-integration",
			References:  v.Reference,
		}
		vulnerabilities = append(vulnerabilities, vuln)
	}

	return vulnerabilities, nil
}

// GetIntegrations 获取所有CI/CD集成配置
func (i *IntegrationController) GetIntegrations(c *gin.Context) {
	var integrations []models.CIIntegration

	if err := utils.DB.Find(&integrations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取集成配置失败: " + err.Error(),
		})
		return
	}

	// 隐藏敏感信息
	for i := range integrations {
		integrations[i].APIKey = "[已隐藏]"
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取集成配置成功",
		"data":    integrations,
	})
}

// CreateIntegration 创建新的CI/CD集成配置
func (i *IntegrationController) CreateIntegration(c *gin.Context) {
	var integration models.CIIntegration

	if err := c.ShouldBindJSON(&integration); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的请求参数: " + err.Error(),
		})
		return
	}

	// 生成API密钥
	integration.APIKey = utils.GenerateRandomString(32)
	integration.CreatedAt = time.Now()
	integration.UpdatedAt = time.Now()

	if err := utils.DB.Create(&integration).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建集成配置失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建集成配置成功",
		"data":    integration,
	})
}

// UpdateIntegration 更新CI/CD集成配置
func (i *IntegrationController) UpdateIntegration(c *gin.Context) {
	id := c.Param("id")

	var integration models.CIIntegration
	if err := utils.DB.First(&integration, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "集成配置不存在",
		})
		return
	}

	var updateData struct {
		Name        string `json:"name"`
		Type        string `json:"type"`
		Description string `json:"description"`
		Enabled     bool   `json:"enabled"`
		Config      string `json:"config"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的请求参数: " + err.Error(),
		})
		return
	}

	integration.Name = updateData.Name
	integration.Type = updateData.Type
	integration.Description = updateData.Description
	integration.Enabled = updateData.Enabled
	integration.Config = updateData.Config
	integration.UpdatedAt = time.Now()

	if err := utils.DB.Save(&integration).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新集成配置失败: " + err.Error(),
		})
		return
	}

	// 隐藏API密钥
	integration.APIKey = "[已隐藏]"

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新集成配置成功",
		"data":    integration,
	})
}

// DeleteIntegration 删除CI/CD集成配置
func (i *IntegrationController) DeleteIntegration(c *gin.Context) {
	id := c.Param("id")

	var integration models.CIIntegration
	if err := utils.DB.First(&integration, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "集成配置不存在",
		})
		return
	}

	if err := utils.DB.Delete(&integration).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除集成配置失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除集成配置成功",
	})
}

// RegenerateAPIKey 重新生成API密钥
func (i *IntegrationController) RegenerateAPIKey(c *gin.Context) {
	id := c.Param("id")

	var integration models.CIIntegration
	if err := utils.DB.First(&integration, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "集成配置不存在",
		})
		return
	}

	// 生成新的API密钥
	integration.APIKey = utils.GenerateRandomString(32)
	integration.UpdatedAt = time.Now()

	if err := utils.DB.Save(&integration).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "重新生成API密钥失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "重新生成API密钥成功",
		"data": gin.H{
			"api_key": integration.APIKey,
		},
	})
}

// GetIntegrationHistory 获取集成历史记录
func (i *IntegrationController) GetIntegrationHistory(c *gin.Context) {
	id := c.Param("id")

	var histories []models.IntegrationHistory
	if err := utils.DB.Where("integration_id = ?", id).Order("executed_at DESC").Find(&histories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取集成历史记录失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取集成历史记录成功",
		"data":    histories,
	})
}

// UpdateIntegrationStatus 更新CI/CD集成启用状态
func (i *IntegrationController) UpdateIntegrationStatus(c *gin.Context) {
	id := c.Param("id")

	var integration models.CIIntegration
	if err := utils.DB.First(&integration, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "集成配置不存在",
		})
		return
	}

	var statusData struct {
		Enabled bool `json:"enabled"`
	}

	if err := c.ShouldBindJSON(&statusData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的请求参数: " + err.Error(),
		})
		return
	}

	// 只更新启用状态和更新时间
	updateFields := map[string]interface{}{
		"enabled":    statusData.Enabled,
		"updated_at": time.Now(),
	}

	if err := utils.DB.Model(&integration).Updates(updateFields).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新集成状态失败: " + err.Error(),
		})
		return
	}

	// 创建状态变更历史记录
	statusText := "禁用"
	if statusData.Enabled {
		statusText = "启用"
	}

	history := models.IntegrationHistory{
		IntegrationID:   integration.ID,
		IntegrationType: integration.Type,
		Status:          "success",
		Message:         "手动" + statusText + "集成",
		ExecutedAt:      time.Now(),
	}

	utils.DB.Create(&history)

	// 返回更新后的集成配置（隐藏API密钥）
	integration.Enabled = statusData.Enabled
	integration.APIKey = "[已隐藏]"

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": statusText + "集成成功",
		"data":    integration,
	})
}
