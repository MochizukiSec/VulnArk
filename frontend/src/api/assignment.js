import request from '@/utils/request'

/**
 * 分配漏洞给用户
 * @param {Object} data 分配数据
 * @returns {Promise} promise
 */
export function assignVulnerability(vulnId, data) {
  // 确保vulnId是一个数字
  const id = parseInt(vulnId, 10);
  if (isNaN(id)) {
    return Promise.reject(new Error('无效的漏洞ID'));
  }
  
  // 构建完整URL以便调试
  const baseURL = process.env.VUE_APP_BASE_API || '/api/v1';
  const url = `/vulnerabilities/${id}/assign`;
  const fullURL = baseURL + url;
  
  console.log(`准备发送漏洞分派请求:`, {
    id,
    baseURL,
    url,
    fullURL,
    data
  });
  
  return request({
    url: `/vulnerabilities/${id}/assign`,
    method: 'post',
    data
  })
}

/**
 * 获取漏洞的分配历史
 * @param {number} vulnId 漏洞ID
 * @returns {Promise} promise
 */
export function getVulnerabilityAssignments(vulnId) {
  return request({
    url: `/vulnerabilities/${vulnId}/assignments`,
    method: 'get'
  })
}

/**
 * 获取分配任务详情
 * @param {number} assignmentId 分配ID
 * @returns {Promise} promise
 */
export function getAssignmentDetails(assignmentId) {
  return request({
    url: `/assignments/${assignmentId}`,
    method: 'get'
  })
}

/**
 * 更新分配任务状态
 * @param {number} assignmentId 分配ID
 * @param {Object} data 状态数据
 * @returns {Promise} promise
 */
export function updateAssignmentStatus(assignmentId, data) {
  return request({
    url: `/assignments/${assignmentId}/status`,
    method: 'put',
    data
  })
}

/**
 * 获取我的分配任务列表
 * @param {Object} params 查询参数
 * @returns {Promise} promise
 */
export function getMyAssignments(params) {
  return request({
    url: '/assignments/my',
    method: 'get',
    params
  })
}

/**
 * 获取所有分配任务列表
 * @param {Object} params 查询参数
 * @returns {Promise} promise
 */
export function listAssignments(params) {
  return request({
    url: '/assignments',
    method: 'get',
    params
  })
}

/**
 * 删除分配任务
 * @param {number} assignmentId 分配ID
 * @returns {Promise} promise
 */
export function deleteAssignment(assignmentId) {
  return request({
    url: `/assignments/${assignmentId}`,
    method: 'delete'
  })
} 