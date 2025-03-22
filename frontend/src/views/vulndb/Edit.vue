<template>
  <div class="vulndb-edit-container">
    <div class="page-header">
      <div class="left">
        <h2>编辑漏洞</h2>
      </div>
      <div class="right">
        <el-button @click="goBack">返回</el-button>
        <el-button type="primary" @click="submitForm" :loading="submitting">保存</el-button>
      </div>
    </div>

    <div v-loading="loading">
      <el-card v-if="!loading">
        <el-form 
          ref="formRef" 
          :model="form" 
          :rules="rules" 
          label-position="top" 
          class="vulndb-form"
        >
          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item label="标题" prop="title">
                <el-input v-model="form.title" placeholder="输入漏洞标题"></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="CVE编号" prop="cve">
                <el-input v-model="form.cve" placeholder="例如: CVE-2021-44228"></el-input>
              </el-form-item>
            </el-col>
          </el-row>

          <el-form-item label="漏洞描述" prop="description">
            <el-input 
              v-model="form.description" 
              type="textarea" 
              :rows="5" 
              placeholder="输入详细的漏洞描述"
            ></el-input>
          </el-form-item>

          <el-row :gutter="20">
            <el-col :span="8">
              <el-form-item label="严重程度" prop="severity">
                <el-select v-model="form.severity" placeholder="请选择" style="width: 100%">
                  <el-option
                    v-for="item in severityOptions"
                    :key="item.value"
                    :label="item.label"
                    :value="item.value"
                  ></el-option>
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="CVSS分数" prop="cvss">
                <el-input v-model="form.cvss" placeholder="0.0-10.0"></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="CWE编号" prop="cwe">
                <el-input v-model="form.cwe" placeholder="例如: 79"></el-input>
              </el-form-item>
            </el-col>
          </el-row>

          <el-form-item label="CVSS向量" prop="cvss_vector">
            <el-input v-model="form.cvss_vector" placeholder="例如: CVSS:3.1/AV:N/AC:L/PR:L/UI:N/S:U/C:H/I:H/A:N"></el-input>
          </el-form-item>

          <el-form-item label="有漏洞利用">
            <el-switch 
              v-model="form.exploit_available" 
              active-text="是" 
              inactive-text="否"
            ></el-switch>
          </el-form-item>

          <el-form-item label="利用信息" prop="exploit_info" v-if="form.exploit_available">
            <el-input 
              v-model="form.exploit_info" 
              type="textarea" 
              :rows="3" 
              placeholder="描述漏洞利用方式、PoC或EXP的获取途径等"
            ></el-input>
          </el-form-item>

          <el-form-item label="参考链接" prop="references">
            <el-tag
              :key="index"
              v-for="(reference, index) in form.references"
              closable
              :disable-transitions="false"
              @close="handleRemoveReference(index)"
              style="margin-right: 10px; margin-bottom: 10px"
            >
              {{ reference }}
            </el-tag>
            <el-input
              v-if="inputVisible"
              ref="referenceInputRef"
              v-model="referenceInput"
              class="input-new-tag"
              size="small"
              @keyup.enter="handleAddReference"
              @blur="handleAddReference"
            ></el-input>
            <el-button v-else class="button-new-tag" size="small" @click="showInput">
              + 添加
            </el-button>
          </el-form-item>

          <el-form-item label="标签" prop="tags">
            <el-tag
              :key="index"
              v-for="(tag, index) in form.tags"
              closable
              :disable-transitions="false"
              @close="handleRemoveTag(index)"
              style="margin-right: 10px; margin-bottom: 10px"
            >
              {{ tag }}
            </el-tag>
            <el-input
              v-if="tagInputVisible"
              ref="tagInputRef"
              v-model="tagInput"
              class="input-new-tag"
              size="small"
              @keyup.enter="handleAddTag"
              @blur="handleAddTag"
            ></el-input>
            <el-button v-else class="button-new-tag" size="small" @click="showTagInput">
              + 添加
            </el-button>
          </el-form-item>

          <el-form-item label="受影响产品" prop="affected_products">
            <el-tag
              :key="index"
              v-for="(product, index) in form.affected_products"
              closable
              :disable-transitions="false"
              @close="handleRemoveProduct(index)"
              style="margin-right: 10px; margin-bottom: 10px"
            >
              {{ product }}
            </el-tag>
            <el-input
              v-if="productInputVisible"
              ref="productInputRef"
              v-model="productInput"
              class="input-new-tag"
              size="small"
              @keyup.enter="handleAddProduct"
              @blur="handleAddProduct"
            ></el-input>
            <el-button v-else class="button-new-tag" size="small" @click="showProductInput">
              + 添加
            </el-button>
          </el-form-item>

          <el-form-item label="解决方案" prop="remediation">
            <el-input 
              v-model="form.remediation" 
              type="textarea" 
              :rows="5" 
              placeholder="输入漏洞修复建议或解决方案"
            ></el-input>
          </el-form-item>

          <el-form-item label="发布日期" prop="published_date">
            <el-date-picker
              v-model="form.published_date"
              type="date"
              placeholder="选择日期"
              style="width: 100%"
            ></el-date-picker>
          </el-form-item>
        </el-form>
      </el-card>
    </div>
  </div>
</template>

<script>
import { ref, reactive, nextTick, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getVulnDBDetail, updateVulnDBEntry, getVulnSeverities } from '@/api/vulndb'

export default {
  name: 'VulnDBEdit',
  setup() {
    const route = useRoute()
    const router = useRouter()
    const formRef = ref(null)
    const loading = ref(true)
    const submitting = ref(false)
    const severityOptions = getVulnSeverities()

    // 参考链接输入相关
    const referenceInputRef = ref(null)
    const inputVisible = ref(false)
    const referenceInput = ref('')

    // 标签输入相关
    const tagInputRef = ref(null)
    const tagInputVisible = ref(false)
    const tagInput = ref('')

    // 受影响产品输入相关
    const productInputRef = ref(null)
    const productInputVisible = ref(false)
    const productInput = ref('')

    // 表单数据
    const form = reactive({
      id: '',
      title: '',
      description: '',
      severity: 'medium',
      cvss: '',
      cvss_vector: '',
      cve: '',
      cwe: '',
      exploit_available: false,
      exploit_info: '',
      references: [],
      tags: [],
      affected_products: [],
      remediation: '',
      published_date: new Date()
    })

    // 表单验证规则
    const rules = reactive({
      title: [
        { required: true, message: '请输入漏洞标题', trigger: 'blur' },
        { min: 5, max: 200, message: '标题长度应在5到200个字符之间', trigger: 'blur' }
      ],
      description: [
        { required: true, message: '请输入漏洞描述', trigger: 'blur' }
      ],
      severity: [
        { required: true, message: '请选择严重程度', trigger: 'change' }
      ]
    })

    // 获取漏洞详情
    const getVulnDetail = async () => {
      const id = route.params.id
      if (!id) {
        ElMessage.error('无效的漏洞ID')
        loading.value = false
        router.push('/vulndb')
        return
      }

      loading.value = true
      try {
        const response = await getVulnDBDetail(id)
        if (response.code === 200 && response.data) {
          const data = response.data
          
          // 填充表单数据
          form.id = data.id
          form.title = data.title
          form.description = data.description
          form.severity = data.severity
          form.cvss = data.cvss
          form.cvss_vector = data.cvss_vector
          form.cve = data.cve
          form.cwe = data.cwe
          form.exploit_available = data.exploit_available
          form.exploit_info = data.exploit_info
          form.references = data.references || []
          form.tags = data.tags || []
          form.affected_products = data.affected_products || []
          form.remediation = data.remediation
          
          if (data.published_date) {
            form.published_date = new Date(data.published_date)
          }
        } else {
          ElMessage.error(response.message || '获取漏洞详情失败')
          router.push('/vulndb')
        }
      } catch (error) {
        console.error('获取漏洞详情失败:', error)
        ElMessage.error('获取漏洞详情失败')
        router.push('/vulndb')
      } finally {
        loading.value = false
      }
    }

    // 参考链接相关方法
    const showInput = () => {
      inputVisible.value = true
      nextTick(() => {
        referenceInputRef.value.focus()
      })
    }

    const handleAddReference = () => {
      if (referenceInput.value) {
        if (form.references.indexOf(referenceInput.value) === -1) {
          form.references.push(referenceInput.value)
        }
      }
      inputVisible.value = false
      referenceInput.value = ''
    }

    const handleRemoveReference = (index) => {
      form.references.splice(index, 1)
    }

    // 标签相关方法
    const showTagInput = () => {
      tagInputVisible.value = true
      nextTick(() => {
        tagInputRef.value.focus()
      })
    }

    const handleAddTag = () => {
      if (tagInput.value) {
        if (form.tags.indexOf(tagInput.value) === -1) {
          form.tags.push(tagInput.value)
        }
      }
      tagInputVisible.value = false
      tagInput.value = ''
    }

    const handleRemoveTag = (index) => {
      form.tags.splice(index, 1)
    }

    // 受影响产品相关方法
    const showProductInput = () => {
      productInputVisible.value = true
      nextTick(() => {
        productInputRef.value.focus()
      })
    }

    const handleAddProduct = () => {
      if (productInput.value) {
        if (form.affected_products.indexOf(productInput.value) === -1) {
          form.affected_products.push(productInput.value)
        }
      }
      productInputVisible.value = false
      productInput.value = ''
    }

    const handleRemoveProduct = (index) => {
      form.affected_products.splice(index, 1)
    }

    // 提交表单
    const submitForm = async () => {
      if (!formRef.value) return
      
      await formRef.value.validate(async (valid) => {
        if (!valid) {
          return false
        }
        
        submitting.value = true
        try {
          const response = await updateVulnDBEntry(form.id, form)
          if (response.code === 200) {
            ElMessage.success('漏洞更新成功')
            router.push('/vulndb')
          } else {
            ElMessage.error(response.message || '更新失败')
          }
        } catch (error) {
          console.error('更新漏洞失败:', error)
          ElMessage.error('更新失败')
        } finally {
          submitting.value = false
        }
      })
    }

    // 返回列表
    const goBack = () => {
      router.push('/vulndb')
    }

    onMounted(() => {
      getVulnDetail()
    })

    return {
      formRef,
      form,
      rules,
      loading,
      submitting,
      severityOptions,
      referenceInputRef,
      inputVisible,
      referenceInput,
      tagInputRef,
      tagInputVisible,
      tagInput,
      productInputRef,
      productInputVisible,
      productInput,
      showInput,
      handleAddReference,
      handleRemoveReference,
      showTagInput,
      handleAddTag,
      handleRemoveTag,
      showProductInput,
      handleAddProduct,
      handleRemoveProduct,
      submitForm,
      goBack
    }
  }
}
</script>

<style scoped>
.vulndb-edit-container {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-header h2 {
  margin: 0;
  font-size: 22px;
  color: #303133;
}

.vulndb-form {
  padding: 20px 0;
}

.input-new-tag {
  width: 100px;
  margin-right: 10px;
  vertical-align: bottom;
}

.button-new-tag {
  margin-right: 10px;
  height: 32px;
  line-height: 30px;
  padding-top: 0;
  padding-bottom: 0;
}
</style> 