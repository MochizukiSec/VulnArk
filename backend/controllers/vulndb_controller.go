package controllers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vulnark/vulnark/models"
	"github.com/vulnark/vulnark/utils"
)

// VulnDBController 漏洞库控制器
type VulnDBController struct{}

// ListVulnDBEntries 获取漏洞库条目列表
func (v *VulnDBController) ListVulnDBEntries(c *gin.Context) {
	// 获取查询参数
	keyword := c.Query("keyword")
	cve := c.Query("cve")
	cwe := c.Query("cwe")
	severity := c.Query("severity")
	hasExploit := c.Query("has_exploit")
	tags := c.Query("tags")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// 构建查询
	query := utils.DB.Model(&models.VulnDB{})

	// 条件过滤
	if keyword != "" {
		query = query.Where("title LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if cve != "" {
		query = query.Where("cve LIKE ?", "%"+cve+"%")
	}
	if cwe != "" {
		query = query.Where("cwe LIKE ?", "%"+cwe+"%")
	}
	if severity != "" {
		query = query.Where("severity = ?", severity)
	}
	if hasExploit == "true" {
		query = query.Where("exploit_available = ?", true)
	}
	if tags != "" {
		query = query.Where("tags LIKE ?", "%"+tags+"%")
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 获取分页数据
	var vulnDBEntries []models.VulnDB
	query.Limit(pageSize).Offset((page - 1) * pageSize).Order("published_date DESC").Find(&vulnDBEntries)

	// 转换为响应格式
	var responseItems []map[string]interface{}
	for _, vuln := range vulnDBEntries {
		item := map[string]interface{}{
			"id":                vuln.ID,
			"cve":               vuln.CVE,
			"cwe":               vuln.CWE,
			"title":             vuln.Title,
			"severity":          string(vuln.Severity),
			"cvss":              vuln.CVSS,
			"affected_systems":  vuln.AffectedSystems,
			"exploit_available": vuln.ExploitAvailable,
			"tags":              vuln.Tags,
			"published_date":    vuln.PublishedDate,
		}
		responseItems = append(responseItems, item)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"items": responseItems,
			"total": total,
		},
	})
}

// GetVulnDBByID 获取单个漏洞库条目详情
func (v *VulnDBController) GetVulnDBByID(c *gin.Context) {
	id := c.Param("id")
	log.Printf("获取漏洞库条目详情, ID: %s", id)

	var vulnDBEntry models.VulnDB
	if err := utils.DB.First(&vulnDBEntry, id).Error; err != nil {
		log.Printf("漏洞库条目不存在, ID: %s, 错误: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "漏洞库条目不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": vulnDBEntry,
	})
}

// GetVulnDBByCVE 根据CVE获取漏洞库条目
func (v *VulnDBController) GetVulnDBByCVE(c *gin.Context) {
	cve := c.Param("cve")
	log.Printf("根据CVE获取漏洞库条目, CVE: %s", cve)

	var vulnDBEntry models.VulnDB
	if err := utils.DB.Where("cve = ?", cve).First(&vulnDBEntry).Error; err != nil {
		log.Printf("漏洞库条目不存在, CVE: %s, 错误: %v", cve, err)
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "漏洞库条目不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": vulnDBEntry,
	})
}

// CreateVulnDBEntry 创建漏洞库条目
func (v *VulnDBController) CreateVulnDBEntry(c *gin.Context) {
	// 获取用户ID，确保用户已登录
	userID, exists := c.Get("userID")
	if !exists {
		log.Printf("用户ID未找到")
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "用户ID未找到",
		})
		return
	}
	log.Printf("用户 ID: %v 尝试创建漏洞库条目", userID)

	// 接收请求数据
	var requestData struct {
		CVE              string   `json:"cve"`
		CWE              string   `json:"cwe"`
		Title            string   `json:"title"`
		Description      string   `json:"description"`
		Severity         string   `json:"severity"`
		CVSS             string   `json:"cvss"`
		CVSSVector       string   `json:"cvss_vector"` // 添加CVSS向量字段
		AffectedSystems  string   `json:"affected_systems"`
		AffectedVersions string   `json:"affected_versions"`
		AffectedProducts []string `json:"affected_products"` // 添加受影响产品数组
		Solution         string   `json:"solution"`
		Remediation      string   `json:"remediation"` // 添加修复建议字段
		References       string   `json:"references"`
		ReferencesArray  []string `json:"references_array"` // 添加参考链接数组
		Tags             string   `json:"tags"`
		TagsArray        []string `json:"tags_array"` // 添加标签数组
		ExploitAvailable bool     `json:"exploit_available"`
		ExploitInfo      string   `json:"exploit_info"` // 添加漏洞利用信息
		PublishedDate    string   `json:"published_date"`
		UpdatedDate      string   `json:"updated_date"`
		LastModifiedDate string   `json:"last_modified_date"` // 添加最后修改日期
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		log.Printf("创建漏洞库条目请求参数错误: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	log.Printf("接收到的创建漏洞库条目数据: %+v", requestData)

	// 验证必填字段
	if requestData.Title == "" || requestData.Description == "" || requestData.Severity == "" {
		log.Printf("缺少必要的字段: title=%s, description=%s, severity=%s",
			requestData.Title, requestData.Description, requestData.Severity)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "缺少必要的字段: 标题、描述或严重程度",
		})
		return
	}

	// 处理数组字段
	if len(requestData.AffectedProducts) > 0 && requestData.AffectedSystems == "" {
		requestData.AffectedSystems = strings.Join(requestData.AffectedProducts, ",")
	}

	if len(requestData.ReferencesArray) > 0 && requestData.References == "" {
		requestData.References = strings.Join(requestData.ReferencesArray, ",")
	}

	if len(requestData.TagsArray) > 0 && requestData.Tags == "" {
		requestData.Tags = strings.Join(requestData.TagsArray, ",")
	}

	// 处理修复建议字段，如果solution为空但remediation有值，则使用remediation
	if requestData.Solution == "" && requestData.Remediation != "" {
		requestData.Solution = requestData.Remediation
	}

	// 处理日期字段
	if requestData.UpdatedDate == "" && requestData.LastModifiedDate != "" {
		requestData.UpdatedDate = requestData.LastModifiedDate
	}

	// 检查CVE是否已存在
	if requestData.CVE != "" {
		var count int64
		utils.DB.Model(&models.VulnDB{}).Where("cve = ?", requestData.CVE).Count(&count)
		if count > 0 {
			log.Printf("CVE已存在: %s", requestData.CVE)
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "该CVE已存在于漏洞库中",
			})
			return
		}
	}

	// 处理日期
	var publishedDate time.Time
	var err error
	if requestData.PublishedDate != "" {
		publishedDate, err = time.Parse("2006-01-02", requestData.PublishedDate)
		if err != nil {
			log.Printf("解析发布日期失败: %v", err)
			publishedDate = time.Now() // 默认为当前时间
		}
	} else {
		publishedDate = time.Now()
	}

	// 解析CVSS字符串为float64
	var cvssValue float64
	if requestData.CVSS != "" {
		var err error
		cvssValue, err = strconv.ParseFloat(requestData.CVSS, 64)
		if err != nil {
			log.Printf("解析CVSS值失败: %v", err)
			cvssValue = 0.0 // 设置默认值
		}
	}

	// 创建漏洞库条目
	now := time.Now()
	vulnDBEntry := models.VulnDB{
		CVE:              requestData.CVE,
		CWE:              requestData.CWE,
		Title:            requestData.Title,
		Description:      requestData.Description,
		Severity:         models.Severity(requestData.Severity),
		CVSS:             cvssValue, // 使用转换后的CVSS值
		AffectedSystems:  requestData.AffectedSystems,
		AffectedVersions: requestData.AffectedVersions,
		Solution:         requestData.Solution,
		References:       requestData.References,
		Tags:             requestData.Tags,
		ExploitAvailable: requestData.ExploitAvailable,
		PublishedDate:    publishedDate,
		UpdatedDate:      now,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	// 保存到数据库
	if err := utils.DB.Create(&vulnDBEntry).Error; err != nil {
		log.Printf("创建漏洞库条目失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建漏洞库条目失败",
			"error":   err.Error(),
		})
		return
	}

	log.Printf("漏洞库条目创建成功, ID: %d, 标题: %s", vulnDBEntry.ID, vulnDBEntry.Title)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建漏洞库条目成功",
		"data": gin.H{
			"id": vulnDBEntry.ID,
		},
	})
}

// UpdateVulnDBEntry 更新漏洞库条目
func (v *VulnDBController) UpdateVulnDBEntry(c *gin.Context) {
	// 获取用户ID，确保用户已登录
	userID, exists := c.Get("userID")
	if !exists {
		log.Printf("用户ID未找到")
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "用户ID未找到",
		})
		return
	}

	// 获取漏洞ID
	vulnDBID := c.Param("id")
	if vulnDBID == "" {
		log.Printf("漏洞库条目ID为空")
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "漏洞库条目ID不能为空",
		})
		return
	}

	log.Printf("用户 ID: %v 尝试更新漏洞库条目 ID: %s", userID, vulnDBID)

	// 接收请求数据
	var requestData struct {
		CVE              string   `json:"cve"`
		CWE              string   `json:"cwe"`
		Title            string   `json:"title"`
		Description      string   `json:"description"`
		Severity         string   `json:"severity"`
		CVSS             string   `json:"cvss"`
		CVSSVector       string   `json:"cvss_vector"` // 添加CVSS向量字段
		AffectedSystems  string   `json:"affected_systems"`
		AffectedVersions string   `json:"affected_versions"`
		AffectedProducts []string `json:"affected_products"` // 添加受影响产品数组
		Solution         string   `json:"solution"`
		Remediation      string   `json:"remediation"` // 添加修复建议字段
		References       string   `json:"references"`
		ReferencesArray  []string `json:"references_array"` // 添加参考链接数组
		Tags             string   `json:"tags"`
		TagsArray        []string `json:"tags_array"` // 添加标签数组
		ExploitAvailable bool     `json:"exploit_available"`
		ExploitInfo      string   `json:"exploit_info"` // 添加漏洞利用信息
		PublishedDate    string   `json:"published_date"`
		UpdatedDate      string   `json:"updated_date"`
		LastModifiedDate string   `json:"last_modified_date"` // 添加最后修改日期
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		log.Printf("更新漏洞库条目请求参数错误: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	log.Printf("接收到的更新漏洞库条目数据: %+v", requestData)

	// 验证必填字段
	if requestData.Title == "" || requestData.Description == "" || requestData.Severity == "" {
		log.Printf("缺少必要的字段: title=%s, description=%s, severity=%s",
			requestData.Title, requestData.Description, requestData.Severity)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "缺少必要的字段: 标题、描述或严重程度",
		})
		return
	}

	// 处理数组字段
	if len(requestData.AffectedProducts) > 0 && requestData.AffectedSystems == "" {
		requestData.AffectedSystems = strings.Join(requestData.AffectedProducts, ",")
	}

	if len(requestData.ReferencesArray) > 0 && requestData.References == "" {
		requestData.References = strings.Join(requestData.ReferencesArray, ",")
	}

	if len(requestData.TagsArray) > 0 && requestData.Tags == "" {
		requestData.Tags = strings.Join(requestData.TagsArray, ",")
	}

	// 处理修复建议字段，如果solution为空但remediation有值，则使用remediation
	if requestData.Solution == "" && requestData.Remediation != "" {
		requestData.Solution = requestData.Remediation
	}

	// 处理日期字段
	if requestData.UpdatedDate == "" && requestData.LastModifiedDate != "" {
		requestData.UpdatedDate = requestData.LastModifiedDate
	}

	// 查找要更新的漏洞库条目
	var existingVulnDB models.VulnDB
	if err := utils.DB.First(&existingVulnDB, vulnDBID).Error; err != nil {
		log.Printf("找不到漏洞库条目 ID: %s, 错误: %v", vulnDBID, err)
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "找不到指定的漏洞库条目",
		})
		return
	}

	// 如果CVE变更了，检查新的CVE是否已存在
	if requestData.CVE != "" && requestData.CVE != existingVulnDB.CVE {
		var count int64
		utils.DB.Model(&models.VulnDB{}).Where("cve = ? AND id != ?", requestData.CVE, vulnDBID).Count(&count)
		if count > 0 {
			log.Printf("CVE已被其他条目使用: %s", requestData.CVE)
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "该CVE已被其他漏洞库条目使用",
			})
			return
		}
	}

	// 解析CVSS字符串为float64
	var cvssValue float64
	if requestData.CVSS != "" {
		var err error
		cvssValue, err = strconv.ParseFloat(requestData.CVSS, 64)
		if err != nil {
			log.Printf("解析CVSS值失败: %v", err)
			cvssValue = 0.0 // 设置默认值
		}
	}

	// 解析日期
	var publishedDate time.Time
	if requestData.PublishedDate != "" {
		var err error
		publishedDate, err = time.Parse("2006-01-02", requestData.PublishedDate)
		if err != nil {
			log.Printf("解析发布日期失败: %v", err)
			publishedDate = existingVulnDB.PublishedDate
		}
	} else {
		publishedDate = existingVulnDB.PublishedDate
	}

	var updatedDate time.Time
	if requestData.UpdatedDate != "" {
		var err error
		updatedDate, err = time.Parse("2006-01-02", requestData.UpdatedDate)
		if err != nil {
			log.Printf("解析更新日期失败: %v", err)
			updatedDate = time.Now()
		}
	} else {
		updatedDate = time.Now()
	}

	// 更新漏洞库条目
	updateData := map[string]interface{}{
		"UpdatedAt":   updatedDate,
		"UpdatedDate": updatedDate,
	}

	if requestData.CVE != "" {
		updateData["CVE"] = requestData.CVE
	}
	if requestData.CWE != "" {
		updateData["CWE"] = requestData.CWE
	}
	if requestData.Title != "" {
		updateData["Title"] = requestData.Title
	}
	if requestData.Description != "" {
		updateData["Description"] = requestData.Description
	}
	if requestData.Severity != "" {
		updateData["Severity"] = models.Severity(requestData.Severity)
	}
	if requestData.CVSS != "" {
		updateData["CVSS"] = cvssValue
	}
	if requestData.AffectedSystems != "" {
		updateData["AffectedSystems"] = requestData.AffectedSystems
	}
	if requestData.AffectedVersions != "" {
		updateData["AffectedVersions"] = requestData.AffectedVersions
	}
	if requestData.Solution != "" {
		updateData["Solution"] = requestData.Solution
	}
	if requestData.References != "" {
		updateData["References"] = requestData.References
	}
	if requestData.Tags != "" {
		updateData["Tags"] = requestData.Tags
	}
	if requestData.ExploitAvailable != existingVulnDB.ExploitAvailable {
		updateData["ExploitAvailable"] = requestData.ExploitAvailable
	}
	if requestData.PublishedDate != "" {
		updateData["PublishedDate"] = publishedDate
	}

	// 更新数据
	if err := utils.DB.Model(&existingVulnDB).Updates(updateData).Error; err != nil {
		log.Printf("更新漏洞库条目失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新漏洞库条目失败",
			"error":   err.Error(),
		})
		return
	}

	log.Printf("漏洞库条目更新成功, ID: %d", existingVulnDB.ID)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新漏洞库条目成功",
	})
}

// DeleteVulnDBEntry 删除漏洞库条目
func (v *VulnDBController) DeleteVulnDBEntry(c *gin.Context) {
	// 获取用户ID
	userID, _ := c.Get("userID")
	log.Printf("用户 ID: %v 尝试删除漏洞库条目", userID)

	id := c.Param("id")
	log.Printf("准备删除漏洞库条目, ID: %s", id)

	var vulnDBEntry models.VulnDB
	if err := utils.DB.First(&vulnDBEntry, id).Error; err != nil {
		log.Printf("漏洞库条目不存在, ID: %s, 错误: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "漏洞库条目不存在",
		})
		return
	}

	log.Printf("找到待删除漏洞库条目, ID: %d, CVE: %s, 标题: %s", vulnDBEntry.ID, vulnDBEntry.CVE, vulnDBEntry.Title)

	// 执行删除操作
	if err := utils.DB.Delete(&vulnDBEntry).Error; err != nil {
		log.Printf("删除漏洞库条目失败, ID: %d, 错误: %v", vulnDBEntry.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除漏洞库条目失败",
			"error":   err.Error(),
		})
		return
	}

	log.Printf("漏洞库条目删除成功, ID: %d", vulnDBEntry.ID)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除漏洞库条目成功",
	})
}

// BatchImportVulnDBEntries 批量导入漏洞库条目
func (v *VulnDBController) BatchImportVulnDBEntries(c *gin.Context) {
	// 获取用户ID
	userID, _ := c.Get("userID")
	log.Printf("用户 ID: %v 尝试批量导入漏洞库条目", userID)

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		log.Printf("获取上传文件失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "获取上传文件失败",
			"error":   err.Error(),
		})
		return
	}

	// 检查文件类型
	fileName := file.Filename
	isJSON := strings.HasSuffix(strings.ToLower(fileName), ".json")
	isCSV := strings.HasSuffix(strings.ToLower(fileName), ".csv")

	if !isJSON && !isCSV {
		log.Printf("不支持的文件类型: %s", fileName)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "不支持的文件类型，请上传JSON或CSV文件",
		})
		return
	}

	log.Printf("开始处理批量导入文件: %s, 大小: %d bytes", fileName, file.Size)

	// 打开文件
	openedFile, err := file.Open()
	if err != nil {
		log.Printf("打开文件失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "打开文件失败",
			"error":   err.Error(),
		})
		return
	}
	defer openedFile.Close()

	// 导入结果
	type ImportResult struct {
		Success       int      `json:"success"`
		Failed        int      `json:"failed"`
		FailedDetails []string `json:"failed_details"`
	}
	result := ImportResult{
		Success:       0,
		Failed:        0,
		FailedDetails: []string{},
	}

	// 解析文件并创建漏洞库条目
	now := time.Now()

	if isJSON {
		// 处理JSON文件
		var vulnDBList []struct {
			CVE              string  `json:"cve"`
			CWE              string  `json:"cwe"`
			Title            string  `json:"title"`
			Description      string  `json:"description"`
			Severity         string  `json:"severity"`
			CVSS             float64 `json:"cvss"`
			AffectedSystems  string  `json:"affected_systems"`
			AffectedVersions string  `json:"affected_versions"`
			Solution         string  `json:"solution"`
			References       string  `json:"references"`
			Tags             string  `json:"tags"`
			ExploitAvailable bool    `json:"exploit_available"`
			PublishedDate    string  `json:"published_date"`
		}

		decoder := json.NewDecoder(openedFile)
		if err := decoder.Decode(&vulnDBList); err != nil {
			log.Printf("解析JSON文件失败: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "解析JSON文件失败，请检查格式",
				"error":   err.Error(),
			})
			return
		}

		// 批量创建漏洞库条目
		for i, vulnData := range vulnDBList {
			// 验证必填字段
			if vulnData.Title == "" || vulnData.Description == "" || vulnData.Severity == "" {
				msg := fmt.Sprintf("记录 #%d 缺少必填字段", i+1)
				result.Failed++
				result.FailedDetails = append(result.FailedDetails, msg)
				continue
			}

			// 检查CVE是否已存在
			if vulnData.CVE != "" {
				var count int64
				utils.DB.Model(&models.VulnDB{}).Where("cve = ?", vulnData.CVE).Count(&count)
				if count > 0 {
					msg := fmt.Sprintf("记录 #%d CVE已存在: %s", i+1, vulnData.CVE)
					result.Failed++
					result.FailedDetails = append(result.FailedDetails, msg)
					continue
				}
			}

			// 处理日期
			var publishedDate time.Time
			if vulnData.PublishedDate != "" {
				var err error
				publishedDate, err = time.Parse("2006-01-02", vulnData.PublishedDate)
				if err != nil {
					log.Printf("记录 #%d 解析发布日期失败: %v", i+1, err)
					publishedDate = now // 默认为当前时间
				}
			} else {
				publishedDate = now
			}

			// 创建漏洞库对象
			vulnDBEntry := models.VulnDB{
				CVE:              vulnData.CVE,
				CWE:              vulnData.CWE,
				Title:            vulnData.Title,
				Description:      vulnData.Description,
				Severity:         models.Severity(vulnData.Severity),
				CVSS:             vulnData.CVSS,
				AffectedSystems:  vulnData.AffectedSystems,
				AffectedVersions: vulnData.AffectedVersions,
				Solution:         vulnData.Solution,
				References:       vulnData.References,
				Tags:             vulnData.Tags,
				ExploitAvailable: vulnData.ExploitAvailable,
				PublishedDate:    publishedDate,
				UpdatedDate:      now,
				CreatedAt:        now,
				UpdatedAt:        now,
			}

			// 保存到数据库
			if err := utils.DB.Create(&vulnDBEntry).Error; err != nil {
				log.Printf("创建漏洞库条目 #%d (%s) 失败: %v", i+1, vulnData.Title, err)
				result.Failed++
				result.FailedDetails = append(result.FailedDetails, fmt.Sprintf("创建 %s 失败: %v", vulnData.Title, err))
				continue
			}

			result.Success++
		}
	} else {
		// 处理CSV文件
		reader := csv.NewReader(openedFile)

		// 读取标题行
		headers, err := reader.Read()
		if err != nil {
			log.Printf("读取CSV标题行失败: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "读取CSV标题行失败，请检查格式",
				"error":   err.Error(),
			})
			return
		}

		// 创建标题索引映射
		headerMap := make(map[string]int)
		for i, header := range headers {
			headerMap[strings.ToLower(strings.TrimSpace(header))] = i
		}

		// 逐行读取并处理
		lineNum := 2 // 从第二行开始（标题行为第一行）
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("读取CSV行 #%d 失败: %v", lineNum, err)
				result.Failed++
				result.FailedDetails = append(result.FailedDetails, fmt.Sprintf("读取行 #%d 失败: %v", lineNum, err))
				lineNum++
				continue
			}

			// 安全地获取字段值
			getSafeString := func(record []string, headerMap map[string]int, field string) string {
				if idx, ok := headerMap[field]; ok && idx < len(record) {
					return strings.TrimSpace(record[idx])
				}
				return ""
			}

			// 从CSV记录创建漏洞库条目
			title := getSafeString(record, headerMap, "title")
			description := getSafeString(record, headerMap, "description")
			severity := getSafeString(record, headerMap, "severity")
			cve := getSafeString(record, headerMap, "cve")
			cwe := getSafeString(record, headerMap, "cwe")

			// 验证必填字段
			if title == "" || description == "" || severity == "" {
				msg := fmt.Sprintf("行 #%d 缺少必填字段", lineNum)
				result.Failed++
				result.FailedDetails = append(result.FailedDetails, msg)
				lineNum++
				continue
			}

			// 检查CVE是否已存在
			if cve != "" {
				var count int64
				utils.DB.Model(&models.VulnDB{}).Where("cve = ?", cve).Count(&count)
				if count > 0 {
					msg := fmt.Sprintf("行 #%d CVE已存在: %s", lineNum, cve)
					result.Failed++
					result.FailedDetails = append(result.FailedDetails, msg)
					lineNum++
					continue
				}
			}

			// 处理CVSS，需要转换字符串为浮点数
			cvss := 0.0
			if cvssStr := getSafeString(record, headerMap, "cvss"); cvssStr != "" {
				if cvssVal, err := strconv.ParseFloat(cvssStr, 64); err == nil {
					cvss = cvssVal
				}
			}

			// 处理日期
			publishedDate := now
			if publishedDateStr := getSafeString(record, headerMap, "published_date"); publishedDateStr != "" {
				if date, err := time.Parse("2006-01-02", publishedDateStr); err == nil {
					publishedDate = date
				}
			}

			// 处理布尔值
			exploitAvailable := false
			if exploitStr := getSafeString(record, headerMap, "exploit_available"); exploitStr != "" {
				exploitAvailable = strings.ToLower(exploitStr) == "true" || exploitStr == "1" || strings.ToLower(exploitStr) == "yes"
			}

			// 创建漏洞库对象
			vulnDBEntry := models.VulnDB{
				CVE:              cve,
				CWE:              cwe,
				Title:            title,
				Description:      description,
				Severity:         models.Severity(severity),
				CVSS:             cvss,
				AffectedSystems:  getSafeString(record, headerMap, "affected_systems"),
				AffectedVersions: getSafeString(record, headerMap, "affected_versions"),
				Solution:         getSafeString(record, headerMap, "solution"),
				References:       getSafeString(record, headerMap, "references"),
				Tags:             getSafeString(record, headerMap, "tags"),
				ExploitAvailable: exploitAvailable,
				PublishedDate:    publishedDate,
				UpdatedDate:      now,
				CreatedAt:        now,
				UpdatedAt:        now,
			}

			// 保存到数据库
			if err := utils.DB.Create(&vulnDBEntry).Error; err != nil {
				log.Printf("创建漏洞库条目 #%d (%s) 失败: %v", lineNum, title, err)
				result.Failed++
				result.FailedDetails = append(result.FailedDetails, fmt.Sprintf("创建 %s 失败: %v", title, err))
				lineNum++
				continue
			}

			result.Success++
			lineNum++
		}
	}

	log.Printf("批量导入漏洞库条目完成, 成功: %d, 失败: %d", result.Success, result.Failed)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "批量导入漏洞库条目完成",
		"data":    result,
	})
}
