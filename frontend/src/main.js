import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import axios from 'axios'
import moment from 'moment'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'

// 配置axios默认值
// 在开发环境中使用完整的基础URL，在生产环境中使用相对路径
if (process.env.NODE_ENV === 'development') {
  // 开发环境使用完整URL，包括域名和端口，但不包含/api前缀
  axios.defaults.baseURL = process.env.VUE_APP_API_URL || 'http://localhost:8000'
  console.log('开发环境设置baseURL:', axios.defaults.baseURL)
} else {
  // 生产环境使用相对路径，自动匹配当前域名
  axios.defaults.baseURL = ''
  console.log('生产环境设置baseURL为空')
}

// 记录当前axios配置
console.log('系统环境:', process.env.NODE_ENV)
console.log('Axios基础URL配置:', axios.defaults.baseURL)

axios.interceptors.request.use(config => {
  const token = store.state.auth.token
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
    console.log(`请求添加令牌: ${config.url}`, {
      tokenExists: !!token,
      tokenPrefix: token ? token.substring(0, 10) + '...' : '',
      headers: config.headers
    })
  } else {
    console.warn(`请求无令牌: ${config.url}`, {
      storeAuthState: store.state.auth,
      localStorageToken: localStorage.getItem('token') ? '存在' : '不存在'
    })
    
    // 尝试从localStorage直接获取token作为备选方案
    const backupToken = localStorage.getItem('token')
    if (backupToken) {
      config.headers.Authorization = `Bearer ${backupToken}`
      console.log('已从localStorage直接获取令牌')
    }
  }
  return config
})

// 添加响应拦截器处理错误
axios.interceptors.response.use(
  response => response,
  error => {
    // 输出更详细的错误信息，帮助调试
    console.error('Axios错误:', {
      url: error.config?.url,
      method: error.config?.method,
      baseURL: axios.defaults.baseURL,
      fullURL: axios.defaults.baseURL + (error.config?.url || ''),
      status: error.response?.status,
      statusText: error.response?.statusText,
      data: error.response?.data,
      headers: error.config?.headers,
      message: error.message
    });
    
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

// 注册Element Plus图标
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

// 注册全局实用工具
app.config.globalProperties.$moment = moment
app.config.globalProperties.$axios = axios

// 注册全局过滤器
app.config.globalProperties.$filters = {
  formatDate(value) {
    if (!value) return '未知'
    return moment(value).format('YYYY-MM-DD HH:mm')
  },
  formatDateOnly(value) {
    if (!value) return '未知'
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