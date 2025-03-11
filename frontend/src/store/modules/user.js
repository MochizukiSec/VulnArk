import axios from 'axios'

const state = {
  users: [],
  loading: false,
  error: null,
  pagination: {
    total: 0,
    pages: 0,
    page: 1,
    perPage: 10
  }
}

const getters = {
  allUsers: state => state.users,
  isLoading: state => state.loading,
  error: state => state.error,
  pagination: state => state.pagination
}

const mutations = {
  SET_USERS(state, users) {
    state.users = users
  },
  
  SET_LOADING(state, loading) {
    state.loading = loading
  },
  
  SET_ERROR(state, error) {
    state.error = error
  },
  
  SET_PAGINATION(state, pagination) {
    state.pagination = pagination
  },
  
  ADD_USER(state, user) {
    state.users.push(user)
  },
  
  UPDATE_USER(state, updatedUser) {
    const index = state.users.findIndex(u => u.id === updatedUser.id)
    if (index !== -1) {
      state.users[index] = updatedUser
    }
  },
  
  REMOVE_USER(state, userId) {
    state.users = state.users.filter(u => u.id !== userId)
  }
}

const actions = {
  // 获取所有用户
  async fetchUsers({ commit }, { page = 1, limit = 10, search = '', role = '', status = '' } = {}) {
    commit('SET_LOADING', true)
    commit('SET_ERROR', null)
    
    try {
      const response = await axios.get('/users', {
        params: {
          page,
          perPage: limit,
          search,
          role,
          status
        }
      })
      
      // 处理响应数据，适应不同的后端返回结构
      let users = []
      let totalItems = 0
      let totalPages = 0
      
      // 假设后端可能返回两种格式之一：
      // 1. { users: [...], total: X, pages: Y }
      // 2. { data: [...], meta: { total: X, pages: Y } }
      
      if (response.data.users) {
        // 格式1
        users = response.data.users
        totalItems = response.data.total || 0
        totalPages = response.data.pages || 0
      } else if (response.data.data) {
        // 格式2
        users = response.data.data
        totalItems = response.data.meta?.total || 0
        totalPages = response.data.meta?.pages || 0
      } else {
        // 默认处理，假设直接返回数组
        users = Array.isArray(response.data) ? response.data : []
        totalItems = users.length
      }
      
      const pagination = {
        total: totalItems,
        pages: totalPages,
        page: page,
        perPage: limit
      }
      
      commit('SET_USERS', users)
      commit('SET_PAGINATION', pagination)
      commit('SET_LOADING', false)
      
      return response.data
    } catch (error) {
      const message = error.response?.data?.error || '获取用户列表失败'
      commit('SET_ERROR', message)
      commit('SET_LOADING', false)
      throw error
    }
  },
  
  // 创建新用户
  async createUser({ commit, dispatch }, userData) {
    try {
      commit('SET_LOADING', true)
      
      // 确保提交的数据格式与后端要求一致
      const requestData = {
        username: userData.username,
        email: userData.email,
        password: userData.password,
        first_name: userData.first_name,
        last_name: userData.last_name || ' ', // 确保提供值，因为这是必填字段
        department: userData.department || '',
        role: userData.role || 'user',
        status: userData.status || 'active'
      }
      
      const response = await axios.post('/users', requestData)
      
      // 适应不同的响应格式
      let user = null
      if (response.data.user) {
        user = response.data.user
      } else if (response.data.data) {
        user = response.data.data
      } else {
        // 假设整个响应就是用户对象
        user = response.data
      }
      
      if (user) {
        commit('ADD_USER', user)
      }
      
      commit('SET_LOADING', false)
      commit('SET_ERROR', null)
      
      // 添加成功通知
      dispatch('addNotification', {
        type: 'success',
        message: '用户已成功创建',
        title: '创建成功'
      }, { root: true })
      
      return response.data
    } catch (error) {
      const message = error.response?.data?.error || '创建用户失败'
      commit('SET_ERROR', message)
      commit('SET_LOADING', false)
      
      // 添加错误通知
      dispatch('addNotification', {
        type: 'error',
        message: message,
        title: '创建失败'
      }, { root: true })
      
      throw error
    }
  },
  
  // 更新用户
  async updateUser({ commit, dispatch }, { id, data }) {
    try {
      commit('SET_LOADING', true)
      
      // 确保提交的数据格式与后端要求一致
      const requestData = {}
      
      // 只包含有值的字段
      if (data.username) requestData.username = data.username
      if (data.email) requestData.email = data.email
      if (data.first_name) requestData.first_name = data.first_name
      if (data.last_name) requestData.last_name = data.last_name
      if (data.department !== undefined) requestData.department = data.department
      if (data.role) requestData.role = data.role
      if (data.status) requestData.status = data.status
      if (data.profile_picture) requestData.profile_picture = data.profile_picture
      if (data.new_password) requestData.new_password = data.new_password
      
      const response = await axios.put(`/users/${id}`, requestData)
      
      // 适应不同的响应格式
      let updatedUser = null
      if (response.data.user) {
        updatedUser = { ...response.data.user }
      } else if (response.data.data) {
        updatedUser = { ...response.data.data }
      } else {
        // 假设整个响应就是用户对象或更新成功的消息
        updatedUser = { id, ...data }
      }
      
      // 如果更新当前用户，同时更新auth模块中的用户信息
      const currentUser = JSON.parse(localStorage.getItem('user'))
      if (currentUser && currentUser.id === id) {
        dispatch('auth/updateUserInfo', data, { root: true })
      }
      
      commit('UPDATE_USER', updatedUser)
      commit('SET_LOADING', false)
      commit('SET_ERROR', null)
      
      // 添加成功通知
      dispatch('addNotification', {
        type: 'success',
        message: '用户信息已成功更新',
        title: '更新成功'
      }, { root: true })
      
      return response.data
    } catch (error) {
      const message = error.response?.data?.error || '更新用户失败'
      commit('SET_ERROR', message)
      commit('SET_LOADING', false)
      
      // 添加错误通知
      dispatch('addNotification', {
        type: 'error',
        message: message,
        title: '更新失败'
      }, { root: true })
      
      throw error
    }
  },
  
  // 删除用户
  async deleteUser({ commit, dispatch }, id) {
    try {
      commit('SET_LOADING', true)
      
      await axios.delete(`/users/${id}`)
      
      commit('REMOVE_USER', id)
      commit('SET_LOADING', false)
      commit('SET_ERROR', null)
      
      // 添加成功通知
      dispatch('addNotification', {
        type: 'success',
        message: '用户已成功删除',
        title: '删除成功'
      }, { root: true })
      
      return { success: true }
    } catch (error) {
      const message = error.response?.data?.error || '删除用户失败'
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
  
  // 更新当前用户个人资料
  async updateProfile({ commit, dispatch }, profileData) {
    try {
      commit('SET_LOADING', true)
      
      const response = await axios.put('/users/me', profileData)
      
      // 更新本地存储的用户信息
      dispatch('auth/updateUserInfo', profileData, { root: true })
      
      commit('SET_LOADING', false)
      commit('SET_ERROR', null)
      
      // 添加成功通知
      dispatch('addNotification', {
        type: 'success',
        message: '个人资料已成功更新',
        title: '更新成功'
      }, { root: true })
      
      return response.data
    } catch (error) {
      const message = error.response?.data?.error || '更新个人资料失败'
      commit('SET_ERROR', message)
      commit('SET_LOADING', false)
      
      // 添加错误通知
      dispatch('addNotification', {
        type: 'error',
        message: message,
        title: '更新失败'
      }, { root: true })
      
      throw error
    }
  },
  
  // 搜索用户
  async searchUsers({ commit, dispatch }, query) {
    try {
      commit('SET_LOADING', true)
      
      // 使用与fetchUsers相同的API，但只提供搜索参数
      return dispatch('fetchUsers', { 
        search: query, 
        page: 1, 
        limit: 10 
      })
    } catch (error) {
      commit('SET_ERROR', error.message || '搜索用户失败')
      return []
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