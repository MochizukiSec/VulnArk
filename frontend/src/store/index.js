import { createStore } from 'vuex'
import user from './modules/user'
import app from './modules/app'

export default createStore({
  modules: {
    user,
    app
  },
  getters: {
    token: state => state.user.token,
    userInfo: state => state.user.userInfo,
    userRole: state => state.user.userInfo ? state.user.userInfo.role : undefined,
    sidebar: state => state.app.sidebar,
    device: state => state.app.device
  }
}) 