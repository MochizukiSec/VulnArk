import axios from 'axios'

const API_URL = process.env.VUE_APP_API_URL || ''

// 資產管理API服務
export default {
  // 获取资产列表
  getAssets(params = {}) {
    console.log('获取资产列表:', {
      baseURL: axios.defaults.baseURL,
      endpoint: `${API_URL}/api/assets`,
      params
    })
    return axios.get(`${API_URL}/api/assets`, { params })
  },

  // 获取单个资产详情
  getAssetById(id) {
    return axios.get(`${API_URL}/api/assets/${id}`)
  },

  // 创建新资产
  createAsset(assetData) {
    return axios.post(`${API_URL}/api/assets`, assetData)
  },

  // 更新资产信息
  updateAsset(id, assetData) {
    return axios.put(`${API_URL}/api/assets/${id}`, assetData)
  },

  // 删除资产
  deleteAsset(id) {
    return axios.delete(`${API_URL}/api/assets/${id}`)
  },

  // 获取资产关联的漏洞
  getAssetVulnerabilities(id) {
    return axios.get(`${API_URL}/api/assets/${id}/vulnerabilities`)
  },

  // 为资产添加漏洞
  addVulnerabilityToAsset(assetId, vulnId) {
    console.log(`添加漏洞到资产: ${assetId}, 漏洞ID: ${vulnId}`)
    return axios.post(`${API_URL}/api/assets/${assetId}/vulnerabilities`, {
      vulnerabilityId: vulnId
    })
  },

  // 从资产中移除漏洞
  removeVulnerabilityFromAsset(assetId, vulnId) {
    return axios.delete(`${API_URL}/api/assets/${assetId}/vulnerabilities/${vulnId}`)
  },

  // 为资产添加备注
  addAssetNote(assetId, content) {
    return axios.post(`${API_URL}/api/assets/${assetId}/notes`, { content })
  },
  
  // 批量导入资产
  importAssets(formData) {
    return axios.post(`${API_URL}/api/assets/import`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      },
      timeout: 60000, // 增加超时时间到60秒
      onUploadProgress: progressEvent => {
        // 可以在这里处理上传进度
        console.log('上传进度:', Math.round((progressEvent.loaded * 100) / progressEvent.total) + '%')
      }
    })
  },
  
  // 下载资产导入模板
  downloadImportTemplate() {
    return axios.get(`${API_URL}/api/assets/import-template`, {
      responseType: 'blob'
    })
  }
} 