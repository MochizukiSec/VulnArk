import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import axios from 'axios'
import moment from 'moment'

// 配置axios默认值
axios.defaults.baseURL = process.env.VUE_APP_API_URL || 'http://localhost:8000/api'
axios.interceptors.request.use(config => {
  const token = store.state.auth.token
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// 添加响应拦截器处理401未授权错误
axios.interceptors.response.use(
  response => response,
  error => {
    // 如果请求包含静默错误处理标记，则不执行重定向
    if (error.config && error.config._silentError) {
      console.log('静默处理错误请求')
      return Promise.reject(error);
    }
    
    // 获取当前路径，只在非首页上重定向
    const currentPath = window.location.pathname;
    console.log('当前路径:', currentPath);
    
    // 只有当收到真实的401响应，且不是访问首页时才执行重定向
    if (error.response && 
        error.response.status === 401 && 
        currentPath !== '/') {
      console.log('收到401响应，执行登出操作')
      store.dispatch('auth/logout')
      router.push('/login')
    }
    return Promise.reject(error)
  }
)

// 检查并清除无效的localStorage数据
const userStr = localStorage.getItem('user');
if (userStr === 'undefined' || (userStr && userStr.trim() === '')) {
  console.warn('发现无效的用户数据，正在清除...');
  localStorage.removeItem('user');
  localStorage.removeItem('token');
}

// 创建应用实例
const app = createApp(App)

// 注册全局实用工具
app.config.globalProperties.$moment = moment
app.config.globalProperties.$axios = axios

// 注册全局过滤器
app.config.globalProperties.$filters = {
  formatDate(value) {
    if (!value) return ''
    return moment(value).format('YYYY-MM-DD HH:mm')
  },
  formatDateOnly(value) {
    if (!value) return ''
    return moment(value).format('YYYY-MM-DD')
  },
  capitalize(value) {
    if (!value) return ''
    value = value.toString()
    return value.charAt(0).toUpperCase() + value.slice(1)
  },
  severityClass(severity) {
    switch (severity) {
      case 'critical': return 'danger'
      case 'high': return 'danger'
      case 'medium': return 'warning'
      case 'low': return 'primary'
      case 'info': return 'info'
      default: return 'info'
    }
  },
  statusClass(status) {
    switch (status) {
      case 'open': return 'danger'
      case 'in_progress': return 'warning'
      case 'resolved': return 'success'
      case 'closed': return 'info'
      case 'false_positive': return 'info'
      default: return 'info'
    }
  },
  statusText(status) {
    switch (status) {
      case 'open': return '开放'
      case 'in_progress': return '处理中'
      case 'resolved': return '已解决'
      case 'closed': return '已关闭'
      case 'false_positive': return '误报'
      default: return status
    }
  },
  severityText(severity) {
    switch (severity) {
      case 'critical': return '严重'
      case 'high': return '高危'
      case 'medium': return '中危'
      case 'low': return '低危'
      case 'info': return '信息'
      default: return severity
    }
  }
}

// 使用插件
app.use(store)
app.use(router)
app.use(ElementPlus)

// 挂载应用
app.mount('#app') 