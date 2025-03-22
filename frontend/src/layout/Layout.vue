<template>
  <div class="app-wrapper">
    <!-- 顶部导航 -->
    <div class="header">
      <div class="logo-container">
        <div class="logo-text">VulnArk</div>
        <div class="logo-subtitle">漏洞管理平台</div>
      </div>
      
      <!-- 顶部菜单 -->
      <div class="top-menu">
        <el-menu
          :default-active="activeMenu"
          mode="horizontal"
          background-color="transparent"
          text-color="#e0e0e0"
          active-text-color="#ffffff"
          :router="false"
          class="internationalized-menu"
        >
          <el-menu-item index="/dashboard" @click="navigateTo('/dashboard')" class="menu-item">
            <el-icon><Odometer /></el-icon>
            <span>控制面板</span>
          </el-menu-item>

          <el-menu-item index="/vulnerabilities" @click="navigateTo('/vulnerabilities')" class="menu-item">
            <el-icon><WarningFilled /></el-icon>
            <span>漏洞管理</span>
          </el-menu-item>

          <el-menu-item index="/assets" @click="navigateTo('/assets')" class="menu-item">
            <el-icon><Monitor /></el-icon>
            <span>资产管理</span>
          </el-menu-item>

          <el-menu-item index="/scanning" @click="navigateTo('/scanning')" class="menu-item">
            <el-icon><Search /></el-icon>
            <span>扫描任务</span>
          </el-menu-item>

          <el-menu-item index="/knowledge" @click="navigateTo('/knowledge')" class="menu-item">
            <el-icon><Document /></el-icon>
            <span>知识库</span>
          </el-menu-item>

          <el-menu-item index="/vulndb" @click="navigateTo('/vulndb')" class="menu-item">
            <el-icon><DataAnalysis /></el-icon>
            <span>漏洞库</span>
          </el-menu-item>

          <el-menu-item index="/users" v-if="isAdmin" @click="navigateTo('/users')" class="menu-item">
            <el-icon><User /></el-icon>
            <span>用户管理</span>
          </el-menu-item>

          <el-menu-item index="/integrations" v-if="isAdmin" @click="navigateTo('/integrations')" class="menu-item">
            <el-icon><Link /></el-icon>
            <span>集成管理</span>
          </el-menu-item>

          <el-menu-item index="/settings" v-if="isAdmin" @click="navigateTo('/settings')" class="menu-item">
            <el-icon><Setting /></el-icon>
            <span>系统设置</span>
          </el-menu-item>
        </el-menu>
      </div>
      
      <div class="right">
        <el-dropdown trigger="click" @command="handleCommand" class="user-dropdown">
          <div class="avatar-wrapper">
            <el-avatar :size="36" :src="userAvatar" class="user-avatar">
              {{ userInitials }}
            </el-avatar>
            <span class="username">{{ username }}</span>
            <el-icon class="el-icon--right"><arrow-down /></el-icon>
          </div>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="profile" class="dropdown-item">
                <el-icon><UserFilled /></el-icon>
                <span>个人资料</span>
              </el-dropdown-item>
              <el-dropdown-item v-if="isAdmin" command="settings" class="dropdown-item">
                <el-icon><Setting /></el-icon>
                <span>个人设置</span>
              </el-dropdown-item>
              <el-dropdown-item divided command="logout" class="dropdown-item">
                <el-icon><SwitchButton /></el-icon>
                <span>退出登录</span>
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </div>
    
    <!-- 内容区域 -->
    <div class="main-container">
      <transition name="fade-transform" mode="out-in">
        <router-view :key="$route.path" />
      </transition>
    </div>
  </div>
</template>

<script>
import { computed, inject } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useStore } from 'vuex'
import { ElMessageBox } from 'element-plus'
import {
  Document,
  User,
  UserFilled,
  SwitchButton,
  WarningFilled,
  Odometer,
  Monitor,
  DataAnalysis,
  ArrowDown,
  Setting,
  Search,
  Link
} from '@element-plus/icons-vue'

export default {
  name: 'Layout',
  components: {
    Document,
    User,
    UserFilled,
    SwitchButton,
    WarningFilled,
    Odometer,
    Monitor,
    DataAnalysis,
    ArrowDown,
    Setting,
    Search,
    Link
  },
  setup() {
    const router = useRouter()
    const route = useRoute()
    const store = useStore()
    
    const username = computed(() => {
      const userInfo = store.getters.userInfo
      if (!userInfo) return ''
      
      // 对admin账号进行标准化处理（不区分大小写）
      if (userInfo.username && userInfo.username.toLowerCase() === 'admin') {
        return 'Admin'
      }
      
      return userInfo.username || userInfo.email || ''
    })

    const userAvatar = computed(() => {
      const userInfo = store.getters.userInfo
      return userInfo && userInfo.avatar ? userInfo.avatar : ''
    })

    const userInitials = computed(() => {
      const userInfo = store.getters.userInfo
      if (!userInfo || !userInfo.username) return '?'
      
      const name = userInfo.username.trim()
      if (!name) return '?'
      
      // 获取姓名首字母（如果有空格分隔，取第一个和最后一个名字的首字母）
      const parts = name.split(' ')
      if (parts.length > 1) {
        return (parts[0][0] + parts[parts.length-1][0]).toUpperCase()
      }
      
      // 只有一个名字，取前两个字母或第一个字母
      return (name.length > 1 ? name.substring(0, 2) : name[0]).toUpperCase()
    })

    const isAdmin = computed(() => {
      const role = store.getters.userRole
      return role && (role.toLowerCase() === 'admin')
    })

    const activeMenu = computed(() => {
      return route.meta.activeMenu || route.path
    })

    const navigateTo = (path) => {
      if (route.path !== path) {
        router.push(path)
      }
    }

    const handleCommand = (command) => {
      if (command === 'logout') {
        ElMessageBox.confirm(
          '确定要退出登录吗？',
          '警告',
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }
        ).then(() => {
          store.dispatch('user/logout')
          router.push('/login')
        }).catch(() => {})
      } else if (command === 'profile') {
        router.push('/profile')
      } else if (command === 'settings') {
        router.push('/settings')
      }
    }

    return {
      username,
      userAvatar,
      userInitials,
      isAdmin,
      activeMenu,
      navigateTo,
      handleCommand
    }
  }
}
</script>

<style scoped>
.app-wrapper {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  position: relative;
  background-color: #f5f7fa;
}

.header {
  height: 60px;
  display: flex;
  align-items: center;
  padding: 0 20px;
  background: linear-gradient(90deg, #304156, #3a4d67);
  color: #ffffff;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
  position: sticky;
  top: 0;
  z-index: 100;
  transition: all 0.3s ease;
}

.logo-container {
  display: flex;
  flex-direction: column;
  justify-content: center;
  margin-right: 20px;
  position: relative;
  perspective: 800px;
}

.logo-text {
  font-size: 22px;
  font-weight: bold;
  color: #ffffff;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
  letter-spacing: 1px;
  background: linear-gradient(90deg, #ffffff, #e0e0e0);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  animation: shimmer 3s infinite;
  transform-style: preserve-3d;
}

.logo-subtitle {
  font-size: 12px;
  color: #e0e0e0;
  opacity: 0.8;
  letter-spacing: 0.5px;
}

.top-menu {
  flex: 1;
  overflow-x: auto;
  transition: all 0.3s ease;
}

.top-menu::-webkit-scrollbar {
  display: none;
}

.menu-item {
  transition: all 0.3s;
  position: relative;
  height: 60px;
  padding: 0 15px;
  margin: 0 2px;
  border-radius: 0;
  overflow: hidden;
}

.menu-item:hover {
  background: linear-gradient(to bottom, rgba(255, 255, 255, 0.15), rgba(255, 255, 255, 0.05));
  transform: translateY(-2px);
}

.menu-item::before {
  content: '';
  position: absolute;
  bottom: 0;
  left: 15px;
  right: 15px;
  height: 3px;
  background: linear-gradient(90deg, #4e54c8, #8f94fb);
  transform: scaleX(0);
  transition: transform 0.3s ease;
  border-radius: 3px 3px 0 0;
  opacity: 0;
}

.menu-item:hover::before {
  transform: scaleX(0.8);
  opacity: 0.6;
}

.menu-item.is-active::before {
  transform: scaleX(1);
  opacity: 1;
}

:deep(.el-menu--horizontal) {
  border-bottom: none;
}

:deep(.el-menu--horizontal > .el-menu-item) {
  height: 60px;
  line-height: 60px;
}

:deep(.el-menu--horizontal > .el-menu-item.is-active) {
  border-bottom: none;
  color: #ffffff;
  background: linear-gradient(to bottom, rgba(78, 84, 200, 0.2), rgba(78, 84, 200, 0.1));
}

:deep(.el-menu--horizontal > .el-menu-item:not(.is-disabled):focus, 
       .el-menu--horizontal > .el-menu-item:not(.is-disabled):hover) {
  background: linear-gradient(to bottom, rgba(255, 255, 255, 0.15), rgba(255, 255, 255, 0.05));
}

.right {
  display: flex;
  align-items: center;
}

.user-dropdown {
  cursor: pointer;
  padding: 0 15px;
  height: 60px;
  display: flex;
  align-items: center;
  transition: all 0.3s ease;
  border-radius: 30px;
  margin-left: 10px;
}

.user-dropdown:hover {
  background: linear-gradient(to bottom, rgba(255, 255, 255, 0.15), rgba(255, 255, 255, 0.05));
}

.avatar-wrapper {
  display: flex;
  align-items: center;
  position: relative;
}

.user-avatar {
  margin-right: 8px;
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.2);
  border: 2px solid rgba(255, 255, 255, 0.3);
  transition: all 0.3s cubic-bezier(0.175, 0.885, 0.32, 1.275);
  position: relative;
  z-index: 1;
  background: linear-gradient(135deg, #8f94fb, #4e54c8);
}

.user-dropdown:hover .user-avatar {
  transform: scale(1.1) rotate(5deg);
  border-color: rgba(255, 255, 255, 0.6);
  box-shadow: 0 6px 15px rgba(0, 0, 0, 0.25);
}

.username {
  font-size: 14px;
  margin-right: 8px;
  font-weight: 500;
  letter-spacing: 0.5px;
  transition: all 0.3s ease;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.2);
  position: relative;
}

.user-dropdown:hover .username {
  color: #ffffff;
  transform: translateX(2px);
}

:deep(.el-icon) {
  transition: all 0.3s ease;
}

.user-dropdown:hover :deep(.el-icon) {
  transform: rotate(180deg);
}

.main-container {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
}

:deep(.el-dropdown-menu) {
  border-radius: 10px;
  overflow: hidden;
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.15);
  border: none;
}

:deep(.el-dropdown-menu__item) {
  padding: 10px 20px;
  display: flex;
  align-items: center;
  transition: all 0.3s ease;
}

:deep(.el-dropdown-menu__item:hover) {
  background-color: #f5f7fa;
  transform: translateX(5px);
}

:deep(.el-dropdown-menu__item i), 
:deep(.el-dropdown-menu__item .el-icon) {
  margin-right: 10px;
  color: #4e54c8;
}

:deep(.el-dropdown__popper) {
  margin-top: 5px !important;
}

.dropdown-item {
  display: flex;
  align-items: center;
}

.dropdown-item i, 
.dropdown-item .el-icon {
  margin-right: 10px;
  font-size: 16px;
}

/* 页面过渡效果 */
.fade-transform-enter-active,
.fade-transform-leave-active {
  transition: all 0.5s cubic-bezier(0.645, 0.045, 0.355, 1);
}

.fade-transform-enter-from {
  opacity: 0;
  transform: translateX(-20px);
}

.fade-transform-leave-to {
  opacity: 0;
  transform: translateX(20px);
}

/* 动画 */
@keyframes shimmer {
  0% {
    background-position: -100% 0;
  }
  100% {
    background-position: 200% 0;
  }
}

/* 响应式调整 */
@media (max-width: 768px) {
  .header {
    padding: 0 10px;
  }
  
  .logo-text {
    font-size: 18px;
  }
  
  .logo-subtitle {
    display: none;
  }
  
  .username {
    display: none;
  }
  
  .menu-item {
    padding: 0 10px;
  }
  
  .main-container {
    padding: 15px;
  }
}
</style> 