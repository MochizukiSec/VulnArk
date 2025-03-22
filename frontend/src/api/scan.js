import request from '@/utils/request'

// 获取扫描任务列表
export function getScanTasks(params) {
  return request({
    url: '/scans',
    method: 'get',
    params
  })
}

// 获取扫描任务详情
export function getScanTask(id) {
  return request({
    url: `/scans/${id}`,
    method: 'get'
  })
}

// 创建扫描任务
export function createScanTask(data) {
  return request({
    url: '/scans',
    method: 'post',
    data
  })
}

// 更新扫描任务
export function updateScanTask(id, data) {
  return request({
    url: `/scans/${id}`,
    method: 'put',
    data
  })
}

// 删除扫描任务
export function deleteScanTask(id) {
  return request({
    url: `/scans/${id}`,
    method: 'delete'
  })
}

// 启动扫描任务
export function startScanTask(id) {
  return request({
    url: `/scans/${id}/start`,
    method: 'post'
  })
}

// 取消扫描任务
export function cancelScanTask(id) {
  return request({
    url: `/scans/${id}/cancel`,
    method: 'post'
  })
}

// 获取扫描结果列表
export function getScanResults(id, params) {
  return request({
    url: `/scans/${id}/results`,
    method: 'get',
    params
  })
}

// 导入扫描结果
export function importScanResults(id, data) {
  return request({
    url: `/scans/${id}/import`,
    method: 'post',
    data
  })
}

/**
 * 获取扫描器列表
 * @returns {Promise} promise
 */
export function getScanners() {
  return request({
    url: '/api/v1/scanners',
    method: 'get'
  })
} 