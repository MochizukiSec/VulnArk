<template>
  <div class="asset-list">
    <!-- 页面标题 -->
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">
          <span class="title-icon"><i class="el-icon-office-building"></i></span>
          资产列表
          <span class="title-highlight">管理</span>
        </h1>
        <p class="page-subtitle">查看和管理所有资产信息，监控资产安全状态</p>
      </div>
      <div class="header-actions">
        <el-button type="primary" @click="goToCreate" class="action-btn">
          <i class="el-icon-plus"></i> 添加资产
        </el-button>
      </div>
    </div>

    <!-- 过滤和搜索区域 -->
    <div class="filter-card">
      <el-form :inline="true" :model="searchParams" class="search-form">
        <el-form-item label="类型">
          <el-select v-model="searchParams.type" placeholder="选择资产类型" clearable>
            <el-option label="服务器" value="server" />
            <el-option label="网络设备" value="network_device" />
            <el-option label="应用程序" value="application" />
            <el-option label="数据库" value="database" />
            <el-option label="云资源" value="cloud_resource" />
            <el-option label="物联网设备" value="iot_device" />
            <el-option label="其他" value="other" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="环境">
          <el-select v-model="searchParams.environment" placeholder="选择环境" clearable>
            <el-option label="生产" value="production" />
            <el-option label="测试" value="testing" />
            <el-option label="开发" value="development" />
            <el-option label="预生产" value="staging" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="搜索">
          <el-input 
            v-model="searchParams.searchTerm" 
            placeholder="搜索名称、描述或IP地址"
            clearable
            @keyup.enter="fetchAssets"
          >
            <template #append>
              <el-button icon="el-icon-search" @click="fetchAssets"></el-button>
            </template>
          </el-input>
        </el-form-item>
        
        <el-form-item>
          <el-button type="primary" @click="fetchAssets">筛选</el-button>
          <el-button @click="resetFilters">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 资产列表表格 -->
    <el-card class="list-card" v-loading="loading">
      <el-table 
        :data="assets" 
        style="width: 100%"
        border
        @row-click="handleRowClick"
        :empty-text="loading ? '加载中...' : '未找到符合条件的资产'"
      >
        <el-table-column label="ID" prop="id" width="80" />
        <el-table-column label="名称" prop="name" min-width="150">
          <template #default="scope">
            <el-tooltip :content="scope.row.description" placement="top">
              <router-link :to="{ name: 'AssetDetail', params: { id: scope.row.id } }" class="asset-link">
                {{ scope.row.name }}
              </router-link>
            </el-tooltip>
          </template>
        </el-table-column>
        <el-table-column label="类型" prop="type" width="120">
          <template #default="scope">
            <el-tag :type="getAssetTypeTag(scope.row.type)">
              {{ getAssetTypeLabel(scope.row.type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="IP地址" prop="ipAddress" width="150" />
        <el-table-column label="环境" prop="environment" width="100">
          <template #default="scope">
            <el-tag :type="getEnvironmentTag(scope.row.environment)">
              {{ getEnvironmentLabel(scope.row.environment) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="漏洞数量" prop="vulnerabilityCount" width="100" align="center">
          <template #default="scope">
            <el-badge :value="scope.row.vulnerabilityCount || 0" type="danger" />
          </template>
        </el-table-column>
        <el-table-column label="创建时间" prop="createdAt" width="180">
          <template #default="scope">
            {{ formatDate(scope.row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="scope">
            <el-button 
              size="small" 
              type="primary" 
              plain
              @click.stop="goToDetails(scope.row.id)"
            >
              详情
            </el-button>
            <el-button 
              size="small" 
              type="warning" 
              plain
              @click.stop="goToEdit(scope.row.id)"
            >
              编辑
            </el-button>
            <el-button 
              size="small" 
              type="danger" 
              plain
              @click.stop="confirmDelete(scope.row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:currentPage="currentPage"
          :page-sizes="[10, 20, 50, 100]"
          :page-size="pageSize"
          layout="total, sizes, prev, pager, next, jumper"
          :total="totalAssets"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>
  </div>
</template>

<script>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessageBox, ElMessage } from 'element-plus'
import assetApi from '@/api/assets'
import { formatDate } from '@/utils/format'

export default {
  name: 'AssetList',
  
  setup() {
    const router = useRouter()
    const loading = ref(false)
    const assets = ref([])
    const currentPage = ref(1)
    const pageSize = ref(10)
    const totalAssets = ref(0)
    
    // 搜索参数
    const searchParams = reactive({
      type: '',
      environment: '',
      searchTerm: '',
    })
    
    // 获取资产列表
    const fetchAssets = async () => {
      loading.value = true
      try {
        const params = {
          ...searchParams,
          limit: pageSize.value,
          skip: (currentPage.value - 1) * pageSize.value
        }
        
        const response = await assetApi.getAssets(params)
        assets.value = response.data.assets || []
        totalAssets.value = response.data.total || 0
      } catch (error) {
        console.error('获取资产列表失败:', error)
        ElMessage.error('获取资产列表失败，请重试')
      } finally {
        loading.value = false
      }
    }
    
    // 分页处理
    const handleSizeChange = (val) => {
      pageSize.value = val
      fetchAssets()
    }
    
    const handleCurrentChange = (val) => {
      currentPage.value = val
      fetchAssets()
    }
    
    // 重置筛选条件
    const resetFilters = () => {
      Object.keys(searchParams).forEach(key => {
        searchParams[key] = ''
      })
      fetchAssets()
    }
    
    // 点击行跳转到详情
    const handleRowClick = (row) => {
      router.push({ name: 'AssetDetail', params: { id: row.id } })
    }
    
    // 跳转到创建页面
    const goToCreate = () => {
      router.push({ name: 'AssetCreate' })
    }
    
    // 跳转到详情页面
    const goToDetails = (id) => {
      router.push({ name: 'AssetDetail', params: { id } })
    }
    
    // 跳转到编辑页面
    const goToEdit = (id) => {
      router.push({ name: 'AssetEdit', params: { id } })
    }
    
    // 确认删除资产
    const confirmDelete = (asset) => {
      ElMessageBox.confirm(
        `确定要删除资产 "${asset.name}" 吗？此操作不可恢复。`,
        '警告',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning',
        }
      ).then(() => {
        deleteAsset(asset.id)
      }).catch(() => {
        // 用户取消删除
      })
    }
    
    // 删除资产
    const deleteAsset = async (id) => {
      try {
        await assetApi.deleteAsset(id)
        ElMessage.success('资产已成功删除')
        fetchAssets() // 刷新列表
      } catch (error) {
        console.error('删除资产失败:', error)
        ElMessage.error('删除资产失败，请重试')
      }
    }
    
    // 获取资产类型显示标签
    const getAssetTypeLabel = (type) => {
      const types = {
        server: '服务器',
        network_device: '网络设备',
        application: '应用程序',
        database: '数据库',
        cloud_resource: '云资源',
        iot_device: '物联网设备',
        other: '其他'
      }
      return types[type] || type
    }
    
    // 获取资产类型对应的标签样式
    const getAssetTypeTag = (type) => {
      const tags = {
        server: 'primary',
        network_device: 'success',
        application: 'warning',
        database: 'danger',
        cloud_resource: 'info',
        iot_device: '',
        other: 'info'
      }
      return tags[type] || 'info'
    }
    
    // 获取环境显示标签
    const getEnvironmentLabel = (env) => {
      const environments = {
        production: '生产',
        testing: '测试',
        development: '开发',
        staging: '预生产'
      }
      return environments[env] || env
    }
    
    // 获取环境对应的标签样式
    const getEnvironmentTag = (env) => {
      const tags = {
        production: 'danger',
        testing: 'warning',
        development: 'success',
        staging: 'info'
      }
      return tags[env] || 'info'
    }
    
    onMounted(() => {
      fetchAssets()
    })
    
    return {
      loading,
      assets,
      currentPage,
      pageSize,
      totalAssets,
      searchParams,
      fetchAssets,
      resetFilters,
      handleSizeChange,
      handleCurrentChange,
      handleRowClick,
      goToCreate,
      goToDetails,
      goToEdit,
      confirmDelete,
      getAssetTypeLabel,
      getAssetTypeTag,
      getEnvironmentLabel,
      getEnvironmentTag,
      formatDate
    }
  }
}
</script>

<style scoped>
.asset-list {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 15px;
  border-bottom: 1px solid #ebeef5;
}

.page-title {
  font-size: 24px;
  color: #303133;
  margin: 0 0 8px 0;
  display: flex;
  align-items: center;
}

.title-icon {
  margin-right: 10px;
  font-size: 24px;
  color: #409eff;
}

.title-highlight {
  color: #409eff;
  margin-left: 5px;
}

.page-subtitle {
  color: #909399;
  font-size: 14px;
  margin: 0;
}

.filter-card {
  background: #f5f7fa;
  padding: 20px;
  border-radius: 4px;
  margin-bottom: 20px;
}

.search-form {
  display: flex;
  flex-wrap: wrap;
}

.list-card {
  margin-bottom: 20px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}

.asset-link {
  color: #409eff;
  text-decoration: none;
  font-weight: bold;
}

.asset-link:hover {
  text-decoration: underline;
}
</style> 