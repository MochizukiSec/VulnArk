import request from '@/utils/request'

/**
 * 获取所有CI/CD集成配置
 * @returns {Promise}
 */
export function getIntegrations() {
  return request({
    url: '/integrations',
    method: 'get'
  })
}

/**
 * 创建CI/CD集成配置
 * @param {Object} data - 集成配置数据
 * @returns {Promise}
 */
export function createIntegration(data) {
  return request({
    url: '/integrations',
    method: 'post',
    data
  })
}

/**
 * 更新CI/CD集成配置
 * @param {Number} id - 集成配置ID
 * @param {Object} data - 集成配置数据
 * @returns {Promise}
 */
export function updateIntegration(id, data) {
  return request({
    url: `/integrations/${id}`,
    method: 'put',
    data
  })
}

/**
 * 删除CI/CD集成配置
 * @param {Number} id - 集成配置ID
 * @returns {Promise}
 */
export function deleteIntegration(id) {
  return request({
    url: `/integrations/${id}`,
    method: 'delete'
  })
}

/**
 * 重新生成API密钥
 * @param {Number} id - 集成配置ID
 * @returns {Promise}
 */
export function regenerateApiKey(id) {
  return request({
    url: `/integrations/${id}/api-key/regenerate`,
    method: 'post'
  })
}

/**
 * 获取集成历史记录
 * @param {Number} id - 集成配置ID
 * @returns {Promise}
 */
export function getIntegrationHistory(id) {
  return request({
    url: `/integrations/${id}/history`,
    method: 'get'
  })
}

/**
 * 切换集成启用状态
 * @param {Number} id - 集成配置ID
 * @param {Boolean} enabled - 是否启用
 * @returns {Promise}
 */
export function toggleIntegrationStatus(id, enabled) {
  return request({
    url: `/integrations/${id}/status`,
    method: 'put',
    data: { enabled }
  })
} 