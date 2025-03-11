<template>
  <div class="analysis-parameters">
    <el-form :model="form" label-width="120px" label-position="top" :rules="rules" ref="analysisForm">
      <!-- 趋势预测参数 -->
      <template v-if="analysisType === 'trend_prediction'">
        <el-form-item label="预测时间范围" prop="timeRange">
          <el-select v-model="form.timeRange" placeholder="请选择预测时间范围" style="width: 100%">
            <el-option label="未来7天" value="未来7天"></el-option>
            <el-option label="未来14天" value="未来14天"></el-option>
            <el-option label="未来30天" value="未来30天"></el-option>
            <el-option label="未来90天" value="未来90天"></el-option>
          </el-select>
        </el-form-item>
        
        <el-form-item label="分析维度" prop="dimensions">
          <el-checkbox-group v-model="form.dimensions">
            <el-checkbox label="severity">按严重程度</el-checkbox>
            <el-checkbox label="status">按状态</el-checkbox>
            <el-checkbox label="category">按类别</el-checkbox>
            <el-checkbox label="system">按系统</el-checkbox>
          </el-checkbox-group>
        </el-form-item>
        
        <el-form-item label="数据来源" prop="dataSources">
          <el-checkbox-group v-model="form.dataSources">
            <el-checkbox label="internal">内部漏洞数据</el-checkbox>
            <el-checkbox label="external">外部威胁情报</el-checkbox>
          </el-checkbox-group>
        </el-form-item>
        
        <el-form-item label="置信度阈值" prop="confidenceThreshold">
          <el-slider v-model="form.confidenceThreshold" :step="5" :marks="{
            0: '0%',
            25: '25%',
            50: '50%',
            75: '75%',
            100: '100%'
          }"></el-slider>
        </el-form-item>
      </template>
      
      <!-- 资源优化参数 -->
      <template v-if="analysisType === 'resource_optimization'">
        <el-form-item label="优化目标" prop="optimizationGoal">
          <el-radio-group v-model="form.optimizationGoal">
            <el-radio label="risk_reduction">最大化风险降低</el-radio>
            <el-radio label="efficiency">最大化资源效率</el-radio>
            <el-radio label="balanced">平衡风险与效率</el-radio>
          </el-radio-group>
        </el-form-item>
        
        <el-form-item label="考虑因素" prop="factors">
          <el-checkbox-group v-model="form.factors">
            <el-checkbox label="cost">成本因素</el-checkbox>
            <el-checkbox label="time">时间因素</el-checkbox>
            <el-checkbox label="expertise">专业知识要求</el-checkbox>
            <el-checkbox label="business_impact">业务影响</el-checkbox>
          </el-checkbox-group>
        </el-form-item>
        
        <el-form-item label="资源类型" prop="resourceTypes">
          <el-checkbox-group v-model="form.resourceTypes">
            <el-checkbox label="human">人力资源</el-checkbox>
            <el-checkbox label="technical">技术资源</el-checkbox>
            <el-checkbox label="financial">财务资源</el-checkbox>
          </el-checkbox-group>
        </el-form-item>
        
        <el-form-item label="历史数据范围" prop="historyTimeRange">
          <el-select v-model="form.historyTimeRange" placeholder="请选择历史数据范围" style="width: 100%">
            <el-option label="过去30天" value="过去30天"></el-option>
            <el-option label="过去90天" value="过去90天"></el-option>
            <el-option label="过去180天" value="过去180天"></el-option>
            <el-option label="过去365天" value="过去365天"></el-option>
          </el-select>
        </el-form-item>
      </template>
      
      <!-- 异常检测参数 -->
      <template v-if="analysisType === 'anomaly_detection'">
        <el-form-item label="分析时间范围" prop="timeRange">
          <el-select v-model="form.timeRange" placeholder="请选择分析时间范围" style="width: 100%">
            <el-option label="过去3天" value="过去3天"></el-option>
            <el-option label="过去7天" value="过去7天"></el-option>
            <el-option label="过去14天" value="过去14天"></el-option>
            <el-option label="过去30天" value="过去30天"></el-option>
          </el-select>
        </el-form-item>
        
        <el-form-item label="检测灵敏度" prop="sensitivity">
          <el-slider v-model="form.sensitivity" :min="1" :max="10" :step="1" :marks="{
            1: '低',
            5: '中',
            10: '高'
          }"></el-slider>
        </el-form-item>
        
        <el-form-item label="异常类型" prop="anomalyTypes">
          <el-checkbox-group v-model="form.anomalyTypes">
            <el-checkbox label="spike">突增异常</el-checkbox>
            <el-checkbox label="drop">突降异常</el-checkbox>
            <el-checkbox label="pattern">模式异常</el-checkbox>
            <el-checkbox label="seasonal">季节性异常</el-checkbox>
          </el-checkbox-group>
        </el-form-item>
        
        <el-form-item label="分析范围" prop="scopes">
          <el-checkbox-group v-model="form.scopes">
            <el-checkbox label="all">全部漏洞</el-checkbox>
            <el-checkbox label="high">仅高危漏洞</el-checkbox>
            <el-checkbox label="web">Web应用漏洞</el-checkbox>
            <el-checkbox label="network">网络设备漏洞</el-checkbox>
            <el-checkbox label="system">系统漏洞</el-checkbox>
          </el-checkbox-group>
        </el-form-item>
      </template>
      
      <el-form-item label="分析说明" prop="description">
        <el-input 
          type="textarea" 
          v-model="form.description"
          :rows="3"
          placeholder="请输入分析说明（可选）">
        </el-input>
      </el-form-item>
    </el-form>
    
    <div class="form-actions">
      <el-button @click="cancel">取消</el-button>
      <el-button type="primary" @click="submitForm" :loading="loading">开始分析</el-button>
    </div>
  </div>
</template>

<script>
export default {
  name: 'AnalysisParameters',
  props: {
    analysisType: {
      type: String,
      required: true,
      validator: value => ['trend_prediction', 'resource_optimization', 'anomaly_detection'].includes(value)
    }
  },
  data() {
    return {
      loading: false,
      form: {
        // 通用参数
        description: '',
        
        // 趋势预测参数
        timeRange: '未来30天',
        dimensions: ['severity', 'status'],
        dataSources: ['internal'],
        confidenceThreshold: 70,
        
        // 资源优化参数
        optimizationGoal: 'balanced',
        factors: ['cost', 'time'],
        resourceTypes: ['human', 'technical'],
        historyTimeRange: '过去90天',
        
        // 异常检测参数
        sensitivity: 5,
        anomalyTypes: ['spike', 'pattern'],
        scopes: ['all']
      },
      rules: {
        // 趋势预测规则
        timeRange: [
          { required: true, message: '请选择预测时间范围', trigger: 'change' }
        ],
        dimensions: [
          { type: 'array', required: true, message: '请至少选择一个分析维度', trigger: 'change' }
        ],
        
        // 资源优化规则
        optimizationGoal: [
          { required: true, message: '请选择优化目标', trigger: 'change' }
        ],
        
        // 异常检测规则
        anomalyTypes: [
          { type: 'array', required: true, message: '请至少选择一种异常类型', trigger: 'change' }
        ]
      }
    }
  },
  methods: {
    submitForm() {
      this.$refs.analysisForm.validate(valid => {
        if (valid) {
          this.loading = true
          
          // 根据分析类型获取相关参数
          const parameters = this.getParameters()
          
          // 发射事件，传递参数
          this.$emit('start', parameters)
          
          // 注意：loading状态会在父组件中重置
        } else {
          this.$message.warning('请完善分析参数')
          return false
        }
      })
    },
    cancel() {
      this.$emit('cancel')
    },
    getParameters() {
      // 基础参数
      const commonParams = {
        description: this.form.description
      }
      
      // 根据分析类型返回相应参数
      switch (this.analysisType) {
        case 'trend_prediction':
          return {
            ...commonParams,
            timeRange: this.form.timeRange,
            dimensions: this.form.dimensions,
            dataSources: this.form.dataSources,
            confidenceThreshold: this.form.confidenceThreshold
          }
        
        case 'resource_optimization':
          return {
            ...commonParams,
            optimizationGoal: this.form.optimizationGoal,
            factors: this.form.factors,
            resourceTypes: this.form.resourceTypes,
            historyTimeRange: this.form.historyTimeRange
          }
        
        case 'anomaly_detection':
          return {
            ...commonParams,
            timeRange: this.form.timeRange,
            sensitivity: this.form.sensitivity,
            anomalyTypes: this.form.anomalyTypes,
            scopes: this.form.scopes
          }
        
        default:
          return commonParams
      }
    }
  }
}
</script>

<style lang="scss" scoped>
.analysis-parameters {
  padding: 10px;
  
  .form-actions {
    display: flex;
    justify-content: flex-end;
    margin-top: 30px;
    
    .el-button {
      margin-left: 10px;
    }
  }
}
</style> 