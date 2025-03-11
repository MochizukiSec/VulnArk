package models

import (
	"time"
)

// VulnDBEntry 表示漏洞库中的条目
type VulnDBEntry struct {
	ID               string    `json:"id" bson:"_id,omitempty"`
	CveId            string    `json:"cveId" bson:"cveId"`
	Title            string    `json:"title" bson:"title"`
	Description      string    `json:"description" bson:"description"`
	Cvss             float64   `json:"cvss" bson:"cvss"`
	Severity         string    `json:"severity" bson:"severity"`
	AffectedSystems  string    `json:"affectedSystems" bson:"affectedSystems"`
	Solution         string    `json:"solution" bson:"solution"`
	PublishedDate    string    `json:"publishedDate" bson:"publishedDate"`
	LastModifiedDate string    `json:"lastModifiedDate" bson:"lastModifiedDate"`
	References       []string  `json:"references" bson:"references"`
	Tags             []string  `json:"tags" bson:"tags"`
	CreatedAt        time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt" bson:"updatedAt"`
}

// VulnDBSearch 表示漏洞库搜索参数
type VulnDBSearch struct {
	Query     string `form:"q"`
	Severity  string `form:"severity"`
	Year      string `form:"year"`
	SortBy    string `form:"sortBy" binding:"omitempty,oneof=publishedDate lastModifiedDate cveId cvss"`
	SortOrder string `form:"sortOrder" binding:"omitempty,oneof=asc desc"`
	Page      int    `form:"page,default=1"`
	PerPage   int    `form:"perPage,default=20"`
}
