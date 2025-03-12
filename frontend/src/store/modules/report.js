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
      console.log('获取报告列表:', {
        baseURL: axios.defaults.baseURL,
        endpoint: '/api/reports',
        params: { page, limit, type, format, search }
      })
      
      const response = await axios.get('/api/reports', {
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
      console.log('获取报告详情:', {
        baseURL: axios.defaults.baseURL,
        endpoint: `/api/reports/${id}`
      })
      
      const response = await axios.get(`/api/reports/${id}`)
      
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
      // 获取令牌
      const token = localStorage.getItem('token')
      console.log('创建报告前令牌检查:', {
        tokenFromLocalStorage: !!token
      })
      
      if (!token) {
        throw new Error('未找到授权令牌，请重新登录')
      }
      
      console.log('创建报告:', {
        baseURL: axios.defaults.baseURL,
        endpoint: '/api/reports',
        data: reportData,
        hasToken: !!token
      })
      
      // 准备请求数据
      const requestData = {
        name: reportData.name,
        type: reportData.type,
        format: reportData.format,
        description: reportData.description || '',
        start_date: reportData.start_date,
        end_date: reportData.end_date,
        severities: reportData.severities,
        statuses: reportData.statuses || ['open', 'in_progress', 'resolved'] // 默认包含所有状态
      }
      
      // 发送请求，显式添加令牌到请求头
      const response = await axios.post('/api/reports', requestData, {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      })
      
      // 处理响应
      let newReport = null
      if (response.data.data) {
        newReport = response.data.data
      } else if (response.data) {
        newReport = response.data
      }
      
      // 添加到状态中
      if (newReport) {
        commit('ADD_REPORT', newReport)
        
        // 添加成功通知
        dispatch('addNotification', {
          type: 'success',
          message: '报告创建请求已提交，正在生成报告',
          title: '请求成功'
        }, { root: true })
      }
      
      commit('SET_GENERATING', false)
      
      // 启动轮询检查报告状态
      dispatch('pollReportStatus', newReport.id)
      
      return newReport
    } catch (error) {
      console.error('创建报告详细错误:', error)
      
      // 获取更详细的错误信息
      let errorMessage = '创建报告失败'
      
      if (error.response) {
        // 服务器返回了错误状态码
        if (error.response.status === 401) {
          errorMessage = '身份验证失败，请重新登录'
        } else {
          errorMessage = error.response.data?.error || '服务器返回错误: ' + error.response.status
        }
      } else if (error.request) {
        // 请求已发送但没有收到响应
        errorMessage = '服务器无响应，请检查网络连接'
      } else {
        // 请求设置时出现问题
        errorMessage = error.message || '请求设置错误'
      }
      
      commit('SET_ERROR', errorMessage)
      commit('SET_GENERATING', false)
      
      // 添加错误通知
      dispatch('addNotification', {
        type: 'error',
        message: errorMessage,
        title: '创建失败'
      }, { root: true })
      
      throw error
    }
  },
  
  // 轮询检查报告状态
  async pollReportStatus({ dispatch }, reportId) {
    // 最多轮询10次，每3秒一次
    let attempts = 0;
    const maxAttempts = 10;
    const pollInterval = 3000; // 3秒
    
    const poll = async () => {
      if (attempts >= maxAttempts) {
        return;
      }
      
      attempts++;
      
      try {
        const report = await dispatch('fetchReportById', reportId);
        
        // 如果报告已完成或失败，停止轮询
        if (report.status === 'completed' || report.status === 'failed') {
          // 更新UI
          dispatch('fetchReports');
          
          // 根据状态发送通知
          if (report.status === 'completed') {
            dispatch('addNotification', {
              type: 'success',
              message: `报告 "${report.name}" 已成功生成`,
              title: '报告已完成'
            }, { root: true });
          } else {
            dispatch('addNotification', {
              type: 'error',
              message: `报告 "${report.name}" 生成失败`,
              title: '报告生成失败'
            }, { root: true });
          }
          
          return;
        }
        
        // 继续轮询
        setTimeout(poll, pollInterval);
      } catch (error) {
        console.error('轮询报告状态失败:', error);
        setTimeout(poll, pollInterval);
      }
    };
    
    // 开始轮询
    setTimeout(poll, pollInterval);
  },
  
  // 删除报告
  async deleteReport({ commit, dispatch }, id) {
    commit('SET_LOADING', true)
    commit('SET_ERROR', null)
    
    try {
      console.log('删除报告:', {
        baseURL: axios.defaults.baseURL,
        endpoint: `/api/reports/${id}`
      })
      
      await axios.delete(`/api/reports/${id}`)
      
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
      console.log('生成摘要报告:', {
        baseURL: axios.defaults.baseURL,
        endpoint: '/api/reports/summary',
        params: { start_date: startDate, end_date: endDate }
      })
      
      const response = await axios.get('/api/reports/summary', {
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
      console.log('生成详细报告:', {
        baseURL: axios.defaults.baseURL,
        endpoint: '/api/reports/detailed',
        params: { start_date: startDate, end_date: endDate, severity, status }
      })
      
      const response = await axios.get('/api/reports/detailed', {
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