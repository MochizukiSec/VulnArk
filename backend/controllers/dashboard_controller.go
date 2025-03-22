package controllers

import (
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vulnark/vulnark/models"
	"github.com/vulnark/vulnark/utils"
)

// DashboardController 控制器结构体
type DashboardController struct{}

// 仪表盘统计数据
type DashboardStats struct {
	TotalVulns         int     `json:"totalVulns"`
	CriticalVulns      int     `json:"criticalVulns"`
	TotalAssets        int     `json:"totalAssets"`
	FixedRate          float64 `json:"fixedRate"`
	VulnChangeRate     float64 `json:"vulnChangeRate"`
	CriticalChangeRate float64 `json:"criticalChangeRate"`
	AssetChangeRate    float64 `json:"assetChangeRate"`
	FixedChangeRate    float64 `json:"fixedChangeRate"`
}

// 漏洞趋势数据
type VulnTrendsData struct {
	Dates []string `json:"dates"`
	New   []int    `json:"new"`
	Fixed []int    `json:"fixed"`
}

// 严重程度分布数据
type SeverityDistData struct {
	Severity string `json:"severity"`
	Count    int    `json:"count"`
}

// 资产漏洞分布数据
type AssetVulnData struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// 活动记录
type ActivityData struct {
	Type     string `json:"type"`
	Content  string `json:"content"`
	Username string `json:"username"`
	Time     string `json:"time"`
}

// 优先修复漏洞数据
type PriorityVuln struct {
	ID           uint   `json:"id"`
	Title        string `json:"title"`
	Severity     string `json:"severity"`
	Status       string `json:"status"`
	Asset        string `json:"asset"`
	DiscoveredAt string `json:"discoveredAt"`
}

// GetDashboardStats 获取仪表盘统计数据
func (dc *DashboardController) GetDashboardStats(c *gin.Context) {
	var stats DashboardStats

	// 获取总漏洞数
	if err := utils.DB.Model(&models.Vulnerability{}).Count(&stats.TotalVulns).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取总漏洞数失败",
			"error":   err.Error(),
		})
		return
	}

	// 获取严重漏洞数
	if err := utils.DB.Model(&models.Vulnerability{}).Where("severity = ?", "critical").Count(&stats.CriticalVulns).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取严重漏洞数失败",
			"error":   err.Error(),
		})
		return
	}

	// 获取总资产数
	if err := utils.DB.Model(&models.Asset{}).Count(&stats.TotalAssets).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取总资产数失败",
			"error":   err.Error(),
		})
		return
	}

	// 计算修复率
	var fixedVulns int
	if err := utils.DB.Model(&models.Vulnerability{}).Where("status = ?", "fixed").Count(&fixedVulns).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取已修复漏洞数失败",
			"error":   err.Error(),
		})
		return
	}

	if stats.TotalVulns > 0 {
		stats.FixedRate = float64(fixedVulns) / float64(stats.TotalVulns) * 100
	}

	// 计算变化率
	// 获取上个月结束时间（本月第一天的前一秒）
	lastMonthEnd := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Local).Add(-time.Second)

	// 上个月总漏洞数
	var lastMonthTotalVulns int
	utils.DB.Model(&models.Vulnerability{}).Where("created_at <= ?", lastMonthEnd).Count(&lastMonthTotalVulns)

	// 上个月严重漏洞数
	var lastMonthCriticalVulns int
	utils.DB.Model(&models.Vulnerability{}).Where("severity = ? AND created_at <= ?", "critical", lastMonthEnd).Count(&lastMonthCriticalVulns)

	// 上个月总资产数
	var lastMonthTotalAssets int
	utils.DB.Model(&models.Asset{}).Where("created_at <= ?", lastMonthEnd).Count(&lastMonthTotalAssets)

	// 上个月修复率
	var lastMonthFixedVulns int
	utils.DB.Model(&models.Vulnerability{}).Where("status = ? AND updated_at <= ?", "fixed", lastMonthEnd).Count(&lastMonthFixedVulns)

	var lastMonthFixedRate float64
	if lastMonthTotalVulns > 0 {
		lastMonthFixedRate = float64(lastMonthFixedVulns) / float64(lastMonthTotalVulns) * 100
	}

	// 计算变化率
	if lastMonthTotalVulns > 0 {
		stats.VulnChangeRate = float64(stats.TotalVulns-lastMonthTotalVulns) / float64(lastMonthTotalVulns) * 100
	} else {
		stats.VulnChangeRate = 100 // 如果上个月为0，则变化率为100%
	}

	if lastMonthCriticalVulns > 0 {
		stats.CriticalChangeRate = float64(stats.CriticalVulns-lastMonthCriticalVulns) / float64(lastMonthCriticalVulns) * 100
	} else {
		stats.CriticalChangeRate = 100
	}

	if lastMonthTotalAssets > 0 {
		stats.AssetChangeRate = float64(stats.TotalAssets-lastMonthTotalAssets) / float64(lastMonthTotalAssets) * 100
	} else {
		stats.AssetChangeRate = 100
	}

	if lastMonthFixedRate > 0 {
		stats.FixedChangeRate = stats.FixedRate - lastMonthFixedRate
	} else {
		stats.FixedChangeRate = stats.FixedRate
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取仪表盘统计数据成功",
		"data":    stats,
	})
}

// GetVulnTrends 获取漏洞趋势数据
func (dc *DashboardController) GetVulnTrends(c *gin.Context) {
	period := c.DefaultQuery("period", "month")

	var trendsData VulnTrendsData

	// 根据时间周期生成日期范围
	var startDate time.Time
	var dateFormat string
	var dateStep int
	var dateCount int

	switch period {
	case "week":
		// 显示过去7天的数据，每天一个数据点
		startDate = time.Now().AddDate(0, 0, -6) // 从6天前开始（包括今天共7天）
		dateFormat = "01-02"                     // 月-日格式
		dateStep = 1                             // 步长1天
		dateCount = 7                            // 总共7个数据点
	case "month":
		startDate = time.Now().AddDate(0, 0, -28)
		dateFormat = "01-02"
		dateStep = 7
		dateCount = 5
	case "quarter":
		startDate = time.Now().AddDate(0, -2, 0)
		dateFormat = "01月"
		dateStep = 30
		dateCount = 3
	default:
		startDate = time.Now().AddDate(0, 0, -28)
		dateFormat = "01-02"
		dateStep = 7
		dateCount = 5
	}

	// 生成日期数组和查询每个时间段的数据
	for i := 0; i < dateCount; i++ {
		currentDate := startDate.AddDate(0, 0, i*dateStep)
		nextDate := startDate.AddDate(0, 0, (i+1)*dateStep)
		if i == dateCount-1 {
			nextDate = time.Now().AddDate(0, 0, 1) // 最后一个时间段包括今天
		}

		// 使用北京时区的时间范围进行查询
		cstCurrentDate := currentDate.In(utils.CSTZone)
		cstNextDate := nextDate.In(utils.CSTZone)
		trendsData.Dates = append(trendsData.Dates, cstCurrentDate.Format(dateFormat))

		// 查询这个时间段内新增的漏洞数量
		var newCount int
		utils.DB.Model(&models.Vulnerability{}).
			Where("created_at >= ? AND created_at < ? AND deleted_at IS NULL",
				cstCurrentDate.Format("2006-01-02")+" 00:00:00",
				cstNextDate.Format("2006-01-02")+" 00:00:00").
			Count(&newCount)
		trendsData.New = append(trendsData.New, newCount)

		// 查询这个时间段内修复的漏洞数量
		var fixedCount int
		utils.DB.Model(&models.Vulnerability{}).
			Where("status = ? AND updated_at >= ? AND updated_at < ? AND deleted_at IS NULL",
				"fixed",
				cstCurrentDate.Format("2006-01-02")+" 00:00:00",
				cstNextDate.Format("2006-01-02")+" 00:00:00").
			Count(&fixedCount)
		trendsData.Fixed = append(trendsData.Fixed, fixedCount)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取漏洞趋势数据成功",
		"data":    trendsData,
	})
}

// GetSeverityDistribution 获取漏洞严重程度分布
func (dc *DashboardController) GetSeverityDistribution(c *gin.Context) {
	var severityDist []SeverityDistData

	// 查询各严重程度的漏洞数量
	severityLevels := []string{"critical", "high", "medium", "low", "info"}

	for _, severity := range severityLevels {
		var count int
		// 显式指定只查询未删除的漏洞
		if err := utils.DB.Model(&models.Vulnerability{}).
			Where("severity = ? AND deleted_at IS NULL", severity).
			Count(&count).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "获取漏洞严重程度分布失败",
				"error":   err.Error(),
			})
			return
		}

		severityDist = append(severityDist, SeverityDistData{
			Severity: severity,
			Count:    count,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取漏洞严重程度分布成功",
		"data":    severityDist,
	})
}

// GetAssetVulnDistribution 获取资产漏洞分布
func (dc *DashboardController) GetAssetVulnDistribution(c *gin.Context) {
	var assetVulnDist []AssetVulnData

	// 执行原始SQL查询来获取每个资产的漏洞数量
	rows, err := utils.DB.Raw(`
		SELECT a.name, COUNT(DISTINCT va.vulnerability_id) as count
		FROM assets a
		LEFT JOIN vulnerability_assets va ON a.id = va.asset_id
		GROUP BY a.id
		ORDER BY count DESC
		LIMIT 5
	`).Rows()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取资产漏洞分布失败: " + err.Error(),
		})
		return
	}
	defer rows.Close()

	// 处理查询结果
	for rows.Next() {
		var name string
		var count int
		if err := rows.Scan(&name, &count); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "处理查询结果失败: " + err.Error(),
			})
			return
		}
		assetVulnDist = append(assetVulnDist, AssetVulnData{
			Name:  name,
			Count: count,
		})
	}

	// 如果没有关联数据，至少返回一些资产
	if len(assetVulnDist) == 0 {
		// 查询最近添加的资产，附加0漏洞数量
		var assets []models.Asset
		if err := utils.DB.Order("created_at DESC").Limit(5).Find(&assets).Error; err == nil {
			for _, asset := range assets {
				assetVulnDist = append(assetVulnDist, AssetVulnData{
					Name:  asset.Name,
					Count: 0,
				})
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取资产漏洞分布成功",
		"data":    assetVulnDist,
	})
}

// GetPriorityVulns 获取优先修复漏洞列表
func (dc *DashboardController) GetPriorityVulns(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "5")
	limit, _ := strconv.Atoi(limitStr)
	if limit <= 0 {
		limit = 5
	}

	var priorityVulns []PriorityVuln

	// 查询需要优先修复的漏洞
	var vulnerabilities []models.Vulnerability
	query := utils.DB.Where("status IN (?, ?)",
		string(models.StatusNew), string(models.StatusVerified))
	query = query.Where("severity IN (?, ?)",
		string(models.SeverityCritical), string(models.SeverityHigh))
	query = query.Order("CASE severity " +
		"WHEN '" + string(models.SeverityCritical) + "' THEN 1 " +
		"WHEN '" + string(models.SeverityHigh) + "' THEN 2 " +
		"ELSE 3 END, created_at DESC")
	query = query.Limit(limit)

	if err := query.Find(&vulnerabilities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取优先修复漏洞失败: " + err.Error(),
		})
		return
	}

	// 如果没有优先漏洞，获取最近的漏洞
	if len(vulnerabilities) == 0 {
		if err := utils.DB.Order("created_at DESC").Limit(limit).Find(&vulnerabilities).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "获取漏洞失败: " + err.Error(),
			})
			return
		}
	}

	// 转换为响应数据结构
	for _, vuln := range vulnerabilities {
		// 查找关联的资产
		var assets []models.Asset
		// 使用预加载替代关联查询
		utils.DB.Raw(`
			SELECT a.* FROM assets a
			JOIN vulnerability_assets va ON a.id = va.asset_id
			WHERE va.vulnerability_id = ?
			LIMIT 1
		`, vuln.ID).Scan(&assets)

		var assetName string
		if len(assets) > 0 {
			assetName = assets[0].Name
		} else {
			assetName = "未关联资产"
		}

		priorityVulns = append(priorityVulns, PriorityVuln{
			ID:           vuln.ID,
			Title:        vuln.Title,
			Severity:     string(vuln.Severity),
			Status:       string(vuln.Status),
			Asset:        assetName,
			DiscoveredAt: utils.FormatTimeCST(vuln.CreatedAt), // 使用创建时间作为发现时间，转换为北京时区
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取优先修复漏洞列表成功",
		"data":    priorityVulns,
	})
}

// GetRecentActivities 获取最近活动
func (dc *DashboardController) GetRecentActivities(c *gin.Context) {
	// 从请求中获取limit参数，默认为5
	limitStr := c.DefaultQuery("limit", "5")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 5
	}

	// 查询最近的漏洞活动
	var activities []ActivityData
	var vulns []models.Vulnerability

	// 查询最近创建或更新的漏洞
	if err := utils.DB.Order("updated_at DESC").Limit(limit).Find(&vulns).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询活动信息失败"})
		return
	}

	// 获取用户信息
	userMap := make(map[uint]string)
	var users []models.User
	if err := utils.DB.Select("id, username").Find(&users).Error; err == nil {
		for _, user := range users {
			userMap[user.ID] = user.Username
		}
	}

	// 转换为响应数据结构
	for _, vuln := range vulns {
		// 获取关联资产
		var assets []models.Asset
		// 使用Raw SQL替代Table查询
		utils.DB.Raw(`
			SELECT a.* FROM assets a
			JOIN vulnerability_assets va ON a.id = va.asset_id
			WHERE va.vulnerability_id = ?
		`, vuln.ID).Scan(&assets)

		// 提取资产名称
		var assetNames []string
		for _, asset := range assets {
			assetNames = append(assetNames, asset.Name)
		}

		username := "system"
		if u, exists := userMap[vuln.ReportedBy]; exists && u != "" {
			username = u
		}

		// 添加到活动列表
		activity := ActivityData{
			Type:     "vulnerability",
			Content:  "更新了漏洞 " + vuln.Title + " 的状态为" + getStatusLabel(string(vuln.Status)),
			Username: username,
			Time:     utils.FormatTimeCST(vuln.UpdatedAt),
		}
		activities = append(activities, activity)
	}

	// 查询最近的资产活动
	var assets []models.Asset
	assetLimit := limit - len(activities)
	if assetLimit > 0 {
		if err := utils.DB.Order("updated_at DESC").Limit(assetLimit).Find(&assets).Error; err == nil {
			for _, asset := range assets {
				// 添加到活动列表
				activity := ActivityData{
					Type:     "asset",
					Content:  "更新了资产 <strong>" + asset.Name + "</strong> 的信息",
					Username: "admin", // 默认为admin
					Time:     utils.FormatTimeCST(asset.UpdatedAt),
				}
				activities = append(activities, activity)
			}
		}
	}

	// 查询最近新增的用户
	var newUsers []models.User
	userLimit := limit - len(activities)
	if userLimit > 0 {
		if err := utils.DB.Order("created_at DESC").Limit(userLimit).Find(&newUsers).Error; err == nil {
			for _, user := range newUsers {
				// 添加到活动列表
				activity := ActivityData{
					Type:     "user",
					Content:  "新增用户 <strong>" + user.Username + "</strong>",
					Username: "admin", // 默认为admin
					Time:     utils.FormatTimeCST(user.CreatedAt),
				}
				activities = append(activities, activity)
			}
		}
	}

	// 查询系统设置变更
	var settings []models.Setting
	settingLimit := limit - len(activities)
	if settingLimit > 0 {
		if err := utils.DB.Order("updated_at DESC").Limit(settingLimit).Find(&settings).Error; err == nil {
			for _, setting := range settings {
				// 添加到活动列表
				activity := ActivityData{
					Type:     "setting",
					Content:  "更新了系统设置 <strong>" + setting.Name + "</strong>",
					Username: "admin", // 默认为admin
					Time:     utils.FormatTimeCST(setting.UpdatedAt),
				}
				activities = append(activities, activity)
			}
		}
	}

	// 查询最近扫描任务
	var scanTasks []models.ScanTask
	scanLimit := limit - len(activities)
	if scanLimit > 0 {
		if err := utils.DB.Order("created_at DESC").Limit(scanLimit).Find(&scanTasks).Error; err == nil {
			for _, task := range scanTasks {
				// 获取操作用户
				username := "system"
				if u, exists := userMap[task.CreatedBy]; exists && u != "" {
					username = u
				}

				// 添加到活动列表
				activity := ActivityData{
					Type:     "scan",
					Content:  "创建了扫描任务 <strong>" + task.Name + "</strong>",
					Username: username,
					Time:     utils.FormatTimeCST(task.CreatedAt),
				}
				activities = append(activities, activity)
			}
		}
	}

	// 查询知识库活动
	var knowledgeItems []models.Knowledge
	knowledgeLimit := limit - len(activities)
	if knowledgeLimit > 0 {
		if err := utils.DB.Order("updated_at DESC").Limit(knowledgeLimit).Find(&knowledgeItems).Error; err == nil {
			for _, item := range knowledgeItems {
				// 获取操作用户
				username := "system"
				if u, exists := userMap[item.AuthorID]; exists && u != "" {
					username = u
				}

				var action string
				if item.CreatedAt.Equal(item.UpdatedAt) {
					action = "新增了知识库条目"
				} else {
					action = "更新了知识库条目"
				}

				// 添加到活动列表
				activity := ActivityData{
					Type:     "knowledge",
					Content:  action + " <strong>" + item.Title + "</strong>",
					Username: username,
					Time:     utils.FormatTimeCST(item.UpdatedAt),
				}
				activities = append(activities, activity)
			}
		}
	}

	// 查询漏洞库活动
	var vulnDBEntries []models.VulnDB
	vulnDBLimit := limit - len(activities)
	if vulnDBLimit > 0 {
		if err := utils.DB.Order("updated_at DESC").Limit(vulnDBLimit).Find(&vulnDBEntries).Error; err == nil {
			for _, entry := range vulnDBEntries {
				// 这里无法直接获取创建者信息，使用默认系统用户
				username := "system"

				var action string
				if entry.CreatedAt.Equal(entry.UpdatedAt) {
					action = "新增了漏洞库条目"
				} else {
					action = "更新了漏洞库条目"
				}

				// 添加到活动列表
				activity := ActivityData{
					Type:     "vulndb",
					Content:  action + " <strong>" + entry.Title + "</strong>",
					Username: username,
					Time:     utils.FormatTimeCST(entry.UpdatedAt),
				}
				activities = append(activities, activity)
			}
		}
	}

	// 对所有活动按时间排序
	sort.Slice(activities, func(i, j int) bool {
		timeI, _ := time.Parse("2006-01-02 15:04:05", activities[i].Time)
		timeJ, _ := time.Parse("2006-01-02 15:04:05", activities[j].Time)
		return timeI.After(timeJ)
	})

	// 如果活动超过了请求的限制，截取前limit条
	if len(activities) > limit {
		activities = activities[:limit]
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取最近活动成功",
		"data":    activities,
	})
}

// 获取状态标签文本
func getStatusLabel(status string) string {
	statusMap := map[string]string{
		"new":            "新发现",
		"verified":       "已验证",
		"in_progress":    "处理中",
		"fixed":          "已修复",
		"closed":         "已关闭",
		"false_positive": "误报",
	}

	if label, exists := statusMap[status]; exists {
		return label
	}
	return status
}
