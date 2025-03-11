import axios from 'axios'

const state = {
  reports: [],
  loading: false,
  error: null,
  currentReport: null,
  generating: false,
  summaryReport: null,
  detailedReport: null,
  pagination: {
    total: 0,
    pages: 0,
    page: 1,
    perPage: 10
  }
}

const getters = {
  allReports: state => state.reports,
  isLoading: state => state.loading,
  isGenerating: state => state.generating,
  error: state => state.error,
  currentReport: state => state.currentReport,
  summaryReport: state => state.summaryReport,
  detailedReport: state => state.detailedReport,
  pagination: state => state.pagination
}

const mutations = {
  SET_REPORTS(state, reports) {
    state.reports = reports
  },
  
  SET_LOADING(state, loading) {
    state.loading = loading
  },
  
  SET_GENERATING(state, generating) {
    state.generating = generating
  },
  
  SET_ERROR(state, error) {
    state.error = error
  },
  
  SET_PAGINATION(state, pagination) {
    state.pagination = pagination
  },
  
  SET_CURRENT_REPORT(state, report) {
    state.currentReport = report
  },
  
  SET_SUMMARY_REPORT(state, report) {
    state.summaryReport = report
  },
  
  SET_DETAILED_REPORT(state, report) {
    state.detailedReport = report
  },
  
  ADD_REPORT(state, report) {
    state.reports.unshift(report)
  },
  
  REMOVE_REPORT(state, reportId) {
    state.reports = state.reports.filter(r => r.id !== reportId)
  }
}

const actions = {
  // 获取所有报告
  async fetchReports({ commit }, { page = 1, limit = 10, type = '', format = '', search = '' } = {}) {
    commit('SET_LOADING', true)
    commit('SET_ERROR', null)
    
    try {
      const response = await axios.get('/reports', {
        params: {
          page,
          limit,
          type,
          format,
          search
        }
      })
      
      // 处理响应数据
      let reports = []
      let totalItems = 0
      let totalPages = 0
      
      // 处理不同的响应格式
      if (response.data.reports) {
        reports = response.data.reports
        totalItems = response.data.total || 0
        totalPages = response.data.pages || 0
      } else if (response.data.data) {
        reports = response.data.data
        totalItems = response.data.meta?.total || 0
        totalPages = response.data.meta?.pages || 0
      } else {
        reports = Array.isArray(response.data) ? response.data : []
      }
      
      const pagination = {
        total: totalItems,
        pages: totalPages,
        page: page,
        perPage: limit
      }
      
      commit('SET_REPORTS', reports)
      commit('SET_PAGINATION', pagination)
      commit('SET_LOADING', false)
      
      return response.data
    } catch (error) {
      const message = error.response?.data?.error || '获取报告列表失败'
      commit('SET_ERROR', message)
      commit('SET_LOADING', false)
      throw error
    }
  },
  
  // 获取单个报告详情
  async fetchReportById({ commit }, id) {
    commit('SET_LOADING', true)
    commit('SET_ERROR', null)
    
    try {
      const response = await axios.get(`/reports/${id}`)
      
      let report = null
      if (response.data.report) {
        report = response.data.report
      } else {
        report = response.data
      }
      
      commit('SET_CURRENT_REPORT', report)
      commit('SET_LOADING', false)
      
      return report
    } catch (error) {
      const message = error.response?.data?.error || '获取报告详情失败'
      commit('SET_ERROR', message)
      commit('SET_LOADING', false)
      throw error
    }
  },
  
  // 创建新报告
  async createReport({ commit, dispatch }, reportData) {
    commit('SET_GENERATING', true)
    commit('SET_ERROR', null)
    
    try {
      const response = await axios.post('/reports', reportData)
      
      let newReport = null
      if (response.data.report) {
        newReport = response.data.report
      } else {
        newReport = response.data
      }
      
      commit('ADD_REPORT', newReport)
      commit('SET_GENERATING', false)
      
      // 添加成功通知
      dispatch('addNotification', {
        type: 'success',
        message: '报告创建请求已提交，正在生成报告',
        title: '请求成功'
      }, { root: true })
      
      return newReport
    } catch (error) {
      const message = error.response?.data?.error || '创建报告失败'
      commit('SET_ERROR', message)
      commit('SET_GENERATING', false)
      
      // 添加错误通知
      dispatch('addNotification', {
        type: 'error',
        message: message,
        title: '创建失败'
      }, { root: true })
      
      throw error
    }
  },
  
  // 删除报告
  async deleteReport({ commit, dispatch }, id) {
    commit('SET_LOADING', true)
    commit('SET_ERROR', null)
    
    try {
      await axios.delete(`/reports/${id}`)
      
      commit('REMOVE_REPORT', id)
      commit('SET_LOADING', false)
      
      // 添加成功通知
      dispatch('addNotification', {
        type: 'success',
        message: '报告已成功删除',
        title: '删除成功'
      }, { root: true })
      
      return { success: true }
    } catch (error) {
      const message = error.response?.data?.error || '删除报告失败'
      commit('SET_ERROR', message)
      commit('SET_LOADING', false)
      
      // 添加错误通知
      dispatch('addNotification', {
        type: 'error',
        message: message,
        title: '删除失败'
      }, { root: true })
      
      throw error
    }
  },
  
  // 生成摘要报告
  async generateSummaryReport({ commit, dispatch }, { startDate, endDate }) {
    commit('SET_GENERATING', true)
    commit('SET_ERROR', null)
    
    try {
      const response = await axios.get('/reports/summary', {
        params: {
          start_date: startDate,
          end_date: endDate
        }
      })
      
      let report = null
      if (response.data.report) {
        report = response.data.report
      } else {
        report = response.data
      }
      
      commit('SET_SUMMARY_REPORT', report)
      commit('SET_GENERATING', false)
      
      return report
    } catch (error) {
      const message = error.response?.data?.error || '生成摘要报告失败'
      commit('SET_ERROR', message)
      commit('SET_GENERATING', false)
      
      // 添加错误通知
      dispatch('addNotification', {
        type: 'error',
        message: message,
        title: '生成失败'
      }, { root: true })
      
      throw error
    }
  },
  
  // 生成详细报告
  async generateDetailedReport({ commit, dispatch }, { startDate, endDate, severity, status }) {
    commit('SET_GENERATING', true)
    commit('SET_ERROR', null)
    
    try {
      const response = await axios.get('/reports/detailed', {
        params: {
          start_date: startDate,
          end_date: endDate,
          severity,
          status
        }
      })
      
      let report = null
      if (response.data.report) {
        report = response.data.report
      } else {
        report = response.data
      }
      
      commit('SET_DETAILED_REPORT', report)
      commit('SET_GENERATING', false)
      
      return report
    } catch (error) {
      const message = error.response?.data?.error || '生成详细报告失败'
      commit('SET_ERROR', message)
      commit('SET_GENERATING', false)
      
      // 添加错误通知
      dispatch('addNotification', {
        type: 'error',
        message: message,
        title: '生成失败'
      }, { root: true })
      
      throw error
    }
  }
}

export default {
  namespaced: true,
  state,
  getters,
  mutations,
  actions
} 