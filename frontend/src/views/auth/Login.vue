<template>
  <div class="login-container">
    <div class="gradient-background"></div>
    <div class="login-form-container">
      <div class="login-header">
        <div class="logo-text">
          <span>VulnArk</span>
        </div>
        <h1 class="title">漏洞管理平台</h1>
        <p class="subtitle">安全管理 · 智能分析 · 高效防护</p>
      </div>

      <el-form 
        ref="loginForm" 
        :model="loginData" 
        :rules="rules" 
        @submit.prevent="handleLogin"
        class="login-form"
      >
        <el-alert 
          v-if="error" 
          :title="error" 
          type="error" 
          show-icon 
          :closable="false"
          class="login-alert"
        />

        <el-form-item prop="email" class="custom-form-item">
          <div class="input-with-icon">
            <svg class="input-icon email-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"></path>
              <polyline points="22,6 12,13 2,6"></polyline>
            </svg>
            <el-input 
              v-model="loginData.email" 
              placeholder="请输入邮箱" 
              type="email"
              class="custom-input"
            />
          </div>
        </el-form-item>

        <el-form-item prop="password" class="custom-form-item">
          <div class="input-with-icon">
            <svg class="input-icon password-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect>
              <path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
            </svg>
            <el-input 
              v-model="loginData.password" 
              placeholder="请输入密码" 
              type="password"
              show-password
              class="custom-input"
            />
          </div>
        </el-form-item>

        <div class="login-options">
          <el-checkbox v-model="rememberMe">记住我</el-checkbox>
        </div>

        <el-form-item>
          <button 
            type="button" 
            class="login-button"
            :class="{ 'loading': loading }"
            :disabled="loading"
            @click="submitForm"
          >
            <span class="button-text">{{ loading ? '登录中...' : '登录' }}</span>
            <i class="el-icon-right" v-if="!loading"></i>
            <span class="loading-spinner" v-if="loading"></span>
          </button>
        </el-form-item>
      </el-form>

      <div class="login-footer">
        <p class="copyright">© {{ new Date().getFullYear() }} 漏洞管理平台</p>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, reactive, computed, onMounted } from 'vue'
import { useStore } from 'vuex'
import { useRouter, useRoute } from 'vue-router'

export default {
  name: 'Login',
  
  setup() {
    const store = useStore()
    const router = useRouter()
    const route = useRoute()
    const loginForm = ref(null)
    const rememberMe = ref(false)
    
    // 表单数据
    const loginData = reactive({
      email: '',
      password: ''
    })
    
    // 表单验证规则
    const rules = {
      email: [
        { required: true, message: '请输入邮箱', trigger: 'blur' },
        { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
      ],
      password: [
        { required: true, message: '请输入密码', trigger: 'blur' },
        { min: 8, message: '密码长度不能少于8个字符', trigger: 'blur' }
      ]
    }
    
    // 从store中获取状态
    const loading = computed(() => store.getters['auth/authLoading'])
    const error = computed(() => store.getters['auth/authError'])
    
    // 提交表单
    const submitForm = async () => {
      if (!loginForm.value) return
      
      await loginForm.value.validate(async (valid) => {
        if (valid) {
          await handleLogin()
        }
      })
    }
    
    // 处理登录逻辑
    const handleLogin = async () => {
      const result = await store.dispatch('auth/login', loginData)
      
      if (result.success) {
        // 重定向到指定页面或默认到仪表盘
        const redirectPath = route.query.redirect || '/dashboard'
        router.push(redirectPath)
      }
    }
    
    // 初始化时清除错误
    onMounted(() => {
      store.dispatch('auth/authError', null)
    })
    
    return {
      loginForm,
      loginData,
      rules,
      loading,
      error,
      rememberMe,
      submitForm,
      handleLogin
    }
  }
}
</script>

<style lang="scss" scoped>
.login-container {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100vh;
  position: relative;
  overflow: hidden;
  background-color: #f0f2f5;
}

.gradient-background {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, #6366f1, #8b5cf6, #ec4899);
  opacity: 0.8;
  z-index: 0;
  background-size: 300% 300%;
  animation: gradient-shift 15s ease infinite;
}

@keyframes gradient-shift {
  0% {
    background-position: 0% 50%;
  }
  50% {
    background-position: 100% 50%;
  }
  100% {
    background-position: 0% 50%;
  }
}

.login-form-container {
  width: 100%;
  max-width: 440px;
  padding: 50px 40px;
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
  position: relative;
  z-index: 1;
  transform: translateY(0);
  transition: all 0.4s ease;
  
  &:hover {
    transform: translateY(-5px);
    box-shadow: 0 15px 35px rgba(0, 0, 0, 0.15);
  }
}

.login-header {
  text-align: center;
  margin-bottom: 40px;
  
  .logo-text {
    width: 120px;
    height: 80px;
    margin: 0 auto 20px;
    background: linear-gradient(135deg, #6366f1, #8b5cf6);
    color: white;
    border-radius: 20px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 24px;
    font-weight: bold;
    box-shadow: 0 10px 20px rgba(99, 102, 241, 0.3);
    position: relative;
    overflow: hidden;
    
    &::after {
      content: '';
      position: absolute;
      top: 0;
      left: -50%;
      width: 200%;
      height: 100%;
      background: rgba(255, 255, 255, 0.1);
      transform: skewX(45deg);
      transition: all 0.6s ease;
    }
    
    &:hover::after {
      left: 100%;
    }
  }
  
  .title {
    font-size: 28px;
    font-weight: 700;
    color: #1f2937;
    margin: 15px 0 10px;
    letter-spacing: 0.5px;
  }
  
  .subtitle {
    font-size: 16px;
    color: #6b7280;
    margin: 0;
    letter-spacing: 1px;
  }
}

.custom-form-item {
  margin-bottom: 25px;
}

.input-with-icon {
  position: relative;
  display: flex;
  align-items: center;
  width: 100%;
}

.input-icon {
  position: absolute;
  left: 18px;
  width: 18px;
  height: 18px;
  color: #6b7280;
  z-index: 2;
  stroke-width: 1.5;
  transition: all 0.3s ease;
  pointer-events: none;
}

.custom-input {
  width: 100%;
  
  :deep(.el-input__inner) {
    height: 50px;
    line-height: 50px;
    border-radius: 10px;
    background-color: #f9fafb;
    border: 1px solid #e5e7eb;
    padding-left: 48px;
    transition: all 0.3s ease;
    font-size: 16px;
    
    &:hover {
      border-color: #6366f1;
      
      & + .input-icon {
        color: #6366f1;
      }
    }
    
    &:focus {
      border-color: #6366f1;
      background-color: #ffffff;
      box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
    }

    &::placeholder {
      color: #9ca3af;
    }
  }
  
  :deep(.el-input__prefix) {
    display: none;
  }
  
  :deep(.el-input__suffix) {
    right: 15px;
  }
}

.input-with-icon:focus-within .input-icon {
  color: #4f46e5;
}

.login-options {
  display: flex;
  justify-content: flex-start;
  align-items: center;
  margin-bottom: 30px;
  
  :deep(.el-checkbox__input.is-checked) {
    .el-checkbox__inner {
      background-color: #6366f1;
      border-color: #6366f1;
    }
  }
  
  :deep(.el-checkbox__input.is-checked + .el-checkbox__label) {
    color: #4f46e5;
  }
  
  :deep(.el-checkbox__inner:hover) {
    border-color: #6366f1;
  }

  :deep(.el-checkbox__label) {
    font-size: 15px;
    color: #4b5563;
  }
}

.login-alert {
  margin-bottom: 25px;
  border-radius: 10px;
}

.login-button {
  width: 100%;
  height: 52px;
  background: linear-gradient(to right, #6366f1, #8b5cf6);
  border: none;
  border-radius: 10px;
  color: white;
  font-size: 17px;
  font-weight: 600;
  letter-spacing: 1px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
  box-shadow: 0 4px 15px rgba(99, 102, 241, 0.3);
  
  &:hover {
    background: linear-gradient(to right, #4f46e5, #7c3aed);
    box-shadow: 0 6px 20px rgba(99, 102, 241, 0.4);
  }
  
  &:active {
    transform: translateY(2px);
    box-shadow: 0 2px 10px rgba(99, 102, 241, 0.3);
  }
  
  &:disabled {
    opacity: 0.7;
    cursor: not-allowed;
  }
  
  .button-text {
    margin-right: 8px;
  }
  
  .loading-spinner {
    width: 20px;
    height: 20px;
    border: 2px solid rgba(255, 255, 255, 0.3);
    border-radius: 50%;
    border-top-color: #ffffff;
    animation: spin 0.8s linear infinite;
  }
  
  @keyframes spin {
    0% { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
  }
}

.login-footer {
  text-align: center;
  margin-top: 40px;
  
  .copyright {
    font-size: 14px;
    color: #6b7280;
    margin: 0;
  }
}

@media (max-width: 480px) {
  .login-form-container {
    padding: 40px 25px;
    margin: 0 20px;
    max-width: 100%;
  }
  
  .login-header {
    .logo-text {
      width: 110px;
      height: 70px;
      font-size: 20px;
    }
    
    .title {
      font-size: 24px;
    }
    
    .subtitle {
      font-size: 14px;
    }
  }
}
</style>
