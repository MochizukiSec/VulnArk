package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AIAnalysisType 定义AI分析的类型
type AIAnalysisType string

const (
	TrendPrediction      AIAnalysisType = "trend_prediction"
	ResourceOptimization AIAnalysisType = "resource_optimization"
	AnomalyDetection     AIAnalysisType = "anomaly_detection"
)

// AIAnalysis 定义AI分析结果的数据模型
type AIAnalysis struct {
	ID              primitive.ObjectID     `bson:"_id" json:"id"`
	Type            AIAnalysisType         `bson:"type" json:"type"`
	Title           string                 `bson:"title" json:"title"`
	Description     string                 `bson:"description" json:"description"`
	AnalysisData    interface{}            `bson:"analysis_data" json:"analysisData"`
	Recommendations []string               `bson:"recommendations" json:"recommendations"`
	Confidence      float64                `bson:"confidence" json:"confidence"`
	CreatedAt       time.Time              `bson:"created_at" json:"createdAt"`
	UpdatedAt       time.Time              `bson:"updated_at" json:"updatedAt"`
	Parameters      map[string]interface{} `bson:"parameters" json:"parameters"`
}

// TrendPredictionData 漏洞趋势预测数据
type TrendPredictionData struct {
	TimeRange       string            `json:"timeRange"`       // 预测时间范围，如"未来30天"
	PredictedCounts map[string]int    `json:"predictedCounts"` // 按类别/严重程度预测的漏洞数量
	TrendFactors    []TrendFactor     `json:"trendFactors"`    // 影响趋势的因素
	HistoricalData  []HistoricalPoint `json:"historicalData"`  // 用于预测的历史数据点
}

// TrendFactor 影响趋势的因素
type TrendFactor struct {
	Factor      string  `json:"factor"`
	Impact      float64 `json:"impact"`     // 正值表示增加，负值表示减少
	Confidence  float64 `json:"confidence"` // 0-1之间，表示置信度
	Description string  `json:"description"`
}

// HistoricalPoint 历史数据点
type HistoricalPoint struct {
	Date  time.Time `json:"date"`
	Count int       `json:"count"`
	Label string    `json:"label,omitempty"` // 可选的分类标签
}

// ResourceOptimizationData 资源优化建议数据
type ResourceOptimizationData struct {
	CurrentAllocation     map[string]float64 `json:"currentAllocation"`     // 当前资源分配
	RecommendedAllocation map[string]float64 `json:"recommendedAllocation"` // 建议的资源分配
	PotentialImprovement  float64            `json:"potentialImprovement"`  // 预计的改进百分比
	BottlenecksIdentified []Bottleneck       `json:"bottlenecksIdentified"` // 识别的瓶颈
}

// Bottleneck 资源瓶颈
type Bottleneck struct {
	Area        string `json:"area"`     // 瓶颈区域
	Severity    string `json:"severity"` // 严重程度
	Description string `json:"description"`
	Solution    string `json:"solution"`
}

// AnomalyDetectionData 异常检测数据
type AnomalyDetectionData struct {
	AnomaliesDetected []Anomaly   `json:"anomaliesDetected"` // 检测到的异常
	TimeRange         string      `json:"timeRange"`         // 分析的时间范围
	BaselineData      interface{} `json:"baselineData"`      // 基准数据
}

// Anomaly 异常数据
type Anomaly struct {
	Type         string    `json:"type"`     // 异常类型
	Severity     string    `json:"severity"` // 严重程度
	Description  string    `json:"description"`
	DetectedAt   time.Time `json:"detectedAt"`
	AffectedArea string    `json:"affectedArea"`
	Score        float64   `json:"score"` // 异常分数，越高越异常
}

// AIAnalysisRequest 定义AI分析请求的参数
type AIAnalysisRequest struct {
	Type       AIAnalysisType         `json:"type" binding:"required"`
	Parameters map[string]interface{} `json:"parameters"`
}

// AIAnalysisResponse 定义AI分析响应
type AIAnalysisResponse struct {
	ID              primitive.ObjectID `json:"id"`
	Type            AIAnalysisType     `json:"type"`
	Title           string             `json:"title"`
	Description     string             `json:"description"`
	Recommendations []string           `json:"recommendations"`
	Confidence      float64            `json:"confidence"`
	CreatedAt       time.Time          `json:"createdAt"`
	AnalysisData    interface{}        `json:"analysisData"`
}
