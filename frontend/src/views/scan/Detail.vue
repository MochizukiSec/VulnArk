<template>
  <div class="scan-detail-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <h2>{{ $t('scan.taskDetail') }}</h2>
          <div class="header-actions">
            <el-button @click="navigateToResults" v-if="scanTask.status === 'completed'">
              {{ $t('scan.results') }}
            </el-button>
            <el-button @click="goBack">{{ $t('common.back') }}</el-button>
          </div>
        </div>
      </template>

      <div v-if="!loading && scanTask">
        <!-- 状态信息和快捷操作 -->
        <div class="status-action-bar">
          <div class="status-info">
            <el-tag :type="getStatusType(scanTask.status)">
              {{ $t(`scan.status.${scanTask.status}`) }}
            </el-tag>
            <span class="scan-id">ID: {{ scanTask.id }}</span>
          </div>
          
          <div class="quick-actions">
            <el-button-group>
              <el-button
                v-if="canStart"
                size="small"
                type="success"
                @click="startTask"
                :loading="actionLoading"
                :title="$t('scan.runNow')"
              >
                <el-icon><VideoPlay /></el-icon>
                {{ $t('scan.runNow') }}
              </el-button>
              
              <el-button
                v-if="canCancel"
                size="small"
                type="warning"
                @click="cancelTask"
                :loading="actionLoading"
                :title="$t('scan.cancel')"
              >
                <el-icon><VideoPause /></el-icon>
                {{ $t('scan.cancel') }}
              </el-button>
              
              <el-button
                v-if="canEdit"
                size="small"
                type="primary"
                @click="navigateToEdit"
                :title="$t('common.edit')"
              >
                <el-icon><Edit /></el-icon>
                {{ $t('common.edit') }}
              </el-button>
              
              <el-button
                v-if="canDelete"
                size="small"
                type="danger"
                @click="confirmDelete"
                :loading="actionLoading"
                :title="$t('common.delete')"
              >
                <el-icon><Delete /></el-icon>
                {{ $t('common.delete') }}
              </el-button>
            </el-button-group>
          </div>
        </div>

        <!-- 进度显示 - 当任务正在运行时显示 -->
        <div v-if="scanTask.status === 'running'" class="progress-section">
          <div class="progress-header">
            <h3>{{ $t('scan.scanProgress') }}</h3>
            <span class="progress-percentage">{{ progressPercentage }}%</span>
          </div>
          
          <el-progress 
            :percentage="progressPercentage" 
            :status="progressStatus"
            :stroke-width="20"
            :format="() => `${progressDetails}`"
          />
          
          <div class="realtime-updates" v-if="progressLogs.length > 0">
            <div class="log-header">
              <h4>{{ $t('scan.realtimeUpdates') }}</h4>
              <el-switch v-model="autoScroll" :active-text="$t('scan.autoScroll')" />
            </div>
            
            <div class="log-container" ref="logContainer">
              <div v-for="(log, index) in progressLogs" :key="index" class="log-entry">
                <span class="log-time">{{ formatTime(log.timestamp) }}</span>
                <span class="log-message" :class="'log-' + log.level">{{ log.message }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- 任务详情信息 -->
        <el-descriptions :column="2" border>
          <el-descriptions-item :label="$t('common.name')">
            {{ scanTask.name }}
          </el-descriptions-item>
          
          <el-descriptions-item :label="$t('scan.scannerType.label')">
            {{ $t(`scan.scannerType.${scanTask.type}`) }}
          </el-descriptions-item>
          
          <el-descriptions-item :label="$t('common.description')" :span="2">
            {{ scanTask.description || $t('common.noData') }}
          </el-descriptions-item>
          
          <el-descriptions-item :label="$t('scan.scannerUrl')">
            {{ scanTask.scanner_url || $t('common.noData') }}
          </el-descriptions-item>
          
          <el-descriptions-item :label="$t('scan.scannerUsername')">
            {{ scanTask.scanner_username || $t('common.noData') }}
          </el-descriptions-item>
        </el-descriptions>

        <!-- 目标信息 -->
        <div class="section-title">
          <h3>{{ $t('scan.targetConfig') }}</h3>
        </div>
        
        <el-descriptions :column="1" border>
          <el-descriptions-item :label="$t('scan.targetIPs')" v-if="scanTask.target_ips">
            <pre>{{ scanTask.target_ips }}</pre>
          </el-descriptions-item>
          
          <el-descriptions-item :label="$t('scan.targetURLs')" v-if="scanTask.target_urls">
            <pre>{{ scanTask.target_urls }}</pre>
          </el-descriptions-item>
          
          <el-descriptions-item :label="$t('scan.targetAssets')" v-if="scanTask.target_assets">
            <div class="target-assets">
              <el-tag 
                v-for="(assetId, index) in scanTask.target_assets.split(',')" 
                :key="index"
                class="asset-tag"
              >
                {{ getAssetName(assetId) || assetId }}
              </el-tag>
            </div>
          </el-descriptions-item>
          
          <el-descriptions-item :label="$t('scan.scanParameters')" v-if="scanTask.scan_parameters">
            <pre>{{ scanTask.scan_parameters }}</pre>
          </el-descriptions-item>
        </el-descriptions>

        <!-- 计划信息 -->
        <div class="section-title">
          <h3>{{ $t('scan.scheduleConfig') }}</h3>
        </div>
        
        <el-descriptions :column="2" border>
          <el-descriptions-item :label="$t('scan.scheduleAt')">
            {{ scanTask.scheduled_at ? formatDate(scanTask.scheduled_at) : $t('common.noData') }}
          </el-descriptions-item>
          
          <el-descriptions-item :label="$t('scan.isRecurring')">
            <el-tag :type="scanTask.is_recurring ? 'success' : 'info'">
              {{ scanTask.is_recurring ? $t('common.yes') : $t('common.no') }}
            </el-tag>
          </el-descriptions-item>
          
          <el-descriptions-item :label="$t('scan.cronSchedule')" v-if="scanTask.is_recurring">
            {{ scanTask.cron_schedule || $t('common.noData') }}
          </el-descriptions-item>
        </el-descriptions>

        <!-- 结果信息 -->
        <div class="section-title">
          <h3>{{ $t('scan.resultSummary') }}</h3>
        </div>
        
        <el-descriptions :column="2" border>
          <el-descriptions-item :label="$t('scan.startedAt')">
            {{ scanTask.started_at ? formatDate(scanTask.started_at) : $t('common.noData') }}
          </el-descriptions-item>
          
          <el-descriptions-item :label="$t('scan.completedAt')">
            {{ scanTask.completed_at ? formatDate(scanTask.completed_at) : $t('common.noData') }}
          </el-descriptions-item>
          
          <el-descriptions-item :label="$t('scan.duration')" v-if="scanTask.started_at && scanTask.completed_at">
            {{ calculateDuration(scanTask.started_at, scanTask.completed_at) }}
          </el-descriptions-item>
          
          <el-descriptions-item :label="$t('scan.vulnerabilityCount')" v-if="scanTask.status === 'completed'">
            <span class="vuln-count">{{ scanTask.total_vulnerabilities }}</span>
          </el-descriptions-item>
          
          <el-descriptions-item :label="$t('scan.criticalCount')" v-if="scanTask.status === 'completed' && scanTask.critical_vulnerabilities > 0">
            <el-tag type="danger">{{ scanTask.critical_vulnerabilities }}</el-tag>
          </el-descriptions-item>
          
          <el-descriptions-item :label="$t('scan.highCount')" v-if="scanTask.status === 'completed' && scanTask.high_vulnerabilities > 0">
            <el-tag type="error">{{ scanTask.high_vulnerabilities }}</el-tag>
          </el-descriptions-item>
          
          <el-descriptions-item :label="$t('scan.mediumCount')" v-if="scanTask.status === 'completed' && scanTask.medium_vulnerabilities > 0">
            <el-tag type="warning">{{ scanTask.medium_vulnerabilities }}</el-tag>
          </el-descriptions-item>
          
          <el-descriptions-item :label="$t('scan.lowCount')" v-if="scanTask.status === 'completed' && scanTask.low_vulnerabilities > 0">
            <el-tag type="info">{{ scanTask.low_vulnerabilities }}</el-tag>
          </el-descriptions-item>
          
          <el-descriptions-item :label="$t('scan.resultSummary')" :span="2" v-if="scanTask.result_summary">
            {{ scanTask.result_summary }}
          </el-descriptions-item>
        </el-descriptions>
      </div>
      
      <div v-else-if="!loading && !scanTask" class="not-found">
        <el-empty :description="$t('common.notFound')" />
      </div>
    </el-card>
  </div>
</template>

<script>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { VideoPlay, VideoPause, Edit, Delete } from '@element-plus/icons-vue'
import api from '@/api'
import wsService from '@/utils/websocket'

export default {
  name: 'ScanTaskDetail',
  components: {
    VideoPlay,
    VideoPause,
    Edit,
    Delete
  },
  setup() {
    const { t } = useI18n()
    const route = useRoute()
    const router = useRouter()
    
    const scanTask = ref(null)
    const loading = ref(true)
    const actionLoading = ref(false)
    const assets = ref([])
    const logContainer = ref(null)
    
    // 进度相关状态
    const progress = ref(0)
    const progressStage = ref('')
    const progressDetails = ref('')
    const progressLogs = ref([])
    const autoScroll = ref(true)
    
    // 获取ID参数
    const taskId = route.params.id
    
    // 计算属性：状态相关按钮显示
    const canStart = computed(() => {
      return scanTask.value && (
        scanTask.value.status === 'created' || 
        scanTask.value.status === 'failed' || 
        scanTask.value.status === 'cancelled' || 
        scanTask.value.status === 'completed'
      )
    })
    
    // 计算进度状态
    const progressStatus = computed(() => {
      if (scanTask.value?.status === 'failed') return 'exception'
      if (progress.value >= 100) return 'success'
      return ''
    })
    
    // 计算百分比，确保在0-100范围内
    const progressPercentage = computed(() => {
      return Math.min(100, Math.max(0, progress.value))
    })
    
    const canCancel = computed(() => {
      return scanTask.value && (
        scanTask.value.status === 'queued' || 
        scanTask.value.status === 'running'
      )
    })
    
    const canEdit = computed(() => {
      return scanTask.value && (
        scanTask.value.status !== 'queued' && 
        scanTask.value.status !== 'running'
      )
    })
    
    const canDelete = computed(() => {
      return scanTask.value && (
        scanTask.value.status !== 'queued' && 
        scanTask.value.status !== 'running'
      )
    })
    
    // 加载扫描任务详情
    const loadTask = async () => {
      loading.value = true
      try {
        const res = await api.get(`/scans/${taskId}`)
        scanTask.value = res.data
      } catch (error) {
        console.error('加载扫描任务详情失败:', error)
        ElMessage.error(t('common.error') + ': ' + (error.response?.data?.message || error.message))
      } finally {
        loading.value = false
      }
    }
    
    // 加载资产列表
    const loadAssets = async () => {
      try {
        const res = await api.get('/assets')
        assets.value = res.data.assets || []
      } catch (error) {
        console.error('加载资产列表失败:', error)
      }
    }
    
    // 获取资产名称
    const getAssetName = (assetId) => {
      const asset = assets.value.find(a => a.id.toString() === assetId.toString())
      return asset ? asset.name : null
    }
    
    // 操作函数
    const startTask = async () => {
      try {
        await api.post(`/scans/${taskId}/start`)
        ElMessage.success(t('scan.runSuccess'))
        loadTask()
      } catch (error) {
        console.error('启动任务失败:', error)
        ElMessage.error(t('common.error') + ': ' + (error.response?.data?.message || error.message))
      }
    }
    
    const cancelTask = async () => {
      try {
        await ElMessageBox.confirm(
          t('scan.confirmCancel'),
          t('common.warning'),
          {
            confirmButtonText: t('common.confirm'),
            cancelButtonText: t('common.cancel'),
            type: 'warning'
          }
        )
        
        await api.post(`/scans/${taskId}/cancel`)
        ElMessage.success(t('scan.cancelSuccess'))
        loadTask()
      } catch (error) {
        if (error !== 'cancel') {
          console.error('取消任务失败:', error)
          ElMessage.error(t('common.error') + ': ' + (error.response?.data?.message || error.message))
        }
      }
    }
    
    const confirmDelete = async () => {
      try {
        await ElMessageBox.confirm(
          t('scan.confirmDelete'),
          t('common.warning'),
          {
            confirmButtonText: t('common.confirm'),
            cancelButtonText: t('common.cancel'),
            type: 'warning'
          }
        )
        
        await api.delete(`/scans/${taskId}`)
        ElMessage.success(t('scan.deleteSuccess'))
        router.push('/scanning')
      } catch (error) {
        if (error !== 'cancel') {
          console.error('删除任务失败:', error)
          ElMessage.error(t('common.error') + ': ' + (error.response?.data?.message || error.message))
        }
      }
    }
    
    // 导航函数
    const navigateToEdit = () => {
      router.push(`/scanning/edit/${taskId}`)
    }
    
    const navigateToResults = () => {
      router.push(`/scanning/${taskId}/results`)
    }
    
    const goBack = () => {
      router.push('/scanning')
    }
    
    // 工具函数
    const formatDate = (dateString) => {
      if (!dateString) return '-'
      const date = new Date(dateString)
      return date.toLocaleString()
    }
    
    const calculateDuration = (startTime, endTime) => {
      const start = new Date(startTime)
      const end = new Date(endTime)
      const diff = end - start
      
      const seconds = Math.floor(diff / 1000)
      if (seconds < 60) return `${seconds}秒`
      
      const minutes = Math.floor(seconds / 60)
      if (minutes < 60) return `${minutes}分钟${seconds % 60}秒`
      
      const hours = Math.floor(minutes / 60)
      return `${hours}小时${minutes % 60}分钟${seconds % 60}秒`
    }
    
    const getStatusType = (status) => {
      switch (status) {
        case 'completed': return 'success'
        case 'running': return 'primary'
        case 'queued': return 'info'
        case 'created': return ''
        case 'failed': return 'danger'
        case 'cancelled': return 'warning'
        default: return ''
      }
    }
    
    // 处理WebSocket消息
    const handleScanTaskUpdate = (data) => {
      // 只处理当前扫描任务的更新
      if (data.id && data.id.toString() === taskId.toString()) {
        console.log('收到扫描任务更新:', data)
        
        // 更新进度信息
        if (data.progress) {
          progress.value = data.progress
        }
        
        if (data.stage) {
          progressStage.value = data.stage
        }
        
        if (data.details) {
          progressDetails.value = data.details
        }
        
        // 添加日志条目
        if (data.log) {
          progressLogs.value.push({
            timestamp: new Date(),
            level: data.log_level || 'info',
            message: data.log
          })
          
          // 如果日志太多，移除最早的条目
          if (progressLogs.value.length > 100) {
            progressLogs.value = progressLogs.value.slice(-100)
          }
          
          // 自动滚动到底部
          if (autoScroll.value) {
            scrollToBottom()
          }
        }
        
        // 如果任务状态发生变化且不是running，刷新任务详情
        if (data.status && 
            scanTask.value && 
            data.status !== scanTask.value.status && 
            data.status !== 'running') {
          loadTask()
        } else if (data.status) {
          // 直接更新状态以避免频繁刷新
          scanTask.value.status = data.status
        }
        
        // 如果发现新漏洞，更新计数
        if (data.new_vulnerabilities) {
          scanTask.value.total_vulnerabilities = (scanTask.value.total_vulnerabilities || 0) + data.new_vulnerabilities.length
          
          // 根据严重程度更新各类计数
          data.new_vulnerabilities.forEach(vuln => {
            switch (vuln.severity) {
              case 'critical':
                scanTask.value.critical_vulnerabilities = (scanTask.value.critical_vulnerabilities || 0) + 1
                break
              case 'high':
                scanTask.value.high_vulnerabilities = (scanTask.value.high_vulnerabilities || 0) + 1
                break
              case 'medium':
                scanTask.value.medium_vulnerabilities = (scanTask.value.medium_vulnerabilities || 0) + 1
                break
              case 'low':
                scanTask.value.low_vulnerabilities = (scanTask.value.low_vulnerabilities || 0) + 1
                break
            }
          })
        }
      }
    }
    
    // 滚动日志到底部
    const scrollToBottom = () => {
      if (logContainer.value) {
        setTimeout(() => {
          logContainer.value.scrollTop = logContainer.value.scrollHeight
        }, 50)
      }
    }
    
    // 格式化时间
    const formatTime = (timestamp) => {
      return new Date(timestamp).toLocaleTimeString()
    }
    
    // 组件挂载时
    onMounted(() => {
      // 建立WebSocket连接
      if (!wsService.isConnected) {
        wsService.connect()
      }
      
      // 添加事件监听
      wsService.addEventListener('scan_task_update', handleScanTaskUpdate)
      
      // 加载初始数据
      Promise.all([loadTask(), loadAssets()])
    })
    
    // 组件卸载时
    onUnmounted(() => {
      // 移除事件监听
      wsService.removeEventListener('scan_task_update', handleScanTaskUpdate)
    })
    
    return {
      scanTask,
      loading,
      actionLoading,
      taskId,
      assets,
      logContainer,
      progress,
      progressStage,
      progressDetails,
      progressLogs,
      autoScroll,
      progressStatus,
      progressPercentage,
      canStart,
      canCancel,
      canEdit,
      canDelete,
      startTask,
      cancelTask,
      confirmDelete,
      navigateToEdit,
      navigateToResults,
      goBack,
      getAssetName,
      formatDate,
      formatTime,
      calculateDuration,
      getStatusType,
      scrollToBottom
    }
  }
}
</script>

<style scoped>
.scan-detail-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.section-title {
  margin: 20px 0 10px 0;
  padding-bottom: 10px;
  border-bottom: 1px solid #eee;
}

.status-action-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.status-info {
  display: flex;
  align-items: center;
}

.scan-id {
  margin-left: 10px;
  color: #909399;
  font-size: 14px;
}

.target-assets {
  display: flex;
  flex-wrap: wrap;
  gap: 5px;
}

.asset-tag {
  margin: 2px;
}

.vuln-count {
  font-weight: bold;
  font-size: 16px;
}

.not-found {
  padding: 40px 0;
}

pre {
  white-space: pre-wrap;
  word-wrap: break-word;
  margin: 0;
  padding: 0;
}

.progress-section {
  margin-bottom: 20px;
}

.progress-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.progress-percentage {
  font-weight: bold;
  font-size: 16px;
}

.realtime-updates {
  margin-top: 10px;
}

.log-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.log-container {
  height: 200px;
  overflow-y: auto;
  padding: 10px;
  border: 1px solid #eee;
}

.log-entry {
  margin-bottom: 5px;
}

.log-time {
  font-size: 12px;
  color: #909399;
}

.log-message {
  font-size: 14px;
}

.log-info {
  color: #606266;
}

.log-warning {
  color: #E6A23C;
}

.log-error {
  color: #F56C6C;
}
</style> 