<template>
  <div class="edit-scan-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <h2>{{ $t('scan.editTask') }}</h2>
          <el-button @click="goBack">{{ $t('common.back') }}</el-button>
        </div>
      </template>

      <div v-if="!loading">
        <el-form 
          ref="scanForm" 
          :model="scanForm" 
          :rules="rules" 
          label-width="120px"
          label-position="left"
        >
          <!-- 基本信息 -->
          <el-form-item :label="$t('common.name')" prop="name">
            <el-input v-model="scanForm.name" />
          </el-form-item>

          <el-form-item :label="$t('common.description')" prop="description">
            <el-input v-model="scanForm.description" type="textarea" :rows="3" />
          </el-form-item>

          <el-form-item :label="$t('scan.scannerType.label')" prop="type">
            <el-select v-model="scanForm.type" style="width: 100%">
              <el-option
                v-for="item in scannerTypes"
                :key="item.value"
                :label="$t(`scan.scannerType.${item.value}`)"
                :value="item.value"
              />
            </el-select>
          </el-form-item>

          <!-- 扫描器配置 -->
          <div class="section-header">
            <h3>{{ $t('scan.scannerConfig') }}</h3>
          </div>

          <el-form-item :label="$t('scan.scannerUrl')" prop="scanner_url">
            <el-input v-model="scanForm.scanner_url" placeholder="https://scanner.example.com" />
          </el-form-item>

          <el-form-item :label="$t('scan.scannerApiKey')" prop="scanner_api_key">
            <el-input v-model="scanForm.scanner_api_key" show-password placeholder="Leave empty to keep current API key" />
          </el-form-item>

          <el-form-item :label="$t('scan.scannerUsername')" prop="scanner_username">
            <el-input v-model="scanForm.scanner_username" />
          </el-form-item>

          <el-form-item :label="$t('scan.scannerPassword')" prop="scanner_password">
            <el-input v-model="scanForm.scanner_password" show-password placeholder="Leave empty to keep current password" />
          </el-form-item>

          <!-- 目标配置 -->
          <div class="section-header">
            <h3>{{ $t('scan.targetConfig') }}</h3>
          </div>

          <el-form-item :label="$t('scan.targetIPs')" prop="target_ips">
            <el-input 
              v-model="scanForm.target_ips" 
              type="textarea" 
              :rows="3"
              :placeholder="$t('scan.targetIPsPlaceholder')"
            />
          </el-form-item>

          <el-form-item :label="$t('scan.targetURLs')" prop="target_urls">
            <el-input 
              v-model="scanForm.target_urls" 
              type="textarea" 
              :rows="3"
              :placeholder="$t('scan.targetURLsPlaceholder')"
            />
          </el-form-item>

          <el-form-item :label="$t('scan.targetAssets')" prop="target_assets">
            <el-select 
              v-model="scanForm.target_assets" 
              multiple 
              filterable 
              style="width: 100%"
            >
              <el-option
                v-for="asset in assets"
                :key="asset.id"
                :label="asset.name"
                :value="asset.id.toString()"
              />
            </el-select>
          </el-form-item>

          <el-form-item :label="$t('scan.scanParameters')" prop="scan_parameters">
            <el-input
              v-model="scanForm.scan_parameters"
              type="textarea"
              :rows="5"
              :placeholder="$t('scan.scanParametersPlaceholder')"
            />
          </el-form-item>

          <!-- 计划配置 -->
          <div class="section-header">
            <h3>{{ $t('scan.scheduleConfig') }}</h3>
          </div>

          <el-form-item :label="$t('scan.scheduleAt')" prop="scheduled_at">
            <el-date-picker
              v-model="scanForm.scheduled_at"
              type="datetime"
              :placeholder="$t('scan.scheduleAt')"
              format="YYYY-MM-DD HH:mm:ss"
              style="width: 100%"
            />
          </el-form-item>

          <el-form-item :label="$t('scan.isRecurring')" prop="is_recurring">
            <el-switch v-model="scanForm.is_recurring" />
          </el-form-item>

          <el-form-item 
            v-if="scanForm.is_recurring" 
            :label="$t('scan.cronSchedule')" 
            prop="cron_schedule"
          >
            <el-input 
              v-model="scanForm.cron_schedule" 
              :placeholder="$t('scan.cronSchedulePlaceholder')" 
            />
          </el-form-item>

          <!-- 提交按钮 -->
          <el-form-item>
            <el-button type="primary" @click="submitForm" :loading="submitting">
              {{ $t('common.save') }}
            </el-button>
            <el-button @click="resetForm">{{ $t('common.reset') }}</el-button>
          </el-form-item>
        </el-form>
      </div>
    </el-card>
  </div>
</template>

<script>
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import api from '@/api'

export default {
  name: 'EditScanTask',
  setup() {
    const { t } = useI18n()
    const route = useRoute()
    const router = useRouter()
    const formRef = ref(null)
    const loading = ref(true)
    const submitting = ref(false)
    const assets = ref([])
    
    // 获取任务ID
    const taskId = route.params.id
    
    // 扫描器类型
    const scannerTypes = [
      { value: 'nessus', label: 'Nessus' },
      { value: 'xray', label: 'Xray' },
      { value: 'awvs', label: 'AWVS' },
      { value: 'zap', label: 'OWASP ZAP' },
      { value: 'custom', label: 'Custom' }
    ]
    
    // 表单数据
    const formData = reactive({
      name: '',
      description: '',
      type: '',
      scanner_url: '',
      scanner_api_key: '',
      scanner_username: '',
      scanner_password: '',
      target_ips: '',
      target_urls: '',
      target_assets: [],
      scan_parameters: '',
      scheduled_at: null,
      is_recurring: false,
      cron_schedule: ''
    })
    
    // 表单验证规则
    const rules = reactive({
      name: [
        { required: true, message: t('common.required', { field: t('common.name') }), trigger: 'blur' },
        { min: 3, max: 50, message: t('common.lengthLimit', { min: 3, max: 50 }), trigger: 'blur' }
      ],
      type: [
        { required: true, message: t('common.required', { field: t('scan.scannerType.label') }), trigger: 'change' }
      ]
    })
    
    // 加载任务数据
    const loadTaskData = async () => {
      loading.value = true
      try {
        const res = await api.get(`/scans/${taskId}`)
        const taskData = res.data
        
        // 填充表单数据
        Object.keys(formData).forEach(key => {
          if (key in taskData) {
            formData[key] = taskData[key]
          }
        })
        
        // 处理目标类型
        if (taskData.target_assets) {
          formData.target_type = 'assets'
          formData.target_assets = taskData.target_assets.split(',').map(id => parseInt(id.trim()))
        } else if (taskData.target_urls) {
          formData.target_type = 'urls'
        } else if (taskData.target_ips) {
          formData.target_type = 'ips'
        }
        
        // 加载资产列表
        await loadAssets()
      } catch (error) {
        console.error('加载任务数据失败:', error)
        ElMessage.error(t('common.error') + ': ' + (error.response?.data?.message || error.message))
      } finally {
        loading.value = false
      }
    }
    
    // 加载资产列表
    const loadAssets = async () => {
      try {
        const res = await api.get('/assets')
        assets.value = res.data.assets || []
      } catch (error) {
        console.error('加载资产列表失败:', error)
        ElMessage.error(t('common.error') + ': ' + (error.response?.data?.message || error.message))
      }
    }
    
    // 表单提交
    const submitForm = async () => {
      if (!formRef.value) return
      
      await formRef.value.validate(async (valid) => {
        if (valid) {
          loading.value = true
          try {
            const submitData = { ...formData }
            
            // 处理资产ID
            if (formData.target_type === 'assets' && formData.target_assets.length) {
              submitData.target_assets = formData.target_assets.join(',')
            }
            
            // 提交数据
            await api.put(`/scans/${taskId}`, submitData)
            ElMessage.success(t('scan.updateSuccess'))
            router.push({ name: 'ScanDetail', params: { id: taskId } })
          } catch (error) {
            console.error('更新扫描任务失败:', error)
            ElMessage.error(t('common.error') + ': ' + (error.response?.data?.message || error.message))
          } finally {
            loading.value = false
          }
        } else {
          ElMessage.warning(t('common.formError'))
          return false
        }
      })
    }
    
    // 重置表单
    const resetForm = () => {
      loadTaskData() // 重新加载任务数据
    }
    
    // 返回任务详情
    const goBack = () => {
      router.push(`/scanning/${taskId}`)
    }
    
    onMounted(async () => {
      await loadTaskData()
    })
    
    return {
      formRef,
      formData,
      rules,
      loading,
      submitting,
      assets,
      scannerTypes,
      taskId,
      submitForm,
      resetForm,
      goBack
    }
  }
}
</script>

<style scoped>
.edit-scan-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.section-header {
  margin: 20px 0 10px 0;
  padding-bottom: 10px;
  border-bottom: 1px solid #eee;
}
</style> 