<template>
  <div class="asset-index-container">
    <div class="asset-header">
      <h2 class="page-title">资产管理</h2>
      <div class="header-actions">
        <el-input
          v-model="queryParams.keyword"
          placeholder="输入关键词搜索"
          class="search-input"
          clearable
          @clear="handleSearch"
          @keyup.enter="handleSearch"
        >
          <template #append>
            <el-button @click="handleSearch">
              <el-icon><Search /></el-icon>
            </el-button>
          </template>
        </el-input>

        <el-button type="primary" @click="handleAdd" class="header-btn">
          <el-icon><Plus /></el-icon>添加资产
        </el-button>

        <el-button type="success" @click="handleImportClick" class="header-btn">
          <el-icon><Upload /></el-icon>批量导入
        </el-button>

        <el-button type="info" @click="handleExport" class="header-btn">
          <el-icon><Download /></el-icon>导出
        </el-button>
        
        <el-button 
          type="danger" 
          @click="handleBatchDelete" 
          class="header-btn"
          :disabled="selectedAssets.length === 0"
        >
          <el-icon><Delete /></el-icon>批量删除
        </el-button>
      </div>
    </div>

    <!-- 数据表格 -->
    <el-card class="asset-table-card" shadow="hover">
      <el-table
        v-loading="loading"
        :data="assetList"
        border
        style="width: 100%"
        @selection-change="handleSelectionChange"
        :highlight-current-row="true"
        stripe
        class="asset-table"
      >
        <el-table-column type="selection" width="50" />
        <el-table-column type="index" label="序号" width="50" />
        <el-table-column prop="name" label="资产名称" min-width="120" show-overflow-tooltip />
        <el-table-column prop="ip" label="IP地址" width="120" show-overflow-tooltip />
        <el-table-column prop="type" label="资产类型" width="100">
          <template #default="scope">
            <el-tag :type="getAssetTypeTag(scope.row.type)" effect="light" size="small">{{ scope.row.type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="80">
          <template #default="scope">
            <el-tag :type="scope.row.status === 'active' ? 'success' : 'danger'" effect="dark" size="small">
              {{ scope.row.status === 'active' ? '在线' : '离线' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="department" label="所属部门" width="100" show-overflow-tooltip />
        <el-table-column prop="createdAt" label="创建时间" width="150" show-overflow-tooltip>
          <template #default="scope">
            {{ formatTime(scope.row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column fixed="right" label="操作" width="180">
          <template #default="scope">
            <div class="operation-buttons">
              <el-dropdown trigger="click">
                <el-button type="primary" size="small" class="operation-dropdown-btn">
                  <el-icon><More /></el-icon>
                  <span>操作</span>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item @click="handleView(scope.row)" class="dropdown-item">
                      <el-icon><View /></el-icon>查看
                    </el-dropdown-item>
                    <el-dropdown-item @click="handleEdit(scope.row)" class="dropdown-item">
                      <el-icon><Edit /></el-icon>编辑
                    </el-dropdown-item>
                    <el-dropdown-item @click="handleDelete(scope.row)" class="dropdown-item danger-item">
                      <el-icon><Delete /></el-icon>删除
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          :current-page="queryParams.page"
          :page-size="queryParams.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          :total="total"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
          background
        />
      </div>
    </el-card>

    <!-- 导入模板下载提示对话框 -->
    <el-dialog
      v-model="importDialogVisible"
      title="批量导入资产"
      width="500px"
    >
      <p>请下载模板文件，按照模板格式填写资产信息后上传。</p>
      <div class="import-template">
        <el-button-group>
          <el-button type="primary" @click="downloadTemplate('csv')">下载CSV模板</el-button>
          <el-button type="primary" @click="downloadTemplate('json')">下载JSON模板</el-button>
        </el-button-group>
      </div>
      <el-divider>或直接上传</el-divider>
      <el-upload
        class="upload-area"
        drag
        action=""
        :http-request="handleBatchImport"
        :show-file-list="false"
        accept=".xlsx,.csv,.json"
      >
        <el-icon class="el-icon--upload"><upload-filled /></el-icon>
        <div class="el-upload__text">
          将文件拖到此处，或<em>点击上传</em>
        </div>
        <template #tip>
          <div class="el-upload__tip">
            支持 xlsx、csv、json 格式文件，文件大小不超过 2MB
          </div>
        </template>
      </el-upload>
    </el-dialog>
  </div>
</template>

<script>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  Search, Plus, Upload, Download, UploadFilled,
  View, Edit, Delete, More
} from '@element-plus/icons-vue'
import { 
  getAssetList, deleteAsset, batchImportAssets, exportAssets, batchDeleteAssets 
} from '@/api/asset'
import { formatTime } from '@/utils/time'

export default {
  name: 'AssetIndex',
  components: {
    Search, Plus, Upload, Download, UploadFilled,
    View, Edit, Delete, More
  },
  setup() {
    const router = useRouter()
    const loading = ref(false)
    const assetList = ref([])
    const total = ref(0)
    const selectedAssets = ref([])
    const importDialogVisible = ref(false)

    // 查询参数
    const queryParams = reactive({
      page: 1,
      pageSize: 10,
      keyword: '',
      status: '',
      type: ''
    })

    // 获取资产列表
    const getList = async () => {
      try {
        loading.value = true
        const response = await getAssetList(queryParams)
        assetList.value = response.data.items || []
        total.value = response.data.total || 0
      } catch (error) {
        console.error('获取资产列表失败:', error)
        ElMessage.error('获取资产列表失败')
      } finally {
        loading.value = false
      }
    }

    // 获取资产类型对应的标签类型
    const getAssetTypeTag = (type) => {
      const typeMap = {
        'server': 'danger',
        'website': 'success',
        'database': 'warning',
        'network': 'info',
        'application': 'primary'
      }
      return typeMap[type] || ''
    }

    // 处理搜索
    const handleSearch = () => {
      queryParams.page = 1
      getList()
    }

    // 添加资产
    const handleAdd = () => {
      router.push('/assets/add')
    }

    // 查看资产详情
    const handleView = (row) => {
      router.push(`/assets/detail/${row.id}`)
    }

    // 编辑资产
    const handleEdit = (row) => {
      router.push(`/assets/edit/${row.id}`)
    }

    // 删除资产
    const handleDelete = (row) => {
      ElMessageBox.confirm(
        `确定要删除资产 "${row.name}" 吗？`,
        '警告',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }
      ).then(async () => {
        try {
          await deleteAsset(row.id)
          ElMessage.success('删除成功')
          getList()
        } catch (error) {
          console.error('删除资产失败:', error)
          ElMessage.error('删除资产失败')
        }
      }).catch(() => {
        // 取消删除
      })
    }

    // 多选框选中数据
    const handleSelectionChange = (selection) => {
      selectedAssets.value = selection
    }

    // 批量删除
    const handleBatchDelete = () => {
      if (selectedAssets.value.length === 0) {
        ElMessage.warning('请先选择要删除的资产')
        return
      }

      const assetNames = selectedAssets.value.map(item => item.name).join('、')
      const confirmMessage = selectedAssets.value.length > 3 
        ? `确定要删除选中的 ${selectedAssets.value.length} 个资产吗？包括 ${assetNames.substring(0, 30)}... 等`
        : `确定要删除资产 ${assetNames} 吗？`

      ElMessageBox.confirm(
        confirmMessage,
        '批量删除确认',
        {
          confirmButtonText: '确定删除',
          cancelButtonText: '取消',
          type: 'warning'
        }
      ).then(async () => {
        try {
          const ids = selectedAssets.value.map(item => item.id)
          await batchDeleteAssets(ids)
          ElMessage.success(`成功删除 ${ids.length} 个资产`)
          // 刷新列表
          getList()
        } catch (error) {
          console.error('批量删除资产失败:', error)
          ElMessage.error('批量删除资产失败')
        }
      }).catch(() => {
        // 取消删除
      })
    }

    // 分页大小改变
    const handleSizeChange = (size) => {
      queryParams.pageSize = size
      getList()
    }

    // 页码改变
    const handleCurrentChange = (page) => {
      queryParams.page = page
      getList()
    }

    // 显示批量导入对话框
    const handleImportClick = () => {
      importDialogVisible.value = true
    }

    // 下载导入模板
    const downloadTemplate = (templateType) => {
      // 这里应该调用后端接口下载模板，这里简单示例使用前端生成
      const link = document.createElement('a')
      link.href = `/templates/asset_import_template.${templateType}` // 假设有这个模板文件
      link.setAttribute('download', `资产导入模板.${templateType}`)
      
      // 添加错误处理
      link.onerror = () => {
        ElMessage.error(`下载模板失败，请联系管理员`)
      }
      
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
    }

    // 处理批量导入
    const handleBatchImport = async (options) => {
      try {
        const formData = new FormData()
        formData.append('file', options.file)
        
        loading.value = true
        const response = await batchImportAssets(formData)
        
        if (response.success) {
          ElMessage.success(`成功导入 ${response.data.successCount || 0} 条数据`)
          importDialogVisible.value = false
          getList()
        } else {
          ElMessage.warning(response.message || '部分数据导入失败，请检查数据格式')
        }
      } catch (error) {
        console.error('批量导入失败:', error)
        ElMessage.error('批量导入失败，请检查文件格式')
      } finally {
        loading.value = false
      }
    }

    // 导出资产数据
    const handleExport = async () => {
      try {
        loading.value = true
        const response = await exportAssets(queryParams)
        
        // 创建Blob对象并下载
        const blob = new Blob([response], { type: 'application/vnd.ms-excel' })
        const link = document.createElement('a')
        link.href = URL.createObjectURL(blob)
        link.setAttribute('download', `资产列表_${new Date().getTime()}.xlsx`)
        document.body.appendChild(link)
        link.click()
        document.body.removeChild(link)
        
        ElMessage.success('导出成功')
      } catch (error) {
        console.error('导出失败:', error)
        ElMessage.error('导出失败')
      } finally {
        loading.value = false
      }
    }

    // 页面加载时获取数据
    onMounted(() => {
      getList()
    })

    return {
      loading,
      assetList,
      total,
      selectedAssets,
      importDialogVisible,
      queryParams,
      getAssetTypeTag,
      handleSearch,
      handleAdd,
      handleView,
      handleEdit,
      handleDelete,
      handleSelectionChange,
      handleSizeChange,
      handleCurrentChange,
      handleImportClick,
      downloadTemplate,
      handleBatchImport,
      handleExport,
      handleBatchDelete,
      formatTime
    }
  }
}
</script>

<style scoped>
.asset-index-container {
  padding: 0;
}

.page-title {
  margin: 0;
  font-size: 22px;
  color: #303133;
  position: relative;
  font-weight: 600;
}

.asset-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  flex-wrap: wrap;
  padding: 10px 0;
}

.header-actions {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
  margin-top: 10px;
}

.search-input {
  width: 300px;
  transition: all 0.3s;
}

.search-input:focus-within {
  box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.2);
}

.header-btn {
  transition: all 0.3s;
}

.header-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.asset-table-card {
  margin-bottom: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.05);
}

.asset-table {
  margin-bottom: 15px;
}

/* 表格行悬停效果 */
.asset-table :deep(.el-table__row:hover) {
  background-color: #f5f7fa !important;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.import-template {
  margin: 20px 0;
  text-align: center;
}

.upload-area {
  display: flex;
  justify-content: center;
  transition: all 0.3s;
}

.upload-area:hover {
  transform: translateY(-2px);
}

.operation-buttons {
  display: flex;
  justify-content: center;
  gap: 5px;
}

.operation-dropdown-btn {
  transition: all 0.3s;
  border-radius: 4px;
  min-width: 90px;
}

.operation-dropdown-btn:hover {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.operation-buttons .el-icon {
  margin-right: 4px;
}

/* 下拉菜单样式 */
:deep(.el-dropdown-menu__item) {
  display: flex;
  align-items: center;
  gap: 5px;
  padding: 8px 16px;
  transition: all 0.2s;
}

:deep(.el-dropdown-menu__item:hover) {
  background-color: #ecf5ff;
}

:deep(.danger-item:hover) {
  background-color: #fef0f0;
  color: #f56c6c;
}

:deep(.el-tag) {
  border-radius: 4px;
  padding: 0 8px;
  font-weight: 500;
  transition: all 0.2s;
}

:deep(.el-tag:hover) {
  transform: scale(1.05);
}

@media screen and (max-width: 768px) {
  .asset-header {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .header-actions {
    width: 100%;
    justify-content: space-between;
  }
  
  .search-input {
    width: 100%;
    margin-bottom: 10px;
  }
}
</style> 