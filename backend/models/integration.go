package models

import (
	"time"
)

// CIIntegration CI/CD集成配置
type CIIntegration struct {
	ID          uint       `json:"id" gorm:"primary_key"`
	Name        string     `json:"name" gorm:"type:varchar(100);not null"`
	Type        string     `json:"type" gorm:"type:varchar(50);not null"` // jenkins, gitlab, github, custom
	Description string     `json:"description" gorm:"type:text"`
	APIKey      string     `json:"api_key" gorm:"type:varchar(64);unique_index;not null"`
	Enabled     bool       `json:"enabled" gorm:"default:true"`
	Config      string     `json:"config" gorm:"type:text"` // JSON格式的额外配置
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"-" gorm:"index"`
}

// TableName 指定表名
func (CIIntegration) TableName() string {
	return "ci_integrations"
}

// IntegrationHistory CI/CD集成历史记录
type IntegrationHistory struct {
	ID              uint      `json:"id" gorm:"primary_key"`
	IntegrationID   uint      `json:"integration_id" gorm:"index;not null"`
	IntegrationType string    `json:"integration_type" gorm:"type:varchar(50);not null"`
	Status          string    `json:"status" gorm:"type:varchar(20);not null"` // success, failed
	Message         string    `json:"message" gorm:"type:text"`
	TotalRecords    int       `json:"total_records" gorm:"default:0"`
	SuccessCount    int       `json:"success_count" gorm:"default:0"`
	ErrorCount      int       `json:"error_count" gorm:"default:0"`
	ExecutedAt      time.Time `json:"executed_at"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// TableName 指定表名
func (IntegrationHistory) TableName() string {
	return "integration_histories"
}
