<template>
  <div class="integration-detail-container">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <span>集成详情</span>
          <div class="header-actions">
            <el-button @click="goBack" class="action-btn">
              <el-icon><Back /></el-icon>
              <span>返回</span>
            </el-button>
            <el-button type="primary" @click="showApiKey" class="action-btn">
              <el-icon><Key /></el-icon>
              <span>API密钥</span>
            </el-button>
          </div>
        </div>
      </template>
      
      <div v-loading="loading">
        <template v-if="integration.id">
          <el-descriptions title="基本信息" :column="2" border>
            <el-descriptions-item label="名称">{{ integration.name }}</el-descriptions-item>
            <el-descriptions-item label="类型">
              <el-tag :type="getIntegrationType(integration.type).tagType">
                {{ getIntegrationType(integration.type).label }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="描述">{{ integration.description || '无' }}</el-descriptions-item>
            <el-descriptions-item label="状态">
              <el-tag :type="integration.enabled ? 'success' : 'info'">
                {{ integration.enabled ? '启用' : '禁用' }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="创建时间">{{ integration.created_at }}</el-descriptions-item>
            <el-descriptions-item label="更新时间">{{ integration.updated_at }}</el-descriptions-item>
          </el-descriptions>
          
          <el-divider content-position="left">配置信息</el-divider>
          
          <template v-if="integration.type === 'jenkins'">
            <el-descriptions title="Jenkins配置" :column="1" border>
              <el-descriptions-item label="Jenkins URL">{{ parsedConfig.url || '未设置' }}</el-descriptions-item>
              <el-descriptions-item label="用户名">{{ parsedConfig.username || '未设置' }}</el-descriptions-item>
              <el-descriptions-item label="API Token">******** (已隐藏)</el-descriptions-item>
            </el-descriptions>
          </template>
          
          <template v-else-if="integration.type === 'gitlab'">
            <el-descriptions title="GitLab配置" :column="1" border>
              <el-descriptions-item label="GitLab URL">{{ parsedConfig.url || '未设置' }}</el-descriptions-item>
              <el-descriptions-item label="项目ID">{{ parsedConfig.project_id || '未设置' }}</el-descriptions-item>
              <el-descriptions-item label="私有Token">******** (已隐藏)</el-descriptions-item>
            </el-descriptions>
          </template>
          
          <template v-else-if="integration.type === 'github'">
            <el-descriptions title="GitHub配置" :column="1" border>
              <el-descriptions-item label="仓库所有者">{{ parsedConfig.owner || '未设置' }}</el-descriptions-item>
              <el-descriptions-item label="仓库名称">{{ parsedConfig.repo || '未设置' }}</el-descriptions-item>
              <el-descriptions-item label="访问令牌">******** (已隐藏)</el-descriptions-item>
            </el-descriptions>
          </template>
          
          <template v-else>
            <el-descriptions title="自定义配置" :column="1" border>
              <el-descriptions-item label="配置数据">
                <pre>{{ formatJson(integration.config) }}</pre>
              </el-descriptions-item>
            </el-descriptions>
          </template>
          
          <el-divider content-position="left">使用说明</el-divider>
          
          <el-alert
            title="如何设置CI/CD集成"
            type="info"
            :closable="false"
            description="在CI/CD流水线中添加一个步骤，将安全扫描的结果发送到VulnArk系统。"
            show-icon
            class="mb-3"
          />
          
          <el-tabs>
            <el-tab-pane label="Jenkins" name="jenkins">
              <div class="code-container">
                <el-alert
                  type="info"
                  :closable="false"
                  show-icon
                >
                  <template #title>
                    在Jenkins Pipeline中添加以下步骤:
                  </template>
                </el-alert>
                <pre><code>pipeline {
  agent any
  stages {
    stage('Security Scan') {
      steps {
        // 执行安全扫描
        sh 'security-scanner --output results.json'
        
        // 发送结果到VulnArk
        sh '''
          curl -X POST \
            {{ apiEndpoint }}/api/v1/webhooks/{{ integration.type }} \
            -H 'Content-Type: application/json' \
            -H 'X-API-Key: YOUR_API_KEY' \
            -d @results.json
        '''
      }
    }
  }
}</code></pre>
              </div>
            </el-tab-pane>
            
            <el-tab-pane label="GitLab CI" name="gitlab">
              <div class="code-container">
                <el-alert
                  type="info"
                  :closable="false"
                  show-icon
                >
                  <template #title>
                    在.gitlab-ci.yml文件中添加以下配置:
                  </template>
                </el-alert>
                <pre><code>security_scan:
  stage: test
  script:
    - security-scanner --output results.json
    - |
      curl -X POST \
        {{ apiEndpoint }}/api/v1/webhooks/{{ integration.type }} \
        -H 'Content-Type: application/json' \
        -H 'X-API-Key: YOUR_API_KEY' \
        -d @results.json
  artifacts:
    paths:
      - results.json</code></pre>
              </div>
            </el-tab-pane>
            
            <el-tab-pane label="GitHub Actions" name="github">
              <div class="code-container">
                <el-alert
                  type="info"
                  :closable="false"
                  show-icon
                >
                  <template #title>
                    在GitHub Actions工作流文件中添加以下步骤:
                  </template>
                </el-alert>
                <pre><code>name: Security Scan

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  security:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      
      - name: Run security scan
        run: |
          # 安装扫描工具
          npm install -g security-scanner
          
          # 运行扫描并输出结果
          security-scanner --output results.json
          
      - name: Send results to VulnArk
        run: |
          curl -X POST \
            {{ apiEndpoint }}/api/v1/webhooks/{{ integration.type }} \
            -H 'Content-Type: application/json' \
            -H 'X-API-Key: ${{ secrets.VULNARK_API_KEY }}' \
            -d @results.json</code></pre>
              </div>
            </el-tab-pane>
          </el-tabs>
          
          <el-divider content-position="left">历史记录</el-divider>
          
          <el-table :data="history" border style="width: 100%">
            <el-table-column prop="executed_at" label="执行时间" width="180" />
            <el-table-column prop="status" label="状态" width="100">
              <template #default="scope">
                <el-tag :type="scope.row.status === 'success' ? 'success' : 'danger'">
                  {{ scope.row.status === 'success' ? '成功' : '失败' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="total_records" label="记录数" width="100" />
            <el-table-column prop="success_count" label="成功数" width="100" />
            <el-table-column prop="error_count" label="失败数" width="100" />
            <el-table-column prop="message" label="消息" min-width="200" show-overflow-tooltip />
          </el-table>
          
          <div v-if="history.length === 0" class="empty-tip">
            <el-empty description="暂无历史记录" />
          </div>
        </template>
        
        <div v-else class="not-found">
          <el-empty description="未找到集成配置" />
          <el-button type="primary" @click="goBack" class="mt-3">返回列表</el-button>
        </div>
      </div>
    </el-card>
    
    <!-- API密钥对话框 -->
    <el-dialog
      v-model="apiKeyDialogVisible"
      title="API密钥"
      width="550px"
    >
      <div class="api-key-container">
        <p>集成名称: <strong>{{ integration.name }}</strong></p>
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
  {{ apiEndpoint }}/api/v1/webhooks/{{ integration.type }} \
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
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { 
  Back, Key, CopyDocument, 
  Link, SetUp, Guide, Document 
} from '@element-plus/icons-vue'
import { 
  getIntegrations, 
  getIntegrationHistory, 
  regenerateApiKey as regenerateApiKeyApi 
} from '@/api/integration'

const router = useRouter()
const route = useRoute()
const loading = ref(false)
const integration = reactive({})
const history = ref([])
const apiKeyDialogVisible = ref(false)
const apiKey = ref('')
const regenerating = ref(false)

// 获取当前环境的API根URL
const apiEndpoint = window.location.origin

// 从URL中获取集成ID
const integrationId = computed(() => route.params.id)

// 解析后的配置
const parsedConfig = computed(() => {
  try {
    return integration.config ? JSON.parse(integration.config) : {}
  } catch (e) {
    return {}
  }
})

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

// 获取集成详情
const fetchIntegrationDetail = async () => {
  loading.value = true
  try {
    const response = await getIntegrations()
    if (response.code === 200) {
      const found = response.data.find(item => item.id == integrationId.value)
      if (found) {
        Object.assign(integration, found)
      }
    } else {
      ElMessage.error(response.message || '获取集成详情失败')
    }
  } catch (error) {
    console.error('获取集成详情失败:', error)
    ElMessage.error('获取集成详情失败')
  } finally {
    loading.value = false
  }
}

// 获取集成历史记录
const fetchIntegrationHistory = async () => {
  try {
    const response = await getIntegrationHistory(integrationId.value)
    if (response.code === 200) {
      history.value = response.data || []
    } else {
      ElMessage.error(response.message || '获取历史记录失败')
    }
  } catch (error) {
    console.error('获取历史记录失败:', error)
    ElMessage.error('获取历史记录失败')
  }
}

// 显示API密钥对话框
const showApiKey = () => {
  apiKey.value = ''
  apiKeyDialogVisible.value = true
}

// 重新生成API密钥
const regenerateApiKey = async () => {
  try {
    regenerating.value = true
    const response = await regenerateApiKeyApi(integration.id)
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

// 返回列表页
const goBack = () => {
  router.push('/integrations')
}

// 格式化JSON
const formatJson = (jsonString) => {
  try {
    const parsed = JSON.parse(jsonString)
    return JSON.stringify(parsed, null, 2)
  } catch (e) {
    return jsonString
  }
}

// 初始化
onMounted(() => {
  fetchIntegrationDetail()
  fetchIntegrationHistory()
})
</script>

<style scoped>
.integration-detail-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-actions {
  display: flex;
  gap: 8px;
}

.action-btn {
  display: flex;
  align-items: center;
}

.action-btn .el-icon {
  margin-right: 4px;
}

.empty-tip, .not-found {
  margin: 40px 0;
  text-align: center;
}

.mb-3 {
  margin-bottom: 12px;
}

.mt-3 {
  margin-top: 12px;
}

.code-container {
  background-color: #f6f8fa;
  border-radius: 4px;
  padding: 16px;
  margin: 16px 0;
}

.code-container pre {
  margin: 16px 0 0 0;
  padding: 0;
  overflow-x: auto;
  font-family: monospace;
  white-space: pre-wrap;
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

.api-key-container {
  padding: 10px;
}
</style> 