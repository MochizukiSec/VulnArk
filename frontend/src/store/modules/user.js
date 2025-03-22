import { login, getUserInfo, updateUserInfo } from '@/api/user'
import { getToken, setToken, removeToken } from '@/utils/auth'

// 从本地存储恢复用户信息
const loadUserInfo = () => {
  try {
    const userInfoStr = localStorage.getItem('vuex_user_info')
    return userInfoStr ? JSON.parse(userInfoStr) : null
  } catch (e) {
    console.error('恢复用户信息失败:', e)
    return null
  }
}

// 保存用户信息到本地存储
const saveUserInfo = (userInfo) => {
  try {
    if (userInfo) {
      localStorage.setItem('vuex_user_info', JSON.stringify(userInfo))
    } else {
      localStorage.removeItem('vuex_user_info')
    }
  } catch (e) {
    console.error('保存用户信息失败:', e)
  }
}

const state = {
  token: getToken(),
  userInfo: loadUserInfo() // 初始化时从本地存储恢复
}

const mutations = {
  SET_TOKEN: (state, token) => {
    state.token = token
  },
  SET_USER_INFO: (state, userInfo) => {
    state.userInfo = userInfo
    saveUserInfo(userInfo) // 保存用户信息到本地存储
  },
  CLEAR_USER: (state) => {
    state.token = null
    state.userInfo = null
    saveUserInfo(null) // 清除本地存储的用户信息
  },
  CLEAR_USER_INFO: (state) => {
    // 只清除用户信息，不清除token
    state.userInfo = null
    saveUserInfo(null) // 清除本地存储的用户信息
  }
}

const actions = {
  // 用户登录
  login({ commit }, userInfo) {
    const { username, password } = userInfo
    return new Promise((resolve, reject) => {
      login({ username: username.trim(), password: password })
        .then(response => {
          const { data } = response
          commit('SET_TOKEN', data.token)
          commit('SET_USER_INFO', {
            id: data.user_id,
            username: data.username,
            email: data.email,
            realName: data.real_name,
            role: data.role,
            avatar: data.avatar
          })
          setToken(data.token)
          resolve()
        })
        .catch(error => {
          reject(error)
        })
    })
  },

  // 获取用户信息
  getUserInfo({ commit, state }) {
    return new Promise((resolve, reject) => {
      if (!state.token) {
        reject('无效的token')
        return
      }
      
      // 如果已经有用户信息，直接返回
      if (state.userInfo) {
        resolve(state.userInfo)
        return
      }
      
      getUserInfo()
        .then(response => {
          const { data } = response
          if (!data) {
            reject('验证失败，请重新登录')
            return
          }
          
          commit('SET_USER_INFO', {
            id: data.id,
            username: data.username,
            email: data.email,
            realName: data.real_name,
            role: data.role,
            avatar: data.avatar
          })
          resolve(data)
        })
        .catch(error => {
          reject(error)
        })
    })
  },

  // 更新用户信息
  updateUserInfo({ commit }, userData) {
    return new Promise((resolve, reject) => {
      updateUserInfo(userData)
        .then(response => {
          if (response.code === 200 && response.data) {
            // 更新本地存储的用户信息
            commit('SET_USER_INFO', {
              id: response.data.id,
              username: response.data.username,
              email: response.data.email,
              realName: response.data.real_name,
              role: response.data.role,
              avatar: response.data.avatar
            })
          }
          resolve(response)
        })
        .catch(error => {
          reject(error)
        })
    })
  },

  // 用户退出登录
  logout({ commit }) {
    return new Promise(resolve => {
      commit('CLEAR_USER')
      removeToken()
      resolve()
    })
  },
  
  // 清除令牌
  resetToken({ commit }) {
    return new Promise(resolve => {
      commit('CLEAR_USER')
      removeToken()
      resolve()
    })
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
} 