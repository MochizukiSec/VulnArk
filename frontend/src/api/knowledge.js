import request from '@/utils/request'

/**
 * 获取知识库列表
 * @param {Object} params 查询参数
 * @returns {Promise} promise
 */
export function getKnowledgeList(params) {
  return request({
    url: '/knowledge',
    method: 'get',
    params
  })
}

/**
 * 获取知识库详情
 * @param {String} id 知识库ID
 * @returns {Promise} promise
 */
export function getKnowledgeDetail(id) {
  return request({
    url: `/knowledge/${id}`,
    method: 'get'
  })
}

/**
 * 创建知识库
 * @param {Object} data 知识库数据
 * @returns {Promise} promise
 */
export function createKnowledge(data) {
  return request({
    url: '/knowledge',
    method: 'post',
    data
  })
}

/**
 * 更新知识库
 * @param {String} id 知识库ID
 * @param {Object} data 更新数据
 * @returns {Promise} promise
 */
export function updateKnowledge(id, data) {
  return request({
    url: `/knowledge/${id}`,
    method: 'put',
    data
  })
}

/**
 * 删除知识库
 * @param {String} id 知识库ID
 * @returns {Promise} promise
 */
export function deleteKnowledge(id) {
  return request({
    url: `/knowledge/${id}`,
    method: 'delete'
  })
}

/**
 * 获取知识库类型列表
 * @returns {Promise} promise
 */
export function getKnowledgeTypes() {
  return request({
    url: '/knowledge/types',
    method: 'get'
  })
}

/**
 * 获取知识库分类列表
 * @returns {Promise} promise
 */
export function getKnowledgeCategories() {
  return request({
    url: '/knowledge/categories',
    method: 'get'
  })
} 