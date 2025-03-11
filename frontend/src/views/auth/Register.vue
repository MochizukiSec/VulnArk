<template>
  <div class="register-container">
    <!-- 背景装饰 -->
    <div class="bg-decoration">
      <div class="bg-shape shape-1"></div>
      <div class="bg-shape shape-2"></div>
      <div class="bg-shape shape-3"></div>
      <div class="bg-dots"></div>
    </div>
    
    <div class="register-card">
      <div class="register-header">
        <h1 class="register-title">漏洞管理平台</h1>
        <p class="register-subtitle">创建新账号</p>
      </div>
      
      <el-form 
        ref="registerForm" 
        :model="registerData" 
        :rules="rules" 
        label-position="top" 
        @submit.prevent="handleRegister"
      >
        <el-alert 
          v-if="error" 
          :title="error" 
          type="error" 
          show-icon 
          :closable="false"
          style="margin-bottom: 20px"
        />
        
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="姓" prop="lastName">
              <el-input 
                v-model="registerData.lastName" 
                placeholder="请输入姓"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="名" prop="firstName">
              <el-input 
                v-model="registerData.firstName" 
                placeholder="请输入名"
              />
            </el-form-item>
          </el-col>
        </el-row>
        
        <el-form-item label="用户名" prop="username">
          <el-input 
            v-model="registerData.username" 
            placeholder="请输入用户名"
          />
        </el-form-item>
        
        <el-form-item label="邮箱" prop="email">
          <el-input 
            v-model="registerData.email" 
            type="email"
            placeholder="请输入邮箱"
          />
        </el-form-item>
        
        <el-form-item label="部门" prop="department">
          <el-input 
            v-model="registerData.department" 
            placeholder="请输入部门(可选)"
          />
        </el-form-item>
        
        <el-form-item label="密码" prop="password">
          <el-input 
            v-model="registerData.password" 
            type="password"
            placeholder="请输入密码"
            show-password
          />
        </el-form-item>
        
        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input 
            v-model="registerData.confirmPassword" 
            type="password"
            placeholder="请再次输入密码"
            show-password
          />
        </el-form-item>
        
        <div class="form-actions">
          <el-button 
            type="primary" 
            :loading="loading" 
            @click="submitForm" 
            class="submit-button"
          >
            注册账号 <i class="el-icon-arrow-right"></i>
          </el-button>
        </div>
        
        <div class="form-footer">
          <p>已有账号？<router-link to="/login">返回登录</router-link></p>
          <p class="home-link"><router-link to="/">返回首页</router-link></p>
        </div>
      </el-form>
    </div>
  </div>
</template>

<script>
import { ref, reactive, computed } from 'vue'
import { useStore } from 'vuex'
import { useRouter } from 'vue-router'

export default {
  name: 'Register',
  
  setup() {
    const store = useStore()
    const router = useRouter()
    const registerForm = ref(null)
    
    // 表单数据
    const registerData = reactive({
      firstName: '',
      lastName: '',
      username: '',
      email: '',
      department: '',
      password: '',
      confirmPassword: ''
    })
    
    // 表单验证规则
    const rules = {
      firstName: [
        { required: true, message: '请输入名', trigger: 'blur' }
      ],
      lastName: [
        { required: true, message: '请输入姓', trigger: 'blur' }
      ],
      username: [
        { required: true, message: '请输入用户名', trigger: 'blur' },
        { min: 3, max: 30, message: '用户名长度在3到30个字符之间', trigger: 'blur' }
      ],
      email: [
        { required: true, message: '请输入邮箱', trigger: 'blur' },
        { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
      ],
      password: [
        { required: true, message: '请输入密码', trigger: 'blur' },
        { min: 8, message: '密码长度不能少于8个字符', trigger: 'blur' }
      ],
      confirmPassword: [
        { required: true, message: '请再次输入密码', trigger: 'blur' },
        {
          validator: (rule, value, callback) => {
            if (value !== registerData.password) {
              callback(new Error('两次输入的密码不一致'))
            } else {
              callback()
            }
          },
          trigger: 'blur'
        }
      ]
    }
    
    // 从store中获取状态
    const loading = computed(() => store.getters['auth/authLoading'])
    const error = computed(() => store.getters['auth/authError'])
    
    // 提交表单
    const submitForm = async () => {
      if (!registerForm.value) return
      
      await registerForm.value.validate(async (valid) => {
        if (valid) {
          await handleRegister()
        }
      })
    }
    
    // 处理注册逻辑
    const handleRegister = async () => {
      // 创建一个符合后端期望的数据结构
      const userData = {
        firstName: registerData.firstName,
        lastName: registerData.lastName,
        username: registerData.username,
        email: registerData.email,
        department: registerData.department,
        password: registerData.password
      }
      
      const result = await store.dispatch('auth/register', userData)
      
      if (result.success) {
        router.push('/dashboard')
      }
    }
    
    return {
      registerForm,
      registerData,
      rules,
      loading,
      error,
      submitForm,
      handleRegister
    }
  }
}
</script>

<style lang="scss" scoped>
.register-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(125deg, #000428 0%, #004e92 100%);
  padding: 20px;
  position: relative;
  overflow: hidden;
  
  &::before, &::after {
    content: '';
    position: absolute;
    border-radius: 50%;
    filter: blur(40px);
    pointer-events: none;
  }
  
  &::before {
    width: 30vw;
    height: 30vw;
    background: rgba(255, 0, 128, 0.2);
    top: -5%;
    right: -5%;
  }
  
  &::after {
    width: 35vw;
    height: 35vw;
    background: rgba(0, 255, 255, 0.2);
    bottom: -10%;
    left: -10%;
  }
}

/* 添加背景装饰 */
.bg-decoration {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  z-index: 0;
  overflow: hidden;
}

.bg-shape {
  position: absolute;
  border-radius: 50%;
  filter: blur(60px);
  opacity: 0.3;
  pointer-events: none;
}

.shape-1 {
  width: 500px;
  height: 500px;
  background: linear-gradient(135deg, #ff00cc, transparent);
  top: -200px;
  right: -100px;
  animation: float 20s ease-in-out infinite alternate;
}

.shape-2 {
  width: 600px;
  height: 600px;
  background: linear-gradient(135deg, #00ffff, transparent);
  bottom: -300px;
  left: -200px;
  animation: float 25s ease-in-out infinite alternate;
  animation-delay: -5s;
}

.shape-3 {
  width: 400px;
  height: 400px;
  background: linear-gradient(135deg, #9900ff, transparent);
  top: 40%;
  right: 30%;
  animation: float 18s ease-in-out infinite alternate;
  animation-delay: -10s;
}

.bg-dots {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-image: radial-gradient(rgba(255, 255, 255, 0.1) 1px, transparent 1px);
  background-size: 30px 30px;
  opacity: 0.2;
  pointer-events: none;
}

@keyframes float {
  0% {
    transform: translate(0, 0) rotate(0deg);
  }
  100% {
    transform: translate(50px, 50px) rotate(10deg);
  }
}

.register-card {
  width: 100%;
  max-width: 650px;
  padding: 40px;
  background: rgba(20, 20, 40, 0.5);
  backdrop-filter: blur(10px);
  border-radius: 24px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.4), 
              0 0 0 1px rgba(255, 255, 255, 0.1) inset, 
              0 0 30px rgba(0, 255, 255, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.1);
  position: relative;
  z-index: 1;
  animation: fadeIn 0.8s ease-out;
}

.register-header {
  text-align: center;
  margin-bottom: 30px;
}

.register-title {
  font-size: 32px;
  font-weight: 800;
  margin-bottom: 16px;
  line-height: 1.2;
  letter-spacing: -0.5px;
  background: linear-gradient(to right, #ffffff, #00ffff);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  text-shadow: 0 2px 10px rgba(0, 255, 255, 0.3);
  position: relative;
  display: inline-block;
}

.register-subtitle {
  font-size: 18px;
  color: rgba(255, 255, 255, 0.85);
  margin-bottom: 10px;
  font-weight: 400;
}

.form-actions {
  margin-top: 40px;
}

.submit-button {
  width: 100%;
  height: 50px;
  font-size: 16px;
  font-weight: 600;
  letter-spacing: 0.8px;
  color: white;
  border: none;
  border-radius: 16px;
  background: linear-gradient(90deg, #ff00cc, #00ffff);
  cursor: pointer;
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
  box-shadow: 0 10px 20px rgba(0, 0, 0, 0.2);
  
  &:hover {
    transform: translateY(-3px);
    box-shadow: 0 15px 30px rgba(0, 0, 0, 0.3),
                0 0 20px rgba(0, 255, 255, 0.3);
  }
  
  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: -100%;
    width: 100%;
    height: 100%;
    background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
    transition: all 0.6s ease;
    z-index: 1;
  }
  
  &:hover::before {
    left: 100%;
  }
}

.form-footer {
  margin-top: 30px;
  text-align: center;
  color: rgba(255, 255, 255, 0.7);
  
  a {
    color: #00ffff;
    text-decoration: none;
    font-weight: 600;
    position: relative;
    
    &::after {
      content: '';
      position: absolute;
      bottom: -2px;
      left: 0;
      width: 100%;
      height: 1px;
      background: linear-gradient(90deg, #ff00cc, #00ffff);
      transform: scaleX(0);
      transform-origin: right;
      transition: transform 0.3s ease;
    }
    
    &:hover::after {
      transform: scaleX(1);
      transform-origin: left;
    }
  }
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 美化表单控件 */
:deep(.el-form-item__label) {
  color: rgba(255, 255, 255, 0.9);
  font-weight: 500;
}

:deep(.el-input__inner) {
  background: rgba(0, 0, 0, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 12px;
  color: rgba(255, 255, 255, 0.9);
  padding: 12px 15px;
  height: 50px;
  transition: all 0.3s ease;
  
  &:focus {
    background: rgba(0, 0, 0, 0.3);
    border-color: rgba(0, 255, 255, 0.5);
    box-shadow: 0 0 15px rgba(0, 255, 255, 0.2);
  }
  
  &::placeholder {
    color: rgba(255, 255, 255, 0.5);
  }
}

:deep(.el-input__suffix) {
  color: rgba(255, 255, 255, 0.7);
}

:deep(.el-alert) {
  background: rgba(255, 73, 73, 0.2);
  color: #ff9999;
  border-radius: 12px;
  border: 1px solid rgba(255, 73, 73, 0.3);
  
  .el-alert__title {
    color: #ffc0c0;
  }
  
  .el-alert__icon {
    color: #ff9999;
  }
}

/* 移动端适配 */
@media (max-width: 768px) {
  .register-card {
    padding: 30px 20px;
  }
  
  .register-title {
    font-size: 28px;
  }
  
  .register-subtitle {
    font-size: 16px;
  }
  
  :deep(.el-input__inner) {
    height: 45px;
  }
}
</style> 