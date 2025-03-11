<template>
  <div class="analysis-history">
    <div v-if="loading" class="loading-container">
      <div class="loading-text">正在加载分析历史...</div>
      <el-skeleton :rows="5" animated></el-skeleton>
    </div>
    
    <div v-else>
      <!-- 过滤器 -->
      <div class="history-filters">
        <el-form :inline="true" size="small">
          <el-form-item label="时间范围">
            <el-date-picker
              v-model="dateRange"
              type="daterange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              format="yyyy-MM-dd"
              value-format="yyyy-MM-dd"
              @change="handleFilterChange">
            </el-date-picker>
          </el-form-item>
          
          <el-form-item label="置信度">
            <el-select 
              v-model="confidenceFilter" 
              placeholder="选择置信度"
              @change="handleFilterChange">
              <el-option label="全部" value=""></el-option>
              <el-option label="高置信度 (>80%)" value="high"></el-option>
              <el-option label="中置信度 (60-80%)" value="medium"></el-option>
              <el-option label="低置信度 (<60%)" value="low"></el-option>
            </el-select>
          </el-form-item>
          
          <el-form-item>
            <el-button type="primary" icon="el-icon-search" @click="handleFilterChange">筛选</el-button>
            <el-button icon="el-icon-refresh" @click="resetFilters">重置</el-button>
          </el-form-item>
        </el-form>
      </div>
      
      <!-- 分析列表 -->
      <div class="history-list">
        <el-table
          v-loading="loading"
          :data="filteredAnalyses"
          style="width: 100%"
          border
          stripe>
          <el-table-column prop="title" label="分析标题" min-width="200">
            <template v-slot="scope">
              <el-link type="primary" @click="viewAnalysisDetail(scope.row.id)">
                {{ scope.row.title }}
              </el-link>
            </template>
          </el-table-column>
          <el-table-column label="分析类型" width="150">
            <template v-slot="scope">
              <el-tag :type="getAnalysisTypeTag(scope.row.type)" size="medium">
                {{ getAnalysisTypeLabel(scope.row.type) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="创建时间" width="180">
            <template v-slot="scope">
              {{ formatDate(scope.row.createdAt) }}
            </template>
          </el-table-column>
          <el-table-column label="置信度" width="150">
            <template v-slot="scope">
              <el-progress 
                :percentage="scope.row.confidence * 100" 
                :color="getConfidenceColor(scope.row.confidence)"
                :format="percentFormat">
              </el-progress>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="150" fixed="right">
            <template v-slot="scope">
              <el-button 
                type="text" 
                @click="viewAnalysisDetail(scope.row.id)">
                查看详情
              </el-button>
              <el-button 
                type="text" 
                @click="deleteAnalysis(scope.row.id)">
                删除
              </el-button>
            </template>
          </el-table-column>
        </el-table>
        
        <!-- 分页器 -->
        <div class="pagination-container">
          <el-pagination
            background
            layout="total, sizes, prev, pager, next, jumper"
            v-model:current-page="currentPage"
            v-model:page-size="pageSize"
            :total="total"
            @current-change="handlePageChange"
            @size-change="handleSizeChange">
          </el-pagination>
        </div>
      </div>
      
      <!-- 空状态 -->
      <el-empty
        v-if="filteredAnalyses.length === 0"
        description="暂无符合条件的分析记录">
      </el-empty>
    </div>
  </div>
</template>

<script>
import { formatDate } from '@/utils/format'

export default {
  name: 'AnalysisHistory',
  props: {
    analysisType: {
      type: String,
      default: ''
    }
  },
  data() {
    return {
      loading: false,
      analyses: [],
      dateRange: null,
      confidenceFilter: '',
      currentPage: 1,
      pageSize: 10,
      total: 0
    }
  },
  computed: {
    filteredAnalyses() {
      let result = [...this.analyses]
      
      // 筛选分析类型
      if (this.analysisType) {
        result = result.filter(item => item.type === this.analysisType)
      }
      
      // 筛选日期范围
      if (this.dateRange && this.dateRange.length === 2) {
        const startDate = new Date(this.dateRange[0])
        const endDate = new Date(this.dateRange[1])
        endDate.setHours(23, 59, 59, 999) // 设置为当天结束时间
        
        result = result.filter(item => {
          const itemDate = new Date(item.createdAt)
          return itemDate >= startDate && itemDate <= endDate
        })
      }
      
      // 筛选置信度
      if (this.confidenceFilter) {
        switch (this.confidenceFilter) {
          case 'high':
            result = result.filter(item => item.confidence > 0.8)
            break
          case 'medium':
            result = result.filter(item => item.confidence >= 0.6 && item.confidence <= 0.8)
            break
          case 'low':
            result = result.filter(item => item.confidence < 0.6)
            break
        }
      }
      
      return result
    }
  },
  created() {
    this.fetchAnalysisHistory()
  },
  methods: {
    formatDate,
    async fetchAnalysisHistory() {
      this.loading = true
      
      try {
        // 实际项目中应调用API获取历史记录
        // const { data } = await this.$api.ai.getAnalysisHistory({
        //   type: this.analysisType,
        //   page: this.currentPage,
        //   pageSize: this.pageSize
        // })
        // this.analyses = data.analyses
        // this.total = data.total
        
        // 模拟数据
        await new Promise(resolve => setTimeout(resolve, 1000))
        
        const mockData = [
          {
            id: '61e8f5c3b32f6c001a9b3d01',
            title: '未来30天漏洞趋势预测',
            type: 'trend_prediction',
            confidence: 0.85,
            createdAt: new Date(Date.now() - 86400000 * 2)
          },
          {
            id: '61e8f5c3b32f6c001a9b3d02',
            title: '安全资源优化配置建议',
            type: 'resource_optimization',
            confidence: 0.78,
            createdAt: new Date(Date.now() - 86400000 * 5)
          },
          {
            id: '61e8f5c3b32f6c001a9b3d03',
            title: '过去7天漏洞异常检测',
            type: 'anomaly_detection',
            confidence: 0.82,
            createdAt: new Date(Date.now() - 86400000 * 1)
          },
          {
            id: '61e8f5c3b32f6c001a9b3d04',
            title: '未来14天漏洞趋势预测',
            type: 'trend_prediction',
            confidence: 0.72,
            createdAt: new Date(Date.now() - 86400000 * 10)
          },
          {
            id: '61e8f5c3b32f6c001a9b3d05',
            title: '安全资源季度优化报告',
            type: 'resource_optimization',
            confidence: 0.91,
            createdAt: new Date(Date.now() - 86400000 * 15)
          }
        ]
        
        this.analyses = mockData
        this.total = mockData.length
      } catch (error) {
        this.$message.error('获取分析历史记录失败: ' + error.message)
      } finally {
        this.loading = false
      }
    },
    handleFilterChange() {
      this.currentPage = 1
      this.fetchAnalysisHistory()
    },
    resetFilters() {
      this.dateRange = null
      this.confidenceFilter = ''
      this.currentPage = 1
      this.fetchAnalysisHistory()
    },
    handlePageChange(page) {
      this.currentPage = page
      this.fetchAnalysisHistory()
    },
    handleSizeChange(size) {
      this.pageSize = size
      this.currentPage = 1
      this.fetchAnalysisHistory()
    },
    viewAnalysisDetail(id) {
      this.$emit('view-detail', id)
    },
    async deleteAnalysis(id) {
      try {
        await this.$confirm('确定要删除这个分析结果吗？', '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        })
        
        // 调用API删除分析结果
        // await this.$api.ai.deleteAnalysis(id)
        
        this.$message.success('删除成功')
        this.analyses = this.analyses.filter(item => item.id !== id)
      } catch (error) {
        if (error !== 'cancel') {
          this.$message.error('删除失败: ' + error.message)
        }
      }
    },
    getAnalysisTypeLabel(type) {
      const typeMap = {
        'trend_prediction': '漏洞趋势预测',
        'resource_optimization': '安全资源优化',
        'anomaly_detection': '异常检测'
      }
      return typeMap[type] || '未知类型'
    },
    getAnalysisTypeTag(type) {
      const typeMap = {
        'trend_prediction': 'primary',
        'resource_optimization': 'success',
        'anomaly_detection': 'warning'
      }
      return typeMap[type] || 'info'
    },
    getConfidenceColor(confidence) {
      if (confidence >= 0.8) return '#67C23A'
      if (confidence >= 0.6) return '#E6A23C'
      return '#F56C6C'
    },
    percentFormat(percentage) {
      return percentage.toFixed(0) + '%'
    }
  }
}
</script>

<style lang="scss" scoped>
.analysis-history {
  .loading-container {
    padding: 20px;
    
    .loading-text {
      margin-bottom: 15px;
      color: #909399;
    }
  }
  
  .history-filters {
    margin-bottom: 20px;
  }
  
  .history-list {
    margin-bottom: 20px;
  }
  
  .pagination-container {
    margin-top: 20px;
    display: flex;
    justify-content: center;
  }
}
</style> 