import { createRouter, createWebHistory } from 'vue-router'
import store from '../store'
import { getToken } from '../utils/auth'
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'

// 路由配置
const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue'),
    meta: { title: 'login.title', requiresAuth: false }
  },
  
  {
    path: '/',
    component: () => import('../layout/Layout.vue'),
    redirect: '/dashboard',
    meta: { requiresAuth: true },
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('../views/Dashboard.vue'),
        meta: { title: 'menu.dashboard', icon: 'Odometer' }
      },
      {
        path: 'vulnerabilities',
        name: 'Vulnerabilities',
        component: () => import('../views/vulnerability/Index.vue'),
        meta: { title: 'menu.vulnerabilities', icon: 'warning' }
      },
      {
        path: 'vulnerabilities/add',
        name: 'AddVulnerability',
        component: () => import('../views/vulnerability/Add.vue'),
        meta: { title: 'vulnerability.add', activeMenu: '/vulnerabilities' }
      },
      {
        path: 'vulnerabilities/:id',
        name: 'VulnerabilityDetail',
        component: () => import('../views/vulnerability/Detail.vue'),
        meta: { title: 'vulnerability.detail', activeMenu: '/vulnerabilities' }
      },
      {
        path: 'assets',
        name: 'Assets',
        component: () => import('../views/asset/Index.vue'),
        meta: { title: 'menu.assets', icon: 'Monitor' }
      },
      {
        path: 'assets/add',
        name: 'AddAsset',
        component: () => import('../views/asset/Add.vue'),
        meta: { title: 'asset.add', activeMenu: '/assets' }
      },
      {
        path: 'assets/detail/:id',
        name: 'AssetDetail',
        component: () => import('../views/asset/Detail.vue'),
        meta: { title: 'asset.detail', activeMenu: '/assets' }
      },
      {
        path: 'assets/edit/:id',
        name: 'AssetEdit',
        component: () => import('../views/asset/Edit.vue'),
        meta: { title: 'asset.edit', activeMenu: '/assets' }
      },
      {
        path: 'knowledge',
        name: 'Knowledge',
        component: () => import('../views/knowledge/Index.vue'),
        meta: { title: 'menu.knowledge', icon: 'Document' }
      },
      {
        path: 'knowledge/add',
        name: 'AddKnowledge',
        component: () => import('../views/knowledge/Edit.vue'),
        meta: { title: 'knowledge.add', activeMenu: '/knowledge' }
      },
      {
        path: 'knowledge/edit/:id',
        name: 'EditKnowledge',
        component: () => import('../views/knowledge/Edit.vue'),
        meta: { title: 'knowledge.edit', activeMenu: '/knowledge' }
      },
      {
        path: 'knowledge/:id',
        name: 'KnowledgeDetail',
        component: () => import('../views/knowledge/Detail.vue'),
        meta: { title: 'knowledge.detail', activeMenu: '/knowledge' }
      },
      {
        path: 'vulndb',
        name: 'VulnDB',
        component: () => import('../views/vulndb/Index.vue'),
        meta: { title: 'menu.vulndb', icon: 'DataAnalysis' }
      },
      {
        path: 'vulndb/id/:id',
        name: 'VulnDBDetail',
        component: () => import('../views/vulndb/Detail.vue'),
        meta: { title: 'vulndb.detail', activeMenu: '/vulndb' }
      },
      {
        path: 'vulndb/add',
        name: 'AddVulnDB',
        component: () => import('../views/vulndb/Add.vue'),
        meta: { title: 'vulndb.add', activeMenu: '/vulndb' }
      },
      {
        path: 'vulndb/edit/:id',
        name: 'EditVulnDB',
        component: () => import('../views/vulndb/Edit.vue'),
        meta: { title: 'vulndb.edit', activeMenu: '/vulndb' }
      },
      {
        path: 'users',
        name: 'Users',
        component: () => import('../views/user/Index.vue'),
        meta: { title: 'menu.users', icon: 'User', roles: ['admin'] }
      },
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('../views/Profile.vue'),
        meta: { title: 'menu.profile', icon: 'UserFilled' }
      },
      {
        path: 'assignments',
        name: 'Assignments',
        component: () => import('../views/assignment/MyAssignments.vue'),
        meta: { title: 'menu.myTasks', icon: 'List' }
      },
      {
        path: 'settings',
        name: 'Settings',
        component: () => import('../views/Settings.vue'),
        meta: { title: 'menu.settings', icon: 'Setting', roles: ['admin'] }
      },
      {
        path: 'scanning',
        name: 'Scanning',
        component: () => import('../views/scan/Index.vue'),
        meta: { title: 'menu.scanning', icon: 'Monitor' }
      },
      {
        path: 'scanning/add',
        name: 'AddScanTask',
        component: () => import('../views/scan/Add.vue'),
        meta: { title: 'scan.createTask', activeMenu: '/scanning' }
      },
      {
        path: 'scanning/:id',
        name: 'ScanTaskDetail',
        component: () => import('../views/scan/Detail.vue'),
        meta: { title: 'scan.taskDetail', activeMenu: '/scanning' }
      },
      {
        path: 'scanning/edit/:id',
        name: 'EditScanTask',
        component: () => import('../views/scan/Edit.vue'),
        meta: { title: 'scan.editTask', activeMenu: '/scanning' }
      },
      {
        path: 'scanning/:id/results',
        name: 'ScanResults',
        component: () => import('../views/scan/Results.vue'),
        meta: { title: 'scan.results', activeMenu: '/scanning' }
      },
      {
        path: 'integrations',
        name: 'Integrations',
        component: () => import('../views/integration/Index.vue'),
        meta: { title: 'CI/CD集成', icon: 'Link', roles: ['admin'] }
      },
      {
        path: 'integrations/add',
        name: 'AddIntegration',
        component: () => import('../views/integration/Add.vue'),
        meta: { title: '添加集成', activeMenu: '/integrations', roles: ['admin'] }
      },
      {
        path: 'integrations/:id',
        name: 'IntegrationDetail',
        component: () => import('../views/integration/Detail.vue'),
        meta: { title: '集成详情', activeMenu: '/integrations', roles: ['admin'] }
      }
    ]
  },
  {
    path: '/404',
    component: () => import('../views/error/404.vue'),
    meta: { title: 'common.notFound', requiresAuth: false }
  },
  {
    path: '/:catchAll(.*)',
    redirect: '/404'
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 进度条配置
NProgress.configure({ showSpinner: false })

// 路由前置守卫
router.beforeEach(async (to, from, next) => {
  // 启动进度条
  NProgress.start()

  // 设置页面标题
  const titleMap = {
    'menu.dashboard': '控制面板',
    'menu.vulnerabilities': '漏洞管理',
    'menu.assets': '资产管理',
    'menu.scanning': '扫描任务',
    'menu.knowledge': '知识库',
    'menu.vulndb': '漏洞库',
    'menu.users': '用户管理',
    'menu.profile': '个人中心',
    'menu.settings': '系统设置',
    'menu.myTasks': '我的任务',
    'scan.createTask': '创建扫描任务',
    'scan.taskDetail': '任务详情',
    'scan.editTask': '编辑任务',
    'scan.results': '扫描结果',
    'vulnerability.add': '添加漏洞',
    'vulnerability.detail': '漏洞详情',
    'asset.add': '添加资产',
    'asset.detail': '资产详情',
    'asset.edit': '编辑资产',
    'knowledge.add': '添加知识',
    'knowledge.edit': '编辑知识',
    'knowledge.detail': '知识详情',
    'vulndb.detail': '漏洞详情',
    'vulndb.add': '添加漏洞',
    'vulndb.edit': '编辑漏洞',
    'login.title': '登录',
    'common.notFound': '页面不存在'
  }
  
  const title = to.meta.title
    ? (titleMap[to.meta.title] || to.meta.title)
    : 'VulnArk'
  document.title = `${title} - VulnArk`

  // 检查是否需要登录
  const hasToken = getToken()
  
  console.log(`路由访问: ${to.path}, 是否需要认证: ${to.meta.requiresAuth !== false}, 是否有token: ${!!hasToken}`)

  if (to.meta.requiresAuth === undefined || to.meta.requiresAuth) {  // 默认需要认证
    if (hasToken) {
      // 如果有token，但是没有用户信息，则获取用户信息
      if (!store.getters.userInfo) {
        try {
          console.log('正在获取用户信息...')
          await store.dispatch('user/getUserInfo')
          console.log('获取用户信息成功')
        } catch (error) {
          // 如果获取用户信息失败，清除token并重定向到登录页
          console.error('获取用户信息失败', error)
          await store.dispatch('user/logout')
          next(`/login?redirect=${to.path}`)
          NProgress.done()
          return
        }
      } else {
        console.log('已有用户信息，无需再次获取')
      }

      // 检查访问权限
      if (to.meta.roles && to.meta.roles.length > 0) {
        const userRole = store.getters.userRole
        console.log('检查访问权限:', {
          path: to.path,
          requiredRoles: to.meta.roles,
          userRole: userRole
        })
        
        // 特殊处理admin角色，确保不区分大小写
        const hasRole = userRole && 
          (to.meta.roles.includes(userRole) || 
           (userRole.toLowerCase() === 'admin' && to.meta.roles.includes('admin')))
        
        if (!hasRole) {
          console.log('无权限访问:', to.path)
          // 无权限访问，重定向到dashboard
          next('/dashboard')
          NProgress.done()
          return
        }
      }

      next()
    } else {
      console.log('未登录，重定向到登录页')
      // 没有token，重定向到登录页
      next(`/login?redirect=${to.path}`)
      NProgress.done()
    }
  } else {
    // 不需要登录的页面
    if (hasToken && (to.path === '/login')) {
      // 已经登录，重定向到首页
      next('/')
    } else {
      next()
    }
  }
})

// 路由后置守卫
router.afterEach(() => {
  // 关闭进度条
  NProgress.done()
})

export default router 