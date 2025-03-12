package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 报告状态常量
const (
	ReportStatusPending   = "pending"
	ReportStatusCompleted = "completed"
	ReportStatusFailed    = "failed"
)

// 报告格式常量
const (
	ReportFormatPDF   = "pdf"
	ReportFormatExcel = "excel"
	ReportFormatWord  = "word"
	ReportFormatHTML  = "html"
	ReportFormatText  = "text"
)

// 报告类型常量
const (
	ReportTypeSummary    = "summary"
	ReportTypeDetailed   = "detailed"
	ReportTypeCompliance = "compliance"
	ReportTypeTrend      = "trend"
)

// Report 报告模型
type Report struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Type        string             `json:"type" bson:"type"`
	Format      string             `json:"format" bson:"format"`
	Status      string             `json:"status" bson:"status"`
	Description string             `json:"description" bson:"description"`
	FileURL     string             `json:"file_url,omitempty" bson:"file_url,omitempty"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	CreatedBy   ReportUser         `json:"created_by" bson:"created_by"`
	Metadata    ReportMetadata     `json:"metadata" bson:"metadata"`
}

// ReportUser 报告创建者信息
type ReportUser struct {
	ID   primitive.ObjectID `json:"id" bson:"id"`
	Name string             `json:"name,omitempty" bson:"name,omitempty"`
}

// ReportMetadata 报告元数据
type ReportMetadata struct {
	StartDate  time.Time `json:"start_date,omitempty" bson:"start_date,omitempty"`
	EndDate    time.Time `json:"end_date,omitempty" bson:"end_date,omitempty"`
	Severities []string  `json:"severities,omitempty" bson:"severities,omitempty"`
	Statuses   []string  `json:"statuses,omitempty" bson:"statuses,omitempty"`
}

// CreateReportRequest 创建报告请求
type CreateReportRequest struct {
	Name        string    `json:"name" binding:"required"`
	Type        string    `json:"type" binding:"required"`
	Format      string    `json:"format" binding:"required"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Severities  []string  `json:"severities"`
	Statuses    []string  `json:"statuses"`
}

// SummaryReport 摘要报告
type SummaryReport struct {
	TotalVulnerabilities    int                       `json:"total_vulnerabilities"`
	NewVulnerabilities      int                       `json:"new_vulnerabilities"`
	ResolvedVulnerabilities int                       `json:"resolved_vulnerabilities"`
	SeverityDistribution    map[string]int            `json:"severity_distribution"`
	StatusDistribution      map[string]int            `json:"status_distribution"`
	TopVulnerableAssets     []AssetVulnCount          `json:"top_vulnerable_assets"`
	VulnerabilityTrend      []DailyVulnerabilityCount `json:"vulnerability_trend"`
	AverageTimeToResolution float64                   `json:"average_time_to_resolution"`
	StartDate               time.Time                 `json:"start_date"`
	EndDate                 time.Time                 `json:"end_date"`
	GeneratedAt             time.Time                 `json:"generated_at"`
}

// DetailedReport 详细报告
type DetailedReport struct {
	SummaryInfo      SummaryReport     `json:"summary_info"`
	Vulnerabilities  []Vulnerability   `json:"vulnerabilities"`
	AffectedAssets   []Asset           `json:"affected_assets"`
	StartDate        time.Time         `json:"start_date"`
	EndDate          time.Time         `json:"end_date"`
	GeneratedAt      time.Time         `json:"generated_at"`
	FilterConditions map[string]string `json:"filter_conditions"`
}

// AssetVulnCount 资产漏洞计数
type AssetVulnCount struct {
	AssetID            primitive.ObjectID `json:"asset_id"`
	AssetName          string             `json:"asset_name"`
	AssetType          string             `json:"asset_type"`
	VulnerabilityCount int                `json:"vulnerability_count"`
}

// DailyVulnerabilityCount 每日漏洞计数
type DailyVulnerabilityCount struct {
	Date          time.Time `json:"date"`
	NewCount      int       `json:"new_count"`
	ResolvedCount int       `json:"resolved_count"`
	TotalCount    int       `json:"total_count"`
}
