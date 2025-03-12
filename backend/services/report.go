package services

import (
	"context"
	"errors"
	"fmt"
	"time"
	"vuln-management/models"

	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ReportService 报告服务接口
type ReportService interface {
	CreateReport(ctx context.Context, req models.CreateReportRequest, userID primitive.ObjectID, userName string) (*models.Report, error)
	GetAllReports(ctx context.Context, page, pageSize int) ([]models.Report, int64, error)
	GetReportByID(ctx context.Context, id string) (*models.Report, error)
	DeleteReport(ctx context.Context, id string) error
	GenerateSummaryReport(ctx context.Context, reportID string, metadata models.ReportMetadata) error
	GenerateDetailedReport(ctx context.Context, reportID string, metadata models.ReportMetadata) error
}

// MongoReportService MongoDB报告服务实现
type MongoReportService struct {
	reportCollection *mongo.Collection
	vulnService      *VulnerabilityService
	assetService     *AssetService
	userService      *UserService
	database         *mongo.Database
}

// NewMongoReportService 创建新的MongoDB报告服务
func NewMongoReportService(db *mongo.Database, vulnService *VulnerabilityService, assetService *AssetService, userService *UserService) ReportService {
	return &MongoReportService{
		reportCollection: db.Collection("reports"),
		vulnService:      vulnService,
		assetService:     assetService,
		userService:      userService,
		database:         db,
	}
}

// NewReportService 创建报告服务（为了兼容旧代码）
func NewReportService(db *mongo.Database) ReportService {
	vulnService := NewVulnerabilityService(db)
	assetService := NewAssetService(db)
	userService := NewUserService(db)
	return &MongoReportService{
		reportCollection: db.Collection("reports"),
		vulnService:      vulnService,
		assetService:     assetService,
		userService:      userService,
		database:         db,
	}
}

// CreateReport 创建报告
func (s *MongoReportService) CreateReport(ctx context.Context, req models.CreateReportRequest, userID primitive.ObjectID, userName string) (*models.Report, error) {
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
			ID:   userID,
			Name: userName,
		},
		Metadata: models.ReportMetadata{
			StartDate:  req.StartDate,
			EndDate:    req.EndDate,
			Severities: req.Severities,
			Statuses:   req.Statuses,
		},
	}

	_, err := s.reportCollection.InsertOne(ctx, report)
	if err != nil {
		return nil, fmt.Errorf("创建报告失败: %w", err)
	}

	// 异步生成报告
	go func() {
		// 创建新的context，因为原始的ctx会被取消
		bgCtx := context.Background()
		var genErr error

		switch report.Type {
		case models.ReportTypeSummary:
			genErr = s.GenerateSummaryReport(bgCtx, report.ID.Hex(), report.Metadata)
		case models.ReportTypeDetailed:
			genErr = s.GenerateDetailedReport(bgCtx, report.ID.Hex(), report.Metadata)
		// 可以添加其他报告类型的处理
		default:
			genErr = errors.New("不支持的报告类型")
		}

		// 更新报告状态
		updateData := bson.M{
			"$set": bson.M{
				"status":     models.ReportStatusCompleted,
				"updated_at": time.Now(),
			},
		}

		if genErr != nil {
			updateData = bson.M{
				"$set": bson.M{
					"status":     models.ReportStatusFailed,
					"updated_at": time.Now(),
				},
			}
		}

		_, _ = s.reportCollection.UpdateByID(bgCtx, report.ID, updateData)
	}()

	return &report, nil
}

// GetAllReports 获取所有报告
func (s *MongoReportService) GetAllReports(ctx context.Context, page, pageSize int) ([]models.Report, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	skip := (page - 1) * pageSize
	limit := int64(pageSize)

	options := options.Find().
		SetSkip(int64(skip)).
		SetLimit(limit).
		SetSort(bson.M{"created_at": -1})

	cursor, err := s.reportCollection.Find(ctx, bson.M{}, options)
	if err != nil {
		return nil, 0, fmt.Errorf("查询报告失败: %w", err)
	}
	defer cursor.Close(ctx)

	var reports []models.Report
	if err := cursor.All(ctx, &reports); err != nil {
		return nil, 0, fmt.Errorf("解析报告数据失败: %w", err)
	}

	// 获取总数
	total, err := s.reportCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, fmt.Errorf("计算报告总数失败: %w", err)
	}

	return reports, total, nil
}

// GetReportByID 根据ID获取报告
func (s *MongoReportService) GetReportByID(ctx context.Context, id string) (*models.Report, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("无效的报告ID: %w", err)
	}

	var report models.Report
	err = s.reportCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&report)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("报告不存在")
		}
		return nil, fmt.Errorf("查询报告失败: %w", err)
	}

	return &report, nil
}

// DeleteReport 删除报告
func (s *MongoReportService) DeleteReport(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("无效的报告ID: %w", err)
	}

	res, err := s.reportCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return fmt.Errorf("删除报告失败: %w", err)
	}

	if res.DeletedCount == 0 {
		return fmt.Errorf("报告不存在")
	}

	return nil
}

// GenerateSummaryReport 生成摘要报告
func (s *MongoReportService) GenerateSummaryReport(ctx context.Context, reportID string, metadata models.ReportMetadata) error {
	objID, err := primitive.ObjectIDFromHex(reportID)
	if err != nil {
		return fmt.Errorf("无效的报告ID: %w", err)
	}

	// 构建查询条件
	query := bson.M{}
	if !metadata.StartDate.IsZero() && !metadata.EndDate.IsZero() {
		query["created_at"] = bson.M{
			"$gte": metadata.StartDate,
			"$lte": metadata.EndDate,
		}
	}

	if len(metadata.Severities) > 0 {
		query["severity"] = bson.M{"$in": metadata.Severities}
	}

	if len(metadata.Statuses) > 0 {
		query["status"] = bson.M{"$in": metadata.Statuses}
	}

	// 1. 计算漏洞总数
	totalVulns, err := s.database.Collection("vulnerabilities").CountDocuments(ctx, query)
	if err != nil {
		return fmt.Errorf("计算漏洞总数失败: %w", err)
	}

	// 2. 计算新增漏洞数量（过去7天内）
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)
	newVulnQuery := bson.M{"created_at": bson.M{"$gte": sevenDaysAgo}}
	for k, v := range query {
		newVulnQuery[k] = v
	}
	newVulns, err := s.database.Collection("vulnerabilities").CountDocuments(ctx, newVulnQuery)
	if err != nil {
		return fmt.Errorf("计算新增漏洞数失败: %w", err)
	}

	// 3. 计算已解决漏洞数量
	resolvedQuery := bson.M{"status": "resolved"}
	for k, v := range query {
		if k != "status" { // 避免覆盖status条件
			resolvedQuery[k] = v
		}
	}
	resolvedVulns, err := s.database.Collection("vulnerabilities").CountDocuments(ctx, resolvedQuery)
	if err != nil {
		return fmt.Errorf("计算已解决漏洞数失败: %w", err)
	}

	// 生成Excel报告
	if err := s.generateExcelReport(reportID, int(totalVulns), int(newVulns), int(resolvedVulns), metadata); err != nil {
		return fmt.Errorf("生成Excel报告失败: %w", err)
	}

	// 更新报告状态和文件URL
	fileURL := fmt.Sprintf("/api/reports/download/%s.xlsx", reportID)
	_, err = s.reportCollection.UpdateByID(ctx, objID, bson.M{
		"$set": bson.M{
			"file_url":   fileURL,
			"updated_at": time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("更新报告状态失败: %w", err)
	}

	return nil
}

// GenerateDetailedReport 生成详细报告
func (s *MongoReportService) GenerateDetailedReport(ctx context.Context, reportID string, metadata models.ReportMetadata) error {
	// 实现详细报告的生成逻辑
	// 类似于GenerateSummaryReport，但包含更多详细信息

	// 此处省略具体实现，可以参考摘要报告的实现方式

	return nil
}

// generateExcelReport 生成Excel格式报告
func (s *MongoReportService) generateExcelReport(reportID string, totalVulns, newVulns, resolvedVulns int, metadata models.ReportMetadata) error {
	f := excelize.NewFile()

	// 创建摘要工作表
	sheetName := "漏洞摘要"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return err
	}

	// 设置标题和样式
	f.SetCellValue(sheetName, "A1", "漏洞管理摘要报告")
	f.SetCellValue(sheetName, "A3", "统计时间范围:")
	f.SetCellValue(sheetName, "B3", fmt.Sprintf("%s 至 %s",
		metadata.StartDate.Format("2006-01-02"),
		metadata.EndDate.Format("2006-01-02")))

	// 添加统计数据
	f.SetCellValue(sheetName, "A5", "漏洞总数:")
	f.SetCellValue(sheetName, "B5", totalVulns)

	f.SetCellValue(sheetName, "A6", "新增漏洞:")
	f.SetCellValue(sheetName, "B6", newVulns)

	f.SetCellValue(sheetName, "A7", "已解决漏洞:")
	f.SetCellValue(sheetName, "B7", resolvedVulns)

	// 设置默认工作表
	f.SetActiveSheet(index)

	// 保存Excel文件
	if err := f.SaveAs(fmt.Sprintf("./reports/%s.xlsx", reportID)); err != nil {
		return err
	}

	return nil
}
