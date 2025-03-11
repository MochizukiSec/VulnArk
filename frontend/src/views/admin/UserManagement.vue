<template>
  <div class="user-management-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">
          <span class="title-icon"><i class="el-icon-user"></i></span>
          用户管理
          <span class="title-highlight">中心</span>
        </h1>
        <p class="page-subtitle">管理系统用户、权限和组织结构</p>
      </div>
      <div class="header-actions">
        <el-button type="primary" @click="showAddUserDialog" class="add-user-btn">
          <i class="el-icon-plus"></i> 添加用户
        </el-button>
      </div>
    </div>

    <!-- 用户筛选和搜索 -->
    <div class="filter-container">
      <el-card shadow="hover" class="filter-card">
        <template #header>
          <div class="card-header">
            <h3 class="filter-title">
              <i class="el-icon-search"></i> 高级筛选
            </h3>
          </div>
        </template>
        <el-form :inline="true" :model="filterForm" class="filter-form">
          <el-form-item label="状态">
            <el-select v-model="filterForm.status" placeholder="选择状态" clearable class="filter-select">
              <el-option label="已激活" value="active">
                <template #default>
                  <div class="status-option">
                    <span class="status-dot active"></span>
                    <span>已激活</span>
                  </div>
                </template>
              </el-option>
              <el-option label="未激活" value="inactive">
                <template #default>
                  <div class="status-option">
                    <span class="status-dot inactive"></span>
                    <span>未激活</span>
                  </div>
                </template>
              </el-option>
              <el-option label="已禁用" value="disabled">
                <template #default>
                  <div class="status-option">
                    <span class="status-dot disabled"></span>
                    <span>已禁用</span>
                  </div>
                </template>
              </el-option>
            </el-select>
          </el-form-item>
          
          <el-form-item label="角色">
            <el-select v-model="filterForm.role" placeholder="选择角色" clearable class="filter-select">
              <el-option label="管理员" value="admin">
                <template #default>
                  <div class="role-option">
                    <i class="el-icon-s-tools role-icon admin"></i>
                    <span>管理员</span>
                  </div>
                </template>
              </el-option>
              <el-option label="普通用户" value="user">
                <template #default>
                  <div class="role-option">
                    <i class="el-icon-user role-icon user"></i>
                    <span>普通用户</span>
                  </div>
                </template>
              </el-option>
              <el-option label="只读用户" value="readonly">
                <template #default>
                  <div class="role-option">
                    <i class="el-icon-view role-icon readonly"></i>
                    <span>只读用户</span>
                  </div>
                </template>
              </el-option>
            </el-select>
          </el-form-item>
          
          <el-form-item label="搜索">
            <el-input 
              v-model="filterForm.search" 
              placeholder="搜索用户名、邮箱或姓名"
              clearable
              prefix-icon="el-icon-search"
              class="search-input"
            />
          </el-form-item>
          
          <el-form-item class="filter-buttons">
            <el-button type="primary" @click="searchUsers" class="search-btn">
              <i class="el-icon-search"></i> 搜索
            </el-button>
            <el-button @click="resetFilter" class="reset-btn">
              <i class="el-icon-refresh"></i> 重置
            </el-button>
          </el-form-item>
        </el-form>
      </el-card>
    </div>

    <!-- 用户列表 -->
    <div class="user-table-container">
      <el-card shadow="hover" class="user-table-card">
        <template #header>
          <div class="card-header">
            <h3 class="user-table-title">
              <i class="el-icon-user-solid"></i> 用户列表
            </h3>
            <div class="table-actions">
              <el-tooltip content="刷新列表" placement="top">
                <el-button circle size="small" @click="searchUsers">
                  <i class="el-icon-refresh"></i>
                </el-button>
              </el-tooltip>
            </div>
          </div>
        </template>
        <div v-loading="loading">
          <el-table 
            :data="users" 
            style="width: 100%"
            class="user-table"
            :header-cell-style="{ background: '#f5f7fa', color: '#606266' }"
          >
            <el-table-column prop="username" label="用户名" min-width="120">
              <template #default="scope">
                <div class="user-info">
                  <el-avatar :size="32" :src="scope.row.profile_picture" class="user-avatar">
                    {{ scope.row.username.charAt(0).toUpperCase() }}
                  </el-avatar>
                  <span class="username">{{ scope.row.username }}</span>
                </div>
              </template>
            </el-table-column>
            
            <el-table-column prop="name" label="姓名" min-width="120">
              <template #default="scope">
                <span>{{ (scope.row.first_name || '') + ' ' + (scope.row.last_name || '') }}</span>
              </template>
            </el-table-column>
            
            <el-table-column prop="email" label="邮箱" min-width="180">
              <template #default="scope">
                <div class="email-info">
                  <i class="el-icon-message email-icon"></i>
                  <span>{{ scope.row.email }}</span>
                </div>
              </template>
            </el-table-column>
            
            <el-table-column prop="role" label="角色" width="120">
              <template #default="scope">
                <el-tag 
                  :type="getRoleTagType(scope.row.role)" 
                  effect="light"
                  class="role-tag"
                >
                  <i class="role-icon" :class="getRoleIcon(scope.row.role)"></i>
                  {{ getRoleDisplayName(scope.row.role) }}
                </el-tag>
              </template>
            </el-table-column>
            
            <el-table-column prop="status" label="状态" width="100">
              <template #default="scope">
                <el-tag 
                  :type="getStatusTagType(scope.row.status)" 
                  effect="light"
                  class="status-tag"
                >
                  <span class="status-dot" :class="scope.row.status"></span>
                  {{ getStatusDisplayName(scope.row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            
            <el-table-column prop="last_login" label="最后登录" width="180">
              <template #default="scope">
                <div class="login-info">
                  <i class="el-icon-time login-icon"></i>
                  <span>{{ formatDateTime(scope.row.last_login) }}</span>
                </div>
              </template>
            </el-table-column>
            
            <el-table-column label="操作" width="220" fixed="right">
              <template #default="scope">
                <div class="table-actions">
                  <el-button 
                    size="small" 
                    type="primary" 
                    @click="showEditUserDialog(scope.row)"
                    class="action-btn edit-btn"
                  >
                    <i class="el-icon-edit"></i> 编辑
                  </el-button>
                  
                  <el-button 
                    size="small" 
                    :type="scope.row.status === 'active' ? 'warning' : 'success'" 
                    @click="toggleUserStatus(scope.row)"
                    class="action-btn toggle-btn"
                  >
                    <i :class="scope.row.status === 'active' ? 'el-icon-lock' : 'el-icon-unlock'"></i>
                    {{ scope.row.status === 'active' ? '禁用' : '启用' }}
                  </el-button>
                  
                  <el-button 
                    size="small" 
                    type="danger" 
                    @click="confirmDeleteUser(scope.row)"
                    class="action-btn delete-btn"
                  >
                    <i class="el-icon-delete"></i> 删除
                  </el-button>
                </div>
              </template>
            </el-table-column>
          </el-table>
          
          <!-- 分页 -->
          <div class="pagination-container">
            <el-pagination
              v-model:currentPage="pagination.page"
              v-model:page-size="pagination.pageSize"
              :page-sizes="[10, 20, 50, 100]"
              layout="total, sizes, prev, pager, next, jumper"
              :total="pagination.total"
              @size-change="handleSizeChange"
              @current-change="handlePageChange"
              background
              class="pagination"
            />
          </div>
        </div>
      </el-card>
    </div>

    <!-- 添加/编辑用户对话框 -->
    <el-dialog
      :title="isEditMode ? '编辑用户' : '添加用户'"
      v-model="userDialogVisible"
      width="520px"
      class="user-dialog"
      destroy-on-close
    >
      <div class="dialog-header" v-if="isEditMode">
        <div class="user-preview">
          <el-avatar :size="60" :src="userForm.profile_picture" class="preview-avatar">
            {{ userForm.username ? userForm.username.charAt(0).toUpperCase() : 'U' }}
          </el-avatar>
          <div class="preview-info">
            <h3>{{ userForm.username || '新用户' }}</h3>
            <p>{{ userForm.email || '无邮箱信息' }}</p>
          </div>
        </div>
      </div>

      <el-form 
        ref="userFormRef"
        :model="userForm" 
        :rules="userFormRules" 
        label-width="100px"
        class="user-form"
      >
        <el-form-item label="用户名" prop="username">
          <el-input 
            v-model="userForm.username" 
            :disabled="isEditMode"
            prefix-icon="el-icon-user"
            placeholder="请输入用户名"
          />
          <div class="form-tip" v-if="!isEditMode">用户名一旦创建无法修改</div>
        </el-form-item>
        
        <el-form-item label="姓名" prop="name">
          <el-input 
            v-model="userForm.name"
            prefix-icon="el-icon-user"
            placeholder="请输入姓名"
          />
        </el-form-item>
        
        <el-form-item label="邮箱" prop="email">
          <el-input 
            v-model="userForm.email" 
            type="email"
            prefix-icon="el-icon-message"
            placeholder="请输入邮箱地址"
          />
        </el-form-item>
        
        <el-form-item label="角色" prop="role">
          <el-select v-model="userForm.role" style="width: 100%" placeholder="请选择用户角色">
            <el-option label="管理员" value="admin">
              <template #default>
                <div class="role-option">
                  <i class="el-icon-s-tools role-icon admin"></i>
                  <span>管理员</span>
                </div>
              </template>
            </el-option>
            <el-option label="普通用户" value="user">
              <template #default>
                <div class="role-option">
                  <i class="el-icon-user role-icon user"></i>
                  <span>普通用户</span>
                </div>
              </template>
            </el-option>
            <el-option label="只读用户" value="readonly">
              <template #default>
                <div class="role-option">
                  <i class="el-icon-view role-icon readonly"></i>
                  <span>只读用户</span>
                </div>
              </template>
            </el-option>
          </el-select>
          <div class="form-tip">管理员拥有系统所有权限</div>
        </el-form-item>
        
        <el-form-item label="状态" prop="status">
          <el-select v-model="userForm.status" style="width: 100%" placeholder="请选择账户状态">
            <el-option label="已激活" value="active">
              <template #default>
                <div class="status-option">
                  <span class="status-dot active"></span>
                  <span>已激活</span>
                </div>
              </template>
            </el-option>
            <el-option label="未激活" value="inactive">
              <template #default>
                <div class="status-option">
                  <span class="status-dot inactive"></span>
                  <span>未激活</span>
                </div>
              </template>
            </el-option>
            <el-option label="已禁用" value="disabled">
              <template #default>
                <div class="status-option">
                  <span class="status-dot disabled"></span>
                  <span>已禁用</span>
                </div>
              </template>
            </el-option>
          </el-select>
        </el-form-item>
        
        <div v-if="!isEditMode" class="password-section">
          <div class="section-header">
            <i class="el-icon-lock"></i>
            <span>密码设置</span>
          </div>
          
          <el-form-item label="密码" prop="password">
            <el-input 
              v-model="userForm.password" 
              type="password"
              prefix-icon="el-icon-lock"
              placeholder="请输入密码"
              show-password
            />
          </el-form-item>
          
          <el-form-item label="确认密码" prop="confirmPassword">
            <el-input 
              v-model="userForm.confirmPassword" 
              type="password"
              prefix-icon="el-icon-lock"
              placeholder="请再次输入密码"
              show-password
            />
          </el-form-item>
        </div>
      </el-form>
      
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="userDialogVisible = false" class="cancel-btn">取消</el-button>
          <el-button type="primary" @click="submitUserForm" :loading="submitting" class="submit-btn">
            {{ isEditMode ? '保存修改' : '创建用户' }}
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useStore } from 'vuex'

export default {
  name: 'UserManagement',
  
  setup() {
    const store = useStore()
    
    // 加载和提交状态
    const loading = ref(false)
    const submitting = ref(false)
    
    // 筛选表单
    const filterForm = reactive({
      status: '',
      role: '',
      search: ''
    })
    
    // 分页信息
    const pagination = reactive({
      page: 1,
      pageSize: 10,
      total: 0
    })
    
    // 用户表单对话框
    const userDialogVisible = ref(false)
    const userFormRef = ref(null)
    const isEditMode = ref(false)
    
    // 用户表单
    const userForm = reactive({
      id: '',
      username: '',
      name: '',
      email: '',
      role: 'user',
      status: 'active',
      password: '',
      confirmPassword: ''
    })
    
    // 表单验证规则
    const userFormRules = {
      username: [
        { required: true, message: '请输入用户名', trigger: 'blur' },
        { min: 3, max: 20, message: '用户名长度在3到20个字符之间', trigger: 'blur' }
      ],
      name: [
        { required: true, message: '请输入姓名', trigger: 'blur' }
      ],
      email: [
        { required: true, message: '请输入邮箱', trigger: 'blur' },
        { type: 'email', message: '请输入有效的邮箱地址', trigger: 'blur' }
      ],
      password: [
        { required: true, message: '请输入密码', trigger: 'blur' },
        { min: 6, message: '密码至少需要6个字符', trigger: 'blur' }
      ],
      confirmPassword: [
        { required: true, message: '请确认密码', trigger: 'blur' },
        {
          validator: (rule, value, callback) => {
            if (value !== userForm.password) {
              callback(new Error('两次输入的密码不一致'))
            } else {
              callback()
            }
          },
          trigger: 'blur'
        }
      ]
    }
    
    // 用户数据来自Vuex store
    const users = computed(() => store.getters['user/allUsers'])
    
    // 获取角色显示名称
    const getRoleDisplayName = (role) => {
      const roleMap = {
        admin: '管理员',
        user: '普通用户',
        readonly: '只读用户'
      }
      return roleMap[role] || role
    }
    
    // 获取角色标签类型
    const getRoleTagType = (role) => {
      const typeMap = {
        admin: 'danger',
        user: 'primary',
        readonly: 'info'
      }
      return typeMap[role] || ''
    }
    
    // 获取角色图标
    const getRoleIcon = (role) => {
      const iconMap = {
        admin: 'el-icon-s-tools',
        user: 'el-icon-user',
        readonly: 'el-icon-view'
      }
      return iconMap[role] || 'el-icon-user'
    }
    
    // 获取状态显示名称
    const getStatusDisplayName = (status) => {
      const statusMap = {
        active: '已激活',
        inactive: '未激活',
        disabled: '已禁用'
      }
      return statusMap[status] || status
    }
    
    // 获取状态标签类型
    const getStatusTagType = (status) => {
      const typeMap = {
        active: 'success',
        inactive: 'warning',
        disabled: 'info'
      }
      return typeMap[status] || ''
    }
    
    // 格式化日期时间
    const formatDateTime = (dateTimeString) => {
      if (!dateTimeString) return '从未登录'
      
      try {
        const date = new Date(dateTimeString)
        if (isNaN(date.getTime())) return '无效日期'
        
        return date.toLocaleDateString('zh-CN', {
          year: 'numeric',
          month: '2-digit',
          day: '2-digit',
          hour: '2-digit',
          minute: '2-digit'
        })
      } catch (error) {
        console.error('日期格式化错误:', error)
        return '无效日期'
      }
    }
    
    // 搜索用户
    const searchUsers = async () => {
      loading.value = true
      
      try {
        await store.dispatch('user/fetchUsers', {
          page: pagination.page,
          limit: pagination.pageSize,
          search: filterForm.search,
          status: filterForm.status,
          role: filterForm.role
        })
        
        // 更新分页数据
        const storePagination = store.getters['user/pagination']
        pagination.total = storePagination.total
      } catch (error) {
        ElMessage.error(`加载用户数据失败: ${error.message || '未知错误'}`)
      } finally {
        loading.value = false
      }
    }
    
    // 在页面加载时获取用户数据
    onMounted(() => {
      searchUsers()
    })
    
    // 处理分页变化
    const handlePageChange = (newPage) => {
      pagination.page = newPage
      searchUsers()
    }
    
    // 处理每页显示数量变化
    const handleSizeChange = (newSize) => {
      pagination.pageSize = newSize
      pagination.page = 1 // 重置到第一页
      searchUsers()
    }
    
    // 显示添加用户对话框
    const showAddUserDialog = () => {
      isEditMode.value = false
      // 重置表单
      Object.assign(userForm, {
        id: '',
        username: '',
        name: '',
        email: '',
        role: 'user',
        status: 'active',
        password: '',
        confirmPassword: ''
      })
      
      userDialogVisible.value = true
    }
    
    // 显示编辑用户对话框
    const showEditUserDialog = (user) => {
      isEditMode.value = true
      // 填充表单数据
      Object.assign(userForm, {
        id: user.id,
        username: user.username,
        first_name: user.first_name || '',
        last_name: user.last_name || '',
        name: `${user.first_name || ''} ${user.last_name || ''}`.trim(),
        email: user.email,
        role: user.role,
        status: user.status,
        password: '',
        confirmPassword: ''
      })
      
      userDialogVisible.value = true
    }
    
    // 提交用户表单
    const submitUserForm = async () => {
      // eslint-disable-next-line no-unused-vars
      await userFormRef.value.validate(async (valid, fields) => {
        if (valid) {
          submitting.value = true
          
          try {
            // 从名称中提取姓和名
            const nameParts = userForm.name.trim().split(' ');
            const firstName = nameParts[0] || '';
            const lastName = nameParts.length > 1 ? nameParts.slice(1).join(' ') : '';
            
            // 准备提交数据
            const userData = {
              username: userForm.username,
              email: userForm.email,
              first_name: firstName,  // 使用下划线格式与后端匹配
              last_name: lastName || ' ', // 提供默认值，因为后端要求必填
              department: '', // 提供默认值
              role: userForm.role,
              status: userForm.status
            }
            
            // 对于新用户，需要添加密码
            if (!isEditMode.value) {
              userData.password = userForm.password
            } else if (userForm.password) {
              // 编辑时，仅在提供密码的情况下更新密码
              userData.new_password = userForm.password
            }
            
            console.log('提交的用户数据:', userData) // 调试用
            
            if (isEditMode.value) {
              // 更新用户
              await store.dispatch('user/updateUser', {
                id: userForm.id,
                data: userData
              })
              ElMessage.success('用户信息已成功更新')
            } else {
              // 创建新用户
              await store.dispatch('user/createUser', userData)
              ElMessage.success('用户已成功创建')
            }
            
            userDialogVisible.value = false
            // 刷新用户列表
            searchUsers()
          } catch (error) {
            console.error('操作失败:', error)
            ElMessage.error(`操作失败: ${error.response?.data?.error || error.message || '未知错误'}`)
          } finally {
            submitting.value = false
          }
        }
      })
    }
    
    // 删除用户
    const confirmDeleteUser = (user) => {
      ElMessageBox.confirm(
        `确定要删除用户 ${user.name || user.username} 吗？此操作不可恢复。`,
        '确认删除',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }
      ).then(async () => {
        loading.value = true
        
        try {
          await store.dispatch('user/deleteUser', user.id)
          ElMessage.success('用户已成功删除')
          // 刷新用户列表
          searchUsers()
        } catch (error) {
          ElMessage.error(`删除失败: ${error.response?.data?.error || error.message || '未知错误'}`)
        } finally {
          loading.value = false
        }
      }).catch(() => {
        // 用户取消删除
      })
    }
    
    // 切换用户状态
    const toggleUserStatus = async (user) => {
      loading.value = true
      
      try {
        const newStatus = user.status === 'active' ? 'disabled' : 'active'
        
        await store.dispatch('user/updateUser', {
          id: user.id,
          data: { status: newStatus }
        })
        
        ElMessage.success(`用户状态已更新为${getStatusDisplayName(newStatus)}`)
        // 刷新用户列表
        searchUsers()
      } catch (error) {
        console.error('状态更新失败:', error)
        ElMessage.error(`更新状态失败: ${error.response?.data?.error || error.message || '未知错误'}`)
      } finally {
        loading.value = false
      }
    }
    
    // 重置筛选条件
    const resetFilter = () => {
      Object.assign(filterForm, {
        status: '',
        role: '',
        search: ''
      })
      searchUsers()
    }
    
    return {
      users,
      loading,
      submitting,
      filterForm,
      pagination,
      userDialogVisible,
      userFormRef,
      userForm,
      userFormRules,
      isEditMode,
      getRoleDisplayName,
      getRoleTagType,
      getRoleIcon,
      getStatusDisplayName,
      getStatusTagType,
      formatDateTime,
      searchUsers,
      resetFilter,
      handlePageChange,
      handleSizeChange,
      showAddUserDialog,
      showEditUserDialog,
      submitUserForm,
      confirmDeleteUser,
      toggleUserStatus
    }
  }
}
</script>

<style lang="scss" scoped>
.user-management-page {
  padding: 24px;
  animation: fadeIn 0.6s ease;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  flex-wrap: wrap;
  gap: 16px;
  
  .header-content {
    flex: 1;
  }
  
  .page-title {
    font-size: 28px;
    font-weight: 600;
    margin: 0 0 8px 0;
    display: flex;
    align-items: center;
    color: #303133;
    
    .title-icon {
      margin-right: 12px;
      font-size: 28px;
      color: #409EFF;
    }
    
    .title-highlight {
      background: linear-gradient(120deg, #409EFF, #53A8FF);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      margin-left: 8px;
      font-weight: 700;
    }
  }
  
  .page-subtitle {
    color: #606266;
    font-size: 16px;
    margin: 0;
  }
  
  .header-actions {
    .add-user-btn {
      border-radius: 8px;
      padding: 10px 20px;
      font-weight: 500;
      transition: all 0.3s;
      
      i {
        margin-right: 6px;
      }
      
      &:hover {
        transform: translateY(-2px);
        box-shadow: 0 4px 12px rgba(64, 158, 255, 0.2);
      }
    }
  }
}

.filter-container,
.user-table-container {
  margin-bottom: 24px;
}

.filter-card,
.user-table-card {
  border-radius: 12px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.05);
  transition: all 0.3s;
  position: relative;
  overflow: hidden;
  
  &:hover {
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.09);
  }
  
  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 4px;
    background: linear-gradient(90deg, #409EFF, #53A8FF);
  }
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 4px 0;
}

.filter-title {
  font-size: 16px;
  font-weight: 600;
  margin: 0;
  color: #303133;
  display: flex;
  align-items: center;
  
  i {
    margin-right: 8px;
    font-size: 16px;
    color: #409EFF;
  }
}

.filter-form {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  padding: 8px 0;
  
  .filter-select,
  .search-input {
    min-width: 200px;
    
    :deep(.el-input__inner) {
      border-radius: 8px;
      transition: all 0.3s;
      
      &:focus {
        box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.2);
      }
    }
  }
  
  .filter-buttons {
    margin-left: auto;
    display: flex;
    gap: 12px;
    
    .search-btn,
    .reset-btn {
      border-radius: 8px;
      padding: 10px 16px;
      transition: all 0.3s;
      
      i {
        margin-right: 6px;
      }
      
      &:hover {
        transform: translateY(-2px);
      }
    }
    
    .search-btn:hover {
      box-shadow: 0 4px 12px rgba(64, 158, 255, 0.2);
    }
  }
}

.status-option,
.role-option {
  display: flex;
  align-items: center;
  
  .status-dot {
    width: 10px;
    height: 10px;
    border-radius: 50%;
    margin-right: 8px;
    
    &.active {
      background-color: #67C23A;
    }
    
    &.inactive {
      background-color: #E6A23C;
    }
    
    &.disabled {
      background-color: #909399;
    }
  }
  
  .role-icon {
    margin-right: 8px;
    font-size: 16px;
    
    &.admin {
      color: #F56C6C;
    }
    
    &.user {
      color: #409EFF;
    }
    
    &.readonly {
      color: #909399;
    }
  }
}

.user-table-title {
  font-size: 16px;
  font-weight: 600;
  margin: 0;
  color: #303133;
  display: flex;
  align-items: center;
  
  i {
    margin-right: 8px;
    font-size: 16px;
    color: #409EFF;
  }
}

.user-table {
  margin-top: 8px;
  border-radius: 8px;
  overflow: hidden;
  
  :deep(.el-table__header) {
    th {
      font-weight: 600;
      padding: 12px 0;
    }
  }
  
  :deep(.el-table__row) {
    transition: all 0.2s;
    
    &:hover {
      background-color: #f5f7fc !important;
    }
  }
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  
  .user-avatar {
    background: linear-gradient(135deg, #a8edea 0%, #fed6e3 100%);
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: bold;
    color: #606266;
  }
  
  .username {
    font-weight: 500;
    color: #303133;
  }
}

.email-info,
.login-info {
  display: flex;
  align-items: center;
  
  .email-icon,
  .login-icon {
    margin-right: 6px;
    color: #909399;
    font-size: 14px;
  }
}

.role-tag,
.status-tag {
  display: flex;
  align-items: center;
  padding: 4px 8px;
  border-radius: 4px;
  
  .role-icon {
    margin-right: 4px;
    font-size: 14px;
  }
  
  .status-dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    margin-right: 6px;
    
    &.active {
      background-color: #67C23A;
    }
    
    &.inactive {
      background-color: #E6A23C;
    }
    
    &.disabled {
      background-color: #909399;
    }
  }
}

.table-actions {
  display: flex;
  gap: 8px;
  
  .action-btn {
    padding: 4px 8px;
    border-radius: 4px;
    transition: all 0.3s;
    
    i {
      margin-right: 4px;
    }
    
    &:hover {
      opacity: 0.9;
      transform: translateY(-1px);
    }
    
    &.edit-btn {
      background-color: #409EFF;
      border-color: #409EFF;
    }
    
    &.toggle-btn {
      &.el-button--warning {
        background-color: #E6A23C;
        border-color: #E6A23C;
      }
      
      &.el-button--success {
        background-color: #67C23A;
        border-color: #67C23A;
      }
    }
    
    &.delete-btn {
      background-color: #F56C6C;
      border-color: #F56C6C;
    }
  }
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
  
  .pagination {
    :deep(.el-pagination__total) {
      font-weight: 500;
    }
    
    :deep(.el-pagination__sizes) {
      margin-right: 16px;
    }
    
    :deep(.btn-prev, .btn-next, .el-pager li) {
      transition: all 0.3s;
      
      &:hover {
        transform: translateY(-2px);
      }
    }
    
    :deep(.el-pager li.active) {
      font-weight: 700;
    }
  }
}

@media (max-width: 768px) {
  .user-management-page {
    padding: 16px;
  }
  
  .page-header {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .filter-form {
    flex-direction: column;
    gap: 12px;
    
    .filter-buttons {
      margin-left: 0;
      width: 100%;
      
      .search-btn,
      .reset-btn {
        flex: 1;
      }
    }
  }
  
  .el-form-item {
    margin-bottom: 10px;
  }
  
  .table-actions {
    flex-direction: column;
    gap: 4px;
    
    .action-btn {
      width: 100%;
    }
  }
}

/* 用户对话框样式 */
:deep(.user-dialog) {
  .el-dialog__header {
    padding: 20px;
    margin-right: 0;
    border-bottom: 1px solid #EBEEF5;
    
    .el-dialog__title {
      font-size: 18px;
      font-weight: 600;
      color: #303133;
    }
  }
  
  .el-dialog__body {
    padding: 24px;
  }
  
  .el-dialog__footer {
    padding: 16px 24px;
    border-top: 1px solid #EBEEF5;
  }
}

.dialog-header {
  margin-bottom: 24px;
  
  .user-preview {
    display: flex;
    align-items: center;
    
    .preview-avatar {
      background: linear-gradient(135deg, #a8edea 0%, #fed6e3 100%);
      display: flex;
      align-items: center;
      justify-content: center;
      font-weight: bold;
      color: #606266;
      margin-right: 16px;
      font-size: 20px;
    }
    
    .preview-info {
      h3 {
        margin: 0 0 4px 0;
        font-size: 16px;
        color: #303133;
      }
      
      p {
        margin: 0;
        color: #909399;
        font-size: 14px;
      }
    }
  }
}

.user-form {
  .form-tip {
    color: #909399;
    font-size: 12px;
    margin-top: 4px;
  }
  
  :deep(.el-input__inner) {
    border-radius: 8px;
    
    &:focus {
      box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.1);
    }
  }
  
  :deep(.el-select .el-input__inner) {
    border-radius: 8px;
  }
}

.password-section {
  margin-top: 16px;
  border-top: 1px dashed #EBEEF5;
  padding-top: 16px;
  
  .section-header {
    display: flex;
    align-items: center;
    margin-bottom: 16px;
    color: #606266;
    font-weight: 600;
    
    i {
      color: #F56C6C;
      margin-right: 6px;
    }
  }
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  
  .cancel-btn,
  .submit-btn {
    min-width: 100px;
    border-radius: 8px;
    transition: all 0.3s;
    
    &:hover {
      transform: translateY(-2px);
    }
  }
  
  .submit-btn:hover {
    box-shadow: 0 4px 12px rgba(64, 158, 255, 0.2);
  }
}
</style> 