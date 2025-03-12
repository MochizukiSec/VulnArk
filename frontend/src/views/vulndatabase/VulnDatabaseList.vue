<template>
  <div class="vulndb-container">
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">漏洞库</h1>
        <p class="page-subtitle">搜索和浏览常见漏洞数据库，支持CVE查询和多种过滤条件</p>
      </div>
      <div class="header-actions">
        <el-button type="success" icon="el-icon-plus" @click="goToCreate">新增漏洞</el-button>
        <el-button type="primary" icon="el-icon-refresh" @click="fetchVulnerabilities">刷新</el-button>
      </div>
    </div>

    <!-- 搜索和过滤区域 -->
    <el-card class="filter-card" shadow="hover">
      <div class="search-filters">
        <div class="search-box">
          <el-input
            v-model="searchParams.searchTerm"
            placeholder="搜索CVE ID、标题或关键词"
            clearable
            prefix-icon="el-icon-search"
            @keyup.enter="handleSearch"
          >
            <template #append>
              <el-button @click="handleSearch">搜索</el-button>
            </template>
          </el-input>
        </div>
        
        <div class="filter-options">
          <el-select v-model="searchParams.year" placeholder="年份" clearable style="width: 120px;">
            <el-option
              v-for="year in yearOptions"
              :key="year"
              :label="year ? year.toString() : '全部'"
              :value="year"
            />
          </el-select>
          
          <el-select v-model="searchParams.severity" placeholder="严重程度" clearable style="width: 120px;">
            <el-option
              v-for="option in severityOptions"
              :key="option.value"
              :label="option.label"
              :value="option.value"
            />
          </el-select>
          
          <el-select v-model="searchParams.cvssRange" placeholder="CVSS评分" clearable style="width: 140px;">
            <el-option
              v-for="option in cvssRangeOptions"
              :key="option.value"
              :label="option.label"
              :value="option.value"
            />
          </el-select>
          
          <el-select v-model="searchParams.sortBy" placeholder="排序字段" style="width: 140px;">
            <el-option
              v-for="option in sortOptions"
              :key="option.value"
              :label="option.label"
              :value="option.value"
            />
          </el-select>
          
          <el-select v-model="searchParams.sortOrder" placeholder="排序方式" style="width: 100px;">
            <el-option key="asc" label="升序" value="asc" />
            <el-option key="desc" label="降序" value="desc" />
          </el-select>
          
          <el-button type="primary" @click="handleSearch">应用筛选</el-button>
          <el-button @click="resetFilters">重置</el-button>
        </div>
      </div>
    </el-card>

    <!-- 漏洞列表 -->
    <el-card shadow="hover" class="vulndb-table-card">
      <div class="table-header">
        <span class="table-title">漏洞列表</span>
        <div class="table-actions">
          <el-select v-model="pageSize" placeholder="每页显示" style="width: 120px;">
            <el-option :value="10" label="10条/页" />
            <el-option :value="20" label="20条/页" />
            <el-option :value="50" label="50条/页" />
            <el-option :value="100" label="100条/页" />
          </el-select>
        </div>
      </div>
      
      <el-table
        v-loading="loading"
        :data="vulnerabilities"
        stripe
        border
        style="width: 100%"
        class="vulndb-table"
        @row-click="handleRowClick"
      >
        <el-table-column prop="cveId" label="CVE ID" width="150" sortable>
          <template #default="scope">
            <el-link type="primary" @click.stop="goToDetail(scope.row.cveId)">
              {{ scope.row.cveId }}
            </el-link>
          </template>
        </el-table-column>
        
        <el-table-column prop="title" label="漏洞标题" min-width="200">
          <template #default="scope">
            <div class="vuln-title">
              <span class="title-text">{{ scope.row.title }}</span>
              <el-tag v-if="scope.row.tags && scope.row.tags.length" 
                size="small" 
                effect="plain" 
                class="tag-item">
                {{ scope.row.tags[0] }}
              </el-tag>
            </div>
          </template>
        </el-table-column>
        
        <el-table-column prop="cvss" label="CVSS" width="100" sortable>
          <template #default="scope">
            <div class="cvss-score">
              <el-progress 
                :percentage="scope.row.cvss * 10" 
                :color="getCvssColor(scope.row.cvss)"
                :show-text="false"
                :stroke-width="4"
                class="cvss-progress"
              />
              <span class="cvss-value">{{ scope.row.cvss.toFixed(1) }}</span>
            </div>
          </template>
        </el-table-column>
        
        <el-table-column prop="severity" label="严重程度" width="120" sortable>
          <template #default="scope">
            <el-tag :type="getSeverityType(scope.row.severity)" effect="dark" size="small">
              {{ getSeverityText(scope.row.severity) }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="publishedDate" label="发布日期" width="120" sortable>
          <template #default="scope">
            {{ formatDate(scope.row.publishedDate) }}
          </template>
        </el-table-column>
        
        <el-table-column prop="lastModifiedDate" label="更新日期" width="120" sortable>
          <template #default="scope">
            {{ formatDate(scope.row.lastModifiedDate) }}
          </template>
        </el-table-column>
        
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="scope">
            <el-button 
              type="primary" 
              size="small" 
              plain 
              icon="el-icon-view" 
              @click.stop="goToDetail(scope.row.cveId)">
              详情
            </el-button>
            <el-button 
              type="success" 
              size="small" 
              plain 
              icon="el-icon-plus" 
              @click.stop="importToVulnerability(scope.row)"
              title="导入到漏洞管理">
              导入
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          background
          layout="total, prev, pager, next, jumper"
          :total="total"
          :page-size="pageSize"
          :current-page="currentPage"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>
  </div>
</template>

<script>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import axios from 'axios'

export default {
  name: 'VulnDatabaseList',
  
  setup() {
    const router = useRouter()
    const loading = ref(false)
    const vulnerabilities = ref([])
    const total = ref(0)
    const currentPage = ref(1)
    const pageSize = ref(20)
    
    // 搜索参数
    const searchParams = reactive({
      searchTerm: '',
      year: null,
      severity: '',
      cvssRange: '',
      sortBy: 'publishedDate',
      sortOrder: 'desc'
    })
    
    // 年份选项，从2000年到当前年份
    const yearOptions = computed(() => {
      const currentYear = new Date().getFullYear()
      const years = []
      for (let year = currentYear; year >= 2000; year--) {
        years.push(year)
      }
      return years
    })
    
    // 严重程度选项
    const severityOptions = [
      { value: 'critical', label: '严重' },
      { value: 'high', label: '高危' },
      { value: 'medium', label: '中危' },
      { value: 'low', label: '低危' },
      { value: 'info', label: '信息' }
    ]
    
    // CVSS评分范围选项
    const cvssRangeOptions = [
      { value: '9-10', label: '9.0 - 10.0' },
      { value: '7-8.9', label: '7.0 - 8.9' },
      { value: '4-6.9', label: '4.0 - 6.9' },
      { value: '0-3.9', label: '0.0 - 3.9' }
    ]
    
    // 排序选项
    const sortOptions = [
      { value: 'publishedDate', label: '发布日期' },
      { value: 'lastModifiedDate', label: '更新日期' },
      { value: 'cveId', label: 'CVE ID' },
      { value: 'cvss', label: 'CVSS评分' }
    ]
    
    // 获取漏洞库数据
    const fetchVulnerabilities = async () => {
      try {
        loading.value = true
        
        // 构建查询参数
        const params = {
          page: currentPage.value,
          perPage: pageSize.value,
          searchTerm: searchParams.searchTerm,
          severity: searchParams.severity,
          year: searchParams.year,
          cvssRange: searchParams.cvssRange,
          sortBy: searchParams.sortBy,
          sortOrder: searchParams.sortOrder
        }
        
        console.log('请求漏洞库数据:', {
          baseURL: axios.defaults.baseURL,
          endpoint: '/api/vulndatabase',
          params
        })
        
        const response = await axios.get('/api/vulndatabase', { params })
        
        console.log('成功获取漏洞库数据:', response.data)
        vulnerabilities.value = response.data.items || []
        total.value = response.data.total || 0
        
        if (vulnerabilities.value.length === 0 && total.value > 0 && currentPage.value > 1) {
          // 如果当前页没有数据但总数大于0，回到第一页
          currentPage.value = 1
          fetchVulnerabilities()
        }
      } catch (error) {
        console.error('获取漏洞库数据失败:', error)
        ElMessage.error('获取漏洞库数据失败，请稍后重试')
        vulnerabilities.value = []
        total.value = 0
      } finally {
        loading.value = false
      }
    }
    
    // 搜索处理
    const handleSearch = () => {
      console.log('执行搜索，搜索参数:', searchParams)
      currentPage.value = 1 // 重置到第一页
      fetchVulnerabilities()
    }
    
    // 重置过滤条件
    const resetFilters = () => {
      searchParams.searchTerm = ''
      searchParams.year = null
      searchParams.severity = ''
      searchParams.cvssRange = ''
      searchParams.sortBy = 'publishedDate'
      searchParams.sortOrder = 'desc'
      currentPage.value = 1
      fetchVulnerabilities()
    }
    
    // 分页处理
    const handleCurrentChange = (page) => {
      currentPage.value = page
      fetchVulnerabilities()
    }
    
    // 行点击处理
    const handleRowClick = (row) => {
      goToDetail(row.cveId)
    }
    
    // 跳转到新增漏洞页面
    const goToCreate = () => {
      router.push('/vulndatabase/create')
    }
    
    // 跳转到详情页面
    const goToDetail = (cveId) => {
      router.push(`/vulndatabase/${cveId}`)
    }
    
    // 导入到漏洞管理
    const importToVulnerability = (vuln) => {
      ElMessageBox.confirm(
        `确定要将 ${vuln.cveId} 导入到漏洞管理系统吗？`,
        '导入确认',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'info'
        }
      ).then(async () => {
        try {
          // 调用API将漏洞导入到漏洞管理系统
          const vulnerabilityData = {
            title: vuln.title,
            description: vuln.description,
            cve: vuln.cveId,
            cvss: vuln.cvss,
            severity: vuln.severity,
            remediation: vuln.remediation || '',
            references: vuln.references || [],
            tags: vuln.tags || []
          }
          
          await axios.post('/vulnerabilities/import-from-vulndb', {
            vulnerability: vulnerabilityData
          })
          
          ElMessage.success(`成功导入 ${vuln.cveId} 到漏洞管理系统`)
        } catch (error) {
          console.error('导入漏洞失败:', error)
          ElMessage.error('导入漏洞失败，请稍后重试')
        }
      }).catch(() => {})
    }
    
    // 格式化日期
    const formatDate = (dateString) => {
      if (!dateString) return '未知'
      
      try {
        const date = new Date(dateString)
        // 检查日期是否有效
        if (isNaN(date.getTime())) {
          return '日期格式错误'
        }
        
        return new Intl.DateTimeFormat('zh-CN', {
          year: 'numeric',
          month: '2-digit',
          day: '2-digit'
        }).format(date)
      } catch (error) {
        console.error('日期格式化错误:', error, dateString)
        return '日期格式错误'
      }
    }
    
    // 获取严重程度对应的ElementUI标签类型
    const getSeverityType = (severity) => {
      const types = {
        critical: 'danger',
        high: 'warning',
        medium: 'warning',
        low: 'success',
        info: 'info'
      }
      return types[severity] || 'info'
    }
    
    // 获取严重程度的中文描述
    const getSeverityText = (severity) => {
      const texts = {
        critical: '严重',
        high: '高危',
        medium: '中危',
        low: '低危',
        info: '信息'
      }
      return texts[severity] || '未知'
    }
    
    // 获取CVSS评分对应的颜色
    const getCvssColor = (score) => {
      if (score >= 9.0) return '#F56C6C' // 严重
      if (score >= 7.0) return '#E6A23C' // 高危
      if (score >= 4.0) return '#F0C050' // 中危
      return '#67C23A' // 低危
    }
    
    onMounted(() => {
      // 使用真实API获取数据
      fetchVulnerabilities()
    })
    
    return {
      vulnerabilities,
      loading,
      total,
      currentPage,
      pageSize,
      searchParams,
      yearOptions,
      severityOptions,
      cvssRangeOptions,
      sortOptions,
      handleSearch,
      resetFilters,
      fetchVulnerabilities,
      handleCurrentChange,
      handleRowClick,
      goToCreate,
      goToDetail,
      importToVulnerability,
      formatDate,
      getSeverityType,
      getSeverityText,
      getCvssColor
    }
  }
}
</script>

<style lang="scss" scoped>
.vulndb-container {
  padding: 20px;
  background-color: #f6f8fb;
  min-height: 100vh;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  padding-bottom: 20px;
  border-bottom: 1px solid #e0e6ed;
  position: relative;
}

.page-header::after {
  content: '';
  position: absolute;
  bottom: -1px;
  left: 0;
  width: 80px;
  height: 3px;
  background: linear-gradient(90deg, #409EFF, #67C23A);
  border-radius: 2px;
}

.page-title {
  font-size: 26px;
  font-weight: 600;
  color: #303133;
  margin: 0 0 10px 0;
  display: flex;
  align-items: center;
}

.page-subtitle {
  color: #606266;
  font-size: 15px;
  margin: 0;
}

.header-actions {
  display: flex;
  gap: 12px;
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

.filter-card {
  margin-bottom: 24px;
  background: #ffffff;
  border-radius: 10px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.05);
  transition: all 0.3s ease;
  
  &:hover {
    box-shadow: 0 6px 20px rgba(0, 0, 0, 0.08);
  }
  
  .search-filters {
    display: flex;
    flex-direction: column;
    gap: 16px;
    
    .search-box {
      flex: 1;
      
      .el-input {
        transition: all 0.3s ease;
        
        &:focus-within {
          transform: translateY(-2px);
        }
      }
    }
    
    .filter-options {
      display: flex;
      flex-wrap: wrap;
      gap: 12px;
      align-items: center;
    }
    
    .el-button {
      transition: all 0.3s ease;
      
      &:hover {
        transform: translateY(-2px);
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
      }
    }
  }
}

.vulndb-table-card {
  background: #fff;
  border-radius: 10px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.05);
  transition: all 0.3s ease;
  overflow: hidden;
  
  &:hover {
    box-shadow: 0 6px 20px rgba(0, 0, 0, 0.08);
  }
  
  .table-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
    padding-bottom: 12px;
    border-bottom: 1px solid #ebeef5;
    
    .table-title {
      font-size: 18px;
      font-weight: 600;
      color: #303133;
      position: relative;
      
      &::after {
        content: '';
        position: absolute;
        bottom: -12px;
        left: 0;
        width: 40px;
        height: 2px;
        background-color: #409EFF;
      }
    }
    
    .table-actions {
      display: flex;
      gap: 12px;
    }
  }
  
  .vulndb-table {
    margin-bottom: 20px;
    
    :deep(.el-table__row) {
      transition: all 0.2s ease;
      
      &:hover {
        background-color: #f0f7ff !important;
        transform: translateY(-1px);
      }
    }
    
    .vuln-title {
      display: flex;
      align-items: center;
      gap: 8px;
      
      .title-text {
        font-weight: 500;
        transition: all 0.3s;
      }
      
      .tag-item {
        font-size: 11px;
        transition: all 0.3s;
        
        &:hover {
          transform: scale(1.05);
        }
      }
    }
    
    .cvss-score {
      display: flex;
      align-items: center;
      gap: 8px;
      
      .cvss-progress {
        width: 60px;
        transition: all 0.3s;
        
        &:hover {
          transform: scaleX(1.1);
        }
      }
      
      .cvss-value {
        font-weight: 600;
        font-size: 13px;
        background: #f5f7fa;
        padding: 2px 6px;
        border-radius: 10px;
        transition: all 0.3s;
        
        &:hover {
          background: #e6ebf5;
        }
      }
    }
    
    :deep(.el-button) {
      transition: all 0.3s ease;
      
      &:hover {
        transform: translateY(-2px);
        box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
      }
    }
    
    :deep(.el-tag) {
      transition: all 0.3s ease;
      
      &:hover {
        transform: scale(1.05);
      }
    }
    
    :deep(.el-link) {
      transition: all 0.3s;
      font-family: 'Courier New', monospace;
      
      &:hover {
        transform: translateY(-1px);
      }
    }
  }
}

.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 20px;
  padding: 16px;
  border-top: 1px solid #ebeef5;
  
  :deep(.el-pagination) {
    transition: all 0.3s;
    
    button {
      transition: all 0.3s;
      
      &:hover {
        transform: translateY(-1px);
      }
    }
  }
}

/* 响应式调整 */
@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    align-items: flex-start;
    
    .header-actions {
      margin-top: 15px;
      width: 100%;
      justify-content: flex-start;
      
      .el-button {
        flex-grow: 1;
      }
    }
  }
  
  .filter-options {
    flex-direction: column;
    align-items: stretch !important;
    width: 100%;
    
    .el-select {
      width: 100% !important;
      margin-bottom: 8px;
    }
    
    .el-button {
      width: 100%;
      margin-bottom: 8px;
    }
  }
  
  .table-header {
    flex-direction: column;
    align-items: flex-start;
    
    .table-actions {
      margin-top: 10px;
      width: 100%;
    }
  }
  
  .vulndb-table {
    :deep(.el-button) {
      padding: 7px 10px;
      font-size: 12px;
    }
  }
}
</style> 