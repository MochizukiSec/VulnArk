package models

import (
	"time"
)

// VulnDB 漏洞库条目模型
type VulnDB struct {
	ID uint `json:"id" gorm:"primary_key"`

	CVE         string `json:"cve" gorm:"type:varchar(50);unique_index"`
	CWE         string `json:"cwe" gorm:"type:varchar(50)"`
	Title       string `json:"title" gorm:"type:varchar(255);not null"`
	Description string `json:"description" gorm:"type:text;not null"`

	Severity Severity `json:"severity" gorm:"type:varchar(20);not null"`
	CVSS     float64  `json:"cvss" gorm:"type:float"`

	AffectedSystems  string `json:"affected_systems" gorm:"type:text"`
	AffectedVersions string `json:"affected_versions" gorm:"type:text"`

	Solution   string `json:"solution" gorm:"type:text"`
	References string `json:"references" gorm:"type:text"`

	Tags             string `json:"tags" gorm:"type:varchar(255)"`
	ExploitAvailable bool   `json:"exploit_available" gorm:"default:false"`

	PublishedDate time.Time `json:"published_date"`
	UpdatedDate   time.Time `json:"updated_date"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-" gorm:"index"`
}

// TableName 指定表名
func (VulnDB) TableName() string {
	return "vulndb"
}

// IsCritical 判断是否为严重漏洞
func (v *VulnDB) IsCritical() bool {
	return v.Severity == SeverityCritical || v.CVSS >= 9.0
}

// IsRecent 判断是否为最近发布的漏洞 (30天内)
func (v *VulnDB) IsRecent() bool {
	return time.Since(v.PublishedDate).Hours() < 24*30
}

// HasExploit 判断是否有可用的利用程序
func (v *VulnDB) HasExploit() bool {
	return v.ExploitAvailable
}

// GetCVELink 获取CVE链接
func (v *VulnDB) GetCVELink() string {
	if v.CVE == "" {
		return ""
	}
	return "https://cve.mitre.org/cgi-bin/cvename.cgi?name=" + v.CVE
}

// GetCWELink 获取CWE链接
func (v *VulnDB) GetCWELink() string {
	if v.CWE == "" {
		return ""
	}
	return "https://cwe.mitre.org/data/definitions/" + v.CWE + ".html"
}
