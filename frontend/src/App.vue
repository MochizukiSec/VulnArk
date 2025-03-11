<template>
  <div id="app">
    <el-config-provider>
      <router-view v-if="!loading" />
      <div v-else class="loading-container">
        <el-empty description="加载中...">
          <el-skeleton style="width: 100%" :rows="8" animated />
        </el-empty>
      </div>
    </el-config-provider>
  </div>
</template>

<script>
import { onMounted, ref, computed } from 'vue'
import { useStore } from 'vuex'
import { useRouter } from 'vue-router'

export default {
  name: 'App',
  setup() {
    const store = useStore()
    const router = useRouter()
    const loading = ref(false)
    
    const isAuthenticated = computed(() => store.getters['auth/isAuthenticated'])
    
    onMounted(async () => {
      try {
        console.log('当前路径:', window.location.pathname);
        
        // 如果当前路径是根路径，不执行其他操作
        if (window.location.pathname === '/') {
          console.log('首页路径，不执行其他操作');
          return;
        }
        
        // 只有非首页才会有加载状态
        loading.value = true;
        
        // 尝试从本地存储恢复用户会话
        await store.dispatch('auth/restoreSession')
        
        // 其他情况，如果未认证且需要认证，则重定向到登录页面
        if (!isAuthenticated.value && 
            router.currentRoute.value.name !== 'Login' && 
            router.currentRoute.value.meta.requiresAuth) {
          console.log('未认证用户访问需要认证的页面，重定向到登录页面')
          router.push('/login')
        }
      } catch (error) {
        console.error('初始化错误:', error)
      } finally {
        // 无论发生什么，都完成加载
        loading.value = false
      }
    })
    
    return {
      loading
    }
  }
}
</script>

<style lang="scss">
/* 重置样式 */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html, body {
  height: auto;
  min-height: 100%;
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  font-size: 16px;
  color: #333;
  background: linear-gradient(135deg, #f8f9fa 0%, #edf1f7 100%);
  overflow-x: hidden;
  overflow-y: auto;
  position: relative;
  
  &::before {
    content: '';
    position: fixed;
    top: 0;
    right: 0;
    bottom: 0;
    left: 0;
    background: radial-gradient(circle at 5% 95%, rgba(56, 128, 255, 0.05) 0%, transparent 50%),
                radial-gradient(circle at 95% 5%, rgba(255, 56, 155, 0.05) 0%, transparent 50%);
    pointer-events: none;
    z-index: -1;
  }
}

#app {
  min-height: 100vh;
  height: auto;
  display: flex;
  flex-direction: column;
  overflow-y: auto;
}

.loading-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(10px);
}

/* 页面过渡动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s, transform 0.3s;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(10px);
}

/* 现代化主题样式覆盖 */
.el-card {
  border-radius: 16px;
  border: none;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.05);
  margin-bottom: 20px;
  transition: all 0.3s;
  overflow: hidden;
  
  &:hover {
    box-shadow: 0 15px 35px rgba(0, 0, 0, 0.1);
    transform: translateY(-3px);
  }
  
  .el-card__header {
    padding: 20px;
    border-bottom: 1px solid rgba(0, 0, 0, 0.03);
  }
  
  .el-card__body {
    padding: 20px;
  }
}

.el-button {
  border-radius: 10px;
  font-weight: 600;
  transition: all 0.3s;
  padding: 12px 24px;
  
  &:hover {
    transform: translateY(-2px);
  }
  
  &.el-button--primary {
    background: linear-gradient(90deg, #36d1dc, #5b86e5);
    border: none;
    box-shadow: 0 5px 15px rgba(91, 134, 229, 0.2);
    
    &:hover {
      box-shadow: 0 8px 20px rgba(91, 134, 229, 0.3);
    }
  }
  
  &.el-button--success {
    background: linear-gradient(90deg, #67C23A, #42b883);
    border: none;
    box-shadow: 0 5px 15px rgba(66, 184, 131, 0.2);
    
    &:hover {
      box-shadow: 0 8px 20px rgba(66, 184, 131, 0.3);
    }
  }
  
  &.el-button--warning {
    background: linear-gradient(90deg, #E6A23C, #ff9d47);
    border: none;
    box-shadow: 0 5px 15px rgba(255, 157, 71, 0.2);
    
    &:hover {
      box-shadow: 0 8px 20px rgba(255, 157, 71, 0.3);
    }
  }
  
  &.el-button--danger {
    background: linear-gradient(90deg, #F56C6C, #ff3366);
    border: none;
    box-shadow: 0 5px 15px rgba(255, 51, 102, 0.2);
    
    &:hover {
      box-shadow: 0 8px 20px rgba(255, 51, 102, 0.3);
    }
  }
}

.el-table {
  border-radius: 12px;
  box-shadow: 0 5px 20px rgba(0, 0, 0, 0.03);
  overflow: hidden;
  
  th {
    background-color: rgba(0, 0, 0, 0.02) !important;
    font-weight: 600;
    padding: 12px 0;
    border-bottom: 1px solid rgba(0, 0, 0, 0.03);
  }
  
  td {
    padding: 12px 0;
    border-bottom: 1px solid rgba(0, 0, 0, 0.03);
  }
  
  tr:hover > td {
    background-color: rgba(0, 0, 0, 0.02) !important;
  }
}

.el-dialog {
  border-radius: 20px;
  overflow: hidden;
  box-shadow: 0 15px 40px rgba(0, 0, 0, 0.1);
  
  .el-dialog__header {
    padding: 20px;
    background: rgba(0, 0, 0, 0.02);
    border-bottom: 1px solid rgba(0, 0, 0, 0.03);
  }
  
  .el-dialog__title {
    font-weight: 700;
    font-size: 18px;
    color: #303133;
  }
  
  .el-dialog__body {
    padding: 30px 20px;
  }
  
  .el-dialog__footer {
    padding: 15px 20px;
    border-top: 1px solid rgba(0, 0, 0, 0.03);
  }
}

.el-tag {
  border-radius: 100px;
  padding: 0 12px;
  height: 28px;
  line-height: 26px;
  font-weight: 600;
  font-size: 12px;
  border: none;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.05);
}

.el-input__inner {
  border-radius: 10px;
  padding: 12px 15px;
  height: 48px;
  transition: all 0.3s;
  
  &:focus {
    box-shadow: 0 0 15px rgba(64, 158, 255, 0.1);
  }
}

.el-select .el-input__inner {
  padding-right: 30px;
}

.el-textarea__inner {
  border-radius: 10px;
  padding: 12px 15px;
}

.el-pagination {
  margin-top: 20px;
  padding: 10px;
  border-radius: 10px;
  background: rgba(255, 255, 255, 0.7);
  justify-content: center;
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.03);
  backdrop-filter: blur(5px);
}

/* 自定义元素样式 */
.page-header {
  margin-bottom: 30px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  position: relative;
}

.page-title {
  font-size: 28px;
  font-weight: 700;
  color: #303133;
  letter-spacing: -0.5px;
  margin-bottom: 8px;
  background: linear-gradient(90deg, #1a1d3f, #303F9F);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  position: relative;
  
  &::after {
    content: '';
    position: absolute;
    bottom: -8px;
    left: 0;
    width: 50px;
    height: 3px;
    background: linear-gradient(90deg, #36d1dc, #5b86e5);
    border-radius: 3px;
  }
}

.page-subtitle {
  font-size: 16px;
  color: #606266;
  margin-top: 10px;
  font-weight: 400;
  max-width: 600px;
  line-height: 1.5;
}

.section-title {
  font-size: 20px;
  font-weight: 700;
  margin-bottom: 20px;
  color: #303133;
  position: relative;
  display: inline-block;
  
  &::after {
    content: '';
    position: absolute;
    bottom: -6px;
    left: 0;
    width: 30px;
    height: 2px;
    background: linear-gradient(90deg, #36d1dc, #5b86e5);
    border-radius: 2px;
  }
}

/* 状态标签样式 */
.status-label {
  text-transform: capitalize;
  font-weight: 600;
  padding: 5px 12px;
  border-radius: 100px;
  font-size: 12px;
  display: inline-block;
  box-shadow: 0 3px 8px rgba(0, 0, 0, 0.1);
}

/* 严重性标签样式 */
.severity-label {
  text-transform: capitalize;
  font-weight: 600;
  padding: 5px 12px;
  border-radius: 100px;
  font-size: 12px;
  display: inline-block;
  box-shadow: 0 3px 8px rgba(0, 0, 0, 0.1);
}

.critical-bg {
  background: linear-gradient(90deg, #F56C6C, #ff3366);
  color: white;
}

.high-bg {
  background: linear-gradient(90deg, #E6A23C, #ff9d47);
  color: white;
}

.medium-bg {
  background: linear-gradient(90deg, #F2D864, #f9d423);
  color: #333;
}

.low-bg {
  background: linear-gradient(90deg, #67C23A, #42b883);
  color: white;
}

.info-bg {
  background: linear-gradient(90deg, #909399, #666a73);
  color: white;
}

/* 响应式调整 */
@media (max-width: 767px) {
  .page-title {
    font-size: 24px;
  }
  
  .page-subtitle {
    font-size: 14px;
  }
  
  .section-title {
    font-size: 18px;
  }
  
  .el-button {
    padding: 10px 20px;
  }
  
  .el-input__inner {
    height: 44px;
  }
}
</style> 