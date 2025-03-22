package controllers

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/vulnark/vulnark/models"
	"github.com/vulnark/vulnark/utils"
)

// AssetController 资产控制器
type AssetController struct{}

// ListAssets 获取资产列表
func (a *AssetController) ListAssets(c *gin.Context) {
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	keyword := c.Query("keyword")
	status := c.Query("status")
	assetType := c.Query("type")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	// 构建查询
	query := utils.DB.Model(&models.Asset{})

	// 添加过滤条件
	if keyword != "" {
		query = query.Where("name LIKE ? OR ip_address LIKE ? OR department LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if assetType != "" {
		query = query.Where("type = ?", assetType)
	}

	// 获取总记录数
	var total int64
	query.Count(&total)

	// 查询分页数据
	var assets []models.Asset
	query.Limit(pageSize).Offset((page - 1) * pageSize).Order("created_at DESC").Find(&assets)

	// 转换为前端需要的格式
	var responseItems []map[string]interface{}
	for _, asset := range assets {
		item := map[string]interface{}{
			"id":              asset.ID,
			"name":            asset.Name,
			"type":            string(asset.Type),
			"status":          string(asset.Status),
			"ip":              asset.IPAddress,
			"department":      asset.Department,
			"owner":           asset.Owner,
			"operatingSystem": asset.OS,
			"version":         asset.Version,
			"url":             asset.URL,
			"notes":           asset.Notes,
			"createdAt":       asset.CreatedAt,
			"updatedAt":       asset.UpdatedAt,
		}

		// 处理标签
		if asset.Tags != "" {
			item["tags"] = strings.Split(asset.Tags, ",")
		} else {
			item["tags"] = []string{}
		}

		responseItems = append(responseItems, item)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"items": responseItems,
			"total": total,
		},
		"message": "获取资产列表成功",
	})
}

// GetAssetByID 获取单个资产详情
func (a *AssetController) GetAssetByID(c *gin.Context) {
	id := c.Param("id")

	var asset models.Asset
	if err := utils.DB.First(&asset, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "资产不存在",
		})
		return
	}

	// 创建前端需要的响应数据结构
	response := struct {
		ID              uint      `json:"id"`
		Name            string    `json:"name"`
		Type            string    `json:"type"`
		Status          string    `json:"status"`
		IP              string    `json:"ip"`
		Department      string    `json:"department"`
		Owner           string    `json:"owner"`
		OperatingSystem string    `json:"operatingSystem"`
		Version         string    `json:"version"`
		URL             string    `json:"url"`
		Tags            []string  `json:"tags"`
		Notes           string    `json:"notes"`
		CreatedAt       time.Time `json:"createdAt"`
		UpdatedAt       time.Time `json:"updatedAt"`
	}{
		ID:              asset.ID,
		Name:            asset.Name,
		Type:            string(asset.Type),
		Status:          string(asset.Status),
		IP:              asset.IPAddress,
		Department:      asset.Department,
		Owner:           asset.Owner,
		OperatingSystem: asset.OS,
		Version:         asset.Version,
		URL:             asset.URL,
		Notes:           asset.Notes,
		CreatedAt:       asset.CreatedAt,
		UpdatedAt:       asset.UpdatedAt,
	}

	// 处理标签字符串转为数组
	if asset.Tags != "" {
		response.Tags = strings.Split(asset.Tags, ",")
	} else {
		response.Tags = []string{}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    response,
		"message": "获取资产详情成功",
	})
}

// CreateAsset 创建新资产
func (a *AssetController) CreateAsset(c *gin.Context) {
	// 打印权限信息
	userID, _ := c.Get("userID")
	log.Printf("用户 ID: %v 尝试创建资产", userID)

	// 接收前端数据
	var requestData struct {
		Name            string   `json:"name"`
		IP              string   `json:"ip"`
		Type            string   `json:"type"`
		Status          string   `json:"status"`
		Department      string   `json:"department"`
		Owner           string   `json:"owner"`
		OperatingSystem string   `json:"operatingSystem"`
		Version         string   `json:"version"`
		URL             string   `json:"url"`
		MacAddress      string   `json:"macAddress"`
		Tags            []string `json:"tags"`
		Notes           string   `json:"notes"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		log.Printf("请求参数绑定错误: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	// 打印接收到的数据，用于调试
	log.Printf("接收到的资产数据: %+v", requestData)

	// 验证必填字段
	if requestData.Name == "" || requestData.IP == "" || requestData.Type == "" {
		log.Printf("缺少必要的字段: name=%s, ip=%s, type=%s",
			requestData.Name, requestData.IP, requestData.Type)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "缺少必要的字段: 名称、IP地址或资产类型",
		})
		return
	}

	// 生成唯一标识符
	identifier := requestData.Name + "-" + requestData.IP

	// 输出将要使用的标识符
	log.Printf("生成的资产标识符: %s", identifier)

	// 首先检查IP地址是否已经存在
	var existingAssetByIP models.Asset
	resultIP := utils.DB.Where("ip_address = ?", requestData.IP).First(&existingAssetByIP)
	if resultIP.Error == nil {
		// IP地址已存在
		log.Printf("IP地址已存在: %s, 资产名称: %s, ID: %d",
			requestData.IP, existingAssetByIP.Name, existingAssetByIP.ID)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "IP地址 " + requestData.IP + " 已被资产 '" + existingAssetByIP.Name + "' 使用，请使用其他IP地址或更新现有资产",
		})
		return
	} else if !gorm.IsRecordNotFoundError(resultIP.Error) {
		// 查询出错
		log.Printf("查询IP地址时出错: %v", resultIP.Error)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询资产时出错: " + resultIP.Error.Error(),
		})
		return
	}

	// 检查标识符是否已存在
	var existingAsset models.Asset
	result := utils.DB.Where("identifier = ?", identifier).First(&existingAsset)
	if result.Error == nil {
		// 标识符已存在
		log.Printf("资产标识符已存在: %s, ID: %d", identifier, existingAsset.ID)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "具有相同名称和IP地址的资产已存在，资产标识符 '" + identifier + "' 重复",
		})
		return
	} else {
		// 检查错误类型 - gorm v1 使用 gorm.IsRecordNotFoundError
		if !gorm.IsRecordNotFoundError(result.Error) {
			// 其他数据库错误
			log.Printf("查询资产时出错: %v", result.Error)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "查询资产时出错: " + result.Error.Error(),
			})
			return
		}
	}

	// 设置资产状态默认值
	assetStatus := models.AssetStatus(requestData.Status)
	if assetStatus == "" {
		assetStatus = models.AssetStatusActive
	}

	// 创建资产对象
	asset := models.Asset{
		Name:       requestData.Name,
		Type:       models.AssetType(requestData.Type),
		Status:     assetStatus,
		IPAddress:  requestData.IP,
		OS:         requestData.OperatingSystem,
		Version:    requestData.Version,
		URL:        requestData.URL,
		Department: requestData.Department,
		Owner:      requestData.Owner,
		Notes:      requestData.Notes,
		Identifier: identifier,              // 自动生成标识符
		Importance: models.ImportanceMedium, // 默认重要性
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// 处理标签
	if len(requestData.Tags) > 0 {
		// 将标签数组转为逗号分隔的字符串
		asset.Tags = strings.Join(requestData.Tags, ",")
	}

	// 保存到数据库
	log.Printf("准备创建资产: %+v", asset)
	err := utils.DB.Create(&asset).Error
	if err != nil {
		log.Printf("创建资产失败: %v", err)

		// 提供更详细的错误信息
		errMsg := "创建资产失败: " + err.Error()

		// 检查是否是唯一键冲突
		if strings.Contains(err.Error(), "Duplicate entry") && strings.Contains(err.Error(), "uix_assets_identifier") {
			errMsg = "资产标识符重复，请修改名称或IP地址"
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": errMsg,
		})
		return
	}

	// 发送资产创建通知
	go func() {
		notificationManager, err := utils.NewNotificationManager()
		if err != nil {
			log.Printf("获取通知管理器失败: %v", err)
			return
		}
		notificationManager.SendAssetNotification(utils.EventAssetCreate, &asset)
	}()

	// 返回成功信息
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "资产创建成功",
		"data": gin.H{
			"id": asset.ID,
		},
	})
}

// UpdateAsset 更新资产
func (a *AssetController) UpdateAsset(c *gin.Context) {
	// 获取用户ID
	userID, _ := c.Get("userID")
	log.Printf("用户 ID: %v 尝试更新资产", userID)

	id := c.Param("id")
	log.Printf("准备更新资产，ID: %s", id)

	var asset models.Asset
	if err := utils.DB.First(&asset, id).Error; err != nil {
		log.Printf("资产不存在，ID: %s, 错误: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "资产不存在",
		})
		return
	}

	// 接收前端数据
	var requestData struct {
		Name            string   `json:"name"`
		IP              string   `json:"ip"`
		Type            string   `json:"type"`
		Status          string   `json:"status"`
		Department      string   `json:"department"`
		Owner           string   `json:"owner"`
		OperatingSystem string   `json:"operatingSystem"`
		Version         string   `json:"version"`
		URL             string   `json:"url"`
		MacAddress      string   `json:"macAddress"`
		Tags            []string `json:"tags"`
		Notes           string   `json:"notes"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		log.Printf("更新资产请求参数错误: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	log.Printf("接收到的更新数据: %+v", requestData)

	// 验证必填字段
	if requestData.Name == "" || requestData.IP == "" || requestData.Type == "" {
		log.Printf("缺少必要的字段: name=%s, ip=%s, type=%s",
			requestData.Name, requestData.IP, requestData.Type)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "缺少必要的字段: 名称、IP地址或资产类型",
		})
		return
	}

	// 如果IP地址已更改，检查新IP是否已被使用
	if requestData.IP != asset.IPAddress {
		var existingAssetByIP models.Asset
		resultIP := utils.DB.Where("ip_address = ? AND id != ?", requestData.IP, asset.ID).First(&existingAssetByIP)
		if resultIP.Error == nil {
			// IP地址已存在于其他资产
			log.Printf("IP地址已被其他资产使用: %s, 资产名称: %s, ID: %d",
				requestData.IP, existingAssetByIP.Name, existingAssetByIP.ID)
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "IP地址 " + requestData.IP + " 已被资产 '" + existingAssetByIP.Name + "' 使用，请使用其他IP地址",
			})
			return
		} else if !gorm.IsRecordNotFoundError(resultIP.Error) {
			// 查询出错
			log.Printf("查询IP地址时出错: %v", resultIP.Error)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "查询资产时出错: " + resultIP.Error.Error(),
			})
			return
		}
	}

	// 设置资产状态默认值
	assetStatus := models.AssetStatus(requestData.Status)
	if assetStatus == "" {
		assetStatus = models.AssetStatusActive
	}

	// 原始标识符
	oldIdentifier := asset.Identifier

	// 新标识符
	newIdentifier := requestData.Name + "-" + requestData.IP

	// 如果标识符已更改，检查新标识符是否已被使用
	if newIdentifier != oldIdentifier {
		var existingAsset models.Asset
		result := utils.DB.Where("identifier = ? AND id != ?", newIdentifier, asset.ID).First(&existingAsset)
		if result.Error == nil {
			// 标识符已存在于其他资产
			log.Printf("资产标识符已被其他资产使用: %s, ID: %d", newIdentifier, existingAsset.ID)
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "具有相同名称和IP地址的资产已存在，资产标识符 '" + newIdentifier + "' 重复",
			})
			return
		} else if !gorm.IsRecordNotFoundError(result.Error) {
			// 其他数据库错误
			log.Printf("查询资产时出错: %v", result.Error)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "查询资产时出错: " + result.Error.Error(),
			})
			return
		}
	}

	// 更新资产对象
	updateData := models.Asset{
		ID:         asset.ID,
		Name:       requestData.Name,
		Type:       models.AssetType(requestData.Type),
		Status:     assetStatus,
		IPAddress:  requestData.IP,
		OS:         requestData.OperatingSystem,
		Version:    requestData.Version,
		URL:        requestData.URL,
		Department: requestData.Department,
		Owner:      requestData.Owner,
		Notes:      requestData.Notes,
		Identifier: newIdentifier, // 使用新标识符
		UpdatedAt:  time.Now(),
	}

	// 处理标签
	if len(requestData.Tags) > 0 {
		// 将标签数组转为逗号分隔的字符串
		updateData.Tags = strings.Join(requestData.Tags, ",")
	}

	log.Printf("准备更新资产，资产ID: %d, 新标识符: %s", asset.ID, newIdentifier)
	if err := utils.DB.Model(&asset).Updates(updateData).Error; err != nil {
		log.Printf("更新资产失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新资产失败: " + err.Error(),
		})
		return
	}

	// 发送资产更新通知
	go func() {
		notificationManager, err := utils.NewNotificationManager()
		if err != nil {
			log.Printf("获取通知管理器失败: %v", err)
			return
		}
		notificationManager.SendAssetNotification(utils.EventAssetUpdate, &asset)
	}()

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "资产更新成功",
		"data": gin.H{
			"id": asset.ID,
		},
	})
}

// DeleteAsset 删除资产
func (a *AssetController) DeleteAsset(c *gin.Context) {
	// 获取用户ID
	userID, _ := c.Get("userID")
	log.Printf("用户 ID: %v 尝试删除资产", userID)

	id := c.Param("id")
	log.Printf("准备删除资产，ID: %s", id)

	var asset models.Asset
	if err := utils.DB.First(&asset, id).Error; err != nil {
		log.Printf("资产不存在，ID: %s, 错误: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "资产不存在",
		})
		return
	}

	log.Printf("找到待删除资产，ID: %d, 名称: %s, 标识符: %s",
		asset.ID, asset.Name, asset.Identifier)

	// 获取资产信息用于发送通知
	var assetInfo models.Asset
	if err := utils.DB.First(&assetInfo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "资产不存在",
		})
		return
	}

	// 执行删除
	if err := utils.DB.Delete(&asset).Error; err != nil {
		log.Printf("删除资产失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除资产失败: " + err.Error(),
		})
		return
	}

	// 发送资产删除通知
	go func() {
		notificationManager, err := utils.NewNotificationManager()
		if err != nil {
			log.Printf("获取通知管理器失败: %v", err)
			return
		}
		notificationManager.SendAssetNotification(utils.EventAssetDelete, &assetInfo)
	}()

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "资产删除成功",
	})
}

// BatchImportAssets 批量导入资产
func (a *AssetController) BatchImportAssets(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请选择要上传的文件",
		})
		return
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "无法打开上传的文件",
		})
		return
	}
	defer src.Close()

	// 读取文件内容
	fileBytes, err := io.ReadAll(src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "读取文件内容失败",
		})
		return
	}

	// 初始化结果
	successCount := 0
	failCount := 0
	errors := []string{}
	var assets []models.Asset

	fileName := file.Filename
	fileExt := filepath.Ext(fileName)

	// 根据文件扩展名处理不同类型的文件
	switch strings.ToLower(fileExt) {
	case ".csv":
		// 处理CSV文件
		reader := csv.NewReader(bytes.NewReader(fileBytes))
		records, err := reader.ReadAll()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "CSV文件格式不正确",
			})
			return
		}

		// 确保有数据并且有头行
		if len(records) < 2 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "CSV文件必须包含标题行和至少一行数据",
			})
			return
		}

		// 从第二行开始解析数据（跳过标题行）
		for i, record := range records {
			if i == 0 {
				// 跳过标题行
				continue
			}

			// 确保行中有足够的字段
			if len(record) < 5 {
				failCount++
				errors = append(errors, fmt.Sprintf("第%d行数据不完整", i+1))
				continue
			}

			// 创建资产对象
			asset := models.Asset{
				Name:        record[0],
				IPAddress:   record[1],                                  // IP地址
				Identifier:  fmt.Sprintf("%s_%s", record[0], record[1]), // 使用名称和IP作为标识符
				Type:        models.AssetType(record[3]),                // 资产类型
				Department:  record[4],                                  // 部门
				Owner:       record[5],                                  // 负责人
				Status:      models.AssetStatus(record[6]),              // 状态
				Importance:  models.AssetImportance(record[7]),          // 重要性/风险等级
				Description: record[8],                                  // 描述
			}

			assets = append(assets, asset)
		}

	case ".json":
		// 处理JSON文件
		var jsonAssets []struct {
			Name        string `json:"name"`
			IPAddress   string `json:"ip"`
			Identifier  string `json:"identifier,omitempty"`
			Type        string `json:"type"`
			Department  string `json:"department"`
			Owner       string `json:"owner"`
			Status      string `json:"status"`
			Importance  string `json:"risk_level"`
			Description string `json:"description"`
		}

		if err := json.Unmarshal(fileBytes, &jsonAssets); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "JSON文件格式不正确: " + err.Error(),
			})
			return
		}

		// 遍历JSON数据并创建资产
		for i, jsonAsset := range jsonAssets {
			if jsonAsset.Name == "" || jsonAsset.IPAddress == "" || jsonAsset.Type == "" {
				failCount++
				errors = append(errors, fmt.Sprintf("第%d条数据缺少必要字段", i+1))
				continue
			}

			// 生成标识符
			identifier := jsonAsset.Identifier
			if identifier == "" {
				identifier = fmt.Sprintf("%s_%s", jsonAsset.Name, jsonAsset.IPAddress)
			}

			asset := models.Asset{
				Name:        jsonAsset.Name,
				IPAddress:   jsonAsset.IPAddress,
				Identifier:  identifier,
				Type:        models.AssetType(jsonAsset.Type),
				Department:  jsonAsset.Department,
				Owner:       jsonAsset.Owner,
				Status:      models.AssetStatus(jsonAsset.Status),
				Importance:  models.AssetImportance(jsonAsset.Importance),
				Description: jsonAsset.Description,
			}

			assets = append(assets, asset)
		}

	case ".xlsx":
		// Excel文件暂不支持，可以引入外部库实现
		c.JSON(http.StatusNotImplemented, gin.H{
			"code":    501,
			"message": "Excel文件导入功能尚未实现",
		})
		return

	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "不支持的文件格式，请上传CSV或JSON文件",
		})
		return
	}

	// 批量保存资产
	for i, asset := range assets {
		// 验证必填字段
		if asset.Name == "" || asset.IPAddress == "" || string(asset.Type) == "" {
			failCount++
			errors = append(errors, fmt.Sprintf("第%d条数据缺少必要字段", i+1))
			continue
		}

		// 设置默认值
		if asset.Status == "" {
			asset.Status = models.AssetStatusActive
		}
		if asset.Importance == "" {
			asset.Importance = models.ImportanceMedium
		}

		// 创建记录
		if err := utils.DB.Create(&asset).Error; err != nil {
			failCount++
			errors = append(errors, fmt.Sprintf("保存第%d条数据失败: %v", i+1, err))
		} else {
			successCount++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"successCount": successCount,
			"failCount":    failCount,
			"errors":       errors,
			"fileName":     fileName,
			"fileSize":     file.Size,
		},
		"message": fmt.Sprintf("成功导入%d条数据，失败%d条", successCount, failCount),
		"success": true,
	})
}

// ExportAssets 导出资产列表
func (a *AssetController) ExportAssets(c *gin.Context) {
	// TODO: 实现导出逻辑
	// 这里简单返回一个临时响应

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "导出功能正在开发中",
	})
}

// GetAssetVulnerabilities 获取资产关联的漏洞
func (a *AssetController) GetAssetVulnerabilities(c *gin.Context) {
	assetID := c.Param("id")

	log.Printf("开始获取资产ID %s 的关联漏洞", assetID)

	// 检查资产是否存在
	var asset models.Asset
	if err := utils.DB.First(&asset, assetID).Error; err != nil {
		log.Printf("资产ID %s 不存在: %v", assetID, err)
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "资产不存在",
		})
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	log.Printf("查询资产漏洞参数: assetID=%s, page=%d, pageSize=%d", assetID, page, pageSize)

	// 尝试通过直接SQL查询获取关联的漏洞
	var vulnerabilities []models.Vulnerability
	var total int64

	query := `
		SELECT v.* FROM vulnerabilities v 
		JOIN vulnerability_assets va ON v.id = va.vulnerability_id 
		WHERE va.asset_id = ? AND v.deleted_at IS NULL
		ORDER BY v.created_at DESC LIMIT ? OFFSET ?
	`

	countQuery := `
		SELECT COUNT(*) FROM vulnerabilities v 
		JOIN vulnerability_assets va ON v.id = va.vulnerability_id 
		WHERE va.asset_id = ? AND v.deleted_at IS NULL
	`

	// 查询总数
	err := utils.DB.Raw(countQuery, assetID).Count(&total).Error
	if err != nil {
		// 检查是否是表不存在的错误
		if strings.Contains(err.Error(), "doesn't exist") ||
			strings.Contains(err.Error(), "not found") ||
			strings.Contains(err.Error(), "Unknown table") {
			log.Printf("可能是关联表不存在: %v", err)
			c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"message": "资产没有关联漏洞",
				"data": map[string]interface{}{
					"total":    0,
					"pageSize": pageSize,
					"page":     page,
					"items":    []interface{}{},
				},
			})
		} else {
			log.Printf("查询漏洞总数失败: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "查询漏洞总数失败: " + err.Error(),
			})
		}
		return
	}

	if total == 0 {
		// 没有关联漏洞
		log.Printf("资产ID %s 没有关联漏洞", assetID)
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "资产没有关联漏洞",
			"data": map[string]interface{}{
				"total":    0,
				"pageSize": pageSize,
				"page":     page,
				"items":    []interface{}{},
			},
		})
		return
	}

	// 查询漏洞详情
	err = utils.DB.Raw(query, assetID, pageSize, offset).Scan(&vulnerabilities).Error
	if err != nil {
		log.Printf("查询漏洞详情失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取资产关联漏洞失败: " + err.Error(),
		})
		return
	}

	log.Printf("成功查询到资产ID %s 的关联漏洞 %d 条", assetID, len(vulnerabilities))

	// 构建响应数据
	var result []map[string]interface{}
	for _, vuln := range vulnerabilities {
		result = append(result, map[string]interface{}{
			"id":           vuln.ID,
			"title":        vuln.Title,
			"severity":     vuln.Severity,
			"status":       vuln.Status,
			"cve_id":       vuln.CVE,
			"discoveredAt": utils.FormatTimeCST(vuln.CreatedAt),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取资产关联漏洞成功",
		"data": map[string]interface{}{
			"total":    total,
			"pageSize": pageSize,
			"page":     page,
			"items":    result,
		},
	})
}

// BatchDeleteAssets 批量删除资产
func (a *AssetController) BatchDeleteAssets(c *gin.Context) {
	// 获取用户ID
	userID, _ := c.Get("userID")
	log.Printf("用户 ID: %v 尝试批量删除资产", userID)

	// 解析请求体中的资产ID列表
	var requestBody struct {
		IDs []uint `json:"ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Printf("解析请求参数失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	if len(requestBody.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "未提供要删除的资产ID",
		})
		return
	}

	log.Printf("准备批量删除资产，ID数量: %d, IDs: %v", len(requestBody.IDs), requestBody.IDs)

	// 查询这些资产是否存在
	var assets []models.Asset
	if err := utils.DB.Where("id IN (?)", requestBody.IDs).Find(&assets).Error; err != nil {
		log.Printf("查询资产失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询资产失败: " + err.Error(),
		})
		return
	}

	// 检查是否所有资产都存在
	if len(assets) != len(requestBody.IDs) {
		log.Printf("部分资产不存在，请求ID数量: %d, 找到资产数量: %d", len(requestBody.IDs), len(assets))
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "部分资产不存在",
		})
		return
	}

	// 执行批量删除
	if err := utils.DB.Where("id IN (?)", requestBody.IDs).Delete(&models.Asset{}).Error; err != nil {
		log.Printf("批量删除资产失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "批量删除资产失败: " + err.Error(),
		})
		return
	}

	// 发送资产删除通知
	go func() {
		notificationManager, err := utils.NewNotificationManager()
		if err != nil {
			log.Printf("获取通知管理器失败: %v", err)
			return
		}

		for _, asset := range assets {
			notificationManager.SendAssetNotification(utils.EventAssetDelete, &asset)
		}
	}()

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": fmt.Sprintf("成功删除 %d 个资产", len(requestBody.IDs)),
		"data": gin.H{
			"deleted_count": len(requestBody.IDs),
		},
	})
}
