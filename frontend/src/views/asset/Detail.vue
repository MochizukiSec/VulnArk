<template>
  <div class="asset-detail-container">
    <div class="page-header">
      <div class="left">
        <el-button link @click="goBack">
          <el-icon><Back /></el-icon>
          返回
        </el-button>
        <h2>资产详情</h2>
      </div>
      <div class="right">
        <el-button type="primary" @click="handleEdit">
          <el-icon><Edit /></el-icon>
          编辑
        </el-button>
      </div>
    </div>

    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>基本信息</span>
        </div>
      </template>
      
      <div class="detail-content" v-if="assetInfo">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="资产名称">{{ assetInfo.name }}</el-descriptions-item>
          <el-descriptions-item label="IP地址">{{ assetInfo.ip }}</el-descriptions-item>
          <el-descriptions-item label="资产类型">
            <el-tag :type="getAssetTypeTag(assetInfo.type)">{{ assetInfo.type }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="assetInfo.status === 'active' ? 'success' : 'danger'">
              {{ assetInfo.status === 'active' ? '在线' : '离线' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="所属部门">{{ assetInfo.department }}</el-descriptions-item>
          <el-descriptions-item label="负责人">{{ assetInfo.owner }}</el-descriptions-item>
          <el-descriptions-item label="操作系统">{{ assetInfo.operatingSystem }}</el-descriptions-item>
          <el-descriptions-item label="版本">{{ assetInfo.version }}</el-descriptions-item>
          <el-descriptions-item label="URL地址" v-if="assetInfo.type === 'website'">
            <a :href="assetInfo.url" target="_blank">{{ assetInfo.url }}</a>
          </el-descriptions-item>
          <el-descriptions-item label="MAC地址">{{ assetInfo.macAddress }}</el-descriptions-item>
          <el-descriptions-item label="创建时间">{{ assetInfo.createdAt }}</el-descriptions-item>
          <el-descriptions-item label="更新时间">{{ assetInfo.updatedAt }}</el-descriptions-item>
          <el-descriptions-item label="标签" :span="2">
            <el-tag
              v-for="tag in assetInfo.tags"
              :key="tag"
              class="tag-item"
            >
              {{ tag }}
            </el-tag>
            <span v-if="!assetInfo.tags || assetInfo.tags.length === 0">无</span>
          </el-descriptions-item>
          <el-descriptions-item label="备注" :span="2">
            {{ assetInfo.notes || '无' }}
          </el-descriptions-item>
        </el-descriptions>
      </div>
      
      <div class="empty-data" v-else>
        <el-empty description="未找到资产信息" />
      </div>
    </el-card>

    <!-- 相关漏洞列表 -->
    <el-card class="vuln-list" v-loading="vulnLoading">
      <template #header>
        <div class="card-header">
          <span>相关漏洞</span>
        </div>
      </template>
      
      <el-table
        :data="vulnList"
        border
        style="width: 100%"
      >
        <el-table-column prop="title" label="漏洞名称" min-width="200" show-overflow-tooltip />
        <el-table-column prop="severity" label="严重等级" width="100">
          <template #default="scope">
            <el-tag :type="getSeverityType(scope.row.severity)">
              {{ scope.row.severity }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="getStatusType(scope.row.status)">
              {{ scope.row.status }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="discoveredAt" label="发现时间" min-width="170" />
        <el-table-column fixed="right" label="操作" width="100">
          <template #default="scope">
            <el-button type="primary" link @click="handleViewVuln(scope.row)">查看</el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <div class="empty-data" v-if="!vulnList || vulnList.length === 0">
        <el-empty description="暂无相关漏洞" />
      </div>
      
      <!-- 分页组件 -->
      <div class="pagination-container" v-if="vulnPagination.total > 0">
        <el-pagination
          :current-page="vulnPagination.currentPage"
          :page-size="vulnPagination.pageSize"
          :total="vulnPagination.total"
          layout="total, prev, pager, next"
          @current-change="handlePageChange"
        />
      </div>
    </el-card>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Back, Edit } from '@element-plus/icons-vue'
import { getAssetDetail, getAssetVulnerabilities } from '@/api/asset'

export default {
  name: 'AssetDetail',
  components: {
    Back,
    Edit
  },
  setup() {
    const router = useRouter()
    const route = useRoute()
    const loading = ref(false)
    const vulnLoading = ref(false)
    const assetInfo = ref(null)
    const vulnList = ref([])
    const vulnPagination = ref({
      currentPage: 1,
      pageSize: 10,
      total: 0
    })

    // 获取资产ID
    const assetId = route.params.id

    // 获取资产详情
    const getAssetInfo = async () => {
      try {
        loading.value = true
        const response = await getAssetDetail(assetId)
        assetInfo.value = response.data
        
        // 获取相关漏洞
        getRelatedVulnerabilities()
      } catch (error) {
        console.error('获取资产详情失败:', error)
        ElMessage.error('获取资产详情失败')
      } finally {
        loading.value = false
      }
    }

    // 获取相关漏洞
    const getRelatedVulnerabilities = async () => {
      try {
        vulnLoading.value = true
        
        const params = {
          page: vulnPagination.value.currentPage,
          pageSize: vulnPagination.value.pageSize
        }
        
        console.log('正在请求资产关联漏洞', { assetId, params })
        
        const response = await getAssetVulnerabilities(assetId, params)
        console.log('获取资产关联漏洞响应:', response)
        
        if (response && response.code === 200) {
          vulnList.value = response.data.items || []
          vulnPagination.value.total = response.data.total || 0
          
          if (vulnList.value.length === 0) {
            ElMessage.info(response.message || '该资产暂无关联漏洞')
          }
        } else {
          vulnList.value = []
          vulnPagination.value.total = 0
          ElMessage.warning(response.message || '获取漏洞列表返回空数据')
        }
      } catch (error) {
        console.error('获取相关漏洞失败:', error)
        
        // 检查是否有响应数据
        if (error.response) {
          const statusCode = error.response.status
          const errorData = error.response.data || {}
          
          console.log('错误响应详情:', { statusCode, errorData })
          
          if (statusCode === 404) {
            ElMessage.warning('找不到指定资产，请刷新页面重试')
          } else if (statusCode === 500) {
            ElMessage.error(`服务器内部错误: ${errorData.message || '请联系管理员'}`)
          } else {
            ElMessage.error(`获取相关漏洞失败: ${errorData.message || error.message || '未知错误'}`)
          }
        } else {
          // 网络错误或请求被取消
          ElMessage.error(`获取相关漏洞失败: ${error.message || '网络错误'}`)
        }
        
        vulnList.value = []
        vulnPagination.value.total = 0
      } finally {
        vulnLoading.value = false
      }
    }

    // 处理分页变化
    const handlePageChange = (page) => {
      vulnPagination.value.currentPage = page
      getRelatedVulnerabilities()
    }

    // 获取资产类型对应的标签类型
    const getAssetTypeTag = (type) => {
      const typeMap = {
        'server': 'danger',
        'website': 'success',
        'database': 'warning',
        'network': 'info',
        'application': 'primary'
      }
      return typeMap[type] || ''
    }

    // 获取漏洞严重等级对应的标签类型
    const getSeverityType = (severity) => {
      const severityMap = {
        'critical': 'danger',
        'high': 'danger',
        'medium': 'warning',
        'low': 'info',
        'info': 'info'
      }
      return severityMap[severity] || ''
    }

    // 获取漏洞状态对应的标签类型
    const getStatusType = (status) => {
      const statusMap = {
        'open': 'danger',
        'in_progress': 'warning',
        'fixed': 'success',
        'closed': 'info'
      }
      return statusMap[status] || ''
    }

    // 返回列表页
    const goBack = () => {
      router.push('/assets')
    }

    // 编辑资产
    const handleEdit = () => {
      router.push(`/assets/edit/${assetId}`)
    }

    // 查看漏洞详情
    const handleViewVuln = (row) => {
      router.push(`/vulnerabilities/${row.id}`)
    }

    // 页面加载时获取数据
    onMounted(() => {
      getAssetInfo()
    })

    return {
      loading,
      vulnLoading,
      assetInfo,
      vulnList,
      vulnPagination,
      getAssetTypeTag,
      getSeverityType,
      getStatusType,
      goBack,
      handleEdit,
      handleViewVuln,
      handlePageChange
    }
  }
}
</script>

<style scoped>
.asset-detail-container {
  padding: 0;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-header .left {
  display: flex;
  align-items: center;
}

.page-header h2 {
  margin: 0 0 0 10px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.tag-item {
  margin-right: 8px;
  margin-bottom: 8px;
}

.empty-data {
  padding: 20px 0;
}

.vuln-list {
  margin-top: 20px;
}

.pagination-container {
  margin-top: 20px;
  text-align: right;
}
</style> 