<template>
  <div id="app">
    <router-view v-slot="{ Component }">
      <transition name="fade" mode="out-in">
        <component :is="Component" />
      </transition>
    </router-view>
  </div>
</template>

<script>
export default {
  name: 'App'
}
</script>

<style>
html, body, #app {
  margin: 0;
  padding: 0;
  height: 100%;
  font-family: 'Helvetica Neue', Helvetica, 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

.el-loading-mask {
  z-index: 9999 !important;
}

#app {
  font-family: 'Helvetica Neue', Helvetica, 'PingFang SC', 'Microsoft YaHei', Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
  min-height: 100vh;
  background-color: #f5f7fa;
}

/* 全局过渡动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* 滑动过渡效果 */
.slide-fade-enter-active,
.slide-fade-leave-active {
  transition: all 0.3s ease;
}
.slide-fade-enter-from,
.slide-fade-leave-to {
  transform: translateX(10px);
  opacity: 0;
}

/* 缩放过渡效果 */
.scale-enter-active,
.scale-leave-active {
  transition: all 0.3s ease;
}
.scale-enter-from,
.scale-leave-to {
  transform: scale(0.98);
  opacity: 0;
}

/* 改进全局按钮样式 */
.el-button {
  transition: all 0.2s ease;
}
.el-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}
.el-button--primary:hover {
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.25);
}
.el-button--success:hover {
  box-shadow: 0 4px 12px rgba(103, 194, 58, 0.25);
}
.el-button--danger:hover {
  box-shadow: 0 4px 12px rgba(245, 108, 108, 0.25);
}
.el-button--warning:hover {
  box-shadow: 0 4px 12px rgba(230, 162, 60, 0.25);
}

/* 卡片样式优化 */
.el-card {
  transition: all 0.3s ease;
  border-radius: 8px;
  border: none;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
}
.el-card:hover {
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.08);
}
</style> 