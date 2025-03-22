<template>
  <div class="knowledge-detail-container" v-loading="loading">
    <el-card v-if="knowledgeDetail.id">
      <div class="knowledge-header">
        <h1 class="title">{{ knowledgeDetail.title }}</h1>
        <div class="meta">
          <el-tag>{{ knowledgeDetail.type }}</el-tag>
          <span class="author">作者: {{ knowledgeDetail.author }}</span>
          <span class="date">发布时间: {{ formatDate(knowledgeDetail.created_at) }}</span>
          <span class="views">浏览次数: {{ knowledgeDetail.view_count }}</span>
        </div>
        <div class="tags-categories" v-if="knowledgeDetail.tags || knowledgeDetail.categories">
          <div v-if="knowledgeDetail.categories" class="categories">
            分类: {{ knowledgeDetail.categories }}
          </div>
          <div v-if="knowledgeDetail.tags" class="tags">
            标签: {{ knowledgeDetail.tags }}
          </div>
        </div>
        <div class="actions">
          <el-button type="primary" @click="handleEdit">编辑</el-button>
          <el-button @click="goBack">返回</el-button>
        </div>
      </div>
      
      <el-divider />
      
      <div class="knowledge-content">
        <pre>{{ knowledgeDetail.content }}</pre>
      </div>
      
      <el-divider v-if="knowledgeDetail.references || knowledgeDetail.attachments" />
      
      <div v-if="knowledgeDetail.references" class="references">
        <h3>参考资料</h3>
        <pre>{{ knowledgeDetail.references }}</pre>
      </div>
      
      <div v-if="knowledgeDetail.attachments" class="attachments">
        <h3>附件</h3>
        <pre>{{ knowledgeDetail.attachments }}</pre>
      </div>
      
      <div v-if="knowledgeDetail.related_vuln_types" class="related-vuln-types">
        <h3>相关漏洞类型</h3>
        <div>{{ knowledgeDetail.related_vuln_types }}</div>
      </div>
    </el-card>
    
    <div v-else-if="!loading" class="not-found">
      <el-empty description="未找到知识条目"></el-empty>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getKnowledgeDetail } from '@/api/knowledge'
import { formatDate } from '@/utils/format'

export default {
  name: 'KnowledgeDetail',
  setup() {
    const route = useRoute()
    const router = useRouter()
    const loading = ref(false)
    const knowledgeDetail = ref({})

    // 获取知识详情
    const fetchDetail = async (id) => {
      loading.value = true
      try {
        const response = await getKnowledgeDetail(id)
        if (response.code === 200 && response.data) {
          knowledgeDetail.value = response.data
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

    // 编辑知识条目
    const handleEdit = () => {
      router.push(`/knowledge/edit/${route.params.id}`)
    }

    // 返回列表
    const goBack = () => {
      router.push('/knowledge')
    }

    onMounted(() => {
      const id = route.params.id
      if (id) {
        fetchDetail(id)
      }
    })

    return {
      loading,
      knowledgeDetail,
      formatDate,
      handleEdit,
      goBack
    }
  }
}
</script>

<style scoped>
.knowledge-detail-container {
  padding: 20px;
}

.knowledge-header {
  margin-bottom: 20px;
}

.title {
  font-size: 24px;
  font-weight: bold;
  margin-bottom: 15px;
  color: #303133;
}

.meta {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  margin-bottom: 10px;
  color: #606266;
  font-size: 14px;
}

.meta > * {
  margin-right: 15px;
  margin-bottom: 5px;
}

.tags-categories {
  display: flex;
  flex-wrap: wrap;
  margin-bottom: 15px;
  font-size: 14px;
  color: #606266;
}

.tags-categories > * {
  margin-right: 20px;
  margin-bottom: 5px;
}

.actions {
  margin-top: 15px;
}

.knowledge-content {
  margin: 20px 0;
}

.knowledge-content pre,
.references pre,
.attachments pre {
  white-space: pre-wrap;
  word-wrap: break-word;
  background-color: #f8f8f8;
  padding: 15px;
  border-radius: 4px;
  font-family: inherit;
  line-height: 1.6;
  color: #303133;
  margin: 0;
}

h3 {
  font-size: 18px;
  color: #303133;
  margin: 20px 0 10px 0;
}

.references,
.attachments,
.related-vuln-types {
  margin-top: 15px;
}

.not-found {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 300px;
}
</style> 