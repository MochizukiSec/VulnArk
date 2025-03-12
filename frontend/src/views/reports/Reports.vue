<template>
  <div class="reports-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">
          <span class="title-icon"><i class="el-icon-document"></i></span>
          安全报告
          <span class="title-highlight">管理</span>
        </h1>
        <p class="page-subtitle">查看、生成和导出系统安全漏洞相关报告</p>
      </div>
    </div>

    <!-- 报告生成卡片 -->
    <div class="report-card-container">
      <el-card shadow="hover" class="report-card">
        <template #header>
          <div class="card-header">
            <h2 class="section-title">
              <i class="el-icon-document-add"></i>
              报告生成
            </h2>
          </div>
        </template>

        <div class="report-options">
          <el-form :model="reportForm" label-position="top">
            <el-row :gutter="24">
              <el-col :xs="24" :sm="12">
                <el-form-item label="报告类型">
                  <el-select v-model="reportForm.type" placeholder="选择报告类型" style="width: 100%">
                    <el-option-group label="常规报告">
                      <el-option label="漏洞摘要报告" value="summary">
                        <template #default>
                          <div class="report-option">
                            <i class="el-icon-document option-icon summary"></i>
                            <span>漏洞摘要报告</span>
                          </div>
                        </template>
                      </el-option>
                      <el-option label="详细漏洞报告" value="detailed">
                        <template #default>
                          <div class="report-option">
                            <i class="el-icon-document-checked option-icon detailed"></i>
                            <span>详细漏洞报告</span>
                          </div>
                        </template>
                      </el-option>
                    </el-option-group>
                    <el-option-group label="专项报告">
                      <el-option label="合规性报告" value="compliance">
                        <template #default>
                          <div class="report-option">
                            <i class="el-icon-collection-check option-icon compliance"></i>
                            <span>合规性报告</span>
                          </div>
                        </template>
                      </el-option>
                      <el-option label="趋势分析报告" value="trend">
                        <template #default>
                          <div class="report-option">
                            <i class="el-icon-data-line option-icon trend"></i>
                            <span>趋势分析报告</span>
                          </div>
                        </template>
                      </el-option>
                    </el-option-group>
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :xs="24" :sm="12">
                <el-form-item label="报告格式">
                  <el-select v-model="reportForm.format" placeholder="选择输出格式" style="width: 100%">
                    <el-option label="PDF文档" value="pdf">
                      <template #default>
                        <div class="report-option">
                          <i class="el-icon-document option-icon pdf"></i>
                          <span>PDF文档</span>
                        </div>
                      </template>
                    </el-option>
                    <el-option label="Excel工作表" value="excel">
                      <template #default>
                        <div class="report-option">
                          <i class="el-icon-document option-icon excel"></i>
                          <span>Excel工作表</span>
                        </div>
                      </template>
                    </el-option>
                    <el-option label="Word文档" value="word">
                      <template #default>
                        <div class="report-option">
                          <i class="el-icon-document option-icon word"></i>
                          <span>Word文档</span>
                        </div>
                      </template>
                    </el-option>
                    <el-option label="HTML网页" value="html">
                      <template #default>
                        <div class="report-option">
                          <i class="el-icon-document option-icon html"></i>
                          <span>HTML网页</span>
                        </div>
                      </template>
                    </el-option>
                    <el-option label="文本文件" value="text">
                      <template #default>
                        <div class="report-option">
                          <i class="el-icon-document option-icon text"></i>
                          <span>文本文件</span>
                        </div>
                      </template>
                    </el-option>
                  </el-select>
                </el-form-item>
              </el-col>
            </el-row>

            <el-form-item label="报告时间范围">
              <el-date-picker
                v-model="reportForm.dateRange"
                type="daterange"
                range-separator="至"
                start-placeholder="开始日期"
                end-placeholder="结束日期"
                style="width: 100%"
                class="date-range-picker"
              />
            </el-form-item>

            <el-form-item label="包含的严重程度">
              <div class="severity-checkboxes">
                <el-checkbox-group v-model="reportForm.severities">
                  <el-checkbox label="critical">
                    <span class="severity-checkbox">
                      <span class="severity-dot critical"></span>
                      <span>严重</span>
                    </span>
                  </el-checkbox>
                  <el-checkbox label="high">
                    <span class="severity-checkbox">
                      <span class="severity-dot high"></span>
                      <span>高危</span>
                    </span>
                  </el-checkbox>
                  <el-checkbox label="medium">
                    <span class="severity-checkbox">
                      <span class="severity-dot medium"></span>
                      <span>中危</span>
                    </span>
                  </el-checkbox>
                  <el-checkbox label="low">
                    <span class="severity-checkbox">
                      <span class="severity-dot low"></span>
                      <span>低危</span>
                    </span>
                  </el-checkbox>
                  <el-checkbox label="info">
                    <span class="severity-checkbox">
                      <span class="severity-dot info"></span>
                      <span>信息</span>
                    </span>
                  </el-checkbox>
                </el-checkbox-group>
              </div>
            </el-form-item>

            <el-form-item>
              <el-button type="primary" @click="generateReport" :loading="generating" class="generate-btn">
                <i class="el-icon-document-add"></i> 生成报告
              </el-button>
            </el-form-item>
          </el-form>
        </div>
      </el-card>
    </div>

    <!-- 报告列表 -->
    <div class="reports-list-container">
      <el-card shadow="hover" class="reports-list-card">
        <template #header>
          <div class="card-header">
            <h2 class="section-title">
              <i class="el-icon-document-copy"></i>
              报告历史
            </h2>
            
            <!-- 筛选工具栏 -->
            <div class="filter-toolbar">
              <el-input
                v-model="filterForm.search"
                placeholder="搜索报告名称"
                clearable
                prefix-icon="el-icon-search"
                class="search-input"
                @keyup.enter="searchReports"
              />
              
              <el-select v-model="filterForm.type" placeholder="报告类型" clearable @change="searchReports">
                <el-option label="摘要报告" value="summary" />
                <el-option label="详细报告" value="detailed" />
                <el-option label="合规报告" value="compliance" />
                <el-option label="趋势分析" value="trend" />
              </el-select>
              
              <el-select v-model="filterForm.format" placeholder="文件格式" clearable @change="searchReports">
                <el-option label="PDF" value="pdf" />
                <el-option label="Excel" value="excel" />
                <el-option label="Word" value="word" />
                <el-option label="HTML" value="html" />
              </el-select>
              
              <el-button type="primary" @click="searchReports">
                <i class="el-icon-search"></i>
                搜索
              </el-button>
              
              <el-button @click="resetFilter">
                <i class="el-icon-refresh"></i>
                重置
              </el-button>
            </div>
          </div>
        </template>
        
        <div v-loading="loading" class="reports-list">
          <!-- 空状态 -->
          <el-empty
            v-if="!loading && reports.length === 0"
            description="暂无报告记录"
            class="empty-state"
          >
            <el-button type="primary" @click="scrollToGenerator">
              创建第一份报告
            </el-button>
          </el-empty>
          
          <!-- 报告列表 -->
          <div v-else class="reports-grid">
            <div v-for="report in reports" :key="report.id" class="report-card-item">
              <div class="report-header">
                <el-tag :type="getReportTagType(report.type)" effect="light" class="report-type">
                  {{ getReportTypeText(report.type) }}
                </el-tag>
                <div class="report-format">
                  {{ report.format ? report.format.toUpperCase() : '' }}
                </div>
              </div>
              
              <h3 class="report-name">{{ report.name }}</h3>
              
              <div class="report-date">
                <i class="el-icon-date"></i>
                <span>{{ formatDate(report.created_at) }}</span>
              </div>
              
              <div class="report-actions">
                <el-button type="primary" size="small" :disabled="!report.file_url" @click="downloadReport(report)">
                  <i class="el-icon-download"></i> 下载
                </el-button>
                <el-button type="danger" size="small" @click="confirmDelete(report)">
                  <i class="el-icon-delete"></i> 删除
                </el-button>
              </div>
              
              <div v-if="report.status === 'pending'" class="report-status pending">
                <i class="el-icon-loading"></i> 生成中...
              </div>
              <div v-else-if="report.status === 'failed'" class="report-status failed">
                <i class="el-icon-warning-outline"></i> 生成失败
              </div>
            </div>
          </div>
          
          <!-- 分页 -->
          <div class="pagination-container">
            <el-pagination
              :current-page="pagination.page"
              :page-size="pagination.pageSize"
              :page-sizes="[10, 20, 50]"
              layout="total, sizes, prev, pager, next, jumper"
              :total="pagination.total"
              @size-change="handleSizeChange"
              @current-change="handlePageChange"
              background
            />
          </div>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script>
import { reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useStore } from 'vuex'

export default {
  name: 'Reports',
  
  setup() {
    const store = useStore()
    
    // 加载和生成状态
    const loading = computed(() => store.getters['report/isLoading'])
    const generating = computed(() => store.getters['report/isGenerating'])
    
    // 分页信息
    const pagination = reactive({
      page: 1,
      pageSize: 10,
      total: 0
    })
    
    // 筛选条件
    const filterForm = reactive({
      type: '',
      format: '',
      search: ''
    })
    
    // 报告生成表单
    const reportForm = reactive({
      name: '',
      type: 'summary',
      format: 'pdf',
      dateRange: [],
      severities: ['critical', 'high', 'medium'],
      description: ''
    })
    
    // 从store获取报告数据
    const reports = computed(() => store.getters['report/allReports'])
    
    // 加载报告列表
    const loadReports = async () => {
      try {
        await store.dispatch('report/fetchReports', {
          page: pagination.page,
          limit: pagination.pageSize,
          type: filterForm.type,
          format: filterForm.format,
          search: filterForm.search
        })
        
        // 更新分页数据
        const storePagination = store.getters['report/pagination']
        pagination.total = storePagination.total
      } catch (error) {
        console.error('加载报告失败:', error)
      }
    }
    
    // 在组件挂载时加载报告
    onMounted(() => {
      loadReports()
    })
    
    // 搜索报告
    const searchReports = () => {
      pagination.page = 1 // 重置到第一页
      loadReports()
    }
    
    // 重置筛选条件
    const resetFilter = () => {
      Object.assign(filterForm, {
        type: '',
        format: '',
        search: ''
      })
      searchReports()
    }
    
    // 处理分页变化
    const handlePageChange = (newPage) => {
      pagination.page = newPage
      loadReports()
    }
    
    // 处理每页显示数量变化
    const handleSizeChange = (newSize) => {
      pagination.pageSize = newSize
      pagination.page = 1 // 重置到第一页
      loadReports()
    }
    
    // 获取报告类型文本
    const getReportTypeText = (type) => {
      const typeMap = {
        summary: '摘要报告',
        detailed: '详细报告',
        compliance: '合规报告',
        trend: '趋势分析'
      }
      return typeMap[type] || type
    }
    
    // 获取报告标签类型
    const getReportTagType = (type) => {
      const typeMap = {
        summary: '',
        detailed: 'success',
        compliance: 'warning',
        trend: 'info'
      }
      return typeMap[type] || ''
    }
    
    // 格式化日期
    const formatDate = (dateString) => {
      if (!dateString) return '无日期'
      
      try {
        const date = new Date(dateString)
        if (isNaN(date.getTime())) return '无效日期'
        
        return date.toLocaleDateString('zh-CN', {
          year: 'numeric',
          month: '2-digit',
          day: '2-digit',
          hour: '2-digit',
          minute: '2-digit'
        })
      } catch (error) {
        console.error('日期格式化错误:', error)
        return '无效日期'
      }
    }
    
    // 滚动到报告生成器
    const scrollToGenerator = () => {
      document.querySelector('.report-card').scrollIntoView({ 
        behavior: 'smooth',
        block: 'start'
      })
    }
    
    // 生成报告
    const generateReport = async () => {
      if (!reportForm.dateRange || reportForm.dateRange.length !== 2) {
        ElMessage.warning('请选择报告时间范围')
        return
      }
      
      if (reportForm.severities.length === 0) {
        ElMessage.warning('请至少选择一个严重程度级别')
        return
      }
      
      // 检查用户是否已登录并尝试刷新token状态
      if (!store.getters['auth/isAuthenticated']) {
        // 尝试刷新token状态
        const refreshed = await store.dispatch('auth/refreshTokenState')
        
        if (!refreshed) {
          ElMessage.error('您需要登录才能生成报告')
          // 可以选择重定向到登录页面
          return
        }
      }
      
      if (!reportForm.name) {
        // 自动生成报告名称
        const reportTypes = {
          summary: '漏洞摘要报告',
          detailed: '详细漏洞报告',
          compliance: '合规性报告',
          trend: '趋势分析报告'
        }
        reportForm.name = `${new Date().getFullYear()}年${reportTypes[reportForm.type]}`
      }
      
      try {
        // 直接从localStorage获取token进行检查
        const token = localStorage.getItem('token')
        if (!token) {
          ElMessage.error('未找到授权令牌，请重新登录')
          return
        }
        
        // 准备请求数据
        const reportData = {
          name: reportForm.name,
          type: reportForm.type,
          format: reportForm.format,
          start_date: reportForm.dateRange[0],
          end_date: reportForm.dateRange[1],
          severities: reportForm.severities,
          description: reportForm.description
        }
        
        // 创建报告
        await store.dispatch('report/createReport', reportData)
        
        // 显示成功消息
        ElMessage.success({
          message: '报告生成请求已提交，后端功能待完善',
          duration: 5000
        })
        
        // 重置表单
        reportForm.name = ''
        reportForm.description = ''
        
        // 刷新报告列表
        loadReports()
      } catch (error) {
        console.error('生成报告失败:', error)
        
        // 检查是否是授权错误
        if (error.response && error.response.status === 401) {
          ElMessage.error({
            message: '身份验证失败，请重新登录后再试',
            duration: 5000
          })
          // 可以选择重定向到登录页面
          // router.push('/login')
        } else {
          ElMessage.error({
            message: '报告生成失败: ' + (error.response?.data?.error || '后端功能未完全实现，请等待功能完善'),
            duration: 5000
          })
        }
      }
    }
    
    // 下载报告
    const downloadReport = (report) => {
      if (report.file_url) {
        try {
          // 获取完整的URL
          let downloadUrl = report.file_url;
          
          // 如果不是绝对URL，添加API基础路径
          if (!downloadUrl.startsWith('http') && !downloadUrl.startsWith('/api')) {
            downloadUrl = '/api' + downloadUrl;
          }
          
          // 打开链接
          window.open(downloadUrl, '_blank');
        } catch (error) {
          console.error('下载报告失败:', error);
          ElMessage.error('下载报告失败，请稍后再试');
        }
      } else {
        ElMessage.warning('报告文件尚未生成完成，请稍后再试');
      }
    }
    
    // 确认删除
    const confirmDelete = (report) => {
      ElMessageBox.confirm(
        `确定要删除报告"${report.name}"吗？`,
        '警告',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }
      ).then(async () => {
        try {
          await store.dispatch('report/deleteReport', report.id)
          loadReports() // 刷新列表
        } catch (error) {
          console.error('删除报告失败:', error)
        }
      }).catch(() => {
        // 用户取消删除操作
      })
    }
    
    // 生成摘要报告
    const generateSummaryReport = async (startDate, endDate) => {
      try {
        await store.dispatch('report/generateSummaryReport', { 
          startDate: startDate?.toISOString(), 
          endDate: endDate?.toISOString() 
        })
      } catch (error) {
        console.error('生成摘要报告失败:', error)
      }
    }
    
    // 生成详细报告
    const generateDetailedReport = async (startDate, endDate, severity, status) => {
      try {
        await store.dispatch('report/generateDetailedReport', { 
          startDate: startDate?.toISOString(), 
          endDate: endDate?.toISOString(),
          severity,
          status
        })
      } catch (error) {
        console.error('生成详细报告失败:', error)
      }
    }
    
    return {
      loading,
      generating,
      pagination,
      filterForm,
      reportForm,
      reports,
      getReportTypeText,
      getReportTagType,
      formatDate,
      scrollToGenerator,
      generateReport,
      downloadReport,
      confirmDelete,
      searchReports,
      resetFilter,
      handlePageChange,
      handleSizeChange,
      generateSummaryReport,
      generateDetailedReport
    }
  }
}
</script>

<style lang="scss" scoped>
.reports-page {
  padding: 24px;
  animation: fadeIn 0.6s ease;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  
  .header-content {
    flex: 1;
  }
  
  .page-title {
    font-size: 28px;
    font-weight: 600;
    margin: 0 0 8px 0;
    display: flex;
    align-items: center;
    color: #303133;
    
    .title-icon {
      margin-right: 12px;
      font-size: 28px;
      color: #409EFF;
    }
    
    .title-highlight {
      background: linear-gradient(120deg, #409EFF, #53A8FF);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      margin-left: 8px;
      font-weight: 700;
    }
  }
  
  .page-subtitle {
    color: #606266;
    font-size: 16px;
    margin: 0;
  }
}

.report-card-container,
.reports-list-container {
  margin-bottom: 24px;
}

.report-card,
.reports-list-card {
  border-radius: 12px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.05);
  transition: all 0.3s;
  position: relative;
  overflow: hidden;
  
  &:hover {
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.09);
  }
  
  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 4px;
    background: linear-gradient(90deg, #409EFF, #53A8FF);
  }
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 4px 0;
}

.section-title {
  font-size: 18px;
  font-weight: 600;
  margin: 0;
  color: #303133;
  display: flex;
  align-items: center;
  
  i {
    margin-right: 8px;
    font-size: 18px;
    color: #409EFF;
  }
}

.report-options {
  padding: 16px 0;
}

.report-option {
  display: flex;
  align-items: center;
  
  .option-icon {
    margin-right: 8px;
    font-size: 16px;
    width: 20px;
    height: 20px;
    text-align: center;
    line-height: 20px;
    
    &.summary {
      color: #409EFF;
    }
    
    &.detailed {
      color: #67C23A;
    }
    
    &.compliance {
      color: #E6A23C;
    }
    
    &.trend {
      color: #9c27b0;
    }
    
    &.pdf {
      color: #F56C6C;
    }
    
    &.excel {
      color: #67C23A;
    }
    
    &.word {
      color: #409EFF;
    }
    
    &.html {
      color: #E6A23C;
    }
    
    &.text {
      color: #909399;
    }
  }
}

.date-range-picker {
  :deep(.el-input__inner) {
    border-radius: 8px;
  }
}

.severity-checkboxes {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  
  .severity-checkbox {
    display: flex;
    align-items: center;
    
    .severity-dot {
      width: 10px;
      height: 10px;
      border-radius: 50%;
      margin-right: 8px;
      
      &.critical {
        background-color: #F44336;
      }
      
      &.high {
        background-color: #FF9800;
      }
      
      &.medium {
        background-color: #FFEB3B;
      }
      
      &.low {
        background-color: #4CAF50;
      }
      
      &.info {
        background-color: #2196F3;
      }
    }
  }
}

.generate-btn {
  border-radius: 8px;
  padding: 12px 24px;
  font-weight: 500;
  transition: all 0.3s;
  
  i {
    margin-right: 6px;
  }
  
  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(64, 158, 255, 0.2);
  }
}

.reports-list-container {
  .filter-toolbar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
    
    .search-input {
      width: 200px;
    }
  }
}

.reports-list {
  .empty-state {
    padding: 40px 0;
    text-align: center;
    
    p {
      color: #909399;
      margin: 8px 0 16px;
    }
  }
  
  .reports-grid {
    display: flex;
    flex-wrap: wrap;
    gap: 16px;
    
    .report-card-item {
      width: calc(33.33% - 16px);
      padding: 16px;
      border-radius: 8px;
      box-shadow: 0 4px 16px rgba(0, 0, 0, 0.05);
      transition: all 0.3s;
      
      &:hover {
        box-shadow: 0 8px 24px rgba(0, 0, 0, 0.09);
      }
      
      .report-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 16px;
        
        .report-type {
          margin-right: 8px;
        }
        
        .report-format {
          padding: 2px 8px;
          border-radius: 4px;
          font-size: 12px;
          font-weight: 500;
          text-transform: uppercase;
          background-color: #f0f0f0;
        }
      }
      
      .report-name {
        font-size: 18px;
        font-weight: 600;
        margin: 8px 0;
      }
      
      .report-date {
        display: flex;
        align-items: center;
        margin-top: 8px;
        margin-bottom: 8px;
        
        .el-icon-date {
          margin-right: 6px;
          color: #909399;
        }
      }
      
      .report-actions {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-top: 16px;
        
        .action-btn {
          padding: 4px 8px;
          border-radius: 4px;
          transition: all 0.3s;
          
          i {
            margin-right: 4px;
          }
          
          &:hover {
            opacity: 0.9;
            transform: translateY(-1px);
          }
          
          &.download-btn {
            background-color: #409EFF;
            border-color: #409EFF;
          }
          
          &.delete-btn {
            background-color: #F56C6C;
            border-color: #F56C6C;
          }
        }
      }
      
      .report-status {
        margin-top: 8px;
        padding: 4px 8px;
        border-radius: 4px;
        font-size: 12px;
        font-weight: 500;
        
        &.pending {
          background-color: #f0f9eb;
          color: #67C23A;
        }
        
        &.failed {
          background-color: #fef2f2;
          color: #F56C6C;
        }
      }
    }
  }
  
  .pagination-container {
    display: flex;
    justify-content: center;
    margin-top: 16px;
  }
}

// 响应式调整
@media (max-width: 768px) {
  .reports-page {
    padding: 16px;
  }
  
  .page-header {
    .page-title {
      font-size: 24px;
      
      .title-icon {
        font-size: 24px;
      }
    }
    
    .page-subtitle {
      font-size: 14px;
    }
  }
  
  .severity-checkboxes {
    gap: 12px;
  }
  
  .reports-list-container {
    .filter-toolbar {
      flex-direction: column;
      gap: 8px;
      
      .search-input {
        width: 100%;
      }
    }
  }
}
</style> 