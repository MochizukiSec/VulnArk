package models

import (
	"time"
)

// ScannerType 扫描器类型
type ScannerType string

const (
	ScannerTypeNessus ScannerType = "nessus" // Nessus扫描器
	ScannerTypeXray   ScannerType = "xray"   // Xray扫描器
	ScannerTypeAwvs   ScannerType = "awvs"   // AWVS扫描器
	ScannerTypeZap    ScannerType = "zap"    // OWASP ZAP扫描器
	ScannerTypeCustom ScannerType = "custom" // 自定义扫描器
)

// ScanTaskStatus 扫描任务状态
type ScanTaskStatus string

const (
	ScanTaskStatusCreated   ScanTaskStatus = "created"   // 已创建
	ScanTaskStatusQueued    ScanTaskStatus = "queued"    // 已排队
	ScanTaskStatusRunning   ScanTaskStatus = "running"   // 运行中
	ScanTaskStatusCompleted ScanTaskStatus = "completed" // 已完成
	ScanTaskStatusFailed    ScanTaskStatus = "failed"    // 失败
	ScanTaskStatusCancelled ScanTaskStatus = "cancelled" // 已取消
)

// ScanTask 扫描任务模型
type ScanTask struct {
	ID          uint           `json:"id" gorm:"primary_key"`
	Name        string         `json:"name" gorm:"type:varchar(255);not null"`
	Description string         `json:"description" gorm:"type:text"`
	Type        ScannerType    `json:"type" gorm:"type:varchar(50);not null"`
	Status      ScanTaskStatus `json:"status" gorm:"type:varchar(20);not null;default:'created'"`

	// 扫描配置
	ScannerURL      string `json:"scanner_url" gorm:"type:varchar(255)"`
	ScannerAPIKey   string `json:"scanner_api_key" gorm:"type:varchar(255)"`
	ScannerUsername string `json:"scanner_username" gorm:"type:varchar(100)"`
	ScannerPassword string `json:"scanner_password" gorm:"type:varchar(100)"`

	// 扫描目标
	TargetIPs      string `json:"target_ips" gorm:"type:text"`      // 逗号分隔的IP地址列表
	TargetURLs     string `json:"target_urls" gorm:"type:text"`     // 逗号分隔的URL列表
	TargetAssets   string `json:"target_assets" gorm:"type:text"`   // 逗号分隔的资产ID列表
	ScanParameters string `json:"scan_parameters" gorm:"type:text"` // 扫描参数（JSON格式）

	// 扫描时间
	ScheduledAt  *time.Time `json:"scheduled_at"`                           // 计划扫描时间
	StartedAt    *time.Time `json:"started_at"`                             // 开始扫描时间
	CompletedAt  *time.Time `json:"completed_at"`                           // 完成扫描时间
	IsRecurring  bool       `json:"is_recurring"`                           // 是否是定期扫描
	CronSchedule string     `json:"cron_schedule" gorm:"type:varchar(100)"` // Cron表达式

	// 扫描结果
	TotalVulnerabilities    int    `json:"total_vulnerabilities"`           // 总漏洞数
	CriticalVulnerabilities int    `json:"critical_vulnerabilities"`        // 严重漏洞数
	HighVulnerabilities     int    `json:"high_vulnerabilities"`            // 高危漏洞数
	MediumVulnerabilities   int    `json:"medium_vulnerabilities"`          // 中危漏洞数
	LowVulnerabilities      int    `json:"low_vulnerabilities"`             // 低危漏洞数
	ResultSummary           string `json:"result_summary" gorm:"type:text"` // 结果摘要
	ScanLog                 string `json:"scan_log" gorm:"type:text"`       // 扫描日志

	// 关联记录
	CreatedBy uint `json:"created_by" gorm:"index;not null"` // 创建者ID

	// 时间戳
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-" gorm:"index"`
}

// ScanResult 扫描结果模型
type ScanResult struct {
	ID         uint `json:"id" gorm:"primary_key"`
	ScanTaskID uint `json:"scan_task_id" gorm:"index;not null"` // 关联的扫描任务ID

	VulnerabilityName string   `json:"vulnerability_name" gorm:"type:varchar(255);not null"`
	Description       string   `json:"description" gorm:"type:text"`
	Severity          Severity `json:"severity" gorm:"type:varchar(20);not null"`

	// 漏洞详情
	AffectedURL  string `json:"affected_url" gorm:"type:varchar(255)"`
	AffectedIP   string `json:"affected_ip" gorm:"type:varchar(50)"`
	AffectedPort string `json:"affected_port" gorm:"type:varchar(50)"`
	Detail       string `json:"detail" gorm:"type:text"`

	// 漏洞分类
	Category string  `json:"category" gorm:"type:varchar(100)"`
	CVE      string  `json:"cve" gorm:"type:varchar(50)"`
	CVSS     float64 `json:"cvss" gorm:"type:float"`

	// 解决方案
	Solution   string `json:"solution" gorm:"type:text"`
	References string `json:"references" gorm:"type:text"`

	// 导入状态
	IsImported bool       `json:"is_imported" gorm:"default:false"` // 是否已导入到漏洞库
	ImportedAt *time.Time `json:"imported_at"`                      // 导入时间
	ImportedID uint       `json:"imported_id"`                      // 导入后的漏洞ID

	// 时间戳
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-" gorm:"index"`
}

// TableName 指定表名
func (ScanTask) TableName() string {
	return "scan_tasks"
}

// TableName 指定表名
func (ScanResult) TableName() string {
	return "scan_results"
}

// IsCriticalTask 判断是否为包含严重漏洞的任务
func (s *ScanTask) IsCriticalTask() bool {
	return s.CriticalVulnerabilities > 0
}

// IsCompleted 判断任务是否已完成
func (s *ScanTask) IsCompleted() bool {
	return s.Status == ScanTaskStatusCompleted
}

// IsInProgress 判断任务是否正在进行
func (s *ScanTask) IsInProgress() bool {
	return s.Status == ScanTaskStatusQueued || s.Status == ScanTaskStatusRunning
}
