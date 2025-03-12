<template>
  <div class="asset-import">
    <!-- 页面标题 -->
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">
          <span class="title-icon"><i class="el-icon-upload2"></i></span>
          资产导入
          <span class="title-highlight">工具</span>
        </h1>
        <p class="page-subtitle">批量导入资产数据，支持多种文件格式，快速完成资产管理</p>
      </div>
      <div class="header-actions">
        <el-button @click="cancel" class="action-btn">
          <i class="el-icon-back"></i> 返回列表
        </el-button>
      </div>
    </div>

    <div class="import-card">
      <el-tabs v-model="activeTab" type="border-card" class="import-tabs">
        <!-- 文件导入 -->
        <el-tab-pane label="文件导入" name="file">
          <div class="tab-content">
            <div class="import-description">
              <div class="description-header">
                <i class="el-icon-document format-icon"></i>
                <h3>支持的文件格式</h3>
              </div>
              <div class="format-grid">
                <div class="format-item">
                  <div class="format-icon-wrapper csv">
                    <i class="el-icon-document"></i>
                  </div>
                  <span>CSV文件 (.csv)</span>
                </div>
                <div class="format-item">
                  <div class="format-icon-wrapper excel">
                    <i class="el-icon-document"></i>
                  </div>
                  <span>Excel文件 (.xlsx, .xls)</span>
                </div>
                <div class="format-item">
                  <div class="format-icon-wrapper json">
                    <i class="el-icon-document"></i>
                  </div>
                  <span>JSON文件 (.json)</span>
                </div>
              </div>
              <p class="format-notice">请确保文件包含必要的字段，如资产名称、类型、状态等。
                <a href="#" @click.prevent="downloadTemplate" class="download-link">
                  <i class="el-icon-download"></i> 下载CSV模板
                </a>
                <a href="#" @click.prevent="downloadJsonTemplate" class="download-link" style="margin-left: 15px;">
                  <i class="el-icon-download"></i> 下载JSON模板
                </a>
              </p>
              <el-alert
                type="warning"
                show-icon
                :closable="false"
                class="import-alert"
              >
                <p>导入大量数据可能需要一些时间，请耐心等待并不要刷新页面。</p>
              </el-alert>
            </div>

            <el-upload
              class="upload-area"
              drag
              action=""
              :auto-upload="false"
              :http-request="handleCustomUpload"
              :on-change="handleFileChange"
              :before-upload="beforeImport"
              :file-list="fileList"
              :limit="1"
              ref="uploadRef"
            >
              <i class="el-icon-upload"></i>
              <div class="el-upload__text">将文件拖到此处，或<em>点击上传</em></div>
              <template #tip>
                <div class="el-upload__tip">支持CSV、Excel和JSON文件，文件大小不超过10MB</div>
              </template>
            </el-upload>

            <div class="import-options">
              <div class="options-header">
                <i class="el-icon-setting"></i>
                <h3>导入选项</h3>
              </div>
              <el-form :model="importOptions" label-width="180px">
                <el-form-item label="导入重复项时">
                  <el-radio-group v-model="importOptions.duplicateStrategy">
                    <el-radio label="skip">跳过重复项</el-radio>
                    <el-radio label="update">更新重复项</el-radio>
                    <el-radio label="create_new">创建为新记录</el-radio>
                  </el-radio-group>
                </el-form-item>
                
                <el-form-item label="默认资产类型">
                  <el-select v-model="importOptions.defaultType" placeholder="选择默认资产类型">
                    <el-option label="未指定" value="" />
                    <el-option label="服务器" value="server" />
                    <el-option label="工作站" value="workstation" />
                    <el-option label="网络设备" value="network" />
                    <el-option label="应用系统" value="application" />
                    <el-option label="数据库" value="database" />
                    <el-option label="云服务" value="cloud" />
                    <el-option label="容器" value="container" />
                    <el-option label="物联网设备" value="iot" />
                    <el-option label="其他" value="other" />
                  </el-select>
                  <div class="option-description">当导入数据未指定资产类型时使用</div>
                </el-form-item>
                
                <el-form-item label="默认资产状态">
                  <el-select v-model="importOptions.defaultStatus" placeholder="选择默认状态">
                    <el-option label="活跃" value="active" />
                    <el-option label="维护中" value="maintenance" />
                    <el-option label="已退役" value="decommissioned" />
                    <el-option label="保留" value="reserved" />
                    <el-option label="问题" value="issue" />
                  </el-select>
                  <div class="option-description">当导入数据未指定状态时使用</div>
                </el-form-item>
                
                <el-form-item label="发送通知">
                  <el-switch v-model="importOptions.sendNotifications" />
                  <div class="option-description">完成导入后是否发送通知</div>
                </el-form-item>
              </el-form>
            </div>
          </div>

          <div class="import-actions">
            <el-button @click="cancel" class="cancel-btn">取消</el-button>
            <el-button type="primary" @click="startImport" :loading="importing" :disabled="fileList.length === 0" class="import-btn">
              <i class="el-icon-upload"></i> 开始导入
            </el-button>
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>

    <!-- 进度对话框 -->
    <el-dialog
      title="导入进度"
      v-model="progressDialogVisible"
      :close-on-click-modal="false"
      :close-on-press-escape="false"
      :show-close="importCompleted"
      width="500px"
      class="progress-dialog"
    >
      <div class="progress-content">
        <el-progress 
          :percentage="importProgress" 
          :status="importProgress === 100 ? 'success' : ''"
          :stroke-width="20"
        ></el-progress>
        <div class="progress-message">{{ progressMessage }}</div>
        
        <div v-if="importCompleted" class="import-results">
          <h4>导入统计</h4>
          <div class="results-grid">
            <div class="result-item">
              <div class="result-value success">{{ importResults.successful }}</div>
              <div class="result-label">成功导入</div>
            </div>
            <div class="result-item">
              <div class="result-value warning">{{ importResults.duplicates }}</div>
              <div class="result-label">重复资产</div>
            </div>
            <div class="result-item">
              <div class="result-value danger">{{ importResults.failed }}</div>
              <div class="result-label">导入失败</div>
            </div>
          </div>
        </div>
      </div>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button 
            v-if="!importCompleted" 
            @click="cancelImport"
            class="cancel-import-btn"
          >
            取消导入
          </el-button>
          <el-button 
            v-else
            type="primary" 
            @click="closeProgressDialog"
            class="close-dialog-btn"
          >
            完成
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import assetApi from '@/api/assets'

export default {
  name: 'AssetImport',
  setup() {
    const router = useRouter()
    const uploadRef = ref(null)
    const activeTab = ref('file')
    const fileList = ref([])
    const importing = ref(false)
    const progressDialogVisible = ref(false)
    const importProgress = ref(0)
    const progressMessage = ref('准备导入...')
    const importCompleted = ref(false)
    
    // 导入选项
    const importOptions = reactive({
      duplicateStrategy: 'skip', // 默认跳过重复项
      defaultType: 'server',      // 默认类型
      defaultStatus: 'active',    // 默认状态
      sendNotifications: true,    // 默认发送通知
    })
    
    // 导入结果统计
    const importResults = reactive({
      successful: 0,
      duplicates: 0,
      failed: 0
    })
    
    // 取消并返回列表
    const cancel = () => {
      router.push('/assets')
    }
    
    // 下载模板文件
    const downloadTemplate = () => {
      assetApi.downloadImportTemplate()
        .then(response => {
          // 创建Blob对象
          const blob = new Blob([response.data], { 
            type: response.headers['content-type'] 
          })
          
          // 创建下载链接
          const link = document.createElement('a')
          const url = URL.createObjectURL(blob)
          link.href = url
          
          // 获取文件名
          const contentDisposition = response.headers['content-disposition']
          let filename = '资产导入模板.csv'
          if (contentDisposition) {
            const filenameMatch = contentDisposition.match(/filename=(.+)/)
            if (filenameMatch && filenameMatch[1]) {
              filename = filenameMatch[1].replace(/["']/g, '')
            }
          }
          
          link.setAttribute('download', filename)
          document.body.appendChild(link)
          link.click()
          document.body.removeChild(link)
        })
        .catch(error => {
          console.error('下载模板失败:', error)
          ElMessage.error('下载模板失败，请稍后重试')
        })
    }
    
    // 下载JSON模板文件
    const downloadJsonTemplate = () => {
      // 创建JSON模板内容
      const jsonTemplate = [
        {
          "name": "Web服务器01",
          "description": "主要生产环境Web服务器",
          "type": "server",
          "status": "active",
          "ipAddress": "192.168.1.10",
          "macAddress": "00:1A:2B:3C:4D:5E",
          "location": "北京IDC",
          "owner": "张三",
          "department": "技术部",
          "purchaseDate": "2022-01-15",
          "expiryDate": "2025-01-15",
          "os": "Linux",
          "osVersion": "Ubuntu 22.04 LTS",
          "manufacturer": "Dell",
          "model": "PowerEdge R740",
          "serialNumber": "SN12345678",
          "tags": "web,production"
        },
        {
          "name": "数据库服务器01",
          "description": "核心数据库服务器",
          "type": "database",
          "status": "active",
          "ipAddress": "192.168.1.20",
          "macAddress": "00:1A:2B:3C:4D:5F",
          "location": "上海IDC",
          "owner": "李四",
          "department": "运维部",
          "purchaseDate": "2022-03-20",
          "expiryDate": "2025-03-20",
          "os": "Linux",
          "osVersion": "CentOS 8",
          "manufacturer": "HP",
          "model": "ProLiant DL380",
          "serialNumber": "SN87654321",
          "tags": "database,production"
        }
      ];
      
      // 创建Blob对象并下载
      const blob = new Blob([JSON.stringify(jsonTemplate, null, 2)], { type: 'application/json' });
      const link = document.createElement('a');
      link.href = URL.createObjectURL(blob);
      link.download = '资产导入模板.json';
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
    };
    
    // 文件变更处理
    const handleFileChange = (file) => {
      // 限制只保留最新的文件
      fileList.value = [file]
    }
    
    // 自定义上传处理
    const handleCustomUpload = (options) => {
      // 这个函数不会被直接调用，因为我们在startImport中手动调用API
      console.log('上传选项:', options)
    }
    
    // 上传前验证
    const beforeImport = (file) => {
      // 检查文件大小
      const maxSizeInMB = 10
      const maxSize = maxSizeInMB * 1024 * 1024
      if (file.size > maxSize) {
        ElMessage.error(`文件大小不能超过 ${maxSizeInMB}MB`)
        return false
      }
      
      // 检查文件类型
      const validTypes = [
        'text/csv', 
        'application/vnd.ms-excel', 
        'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
        'application/json'
      ]
      
      const validExtensions = ['.csv', '.xls', '.xlsx', '.json']
      const fileName = file.name.toLowerCase()
      // 检查文件扩展名，用于验证文件类型
      const fileExt = fileName.substring(fileName.lastIndexOf('.'))
      
      const isValidType = validTypes.includes(file.type) || validExtensions.some(ext => fileExt === ext)
      
      if (!isValidType) {
        ElMessage.error(`只支持CSV、Excel和JSON文件，当前文件扩展名: ${fileExt}`)
        return false
      }
      
      return true
    }
    
    // 开始导入
    const startImport = () => {
      ElMessageBox.confirm(
        '确定要开始导入资产数据吗？根据数据量大小，这可能需要一些时间。',
        '开始导入',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'info'
        }
      ).then(() => {
        importing.value = true
        progressDialogVisible.value = true
        importProgress.value = 0
        progressMessage.value = '正在解析文件...'
        
        // 创建FormData对象
        const formData = new FormData()
        
        // 添加导入选项
        formData.append('options', JSON.stringify(importOptions))
        
        // 添加文件
        if (fileList.value.length > 0 && fileList.value[0].raw) {
          formData.append('file', fileList.value[0].raw)
        }
        
        // 设置进度模拟
        startProgressSimulation()
        
        // 发送导入请求
        assetApi.importAssets(formData)
          .then(response => {
            // 处理成功响应
            importing.value = false
            importProgress.value = 100
            progressMessage.value = '导入完成！'
            importCompleted.value = true
            
            // 更新导入统计信息
            if (response && response.data && response.data.stats) {
              importResults.successful = response.data.stats.successful || 0
              importResults.duplicates = response.data.stats.duplicates || 0
              importResults.failed = response.data.stats.failed || 0
            }
            
            ElMessage.success('导入成功')
          })
          .catch(error => {
            // 处理错误
            importing.value = false
            let errorMessage = '导入失败'
            
            if (error.response && error.response.data && error.response.data.error) {
              errorMessage += ': ' + error.response.data.error
            } else if (error.message) {
              errorMessage += ': ' + error.message
            }
            
            progressMessage.value = errorMessage
            importCompleted.value = true
            ElMessage.error(errorMessage)
          })
      }).catch(() => {
        // 用户取消导入
      })
    }
    
    // 模拟进度更新
    const startProgressSimulation = () => {
      let currentProgress = 0
      const interval = setInterval(() => {
        currentProgress += Math.random() * 5
        if (currentProgress >= 90) {
          currentProgress = 90
          clearInterval(interval)
        }
        
        // 只有在导入未完成时更新进度
        if (!importCompleted.value) {
          importProgress.value = Math.floor(currentProgress)
        } else {
          // 导入已完成，清除定时器
          clearInterval(interval)
        }
      }, 800)
    }
    
    // 处理导入成功
    const handleImportSuccess = (response) => {
      importing.value = false
      importProgress.value = 100
      progressMessage.value = '导入完成！'
      importCompleted.value = true
      
      // 更新导入统计信息
      if (response && response.stats) {
        importResults.successful = response.stats.successful || 0
        importResults.duplicates = response.stats.duplicates || 0
        importResults.failed = response.stats.failed || 0
      }
      
      ElMessage.success('导入成功')
    }
    
    // 处理导入错误
    const handleImportError = (error) => {
      importing.value = false
      let errorMessage = '导入失败'
      
      if (error.response && error.response.data && error.response.data.error) {
        errorMessage += ': ' + error.response.data.error
      } else if (error.message) {
        errorMessage += ': ' + error.message
      }
      
      progressMessage.value = errorMessage
      importCompleted.value = true
      ElMessage.error(errorMessage)
      
      // 关闭进度对话框
      progressDialogVisible.value = false
    }
    
    // 关闭进度对话框
    const closeProgressDialog = () => {
      progressDialogVisible.value = false
      importCompleted.value = false
      
      // 成功导入后跳转到资产列表
      if (importProgress.value === 100) {
        router.push('/assets')
      }
    }
    
    // 取消导入
    const cancelImport = () => {
      ElMessageBox.confirm(
        '确定要取消当前导入操作吗？已导入的数据将不会回滚。',
        '取消导入',
        {
          confirmButtonText: '确定',
          cancelButtonText: '继续导入',
          type: 'warning'
        }
      ).then(() => {
        importing.value = false
        progressDialogVisible.value = false
        ElMessage.info('导入已取消')
      }).catch(() => {
        // 用户选择继续导入
      })
    }
    
    return {
      activeTab,
      fileList,
      importing,
      progressDialogVisible,
      importProgress,
      progressMessage,
      importCompleted,
      importOptions,
      importResults,
      uploadRef,
      cancel,
      downloadTemplate,
      downloadJsonTemplate,
      handleFileChange,
      handleCustomUpload,
      beforeImport,
      startImport,
      handleImportSuccess,
      handleImportError,
      closeProgressDialog,
      cancelImport
    }
  }
}
</script>

<style lang="scss" scoped>
.asset-import {
  padding: 24px;
  animation: fadeIn 0.6s ease;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  
  .header-content {
    .page-title {
      font-size: 28px;
      font-weight: 600;
      margin: 0 0 8px 0;
      display: flex;
      align-items: center;
      
      .title-icon {
        margin-right: 12px;
        color: #409EFF;
      }
      
      .title-highlight {
        background: linear-gradient(120deg, #409EFF, #53A8FF);
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
        margin-left: 8px;
        font-weight: 700;
      }
    }
    
    .page-subtitle {
      color: #606266;
      font-size: 16px;
      margin: 0;
    }
  }
  
  .header-actions {
    .action-btn {
      font-weight: 500;
    }
  }
}

.import-card {
  background-color: #fff;
  border-radius: 12px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.05);
  overflow: hidden;
  
  .import-tabs {
    &:deep(.el-tabs__header) {
      margin: 0;
    }
    
    &:deep(.el-tabs__nav) {
      border: none;
    }
    
    &:deep(.el-tabs__item) {
      font-size: 16px;
      height: 60px;
      line-height: 60px;
      padding: 0 24px;
    }
    
    &:deep(.el-tabs__content) {
      padding: 0;
    }
  }
}

.tab-content {
  padding: 30px;
}

.import-description {
  margin-bottom: 30px;
  
  .description-header {
    display: flex;
    align-items: center;
    margin-bottom: 16px;
    
    .format-icon {
      font-size: 22px;
      color: #409EFF;
      margin-right: 10px;
    }
    
    h3 {
      margin: 0;
      font-size: 18px;
      font-weight: 600;
      color: #303133;
    }
  }
}

.format-grid {
  display: flex;
  gap: 20px;
  flex-wrap: wrap;
  margin-bottom: 20px;
  
  .format-item {
    display: flex;
    flex-direction: column;
    align-items: center;
    width: 120px;
    
    .format-icon-wrapper {
      width: 60px;
      height: 60px;
      border-radius: 12px;
      display: flex;
      align-items: center;
      justify-content: center;
      margin-bottom: 10px;
      
      i {
        font-size: 24px;
        color: white;
      }
      
      &.csv {
        background: linear-gradient(135deg, #4CAF50, #8BC34A);
      }
      
      &.excel {
        background: linear-gradient(135deg, #2196F3, #03A9F4);
      }
      
      &.json {
        background: linear-gradient(135deg, #FF9800, #FFC107);
      }
      
      &.xml {
        background: linear-gradient(135deg, #9C27B0, #673AB7);
      }
    }
    
    span {
      font-size: 14px;
      color: #606266;
      text-align: center;
    }
  }
}

.format-notice {
  color: #606266;
  font-size: 14px;
  line-height: 1.5;
  
  .download-link {
    color: #409EFF;
    text-decoration: none;
    margin-left: 5px;
    
    &:hover {
      text-decoration: underline;
    }
    
    i {
      margin-right: 3px;
    }
  }
}

.import-alert {
  margin-top: 16px;
}

.upload-area {
  margin: 30px 0;
  
  &:deep(.el-upload-dragger) {
    width: 100%;
    height: 200px;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    
    i {
      font-size: 48px;
      color: #c0c4cc;
      margin-bottom: 16px;
    }
    
    .el-upload__text {
      font-size: 16px;
      
      em {
        color: #409EFF;
        font-style: normal;
      }
    }
  }
  
  &:deep(.el-upload__tip) {
    text-align: center;
    margin-top: 12px;
    color: #909399;
  }
}

.import-options {
  background-color: #f8f9fa;
  border-radius: 8px;
  padding: 24px;
  
  .options-header {
    display: flex;
    align-items: center;
    margin-bottom: 20px;
    
    i {
      font-size: 20px;
      color: #409EFF;
      margin-right: 10px;
    }
    
    h3 {
      margin: 0;
      font-size: 18px;
      font-weight: 600;
      color: #303133;
    }
  }
  
  .option-description {
    font-size: 12px;
    color: #909399;
    line-height: 1.5;
    margin-top: 4px;
  }
  
  :deep(.el-form-item) {
    margin-bottom: 24px;
    
    &:last-child {
      margin-bottom: 0;
    }
  }
  
  :deep(.el-form-item__label) {
    color: #606266;
    font-weight: 500;
  }
  
  .severity-option {
    display: flex;
    align-items: center;
    
    .severity-dot {
      width: 10px;
      height: 10px;
      border-radius: 50%;
      margin-right: 8px;
      
      &.critical { background-color: #F56C6C; }
      &.high { background-color: #E6A23C; }
      &.medium { background-color: #F2D03B; }
      &.low { background-color: #67C23A; }
      &.info { background-color: #909399; }
    }
  }
}

.import-actions {
  display: flex;
  justify-content: flex-end;
  padding: 24px 30px;
  border-top: 1px solid #EBEEF5;
  
  .cancel-btn {
    margin-right: 12px;
  }
  
  .import-btn {
    min-width: 120px;
    
    i {
      margin-right: 5px;
    }
  }
}

.progress-dialog {
  &:deep(.el-dialog__body) {
    padding: 30px;
  }
  
  &:deep(.el-dialog__header) {
    padding-top: 20px;
    padding-bottom: 20px;
    border-bottom: 1px solid #EBEEF5;
  }
  
  &:deep(.el-dialog__footer) {
    padding-top: 20px;
    padding-bottom: 20px;
    border-top: 1px solid #EBEEF5;
  }
  
  .progress-content {
    .progress-message {
      margin-top: 16px;
      font-size: 16px;
      text-align: center;
      color: #606266;
    }
  }
  
  .import-results {
    margin-top: 30px;
    
    h4 {
      font-size: 16px;
      color: #303133;
      margin: 0 0 16px 0;
      text-align: center;
    }
    
    .results-grid {
      display: grid;
      grid-template-columns: repeat(3, 1fr);
      gap: 16px;
      
      .result-item {
        text-align: center;
        
        .result-value {
          font-size: 28px;
          font-weight: 700;
          line-height: 1.2;
          margin-bottom: 8px;
          
          &.success { color: #67C23A; }
          &.warning { color: #E6A23C; }
          &.danger { color: #F56C6C; }
        }
        
        .result-label {
          font-size: 14px;
          color: #606266;
        }
      }
    }
  }
}

// 针对小屏幕的响应式调整
@media screen and (max-width: 768px) {
  .page-header {
    flex-direction: column;
    align-items: flex-start;
    
    .header-actions {
      margin-top: 16px;
    }
  }
  
  .format-grid {
    gap: 10px;
    
    .format-item {
      width: 100px;
      
      .format-icon-wrapper {
        width: 50px;
        height: 50px;
      }
    }
  }
  
  .import-options {
    .option-description {
      font-size: 11px;
    }
  }
  
  .results-grid {
    grid-template-columns: 1fr !important;
  }
}
</style> 