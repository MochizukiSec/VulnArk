<template>
  <div class="not-found">
    <div class="content">
      <h1>404</h1>
      <h2>页面不存在</h2>
      <p>抱歉，您访问的页面不存在或已被移除</p>
      <el-button type="primary" @click="goHome">返回首页</el-button>
    </div>
  </div>
</template>

<script>
import { useRouter } from 'vue-router'

export default {
  name: 'NotFound',
  setup() {
    const router = useRouter()

    const goHome = () => {
      router.push('/')
    }

    return {
      goHome
    }
  }
}
</script>

<style scoped>
.not-found {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background-color: #f5f7fa;
}

.content {
  text-align: center;
  padding: 30px;
}

h1 {
  font-size: 120px;
  margin: 0;
  color: #409eff;
}

h2 {
  font-size: 30px;
  margin: 20px 0;
  color: #333;
}

p {
  font-size: 16px;
  color: #666;
  margin-bottom: 30px;
}
</style> 