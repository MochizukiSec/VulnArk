package controllers

import (
	"fmt"
	"net/http"

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
