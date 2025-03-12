<template>
  <div class="profile-container">
    <el-row :gutter="20">
      <el-col :xs="24" :md="8">
        <!-- 个人资料卡片 -->
        <el-card shadow="hover" class="profile-card">
          <div class="profile-card-header">
            <el-avatar :size="90" :src="user.profilePicture">
              {{ userInitials }}
            </el-avatar>
            <h2 class="user-name">{{ `${user.firstName} ${user.lastName}` }}</h2>
            <p class="user-role">{{ formatRole(user.role) }}</p>
          </div>
          
          <div class="profile-info">
            <div class="info-item">
              <span class="info-label">用户名</span>
              <span class="info-value">{{ user.username }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">邮箱</span>
              <span class="info-value">{{ user.email }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">部门</span>
              <span class="info-value">{{ user.department || '未设置' }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">加入时间</span>
              <span class="info-value">{{ $filters.formatDateOnly(user.createdAt) }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">上次登录</span>
              <span class="info-value">{{ $filters.formatDateOnly(user.lastLogin) }}</span>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :xs="24" :md="16">
        <!-- 编辑个人资料 -->
        <el-card shadow="hover" class="edit-profile-card">
          <template #header>
            <div class="card-header">
              <h3>个人资料设置</h3>
            </div>
          </template>
          
          <el-form 
            ref="profileForm" 
            :model="profileData" 
            :rules="rules" 
            label-position="top"
          >
            <el-alert 
              v-if="error" 
              :title="error" 
              type="error" 
              show-icon 
              :closable="false"
              style="margin-bottom: 20px"
            />
              
            <el-tabs>
              <el-tab-pane label="基本信息">
                <el-row :gutter="20">
                  <el-col :xs="24" :sm="12">
                    <el-form-item label="姓" prop="lastName">
                      <el-input v-model="profileData.lastName" />
                    </el-form-item>
                  </el-col>
                  <el-col :xs="24" :sm="12">
                    <el-form-item label="名" prop="firstName">
                      <el-input v-model="profileData.firstName" />
                    </el-form-item>
                  </el-col>
                </el-row>
                
                <el-form-item label="部门" prop="department">
                  <el-input v-model="profileData.department" />
                </el-form-item>
                
                <el-form-item label="头像" prop="profilePicture">
                  <el-input v-model="profileData.profilePicture" placeholder="头像URL" />
                </el-form-item>
              </el-tab-pane>
              
              <el-tab-pane label="修改密码">
                <el-form-item label="当前密码" prop="currentPassword">
                  <el-input 
                    v-model="profileData.currentPassword" 
                    type="password"
                    show-password
                  />
                </el-form-item>
                
                <el-form-item label="新密码" prop="newPassword">
                  <el-input 
                    v-model="profileData.newPassword" 
                    type="password"
                    show-password
                  />
                </el-form-item>
                
                <el-form-item label="确认新密码" prop="confirmPassword">
                  <el-input 
                    v-model="profileData.confirmPassword" 
                    type="password"
                    show-password
                  />
                </el-form-item>
              </el-tab-pane>
            </el-tabs>
            
            <div class="form-actions">
              <el-button 
                type="primary" 
                :loading="loading" 
                @click="submitForm"
              >
                保存更改
              </el-button>
            </div>
          </el-form>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import { ref, reactive, computed, onMounted } from 'vue'
import { useStore } from 'vuex'

export default {
  name: 'Profile',
  
  setup() {
    const store = useStore()
    const profileForm = ref(null)
    
    // 获取当前用户
    const user = computed(() => store.getters['auth/currentUser'] || {})
    
    // 用户首字母缩写
    const userInitials = computed(() => {
      if (user.value.firstName && user.value.lastName) {
        return `${user.value.firstName.charAt(0)}${user.value.lastName.charAt(0)}`.toUpperCase()
      }
      return user.value.username ? user.value.username.charAt(0).toUpperCase() : '?'
    })
    
    // 表单数据
    const profileData = reactive({
      firstName: '',
      lastName: '',
      department: '',
      profilePicture: '',
      currentPassword: '',
      newPassword: '',
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
      newPassword: [
        { min: 8, message: '密码长度不能少于8个字符', trigger: 'blur' }
      ],
      confirmPassword: [
        {
          validator: (rule, value, callback) => {
            if (profileData.newPassword && value !== profileData.newPassword) {
              callback(new Error('两次输入的密码不一致'))
            } else {
              callback()
            }
          },
          trigger: 'blur'
        }
      ]
    }
    
    // 加载状态和错误信息
    const loading = computed(() => store.getters['user/isLoading'])
    const error = computed(() => store.getters['user/error'])
    
    // 角色格式化
    const formatRole = (role) => {
      switch (role) {
        case 'admin':
          return '管理员'
        case 'user':
          return '普通用户'
        case 'viewer':
          return '仅查看用户'
        default:
          return role
      }
    }
    
    // 初始化表单
    const initForm = () => {
      // 填充个人资料表单
      profileData.firstName = user.value.firstName || ''
      profileData.lastName = user.value.lastName || ''
      profileData.department = user.value.department || ''
      profileData.profilePicture = user.value.profilePicture || ''
      
      // 清空密码字段
      profileData.currentPassword = ''
      profileData.newPassword = ''
      profileData.confirmPassword = ''
    }
    
    // 提交表单
    const submitForm = async () => {
      if (!profileForm.value) return
      
      await profileForm.value.validate(async (valid) => {
        if (valid) {
          // 构建更新数据
          const updateData = {}
          
          // 只包含修改过的基本信息字段
          if (profileData.firstName !== user.value.firstName) {
            updateData.firstName = profileData.firstName
          }
          
          if (profileData.lastName !== user.value.lastName) {
            updateData.lastName = profileData.lastName
          }
          
          if (profileData.department !== user.value.department) {
            updateData.department = profileData.department
          }
          
          if (profileData.profilePicture !== user.value.profilePicture) {
            updateData.profilePicture = profileData.profilePicture
          }
          
          // 如果有提供密码字段，添加到更新数据
          if (profileData.currentPassword && profileData.newPassword) {
            updateData.currentPassword = profileData.currentPassword
            updateData.newPassword = profileData.newPassword
          }
          
          // 只有在有数据更新时才发送请求
          if (Object.keys(updateData).length > 0) {
            try {
              await store.dispatch('user/updateProfile', updateData)
              initForm() // 重新加载表单
            } catch (error) {
              // 错误已在store中处理
            }
          } else {
            // 如果没有改动
            store.dispatch('addNotification', {
              type: 'info',
              message: '没有检测到任何改动',
              title: '信息'
            })
          }
        }
      })
    }
    
    // 页面加载时初始化
    onMounted(() => {
      // 获取最新的用户信息
      store.dispatch('user/fetchCurrentUser')
        .then(() => {
          console.log('用户信息已更新')
          initForm()
        })
        .catch(error => {
          console.error('获取用户信息失败:', error)
        })
    })
    
    return {
      user,
      userInitials,
      profileForm,
      profileData,
      rules,
      loading,
      error,
      formatRole,
      submitForm
    }
  }
}
</script>

<style lang="scss" scoped>
.profile-container {
  max-width: 1200px;
  margin: 0 auto;
}

.profile-card {
  margin-bottom: 20px;
  border-radius: 8px;
  overflow: hidden;
  
  .profile-card-header {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 24px;
    background-color: #f9fafc;
    border-bottom: 1px solid #f0f2f5;
  }
  
  .user-name {
    margin: 16px 0 4px;
    font-size: 18px;
    font-weight: 600;
  }
  
  .user-role {
    margin: 0;
    font-size: 14px;
    color: #909399;
  }
  
  .profile-info {
    padding: 24px;
  }
  
  .info-item {
    display: flex;
    justify-content: space-between;
    padding: 12px 0;
    border-bottom: 1px solid #f0f2f5;
    
    &:last-child {
      border-bottom: none;
    }
  }
  
  .info-label {
    color: #909399;
    font-size: 14px;
  }
  
  .info-value {
    font-weight: 500;
    font-size: 14px;
  }
}

.edit-profile-card {
  border-radius: 8px;
  overflow: hidden;
  
  .card-header {
    h3 {
      margin: 0;
      font-size: 16px;
      font-weight: 600;
    }
  }
}

.form-actions {
  margin-top: 30px;
  text-align: right;
}

@media (max-width: 768px) {
  .profile-card, .edit-profile-card {
    margin-bottom: 20px;
  }
}
</style> 