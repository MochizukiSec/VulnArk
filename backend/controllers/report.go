package controllers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"vuln-management/config"
	"vuln-management/models"
	"vuln-management/services"
	"vuln-management/utils"
)

// ReportController 报告相关控制器
type ReportController struct {
	reportService services.ReportService
}

// NewReportController 创建报告控制器
func NewReportController(reportService services.ReportService) *ReportController {
	return &ReportController{
		reportService: reportService,
	}
}

// RegisterRoutes 注册路由
func (c *ReportController) RegisterRoutes(router *gin.RouterGroup) {
	reports := router.Group("/reports")
	{
		reports.POST("", c.CreateReport)
		reports.GET("", c.GetAllReports)
		reports.GET("/:id", c.GetReportByID)
		reports.DELETE("/:id", c.DeleteReport)
		reports.GET("/download/:filename", c.DownloadReport)
	}
}

// GetAllReports 获取所有报告
func (c *ReportController) GetAllReports(ctx *gin.Context) {
	// 获取查询参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	reportType := ctx.Query("type")
	format := ctx.Query("format")
	search := ctx.Query("search")

	// 设置查询选项
	findOptions := options.Find()
	findOptions.SetSkip(int64((page - 1) * limit))
	findOptions.SetLimit(int64(limit))
	findOptions.SetSort(bson.M{"created_at": -1}) // 默认按创建时间降序排序

	// 构建查询条件
	filter := bson.M{}
	if reportType != "" {
		filter["type"] = reportType
	}
	if format != "" {
		filter["format"] = format
	}
	if search != "" {
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": search, "$options": "i"}},
			{"description": bson.M{"$regex": search, "$options": "i"}},
		}
	}

	// 获取报告
	reportsCollection := config.GetCollection(config.ReportsCollection)
	ctx2, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 查询总数
	total, err := reportsCollection.CountDocuments(ctx2, filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取报告总数失败"})
		return
	}

	// 查询报告列表
	cursor, err := reportsCollection.Find(ctx2, filter, findOptions)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取报告列表失败"})
		return
	}
	defer cursor.Close(ctx2)

	// 解析报告
	var reports []models.Report
	if err := cursor.All(ctx2, &reports); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "解析报告数据失败"})
		return
	}

	// 计算总页数
	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	// 返回数据
	ctx.JSON(http.StatusOK, gin.H{
		"reports": reports,
		"meta": gin.H{
			"total":     total,
			"page":      page,
			"pages":     totalPages,
			"per_page":  limit,
			"has_next":  page < totalPages,
			"has_prev":  page > 1,
			"next_page": min(page+1, totalPages),
			"prev_page": max(page-1, 1),
		},
	})
}

// GetReportByID 根据ID获取报告
func (c *ReportController) GetReportByID(ctx *gin.Context) {
	reportID, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的报告ID"})
		return
	}

	// 查询报告
	reportsCollection := config.GetCollection(config.ReportsCollection)
	ctx2, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var report models.Report
	err = reportsCollection.FindOne(ctx2, bson.M{"_id": reportID}).Decode(&report)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "报告不存在"})
		return
	}

	ctx.JSON(http.StatusOK, report)
}

// CreateReport 创建报告
func (c *ReportController) CreateReport(ctx *gin.Context) {
	var req models.CreateReportRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 获取当前用户
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	// 创建报告目录（如果不存在）
	if err := ensureReportDirectory(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "创建报告目录失败"})
		return
	}

	// 创建报告对象
	now := time.Now()
	report := models.Report{
		ID:          primitive.NewObjectID(),
		Name:        req.Name,
		Type:        req.Type,
		Format:      req.Format,
		Status:      models.ReportStatusPending,
		Description: req.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
		CreatedBy: models.ReportUser{
			ID:   userID.(primitive.ObjectID),
			Name: getUserName(ctx),
		},
		Metadata: models.ReportMetadata{
			StartDate:  req.StartDate,
			EndDate:    req.EndDate,
			Severities: req.Severities,
			Statuses:   req.Statuses,
		},
	}

	// 保存到数据库
	reportsCollection := config.GetCollection(config.ReportsCollection)
	ctx2, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := reportsCollection.InsertOne(ctx2, report)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("创建报告失败: %v", err)})
		return
	}

	// 异步生成报告
	go generateReport(report)

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "报告创建成功，正在后台生成",
		"data":    report,
	})
}

// DeleteReport 删除报告
func (c *ReportController) DeleteReport(ctx *gin.Context) {
	reportID, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的报告ID"})
		return
	}

	// 删除报告
	reportsCollection := config.GetCollection(config.ReportsCollection)
	ctx2, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := reportsCollection.DeleteOne(ctx2, bson.M{"_id": reportID})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "删除报告失败: " + err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "报告不存在"})
		return
	}

	// 尝试删除报告文件
	filename := fmt.Sprintf("./reports/%s.txt", reportID.Hex())
	_ = os.Remove(filename) // 忽略错误

	// 返回成功
	ctx.JSON(http.StatusOK, gin.H{"message": "报告已成功删除"})
}

// GenerateSummaryReport 生成摘要报告
func (c *ReportController) GenerateSummaryReport(ctx *gin.Context) {
	// 获取请求参数
	var req models.CreateReportRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		req = models.CreateReportRequest{
			Name:   "漏洞摘要报告",
			Type:   models.ReportTypeSummary,
			Format: models.ReportFormatText,
		}
	}

	// 设置默认值
	if req.Name == "" {
		req.Name = "漏洞摘要报告"
	}
	if req.Type == "" {
		req.Type = models.ReportTypeSummary
	}
	if req.Format == "" {
		req.Format = models.ReportFormatText
	}
	if req.StartDate.IsZero() {
		req.StartDate = time.Now().AddDate(0, -1, 0) // 默认过去1个月
	}
	if req.EndDate.IsZero() {
		req.EndDate = time.Now()
	}

	// 调用创建报告接口
	c.CreateReport(ctx)
}

// GenerateDetailedReport 生成详细报告
func (c *ReportController) GenerateDetailedReport(ctx *gin.Context) {
	// 获取请求参数
	var req models.CreateReportRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		req = models.CreateReportRequest{
			Name:   "漏洞详细报告",
			Type:   models.ReportTypeDetailed,
			Format: models.ReportFormatText,
		}
	}

	// 设置默认值
	if req.Name == "" {
		req.Name = "漏洞详细报告"
	}
	if req.Type == "" {
		req.Type = models.ReportTypeDetailed
	}
	if req.Format == "" {
		req.Format = models.ReportFormatText
	}
	if req.StartDate.IsZero() {
		req.StartDate = time.Now().AddDate(0, -1, 0) // 默认过去1个月
	}
	if req.EndDate.IsZero() {
		req.EndDate = time.Now()
	}

	// 调用创建报告接口
	c.CreateReport(ctx)
}

// DownloadReport 下载报告
func (c *ReportController) DownloadReport(ctx *gin.Context) {
	filename := ctx.Param("filename")
	if filename == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "文件名不能为空"})
		return
	}

	// 安全检查：防止路径遍历攻击
	if filepath.Ext(filename) != ".txt" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "不支持的文件类型"})
		return
	}

	// 构建文件路径
	filePath := filepath.Join("./reports", filename)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "报告文件不存在"})
		return
	}

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "无法打开报告文件"})
		return
	}
	defer file.Close()

	// 设置响应头
	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	ctx.Header("Content-Type", "text/plain")
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Expires", "0")
	ctx.Header("Cache-Control", "must-revalidate")
	ctx.Header("Pragma", "public")

	// 将文件内容复制到响应中
	ctx.Status(http.StatusOK)
	_, err = io.Copy(ctx.Writer, file)
	if err != nil {
		utils.LogError(fmt.Sprintf("下载报告文件失败: %v", err))
	}
}

// 辅助函数

// getUserName 获取用户名
func getUserName(ctx *gin.Context) string {
	userName, exists := ctx.Get("user_name")
	if !exists {
		return "系统用户"
	}
	return userName.(string)
}

// ensureReportDirectory 确保报告目录存在
func ensureReportDirectory() error {
	dir := "./reports"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}

// generateReport 异步生成报告
func generateReport(report models.Report) {
	// 创建报告目录（如果不存在）
	if err := ensureReportDirectory(); err != nil {
		updateReportStatus(report.ID, models.ReportStatusFailed, "")
		return
	}

	// 按照报告类型生成不同的报告
	var err error
	switch report.Type {
	case models.ReportTypeSummary:
		err = generateSummaryReport(report)
	case models.ReportTypeDetailed:
		err = generateDetailedReport(report)
	default:
		err = fmt.Errorf("不支持的报告类型: %s", report.Type)
	}

	if err != nil {
		fmt.Printf("生成报告失败: %v\n", err)
		updateReportStatus(report.ID, models.ReportStatusFailed, "")
		return
	}

	// 生成文件URL
	fileURL := fmt.Sprintf("/api/reports/download/%s.txt", report.ID.Hex())
	updateReportStatus(report.ID, models.ReportStatusCompleted, fileURL)
}

// updateReportStatus 更新报告状态
func updateReportStatus(reportID primitive.ObjectID, status string, fileURL string) {
	reportsCollection := config.GetCollection(config.ReportsCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"file_url":   fileURL,
			"updated_at": time.Now(),
		},
	}

	_, err := reportsCollection.UpdateOne(ctx, bson.M{"_id": reportID}, update)
	if err != nil {
		fmt.Printf("更新报告状态失败: %v\n", err)
	}
}

// generateSummaryReport 生成摘要报告
func generateSummaryReport(report models.Report) error {
	// 使用简单文本文件生成报告
	filePath := fmt.Sprintf("./reports/%s.txt", report.ID.Hex())
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("创建报告文件失败: %w", err)
	}
	defer file.Close()

	// 写入报告内容
	content := []string{
		"==========================",
		"    漏洞管理系统 - 摘要报告    ",
		"==========================\n",
		fmt.Sprintf("报告名称: %s", report.Name),
		fmt.Sprintf("生成时间: %s", time.Now().Format("2006-01-02 15:04:05")),
		fmt.Sprintf("时间范围: %s 至 %s",
			report.Metadata.StartDate.Format("2006-01-02"),
			report.Metadata.EndDate.Format("2006-01-02")),
		"==========================\n",
		"漏洞统计信息:",
	}

	// 查询漏洞统计数据
	vulnsCollection := config.GetCollection("vulnerabilities")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 构建查询条件
	query := bson.M{}
	if !report.Metadata.StartDate.IsZero() && !report.Metadata.EndDate.IsZero() {
		query["created_at"] = bson.M{
			"$gte": report.Metadata.StartDate,
			"$lte": report.Metadata.EndDate,
		}
	}

	// 计算漏洞总数
	totalVulns, err := vulnsCollection.CountDocuments(ctx, query)
	if err != nil {
		return fmt.Errorf("计算漏洞总数失败: %w", err)
	}

	// 添加统计数据
	content = append(content, fmt.Sprintf("漏洞总数: %d", totalVulns))

	// 新增漏洞数量（过去7天）
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)
	newVulnQuery := bson.M{"created_at": bson.M{"$gte": sevenDaysAgo}}
	newVulns, _ := vulnsCollection.CountDocuments(ctx, newVulnQuery)
	content = append(content, fmt.Sprintf("新增漏洞(7天): %d", newVulns))

	// 写入文件
	for _, line := range content {
		if _, err := file.WriteString(line + "\n"); err != nil {
			return fmt.Errorf("写入报告文件失败: %w", err)
		}
	}

	return nil
}

// generateDetailedReport 生成详细报告
func generateDetailedReport(report models.Report) error {
	// 使用简单文本文件生成报告
	filePath := fmt.Sprintf("./reports/%s.txt", report.ID.Hex())
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("创建报告文件失败: %w", err)
	}
	defer file.Close()

	// 写入报告内容
	content := []string{
		"==========================",
		"    漏洞管理系统 - 详细报告    ",
		"==========================\n",
		fmt.Sprintf("报告名称: %s", report.Name),
		fmt.Sprintf("生成时间: %s", time.Now().Format("2006-01-02 15:04:05")),
		fmt.Sprintf("时间范围: %s 至 %s",
			report.Metadata.StartDate.Format("2006-01-02"),
			report.Metadata.EndDate.Format("2006-01-02")),
		"==========================\n",
		"漏洞列表:",
		"ID | 名称 | 严重程度 | 状态 | CVE | 发现日期",
		"----------------------------------------",
	}

	// 查询漏洞数据
	vulnsCollection := config.GetCollection("vulnerabilities")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 构建查询条件
	query := bson.M{}
	if !report.Metadata.StartDate.IsZero() && !report.Metadata.EndDate.IsZero() {
		query["created_at"] = bson.M{
			"$gte": report.Metadata.StartDate,
			"$lte": report.Metadata.EndDate,
		}
	}

	if len(report.Metadata.Severities) > 0 {
		query["severity"] = bson.M{"$in": report.Metadata.Severities}
	}

	if len(report.Metadata.Statuses) > 0 {
		query["status"] = bson.M{"$in": report.Metadata.Statuses}
	}

	cursor, err := vulnsCollection.Find(ctx, query)
	if err != nil {
		return fmt.Errorf("查询漏洞数据失败: %w", err)
	}
	defer cursor.Close(ctx)

	// 填充漏洞数据
	var vulnerabilities []models.Vulnerability
	if err := cursor.All(ctx, &vulnerabilities); err != nil {
		return fmt.Errorf("解析漏洞数据失败: %w", err)
	}

	for _, vuln := range vulnerabilities {
		vulnLine := fmt.Sprintf("%s | %s | %s | %s | %s | %s",
			vuln.ID.Hex(),
			vuln.Title,
			vuln.Severity,
			vuln.Status,
			vuln.CVE,
			vuln.CreatedAt.Format("2006-01-02"))
		content = append(content, vulnLine)
	}

	// 写入文件
	for _, line := range content {
		if _, err := file.WriteString(line + "\n"); err != nil {
			return fmt.Errorf("写入报告文件失败: %w", err)
		}
	}

	return nil
}

// 工具函数
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
