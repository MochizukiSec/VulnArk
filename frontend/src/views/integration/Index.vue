<template>
  <div class="integration-container">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <span>CI/CD集成管理</span>
          <el-button type="primary" @click="handleAddIntegration" class="add-button">
            <el-icon><Plus /></el-icon>
            <span>添加集成</span>
          </el-button>
        </div>
      </template>
      
      <div v-loading="loading">
        <div v-if="integrations.length === 0" class="empty-tip">
          <el-empty description="暂无集成配置" />
          <el-button type="primary" @click="handleAddIntegration" class="mt-3 add-button">添加集成</el-button>
        </div>
        
        <el-table 
          v-else 
          :data="integrations" 
          border 
          style="width: 100%"
          :header-cell-style="{background:'#f5f7fa', color:'#606266'}"
          :cell-style="{padding: '8px 0'}"
          class="integration-table"
        >
          <el-table-column prop="name" label="名称" min-width="150" />
          <el-table-column prop="type" label="类型" width="120">
            <template #default="scope">
              <el-tag :type="getIntegrationType(scope.row.type).tagType">
                {{ getIntegrationType(scope.row.type).label }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
          <el-table-column prop="enabled" label="状态" width="100">
            <template #default="scope">
              <el-tag :type="scope.row.enabled ? 'success' : 'info'">
                {{ scope.row.enabled ? '启用' : '禁用' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" min-width="240" fixed="right">
            <template #default="scope">
              <div class="operation-buttons">
                <el-button
                  type="primary"
                  size="small"
                  @click="showIntegrationDetail(scope.row)"
                  class="action-btn"
                >
                  <el-icon><View /></el-icon>
                  详情
                </el-button>
                <el-button
                  type="success"
                  size="small"
                  @click="showApiKey(scope.row)"
                  class="action-btn"
                >
                  <el-icon><Key /></el-icon>
                  密钥
                </el-button>
                <el-button
                  :type="scope.row.enabled ? 'warning' : 'success'"
                  size="small"
                  @click="toggleStatus(scope.row)"
                  class="action-btn"
                >
                  <el-icon>
                    <component :is="scope.row.enabled ? 'VideoPause' : 'VideoPlay'" />
                  </el-icon>
                  {{ scope.row.enabled ? '暂停' : '启用' }}
                </el-button>
                <el-button
                  type="danger"
                  size="small"
                  @click="handleDelete(scope.row)"
                  class="action-btn"
                >
                  <el-icon><Delete /></el-icon>
                  删除
                </el-button>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-card>
    
    <!-- API密钥对话框 -->
    <el-dialog
      v-model="apiKeyDialogVisible"
      title="API密钥"
      width="550px"
    >
      <div class="api-key-container">
        <p>集成名称: <strong>{{ currentIntegration.name }}</strong></p>
        <p class="mb-3">请妥善保管此API密钥，它不会再次显示。</p>
        
        <el-input
          v-if="apiKey"
          v-model="apiKey"
          readonly
          :rows="2"
          type="textarea"
        >
          <template #append>
            <el-button @click="copyApiKey">
              <el-icon><CopyDocument /></el-icon>
            </el-button>
          </template>
        </el-input>
        
        <p class="mt-3">使用示例:</p>
        <el-card class="code-block">
          <pre><code>curl -X POST \
  {{ apiEndpoint }}/api/v1/webhooks/{{ currentIntegration.type }} \
  -H 'Content-Type: application/json' \
  -H 'X-API-Key: {{ apiKey || 'YOUR_API_KEY' }}' \
  -d '{"vulnerabilities": [...]}'
          </code></pre>
        </el-card>
        
        <el-alert
          title="更新API密钥会使当前密钥立即失效"
          type="warning"
          :closable="false"
          show-icon
          class="mt-3"
        />
      </div>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="apiKeyDialogVisible = false">关闭</el-button>
          <el-button 
            type="warning" 
            @click="regenerateApiKey" 
            :loading="regenerating"
          >重新生成</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { 
  Plus, View, Key, Delete, CopyDocument, 
  Link, SetUp, Guide, Document,
  VideoPause, VideoPlay
} from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  getIntegrations, 
  deleteIntegration, 
  regenerateApiKey as regenerateApiKeyAPI,
  toggleIntegrationStatus
} from '@/api/integration'

const router = useRouter()
const loading = ref(false)
const integrations = ref([])
const apiKeyDialogVisible = ref(false)
const apiKey = ref('')
const currentIntegration = reactive({
  id: null,
  name: '',
  type: ''
})
const regenerating = ref(false)

// 获取当前环境的API根URL
const apiEndpoint = window.location.origin

// 获取集成列表
const fetchIntegrations = async () => {
  loading.value = true
  try {
    const response = await getIntegrations()
    if (response.code === 200) {
      integrations.value = response.data || []
    } else {
      ElMessage.error(response.message || '获取集成配置失败')
    }
  } catch (error) {
    console.error('获取集成配置失败:', error)
    ElMessage.error('获取集成配置失败')
  } finally {
    loading.value = false
  }
}

// 获取集成类型对应的标签和图标
const getIntegrationType = (type) => {
  const types = {
    'jenkins': { label: 'Jenkins', tagType: 'danger', icon: SetUp },
    'gitlab': { label: 'GitLab', tagType: 'warning', icon: Guide },
    'github': { label: 'GitHub', tagType: 'success', icon: Link },
    'custom': { label: '自定义', tagType: 'info', icon: Document }
  }
  return types[type] || { label: type, tagType: '', icon: Document }
}

// 添加集成
const handleAddIntegration = () => {
  router.push('/integrations/add')
}

// 查看集成详情
const showIntegrationDetail = (integration) => {
  router.push(`/integrations/${integration.id}`)
}

// 显示API密钥
const showApiKey = async (integration) => {
  currentIntegration.id = integration.id
  currentIntegration.name = integration.name
  currentIntegration.type = integration.type
  
  // 由于安全原因，API返回的集成信息中不包含完整的API密钥
  // 需要用户点击"重新生成"按钮获取新的API密钥
  apiKey.value = ''
  apiKeyDialogVisible.value = true
}

// 重新生成API密钥
const regenerateApiKey = async () => {
  try {
    regenerating.value = true
    const response = await regenerateApiKeyAPI(currentIntegration.id)
    if (response.code === 200) {
      apiKey.value = response.data.api_key
      ElMessage.success('API密钥重新生成成功')
    } else {
      ElMessage.error(response.message || '重新生成API密钥失败')
    }
  } catch (error) {
    console.error('重新生成API密钥失败:', error)
    ElMessage.error('重新生成API密钥失败')
  } finally {
    regenerating.value = false
  }
}

// 复制API密钥
const copyApiKey = () => {
  navigator.clipboard.writeText(apiKey.value)
    .then(() => {
      ElMessage.success('API密钥已复制到剪贴板')
    })
    .catch(() => {
      ElMessage.error('复制失败，请手动复制')
    })
}

// 切换集成状态（启用/禁用）
const toggleStatus = async (integration) => {
  try {
    const newStatus = !integration.enabled
    const actionText = newStatus ? '启用' : '暂停'
    
    const response = await toggleIntegrationStatus(integration.id, newStatus)
    if (response.code === 200) {
      integration.enabled = newStatus
      ElMessage.success(`${actionText}成功`)
    } else {
      ElMessage.error(response.message || `${actionText}失败`)
    }
  } catch (error) {
    console.error('切换状态失败:', error)
    ElMessage.error('切换状态失败')
  }
}

// 删除集成
const handleDelete = (integration) => {
  ElMessageBox.confirm(
    `确定要删除集成 "${integration.name}" 吗？此操作不可恢复。`,
    '警告',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      const response = await deleteIntegration(integration.id)
      if (response.code === 200) {
        ElMessage.success('删除成功')
        fetchIntegrations()
      } else {
        ElMessage.error(response.message || '删除失败')
      }
    } catch (error) {
      console.error('删除集成失败:', error)
      ElMessage.error('删除集成失败')
    }
  }).catch(() => {
    // 用户取消删除
  })
}

// 初始化
onMounted(() => {
  fetchIntegrations()
})
</script>

<style scoped>
.integration-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 5px;
}

.empty-tip {
  margin: 60px 0;
  text-align: center;
}

.operation-buttons {
  display: flex;
  gap: 5px;
  justify-content: flex-start;
  white-space: nowrap;
}

.action-btn {
  display: flex;
  align-items: center;
  padding: 0 8px;
  min-width: 75px;
  justify-content: center;
}

.action-btn .el-icon {
  margin-right: 3px;
}

.api-key-container {
  padding: 10px;
}

.code-block {
  background-color: #f6f8fa;
  margin-top: 10px;
}

.code-block pre {
  margin: 0;
  padding: 10px;
  overflow-x: auto;
  font-family: monospace;
}

.mt-3 {
  margin-top: 12px;
}

.mb-3 {
  margin-bottom: 12px;
}

.add-button {
  padding: 8px 16px;
  font-weight: 500;
}

.integration-table {
  margin-top: 10px;
  border-radius: 4px;
  overflow: hidden;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.05);
}

.integration-table :deep(.el-table__row) {
  height: 55px;
}

.integration-table :deep(.el-button--small) {
  padding: 5px 10px;
  font-size: 12px;
}
</style> 