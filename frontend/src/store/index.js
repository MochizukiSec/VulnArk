import { createStore } from 'vuex'
import auth from './modules/auth'
import vulnerability from './modules/vulnerability'
import dashboard from './modules/dashboard'
import user from './modules/user'
import report from './modules/report'

export default createStore({
  state: {
    loading: false,
    error: null,
    notifications: []
  },
  
  getters: {
    isLoading: state => state.loading,
    error: state => state.error,
    notifications: state => state.notifications
  },
  
  mutations: {
    SET_LOADING(state, loading) {
      state.loading = loading
    },
    
    SET_ERROR(state, error) {
      state.error = error
    },
    
    CLEAR_ERROR(state) {
      state.error = null
    },
    
    ADD_NOTIFICATION(state, notification) {
      state.notifications.push({
        id: Date.now(),
        ...notification,
        read: false,
        timestamp: new Date()
      })
    },
    
    MARK_NOTIFICATION_READ(state, id) {
      const notification = state.notifications.find(n => n.id === id)
      if (notification) {
        notification.read = true
      }
    },
    
    CLEAR_NOTIFICATIONS(state) {
      state.notifications = []
    }
  },
  
  actions: {
    setLoading({ commit }, loading) {
      commit('SET_LOADING', loading)
    },
    
    setError({ commit }, error) {
      commit('SET_ERROR', error)
    },
    
    clearError({ commit }) {
      commit('CLEAR_ERROR')
    },
    
    addNotification({ commit }, notification) {
      commit('ADD_NOTIFICATION', notification)
    },
    
    markNotificationRead({ commit }, id) {
      commit('MARK_NOTIFICATION_READ', id)
    },
    
    clearNotifications({ commit }) {
      commit('CLEAR_NOTIFICATIONS')
    }
  },
  
  modules: {
    auth,
    vulnerability,
    dashboard,
    user,
    report
  }
}) 