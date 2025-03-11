package routes

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"vuln-management/config"
	"vuln-management/models"
	"vuln-management/services"
)

// VulnDatabaseController 漏洞库控制器
type VulnDatabaseController struct {
	service *services.VulnDatabaseService
}

// NewVulnDatabaseController 创建漏洞库控制器实例
func NewVulnDatabaseController(service *services.VulnDatabaseService) *VulnDatabaseController {
	return &VulnDatabaseController{
		service: service,
	}
}

// 初始化一些示例数据
var initialVulnData = []models.VulnDBEntry{
	{
		CveId:            "CVE-2023-1234",
		Title:            "SQL注入漏洞",
		Description:      "该漏洞允许攻击者通过构造特殊的SQL查询语句注入恶意代码，从而获取敏感数据。",
		Cvss:             8.5,
		Severity:         "high",
		AffectedSystems:  "Windows, Linux",
		Solution:         "升级到最新版本或应用补丁。使用参数化查询或预处理语句防止SQL注入。",
		PublishedDate:    "2023-06-15",
		LastModifiedDate: "2023-07-20",
		References:       []string{"https://example.com/cve-2023-1234", "https://nvd.nist.gov/vuln/detail/CVE-2023-1234"},
		Tags:             []string{"injection", "database"},
	},
	{
		CveId:            "CVE-2023-5678",
		Title:            "跨站脚本(XSS)漏洞",
		Description:      "此漏洞允许攻击者在网页中注入恶意JavaScript代码，可能导致会话劫持或信息窃取。",
		Cvss:             6.5,
		Severity:         "medium",
		AffectedSystems:  "Web应用",
		Solution:         "对输入进行验证和转义。使用内容安全策略(CSP)限制脚本执行。",
		PublishedDate:    "2023-08-10",
		LastModifiedDate: "2023-09-05",
		References:       []string{"https://example.com/cve-2023-5678"},
		Tags:             []string{"xss", "web"},
	},
	{
		CveId:            "CVE-2023-9012",
		Title:            "远程代码执行漏洞",
		Description:      "严重的漏洞允许远程攻击者在目标系统上执行任意代码，可能导致完全系统接管。",
		Cvss:             9.8,
		Severity:         "critical",
		AffectedSystems:  "Linux, Windows, macOS",
		Solution:         "立即更新到最新版本。临时缓解措施包括禁用相关服务或限制访问。",
		PublishedDate:    "2023-10-20",
		LastModifiedDate: "2023-10-25",
		References:       []string{"https://example.com/cve-2023-9012", "https://cert.org/advisories/..."},
		Tags:             []string{"rce", "critical"},
	},
}

// SearchVulnDatabase 搜索漏洞库
func (c *VulnDatabaseController) SearchVulnDatabase(ctx *gin.Context) {
	// 获取查询参数
	query := ctx.Query("q")
	severity := ctx.Query("severity")
	yearStr := ctx.Query("year")
	sortBy := ctx.DefaultQuery("sortBy", "publishedDate")
	sortOrder := ctx.DefaultQuery("sortOrder", "desc")
	pageStr := ctx.DefaultQuery("page", "1")
	perPageStr := ctx.DefaultQuery("perPage", "20")

	// 解析分页参数
	page, _ := strconv.Atoi(pageStr)
	perPage, _ := strconv.Atoi(perPageStr)

	if page < 1 {
		page = 1
	}

	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	// 使用服务层搜索漏洞
	items, total, err := c.service.SearchVulnerabilities(query, severity, yearStr, sortBy, sortOrder, page, perPage)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "搜索漏洞数据失败: " + err.Error(),
		})
		return
	}

	// 返回结果
	ctx.JSON(http.StatusOK, gin.H{
		"items": items,
		"total": total,
		"page":  page,
	})
}

// GetVulnerabilityByCveID 根据CVE ID获取漏洞详情
func (c *VulnDatabaseController) GetVulnerabilityByCveID(ctx *gin.Context) {
	cveId := ctx.Param("cveId")

	// 使用服务层获取漏洞详情
	vuln, err := c.service.GetVulnerabilityByCveID(cveId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "漏洞不存在: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, vuln)
}

// CreateVulnerability 创建新漏洞
func (c *VulnDatabaseController) CreateVulnerability(ctx *gin.Context) {
	var newVuln models.VulnDBEntry

	if err := ctx.ShouldBindJSON(&newVuln); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求数据: " + err.Error(),
		})
		return
	}

	// 验证必填字段
	if newVuln.Title == "" || newVuln.Description == "" || newVuln.Severity == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "标题、描述和严重程度为必填项",
		})
		return
	}

	// 如果没有设置最后更新日期，则使用发布日期
	if newVuln.LastModifiedDate == "" {
		newVuln.LastModifiedDate = newVuln.PublishedDate
	}

	// 使用服务层创建漏洞
	createdVuln, err := c.service.CreateVulnerability(newVuln)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "创建漏洞失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, createdVuln)
}

// UpdateVulnerability 更新漏洞信息
func (c *VulnDatabaseController) UpdateVulnerability(ctx *gin.Context) {
	cveId := ctx.Param("cveId")

	var updateData models.VulnDBEntry
	if err := ctx.ShouldBindJSON(&updateData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求数据: " + err.Error(),
		})
		return
	}

	// 使用服务层更新漏洞
	updatedVuln, err := c.service.UpdateVulnerability(cveId, updateData)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "更新漏洞失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, updatedVuln)
}

// DeleteVulnerability 删除漏洞
func (c *VulnDatabaseController) DeleteVulnerability(ctx *gin.Context) {
	cveId := ctx.Param("cveId")

	// 使用服务层删除漏洞
	err := c.service.DeleteVulnerability(cveId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "删除漏洞失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "漏洞已成功删除",
	})
}

// JWTAuth 临时中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

func SetupRouter(r *gin.Engine) {
	// 创建服务和控制器
	db := config.Database
	vulnDBService := services.NewVulnDatabaseService(db)
	vulnDBController := NewVulnDatabaseController(vulnDBService)

	// 导入初始数据
	if err := vulnDBService.ImportInitialData(initialVulnData); err != nil {
		fmt.Println("导入初始漏洞数据失败:", err)
	}

	// 为了方便测试，暂时不使用认证中间件
	apiGroup := r.Group("/api")
	{
		// 漏洞库相关路由
		apiGroup.GET("/vulndatabase", vulnDBController.SearchVulnDatabase)
		apiGroup.GET("/vulndatabase/:cveId", vulnDBController.GetVulnerabilityByCveID)
		apiGroup.POST("/vulndatabase", vulnDBController.CreateVulnerability)
		apiGroup.PUT("/vulndatabase/:cveId", vulnDBController.UpdateVulnerability)
		apiGroup.DELETE("/vulndatabase/:cveId", vulnDBController.DeleteVulnerability)

		// 处理前端重复的/api前缀
		apiGroup.GET("/api/vulndatabase", vulnDBController.SearchVulnDatabase)
		apiGroup.GET("/api/vulndatabase/:cveId", vulnDBController.GetVulnerabilityByCveID)
		apiGroup.POST("/api/vulndatabase", vulnDBController.CreateVulnerability)
		apiGroup.PUT("/api/vulndatabase/:cveId", vulnDBController.UpdateVulnerability)
		apiGroup.DELETE("/api/vulndatabase/:cveId", vulnDBController.DeleteVulnerability)
	}
}
