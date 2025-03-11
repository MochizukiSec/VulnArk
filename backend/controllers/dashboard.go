package controllers

import (
	"context"
	"net/http"
	"sort"
	"time"

	"vuln-management/config"
	"vuln-management/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DashboardData 表示仪表盘响应数据
type DashboardData struct {
	TotalVulnerabilities      int                      `json:"totalVulnerabilities"`
	VulnerabilitiesBySeverity map[string]int           `json:"vulnerabilitiesBySeverity"`
	VulnerabilitiesByStatus   map[string]int           `json:"vulnerabilitiesByStatus"`
	RecentVulnerabilities     []models.Vulnerability   `json:"recentVulnerabilities"`
	TopAffectedSystems        []SystemCount            `json:"topAffectedSystems"`
	VulnerabilitiesByMonth    []MonthlyVulnerabilities `json:"vulnerabilitiesByMonth"`
	RiskScore                 float64                  `json:"riskScore"`               // 新增: 整体风险分数
	TeamVulnerabilities       []TeamVulnerabilities    `json:"teamVulnerabilities"`     // 新增: 按团队统计漏洞
	CriticalVulnerabilities   []models.Vulnerability   `json:"criticalVulnerabilities"` // 新增: 需优先关注的高危漏洞
	RemediationProgress       RemediationProgress      `json:"remediationProgress"`     // 新增: 修复进度
	VulnerabilityTrends       VulnerabilityTrends      `json:"vulnerabilityTrends"`     // 新增: 漏洞趋势
}

// SystemCount 表示系统及其漏洞数量
type SystemCount struct {
	System string `json:"system"`
	Count  int    `json:"count"`
}

// MonthlyVulnerabilities 表示按月统计的漏洞数量
type MonthlyVulnerabilities struct {
	Month      string         `json:"month"`
	Year       int            `json:"year"`
	Count      int            `json:"count"`
	Severities map[string]int `json:"severities"`
}

// TeamVulnerabilities 表示团队及其漏洞统计
type TeamVulnerabilities struct {
	Team       string         `json:"team"`
	Count      int            `json:"count"`
	Severities map[string]int `json:"severities"`
	OpenCount  int            `json:"openCount"`
}

// RemediationProgress 表示漏洞修复进度
type RemediationProgress struct {
	ResolvedCount int     `json:"resolvedCount"`
	TotalCount    int     `json:"totalCount"`
	ProgressRate  float64 `json:"progressRate"`
	AverageDays   float64 `json:"averageDays"`
}

// VulnerabilityTrends 表示漏洞趋势
type VulnerabilityTrends struct {
	NewVulnerabilities      []int    `json:"newVulnerabilities"`
	ResolvedVulnerabilities []int    `json:"resolvedVulnerabilities"`
	TimeLabels              []string `json:"timeLabels"`
	NetChange               []int    `json:"netChange"`
}

// GetDashboardData 生成仪表盘数据
func GetDashboardData(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	vulnsCollection := config.GetCollection(config.VulnerabilitiesCollection)

	// 初始化响应数据
	dashboardData := DashboardData{
		VulnerabilitiesBySeverity: make(map[string]int),
		VulnerabilitiesByStatus:   make(map[string]int),
	}

	// 1. 获取总漏洞数
	total, err := vulnsCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "统计漏洞总数失败", "details": err.Error()})
		return
	}
	dashboardData.TotalVulnerabilities = int(total)

	// 2. 按严重程度统计漏洞
	severityCounts, err := getCountBySeverity(ctx, vulnsCollection)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "按严重程度统计漏洞失败", "details": err.Error()})
		return
	}
	dashboardData.VulnerabilitiesBySeverity = severityCounts

	// 3. 按状态统计漏洞
	statusCounts, err := getCountByStatus(ctx, vulnsCollection)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "按状态统计漏洞失败", "details": err.Error()})
		return
	}
	dashboardData.VulnerabilitiesByStatus = statusCounts

	// 4. 获取最近10个漏洞
	recentVulns, err := getRecentVulnerabilities(ctx, vulnsCollection, 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取最近漏洞失败", "details": err.Error()})
		return
	}
	dashboardData.RecentVulnerabilities = recentVulns

	// 5. 获取受影响最多的系统
	topSystems, err := getTopAffectedSystems(ctx, vulnsCollection, 5)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取受影响系统统计失败", "details": err.Error()})
		return
	}
	dashboardData.TopAffectedSystems = topSystems

	// 6. 获取按月统计的漏洞数量
	monthlyData, err := getVulnerabilitiesByMonth(ctx, vulnsCollection, 6)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取月度统计失败", "details": err.Error()})
		return
	}
	dashboardData.VulnerabilitiesByMonth = monthlyData

	// 7. 计算整体风险分数
	riskScore, err := calculateRiskScore(severityCounts, statusCounts, int(total))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "计算风险分数失败", "details": err.Error()})
		return
	}
	dashboardData.RiskScore = riskScore

	// 8. 获取按团队统计的漏洞数据
	teamVulns, err := getVulnerabilitiesByTeam(ctx, vulnsCollection)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取团队漏洞统计失败", "details": err.Error()})
		return
	}
	dashboardData.TeamVulnerabilities = teamVulns

	// 9. 获取需优先关注的高危漏洞
	criticalVulns, err := getCriticalVulnerabilities(ctx, vulnsCollection, 5)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取高危漏洞失败", "details": err.Error()})
		return
	}
	dashboardData.CriticalVulnerabilities = criticalVulns

	// 10. 获取修复进度
	progress, err := getRemediationProgress(ctx, vulnsCollection)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取修复进度失败", "details": err.Error()})
		return
	}
	dashboardData.RemediationProgress = progress

	// 11. 获取漏洞趋势
	trends, err := getVulnerabilityTrends(ctx, vulnsCollection, 8) // 过去8周的趋势
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取漏洞趋势失败", "details": err.Error()})
		return
	}
	dashboardData.VulnerabilityTrends = trends

	c.JSON(http.StatusOK, dashboardData)
}

// getCountBySeverity 按严重程度统计漏洞数量
func getCountBySeverity(ctx context.Context, collection *mongo.Collection) (map[string]int, error) {
	pipeline := mongo.Pipeline{
		{
			{Key: "$group", Value: bson.M{
				"_id":   "$severity",
				"count": bson.M{"$sum": 1},
			}},
		},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	results := make(map[string]int)
	// 预设所有严重程度，确保即使为0也会显示
	results[models.SeverityCritical] = 0
	results[models.SeverityHigh] = 0
	results[models.SeverityMedium] = 0
	results[models.SeverityLow] = 0
	results[models.SeverityInfo] = 0

	var result struct {
		ID    string `bson:"_id"`
		Count int    `bson:"count"`
	}

	for cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		results[result.ID] = result.Count
	}

	return results, nil
}

// getCountByStatus 按状态统计漏洞数量
func getCountByStatus(ctx context.Context, collection *mongo.Collection) (map[string]int, error) {
	pipeline := mongo.Pipeline{
		{
			{Key: "$group", Value: bson.M{
				"_id":   "$status",
				"count": bson.M{"$sum": 1},
			}},
		},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	results := make(map[string]int)
	// 预设所有状态，确保即使为0也会显示
	results[models.StatusOpen] = 0
	results[models.StatusInProgress] = 0
	results[models.StatusResolved] = 0
	results[models.StatusClosed] = 0
	results[models.StatusFalsePositive] = 0

	var result struct {
		ID    string `bson:"_id"`
		Count int    `bson:"count"`
	}

	for cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		results[result.ID] = result.Count
	}

	return results, nil
}

// getRecentVulnerabilities 获取最近的漏洞
func getRecentVulnerabilities(ctx context.Context, collection *mongo.Collection, limit int) ([]models.Vulnerability, error) {
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(int64(limit))

	cursor, err := collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var vulnerabilities []models.Vulnerability
	if err = cursor.All(ctx, &vulnerabilities); err != nil {
		return nil, err
	}

	return vulnerabilities, nil
}

// getTopAffectedSystems 获取受影响最多的系统
func getTopAffectedSystems(ctx context.Context, collection *mongo.Collection, limit int) ([]SystemCount, error) {
	pipeline := mongo.Pipeline{
		{
			{Key: "$unwind", Value: "$affected_systems"},
		},
		{
			{Key: "$group", Value: bson.M{
				"_id":   "$affected_systems",
				"count": bson.M{"$sum": 1},
			}},
		},
		{
			{Key: "$sort", Value: bson.M{"count": -1}},
		},
		{
			{Key: "$limit", Value: limit},
		},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []SystemCount
	for cursor.Next(ctx) {
		var result struct {
			ID    string `bson:"_id"`
			Count int    `bson:"count"`
		}
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}

		results = append(results, SystemCount{
			System: result.ID,
			Count:  result.Count,
		})
	}

	return results, nil
}

// getVulnerabilitiesByMonth 获取按月统计的漏洞数量
func getVulnerabilitiesByMonth(ctx context.Context, collection *mongo.Collection, months int) ([]MonthlyVulnerabilities, error) {
	// 计算开始日期（当前日期往前推N个月）
	now := time.Now()
	startDate := time.Date(now.Year(), now.Month()-time.Month(months-1), 1, 0, 0, 0, 0, now.Location())

	// 按月统计漏洞
	var results []MonthlyVulnerabilities

	// 生成月份列表（从最早到最近）
	for i := 0; i < months; i++ {
		currentMonth := startDate.AddDate(0, i, 0)
		nextMonth := currentMonth.AddDate(0, 1, 0)

		// 查询该月创建的漏洞总数
		count, err := collection.CountDocuments(ctx, bson.M{
			"created_at": bson.M{
				"$gte": currentMonth,
				"$lt":  nextMonth,
			},
		})
		if err != nil {
			return nil, err
		}

		// 按严重程度统计该月漏洞
		severityCounts := make(map[string]int)

		// 预设所有严重程度，确保即使为0也会显示
		severityCounts[models.SeverityCritical] = 0
		severityCounts[models.SeverityHigh] = 0
		severityCounts[models.SeverityMedium] = 0
		severityCounts[models.SeverityLow] = 0
		severityCounts[models.SeverityInfo] = 0

		// 查询各严重程度的漏洞数量
		for _, severity := range []string{
			models.SeverityCritical,
			models.SeverityHigh,
			models.SeverityMedium,
			models.SeverityLow,
			models.SeverityInfo,
		} {
			severityCount, err := collection.CountDocuments(ctx, bson.M{
				"created_at": bson.M{
					"$gte": currentMonth,
					"$lt":  nextMonth,
				},
				"severity": severity,
			})
			if err != nil {
				return nil, err
			}
			severityCounts[severity] = int(severityCount)
		}

		// 添加到结果
		results = append(results, MonthlyVulnerabilities{
			Month:      currentMonth.Month().String(),
			Year:       currentMonth.Year(),
			Count:      int(count),
			Severities: severityCounts,
		})
	}

	return results, nil
}

// calculateRiskScore 计算整体风险分数
func calculateRiskScore(severityCounts map[string]int, statusCounts map[string]int, totalVulns int) (float64, error) {
	if totalVulns == 0 {
		return 0.0, nil
	}

	// 各严重程度的权重
	severityWeights := map[string]float64{
		models.SeverityCritical: 10.0,
		models.SeverityHigh:     7.5,
		models.SeverityMedium:   5.0,
		models.SeverityLow:      2.5,
		models.SeverityInfo:     0.5,
	}

	// 各状态的权重
	statusWeights := map[string]float64{
		models.StatusOpen:          1.0,
		models.StatusInProgress:    0.7,
		models.StatusResolved:      0.1,
		models.StatusClosed:        0.0,
		models.StatusFalsePositive: 0.0,
	}

	// 计算严重程度得分
	severityScore := 0.0
	for severity, count := range severityCounts {
		weight, exists := severityWeights[severity]
		if exists {
			severityScore += float64(count) * weight
		}
	}

	// 计算状态得分
	statusMultiplier := 0.0
	statusDivisor := 0
	for status, count := range statusCounts {
		weight, exists := statusWeights[status]
		if exists && count > 0 {
			statusMultiplier += float64(count) * weight
			statusDivisor += count
		}
	}

	if statusDivisor == 0 {
		return 0.0, nil
	}

	// 状态系数 (0-1之间)
	statusFactor := statusMultiplier / float64(statusDivisor)

	// 最终风险分数，范围0-100
	maxPossibleScore := float64(totalVulns) * severityWeights[models.SeverityCritical]
	if maxPossibleScore == 0 {
		return 0.0, nil
	}

	riskScore := (severityScore * statusFactor / maxPossibleScore) * 100

	// 确保分数在0-100范围内
	if riskScore > 100 {
		riskScore = 100
	} else if riskScore < 0 {
		riskScore = 0
	}

	return riskScore, nil
}

// getVulnerabilitiesByTeam 按团队统计漏洞
func getVulnerabilitiesByTeam(ctx context.Context, collection *mongo.Collection) ([]TeamVulnerabilities, error) {
	// 模拟团队数据 - 在实际应用中，这应该从数据库中查询
	teams := []string{"安全团队", "开发团队", "运维团队", "产品团队", "测试团队"}

	var results []TeamVulnerabilities

	for _, team := range teams {
		// 模拟查询 - 在实际应用中，应该基于团队标签或自定义字段查询
		// 这里我们使用tags字段模拟团队归属
		filter := bson.M{"tags": team}

		// 获取该团队的漏洞总数
		count, err := collection.CountDocuments(ctx, filter)
		if err != nil {
			return nil, err
		}

		// 按严重程度统计该团队漏洞
		severityCounts := make(map[string]int)
		for _, severity := range []string{
			models.SeverityCritical,
			models.SeverityHigh,
			models.SeverityMedium,
			models.SeverityLow,
			models.SeverityInfo,
		} {
			severityFilter := bson.M{
				"tags":     team,
				"severity": severity,
			}
			severityCount, err := collection.CountDocuments(ctx, severityFilter)
			if err != nil {
				return nil, err
			}
			severityCounts[severity] = int(severityCount)
		}

		// 获取该团队的未解决漏洞数量
		openFilter := bson.M{
			"tags": team,
			"status": bson.M{
				"$in": []string{models.StatusOpen, models.StatusInProgress},
			},
		}
		openCount, err := collection.CountDocuments(ctx, openFilter)
		if err != nil {
			return nil, err
		}

		// 添加团队数据
		results = append(results, TeamVulnerabilities{
			Team:       team,
			Count:      int(count),
			Severities: severityCounts,
			OpenCount:  int(openCount),
		})
	}

	// 按漏洞总数排序
	sort.Slice(results, func(i, j int) bool {
		return results[i].Count > results[j].Count
	})

	return results, nil
}

// getCriticalVulnerabilities 获取需优先关注的高危漏洞
func getCriticalVulnerabilities(ctx context.Context, collection *mongo.Collection, limit int) ([]models.Vulnerability, error) {
	// 查询条件：高危或严重漏洞，状态为未解决，按CVSS分数降序排列
	filter := bson.M{
		"severity": bson.M{
			"$in": []string{models.SeverityCritical, models.SeverityHigh},
		},
		"status": bson.M{
			"$in": []string{models.StatusOpen, models.StatusInProgress},
		},
	}

	options := options.Find().
		SetSort(bson.D{{Key: "cvss", Value: -1}}).
		SetLimit(int64(limit))

	cursor, err := collection.Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var vulnerabilities []models.Vulnerability
	if err = cursor.All(ctx, &vulnerabilities); err != nil {
		return nil, err
	}

	return vulnerabilities, nil
}

// getRemediationProgress 获取修复进度
func getRemediationProgress(ctx context.Context, collection *mongo.Collection) (RemediationProgress, error) {
	var progress RemediationProgress

	// 获取漏洞总数
	totalCount, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return progress, err
	}
	progress.TotalCount = int(totalCount)

	// 获取已解决的漏洞数
	resolvedFilter := bson.M{
		"status": bson.M{
			"$in": []string{models.StatusResolved, models.StatusClosed},
		},
	}
	resolvedCount, err := collection.CountDocuments(ctx, resolvedFilter)
	if err != nil {
		return progress, err
	}
	progress.ResolvedCount = int(resolvedCount)

	// 计算进度百分比
	if progress.TotalCount > 0 {
		progress.ProgressRate = float64(progress.ResolvedCount) / float64(progress.TotalCount) * 100
	}

	// 计算平均修复天数
	// 注意：这需要较复杂的聚合查询，以下是简化版
	pipeline := mongo.Pipeline{
		{
			{Key: "$match", Value: bson.M{
				"status": bson.M{
					"$in": []string{models.StatusResolved, models.StatusClosed},
				},
				"resolved_at": bson.M{"$ne": nil},
			}},
		},
		{
			{Key: "$project", Value: bson.M{
				"resolution_time": bson.M{
					"$divide": []interface{}{
						bson.M{"$subtract": []interface{}{"$resolved_at", "$created_at"}},
						float64(1000 * 60 * 60 * 24), // 转换为天
					},
				},
			}},
		},
		{
			{Key: "$group", Value: bson.M{
				"_id":      nil,
				"avg_days": bson.M{"$avg": "$resolution_time"},
			}},
		},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return progress, err
	}
	defer cursor.Close(ctx)

	type avgResult struct {
		AvgDays float64 `bson:"avg_days"`
	}

	var result avgResult
	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return progress, err
		}
		progress.AverageDays = result.AvgDays
	} else {
		progress.AverageDays = 0
	}

	return progress, nil
}

// getVulnerabilityTrends 获取漏洞趋势
func getVulnerabilityTrends(ctx context.Context, collection *mongo.Collection, weeks int) (VulnerabilityTrends, error) {
	var trends VulnerabilityTrends

	// 初始化数据结构
	trends.NewVulnerabilities = make([]int, weeks)
	trends.ResolvedVulnerabilities = make([]int, weeks)
	trends.NetChange = make([]int, weeks)
	trends.TimeLabels = make([]string, weeks)

	// 计算开始日期（当前日期往前推N周）
	now := time.Now()
	weekStart := now.AddDate(0, 0, -7*weeks+1)
	weekStart = time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 0, 0, 0, 0, weekStart.Location())

	// 生成每周数据
	for i := 0; i < weeks; i++ {
		currentWeekStart := weekStart.AddDate(0, 0, 7*i)
		currentWeekEnd := currentWeekStart.AddDate(0, 0, 7)

		// 设置时间标签
		trends.TimeLabels[i] = currentWeekStart.Format("01/02")

		// 获取该周新增的漏洞数
		newFilter := bson.M{
			"created_at": bson.M{
				"$gte": currentWeekStart,
				"$lt":  currentWeekEnd,
			},
		}
		newCount, err := collection.CountDocuments(ctx, newFilter)
		if err != nil {
			return trends, err
		}
		trends.NewVulnerabilities[i] = int(newCount)

		// 获取该周解决的漏洞数
		resolvedFilter := bson.M{
			"resolved_at": bson.M{
				"$gte": currentWeekStart,
				"$lt":  currentWeekEnd,
			},
		}
		resolvedCount, err := collection.CountDocuments(ctx, resolvedFilter)
		if err != nil {
			return trends, err
		}
		trends.ResolvedVulnerabilities[i] = int(resolvedCount)

		// 计算净变化
		trends.NetChange[i] = trends.NewVulnerabilities[i] - trends.ResolvedVulnerabilities[i]
	}

	return trends, nil
}
