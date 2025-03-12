import axios from 'axios'

const API_URL = process.env.VUE_APP_API_URL || ''

// 漏洞管理API服務
export default {
  // 获取漏洞列表
  getVulnerabilities(params = {}) {
    return axios.get(`${API_URL}/api/vulnerabilities`, { params })
  },

  // 根据ID获取漏洞详情
  getVulnerabilityById(id) {
    return axios.get(`${API_URL}/api/vulnerabilities/${id}`)
  },

  // 创建新漏洞
  createVulnerability(vulnerabilityData) {
    return axios.post(`${API_URL}/api/vulnerabilities`, vulnerabilityData)
  },

  // 更新漏洞信息
  updateVulnerability(id, vulnerabilityData) {
    return axios.put(`${API_URL}/api/vulnerabilities/${id}`, vulnerabilityData)
  },

  // 删除漏洞
  deleteVulnerability(id) {
    return axios.delete(`${API_URL}/api/vulnerabilities/${id}`)
  },

  // 导入漏洞
  importVulnerabilities(importData) {
    return axios.post(`${API_URL}/api/vulnerabilities/import`, importData)
  },
  
  // 从漏洞库导入单个漏洞
  importFromVulnDatabase(vulnerabilityData) {
    // 确保使用存在的后端API端点
    return axios.post(`${API_URL}/api/vulnerabilities/import-from-vulndb`, {
      vulnerability: vulnerabilityData
    })
  },

  // 搜索漏洞
  searchVulnerabilities(params = {}) {
    console.log("搜索漏洞库原始参数:", params)
    
    // 构建请求参数 - 将searchTerm转换为keyword
    const requestParams = {}
    
    // 处理searchTerm参数 - 转换为keyword
    if (params.searchTerm) {
      requestParams.keyword = params.searchTerm
    } else {
      requestParams.keyword = params.keyword || ''
    }
    
    // 处理limit参数 - 转换为pageSize
    if (params.limit) {
      requestParams.pageSize = parseInt(params.limit)
    } else {
      requestParams.pageSize = params.pageSize || params.perPage || 10
    }
    
    // 其他参数
    requestParams.cveId = params.cveId || ''
    requestParams.severity = params.severity || ''
    requestParams.vendor = params.vendor || ''
    requestParams.product = params.product || ''
    requestParams.year = params.year || ''
    requestParams.page = params.page || 1
    requestParams.sortBy = params.sortBy || 'publishedDate'
    requestParams.sortOrder = params.sortOrder || 'desc'
    
    console.log("处理后的请求参数:", requestParams)
    
    // 注意：问题在于请求发送到了不正确的端点
    // 1. 修正：确保使用完整URL路径，包含/api前缀
    // 2. 因为这是漏洞库API，而不是一般漏洞API，所以应该使用/api/vulndatabase
    const apiEndpoint = '/api/vulndatabase'
    console.log("使用API端点:", apiEndpoint)
    
    // 注意：axios在配置了baseURL的情况下，也接受相对路径
    return axios.get(apiEndpoint, { 
      params: requestParams,
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json'
      }
    });
  },
  
  // 获取模拟漏洞数据（当后端API不可用时使用）
  getMockVulnerabilities(keyword = '') {
    console.log("使用模拟数据替代后端API, 关键词:", keyword)
    
    // 创建模拟数据
    const mockData = [
      {
        cveId: "CVE-2023-1234",
        title: "SQL注入漏洞 (模拟数据)",
        description: "这是一个模拟的SQL注入漏洞，用于测试",
        severity: "high",
        cvss: 8.5,
        vendor: "测试厂商",
        product: "测试产品"
      },
      {
        cveId: "CVE-2023-5678", 
        title: "跨站脚本漏洞 (模拟数据)",
        description: "这是一个模拟的XSS漏洞，用于测试",
        severity: "medium",
        cvss: 6.5,
        vendor: "测试厂商",
        product: "测试产品"
      },
      {
        cveId: "CVE-2023-9012", 
        title: "命令注入漏洞 (模拟数据)",
        description: "这是一个模拟的命令注入漏洞，用于测试",
        severity: "critical",
        cvss: 9.5,
        vendor: "测试厂商",
        product: "测试产品"
      }
    ]
    
    // 根据关键字过滤
    const filteredData = keyword 
      ? mockData.filter(v => 
          v.title.toLowerCase().includes(keyword.toLowerCase()) || 
          v.cveId.toLowerCase().includes(keyword.toLowerCase()) ||
          v.description.toLowerCase().includes(keyword.toLowerCase())
        )
      : mockData
    
    // 返回一个Promise以模拟API调用
    return Promise.resolve({
      data: {
        items: filteredData,
        total: filteredData.length
      }
    })
  }
} 