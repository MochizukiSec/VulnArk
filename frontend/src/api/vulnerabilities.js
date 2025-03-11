import axios from 'axios'

const API_URL = process.env.VUE_APP_API_URL || ''

// 漏洞管理API服務
export default {
  // 获取漏洞列表
  getVulnerabilities(params = {}) {
    return axios.get(`${API_URL}/vulnerabilities`, { params })
  },

  // 根据ID获取漏洞详情
  getVulnerabilityById(id) {
    return axios.get(`${API_URL}/vulnerabilities/${id}`)
  },

  // 创建新漏洞
  createVulnerability(vulnerabilityData) {
    return axios.post(`${API_URL}/vulnerabilities`, vulnerabilityData)
  },

  // 更新漏洞信息
  updateVulnerability(id, vulnerabilityData) {
    return axios.put(`${API_URL}/vulnerabilities/${id}`, vulnerabilityData)
  },

  // 删除漏洞
  deleteVulnerability(id) {
    return axios.delete(`${API_URL}/vulnerabilities/${id}`)
  },

  // 导入漏洞
  importVulnerabilities(importData) {
    return axios.post(`${API_URL}/vulnerabilities/import`, importData)
  },

  // 搜索漏洞
  searchVulnerabilities(params = {}) {
    return axios.get(`${API_URL}/vulnerabilities`, { params })
  }
} 