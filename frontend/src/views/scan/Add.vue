<template>
  <div class="scan-add-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <h2>创建扫描任务</h2>
        </div>
      </template>
      
      <el-form :model="formData" :rules="rules" ref="formRef" label-width="100px">
        <!-- 基本信息 -->
        <el-divider content-position="left">基本信息</el-divider>
        
        <el-form-item label="任务名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入任务名称"></el-input>
        </el-form-item>
        
        <el-form-item label="描述" prop="description">
          <el-input type="textarea" v-model="formData.description" placeholder="请输入任务描述"></el-input>
        </el-form-item>
        
        <el-form-item label="任务类型" prop="task_type">
          <el-select v-model="formData.task_type" placeholder="请选择任务类型">
            <el-option label="漏洞扫描" value="vulnerability"></el-option>
            <el-option label="合规检查" value="compliance"></el-option>
            <el-option label="安全评估" value="assessment"></el-option>
          </el-select>
        </el-form-item>
        
        <!-- 目标配置 -->
        <el-divider content-position="left">目标配置</el-divider>
        
        <el-form-item label="扫描目标" prop="targets">
          <el-input
            type="textarea"
            v-model="formData.targets"
            placeholder="输入IP地址、域名或URL，每行一个目标"
            :autosize="{ minRows: 3, maxRows: 5 }"
          ></el-input>
        </el-form-item>
        
        <el-form-item label="目标类型" prop="target_type">
          <el-select v-model="formData.target_type" placeholder="请选择目标类型">
            <el-option label="Web应用" value="web"></el-option>
            <el-option label="服务器" value="server"></el-option>
            <el-option label="网络设备" value="network"></el-option>
            <el-option label="移动应用" value="mobile"></el-option>
          </el-select>
        </el-form-item>
        
        <!-- 扫描配置 -->
        <el-divider content-position="left">扫描配置</el-divider>
        
        <el-form-item label="扫描器" prop="scanner">
          <el-select v-model="formData.scanner" placeholder="请选择扫描器">
            <el-option 
              v-for="scanner in scanners" 
              :key="scanner.id" 
              :label="scanner.name" 
              :value="scanner.id"
            ></el-option>
          </el-select>
        </el-form-item>
        
        <el-form-item label="扫描策略" prop="scan_policy">
          <el-select v-model="formData.scan_policy" placeholder="请选择扫描策略">
            <el-option label="全面扫描" value="full"></el-option>
            <el-option label="快速扫描" value="quick"></el-option>
            <el-option label="自定义" value="custom"></el-option>
          </el-select>
        </el-form-item>
        
        <el-form-item label="计划类型" prop="schedule_type">
          <el-radio-group v-model="formData.schedule_type">
            <el-radio label="once">单次执行</el-radio>
            <el-radio label="recurring">定期执行</el-radio>
          </el-radio-group>
        </el-form-item>
        
        <el-form-item label="开始时间" prop="start_time">
          <el-date-picker
            v-model="formData.start_time"
            type="datetime"
            placeholder="选择开始时间"
          ></el-date-picker>
        </el-form-item>
        
        <el-form-item label="重复间隔" prop="interval" v-if="formData.schedule_type === 'recurring'">
          <el-select v-model="formData.interval" placeholder="请选择重复间隔">
            <el-option label="每天" value="daily"></el-option>
            <el-option label="每周" value="weekly"></el-option>
            <el-option label="每月" value="monthly"></el-option>
          </el-select>
        </el-form-item>
        
        <!-- 高级选项 -->
        <el-divider content-position="left">高级选项</el-divider>
        
        <el-form-item label="优先级" prop="priority">
          <el-select v-model="formData.priority" placeholder="请选择优先级">
            <el-option label="高" value="high"></el-option>
            <el-option label="中" value="medium"></el-option>
            <el-option label="低" value="low"></el-option>
          </el-select>
        </el-form-item>
        
        <el-form-item label="通知" prop="notification">
          <el-checkbox v-model="formData.notify_on_complete">扫描完成时通知</el-checkbox>
        </el-form-item>
        
        <el-form-item label="授权" prop="authorization">
          <el-checkbox v-model="formData.has_authorization">已获得授权</el-checkbox>
          <div class="form-tip">请确保您已获得扫描目标的授权，未经授权的扫描可能违反法律</div>
        </el-form-item>
        
        <el-form-item>
          <el-button type="primary" @click="submitForm" :loading="submitting">
            保存
          </el-button>
          <el-button @click="resetForm">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getScanners, createScanTask } from '@/api/scan'

export default {
  name: 'AddScanTask',
  setup() {
    const router = useRouter()
    const formRef = ref(null)
    const submitting = ref(false)
    const scanners = ref([])
    
    // 表单数据
    const formData = reactive({
      name: '',
      description: '',
      task_type: 'vulnerability',
      targets: '',
      target_type: 'web',
      scanner: '',
      scan_policy: 'full',
      schedule_type: 'once',
      start_time: new Date(Date.now() + 30 * 60 * 1000), // 当前时间30分钟后
      interval: 'weekly',
      priority: 'medium',
      notify_on_complete: true,
      has_authorization: false
    })
    
    // 表单验证规则
    const rules = {
      name: [
        { required: true, message: '请输入任务名称', trigger: 'blur' },
        { min: 3, max: 50, message: '长度在3到50个字符之间', trigger: 'blur' }
      ],
      description: [
        { max: 500, message: '不超过500个字符', trigger: 'blur' }
      ],
      task_type: [
        { required: true, message: '请选择任务类型', trigger: 'change' }
      ],
      targets: [
        { required: true, message: '请输入扫描目标', trigger: 'blur' }
      ],
      target_type: [
        { required: true, message: '请选择目标类型', trigger: 'change' }
      ],
      scanner: [
        { required: true, message: '请选择扫描器', trigger: 'change' }
      ],
      scan_policy: [
        { required: true, message: '请选择扫描策略', trigger: 'change' }
      ],
      start_time: [
        { required: true, message: '请选择开始时间', trigger: 'change' }
      ],
      has_authorization: [
        { 
          validator: (rule, value, callback) => {
            if (!value) {
              callback(new Error('必须确认已获得授权'))
            } else {
              callback()
            }
          }, 
          trigger: 'change' 
        }
      ]
    }
    
    // 获取扫描器列表
    const fetchScanners = async () => {
      try {
        const response = await getScanners()
        if (response.code === 200 && response.data) {
          scanners.value = response.data || []
          
          // 如果有扫描器则默认选择第一个
          if (scanners.value.length > 0) {
            formData.scanner = scanners.value[0].id
          }
        } else {
          ElMessage.warning('没有可用的扫描器')
        }
      } catch (error) {
        console.error('获取扫描器列表失败:', error)
        ElMessage.error('获取扫描器列表失败')
      }
    }
    
    // 提交表单
    const submitForm = async () => {
      if (!formRef.value) return
      
      await formRef.value.validate(async (valid) => {
        if (!valid) return
        
        submitting.value = true
        try {
          // 处理目标格式
          const targets = formData.targets.split('\n').filter(t => t.trim()).map(t => t.trim())
          
          const taskData = {
            ...formData,
            targets
          }
          
          const response = await createScanTask(taskData)
          ElMessage.success('创建任务成功')
          router.push(`/scanning/${response.data.id}`)
        } catch (error) {
          console.error('创建扫描任务失败:', error)
          ElMessage.error(error.response?.data?.message || '创建扫描任务失败')
        } finally {
          submitting.value = false
        }
      })
    }
    
    // 重置表单
    const resetForm = () => {
      if (formRef.value) {
        formRef.value.resetFields()
      }
    }
    
    onMounted(() => {
      fetchScanners()
    })
    
    return {
      formRef,
      formData,
      rules,
      scanners,
      submitting,
      submitForm,
      resetForm
    }
  }
}
</script>

<style scoped>
.scan-add-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header h2 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
}

.form-tip {
  font-size: 12px;
  color: #999;
  margin-top: 5px;
}

.el-divider {
  margin: 20px 0;
}
</style> 