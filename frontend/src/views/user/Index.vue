<template>
  <div class="user-index-container">
    <div class="page-header">
      <div class="left">
        <h2>用户管理</h2>
      </div>
      <div class="right">
        <el-button type="primary" @click="handleAddUser">添加用户</el-button>
      </div>
    </div>

    <!-- 用户列表 -->
    <el-card>
      <el-table :data="userList" v-loading="loading" border style="width: 100%">
        <el-table-column prop="id" label="ID" width="80"></el-table-column>
        <el-table-column prop="username" label="用户名" width="120"></el-table-column>
        <el-table-column prop="email" label="邮箱" min-width="180"></el-table-column>
        <el-table-column prop="real_name" label="姓名" width="120"></el-table-column>
        <el-table-column prop="phone" label="电话" width="120"></el-table-column>
        <el-table-column prop="role" label="角色" width="100">
          <template #default="scope">
            <el-tag :type="getRoleType(scope.row.role)">{{ getRoleLabel(scope.row.role) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="active" label="状态" width="80">
          <template #default="scope">
            <el-tag :type="scope.row.active ? 'success' : 'danger'">
              {{ scope.row.active ? '激活' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="last_login" label="最后登录" width="180">
          <template #default="scope">
            {{ formatDateTime(scope.row.last_login) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="scope">
            <el-button type="primary" link @click="handleChangeRole(scope.row)">修改角色</el-button>
            <el-button type="danger" link @click="handleDelete(scope.row)" :disabled="scope.row.id === currentUserId">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 添加用户对话框 -->
    <el-dialog title="添加用户" v-model="addDialogVisible" width="500px">
      <el-form :model="userForm" :rules="userRules" ref="userFormRef" label-width="100px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="userForm.username" placeholder="请输入用户名"></el-input>
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="userForm.password" type="password" placeholder="请输入密码"></el-input>
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="userForm.email" placeholder="请输入邮箱"></el-input>
        </el-form-item>
        <el-form-item label="姓名" prop="real_name">
          <el-input v-model="userForm.real_name" placeholder="请输入姓名"></el-input>
        </el-form-item>
        <el-form-item label="电话" prop="phone">
          <el-input v-model="userForm.phone" placeholder="请输入电话"></el-input>
        </el-form-item>
        <el-form-item label="角色" prop="role">
          <el-select v-model="userForm.role" placeholder="请选择角色">
            <el-option v-for="role in roleOptions" :key="role.value" :label="role.label" :value="role.value"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="状态" prop="active">
          <el-switch v-model="userForm.active" active-text="激活" inactive-text="禁用"></el-switch>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="addDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitUserForm" :loading="submitting">确定</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 修改角色对话框 -->
    <el-dialog title="修改用户角色" v-model="roleDialogVisible" width="400px">
      <el-form :model="roleForm" ref="roleFormRef" label-width="80px">
        <el-form-item label="用户名">
          <el-input v-model="roleForm.username" disabled></el-input>
        </el-form-item>
        <el-form-item label="角色" prop="role">
          <el-select v-model="roleForm.role" placeholder="请选择角色">
            <el-option v-for="role in roleOptions" :key="role.value" :label="role.label" :value="role.value"></el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="roleDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitRoleChange" :loading="submitting">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { reactive, ref, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getUserList, createUser, deleteUser, changeUserRole } from '@/api/user'
import { useStore } from 'vuex'

export default {
  name: 'UserIndex',
  setup() {
    const store = useStore()
    const loading = ref(false)
    const submitting = ref(false)
    const addDialogVisible = ref(false)
    const roleDialogVisible = ref(false)
    const userFormRef = ref(null)
    const roleFormRef = ref(null)

    // 当前登录用户的ID
    const currentUserId = computed(() => store.getters.userInfo ? store.getters.userInfo.id : 0)

    // 用户列表数据
    const userList = ref([])

    // 用户表单
    const userForm = reactive({
      username: '',
      password: '',
      email: '',
      real_name: '',
      phone: '',
      role: 'viewer',
      active: true
    })

    // 角色表单
    const roleForm = reactive({
      id: null,
      username: '',
      role: ''
    })

    // 表单验证规则
    const userRules = {
      username: [
        { required: true, message: '请输入用户名', trigger: 'blur' },
        { min: 4, max: 20, message: '长度在 4 到 20 个字符', trigger: 'blur' }
      ],
      password: [
        { required: true, message: '请输入密码', trigger: 'blur' },
        { min: 6, message: '密码长度不能小于 6 个字符', trigger: 'blur' }
      ],
      email: [
        { required: true, message: '请输入邮箱地址', trigger: 'blur' },
        { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
      ],
      role: [
        { required: true, message: '请选择角色', trigger: 'change' }
      ]
    }

    // 角色选项
    const roleOptions = [
      { value: 'admin', label: '管理员' },
      { value: 'manager', label: '经理' },
      { value: 'auditor', label: '审计员' },
      { value: 'operator', label: '操作员' },
      { value: 'viewer', label: '浏览者' }
    ]

    // 获取用户列表
    const fetchUserList = async () => {
      loading.value = true
      try {
        const response = await getUserList()
        if (response.code === 200) {
          userList.value = response.data
        } else {
          ElMessage.error(response.message || '获取用户列表失败')
        }
      } catch (error) {
        console.error('获取用户列表失败:', error)
        ElMessage.error('获取用户列表失败')
      } finally {
        loading.value = false
      }
    }

    // 获取角色标签
    const getRoleLabel = (role) => {
      const option = roleOptions.find(item => item.value === role)
      return option ? option.label : role
    }

    // 获取角色标签类型
    const getRoleType = (role) => {
      switch (role) {
        case 'admin': return 'danger'
        case 'manager': return 'warning'
        case 'auditor': return 'info'
        case 'operator': return 'primary'
        case 'viewer': return 'success'
        default: return 'info'
      }
    }

    // 格式化日期时间
    const formatDateTime = (dateStr) => {
      if (!dateStr) return '-'
      const date = new Date(dateStr)
      return date.toLocaleString()
    }

    // 添加用户
    const handleAddUser = () => {
      // 重置表单
      userForm.username = ''
      userForm.password = ''
      userForm.email = ''
      userForm.real_name = ''
      userForm.phone = ''
      userForm.role = 'viewer'
      userForm.active = true
      
      addDialogVisible.value = true
      // 等待DOM更新后设置焦点
      setTimeout(() => {
        if (userFormRef.value) {
          userFormRef.value.resetFields()
        }
      }, 0)
    }

    // 提交用户表单
    const submitUserForm = async () => {
      if (!userFormRef.value) return
      
      await userFormRef.value.validate(async (valid, fields) => {
        if (!valid) {
          console.log('表单验证失败:', fields)
          return
        }
        
        submitting.value = true
        try {
          const response = await createUser(userForm)
          if (response.code === 200) {
            ElMessage.success('用户创建成功')
            addDialogVisible.value = false
            fetchUserList() // 刷新列表
          } else {
            ElMessage.error(response.message || '创建用户失败')
          }
        } catch (error) {
          console.error('创建用户失败:', error)
          ElMessage.error('创建用户失败')
        } finally {
          submitting.value = false
        }
      })
    }

    // 删除用户
    const handleDelete = (row) => {
      // 不允许删除自己
      if (row.id === currentUserId.value) {
        ElMessage.warning('不能删除自己的账号')
        return
      }
      
      ElMessageBox.confirm(`确认删除用户"${row.username}"?`, '警告', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async () => {
        try {
          const response = await deleteUser(row.id)
          if (response.code === 200) {
            ElMessage.success(`成功删除用户"${row.username}"`)
            fetchUserList() // 刷新列表
          } else {
            ElMessage.error(response.message || '删除用户失败')
          }
        } catch (error) {
          console.error('删除用户失败:', error)
          ElMessage.error('删除用户失败')
        }
      }).catch(() => {})
    }

    // 修改用户角色
    const handleChangeRole = (row) => {
      roleForm.id = row.id
      roleForm.username = row.username
      roleForm.role = row.role
      
      roleDialogVisible.value = true
    }

    // 提交角色修改
    const submitRoleChange = async () => {
      submitting.value = true
      try {
        const response = await changeUserRole(roleForm.id, roleForm.role)
        if (response.code === 200) {
          ElMessage.success('角色修改成功')
          roleDialogVisible.value = false
          fetchUserList() // 刷新列表
        } else {
          ElMessage.error(response.message || '修改角色失败')
        }
      } catch (error) {
        console.error('修改角色失败:', error)
        ElMessage.error('修改角色失败')
      } finally {
        submitting.value = false
      }
    }

    onMounted(() => {
      fetchUserList()
    })

    return {
      loading,
      submitting,
      userList,
      userForm,
      roleForm,
      userFormRef,
      roleFormRef,
      addDialogVisible,
      roleDialogVisible,
      roleOptions,
      currentUserId,
      getRoleLabel,
      getRoleType,
      formatDateTime,
      handleAddUser,
      submitUserForm,
      handleDelete,
      handleChangeRole,
      submitRoleChange,
      userRules
    }
  }
}
</script>

<style scoped>
.user-index-container {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  animation: slideDown 0.5s ease-out;
}

.page-header h2 {
  font-size: 24px;
  font-weight: 600;
  color: #303133;
  margin: 0;
  position: relative;
  padding-left: 12px;
}

.page-header h2::before {
  content: '';
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%);
  height: 20px;
  width: 4px;
  background: linear-gradient(to bottom, #4e54c8, #8f94fb);
  border-radius: 2px;
}

.el-card {
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
  margin-bottom: 20px;
  overflow: hidden;
  animation: fadeIn 0.5s ease-out;
}

.el-table {
  border-radius: 4px;
}

.el-tag {
  text-transform: capitalize;
  border-radius: 4px;
  padding: 0 8px;
  height: 24px;
  line-height: 22px;
  transition: all 0.3s;
}

.el-tag:hover {
  transform: translateY(-2px);
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.1);
}

.el-dialog__body {
  padding: 20px 30px 30px;
}

.el-dialog__header {
  background: linear-gradient(90deg, #f5f7fa, #f0f2f5);
  padding: 15px 20px;
  border-bottom: 1px solid #eaeaea;
}

.el-dialog__title {
  font-weight: 600;
  color: #303133;
}

/* 添加用户对话框中的角色选择器 */
.el-select {
  width: 100%;
}

/* 按钮动画效果 */
.el-button [class*="el-icon"] + span {
  margin-left: 6px;
}

.el-button--primary {
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.1);
}

.el-button--danger {
  box-shadow: 0 4px 12px rgba(245, 108, 108, 0.1);
}

/* 分页居中 */
.el-pagination {
  justify-content: center;
  margin-top: 20px;
}

/* 动画 */
@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 表格行悬停效果 */
:deep(.el-table__row) {
  transition: all 0.3s;
}

:deep(.el-table__row:hover) {
  background-color: #f5f7fa;
  transform: translateY(-2px);
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);
}

/* 响应式调整 */
@media (max-width: 768px) {
  .user-index-container {
    padding: 15px;
  }
  
  .page-header h2 {
    font-size: 20px;
  }
}
</style> 