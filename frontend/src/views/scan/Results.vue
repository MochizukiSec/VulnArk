<template>
  <div class="scan-results-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <h2>扫描结果 - {{ taskName }}</h2>
          <div class="header-actions">
            <el-button type="success" @click="importAll" :disabled="importingAll || noResultsToImport">
              导入全部
            </el-button>
            <el-button @click="navigateToTask">返回</el-button>
          </div>
        </div>
      </template>

      <div class="filter-section">
        <el-input
          v-model="searchName"
          placeholder="搜索"
          prefix-icon="Search"
          clearable
          @keyup.enter="loadResults"
          @clear="loadResults"
          style="width: 300px;"
        />

        <el-select
          v-model="filterSeverity"
          placeholder="严重等级"
          clearable
          @change="loadResults"
        >
          <el-option
            v-for="item in severityOptions"
            :key="item.value"
            :label="item.label"
            :value="item.value"
          />
        </el-select>

        <el-select
          v-model="filterImported"
          placeholder="导入状态"
          clearable
          @change="loadResults"
        >
          <el-option
            label="已导入"
            value="true"
          />
          <el-option
            label="未导入"
            value="false"
          />
        </el-select>
      </div>

      <!-- 结果数据表格 -->
      <el-table
        v-loading="loading"
        :data="results"
        border
        style="width: 100%"
        empty-text="暂无数据"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        
        <el-table-column prop="id" label="ID" width="70" />
        
        <el-table-column label="漏洞名称" min-width="200">
          <template #default="{ row }">
            <div class="vuln-name">
              <span>{{ row.vulnerability_name }}</span>
              <el-tag v-if="row.is_imported" size="small" type="success">
                已导入
              </el-tag>
            </div>
          </template>
        </el-table-column>
        
        <el-table-column label="严重等级" width="120">
          <template #default="{ row }">
            <el-tag :type="getSeverityType(row.severity)">
              {{ getSeverityLabel(row.severity) }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column label="受影响资产" min-width="180">
          <template #default="{ row }">
            <div v-if="row.affected_url">
              <strong>URL:</strong> {{ truncateText(row.affected_url, 30) }}
            </div>
            <div v-if="row.affected_ip">
              <strong>IP:</strong> {{ row.affected_ip }}
              <span v-if="row.affected_port">:{{ row.affected_port }}</span>
            </div>
          </template>
        </el-table-column>
        
        <el-table-column label="分类" width="150">
          <template #default="{ row }">
            {{ row.category || '未知' }}
          </template>
        </el-table-column>
        
        <el-table-column label="CVSS" width="100" align="center">
          <template #default="{ row }">
            <span v-if="row.cvss">{{ row.cvss.toFixed(1) }}</span>
            <span v-else>-</span>
          </template>
        </el-table-column>
        
        <el-table-column label="创建时间" width="170">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button
              size="small"
              @click="showDetail(row)"
              type="info"
              title="详情"
            >
              <el-icon><View /></el-icon>
            </el-button>
            
            <el-button
              v-if="!row.is_imported"
              size="small"
              type="success"
              @click="importResult(row.id)"
              title="导入"
            >
              <el-icon><Upload /></el-icon>
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <!-- 多选操作栏 -->
      <div class="batch-actions" v-if="selectedResults.length > 0">
        <el-button type="primary" @click="importSelected" :disabled="importingSelected">
          导入选中项 ({{ selectedResults.length }})
        </el-button>
      </div>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          :current-page="currentPage"
          :page-sizes="[10, 20, 50, 100]"
          :page-size="pageSize"
          layout="total, sizes, prev, pager, next, jumper"
          :total="total"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>
    
    <!-- 结果详情对话框 -->
    <el-dialog
      v-model="detailDialogVisible"
      :title="detailItem ? detailItem.vulnerability_name : ''"
      width="70%"
    >
      <div v-if="detailItem" class="result-detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="严重等级">
            <el-tag :type="getSeverityType(detailItem.severity)">
              {{ getSeverityLabel(detailItem.severity) }}
            </el-tag>
          </el-descriptions-item>
          
          <el-descriptions-item label="分类">
            {{ detailItem.category || '未知' }}
          </el-descriptions-item>
          
          <el-descriptions-item label="CVSS" v-if="detailItem.cvss">
            {{ detailItem.cvss.toFixed(1) }}
          </el-descriptions-item>
          
          <el-descriptions-item label="CVE" v-if="detailItem.cve">
            {{ detailItem.cve }}
          </el-descriptions-item>
          
          <el-descriptions-item label="描述" :span="2">
            {{ detailItem.description }}
          </el-descriptions-item>
          
          <el-descriptions-item label="受影响URL" v-if="detailItem.affected_url" :span="2">
            {{ detailItem.affected_url }}
          </el-descriptions-item>
          
          <el-descriptions-item label="受影响IP" v-if="detailItem.affected_ip">
            {{ detailItem.affected_ip }}
          </el-descriptions-item>
          
          <el-descriptions-item label="受影响端口" v-if="detailItem.affected_port">
            {{ detailItem.affected_port }}
          </el-descriptions-item>
          
          <el-descriptions-item label="复现步骤" v-if="detailItem.detail" :span="2">
            <pre>{{ detailItem.detail }}</pre>
          </el-descriptions-item>
          
          <el-descriptions-item label="解决方案" v-if="detailItem.solution" :span="2">
            <pre>{{ detailItem.solution }}</pre>
          </el-descriptions-item>
          
          <el-descriptions-item label="参考链接" v-if="detailItem.references" :span="2">
            <pre>{{ detailItem.references }}</pre>
          </el-descriptions-item>
          
          <el-descriptions-item label="导入状态">
            <el-tag :type="detailItem.is_imported ? 'success' : 'info'">
              {{ detailItem.is_imported ? '已导入' : '未导入' }}
            </el-tag>
          </el-descriptions-item>
          
          <el-descriptions-item label="导入时间" v-if="detailItem.is_imported && detailItem.imported_at">
            {{ formatDate(detailItem.imported_at) }}
          </el-descriptions-item>
        </el-descriptions>
        
        <div class="dialog-footer">
          <el-button
            v-if="!detailItem.is_imported"
            type="primary"
            @click="importResult(detailItem.id)"
          >
            导入
          </el-button>
          <el-button @click="detailDialogVisible = false">
            关闭
          </el-button>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { View, Upload } from '@element-plus/icons-vue'
import api from '@/utils/request'

// 组件
const components = {
  View,
  Upload
}

const route = useRoute()
const router = useRouter()

// 获取任务ID
const taskId = route.params.id

// 数据状态
const loading = ref(true)
const results = ref([])
const selectedResults = ref([])
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

// 筛选数据
const searchName = ref('')
const filterSeverity = ref('')
const filterImported = ref('')

// 导入状态
const importingSelected = ref(false)
const importingAll = ref(false)

// 详情弹窗
const detailDialogVisible = ref(false)
const detailItem = ref(null)

// 任务信息
const taskName = ref('')

// 可选严重程度
const severityOptions = [
  { value: 'critical', label: '严重' },
  { value: 'high', label: '高危' },
  { value: 'medium', label: '中危' },
  { value: 'low', label: '低危' },
  { value: 'info', label: '信息' }
]

// 计算属性
const noResultsToImport = computed(() => {
  return !results.value.some(r => !r.is_imported)
})

// 加载扫描结果
const loadResults = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value
    }
    
    if (searchName.value) {
      params.search = searchName.value
    }
    
    if (filterSeverity.value) {
      params.severity = filterSeverity.value
    }
    
    if (filterImported.value) {
      params.is_imported = filterImported.value
    }
    
    const res = await api.get(`/scans/${taskId}/results`, { params })
    results.value = res.data.results || []
    total.value = res.data.total || 0
    
    // 加载任务信息以获取任务名称
    if (!taskName.value) {
      loadTaskInfo()
    }
  } catch (error) {
    console.error('加载扫描结果失败:', error)
    ElMessage.error(`错误: ${error.response?.data?.message || error.message}`)
  } finally {
    loading.value = false
  }
}

// 加载任务信息
const loadTaskInfo = async () => {
  try {
    const res = await api.get(`/scans/${taskId}`)
    taskName.value = res.data.name || `任务 #${taskId}`
  } catch (error) {
    console.error('加载任务信息失败:', error)
    taskName.value = `任务 #${taskId}`
  }
}

// 分页处理
const handleSizeChange = (size) => {
  pageSize.value = size
  loadResults()
}

const handleCurrentChange = (page) => {
  currentPage.value = page
  loadResults()
}

// 表格多选
const handleSelectionChange = (selection) => {
  selectedResults.value = selection
}

// 导入操作
const importResult = async (resultId) => {
  try {
    await api.post(`/scans/${taskId}/import`, {
      result_ids: [resultId]
    })
    
    ElMessage.success('导入成功')
    loadResults() // 重新加载数据
  } catch (error) {
    console.error('导入结果失败:', error)
    ElMessage.error(`错误: ${error.response?.data?.message || error.message}`)
  }
}

const importSelected = async () => {
  if (selectedResults.value.length === 0) return
  
  try {
    importingSelected.value = true
    
    // 筛选出未导入的结果
    const idsToImport = selectedResults.value
      .filter(result => !result.is_imported)
      .map(result => result.id)
    
    if (idsToImport.length === 0) {
      ElMessage.warning('所有选中的结果已经导入')
      importingSelected.value = false
      return
    }
    
    await api.post(`/scans/${taskId}/import`, {
      result_ids: idsToImport
    })
    
    ElMessage.success('导入成功')
    loadResults() // 重新加载数据
  } catch (error) {
    console.error('导入选中结果失败:', error)
    ElMessage.error(`错误: ${error.response?.data?.message || error.message}`)
  } finally {
    importingSelected.value = false
  }
}

const importAll = async () => {
  try {
    await ElMessageBox.confirm(
      '确定要导入全部扫描结果吗？',
      '确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'info'
      }
    )
    
    importingAll.value = true
    
    await api.post(`/scans/${taskId}/import`, {})
    
    ElMessage.success('导入成功')
    loadResults() // 重新加载数据
  } catch (error) {
    if (error !== 'cancel') {
      console.error('导入所有结果失败:', error)
      ElMessage.error(`错误: ${error.response?.data?.message || error.message}`)
    }
  } finally {
    importingAll.value = false
  }
}

// 详情对话框
const showDetail = (row) => {
  detailItem.value = row
  detailDialogVisible.value = true
}

// 导航函数
const navigateToTask = () => {
  router.push(`/scanning/${taskId}`)
}

// 工具函数
const formatDate = (dateString) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return date.toLocaleString()
}

const getSeverityType = (severity) => {
  switch (severity) {
    case 'critical': return 'danger'
    case 'high': return 'error'
    case 'medium': return 'warning'
    case 'low': return 'info'
    case 'info': return ''
    default: return ''
  }
}

const getSeverityLabel = (severity) => {
  switch (severity) {
    case 'critical': return '严重'
    case 'high': return '高危'
    case 'medium': return '中危'
    case 'low': return '低危'
    case 'info': return '信息'
    default: return severity
  }
}

const truncateText = (text, maxLength) => {
  if (!text) return ''
  return text.length > maxLength ? text.substring(0, maxLength) + '...' : text
}

onMounted(() => {
  loadResults()
})
</script>

<style scoped>
.scan-results-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.filter-section {
  display: flex;
  gap: 10px;
  margin-bottom: 20px;
}

.vuln-name {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.batch-actions {
  margin-top: 15px;
  padding: 10px;
  background-color: #f5f7fa;
  border-radius: 4px;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}

.result-detail {
  margin-top: 10px;
}

.dialog-footer {
  margin-top: 20px;
  text-align: right;
}

pre {
  white-space: pre-wrap;
  word-wrap: break-word;
  margin: 0;
  padding: 0;
}
</style> 