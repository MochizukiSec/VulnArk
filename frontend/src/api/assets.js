import axios from 'axios'

const API_URL = process.env.VUE_APP_API_URL || ''

// 資產管理API服務
export default {
  // 获取资产列表
  getAssets(params = {}) {
    return axios.get(`${API_URL}/assets`, { params })
  },

  // 获取单个资产详情
  getAssetById(id) {
    return axios.get(`${API_URL}/assets/${id}`)
  },

  // 创建新资产
  createAsset(assetData) {
    return axios.post(`${API_URL}/assets`, assetData)
  },

  // 更新资产信息
  updateAsset(id, assetData) {
    return axios.put(`${API_URL}/assets/${id}`, assetData)
  },

  // 删除资产
  deleteAsset(id) {
    return axios.delete(`${API_URL}/assets/${id}`)
  },

  // 获取资产关联的漏洞
  getAssetVulnerabilities(id) {
    return axios.get(`${API_URL}/assets/${id}/vulnerabilities`)
  },

  // 为资产添加漏洞
  addVulnerabilityToAsset(assetId, vulnId) {
    return axios.post(`${API_URL}/assets/${assetId}/vulnerabilities/${vulnId}`)
  },

  // 从资产中移除漏洞
  removeVulnerabilityFromAsset(assetId, vulnId) {
    return axios.delete(`${API_URL}/assets/${assetId}/vulnerabilities/${vulnId}`)
  },

  // 为资产添加备注
  addAssetNote(assetId, content) {
    return axios.post(`${API_URL}/assets/${assetId}/notes`, { content })
  }
} 