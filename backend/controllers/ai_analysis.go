package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"vuln-management/models"
	"vuln-management/services"
)

// AIAnalysisController 处理AI分析相关的请求
type AIAnalysisController struct {
	aiAnalysisService *services.AIAnalysisService
}

// NewAIAnalysisController 创建新的AI分析控制器
func NewAIAnalysisController(aiAnalysisService *services.AIAnalysisService) *AIAnalysisController {
	return &AIAnalysisController{
		aiAnalysisService: aiAnalysisService,
	}
}

// RunAnalysis 执行AI分析
// @Summary 执行AI分析
// @Description 根据请求参数执行相应类型的AI分析
// @Tags AI分析
// @Accept json
// @Produce json
// @Param request body models.AIAnalysisRequest true "分析请求参数"
// @Success 200 {object} models.APIResponse{data=models.AIAnalysisResponse} "成功"
// @Failure 400 {object} models.APIResponse "请求参数错误"
// @Failure 500 {object} models.APIResponse "服务器内部错误"
// @Router /api/ai-analysis/run [post]
func (c *AIAnalysisController) RunAnalysis(ctx *gin.Context) {
	var req models.AIAnalysisRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数无效",
			"error":   err.Error(),
		})
		return
	}

	// 执行分析
	result, err := c.aiAnalysisService.RunAnalysis(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "执行分析失败",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "分析执行成功",
		"data":    result,
	})
}

// GetAnalysis 获取分析结果
// @Summary 获取特定分析结果
// @Description 根据ID获取特定的分析结果
// @Tags AI分析
// @Accept json
// @Produce json
// @Param id path string true "分析ID"
// @Success 200 {object} models.APIResponse{data=models.AIAnalysisResponse} "成功"
// @Failure 400 {object} models.APIResponse "ID格式无效"
// @Failure 404 {object} models.APIResponse "分析结果不存在"
// @Failure 500 {object} models.APIResponse "服务器内部错误"
// @Router /api/ai-analysis/{id} [get]
func (c *AIAnalysisController) GetAnalysis(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的分析ID",
			"error":   err.Error(),
		})
		return
	}

	// 获取分析结果
	result, err := c.aiAnalysisService.GetAnalysisByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取分析结果失败",
			"error":   err.Error(),
		})
		return
	}

	if result == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "分析结果不存在",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// ListAnalyses 列出分析结果
// @Summary 列出分析结果
// @Description 获取分析结果列表，支持分页和过滤
// @Tags AI分析
// @Accept json
// @Produce json
// @Param page query int false "页码，默认1"
// @Param limit query int false "每页数量，默认10"
// @Param type query string false "分析类型过滤"
// @Success 200 {object} models.APIResponse{data=[]models.AIAnalysisResponse} "成功"
// @Failure 500 {object} models.APIResponse "服务器内部错误"
// @Router /api/ai-analysis [get]
func (c *AIAnalysisController) ListAnalyses(ctx *gin.Context) {
	// 获取查询参数
	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "10")
	typeFilter := ctx.Query("type")

	// 解析分页参数
	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil || limit < 1 {
		limit = 10
	}

	// 计算跳过的记录数
	skip := (page - 1) * limit

	// 获取分析结果列表
	results, total, err := c.aiAnalysisService.ListAnalyses(ctx, limit, skip, typeFilter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取分析结果列表失败",
			"error":   err.Error(),
		})
		return
	}

	// 计算总页数
	totalPages := (total + limit - 1) / limit

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"analyses": results,
			"pagination": gin.H{
				"currentPage": page,
				"totalPages":  totalPages,
				"totalItems":  total,
				"limit":       limit,
			},
		},
	})
}

// RegisterRoutes 注册路由
func (c *AIAnalysisController) RegisterRoutes(router *gin.RouterGroup) {
	aiGroup := router.Group("/ai-analysis")
	{
		aiGroup.POST("/run", c.RunAnalysis)
		aiGroup.GET("/:id", c.GetAnalysis)
		aiGroup.GET("", c.ListAnalyses)
	}
}
