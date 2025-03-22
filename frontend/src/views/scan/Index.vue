<template>
  <div class="scan-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <h2>扫描任务</h2>
          <el-button type="primary" @click="navigateToAdd">
            创建任务
          </el-button>
        </div>
      </template>
      
      <!-- 批量操作和过滤器 -->
      <div class="top-actions">
        <!-- 批量操作按钮 -->
        <div class="batch-actions" v-if="selectedTasks.length > 0">
          <span class="selection-info">已选择 {{ selectedTasks.length }} 项</span>
          <el-button
            type="success"
            size="small"
            :loading="batchActionLoading"
            @click="batchStart"
          >
            批量启动
          </el-button>
          <el-button
            type="danger"
            size="small"
            :loading="batchActionLoading"
            @click="batchDelete"
          >
            批量删除
          </el-button>
          <el-button
            type="info"
            size="small"
            @click="clearSelection"
          >
            清除选择
          </el-button>
        </div>

        <!-- 过滤器 -->
        <div class="filter-container">
          <el-input
            v-model="searchQuery"
            placeholder="搜索"
            prefix-icon="Search"
            style="width: 200px"
            @keyup.enter="handleSearch"
          />
          <el-select 
            v-model="statusFilter" 
            placeholder="按状态筛选"
            style="width: 140px; margin-left: 10px;"
            clearable
            @change="handleSearch"
          >
            <el-option
              v-for="item in statusOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
          <el-button 
            type="primary" 
            icon="Search" 
            @click="handleSearch"
            style="margin-left: 10px;"
          >
            搜索
          </el-button>
          <el-button 
            type="info" 
            icon="RefreshRight"
            @click="resetSearch"
            style="margin-left: 10px;"
          >
            重置
          </el-button>
        </div>
      </div>

      <!-- 任务列表 -->
      <el-table
        v-loading="loading"
        :data="tasks"
        style="width: 100%"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column label="名称" min-width="180">
          <template #default="{ row }">
            <router-link :to="`/scanning/${row.id}`" class="task-name">
              {{ row.name }}
            </router-link>
          </template>
        </el-table-column>
        <el-table-column label="描述" min-width="200">
          <template #default="{ row }">
            {{ row.description || '暂无描述' }}
          </template>
        </el-table-column>
        <el-table-column label="扫描器类型" min-width="120">
          <template #default="{ row }">
            <el-tag>{{ getScannerType(row.scanner_type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="目标类型" min-width="120">
          <template #default="{ row }">
            <el-tag type="success">{{ getTargetType(row.target_type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="状态" min-width="120">
          <template #default="{ row }">
            <scan-status :status="row.status" :progress="row.progress" />
          </template>
        </el-table-column>
        <el-table-column label="创建时间" min-width="180">
          <template #default="{ row }">
            {{ new Date(row.created_at).toLocaleString() }}
          </template>
        </el-table-column>
        <el-table-column label="操作" min-width="250" fixed="right">
          <template #default="{ row }">
            <div class="operation-buttons">
              <el-button
                v-if="['ready', 'failed', 'completed', 'cancelled', 'error'].includes(row.status)"
                size="small"
                type="success"
                :loading="row.actionLoading"
                @click="startTask(row)"
                class="action-btn"
              >
                <el-icon><VideoPlay /></el-icon>
                <span>运行</span>
              </el-button>
              <el-button
                v-if="row.status === 'running'"
                size="small"
                type="warning"
                :loading="row.actionLoading"
                @click="cancelTask(row)"
                class="action-btn"
              >
                <el-icon><VideoPause /></el-icon>
                <span>取消</span>
              </el-button>
              <el-button
                size="small"
                type="primary"
                :disabled="['running', 'cancelling'].includes(row.status)"
                @click="navigateToEdit(row.id)"
                class="action-btn"
              >
                <el-icon><Edit /></el-icon>
                <span>编辑</span>
              </el-button>
              <el-button
                size="small"
                type="info"
                @click="navigateToDetail(row.id)"
                class="action-btn"
              >
                <el-icon><View /></el-icon>
                <span>详情</span>
              </el-button>
              <el-button
                size="small"
                type="danger"
                :loading="row.actionLoading"
                :disabled="['running', 'cancelling'].includes(row.status)"
                @click="deleteTask(row)"
                class="action-btn"
              >
                <el-icon><Delete /></el-icon>
                <span>删除</span>
              </el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          :current-page="pagination.currentPage"
          :page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import ScanStatus from '../../components/ScanStatus.vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  Edit, View, Delete, VideoPlay, VideoPause
} from '@element-plus/icons-vue'
import { 
  getScanTasks, 
  startScanTask, 
  cancelScanTask, 
  deleteScanTask 
} from '@/api/scan'

const router = useRouter()

// 状态选项
const statusOptions = [
  { value: 'ready', label: '就绪' },
  { value: 'running', label: '运行中' },
  { value: 'completed', label: '已完成' },
  { value: 'failed', label: '失败' },
  { value: 'cancelled', label: '已取消' },
  { value: 'error', label: '错误' }
]

// 数据和状态
const tasks = ref([])
const loading = ref(false)
const selectedTasks = ref([])
const batchActionLoading = ref(false)
const searchQuery = ref('')
const statusFilter = ref('')

// 分页
const pagination = reactive({
  currentPage: 1,
  pageSize: 10,
  total: 0
})

// 生命周期钩子
onMounted(() => {
  fetchTasks()
})

// 方法
const fetchTasks = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.currentPage,
      per_page: pagination.pageSize,
      name: searchQuery.value || undefined,
      status: statusFilter.value || undefined
    }
    
    console.log('请求扫描任务参数:', params)
    const response = await getScanTasks(params)
    console.log('扫描任务API响应:', response)
    
    // 更全面的响应格式检查
    if (response) {
      if (response.code === 200 && response.data) {
        let tasksData = []
        let totalCount = 0
        
        // 处理不同的数据结构
        if (Array.isArray(response.data.items)) {
          console.log('使用标准分页格式处理数据')
          tasksData = response.data.items
          totalCount = response.data.total || response.data.items.length
        } else if (Array.isArray(response.data)) {
          console.log('使用数组格式处理数据')
          tasksData = response.data
          totalCount = response.data.length
        } else if (response.data.list && Array.isArray(response.data.list)) {
          // 处理{list: [...], total: number}格式
          console.log('使用list数组格式处理数据')
          tasksData = response.data.list
          totalCount = response.data.total || response.data.list.length
        } else if (typeof response.data === 'object') {
          console.log('尝试将data作为单个任务处理')
          tasksData = [response.data]
          totalCount = 1
        }
        
        // 处理和设置任务数据
        tasks.value = tasksData.map(task => ({
          ...task,
          actionLoading: false
        }))
        pagination.total = totalCount
        console.log('成功加载扫描任务:', tasks.value.length, '条')
      } else {
        console.error('获取扫描任务失败:', response.message || '未知错误')
        tasks.value = []
        pagination.total = 0
        ElMessage.error('获取扫描任务失败: ' + (response.message || '服务器返回错误'))
      }
    } else {
      console.error('获取扫描任务响应为空')
      tasks.value = []
      pagination.total = 0
      ElMessage.error('获取扫描任务失败: 服务器未返回数据')
    }
  } catch (error) {
    console.error('获取扫描任务异常:', error)
    tasks.value = []
    pagination.total = 0
    ElMessage.error('获取扫描任务失败: ' + (error.message || '未知错误'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.currentPage = 1
  fetchTasks()
}

const resetSearch = () => {
  searchQuery.value = ''
  statusFilter.value = ''
  pagination.currentPage = 1
  fetchTasks()
}

const handleSizeChange = (size) => {
  pagination.pageSize = size
  fetchTasks()
}

const handleCurrentChange = (page) => {
  pagination.currentPage = page
  fetchTasks()
}

const handleSelectionChange = (selection) => {
  selectedTasks.value = selection
}

const clearSelection = () => {
  selectedTasks.value = []
}

const navigateToAdd = () => {
  router.push('/scanning/add')
}

const navigateToEdit = (id) => {
  router.push(`/scanning/edit/${id}`)
}

const navigateToDetail = (id) => {
  router.push(`/scanning/${id}`)
}

// 操作函数
const startTask = async (row) => {
  try {
    // 显示加载状态
    row.actionLoading = true;
    
    const response = await startScanTask(row.id)
    ElMessage({
      message: '任务启动成功',
      type: 'success',
      duration: 3000
    })
    fetchTasks()
  } catch (error) {
    console.error('启动任务失败:', error)
    
    // 使用详细的错误处理
    let errorMsg = '启动任务失败'
    if (error.response?.data?.message) {
      errorMsg += `: ${error.response.data.message}`
    } else if (error.message) {
      errorMsg += `: ${error.message}`
    }
    
    ElMessage({
      message: errorMsg,
      type: 'error',
      duration: 5000,
      showClose: true
    })
  } finally {
    // 清除加载状态
    if (row) row.actionLoading = false;
  }
}

const cancelTask = async (row) => {
  try {
    await ElMessageBox.confirm(
      '确定要取消当前运行的任务吗？',
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    // 显示加载状态
    row.actionLoading = true;
    
    const response = await cancelScanTask(row.id)
    ElMessage({
      message: '任务已取消',
      type: 'success',
      duration: 3000
    })
    fetchTasks()
  } catch (error) {
    if (error === 'cancel') {
      // 用户取消了操作
      return
    }
    
    console.error('取消任务失败:', error)
    
    // 详细的错误处理
    let errorMsg = '取消任务失败'
    if (error.response?.data?.message) {
      errorMsg += `: ${error.response.data.message}`
    } else if (error.message) {
      errorMsg += `: ${error.message}`
    }
    
    ElMessage({
      message: errorMsg,
      type: 'error',
      duration: 5000,
      showClose: true
    })
  } finally {
    // 清除加载状态
    if (row) row.actionLoading = false;
  }
}

const deleteTask = async (row) => {
  try {
    await ElMessageBox.confirm(
      '确定要删除此任务吗？此操作不可恢复',
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'danger'
      }
    )
    
    // 显示加载状态
    row.actionLoading = true;
    
    const response = await deleteScanTask(row.id)
    ElMessage({
      message: '任务删除成功',
      type: 'success',
      duration: 3000
    })
    fetchTasks()
  } catch (error) {
    if (error === 'cancel') {
      // 用户取消了操作
      return
    }
    
    console.error('删除任务失败:', error)
    
    // 详细的错误处理
    let errorMsg = '删除任务失败'
    if (error.response?.data?.message) {
      errorMsg += `: ${error.response.data.message}`
    } else if (error.message) {
      errorMsg += `: ${error.message}`
    }
    
    ElMessage({
      message: errorMsg,
      type: 'error',
      duration: 5000,
      showClose: true
    })
  } finally {
    // 清除加载状态
    if (row) row.actionLoading = false;
  }
}

const batchStart = async () => {
  try {
    await ElMessageBox.confirm(
      `确定要批量启动 ${selectedTasks.value.length} 个任务吗？`,
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    batchActionLoading.value = true
    
    // 批量处理所有选择的任务
    const promises = selectedTasks.value.map(task => 
      startScanTask(task.id).catch(error => {
        console.error(`启动任务 ${task.id} 失败:`, error)
        return error // 返回错误以便后续统计
      })
    )
    
    const results = await Promise.all(promises)
    
    // 统计成功和失败的数量
    const successCount = results.filter(r => !(r instanceof Error)).length
    const failCount = results.filter(r => r instanceof Error).length
    
    // 提示结果
    if (failCount === 0) {
      ElMessage({
        message: `成功启动 ${successCount} 个任务`,
        type: 'success',
        duration: 3000
      })
    } else {
      ElMessage({
        message: `${successCount} 个任务启动成功，${failCount} 个任务启动失败`,
        type: 'warning',
        duration: 5000,
        showClose: true
      })
    }
    
    // 刷新任务列表
    fetchTasks()
    // 清除选择
    clearSelection()
  } catch (error) {
    if (error === 'cancel') {
      // 用户取消了操作
      return
    }
    
    console.error('批量启动任务失败:', error)
    ElMessage({
      message: '批量操作失败',
      type: 'error',
      duration: 5000
    })
  } finally {
    batchActionLoading.value = false
  }
}

const batchDelete = async () => {
  try {
    await ElMessageBox.confirm(
      `确定要批量删除 ${selectedTasks.value.length} 个任务吗？此操作不可恢复`,
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'danger'
      }
    )
    
    batchActionLoading.value = true
    
    // 批量处理所有选择的任务
    const promises = selectedTasks.value.map(task => 
      deleteScanTask(task.id).catch(error => {
        console.error(`删除任务 ${task.id} 失败:`, error)
        return error // 返回错误以便后续统计
      })
    )
    
    const results = await Promise.all(promises)
    
    // 统计成功和失败的数量
    const successCount = results.filter(r => !(r instanceof Error)).length
    const failCount = results.filter(r => r instanceof Error).length
    
    // 提示结果
    if (failCount === 0) {
      ElMessage({
        message: `成功删除 ${successCount} 个任务`,
        type: 'success',
        duration: 3000
      })
    } else {
      ElMessage({
        message: `${successCount} 个任务删除成功，${failCount} 个任务删除失败`,
        type: 'warning',
        duration: 5000,
        showClose: true
      })
    }
    
    // 刷新任务列表
    fetchTasks()
    // 清除选择
    clearSelection()
  } catch (error) {
    if (error === 'cancel') {
      // 用户取消了操作
      return
    }
    
    console.error('批量删除任务失败:', error)
    ElMessage({
      message: '批量操作失败',
      type: 'error',
      duration: 5000
    })
  } finally {
    batchActionLoading.value = false
  }
}

// 辅助函数
const getScannerType = (type) => {
  const typeMap = {
    'nessus': 'Nessus',
    'awvs': 'AWVS',
    'xray': 'Xray',
    'zap': 'OWASP ZAP',
    'nuclei': 'Nuclei',
    'custom': '自定义'
  }
  return typeMap[type] || type
}

const getTargetType = (type) => {
  const typeMap = {
    'ip': 'IP地址',
    'url': 'URL',
    'domain': '域名',
    'app': '应用',
    'asset': '资产'
  }
  return typeMap[type] || type
}
</script>

<style lang="scss" scoped>
.scan-container {
  padding: 20px;
  
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    
    h2 {
      margin: 0;
      font-size: 18px;
    }
  }
  
  .top-actions {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
    flex-wrap: wrap;
    gap: 10px;
    
    .batch-actions {
      display: flex;
      align-items: center;
      gap: 10px;
      
      .selection-info {
        font-size: 14px;
        color: #606266;
      }
    }
    
    .filter-container {
      display: flex;
      align-items: center;
      flex-wrap: wrap;
      gap: 10px;
    }
  }
  
  .task-name {
    color: #409eff;
    text-decoration: none;
    font-weight: 500;
    
    &:hover {
      text-decoration: underline;
    }
  }
  
  .pagination-container {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
  }
}

.operation-buttons {
  display: flex;
  flex-wrap: wrap;
  gap: 5px;
  justify-content: center;
}

.action-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s;
}

.action-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.action-btn .el-icon {
  margin-right: 4px;
}

.action-btn:disabled {
  cursor: not-allowed;
  opacity: 0.6;
}

@media screen and (max-width: 768px) {
  .operation-buttons {
    flex-direction: column;
  }
}
</style>