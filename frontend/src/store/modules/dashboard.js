import axios from 'axios'

const state = {
  data: {
    totalVulnerabilities: 0,
    vulnerabilitiesBySeverity: {},
    vulnerabilitiesByStatus: {},
    recentVulnerabilities: [],
    topAffectedSystems: [],
    vulnerabilitiesByMonth: [],
    riskScore: 0,
    teamVulnerabilities: [],
    criticalVulnerabilities: [],
    remediationProgress: {
      resolvedCount: 0,
      totalCount: 0,
      progressRate: 0,
      averageDays: 0
    },
    vulnerabilityTrends: {
      newVulnerabilities: [],
      resolvedVulnerabilities: [],
      timeLabels: [],
      netChange: []
    }
  },
  loading: false,
  error: null
}

const getters = {
  dashboardData: state => state.data,
  isLoading: state => state.loading,
  error: state => state.error,
  
  // 获取总漏洞数
  totalVulnerabilities: state => state.data.totalVulnerabilities || 0,
  
  // 获取严重程度统计
  severityCounts: state => state.data.vulnerabilitiesBySeverity || {},
  
  // 获取状态统计
  statusCounts: state => state.data.vulnerabilitiesByStatus || {},
  
  // 获取风险分数
  riskScore: state => Math.round(state.data.riskScore || 0),
  
  // 获取团队漏洞统计
  teamVulnerabilities: state => state.data.teamVulnerabilities || [],
  
  // 获取需优先关注的高危漏洞
  criticalVulnerabilities: state => state.data.criticalVulnerabilities || [],
  
  // 获取修复进度
  remediationProgress: state => state.data.remediationProgress || {
    resolvedCount: 0,
    totalCount: 0,
    progressRate: 0,
    averageDays: 0
  },
  
  // 获取漏洞趋势
  vulnerabilityTrends: state => state.data.vulnerabilityTrends || {
    newVulnerabilities: [],
    resolvedVulnerabilities: [],
    timeLabels: [],
    netChange: []
  },
  
  // 获取修复率百分比
  remediationRate: state => {
    const progress = state.data.remediationProgress || {}
    if (progress.totalCount > 0) {
      return Math.round(progress.progressRate)
    }
    return 0
  }
}

const mutations = {
  SET_DASHBOARD_DATA(state, data) {
    state.data = data
  },
  
  SET_LOADING(state, loading) {
    state.loading = loading
  },
  
  SET_ERROR(state, error) {
    state.error = error
  }
}

const actions = {
  // 获取仪表盘数据
  async fetchDashboardData({ commit }) {
    try {
      commit('SET_LOADING', true)
      
      const response = await axios.get('/dashboard')
      
      // 解析仪表盘数据
      const dashboardData = response.data
      
      commit('SET_DASHBOARD_DATA', dashboardData)
      commit('SET_LOADING', false)
      
      return dashboardData
    } catch (error) {
      const message = error.response?.data?.error || '获取仪表盘数据失败'
      commit('SET_ERROR', message)
      commit('SET_LOADING', false)
      throw error
    }
  },
  
  // 重置仪表盘数据
  resetDashboardData({ commit }) {
    commit('SET_DASHBOARD_DATA', {
      totalVulnerabilities: 0,
      vulnerabilitiesBySeverity: {},
      vulnerabilitiesByStatus: {},
      recentVulnerabilities: [],
      topAffectedSystems: [],
      vulnerabilitiesByMonth: [],
      riskScore: 0,
      teamVulnerabilities: [],
      criticalVulnerabilities: [],
      remediationProgress: {
        resolvedCount: 0,
        totalCount: 0,
        progressRate: 0,
        averageDays: 0
      },
      vulnerabilityTrends: {
        newVulnerabilities: [],
        resolvedVulnerabilities: [],
        timeLabels: [],
        netChange: []
      }
    })
  }
}

export default {
  namespaced: true,
  state,
  getters,
  mutations,
  actions
} 