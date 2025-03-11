<template>
  <div class="dashboard-container">
    <!-- 页面标题 -->
    <div class="dashboard-header">
      <h1 class="dashboard-title">
        <span class="title-icon"><i class="el-icon-data-analysis"></i></span>
        仪表盘
        <span class="title-highlight">概览</span>
      </h1>
      <p class="dashboard-subtitle">查看系统安全状态和最新漏洞信息，实时掌握安全态势</p>
    </div>
    
    <el-row :gutter="24">
      <!-- 统计卡片 -->
      <el-col :xs="12" :sm="12" :md="6" :lg="6">
        <div class="stat-card red-gradient">
          <div class="stat-card-content">
            <div class="stat-icon-container">
              <i class="el-icon-warning"></i>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ dashboardData.totalVulnerabilities || 0 }}</div>
              <div class="stat-label">总漏洞数</div>
            </div>
          </div>
          <div class="stat-progress"></div>
        </div>
      </el-col>
      
      <el-col :xs="12" :sm="12" :md="6" :lg="6">
        <div class="stat-card orange-gradient">
          <div class="stat-card-content">
            <div class="stat-icon-container">
              <i class="el-icon-timer"></i>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ getStatusCount('open') + getStatusCount('in_progress') }}</div>
              <div class="stat-label">待处理漏洞</div>
            </div>
          </div>
          <div class="stat-progress"></div>
        </div>
      </el-col>
      
      <el-col :xs="12" :sm="12" :md="6" :lg="6">
        <div class="stat-card green-gradient">
          <div class="stat-card-content">
            <div class="stat-icon-container">
              <i class="el-icon-check"></i>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ getStatusCount('resolved') }}</div>
              <div class="stat-label">已解决漏洞</div>
            </div>
          </div>
          <div class="stat-progress"></div>
        </div>
      </el-col>
      
      <el-col :xs="12" :sm="12" :md="6" :lg="6">
        <div class="stat-card blue-gradient">
          <div class="stat-card-content">
            <div class="stat-icon-container">
              <i class="el-icon-document"></i>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ getCriticalHighCount() }}</div>
              <div class="stat-label">高危漏洞</div>
            </div>
          </div>
          <div class="stat-progress"></div>
        </div>
      </el-col>
    </el-row>
    
    <el-row :gutter="24" class="chart-row">
      <!-- 漏洞严重程度统计图 -->
      <el-col :xs="24" :sm="24" :md="12" :lg="12">
        <div class="chart-card">
          <div class="chart-header">
            <h3 class="chart-title">
              <i class="el-icon-pie-chart chart-icon"></i>
              漏洞严重程度分布
            </h3>
          </div>
          <div class="chart-container" v-loading="loading">
            <pie-chart 
              v-if="!loading && hasSeverityData" 
              :chart-data="severityChartData" 
              :options="pieChartOptions" 
            />
            <div v-else-if="!loading && !hasSeverityData" class="no-data">
              <i class="el-icon-data-analysis"></i>
              <p>暂无严重程度数据</p>
              <span class="no-data-tip">添加漏洞后将在此处显示数据</span>
            </div>
          </div>
        </div>
      </el-col>
      
      <!-- 漏洞状态统计图 -->
      <el-col :xs="24" :sm="24" :md="12" :lg="12">
        <div class="chart-card">
          <div class="chart-header">
            <h3 class="chart-title">
              <i class="el-icon-s-data chart-icon"></i>
              漏洞状态分布
            </h3>
          </div>
          <div class="chart-container" v-loading="loading">
            <doughnut-chart 
              v-if="!loading && hasStatusData" 
              :chart-data="statusChartData"
              :options="doughnutChartOptions"
            />
            <div v-else-if="!loading && !hasStatusData" class="no-data">
              <i class="el-icon-data-analysis"></i>
              <p>暂无状态数据</p>
              <span class="no-data-tip">添加漏洞后将在此处显示数据</span>
            </div>
          </div>
        </div>
      </el-col>
    </el-row>
    
    <el-row :gutter="24" class="chart-row">
      <!-- 月度漏洞趋势图 -->
      <el-col :xs="24">
        <div class="chart-card">
          <div class="chart-header">
            <h3 class="chart-title">
              <i class="el-icon-trend chart-icon"></i>
              月度漏洞趋势
            </h3>
          </div>
          <div class="chart-container" v-loading="loading">
            <line-chart 
              v-if="!loading && hasMonthlyData" 
              :chart-data="monthlyChartData"
              :options="lineChartOptions"
            />
            <div v-else-if="!loading && !hasMonthlyData" class="no-data">
              <i class="el-icon-data-analysis"></i>
              <p>暂无趋势数据</p>
              <span class="no-data-tip">需要至少两个月的数据才能显示趋势</span>
            </div>
          </div>
        </div>
      </el-col>
    </el-row>
    
    <el-row :gutter="24" class="chart-row">
      <!-- 风险分数指示器 -->
      <el-col :xs="24" :sm="24" :md="8" :lg="8">
        <div class="chart-card">
          <div class="chart-header">
            <h3 class="chart-title">
              <i class="el-icon-warning chart-icon"></i>
              整体风险评分
            </h3>
          </div>
          <div class="chart-container risk-score-container" v-loading="loading">
            <div class="risk-score-gauge">
              <div class="gauge-container">
                <div class="gauge" :style="{ background: `conic-gradient(${getRiskScoreColor(riskScore)} ${riskScore}%, #f5f5f5 0%)` }">
                  <div class="gauge-center">
                    <div class="gauge-value">{{ riskScore }}</div>
                    <div class="gauge-label">风险分</div>
                  </div>
                </div>
              </div>
              <div class="risk-level">
                <div class="risk-level-title">风险等级：</div>
                <div class="risk-level-value" :class="getRiskLevelClass(riskScore)">
                  {{ getRiskLevelText(riskScore) }}
                </div>
              </div>
            </div>
            <div class="risk-score-description">
              <p>{{ getRiskDescription(riskScore) }}</p>
            </div>
          </div>
        </div>
      </el-col>
      
      <!-- 修复进度 -->
      <el-col :xs="24" :sm="24" :md="16" :lg="16">
        <div class="chart-card">
          <div class="chart-header">
            <h3 class="chart-title">
              <i class="el-icon-s-data chart-icon"></i>
              修复进度
            </h3>
          </div>
          <div class="chart-container remediation-container" v-loading="loading">
            <div class="remediation-stats">
              <div class="remediation-stat-item">
                <div class="stat-value">{{ remediationProgress.resolvedCount }}</div>
                <div class="stat-label">已修复漏洞</div>
              </div>
              <div class="remediation-stat-item">
                <div class="stat-value">{{ remediationProgress.totalCount }}</div>
                <div class="stat-label">总漏洞数</div>
              </div>
              <div class="remediation-stat-item">
                <div class="stat-value">{{ Math.round(remediationProgress.averageDays) }} 天</div>
                <div class="stat-label">平均修复时间</div>
              </div>
            </div>
            <div class="progress-bar-container">
              <div class="progress-label">修复完成率</div>
              <el-progress 
                :percentage="Math.round(remediationProgress.progressRate)" 
                :color="getProgressColor(remediationProgress.progressRate)"
                :stroke-width="16"
                :show-text="false"
              ></el-progress>
              <div class="progress-value">{{ Math.round(remediationProgress.progressRate) }}%</div>
            </div>
          </div>
        </div>
      </el-col>
    </el-row>
    
    <el-row :gutter="24" class="chart-row">
      <!-- 团队漏洞统计 -->
      <el-col :xs="24" :sm="24" :md="12" :lg="12">
        <div class="chart-card">
          <div class="chart-header">
            <h3 class="chart-title">
              <i class="el-icon-s-custom chart-icon"></i>
              团队漏洞统计
            </h3>
          </div>
          <div class="chart-container team-container" v-loading="loading">
            <div v-if="!loading && teamVulnerabilities.length > 0" class="team-list">
              <div 
                v-for="(team, index) in teamVulnerabilities" 
                :key="index" 
                class="team-item"
              >
                <div class="team-info">
                  <div class="team-name">{{ team.team }}</div>
                  <div class="team-count">
                    <span class="count-value">{{ team.count }}</span>
                    <span class="count-label">漏洞</span>
                  </div>
                </div>
                <div class="team-progress">
                  <el-progress 
                    :percentage="getTeamProgressPercentage(team)" 
                    :color="getTeamProgressColor(team)"
                    :stroke-width="8"
                  ></el-progress>
                </div>
                <div class="team-severity">
                  <div class="severity-dots">
                    <span class="severity-dot critical" :style="{ opacity: team.severities.critical > 0 ? 1 : 0.3 }"></span>
                    <span class="severity-dot high" :style="{ opacity: team.severities.high > 0 ? 1 : 0.3 }"></span>
                    <span class="severity-dot medium" :style="{ opacity: team.severities.medium > 0 ? 1 : 0.3 }"></span>
                    <span class="severity-dot low" :style="{ opacity: team.severities.low > 0 ? 1 : 0.3 }"></span>
                  </div>
                </div>
              </div>
            </div>
            <div v-else-if="!loading && teamVulnerabilities.length === 0" class="no-data">
              <i class="el-icon-s-custom"></i>
              <p>暂无团队数据</p>
            </div>
          </div>
        </div>
      </el-col>
      
      <!-- 漏洞趋势图 -->
      <el-col :xs="24" :sm="24" :md="12" :lg="12">
        <div class="chart-card">
          <div class="chart-header">
            <h3 class="chart-title">
              <i class="el-icon-data-line chart-icon"></i>
              漏洞趋势
            </h3>
          </div>
          <div class="chart-container trend-container" v-loading="loading">
            <line-chart 
              v-if="!loading && hasTrendsData" 
              :chart-data="trendsChartData" 
              :options="lineChartOptions" 
            />
            <div v-else-if="!loading && !hasTrendsData" class="no-data">
              <i class="el-icon-data-line"></i>
              <p>暂无趋势数据</p>
              <span class="no-data-tip">数据收集中，请稍后查看</span>
            </div>
          </div>
        </div>
      </el-col>
    </el-row>
    
    <!-- 高危漏洞列表 -->
    <el-row :gutter="24" class="chart-row">
      <el-col :span="24">
        <div class="chart-card">
          <div class="chart-header">
            <h3 class="chart-title">
              <i class="el-icon-warning-outline chart-icon"></i>
              需优先关注的高危漏洞
            </h3>
            <div class="chart-actions">
              <el-button size="small" type="primary" plain @click="viewAllCritical">
                查看全部
              </el-button>
            </div>
          </div>
          <div class="chart-container critical-container" v-loading="loading">
            <el-table 
              v-if="!loading && criticalVulnerabilities.length > 0"
              :data="criticalVulnerabilities" 
              style="width: 100%"
              class="critical-table"
            >
              <el-table-column prop="title" label="漏洞标题" min-width="220">
                <template #default="scope">
                  <router-link 
                    :to="`/vulnerabilities/${scope.row.id}`" 
                    class="vulnerability-link"
                  >
                    {{ scope.row.title }}
                  </router-link>
                </template>
              </el-table-column>
              <el-table-column prop="severity" label="严重程度" width="120">
                <template #default="scope">
                  <el-tag :type="getSeverityType(scope.row.severity)" size="small">
                    {{ getSeverityText(scope.row.severity) }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="cvss" label="CVSS" width="100">
                <template #default="scope">
                  <span class="cvss-score">{{ scope.row.cvss.toFixed(1) }}</span>
                </template>
              </el-table-column>
              <el-table-column prop="status" label="状态" width="120">
                <template #default="scope">
                  <el-tag :type="getStatusType(scope.row.status)" size="small" effect="plain">
                    {{ getStatusText(scope.row.status) }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="createdAt" label="发现时间" width="180">
                <template #default="scope">
                  {{ formatDate(scope.row.createdAt) }}
                </template>
              </el-table-column>
            </el-table>
            <div v-else-if="!loading && criticalVulnerabilities.length === 0" class="no-data">
              <i class="el-icon-warning-outline"></i>
              <p>暂无需优先关注的高危漏洞</p>
              <span class="no-data-tip">恭喜！目前没有需要紧急处理的高危漏洞</span>
            </div>
          </div>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import { computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import PieChart from '@/components/charts/PieChart.vue';
import DoughnutChart from '@/components/charts/DoughnutChart.vue';
import LineChart from '@/components/charts/LineChart.vue';
import { useStore } from 'vuex';

export default {
  name: 'Dashboard',
  components: {
    PieChart,
    DoughnutChart,
    LineChart
  },
  setup() {
    const router = useRouter();
    const store = useStore();
    
    // 从Store获取数据和状态
    const loading = computed(() => store.getters['dashboard/isLoading']);
    const dashboardData = computed(() => store.getters['dashboard/dashboardData']);
    const riskScore = computed(() => store.getters['dashboard/riskScore']);
    const remediationProgress = computed(() => store.getters['dashboard/remediationProgress']);
    const teamVulnerabilities = computed(() => store.getters['dashboard/teamVulnerabilities'] || []);
    const criticalVulnerabilities = computed(() => store.getters['dashboard/criticalVulnerabilities'] || []);
    const vulnerabilityTrends = computed(() => store.getters['dashboard/vulnerabilityTrends']);
    
    // 获取特定状态的数量
    const getStatusCount = (status) => {
      return dashboardData.value.vulnerabilitiesByStatus?.[status] || 0;
    };

    // 获取高危和严重漏洞总数
    const getCriticalHighCount = () => {
      return (dashboardData.value.vulnerabilitiesBySeverity?.critical || 0) + 
             (dashboardData.value.vulnerabilitiesBySeverity?.high || 0);
    };
    
    // 图表配置
    const pieChartOptions = {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: {
          position: 'bottom',
          labels: {
            usePointStyle: true,
            padding: 20
          }
        },
        tooltip: {
          callbacks: {
            label: function(context) {
              const label = context.label || '';
              const value = context.raw || 0;
              const total = context.dataset.data.reduce((a, b) => a + b, 0);
              const percentage = Math.round((value / total) * 100);
              return `${label}: ${value} (${percentage}%)`;
            }
          }
        }
      }
    };
    
    const doughnutChartOptions = {
      responsive: true,
      maintainAspectRatio: false,
      cutout: '70%',
      plugins: {
        legend: {
          position: 'bottom',
          labels: {
            usePointStyle: true,
            padding: 20
          }
        },
        tooltip: {
          callbacks: {
            label: function(context) {
              const label = context.label || '';
              const value = context.raw || 0;
              const total = context.dataset.data.reduce((a, b) => a + b, 0);
              const percentage = Math.round((value / total) * 100);
              return `${label}: ${value} (${percentage}%)`;
            }
          }
        }
      }
    };
    
    const lineChartOptions = {
      responsive: true,
      maintainAspectRatio: false,
      scales: {
        y: {
          beginAtZero: true,
          ticks: {
            precision: 0
          }
        }
      },
      plugins: {
        legend: {
          position: 'bottom',
          labels: {
            usePointStyle: true,
            padding: 20
          }
        }
      }
    };

    // 漏洞状态分布图表数据
    const statusChartData = computed(() => {
      if (!dashboardData.value.vulnerabilitiesByStatus) return null;
      
      return {
        labels: ['未处理', '处理中', '已解决', '已关闭', '误报'],
        datasets: [
          {
            data: [
              dashboardData.value.vulnerabilitiesByStatus.open || 0,
              dashboardData.value.vulnerabilitiesByStatus.in_progress || 0,
              dashboardData.value.vulnerabilitiesByStatus.resolved || 0,
              dashboardData.value.vulnerabilitiesByStatus.closed || 0,
              dashboardData.value.vulnerabilitiesByStatus.false_positive || 0
            ],
            backgroundColor: ['#F44336', '#FF9800', '#4CAF50', '#2196F3', '#9E9E9E'],
            borderWidth: 0
          }
        ]
      };
    });

    // 漏洞严重程度分布图表数据
    const severityChartData = computed(() => {
      if (!dashboardData.value.vulnerabilitiesBySeverity) return null;
      
      return {
        labels: ['严重', '高危', '中危', '低危', '信息'],
        datasets: [
          {
            data: [
              dashboardData.value.vulnerabilitiesBySeverity.critical || 0,
              dashboardData.value.vulnerabilitiesBySeverity.high || 0,
              dashboardData.value.vulnerabilitiesBySeverity.medium || 0,
              dashboardData.value.vulnerabilitiesBySeverity.low || 0,
              dashboardData.value.vulnerabilitiesBySeverity.info || 0
            ],
            backgroundColor: ['#F44336', '#FF9800', '#FFEB3B', '#4CAF50', '#2196F3'],
            borderWidth: 0
          }
        ]
      };
    });

    // 月度漏洞趋势图表数据
    const monthlyChartData = computed(() => {
      if (!dashboardData.value.vulnerabilitiesByMonth || dashboardData.value.vulnerabilitiesByMonth.length === 0) {
        return null;
      }
      
      // 按日期排序
      const sortedData = [...dashboardData.value.vulnerabilitiesByMonth].sort((a, b) => {
        // 先按年份排序
        if (a.year !== b.year) {
          return a.year - b.year;
        }
        
        // 然后按月份名称排序
        const months = ['January', 'February', 'March', 'April', 'May', 'June', 
                      'July', 'August', 'September', 'October', 'November', 'December'];
        return months.indexOf(a.month) - months.indexOf(b.month);
      });
      
      return {
        labels: sortedData.map(item => `${item.month.substr(0, 3)} ${item.year}`),
        datasets: [
          {
            label: '漏洞数量',
            backgroundColor: 'rgba(33, 150, 243, 0.2)',
            borderColor: '#2196F3',
            pointBackgroundColor: '#2196F3',
            pointBorderColor: '#fff',
            pointHoverBackgroundColor: '#fff',
            pointHoverBorderColor: '#2196F3',
            data: sortedData.map(item => item.count),
            tension: 0.4,
            borderWidth: 3,
            fill: true
          }
        ]
      };
    });

    // 风险分数指示器
    const getRiskScoreColor = (score) => {
      if (score <= 20) return '#67C23A'; // 低风险
      if (score <= 40) return '#85CF4E'; // 较低风险
      if (score <= 60) return '#E6A23C'; // 中等风险
      if (score <= 80) return '#F56C6C'; // 高风险
      return '#F44336'; // 严重风险
    };

    const getRiskLevelText = (score) => {
      if (score <= 20) return '低风险';
      if (score <= 40) return '较低风险';
      if (score <= 60) return '中等风险';
      if (score <= 80) return '高风险';
      return '严重风险';
    };

    const getRiskLevelClass = (score) => {
      if (score <= 20) return 'risk-low';
      if (score <= 40) return 'risk-moderate-low';
      if (score <= 60) return 'risk-moderate';
      if (score <= 80) return 'risk-high';
      return 'risk-critical';
    };

    const getRiskDescription = (score) => {
      if (score <= 20) return '当前系统安全风险较低，继续保持良好的安全态势。';
      if (score <= 40) return '当前系统存在少量风险，建议关注并逐步解决。';
      if (score <= 60) return '当前系统风险处于中等水平，需要关注高危漏洞并及时修复。';
      if (score <= 80) return '当前系统风险较高，请优先处理高危和严重漏洞。';
      return '当前系统存在严重安全风险，请立即处理严重漏洞！';
    };

    const getProgressColor = (percentage) => {
      if (percentage >= 80) return '#67C23A';
      if (percentage >= 50) return '#E6A23C';
      return '#F56C6C';
    };

    const viewAllCritical = () => {
      router.push({
        path: '/vulnerabilities',
        query: {
          severity: 'critical,high',
          status: 'open,in_progress'
        }
      });
    };

    const formatDate = (dateString) => {
      if (!dateString) return '';
      const date = new Date(dateString);
      return date.toLocaleDateString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit'
      });
    };

    const getSeverityType = (severity) => {
      const typeMap = {
        critical: 'danger',
        high: 'warning',
        medium: '',
        low: 'success',
        info: 'info'
      };
      return typeMap[severity] || '';
    };

    const getSeverityText = (severity) => {
      const textMap = {
        critical: '严重',
        high: '高危',
        medium: '中危',
        low: '低危',
        info: '信息'
      };
      return textMap[severity] || severity;
    };

    const getStatusType = (status) => {
      const typeMap = {
        open: 'danger',
        in_progress: 'warning',
        resolved: 'success',
        closed: 'info',
        false_positive: ''
      };
      return typeMap[status] || '';
    };

    const getStatusText = (status) => {
      const textMap = {
        open: '未处理',
        in_progress: '处理中',
        resolved: '已解决',
        closed: '已关闭',
        false_positive: '误报'
      };
      return textMap[status] || status;
    };

    const getTeamProgressPercentage = (team) => {
      if (team.count === 0) return 100;
      return Math.round(((team.count - team.openCount) / team.count) * 100);
    };

    const getTeamProgressColor = (team) => {
      const percentage = getTeamProgressPercentage(team);
      if (percentage >= 80) return '#67C23A';
      if (percentage >= 50) return '#E6A23C';
      return '#F56C6C';
    };

    const trendsChartData = computed(() => {
      if (!vulnerabilityTrends.value) return null;
      
      return {
        labels: vulnerabilityTrends.value.timeLabels || [],
        datasets: [
          {
            label: '新增漏洞',
            data: vulnerabilityTrends.value.newVulnerabilities || [],
            backgroundColor: 'rgba(255, 152, 0, 0.2)',
            borderColor: '#FF9800',
            borderWidth: 2,
            tension: 0.4,
            fill: true
          },
          {
            label: '已修复漏洞',
            data: vulnerabilityTrends.value.resolvedVulnerabilities || [],
            backgroundColor: 'rgba(76, 175, 80, 0.2)',
            borderColor: '#4CAF50',
            borderWidth: 2,
            tension: 0.4,
            fill: true
          },
          {
            label: '净变化',
            data: vulnerabilityTrends.value.netChange || [],
            backgroundColor: 'rgba(33, 150, 243, 0.2)',
            borderColor: '#2196F3',
            borderWidth: 2,
            tension: 0.4,
            fill: false,
            borderDash: [5, 5]
          }
        ]
      };
    });

    const hasSeverityData = computed(() => {
      if (!dashboardData.value.vulnerabilitiesBySeverity) return false;
      return Object.values(dashboardData.value.vulnerabilitiesBySeverity).some(count => count > 0);
    });

    const hasStatusData = computed(() => {
      if (!dashboardData.value.vulnerabilitiesByStatus) return false;
      return Object.values(dashboardData.value.vulnerabilitiesByStatus).some(count => count > 0);
    });

    const hasMonthlyData = computed(() => {
      return dashboardData.value.vulnerabilitiesByMonth && 
             dashboardData.value.vulnerabilitiesByMonth.length > 1;
    });

    const hasTrendsData = computed(() => {
      if (!vulnerabilityTrends.value) return false;
      return (
        vulnerabilityTrends.value.timeLabels &&
        vulnerabilityTrends.value.timeLabels.length > 0 &&
        (vulnerabilityTrends.value.newVulnerabilities.some(v => v > 0) || vulnerabilityTrends.value.resolvedVulnerabilities.some(v => v > 0))
      );
    });

    onMounted(async () => {
      try {
        await store.dispatch('dashboard/fetchDashboardData');
      } catch (error) {
        console.error('加载仪表盘数据失败:', error);
      }
    });

    return {
      loading,
      dashboardData,
      statusChartData,
      severityChartData,
      monthlyChartData,
      riskScore,
      remediationProgress,
      teamVulnerabilities,
      criticalVulnerabilities,
      vulnerabilityTrends,
      hasSeverityData,
      hasStatusData,
      hasMonthlyData,
      hasTrendsData,
      getStatusCount,
      getCriticalHighCount,
      trendsChartData,
      pieChartOptions,
      doughnutChartOptions,
      lineChartOptions,
      getTeamProgressPercentage,
      getTeamProgressColor,
      getRiskScoreColor,
      getRiskLevelText,
      getRiskLevelClass,
      getRiskDescription,
      getProgressColor,
      viewAllCritical,
      formatDate,
      getSeverityType,
      getSeverityText,
      getStatusType,
      getStatusText
    };
  }
};
</script>

<style lang="scss" scoped>
.dashboard-container {
  padding: 24px;
  animation: fadeIn 0.6s ease;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

.dashboard-header {
  margin-bottom: 24px;
  
  .dashboard-title {
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
  
  .dashboard-subtitle {
    color: #606266;
    font-size: 16px;
    margin: 0;
  }
}

.el-row {
  margin-bottom: 24px;
  
  &:last-child {
    margin-bottom: 0;
  }
}

.stat-card {
  height: 120px;
  border-radius: 12px;
  padding: 16px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.05);
  transition: all 0.3s ease;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  margin-bottom: 24px;
  position: relative;
  overflow: hidden;
  
  &:hover {
    transform: translateY(-5px);
    box-shadow: 0 8px 20px rgba(0, 0, 0, 0.1);
  }
  
  &.red-gradient {
    background: linear-gradient(135deg, #f5515f, #ff7676);
    color: white;
  }
  
  &.orange-gradient {
    background: linear-gradient(135deg, #ff9966, #ff5e62);
    color: white;
  }
  
  &.green-gradient {
    background: linear-gradient(135deg, #43e97b, #38f9d7);
    color: white;
  }
  
  &.blue-gradient {
    background: linear-gradient(135deg, #2193b0, #6dd5ed);
    color: white;
  }
  
  .stat-card-content {
    display: flex;
    align-items: center;
    z-index: 2;
    position: relative;
  }
  
  .stat-icon-container {
    width: 48px;
    height: 48px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 12px;
    margin-right: 16px;
    background-color: rgba(255, 255, 255, 0.2);
    font-size: 24px;
  }
  
  .stat-info {
    display: flex;
    flex-direction: column;
  }
  
  .stat-value {
    font-size: 28px;
    font-weight: 700;
    line-height: 1.2;
  }
  
  .stat-label {
    font-size: 14px;
    font-weight: 500;
    opacity: 0.8;
  }
  
  .stat-progress {
    position: absolute;
    bottom: 0;
    left: 0;
    height: 4px;
    width: 100%;
    background: rgba(255, 255, 255, 0.2);
    
    &::after {
      content: '';
      position: absolute;
      top: 0;
      left: 0;
      height: 100%;
      width: 35%;
      background: rgba(255, 255, 255, 0.4);
      animation: progressAnim 2s ease-in-out infinite;
    }
  }
  
  @keyframes progressAnim {
    0% { width: 10%; }
    50% { width: 70%; }
    100% { width: 10%; }
  }
}

.chart-row {
  margin-bottom: 24px;
}

.chart-card, .recent-vulns-card, .systems-card {
  background: white;
  border-radius: 12px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.03);
  padding: 20px;
  height: 100%;
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
  margin-bottom: 24px;
  
  &:hover {
    box-shadow: 0 8px 25px rgba(0, 0, 0, 0.08);
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

.chart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  
  .chart-title {
    margin: 0;
    font-size: 16px;
    font-weight: 600;
    color: #303133;
    display: flex;
    align-items: center;
    
    .chart-icon {
      margin-right: 8px;
      color: #409EFF;
    }
  }
}

.chart-container {
  height: 300px;
  position: relative;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  
  .card-title-wrapper {
    display: flex;
    flex-direction: column;
  }
  
  .card-title {
    margin: 0;
    font-size: 16px;
    font-weight: 600;
    color: #303133;
    display: flex;
    align-items: center;
    
    .chart-icon {
      margin-right: 8px;
      color: #409EFF;
    }
  }
  
  .card-subtitle {
    font-size: 12px;
    color: #909399;
    margin-top: 4px;
  }
  
  .view-all-btn {
    font-size: 12px;
    font-weight: 500;
    
    i {
      margin-left: 4px;
    }
  }
}

.vuln-title-link {
  color: #409EFF;
  text-decoration: none;
  transition: color 0.3s ease;
  font-weight: 500;
  
  &:hover {
    color: #66b1ff;
    text-decoration: underline;
  }
}

.date-text {
  display: flex;
  align-items: center;
  color: #909399;
  font-size: 13px;
  
  .date-icon {
    margin-right: 4px;
  }
}

.systems-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
  
  .system-item {
    display: flex;
    align-items: center;
    
    .system-icon {
      width: 40px;
      height: 40px;
      border-radius: 8px;
      display: flex;
      align-items: center;
      justify-content: center;
      margin-right: 12px;
      font-size: 18px;
    }
    
    .system-details {
      flex: 1;
    }
    
    .system-info {
      display: flex;
      justify-content: space-between;
      margin-bottom: 8px;
    }
    
    .system-name {
      font-weight: 500;
      color: #303133;
    }
    
    .system-count {
      font-size: 12px;
      color: #606266;
      font-weight: 500;
    }
  }
}

.no-data {
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  color: #909399;
  
  i {
    font-size: 42px;
    margin-bottom: 16px;
    opacity: 0.5;
  }
  
  p {
    font-size: 16px;
    margin: 0 0 16px 0;
  }
  
  .no-data-tip {
    font-size: 12px;
    opacity: 0.8;
    text-align: center;
    max-width: 200px;
    margin-bottom: 16px;
  }
}

.risk-score-container {
  padding: 20px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  
  .risk-score-gauge {
    display: flex;
    flex-direction: column;
    align-items: center;
    width: 100%;
    margin-bottom: 16px;
    
    .gauge-container {
      position: relative;
      width: 150px;
      height: 150px;
      margin-bottom: 16px;
      
      .gauge {
        width: 100%;
        height: 100%;
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
        position: relative;
        
        &:before {
          content: '';
          position: absolute;
          width: 80%;
          height: 80%;
          background: white;
          border-radius: 50%;
        }
        
        .gauge-center {
          position: relative;
          z-index: 1;
          text-align: center;
          
          .gauge-value {
            font-size: 32px;
            font-weight: 700;
            color: #333;
          }
          
          .gauge-label {
            font-size: 14px;
            color: #666;
          }
        }
      }
    }
    
    .risk-level {
      display: flex;
      align-items: center;
      justify-content: center;
      
      .risk-level-title {
        font-size: 14px;
        color: #666;
        margin-right: 8px;
      }
      
      .risk-level-value {
        font-size: 16px;
        font-weight: 600;
        
        &.risk-low {
          color: #67C23A;
        }
        
        &.risk-moderate-low {
          color: #85CF4E;
        }
        
        &.risk-moderate {
          color: #E6A23C;
        }
        
        &.risk-high {
          color: #F56C6C;
        }
        
        &.risk-critical {
          color: #F44336;
        }
      }
    }
  }
  
  .risk-score-description {
    text-align: center;
    
    p {
      color: #606266;
      font-size: 14px;
      line-height: 1.5;
      margin: 0;
    }
  }
}

.remediation-container {
  padding: 20px;
  
  .remediation-stats {
    display: flex;
    justify-content: space-around;
    margin-bottom: 24px;
    
    .remediation-stat-item {
      text-align: center;
      
      .stat-value {
        font-size: 24px;
        font-weight: 600;
        color: #333;
      }
      
      .stat-label {
        font-size: 14px;
        color: #666;
        margin-top: 5px;
      }
    }
  }
  
  .progress-bar-container {
    .progress-label {
      font-size: 14px;
      color: #666;
      margin-bottom: 8px;
    }
    
    .progress-value {
      font-size: 14px;
      color: #333;
      font-weight: 600;
      margin-top: 8px;
      text-align: right;
    }
  }
}

.team-container {
  padding: 20px;
  
  .team-list {
    .team-item {
      padding: 16px;
      border-radius: 8px;
      background-color: #f9f9f9;
      margin-bottom: 12px;
      
      &:last-child {
        margin-bottom: 0;
      }
      
      .team-info {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 12px;
        
        .team-name {
          font-size: 16px;
          font-weight: 600;
          color: #333;
        }
        
        .team-count {
          .count-value {
            font-size: 16px;
            font-weight: 600;
            color: #409EFF;
          }
          
          .count-label {
            font-size: 12px;
            color: #666;
            margin-left: 4px;
          }
        }
      }
      
      .team-progress {
        margin-bottom: 12px;
      }
      
      .team-severity {
        .severity-dots {
          display: flex;
          
          .severity-dot {
            width: 12px;
            height: 12px;
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
          }
        }
      }
    }
  }
}

.trend-container {
  padding: 20px;
  height: 350px;
}

.critical-container {
  padding: 0;
  
  .critical-table {
    :deep(th) {
      background-color: #f5f7fa;
      color: #606266;
      font-weight: 600;
    }
    
    :deep(tr:hover) {
      background-color: #f5f7fc;
    }
    
    .vulnerability-link {
      color: #409EFF;
      text-decoration: none;
      
      &:hover {
        text-decoration: underline;
      }
    }
    
    .cvss-score {
      font-weight: 600;
      padding: 2px 6px;
      border-radius: 4px;
      background-color: #f0f9eb;
      color: #67C23A;
    }
  }
}

.chart-actions {
  display: flex;
  align-items: center;
}

// 响应式调整
@media (max-width: 768px) {
  .dashboard-container {
    padding: 16px;
  }
  
  .dashboard-header {
    .dashboard-title {
      font-size: 24px;
      
      .title-icon {
        font-size: 24px;
      }
    }
    
    .dashboard-subtitle {
      font-size: 14px;
    }
  }
  
  .stat-card {
    height: 100px;
    
    .stat-icon-container {
      width: 36px;
      height: 36px;
      font-size: 18px;
    }
    
    .stat-value {
      font-size: 22px;
    }
    
    .stat-label {
      font-size: 12px;
    }
  }
  
  .chart-container {
    height: 220px;
  }
}
</style> 