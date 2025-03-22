<template>
  <div class="dashboard-container">
    <!-- 顶部统计卡片 -->
    <el-row :gutter="20" v-loading="statsLoading">
      <el-col :xs="24" :sm="12" :md="6" :lg="6" :xl="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-header">
            <div class="stat-value">{{ stats.totalVulns || 0 }}</div>
            <el-icon class="stat-icon"><Monitor /></el-icon>
          </div>
          <div class="stat-title">总漏洞数量</div>
          <div class="stat-footer" v-if="stats.vulnChangeRate !== undefined">
            <el-tag :type="stats.vulnChangeRate > 0 ? 'danger' : 'success'" size="small">
              {{ stats.vulnChangeRate > 0 ? '+' : '' }}{{ stats.vulnChangeRate }}%
            </el-tag>
            <span class="stat-period">与上月相比</span>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12" :md="6" :lg="6" :xl="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-header">
            <div class="stat-value">{{ stats.criticalVulns || 0 }}</div>
            <el-icon class="stat-icon critical"><CircleCloseFilled /></el-icon>
          </div>
          <div class="stat-title">高危漏洞数量</div>
          <div class="stat-footer" v-if="stats.criticalChangeRate !== undefined">
            <el-tag :type="stats.criticalChangeRate > 0 ? 'danger' : 'success'" size="small">
              {{ stats.criticalChangeRate > 0 ? '+' : '' }}{{ stats.criticalChangeRate }}%
            </el-tag>
            <span class="stat-period">与上月相比</span>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12" :md="6" :lg="6" :xl="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-header">
            <div class="stat-value">{{ stats.totalAssets || 0 }}</div>
            <el-icon class="stat-icon"><Monitor /></el-icon>
          </div>
          <div class="stat-title">总资产数量</div>
          <div class="stat-footer" v-if="stats.assetChangeRate !== undefined">
            <el-tag :type="stats.assetChangeRate > 0 ? 'success' : 'info'" size="small">
              {{ stats.assetChangeRate > 0 ? '+' : '' }}{{ stats.assetChangeRate }}%
            </el-tag>
            <span class="stat-period">与上月相比</span>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12" :md="6" :lg="6" :xl="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-header">
            <div class="stat-value">{{ stats.fixedRate || 0 }}%</div>
            <el-icon class="stat-icon success"><CircleCheckFilled /></el-icon>
          </div>
          <div class="stat-title">已修复漏洞比例</div>
          <div class="stat-footer" v-if="stats.fixedChangeRate !== undefined">
            <el-tag :type="stats.fixedChangeRate > 0 ? 'success' : 'danger'" size="small">
              {{ stats.fixedChangeRate > 0 ? '+' : '' }}{{ stats.fixedChangeRate }}%
            </el-tag>
            <span class="stat-period">与上月相比</span>
          </div>
        </el-card>
      </el-col>

      <!-- 刷新控制 -->
      <el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mt-4">
        <el-card shadow="hover" class="refresh-card">
          <div class="refresh-controls">
            <span>自动刷新: </span>
            <el-switch v-model="autoRefresh" @change="toggleAutoRefresh" />
            <span class="refresh-interval" v-if="autoRefresh">
              <span>刷新间隔: </span>
              <el-select v-model="refreshInterval" size="small" @change="startAutoRefresh">
                <el-option label="30秒" :value="30" />
                <el-option label="60秒" :value="60" />
                <el-option label="5分钟" :value="300" />
                <el-option label="10分钟" :value="600" />
              </el-select>
            </span>
            <el-button type="primary" size="small" :icon="RefreshRight" @click="initDashboard">
              立即刷新
            </el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="mt-4">
      <!-- 漏洞趋势图 -->
      <el-col :xs="24" :sm="24" :md="12" :lg="12" :xl="12">
        <el-card shadow="hover" class="chart-card">
          <template #header>
            <div class="card-header">
              <span>漏洞趋势</span>
              <div class="card-header-actions">
                <el-radio-group v-model="trendPeriod" size="small" @change="fetchTrendsData">
                  <el-radio-button label="7d">7天</el-radio-button>
                  <el-radio-button label="1m">1个月</el-radio-button>
                  <el-radio-button label="3m">3个月</el-radio-button>
                </el-radio-group>
              </div>
            </div>
          </template>
          <div class="chart-container" v-loading="trendChartLoading">
            <div id="vulnTrendChart" style="height: 300px"></div>
          </div>
        </el-card>
      </el-col>

      <!-- 漏洞严重程度分布 -->
      <el-col :xs="24" :sm="24" :md="12" :lg="12" :xl="12">
        <el-card shadow="hover" class="chart-card">
          <template #header>
            <div class="card-header">
              <span>漏洞严重程度分布</span>
            </div>
          </template>
          <div class="chart-container" v-loading="severityChartLoading">
            <div id="severityChart" style="height: 300px"></div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="mt-4">
      <!-- 资产漏洞分布 -->
      <el-col :xs="24" :sm="24" :md="12" :lg="12" :xl="12">
        <el-card shadow="hover" class="chart-card">
          <template #header>
            <div class="card-header">
              <span>资产漏洞分布</span>
              <el-tooltip content="资产漏洞分布说明" placement="top">
                <el-icon><InfoFilled /></el-icon>
              </el-tooltip>
            </div>
          </template>
          <div class="chart-container" v-loading="assetChartLoading">
            <div id="assetVulnChart" style="height: 300px"></div>
          </div>
        </el-card>
      </el-col>

      <!-- 最近活动 -->
      <el-col :xs="24" :sm="24" :md="12" :lg="12" :xl="12">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>最近活动</span>
            </div>
          </template>
          <div class="activities-container" v-loading="activitiesLoading">
            <div v-if="activities.length === 0" class="empty-data">
              <el-empty :description="无数据" :image-size="80" />
            </div>
            <el-timeline v-else>
              <el-timeline-item
                v-for="(activity, index) in activities"
                :key="index"
                :type="getActivityType(activity.type)"
                :timestamp="formatTime(activity.time)"
                :icon="getActivityIcon(activity.type)"
                placement="top"
              >
                <h4 class="activity-title">{{ getActivityTypeLabel(activity.type) }}</h4>
                <div class="activity-content" :data-type="activity.type" v-html="activity.content"></div>
                <div class="activity-user">
                  <el-icon><UserFilled /></el-icon>
                  <span>{{ activity.username }}</span>
                </div>
              </el-timeline-item>
            </el-timeline>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="mt-4">
      <!-- 优先修复漏洞 -->
      <el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>优先修复漏洞</span>
              <div class="card-header-actions">
                <el-button 
                  type="primary" 
                  size="small" 
                  @click="handleMoreVulns"
                  class="action-btn"
                >
                  <el-icon><ArrowRight /></el-icon>
                  <span>查看更多</span>
                </el-button>
              </div>
            </div>
          </template>
          <div v-if="priorityVulns.length === 0 && !priorityVulnsLoading" class="empty-data">
            <el-empty :description="无数据" :image-size="80" />
          </div>
          <el-table :data="priorityVulns" style="width: 100%" v-loading="priorityVulnsLoading">
            <el-table-column prop="title" :label="标题" min-width="200" />
            <el-table-column prop="severity" :label="严重程度" width="100">
              <template #default="scope">
                <el-tag
                  :type="getSeverityType(scope.row.severity)"
                  effect="dark"
                >
                  {{ getSeverityLabel(scope.row.severity) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="status" :label="状态" width="120">
              <template #default="scope">
                <el-tag
                  :type="getStatusType(scope.row.status)"
                >
                  {{ getStatusLabel(scope.row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="asset" :label="资产" min-width="180" />
            <el-table-column prop="discoveredAt" :label="发现时间" width="160">
              <template #default="scope">
                {{ formatDate(scope.row.discoveredAt) }}
              </template>
            </el-table-column>
            <el-table-column :label="操作" width="150" fixed="right">
              <template #default="scope">
                <div class="operation-buttons">
                  <el-button 
                    type="primary" 
                    size="small"
                    @click="handleViewDetail(scope.row.id)"
                    class="action-btn"
                  >
                    <el-icon><View /></el-icon>
                    <span>查看</span>
                  </el-button>
                </div>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { reactive, ref, onMounted, onBeforeUnmount, toRefs } from 'vue'
import { useRouter } from 'vue-router'
import { useStore } from 'vuex'
import { ElMessage } from 'element-plus'
import {
  WarningFilled,
  CircleCloseFilled,
  Monitor,
  CircleCheckFilled,
  InfoFilled,
  RefreshRight,
  ArrowRight,
  View,
  UserFilled,
  Reading,
  MoreFilled
} from '@element-plus/icons-vue'
import * as echarts from 'echarts/core'
import { BarChart, LineChart, PieChart } from 'echarts/charts'
import {
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent
} from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import {
  getDashboardStats,
  getVulnTrends,
  getSeverityDistribution,
  getPriorityVulns,
  getAssetVulnDistribution,
  getRecentActivities
} from '@/api/dashboard'
import { formatTime, formatRelativeTime } from '@/utils/time'

// 初始化echarts组件
echarts.use([
  BarChart,
  LineChart,
  PieChart,
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent,
  CanvasRenderer
])

const router = useRouter()
const store = useStore()

// 数据加载状态
const statsLoading = ref(false)
const trendChartLoading = ref(false)
const severityChartLoading = ref(false)
const assetChartLoading = ref(false)
const priorityVulnsLoading = ref(false)
const activitiesLoading = ref(false)

// 图表实例
let trendChart = null
let severityChart = null
let assetVulnChart = null

// 趋势周期选择
const trendPeriod = ref('7d')

// 自动刷新设置
const autoRefresh = ref(true)
const refreshInterval = ref(60) // 默认60秒刷新一次
let refreshTimer = null

// 切换自动刷新
const toggleAutoRefresh = () => {
  autoRefresh.value = !autoRefresh.value
  if (autoRefresh.value) {
    startAutoRefresh()
  } else {
    stopAutoRefresh()
  }
}

// 开始自动刷新
const startAutoRefresh = () => {
  stopAutoRefresh() // 先清除之前的定时器
  
  // 设置定时刷新
  refreshTimer = setInterval(() => {
    if (autoRefresh.value) {
      fetchStats()
      fetchTrendsData()
      fetchSeverityDistribution()
      fetchAssetVulnDistribution()
      fetchPriorityVulns()
      fetchRecentActivities()
    }
  }, refreshInterval.value * 1000)
}

// 停止自动刷新
const stopAutoRefresh = () => {
  if (refreshTimer) {
    clearInterval(refreshTimer)
    refreshTimer = null
  }
}

// 数据
const state = reactive({
  stats: {
    totalVulns: 0,
    criticalVulns: 0,
    totalAssets: 0,
    fixedRate: 0,
    vulnChangeRate: 0,
    criticalChangeRate: 0,
    assetChangeRate: 0,
    fixedChangeRate: 0
  },
  priorityVulns: [],
  activities: []
})

// 使用toRefs从state中提取响应式属性
const { stats, priorityVulns, activities } = toRefs(state)

// Mock数据生成函数
const getMockTrendData = () => {
  const dates = []
  const newVulns = []
  const fixedVulns = []
  
  // 生成最近30天的日期
  for (let i = 29; i >= 0; i--) {
    const date = new Date()
    date.setDate(date.getDate() - i)
    dates.push(date.toLocaleDateString())
    newVulns.push(Math.floor(Math.random() * 8) + 1) // 1-8的随机数
    fixedVulns.push(Math.floor(Math.random() * 6) + 1) // 1-6的随机数
  }
  
  return { dates, newVulns, fixedVulns }
}

const getMockSeverityData = () => {
  return {
    critical: 18,
    high: 35,
    medium: 42,
    low: 24,
    info: 9
  }
}

const getMockAssetData = () => {
  return [
    { name: 'Web Server', count: 15 },
    { name: 'Database Server', count: 12 },
    { name: 'API Gateway', count: 9 },
    { name: 'Mobile App Backend', count: 7 },
    { name: 'Internal Portal', count: 5 }
  ]
}

const getMockPriorityVulns = () => {
  return [
    {
      id: 1,
      title: 'SQL Injection in Auth Module',
      severity: 'critical',
      status: 'open',
      asset: 'Web Application Server',
      discoveredAt: new Date().toISOString()
    },
    {
      id: 2,
      title: 'Cross-Site Scripting in User Profile',
      severity: 'high',
      status: 'in_progress',
      asset: 'Customer Portal',
      discoveredAt: new Date(Date.now() - 2 * 24 * 60 * 60 * 1000).toISOString()
    },
    {
      id: 3,
      title: 'Insecure Direct Object Reference',
      severity: 'medium',
      status: 'open',
      asset: 'API Gateway',
      discoveredAt: new Date(Date.now() - 5 * 24 * 60 * 60 * 1000).toISOString()
    },
    {
      id: 4,
      title: 'Outdated TLS Version',
      severity: 'high',
      status: 'pending_review',
      asset: 'Load Balancer',
      discoveredAt: new Date(Date.now() - 3 * 24 * 60 * 60 * 1000).toISOString()
    },
    {
      id: 5,
      title: 'Default Admin Credentials',
      severity: 'critical',
      status: 'open',
      asset: 'Database Server',
      discoveredAt: new Date(Date.now() - 1 * 24 * 60 * 60 * 1000).toISOString()
    }
  ]
}

const getMockActivities = () => {
  return [
    {
      type: 'vuln_create',
      content: '发现新漏洞：SQL注入漏洞 (CRITICAL)',
      username: 'admin',
      time: new Date().toLocaleString()
    },
    {
      type: 'status_change',
      content: '漏洞状态更新：XSS漏洞 -> 修复中',
      username: 'security',
      time: new Date(Date.now() - 1 * 60 * 60 * 1000).toLocaleString()
    },
    {
      type: 'comment',
      content: '添加评论：已确认此漏洞影响范围',
      username: 'tester',
      time: new Date(Date.now() - 3 * 60 * 60 * 1000).toLocaleString()
    },
    {
      type: 'asset_create',
      content: '添加新资产：测试服务器',
      username: 'admin',
      time: new Date(Date.now() - 5 * 60 * 60 * 1000).toLocaleString()
    },
    {
      type: 'fix',
      content: '漏洞已修复：目录遍历漏洞',
      username: 'developer',
      time: new Date(Date.now() - 8 * 60 * 60 * 1000).toLocaleString()
    }
  ]
}

// 获取统计数据
const fetchStats = async () => {
  statsLoading.value = true
  try {
    const response = await getDashboardStats()
    if (response && response.code === 200 && response.data) {
      state.stats = response.data
    } else {
      // 显示错误消息
      ElMessage.error('获取仪表盘统计数据失败')
      // 重置为默认值而不是使用mock数据
      state.stats = {
        totalVulns: 0,
        criticalVulns: 0,
        totalAssets: 0,
        fixedRate: 0,
        vulnChangeRate: 0,
        criticalChangeRate: 0,
        assetChangeRate: 0,
        fixedChangeRate: 0
      }
    }
  } catch (error) {
    console.error('获取仪表盘统计数据失败:', error)
    // 显示错误消息
    ElMessage.error('获取仪表盘统计数据失败')
    // 重置为默认值而不是使用mock数据
    state.stats = {
      totalVulns: 0,
      criticalVulns: 0,
      totalAssets: 0,
      fixedRate: 0,
      vulnChangeRate: 0,
      criticalChangeRate: 0,
      assetChangeRate: 0,
      fixedChangeRate: 0
    }
  } finally {
    statsLoading.value = false
  }
}

// 获取漏洞趋势数据
const fetchTrendsData = async () => {
  trendChartLoading.value = true
  try {
    // 将前端周期值映射到后端期望的格式
    const periodMap = {
      '7d': 'week',
      '1m': 'month',
      '3m': 'quarter'
    }
    const backendPeriod = periodMap[trendPeriod.value] || 'month'
    
    const response = await getVulnTrends({ period: backendPeriod })
    if (response && response.code === 200 && response.data) {
      updateTrendChart(response.data)
    } else {
      // 显示错误消息
      ElMessage.error('获取漏洞趋势数据失败')
      // 使用空数据而不是mock数据
      updateTrendChart(null)
    }
  } catch (error) {
    console.error('获取漏洞趋势数据失败:', error)
    // 显示错误消息
    ElMessage.error('获取漏洞趋势数据失败')
    // 使用空数据而不是mock数据
    updateTrendChart(null)
  } finally {
    trendChartLoading.value = false
  }
}

// 获取严重程度分布
const fetchSeverityDistribution = async () => {
  severityChartLoading.value = true
  try {
    const response = await getSeverityDistribution()
    if (response && response.code === 200 && response.data) {
      // 将后端返回的数组格式转换为前端所需的对象格式
      const formattedData = {}
      response.data.forEach(item => {
        formattedData[item.severity] = item.count
      })
      
      // 使用转换后的数据更新图表
      updateSeverityChart(formattedData)
    } else {
      // 显示错误消息
      ElMessage.error('获取漏洞严重程度分布失败')
      // 使用空数据而不是mock数据
      updateSeverityChart(null)
    }
  } catch (error) {
    console.error('获取漏洞严重程度分布失败:', error)
    // 显示错误消息
    ElMessage.error('获取漏洞严重程度分布失败')
    // 使用空数据而不是mock数据
    updateSeverityChart(null)
  } finally {
    severityChartLoading.value = false
  }
}

// 获取资产漏洞分布
const fetchAssetVulnDistribution = async () => {
  assetChartLoading.value = true
  try {
    const response = await getAssetVulnDistribution()
    if (response && response.code === 200 && response.data) {
      updateAssetVulnChart(response.data)
    } else {
      // 显示错误消息
      ElMessage.error('获取资产漏洞分布失败')
      // 使用空数据而不是mock数据
      updateAssetVulnChart(null)
    }
  } catch (error) {
    console.error('获取资产漏洞分布失败:', error)
    // 显示错误消息
    ElMessage.error('获取资产漏洞分布失败')
    // 使用空数据而不是mock数据
    updateAssetVulnChart(null)
  } finally {
    assetChartLoading.value = false
  }
}

// 获取优先修复漏洞
const fetchPriorityVulns = async () => {
  priorityVulnsLoading.value = true
  try {
    const response = await getPriorityVulns({ limit: 5 })
    if (response && response.code === 200 && response.data) {
      state.priorityVulns = response.data
    } else {
      // 显示错误消息
      ElMessage.error('获取优先修复漏洞失败')
      // 使用空数组而不是mock数据
      state.priorityVulns = []
    }
  } catch (error) {
    console.error('获取优先修复漏洞失败:', error)
    // 显示错误消息
    ElMessage.error('获取优先修复漏洞失败')
    // 使用空数组而不是mock数据
    state.priorityVulns = []
  } finally {
    priorityVulnsLoading.value = false
  }
}

// 获取最近活动
const fetchRecentActivities = async () => {
  activitiesLoading.value = true
  try {
    const response = await getRecentActivities(5)
    if (response && response.code === 200 && response.data) {
      state.activities = response.data
    } else {
      // 显示错误消息
      ElMessage.error('获取最近活动失败')
      // 使用空数组而不是mock数据
      state.activities = []
    }
  } catch (error) {
    console.error('获取最近活动失败:', error)
    // 显示错误消息
    ElMessage.error('获取最近活动失败')
    // 使用空数组而不是mock数据
    state.activities = []
  } finally {
    activitiesLoading.value = false
  }
}

// 更新趋势图
const updateTrendChart = (data) => {
  if (!document.getElementById('vulnTrendChart')) {
    return
  }
  
  if (!trendChart) {
    trendChart = echarts.init(document.getElementById('vulnTrendChart'))
  }
  
  // 如果数据为空或未定义，设置空图表
  if (!data || !data.dates || !data.newVulns || !data.fixedVulns || 
      !data.dates.length || !data.newVulns.length || !data.fixedVulns.length) {
    // 检查是否是API返回的不同格式
    if (data && data.dates && data.new && data.fixed && 
        data.dates.length && data.new.length && data.fixed.length) {
      // 将后端API返回的格式转换为前端期望的格式
      data = {
        dates: data.dates,
        newVulns: data.new,
        fixedVulns: data.fixed
      }
    } else {
      // 设置空图表
      trendChart.setOption({
        tooltip: {
          trigger: 'axis'
        },
        legend: {
          data: ['新增漏洞', '已修复漏洞']
        },
        grid: {
          left: '3%',
          right: '4%',
          bottom: '3%',
          containLabel: true
        },
        xAxis: {
          type: 'category',
          boundaryGap: false,
          data: ['']
        },
        yAxis: {
          type: 'value'
        },
        series: [
          {
            name: '新增漏洞',
            type: 'line',
            data: [],
            itemStyle: {
              color: '#F56C6C'
            },
            areaStyle: {
              color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
                { offset: 0, color: 'rgba(245, 108, 108, 0.5)' },
                { offset: 1, color: 'rgba(245, 108, 108, 0.1)' }
              ])
            },
            smooth: true
          },
          {
            name: '已修复漏洞',
            type: 'line',
            data: [],
            itemStyle: {
              color: '#67C23A'
            },
            areaStyle: {
              color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
                { offset: 0, color: 'rgba(103, 194, 58, 0.5)' },
                { offset: 1, color: 'rgba(103, 194, 58, 0.1)' }
              ])
            },
            smooth: true
          }
        ]
      })
      return
    }
  }
  
  trendChart.setOption({
    tooltip: {
      trigger: 'axis',
      formatter: function(params) {
        const date = params[0].axisValue;
        let html = `<div style="font-weight:bold;margin-bottom:5px;">${date}</div>`;
        params.forEach(item => {
          html += `<div style="display:flex;align-items:center;margin:3px 0;">
            <span style="display:inline-block;width:10px;height:10px;background:${item.color};margin-right:5px;border-radius:50%;"></span>
            <span>${item.seriesName}: </span>
            <span style="font-weight:bold;margin-left:5px;">${item.value}</span>
          </div>`;
        });
        return html;
      }
    },
    legend: {
      data: ['新增漏洞', '已修复漏洞']
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: data.dates,
      axisLabel: {
        formatter: function(value) {
          // 针对7天视图，显示星期几
          if (trendPeriod.value === '7d') {
            const dateObj = new Date(new Date().getFullYear() + '-' + value);
            const days = ['周日', '周一', '周二', '周三', '周四', '周五', '周六'];
            return value + '\n' + days[dateObj.getDay()];
          }
          return value;
        }
      }
    },
    yAxis: {
      type: 'value',
      minInterval: 1 // 确保Y轴刻度为整数
    },
    series: [
      {
        name: '新增漏洞',
        type: 'line',
        data: data.newVulns,
        itemStyle: {
          color: '#F56C6C'
        },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(245, 108, 108, 0.5)' },
            { offset: 1, color: 'rgba(245, 108, 108, 0.1)' }
          ])
        },
        smooth: true,
        label: {
          show: trendPeriod.value === '7d', // 只在7天视图显示标签
          position: 'top',
          formatter: '{c}'
        }
      },
      {
        name: '已修复漏洞',
        type: 'line',
        data: data.fixedVulns,
        itemStyle: {
          color: '#67C23A'
        },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(103, 194, 58, 0.5)' },
            { offset: 1, color: 'rgba(103, 194, 58, 0.1)' }
          ])
        },
        smooth: true,
        label: {
          show: trendPeriod.value === '7d', // 只在7天视图显示标签
          position: 'top',
          formatter: '{c}'
        }
      }
    ]
  })
}

// 更新严重程度分布图
const updateSeverityChart = (data) => {
  if (!document.getElementById('severityChart')) {
    return
  }
  
  if (!severityChart) {
    severityChart = echarts.init(document.getElementById('severityChart'))
  }
  
  const severityMap = {
    critical: '严重',
    high: '高危',
    medium: '中危',
    low: '低危',
    info: '信息'
  }
  
  const colorMap = {
    critical: '#F56C6C',
    high: '#E6A23C',
    medium: '#409EFF',
    low: '#67C23A',
    info: '#909399'
  }

  // 如果数据为空或未定义，设置空图表
  if (!data || Object.keys(data).length === 0) {
    const allSeverities = ['critical', 'high', 'medium', 'low', 'info']
    severityChart.setOption({
      tooltip: {
        trigger: 'item',
        formatter: function(params) {
          const percent = params.percent.toFixed(1);
          return `<div style="font-weight:bold;margin-bottom:5px;">${params.name}</div>
                  <div style="display:flex;align-items:center;margin:3px 0;">
                    <span style="display:inline-block;width:10px;height:10px;background:${params.color};margin-right:5px;border-radius:50%;"></span>
                    <span>数量: </span>
                    <span style="font-weight:bold;margin-left:5px;">${params.value}</span>
                  </div>
                  <div style="margin-top:5px;">占比: ${percent}%</div>`;
        }
      },
      legend: {
        orient: 'vertical',
        left: 10,
        data: allSeverities.map(key => severityMap[key]),
        formatter: function(name) {
          // 查找对应的原始严重程度级别
          const severity = Object.keys(severityMap).find(key => severityMap[key] === name);
          if (severity && data[severity] !== undefined) {
            return `${name}: ${data[severity]}`;
          }
          return name;
        }
      },
      series: [
        {
          name: '漏洞严重程度',
          type: 'pie',
          radius: ['50%', '70%'],
          avoidLabelOverlap: false,
          itemStyle: {
            borderRadius: 10,
            borderColor: '#fff',
            borderWidth: 2
          },
          label: {
            show: true,
            position: 'outside',
            formatter: '{b}: {c}',
            fontSize: 12
          },
          emphasis: {
            label: {
              show: true,
              fontSize: 14,
              fontWeight: 'bold'
            },
            itemStyle: {
              shadowBlur: 10,
              shadowOffsetX: 0,
              shadowColor: 'rgba(0, 0, 0, 0.5)'
            }
          },
          labelLine: {
            show: true,
            length: 10,
            length2: 10
          },
          data: allSeverities.map(key => ({
            value: data[key],
            name: severityMap[key],
            itemStyle: {
              color: colorMap[key]
            }
          }))
        }
      ]
    })
    return
  }
  
  // 检查数据类型，如果是数组则转换为对象
  if (Array.isArray(data)) {
    const tempData = {}
    data.forEach(item => {
      if (item.severity && typeof item.count === 'number') {
        tempData[item.severity] = item.count
      }
    })
    data = tempData
  }
  
  // 确保所有严重程度类型都有数据，即使值为0
  const allSeverities = ['critical', 'high', 'medium', 'low', 'info']
  allSeverities.forEach(severity => {
    if (data[severity] === undefined) {
      data[severity] = 0
    }
  })
  
  severityChart.setOption({
    tooltip: {
      trigger: 'item',
      formatter: function(params) {
        const percent = params.percent.toFixed(1);
        return `<div style="font-weight:bold;margin-bottom:5px;">${params.name}</div>
                <div style="display:flex;align-items:center;margin:3px 0;">
                  <span style="display:inline-block;width:10px;height:10px;background:${params.color};margin-right:5px;border-radius:50%;"></span>
                  <span>数量: </span>
                  <span style="font-weight:bold;margin-left:5px;">${params.value}</span>
                </div>
                <div style="margin-top:5px;">占比: ${percent}%</div>`;
      }
    },
    legend: {
      orient: 'vertical',
      left: 10,
      data: allSeverities.map(key => severityMap[key]),
      formatter: function(name) {
        // 查找对应的原始严重程度级别
        const severity = Object.keys(severityMap).find(key => severityMap[key] === name);
        if (severity && data[severity] !== undefined) {
          return `${name}: ${data[severity]}`;
        }
        return name;
      }
    },
    series: [
      {
        name: '漏洞严重程度',
        type: 'pie',
        radius: ['50%', '70%'],
        avoidLabelOverlap: false,
        itemStyle: {
          borderRadius: 10,
          borderColor: '#fff',
          borderWidth: 2
        },
        label: {
          show: true,
          position: 'outside',
          formatter: '{b}: {c}',
          fontSize: 12
        },
        emphasis: {
          label: {
            show: true,
            fontSize: 14,
            fontWeight: 'bold'
          },
          itemStyle: {
            shadowBlur: 10,
            shadowOffsetX: 0,
            shadowColor: 'rgba(0, 0, 0, 0.5)'
          }
        },
        labelLine: {
          show: true,
          length: 10,
          length2: 10
        },
        data: allSeverities.map(key => ({
          value: data[key],
          name: severityMap[key],
          itemStyle: {
            color: colorMap[key]
          }
        }))
      }
    ]
  })
}

// 更新资产漏洞分布图
const updateAssetVulnChart = (data) => {
  if (!document.getElementById('assetVulnChart')) {
    return
  }
  
  if (!assetVulnChart) {
    assetVulnChart = echarts.init(document.getElementById('assetVulnChart'))
  }
  
  // 如果没有数据或者数据为空数组，使用空数据
  if (!data || !Array.isArray(data) || data.length === 0) {
    // 设置空图表
    assetVulnChart.setOption({
      tooltip: {
        trigger: 'axis',
        axisPointer: {
          type: 'shadow'
        }
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true
      },
      xAxis: {
        type: 'value',
        name: '漏洞数量'
      },
      yAxis: {
        type: 'category',
        data: [],
        axisLabel: {
          width: 120,
          overflow: 'truncate'
        }
      },
      series: [
        {
          name: '漏洞数量',
          type: 'bar',
          data: [],
          label: {
            show: true,
            position: 'right'
          }
        }
      ]
    })
    return
  }
  
  // 限制最多显示10个资产
  if (data.length > 10) {
    data = data.slice(0, 10)
  }
  
  assetVulnChart.setOption({
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      }
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'value',
      name: '漏洞数量'
    },
    yAxis: {
      type: 'category',
      data: data.map(item => item.name),
      axisLabel: {
        width: 120,
        overflow: 'truncate'
      }
    },
    series: [
      {
        name: '漏洞数量',
        type: 'bar',
        data: data.map(item => ({
          value: item.count,
          itemStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 1, 0, [
              { offset: 0, color: '#4e54c8' },
              { offset: 1, color: '#8f94fb' }
            ])
          }
        })),
        label: {
          show: true,
          position: 'right'
        }
      }
    ]
  })
}

// 获取严重程度标签样式
const getSeverityType = (severity) => {
  const map = {
    critical: 'danger',
    high: 'warning',
    medium: 'primary',
    low: 'success',
    info: 'info'
  }
  return map[severity] || 'info'
}

// 获取严重程度标签文本
const getSeverityLabel = (severity) => {
  return severity
}

// 获取状态标签样式
const getStatusType = (status) => {
  const map = {
    new: 'danger',
    verified: 'warning',
    in_progress: 'primary',
    fixed: 'success',
    closed: 'info',
    false_positive: 'info'
  }
  return map[status] || 'info'
}

// 获取状态标签文本
const getStatusLabel = (status) => {
  return status
}

// 获取活动类型对应的图标
const getActivityIcon = (type) => {
  const iconMap = {
    'vulnerability': 'WarningFilled',
    'asset': 'Monitor',
    'user': 'UserFilled',
    'setting': 'Setting',
    'scan': 'DataAnalysis',
    'knowledge': 'Reading',
    'vulndb': 'MoreFilled',
    'default': 'InfoFilled'
  }
  return iconMap[type] || iconMap.default
}

// 获取活动类型对应的标签文本
const getActivityTypeLabel = (type) => {
  const labelMap = {
    'vulnerability': '漏洞活动',
    'asset': '资产变更',
    'user': '用户管理',
    'setting': '系统设置',
    'scan': '扫描任务',
    'knowledge': '知识库管理',
    'vulndb': '漏洞库管理',
    'default': '系统活动'
  }
  return labelMap[type] || labelMap.default
}

// 获取活动类型
const getActivityType = (type) => {
  const typeMap = {
    'vulnerability': 'danger',
    'asset': 'success',
    'user': 'primary',
    'setting': 'info',
    'scan': 'warning',
    'knowledge': 'success',
    'vulndb': 'danger',
    'default': 'info'
  }
  return typeMap[type] || typeMap.default
}

// 格式化日期
const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleString()
}

// 查看详情
const handleViewDetail = (id) => {
  router.push(`/vulnerabilities/${id}`)
}

// 查看更多漏洞
const handleMoreVulns = () => {
  router.push('/vulnerabilities')
}

// 初始化图表和数据
const initDashboard = () => {
  fetchStats()
  fetchTrendsData()
  fetchSeverityDistribution()
  fetchAssetVulnDistribution()
  fetchPriorityVulns()
  fetchRecentActivities()
}

// 窗口大小变化时重新调整图表大小
const handleResize = () => {
  if (trendChart) trendChart.resize()
  if (severityChart) severityChart.resize()
  if (assetVulnChart) assetVulnChart.resize()
}

onMounted(() => {
  initDashboard()
  window.addEventListener('resize', handleResize)
  
  // 启动自动刷新
  if (autoRefresh.value) {
    startAutoRefresh()
  }
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', handleResize)
  if (trendChart) trendChart.dispose()
  if (severityChart) severityChart.dispose()
  if (assetVulnChart) assetVulnChart.dispose()
  
  // 清除自动刷新定时器
  stopAutoRefresh()
})
</script>

<style scoped>
.dashboard-container {
  padding: 20px;
  animation: fadeIn 0.6s ease-out;
}

.mt-4 {
  margin-top: 20px;
}

.stat-card {
  min-height: 120px;
  padding: 10px;
  border-radius: 8px;
  overflow: hidden;
  position: relative;
  transition: all 0.3s ease;
  margin-bottom: 20px;
}

.stat-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 4px;
  background: linear-gradient(90deg, #4e54c8, #8f94fb);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
}

.stat-card:hover::before {
  opacity: 1;
}

.stat-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.stat-value {
  font-size: 36px;
  font-weight: bold;
  color: #333;
  background: linear-gradient(90deg, #333, #666);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.stat-icon {
  font-size: 48px;
  color: #4e54c8;
  opacity: 0.8;
}

.stat-icon.critical {
  color: #F56C6C;
}

.stat-icon.success {
  color: #67C23A;
}

.stat-title {
  font-size: 16px;
  color: #606266;
  font-weight: 500;
  margin-bottom: 8px;
}

.stat-footer {
  display: flex;
  align-items: center;
  margin-top: 8px;
}

.stat-period {
  margin-left: 5px;
  font-size: 12px;
  color: #909399;
}

.chart-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 16px;
  font-weight: 500;
}

.card-header-actions {
  display: flex;
  align-items: center;
}

.chart-container {
  width: 100%;
  height: 100%;
}

.activities-container {
  max-height: 400px;
  overflow-y: auto;
}

.activity-title {
  font-size: 14px;
  font-weight: bold;
  margin: 0 0 8px 0;
  color: #606266;
}

.activity-content {
  margin-bottom: 8px;
}

.activity-content strong {
  font-weight: bold;
  color: #409EFF;
}

.activity-content[data-type="knowledge"] strong {
  color: #67C23A;
}

.activity-content[data-type="vulndb"] strong {
  color: #F56C6C;
}

.activity-user {
  display: flex;
  align-items: center;
  font-size: 12px;
  color: #909399;
}

.activity-user .el-icon {
  margin-right: 4px;
  font-size: 14px;
}

.empty-data {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 20px 0;
}

:deep(.el-table) {
  border-radius: 8px;
  overflow: hidden;
}

:deep(.el-table__row) {
  transition: all 0.3s;
}

:deep(.el-table__row:hover) {
  background-color: #f5f7fa;
  transform: translateY(-2px);
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);
}

:deep(.el-card) {
  border-radius: 8px;
  border: none;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
  transition: all 0.3s ease;
}

:deep(.el-card:hover) {
  box-shadow: 0 6px 18px rgba(0, 0, 0, 0.08);
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@media (max-width: 768px) {
  .dashboard-container {
    padding: 15px;
  }
  
  .stat-value {
    font-size: 30px;
  }
  
  .stat-icon {
    font-size: 36px;
  }
}

.activities-content {
  margin-left: 5px;
}

.refresh-card {
  margin-top: 20px;
}

.refresh-controls {
  display: flex;
  align-items: center;
  gap: 15px;
}

.refresh-interval {
  display: flex;
  align-items: center;
  gap: 5px;
}

.action-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s;
}

.action-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.action-btn .el-icon {
  margin-right: 4px;
}

.operation-buttons {
  display: flex;
  justify-content: center;
}

.no-assessment {
  padding: 40px 0;
  text-align: center;
}
</style> 