<template>
  <div class="ai-analysis-container">
    <h1 class="page-title">
      <i class="el-icon-cpu"></i> AI 智能分析
    </h1>
    
    <el-card class="intro-card">
      <div class="intro-content">
        <div class="intro-text">
          <h2>安全数据智能分析</h2>
          <p>利用人工智能技术对漏洞数据进行深度分析，提供趋势预测、资源优化建议和异常检测。帮助您更高效地管理安全资源，提前发现潜在威胁。</p>
          <div class="analysis-metrics">
            <div class="metric">
              <div class="metric-value">{{ analysisCount }}</div>
              <div class="metric-label">分析报告</div>
            </div>
            <div class="metric">
              <div class="metric-value">{{ averageConfidence.toFixed(0) }}%</div>
              <div class="metric-label">平均置信度</div>
            </div>
            <div class="metric">
              <div class="metric-value">{{ improvementRate.toFixed(0) }}%</div>
              <div class="metric-label">效率提升</div>
            </div>
          </div>
        </div>
        <div class="intro-image">
          <img src="@/assets/images/ai-analysis.svg" alt="AI分析" />
        </div>
      </div>
    </el-card>

    <div class="analysis-types">
      <el-row :gutter="20">
        <el-col :xs="24" :sm="24" :md="8" :lg="8" :xl="8" v-for="(analysisType, index) in analysisTypes" :key="index">
          <el-card class="analysis-card" :body-style="{ padding: '0px' }" shadow="hover">
            <div class="analysis-card-header" :class="analysisType.colorClass">
              <i :class="analysisType.icon"></i>
            </div>
            <div class="analysis-card-content">
              <h3>{{ analysisType.title }}</h3>
              <p>{{ analysisType.description }}</p>
              <div class="card-actions">
                <el-button type="primary" size="small" @click="startAnalysis(analysisType.type)">
                  开始分析
                </el-button>
                <el-button size="small" @click="viewHistory(analysisType.type)">
                  历史结果
                </el-button>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <div class="recent-analyses" v-if="recentAnalyses.length > 0">
      <h2 class="section-title">最近分析结果</h2>
      <el-table :data="recentAnalyses" style="width: 100%" border stripe>
        <el-table-column prop="title" label="分析标题" min-width="200">
          <template v-slot="scope">
            <el-link type="primary" @click="viewAnalysisDetail(scope.row.id)">
              {{ scope.row.title }}
            </el-link>
          </template>
        </el-table-column>
        <el-table-column prop="type" label="分析类型" width="150">
          <template v-slot="scope">
            <el-tag :type="getAnalysisTypeTag(scope.row.type)">
              {{ getAnalysisTypeLabel(scope.row.type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="confidence" label="置信度" width="100">
          <template v-slot="scope">
            <el-progress :percentage="scope.row.confidence * 100" :color="getConfidenceColor(scope.row.confidence)"></el-progress>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" width="180">
          <template v-slot="scope">
            {{ formatDate(scope.row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template v-slot="scope">
            <el-button type="text" @click="viewAnalysisDetail(scope.row.id)">查看详情</el-button>
            <el-button type="text" @click="deleteAnalysis(scope.row.id)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- 开始分析对话框 -->
    <el-dialog 
      :title="currentAnalysisType ? `开始${getAnalysisTypeLabel(currentAnalysisType)}分析` : '开始分析'" 
      v-model="analysisDialogVisible"
      width="600px">
      <analysis-parameters 
        v-if="analysisDialogVisible" 
        :analysis-type="currentAnalysisType"
        @start="executeAnalysis"
        @cancel="analysisDialogVisible = false">
      </analysis-parameters>
    </el-dialog>

    <!-- 分析结果详情对话框 -->
    <el-dialog 
      title="分析结果详情" 
      v-model="detailDialogVisible" 
      width="80%"
      fullscreen>
      <analysis-detail 
        v-if="detailDialogVisible && currentAnalysisId" 
        :analysis-id="currentAnalysisId"
        @close="detailDialogVisible = false">
      </analysis-detail>
    </el-dialog>

    <!-- 历史分析结果对话框 -->
    <el-dialog 
      :title="`${currentAnalysisTypeLabel || ''}分析历史`" 
      v-model="historyDialogVisible" 
      width="80%">
      <analysis-history 
        v-if="historyDialogVisible" 
        :analysis-type="currentAnalysisType"
        @view-detail="viewAnalysisDetail"
        @close="historyDialogVisible = false">
      </analysis-history>
    </el-dialog>
  </div>
</template>

<script>
import { formatDate } from '@/utils/format'
import AnalysisParameters from '@/components/ai/AnalysisParameters.vue'
import AnalysisDetail from '@/components/ai/AnalysisDetail.vue'
import AnalysisHistory from '@/components/ai/AnalysisHistory.vue'

export default {
  name: 'AIAnalysis',
  components: {
    AnalysisParameters,
    AnalysisDetail,
    AnalysisHistory
  },
  data() {
    return {
      analysisCount: 42,
      averageConfidence: 87.5,
      improvementRate: 35,
      analysisTypes: [
        {
          title: '漏洞趋势预测',
          description: '基于历史漏洞数据，预测未来安全趋势和漏洞发展模式，帮助提前做好防御准备。',
          icon: 'el-icon-data-line',
          type: 'trend_prediction',
          colorClass: 'blue-header'
        },
        {
          title: '安全资源优化',
          description: '分析当前安全资源分配状况，提供优化建议，提高安全投入的回报率。',
          icon: 'el-icon-s-operation',
          type: 'resource_optimization',
          colorClass: 'green-header'
        },
        {
          title: '异常检测分析',
          description: '识别漏洞数据中的异常模式和潜在安全事件，及早发现攻击迹象。',
          icon: 'el-icon-warning-outline',
          type: 'anomaly_detection',
          colorClass: 'orange-header'
        }
      ],
      recentAnalyses: [],
      loading: false,
      analysisDialogVisible: false,
      detailDialogVisible: false,
      historyDialogVisible: false,
      currentAnalysisType: null,
      currentAnalysisId: null,
      currentAnalysisTypeLabel: ''
    }
  },
  created() {
    this.fetchRecentAnalyses()
  },
  methods: {
    formatDate,
    async fetchRecentAnalyses() {
      this.loading = true
      try {
        // 从API获取最近的分析结果
        // const { data } = await this.$api.ai.getRecentAnalyses()
        // this.recentAnalyses = data.analyses
        
        // 模拟数据
        this.recentAnalyses = [
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
          }
        ]
      } catch (error) {
        this.$message.error('获取分析结果失败: ' + error.message)
      } finally {
        this.loading = false
      }
    },
    startAnalysis(type) {
      this.currentAnalysisType = type
      this.analysisDialogVisible = true
    },
    async executeAnalysis(parameters) {
      this.analysisDialogVisible = false
      this.loading = true
      
      try {
        // 调用API执行分析
        // const { data } = await this.$api.ai.runAnalysis({
        //   type: this.currentAnalysisType,
        //   parameters
        // })
        
        // 模拟分析过程
        await new Promise(resolve => setTimeout(resolve, 2000))
        
        this.$message.success('分析完成！')
        this.fetchRecentAnalyses()
        
        // 显示结果 - 使用parameters确定分析类型
        console.log('使用参数开始分析:', parameters);
        
        // 模拟ID
        const mockIdMap = {
          'trend_prediction': '61e8f5c3b32f6c001a9b3d01',
          'resource_optimization': '61e8f5c3b32f6c001a9b3d02',
          'anomaly_detection': '61e8f5c3b32f6c001a9b3d03'
        }
        
        this.currentAnalysisId = mockIdMap[this.currentAnalysisType]
        this.detailDialogVisible = true
      } catch (error) {
        this.$message.error('执行分析失败: ' + error.message)
      } finally {
        this.loading = false
      }
    },
    viewAnalysisDetail(id) {
      this.currentAnalysisId = id
      this.detailDialogVisible = true
    },
    viewHistory(type) {
      this.currentAnalysisType = type
      this.currentAnalysisTypeLabel = this.getAnalysisTypeLabel(type)
      this.historyDialogVisible = true
    },
    async deleteAnalysis(id) {
      try {
        await this.$confirm('确定要删除这个分析结果吗？', '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        })
        
        // 从API删除分析结果
        // await this.$api.ai.deleteAnalysis(id)
        
        this.$message.success('删除成功')
        this.recentAnalyses = this.recentAnalyses.filter(item => item.id !== id)
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
    }
  }
}
</script>

<style lang="scss" scoped>
.ai-analysis-container {
  padding: 24px;
  
  .page-title {
    margin-bottom: 24px;
    font-size: 24px;
    font-weight: 600;
    display: flex;
    align-items: center;
    color: #303133;
    
    i {
      margin-right: 12px;
      color: #409EFF;
      font-size: 28px;
    }
  }
  
  .intro-card {
    margin-bottom: 30px;
    
    .intro-content {
      display: flex;
      align-items: center;
      
      @media (max-width: 768px) {
        flex-direction: column-reverse;
      }
      
      .intro-text {
        flex: 3;
        padding-right: 30px;
        
        @media (max-width: 768px) {
          padding-right: 0;
          margin-top: 20px;
        }
        
        h2 {
          font-size: 22px;
          margin-top: 0;
          margin-bottom: 16px;
          color: #303133;
        }
        
        p {
          font-size: 16px;
          color: #606266;
          line-height: 1.6;
          margin-bottom: 24px;
        }
        
        .analysis-metrics {
          display: flex;
          margin-top: 20px;
          
          .metric {
            flex: 1;
            text-align: center;
            border-right: 1px solid #EBEEF5;
            padding: 0 15px;
            
            &:last-child {
              border-right: none;
            }
            
            .metric-value {
              font-size: 28px;
              font-weight: 600;
              color: #409EFF;
              margin-bottom: 8px;
            }
            
            .metric-label {
              font-size: 14px;
              color: #909399;
            }
          }
        }
      }
      
      .intro-image {
        flex: 2;
        text-align: center;
        
        img {
          max-width: 100%;
          max-height: 240px;
        }
      }
    }
  }
  
  .analysis-types {
    margin-bottom: 30px;
    
    .analysis-card {
      height: 100%;
      transition: all 0.3s;
      position: relative;
      overflow: hidden;
      margin-bottom: 20px;
      
      &:hover {
        transform: translateY(-5px);
        box-shadow: 0 10px 20px rgba(0, 0, 0, 0.1);
      }
      
      .analysis-card-header {
        height: 120px;
        display: flex;
        align-items: center;
        justify-content: center;
        color: #fff;
        
        i {
          font-size: 48px;
        }
        
        &.blue-header {
          background: linear-gradient(135deg, #1976d2, #64b5f6);
        }
        
        &.green-header {
          background: linear-gradient(135deg, #388e3c, #81c784);
        }
        
        &.orange-header {
          background: linear-gradient(135deg, #e64a19, #ff8a65);
        }
      }
      
      .analysis-card-content {
        padding: 20px;
        
        h3 {
          margin-top: 0;
          margin-bottom: 12px;
          font-size: 18px;
          color: #303133;
        }
        
        p {
          color: #606266;
          font-size: 14px;
          line-height: 1.5;
          margin-bottom: 20px;
          height: 63px;
          overflow: hidden;
        }
        
        .card-actions {
          display: flex;
          justify-content: space-between;
        }
      }
    }
  }
  
  .recent-analyses {
    .section-title {
      font-size: 20px;
      font-weight: 600;
      margin-bottom: 20px;
      position: relative;
      padding-left: 12px;
      
      &:before {
        content: '';
        position: absolute;
        left: 0;
        top: 50%;
        transform: translateY(-50%);
        width: 4px;
        height: 18px;
        background-color: #409EFF;
        border-radius: 2px;
      }
    }
  }
}
</style> 