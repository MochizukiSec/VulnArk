<template>
  <div class="login-container">
    <div class="login-content">
      <div class="login-form-wrapper">
        <transition name="fade">
          <div class="login-form">
            <div class="logo-container">
              <h1 class="login-title">VulnArk</h1>
              <p class="login-subtitle">漏洞管理平台</p>
            </div>
            <el-form :model="loginForm" :rules="loginRules" ref="loginFormRef" label-width="0">
              <el-form-item prop="username">
                <el-input 
                  v-model="loginForm.username" 
                  placeholder="用户名" 
                  prefix-icon="User"
                  @keyup.enter="handleLogin" 
                />
              </el-form-item>
              <el-form-item prop="password">
                <el-input 
                  v-model="loginForm.password" 
                  placeholder="密码" 
                  prefix-icon="Lock" 
                  type="password" 
                  show-password
                  @keyup.enter="handleLogin" 
                />
              </el-form-item>
              <el-form-item>
                <el-button 
                  :loading="loading" 
                  type="primary" 
                  class="login-button" 
                  @click="handleLogin"
                >
                  <span v-if="!loading">登录</span>
                  <span v-else>加载中...</span>
                </el-button>
              </el-form-item>
            </el-form>
          </div>
        </transition>
      </div>
      
      <div class="login-background">
        <div class="bg-shapes">
          <div class="shape shape-1"></div>
          <div class="shape shape-2"></div>
          <div class="shape shape-3"></div>
          <div class="shape shape-4"></div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useStore } from 'vuex'
import { ElMessage } from 'element-plus'

const loginFormRef = ref(null)
const router = useRouter()
const route = useRoute()
const store = useStore()

const loginForm = reactive({
  username: '',
  password: ''
})

const loading = ref(false)

// 验证规则
const loginRules = reactive({
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
})

const handleLogin = () => {
  loginFormRef.value.validate(valid => {
    if (valid) {
      loading.value = true
      store.dispatch('user/login', loginForm).then(() => {
        loading.value = false
        const redirectPath = route.query.redirect || '/'
        router.push(redirectPath)
        ElMessage.success('登录成功')
      }).catch(error => {
        console.error('Login failed:', error)
        loading.value = false
        ElMessage.error(error.message || '登录失败，请检查用户名和密码')
      })
    }
  })
}
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #2b5876 0%, #4e4376 100%);
  overflow: hidden;
}

.login-content {
  position: relative;
  width: 100%;
  max-width: 1200px;
  display: flex;
  justify-content: center;
  z-index: 1;
}

.login-form-wrapper {
  width: 400px;
  position: relative;
  z-index: 2;
}

.login-form {
  background-color: rgba(255, 255, 255, 0.95);
  padding: 40px;
  border-radius: 12px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
  animation: float 6s ease-in-out infinite;
  backdrop-filter: blur(10px);
}

.logo-container {
  text-align: center;
  margin-bottom: 30px;
}

.login-title {
  font-size: 32px;
  font-weight: 600;
  margin: 0;
  background: linear-gradient(90deg, #4e54c8, #8f94fb);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  animation: title-glow 1.5s ease-in-out infinite alternate;
}

.login-subtitle {
  font-size: 16px;
  color: #606266;
  margin-top: 8px;
}

.login-button {
  width: 100%;
  height: 44px;
  border-radius: 22px;
  font-size: 16px;
  font-weight: 500;
  letter-spacing: 1px;
  background: linear-gradient(90deg, #4e54c8, #8f94fb);
  border: none;
  transition: all 0.3s cubic-bezier(0.645, 0.045, 0.355, 1);
}

.login-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(143, 148, 251, 0.5);
}

/* 背景动画 */
.login-background {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  overflow: hidden;
}

.bg-shapes {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
}

.shape {
  position: absolute;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.1);
}

.shape-1 {
  width: 300px;
  height: 300px;
  top: -150px;
  right: -50px;
  animation: floating 15s infinite linear;
}

.shape-2 {
  width: 200px;
  height: 200px;
  bottom: -100px;
  left: 10%;
  animation: floating 20s infinite linear reverse;
}

.shape-3 {
  width: 150px;
  height: 150px;
  top: 30%;
  left: -75px;
  animation: floating 25s infinite linear;
}

.shape-4 {
  width: 100px;
  height: 100px;
  bottom: 20%;
  right: 10%;
  animation: floating 30s infinite linear reverse;
}

@keyframes floating {
  from {
    transform: rotate(0deg) translateY(0) rotate(0deg);
  }
  to {
    transform: rotate(360deg) translateY(-10px) rotate(-360deg);
  }
}

@keyframes float {
  0% {
    transform: translateY(0px);
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
  }
  50% {
    transform: translateY(-10px);
    box-shadow: 0 15px 35px rgba(0, 0, 0, 0.25);
  }
  100% {
    transform: translateY(0px);
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
  }
}

@keyframes title-glow {
  from {
    text-shadow: 0 0 5px rgba(78, 84, 200, 0.3), 0 0 10px rgba(78, 84, 200, 0);
  }
  to {
    text-shadow: 0 0 10px rgba(78, 84, 200, 0.5), 0 0 20px rgba(78, 84, 200, 0.3);
  }
}

/* 响应式设计 */
@media (max-width: 576px) {
  .login-form-wrapper {
    width: 90%;
  }
  
  .login-form {
    padding: 30px 20px;
  }
}
</style> 