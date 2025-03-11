<template>
  <div class="app-container">
    <!-- 顶部导航栏 -->
    <header class="app-header">
      <div class="header-left">
        <h1 class="logo">
          <router-link to="/dashboard">漏洞管理平台</router-link>
        </h1>
        
        <!-- 主导航菜单 -->
        <nav class="main-nav">
          <router-link to="/dashboard" class="nav-item" active-class="nav-active">
            <i class="el-icon-s-home"></i> 仪表盘
          </router-link>
          
          <el-dropdown trigger="hover" class="nav-item-dropdown">
            <div class="nav-item">
              <i class="el-icon-warning"></i> 漏洞管理
              <i class="el-icon-arrow-down nav-arrow"></i>
            </div>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item>
                  <router-link to="/vulnerabilities" class="dropdown-link">漏洞列表</router-link>
                </el-dropdown-item>
                <el-dropdown-item>
                  <router-link to="/vulnerabilities/create" class="dropdown-link">创建漏洞</router-link>
                </el-dropdown-item>
                <el-dropdown-item>
                  <router-link to="/vulnerabilities/import" class="dropdown-link">导入漏洞</router-link>
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
          
          <router-link to="/reports" class="nav-item" active-class="nav-active">
            <i class="el-icon-document"></i> 报告中心
          </router-link>
          
          <router-link to="/admin/users" class="nav-item" active-class="nav-active">
            <i class="el-icon-user"></i> 用户管理
          </router-link>
        </nav>
      </div>
      
      <div class="header-right">
        <!-- 通知 -->
        <el-popover
          placement="bottom"
          trigger="click"
          width="320"
        >
          <template #reference>
            <el-badge :value="unreadNotifications.length" :hidden="!unreadNotifications.length" class="notification-badge">
              <el-button type="text" icon="el-icon-bell" class="notification-btn"></el-button>
            </el-badge>
          </template>
          <div class="notification-menu">
            <div class="notification-header">
              <span class="notification-title">通知</span>
              <el-button 
                v-if="notifications.length" 
                type="text" 
                size="small" 
                @click="clearAllNotifications"
              >
                清除全部
              </el-button>
            </div>
            <div class="notification-list" v-if="notifications.length">
              <div
                v-for="notification in notifications"
                :key="notification.id"
                class="notification-item"
                :class="{ 'is-read': notification.read }"
              >
                <div class="notification-icon" :class="`icon-${notification.type}`">
                  <i :class="getNotificationIcon(notification.type)"></i>
                </div>
                <div class="notification-content">
                  <h4 class="notification-content-title">{{ notification.title }}</h4>
                  <p class="notification-content-msg">{{ notification.message }}</p>
                  <span class="notification-time">{{ formatTime(notification.timestamp) }}</span>
                </div>
                <el-button 
                  type="text" 
                  size="small" 
                  icon="el-icon-close"
                  @click="markNotificationRead(notification.id)"
                  class="notification-close"
                ></el-button>
              </div>
            </div>
            <div v-else class="notification-empty">
              <i class="el-icon-bell"></i>
              <p>暂无通知</p>
            </div>
          </div>
        </el-popover>
        
        <!-- 用户菜单 -->
        <el-dropdown trigger="click" @command="handleCommand">
          <div class="user-profile">
            <el-avatar :size="32" :src="userAvatar">{{ userInitials }}</el-avatar>
            <span class="username">{{ user.firstName }}</span>
            <i class="el-icon-arrow-down"></i>
          </div>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="profile">
                <i class="el-icon-user"></i> 个人资料
              </el-dropdown-item>
              <el-dropdown-item divided command="logout">
                <i class="el-icon-switch-button"></i> 退出登录
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </header>
    
    <!-- 主要内容区域 -->
    <div class="main-content">
      <!-- 页面内容 -->
      <main class="page-content">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </main>
    </div>
  </div>
</template>

<script>
import { computed } from 'vue'
import { useStore } from 'vuex'
import { useRouter, useRoute } from 'vue-router'
import moment from 'moment'

export default {
  name: 'AppLayout',
  
  setup() {
    const store = useStore()
    const router = useRouter()
    const route = useRoute()
    
    // 当前用户
    const user = computed(() => store.getters['auth/currentUser'] || {})
    
    // 判断是否为管理员
    const isAdmin = computed(() => store.getters['auth/isAdmin'])
    
    // 获取当前活动菜单项
    const activeMenu = computed(() => {
      // 根据当前路由路径返回对应的菜单项
      return route.path
    })
    
    // 获取当前页面标题
    const currentPageTitle = computed(() => {
      return route.meta.title || '漏洞管理平台'
    })
    
    // 用户头像和缩写
    const userAvatar = computed(() => user.value.profilePicture || '')
    const userInitials = computed(() => {
      if (user.value.firstName && user.value.lastName) {
        return `${user.value.firstName.charAt(0)}${user.value.lastName.charAt(0)}`
      }
      return user.value.username ? user.value.username.charAt(0).toUpperCase() : '?'
    })
    
    // 获取通知
    const notifications = computed(() => store.getters.notifications || [])
    const unreadNotifications = computed(() => notifications.value.filter(n => !n.read))
    
    // 处理用户菜单指令
    const handleCommand = (command) => {
      switch (command) {
        case 'profile':
          router.push('/profile')
          break
        case 'logout':
          store.dispatch('auth/logout')
          router.push('/login')
          break
      }
    }
    
    // 标记通知为已读
    const markNotificationRead = (id) => {
      store.dispatch('markNotificationRead', id)
    }
    
    // 清除所有通知
    const clearAllNotifications = () => {
      store.dispatch('clearNotifications')
    }
    
    // 格式化时间
    const formatTime = (timestamp) => {
      return moment(timestamp).fromNow()
    }
    
    // 获取通知图标
    const getNotificationIcon = (type) => {
      switch (type) {
        case 'success':
          return 'el-icon-success'
        case 'warning':
          return 'el-icon-warning'
        case 'error':
          return 'el-icon-error'
        case 'info':
        default:
          return 'el-icon-info'
      }
    }
    
    return {
      activeMenu,
      currentPageTitle,
      user,
      isAdmin,
      userAvatar,
      userInitials,
      notifications,
      unreadNotifications,
      handleCommand,
      markNotificationRead,
      clearAllNotifications,
      formatTime,
      getNotificationIcon
    }
  }
}
</script>

<style lang="scss" scoped>
.app-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  width: 100%;
}

/* 顶部导航栏样式 */
.app-header {
  height: 70px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  position: sticky;
  top: 0;
  z-index: 100;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.03);
  transition: all 0.3s;
}

.header-left {
  display: flex;
  align-items: center;
}

.logo {
  font-size: 22px;
  font-weight: 700;
  margin: 0 35px 0 0;
  white-space: nowrap;
  
  a {
    background: linear-gradient(90deg, #36d1dc, #5b86e5);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    text-shadow: 0 0 10px rgba(91, 134, 229, 0.3);
    text-decoration: none;
    transition: all 0.3s;
    
    &:hover {
      text-shadow: 0 0 15px rgba(91, 134, 229, 0.5);
      transform: scale(1.02);
    }
  }
}

/* 主导航菜单 */
.main-nav {
  display: flex;
  align-items: center;
  gap: 5px;
}

.nav-item {
  padding: 8px 16px;
  border-radius: 10px;
  color: #606266;
  text-decoration: none;
  font-weight: 500;
  display: flex;
  align-items: center;
  gap: 6px;
  transition: all 0.3s;
  position: relative;
  cursor: pointer;
  font-size: 15px;
  
  &:hover, &.nav-active {
    color: #5b86e5;
    background: rgba(91, 134, 229, 0.08);
  }
  
  i {
    font-size: 18px;
  }
}

.nav-item-dropdown {
  .nav-arrow {
    font-size: 12px;
    margin-left: 3px;
    color: #999;
    transition: all 0.3s;
  }
  
  &:hover .nav-arrow {
    transform: rotate(180deg);
    color: #5b86e5;
  }
}

.dropdown-link {
  text-decoration: none;
  color: #606266;
  display: block;
  width: 100%;
  transition: all 0.3s;
  
  &:hover {
    color: #5b86e5;
  }
}

.header-right {
  display: flex;
  align-items: center;
  gap: 20px;
}

.user-profile {
  display: flex;
  align-items: center;
  cursor: pointer;
  border-radius: 40px;
  padding: 4px 10px 4px 4px;
  transition: all 0.3s;
  background: rgba(0, 0, 0, 0.02);
  
  &:hover {
    background: rgba(0, 0, 0, 0.05);
    transform: translateY(-1px);
  }
  
  .username {
    margin: 0 8px;
    font-size: 14px;
    font-weight: 500;
    color: #333;
  }
  
  .el-avatar {
    border: 2px solid #fff;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  }
}

/* 主内容区域样式 */
.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  position: relative;
}

/* 页面内容样式 */
.page-content {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
  background-color: #f8f9fa;
  
  &::-webkit-scrollbar {
    width: 8px;
  }
  
  &::-webkit-scrollbar-thumb {
    background: rgba(0, 0, 0, 0.1);
    border-radius: 4px;
  }
}

/* 通知样式 */
.notification-badge {
  margin-right: 10px;
}

.notification-btn {
  font-size: 20px;
  color: #606266;
  transition: all 0.3s;
  
  &:hover {
    color: #409EFF;
    transform: scale(1.1);
  }
}

.notification-menu {
  max-height: 400px;
  overflow-y: auto;
  border-radius: 12px;
  
  &::-webkit-scrollbar {
    width: 4px;
  }
  
  &::-webkit-scrollbar-thumb {
    background: rgba(0, 0, 0, 0.1);
    border-radius: 2px;
  }
}

.notification-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
}

.notification-title {
  font-weight: 600;
  font-size: 16px;
  color: #333;
}

.notification-list {
  padding: 0;
}

.notification-item {
  padding: 15px;
  display: flex;
  align-items: flex-start;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
  position: relative;
  transition: all 0.2s;
  
  &.is-read {
    opacity: 0.6;
  }
  
  &:hover {
    background-color: rgba(0, 0, 0, 0.02);
    
    .notification-close {
      opacity: 1;
    }
  }
}

.notification-icon {
  width: 36px;
  height: 36px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 12px;
  transition: all 0.3s;
  
  &.icon-success {
    background-color: rgba(103, 194, 58, 0.15);
    color: #67c23a;
  }
  
  &.icon-warning {
    background-color: rgba(230, 162, 60, 0.15);
    color: #e6a23c;
  }
  
  &.icon-error {
    background-color: rgba(245, 108, 108, 0.15);
    color: #f56c6c;
  }
  
  &.icon-info {
    background-color: rgba(144, 147, 153, 0.15);
    color: #909399;
  }
  
  i {
    font-size: 18px;
  }
}

.notification-content {
  flex: 1;
}

.notification-content-title {
  margin: 0 0 6px;
  font-size: 14px;
  font-weight: 600;
  color: #333;
}

.notification-content-msg {
  margin: 0 0 6px;
  font-size: 13px;
  color: #666;
  line-height: 1.5;
}

.notification-time {
  font-size: 12px;
  color: #999;
  font-weight: 500;
}

.notification-close {
  position: absolute;
  top: 12px;
  right: 12px;
  opacity: 0;
  transition: opacity 0.2s;
  background-color: rgba(0, 0, 0, 0.05);
  border-radius: 50%;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  
  &:hover {
    background-color: rgba(0, 0, 0, 0.1);
  }
}

.notification-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 0;
  color: #aaa;
  
  i {
    font-size: 32px;
    margin-bottom: 15px;
    opacity: 0.5;
  }
  
  p {
    margin: 0;
    font-size: 14px;
  }
}

/* 响应式适配 */
@media (max-width: 991px) {
  .app-header {
    height: auto;
    padding: 10px 15px;
    flex-wrap: wrap;
  }
  
  .header-left {
    width: 100%;
    justify-content: space-between;
    margin-bottom: 10px;
  }
  
  .main-nav {
    order: 3;
    width: 100%;
    margin-top: 10px;
    overflow-x: auto;
    padding-bottom: 5px;
    
    &::-webkit-scrollbar {
      height: 3px;
    }
    
    &::-webkit-scrollbar-thumb {
      background: rgba(0, 0, 0, 0.1);
      border-radius: 2px;
    }
  }
  
  .header-right {
    order: 2;
  }
  
  .logo {
    margin-right: 0;
  }
}

@media (max-width: 767px) {
  .app-header {
    padding: 10px;
  }
  
  .nav-item {
    padding: 8px 12px;
    font-size: 14px;
  }
  
  .logo {
    font-size: 20px;
  }
  
  .user-profile .username {
    display: none;
  }
}

/* 过渡动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s, transform 0.3s;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(10px);
}
</style> 
