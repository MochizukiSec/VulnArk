<template>
  <div class="asset-detail" v-loading="loading">
    <!-- 返回按钮 -->
    <div class="back-link">
      <el-button @click="$router.back()" icon="el-icon-arrow-left" type="text">返回资产列表</el-button>
    </div>
    
    <!-- 资产详情头部 -->
    <div class="detail-header" v-if="asset">
      <div class="header-left">
        <h1 class="asset-name">
          <span class="asset-icon"><i class="el-icon-office-building"></i></span>
          {{ asset.name }}
        </h1>
        <div class="asset-meta">
          <el-tag :type="getAssetTypeTag(asset.type)" class="meta-tag">
            {{ getAssetTypeLabel(asset.type) }}
          </el-tag>
          <el-tag :type="getEnvironmentTag(asset.environment)" class="meta-tag">
            {{ getEnvironmentLabel(asset.environment) }}
          </el-tag>
          <span class="meta-date">创建于 {{ formatDate(asset.createdAt) }}</span>
        </div>
        <p class="asset-description">{{ asset.description }}</p>
      </div>
      <div class="header-actions">
        <el-button type="primary" @click="goToEdit" class="action-btn">
          <i class="el-icon-edit"></i> 编辑资产
        </el-button>
        <el-button type="danger" @click="confirmDelete" class="action-btn">
          <i class="el-icon-delete"></i> 删除资产
        </el-button>
      </div>
    </div>
    
    <!-- 资产详情内容 -->
    <div class="detail-content" v-if="asset">
      <el-tabs v-model="activeTab" type="border-card">
        <!-- 基本信息 -->
        <el-tab-pane label="基本信息" name="info">
          <el-card class="info-card">
            <div class="info-grid">
              <div class="info-item">
                <h3 class="info-label">ID</h3>
                <p class="info-value">{{ asset.id }}</p>
              </div>
              <div class="info-item">
                <h3 class="info-label">IP地址</h3>
                <p class="info-value">{{ asset.ipAddress || '未设置' }}</p>
              </div>
              <div class="info-item">
                <h3 class="info-label">域名</h3>
                <p class="info-value">{{ asset.domain || '未设置' }}</p>
              </div>
              <div class="info-item">
                <h3 class="info-label">操作系统</h3>
                <p class="info-value">{{ asset.operatingSystem || '未知' }}</p>
              </div>
              <div class="info-item">
                <h3 class="info-label">责任人</h3>
                <p class="info-value">{{ asset.owner || '未指定' }}</p>
              </div>
              <div class="info-item">
                <h3 class="info-label">部门</h3>
                <p class="info-value">{{ asset.department || '未指定' }}</p>
              </div>
              <div class="info-item">
                <h3 class="info-label">MAC地址</h3>
                <p class="info-value">{{ asset.macAddress || '未设置' }}</p>
              </div>
              <div class="info-item">
                <h3 class="info-label">端口</h3>
                <p class="info-value">{{ asset.ports ? asset.ports.join(', ') : '未扫描' }}</p>
              </div>
              <div class="info-item">
                <h3 class="info-label">最后更新</h3>
                <p class="info-value">{{ formatDate(asset.updatedAt) }}</p>
              </div>
            </div>
          </el-card>
          
          <!-- 自定义属性 -->
          <el-card class="custom-attributes-card" v-if="asset.attributes && Object.keys(asset.attributes).length > 0">
            <template #header>
              <div class="card-header">
                <span>自定义属性</span>
              </div>
            </template>
            <el-table :data="attributesArray" style="width: 100%">
              <el-table-column prop="key" label="属性名" width="200" />
              <el-table-column prop="value" label="值" />
            </el-table>
          </el-card>
        </el-tab-pane>
        
        <!-- 关联漏洞 -->
        <el-tab-pane label="关联漏洞" name="vulnerabilities">
          <div class="tab-header">
            <h3 class="tab-title">关联漏洞 <el-badge :value="vulnerabilities.length" type="danger" /></h3>
            <el-button type="primary" size="small" @click="dialogVisible = true">关联新漏洞</el-button>
          </div>
          
          <el-empty description="暂无关联漏洞" v-if="vulnerabilities.length === 0"></el-empty>
          
          <el-table 
            v-else
            :data="vulnerabilities" 
            style="width: 100%"
            border
          >
            <el-table-column label="ID" prop="id" width="80" />
            <el-table-column label="标题" prop="title" min-width="150">
              <template #default="scope">
                <router-link :to="{ name: 'VulnDetail', params: { id: scope.row.id } }" class="vuln-link">
                  {{ scope.row.title }}
                </router-link>
              </template>
            </el-table-column>
            <el-table-column label="CVE编号" prop="cveId" width="120" />
            <el-table-column label="严重程度" prop="severity" width="100">
              <template #default="scope">
                <el-tag :type="getSeverityTag(scope.row.severity)">
                  {{ getSeverityLabel(scope.row.severity) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="状态" prop="status" width="100">
              <template #default="scope">
                <el-tag :type="getStatusTag(scope.row.status)">
                  {{ getStatusLabel(scope.row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="120">
              <template #default="scope">
                <el-button 
                  size="small" 
                  type="danger" 
                  plain
                  @click="confirmRemoveVulnerability(scope.row)"
                >
                  移除
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
        
        <!-- 备注 -->
        <el-tab-pane label="备注" name="notes">
          <div class="tab-header">
            <h3 class="tab-title">资产备注</h3>
            <el-button type="primary" size="small" @click="showAddNoteDialog">添加备注</el-button>
          </div>
          
          <el-empty description="暂无备注" v-if="!asset.notes || asset.notes.length === 0"></el-empty>
          
          <div v-else class="notes-list">
            <el-timeline>
              <el-timeline-item
                v-for="(note, index) in asset.notes"
                :key="index"
                :timestamp="formatDate(note.createdAt)"
                placement="top"
                color="#409EFF"
              >
                <el-card class="note-card">
                  <div class="note-content">{{ note.content }}</div>
                  <div class="note-footer">
                    <span class="note-author">由 {{ note.createdBy.name || note.createdBy.username || '系统' }} 添加</span>
                  </div>
                </el-card>
              </el-timeline-item>
            </el-timeline>
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>
    
    <!-- 关联漏洞对话框 -->
    <el-dialog
      title="关联漏洞"
      v-model="dialogVisible"
      width="600px"
    >
      <el-form :model="searchForm" label-width="80px">
        <el-form-item label="漏洞搜索">
          <el-input 
            v-model="searchForm.searchTerm" 
            placeholder="输入漏洞名称、CVE编号或描述进行搜索"
            clearable
            @keyup.enter="searchVulnerabilities"
          >
            <template #append>
              <el-button icon="el-icon-search" @click="searchVulnerabilities"></el-button>
            </template>
          </el-input>
        </el-form-item>
      </el-form>
      
      <el-table
        v-loading="searchLoading"
        :data="searchResults"
        style="width: 100%"
        height="300"
        @row-click="handleVulnerabilitySelect"
      >
        <el-table-column type="selection" width="55" :selectable="isSelectable" />
        <el-table-column label="标题" prop="title" min-width="150" />
        <el-table-column label="CVE编号" prop="cveId" width="120" />
        <el-table-column label="严重程度" prop="severity" width="100">
          <template #default="scope">
            <el-tag :type="getSeverityTag(scope.row.severity)">
              {{ getSeverityLabel(scope.row.severity) }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="addSelectedVulnerabilities" :disabled="selectedVulnerabilities.length === 0">
            关联 ({{ selectedVulnerabilities.length }})
          </el-button>
        </span>
      </template>
    </el-dialog>
    
    <!-- 添加备注对话框 -->
    <el-dialog
      title="添加备注"
      v-model="noteDialogVisible"
      width="500px"
    >
      <el-form :model="noteForm" label-width="80px">
        <el-form-item label="备注内容" required>
          <el-input 
            v-model="noteForm.content" 
            type="textarea" 
            :rows="4"
            placeholder="请输入备注内容"
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="noteDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="addAssetNote" :disabled="!noteForm.content">添加</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessageBox, ElMessage } from 'element-plus'
import assetApi from '@/api/assets'
import vulnerabilityApi from '@/api/vulnerabilities'
import { formatDate } from '@/utils/helpers'

export default {
  name: 'AssetDetail',
  
  setup() {
    const route = useRoute()
    const router = useRouter()
    const loading = ref(false)
    const asset = ref(null)
    const activeTab = ref('info')
    const vulnerabilities = ref([])
    
    // 关联漏洞对话框
    const dialogVisible = ref(false)
    const searchForm = reactive({
      searchTerm: ''
    })
    const searchLoading = ref(false)
    const searchResults = ref([])
    const selectedVulnerabilities = ref([])
    
    // 备注对话框
    const noteDialogVisible = ref(false)
    const noteForm = reactive({
      content: ''
    })
    
    // 初始化资产详情
    const fetchAssetDetails = async () => {
      const assetId = route.params.id
      if (!assetId) return
      
      loading.value = true
      try {
        // 获取资产详情
        const response = await assetApi.getAssetById(assetId)
        asset.value = response.data
        
        // 获取关联漏洞
        await fetchAssetVulnerabilities()
      } catch (error) {
        console.error('获取资产详情失败:', error)
        ElMessage.error('获取资产详情失败，请重试')
      } finally {
        loading.value = false
      }
    }
    
    // 获取资产关联的漏洞
    const fetchAssetVulnerabilities = async () => {
      try {
        const response = await assetApi.getAssetVulnerabilities(route.params.id)
        vulnerabilities.value = response.data.vulnerabilities || []
      } catch (error) {
        console.error('获取关联漏洞失败:', error)
        ElMessage.error('获取关联漏洞失败，请重试')
      }
    }
    
    // 跳转到编辑页面
    const goToEdit = () => {
      router.push({ name: 'AssetEdit', params: { id: asset.value.id } })
    }
    
    // 确认删除资产
    const confirmDelete = () => {
      ElMessageBox.confirm(
        `确定要删除资产 "${asset.value.name}" 吗？此操作不可恢复。`,
        '警告',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning',
        }
      ).then(() => {
        deleteAsset()
      }).catch(() => {
        // 用户取消删除
      })
    }
    
    // 删除资产
    const deleteAsset = async () => {
      try {
        await assetApi.deleteAsset(asset.value.id)
        ElMessage.success('资产已成功删除')
        router.push({ name: 'AssetList' })
      } catch (error) {
        console.error('删除资产失败:', error)
        ElMessage.error('删除资产失败，请重试')
      }
    }
    
    // 确认移除关联漏洞
    const confirmRemoveVulnerability = (vulnerability) => {
      ElMessageBox.confirm(
        `确定要移除漏洞 "${vulnerability.title}" 与该资产的关联吗？`,
        '确认',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning',
        }
      ).then(() => {
        removeVulnerability(vulnerability.id)
      }).catch(() => {
        // 用户取消移除
      })
    }
    
    // 移除关联漏洞
    const removeVulnerability = async (vulnId) => {
      try {
        await assetApi.removeVulnerabilityFromAsset(route.params.id, vulnId)
        ElMessage.success('已移除漏洞关联')
        await fetchAssetVulnerabilities() // 刷新漏洞列表
      } catch (error) {
        console.error('移除漏洞关联失败:', error)
        ElMessage.error('移除漏洞关联失败，请重试')
      }
    }
    
    // 搜索漏洞
    const searchVulnerabilities = async () => {
      if (!searchForm.searchTerm || searchForm.searchTerm.trim() === '') {
        ElMessage.warning('请输入搜索内容')
        return
      }
      
      searchLoading.value = true
      try {
        // 将原始参数格式传递给API服务，让API服务内部处理参数转换
        const params = {
          searchTerm: searchForm.searchTerm,
          limit: 50
        }
        
        console.log('开始搜索漏洞，原始请求参数:', params)
        
        // 尝试使用后端API搜索
        const response = await vulnerabilityApi.searchVulnerabilities(params)
        console.log('搜索漏洞响应:', response.data)
        
        // 处理搜索结果
        searchResults.value = response.data.items || response.data.vulnerabilities || []
        
        // 过滤掉已关联的漏洞
        searchResults.value = searchResults.value.filter(vuln => {
          // 检查漏洞ID字段，适应不同的API返回格式
          const vulnId = vuln.id || vuln._id || vuln.cveId
          return !vulnerabilities.value.some(v => 
            (v.id === vulnId) || (v.cveId && v.cveId === vuln.cveId)
          )
        })
        
        if (searchResults.value.length === 0) {
          ElMessage.info('没有找到匹配的漏洞或所有匹配漏洞已关联')
        }
      } catch (error) {
        console.error('搜索漏洞失败详情:', error)
        console.error('请求配置:', error.config)
        console.error('响应状态:', error.response?.status)
        console.error('响应数据:', error.response?.data)
        if (error.request) {
          console.error('请求信息:', error.request)
        }
        ElMessage.error(`搜索漏洞失败: ${error.response?.data?.message || error.message || '请重试'}`)
      } finally {
        searchLoading.value = false
      }
    }
    
    // 选择漏洞
    const handleVulnerabilitySelect = (row) => {
      // 检查漏洞ID字段，适应不同的API返回格式
      const vulnId = row.id || row._id || row.cveId
      const index = selectedVulnerabilities.value.findIndex(v => 
        (v.id === vulnId) || (v.cveId && v.cveId === row.cveId)
      )
      
      if (index === -1) {
        selectedVulnerabilities.value.push(row)
      } else {
        selectedVulnerabilities.value.splice(index, 1)
      }
    }
    
    // 判断漏洞是否可选
    const isSelectable = (row) => {
      // 检查漏洞ID字段，适应不同的API返回格式
      const vulnId = row.id || row._id || row.cveId
      return !vulnerabilities.value.some(v => 
        (v.id === vulnId) || (v.cveId && v.cveId === row.cveId)
      )
    }
    
    // 添加选中的漏洞
    const addSelectedVulnerabilities = async () => {
      if (selectedVulnerabilities.value.length === 0) return
      
      try {
        // 处理每个选中的漏洞
        const promises = selectedVulnerabilities.value.map(async vuln => {
          try {
            // 获取漏洞ID，适应不同的API结构
            const vulnId = vuln.id || vuln._id
            
            // 如果是漏洞库返回的结构，先创建漏洞再关联
            if (!vulnId && vuln.cveId) {
              // 需要先导入漏洞，再关联
              console.log('导入漏洞库中的漏洞:', vuln)
              const vulnerabilityData = {
                title: vuln.title,
                description: vuln.description || '',
                cve: vuln.cveId,
                cvss: vuln.cvss || 0,
                severity: vuln.severity || 'medium',
                remediation: vuln.remediation || vuln.solution || '',
                references: vuln.references || [],
                tags: vuln.tags || []
              }
              
              console.log('准备导入漏洞数据:', vulnerabilityData)
              
              // 使用API服务导入漏洞
              const importResponse = await vulnerabilityApi.importFromVulnDatabase(vulnerabilityData)
              console.log('漏洞导入响应:', importResponse.data)
              
              const newVulnId = importResponse.data.id || importResponse.data.vulnerabilityId || importResponse.data._id
              if (newVulnId) {
                console.log('新漏洞ID:', newVulnId, '准备关联到资产:', route.params.id)
                try {
                  const relationResponse = await assetApi.addVulnerabilityToAsset(route.params.id, newVulnId)
                  console.log('漏洞关联响应:', relationResponse)
                  return relationResponse
                } catch (relationError) {
                  console.error('关联新导入漏洞时出错:', relationError.response || relationError)
                  throw relationError
                }
              } else {
                console.error('导入响应中找不到漏洞ID:', importResponse.data)
                throw new Error('导入漏洞失败：未返回有效的漏洞ID')
              }
            } else {
              // 直接关联已有漏洞
              console.log('关联已有漏洞:', vulnId, '到资产:', route.params.id)
              try {
                const relationResponse = await assetApi.addVulnerabilityToAsset(route.params.id, vulnId)
                console.log('漏洞关联响应:', relationResponse)
                return relationResponse
              } catch (relationError) {
                console.error('关联已有漏洞时出错:', relationError.response || relationError)
                throw relationError
              }
            }
          } catch (error) {
            console.error('处理单个漏洞时出错:', error.response || error)
            throw error // 重新抛出错误以便Promise.all捕获
          }
        })
        
        await Promise.all(promises)
        ElMessage.success('漏洞关联成功')
        dialogVisible.value = false
        selectedVulnerabilities.value = []
        searchResults.value = []
        searchForm.searchTerm = ''
        
        // 刷新漏洞列表
        await fetchAssetVulnerabilities()
      } catch (error) {
        console.error('关联漏洞失败详情:', error)
        if (error.response) {
          console.error('响应状态:', error.response.status)
          console.error('响应数据:', error.response.data)
          console.error('请求配置:', error.config)
          if (error.request) {
            console.error('请求信息:', error.request)
          }
        }
        ElMessage.error(`关联漏洞失败: ${error.response?.data?.message || error.message || '服务器内部错误，请检查后端日志'}`)
      }
    }
    
    // 显示添加备注对话框
    const showAddNoteDialog = () => {
      noteForm.content = ''
      noteDialogVisible.value = true
    }
    
    // 添加资产备注
    const addAssetNote = async () => {
      if (!noteForm.content || noteForm.content.trim() === '') {
        ElMessage.warning('请输入备注内容')
        return
      }
      
      try {
        await assetApi.addAssetNote(route.params.id, noteForm.content)
        ElMessage.success('备注添加成功')
        noteDialogVisible.value = false
        
        // 刷新资产详情
        await fetchAssetDetails()
      } catch (error) {
        console.error('添加备注失败:', error)
        ElMessage.error('添加备注失败，请重试')
      }
    }
    
    // 处理自定义属性为数组
    const attributesArray = computed(() => {
      if (!asset.value || !asset.value.attributes) return []
      
      return Object.entries(asset.value.attributes).map(([key, value]) => ({
        key,
        value: typeof value === 'object' ? JSON.stringify(value) : value
      }))
    })
    
    // 获取资产类型显示标签
    const getAssetTypeLabel = (type) => {
      const types = {
        server: '服务器',
        network_device: '网络设备',
        application: '应用程序',
        database: '数据库',
        cloud_resource: '云资源',
        iot_device: '物联网设备',
        other: '其他'
      }
      return types[type] || type
    }
    
    // 获取资产类型对应的标签样式
    const getAssetTypeTag = (type) => {
      const tags = {
        server: 'primary',
        network_device: 'success',
        application: 'warning',
        database: 'danger',
        cloud_resource: 'info',
        iot_device: '',
        other: 'info'
      }
      return tags[type] || 'info'
    }
    
    // 获取环境显示标签
    const getEnvironmentLabel = (env) => {
      const environments = {
        production: '生产',
        testing: '测试',
        development: '开发',
        staging: '预生产'
      }
      return environments[env] || env
    }
    
    // 获取环境对应的标签样式
    const getEnvironmentTag = (env) => {
      const tags = {
        production: 'danger',
        testing: 'warning',
        development: 'success',
        staging: 'info'
      }
      return tags[env] || 'info'
    }
    
    // 获取漏洞严重程度显示标签
    const getSeverityLabel = (severity) => {
      const severities = {
        critical: '严重',
        high: '高危',
        medium: '中危',
        low: '低危',
        info: '信息'
      }
      return severities[severity] || severity
    }
    
    // 获取漏洞严重程度对应的标签样式
    const getSeverityTag = (severity) => {
      const tags = {
        critical: 'danger',
        high: 'warning',
        medium: '',
        low: 'success',
        info: 'info'
      }
      return tags[severity] || 'info'
    }
    
    // 获取漏洞状态显示标签
    const getStatusLabel = (status) => {
      const statuses = {
        open: '开放',
        in_progress: '处理中',
        resolved: '已解决',
        closed: '已关闭',
        false_positive: '误报'
      }
      return statuses[status] || status
    }
    
    // 获取漏洞状态对应的标签样式
    const getStatusTag = (status) => {
      const tags = {
        open: 'danger',
        in_progress: 'warning',
        resolved: 'success',
        closed: 'info',
        false_positive: ''
      }
      return tags[status] || 'info'
    }
    
    onMounted(() => {
      fetchAssetDetails()
    })
    
    return {
      loading,
      asset,
      activeTab,
      vulnerabilities,
      dialogVisible,
      searchForm,
      searchLoading,
      searchResults,
      selectedVulnerabilities,
      noteDialogVisible,
      noteForm,
      attributesArray,
      formatDate,
      fetchAssetDetails,
      goToEdit,
      confirmDelete,
      confirmRemoveVulnerability,
      searchVulnerabilities,
      handleVulnerabilitySelect,
      isSelectable,
      addSelectedVulnerabilities,
      showAddNoteDialog,
      addAssetNote,
      getAssetTypeLabel,
      getAssetTypeTag,
      getEnvironmentLabel,
      getEnvironmentTag,
      getSeverityLabel,
      getSeverityTag,
      getStatusLabel,
      getStatusTag
    }
  }
}
</script>

<style scoped>
.asset-detail {
  padding: 20px;
}

.back-link {
  margin-bottom: 20px;
}

.detail-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 20px;
  padding-bottom: 20px;
  border-bottom: 1px solid #ebeef5;
}

.asset-name {
  font-size: 24px;
  margin: 0 0 15px 0;
  color: #303133;
  display: flex;
  align-items: center;
}

.asset-icon {
  margin-right: 10px;
  font-size: 24px;
  color: #409eff;
}

.asset-meta {
  display: flex;
  align-items: center;
  margin-bottom: 15px;
}

.meta-tag {
  margin-right: 10px;
}

.meta-date {
  color: #909399;
  font-size: 14px;
}

.asset-description {
  color: #606266;
  margin: 0;
  line-height: 1.6;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.info-card {
  margin-bottom: 20px;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
}

.info-item {
  padding: 10px;
}

.info-label {
  color: #909399;
  font-size: 14px;
  margin: 0 0 5px 0;
  font-weight: normal;
}

.info-value {
  color: #303133;
  margin: 0;
  font-size: 16px;
  word-break: break-all;
}

.custom-attributes-card {
  margin-top: 20px;
}

.card-header {
  font-weight: bold;
  color: #303133;
}

.tab-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.tab-title {
  margin: 0;
  font-size: 16px;
  color: #303133;
}

.notes-list {
  padding: 10px 0;
}

.note-card {
  margin-bottom: 10px;
}

.note-content {
  color: #303133;
  line-height: 1.6;
  white-space: pre-line;
}

.note-footer {
  display: flex;
  justify-content: flex-end;
  margin-top: 10px;
  color: #909399;
  font-size: 12px;
}

.vuln-link {
  color: #409eff;
  text-decoration: none;
  font-weight: bold;
}

.vuln-link:hover {
  text-decoration: underline;
}
</style> 