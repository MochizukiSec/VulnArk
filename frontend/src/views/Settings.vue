<template>
  <div class="settings-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>系统设置</span>
        </div>
      </template>
      
      <el-tabs v-model="activeTab" type="border-card">
        <!-- 集成配置 -->
        <el-tab-pane label="集成配置" name="integrations">
          <div class="tab-content">
            <!-- JIRA 集成 -->
            <el-divider content-position="left">JIRA集成</el-divider>
            <el-form :model="integrationForm.jira" label-width="120px">
              <el-form-item label="启用">
                <el-switch v-model="integrationForm.jira.enabled" />
              </el-form-item>
              <el-form-item label="URL" v-show="integrationForm.jira.enabled">
                <el-input v-model="integrationForm.jira.url" placeholder="例如: https://your-domain.atlassian.net" />
              </el-form-item>
              <el-form-item label="API令牌" v-show="integrationForm.jira.enabled">
                <el-input v-model="integrationForm.jira.apiToken" placeholder="API令牌" show-password />
              </el-form-item>
              <el-form-item label="用户名" v-show="integrationForm.jira.enabled">
                <el-input v-model="integrationForm.jira.username" placeholder="用户名" />
              </el-form-item>
              <el-form-item label="默认项目" v-show="integrationForm.jira.enabled">
                <el-input v-model="integrationForm.jira.defaultProject" placeholder="默认项目" />
              </el-form-item>
              <el-form-item v-show="integrationForm.jira.enabled">
                <el-button type="primary" @click="testJiraConnection" :loading="jiraTestLoading">测试连接</el-button>
              </el-form-item>
            </el-form>

            <!-- 微信扫码登录 -->
            <el-divider content-position="left">微信扫码登录</el-divider>
            <el-form :model="integrationForm.wechat" label-width="120px">
              <el-form-item label="启用">
                <el-switch v-model="integrationForm.wechat.enabled" />
              </el-form-item>
              <el-form-item label="AppID" v-show="integrationForm.wechat.enabled">
                <el-input v-model="integrationForm.wechat.appId" placeholder="AppID" />
              </el-form-item>
              <el-form-item label="AppSecret" v-show="integrationForm.wechat.enabled">
                <el-input v-model="integrationForm.wechat.appSecret" placeholder="AppSecret" show-password />
              </el-form-item>
              <el-form-item label="授权范围" v-show="integrationForm.wechat.enabled">
                <el-select v-model="integrationForm.wechat.scope" placeholder="授权范围">
                  <el-option label="基础信息(snsapi_base)" value="snsapi_base" />
                  <el-option label="用户信息(snsapi_userinfo)" value="snsapi_userinfo" />
                </el-select>
              </el-form-item>
              <el-form-item label="回调URL" v-show="integrationForm.wechat.enabled">
                <el-input v-model="integrationForm.wechat.callbackUrl" placeholder="回调URL" />
              </el-form-item>
              <el-form-item label="重定向域名" v-show="integrationForm.wechat.enabled">
                <el-input v-model="integrationForm.wechat.redirectDomain" placeholder="重定向域名" />
              </el-form-item>
              <el-form-item v-show="integrationForm.wechat.enabled">
                <el-button type="primary" @click="testWechatLogin" :loading="wechatTestLoading">测试连接</el-button>
              </el-form-item>
            </el-form>

            <!-- 漏洞库API -->
            <el-divider content-position="left">漏洞库API</el-divider>
            <el-form :model="integrationForm.vulnDb" label-width="120px">
              <el-form-item label="启用">
                <el-switch v-model="integrationForm.vulnDb.enabled" />
              </el-form-item>
              <el-form-item label="服务提供商" v-show="integrationForm.vulnDb.enabled">
                <el-select v-model="integrationForm.vulnDb.provider" placeholder="服务提供商">
                  <el-option label="微步社区" value="weibu" />
                  <el-option label="VulnIQ" value="vulniq" />
                  <el-option label="VulDB" value="vuldb" />
                  <el-option label="其他" value="other" />
                </el-select>
              </el-form-item>
              <el-form-item label="API URL" v-show="integrationForm.vulnDb.enabled">
                <el-input v-model="integrationForm.vulnDb.apiUrl" placeholder="API URL" />
              </el-form-item>
              <el-form-item label="API Key" v-show="integrationForm.vulnDb.enabled">
                <el-input v-model="integrationForm.vulnDb.apiKey" placeholder="API Key" show-password />
              </el-form-item>
              <el-form-item label="API Secret" v-show="integrationForm.vulnDb.enabled">
                <el-input v-model="integrationForm.vulnDb.apiSecret" placeholder="API Secret" show-password />
              </el-form-item>
              <el-form-item label="额外参数" v-show="integrationForm.vulnDb.enabled">
                <el-input 
                  v-model="integrationForm.vulnDb.parameters" 
                  type="textarea" 
                  placeholder="额外参数"
                  :autosize="{ minRows: 2, maxRows: 4 }"
                />
              </el-form-item>
              <el-form-item v-show="integrationForm.vulnDb.enabled">
                <el-button type="primary" @click="testVulnDBConnection" :loading="vulndbTestLoading">测试连接</el-button>
              </el-form-item>
            </el-form>
          </div>
        </el-tab-pane>

        <!-- 消息通知 -->
        <el-tab-pane label="消息通知" name="notifications">
          <div class="tab-content">
            <!-- 企业微信机器人 -->
            <el-divider content-position="left">企业微信机器人</el-divider>
            <el-form :model="notificationForm.workWechat" label-width="120px">
              <el-form-item label="启用">
                <el-switch v-model="notificationForm.workWechat.enabled" />
              </el-form-item>
              <el-form-item label="Webhook URL" v-show="notificationForm.workWechat.enabled">
                <el-input v-model="notificationForm.workWechat.webhookUrl" placeholder="Webhook URL" />
              </el-form-item>
              <el-form-item label="事件" v-show="notificationForm.workWechat.enabled">
                <el-checkbox-group v-model="notificationForm.workWechat.events">
                  <el-checkbox label="资产创建"></el-checkbox>
                  <el-checkbox label="资产更新"></el-checkbox>
                  <el-checkbox label="资产删除"></el-checkbox>
                  <el-checkbox label="漏洞创建"></el-checkbox>
                  <el-checkbox label="漏洞状态变化"></el-checkbox>
                  <el-checkbox label="漏洞更新"></el-checkbox>
                  <el-checkbox label="漏洞删除"></el-checkbox>
                </el-checkbox-group>
              </el-form-item>
              <el-form-item v-show="notificationForm.workWechat.enabled">
                <el-button type="primary" @click="testWorkWechatNotification" :loading="workWechatTestLoading">测试连接</el-button>
              </el-form-item>
            </el-form>

            <!-- 飞书机器人 -->
            <el-divider content-position="left">飞书机器人</el-divider>
            <el-form :model="notificationForm.feishu" label-width="120px">
              <el-form-item label="启用">
                <el-switch v-model="notificationForm.feishu.enabled" />
              </el-form-item>
              <el-form-item label="Webhook URL" v-show="notificationForm.feishu.enabled">
                <el-input v-model="notificationForm.feishu.webhookUrl" placeholder="Webhook URL" />
              </el-form-item>
              <el-form-item label="Secret" v-show="notificationForm.feishu.enabled">
                <el-input v-model="notificationForm.feishu.secret" placeholder="Secret" show-password />
              </el-form-item>
              <el-form-item label="事件" v-show="notificationForm.feishu.enabled">
                <el-checkbox-group v-model="notificationForm.feishu.events">
                  <el-checkbox label="资产创建"></el-checkbox>
                  <el-checkbox label="资产更新"></el-checkbox>
                  <el-checkbox label="资产删除"></el-checkbox>
                  <el-checkbox label="漏洞创建"></el-checkbox>
                  <el-checkbox label="漏洞状态变化"></el-checkbox>
                  <el-checkbox label="漏洞更新"></el-checkbox>
                  <el-checkbox label="漏洞删除"></el-checkbox>
                </el-checkbox-group>
              </el-form-item>
              <el-form-item v-show="notificationForm.feishu.enabled">
                <el-button type="primary" @click="testFeishuNotification" :loading="feishuTestLoading">测试连接</el-button>
              </el-form-item>
            </el-form>

            <!-- 钉钉机器人 -->
            <el-divider content-position="left">钉钉机器人</el-divider>
            <el-form :model="notificationForm.dingtalk" label-width="120px">
              <el-form-item label="启用">
                <el-switch v-model="notificationForm.dingtalk.enabled" />
              </el-form-item>
              <el-form-item label="Webhook URL" v-show="notificationForm.dingtalk.enabled">
                <el-input v-model="notificationForm.dingtalk.webhookUrl" placeholder="Webhook URL" />
              </el-form-item>
              <el-form-item label="Secret" v-show="notificationForm.dingtalk.enabled">
                <el-input v-model="notificationForm.dingtalk.secret" placeholder="Secret" show-password />
              </el-form-item>
              <el-form-item label="事件" v-show="notificationForm.dingtalk.enabled">
                <el-checkbox-group v-model="notificationForm.dingtalk.events">
                  <el-checkbox label="资产创建"></el-checkbox>
                  <el-checkbox label="资产更新"></el-checkbox>
                  <el-checkbox label="资产删除"></el-checkbox>
                  <el-checkbox label="漏洞创建"></el-checkbox>
                  <el-checkbox label="漏洞状态变化"></el-checkbox>
                  <el-checkbox label="漏洞更新"></el-checkbox>
                  <el-checkbox label="漏洞删除"></el-checkbox>
                </el-checkbox-group>
              </el-form-item>
              <el-form-item v-show="notificationForm.dingtalk.enabled">
                <el-button type="primary" @click="testDingtalkNotification" :loading="dingtalkTestLoading">测试连接</el-button>
              </el-form-item>
            </el-form>

            <!-- 邮件通知 -->
            <el-divider content-position="left">邮件通知</el-divider>
            <el-form :model="notificationForm.email" label-width="120px">
              <el-form-item label="启用">
                <el-switch v-model="notificationForm.email.enabled" />
              </el-form-item>
              <el-form-item label="SMTP服务器" v-show="notificationForm.email.enabled">
                <el-input v-model="notificationForm.email.smtpServer" placeholder="SMTP服务器" />
              </el-form-item>
              <el-form-item label="SMTP端口" v-show="notificationForm.email.enabled">
                <el-input-number v-model="notificationForm.email.smtpPort" :min="1" :max="65535" />
              </el-form-item>
              <el-form-item label="发件人邮箱" v-show="notificationForm.email.enabled">
                <el-input v-model="notificationForm.email.fromEmail" placeholder="发件人邮箱" />
              </el-form-item>
              <el-form-item label="用户名" v-show="notificationForm.email.enabled">
                <el-input v-model="notificationForm.email.username" placeholder="用户名" />
              </el-form-item>
              <el-form-item label="密码" v-show="notificationForm.email.enabled">
                <el-input v-model="notificationForm.email.password" placeholder="密码" show-password />
              </el-form-item>
              <el-form-item label="使用SSL" v-show="notificationForm.email.enabled">
                <el-switch v-model="notificationForm.email.useSsl" />
              </el-form-item>
              <el-form-item label="事件" v-show="notificationForm.email.enabled">
                <el-checkbox-group v-model="notificationForm.email.events">
                  <el-checkbox label="资产创建"></el-checkbox>
                  <el-checkbox label="资产更新"></el-checkbox>
                  <el-checkbox label="资产删除"></el-checkbox>
                  <el-checkbox label="漏洞创建"></el-checkbox>
                  <el-checkbox label="漏洞状态变化"></el-checkbox>
                  <el-checkbox label="漏洞更新"></el-checkbox>
                  <el-checkbox label="漏洞删除"></el-checkbox>
                </el-checkbox-group>
              </el-form-item>
              <el-form-item label="收件人" v-show="notificationForm.email.enabled">
                <el-tag
                  v-for="email in notificationForm.email.recipients"
                  :key="email"
                  closable
                  @close="removeRecipient(email)"
                  class="recipient-tag"
                >
                  {{ email }}
                </el-tag>
                <el-input
                  v-if="inputVisible"
                  ref="emailInput"
                  v-model="inputValue"
                  @keyup.enter="addRecipient"
                  @blur="addRecipient"
                  class="input-new-tag"
                />
                <el-button v-else class="button-new-tag" @click="showInput">添加收件人</el-button>
              </el-form-item>
              <el-form-item v-show="notificationForm.email.enabled">
                <el-button type="primary" @click="testEmailNotification" :loading="emailTestLoading">测试连接</el-button>
              </el-form-item>
            </el-form>
          </div>
        </el-tab-pane>

        <!-- AI功能 -->
        <el-tab-pane label="AI功能" name="ai">
          <div class="tab-content">
            <el-divider content-position="left">AI功能</el-divider>
            <el-form :model="aiForm" label-width="120px">
              <el-form-item label="启用">
                <el-switch v-model="aiForm.enabled" />
              </el-form-item>
              <el-form-item label="提供商" v-show="aiForm.enabled">
                <el-select v-model="aiForm.provider">
                  <el-option label="OpenAI" value="openai" />
                  <el-option label="百度文心" value="baidu" />
                  <el-option label="通义千问" value="aliyun" />
                  <el-option label="自定义服务" value="custom" />
                </el-select>
              </el-form-item>
              <el-form-item label="API Key" v-show="aiForm.enabled">
                <el-input v-model="aiForm.apiKey" placeholder="API Key" show-password />
              </el-form-item>
              <el-form-item label="API URL" v-show="aiForm.enabled && aiForm.provider === 'custom'">
                <el-input v-model="aiForm.apiUrl" placeholder="API URL" />
              </el-form-item>
              <el-form-item label="分析选项" v-show="aiForm.enabled">
                <el-checkbox-group v-model="aiForm.analysisOptions">
                  <el-checkbox label="威胁分析"></el-checkbox>
                  <el-checkbox label="漏洞修复"></el-checkbox>
                  <el-checkbox label="风险评估"></el-checkbox>
                  <el-checkbox label="漏洞检测"></el-checkbox>
                </el-checkbox-group>
              </el-form-item>
            </el-form>
          </div>
        </el-tab-pane>

        <!-- 扫描器配置 -->
        <el-tab-pane label="扫描器配置" name="scanners">
          <div class="tab-content">
            <el-divider content-position="left">扫描器配置</el-divider>
            <div class="scanners-list">
              <div class="scanners-header">
                <h3>已配置的扫描器</h3>
                <el-button type="primary" size="small" @click="showAddScannerDialog" class="action-btn">
                  <el-icon><Plus /></el-icon>
                  <span>添加扫描器</span>
                </el-button>
              </div>
              
              <el-table :data="scannersList" border style="width: 100%">
                <el-table-column prop="name" label="名称" min-width="150"></el-table-column>
                <el-table-column prop="type" label="类型" width="120">
                  <template #default="scope">
                    <el-tag>{{ getScannerTypeLabel(scope.row.type) }}</el-tag>
                  </template>
                </el-table-column>
                <el-table-column label="URL/路径" min-width="180">
                  <template #default="scope">
                    <div v-if="scope.row.type === 'xray'">
                      <span v-if="scope.row.binary_path" class="binary-path">
                        <el-tag size="small" type="info">二进制: {{ scope.row.binary_path }}</el-tag>
                      </span>
                      <span v-if="scope.row.url" class="config-path">
                        <el-tag size="small" type="info">配置: {{ scope.row.url }}</el-tag>
                      </span>
                      <span v-if="!scope.row.url && !scope.row.binary_path" class="no-data">
                        -
                      </span>
                    </div>
                    <span v-else>{{ scope.row.url }}</span>
                  </template>
                </el-table-column>
                <el-table-column prop="enabled" label="状态" width="100">
                  <template #default="scope">
                    <el-tag :type="scope.row.enabled ? 'success' : 'info'">
                      {{ scope.row.enabled ? '启用' : '禁用' }}
                    </el-tag>
                  </template>
                </el-table-column>
                <el-table-column label="操作" width="200">
                  <template #default="scope">
                    <div class="operation-buttons">
                      <el-button
                        type="primary"
                        size="small"
                        @click="editScanner(scope.row)"
                        class="action-btn"
                      >
                        <el-icon><Edit /></el-icon>
                        <span>编辑</span>
                      </el-button>
                      <el-button
                        type="danger"
                        size="small"
                        @click="deleteScanner(scope.row)"
                        class="action-btn"
                      >
                        <el-icon><Delete /></el-icon>
                        <span>删除</span>
                      </el-button>
                    </div>
                  </template>
                </el-table-column>
              </el-table>
              
              <div class="empty-tip" v-if="scannersList.length === 0">
                <el-empty description="暂无配置的扫描器"></el-empty>
              </div>
            </div>
          </div>
        </el-tab-pane>
      </el-tabs>
      
      <div class="form-actions">
        <el-button type="primary" @click="saveSettings" :loading="saveLoading">保存</el-button>
        <el-button @click="resetSettings">重置</el-button>
        <el-button type="warning" @click="testNotification">测试漏洞通知</el-button>
      </div>
    </el-card>

    <!-- 添加/编辑扫描器对话框 -->
    <el-dialog
      :title="scannerForm.isEdit ? '编辑扫描器' : '添加扫描器'"
      v-model="scannerDialogVisible"
      width="600px"
    >
      <el-form 
        :model="scannerForm" 
        :rules="scannerRules" 
        ref="scannerFormRef" 
        label-width="100px"
      >
        <el-form-item label="名称" prop="name">
          <el-input v-model="scannerForm.name" placeholder="名称" />
        </el-form-item>
        
        <el-form-item label="类型" prop="type">
          <el-select v-model="scannerForm.type" placeholder="类型">
            <el-option label="Nessus" value="nessus" />
            <el-option label="AWVS" value="awvs" />
            <el-option label="Xray" value="xray" />
            <el-option label="ZAP" value="zap" />
            <el-option label="自定义" value="custom" />
          </el-select>
        </el-form-item>
        
        <el-form-item :label="scannerForm.type === 'xray' ? '配置文件路径' : '服务URL'" prop="url">
          <el-input 
            v-model="scannerForm.url" 
            :placeholder="scannerForm.type === 'xray' ? 'Xray配置文件路径（可选）' : '例如: https://scanner.example.com'" 
          />
          <div class="form-tip" v-if="scannerForm.type === 'xray'">
            如果留空，将使用Xray默认配置。可填写绝对路径，如：/etc/xray/config.yaml
          </div>
        </el-form-item>
        
        <el-form-item label="API Key" prop="api_key" v-if="scannerForm.type !== 'xray'">
          <el-input v-model="scannerForm.api_key" placeholder="API Key" show-password />
        </el-form-item>
        
        <el-form-item label="二进制路径" prop="binary_path" v-if="scannerForm.type === 'xray'">
          <el-input v-model="scannerForm.binary_path" placeholder="Xray可执行文件的绝对路径" />
          <div class="form-tip">例如: /usr/local/bin/xray 或 C:\tools\xray.exe</div>
        </el-form-item>
        
        <el-form-item label="用户名" prop="username" v-if="scannerForm.type !== 'xray'">
          <el-input v-model="scannerForm.username" placeholder="用户名" />
        </el-form-item>
        
        <el-form-item label="密码" prop="password" v-if="scannerForm.type !== 'xray'">
          <el-input v-model="scannerForm.password" placeholder="密码" show-password />
        </el-form-item>
        
        <el-form-item label="验证SSL" prop="verify_ssl">
          <el-switch v-model="scannerForm.verify_ssl" />
        </el-form-item>
        
        <el-form-item label="启用" prop="enabled">
          <el-switch v-model="scannerForm.enabled" />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="scannerDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitScannerForm" :loading="scannerSubmitLoading">
            {{ scannerForm.isEdit ? '保存' : '添加' }}
          </el-button>
          <el-button type="success" @click="testScannerConnection" :loading="scannerTestLoading">
            测试连接
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { reactive, ref, onMounted, nextTick } from 'vue'
import { ElMessage, ElLoading, ElMessageBox } from 'element-plus'
import { getToken } from '@/utils/auth'
import { 
  getSettings, 
  saveSettings as saveSettingsApi,
  testJiraConnection as testJiraConnectionApi,
  testWechatLogin as testWechatLoginApi,
  testWorkWechatBot as testWorkWechatBotApi,
  testFeishuBot as testFeishuBotApi,
  testDingtalkBot as testDingtalkBotApi,
  testEmailNotification as testEmailNotificationApi,
  testAiService as testAiServiceApi,
  testVulnerabilityNotification as testVulnerabilityNotificationApi,
  testVulnDBConnection as testVulnDBConnectionApi
} from '@/api/settings'
import store from '@/store'
import { Plus, Edit, Delete } from '@element-plus/icons-vue'

export default {
  name: 'Settings',
  components: {
    Plus, Edit, Delete
  },
  setup() {
    const activeTab = ref('integrations')
    const saveLoading = ref(false)
    const jiraTestLoading = ref(false)
    const workWechatTestLoading = ref(false)
    const emailTestLoading = ref(false)
    const wechatTestLoading = ref(false)
    const feishuTestLoading = ref(false)
    const dingtalkTestLoading = ref(false)
    const vulndbTestLoading = ref(false)
    const aiTestLoading = ref(false)
    
    // 邮件收件人输入
    const inputVisible = ref(false)
    const inputValue = ref('')
    const emailInput = ref(null)
    
    // 集成配置表单
    const integrationForm = reactive({
      jira: {
        enabled: false,
        url: '',
        apiToken: '',
        username: '',
        defaultProject: ''
      },
      wechat: {
        enabled: false,
        appId: '',
        appSecret: '',
        callbackUrl: '',
        scope: 'snsapi_base',
        redirectDomain: ''
      },
      vulnDb: {
        enabled: false,
        provider: '',
        apiUrl: '',
        apiKey: '',
        apiSecret: '',
        parameters: ''
      }
    })
    
    // 消息通知表单
    const notificationForm = reactive({
      workWechat: {
        enabled: false,
        webhookUrl: '',
        events: []
      },
      feishu: {
        enabled: false,
        webhookUrl: '',
        secret: '',
        events: []
      },
      dingtalk: {
        enabled: false,
        webhookUrl: '',
        secret: '',
        events: []
      },
      email: {
        enabled: false,
        smtpServer: '',
        smtpPort: 25,
        fromEmail: '',
        username: '',
        password: '',
        useSsl: true,
        events: [],
        recipients: []
      }
    })
    
    // AI功能表单
    const aiForm = reactive({
      enabled: false,
      provider: 'openai',
      apiKey: '',
      apiUrl: '',
      analysisOptions: []
    })
    
    // 通用配置表单
    const generalForm = reactive({
      // 如果需要通用配置，添加相应字段
    })
    
    // 扫描器配置
    const scannersList = ref([])
    const scannerDialogVisible = ref(false)
    const scannerFormRef = ref(null)
    const scannerSubmitLoading = ref(false)
    const scannerTestLoading = ref(false)
    
    // 扫描器表单
    const scannerForm = reactive({
      id: null,
      name: '',
      type: '',
      url: '',
      api_key: '',
      username: '',
      password: '',
      binary_path: '',
      verify_ssl: true,
      enabled: true,
      isEdit: false
    })
    
    // 扫描器表单验证规则
    const scannerRules = {
      name: [
        { required: true, message: '请输入扫描器名称', trigger: 'blur' },
        { min: 2, max: 50, message: '长度在2到50个字符之间', trigger: 'blur' }
      ],
      type: [
        { required: true, message: '请选择扫描器类型', trigger: 'change' }
      ],
      url: [
        { 
          validator: (rule, value, callback) => {
            if (scannerForm.type !== 'xray' && !value) {
              callback(new Error('请输入服务URL'))
            } else if (scannerForm.type !== 'xray' && value && !value.match(/^https?:\/\/.+/)) {
              callback(new Error('URL格式不正确，应以http://或https://开头'))
            } else {
              callback()
            }
          }, 
          trigger: 'blur' 
        }
      ],
      binary_path: [
        { 
          required: true, 
          message: '请输入Xray二进制文件的路径', 
          trigger: 'blur',
          validator: (rule, value, callback) => {
            if (scannerForm.type === 'xray' && !value) {
              callback(new Error('请输入Xray二进制文件的路径'))
            } else {
              callback()
            }
          }
        }
      ]
    }
    
    // 获取扫描器类型标签
    const getScannerTypeLabel = (type) => {
      const types = {
        'nessus': 'Nessus',
        'awvs': 'AWVS',
        'xray': 'Xray',
        'zap': 'ZAP',
        'custom': '自定义'
      }
      return types[type] || type
    }
    
    // 获取扫描器列表
    const fetchScanners = async () => {
      try {
        const response = await getSettings()
        if (response.code === 200 && response.data && response.data.scanners) {
          scannersList.value = response.data.scanners
        }
      } catch (error) {
        console.error('获取扫描器列表失败:', error)
        ElMessage.error('获取扫描器列表失败')
      }
    }
    
    // 显示输入框
    const showInput = () => {
      inputVisible.value = true
      nextTick(() => {
        emailInput.value.focus()
      })
    }
    
    // 添加收件人
    const addRecipient = () => {
      const email = inputValue.value.trim()
      if (!email) return
      
      const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
      if (!emailRegex.test(email)) {
        ElMessage.warning('请输入有效的电子邮件地址')
        return
      }
      
      if (!notificationForm.email.recipients.includes(email)) {
        notificationForm.email.recipients.push(email)
      }
      inputValue.value = ''
    }
    
    // 移除收件人
    const removeRecipient = (email) => {
      const index = notificationForm.email.recipients.indexOf(email)
      if (index !== -1) {
        notificationForm.email.recipients.splice(index, 1)
      }
    }
    
    // 获取设置
    const fetchSettings = async () => {
      try {
        console.log('开始获取系统设置...')
        
        // 检查用户信息和token
        const token = getToken()
        if (!token) {
          console.error('获取设置失败: 未找到认证令牌')
          ElMessage.error('未找到认证令牌')
          return
        }
        
        if (!store.getters.userInfo) {
          console.log('用户信息不存在，尝试获取...')
          try {
            await store.dispatch('user/getUserInfo')
            console.log('用户信息获取成功')
          } catch (error) {
            console.error('获取用户信息失败:', error)
            ElMessage.error('获取用户信息失败')
            return
          }
        }
        
        console.log('发送获取设置请求...')
        const response = await getSettings()
        
        if (response.code === 200 && response.data) {
          console.log('设置数据获取成功:', response.data)
          // 更新表单数据
          const settings = response.data
          
          // 更新集成配置
          if (settings.integrations) {
            if (settings.integrations.jira) {
              Object.assign(integrationForm.jira, settings.integrations.jira)
            }
            if (settings.integrations.wechat) {
              Object.assign(integrationForm.wechat, settings.integrations.wechat)
            }
            if (settings.integrations.vulnDb) {
              Object.assign(integrationForm.vulnDb, settings.integrations.vulnDb)
            }
          }
          
          // 更新消息通知配置
          if (settings.notifications) {
            if (settings.notifications.workWechat) {
              Object.assign(notificationForm.workWechat, settings.notifications.workWechat)
            }
            if (settings.notifications.feishu) {
              Object.assign(notificationForm.feishu, settings.notifications.feishu)
            }
            if (settings.notifications.dingtalk) {
              Object.assign(notificationForm.dingtalk, settings.notifications.dingtalk)
            }
            if (settings.notifications.email) {
              Object.assign(notificationForm.email, settings.notifications.email)
            }
          }
          
          // 更新AI设置
          if (settings.ai) {
            Object.assign(aiForm, settings.ai)
          }
          
          ElMessage.success('设置加载成功')
        } else {
          ElMessage.warning('没有获取到有效的设置数据')
        }
      } catch (error) {
        console.error('获取设置失败:', error)
        ElMessage.error('获取设置失败')
      }
    }
    
    // 测试JIRA连接
    const testJiraConnection = async () => {
      if (!integrationForm.jira.url || !integrationForm.jira.apiToken || !integrationForm.jira.username) {
        ElMessage.error('请填写完整的JIRA配置信息')
        return
      }
      
      jiraTestLoading.value = true
      try {
        const response = await testJiraConnectionApi(integrationForm.jira)
        if (response.success) {
          ElMessage.success('JIRA连接测试成功')
        } else {
          ElMessage.error(response.message || 'JIRA连接测试失败')
        }
      } catch (error) {
        ElMessage.error('JIRA连接测试失败: ' + error.message)
      } finally {
        jiraTestLoading.value = false
      }
    }
    
    // 测试企业微信机器人
    const testWorkWechatNotification = async () => {
      if (!notificationForm.workWechat.webhookUrl) {
        ElMessage.error('请填写企业微信机器人Webhook URL')
        return
      }
      
      workWechatTestLoading.value = true
      try {
        const response = await testWorkWechatBotApi(notificationForm.workWechat)
        if (response.success) {
          ElMessage.success('企业微信机器人测试成功')
        } else {
          ElMessage.error(response.message || '企业微信机器人测试失败')
        }
      } catch (error) {
        ElMessage.error('企业微信机器人测试失败: ' + error.message)
      } finally {
        workWechatTestLoading.value = false
      }
    }
    
    // 测试微信扫码登录
    const testWechatLogin = async () => {
      if (!integrationForm.wechat.appId || !integrationForm.wechat.appSecret) {
        ElMessage.error('请填写完整的微信配置信息')
        return
      }
      
      wechatTestLoading.value = true
      try {
        const response = await testWechatLoginApi(integrationForm.wechat)
        if (response.success) {
          ElMessage.success('微信配置测试成功')
        } else {
          ElMessage.error(response.message || '微信配置测试失败')
        }
      } catch (error) {
        ElMessage.error('微信配置测试失败: ' + error.message)
      } finally {
        wechatTestLoading.value = false
      }
    }
    
    // 测试飞书机器人
    const testFeishuNotification = async () => {
      if (!notificationForm.feishu.webhookUrl) {
        ElMessage.error('请填写飞书机器人Webhook URL')
        return
      }
      
      feishuTestLoading.value = true
      try {
        const response = await testFeishuBotApi(notificationForm.feishu)
        if (response.success) {
          ElMessage.success('飞书机器人测试成功')
        } else {
          ElMessage.error(response.message || '飞书机器人测试失败')
        }
      } catch (error) {
        ElMessage.error('飞书机器人测试失败: ' + error.message)
      } finally {
        feishuTestLoading.value = false
      }
    }
    
    // 测试钉钉机器人
    const testDingtalkNotification = async () => {
      if (!notificationForm.dingtalk.webhookUrl) {
        ElMessage.error('请填写钉钉机器人Webhook URL')
        return
      }
      
      dingtalkTestLoading.value = true
      try {
        const response = await testDingtalkBotApi(notificationForm.dingtalk)
        if (response.success) {
          ElMessage.success('钉钉机器人测试成功')
        } else {
          ElMessage.error(response.message || '钉钉机器人测试失败')
        }
      } catch (error) {
        ElMessage.error('钉钉机器人测试失败: ' + error.message)
      } finally {
        dingtalkTestLoading.value = false
      }
    }
    
    // 测试邮件通知
    const testEmailNotification = async () => {
      if (!notificationForm.email.smtpServer || !notificationForm.email.smtpPort || 
          !notificationForm.email.fromEmail) {
        ElMessage.error('请填写完整的邮件配置信息')
        return
      }
      
      emailTestLoading.value = true
      try {
        const response = await testEmailNotificationApi(notificationForm.email)
        if (response.success) {
          ElMessage.success('邮件配置测试成功')
        } else {
          ElMessage.error(response.message || '邮件配置测试失败')
        }
      } catch (error) {
        ElMessage.error('邮件配置测试失败: ' + error.message)
      } finally {
        emailTestLoading.value = false
      }
    }
    
    // 测试AI服务
    const testAiService = async () => {
      if (!aiForm.provider || !aiForm.apiKey) {
        ElMessage.error('请填写完整的AI服务配置信息')
        return
      }
      
      aiTestLoading.value = true
      try {
        const response = await testAiServiceApi(aiForm)
        if (response.success) {
          ElMessage.success('AI服务配置测试成功')
        } else {
          ElMessage.error(response.message || 'AI服务配置测试失败')
        }
      } catch (error) {
        ElMessage.error('AI服务配置测试失败: ' + error.message)
      } finally {
        aiTestLoading.value = false
      }
    }
    
    // 测试漏洞库连接
    const testVulnDBConnection = async () => {
      if (!integrationForm.vulnDb.apiUrl || !integrationForm.vulnDb.apiKey) {
        ElMessage.error('请填写漏洞库API URL和API Key')
        return
      }
      
      vulndbTestLoading.value = true
      try {
        const response = await testVulnDBConnectionApi(integrationForm.vulnDb)
        if (response.success) {
          ElMessage.success('漏洞库连接测试成功')
        } else {
          ElMessage.error(response.message || '漏洞库连接测试失败')
        }
      } catch (error) {
        ElMessage.error('漏洞库连接测试失败: ' + error.message)
      } finally {
        vulndbTestLoading.value = false
      }
    }
    
    // 保存设置
    const saveSettings = async () => {
      saveLoading.value = true
      
      try {
        const data = {
          integrations: {
            jira: integrationForm.jira,
            wechat: integrationForm.wechat,
            workWechat: notificationForm.workWechat,
            dingtalk: notificationForm.dingtalk,
            feishu: notificationForm.feishu,
            vulnDb: integrationForm.vulnDb
          },
          notifications: {
            email: notificationForm.email
          },
          general: generalForm,
          ai: aiForm
        }
        
        const response = await saveSettingsApi(data)
        if (response.success) {
          ElMessage.success('设置保存成功')
          // 更新store中的设置
          if (store.getters.settings) {
            store.commit('SET_SETTINGS', data)
          }
        } else {
          ElMessage.error(response.message || '设置保存失败')
        }
      } catch (error) {
        ElMessage.error('设置保存失败: ' + error.message)
      } finally {
        saveLoading.value = false
      }
    }
    
    // 重置设置
    const resetSettings = () => {
      fetchSettings()
      ElMessage.info('设置重置成功')
    }
    
    // 测试漏洞通知
    const testNotification = async () => {
      const loading = ElLoading.service({
        lock: true,
        text: '正在发送测试通知...',
        background: 'rgba(0, 0, 0, 0.7)'
      })
      
      try {
        const response = await testVulnerabilityNotificationApi()
        if (response.code === 200) {
          ElMessage.success('测试通知已发送!')
        } else {
          ElMessage.error(response.message || '发送测试通知失败')
        }
      } catch (error) {
        console.error('发送测试通知失败:', error)
        ElMessage.error('发送测试通知失败: ' + (error.message || '未知错误'))
      } finally {
        loading.close()
      }
    }
    
    // 显示添加扫描器对话框
    const showAddScannerDialog = () => {
      // 重置表单
      Object.keys(scannerForm).forEach(key => {
        if (key !== 'verify_ssl' && key !== 'enabled') {
          scannerForm[key] = ''
        }
      })
      scannerForm.verify_ssl = true
      scannerForm.enabled = true
      scannerForm.isEdit = false
      
      // 显示对话框
      scannerDialogVisible.value = true
      
      // 重置表单验证
      nextTick(() => {
        if (scannerFormRef.value) {
          scannerFormRef.value.resetFields()
        }
      })
    }
    
    // 编辑扫描器
    const editScanner = (scanner) => {
      // 填充表单
      Object.keys(scannerForm).forEach(key => {
        if (key in scanner) {
          scannerForm[key] = scanner[key]
        }
      })
      scannerForm.isEdit = true
      
      // 显示对话框
      scannerDialogVisible.value = true
      
      // 重置表单验证
      nextTick(() => {
        if (scannerFormRef.value) {
          scannerFormRef.value.resetFields()
        }
      })
    }
    
    // 删除扫描器
    const deleteScanner = (scanner) => {
      ElMessageBox.confirm(
        `确定要删除扫描器 "${scanner.name}" 吗？`,
        '警告',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }
      ).then(async () => {
        try {
          // 调用API删除扫描器
          await deleteScanner(scanner.id)
          ElMessage.success('删除成功')
          fetchScanners() // 刷新列表
        } catch (error) {
          console.error('删除扫描器失败:', error)
          ElMessage.error('删除扫描器失败')
        }
      }).catch(() => {})
    }
    
    // 提交扫描器表单
    const submitScannerForm = async () => {
      if (!scannerFormRef.value) return
      
      await scannerFormRef.value.validate(async (valid) => {
        if (!valid) return
        
        scannerSubmitLoading.value = true
        try {
          if (scannerForm.isEdit) {
            // 调用API更新扫描器
            await updateScanner(scannerForm)
            ElMessage.success('更新扫描器成功')
          } else {
            // 调用API添加扫描器
            await addScanner(scannerForm)
            ElMessage.success('添加扫描器成功')
          }
          
          // 关闭对话框并刷新列表
          scannerDialogVisible.value = false
          fetchScanners()
        } catch (error) {
          console.error('保存扫描器失败:', error)
          ElMessage.error(error.response?.data?.message || '保存扫描器失败')
        } finally {
          scannerSubmitLoading.value = false
        }
      })
    }
    
    // 测试扫描器连接
    const testScannerConnection = async () => {
      if (!scannerFormRef.value) return
      
      await scannerFormRef.value.validate(async (valid) => {
        if (!valid) return
        
        scannerTestLoading.value = true
        try {
          // 调用API测试扫描器连接
          await testScannerConnection(scannerForm)
          ElMessage.success('连接测试成功')
        } catch (error) {
          console.error('连接测试失败:', error)
          ElMessage.error(error.response?.data?.message || '连接测试失败')
        } finally {
          scannerTestLoading.value = false
        }
      })
    }
    
    // 初始化加载设置
    onMounted(async () => {
      console.log('设置页面挂载，准备初始化...')
      const token = getToken()
      
      if (!token) {
        console.error('设置页面初始化: 未找到令牌')
        ElMessage.error('未找到认证令牌')
        return
      }
      
      if (!store.getters.userInfo) {
        console.log('设置页面初始化: 用户信息不存在，尝试获取...')
        try {
          await store.dispatch('user/getUserInfo')
          console.log('设置页面初始化: 用户信息获取成功')
        } catch (error) {
          console.error('设置页面初始化: 获取用户信息失败', error)
          ElMessage.error('获取用户信息失败')
          return
        }
      }
      
      // 加载设置
      await fetchSettings()
      console.log('设置页面初始化完成')
      fetchScanners() // 获取扫描器列表
    })
    
    return {
      activeTab,
      saveLoading,
      jiraTestLoading,
      workWechatTestLoading,
      emailTestLoading,
      wechatTestLoading,
      feishuTestLoading,
      dingtalkTestLoading,
      vulndbTestLoading,
      aiTestLoading,
      integrationForm,
      notificationForm,
      aiForm,
      inputVisible,
      inputValue,
      emailInput,
      showInput,
      addRecipient,
      removeRecipient,
      testJiraConnection,
      testWorkWechatNotification,
      testWechatLogin,
      testFeishuNotification,
      testDingtalkNotification,
      testEmailNotification,
      testAiService,
      testVulnDBConnection,
      saveSettings,
      resetSettings,
      testNotification,
      scannersList,
      scannerDialogVisible,
      scannerFormRef,
      scannerForm,
      scannerRules,
      scannerSubmitLoading,
      scannerTestLoading,
      showAddScannerDialog,
      editScanner,
      deleteScanner,
      submitScannerForm,
      testScannerConnection,
      getScannerTypeLabel
    }
  }
}
</script>

<style scoped>
.settings-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.tab-content {
  padding: 20px 0;
}

.form-actions {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}

.recipient-tag {
  margin-right: 10px;
  margin-bottom: 10px;
}

.input-new-tag {
  width: 200px;
  margin-right: 10px;
  vertical-align: bottom;
}

.button-new-tag {
  height: 32px;
  line-height: 30px;
  padding-top: 0;
  padding-bottom: 0;
}

/* 添加高级渐变背景和动画 */
.settings-container {
  background: linear-gradient(to bottom right, #f5f7fa, #c3cfe2);
  animation: fadeIn 0.5s ease-in-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 添加阴影效果 */
.el-card {
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
}

.el-card:hover {
  box-shadow: 0 15px 40px rgba(0, 0, 0, 0.15);
}

/* Tab样式优化 */
:deep(.el-tabs__item) {
  padding: 0 20px;
  height: 50px;
  line-height: 50px;
  font-size: 15px;
  transition: all 0.3s;
}

:deep(.el-tabs__item.is-active) {
  font-weight: bold;
  background: linear-gradient(90deg, #4e54c8, #8f94fb);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

:deep(.el-tabs__active-bar) {
  background: linear-gradient(90deg, #4e54c8, #8f94fb);
  height: 3px;
}

/* 表单元素样式 */
:deep(.el-input__inner) {
  border-radius: 6px;
}

:deep(.el-button--primary) {
  background: linear-gradient(90deg, #4e54c8, #8f94fb);
  border: none;
  border-radius: 6px;
  transition: all 0.3s;
}

:deep(.el-button--primary:hover) {
  transform: translateY(-2px);
  box-shadow: 0 6px 12px rgba(78, 84, 200, 0.25);
}

/* 响应式调整 */
@media (max-width: 768px) {
  .settings-container {
    padding: 10px;
  }
  
  .tab-content {
    padding: 15px 0;
  }
  
  :deep(.el-form-item__label) {
    width: 100% !important;
    text-align: left;
    margin-bottom: 8px;
  }
  
  :deep(.el-form-item__content) {
    margin-left: 0 !important;
  }
}
</style> 