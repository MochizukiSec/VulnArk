import axios from 'axios'
import { ElMessage } from 'element-plus'

const state = {
  vulnerabilities: [],
  currentVulnerability: null,
  loading: false,
  error: null,
  searchParams: {
    cveId: '',
    severity: '',
    vendor: '',
    product: '',
    yearFrom: '',
    yearTo: '',
    searchTerm: '',
    sortBy: 'published',
    sortOrder: 'desc',
    page: 1,
    perPage: 20
  },
  pagination: {
    total: 0,
    pages: 0
  }
}

const getters = {
  dbVulnerabilities: state => state.vulnerabilities,
  currentDbVulnerability: state => state.currentVulnerability,
  isLoading: state => state.loading,
  error: state => state.error,
  searchParams: state => state.searchParams,
  pagination: state => state.pagination
}

const mutations = {
  SET_DB_VULNERABILITIES(state, vulnerabilities) {
    state.vulnerabilities = vulnerabilities
  },
  
  SET_CURRENT_DB_VULNERABILITY(state, vulnerability) {
    state.currentVulnerability = vulnerability
  },
  
  CLEAR_CURRENT_DB_VULNERABILITY(state) {
    state.currentVulnerability = null
  },
  
  SET_LOADING(state, loading) {
    state.loading = loading
  },
  
  SET_ERROR(state, error) {
    state.error = error
  },
  
  SET_SEARCH_PARAMS(state, params) {
    state.searchParams = { ...state.searchParams, ...params }
  },
  
  SET_PAGINATION(state, pagination) {
    state.pagination = pagination
  },
  
  RESET_SEARCH_PARAMS(state) {
    state.searchParams = {
      cveId: '',
      severity: '',
      vendor: '',
      product: '',
      yearFrom: '',
      yearTo: '',
      searchTerm: '',
      sortBy: 'published',
      sortOrder: 'desc',
      page: 1,
      perPage: 20
    }
  }
}

const actions = {
  // 获取漏洞库数据
  async fetchVulnerabilities({ commit, state }, payload = {}) {
    try {
      commit('SET_LOADING', true)
      commit('SET_ERROR', null)
      
      const params = { ...state.searchParams, ...payload }
      
      console.log('获取漏洞库列表:', {
        baseURL: axios.defaults.baseURL,
        endpoint: '/api/vulndatabase',
        params
      })
      
      const response = await axios.get('/api/vulndatabase', { params })
      
      // 更新状态
      commit('SET_DB_VULNERABILITIES', response.data.vulnerabilities)
      commit('SET_PAGINATION', {
        total: response.data.total,
        pages: Math.ceil(response.data.total / state.searchParams.perPage)
      })
      
      return response.data.vulnerabilities
    } catch (error) {
      commit('SET_ERROR', error.response?.data?.error || '获取漏洞库数据失败')
      throw error
    } finally {
      commit('SET_LOADING', false)
    }
  },
  
  // 获取单个漏洞详情
  async fetchVulnerabilityByCveId({ commit }, cveId) {
    try {
      commit('SET_LOADING', true)
      commit('SET_ERROR', null)
      
      // 调用API获取漏洞详情
      const response = await axios.get(`/api/vulndatabase/${cveId}`)
      
      // 更新状态
      commit('SET_CURRENT_DB_VULNERABILITY', response.data)
      
      return response.data
    } catch (error) {
      commit('SET_ERROR', error.response?.data?.error || '获取漏洞详情失败')
      throw error
    } finally {
      commit('SET_LOADING', false)
    }
  },
  
  // 更新搜索参数
  updateSearchParams({ commit, dispatch }, params) {
    // 更新参数时重置页码为1
    if (Object.keys(params).some(key => key !== 'page')) {
      params.page = 1
    }
    
    commit('SET_SEARCH_PARAMS', params)
    return dispatch('fetchVulnerabilities')
  },
  
  // 重置搜索参数
  resetSearchParams({ commit, dispatch }) {
    commit('RESET_SEARCH_PARAMS')
    return dispatch('fetchVulnerabilities')
  },
  
  // 清除当前漏洞
  clearCurrentVulnerability({ commit }) {
    commit('CLEAR_CURRENT_DB_VULNERABILITY')
  },
  
  // 创建新漏洞
  async createVulnerability({ commit }, vulnerabilityData) {
    try {
      commit('SET_LOADING', true)
      commit('SET_ERROR', null)
      
      // 调用API创建漏洞
      const response = await axios.post('/api/vulndatabase', vulnerabilityData)
      
      // 更新状态
      ElMessage.success('漏洞创建成功')
      
      return response.data
    } catch (error) {
      commit('SET_ERROR', error.response?.data?.error || '创建漏洞失败')
      ElMessage.error(error.response?.data?.error || '创建漏洞失败')
      throw error
    } finally {
      commit('SET_LOADING', false)
    }
  },
  
  // 更新漏洞
  async updateVulnerability({ commit }, { cveId, vulnerabilityData }) {
    try {
      commit('SET_LOADING', true)
      commit('SET_ERROR', null)
      
      // 调用API更新漏洞
      const response = await axios.put(`/api/vulndatabase/${cveId}`, vulnerabilityData)
      
      // 更新状态
      if (response.data) {
        commit('SET_CURRENT_DB_VULNERABILITY', response.data)
      }
      
      ElMessage.success('漏洞更新成功')
      
      return response.data
    } catch (error) {
      commit('SET_ERROR', error.response?.data?.error || '更新漏洞失败')
      ElMessage.error(error.response?.data?.error || '更新漏洞失败')
      throw error
    } finally {
      commit('SET_LOADING', false)
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