package controllers

import (
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// 漏洞模型
type Vulnerability struct {
	ID               string     `json:"id"`
	CVEId            string     `json:"cveId"`
	Title            string     `json:"title"`
	Description      string     `json:"description"`
	Severity         string     `json:"severity"`
	CVSS             float64    `json:"cvss"`
	CVSSVector       string     `json:"cvssVector"`
	Vendor           string     `json:"vendor"`
	Product          string     `json:"product"`
	AffectedVersions string     `json:"affectedVersions"`
	Solution         string     `json:"solution"`
	Exploit          string     `json:"exploit"`
	References       []string   `json:"references"`
	Tags             []string   `json:"tags"`
	PublishedDate    *time.Time `json:"publishedDate"`
	LastModifiedDate *time.Time `json:"lastModifiedDate"`
}

// VulnDatabaseController 漏洞库控制器
type VulnDatabaseController struct{}

// VulnDatabaseSearchParams 漏洞库搜索参数
type VulnDatabaseSearchParams struct {
	Keyword   string `form:"keyword"`
	CveID     string `form:"cveId"`
	Severity  string `form:"severity"`
	Vendor    string `form:"vendor"`
	Product   string `form:"product"`
	Year      string `form:"year"`
	SortBy    string `form:"sortBy"`
	SortOrder string `form:"sortOrder"`
	Page      int    `form:"page,default=1"`
	PageSize  int    `form:"pageSize,default=10"`
}

// VulnDatabaseResponse NVD API 响应结构
type VulnDatabaseResponse struct {
	ResultsPerPage  int                 `json:"resultsPerPage"`
	StartIndex      int                 `json:"startIndex"`
	TotalResults    int                 `json:"totalResults"`
	Format          string              `json:"format"`
	Version         string              `json:"version"`
	Vulnerabilities []VulnerabilityItem `json:"vulnerabilities"`
}

// VulnerabilityItem 漏洞项
type VulnerabilityItem struct {
	CVE CVEItem `json:"cve"`
}

// CVEItem CVE项
type CVEItem struct {
	ID               string                 `json:"id"`
	SourceIdentifier string                 `json:"sourceIdentifier"`
	Published        string                 `json:"published"`
	LastModified     string                 `json:"lastModified"`
	VulnStatus       string                 `json:"vulnStatus"`
	Descriptions     []DescriptionItem      `json:"descriptions"`
	Metrics          MetricsItem            `json:"metrics"`
	Configurations   []ConfigurationItem    `json:"configurations"`
	References       []ReferenceItem        `json:"references"`
	VulnerableConfig []VulnerableConfigItem `json:"vulnerableConfigurations,omitempty"`
	Solutions        []SolutionItem         `json:"solutions,omitempty"`
}

// DescriptionItem 描述项
type DescriptionItem struct {
	Lang  string `json:"lang"`
	Value string `json:"value"`
}

// MetricsItem 指标项
type MetricsItem struct {
	CvssMetricV31 []CvssMetricItem `json:"cvssMetricV31"`
	CvssMetricV30 []CvssMetricItem `json:"cvssMetricV30"`
	CvssMetricV2  []CvssMetricItem `json:"cvssMetricV2"`
}

// CvssMetricItem CVSS指标项
type CvssMetricItem struct {
	Source              string   `json:"source"`
	Type                string   `json:"type"`
	CvssData            CvssData `json:"cvssData"`
	ExploitabilityScore float64  `json:"exploitabilityScore"`
	ImpactScore         float64  `json:"impactScore"`
	BaseSeverity        string   `json:"baseSeverity"`
	BaseScore           float64  `json:"baseScore"`
}

// CvssData CVSS数据
type CvssData struct {
	Version               string  `json:"version"`
	VectorString          string  `json:"vectorString"`
	BaseScore             float64 `json:"baseScore"`
	BaseSeverity          string  `json:"baseSeverity"`
	AttackVector          string  `json:"attackVector,omitempty"`
	AttackComplexity      string  `json:"attackComplexity,omitempty"`
	PrivilegesRequired    string  `json:"privilegesRequired,omitempty"`
	UserInteraction       string  `json:"userInteraction,omitempty"`
	Scope                 string  `json:"scope,omitempty"`
	ConfidentialityImpact string  `json:"confidentialityImpact,omitempty"`
	IntegrityImpact       string  `json:"integrityImpact,omitempty"`
	AvailabilityImpact    string  `json:"availabilityImpact,omitempty"`
}

// ConfigurationItem 配置项
type ConfigurationItem struct {
	Nodes []NodeItem `json:"nodes"`
}

// NodeItem 节点项
type NodeItem struct {
	Operator string     `json:"operator"`
	Negate   bool       `json:"negate"`
	CPEMatch []CPEMatch `json:"cpeMatch"`
	Children []NodeItem `json:"children,omitempty"`
}

// CPEMatch CPE匹配
type CPEMatch struct {
	Vulnerable            bool   `json:"vulnerable"`
	CPE23URI              string `json:"cpe23Uri"`
	VersionStartIncluding string `json:"versionStartIncluding,omitempty"`
	VersionEndExcluding   string `json:"versionEndExcluding,omitempty"`
	VersionEndIncluding   string `json:"versionEndIncluding,omitempty"`
}

// ReferenceItem 参考项
type ReferenceItem struct {
	URL    string   `json:"url"`
	Source string   `json:"source"`
	Tags   []string `json:"tags"`
}

// VulnerableConfigItem 易受攻击的配置项
type VulnerableConfigItem struct {
	CPE      string `json:"cpe"`
	Versions []struct {
		Version string `json:"version"`
		Status  string `json:"status"`
	} `json:"versions"`
}

// SolutionItem 解决方案项
type SolutionItem struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// FrontendVulnerability 前端漏洞结构体
type FrontendVulnerability struct {
	CveId            string   `json:"cveId"`
	Title            string   `json:"title"`
	Description      string   `json:"description"`
	Severity         string   `json:"severity"`
	CvssScore        float64  `json:"cvssScore"`
	CvssVector       string   `json:"cvssVector"`
	CvssVersion      string   `json:"cvssVersion"`
	Published        string   `json:"published"`
	LastModified     string   `json:"lastModified"`
	Vendor           string   `json:"vendor"`
	Product          string   `json:"product"`
	AffectedVersions string   `json:"affectedVersions"`
	Solution         string   `json:"solution"`
	References       []string `json:"references"`
}

// 获取NVD API密钥环境变量
var nvdApiKey = "API_KEY" // 替换为实际的NVD API密钥

// 模拟数据存储
var vulnerabilities = []Vulnerability{}
var lastID = 0

// 初始化一些模拟数据
func init() {
	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)

	// 添加一些模拟数据
	vulnerabilities = append(vulnerabilities, Vulnerability{
		ID:               "1",
		CVEId:            "CVE-2023-1234",
		Title:            "SQL Injection in Example Application",
		Description:      "A SQL injection vulnerability in the login form allows attackers to bypass authentication.",
		Severity:         "high",
		CVSS:             8.5,
		CVSSVector:       "CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:H/I:H/A:H",
		Vendor:           "Example Corp",
		Product:          "Example App",
		AffectedVersions: "1.0 - 2.3",
		Solution:         "Upgrade to version 2.4 or apply the security patch.",
		Exploit:          "' OR 1=1 --",
		References:       []string{"https://example.com/advisory/123"},
		Tags:             []string{"injection", "authentication"},
		PublishedDate:    &yesterday,
		LastModifiedDate: &now,
	})

	vulnerabilities = append(vulnerabilities, Vulnerability{
		ID:               "2",
		CVEId:            "CVE-2023-5678",
		Title:            "Cross-site Scripting in Example Blog",
		Description:      "A stored XSS vulnerability in the comment section allows attackers to inject malicious scripts.",
		Severity:         "medium",
		CVSS:             6.5,
		CVSSVector:       "CVSS:3.1/AV:N/AC:L/PR:L/UI:R/S:U/C:H/I:H/A:L",
		Vendor:           "Example Corp",
		Product:          "Example Blog",
		AffectedVersions: "3.0 - 3.2",
		Solution:         "Upgrade to version 3.3 or apply the security patch.",
		Exploit:          "<script>alert('XSS')</script>",
		References:       []string{"https://example.com/advisory/456"},
		Tags:             []string{"xss", "injection"},
		PublishedDate:    &yesterday,
		LastModifiedDate: &now,
	})

	lastID = 2
}

// SearchVulnDatabase 搜索漏洞库
func (c *VulnDatabaseController) SearchVulnDatabase(ctx *gin.Context) {
	// 获取分页参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("perPage", "20"))
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	// 过滤和搜索参数
	query := ctx.Query("q")
	year := ctx.Query("year")
	severity := ctx.Query("severity")
	sortBy := ctx.DefaultQuery("sortBy", "publishedDate")
	sortOrder := ctx.DefaultQuery("sortOrder", "desc")

	// 进行过滤（这里简化处理，实际应用应进行更复杂的过滤）
	var filtered []Vulnerability
	for _, v := range vulnerabilities {
		// 实现简单搜索
		if query != "" {
			// 关键词搜索：标题、描述、CVE ID
			searchLower := strings.ToLower(query)
			titleLower := strings.ToLower(v.Title)
			descLower := strings.ToLower(v.Description)
			cveLower := strings.ToLower(v.CVEId)

			if !strings.Contains(titleLower, searchLower) &&
				!strings.Contains(descLower, searchLower) &&
				!strings.Contains(cveLower, searchLower) {
				continue
			}
		}

		// 严重程度过滤
		if severity != "" && v.Severity != severity {
			continue
		}

		// 年份过滤
		if year != "" {
			if v.PublishedDate == nil || !strings.Contains(v.CVEId, year) {
				continue
			}
		}

		// 通过所有过滤条件，添加到结果中
		filtered = append(filtered, v)
	}

	// 排序
	sort.Slice(filtered, func(i, j int) bool {
		var aVal, bVal interface{}

		// 根据排序字段获取值
		switch sortBy {
		case "publishedDate":
			aVal, bVal = filtered[i].PublishedDate, filtered[j].PublishedDate
		case "lastModifiedDate":
			aVal, bVal = filtered[i].LastModifiedDate, filtered[j].LastModifiedDate
		case "cveId":
			aVal, bVal = filtered[i].CVEId, filtered[j].CVEId
		case "cvss":
			aVal, bVal = filtered[i].CVSS, filtered[j].CVSS
		default:
			aVal, bVal = filtered[i].PublishedDate, filtered[j].PublishedDate
		}

		// 根据排序方向比较
		if sortOrder == "asc" {
			// 升序
			switch v1 := aVal.(type) {
			case *time.Time:
				v2 := bVal.(*time.Time)
				if v1 == nil {
					return true
				}
				if v2 == nil {
					return false
				}
				return v1.Before(*v2)
			case string:
				return v1 < bVal.(string)
			case float64:
				return v1 < bVal.(float64)
			default:
				return false
			}
		} else {
			// 降序
			switch v1 := aVal.(type) {
			case *time.Time:
				v2 := bVal.(*time.Time)
				if v1 == nil {
					return false
				}
				if v2 == nil {
					return true
				}
				return v1.After(*v2)
			case string:
				return v1 > bVal.(string)
			case float64:
				return v1 > bVal.(float64)
			default:
				return false
			}
		}
	})

	// 计算分页
	startIndex := (page - 1) * perPage
	endIndex := startIndex + perPage

	if startIndex >= len(filtered) {
		startIndex = 0
		endIndex = 0
	}

	if endIndex > len(filtered) {
		endIndex = len(filtered)
	}

	// 返回结果
	ctx.JSON(http.StatusOK, gin.H{
		"items":   filtered[startIndex:endIndex],
		"total":   len(filtered),
		"page":    page,
		"perPage": perPage,
	})
}

// GetVulnerabilityByCveID 获取单个漏洞详情
func (c *VulnDatabaseController) GetVulnerabilityByCveID(ctx *gin.Context) {
	cveId := ctx.Param("cveId")

	for _, v := range vulnerabilities {
		if v.CVEId == cveId {
			ctx.JSON(http.StatusOK, v)
			return
		}
	}

	ctx.JSON(http.StatusNotFound, gin.H{"error": "Vulnerability not found"})
}

// CreateVulnerability 创建新漏洞
func (c *VulnDatabaseController) CreateVulnerability(ctx *gin.Context) {
	var vuln Vulnerability
	if err := ctx.ShouldBindJSON(&vuln); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 生成ID
	lastID++
	vuln.ID = strconv.Itoa(lastID)

	// 设置时间
	now := time.Now()
	vuln.PublishedDate = &now
	vuln.LastModifiedDate = &now

	// 保存
	vulnerabilities = append(vulnerabilities, vuln)

	ctx.JSON(http.StatusCreated, vuln)
}

// UpdateVulnerability 更新漏洞
func (c *VulnDatabaseController) UpdateVulnerability(ctx *gin.Context) {
	cveId := ctx.Param("cveId")

	var update Vulnerability
	if err := ctx.ShouldBindJSON(&update); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, v := range vulnerabilities {
		if v.CVEId == cveId {
			// 保留ID和发布日期
			update.ID = v.ID
			update.PublishedDate = v.PublishedDate

			// 更新最后修改时间
			now := time.Now()
			update.LastModifiedDate = &now

			// 更新
			vulnerabilities[i] = update

			ctx.JSON(http.StatusOK, update)
			return
		}
	}

	ctx.JSON(http.StatusNotFound, gin.H{"error": "Vulnerability not found"})
}

// DeleteVulnerability 删除漏洞
func (c *VulnDatabaseController) DeleteVulnerability(ctx *gin.Context) {
	cveId := ctx.Param("cveId")

	for i, v := range vulnerabilities {
		if v.CVEId == cveId {
			// 删除
			vulnerabilities = append(vulnerabilities[:i], vulnerabilities[i+1:]...)

			ctx.JSON(http.StatusOK, gin.H{"message": "Vulnerability deleted"})
			return
		}
	}

	ctx.JSON(http.StatusNotFound, gin.H{"error": "Vulnerability not found"})
}

// convertToFrontendVulnerability 将NVD漏洞数据转换为前端格式
func convertToFrontendVulnerability(vuln VulnerabilityItem) FrontendVulnerability {
	cve := vuln.CVE

	// 提取标题和描述（优先中文，否则使用英文）
	title := ""
	description := ""
	for _, desc := range cve.Descriptions {
		if desc.Lang == "zh" {
			if title == "" && len(desc.Value) < 100 {
				title = desc.Value
			} else if description == "" {
				description = desc.Value
			}
		}
	}

	// 如果没有找到中文，则使用英文
	if title == "" || description == "" {
		for _, desc := range cve.Descriptions {
			if desc.Lang == "en" {
				if title == "" && len(desc.Value) < 100 {
					title = desc.Value
				} else if description == "" {
					description = desc.Value
				}
			}
		}
	}

	// 如果仍然没有标题，使用CVE ID作为标题
	if title == "" {
		title = cve.ID
	}

	// 获取CVSS信息
	severity := "UNKNOWN"
	cvssScore := 0.0
	cvssVector := ""
	cvssVersion := ""

	// 优先使用CVSS v3.1
	if len(cve.Metrics.CvssMetricV31) > 0 {
		metric := cve.Metrics.CvssMetricV31[0]
		severity = metric.BaseSeverity
		cvssScore = metric.BaseScore
		cvssVector = metric.CvssData.VectorString
		cvssVersion = metric.CvssData.Version
	} else if len(cve.Metrics.CvssMetricV30) > 0 {
		metric := cve.Metrics.CvssMetricV30[0]
		severity = metric.BaseSeverity
		cvssScore = metric.BaseScore
		cvssVector = metric.CvssData.VectorString
		cvssVersion = metric.CvssData.Version
	} else if len(cve.Metrics.CvssMetricV2) > 0 {
		metric := cve.Metrics.CvssMetricV2[0]
		severity = metric.BaseSeverity
		cvssScore = metric.BaseScore
		cvssVector = metric.CvssData.VectorString
		cvssVersion = metric.CvssData.Version
	}

	// 获取厂商和产品信息
	vendor := ""
	product := ""
	affectedVersions := ""

	if len(cve.Configurations) > 0 {
		for _, config := range cve.Configurations {
			if len(config.Nodes) > 0 {
				for _, node := range config.Nodes {
					processCPENodes(node, &vendor, &product, &affectedVersions)
				}
			}
		}
	}

	// 获取解决方案信息
	solution := ""
	if len(cve.Solutions) > 0 {
		solution = cve.Solutions[0].Value
	}

	// 获取参考链接
	references := make([]string, 0)
	for _, ref := range cve.References {
		references = append(references, ref.URL)
	}

	return FrontendVulnerability{
		CveId:            cve.ID,
		Title:            title,
		Description:      description,
		Severity:         severity,
		CvssScore:        cvssScore,
		CvssVector:       cvssVector,
		CvssVersion:      cvssVersion,
		Published:        cve.Published,
		LastModified:     cve.LastModified,
		Vendor:           vendor,
		Product:          product,
		AffectedVersions: affectedVersions,
		Solution:         solution,
		References:       references,
	}
}

// processCPENodes 处理CPE节点以提取厂商和产品信息
func processCPENodes(node NodeItem, vendor *string, product *string, affectedVersions *string) {
	// 处理当前节点的CPE匹配
	for _, cpeMatch := range node.CPEMatch {
		if cpeMatch.Vulnerable {
			cpeParts := strings.Split(cpeMatch.CPE23URI, ":")
			if len(cpeParts) > 5 {
				// CPE格式: cpe:2.3:a|o|h:vendor:product:version:...
				if *vendor == "" && cpeParts[3] != "*" {
					*vendor = cpeParts[3]
				}
				if *product == "" && cpeParts[4] != "*" {
					*product = cpeParts[4]
				}

				// 提取版本范围信息
				versionInfo := ""
				if cpeParts[5] != "*" {
					versionInfo = cpeParts[5]
				}

				if cpeMatch.VersionStartIncluding != "" {
					if versionInfo != "" {
						versionInfo = cpeMatch.VersionStartIncluding + " 及以上"
					} else {
						versionInfo = "≥ " + cpeMatch.VersionStartIncluding
					}
				}

				if cpeMatch.VersionEndExcluding != "" {
					if versionInfo != "" {
						versionInfo += " 且 < " + cpeMatch.VersionEndExcluding
					} else {
						versionInfo = "< " + cpeMatch.VersionEndExcluding
					}
				}

				if cpeMatch.VersionEndIncluding != "" {
					if versionInfo != "" {
						versionInfo += " 且 ≤ " + cpeMatch.VersionEndIncluding
					} else {
						versionInfo = "≤ " + cpeMatch.VersionEndIncluding
					}
				}

				if versionInfo != "" && *affectedVersions == "" {
					*affectedVersions = versionInfo
				}
			}
		}
	}

	// 递归处理子节点
	for _, child := range node.Children {
		processCPENodes(child, vendor, product, affectedVersions)
	}
}
