<template>
  <div class="analysis-detail">
    <div v-if="loading" class="loading-container">
      <el-skeleton style="width: 100%" animated>
        <template v-slot:template>
          <el-skeleton-item variant="text" style="width: 50%; height: 40px; margin-bottom: 20px;"></el-skeleton-item>
          <el-skeleton-item variant="text" style="width: 80%; height: 20px; margin-bottom: 10px;"></el-skeleton-item>
          <el-skeleton-item variant="text" style="width: 70%; height: 20px; margin-bottom: 40px;"></el-skeleton-item>
          
          <div style="display: flex; justify-content: space-between; margin-bottom: 20px;">
            <el-skeleton-item variant="text" style="width: 30%; height: 100px;"></el-skeleton-item>
            <el-skeleton-item variant="text" style="width: 30%; height: 100px;"></el-skeleton-item>
            <el-skeleton-item variant="text" style="width: 30%; height: 100px;"></el-skeleton-item>
          </div>
          
          <el-skeleton-item variant="p" style="width: 100%; height: 300px;"></el-skeleton-item>
        </template>
      </el-skeleton>
    </div>
    
    <div v-else-if="!analysis" class="not-found">
      <el-empty description="未找到分析结果"></el-empty>
    </div>
    
    <div v-else class="analysis-content">
      <!-- 标题和基本信息 -->
      <div class="analysis-header">
        <h1 class="analysis-title">{{ analysis.title }}</h1>
        <div class="analysis-meta">
          <div class="meta-item">
            <i class="el-icon-time"></i>
            <span>{{ formatDate(analysis.createdAt) }}</span>
          </div>
          <div class="meta-item">
            <i class="el-icon-data-analysis"></i>
            <span>{{ getAnalysisTypeLabel(analysis.type) }}</span>
          </div>
          <div class="meta-item">
            <i class="el-icon-pie-chart"></i>
            <span>置信度: {{ (analysis.confidence * 100).toFixed(0) }}%</span>
          </div>
        </div>
        <div class="analysis-description">
          {{ analysis.description }}
        </div>
      </div>
      
      <!-- 分析结果总结 -->
      <el-card class="summary-card">
        <template v-slot:header>
          <div class="summary-header">
            <span><i class="el-icon-s-data"></i> 分析结果总结</span>
          </div>
        </template>
        <div class="summary-content">
          <div class="recommendations">
            <h3>主要发现与建议</h3>
            <ul class="recommendation-list">
              <li v-for="(recommendation, index) in analysis.recommendations" :key="index">
                {{ recommendation }}
              </li>
            </ul>
          </div>
        </div>
      </el-card>
      
      <!-- 根据分析类型展示不同的详细内容 -->
      <component 
        :is="getDetailComponent()" 
        :analysis-data="analysis.analysisData">
      </component>
      
      <!-- 操作按钮 -->
      <div class="action-buttons">
        <el-button plain @click="$emit('close')">关闭</el-button>
        <el-button plain icon="el-icon-download" @click="downloadReport">下载报告</el-button>
        <el-button type="primary" icon="el-icon-share" @click="shareAnalysis">分享结果</el-button>
      </div>
    </div>
  </div>
</template>

<script>
import { formatDate } from '@/utils/format'
import TrendPredictionDetail from './details/TrendPredictionDetail.vue'
import ResourceOptimizationDetail from './details/ResourceOptimizationDetail.vue'
import AnomalyDetectionDetail from './details/AnomalyDetectionDetail.vue'

export default {
  name: 'AnalysisDetail',
  components: {
    TrendPredictionDetail,
    ResourceOptimizationDetail,
    AnomalyDetectionDetail
  },
  props: {
    analysisId: {
      type: String,
      required: true
    }
  },
  data() {
    return {
      loading: true,
      analysis: null
    }
  },
  created() {
    this.fetchAnalysisDetail()
  },
  methods: {
    formatDate,
    async fetchAnalysisDetail() {
      this.loading = true
      
      try {
        // 实际项目中应调用API获取详情
        // const { data } = await this.$api.ai.getAnalysisById(this.analysisId)
        // this.analysis = data
        
        // 模拟数据
        await new Promise(resolve => setTimeout(resolve, 1000))
        
        // 根据ID使用不同的模拟数据
        switch (this.analysisId) {
          case '61e8f5c3b32f6c001a9b3d01':
            this.analysis = this.getTrendPredictionData()
            break
          case '61e8f5c3b32f6c001a9b3d02':
            this.analysis = this.getResourceOptimizationData()
            break
          case '61e8f5c3b32f6c001a9b3d03':
            this.analysis = this.getAnomalyDetectionData()
            break
          default:
            this.analysis = null
        }
      } catch (error) {
        this.$message.error('获取分析详情失败: ' + error.message)
        this.analysis = null
      } finally {
        this.loading = false
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
    getDetailComponent() {
      if (!this.analysis) return null
      
      switch(this.analysis.type) {
        case 'trend_prediction':
          return 'trend-prediction-detail'
        case 'resource_optimization':
          return 'resource-optimization-detail'
        case 'anomaly_detection':
          return 'anomaly-detection-detail'
        default:
          return null
      }
    },
    downloadReport() {
      this.$message.info('报告下载功能正在开发中...')
    },
    shareAnalysis() {
      this.$message.info('分享功能正在开发中...')
    },
    // 模拟数据生成方法
    getTrendPredictionData() {
      const now = new Date()
      const historicalData = Array.from({ length: 90 }, (_, i) => {
        const date = new Date(now)
        date.setDate(date.getDate() - 90 + i)
        return {
          date,
          count: 10 + Math.floor(Math.random() * 40),
          label: i % 7 === 0 ? '周期性峰值' : undefined
        }
      })
      
      return {
        id: '61e8f5c3b32f6c001a9b3d01',
        title: '未来30天漏洞趋势预测',
        type: 'trend_prediction',
        description: '基于历史数据的未来30天安全漏洞趋势预测分析',
        confidence: 0.85,
        createdAt: new Date(Date.now() - 86400000 * 2),
        recommendations: [
          '增加对高风险类别漏洞的监控频率',
          '优先分配资源修复预测增长最快的漏洞类型',
          '为可能增长的漏洞类型提前制定应对策略',
          '关注影响因素中的主要驱动因素，采取预防措施'
        ],
        analysisData: {
          timeRange: '未来30天',
          predictedCounts: {
            '高危': 120,
            '中危': 180,
            '低危': 90,
            '信息': 45
          },
          trendFactors: [
            {
              factor: '新版本软件发布',
              impact: 0.25,
              confidence: 0.85,
              description: '新发布的软件版本通常会引入新的漏洞'
            },
            {
              factor: '安全团队扩展',
              impact: -0.15,
              confidence: 0.78,
              description: '安全团队扩展可能改善漏洞检测和修复效率'
            },
            {
              factor: '第三方组件使用增加',
              impact: 0.30,
              confidence: 0.92,
              description: '更多第三方组件可能引入更多潜在漏洞'
            },
            {
              factor: '安全培训计划',
              impact: -0.18,
              confidence: 0.75,
              description: '安全意识培训可能减少开发引入的漏洞'
            }
          ],
          historicalData
        }
      }
    },
    getResourceOptimizationData() {
      return {
        id: '61e8f5c3b32f6c001a9b3d02',
        title: '安全资源优化配置建议',
        type: 'resource_optimization',
        description: '基于当前漏洞状态和历史数据的安全资源分配优化建议',
        confidence: 0.78,
        createdAt: new Date(Date.now() - 86400000 * 5),
        recommendations: [
          '将更多资源分配给高危漏洞的修复，以降低整体风险',
          '减少对低危漏洞的资源投入，提高资源利用效率',
          '解决已识别的资源瓶颈，特别是团队协作问题',
          '考虑引入自动化工具，减少手动任务所需的资源'
        ],
        analysisData: {
          currentAllocation: {
            '高危漏洞修复': 30,
            '中危漏洞修复': 25,
            '低危漏洞修复': 20,
            '漏洞扫描': 15,
            '安全培训': 10
          },
          recommendedAllocation: {
            '高危漏洞修复': 40,
            '中危漏洞修复': 25,
            '低危漏洞修复': 10,
            '漏洞扫描': 20,
            '安全培训': 5
          },
          potentialImprovement: 25.5,
          bottlenecksIdentified: [
            {
              area: '漏洞修复流程',
              severity: '高',
              description: '开发团队和安全团队之间的协作效率低',
              solution: '实施更有效的团队协作工具和流程'
            },
            {
              area: '漏洞验证',
              severity: '中',
              description: '漏洞修复验证过程耗时较长',
              solution: '引入自动化测试和验证工具'
            },
            {
              area: '低危漏洞积压',
              severity: '低',
              description: '低危漏洞数量较多但优先级不足',
              solution: '批量修复策略或考虑接受部分风险'
            }
          ]
        }
      }
    },
    getAnomalyDetectionData() {
      return {
        id: '61e8f5c3b32f6c001a9b3d03',
        title: '过去7天漏洞异常检测',
        type: 'anomaly_detection',
        description: '检测漏洞数据中的异常模式和潜在安全事件',
        confidence: 0.82,
        createdAt: new Date(Date.now() - 86400000 * 1),
        recommendations: [
          '调查检测到的高严重度异常，确认是否存在安全事件',
          '关注Web应用安全区域的异常增长，可能表明新的攻击模式',
          '审查数据库安全相关的漏洞类型，加强相关防御措施',
          '考虑增加受影响区域的监控频率'
        ],
        analysisData: {
          anomaliesDetected: [
            {
              type: '突增',
              severity: '高',
              description: 'Web应用关键漏洞数量突然增加',
              detectedAt: new Date(Date.now() - 86400000 * 2),
              affectedArea: 'Web应用安全',
              score: 0.92
            },
            {
              type: '异常模式',
              severity: '中',
              description: '数据库相关漏洞呈周期性出现',
              detectedAt: new Date(Date.now() - 86400000 * 5),
              affectedArea: '数据库安全',
              score: 0.78
            },
            {
              type: '偏差',
              severity: '低',
              description: '认证模块漏洞修复率低于基准',
              detectedAt: new Date(Date.now() - 86400000 * 3),
              affectedArea: '身份认证',
              score: 0.65
            }
          ],
          timeRange: '过去7天',
          baselineData: {
            dailyAverage: 8.5,
            weeklyTrend: '稳定',
            normalRange: [5.0, 12.0]
          }
        }
      }
    }
  }
}
</script>

<style lang="scss" scoped>
.analysis-detail {
  .loading-container {
    padding: 20px;
  }
  
  .not-found {
    padding: 40px 0;
  }
  
  .analysis-content {
    .analysis-header {
      margin-bottom: 30px;
      
      .analysis-title {
        font-size: 24px;
        font-weight: 600;
        margin-bottom: 12px;
        color: #303133;
      }
      
      .analysis-meta {
        display: flex;
        margin-bottom: 16px;
        
        .meta-item {
          display: flex;
          align-items: center;
          margin-right: 24px;
          color: #606266;
          
          i {
            margin-right: 8px;
            font-size: 16px;
          }
        }
      }
      
      .analysis-description {
        font-size: 16px;
        color: #606266;
        line-height: 1.6;
      }
    }
    
    .summary-card {
      margin-bottom: 24px;
      
      .summary-header {
        display: flex;
        align-items: center;
        
        i {
          margin-right: 10px;
          font-size: 18px;
        }
      }
      
      .summary-content {
        .recommendations {
          h3 {
            font-size: 16px;
            margin-bottom: 12px;
            color: #303133;
          }
          
          .recommendation-list {
            margin: 0;
            padding-left: 20px;
            
            li {
              margin-bottom: 10px;
              color: #606266;
              line-height: 1.5;
            }
          }
        }
      }
    }
    
    .action-buttons {
      margin-top: 40px;
      display: flex;
      justify-content: flex-end;
      gap: 12px;
    }
  }
}
</style> 