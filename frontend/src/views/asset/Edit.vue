<template>
  <div class="asset-edit-container">
    <div class="page-header">
      <div class="left">
        <el-button link @click="goBack">
          <el-icon><Back /></el-icon>
          返回
        </el-button>
        <h2>编辑资产</h2>
      </div>
    </div>

    <el-card class="form-card">
      <el-form 
        ref="assetFormRef" 
        :model="assetForm" 
        :rules="rules" 
        label-width="100px"
        v-loading="loading"
      >
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="资产名称" prop="name">
              <el-input v-model="assetForm.name" placeholder="请输入资产名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="IP地址" prop="ip">
              <el-input v-model="assetForm.ip" placeholder="请输入IP地址" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="资产类型" prop="type">
              <el-select v-model="assetForm.type" placeholder="请选择资产类型" style="width: 100%">
                <el-option label="服务器" value="server" />
                <el-option label="网站" value="website" />
                <el-option label="数据库" value="database" />
                <el-option label="网络设备" value="network" />
                <el-option label="应用程序" value="application" />
                <el-option label="其他" value="other" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="状态" prop="status">
              <el-select v-model="assetForm.status" placeholder="请选择状态" style="width: 100%">
                <el-option label="在线" value="active" />
                <el-option label="离线" value="inactive" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="所属部门" prop="department">
              <el-input v-model="assetForm.department" placeholder="请输入所属部门" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="负责人" prop="owner">
              <el-input v-model="assetForm.owner" placeholder="请输入负责人" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="操作系统" prop="operatingSystem">
              <el-input v-model="assetForm.operatingSystem" placeholder="请输入操作系统" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="版本" prop="version">
              <el-input v-model="assetForm.version" placeholder="请输入版本信息" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="URL地址" prop="url" v-if="assetForm.type === 'website'">
          <el-input v-model="assetForm.url" placeholder="请输入URL地址" />
        </el-form-item>

        <el-form-item label="MAC地址" prop="macAddress">
          <el-input v-model="assetForm.macAddress" placeholder="请输入MAC地址" />
        </el-form-item>

        <el-form-item label="标签" prop="tags">
          <el-tag
            v-for="tag in assetForm.tags"
            :key="tag"
            class="tag-item"
            closable
            @close="removeTag(tag)"
          >
            {{ tag }}
          </el-tag>
          <el-input
            v-if="inputTagVisible"
            ref="tagInputRef"
            v-model="inputTagValue"
            class="tag-input"
            size="small"
            @keyup.enter="handleInputTagConfirm"
            @blur="handleInputTagConfirm"
          />
          <el-button v-else class="button-new-tag" size="small" @click="showTagInput">
            + 新标签
          </el-button>
        </el-form-item>

        <el-form-item label="备注" prop="notes">
          <el-input
            v-model="assetForm.notes"
            type="textarea"
            :rows="3"
            placeholder="请输入备注信息"
          />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="submitForm">保存</el-button>
          <el-button @click="resetForm">重置</el-button>
          <el-button @click="goBack">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script>
import { ref, reactive, nextTick, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Back } from '@element-plus/icons-vue'
import { getAssetDetail, updateAsset } from '@/api/asset'

export default {
  name: 'AssetEdit',
  components: {
    Back
  },
  setup() {
    const router = useRouter()
    const route = useRoute()
    const assetFormRef = ref(null)
    const loading = ref(false)
    const assetId = route.params.id
    
    // 表单数据
    const assetForm = reactive({
      id: '',
      name: '',
      ip: '',
      type: '',
      status: 'active',
      department: '',
      owner: '',
      operatingSystem: '',
      version: '',
      url: '',
      macAddress: '',
      tags: [],
      notes: ''
    })

    // 表单验证规则
    const rules = {
      name: [
        { required: true, message: '请输入资产名称', trigger: 'blur' },
        { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
      ],
      ip: [
        { required: true, message: '请输入IP地址', trigger: 'blur' },
        { pattern: /^(\d{1,3}\.){3}\d{1,3}$/, message: 'IP地址格式不正确', trigger: 'blur' }
      ],
      type: [
        { required: true, message: '请选择资产类型', trigger: 'change' }
      ],
      status: [
        { required: true, message: '请选择状态', trigger: 'change' }
      ],
      department: [
        { required: true, message: '请输入所属部门', trigger: 'blur' }
      ]
    }

    // 标签输入相关
    const tagInputRef = ref(null)
    const inputTagVisible = ref(false)
    const inputTagValue = ref('')

    // 获取资产详情
    const getAssetInfo = async () => {
      try {
        loading.value = true
        const response = await getAssetDetail(assetId)
        if (response.code === 200 && response.data) {
          const data = response.data
          assetForm.id = data.id
          assetForm.name = data.name
          assetForm.ip = data.ip
          assetForm.type = data.type
          assetForm.status = data.status
          assetForm.department = data.department
          assetForm.owner = data.owner
          assetForm.operatingSystem = data.operatingSystem
          assetForm.version = data.version
          assetForm.url = data.url
          assetForm.tags = data.tags || []
          assetForm.notes = data.notes
        } else {
          ElMessage.error('获取资产详情失败')
          router.push('/assets')
        }
      } catch (error) {
        console.error('获取资产详情失败:', error)
        ElMessage.error('获取资产详情失败')
        router.push('/assets')
      } finally {
        loading.value = false
      }
    }

    // 显示标签输入框
    const showTagInput = () => {
      inputTagVisible.value = true
      nextTick(() => {
        tagInputRef.value.focus()
      })
    }

    // 确认添加标签
    const handleInputTagConfirm = () => {
      const value = inputTagValue.value.trim()
      if (value) {
        if (!assetForm.tags.includes(value)) {
          assetForm.tags.push(value)
        }
      }
      inputTagVisible.value = false
      inputTagValue.value = ''
    }

    // 移除标签
    const removeTag = (tag) => {
      assetForm.tags = assetForm.tags.filter(t => t !== tag)
    }

    // 提交表单
    const submitForm = async () => {
      if (!assetFormRef.value) return
      
      await assetFormRef.value.validate(async (valid) => {
        if (valid) {
          try {
            loading.value = true
            const response = await updateAsset(assetId, assetForm)
            ElMessage.success('更新资产成功')
            router.push('/assets')
          } catch (error) {
            console.error('更新资产失败:', error)
            ElMessage.error('更新资产失败，请重试')
          } finally {
            loading.value = false
          }
        } else {
          return false
        }
      })
    }

    // 重置表单
    const resetForm = () => {
      // 重新获取原始数据
      getAssetInfo()
    }

    // 返回列表页
    const goBack = () => {
      router.push('/assets')
    }

    // 页面加载时获取资产详情
    onMounted(() => {
      getAssetInfo()
    })

    return {
      assetFormRef,
      assetForm,
      rules,
      loading,
      tagInputRef,
      inputTagVisible,
      inputTagValue,
      showTagInput,
      handleInputTagConfirm,
      removeTag,
      submitForm,
      resetForm,
      goBack
    }
  }
}
</script>

<style scoped>
.asset-edit-container {
  padding: 0;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-header .left {
  display: flex;
  align-items: center;
}

.page-header h2 {
  margin: 0 0 0 10px;
}

.form-card {
  margin-bottom: 20px;
}

.tag-item {
  margin-right: 10px;
  margin-bottom: 10px;
}

.tag-input {
  width: 90px;
  margin-left: 10px;
  vertical-align: bottom;
}

.button-new-tag {
  margin-left: 10px;
  height: 32px;
  padding-top: 0;
  padding-bottom: 0;
}
</style> 