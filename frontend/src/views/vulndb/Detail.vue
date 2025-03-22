<template>
  <div class="vulndb-detail-container">
    <div v-loading="loading">
      <div class="page-header">
        <div class="left">
          <h2>漏洞详情</h2>
        </div>
        <div class="right">
          <el-button @click="goBack">返回</el-button>
          <el-button type="primary" @click="handleEdit" v-if="vulnDetail">编辑</el-button>
        </div>
      </div>

      <div v-if="vulnDetail">
        <el-card class="basic-info-card">
          <template #header>
            <div class="vuln-header">
              <h3>
                {{ vulnDetail.title }}
                <el-tag 
                  :type="getSeverityType(vulnDetail.severity)" 
                  size="small" 
                  style="margin-left: 10px;"
                >
                  {{ getSeverityLabel(vulnDetail.severity) }}
                </el-tag>
              </h3>
              <div class="vuln-meta">
                <span v-if="vulnDetail.cve">
                  <strong>CVE编号：</strong>
                  <a :href="'https://cve.mitre.org/cgi-bin/cvename.cgi?name=' + vulnDetail.cve" target="_blank" class="link-type">
                    {{ vulnDetail.cve }}
                  </a>
                </span>
                <span v-if="vulnDetail.cwe">
                  <strong>CWE：</strong>
                  <a :href="'https://cwe.mitre.org/data/definitions/' + vulnDetail.cwe + '.html'" target="_blank" class="link-type">
                    {{ vulnDetail.cwe }}
                  </a>
                </span>
                <span v-if="vulnDetail.cvss">
                  <strong>CVSS：</strong>{{ vulnDetail.cvss }}
                </span>
                <span v-if="vulnDetail.published_date">
                  <strong>发布日期：</strong>{{ formatDate(vulnDetail.published_date) }}
                </span>
                <span v-if="vulnDetail.last_modified_date">
                  <strong>最后修改：</strong>{{ formatDate(vulnDetail.last_modified_date) }}
                </span>
              </div>
              <div class="vuln-tags" v-if="vulnDetail.tags && vulnDetail.tags.length > 0">
                <strong>标签：</strong>
                <el-tag 
                  v-for="tag in vulnDetail.tags" 
                  :key="tag" 
                  size="small" 
                  effect="plain"
                  style="margin-right: 5px;"
                >
                  {{ tag }}
                </el-tag>
              </div>
            </div>
          </template>

          <div class="section">
            <h4>漏洞描述</h4>
            <p v-html="formattedDescription"></p>
          </div>

          <div class="section" v-if="vulnDetail.cvss_vector">
            <h4>CVSS向量</h4>
            <el-tag type="info">{{ vulnDetail.cvss_vector }}</el-tag>
          </div>

          <div class="section" v-if="vulnDetail.exploit_available">
            <h4>利用代码</h4>
            <el-alert
              title="存在已知的漏洞利用代码，请及时修复！"
              type="error"
              :closable="false"
              show-icon
            >
              <div v-if="vulnDetail.exploit_info">{{ vulnDetail.exploit_info }}</div>
            </el-alert>
          </div>

          <div class="section" v-if="vulnDetail.affected_products && vulnDetail.affected_products.length > 0">
            <h4>受影响产品</h4>
            <el-tag 
              v-for="product in vulnDetail.affected_products" 
              :key="product" 
              style="margin-right: 5px; margin-bottom: 5px;"
            >
              {{ product }}
            </el-tag>
          </div>

          <div class="section" v-if="vulnDetail.remediation">
            <h4>解决方案</h4>
            <p v-html="formattedRemediation"></p>
          </div>

          <div class="section" v-if="vulnDetail.references && vulnDetail.references.length > 0">
            <h4>参考资料</h4>
            <ul class="reference-list">
              <li v-for="(ref, index) in vulnDetail.references" :key="index">
                <a :href="ref" target="_blank" class="link-type">{{ ref }}</a>
              </li>
            </ul>
          </div>
        </el-card>
      </div>

      <div v-else-if="!loading" class="not-found">
        <el-empty description="未找到漏洞信息"></el-empty>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getVulnDBDetail, getVulnSeverities } from '@/api/vulndb'

export default {
  name: 'VulnDBDetail',
  setup() {
    const route = useRoute()
    const router = useRouter()
    const loading = ref(true)
    const vulnDetail = ref(null)
    const severityOptions = getVulnSeverities()

    // 获取漏洞详情
    const getVulnDetail = async () => {
      const id = route.params.id
      if (!id) {
        ElMessage.error('无效的漏洞ID')
        loading.value = false
        return
      }

      loading.value = true
      try {
        const response = await getVulnDBDetail(id)
        if (response.code === 200 && response.data) {
          vulnDetail.value = response.data
        } else {
          ElMessage.error(response.message || '获取漏洞详情失败')
        }
      } catch (error) {
        console.error('获取漏洞详情失败:', error)
        ElMessage.error('获取漏洞详情失败')
      } finally {
        loading.value = false
      }
    }

    // 格式化日期
    const formatDate = (dateStr) => {
      if (!dateStr) return '-'
      const date = new Date(dateStr)
      return date.toLocaleDateString('zh-CN')
    }

    // 获取严重程度类型
    const getSeverityType = (severity) => {
      const typeMap = {
        critical: 'danger',
        high: 'danger',
        medium: 'warning',
        low: 'success',
        info: 'info'
      }
      return typeMap[severity] || 'info'
    }

    // 获取严重程度标签
    const getSeverityLabel = (severity) => {
      const labelMap = {
        critical: '严重',
        high: '高危',
        medium: '中危',
        low: '低危',
        info: '信息'
      }
      return labelMap[severity] || severity
    }

    // 处理描述和修复方案的格式
    const formattedDescription = computed(() => {
      return vulnDetail.value?.description?.replace(/\n/g, '<br>') || ''
    })

    const formattedRemediation = computed(() => {
      return vulnDetail.value?.remediation?.replace(/\n/g, '<br>') || ''
    })

    // 返回上一页
    const goBack = () => {
      router.back()
    }

    // 跳转到编辑页面
    const handleEdit = () => {
      router.push(`/vulndb/edit/${vulnDetail.value.id}`)
    }

    onMounted(() => {
      getVulnDetail()
    })

    return {
      loading,
      vulnDetail,
      formatDate,
      getSeverityType,
      getSeverityLabel,
      formattedDescription,
      formattedRemediation,
      goBack,
      handleEdit
    }
  }
}
</script>

<style scoped>
.vulndb-detail-container {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-header h2 {
  margin: 0;
  font-size: 22px;
  color: #303133;
}

.basic-info-card {
  margin-bottom: 20px;
}

.vuln-header h3 {
  margin-top: 0;
  margin-bottom: 15px;
  font-size: 18px;
  color: #303133;
}

.vuln-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 15px;
  margin-bottom: 10px;
}

.vuln-tags {
  margin-top: 10px;
}

.section {
  margin-bottom: 20px;
}

.section h4 {
  font-size: 16px;
  color: #303133;
  margin-bottom: 10px;
  font-weight: 600;
}

.reference-list {
  padding-left: 20px;
  margin: 0;
}

.reference-list li {
  margin-bottom: 5px;
}

.link-type {
  color: #409EFF;
  text-decoration: none;
}

.link-type:hover {
  text-decoration: underline;
}

.not-found {
  padding: 40px 0;
  text-align: center;
}
</style> 