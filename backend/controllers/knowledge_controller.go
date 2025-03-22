package controllers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vulnark/vulnark/models"
	"github.com/vulnark/vulnark/utils"
)

// KnowledgeController 知识库控制器
type KnowledgeController struct{}

// ListKnowledgeItems 获取知识库列表
func (k *KnowledgeController) ListKnowledgeItems(c *gin.Context) {
	// 获取查询参数
	keyword := c.Query("keyword")
	typeFilter := c.Query("type")
	category := c.Query("category")
	tags := c.Query("tags")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// 构建查询
	query := utils.DB.Model(&models.Knowledge{})

	// 条件过滤
	if keyword != "" {
		query = query.Where("title LIKE ? OR content LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if typeFilter != "" {
		query = query.Where("type = ?", typeFilter)
	}
	if category != "" {
		query = query.Where("categories LIKE ?", "%"+category+"%")
	}
	if tags != "" {
		query = query.Where("tags LIKE ?", "%"+tags+"%")
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 获取分页数据
	var knowledgeItems []models.Knowledge
	query.Limit(pageSize).Offset((page - 1) * pageSize).Order("created_at DESC").Find(&knowledgeItems)

	// 转换为响应格式
	var responseItems []map[string]interface{}
	for _, item := range knowledgeItems {
		respItem := map[string]interface{}{
			"id":         item.ID,
			"title":      item.Title,
			"type":       string(item.Type),
			"author":     item.Author,
			"tags":       item.Tags,
			"categories": item.Categories,
			"view_count": item.ViewCount,
			"created_at": item.CreatedAt,
			"updated_at": item.UpdatedAt,
		}
		responseItems = append(responseItems, respItem)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"items": responseItems,
			"total": total,
		},
	})
}

// GetKnowledgeByID 获取知识库详情
func (k *KnowledgeController) GetKnowledgeByID(c *gin.Context) {
	id := c.Param("id")
	log.Printf("获取知识库详情, ID: %s", id)

	var knowledgeItem models.Knowledge
	if err := utils.DB.First(&knowledgeItem, id).Error; err != nil {
		log.Printf("知识库不存在, ID: %s, 错误: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "知识库不存在",
		})
		return
	}

	// 增加查看次数
	knowledgeItem.IncrementViewCount()
	utils.DB.Save(&knowledgeItem)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": knowledgeItem,
	})
}

// CreateKnowledge 创建知识库
func (k *KnowledgeController) CreateKnowledge(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Printf("用户ID未找到")
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "用户ID未找到",
		})
		return
	}
	log.Printf("用户 ID: %v 尝试创建知识库", userID)

	// 接收请求数据
	var requestData struct {
		Title            string `json:"title"`
		Content          string `json:"content"`
		Type             string `json:"type"`
		Tags             string `json:"tags"`
		Categories       string `json:"categories"`
		RelatedVulnTypes string `json:"related_vuln_types"`
		Author           string `json:"author"`
		References       string `json:"references"`
		Attachments      string `json:"attachments"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		log.Printf("创建知识库请求参数错误: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	log.Printf("接收到的创建知识库数据: %+v", requestData)

	// 验证必填字段
	if requestData.Title == "" || requestData.Content == "" || requestData.Type == "" {
		log.Printf("缺少必要的字段: title=%s, content=%s, type=%s",
			requestData.Title, requestData.Content, requestData.Type)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "缺少必要的字段: 标题、内容或类型",
		})
		return
	}

	// 从userID获取作者ID
	var authorID uint
	switch v := userID.(type) {
	case float64:
		authorID = uint(v)
	case int:
		authorID = uint(v)
	case uint:
		authorID = v
	case int64:
		authorID = uint(v)
	case uint64:
		authorID = uint(v)
	default:
		log.Printf("无法转换用户ID类型: %T, 值: %v", userID, userID)
		authorID = 0
	}

	// 创建知识库对象
	now := time.Now()
	knowledgeItem := models.Knowledge{
		Title:            requestData.Title,
		Content:          requestData.Content,
		Type:             models.KnowledgeType(requestData.Type),
		Tags:             requestData.Tags,
		Categories:       requestData.Categories,
		RelatedVulnTypes: requestData.RelatedVulnTypes,
		Author:           requestData.Author,
		AuthorID:         authorID,
		ViewCount:        0,
		References:       requestData.References,
		Attachments:      requestData.Attachments,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	// 保存到数据库
	if err := utils.DB.Create(&knowledgeItem).Error; err != nil {
		log.Printf("创建知识库失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建知识库失败",
			"error":   err.Error(),
		})
		return
	}

	log.Printf("知识库创建成功, ID: %d, 标题: %s", knowledgeItem.ID, knowledgeItem.Title)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建知识库成功",
		"data": gin.H{
			"id": knowledgeItem.ID,
		},
	})
}

// UpdateKnowledge 更新知识库
func (k *KnowledgeController) UpdateKnowledge(c *gin.Context) {
	// 获取用户ID
	userID, _ := c.Get("userID")
	log.Printf("用户 ID: %v 尝试更新知识库", userID)

	id := c.Param("id")
	log.Printf("准备更新知识库, ID: %s", id)

	var knowledgeItem models.Knowledge
	if err := utils.DB.First(&knowledgeItem, id).Error; err != nil {
		log.Printf("知识库不存在, ID: %s, 错误: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "知识库不存在",
		})
		return
	}

	// 接收请求数据
	var requestData struct {
		Title            string `json:"title"`
		Content          string `json:"content"`
		Type             string `json:"type"`
		Tags             string `json:"tags"`
		Categories       string `json:"categories"`
		RelatedVulnTypes string `json:"related_vuln_types"`
		Author           string `json:"author"`
		References       string `json:"references"`
		Attachments      string `json:"attachments"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		log.Printf("更新知识库请求参数错误: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	log.Printf("接收到的更新知识库数据: %+v", requestData)

	// 更新知识库对象
	updateData := map[string]interface{}{
		"UpdatedAt": time.Now(),
	}

	if requestData.Title != "" {
		updateData["Title"] = requestData.Title
	}
	if requestData.Content != "" {
		updateData["Content"] = requestData.Content
	}
	if requestData.Type != "" {
		updateData["Type"] = models.KnowledgeType(requestData.Type)
	}
	if requestData.Tags != "" {
		updateData["Tags"] = requestData.Tags
	}
	if requestData.Categories != "" {
		updateData["Categories"] = requestData.Categories
	}
	if requestData.RelatedVulnTypes != "" {
		updateData["RelatedVulnTypes"] = requestData.RelatedVulnTypes
	}
	if requestData.Author != "" {
		updateData["Author"] = requestData.Author
	}
	if requestData.References != "" {
		updateData["References"] = requestData.References
	}
	if requestData.Attachments != "" {
		updateData["Attachments"] = requestData.Attachments
	}

	// 更新知识库
	if err := utils.DB.Model(&knowledgeItem).Updates(updateData).Error; err != nil {
		log.Printf("更新知识库失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新知识库失败",
			"error":   err.Error(),
		})
		return
	}

	log.Printf("知识库更新成功, ID: %d", knowledgeItem.ID)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新知识库成功",
	})
}

// DeleteKnowledge 删除知识库
func (k *KnowledgeController) DeleteKnowledge(c *gin.Context) {
	// 获取用户ID
	userID, _ := c.Get("userID")
	log.Printf("用户 ID: %v 尝试删除知识库", userID)

	id := c.Param("id")
	log.Printf("准备删除知识库, ID: %s", id)

	var knowledgeItem models.Knowledge
	if err := utils.DB.First(&knowledgeItem, id).Error; err != nil {
		log.Printf("知识库不存在, ID: %s, 错误: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "知识库不存在",
		})
		return
	}

	log.Printf("找到待删除知识库, ID: %d, 标题: %s", knowledgeItem.ID, knowledgeItem.Title)

	// 执行删除操作
	if err := utils.DB.Delete(&knowledgeItem).Error; err != nil {
		log.Printf("删除知识库失败, ID: %d, 错误: %v", knowledgeItem.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除知识库失败",
			"error":   err.Error(),
		})
		return
	}

	log.Printf("知识库删除成功, ID: %d", knowledgeItem.ID)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除知识库成功",
	})
}

// GetKnowledgeTypes 获取知识库类型列表
func (k *KnowledgeController) GetKnowledgeTypes(c *gin.Context) {
	types := []map[string]string{
		{"value": string(models.KnowledgeTypeTutorial), "label": "教程"},
		{"value": string(models.KnowledgeTypeGuide), "label": "指南"},
		{"value": string(models.KnowledgeTypeReference), "label": "参考资料"},
		{"value": string(models.KnowledgeTypeMitigation), "label": "缓解策略"},
		{"value": string(models.KnowledgeTypeProc), "label": "操作流程"},
		{"value": string(models.KnowledgeTypeOther), "label": "其他"},
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": types,
	})
}

// GetKnowledgeCategories 获取知识库分类列表
func (k *KnowledgeController) GetKnowledgeCategories(c *gin.Context) {
	// 从数据库中获取所有不同的分类
	var categories []string
	utils.DB.Model(&models.Knowledge{}).Pluck("DISTINCT categories", &categories)

	// 处理和提取唯一的分类
	uniqueCategories := make(map[string]bool)
	var resultCategories []string

	for _, categoryStr := range categories {
		// 处理逗号分隔的分类
		// 这里实现可能需要基于数据库中存储分类的格式来调整
		// 假设分类是以逗号分隔的字符串
		// TODO: 实现分类提取逻辑

		// 简单实现：直接添加不为空的分类
		if categoryStr != "" {
			if _, exists := uniqueCategories[categoryStr]; !exists {
				uniqueCategories[categoryStr] = true
				resultCategories = append(resultCategories, categoryStr)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": resultCategories,
	})
}
