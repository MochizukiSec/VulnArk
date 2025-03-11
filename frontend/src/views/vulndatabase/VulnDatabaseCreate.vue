<template>
  <div class="vuln-create-container">
    <el-page-header :title="'漏洞库'" @back="goBack">
      <template #content>
        <span class="header-title">创建新漏洞</span>
      </template>
    </el-page-header>
    
    <div class="form-container">
      <el-form 
        ref="vulnFormRef" 
        :model="vulnForm" 
        :rules="rules" 
        label-position="top" 
        class="vuln-form"
        v-loading="loading"
      >
        <el-card class="form-card">
          <template #header>
            <div class="card-header">
              <h3>基本信息</h3>
              <div class="required-note">* 必填项</div>
            </div>
          </template>
          
          <el-row :gutter="20">
            <el-col :sm="24" :md="12">
              <el-form-item label="CVE ID" prop="cveId">
                <el-input 
                  v-model="vulnForm.cveId" 
                  placeholder="例如: CVE-2023-1234" 
                />
              </el-form-item>
            </el-col>
            
            <el-col :sm="24" :md="12">
              <el-form-item label="标题" prop="title" required>
                <el-input 
                  v-model="vulnForm.title" 
                  placeholder="输入漏洞标题" 
                />
              </el-form-item>
            </el-col>
            
            <el-col :sm="24" :md="12">
              <el-form-item label="CVSS评分" prop="cvss">
                <el-input-number 
                  v-model="vulnForm.cvss" 
                  :min="0" 
                  :max="10" 
                  :precision="1" 
                  :step="0.1" 
                  controls-position="right" 
                  style="width: 100%;"
                />
              </el-form-item>
            </el-col>
            
            <el-col :sm="24" :md="12">
              <el-form-item label="严重程度" prop="severity" required>
                <el-select v-model="vulnForm.severity" placeholder="选择严重程度" style="width: 100%;">
                  <el-option label="严重" value="critical" />
                  <el-option label="高危" value="high" />
                  <el-option label="中危" value="medium" />
                  <el-option label="低危" value="low" />
                  <el-option label="信息" value="info" />
                </el-select>
              </el-form-item>
            </el-col>
          </el-row>
        </el-card>
        
        <el-card class="form-card">
          <template #header>
            <div class="card-header">
              <h3>详细描述</h3>
            </div>
          </template>
          
          <el-form-item label="漏洞描述" prop="description" required>
            <el-input 
              v-model="vulnForm.description" 
              type="textarea" 
              :rows="5" 
              placeholder="请详细描述此漏洞，包括其影响、利用方式等" 
            />
          </el-form-item>
          
          <el-form-item label="受影响系统" prop="affectedSystems">
            <el-input 
              v-model="vulnForm.affectedSystems" 
              placeholder="列出受影响的操作系统、软件或硬件" 
            />
          </el-form-item>
          
          <el-form-item label="解决方案" prop="solution">
            <el-input 
              v-model="vulnForm.solution" 
              type="textarea" 
              :rows="3" 
              placeholder="提供修复或缓解此漏洞的建议措施" 
            />
          </el-form-item>
        </el-card>
        
        <el-card class="form-card">
          <template #header>
            <div class="card-header">
              <h3>时间信息</h3>
            </div>
          </template>
          
          <el-row :gutter="20">
            <el-col :sm="24" :md="12">
              <el-form-item label="发布日期" prop="publishedDate">
                <el-date-picker 
                  v-model="vulnForm.publishedDate" 
                  type="date" 
                  placeholder="选择发布日期" 
                  style="width: 100%;"
                  value-format="YYYY-MM-DD"
                />
              </el-form-item>
            </el-col>
            
            <el-col :sm="24" :md="12">
              <el-form-item label="更新日期" prop="lastModifiedDate">
                <el-date-picker 
                  v-model="vulnForm.lastModifiedDate" 
                  type="date" 
                  placeholder="选择最后更新日期" 
                  style="width: 100%;"
                  value-format="YYYY-MM-DD"
                />
              </el-form-item>
            </el-col>
          </el-row>
        </el-card>
        
        <el-card class="form-card">
          <template #header>
            <div class="card-header">
              <h3>参考资料</h3>
            </div>
          </template>
          
          <el-form-item label="参考链接" prop="references">
            <div v-for="(reference, index) in vulnForm.references" :key="index" class="reference-item">
              <el-input 
                v-model="vulnForm.references[index]" 
                placeholder="输入参考链接URL" 
                class="reference-input"
              >
                <template #append>
                  <el-button @click="removeReference(index)" type="danger">
                    <el-icon><Delete /></el-icon>
                  </el-button>
                </template>
              </el-input>
            </div>
            
            <el-button @click="addReference" type="primary" plain class="add-reference-btn">
              <el-icon><Plus /></el-icon> 添加参考链接
            </el-button>
          </el-form-item>
        </el-card>
        
        <div class="form-actions">
          <el-button @click="goBack">取消</el-button>
          <el-button type="primary" @click="submitForm" :loading="submitting">保存漏洞</el-button>
        </div>
      </el-form>
    </div>
  </div>
</template>

<script>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Plus, Delete } from '@element-plus/icons-vue'
import axios from 'axios'

export default {
  name: 'VulnDatabaseCreate',
  
  components: {
    Plus,
    Delete
  },
  
  setup() {
    const router = useRouter()
    const vulnFormRef = ref(null)
    const loading = ref(false)
    const submitting = ref(false)
    
    // 表单数据
    const vulnForm = reactive({
      cveId: '',
      title: '',
      description: '',
      cvss: 5.0,
      severity: '',
      affectedSystems: '',
      solution: '',
      publishedDate: '',
      lastModifiedDate: '',
      references: ['']
    })
    
    // 表单验证规则
    const rules = {
      title: [
        { required: true, message: '请输入漏洞标题', trigger: 'blur' },
        { min: 5, max: 100, message: '标题长度应在5-100个字符之间', trigger: 'blur' }
      ],
      description: [
        { required: true, message: '请输入漏洞描述', trigger: 'blur' },
        { min: 10, message: '描述应至少包含10个字符', trigger: 'blur' }
      ],
      severity: [
        { required: true, message: '请选择严重程度', trigger: 'change' }
      ],
      cveId: [
        { pattern: /^CVE-\d{4}-\d{4,}$/, message: 'CVE ID格式应为CVE-YYYY-XXXXX', trigger: 'blur' }
      ],
      cvss: [
        { type: 'number', min: 0, max: 10, message: 'CVSS评分应在0-10之间', trigger: 'blur' }
      ],
      references: [
        { 
          type: 'array', 
          validator: (rule, value, callback) => {
            const invalidUrls = value.filter(url => url && !isValidUrl(url))
            if (invalidUrls.length > 0) {
              callback(new Error('请输入有效的URL地址'))
            } else {
              callback()
            }
          }, 
          trigger: 'blur' 
        }
      ]
    }
    
    // 验证URL格式
    const isValidUrl = (string) => {
      try {
        new URL(string)
        return true
      } catch (_) {
        return false
      }
    }
    
    // 添加参考链接
    const addReference = () => {
      vulnForm.references.push('')
    }
    
    // 移除参考链接
    const removeReference = (index) => {
      vulnForm.references.splice(index, 1)
      if (vulnForm.references.length === 0) {
        vulnForm.references.push('')
      }
    }
    
    // 返回漏洞库列表
    const goBack = () => {
      router.push('/vulndatabase')
    }
    
    // 提交表单
    const submitForm = async () => {
      // 移除空的参考链接
      vulnForm.references = vulnForm.references.filter(ref => ref.trim() !== '')
      
      if (!vulnForm.lastModifiedDate) {
        vulnForm.lastModifiedDate = vulnForm.publishedDate
      }
      
      submitting.value = true
      
      console.log('准备提交表单数据:', vulnForm);
      console.log('提交URL:', axios.defaults.baseURL + '/api/vulndatabase');
      
      try {
        console.log('创建漏洞，URL:', axios.defaults.baseURL + '/api/vulndatabase');
        console.log('请求数据:', JSON.stringify(vulnForm, null, 2));
        const response = await axios.post('/api/vulndatabase', vulnForm);
        console.log('创建漏洞成功:', response.data);
        ElMessage.success('成功创建漏洞');
        router.push(`/vulndatabase/${response.data.cveId || response.data.id}`);
      } catch (error) {
        console.error('创建漏洞失败，详细信息:', {
          error: error,
          message: error.message,
          config: error.config,
          response: error.response ? {
            status: error.response.status,
            statusText: error.response.statusText,
            data: error.response.data
          } : '无响应数据'
        });
        ElMessage.error('创建漏洞失败，请稍后重试');
      } finally {
        submitting.value = false;
      }
    }
    
    return {
      vulnFormRef,
      vulnForm,
      rules,
      loading,
      submitting,
      addReference,
      removeReference,
      goBack,
      submitForm,
      Plus,
      Delete
    }
  }
}
</script>

<style scoped>
.vuln-create-container {
  padding: 20px;
}

.header-title {
  font-size: 16px;
  font-weight: 500;
}

.form-container {
  margin-top: 20px;
}

.form-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 500;
}

.required-note {
  color: #f56c6c;
  font-size: 14px;
}

.reference-item {
  margin-bottom: 10px;
}

.reference-input {
  width: 100%;
}

.add-reference-btn {
  margin-top: 10px;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 20px;
}

@media (max-width: 768px) {
  .form-actions {
    flex-direction: column;
  }
  
  .form-actions .el-button {
    width: 100%;
    margin-left: 0;
    margin-bottom: 10px;
  }
}
</style> 