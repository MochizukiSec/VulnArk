<template>
  <div class="resource-optimization-detail">
    <!-- 资源分配对比卡片 -->
    <el-card class="detail-card">
      <template v-slot:header>
        <div class="card-header">
          <span><i class="el-icon-s-operation"></i> 资源分配对比</span>
        </div>
      </template>
      <div class="allocation-comparison">
        <div class="comparison-charts">
          <div class="chart-wrapper">
            <h3 class="chart-title">当前资源分配</h3>
            <ve-pie
              :data="currentAllocationData"
              :settings="pieSettings"
              height="300px">
            </ve-pie>
          </div>
          <div class="comparison-arrow">
            <i class="el-icon-right"></i>
          </div>
          <div class="chart-wrapper">
            <h3 class="chart-title">推荐资源分配</h3>
            <ve-pie
              :data="recommendedAllocationData"
              :settings="pieSettings"
              height="300px">
            </ve-pie>
          </div>
        </div>
        
        <div class="comparison-table">
          <el-table :data="comparisonTableData" style="width: 100%">
            <el-table-column prop="category" label="资源类别"></el-table-column>
            <el-table-column prop="current" label="当前分配">
              <template v-slot="scope">
                {{ scope.row.current }}%
              </template>
            </el-table-column>
            <el-table-column prop="recommended" label="推荐分配">
              <template v-slot="scope">
                {{ scope.row.recommended }}%
              </template>
            </el-table-column>
            <el-table-column label="变化">
              <template v-slot="scope">
                <span :class="getChangeClass(scope.row.change)">
                  {{ formatChange(scope.row.change) }}
                </span>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </div>
    </el-card>
    
    <!-- 潜在改进卡片 -->
    <el-card class="detail-card">
      <template v-slot:header>
        <div class="card-header">
          <span><i class="el-icon-data-analysis"></i> 潜在改进效果</span>
        </div>
      </template>
      <div class="improvement-info">
        <div class="improvement-metric">
          <div class="metric-label">预计效率提升</div>
          <div class="metric-value">{{ analysisData.potentialImprovement.toFixed(1) }}%</div>
          <el-progress :percentage="analysisData.potentialImprovement" :format="percentFormat" status="success"></el-progress>
        </div>
        <div class="improvement-description">
          <p>通过优化资源分配，预计可以提高安全团队的整体效率，减少漏洞修复时间，提升安全防御能力。</p>
        </div>
      </div>
    </el-card>
    
    <!-- 瓶颈分析卡片 -->
    <el-card class="detail-card">
      <template v-slot:header>
        <div class="card-header">
          <span><i class="el-icon-warning-outline"></i> 已识别的资源瓶颈</span>
        </div>
      </template>
      <div class="bottlenecks">
        <el-collapse accordion>
          <el-collapse-item v-for="(bottleneck, index) in analysisData.bottlenecksIdentified" :key="index">
            <template v-slot:title>
              <div class="bottleneck-title">
                <el-tag :type="getSeverityType(bottleneck.severity)" size="medium">{{ bottleneck.severity }}</el-tag>
                <span class="bottleneck-area">{{ bottleneck.area }}</span>
              </div>
            </template>
            <div class="bottleneck-content">
              <p class="bottleneck-description">
                <strong>问题描述：</strong> {{ bottleneck.description }}
              </p>
              <p class="bottleneck-solution">
                <strong>推荐解决方案：</strong> {{ bottleneck.solution }}
              </p>
            </div>
          </el-collapse-item>
        </el-collapse>
      </div>
    </el-card>
  </div>
</template>

<script>
export default {
  name: 'ResourceOptimizationDetail',
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
        offsetY: '60px',
        label: {
          formatter: '{b}: {d}%'
        }
      }
    }
  },
  computed: {
    currentAllocationData() {
      return this.transformAllocationData(this.analysisData.currentAllocation)
    },
    recommendedAllocationData() {
      return this.transformAllocationData(this.analysisData.recommendedAllocation)
    },
    comparisonTableData() {
      const data = []
      
      // 获取所有资源类别
      const categories = [...new Set([
        ...Object.keys(this.analysisData.currentAllocation),
        ...Object.keys(this.analysisData.recommendedAllocation)
      ])]
      
      // 构建表格数据
      categories.forEach(category => {
        const current = this.analysisData.currentAllocation[category] || 0
        const recommended = this.analysisData.recommendedAllocation[category] || 0
        const change = recommended - current
        
        data.push({
          category,
          current,
          recommended,
          change
        })
      })
      
      return data
    }
  },
  methods: {
    transformAllocationData(allocation) {
      const rows = Object.entries(allocation).map(([category, value]) => ({
        category,
        percentage: value
      }))
      
      return {
        columns: ['category', 'percentage'],
        rows
      }
    },
    formatChange(value) {
      return value > 0 ? `+${value}%` : `${value}%`
    },
    getChangeClass(value) {
      if (value > 0) return 'change-increase'
      if (value < 0) return 'change-decrease'
      return 'change-neutral'
    },
    getSeverityType(severity) {
      const types = {
        '高': 'danger',
        '中': 'warning',
        '低': 'info'
      }
      return types[severity] || 'info'
    },
    percentFormat(percentage) {
      return percentage.toFixed(1) + '%'
    }
  }
}
</script>

<style lang="scss" scoped>
.resource-optimization-detail {
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
    
    .allocation-comparison {
      .comparison-charts {
        display: flex;
        align-items: center;
        justify-content: space-between;
        margin-bottom: 30px;
        
        @media (max-width: 992px) {
          flex-direction: column;
        }
        
        .chart-wrapper {
          flex: 1;
          
          .chart-title {
            text-align: center;
            font-size: 16px;
            font-weight: 500;
            margin-bottom: 10px;
            color: #303133;
          }
        }
        
        .comparison-arrow {
          display: flex;
          justify-content: center;
          align-items: center;
          padding: 0 20px;
          
          @media (max-width: 992px) {
            transform: rotate(90deg);
            padding: 20px 0;
          }
          
          i {
            font-size: 32px;
            color: #909399;
          }
        }
      }
      
      .comparison-table {
        margin-top: 20px;
        
        .change-increase {
          color: #67C23A;
          font-weight: bold;
        }
        
        .change-decrease {
          color: #F56C6C;
          font-weight: bold;
        }
        
        .change-neutral {
          color: #909399;
        }
      }
    }
    
    .improvement-info {
      display: flex;
      
      @media (max-width: 768px) {
        flex-direction: column;
      }
      
      .improvement-metric {
        flex: 1;
        padding: 20px;
        background-color: #F5F7FA;
        border-radius: 8px;
        margin-right: 20px;
        
        @media (max-width: 768px) {
          margin-right: 0;
          margin-bottom: 20px;
        }
        
        .metric-label {
          font-size: 16px;
          margin-bottom: 10px;
          color: #606266;
        }
        
        .metric-value {
          font-size: 36px;
          font-weight: 700;
          color: #67C23A;
          margin-bottom: 15px;
        }
      }
      
      .improvement-description {
        flex: 2;
        
        p {
          font-size: 16px;
          line-height: 1.6;
          color: #606266;
        }
      }
    }
    
    .bottlenecks {
      .bottleneck-title {
        display: flex;
        align-items: center;
        
        .bottleneck-area {
          margin-left: 10px;
          font-weight: 500;
        }
      }
      
      .bottleneck-content {
        padding: 10px 0;
        
        .bottleneck-description,
        .bottleneck-solution {
          margin: 10px 0;
          line-height: 1.5;
          color: #606266;
        }
      }
    }
  }
}
</style> 