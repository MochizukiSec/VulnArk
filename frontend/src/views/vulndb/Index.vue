<template>
  <div class="vulndb-container">
    <div class="page-header">
      <div class="left">
        <h2>漏洞库</h2>
      </div>
      <div class="right">
        <el-button type="primary" @click="handleAddVulnDB">添加漏洞</el-button>
        <el-button type="success" @click="handleImport">批量导入</el-button>
      </div>
    </div>

    <!-- 搜索过滤 -->
    <el-card class="filter-container">
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="标题或描述"></el-input>
        </el-form-item>
        <el-form-item label="CVE编号">
          <el-input v-model="searchForm.cve" placeholder="例如: CVE-2021-44228"></el-input>
        </el-form-item>
        <el-form-item label="CWE">
          <el-input v-model="searchForm.cwe" placeholder="CWE"></el-input>
        </el-form-item>
        <el-form-item label="严重程度">
          <el-select v-model="searchForm.severity" placeholder="所有级别" clearable>
            <el-option
              v-for="item in severityOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            ></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="利用代码">
          <el-select v-model="searchForm.hasExploit" placeholder="所有状态" clearable>
            <el-option label="有" value="true"></el-option>
            <el-option label="无" value="false"></el-option>
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

    <!-- 漏洞库列表 -->
    <el-card>
      <el-table :data="vulnDBList" v-loading="loading" border style="width: 100%">
        <el-table-column prop="cve" label="CVE编号" min-width="130">
          <template #default="scope">
            <a :href="'https://cve.mitre.org/cgi-bin/cvename.cgi?name=' + scope.row.cve" target="_blank" class="link-type" v-if="scope.row.cve">
              {{ scope.row.cve }}
            </a>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="title" label="标题" min-width="200">
          <template #default="scope">
            <router-link :to="`/vulndb/id/${scope.row.id}`" class="link-type">
              {{ scope.row.title }}
            </router-link>
          </template>
        </el-table-column>
        <el-table-column prop="severity" label="严重程度" width="100">
          <template #default="scope">
            <el-tag :type="getSeverityType(scope.row.severity)">{{ getSeverityLabel(scope.row.severity) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="cvss" label="CVSS" width="90"></el-table-column>
        <el-table-column prop="cwe" label="CWE" width="100">
          <template #default="scope">
            <a :href="'https://cwe.mitre.org/data/definitions/' + scope.row.cwe + '.html'" target="_blank" class="link-type" v-if="scope.row.cwe">
              {{ scope.row.cwe }}
            </a>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="exploit_available" label="利用代码" width="80">
          <template #default="scope">
            <el-tag type="danger" v-if="scope.row.exploit_available">有</el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="published_date" label="发布日期" width="130">
          <template #default="scope">
            {{ formatDate(scope.row.published_date) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="scope">
            <el-button type="primary" link @click="handleEdit(scope.row)">编辑</el-button>
            <el-button type="danger" link @click="handleDelete(scope.row)">删除</el-button>
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

    <!-- 批量导入对话框 -->
    <el-dialog title="导入漏洞库" v-model="importVisible" width="500px">
      <el-upload
        class="upload-container"
        drag
        action="#"
        :auto-upload="false"
        :file-list="fileList"
        :on-change="handleFileChange"
        :limit="1"
      >
        <el-icon class="el-icon--upload"><upload-filled /></el-icon>
        <div class="el-upload__text">
          将文件拖到此处，或<em>点击上传</em>
        </div>
        <template #tip>
          <div class="el-upload__tip">
            支持 .json 或 .csv 格式文件，最大10MB
          </div>
        </template>
      </el-upload>
      <div class="template-links">
        <p>下载模板:</p>
        <el-link type="primary" href="/templates/vulndb_template.json" target="_blank">JSON模板</el-link>
        <el-link type="primary" href="/templates/vulndb_template.csv" target="_blank">CSV模板</el-link>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="importVisible = false">取消</el-button>
          <el-button type="primary" @click="submitImport" :loading="importing">导入</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { reactive, ref, onMounted, toRefs } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessageBox, ElMessage } from 'element-plus'
import { getVulnDBList, deleteVulnDBEntry, batchImportVulnDB, getVulnSeverities } from '@/api/vulndb'
import { UploadFilled } from '@element-plus/icons-vue'

export default {
  name: 'VulnDBIndex',
  components: {
    UploadFilled
  },
  setup() {
    const router = useRouter()
    const loading = ref(false)
    const total = ref(0)
    const pageSize = ref(10)
    const currentPage = ref(1)
    const importVisible = ref(false)
    const importing = ref(false)
    const fileList = ref([])
    const severityOptions = getVulnSeverities()

    const state = reactive({
      vulnDBList: [],
      searchForm: {
        keyword: '',
        cve: '',
        cwe: '',
        severity: '',
        hasExploit: '',
        tags: ''
      }
    })

    // 获取漏洞库列表
    const getList = async () => {
      loading.value = true
      try {
        const params = {
          page: currentPage.value,
          page_size: pageSize.value,
          keyword: state.searchForm.keyword,
          cve: state.searchForm.cve,
          cwe: state.searchForm.cwe,
          severity: state.searchForm.severity,
          has_exploit: state.searchForm.hasExploit,
          tags: state.searchForm.tags
        }
        
        const response = await getVulnDBList(params)
        if (response.code === 200 && response.data) {
          state.vulnDBList = response.data.items
          total.value = response.data.total
        } else {
          ElMessage.error('获取漏洞库列表失败')
        }
      } catch (error) {
        console.error('获取漏洞库列表失败:', error)
        ElMessage.error('获取漏洞库列表失败')
      } finally {
        loading.value = false
      }
    }

    // 获取严重程度标签样式
    const getSeverityType = (severity) => {
      const typeMap = {
        critical: 'danger',
        high: 'danger',
        medium: 'warning',
        low: 'success',
        info: 'info'
      }
      return typeMap[severity] || 'info'
    }

    // 获取严重程度标签文本
    const getSeverityLabel = (severity) => {
      const labelMap = {
        critical: '严重',
        high: '高危',
        medium: '中危',
        low: '低危',
        info: '信息'
      }
      return labelMap[severity] || severity
    }

    // 格式化日期
    const formatDate = (dateStr) => {
      if (!dateStr) return '-'
      const date = new Date(dateStr)
      return date.toLocaleDateString('zh-CN')
    }

    // 添加漏洞库条目
    const handleAddVulnDB = () => {
      router.push('/vulndb/add')
    }

    // 编辑漏洞库条目
    const handleEdit = (row) => {
      router.push(`/vulndb/edit/${row.id}`)
    }

    // 删除漏洞库条目
    const handleDelete = (row) => {
      ElMessageBox.confirm(`确定要删除漏洞 "${row.title}" 吗？`, '警告', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async () => {
        try {
          const response = await deleteVulnDBEntry(row.id)
          if (response.code === 200) {
            ElMessage.success('删除成功')
            getList() // 刷新列表
          } else {
            ElMessage.error(response.message || '删除失败')
          }
        } catch (error) {
          console.error('删除漏洞库条目失败:', error)
          ElMessage.error('删除失败')
        }
      }).catch(() => {})
    }

    // 搜索
    const handleSearch = () => {
      currentPage.value = 1
      getList()
    }

    // 重置搜索
    const handleReset = () => {
      state.searchForm.keyword = ''
      state.searchForm.cve = ''
      state.searchForm.cwe = ''
      state.searchForm.severity = ''
      state.searchForm.hasExploit = ''
      state.searchForm.tags = ''
      handleSearch()
    }

    // 处理文件变更
    const handleFileChange = (file) => {
      fileList.value = [file]
    }

    // 处理导入
    const handleImport = () => {
      importVisible.value = true
      fileList.value = []
    }

    // 提交导入
    const submitImport = async () => {
      if (fileList.value.length === 0) {
        ElMessage.warning('请先选择文件')
        return
      }

      importing.value = true
      try {
        const formData = new FormData()
        formData.append('file', fileList.value[0].raw)
        
        const response = await batchImportVulnDB(formData)
        if (response.code === 200) {
          ElMessage.success(`成功导入${response.data.imported || 0}条记录`)
          importVisible.value = false
          getList() // 刷新列表
        } else {
          ElMessage.error(response.message || '导入失败')
        }
      } catch (error) {
        console.error('导入漏洞库条目失败:', error)
        ElMessage.error('导入失败')
      } finally {
        importing.value = false
      }
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
    })

    return {
      ...toRefs(state),
      loading,
      total,
      pageSize,
      currentPage,
      importVisible,
      importing,
      fileList,
      severityOptions,
      getSeverityType,
      getSeverityLabel,
      formatDate,
      handleAddVulnDB,
      handleEdit,
      handleDelete,
      handleSearch,
      handleReset,
      handleImport,
      handleFileChange,
      submitImport,
      handleSizeChange,
      handleCurrentChange
    }
  }
}
</script>

<style scoped>
.vulndb-container {
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

.upload-container {
  margin-bottom: 20px;
}

.template-links {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 15px;
  margin-top: 10px;
}

.template-links p {
  margin: 0;
}
</style>