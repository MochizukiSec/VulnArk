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
      
      const loginPath = '/api/auth/login'
      // 在登录时记录详细信息，有助于调试
      console.log('尝试登录:', {
        baseURL: axios.defaults.baseURL,
        email: credentials.email,
        loginPath: loginPath,
        fullURL: `${axios.defaults.baseURL}${loginPath}`
      });
      
      const response = await axios.post(loginPath, credentials)
      const { token, user } = response.data
      
      // 保存到本地存储
      localStorage.setItem('token', token)
      localStorage.setItem('user', JSON.stringify(user))
      
      console.log('登录成功，获取并保存令牌:', {
        tokenReceived: !!token,
        tokenPrefix: token ? token.substring(0, 10) + '...' : '',
        userReceived: !!user,
        localStorageSaved: {
          token: !!localStorage.getItem('token'),
          user: !!localStorage.getItem('user')
        }
      })
      
      // 更新状态
      commit('AUTH_SUCCESS', { token, user })
      
      // 设置axios默认Authorization头
      axios.defaults.headers.common['Authorization'] = `Bearer ${token}`
      
      // 获取完整的用户信息
      try {
        await dispatch('user/fetchCurrentUser', null, { root: true })
      } catch (error) {
        console.warn('无法获取完整的用户信息，将使用登录返回的基本信息')
      }
      
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
      
      // 添加失败通知
      dispatch('addNotification', {
        type: 'error',
        message: message,
        title: '登录失败'
      }, { root: true }).catch(err => {
        console.error('无法添加通知:', err);
      });
      
      return { success: false, message }
    }
  },
  
  // 刷新token状态 - 用于解决授权问题
  refreshTokenState({ commit }) {
    try {
      // 尝试从localStorage重新获取token
      const token = localStorage.getItem('token')
      let user = null
      
      try {
        const userStr = localStorage.getItem('user')
        if (userStr) {
          user = JSON.parse(userStr)
        }
      } catch (e) {
        console.error('解析用户数据失败:', e)
      }
      
      if (token) {
        // 重新设置到store中
        commit('AUTH_SUCCESS', { token, user })
        
        // 重新设置axios默认头
        axios.defaults.headers.common['Authorization'] = `Bearer ${token}`
        
        console.log('令牌状态已刷新:', {
          tokenRefreshed: true,
          userRefreshed: !!user
        })
        
        return true
      }
      
      return false
    } catch (error) {
      console.error('刷新令牌状态失败:', error)
      return false
    }
  },
  
  // 添加明确的authError处理
  authError({ commit, dispatch }, message) {
    commit('AUTH_ERROR', message);
    
    // 添加错误通知
    dispatch('addNotification', {
      type: 'error',
      message: message || '认证错误',
      title: '认证失败'
    }, { root: true }).catch(err => {
      console.error('无法添加通知:', err);
    });
    
    return { success: false, message };
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
          response = await axios.get('/api/users/me', {
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
        response = await axios.get('/api/users/me');
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