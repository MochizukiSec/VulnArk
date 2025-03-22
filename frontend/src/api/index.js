import request from '@/utils/request'

// 通用API接口，提供基本HTTP方法
const api = {
  // GET请求
  get(url, config) {
    return request({
      url,
      method: 'get',
      ...config
    })
  },
  
  // POST请求
  post(url, data, config) {
    return request({
      url,
      method: 'post',
      data,
      ...config
    })
  },
  
  // PUT请求
  put(url, data, config) {
    return request({
      url,
      method: 'put',
      data,
      ...config
    })
  },
  
  // DELETE请求
  delete(url, config) {
    return request({
      url,
      method: 'delete',
      ...config
    })
  }
}

export default api 