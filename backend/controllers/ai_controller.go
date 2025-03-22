package controllers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// VulnerabilityContext 表示发送给AI风险评估的漏洞上下文
type VulnerabilityContext struct {
	ID               string        `json:"id"`
	Title            string        `json:"title"`
	CVEID            string        `json:"cve_id"`
	OriginalSeverity string        `json:"original_severity"`
	Description      string        `json:"description"`
	Assets           []AssetSimple `json:"assets"`
}

// AssetSimple 表示简化的资产信息
type AssetSimple struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Environment string `json:"environment"`
	Importance  string `json:"importance"`
}

// RiskFactor 表示影响风险的因素
type RiskFactor struct {
	Name        string `json:"name"`
	Impact      string `json:"impact"`
	Description string `json:"description"`
}

// Recommendation 表示缓解建议
type Recommendation struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
}

// RiskAssessmentResponse 表示AI风险评估结果
type RiskAssessmentResponse struct {
	VulnerabilityID string           `json:"vulnerability_id"`
	ContextualScore float64          `json:"contextual_score"`
	RiskFactors     []RiskFactor     `json:"risk_factors"`
	Analysis        string           `json:"analysis"`
	Recommendations []Recommendation `json:"recommendations"`
	Timestamp       time.Time        `json:"timestamp"`
}

// PerformRiskAssessment 执行AI风险评估
func PerformRiskAssessment(c *gin.Context) {
	// 记录请求详情
	fmt.Println("\n[AI风险评估] ================ 收到新请求 ================")
	fmt.Printf("[AI风险评估] 请求方法: %s, 路径: %s\n", c.Request.Method, c.Request.URL.Path)
	fmt.Printf("[AI风险评估] 请求头:\n")
	for key, values := range c.Request.Header {
		fmt.Printf("[AI风险评估]   %s: %s\n", key, values)
	}

	// 尝试读取请求体
	var rawBody []byte
	if c.Request.Body != nil {
		rawBody, _ = ioutil.ReadAll(c.Request.Body)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(rawBody))
		fmt.Printf("[AI风险评估] 原始请求体: %s\n", string(rawBody))
	}

	// 解析JSON
	var vulnContext VulnerabilityContext
	if err := c.ShouldBindJSON(&vulnContext); err != nil {
		fmt.Printf("[AI风险评估] ❌ JSON绑定失败: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": fmt.Sprintf("无效的请求数据: %v", err),
			"data":    nil,
		})
		return
	}

	fmt.Printf("[AI风险评估] ✅ 成功解析请求数据: %+v\n", vulnContext)

	// 确保ID不为空，如果为空则使用默认值
	if vulnContext.ID == "" {
		vulnContext.ID = "default-id"
		fmt.Println("[AI风险评估] ⚠️ ID为空，使用默认值")
	}

	// 模拟AI风险评估结果
	response := RiskAssessmentResponse{
		VulnerabilityID: vulnContext.ID,
		ContextualScore: calculateContextualScore(vulnContext),
		RiskFactors: []RiskFactor{
			{
				Name:        "资产重要性",
				Impact:      "高",
				Description: "受影响资产包含多个重要系统",
			},
			{
				Name:        "漏洞严重程度",
				Impact:      "中",
				Description: "基于原始评分和环境因素综合评估",
			},
		},
		Analysis: generateAnalysis(vulnContext),
		Recommendations: []Recommendation{
			{
				Title:       "立即修复",
				Description: "建议在24小时内完成修复",
				Priority:    "高",
			},
			{
				Title:       "加强监控",
				Description: "在修复完成前加强对相关资产的监控",
				Priority:    "中",
			},
		},
		Timestamp: time.Now(),
	}

	fmt.Println("[AI风险评估] ✅ 成功生成评估结果")

	// 返回标准格式的响应
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    response,
	})
	fmt.Println("[AI风险评估] ================ 请求处理完成 ================\n")
}

// calculateContextualScore 计算上下文风险分数
func calculateContextualScore(ctx VulnerabilityContext) float64 {
	// 基础分数（根据原始严重程度）
	var baseScore float64
	switch ctx.OriginalSeverity {
	case "critical", "严重":
		baseScore = 9.0
	case "high", "高危", "高":
		baseScore = 7.0
	case "medium", "中危", "中":
		baseScore = 5.0
	case "low", "低危", "低":
		baseScore = 3.0
	case "info", "信息":
		baseScore = 1.0
	default:
		baseScore = 5.0
	}

	// 资产重要性调整
	importanceMultiplier := 1.0
	for _, asset := range ctx.Assets {
		switch asset.Importance {
		case "关键":
			importanceMultiplier = 1.2
		case "重要":
			importanceMultiplier = 1.1
		case "一般":
			importanceMultiplier = 1.0
		}
	}

	score := baseScore * importanceMultiplier
	if score > 10.0 {
		score = 10.0
	}

	return score
}

// generateAnalysis 生成风险分析报告
func generateAnalysis(ctx VulnerabilityContext) string {
	var analysis string
	if len(ctx.Assets) > 0 {
		analysis = "该漏洞影响了多个重要资产，包括："
		for _, asset := range ctx.Assets {
			analysis += "\n- " + asset.Name + "（" + asset.Type + "）"
		}
		analysis += "\n\n基于资产重要性和漏洞特征，建议优先处理。"
	} else {
		analysis = "暂无受影响资产信息，建议进一步评估影响范围。"
	}
	return analysis
}
