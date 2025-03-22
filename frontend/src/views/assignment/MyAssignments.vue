<template>
  <div class="app-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>我的漏洞任务</span>
          <div>
            <el-radio-group v-model="filterStatus" size="small" @change="handleStatusChange">
              <el-radio-button label="">全部</el-radio-button>
              <el-radio-button label="pending">待处理</el-radio-button>
              <el-radio-button label="accepted">已接受</el-radio-button>
              <el-radio-button label="fixed">已修复</el-radio-button>
              <el-radio-button label="rejected">已拒绝</el-radio-button>
              <el-radio-button label="closed">已关闭</el-radio-button>
            </el-radio-group>
          </div>
        </div>
      </template>

      <!-- 任务列表 -->
      <el-table 
        :data="assignments" 
        style="width: 100%"
        v-loading="loading"
        border
      >
        <el-table-column label="漏洞名称" min-width="200">
          <template #default="scope">
            <router-link 
              :to="'/vulnerabilities/' + scope.row.vulnerability?.id" 
              class="link-type"
            >
              {{ scope.row.vulnerability?.title || '未知漏洞' }}
            </router-link>
          </template>
        </el-table-column>
        
        <el-table-column label="状态" width="100">
          <template #default="scope">
            <el-tag
              :type="getStatusType(scope.row.status)"
              effect="dark"
            >
              {{ getStatusText(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column label="优先级" width="100">
          <template #default="scope">
            <el-tag :type="getPriorityType(scope.row.priority)" size="small">
              {{ getPriorityText(scope.row.priority) }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column label="分派人" width="150">
          <template #default="scope">
            {{ scope.row.assigned_by?.real_name || scope.row.assigned_by?.username || '未知' }}
          </template>
        </el-table-column>
        
        <el-table-column label="截止日期" width="120">
          <template #default="scope">
            <span :class="{'text-danger': isOverdue(scope.row.due_date)}">
              {{ formatDate(scope.row.due_date) }}
              <el-tag v-if="isOverdue(scope.row.due_date)" type="danger" size="small">已逾期</el-tag>
            </span>
          </template>
        </el-table-column>
        
        <el-table-column label="分派时间" width="180">
          <template #default="scope">
            {{ formatDateTime(scope.row.created_at) }}
          </template>
        </el-table-column>
        
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="scope">
            <el-button
              size="small"
              type="primary"
              @click="handleViewDetail(scope.row)"
              plain
            >
              查看详情
            </el-button>
            
            <el-button
              v-if="scope.row.status === 'pending'"
              size="small"
              type="success"
              @click="handleUpdateStatus(scope.row, 'accepted')"
              plain
            >
              接受
            </el-button>
            
            <el-button
              v-if="scope.row.status === 'pending'"
              size="small"
              type="danger"
              @click="handleUpdateStatus(scope.row, 'rejected')"
              plain
            >
              拒绝
            </el-button>
            
            <el-button
              v-if="scope.row.status === 'accepted'"
              size="small"
              type="primary"
              @click="handleUpdateStatus(scope.row, 'fixed')"
              plain
            >
              标记已修复
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          background
          layout="prev, pager, next, jumper"
          :page-size="pageSize"
          :total="total"
          :current-page="currentPage"
          @current-change="handlePageChange"
        />
      </div>
    </el-card>
    
    <!-- 任务详情对话框 -->
    <el-dialog
      title="任务详情"
      v-model="detailDialogVisible"
      width="700px"
    >
      <template v-if="currentAssignment">
        <div class="detail-header">
          <h3>{{ currentAssignment.vulnerability?.title || '未知漏洞' }}</h3>
          <el-tag
            :type="getStatusType(currentAssignment.status)"
            effect="dark"
          >
            {{ getStatusText(currentAssignment.status) }}
          </el-tag>
        </div>
        
        <el-descriptions :column="2" border>
          <el-descriptions-item label="漏洞ID">
            {{ currentAssignment.vulnerability?.id || '未知' }}
          </el-descriptions-item>
          <el-descriptions-item label="优先级">
            <el-tag :type="getPriorityType(currentAssignment.priority)">
              {{ getPriorityText(currentAssignment.priority) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="分派人">
            {{ currentAssignment.assigned_by?.real_name || currentAssignment.assigned_by?.username || '未知' }}
          </el-descriptions-item>
          <el-descriptions-item label="分派时间">
            {{ formatDateTime(currentAssignment.created_at) }}
          </el-descriptions-item>
          <el-descriptions-item label="截止日期">
            <span :class="{'text-danger': isOverdue(currentAssignment.due_date)}">
              {{ formatDate(currentAssignment.due_date) }}
              <el-tag v-if="isOverdue(currentAssignment.due_date)" type="danger" size="small">已逾期</el-tag>
            </span>
          </el-descriptions-item>
          <el-descriptions-item label="更新时间">
            {{ formatDateTime(currentAssignment.updated_at) }}
          </el-descriptions-item>
        </el-descriptions>
        
        <div class="detail-content">
          <div v-if="currentAssignment.notes" class="detail-section">
            <h4>分派说明</h4>
            <p>{{ currentAssignment.notes }}</p>
          </div>
          
          <div v-if="currentAssignment.response" class="detail-section">
            <h4>处理反馈</h4>
            <p>{{ currentAssignment.response }}</p>
          </div>
          
          <div class="detail-section">
            <h4>漏洞详情</h4>
            <el-button 
              type="primary" 
              size="small" 
              @click="viewVulnerability(currentAssignment.vulnerability?.id)"
            >
              查看漏洞详情页
            </el-button>
          </div>

          <div v-if="currentAssignment.status === 'pending'" class="action-section">
            <h4>处理任务</h4>
            <div class="action-buttons">
              <el-button 
                type="success" 
                @click="showResponseDialog('accepted')"
                class="action-btn"
              >
                <el-icon><Check /></el-icon>
                <span>接受</span>
              </el-button>
              <el-button 
                type="danger" 
                @click="showResponseDialog('rejected')"
                class="action-btn"
              >
                <el-icon><Close /></el-icon>
                <span>拒绝</span>
              </el-button>
            </div>
          </div>
          
          <div v-if="currentAssignment.status === 'accepted'" class="action-section">
            <h4>处理任务</h4>
            <div class="action-buttons">
              <el-button 
                type="primary" 
                @click="showResponseDialog('fixed')"
                class="action-btn"
              >
                <el-icon><SuccessFilled /></el-icon>
                <span>标记为已修复</span>
              </el-button>
            </div>
          </div>
        </div>
      </template>
    </el-dialog>
    
    <!-- 更新状态对话框 -->
    <el-dialog
      title="更新任务状态"
      v-model="responseDialogVisible"
      width="500px"
    >
      <el-form :model="responseForm" ref="responseFormRef" label-width="100px">
        <el-form-item label="状态">
          <el-tag :type="getStatusType(responseForm.status)" effect="dark">
            {{ getStatusText(responseForm.status) }}
          </el-tag>
        </el-form-item>
        
        <el-form-item label="备注信息" prop="response">
          <el-input
            v-model="responseForm.response"
            type="textarea"
            :rows="4"
            :placeholder="getResponsePlaceholder(responseForm.status)"
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="responseDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitStatusUpdate" :loading="submitLoading">提交</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Check, Close, SuccessFilled } from '@element-plus/icons-vue'
import { getMyAssignments, getAssignmentDetails, updateAssignmentStatus } from '@/api/assignment'

export default {
  name: 'MyAssignments',
  components: {
    Check,
    Close,
    SuccessFilled
  },
  setup() {
    const router = useRouter()
    const loading = ref(false)
    const assignments = ref([])
    const total = ref(0)
    const currentPage = ref(1)
    const pageSize = ref(10)
    const filterStatus = ref('')
    
    const detailDialogVisible = ref(false)
    const responseDialogVisible = ref(false)
    const submitLoading = ref(false)
    const currentAssignment = ref(null)
    const responseFormRef = ref(null)
    
    const responseForm = reactive({
      id: null,
      status: '',
      response: '',
      comment: ''
    })
    
    // 获取任务列表
    const fetchAssignments = async () => {
      loading.value = true
      try {
        const response = await getMyAssignments({
          page: currentPage.value,
          pageSize: pageSize.value,
          status: filterStatus.value
        })
        
        if (response.code === 200) {
          assignments.value = response.data
          total.value = response.total || 0
        } else {
          ElMessage.error(response.message || '获取任务列表失败')
        }
      } catch (error) {
        console.error('获取任务列表失败:', error)
        ElMessage.error('获取任务列表失败')
      } finally {
        loading.value = false
      }
    }
    
    // 获取任务详情
    const fetchAssignmentDetails = async (id) => {
      loading.value = true
      try {
        const response = await getAssignmentDetails(id)
        if (response.code === 200) {
          currentAssignment.value = response.data.assignment
          detailDialogVisible.value = true
        } else {
          ElMessage.error(response.message || '获取任务详情失败')
        }
      } catch (error) {
        console.error('获取任务详情失败:', error)
        ElMessage.error('获取任务详情失败')
      } finally {
        loading.value = false
      }
    }
    
    // 处理状态变更
    const handleStatusChange = () => {
      currentPage.value = 1
      fetchAssignments()
    }
    
    // 处理分页变更
    const handlePageChange = (page) => {
      currentPage.value = page
      fetchAssignments()
    }
    
    // 查看详情
    const handleViewDetail = (row) => {
      fetchAssignmentDetails(row.id)
    }
    
    // 查看漏洞详情
    const viewVulnerability = (id) => {
      if (!id) return
      router.push(`/vulnerabilities/${id}`)
    }
    
    // 显示回复对话框
    const showResponseDialog = (status) => {
      responseForm.id = currentAssignment.value.id
      responseForm.status = status
      responseForm.response = ''
      responseForm.comment = getResponsePlaceholder(status)
      
      responseDialogVisible.value = true
    }
    
    // 处理状态更新
    const handleUpdateStatus = (row, status) => {
      responseForm.id = row.id
      responseForm.status = status
      responseForm.response = ''
      responseForm.comment = getResponsePlaceholder(status)
      
      responseDialogVisible.value = true
    }
    
    // 提交状态更新
    const submitStatusUpdate = async () => {
      submitLoading.value = true
      try {
        const response = await updateAssignmentStatus(responseForm.id, {
          status: responseForm.status,
          response: responseForm.response,
          comment: responseForm.comment || getResponsePlaceholder(responseForm.status)
        })
        
        if (response.code === 200) {
          ElMessage.success('状态更新成功')
          responseDialogVisible.value = false
          detailDialogVisible.value = false
          fetchAssignments()
        } else {
          ElMessage.error(response.message || '状态更新失败')
        }
      } catch (error) {
        console.error('状态更新失败:', error)
        ElMessage.error('状态更新失败')
      } finally {
        submitLoading.value = false
      }
    }
    
    // 根据状态获取提示文字
    const getResponsePlaceholder = (status) => {
      const placeholders = {
        'accepted': '接受任务的备注信息，例如预计完成时间等',
        'rejected': '拒绝任务的原因',
        'fixed': '修复的具体方法和过程描述'
      }
      return placeholders[status] || '请输入处理备注'
    }
    
    // 格式化日期
    const formatDate = (dateString) => {
      if (!dateString) return '未设置'
      const date = new Date(dateString)
      return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
    }
    
    // 格式化日期时间
    const formatDateTime = (dateString) => {
      if (!dateString) return ''
      const date = new Date(dateString)
      return `${formatDate(dateString)} ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`
    }
    
    // 判断是否逾期
    const isOverdue = (dateString) => {
      if (!dateString) return false
      const dueDate = new Date(dateString)
      return dueDate < new Date() && new Date(dateString).toDateString() !== new Date().toDateString()
    }
    
    // 获取状态对应的类型
    const getStatusType = (status) => {
      const types = {
        'pending': 'info',
        'accepted': 'success',
        'rejected': 'danger',
        'fixed': 'primary',
        'closed': 'warning'
      }
      return types[status] || 'info'
    }
    
    // 获取状态文本
    const getStatusText = (status) => {
      const texts = {
        'pending': '待处理',
        'accepted': '已接受',
        'rejected': '已拒绝',
        'fixed': '已修复',
        'closed': '已关闭'
      }
      return texts[status] || '未知状态'
    }
    
    // 获取优先级对应的类型
    const getPriorityType = (priority) => {
      const types = {
        1: 'info',
        2: 'info',
        3: 'warning',
        4: 'warning',
        5: 'danger'
      }
      return types[priority] || 'info'
    }
    
    // 获取优先级文本
    const getPriorityText = (priority) => {
      const texts = {
        1: '低',
        2: '中低',
        3: '中',
        4: '中高',
        5: '高'
      }
      return texts[priority] || '未知'
    }
    
    onMounted(() => {
      fetchAssignments()
    })
    
    return {
      loading,
      assignments,
      total,
      currentPage,
      pageSize,
      filterStatus,
      detailDialogVisible,
      responseDialogVisible,
      submitLoading,
      currentAssignment,
      responseForm,
      responseFormRef,
      handleStatusChange,
      handlePageChange,
      handleViewDetail,
      viewVulnerability,
      showResponseDialog,
      handleUpdateStatus,
      submitStatusUpdate,
      getResponsePlaceholder,
      formatDate,
      formatDateTime,
      isOverdue,
      getStatusType,
      getStatusText,
      getPriorityType,
      getPriorityText
    }
  }
}
</script>

<style scoped>
.app-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}

.link-type {
  color: #409EFF;
  text-decoration: none;
}

.link-type:hover {
  text-decoration: underline;
}

.text-danger {
  color: #f56c6c;
}

.detail-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.detail-header h3 {
  margin: 0;
}

.detail-content {
  margin-top: 20px;
}

.detail-section {
  margin-bottom: 20px;
}

.detail-section h4 {
  font-size: 16px;
  color: #606266;
  margin-bottom: 10px;
}

.action-section {
  margin-top: 30px;
  border-top: 1px solid #ebeef5;
  padding-top: 20px;
}

.action-buttons {
  display: flex;
  justify-content: flex-start;
  gap: 12px;
  margin-top: 10px;
}

.action-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s;
}

.action-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.action-btn .el-icon {
  margin-right: 4px;
}
</style> 