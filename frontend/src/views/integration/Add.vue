<template>
  <div class="add-integration-container">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <span>添加CI/CD集成</span>
          <el-button @click="goBack" class="action-btn">
            <el-icon><Back /></el-icon>
            <span>返回</span>
          </el-button>
        </div>
      </template>
      
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="100px"
        class="integration-form"
      >
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入集成名称" />
        </el-form-item>
        
        <el-form-item label="类型" prop="type">
          <el-select v-model="form.type" placeholder="请选择集成类型" style="width: 100%">
            <el-option
              v-for="item in integrationTypes"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            >
              <span style="display: flex; align-items: center;">
                <el-icon class="mr-2">
                  <component :is="item.icon"></component>
                </el-icon>
                {{ item.label }}
              </span>
            </el-option>
          </el-select>
        </el-form-item>
        
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="3"
            placeholder="请输入集成描述"
          />
        </el-form-item>
        
        <el-form-item label="状态" prop="enabled">
          <el-switch v-model="form.enabled" />
        </el-form-item>
        
        <!-- Jenkins配置 -->
        <template v-if="form.type === 'jenkins'">
          <el-divider content-position="left">Jenkins配置</el-divider>
          <el-form-item label="Jenkins URL" prop="jenkins_url">
            <el-input v-model="jenkinsConfig.url" placeholder="例如: https://jenkins.example.com" />
          </el-form-item>
          <el-form-item label="用户名" prop="jenkins_username">
            <el-input v-model="jenkinsConfig.username" placeholder="Jenkins用户名" />
          </el-form-item>
          <el-form-item label="API Token" prop="jenkins_token">
            <el-input v-model="jenkinsConfig.token" show-password placeholder="Jenkins API Token" />
          </el-form-item>
        </template>
        
        <!-- GitLab配置 -->
        <template v-if="form.type === 'gitlab'">
          <el-divider content-position="left">GitLab配置</el-divider>
          <el-form-item label="GitLab URL" prop="gitlab_url">
            <el-input v-model="gitlabConfig.url" placeholder="例如: https://gitlab.example.com" />
          </el-form-item>
          <el-form-item label="私有Token" prop="gitlab_token">
            <el-input v-model="gitlabConfig.token" show-password placeholder="GitLab私有令牌" />
          </el-form-item>
          <el-form-item label="项目ID" prop="gitlab_project_id">
            <el-input v-model="gitlabConfig.project_id" placeholder="GitLab项目ID" />
          </el-form-item>
        </template>
        
        <!-- GitHub配置 -->
        <template v-if="form.type === 'github'">
          <el-divider content-position="left">GitHub配置</el-divider>
          <el-form-item label="仓库所有者" prop="github_owner">
            <el-input v-model="githubConfig.owner" placeholder="例如: octocat" />
          </el-form-item>
          <el-form-item label="仓库名称" prop="github_repo">
            <el-input v-model="githubConfig.repo" placeholder="例如: hello-world" />
          </el-form-item>
          <el-form-item label="访问令牌" prop="github_token">
            <el-input v-model="githubConfig.token" show-password placeholder="GitHub Personal Access Token" />
          </el-form-item>
        </template>
        
        <!-- 自定义配置 -->
        <template v-if="form.type === 'custom'">
          <el-divider content-position="left">自定义配置</el-divider>
          <el-form-item label="配置项" prop="custom_config">
            <el-input
              v-model="customConfig.config"
              type="textarea"
              :rows="5"
              placeholder="请输入JSON格式的配置数据"
            />
          </el-form-item>
        </template>
        
        <el-form-item>
          <el-button type="primary" @click="submitForm" :loading="submitting">创建集成</el-button>
          <el-button @click="resetForm">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { 
  Back, SetUp, Guide, Link, Document
} from '@element-plus/icons-vue'
import { createIntegration } from '@/api/integration'

const router = useRouter()
const formRef = ref(null)
const submitting = ref(false)

// 表单数据
const form = reactive({
  name: '',
  type: '',
  description: '',
  enabled: true,
  config: ''
})

// Jenkins配置
const jenkinsConfig = reactive({
  url: '',
  username: '',
  token: ''
})

// GitLab配置
const gitlabConfig = reactive({
  url: '',
  token: '',
  project_id: ''
})

// GitHub配置
const githubConfig = reactive({
  owner: '',
  repo: '',
  token: ''
})

// 自定义配置
const customConfig = reactive({
  config: ''
})

// 当表单提交时，根据选择的集成类型自动填充config字段
const updateConfig = () => {
  switch (form.type) {
    case 'jenkins':
      form.config = JSON.stringify(jenkinsConfig)
      break
    case 'gitlab':
      form.config = JSON.stringify(gitlabConfig)
      break
    case 'github':
      form.config = JSON.stringify(githubConfig)
      break
    case 'custom':
      form.config = customConfig.config
      break
    default:
      form.config = '{}'
  }
}

// 集成类型列表
const integrationTypes = [
  { label: 'Jenkins', value: 'jenkins', icon: SetUp },
  { label: 'GitLab', value: 'gitlab', icon: Guide },
  { label: 'GitHub', value: 'github', icon: Link },
  { label: '自定义', value: 'custom', icon: Document }
]

// 表单验证规则
const rules = {
  name: [
    { required: true, message: '请输入集成名称', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在2到50个字符之间', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请选择集成类型', trigger: 'change' }
  ]
}

// 提交表单
const submitForm = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    
    // 更新配置
    updateConfig()
    
    submitting.value = true
    try {
      const response = await createIntegration(form)
      if (response.code === 200) {
        ElMessage.success('创建集成成功')
        router.push('/integrations')
      } else {
        ElMessage.error(response.message || '创建集成失败')
      }
    } catch (error) {
      console.error('创建集成失败:', error)
      ElMessage.error('创建集成失败')
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
  
  // 重置配置
  Object.keys(jenkinsConfig).forEach(key => jenkinsConfig[key] = '')
  Object.keys(gitlabConfig).forEach(key => gitlabConfig[key] = '')
  Object.keys(githubConfig).forEach(key => githubConfig[key] = '')
  customConfig.config = ''
}

// 返回列表页
const goBack = () => {
  router.push('/integrations')
}
</script>

<style scoped>
.add-integration-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.integration-form {
  margin-top: 20px;
  max-width: 800px;
}

.action-btn {
  display: flex;
  align-items: center;
}

.action-btn .el-icon {
  margin-right: 4px;
}

.mr-2 {
  margin-right: 8px;
}
</style> 