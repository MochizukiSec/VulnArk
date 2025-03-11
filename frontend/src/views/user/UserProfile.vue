<template>
  <div class="user-profile-container">
    <div class="page-header">
      <h1 class="page-title">个人资料</h1>
      <p class="page-subtitle">查看和管理您的个人账户信息</p>
    </div>

    <el-row :gutter="20">
      <!-- 左侧个人信息卡片 -->
      <el-col :xs="24" :sm="24" :md="8" :lg="7" :xl="6">
        <el-card class="profile-card" shadow="hover">
          <div class="profile-header">
            <div class="profile-avatar">
              <el-avatar 
                :size="100" 
                :src="userAvatar"
                icon="el-icon-user">
              </el-avatar>
              <div class="avatar-upload">
                <el-button type="primary" size="small" round plain @click="handleAvatarClick">
                  <i class="el-icon-camera"></i> 更换头像
                </el-button>
                <input ref="avatarInput" type="file" accept="image/*" style="display: none" @change="handleAvatarUpload" />
              </div>
            </div>
            <div class="profile-name">
              <h2>{{ userFullName }}</h2>
              <span class="user-role">{{ userRoleText }}</span>
            </div>
          </div>
          <div class="profile-info-list">
            <div class="info-item">
              <i class="el-icon-user"></i>
              <span class="item-label">用户名:</span>
              <span class="item-value">{{ getUserField('username', [], '未设置') }}</span>
            </div>
            <div class="info-item">
              <i class="el-icon-message"></i>
              <span class="item-label">邮箱:</span>
              <span class="item-value">{{ getUserField('email', [], '未设置') }}</span>
            </div>
            <div class="info-item">
              <i class="el-icon-office-building"></i>
              <span class="item-label">部门:</span>
              <span class="item-value">{{ getUserField('department', [], '未设置') }}</span>
            </div>
            <div class="info-item">
              <i class="el-icon-date"></i>
              <span class="item-label">注册时间:</span>
              <span class="item-value">{{ formatDate(getUserField('createdAt', ['created_at', 'createTime'])) }}</span>
            </div>
            <div class="info-item">
              <i class="el-icon-time"></i>
              <span class="item-label">上次登录:</span>
              <span class="item-value">{{ formatDate(getUserField('lastLogin', ['last_login', 'loginTime'])) }}</span>
            </div>
          </div>
          <div class="profile-actions">
            <el-button 
              type="danger" 
              round 
              plain 
              size="small" 
              @click="logout"
              class="logout-btn">
              <i class="el-icon-switch-button"></i> 退出登录
            </el-button>
          </div>
        </el-card>
      </el-col>

      <!-- 右侧编辑表单和安全设置 -->
      <el-col :xs="24" :sm="24" :md="16" :lg="17" :xl="18">
        <el-tabs v-model="activeTab" class="profile-tabs">
          <el-tab-pane label="个人信息" name="basic">
            <el-card shadow="hover">
              <template #header>
                <div class="card-header">
                  <span>个人信息</span>
                  <el-button 
                    type="primary" 
                    plain 
                    size="small" 
                    @click="submitProfileForm"
                    :loading="loading.profile"
                  >保存更改</el-button>
                </div>
              </template>

              <el-form 
                ref="profileForm" 
                :model="profileData" 
                :rules="profileRules" 
                label-position="top" 
                @submit.prevent="submitProfileForm"
              >
                <el-row :gutter="20">
                  <el-col :xs="24" :sm="12">
                    <el-form-item label="名字" prop="firstName">
                      <el-input v-model="profileData.firstName" placeholder="请输入名字" />
                    </el-form-item>
                  </el-col>
                  <el-col :xs="24" :sm="12">
                    <el-form-item label="姓氏" prop="lastName">
                      <el-input v-model="profileData.lastName" placeholder="请输入姓氏" />
                    </el-form-item>
                  </el-col>
                </el-row>

                <el-form-item label="用户名" prop="username">
                  <el-input v-model="profileData.username" placeholder="请输入用户名" />
                </el-form-item>

                <el-form-item label="邮箱" prop="email">
                  <el-input v-model="profileData.email" type="email" placeholder="请输入邮箱" />
                </el-form-item>

                <el-form-item label="部门" prop="department">
                  <el-input v-model="profileData.department" placeholder="请输入部门名称" />
                </el-form-item>
              </el-form>
            </el-card>
          </el-tab-pane>

          <el-tab-pane label="安全设置" name="security">
            <el-card shadow="hover">
              <template #header>
                <div class="card-header">
                  <span>修改密码</span>
                  <el-button 
                    type="primary" 
                    plain 
                    size="small"
                    @click="submitPasswordForm"
                    :loading="loading.password"
                  >更新密码</el-button>
                </div>
              </template>

              <el-form 
                ref="passwordForm" 
                :model="passwordData" 
                :rules="passwordRules"  
                label-position="top"
                @submit.prevent="submitPasswordForm"
              >
                <el-form-item label="当前密码" prop="currentPassword">
                  <el-input 
                    v-model="passwordData.currentPassword" 
                    type="password" 
                    show-password
                    placeholder="请输入当前密码" 
                  />
                </el-form-item>

                <el-form-item label="新密码" prop="newPassword">
                  <el-input 
                    v-model="passwordData.newPassword" 
                    type="password" 
                    show-password
                    placeholder="请输入新密码（至少6个字符）" 
                  />
                </el-form-item>

                <el-form-item label="确认新密码" prop="confirmPassword">
                  <el-input 
                    v-model="passwordData.confirmPassword" 
                    type="password" 
                    show-password
                    placeholder="请再次输入新密码" 
                  />
                </el-form-item>
              </el-form>
            </el-card>

            <el-card shadow="hover" class="mt-20">
              <template #header>
                <div class="card-header">
                  <span>账户安全</span>
                </div>
              </template>

              <div class="security-items">
                <div class="security-item">
                  <div class="security-item-content">
                    <div class="security-item-icon">
                      <i class="el-icon-lock"></i>
                    </div>
                    <div class="security-item-text">
                      <h3>双因素认证</h3>
                      <p>增强账户安全性，需要额外的验证步骤才能登录</p>
                    </div>
                  </div>
                  <el-switch v-model="securitySettings.twoFactor" :disabled="true" />
                </div>

                <div class="security-item">
                  <div class="security-item-content">
                    <div class="security-item-icon">
                      <i class="el-icon-bell"></i>
                    </div>
                    <div class="security-item-text">
                      <h3>登录通知</h3>
                      <p>当有新设备登录您的账户时发送通知</p>
                    </div>
                  </div>
                  <el-switch v-model="securitySettings.loginNotifications" :disabled="true" />
                </div>
              </div>
            </el-card>
          </el-tab-pane>

          <el-tab-pane label="通知设置" name="notifications">
            <el-card shadow="hover">
              <template #header>
                <div class="card-header">
                  <span>通知偏好</span>
                  <el-button type="primary" plain size="small" :disabled="true">保存设置</el-button>
                </div>
              </template>

              <div class="notification-item">
                <div class="notification-text">
                  <h3>漏洞更新通知</h3>
                  <p>当有新漏洞添加或现有漏洞状态变更时通知我</p>
                </div>
                <el-switch v-model="notificationSettings.vulnerabilities" :disabled="true" />
              </div>

              <div class="notification-item">
                <div class="notification-text">
                  <h3>周报告摘要</h3>
                  <p>每周发送漏洞管理工作摘要</p>
                </div>
                <el-switch v-model="notificationSettings.weeklyReport" :disabled="true" />
              </div>

              <div class="notification-item">
                <div class="notification-text">
                  <h3>系统公告</h3>
                  <p>接收关于系统更新和维护的消息</p>
                </div>
                <el-switch v-model="notificationSettings.systemAnnouncements" :disabled="true" />
              </div>
            </el-card>
          </el-tab-pane>
        </el-tabs>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import { ref, reactive, computed, onMounted } from 'vue'
import { useStore } from 'vuex'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import axios from 'axios'

export default {
  name: 'UserProfile',
  
  setup() {
    const store = useStore()
    const router = useRouter()
    const profileForm = ref(null)
    const passwordForm = ref(null)
    const avatarInput = ref(null)
    
    const activeTab = ref('basic')
    
    const user = computed(() => {
      const currentUser = store.getters['auth/currentUser'] || {}
      console.log('当前用户信息:', currentUser)
      return currentUser
    })
    
    // 辅助函数 - 从用户对象获取字段值（支持多种可能的字段名格式）
    const getUserField = (fieldName, aliases = [], defaultValue = '') => {
      if (!user.value) return defaultValue
      
      // 检查主字段名
      if (user.value[fieldName] !== undefined) {
        return user.value[fieldName]
      }
      
      // 检查别名字段
      for (const alias of aliases) {
        if (user.value[alias] !== undefined) {
          return user.value[alias]
        }
      }
      
      return defaultValue
    }
    
    const userFullName = computed(() => {
      const firstName = getUserField('firstName', ['first_name', 'firstname'])
      const lastName = getUserField('lastName', ['last_name', 'lastname'])
      
      const fullName = `${firstName} ${lastName}`.trim()
      
      if (fullName) {
        return fullName
      }
      
      return getUserField('username', ['name'], '用户')
    })

    const userRoleText = computed(() => {
      const roleMap = {
        admin: '管理员',
        user: '普通用户',
        viewer: '只读用户'
      }
      const role = getUserField('role', [])
      return roleMap[role] || '用户'
    })
    
    const userAvatar = computed(() => {
      return getUserField('profilePicture', ['profile_picture', 'avatar'], '')
    })
    
    // 表单数据
    const profileData = reactive({
      firstName: '',
      lastName: '',
      username: '',
      email: '',
      department: ''
    })
    
    // 密码表单
    const passwordData = reactive({
      currentPassword: '',
      newPassword: '',
      confirmPassword: ''
    })
    
    // 安全设置
    const securitySettings = reactive({
      twoFactor: false,
      loginNotifications: true
    })
    
    // 通知设置
    const notificationSettings = reactive({
      vulnerabilities: true,
      weeklyReport: true,
      systemAnnouncements: true
    })
    
    // 加载状态
    const loading = reactive({
      profile: false,
      password: false,
      avatar: false
    })
    
    // 表单验证规则
    const profileRules = {
      firstName: [
        { required: true, message: '请输入名字', trigger: 'blur' }
      ],
      lastName: [
        { required: true, message: '请输入姓氏', trigger: 'blur' }
      ],
      username: [
        { required: true, message: '请输入用户名', trigger: 'blur' },
        { min: 3, max: 20, message: '用户名长度必须在3到20个字符之间', trigger: 'blur' }
      ],
      email: [
        { required: true, message: '请输入邮箱', trigger: 'blur' },
        { type: 'email', message: '请输入有效的邮箱地址', trigger: 'blur' }
      ]
    }
    
    const passwordRules = {
      currentPassword: [
        { required: true, message: '请输入当前密码', trigger: 'blur' }
      ],
      newPassword: [
        { required: true, message: '请输入新密码', trigger: 'blur' },
        { min: 6, message: '密码长度至少为6个字符', trigger: 'blur' }
      ],
      confirmPassword: [
        { required: true, message: '请再次输入新密码', trigger: 'blur' },
        {
          validator: (rule, value, callback) => {
            if (value !== passwordData.newPassword) {
              callback(new Error('两次输入的密码不一致'))
            } else {
              callback()
            }
          },
          trigger: 'blur'
        }
      ]
    }
    
    // 初始化表单数据
    const initFormData = () => {
      if (!user.value) return
      console.log('初始化表单数据，用户信息:', user.value)
      
      profileData.firstName = getUserField('firstName', ['first_name', 'firstname'])
      profileData.lastName = getUserField('lastName', ['last_name', 'lastname'])
      profileData.username = getUserField('username', [])
      profileData.email = getUserField('email', [])
      profileData.department = getUserField('department', [])
    }
    
    // 提交个人信息表单
    const submitProfileForm = async () => {
      if (!profileForm.value) return
      
      try {
        await profileForm.value.validate()
        
        loading.profile = true
        
        // 调用API更新个人信息
        const response = await axios.put('/users/me', {
          firstName: profileData.firstName,
          lastName: profileData.lastName,
          department: profileData.department
        })
        
        if (response.data) {
          // 更新 Vuex 存储的用户信息
          store.dispatch('auth/updateUserInfo', {
            first_name: profileData.firstName,
            last_name: profileData.lastName,
            firstName: profileData.firstName,
            lastName: profileData.lastName,
            department: profileData.department
          })
          
          ElMessage({
            type: 'success',
            message: '个人信息更新成功'
          })
        }
      } catch (error) {
        console.error('更新个人信息失败:', error)
        ElMessage({
          type: 'error',
          message: error.response?.data?.error || '更新个人信息失败，请稍后重试'
        })
      } finally {
        loading.profile = false
      }
    }
    
    // 提交密码修改表单
    const submitPasswordForm = async () => {
      if (!passwordForm.value) return
      
      try {
        await passwordForm.value.validate()
        
        loading.password = true
        
        // 调用API更新密码
        const response = await axios.put('/users/me', {
          currentPassword: passwordData.currentPassword,
          newPassword: passwordData.newPassword
        })
        
        if (response.data) {
          // 清空表单
          passwordData.currentPassword = ''
          passwordData.newPassword = ''
          passwordData.confirmPassword = ''
          
          ElMessage({
            type: 'success',
            message: '密码修改成功'
          })
        }
      } catch (error) {
        console.error('修改密码失败:', error)
        ElMessage({
          type: 'error',
          message: error.response?.data?.error || '修改密码失败，请确认当前密码是否正确'
        })
      } finally {
        loading.password = false
      }
    }
    
    // 头像上传相关
    const handleAvatarClick = () => {
      avatarInput.value.click()
    }
    
    const handleAvatarUpload = async (event) => {
      const file = event.target.files[0]
      if (!file) return
      
      // 文件类型和大小验证
      const isImage = file.type.startsWith('image/')
      const maxSizeMB = 2
      const maxSize = maxSizeMB * 1024 * 1024
      
      if (!isImage) {
        ElMessage.error('请上传图片文件')
        return
      }
      
      if (file.size > maxSize) {
        ElMessage.error(`图片大小不能超过${maxSizeMB}MB`)
        return
      }
      
      // 暂时显示提示，因为后端API尚未实现
      ElMessage({
        type: 'info',
        message: '头像上传功能即将推出，敬请期待！'
      })
      
      // 清空文件输入以允许再次选择相同文件
      event.target.value = ''
      
      /* 
      // 以下代码暂时注释，等待后端API实现
      try {
        loading.avatar = true
        
        const formData = new FormData()
        formData.append('avatar', file)
        
        // 调用API上传头像
        const response = await axios.post('/users/upload-avatar', formData, {
          headers: {
            'Content-Type': 'multipart/form-data'
          }
        })
        
        if (response.data.success) {
          // 更新用户信息中的头像
          store.dispatch('auth/updateUserInfo', {
            profilePicture: response.data.profilePicture
          })
          
          ElMessage({
            type: 'success',
            message: '头像上传成功'
          })
        }
      } catch (error) {
        console.error('头像上传失败:', error)
        ElMessage({
          type: 'error',
          message: error.response?.data?.error || '头像上传失败，请稍后重试'
        })
      } finally {
        loading.avatar = false
        // 清空文件输入以允许再次选择相同文件
        event.target.value = ''
      }
      */
    }
    
    // 时间格式化
    const formatDate = (dateString) => {
      if (!dateString) return '未知'
      
      const date = new Date(dateString)
      return new Intl.DateTimeFormat('zh-CN', {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      }).format(date)
    }
    
    // 退出登录
    const logout = () => {
      ElMessageBox.confirm('确定要退出登录吗?', '退出确认', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        store.dispatch('auth/logout')
        router.push('/login')
        ElMessage({
          type: 'success',
          message: '已安全退出系统'
        })
      }).catch(() => {})
    }
    
    onMounted(() => {
      initFormData()
    })
    
    return {
      user,
      userFullName,
      userRoleText,
      userAvatar,
      getUserField,
      activeTab,
      profileForm,
      profileData,
      passwordForm,
      passwordData,
      securitySettings,
      notificationSettings,
      loading,
      profileRules,
      passwordRules,
      avatarInput,
      formatDate,
      submitProfileForm,
      submitPasswordForm,
      handleAvatarClick,
      handleAvatarUpload,
      logout
    }
  }
}
</script>

<style lang="scss" scoped>
.user-profile-container {
  padding: 20px;
  max-width: 1400px;
  margin: 0 auto;
}

.profile-card {
  transition: all 0.3s ease;
  height: 100%;
  
  &:hover {
    transform: translateY(-5px);
  }
}

.profile-header {
  text-align: center;
  padding-bottom: 20px;
  margin-bottom: 20px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
}

.profile-avatar {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 15px;
  margin-bottom: 20px;
  
  .avatar-upload {
    opacity: 0;
    transition: opacity 0.2s ease;
  }
  
  &:hover .avatar-upload {
    opacity: 1;
  }
}

.profile-name {
  h2 {
    margin: 0 0 10px;
    font-size: 22px;
    color: #303133;
  }
  
  .user-role {
    background: linear-gradient(90deg, #36d1dc, #5b86e5);
    color: white;
    font-size: 12px;
    font-weight: 600;
    padding: 3px 12px;
    border-radius: 100px;
    display: inline-block;
  }
}

.profile-info-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
  margin-bottom: 20px;
}

.info-item {
  display: flex;
  align-items: center;
  gap: 10px;
  
  i {
    font-size: 18px;
    color: #409EFF;
    width: 24px;
    height: 24px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: rgba(64, 158, 255, 0.1);
    border-radius: 50%;
  }
  
  .item-label {
    color: #606266;
    font-size: 14px;
    font-weight: 500;
    width: 60px;
  }
  
  .item-value {
    color: #303133;
    font-weight: 500;
    font-size: 14px;
    flex: 1;
  }
}

.profile-actions {
  display: flex;
  justify-content: center;
  margin-top: 30px;
}

.logout-btn {
  transition: all 0.3s ease;
  
  &:hover {
    background: rgba(245, 108, 108, 0.1);
    transform: translateY(-2px);
  }
}

.profile-tabs {
  .el-tabs__header {
    margin-bottom: 20px;
  }
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  
  span {
    font-size: 16px;
    font-weight: 600;
    color: #303133;
  }
}

.security-items {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.security-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: rgba(245, 245, 245, 0.3);
  padding: 15px 20px;
  border-radius: 12px;
  transition: all 0.3s ease;
  
  &:hover {
    background: rgba(245, 245, 245, 0.5);
  }
  
  .security-item-content {
    display: flex;
    align-items: center;
    gap: 15px;
  }
  
  .security-item-icon {
    width: 40px;
    height: 40px;
    background: linear-gradient(135deg, #3a7bd5, #00d2ff);
    border-radius: 10px;
    display: flex;
    align-items: center;
    justify-content: center;
    
    i {
      color: white;
      font-size: 20px;
    }
  }
  
  .security-item-text {
    h3 {
      margin: 0 0 4px;
      font-size: 16px;
      color: #303133;
    }
    
    p {
      margin: 0;
      font-size: 13px;
      color: #606266;
    }
  }
}

.notification-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 0;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
  
  &:last-child {
    border-bottom: none;
  }
  
  .notification-text {
    h3 {
      margin: 0 0 4px;
      font-size: 16px;
      color: #303133;
    }
    
    p {
      margin: 0;
      font-size: 13px;
      color: #606266;
    }
  }
}

.mt-20 {
  margin-top: 20px;
}

/* 响应式调整 */
@media (max-width: 768px) {
  .user-profile-container {
    padding: 10px;
  }
  
  .profile-card {
    margin-bottom: 20px;
  }
  
  .profile-name h2 {
    font-size: 20px;
  }
  
  .security-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 15px;
    
    .el-switch {
      margin-left: 55px;
    }
  }
}
</style> 