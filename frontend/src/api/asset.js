import request from '@/utils/request'

/**
 * 获取资产列表
 * @param {Object} params 查询参数
 * @returns {Promise} promise
 */
export function getAssetList(params) {
  return request({
    url: '/assets',
    method: 'get',
    params
  })
}

/**
 * 获取资产详情
 * @param {String} id 资产ID
 * @returns {Promise} promise
 */
export function getAssetDetail(id) {
  return request({
    url: `/assets/${id}`,
    method: 'get'
  })
}

/**
 * 获取资产关联的漏洞
 * @param {String} assetId 资产ID
 * @param {Object} params 查询参数
 * @returns {Promise} promise
 */
export function getAssetVulnerabilities(assetId, params) {
  return request({
    url: `/assets/${assetId}/vulnerabilities`,
    method: 'get',
    params
  })
}

/**
 * 添加资产
 * @param {Object} data 资产数据
 * @returns {Promise} promise
 */
export function addAsset(data) {
  return request({
    url: '/assets',
    method: 'post',
    data
  })
}

/**
 * 更新资产
 * @param {String} id 资产ID
 * @param {Object} data 资产数据
 * @returns {Promise} promise
 */
export function updateAsset(id, data) {
  return request({
    url: `/assets/${id}`,
    method: 'put',
    data
  })
}

/**
 * 删除资产
 * @param {String} id 资产ID
 * @returns {Promise} promise
 */
export function deleteAsset(id) {
  return request({
    url: `/assets/${id}`,
    method: 'delete'
  })
}

/**
 * 批量导入资产
 * @param {FormData} formData 包含文件的表单数据
 * @returns {Promise} promise
 */
export function batchImportAssets(formData) {
  return request({
    url: '/assets/import',
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

/**
 * 导出资产列表
 * @param {Object} params 查询参数
 * @returns {Promise} promise
 */
export function exportAssets(params) {
  return request({
    url: '/assets/export',
    method: 'get',
    params,
    responseType: 'blob'
  })
}

/**
 * 获取所有资产（简化列表，用于关联选择）
 * @returns {Promise} promise
 */
export function getAllAssets() {
  return request({
    url: '/assets',
    method: 'get',
    params: {
      page_size: 100 // 获取较多记录用于选择
    }
  })
}

/**
 * 批量删除资产
 * @param {Array} ids 资产ID数组
 * @returns {Promise} promise
 */
export function batchDeleteAssets(ids) {
  return request({
    url: '/assets/batch-delete',
    method: 'post',
    data: { ids }
  })
} 