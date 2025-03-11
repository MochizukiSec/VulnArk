// asset.js - 资产管理模块
import { ElMessage } from 'element-plus'

export default {
  namespaced: true,
  
  state: {
    assets: [],
    asset: null,
    loading: false,
    error: null,
    pagination: {
      total: 0,
      page: 1,
      perPage: 10
    }
  },
  
  getters: {
    assets: state => state.assets,
    asset: state => state.asset,
    isLoading: state => state.loading,
    error: state => state.error,
    pagination: state => state.pagination
  },
  
  mutations: {
    SET_ASSETS(state, assets) {
      state.assets = assets
    },
    
    SET_ASSET(state, asset) {
      state.asset = asset
    },
    
    ADD_ASSET(state, asset) {
      state.assets.unshift(asset)
    },
    
    UPDATE_ASSET(state, updatedAsset) {
      const index = state.assets.findIndex(a => a.id === updatedAsset.id)
      if (index !== -1) {
        state.assets.splice(index, 1, updatedAsset)
      }
      if (state.asset && state.asset.id === updatedAsset.id) {
        state.asset = updatedAsset
      }
    },
    
    REMOVE_ASSET(state, id) {
      state.assets = state.assets.filter(a => a.id !== id)
      if (state.asset && state.asset.id === id) {
        state.asset = null
      }
    },
    
    SET_LOADING(state, loading) {
      state.loading = loading
    },
    
    SET_ERROR(state, error) {
      state.error = error
    },
    
    SET_PAGINATION(state, pagination) {
      state.pagination = { ...state.pagination, ...pagination }
    }
  },
  
  actions: {
    // 获取资产列表
    async fetchAssets({ commit }, { page = 1, perPage = 10, search = '', filters = {} } = {}) {
      commit('SET_LOADING', true)
      
      try {
        // 模拟API请求
        console.log('Fetching assets with params:', { page, perPage, search, filters })
        
        // 使用模拟数据
        setTimeout(() => {
          const assets = []
          const total = 0
          
          commit('SET_ASSETS', assets)
          commit('SET_PAGINATION', { total, page, perPage })
          commit('SET_LOADING', false)
        }, 500)
        
        // 实际项目中的API调用
        // const response = await fetch(`/api/assets?page=${page}&perPage=${perPage}&search=${search}&...`)
        // const data = await response.json()
        // commit('SET_ASSETS', data.items)
        // commit('SET_PAGINATION', { 
        //   total: data.total,
        //   page: data.page,
        //   perPage: data.perPage
        // })
      } catch (error) {
        console.error('获取资产列表失败:', error)
        commit('SET_ERROR', '获取资产列表失败')
        ElMessage.error('获取资产列表失败，请稍后重试')
      } finally {
        commit('SET_LOADING', false)
      }
    },
    
    // 获取单个资产详情
    async fetchAssetById({ commit }, id) {
      commit('SET_LOADING', true)
      
      try {
        // 模拟API请求
        console.log('Fetching asset with id:', id)
        
        // 使用模拟数据
        setTimeout(() => {
          const asset = {
            id,
            name: `资产 ${id}`,
            ip: '192.168.1.1',
            type: 'server',
            os: 'Linux',
            owner: '系统管理员',
            status: 'active',
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString()
          }
          
          commit('SET_ASSET', asset)
          commit('SET_LOADING', false)
        }, 500)
        
        // 实际项目中的API调用
        // const response = await fetch(`/api/assets/${id}`)
        // const data = await response.json()
        // commit('SET_ASSET', data)
      } catch (error) {
        console.error('获取资产详情失败:', error)
        commit('SET_ERROR', '获取资产详情失败')
        ElMessage.error('获取资产详情失败，请稍后重试')
      } finally {
        commit('SET_LOADING', false)
      }
    },
    
    // 创建资产
    async createAsset({ commit }, assetData) {
      commit('SET_LOADING', true)
      
      try {
        // 模拟API请求
        console.log('Creating asset with data:', assetData)
        
        // 使用模拟数据
        setTimeout(() => {
          const newAsset = {
            id: Date.now().toString(),
            ...assetData,
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString()
          }
          
          commit('ADD_ASSET', newAsset)
          commit('SET_LOADING', false)
          ElMessage.success('资产创建成功')
          return newAsset
        }, 500)
        
        // 实际项目中的API调用
        // const response = await fetch('/api/assets', {
        //   method: 'POST',
        //   headers: {
        //     'Content-Type': 'application/json'
        //   },
        //   body: JSON.stringify(assetData)
        // })
        // const data = await response.json()
        // commit('ADD_ASSET', data)
        // ElMessage.success('资产创建成功')
        // return data
      } catch (error) {
        console.error('创建资产失败:', error)
        commit('SET_ERROR', '创建资产失败')
        ElMessage.error('创建资产失败，请稍后重试')
      } finally {
        commit('SET_LOADING', false)
      }
    },
    
    // 更新资产
    async updateAsset({ commit }, { id, data }) {
      commit('SET_LOADING', true)
      
      try {
        // 模拟API请求
        console.log('Updating asset:', { id, data })
        
        // 使用模拟数据
        setTimeout(() => {
          const updatedAsset = {
            id,
            ...data,
            updatedAt: new Date().toISOString()
          }
          
          commit('UPDATE_ASSET', updatedAsset)
          commit('SET_LOADING', false)
          ElMessage.success('资产更新成功')
          return updatedAsset
        }, 500)
        
        // 实际项目中的API调用
        // const response = await fetch(`/api/assets/${id}`, {
        //   method: 'PUT',
        //   headers: {
        //     'Content-Type': 'application/json'
        //   },
        //   body: JSON.stringify(data)
        // })
        // const updatedAsset = await response.json()
        // commit('UPDATE_ASSET', updatedAsset)
        // ElMessage.success('资产更新成功')
        // return updatedAsset
      } catch (error) {
        console.error('更新资产失败:', error)
        commit('SET_ERROR', '更新资产失败')
        ElMessage.error('更新资产失败，请稍后重试')
      } finally {
        commit('SET_LOADING', false)
      }
    },
    
    // 删除资产
    async deleteAsset({ commit }, id) {
      commit('SET_LOADING', true)
      
      try {
        // 模拟API请求
        console.log('Deleting asset with id:', id)
        
        // 使用模拟数据
        setTimeout(() => {
          commit('REMOVE_ASSET', id)
          commit('SET_LOADING', false)
          ElMessage.success('资产删除成功')
        }, 500)
        
        // 实际项目中的API调用
        // await fetch(`/api/assets/${id}`, {
        //   method: 'DELETE'
        // })
        // commit('REMOVE_ASSET', id)
        // ElMessage.success('资产删除成功')
      } catch (error) {
        console.error('删除资产失败:', error)
        commit('SET_ERROR', '删除资产失败')
        ElMessage.error('删除资产失败，请稍后重试')
      } finally {
        commit('SET_LOADING', false)
      }
    }
  }
} 