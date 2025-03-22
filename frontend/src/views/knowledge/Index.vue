<template>
  <div class="knowledge-container">
    <div class="page-header">
      <div class="left">
        <h2>知识库管理</h2>
      </div>
      <div class="right">
        <el-button type="primary" @click="handleAddKnowledge">添加知识</el-button>
      </div>
    </div>

    <!-- 搜索过滤 -->
    <el-card class="filter-container">
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="标题或内容"></el-input>
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="searchForm.type" placeholder="所有类型" clearable>
            <el-option
              v-for="item in typeOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            ></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="分类">
          <el-select v-model="searchForm.category" placeholder="所有分类" clearable>
            <el-option
              v-for="category in categories"
              :key="category"
              :label="category"
              :value="category"
            ></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="标签">
          <el-input v-model="searchForm.tags" placeholder="标签"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 知识库列表 -->
    <el-card>
      <el-table :data="knowledgeList" v-loading="loading" border style="width: 100%">
        <el-table-column prop="title" label="标题" min-width="200">
          <template #default="scope">
            <router-link :to="`/knowledge/${scope.row.id}`" class="link-type">
              {{ scope.row.title }}
            </router-link>
          </template>
        </el-table-column>
        <el-table-column prop="type" label="类型" width="120">
          <template #default="scope">
            <el-tag>{{ getTypeLabel(scope.row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="author" label="作者" width="120"></el-table-column>
        <el-table-column prop="view_count" label="浏览次数" width="100"></el-table-column>
        <el-table-column prop="categories" label="分类" width="150"></el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="160"></el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="scope">
            <div class="operation-buttons">
              <el-button
                type="primary"
                size="small"
                @click="handleEdit(scope.row)"
                class="action-btn"
              >
                <el-icon><Edit /></el-icon>
                <span>编辑</span>
              </el-button>
              <el-button
                type="danger"
                size="small"
                @click="handleDelete(scope.row)"
                class="action-btn"
              >
                <el-icon><Delete /></el-icon>
                <span>删除</span>
              </el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          background
          layout="total, sizes, prev, pager, next, jumper"
          :total="total"
          :page-size="pageSize"
          :current-page="currentPage"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        ></el-pagination>
      </div>
    </el-card>
  </div>
</template>

<script>
import { reactive, ref, onMounted, toRefs } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessageBox, ElMessage } from 'element-plus'
import { Edit, Delete } from '@element-plus/icons-vue'
import { getKnowledgeList, deleteKnowledge, getKnowledgeTypes, getKnowledgeCategories } from '@/api/knowledge'

export default {
  name: 'KnowledgeList',
  components: {
    Edit,
    Delete
  },
  setup() {
    const router = useRouter()
    const loading = ref(false)
    const total = ref(0)
    const pageSize = ref(10)
    const currentPage = ref(1)
    const typeOptions = ref([])
    const categories = ref([])

    const state = reactive({
      knowledgeList: [],
      searchForm: {
        keyword: '',
        type: '',
        category: '',
        tags: ''
      }
    })

    // 获取知识库列表
    const getList = async () => {
      loading.value = true
      try {
        const params = {
          page: currentPage.value,
          page_size: pageSize.value,
          keyword: state.searchForm.keyword,
          type: state.searchForm.type,
          category: state.searchForm.category,
          tags: state.searchForm.tags
        }
        
        const response = await getKnowledgeList(params)
        if (response.code === 200 && response.data) {
          state.knowledgeList = response.data.items
          total.value = response.data.total
        } else {
          ElMessage.error('获取知识库列表失败')
        }
      } catch (error) {
        console.error('获取知识库列表失败:', error)
        ElMessage.error('获取知识库列表失败')
      } finally {
        loading.value = false
      }
    }

    // 获取类型选项
    const fetchTypes = async () => {
      try {
        const response = await getKnowledgeTypes()
        if (response.code === 200 && response.data) {
          typeOptions.value = response.data
        }
      } catch (error) {
        console.error('获取知识库类型失败:', error)
      }
    }

    // 获取分类
    const fetchCategories = async () => {
      try {
        const response = await getKnowledgeCategories()
        if (response.code === 200 && response.data) {
          categories.value = response.data
        }
      } catch (error) {
        console.error('获取知识库分类失败:', error)
      }
    }

    // 获取类型标签
    const getTypeLabel = (type) => {
      const option = typeOptions.value.find(item => item.value === type)
      return option ? option.label : type
    }

    // 添加知识条目
    const handleAddKnowledge = () => {
      router.push('/knowledge/add')
    }

    // 编辑知识条目
    const handleEdit = (row) => {
      router.push(`/knowledge/edit/${row.id}`)
    }

    // 删除知识条目
    const handleDelete = (row) => {
      ElMessageBox.confirm(`确定要删除 "${row.title}" 吗？`, '警告', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async () => {
        try {
          const response = await deleteKnowledge(row.id)
          if (response.code === 200) {
            ElMessage.success(`成功删除 "${row.title}"`)
            getList() // 刷新列表
          } else {
            ElMessage.error(response.message || '删除失败')
          }
        } catch (error) {
          console.error('删除知识条目失败:', error)
          ElMessage.error('删除失败')
        }
      }).catch(() => {})
    }

    // 搜索
    const handleSearch = () => {
      currentPage.value = 1
      getList()
    }

    // 重置
    const handleReset = () => {
      state.searchForm.keyword = ''
      state.searchForm.type = ''
      state.searchForm.category = ''
      state.searchForm.tags = ''
      handleSearch()
    }

    // 处理分页大小变化
    const handleSizeChange = (size) => {
      pageSize.value = size
      getList()
    }

    // 处理页码变化
    const handleCurrentChange = (page) => {
      currentPage.value = page
      getList()
    }

    onMounted(() => {
      getList()
      fetchTypes()
      fetchCategories()
    })

    return {
      ...toRefs(state),
      loading,
      total,
      pageSize,
      currentPage,
      typeOptions,
      categories,
      getTypeLabel,
      handleAddKnowledge,
      handleEdit,
      handleDelete,
      handleSearch,
      handleReset,
      handleSizeChange,
      handleCurrentChange
    }
  }
}
</script>

<style scoped>
.knowledge-container {
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

.filter-container {
  margin-bottom: 20px;
}

.search-form {
  display: flex;
  flex-wrap: wrap;
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

.operation-buttons {
  display: flex;
  justify-content: space-between;
  gap: 8px;
}

.action-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 6px 10px;
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