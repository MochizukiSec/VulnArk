import request from '@/utils/request'

/**
 * 获取系统设置
 * @returns {Promise} promise
 */
export function getSettings() {
  return request({
    url: '/settings',
    method: 'get'
  })
}

/**
 * 保存系统设置
 * @param {Object} data 设置数据
 * @returns {Promise} promise
 */
export function saveSettings(data) {
  return request({
    url: '/settings',
    method: 'put',
    data
  })
}

/**
 * 测试JIRA连接
 * @param {Object} data JIRA配置信息
 * @returns {Promise} promise
 */
export function testJiraConnection(data) {
  return request({
    url: '/settings/test/jira',
    method: 'post',
    data
  })
}

/**
 * 测试微信登录配置
 * @param {Object} data 微信登录配置信息
 * @returns {Promise} promise
 */
export function testWechatLogin(data) {
  return request({
    url: '/settings/test/wechat-login',
    method: 'post',
    data
  })
}

/**
 * 测试企业微信机器人
 * @param {Object} data 企业微信配置信息
 * @returns {Promise} promise
 */
export function testWorkWechatBot(data) {
  return request({
    url: '/settings/test/work-wechat',
    method: 'post',
    data
  })
}

/**
 * 测试飞书机器人
 * @param {Object} data 飞书配置信息
 * @returns {Promise} promise
 */
export function testFeishuBot(data) {
  return request({
    url: '/settings/test/feishu',
    method: 'post',
    data
  })
}

/**
 * 测试钉钉机器人
 * @param {Object} data 钉钉配置信息
 * @returns {Promise} promise
 */
export function testDingtalkBot(data) {
  return request({
    url: '/settings/test/dingtalk',
    method: 'post',
    data
  })
}

/**
 * 测试邮件发送
 * @param {Object} data 邮件配置信息
 * @returns {Promise} promise
 */
export function testEmailNotification(data) {
  return request({
    url: '/settings/test/email',
    method: 'post',
    data
  })
}

/**
 * 测试AI服务连接
 * @param {Object} data AI服务配置信息
 * @returns {Promise} promise
 */
export function testAiService(data) {
  return request({
    url: '/settings/test/ai',
    method: 'post',
    data
  })
}

/**
 * 测试漏洞库API连接
 * @param {Object} data 漏洞库API配置信息
 * @returns {Promise} promise
 */
export function testVulnDBConnection(data) {
  return request({
    url: '/settings/test/vulndb',
    method: 'post',
    data
  })
}

/**
 * 测试漏洞通知
 * @returns {Promise} promise
 */
export function testVulnerabilityNotification() {
  return request({
    url: '/settings/test/notification/vulnerability',
    method: 'post'
  })
}

/**
 * 添加扫描器
 * @param {Object} data - 扫描器数据
 * @returns {Promise}
 */
export function addScanner(data) {
  return request({
    url: '/api/v1/scanners',
    method: 'post',
    data
  })
}

/**
 * 更新扫描器
 * @param {Object} data - 扫描器数据
 * @returns {Promise}
 */
export function updateScanner(data) {
  return request({
    url: `/api/v1/scanners/${data.id}`,
    method: 'put',
    data
  })
}

/**
 * 删除扫描器
 * @param {Number} id - 扫描器ID
 * @returns {Promise}
 */
export function deleteScanner(id) {
  return request({
    url: `/api/v1/scanners/${id}`,
    method: 'delete'
  })
}

/**
 * 测试扫描器连接
 * @param {Object} data - 扫描器数据
 * @returns {Promise}
 */
export function testScannerConnection(data) {
  return request({
    url: '/api/v1/scanners/test',
    method: 'post',
    data
  })
} 