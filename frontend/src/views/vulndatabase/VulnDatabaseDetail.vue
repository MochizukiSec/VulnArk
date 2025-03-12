<template>
  <div class="vuln-detail-container">
    <div class="header-section">
      <div class="vuln-detail-header">
        <el-page-header :title="'漏洞数据库'" @back="goBack">
          <template #content>
            <div class="header-content">
              <span class="cve-id">{{ vulnerability.cveId }}</span>
              <el-tag 
                v-if="vulnerability.severity" 
                :type="getSeverityType(vulnerability.severity)" 
                size="large" 
                class="severity-tag"
              >
                {{ getSeverityLabel(vulnerability.severity) }}
              </el-tag>
            </div>
          </template>
        </el-page-header>
        
        <div class="header-actions">
          <el-button type="primary" @click="importVulnerability">
            <el-icon><DocumentAdd /></el-icon>
            添加到漏洞列表
          </el-button>
        </div>
      </div>
      
      <h1 class="vuln-title">{{ vulnerability.title }}</h1>
    </div>
    
    <el-skeleton :loading="loading" animated :rows="10" v-if="loading">
    </el-skeleton>
    
    <div v-else class="vuln-detail-content">
      <el-row :gutter="20">
        <el-col :span="24">
          <el-card class="info-card">
            <template #header>
              <div class="card-header">
                <h3>基本信息</h3>
              </div>
            </template>
            
            <el-descriptions :column="3" border>
              <el-descriptions-item label="CVE ID">
                <el-link 
                  type="primary" 
                  :href="`https://cve.mitre.org/cgi-bin/cvename.cgi?name=${vulnerability.cveId}`" 
                  target="_blank"
                >
                  {{ vulnerability.cveId }}
                </el-link>
              </el-descriptions-item>
              
              <el-descriptions-item label="CVSS 评分">
                <span class="cvss-score" :class="getCvssScoreClass(vulnerability.cvss)">
                  {{ vulnerability.cvss }}
                </span>
              </el-descriptions-item>
              
              <el-descriptions-item label="严重程度">
                <el-tag 
                  :type="getSeverityType(vulnerability.severity)" 
                  size="default"
                >
                  {{ getSeverityLabel(vulnerability.severity) }}
                </el-tag>
              </el-descriptions-item>
              
              <el-descriptions-item label="发布日期">
                {{ formatDate(vulnerability.publishedDate) }}
              </el-descriptions-item>
              
              <el-descriptions-item label="最后更新日期">
                {{ formatDate(vulnerability.lastModifiedDate) }}
              </el-descriptions-item>
              
              <el-descriptions-item label="受影响系统">
                {{ vulnerability.affectedSystems || '未提供' }}
              </el-descriptions-item>
            </el-descriptions>
          </el-card>
        </el-col>
        
        <el-col :span="24" class="mt-4">
          <el-card class="info-card">
            <template #header>
              <div class="card-header">
                <h3>漏洞描述</h3>
              </div>
            </template>
            
            <div class="description-content">
              <p v-if="vulnerability.description">{{ vulnerability.description }}</p>
              <p v-else>暂无描述信息</p>
            </div>
          </el-card>
        </el-col>
        
        <el-col :span="24" class="mt-4">
          <el-card class="info-card">
            <template #header>
              <div class="card-header">
                <h3>解决方案</h3>
              </div>
            </template>
            
            <div class="solution-content">
              <p v-if="vulnerability.solution">{{ vulnerability.solution }}</p>
              <p v-else>暂无解决方案信息</p>
            </div>
          </el-card>
        </el-col>
        
        <el-col :span="24" class="mt-4">
          <el-card class="info-card">
            <template #header>
              <div class="card-header">
                <h3>参考链接</h3>
              </div>
            </template>
            
            <div class="references-content">
              <template v-if="vulnerability.references && vulnerability.references.length">
                <ul class="reference-list">
                  <li v-for="(ref, index) in vulnerability.references" :key="index">
                    <el-link :href="ref" type="primary" target="_blank">{{ ref }}</el-link>
                  </li>
                </ul>
              </template>
              <p v-else>暂无参考链接</p>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { DocumentAdd } from '@element-plus/icons-vue'
import axios from 'axios'

export default {
  name: 'VulnDatabaseDetail',
  
  components: {
    DocumentAdd
  },
  
  setup() {
    const route = useRoute()
    const router = useRouter()
    
    const loading = ref(true)
    const vulnerability = ref({})
    
    // 获取漏洞详情
    const fetchVulnerabilityDetail = async () => {
      const cveId = route.params.cveId
      
      if (!cveId) {
        ElMessage.error('未提供有效的CVE ID')
        goBack()
        return
      }
      
      loading.value = true
      
      try {
        console.log('获取漏洞详情:', {
          baseURL: axios.defaults.baseURL,
          endpoint: `/api/vulndatabase/${cveId}`,
          cveId
        })
        
        const response = await axios.get(`/api/vulndatabase/${cveId}`)
        console.log('获取漏洞详情成功:', response.data)
        vulnerability.value = response.data
      } catch (error) {
        console.error('获取漏洞详情失败:', error)
        console.error('详细错误信息:', {
          message: error.message,
          response: error.response ? {
            status: error.response.status,
            data: error.response.data
          } : '无响应'
        })
        ElMessage.error('获取漏洞详情失败，请稍后重试')
      } finally {
        loading.value = false
      }
    }
    
    // 格式化日期
    const formatDate = (dateString) => {
      if (!dateString) return '未知'
      
      try {
        const date = new Date(dateString)
        return date.toLocaleDateString('zh-CN', {
          year: 'numeric',
          month: 'long',
          day: 'numeric'
        })
      } catch (e) {
        return dateString
      }
    }
    
    // 获取严重程度类型
    const getSeverityType = (severity) => {
      const types = {
        'critical': 'danger',
        'high': 'danger',
        'medium': 'warning',
        'low': 'info',
        'info': 'info'
      }
      return types[severity] || 'info'
    }
    
    // 获取严重程度标签
    const getSeverityLabel = (severity) => {
      const labels = {
        'critical': '严重',
        'high': '高危',
        'medium': '中危',
        'low': '低危',
        'info': '信息'
      }
      return labels[severity] || '未知'
    }
    
    // 获取CVSS评分样式类
    const getCvssScoreClass = (score) => {
      if (score >= 9.0) return 'critical'
      if (score >= 7.0) return 'high'
      if (score >= 4.0) return 'medium'
      return 'low'
    }
    
    // 返回上一页
    const goBack = () => {
      router.push('/vulndatabase')
    }
    
    // 导入漏洞到系统
    const importVulnerability = async () => {
      try {
        const vuln = vulnerability.value
        
        const importData = {
          title: vuln.title,
          description: vuln.description,
          cve: vuln.cveId,
          cvss: vuln.cvss,
          severity: vuln.severity,
          affectedSystems: vuln.affectedSystems,
          references: vuln.references?.join('\n') || '',
          remediation: vuln.solution
        }
        
        await ElMessageBox.confirm(
          '确定要将此漏洞添加到您的漏洞管理列表中吗？',
          '导入漏洞',
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'info'
          }
        )
        
        const response = await axios.post('/api/vulnerabilities', importData)
        
        ElMessage.success('漏洞已成功导入到您的漏洞管理列表')
        
        // 可选：导入后跳转到漏洞详情页
        router.push(`/vulnerabilities/${response.data.id}`)
      } catch (error) {
        if (error === 'cancel') return
        
        console.error('导入漏洞失败:', error)
        ElMessage.error('导入漏洞失败，请稍后重试')
      }
    }
    
    onMounted(() => {
      fetchVulnerabilityDetail()
    })
    
    return {
      loading,
      vulnerability,
      formatDate,
      getSeverityType,
      getSeverityLabel,
      getCvssScoreClass,
      goBack,
      importVulnerability
    }
  }
}
</script>

<style scoped>
.vuln-detail-container {
  padding: 20px;
  background-color: #f6f8fb;
  min-height: 100vh;
}

.header-section {
  background-color: white;
  border-radius: 10px;
  padding: 20px;
  margin-bottom: 24px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.05);
  transition: all 0.3s ease;
}

.header-section:hover {
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.08);
}

.vuln-detail-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  padding-bottom: 16px;
  border-bottom: 1px solid #ebeef5;
}

.header-content {
  display: flex;
  align-items: center;
  gap: 12px;
}

.cve-id {
  font-family: 'Courier New', monospace;
  font-size: 16px;
  font-weight: 600;
  background-color: #f5f7fa;
  padding: 4px 10px;
  border-radius: 4px;
  transition: all 0.3s;
}

.cve-id:hover {
  background-color: #e4e7ed;
  transform: translateY(-2px);
}

.vuln-title {
  font-size: 26px;
  margin-top: 5px;
  margin-bottom: 5px;
  color: #303133;
  font-weight: 600;
  line-height: 1.4;
  position: relative;
  display: inline-block;
}

.vuln-title::after {
  content: '';
  position: absolute;
  bottom: -4px;
  left: 0;
  width: 60px;
  height: 3px;
  background: linear-gradient(90deg, #409EFF, #67C23A);
  border-radius: 2px;
  transition: width 0.3s ease;
}

.vuln-title:hover::after {
  width: 100%;
}

.header-actions .el-button {
  transition: all 0.3s ease;
  border-radius: 6px;
  padding: 10px 16px;
  font-weight: 500;
}

.header-actions .el-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.info-card {
  margin-bottom: 24px;
  border-radius: 10px;
  overflow: hidden;
  transition: all 0.3s ease;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.05);
}

.info-card:hover {
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.08);
  transform: translateY(-2px);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 20px;
  border-bottom: 1px solid #ebeef5;
  background-color: #fafafa;
}

.card-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #303133;
  position: relative;
  padding-left: 16px;
}

.card-header h3::before {
  content: '';
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%);
  width: 4px;
  height: 16px;
  background: linear-gradient(to bottom, #409EFF, #67C23A);
  border-radius: 2px;
}

.description-content,
.solution-content,
.references-content {
  padding: 16px 20px;
  line-height: 1.8;
  color: #606266;
}

.reference-list {
  list-style-type: none;
  padding-left: 0;
  margin: 0;
}

.reference-list li {
  margin-bottom: 12px;
  padding: 8px 12px;
  background-color: #f5f7fa;
  border-radius: 4px;
  transition: all 0.3s;
}

.reference-list li:hover {
  background-color: #eef5fe;
  transform: translateX(4px);
}

.reference-list li:last-child {
  margin-bottom: 0;
}

.reference-list .el-link {
  transition: all 0.3s;
  word-break: break-all;
}

.reference-list .el-link:hover {
  transform: translateY(-1px);
}

.cvss-score {
  font-weight: 600;
  padding: 4px 10px;
  border-radius: 12px;
  display: inline-block;
  transition: all 0.3s;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
}

.cvss-score:hover {
  transform: scale(1.05);
}

.cvss-score.critical {
  background: linear-gradient(135deg, #F56C6C, #d63131);
  color: white;
}

.cvss-score.high {
  background: linear-gradient(135deg, #E6A23C, #c4811a);
  color: white;
}

.cvss-score.medium {
  background: linear-gradient(135deg, #E6A23C, #c4811a);
  color: white;
}

.cvss-score.low {
  background: linear-gradient(135deg, #67C23A, #4e9a2a);
  color: white;
}

.severity-tag {
  text-transform: capitalize;
  transition: all 0.3s;
}

.severity-tag:hover {
  transform: scale(1.05);
}

:deep(.el-descriptions__label) {
  font-weight: 600;
  color: #606266;
}

:deep(.el-descriptions__content) {
  font-size: 14px;
}

:deep(.el-page-header__icon) {
  transition: all 0.3s;
}

:deep(.el-page-header__icon:hover) {
  transform: translateX(-4px);
}

:deep(.el-page-header__title) {
  font-weight: 600;
  color: #606266;
}

.mt-4 {
  margin-top: 16px;
}

/* 动画效果 */
.vuln-detail-content {
  animation: fadeInUp 0.5s ease-out;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 手机适配 */
@media (max-width: 768px) {
  .vuln-detail-header {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .header-actions {
    margin-top: 16px;
    width: 100%;
  }
  
  .header-actions .el-button {
    width: 100%;
  }
  
  .vuln-title {
    font-size: 22px;
  }
  
  :deep(.el-descriptions-item) {
    margin-bottom: 8px;
  }
  
  .card-header h3 {
    font-size: 16px;
  }
}
</style>