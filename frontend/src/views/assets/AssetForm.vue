<template>
  <div class="asset-form" v-loading="loading">
    <!-- 返回按钮 -->
    <div class="back-link">
      <el-button @click="$router.back()" icon="el-icon-arrow-left" type="text">返回</el-button>
    </div>
    
    <!-- 页面标题 -->
    <h1 class="page-title">
      {{ isEdit ? '编辑资产' : '创建资产' }}
    </h1>
    
    <!-- 资产表单 -->
    <el-card class="form-card">
      <el-form 
        ref="formRef" 
        :model="formData" 
        :rules="rules" 
        label-width="100px"
        label-position="top"
      >
        <!-- 基本信息 -->
        <h2 class="form-section-title">基本信息</h2>
        
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="资产名称" prop="name">
              <el-input v-model="formData.name" placeholder="请输入资产名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="资产类型" prop="type">
              <el-select v-model="formData.type" placeholder="请选择资产类型" style="width: 100%">
                <el-option label="服务器" value="server" />
                <el-option label="网络设备" value="network_device" />
                <el-option label="应用程序" value="application" />
                <el-option label="数据库" value="database" />
                <el-option label="云资源" value="cloud_resource" />
                <el-option label="物联网设备" value="iot_device" />
                <el-option label="其他" value="other" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        
        <el-form-item label="描述" prop="description">
          <el-input 
            v-model="formData.description" 
            type="textarea" 
            :rows="3" 
            placeholder="请输入资产描述"
          />
        </el-form-item>
        
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="IP地址" prop="ipAddress">
              <el-input v-model="formData.ipAddress" placeholder="例如: 192.168.1.1" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="域名" prop="domain">
              <el-input v-model="formData.domain" placeholder="例如: example.com" />
            </el-form-item>
          </el-col>
        </el-row>
        
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="MAC地址" prop="macAddress">
              <el-input v-model="formData.macAddress" placeholder="例如: 00:1A:2B:3C:4D:5E" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="操作系统" prop="operatingSystem">
              <el-input v-model="formData.operatingSystem" placeholder="例如: Windows Server 2019" />
            </el-form-item>
          </el-col>
        </el-row>
        
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="责任人" prop="owner">
              <el-input v-model="formData.owner" placeholder="请输入责任人姓名" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="部门" prop="department">
              <el-input v-model="formData.department" placeholder="请输入所属部门" />
            </el-form-item>
          </el-col>
        </el-row>
        
        <el-form-item label="环境" prop="environment">
          <el-radio-group v-model="formData.environment">
            <el-radio label="production">生产环境</el-radio>
            <el-radio label="testing">测试环境</el-radio>
            <el-radio label="development">开发环境</el-radio>
            <el-radio label="staging">预生产环境</el-radio>
          </el-radio-group>
        </el-form-item>
        
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="formData.status">
            <el-radio label="active">活动</el-radio>
            <el-radio label="inactive">不活动</el-radio>
          </el-radio-group>
        </el-form-item>
        
        <el-form-item label="开放端口" prop="ports">
          <el-tag
            v-for="port in formData.ports"
            :key="port"
            closable
            @close="removePort(port)"
            class="port-tag"
          >
            {{ port }}
          </el-tag>
          <el-input
            v-if="portInputVisible"
            ref="portInputRef"
            v-model="portInput"
            class="port-input"
            size="small"
            @keyup.enter="addPort"
            @blur="addPort"
          />
          <el-button v-else size="small" @click="showPortInput">+ 添加端口</el-button>
        </el-form-item>
        
        <!-- 自定义属性 -->
        <h2 class="form-section-title">自定义属性</h2>
        <p class="section-description">添加此资产的自定义属性</p>
        
        <div v-for="(attribute, index) in attributes" :key="index" class="attribute-item">
          <el-row :gutter="20">
            <el-col :span="10">
              <el-form-item :label="`属性名 ${index + 1}`" :prop="`attributes[${index}].key`" :rules="attributeRules.key">
                <el-input v-model="attribute.key" placeholder="请输入属性名" />
              </el-form-item>
            </el-col>
            <el-col :span="10">
              <el-form-item :label="`属性值 ${index + 1}`" :prop="`attributes[${index}].value`" :rules="attributeRules.value">
                <el-input v-model="attribute.value" placeholder="请输入属性值" />
              </el-form-item>
            </el-col>
            <el-col :span="4" class="delete-col">
              <el-button 
                type="danger" 
                icon="el-icon-delete" 
                circle 
                @click="removeAttribute(index)"
                size="small"
                class="delete-btn"
              />
            </el-col>
          </el-row>
        </div>
        
        <el-button type="dashed" icon="el-icon-plus" @click="addAttribute" class="add-attribute-btn">
          添加自定义属性
        </el-button>
        
        <!-- 提交按钮 -->
        <div class="form-actions">
          <el-button @click="$router.back()">取消</el-button>
          <el-button type="primary" @click="submitForm">{{ isEdit ? '保存修改' : '创建资产' }}</el-button>
        </div>
      </el-form>
    </el-card>
  </div>
</template>

<script>
import { ref, reactive, computed, nextTick, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import assetApi from '@/api/assets'

export default {
  name: 'AssetForm',
  
  setup() {
    const route = useRoute()
    const router = useRouter()
    const formRef = ref(null)
    const loading = ref(false)
    
    // 判断是编辑还是创建
    const isEdit = computed(() => !!route.params.id)
    
    // 表单数据
    const formData = reactive({
      name: '',
      type: '',
      description: '',
      ipAddress: '',
      domain: '',
      macAddress: '',
      operatingSystem: '',
      owner: '',
      department: '',
      environment: 'production',
      status: 'active',
      ports: []
    })
    
    // 端口输入
    const portInputVisible = ref(false)
    const portInputRef = ref(null)
    const portInput = ref('')
    
    // 自定义属性
    const attributes = ref([])
    
    // 表单验证规则
    const rules = {
      name: [
        { required: true, message: '请输入资产名称', trigger: 'blur' },
        { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
      ],
      type: [
        { required: true, message: '请选择资产类型', trigger: 'change' }
      ],
      status: [
        { required: true, message: '请选择资产状态', trigger: 'change' }
      ],
      description: [
        { max: 500, message: '描述不能超过500个字符', trigger: 'blur' }
      ],
      ipAddress: [
        { pattern: /^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/, message: '请输入有效的IP地址', trigger: 'blur' }
      ],
      domain: [
        { pattern: /^([a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$/, message: '请输入有效的域名', trigger: 'blur' }
      ],
      macAddress: [
        { pattern: /^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$/, message: '请输入有效的MAC地址', trigger: 'blur' }
      ],
      environment: [
        { required: true, message: '请选择环境', trigger: 'change' }
      ]
    }
    
    // 自定义属性验证规则
    const attributeRules = {
      key: [
        { required: true, message: '请输入属性名', trigger: 'blur' },
        { max: 50, message: '属性名不能超过50个字符', trigger: 'blur' }
      ],
      value: [
        { required: true, message: '请输入属性值', trigger: 'blur' },
        { max: 200, message: '属性值不能超过200个字符', trigger: 'blur' }
      ]
    }
    
    // 获取资产详情
    const fetchAssetDetails = async () => {
      if (!isEdit.value) return
      
      loading.value = true
      try {
        const response = await assetApi.getAssetById(route.params.id)
        const asset = response.data
        
        // 填充表单数据
        Object.keys(formData).forEach(key => {
          if (key in asset) {
            formData[key] = asset[key]
          }
        })
        
        // 处理自定义属性
        if (asset.attributes) {
          attributes.value = Object.entries(asset.attributes).map(([key, value]) => ({
            key,
            value: typeof value === 'object' ? JSON.stringify(value) : value
          }))
        }
      } catch (error) {
        console.error('获取资产详情失败:', error)
        ElMessage.error('获取资产详情失败，请重试')
      } finally {
        loading.value = false
      }
    }
    
    // 显示端口输入框
    const showPortInput = () => {
      portInputVisible.value = true
      nextTick(() => {
        portInputRef.value.focus()
      })
    }
    
    // 添加端口
    const addPort = () => {
      if (portInput.value && portInput.value.trim()) {
        const port = parseInt(portInput.value.trim())
        if (!isNaN(port) && port > 0 && port <= 65535) {
          if (!formData.ports.includes(port)) {
            formData.ports.push(port)
          }
        } else {
          ElMessage.warning('请输入有效的端口号(1-65535)')
        }
      }
      portInputVisible.value = false
      portInput.value = ''
    }
    
    // 移除端口
    const removePort = (port) => {
      formData.ports = formData.ports.filter(p => p !== port)
    }
    
    // 添加自定义属性
    const addAttribute = () => {
      attributes.value.push({
        key: '',
        value: ''
      })
    }
    
    // 移除自定义属性
    const removeAttribute = (index) => {
      attributes.value.splice(index, 1)
    }
    
    // 提交表单
    const submitForm = async () => {
      if (!formRef.value) return
      
      await formRef.value.validate(async (valid) => {
        if (!valid) {
          ElMessage.error('请正确填写表单信息')
          return
        }
        
        // 处理自定义属性
        const attributesObj = {}
        attributes.value.forEach(attr => {
          if (attr.key && attr.value) {
            attributesObj[attr.key] = attr.value
          }
        })
        
        const assetData = {
          ...formData,
          attributes: attributesObj
        }
        
        loading.value = true
        try {
          if (isEdit.value) {
            // 编辑资产
            await assetApi.updateAsset(route.params.id, assetData)
            ElMessage.success('资产更新成功')
          } else {
            // 创建资产
            await assetApi.createAsset(assetData)
            ElMessage.success('资产创建成功')
          }
          
          // 返回资产列表
          router.push({ name: 'AssetList' })
        } catch (error) {
          console.error('操作失败:', error)
          ElMessage.error(`${isEdit.value ? '更新' : '创建'}资产失败，请重试`)
        } finally {
          loading.value = false
        }
      })
    }
    
    onMounted(() => {
      fetchAssetDetails()
    })
    
    return {
      formRef,
      formData,
      rules,
      attributeRules,
      loading,
      isEdit,
      attributes,
      portInputVisible,
      portInputRef,
      portInput,
      showPortInput,
      addPort,
      removePort,
      addAttribute,
      removeAttribute,
      submitForm
    }
  }
}
</script>

<style scoped>
.asset-form {
  padding: 20px;
}

.back-link {
  margin-bottom: 10px;
}

.page-title {
  font-size: 24px;
  color: #303133;
  margin-bottom: 20px;
}

.form-card {
  margin-bottom: 20px;
}

.form-section-title {
  font-size: 18px;
  color: #303133;
  margin: 20px 0 10px 0;
  padding-bottom: 10px;
  border-bottom: 1px solid #ebeef5;
}

.section-description {
  color: #909399;
  font-size: 14px;
  margin-bottom: 15px;
}

.port-tag {
  margin-right: 10px;
  margin-bottom: 10px;
}

.port-input {
  width: 100px;
  margin-right: 10px;
  vertical-align: bottom;
}

.attribute-item {
  margin-bottom: 10px;
  padding-bottom: 10px;
  border-bottom: 1px dashed #ebeef5;
}

.delete-col {
  display: flex;
  align-items: flex-end;
  justify-content: center;
  margin-bottom: 22px;
}

.delete-btn {
  margin-bottom: 10px;
}

.add-attribute-btn {
  margin-bottom: 20px;
  width: 100%;
  border-style: dashed;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
  gap: 10px;
}
</style> 