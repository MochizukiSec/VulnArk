<template>
  <div class="profile-container">
    <el-card class="profile-card">
      <template #header>
        <div class="card-header">
          <span>个人资料</span>
          <el-button type="primary" size="small" @click="handleEdit">编辑</el-button>
        </div>
      </template>
      <div class="profile-content">
        <div class="avatar-container">
          <el-avatar :size="100" :src="userInfo.avatar || ''" />
          <div class="user-role">{{ getRoleName(userInfo.role) }}</div>
        </div>
        <div class="info-container">
          <el-descriptions :column="1" border>
            <el-descriptions-item label="用户名">{{ userInfo.username }}</el-descriptions-item>
            <el-descriptions-item label="姓名">{{ userInfo.real_name || '-' }}</el-descriptions-item>
            <el-descriptions-item label="邮箱">{{ userInfo.email }}</el-descriptions-item>
            <el-descriptions-item label="电话">{{ userInfo.phone || '-' }}</el-descriptions-item>
            <el-descriptions-item label="最后登录">{{ userInfo.last_login || '-' }}</el-descriptions-item>
          </el-descriptions>
        </div>
      </div>
    </el-card>

    <!-- 编辑个人资料对话框 -->
    <el-dialog title="编辑个人资料" v-model="dialogVisible" width="500px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="form.username" disabled></el-input>
        </el-form-item>
        <el-form-item label="姓名" prop="real_name">
          <el-input v-model="form.real_name"></el-input>
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="form.email"></el-input>
        </el-form-item>
        <el-form-item label="电话" prop="phone">
          <el-input v-model="form.phone"></el-input>
        </el-form-item>
        <el-form-item label="头像">
          <el-upload
            action="#"
            list-type="picture-card"
            :auto-upload="false"
            :limit="1"
            :on-change="handleAvatarChange"
          >
            <el-icon><Plus /></el-icon>
          </el-upload>
        </el-form-item>
        
        <!-- 修改密码部分 -->
        <el-divider>修改密码</el-divider>
        <el-form-item label="旧密码" prop="old_password">
          <el-input v-model="form.old_password" type="password" show-password></el-input>
        </el-form-item>
        <el-form-item label="新密码" prop="new_password">
          <el-input v-model="form.new_password" type="password" show-password></el-input>
        </el-form-item>
        <el-form-item label="确认密码" prop="confirm_password">
          <el-input v-model="form.confirm_password" type="password" show-password></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSave" :loading="loading">保存</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { reactive, ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getUserInfo, updateProfile, updatePassword } from '@/api/user'
import { Plus } from '@element-plus/icons-vue'

export default {
  name: 'UserProfile',
  components: {
    Plus
  },
  setup() {
    const userInfo = reactive({
      id: '',
      username: '',
      real_name: '',
      email: '',
      phone: '',
      role: '',
      avatar: '',
      last_login: ''
    })
    
    const dialogVisible = ref(false)
    const loading = ref(false)
    const formRef = ref(null)
    
    // 表单数据
    const form = reactive({
      username: '',
      real_name: '',
      email: '',
      phone: '',
      avatar: '',
      old_password: '',
      new_password: '',
      confirm_password: ''
    })
    
    // 表单校验规则
    const rules = reactive({
      email: [
        { type: 'email', message: '请输入有效的电子邮件地址', trigger: 'blur' }
      ],
      phone: [
        { pattern: /^1[3-9]\d{9}$/, message: '请输入有效的手机号码', trigger: 'blur' }
      ],
      new_password: [
        { min: 6, message: '密码长度不能少于6个字符', trigger: 'blur' }
      ],
      confirm_password: [
        { validator: validateConfirmPassword, trigger: 'blur' }
      ]
    })
    
    // 确认密码校验
    function validateConfirmPassword(rule, value, callback) {
      if (value === '') {
        callback()
      } else if (value !== form.new_password) {
        callback(new Error('两次输入密码不一致'))
      } else {
        callback()
      }
    }
    
    // 获取用户角色名称
    const getRoleName = (role) => {
      switch (role) {
        case 'admin':
          return '系统管理员'
        case 'security_engineer':
          return '安全工程师'
        case 'asset_manager':
          return '资产管理员'
        case 'reporter':
          return '报告查看者'
        case 'user':
          return '普通用户'
        default:
          return '未知角色'
      }
    }
    
    // 获取用户信息
    const fetchUserInfo = async () => {
      try {
        const response = await getUserInfo()
        if (response.code === 200) {
          Object.assign(userInfo, response.data)
        } else {
          ElMessage.error('获取用户信息失败: ' + (response.message || '未知错误'))
        }
      } catch (error) {
        console.error('获取用户信息失败', error)
        ElMessage.error('获取用户信息失败: ' + (error.message || '未知错误'))
      }
    }
    
    // 处理编辑按钮点击
    const handleEdit = () => {
      Object.assign(form, {
        username: userInfo.username,
        real_name: userInfo.real_name || '',
        email: userInfo.email || '',
        phone: userInfo.phone || '',
        avatar: userInfo.avatar || '',
        old_password: '',
        new_password: '',
        confirm_password: ''
      })
      dialogVisible.value = true
    }
    
    // 处理头像变更
    const handleAvatarChange = (file) => {
      // 这里可以处理头像上传逻辑
      const isImage = file.raw.type.startsWith('image/')
      const isLt2M = file.raw.size / 1024 / 1024 < 2
      
      if (!isImage) {
        ElMessage.error('上传头像图片只能是图片格式!')
        return false
      }
      
      if (!isLt2M) {
        ElMessage.error('上传头像图片大小不能超过 2MB!')
        return false
      }
      
      // 可以在这里将图片转为base64存储或上传到服务器
      const reader = new FileReader()
      reader.readAsDataURL(file.raw)
      reader.onload = () => {
        form.avatar = reader.result
      }
    }
    
    // 处理保存按钮点击
    const handleSave = async () => {
      if (!formRef.value) return
      
      try {
        await formRef.value.validate()
        
        loading.value = true
        // 分开处理个人资料更新和密码更新
        try {
          // 更新个人资料
          const profileData = {
            real_name: form.real_name,
            email: form.email,
            phone: form.phone,
            avatar: form.avatar
          }
          
          const profileResponse = await updateProfile(profileData)
          
          if (profileResponse.code !== 200) {
            ElMessage.error('更新个人资料失败: ' + (profileResponse.message || '未知错误'))
            return
          }
          
          // 如果输入了旧密码，则尝试更新密码
          if (form.old_password) {
            if (!form.new_password) {
              ElMessage.warning('请输入新密码')
              return
            }
            
            const passwordData = {
              old_password: form.old_password,
              new_password: form.new_password
            }
            
            const passwordResponse = await updatePassword(passwordData)
            
            if (passwordResponse.code !== 200) {
              ElMessage.error('密码更新失败: ' + (passwordResponse.message || '未知错误'))
              return
            }
            
            ElMessage.success('个人资料和密码已更新')
          } else {
            ElMessage.success('个人资料已更新')
          }
          
          // 重新获取用户信息
          await fetchUserInfo()
          dialogVisible.value = false
        } catch (error) {
          console.error('保存失败', error)
          ElMessage.error('保存失败: ' + (error.message || '未知错误'))
        }
      } catch (error) {
        console.error('表单验证失败', error)
      } finally {
        loading.value = false
      }
    }
    
    // 组件挂载时获取用户信息
    onMounted(() => {
      fetchUserInfo()
    })
    
    return {
      userInfo,
      dialogVisible,
      form,
      rules,
      loading,
      formRef,
      getRoleName,
      handleEdit,
      handleSave,
      handleAvatarChange
    }
  }
}
</script>

<style scoped>
.profile-container {
  padding: 20px;
  animation: fadeIn 0.6s ease-out;
}

.profile-card {
  max-width: 800px;
  margin: 0 auto;
  border-radius: 12px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.08);
  overflow: hidden;
  transition: all 0.3s;
}

.profile-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 12px 30px rgba(0, 0, 0, 0.12);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 20px;
  background: linear-gradient(to right, #f8f9fa, #eef2f5);
  border-bottom: 1px solid #eaedf1;
}

.card-header span {
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

.profile-content {
  display: flex;
  flex-wrap: wrap;
  gap: 30px;
  padding: 30px;
}

.avatar-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 15px;
}

.el-avatar {
  border: 4px solid rgba(245, 247, 250, 0.8);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  cursor: pointer;
  transition: all 0.3s ease;
}

.el-avatar:hover {
  transform: scale(1.05);
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.15);
}

.user-role {
  background: linear-gradient(to right, #4e54c8, #8f94fb);
  color: white;
  padding: 6px 12px;
  border-radius: 20px;
  font-size: 14px;
  font-weight: 500;
  letter-spacing: 0.5px;
  box-shadow: 0 4px 10px rgba(78, 84, 200, 0.3);
  animation: pulse 2s infinite;
}

.info-container {
  flex: 1;
  min-width: 300px;
}

/* 描述列表样式 */
:deep(.el-descriptions__label) {
  color: #606266;
}

:deep(.el-descriptions__content) {
  color: #303133;
  font-weight: 500;
}

:deep(.el-descriptions-item__cell) {
  padding: 16px 24px;
}

/* 对话框样式 */
.el-dialog {
  border-radius: 12px;
  overflow: hidden;
}

.el-dialog__header {
  padding: 20px;
  background: linear-gradient(to right, #f8f9fa, #eef2f5);
  border-bottom: 1px solid #eaedf1;
}

.el-dialog__title {
  font-weight: 600;
  color: #303133;
}

.el-dialog__body {
  padding: 30px;
}

.el-dialog__footer {
  padding: 15px 20px 20px;
  border-top: 1px solid #eaedf1;
}

/* 表单样式 */
.el-form-item {
  margin-bottom: 22px;
}

:deep(.el-input__inner) {
  border-radius: 8px;
}

:deep(.el-form-item__label) {
  font-weight: 500;
}

/* 按钮样式 */
.el-button {
  border-radius: 8px;
  padding: 10px 20px;
  font-weight: 500;
  transition: all 0.3s;
}

.el-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.el-button--primary {
  background: linear-gradient(to right, #4e54c8, #8f94fb);
  border: none;
}

.el-button--primary:hover {
  background: linear-gradient(to right, #4e54c8, #8f94fb);
  box-shadow: 0 4px 12px rgba(78, 84, 200, 0.25);
}

/* 动画 */
@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes pulse {
  0% {
    box-shadow: 0 4px 10px rgba(78, 84, 200, 0.3);
  }
  50% {
    box-shadow: 0 4px 15px rgba(78, 84, 200, 0.5);
  }
  100% {
    box-shadow: 0 4px 10px rgba(78, 84, 200, 0.3);
  }
}

/* 响应式调整 */
@media (max-width: 768px) {
  .profile-content {
    padding: 20px;
    gap: 20px;
  }
  
  .card-header {
    padding: 12px 15px;
  }
  
  .card-header span {
    font-size: 16px;
  }
  
  :deep(.el-descriptions-item__cell) {
    padding: 12px 15px;
  }
}
</style> 