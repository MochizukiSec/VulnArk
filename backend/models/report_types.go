package models

import "time"

// SystemCount 用于统计受影响系统数量
type SystemCount struct {
	System string `json:"system"`
	Count  int    `json:"count"`
}

// SummaryReportData 摘要报告数据结构
type SummaryReportData struct {
	GeneratedAt             time.Time      `json:"generatedAt"`
	TotalVulnerabilities    int            `json:"totalVulnerabilities"`
	OpenVulnerabilities     int            `json:"openVulnerabilities"`
	ResolvedVulnerabilities int            `json:"resolvedVulnerabilities"`
	SeverityCounts          map[string]int `json:"severityCounts"`
	StatusCounts            map[string]int `json:"statusCounts"`
	TopAffectedSystems      []SystemCount  `json:"topAffectedSystems"`
}

// DetailedReportData 详细报告数据结构
type DetailedReportData struct {
	GeneratedAt     time.Time         `json:"generatedAt"`
	Vulnerabilities []Vulnerability   `json:"vulnerabilities"`
	Stats           SummaryReportData `json:"stats"`
}
