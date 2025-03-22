package models

import (
	"time"
)

// 资产类型
type AssetType string

const (
	AssetTypeHost        AssetType = "host"        // 主机
	AssetTypeWebsite     AssetType = "website"     // 网站
	AssetTypeDatabase    AssetType = "database"    // 数据库
	AssetTypeApplication AssetType = "application" // 应用程序
	AssetTypeServer      AssetType = "server"      // 服务器
	AssetTypeNetwork     AssetType = "network"     // 网络设备
	AssetTypeCloud       AssetType = "cloud"       // 云服务
	AssetTypeIoT         AssetType = "iot"         // 物联网设备
	AssetTypeOther       AssetType = "other"       // 其他
)

// 资产状态
type AssetStatus string

const (
	AssetStatusActive   AssetStatus = "active"   // 活跃
	AssetStatusInactive AssetStatus = "inactive" // 非活跃
	AssetStatusArchived AssetStatus = "archived" // 已归档
)

// 资产重要性级别
type AssetImportance string

const (
	ImportanceCritical AssetImportance = "critical" // 关键
	ImportanceHigh     AssetImportance = "high"     // 高
	ImportanceMedium   AssetImportance = "medium"   // 中
	ImportanceLow      AssetImportance = "low"      // 低
)

// Asset 资产模型
type Asset struct {
	ID         uint        `json:"id" gorm:"primary_key"`
	Name       string      `json:"name" gorm:"type:varchar(255);not null"`
	Type       AssetType   `json:"type" gorm:"type:varchar(50);not null"`
	Identifier string      `json:"identifier" gorm:"type:varchar(255);unique_index;not null"`
	Status     AssetStatus `json:"status" gorm:"type:varchar(20);not null;default:'active'"`

	Description string `json:"description" gorm:"type:text"`
	URL         string `json:"url" gorm:"type:varchar(255)"`
	IPAddress   string `json:"ip_address" gorm:"type:varchar(50)"`
	Port        string `json:"port" gorm:"type:varchar(100)"`
	OS          string `json:"os" gorm:"type:varchar(100)"`
	Version     string `json:"version" gorm:"type:varchar(50)"`

	Owner      string `json:"owner" gorm:"type:varchar(100)"`
	Department string `json:"department" gorm:"type:varchar(100)"`
	Location   string `json:"location" gorm:"type:varchar(100)"`

	Importance AssetImportance `json:"importance" gorm:"type:varchar(20);not null;default:'medium'"`

	Tags string `json:"tags" gorm:"type:varchar(255)"`

	LastScan *time.Time `json:"last_scan"`

	Vulnerabilities []Vulnerability `json:"vulnerabilities" gorm:"many2many:vulnerability_assets;"`

	CustomFields string `json:"custom_fields" gorm:"type:text"`
	Notes        string `json:"notes" gorm:"type:text"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-" gorm:"index"`
}

// TableName 指定表名
func (Asset) TableName() string {
	return "assets"
}

// IsCritical 判断是否为关键资产
func (a *Asset) IsCritical() bool {
	return a.Importance == ImportanceCritical
}

// IsActive 判断是否为活跃资产
func (a *Asset) IsActive() bool {
	return a.Status == AssetStatusActive
}

// GetVulnCount 获取漏洞数量
func (a *Asset) GetVulnCount() int {
	return len(a.Vulnerabilities)
}

// UpdateLastScan 更新最后扫描时间
func (a *Asset) UpdateLastScan() {
	now := time.Now()
	a.LastScan = &now
}

// SetArchived 归档资产
func (a *Asset) SetArchived() {
	a.Status = AssetStatusArchived
}
