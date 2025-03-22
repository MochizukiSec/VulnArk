import request from '@/utils/request'

/**
 * 获取仪表盘统计数据
 * @returns {Promise} Promise对象
 */
export function getDashboardStats() {
  return request({
    url: '/dashboard/stats',
    method: 'get',
    params: {
      _t: Date.now() // 添加时间戳，避免缓存
    }
  })
}

/**
 * 获取漏洞趋势数据
 * @param {Object} params 查询参数，例如时间范围
 * @returns {Promise} Promise对象
 */
export function getVulnTrends(params) {
  return request({
    url: '/dashboard/vuln-trends',
    method: 'get',
    params: {
      ...params,
      _t: Date.now() // 添加时间戳，避免缓存
    }
  })
}

/**
 * 获取漏洞严重程度分布
 * @returns {Promise} Promise对象
 */
export function getSeverityDistribution() {
  return request({
    url: '/dashboard/severity-distribution',
    method: 'get',
    params: {
      _t: Date.now() // 添加时间戳，避免缓存
    }
  })
}

/**
 * 获取优先修复漏洞列表
 * @param {Object} params 查询参数
 * @returns {Promise} Promise对象
 */
export function getPriorityVulns(params) {
  return request({
    url: '/dashboard/priority-vulns',
    method: 'get',
    params: {
      ...params,
      _t: Date.now() // 添加时间戳，避免缓存
    }
  })
}

/**
 * 获取资产漏洞分布
 * @returns {Promise} Promise对象
 */
export function getAssetVulnDistribution() {
  return request({
    url: '/dashboard/asset-vuln-distribution',
    method: 'get',
    params: {
      _t: Date.now() // 添加时间戳，避免缓存
    }
  })
}

/**
 * 获取最近活动
 * @param {Number} limit 限制条数
 * @returns {Promise} Promise对象
 */
export function getRecentActivities(limit = 5) {
  return request({
    url: '/dashboard/recent-activities',
    method: 'get',
    params: { 
      limit,
      _t: Date.now() // 添加时间戳，避免缓存
    }
  })
} 