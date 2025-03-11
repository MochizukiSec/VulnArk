package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 报告类型常量
const (
	ReportTypeSummary    = "summary"    // 摘要报告
	ReportTypeDetailed   = "detailed"   // 详细报告
	ReportTypeCompliance = "compliance" // 合规报告
	ReportTypeTrend      = "trend"      // 趋势分析
)

// 报告格式常量
const (
	ReportFormatPDF   = "pdf"   // PDF格式
	ReportFormatExcel = "excel" // Excel格式
	ReportFormatWord  = "word"  // Word格式
	ReportFormatHTML  = "html"  // HTML格式
)

// Report 报告模型
type Report struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name            string             `bson:"name" json:"name"`                           // 报告名称
	Type            string             `bson:"type" json:"type"`                           // 报告类型
	Format          string             `bson:"format" json:"format"`                       // 报告格式
	DateRange       [2]time.Time       `bson:"date_range" json:"date_range"`               // 报告时间范围
	Severities      []string           `bson:"severities" json:"severities"`               // 包含的漏洞严重程度
	FileURL         string             `bson:"file_url" json:"file_url"`                   // 报告文件URL
	GeneratedBy     primitive.ObjectID `bson:"generated_by" json:"generated_by"`           // 生成者ID
	GeneratedByInfo UserBasicInfo      `bson:"generated_by_info" json:"generated_by_info"` // 生成者基本信息
	Status          string             `bson:"status" json:"status"`                       // 状态：pending, completed, failed
	Description     string             `bson:"description" json:"description"`             // 报告描述
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
}

// UserBasicInfo 用户基本信息（用于嵌入到其他文档中）
type UserBasicInfo struct {
	ID        primitive.ObjectID `bson:"id" json:"id"`
	Username  string             `bson:"username" json:"username"`
	FirstName string             `bson:"first_name" json:"first_name"`
	LastName  string             `bson:"last_name" json:"last_name"`
}

// ReportRequest 创建报告的请求参数
type ReportRequest struct {
	Name        string    `json:"name" binding:"required"`
	Type        string    `json:"type" binding:"required,oneof=summary detailed compliance trend"`
	Format      string    `json:"format" binding:"required,oneof=pdf excel word html"`
	StartDate   time.Time `json:"start_date" binding:"required"`
	EndDate     time.Time `json:"end_date" binding:"required"`
	Severities  []string  `json:"severities" binding:"required"`
	Description string    `json:"description"`
}

// ReportResponse 报告响应数据
type ReportResponse struct {
	ID              string       `json:"id"`
	Name            string       `json:"name"`
	Type            string       `json:"type"`
	Format          string       `json:"format"`
	DateRange       [2]time.Time `json:"date_range"`
	Severities      []string     `json:"severities"`
	FileURL         string       `json:"file_url"`
	GeneratedBy     string       `json:"generated_by"`
	GeneratedByName string       `json:"generated_by_name"`
	Status          string       `json:"status"`
	Description     string       `json:"description"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
}

// ToResponse 将报告模型转换为响应数据
func (r *Report) ToResponse() ReportResponse {
	return ReportResponse{
		ID:              r.ID.Hex(),
		Name:            r.Name,
		Type:            r.Type,
		Format:          r.Format,
		DateRange:       r.DateRange,
		Severities:      r.Severities,
		FileURL:         r.FileURL,
		GeneratedBy:     r.GeneratedBy.Hex(),
		GeneratedByName: r.GeneratedByInfo.FirstName + " " + r.GeneratedByInfo.LastName,
		Status:          r.Status,
		Description:     r.Description,
		CreatedAt:       r.CreatedAt,
		UpdatedAt:       r.UpdatedAt,
	}
}
