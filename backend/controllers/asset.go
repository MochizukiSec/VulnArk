package controllers

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"vuln-management/models"
	"vuln-management/services"
)

// AssetController 处理资产相关的HTTP请求
type AssetController struct {
	assetService *services.AssetService
}

// NewAssetController 创建资产控制器
func NewAssetController(assetService *services.AssetService) *AssetController {
	return &AssetController{
		assetService: assetService,
	}
}

// CreateAsset 创建新资产
func (c *AssetController) CreateAsset(ctx *gin.Context) {
	var createAsset models.AssetCreate
	if err := ctx.ShouldBindJSON(&createAsset); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从上下文获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}
	objID, _ := primitive.ObjectIDFromHex(userID.(string))

	asset, err := c.assetService.CreateAsset(ctx, createAsset, objID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, asset)
}

// GetAssetByID 根据ID获取资产
func (c *AssetController) GetAssetByID(ctx *gin.Context) {
	id := ctx.Param("id")

	response, err := c.assetService.GetAssetResponse(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("资产未找到: %s", err.Error())})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// UpdateAsset 更新资产信息
func (c *AssetController) UpdateAsset(ctx *gin.Context) {
	id := ctx.Param("id")
	var update models.AssetUpdate
	if err := ctx.ShouldBindJSON(&update); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从上下文获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}
	objID, _ := primitive.ObjectIDFromHex(userID.(string))

	asset, err := c.assetService.UpdateAsset(ctx, id, update, objID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, asset)
}

// DeleteAsset 删除资产
func (c *AssetController) DeleteAsset(ctx *gin.Context) {
	id := ctx.Param("id")

	err := c.assetService.DeleteAsset(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "资产已删除"})
}

// SearchAssets 搜索资产
func (c *AssetController) SearchAssets(ctx *gin.Context) {
	var searchParams models.AssetSearchParams
	if err := ctx.ShouldBindQuery(&searchParams); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := c.assetService.SearchAssets(ctx, searchParams)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// AddVulnerabilityToAsset 为资产添加漏洞
func (c *AssetController) AddVulnerabilityToAsset(ctx *gin.Context) {
	assetID := ctx.Param("id")
	vulnID := ctx.Param("vulnId")

	// 如果URL中没有vulnId参数，则尝试从请求体获取
	if vulnID == "" {
		var requestBody struct {
			VulnerabilityId string `json:"vulnerabilityId"`
		}
		if err := ctx.ShouldBindJSON(&requestBody); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求体: " + err.Error()})
			return
		}

		if requestBody.VulnerabilityId == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "漏洞ID不能为空"})
			return
		}

		vulnID = requestBody.VulnerabilityId
	}

	// 记录API调用信息，帮助调试
	fmt.Printf("AddVulnerabilityToAsset - 资产ID: %s, 漏洞ID: %s\n", assetID, vulnID)

	err := c.assetService.AddVulnerabilityToAsset(ctx, assetID, vulnID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "漏洞已添加到资产"})
}

// RemoveVulnerabilityFromAsset 从资产中移除漏洞
func (c *AssetController) RemoveVulnerabilityFromAsset(ctx *gin.Context) {
	assetID := ctx.Param("id")
	vulnID := ctx.Param("vulnId")

	err := c.assetService.RemoveVulnerabilityFromAsset(ctx, assetID, vulnID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "漏洞已从资产中移除"})
}

// GetAssetVulnerabilities 获取资产关联的所有漏洞
func (c *AssetController) GetAssetVulnerabilities(ctx *gin.Context) {
	assetID := ctx.Param("id")

	vulnerabilities, err := c.assetService.GetAssetVulnerabilities(ctx, assetID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"vulnerabilities": vulnerabilities})
}

// AddAssetNote 为资产添加备注
func (c *AssetController) AddAssetNote(ctx *gin.Context) {
	assetID := ctx.Param("id")

	var noteData struct {
		Content string `json:"content" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&noteData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从上下文获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}
	objID, _ := primitive.ObjectIDFromHex(userID.(string))

	err := c.assetService.AddAssetNote(ctx, assetID, noteData.Content, objID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "备注已添加"})
}

// ImportAssets 批量导入资产
func (c *AssetController) ImportAssets(ctx *gin.Context) {
	// 获取当前用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	// 转换用户ID
	createdByID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "无效的用户ID", "details": err.Error()})
		return
	}

	// 解析multipart表单
	if err := ctx.Request.ParseMultipartForm(10 << 20); err != nil { // 10 MB 限制
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "解析表单数据失败", "details": err.Error()})
		return
	}

	// 获取导入选项
	optionsStr := ctx.Request.FormValue("options")
	var options models.AssetImportOptions
	if optionsStr != "" {
		if err := json.Unmarshal([]byte(optionsStr), &options); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "解析导入选项失败", "details": err.Error()})
			return
		}
	} else {
		// 设置默认选项
		options = models.AssetImportOptions{
			DuplicateStrategy: "skip",
			DefaultType:       "server",
			DefaultStatus:     "active",
			SendNotifications: true,
		}
	}

	// 获取上传的文件
	file, fileHeader, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "获取上传文件失败", "details": err.Error()})
		return
	}
	defer file.Close()

	// 导入结果统计
	stats := models.AssetImportStats{}

	// 根据文件类型处理
	fileExt := strings.ToLower(fileHeader.Filename[strings.LastIndex(fileHeader.Filename, ".")+1:])

	// 读取CSV文件
	if fileExt == "csv" {
		reader := csv.NewReader(file)

		// 读取CSV头部
		headers, err := reader.Read()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "读取CSV头部失败: " + err.Error()})
			return
		}

		// 将标题转换为小写并创建索引映射
		headerMap := make(map[string]int)
		for i, header := range headers {
			// 移除引号并清理标题
			cleanHeader := strings.ToLower(strings.Trim(header, "\"' "))
			headerMap[cleanHeader] = i
		}

		// 必需字段检查
		requiredFields := []string{"name", "type", "status"}
		for _, field := range requiredFields {
			if _, ok := headerMap[field]; !ok {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": fmt.Sprintf("CSV文件缺少必需的列: %s", field),
				})
				return
			}
		}

		// 逐行处理数据
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				stats.Failed++
				continue
			}

			// 创建资产对象
			asset := models.Asset{
				ID:              primitive.NewObjectID(),
				CreatedBy:       createdByID,
				UpdatedBy:       createdByID,
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
				Vulnerabilities: []primitive.ObjectID{},
				Tags:            []string{},
				Notes:           []models.AssetNote{},
				CustomFields:    make(map[string]interface{}),
			}

			// 设置基本字段
			if idx, ok := headerMap["name"]; ok && idx < len(record) {
				asset.Name = strings.Trim(record[idx], "\"' ")
			}

			if idx, ok := headerMap["description"]; ok && idx < len(record) {
				asset.Description = strings.Trim(record[idx], "\"' ")
			}

			if idx, ok := headerMap["type"]; ok && idx < len(record) {
				assetType := strings.Trim(record[idx], "\"' ")
				if assetType == "" {
					assetType = options.DefaultType
				}
				asset.Type = assetType
			} else {
				asset.Type = options.DefaultType
			}

			if idx, ok := headerMap["status"]; ok && idx < len(record) {
				status := strings.Trim(record[idx], "\"' ")
				if status == "" {
					status = options.DefaultStatus
				}
				asset.Status = status
			} else {
				asset.Status = options.DefaultStatus
			}

			// 设置其他可选字段
			if idx, ok := headerMap["ipaddress"]; ok && idx < len(record) {
				asset.IPAddress = strings.Trim(record[idx], "\"' ")
			}

			if idx, ok := headerMap["macaddress"]; ok && idx < len(record) {
				asset.MACAddress = strings.Trim(record[idx], "\"' ")
			}

			if idx, ok := headerMap["location"]; ok && idx < len(record) {
				asset.Location = strings.Trim(record[idx], "\"' ")
			}

			if idx, ok := headerMap["owner"]; ok && idx < len(record) {
				asset.Owner = strings.Trim(record[idx], "\"' ")
			}

			if idx, ok := headerMap["department"]; ok && idx < len(record) {
				asset.Department = strings.Trim(record[idx], "\"' ")
			}

			if idx, ok := headerMap["purchasedate"]; ok && idx < len(record) {
				dateStr := strings.Trim(record[idx], "\"' ")
				if dateStr != "" {
					date, err := time.Parse("2006-01-02", dateStr)
					if err == nil {
						asset.PurchaseDate = &date
					}
				}
			}

			if idx, ok := headerMap["expirydate"]; ok && idx < len(record) {
				dateStr := strings.Trim(record[idx], "\"' ")
				if dateStr != "" {
					date, err := time.Parse("2006-01-02", dateStr)
					if err == nil {
						asset.ExpiryDate = &date
					}
				}
			}

			if idx, ok := headerMap["os"]; ok && idx < len(record) {
				asset.OS = strings.Trim(record[idx], "\"' ")
			}

			if idx, ok := headerMap["osversion"]; ok && idx < len(record) {
				asset.OSVersion = strings.Trim(record[idx], "\"' ")
			}

			if idx, ok := headerMap["manufacturer"]; ok && idx < len(record) {
				asset.Manufacturer = strings.Trim(record[idx], "\"' ")
			}

			if idx, ok := headerMap["model"]; ok && idx < len(record) {
				asset.Model = strings.Trim(record[idx], "\"' ")
			}

			if idx, ok := headerMap["serialnumber"]; ok && idx < len(record) {
				asset.SerialNumber = strings.Trim(record[idx], "\"' ")
			}

			if idx, ok := headerMap["tags"]; ok && idx < len(record) {
				tagsStr := strings.Trim(record[idx], "\"' ")
				if tagsStr != "" {
					tags := strings.Split(tagsStr, ",")
					for i, tag := range tags {
						tags[i] = strings.TrimSpace(tag)
					}
					asset.Tags = tags
				}
			}

			// 检查是否已存在相同名称的资产
			dbCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			count, err := c.assetService.CountAssetsByName(dbCtx, asset.Name)
			if err != nil {
				stats.Failed++
				continue
			}

			if count > 0 {
				// 处理重复项
				if options.DuplicateStrategy == "skip" {
					stats.Duplicates++
					continue
				} else if options.DuplicateStrategy == "update" {
					// 更新已有资产
					existingAsset, err := c.assetService.GetAssetByName(dbCtx, asset.Name)
					if err != nil {
						stats.Failed++
						continue
					}

					// 保留原资产ID和漏洞关联
					asset.ID = existingAsset.ID
					asset.Vulnerabilities = existingAsset.Vulnerabilities

					// 使用AssetService更新资产
					updateData := models.AssetUpdate{
						Name:         &asset.Name,
						Description:  &asset.Description,
						Type:         &asset.Type,
						Status:       &asset.Status,
						IPAddress:    &asset.IPAddress,
						MACAddress:   &asset.MACAddress,
						Location:     &asset.Location,
						Owner:        &asset.Owner,
						Department:   &asset.Department,
						PurchaseDate: asset.PurchaseDate,
						ExpiryDate:   asset.ExpiryDate,
						OS:           &asset.OS,
						OSVersion:    &asset.OSVersion,
						Manufacturer: &asset.Manufacturer,
						Model:        &asset.Model,
						SerialNumber: &asset.SerialNumber,
						Tags:         asset.Tags,
					}

					// 将字符串ID转换为ObjectID
					assetIDStr := existingAsset.ID.Hex()
					_, err = c.assetService.UpdateAsset(dbCtx, assetIDStr, updateData, createdByID)
					if err != nil {
						stats.Failed++
						continue
					}

					stats.Successful++
					continue
				}
				// 默认策略 "create_new" - 直接继续后面的创建流程
			}

			// 创建新资产
			_, err = c.assetService.CreateAsset(dbCtx, models.AssetCreate{
				Name:         asset.Name,
				Description:  asset.Description,
				Type:         asset.Type,
				Status:       asset.Status,
				IPAddress:    asset.IPAddress,
				MACAddress:   asset.MACAddress,
				Location:     asset.Location,
				Owner:        asset.Owner,
				Department:   asset.Department,
				PurchaseDate: asset.PurchaseDate,
				ExpiryDate:   asset.ExpiryDate,
				OS:           asset.OS,
				OSVersion:    asset.OSVersion,
				Manufacturer: asset.Manufacturer,
				Model:        asset.Model,
				SerialNumber: asset.SerialNumber,
				Tags:         asset.Tags,
			}, createdByID)
			if err != nil {
				stats.Failed++
				continue
			}

			stats.Successful++
		}

		// 返回导入结果
		ctx.JSON(http.StatusOK, models.AssetImportResponse{
			Message: "资产导入完成",
			Stats:   stats,
		})
		return
	} else if fileExt == "json" {
		// 处理JSON文件
		var assets []struct {
			Name         string `json:"name"`
			Description  string `json:"description"`
			Type         string `json:"type"`
			Status       string `json:"status"`
			IPAddress    string `json:"ipAddress"`
			MACAddress   string `json:"macAddress"`
			Location     string `json:"location"`
			Owner        string `json:"owner"`
			Department   string `json:"department"`
			PurchaseDate string `json:"purchaseDate"`
			ExpiryDate   string `json:"expiryDate"`
			OS           string `json:"os"`
			OSVersion    string `json:"osVersion"`
			Manufacturer string `json:"manufacturer"`
			Model        string `json:"model"`
			SerialNumber string `json:"serialNumber"`
			Tags         string `json:"tags"`
		}

		// 解析JSON数据
		jsonData, err := io.ReadAll(file)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "读取JSON文件失败: " + err.Error()})
			return
		}

		if err := json.Unmarshal(jsonData, &assets); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "解析JSON数据失败: " + err.Error()})
			return
		}

		// 处理每个资产数据
		for _, assetData := range assets {
			// 创建资产对象
			asset := models.Asset{
				ID:              primitive.NewObjectID(),
				CreatedBy:       createdByID,
				UpdatedBy:       createdByID,
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
				Vulnerabilities: []primitive.ObjectID{},
				Notes:           []models.AssetNote{},
				CustomFields:    make(map[string]interface{}),
			}

			// 设置基本字段
			asset.Name = assetData.Name
			asset.Description = assetData.Description

			// 处理类型字段
			if assetData.Type != "" {
				asset.Type = assetData.Type
			} else {
				asset.Type = options.DefaultType
			}

			// 处理状态字段
			if assetData.Status != "" {
				asset.Status = assetData.Status
			} else {
				asset.Status = options.DefaultStatus
			}

			// 设置其他字段
			asset.IPAddress = assetData.IPAddress
			asset.MACAddress = assetData.MACAddress
			asset.Location = assetData.Location
			asset.Owner = assetData.Owner
			asset.Department = assetData.Department

			// 处理日期字段
			if assetData.PurchaseDate != "" {
				date, err := time.Parse("2006-01-02", assetData.PurchaseDate)
				if err == nil {
					asset.PurchaseDate = &date
				}
			}

			if assetData.ExpiryDate != "" {
				date, err := time.Parse("2006-01-02", assetData.ExpiryDate)
				if err == nil {
					asset.ExpiryDate = &date
				}
			}

			asset.OS = assetData.OS
			asset.OSVersion = assetData.OSVersion
			asset.Manufacturer = assetData.Manufacturer
			asset.Model = assetData.Model
			asset.SerialNumber = assetData.SerialNumber

			// 处理标签
			if assetData.Tags != "" {
				tags := strings.Split(assetData.Tags, ",")
				for i, tag := range tags {
					tags[i] = strings.TrimSpace(tag)
				}
				asset.Tags = tags
			} else {
				asset.Tags = []string{}
			}

			// 检查是否存在重复名称的资产
			dbCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			count, err := c.assetService.CountAssetsByName(dbCtx, asset.Name)
			if err != nil {
				stats.Failed++
				continue
			}

			if count > 0 {
				// 处理重复项
				if options.DuplicateStrategy == "skip" {
					stats.Duplicates++
					continue
				} else if options.DuplicateStrategy == "update" {
					// 更新已有资产
					existingAsset, err := c.assetService.GetAssetByName(dbCtx, asset.Name)
					if err != nil {
						stats.Failed++
						continue
					}

					// 保留原资产ID和漏洞关联
					asset.ID = existingAsset.ID
					asset.Vulnerabilities = existingAsset.Vulnerabilities

					// 使用AssetService更新资产
					updateData := models.AssetUpdate{
						Name:         &asset.Name,
						Description:  &asset.Description,
						Type:         &asset.Type,
						Status:       &asset.Status,
						IPAddress:    &asset.IPAddress,
						MACAddress:   &asset.MACAddress,
						Location:     &asset.Location,
						Owner:        &asset.Owner,
						Department:   &asset.Department,
						PurchaseDate: asset.PurchaseDate,
						ExpiryDate:   asset.ExpiryDate,
						OS:           &asset.OS,
						OSVersion:    &asset.OSVersion,
						Manufacturer: &asset.Manufacturer,
						Model:        &asset.Model,
						SerialNumber: &asset.SerialNumber,
						Tags:         asset.Tags,
					}

					// 将字符串ID转换为ObjectID
					assetIDStr := existingAsset.ID.Hex()
					_, err = c.assetService.UpdateAsset(dbCtx, assetIDStr, updateData, createdByID)
					if err != nil {
						stats.Failed++
						continue
					}

					stats.Successful++
					continue
				}
				// 默认策略 "create_new" - 直接继续后面的创建流程
			}

			// 创建新资产
			_, err = c.assetService.CreateAsset(dbCtx, models.AssetCreate{
				Name:         asset.Name,
				Description:  asset.Description,
				Type:         asset.Type,
				Status:       asset.Status,
				IPAddress:    asset.IPAddress,
				MACAddress:   asset.MACAddress,
				Location:     asset.Location,
				Owner:        asset.Owner,
				Department:   asset.Department,
				PurchaseDate: asset.PurchaseDate,
				ExpiryDate:   asset.ExpiryDate,
				OS:           asset.OS,
				OSVersion:    asset.OSVersion,
				Manufacturer: asset.Manufacturer,
				Model:        asset.Model,
				SerialNumber: asset.SerialNumber,
				Tags:         asset.Tags,
			}, createdByID)
			if err != nil {
				stats.Failed++
				continue
			}

			stats.Successful++
		}

		// 返回导入结果
		ctx.JSON(http.StatusOK, models.AssetImportResponse{
			Message: "资产导入完成",
			Stats:   stats,
		})
		return
	} else if fileExt == "xlsx" || fileExt == "xls" {
		// 处理Excel文件
		// 注意: 完整实现需要引入Excel处理库如 "github.com/360EntSecGroup-Skylar/excelize"
		// 这里提供结构化的伪代码
		ctx.JSON(http.StatusOK, models.AssetImportResponse{
			Message: "Excel格式导入功能即将推出，请使用CSV或JSON格式",
			Stats: models.AssetImportStats{
				Successful: 0,
				Duplicates: 0,
				Failed:     0,
			},
		})
		return
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "不支持的文件格式"})
		return
	}
}

// GetImportTemplate 获取资产导入模板
func (c *AssetController) GetImportTemplate(ctx *gin.Context) {
	// 创建CSV模板
	var buffer bytes.Buffer
	writer := csv.NewWriter(&buffer)

	// 写入CSV标题行
	headers := []string{
		"name", "description", "type", "status",
		"ipAddress", "macAddress", "location", "owner",
		"department", "purchaseDate", "expiryDate",
		"os", "osVersion", "manufacturer", "model",
		"serialNumber", "tags",
	}
	writer.Write(headers)

	// 写入示例数据
	examples := [][]string{
		{
			"Web服务器01", "主要生产环境Web服务器", "server", "active",
			"192.168.1.10", "00:1A:2B:3C:4D:5E", "北京IDC", "张三",
			"技术部", "2022-01-15", "2025-01-15",
			"Linux", "Ubuntu 22.04 LTS", "Dell", "PowerEdge R740",
			"SN12345678", "web,production",
		},
		{
			"数据库服务器01", "核心数据库服务器", "database", "active",
			"192.168.1.20", "00:1A:2B:3C:4D:5F", "上海IDC", "李四",
			"运维部", "2022-03-20", "2025-03-20",
			"Linux", "CentOS 8", "HP", "ProLiant DL380",
			"SN87654321", "database,production",
		},
		{
			"办公电脑A-101", "市场部办公电脑", "workstation", "active",
			"192.168.10.101", "00:1B:44:33:22:11", "总部5楼", "王五",
			"市场部", "2023-05-10", "2026-05-10",
			"Windows", "Windows 11 Pro", "Lenovo", "ThinkPad X1",
			"LT123456", "office,marketing",
		},
	}

	for _, example := range examples {
		writer.Write(example)
	}

	writer.Flush()

	// 设置响应头
	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Disposition", "attachment; filename=资产导入模板.csv")
	ctx.Header("Content-Type", "text/csv; charset=utf-8")
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Expires", "0")
	ctx.Header("Pragma", "public")

	// 返回文件内容
	ctx.Data(http.StatusOK, "text/csv", buffer.Bytes())
}
