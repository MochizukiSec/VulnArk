import request from '@/utils/request'

/**
 * 获取漏洞库列表
 * @param {Object} params 查询参数
 * @returns {Promise} promise
 */
export function getVulnDBList(params) {
  return request({
    url: '/vulndb',
    method: 'get',
    params
  })
}

/**
 * 获取漏洞库详情
 * @param {String} id 漏洞库条目ID
 * @returns {Promise} promise
 */
export function getVulnDBDetail(id) {
  return request({
    url: `/vulndb/id/${id}`,
    method: 'get'
  })
}

/**
 * 根据CVE获取漏洞库条目
 * @param {String} cve CVE编号
 * @returns {Promise} promise
 */
export function getVulnDBByCVE(cve) {
  return request({
    url: `/vulndb/cve/${cve}`,
    method: 'get'
  })
}

/**
 * 创建漏洞库条目
 * @param {Object} data 漏洞库条目数据
 * @returns {Promise} promise
 */
export function createVulnDBEntry(data) {
  // 转换字段以匹配后端期望的格式
  const serverData = {
    ...data,
    // 将remediation转换为solution
    solution: data.remediation,
    // 将affected_products数组转换为affected_systems字符串
    affected_systems: Array.isArray(data.affected_products) ? data.affected_products.join(',') : data.affected_products,
    // 将references数组转换为字符串
    references: Array.isArray(data.references) ? data.references.join(',') : data.references,
    // 将tags数组转换为字符串
    tags: Array.isArray(data.tags) ? data.tags.join(',') : data.tags,
    // 将last_modified_date转换为updated_date
    updated_date: data.last_modified_date
  }

  return request({
    url: '/vulndb',
    method: 'post',
    data: serverData
  })
}

/**
 * 更新漏洞库条目
 * @param {String} id 漏洞库条目ID
 * @param {Object} data 更新数据
 * @returns {Promise} promise
 */
export function updateVulnDBEntry(id, data) {
  // 转换字段以匹配后端期望的格式
  const serverData = {
    ...data,
    // 将remediation转换为solution
    solution: data.remediation,
    // 将affected_products数组转换为affected_systems字符串
    affected_systems: Array.isArray(data.affected_products) ? data.affected_products.join(',') : data.affected_products,
    // 将references数组转换为字符串
    references: Array.isArray(data.references) ? data.references.join(',') : data.references,
    // 将tags数组转换为字符串
    tags: Array.isArray(data.tags) ? data.tags.join(',') : data.tags,
    // 将last_modified_date转换为updated_date
    updated_date: data.last_modified_date
  }

  return request({
    url: `/vulndb/id/${id}`,
    method: 'put',
    data: serverData
  })
}

/**
 * 删除漏洞库条目
 * @param {String} id 漏洞库条目ID
 * @returns {Promise} promise
 */
export function deleteVulnDBEntry(id) {
  return request({
    url: `/vulndb/id/${id}`,
    method: 'delete'
  })
}

/**
 * 批量导入漏洞库条目
 * @param {FormData} formData 包含文件的表单数据
 * @returns {Promise} promise
 */
export function batchImportVulnDB(formData) {
  return request({
    url: '/vulndb/import',
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

/**
 * 获取漏洞严重程度列表
 * @returns {Array} 严重程度列表
 */
export function getVulnSeverities() {
  return [
    { value: 'critical', label: '严重' },
    { value: 'high', label: '高危' },
    { value: 'medium', label: '中危' },
    { value: 'low', label: '低危' },
    { value: 'info', label: '信息' }
  ]
} 