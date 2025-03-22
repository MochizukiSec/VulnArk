<template>
  <div class="knowledge-edit-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ isEdit ? '编辑知识' : '添加知识' }}</span>
        </div>
      </template>
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px" v-loading="loading">
        <el-form-item label="标题" prop="title">
          <el-input v-model="form.title" placeholder="请输入标题"></el-input>
        </el-form-item>
        <el-form-item label="类型" prop="type">
          <el-select v-model="form.type" placeholder="请选择类型">
            <el-option
              v-for="item in typeOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            ></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="作者" prop="author">
          <el-input v-model="form.author" placeholder="请输入作者"></el-input>
        </el-form-item>
        <el-form-item label="分类" prop="categories">
          <el-input v-model="form.categories" placeholder="请输入分类，多个分类用逗号分隔"></el-input>
        </el-form-item>
        <el-form-item label="标签" prop="tags">
          <el-input v-model="form.tags" placeholder="请输入标签，多个标签用逗号分隔"></el-input>
        </el-form-item>
        <el-form-item label="相关漏洞类型" prop="related_vuln_types">
          <el-input v-model="form.related_vuln_types" placeholder="请输入相关的漏洞类型，多个类型用逗号分隔"></el-input>
        </el-form-item>
        <el-form-item label="内容" prop="content">
          <div class="editor-container">
            <div class="editor-tabs">
              <el-tabs v-model="activeTab">
                <el-tab-pane label="编辑" name="edit">
                  <el-input
                    v-model="form.content"
                    type="textarea"
                    :rows="15"
                    placeholder="支持Markdown格式"
                  ></el-input>
                </el-tab-pane>
                <el-tab-pane label="预览" name="preview">
                  <div class="preview-content" v-html="formattedContent"></div>
                </el-tab-pane>
              </el-tabs>
            </div>
          </div>
        </el-form-item>
        <el-form-item label="参考资料" prop="references">
          <el-input
            v-model="form.references"
            type="textarea"
            :rows="4"
            placeholder="请输入参考资料链接，每行一个"
          ></el-input>
        </el-form-item>
        <el-form-item label="附件" prop="attachments">
          <el-input
            v-model="form.attachments"
            type="textarea"
            :rows="4"
            placeholder="请输入附件链接，每行一个"
          ></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="submitForm" :loading="saving">保存</el-button>
          <el-button @click="goBack">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script>
import { reactive, ref, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getKnowledgeDetail, createKnowledge, updateKnowledge, getKnowledgeTypes } from '@/api/knowledge'
import marked from 'marked'
import DOMPurify from 'dompurify'

export default {
  name: 'KnowledgeEdit',
  setup() {
    const router = useRouter()
    const route = useRoute()
    const formRef = ref(null)
    const loading = ref(false)
    const saving = ref(false)
    const activeTab = ref('edit')
    const typeOptions = ref([])
    
    // 判断是编辑还是新增
    const isEdit = computed(() => route.params.id !== undefined)
    
    // 表单数据
    const form = reactive({
      title: '',
      type: '',
      author: '',
      categories: '',
      tags: '',
      content: '',
      references: '',
      attachments: '',
      related_vuln_types: ''
    })
    
    // 表单验证规则
    const rules = {
      title: [
        { required: true, message: '请输入标题', trigger: 'blur' },
        { min: 2, max: 100, message: '长度在2到100个字符之间', trigger: 'blur' }
      ],
      type: [
        { required: true, message: '请选择类型', trigger: 'change' }
      ],
      content: [
        { required: true, message: '请输入内容', trigger: 'blur' }
      ]
    }
    
    // 获取知识类型选项
    const fetchTypeOptions = async () => {
      try {
        const response = await getKnowledgeTypes()
        if (response.code === 200 && response.data) {
          typeOptions.value = response.data
        }
      } catch (error) {
        console.error('获取知识类型失败:', error)
        ElMessage.error('获取知识类型失败')
      }
    }
    
    // 如果是编辑模式，获取详情数据
    const fetchDetail = async (id) => {
      loading.value = true
      try {
        const response = await getKnowledgeDetail(id)
        if (response.code === 200 && response.data) {
          const data = response.data
          Object.keys(form).forEach(key => {
            if (data[key] !== undefined) {
              form[key] = data[key]
            }
          })
        } else {
          ElMessage.error('获取知识详情失败')
        }
      } catch (error) {
        console.error('获取知识详情失败:', error)
        ElMessage.error('获取知识详情失败')
      } finally {
        loading.value = false
      }
    }
    
    // 预览内容格式化（Markdown转HTML）
    const formattedContent = computed(() => {
      if (!form.content) return ''
      const html = marked(form.content)
      return DOMPurify.sanitize(html)
    })
    
    // 提交表单
    const submitForm = async () => {
      if (!formRef.value) return
      
      await formRef.value.validate(async (valid) => {
        if (!valid) return
        
        saving.value = true
        try {
          const data = { ...form }
          let response
          
          if (isEdit.value) {
            response = await updateKnowledge(route.params.id, data)
            ElMessage.success('更新成功')
          } else {
            response = await createKnowledge(data)
            ElMessage.success('创建成功')
          }
          
          if (response.code === 200) {
            router.push('/knowledge')
          }
        } catch (error) {
          console.error(isEdit.value ? '更新知识失败:' : '创建知识失败:', error)
          ElMessage.error(isEdit.value ? '更新失败' : '创建失败')
        } finally {
          saving.value = false
        }
      })
    }
    
    // 返回
    const goBack = () => {
      router.push('/knowledge')
    }
    
    onMounted(() => {
      fetchTypeOptions()
      
      if (isEdit.value) {
        fetchDetail(route.params.id)
      }
    })
    
    return {
      formRef,
      form,
      rules,
      loading,
      saving,
      activeTab,
      typeOptions,
      isEdit,
      formattedContent,
      submitForm,
      goBack
    }
  }
}
</script>

<style scoped>
.knowledge-edit-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header span {
  font-size: 18px;
  font-weight: bold;
}

.editor-container {
  border: 1px solid #dcdfe6;
  border-radius: 4px;
}

.editor-tabs {
  padding: 10px;
}

.preview-content {
  padding: 10px;
  min-height: 300px;
  max-height: 500px;
  overflow-y: auto;
  line-height: 1.6;
  color: #303133;
  background-color: #f8f8f8;
  border-radius: 4px;
}

:deep(.el-form-item__label) {
  font-weight: 500;
}
</style> 