<template>
  <div class="trend-prediction-detail">
    <!-- 预测结果卡片 -->
    <el-card class="detail-card">
      <template v-slot:header>
        <div class="card-header">
          <span><i class="el-icon-data-line"></i> 漏洞数量趋势预测 ({{ analysisData.timeRange }})</span>
        </div>
      </template>
      <div class="prediction-chart">
        <div class="chart-container">
          <ve-line 
            :data="chartData" 
            :settings="chartSettings"
            :extend="chartExtend"
            height="400px">
          </ve-line>
        </div>
      </div>
    </el-card>
    
    <!-- 预测细分分析卡片 -->
    <el-card class="detail-card">
      <template v-slot:header>
        <div class="card-header">
          <span><i class="el-icon-pie-chart"></i> 预测结果细分</span>
        </div>
      </template>
      <div class="prediction-breakdown">
        <div class="breakdown-chart">
          <ve-pie
            :data="breakdownChartData"
            :settings="pieSettings"
            height="300px">
          </ve-pie>
        </div>
        <div class="breakdown-stats">
          <div v-for="(value, key) in analysisData.predictedCounts" :key="key" class="stat-item">
            <div class="stat-label">{{ key }}</div>
            <div class="stat-value">{{ value }}</div>
            <div class="stat-bar">
              <div class="stat-bar-inner" :style="{ width: `${calculatePercentage(value)}%`, backgroundColor: getSeverityColor(key) }"></div>
            </div>
          </div>
        </div>
      </div>
    </el-card>
    
    <!-- 影响因素卡片 -->
    <el-card class="detail-card">
      <template v-slot:header>
        <div class="card-header">
          <span><i class="el-icon-s-marketing"></i> 趋势影响因素</span>
        </div>
      </template>
      <div class="impact-factors">
        <el-table :data="analysisData.trendFactors" style="width: 100%">
          <el-table-column prop="factor" label="影响因素" min-width="150"></el-table-column>
          <el-table-column label="影响程度" width="180">
            <template v-slot="scope">
              <div class="impact-gauge">
                <div class="impact-value" :class="{ 'negative': scope.row.impact < 0 }">
                  {{ formatImpact(scope.row.impact) }}
                </div>
                <div class="impact-bar-container">
                  <div 
                    class="impact-bar" 
                    :class="{ 'negative': scope.row.impact < 0 }"
                    :style="{ width: `${Math.abs(scope.row.impact * 100)}%` }">
                  </div>
                </div>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="置信度" width="120">
            <template v-slot="scope">
              <el-progress :percentage="scope.row.confidence * 100" :format="percentFormat"></el-progress>
            </template>
          </el-table-column>
          <el-table-column prop="description" label="说明" min-width="250"></el-table-column>
        </el-table>
      </div>
    </el-card>
  </div>
</template>

<script>
export default {
  name: 'TrendPredictionDetail',
  props: {
    analysisData: {
      type: Object,
      required: true
    }
  },
  data() {
    return {
      chartSettings: {
        labelMap: {
          'count': '漏洞数量'
        },
        yAxisName: ['数量'],
        area: true,
        smooth: true
      },
      chartExtend: {
        series: {
          areaStyle: {
            color: {
              type: 'linear',
              x: 0,
              y: 0,
              x2: 0,
              y2: 1,
              colorStops: [{
                offset: 0, color: 'rgba(64, 158, 255, 0.8)' 
              }, {
                offset: 1, color: 'rgba(64, 158, 255, 0.1)'
              }]
            }
          },
          lineStyle: {
            width: 3
          },
          itemStyle: {
            borderWidth: 2
          }
        },
        xAxis: {
          axisLabel: {
            showMaxLabel: true
          }
        }
      },
      pieSettings: {
        radius: '65%',
        offsetY: '60px'
      }
    }
  },
  computed: {
    chartData() {
      // 准备历史数据
      const historyRows = this.analysisData.historicalData.map(item => ({
        'date': this.formatDateShort(item.date),
        'count': item.count
      }))
      
      // 生成预测数据
      const predictionDates = this.generateFutureDates()
      
      // 计算预测趋势值（简单示例）
      const predictedTotal = Object.values(this.analysisData.predictedCounts).reduce((a, b) => a + b, 0)
      const lastHistoryValue = historyRows[historyRows.length - 1]?.count || 0
      
      const predictionSteps = predictionDates.length
      const predictionIncrement = (predictedTotal / predictionSteps - lastHistoryValue) / predictionSteps
      
      const predictionRows = predictionDates.map((date, index) => ({
        'date': this.formatDateShort(date),
        'count': Math.round(lastHistoryValue + predictionIncrement * (index + 1))
      }))
      
      // 合并数据，并标记历史和预测的分界点
      const boundary = historyRows.length
      const rows = [...historyRows, ...predictionRows]
      
      return {
        columns: ['date', 'count'],
        rows,
        boundary
      }
    },
    breakdownChartData() {
      const rows = Object.entries(this.analysisData.predictedCounts).map(([severity, count]) => ({
        'severity': severity,
        'count': count
      }))
      
      return {
        columns: ['severity', 'count'],
        rows
      }
    }
  },
  methods: {
    formatDateShort(date) {
      const d = new Date(date)
      return `${d.getMonth() + 1}/${d.getDate()}`
    },
    generateFutureDates() {
      // 根据时间范围生成未来日期
      const days = this.getTimeRangeDays()
      const today = new Date()
      
      return Array.from({ length: days }, (_, i) => {
        const date = new Date(today)
        date.setDate(date.getDate() + i + 1)
        return date
      })
    },
    getTimeRangeDays() {
      // 将时间范围转换为天数
      const range = this.analysisData.timeRange
      
      if (range.includes('7')) return 7
      if (range.includes('14')) return 14
      if (range.includes('30')) return 30
      if (range.includes('90')) return 90
      
      return 30 // 默认30天
    },
    formatImpact(impact) {
      const value = (impact * 100).toFixed(0)
      return impact >= 0 ? `+${value}%` : `${value}%`
    },
    calculatePercentage(value) {
      const total = Object.values(this.analysisData.predictedCounts).reduce((a, b) => a + b, 0)
      return total > 0 ? (value / total * 100) : 0
    },
    getSeverityColor(severity) {
      const colorMap = {
        '高危': '#F56C6C',
        '中危': '#E6A23C',
        '低危': '#67C23A',
        '信息': '#909399'
      }
      return colorMap[severity] || '#409EFF'
    },
    percentFormat(percentage) {
      return percentage.toFixed(0) + '%'
    }
  }
}
</script>

<style lang="scss" scoped>
.trend-prediction-detail {
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
    
    .prediction-chart {
      padding: 10px 0;
      
      .chart-container {
        height: 400px;
      }
    }
    
    .prediction-breakdown {
      display: flex;
      
      @media (max-width: 992px) {
        flex-direction: column;
      }
      
      .breakdown-chart {
        flex: 1;
      }
      
      .breakdown-stats {
        flex: 1;
        padding: 20px;
        
        .stat-item {
          margin-bottom: 15px;
          
          .stat-label {
            font-weight: 500;
            margin-bottom: 5px;
            display: flex;
            justify-content: space-between;
          }
          
          .stat-value {
            font-size: 20px;
            font-weight: 600;
            margin-bottom: 8px;
          }
          
          .stat-bar {
            height: 8px;
            background-color: #EBEEF5;
            border-radius: 4px;
            overflow: hidden;
            
            .stat-bar-inner {
              height: 100%;
              border-radius: 4px;
            }
          }
        }
      }
    }
    
    .impact-factors {
      .impact-gauge {
        display: flex;
        align-items: center;
        
        .impact-value {
          width: 60px;
          text-align: right;
          padding-right: 10px;
          font-weight: 500;
          color: #67C23A;
          
          &.negative {
            color: #F56C6C;
          }
        }
        
        .impact-bar-container {
          flex: 1;
          height: 16px;
          background-color: #F5F7FA;
          border-radius: 8px;
          position: relative;
          overflow: hidden;
          
          .impact-bar {
            height: 100%;
            background-color: #67C23A;
            border-radius: 8px;
            
            &.negative {
              background-color: #F56C6C;
            }
          }
        }
      }
    }
  }
}
</style> 