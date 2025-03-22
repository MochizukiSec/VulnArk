<template>
  <div class="scan-status">
    <el-tag :type="statusType" class="status-tag">
      {{ statusText }}
    </el-tag>
    
    <el-progress 
      v-if="showProgress" 
      :percentage="progress || 0" 
      :status="progressStatus"
      :stroke-width="8"
      class="status-progress"
    />
  </div>
</template>

<script>
export default {
  name: 'ScanStatus',
  props: {
    status: {
      type: String,
      required: true
    },
    progress: {
      type: Number,
      default: 0
    }
  },
  computed: {
    statusType() {
      switch (this.status) {
        case 'ready': return 'info'
        case 'running': return 'primary'
        case 'completed': return 'success'
        case 'failed': return 'danger'
        case 'cancelled': return 'warning'
        case 'error': return 'danger'
        default: return 'info'
      }
    },
    statusText() {
      switch (this.status) {
        case 'ready': return '就绪'
        case 'running': return '运行中'
        case 'completed': return '已完成'
        case 'failed': return '失败'
        case 'cancelled': return '已取消'
        case 'error': return '错误'
        default: return this.status
      }
    },
    showProgress() {
      return this.status === 'running' && this.progress !== undefined
    },
    progressStatus() {
      if (this.progress >= 100) return 'success'
      return ''
    }
  }
}
</script>

<style scoped>
.scan-status {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 8px;
}

.status-tag {
  font-size: 13px;
}

.status-progress {
  width: 100%;
  margin-top: 5px;
}
</style> 