package services

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"vuln-management/models"
)

// AIAnalysisService 提供AI分析功能的服务
type AIAnalysisService struct {
	db                      *mongo.Database
	vulnerabilityCollection string
	aiAnalysisCollection    string
}

// NewAIAnalysisService 创建新的AI分析服务
func NewAIAnalysisService(db *mongo.Database) *AIAnalysisService {
	return &AIAnalysisService{
		db:                      db,
		vulnerabilityCollection: "vulnerabilities",
		aiAnalysisCollection:    "ai_analysis",
	}
}

// RunAnalysis 执行AI分析
func (s *AIAnalysisService) RunAnalysis(ctx context.Context, req models.AIAnalysisRequest) (*models.AIAnalysisResponse, error) {
	var (
		analysis models.AIAnalysis
		err      error
	)

	// 根据分析类型执行相应的分析
	switch req.Type {
	case models.TrendPrediction:
		analysis, err = s.predictTrends(ctx, req.Parameters)
	case models.ResourceOptimization:
		analysis, err = s.optimizeResources(ctx, req.Parameters)
	case models.AnomalyDetection:
		analysis, err = s.detectAnomalies(ctx, req.Parameters)
	default:
		return nil, errors.New("不支持的分析类型")
	}

	if err != nil {
		return nil, err
	}

	// 保存分析结果到数据库
	if err := s.saveAnalysis(ctx, &analysis); err != nil {
		return nil, err
	}

	// 返回响应
	return &models.AIAnalysisResponse{
		ID:              analysis.ID,
		Type:            analysis.Type,
		Title:           analysis.Title,
		Description:     analysis.Description,
		Recommendations: analysis.Recommendations,
		Confidence:      analysis.Confidence,
		CreatedAt:       analysis.CreatedAt,
		AnalysisData:    analysis.AnalysisData,
	}, nil
}

// predictTrends 预测漏洞趋势
func (s *AIAnalysisService) predictTrends(ctx context.Context, params map[string]interface{}) (models.AIAnalysis, error) {
	// 获取时间范围参数，默认30天
	timeRange := "未来30天"
	if val, ok := params["timeRange"]; ok {
		if strVal, ok := val.(string); ok {
			timeRange = strVal
		}
	}

	// 从数据库获取历史漏洞数据
	historicalData, err := s.getHistoricalVulnerabilityData(ctx)
	if err != nil {
		return models.AIAnalysis{}, err
	}

	// 使用AI模型预测未来趋势
	// 这里使用简单的模拟算法，实际项目中应替换为真实的ML模型
	prediction := s.simulateTrendPrediction(historicalData, timeRange)

	// 创建分析结果
	now := time.Now()
	analysis := models.AIAnalysis{
		ID:           primitive.NewObjectID(),
		Type:         models.TrendPrediction,
		Title:        fmt.Sprintf("%s漏洞趋势预测", timeRange),
		Description:  fmt.Sprintf("基于历史数据的%s安全漏洞趋势预测分析", timeRange),
		AnalysisData: prediction,
		Recommendations: []string{
			"增加对高风险类别漏洞的监控频率",
			"优先分配资源修复预测增长最快的漏洞类型",
			"为可能增长的漏洞类型提前制定应对策略",
			"关注影响因素中的主要驱动因素，采取预防措施",
		},
		Confidence: 0.85, // 模拟的置信度
		CreatedAt:  now,
		UpdatedAt:  now,
		Parameters: params,
	}

	return analysis, nil
}

// optimizeResources 优化资源分配
func (s *AIAnalysisService) optimizeResources(ctx context.Context, params map[string]interface{}) (models.AIAnalysis, error) {
	// 获取当前资源分配情况，或使用默认值
	currentAllocation := map[string]float64{
		"高危漏洞修复": 30,
		"中危漏洞修复": 25,
		"低危漏洞修复": 20,
		"漏洞扫描":   15,
		"安全培训":   10,
	}

	// 根据历史数据和当前状态计算推荐的资源分配
	// 这里使用简单的模拟算法，实际项目中应替换为真实的优化算法
	optimization := s.simulateResourceOptimization(ctx, currentAllocation)

	// 创建分析结果
	now := time.Now()
	analysis := models.AIAnalysis{
		ID:           primitive.NewObjectID(),
		Type:         models.ResourceOptimization,
		Title:        "安全资源优化配置建议",
		Description:  "基于当前漏洞状态和历史数据的安全资源分配优化建议",
		AnalysisData: optimization,
		Recommendations: []string{
			"将更多资源分配给高危漏洞的修复，以降低整体风险",
			"减少对低危漏洞的资源投入，提高资源利用效率",
			"解决已识别的资源瓶颈，特别是团队协作问题",
			"考虑引入自动化工具，减少手动任务所需的资源",
		},
		Confidence: 0.78, // 模拟的置信度
		CreatedAt:  now,
		UpdatedAt:  now,
		Parameters: params,
	}

	return analysis, nil
}

// detectAnomalies 检测异常
func (s *AIAnalysisService) detectAnomalies(ctx context.Context, params map[string]interface{}) (models.AIAnalysis, error) {
	// 获取时间范围参数，默认7天
	timeRange := "过去7天"
	if val, ok := params["timeRange"]; ok {
		if strVal, ok := val.(string); ok {
			timeRange = strVal
		}
	}

	// 从数据库获取漏洞数据
	vulnerabilityData, err := s.getRecentVulnerabilityData(ctx, timeRange)
	if err != nil {
		return models.AIAnalysis{}, err
	}

	// 使用AI算法检测异常
	// 这里使用简单的模拟算法，实际项目中应替换为真实的异常检测算法
	anomalyResult := s.simulateAnomalyDetection(vulnerabilityData, timeRange)

	// 创建分析结果
	now := time.Now()
	analysis := models.AIAnalysis{
		ID:           primitive.NewObjectID(),
		Type:         models.AnomalyDetection,
		Title:        fmt.Sprintf("%s漏洞异常检测", timeRange),
		Description:  "检测漏洞数据中的异常模式和潜在安全事件",
		AnalysisData: anomalyResult,
		Recommendations: []string{
			"调查检测到的高严重度异常，确认是否存在安全事件",
			"关注特定区域的异常增长，可能表明新的攻击模式",
			"审查异常检测到的漏洞类型，加强相关防御措施",
			"考虑增加受影响区域的监控频率",
		},
		Confidence: 0.82, // 模拟的置信度
		CreatedAt:  now,
		UpdatedAt:  now,
		Parameters: params,
	}

	return analysis, nil
}

// saveAnalysis 保存分析结果到数据库
func (s *AIAnalysisService) saveAnalysis(ctx context.Context, analysis *models.AIAnalysis) error {
	_, err := s.db.Collection(s.aiAnalysisCollection).InsertOne(ctx, analysis)
	return err
}

// GetAnalysisByID 通过ID获取分析结果
func (s *AIAnalysisService) GetAnalysisByID(ctx context.Context, id primitive.ObjectID) (*models.AIAnalysisResponse, error) {
	var analysis models.AIAnalysis
	err := s.db.Collection(s.aiAnalysisCollection).FindOne(ctx, bson.M{"_id": id}).Decode(&analysis)
	if err != nil {
		return nil, err
	}

	return &models.AIAnalysisResponse{
		ID:              analysis.ID,
		Type:            analysis.Type,
		Title:           analysis.Title,
		Description:     analysis.Description,
		Recommendations: analysis.Recommendations,
		Confidence:      analysis.Confidence,
		CreatedAt:       analysis.CreatedAt,
		AnalysisData:    analysis.AnalysisData,
	}, nil
}

// ListAnalyses 列出分析结果
func (s *AIAnalysisService) ListAnalyses(ctx context.Context, limit, skip int64, typeFilter string) ([]models.AIAnalysisResponse, int64, error) {
	// 准备查询条件
	filter := bson.M{}
	if typeFilter != "" {
		filter["type"] = typeFilter
	}

	// 查询选项
	findOptions := options.Find()
	findOptions.SetLimit(limit)
	findOptions.SetSkip(skip)
	findOptions.SetSort(bson.M{"created_at": -1}) // 按创建时间降序

	// 执行查询
	cursor, err := s.db.Collection(s.aiAnalysisCollection).Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	// 获取总数
	total, err := s.db.Collection(s.aiAnalysisCollection).CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// 解析结果
	var analyses []models.AIAnalysisResponse
	for cursor.Next(ctx) {
		var analysis models.AIAnalysis
		if err := cursor.Decode(&analysis); err != nil {
			return nil, 0, err
		}

		analyses = append(analyses, models.AIAnalysisResponse{
			ID:              analysis.ID,
			Type:            analysis.Type,
			Title:           analysis.Title,
			Description:     analysis.Description,
			Recommendations: analysis.Recommendations,
			Confidence:      analysis.Confidence,
			CreatedAt:       analysis.CreatedAt,
			AnalysisData:    analysis.AnalysisData,
		})
	}

	if err := cursor.Err(); err != nil {
		return nil, 0, err
	}

	return analyses, total, nil
}

// ----- 辅助方法 -----

// getHistoricalVulnerabilityData 获取历史漏洞数据
func (s *AIAnalysisService) getHistoricalVulnerabilityData(ctx context.Context) ([]models.HistoricalPoint, error) {
	// 创建时间范围 - 过去90天
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -90)

	// 从数据库聚合数据
	// 实际项目中应从数据库聚合真实数据，这里使用模拟数据
	var historicalData []models.HistoricalPoint

	// 每天生成一个数据点
	currentDate := startDate
	for currentDate.Before(endDate) || currentDate.Equal(endDate) {
		// 模拟数据
		count := 10 + rand.Intn(40) // 随机波动

		// 添加特定的模式
		dayFactor := float64(currentDate.Day()) / 30.0 * 20                     // 月内模式
		weekFactor := math.Sin(float64(currentDate.Weekday())/7.0*math.Pi) * 15 // 周模式
		count = count + int(dayFactor) + int(weekFactor)

		// 确保非负数
		if count < 0 {
			count = 0
		}

		historicalData = append(historicalData, models.HistoricalPoint{
			Date:  currentDate,
			Count: count,
		})

		currentDate = currentDate.AddDate(0, 0, 1) // 下一天
	}

	return historicalData, nil
}

// simulateTrendPrediction 模拟趋势预测
func (s *AIAnalysisService) simulateTrendPrediction(historicalData []models.HistoricalPoint, timeRange string) models.TrendPredictionData {
	// 简单分析历史数据趋势
	// 在实际项目中，这应该是使用机器学习模型进行预测

	// 根据历史数据生成预测
	predictedCounts := map[string]int{
		"高危": 120,
		"中危": 180,
		"低危": 90,
		"信息": 45,
	}

	// 影响因素
	trendFactors := []models.TrendFactor{
		{
			Factor:      "新版本软件发布",
			Impact:      0.25,
			Confidence:  0.85,
			Description: "新发布的软件版本通常会引入新的漏洞",
		},
		{
			Factor:      "安全团队扩展",
			Impact:      -0.15,
			Confidence:  0.78,
			Description: "安全团队扩展可能改善漏洞检测和修复效率",
		},
		{
			Factor:      "第三方组件使用增加",
			Impact:      0.30,
			Confidence:  0.92,
			Description: "更多第三方组件可能引入更多潜在漏洞",
		},
		{
			Factor:      "安全培训计划",
			Impact:      -0.18,
			Confidence:  0.75,
			Description: "安全意识培训可能减少开发引入的漏洞",
		},
	}

	return models.TrendPredictionData{
		TimeRange:       timeRange,
		PredictedCounts: predictedCounts,
		TrendFactors:    trendFactors,
		HistoricalData:  historicalData,
	}
}

// getRecentVulnerabilityData 获取最近的漏洞数据
func (s *AIAnalysisService) getRecentVulnerabilityData(ctx context.Context, timeRange string) ([]map[string]interface{}, error) {
	// 实际项目中，应从数据库获取真实数据
	// 这里返回模拟数据
	var data []map[string]interface{}
	for i := 0; i < 50; i++ {
		data = append(data, map[string]interface{}{
			"id":        primitive.NewObjectID().Hex(),
			"title":     fmt.Sprintf("测试漏洞 %d", i),
			"severity":  []string{"高", "中", "低"}[rand.Intn(3)],
			"status":    []string{"待修复", "已修复", "已确认", "已关闭"}[rand.Intn(4)],
			"createdAt": time.Now().AddDate(0, 0, -rand.Intn(30)),
		})
	}
	return data, nil
}

// simulateAnomalyDetection 模拟异常检测
func (s *AIAnalysisService) simulateAnomalyDetection(data []map[string]interface{}, timeRange string) models.AnomalyDetectionData {
	// 模拟检测到的异常
	anomalies := []models.Anomaly{
		{
			Type:         "突增",
			Severity:     "高",
			Description:  "Web应用关键漏洞数量突然增加",
			DetectedAt:   time.Now().AddDate(0, 0, -2),
			AffectedArea: "Web应用安全",
			Score:        0.92,
		},
		{
			Type:         "异常模式",
			Severity:     "中",
			Description:  "数据库相关漏洞呈周期性出现",
			DetectedAt:   time.Now().AddDate(0, 0, -5),
			AffectedArea: "数据库安全",
			Score:        0.78,
		},
		{
			Type:         "偏差",
			Severity:     "低",
			Description:  "认证模块漏洞修复率低于基准",
			DetectedAt:   time.Now().AddDate(0, 0, -3),
			AffectedArea: "身份认证",
			Score:        0.65,
		},
	}

	return models.AnomalyDetectionData{
		AnomaliesDetected: anomalies,
		TimeRange:         timeRange,
		BaselineData: map[string]interface{}{
			"dailyAverage": 8.5,
			"weeklyTrend":  "稳定",
			"normalRange":  []float64{5.0, 12.0},
		},
	}
}

// simulateResourceOptimization 模拟资源优化
func (s *AIAnalysisService) simulateResourceOptimization(ctx context.Context, currentAllocation map[string]float64) models.ResourceOptimizationData {
	// 模拟推荐的资源分配
	recommendedAllocation := map[string]float64{
		"高危漏洞修复": 40, // 增加
		"中危漏洞修复": 25, // 不变
		"低危漏洞修复": 10, // 减少
		"漏洞扫描":   20, // 增加
		"安全培训":   5,  // 减少
	}

	// 计算潜在改进
	potentialImprovement := 25.5

	// 模拟识别的瓶颈
	bottlenecks := []models.Bottleneck{
		{
			Area:        "漏洞修复流程",
			Severity:    "高",
			Description: "开发团队和安全团队之间的协作效率低",
			Solution:    "实施更有效的团队协作工具和流程",
		},
		{
			Area:        "漏洞验证",
			Severity:    "中",
			Description: "漏洞修复验证过程耗时较长",
			Solution:    "引入自动化测试和验证工具",
		},
		{
			Area:        "低危漏洞积压",
			Severity:    "低",
			Description: "低危漏洞数量较多但优先级不足",
			Solution:    "批量修复策略或考虑接受部分风险",
		},
	}

	return models.ResourceOptimizationData{
		CurrentAllocation:     currentAllocation,
		RecommendedAllocation: recommendedAllocation,
		PotentialImprovement:  potentialImprovement,
		BottlenecksIdentified: bottlenecks,
	}
}
