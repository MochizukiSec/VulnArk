import axios from 'axios'

const state = {
  token: localStorage.getItem('token') || null,
  user: (() => {
    const userStr = localStorage.getItem('user');
    if (userStr && userStr !== 'undefined') {
      try {
        return JSON.parse(userStr);
      } catch (e) {
        console.error('解析用户数据失败:', e);
        return null;
      }
    }
    return null;
  })(),
  loading: false,
  error: null
}

const getters = {
  isAuthenticated: state => !!state.token,
  isAdmin: state => state.user && state.user.role === 'admin',
  currentUser: state => state.user,
  authError: state => state.error,
  authLoading: state => state.loading
}

const mutations = {
  AUTH_REQUEST(state) {
    state.loading = true
    state.error = null
  },
  
  AUTH_SUCCESS(state, { token, user }) {
    state.token = token
    state.user = user
    state.loading = false
    state.error = null
  },
  
  AUTH_ERROR(state, error) {
    state.loading = false
    state.error = error
  },
  
  LOGOUT(state) {
    state.token = null
    state.user = null
  },
  
  UPDATE_USER(state, user) {
    state.user = { ...state.user, ...user }
  }
}

const actions = {
  // 用户登录
  async login({ commit, dispatch }, credentials) {
    try {
      commit('AUTH_REQUEST')
      
      // 在登录时记录详细信息，有助于调试
      console.log('尝试登录:', {
        baseURL: axios.defaults.baseURL,
        fullURL: axios.defaults.baseURL + '/api/auth/login',
        email: credentials.email
      });
      
      const response = await axios.post('/api/auth/login', credentials)
      const { token, user } = response.data
      
      // 保存到本地存储
      localStorage.setItem('token', token)
      localStorage.setItem('user', JSON.stringify(user))
      
      // 更新状态
      commit('AUTH_SUCCESS', { token, user })
      
      // 设置axios默认Authorization头
      axios.defaults.headers.common['Authorization'] = `Bearer ${token}`
      
      // 添加欢迎通知
      dispatch('addNotification', {
        type: 'success',
        message: `欢迎回来，${user.firstName}！`,
        title: '登录成功'
      }, { root: true })
      
      return { success: true }
    } catch (error) {
      const message = error.response?.data?.error || '登录失败，请检查您的凭据'
      commit('AUTH_ERROR', message)
      return { success: false, message }
    }
  },
  
  // 恢复会话
  async restoreSession({ commit, state }) {
    if (!state.token) return { success: false }
    
    try {
      // 设置axios默认Authorization头
      axios.defaults.headers.common['Authorization'] = `Bearer ${state.token}`
      
      // 尝试获取当前用户信息以验证令牌
      // 如果是在首页，使用静默方式处理错误，不触发全局拦截器
      let response;
      
      if (window.location.pathname === '/') {
        try {
          response = await axios.get('/users/me', {
            _silentError: true // 自定义标记，用于标识静默处理错误
          });
        } catch (e) {
          console.log('在首页恢复会话失败，静默处理');
          // 在首页上，我们不希望显示错误或触发重定向
          commit('LOGOUT');
          localStorage.removeItem('token');
          localStorage.removeItem('user');
          delete axios.defaults.headers.common['Authorization'];
          return { success: false };
        }
      } else {
        response = await axios.get('/users/me');
      }
      
      console.log('获取用户信息成功，用户数据:', response.data);
      
      // 更新用户状态
      if (response.data) {
        commit('AUTH_SUCCESS', { 
          token: state.token,
          user: response.data
        });
        
        // 更新本地存储
        localStorage.setItem('user', JSON.stringify(response.data));
      }
      
      return { success: true }
    } catch (error) {
      console.error('恢复会话失败:', error)
      // 如果令牌无效，退出登录
      commit('LOGOUT')
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      delete axios.defaults.headers.common['Authorization']
      return { success: false }
    }
  },
  
  // 用户退出
  logout({ commit }) {
    // 清除本地存储
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    
    // 移除axios默认Authorization头
    delete axios.defaults.headers.common['Authorization']
    
    // 更新状态
    commit('LOGOUT')
  },
  
  // 更新用户信息
  updateUserInfo({ commit }, user) {
    // 更新本地存储
    localStorage.setItem('user', JSON.stringify({ ...state.user, ...user }))
    
    // 更新状态
    commit('UPDATE_USER', user)
  }
}

export default {
  namespaced: true,
  state,
  getters,
  mutations,
  actions
} 