import axios from 'axios'
import { ElMessage, ElMessageBox } from 'element-plus'
import store from '@/store'
import { getToken } from '@/utils/auth'

// 创建axios实例
const service = axios.create({
  baseURL: process.env.VUE_APP_BASE_API || '/api/v1', // url = base url + request url
  timeout: 15000 // 请求超时时间
})

// 请求拦截器
service.interceptors.request.use(
  config => {
    // 在发送请求之前设置token
    const token = getToken()
    if (token) {
      config.headers['Authorization'] = 'Bearer ' + token
      console.log('已为请求添加认证token')
    } else {
      console.warn('未找到认证token，请求可能会被拒绝')
    }
    
    // 确保Content-Type字段存在
    if (config.method.toLowerCase() === 'post' || config.method.toLowerCase() === 'put') {
      config.headers['Content-Type'] = config.headers['Content-Type'] || 'application/json'
    }
    
    // 详细日志
    console.log('准备发送请求:', {
      method: config.method,
      url: config.url,
      baseURL: config.baseURL,
      完整路径: config.baseURL + config.url,
      headers: {
        ...config.headers,
        Authorization: config.headers.Authorization ? 
          config.headers.Authorization.substring(0, 20) + '...' : 
          '未设置认证token'
      },
      data: config.data,
      params: config.params
    })
    
    return config
  },
  error => {
    console.log(error)
    return Promise.reject(error)
  }
)

// 响应拦截器
service.interceptors.response.use(
  response => {
    const res = response.data
    const url = response.config.url

    // 针对设置保存接口的特殊处理
    if (url && url.includes('/settings') && response.config.method === 'put') {
      // 设置保存接口直接返回成功响应，不按code判断错误
      return res
    }

    // 如果返回的状态码不是200，则判断为错误
    if (res.code !== 200) {
      ElMessage({
        message: res.message || '错误',
        type: 'error',
        duration: 5 * 1000
      })

      // 401: 未登录或token过期
      if (res.code === 401) {
        // 提示用户重新登录
        ElMessageBox.confirm('您已登出，可以取消继续留在该页面，或者重新登录', '确认登出', {
          confirmButtonText: '重新登录',
          cancelButtonText: '取消',
          type: 'warning'
        }).then(() => {
          store.dispatch('user/resetToken').then(() => {
            location.reload()
          })
        })
      }
      return Promise.reject(new Error(res.message || '错误'))
    } else {
      return res
    }
  },
  error => {
    console.log('响应错误：' + error)
    
    // 详细错误日志
    if (error.response) {
      console.error('错误响应详情:', {
        status: error.response.status,
        statusText: error.response.statusText,
        data: error.response.data,
        headers: error.response.headers,
        config: {
          method: error.response.config.method,
          url: error.response.config.url,
          baseURL: error.response.config.baseURL,
          完整URL: error.response.config.baseURL + error.response.config.url,
          data: error.response.config.data
        }
      })
      
      // 自动处理401错误
      if (error.response.status === 401) {
        // 清除用户状态并跳转到登录页
        store.dispatch('user/resetToken').then(() => {
          // 保存当前路径以便登录后重定向回来
          const currentPath = window.location.pathname
          window.location.href = `/login?redirect=${encodeURIComponent(currentPath)}`
        })
      }
    }
    
    ElMessage({
      message: error.message || '请求错误',
      type: 'error',
      duration: 5 * 1000
    })
    return Promise.reject(error)
  }
)

export default service 