<template>
  <div class="anomaly-detection-detail">
    <!-- 异常总览卡片 -->
    <el-card class="detail-card">
      <template v-slot:header>
        <div class="card-header">
          <span><i class="el-icon-warning-outline"></i> 异常检测总览 ({{ analysisData.timeRange }})</span>
        </div>
      </template>
      <div class="anomaly-overview">
        <div class="overview-metrics">
          <div class="metric-item">
            <div class="metric-icon">
              <i class="el-icon-warning"></i>
            </div>
            <div class="metric-info">
              <div class="metric-value">{{ anomaliesCount }}</div>
              <div class="metric-label">检测到的异常</div>
            </div>
          </div>
          <div class="metric-item">
            <div class="metric-icon">
              <i class="el-icon-bell"></i>
            </div>
            <div class="metric-info">
              <div class="metric-value">{{ highSeverityCount }}</div>
              <div class="metric-label">高严重度异常</div>
            </div>
          </div>
          <div class="metric-item">
            <div class="metric-icon">
              <i class="el-icon-data-line"></i>
            </div>
            <div class="metric-info">
              <div class="metric-value">{{ analysisData.baselineData.dailyAverage }}</div>
              <div class="metric-label">日均漏洞数</div>
            </div>
          </div>
        </div>
        <div class="overview-description">
          <p>在{{ analysisData.timeRange }}内，异常检测系统共发现 {{ anomaliesCount }} 个异常模式，其中 {{ highSeverityCount }} 个为高严重度异常。系统的标准漏洞数量在 {{ formatRange(analysisData.baselineData.normalRange) }} 之间，日均为 {{ analysisData.baselineData.dailyAverage }} 个漏洞。</p>
          <p>本次分析使用的基线数据显示，整体趋势为 <strong>{{ analysisData.baselineData.weeklyTrend }}</strong>。</p>
        </div>
      </div>
    </el-card>
    
    <!-- 详细异常卡片 -->
    <el-card class="detail-card">
      <template v-slot:header>
        <div class="card-header">
          <span><i class="el-icon-view"></i> 详细异常分析</span>
        </div>
      </template>
      <div class="anomaly-details">
        <el-table :data="sortedAnomalies" style="width: 100%">
          <el-table-column label="异常类型" width="120">
            <template v-slot="scope">
              <el-tag :type="getAnomalyTypeTag(scope.row.type)">{{ scope.row.type }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="严重度" width="100">
            <template v-slot="scope">
              <el-tag :type="getSeverityType(scope.row.severity)" size="medium">{{ scope.row.severity }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="affectedArea" label="影响区域" width="150"></el-table-column>
          <el-table-column prop="description" label="描述" min-width="250"></el-table-column>
          <el-table-column label="检测时间" width="180">
            <template v-slot="scope">
              {{ formatDate(scope.row.detectedAt) }}
            </template>
          </el-table-column>
          <el-table-column label="异常分数" width="150">
            <template v-slot="scope">
              <el-progress :percentage="scope.row.score * 100" :format="percentFormat" :color="getScoreColor(scope.row.score)"></el-progress>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-card>
    
    <!-- 异常分布卡片 -->
    <el-card class="detail-card">
      <template v-slot:header>
        <div class="card-header">
          <span><i class="el-icon-s-marketing"></i> 异常分布分析</span>
        </div>
      </template>
      <div class="anomaly-distribution">
        <div class="distribution-charts">
          <div class="chart-item">
            <h3 class="chart-title">按严重度分布</h3>
            <ve-pie
              :data="severityDistributionData"
              :settings="pieSettings"
              height="300px">
            </ve-pie>
          </div>
          <div class="chart-item">
            <h3 class="chart-title">按影响区域分布</h3>
            <ve-pie
              :data="areaDistributionData"
              :settings="pieSettings"
              height="300px">
            </ve-pie>
          </div>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script>
export default {
  name: 'AnomalyDetectionDetail',
  props: {
    analysisData: {
      type: Object,
      required: true
    }
  },
  data() {
    return {
      pieSettings: {
        radius: '65%',
        offsetY: '20px',
        label: {
          formatter: '{b}: {c} ({d}%)'
        }
      }
    }
  },
  computed: {
    anomaliesCount() {
      return this.analysisData.anomaliesDetected.length
    },
    highSeverityCount() {
      return this.analysisData.anomaliesDetected.filter(anomaly => 
        anomaly.severity === '高'
      ).length
    },
    sortedAnomalies() {
      // 创建数组的副本并排序，而不是直接修改props
      return [...this.analysisData.anomaliesDetected].sort((a, b) => b.score - a.score)
    },
    severityDistributionData() {
      // 按严重度聚合数据
      const distribution = this.analysisData.anomaliesDetected.reduce((acc, anomaly) => {
        const severity = anomaly.severity
        if (!acc[severity]) acc[severity] = 0
        acc[severity]++
        return acc
      }, {})
      
      // 转换为图表需要的格式
      const rows = Object.entries(distribution).map(([severity, count]) => ({
        severity,
        count
      }))
      
      return {
        columns: ['severity', 'count'],
        rows
      }
    },
    areaDistributionData() {
      // 按影响区域聚合数据
      const distribution = this.analysisData.anomaliesDetected.reduce((acc, anomaly) => {
        const area = anomaly.affectedArea
        if (!acc[area]) acc[area] = 0
        acc[area]++
        return acc
      }, {})
      
      // 转换为图表需要的格式
      const rows = Object.entries(distribution).map(([area, count]) => ({
        area,
        count
      }))
      
      return {
        columns: ['area', 'count'],
        rows
      }
    }
  },
  methods: {
    formatDate(date) {
      const d = new Date(date)
      return `${d.getFullYear()}-${d.getMonth() + 1}-${d.getDate()} ${d.getHours()}:${d.getMinutes().toString().padStart(2, '0')}`
    },
    getAnomalyTypeTag(type) {
      const typeMap = {
        '突增': 'danger',
        '突降': 'warning',
        '异常模式': 'primary',
        '偏差': 'info'
      }
      return typeMap[type] || 'info'
    },
    getSeverityType(severity) {
      const types = {
        '高': 'danger',
        '中': 'warning',
        '低': 'info'
      }
      return types[severity] || 'info'
    },
    getScoreColor(score) {
      if (score >= 0.8) return '#F56C6C'
      if (score >= 0.6) return '#E6A23C'
      return '#909399'
    },
    percentFormat(percentage) {
      return percentage.toFixed(0) + '%'
    },
    formatRange(range) {
      if (!range || !Array.isArray(range) || range.length !== 2) return '未知'
      return `${range[0]} - ${range[1]}`
    }
  }
}
</script>

<style lang="scss" scoped>
.anomaly-detection-detail {
  .detail-card {
    margin-bottom: 24px;
    
    .card-header {
      display: flex;
      align-items: center;
      
      i {
        margin-right: 10px;
        font-size: 18px;
      }
    }
    
    .anomaly-overview {
      .overview-metrics {
        display: flex;
        margin-bottom: 20px;
        
        @media (max-width: 768px) {
          flex-direction: column;
        }
        
        .metric-item {
          flex: 1;
          display: flex;
          align-items: center;
          padding: 15px;
          border-radius: 8px;
          background-color: #F5F7FA;
          margin-right: 15px;
          
          @media (max-width: 768px) {
            margin-right: 0;
            margin-bottom: 15px;
          }
          
          &:last-child {
            margin-right: 0;
          }
          
          .metric-icon {
            width: 50px;
            height: 50px;
            border-radius: 50%;
            background-color: #EDF2FC;
            display: flex;
            align-items: center;
            justify-content: center;
            margin-right: 15px;
            
            i {
              font-size: 24px;
              color: #409EFF;
            }
            
            &:nth-child(2) i {
              color: #E6A23C;
            }
            
            &:nth-child(3) i {
              color: #67C23A;
            }
          }
          
          .metric-info {
            .metric-value {
              font-size: 24px;
              font-weight: 600;
              color: #303133;
              margin-bottom: 5px;
            }
            
            .metric-label {
              font-size: 14px;
              color: #606266;
            }
          }
        }
      }
      
      .overview-description {
        p {
          font-size: 14px;
          line-height: 1.6;
          color: #606266;
          margin-bottom: 10px;
          
          &:last-child {
            margin-bottom: 0;
          }
        }
      }
    }
    
    .anomaly-distribution {
      .distribution-charts {
        display: flex;
        
        @media (max-width: 992px) {
          flex-direction: column;
        }
        
        .chart-item {
          flex: 1;
          
          .chart-title {
            text-align: center;
            font-size: 16px;
            font-weight: 500;
            margin-bottom: 10px;
            color: #303133;
          }
        }
      }
    }
  }
}
</style> 