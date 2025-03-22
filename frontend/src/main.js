import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import ElementPlus from 'element-plus'
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import 'element-plus/dist/index.css'
import './assets/styles/index.scss'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import Cookies from 'js-cookie'
import { v4 as uuidv4 } from 'uuid'
import { getToken } from '@/utils/auth'

// 添加crypto.randomUUID polyfill
if (!crypto.randomUUID) {
  crypto.randomUUID = () => uuidv4()
}

// 固定使用中文
document.title = 'VulnArk - 漏洞管理平台'
document.documentElement.lang = 'zh-CN'

// 创建模拟的i18n对象，用于解决使用了useI18n()的页面报错
const mockI18n = {
  global: {
    t: (key) => {
      // 提取最后一个部分作为中文文本的最佳猜测
      const parts = key.split('.');
      return parts[parts.length - 1] || key;
    }
  }
}

// 尝试恢复用户状态
const initUserState = async () => {
  const token = getToken()
  if (token && !store.getters.userInfo) {
    console.log('应用启动：检测到token，尝试获取用户信息')
    try {
      await store.dispatch('user/getUserInfo')
      console.log('应用启动：用户信息获取成功')
    } catch (error) {
      console.error('应用启动：获取用户信息失败', error)
      await store.dispatch('user/logout')
      
      // 检查当前路径是否需要登录
      const currentPath = window.location.pathname
      if (currentPath !== '/login' && currentPath !== '/404') {
        console.log('应用启动：用户未登录，跳转到登录页')
        window.location.href = `/login?redirect=${encodeURIComponent(currentPath)}`
        return false
      }
    }
  } else if (!token) {
    // 检查当前路径是否需要登录
    const currentPath = window.location.pathname
    if (currentPath !== '/login' && currentPath !== '/404') {
      console.log('应用启动：用户未登录，跳转到登录页')
      window.location.href = `/login?redirect=${encodeURIComponent(currentPath)}`
      return false
    }
  }
  return true
}

// 创建应用实例
const app = createApp(App)

// 注册ElementPlus图标
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

// 添加模拟的i18n
app.config.globalProperties.$i18n = mockI18n
app.provide('i18n', {
  t: mockI18n.global.t
})

// 配置ElementPlus语言
const locale = zhCn

app.use(store)
   .use(router)
   .use(ElementPlus, {
     size: 'medium',
     locale
   })

// 初始化用户状态
initUserState().then(() => {
  app.mount('#app')
  console.log('应用启动完成')
}).catch(error => {
  console.error('应用启动错误:', error)
  app.mount('#app')
})