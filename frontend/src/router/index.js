import { createRouter, createWebHistory } from 'vue-router'
import store from '../store'

// 懒加载组件
const Home = () => import('../views/Home.vue')
const Login = () => import('../views/auth/Login.vue')
const Register = () => import('../views/auth/Register.vue')
const Dashboard = () => import('../views/Dashboard.vue')
const Reports = () => import('../views/reports/Reports.vue')
const UserManagement = () => import('../views/admin/UserManagement.vue')
const NotFound = () => import('../views/NotFound.vue')
const AppLayout = () => import('../components/layout/AppLayout.vue')
const VulnerabilityCreate = () => import('../views/vulnerabilities/VulnerabilityCreate.vue')
const VulnerabilityImport = () => import('../views/vulnerabilities/VulnerabilityImport.vue')

// 将未使用的导入改为使用现有的文件
const UserProfile = () => import('../views/user/UserProfile.vue')
const VulnList = () => import('../views/vulnerabilities/VulnerabilityList.vue')
const VulnDetail = () => import('../views/vulnerabilities/VulnerabilityDetails.vue')
const ReportGeneration = () => import('../views/reports/Reports.vue')
const ReportHistory = () => import('../views/reports/Reports.vue')
const AIAnalysis = () => import('../views/ai/AIAnalysis.vue')

// 路由配置
const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home,
    meta: { requiresAuth: false, title: '首页' }
  },
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: { requiresAuth: false, title: '登录' }
  },
  {
    path: '/register',
    name: 'Register',
    component: Register,
    meta: { requiresAuth: false, title: '注册' }
  },
  // 嵌套在AppLayout内的路由
  {
    path: '/',
    component: AppLayout,
    meta: { requiresAuth: true },
    children: [
      {
        path: '/dashboard',
        name: 'Dashboard',
        component: Dashboard,
        meta: { requiresAuth: true, title: '仪表盘' }
      },
      {
        path: '/vulnerabilities',
        name: 'VulnList',
        component: VulnList,
        meta: { requiresAuth: true, title: '漏洞管理' }
      },
      {
        path: '/vulnerabilities/create',
        name: 'VulnerabilityCreate',
        component: VulnerabilityCreate,
        meta: { requiresAuth: true, title: '创建漏洞' }
      },
      {
        path: '/vulnerabilities/import',
        name: 'VulnerabilityImport',
        component: VulnerabilityImport,
        meta: { requiresAuth: true, title: '导入漏洞' }
      },
      {
        path: '/vulnerabilities/:id',
        name: 'VulnDetail',
        component: VulnDetail,
        meta: { requiresAuth: true, title: '漏洞详情' },
        props: true
      },
      {
        path: '/reports',
        name: 'Reports',
        component: Reports,
        meta: { requiresAuth: true, title: '报告' }
      },
      {
        path: '/profile',
        name: 'UserProfile',
        component: UserProfile,
        meta: { requiresAuth: true, title: '个人资料' }
      },
      {
        path: '/admin/users',
        name: 'UserManagement',
        component: UserManagement,
        meta: { requiresAuth: true, title: '用户管理' }
      },
      {
        path: '/reports/generate',
        name: 'ReportGeneration',
        component: ReportGeneration,
        meta: { requiresAuth: true, title: '生成报告' }
      },
      {
        path: '/reports/history',
        name: 'ReportHistory',
        component: ReportHistory,
        meta: { requiresAuth: true, title: '报告历史' }
      },
      {
        path: '/ai-analysis',
        name: 'AIAnalysis',
        component: AIAnalysis,
        meta: { requiresAuth: true, title: 'AI智能分析' }
      }
    ]
  },
  {
    path: '/:catchAll(.*)',
    name: 'NotFound',
    component: NotFound,
    meta: { title: '页面未找到' }
  }
]

// 创建路由器
const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
  scrollBehavior() {
    // 始终滚动到顶部
    return { top: 0 }
  }
})

// 路由守卫，检查认证和权限
router.beforeEach((to, from, next) => {
  // 增加调试输出
  console.log('路由变化:', { 
    to: to.path, 
    name: to.name, 
    requiresAuth: to.meta.requiresAuth,
    pathname: window.location.pathname
  });
  
  // 设置页面标题
  document.title = `${to.meta.title || '漏洞管理平台'} - 漏洞管理系统`
  
  // 检查用户是否已登录
  const isAuthenticated = store.getters['auth/isAuthenticated']
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth)
  const requiresAdmin = to.matched.some(record => record.meta.requiresAdmin)
  
  // 如果路由是首页，始终允许访问
  if (to.path === '/' || to.name === 'Home') {
    console.log('检测到首页路由，允许访问');
    next()
    return
  }
  
  // 处理需要认证的路由
  if (requiresAuth && !isAuthenticated) {
    next({ name: 'Login', query: { redirect: to.fullPath } })
    return
  }
  
  // 处理需要管理员权限的路由
  if (requiresAdmin && !store.getters['auth/isAdmin']) {
    next({ name: 'Dashboard' })
    return
  }
  
  // 如果用户已登录并尝试访问登录/注册页，重定向到仪表盘
  if (isAuthenticated && (to.name === 'Login' || to.name === 'Register')) {
    next({ name: 'Dashboard' })
    return
  }
  
  // 继续导航
  next()
})

export default router 