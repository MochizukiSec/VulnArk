# VulnArk - 安全漏洞管理平台

<div align="center">
  <img src="docs/images/logo.svg" alt="VulnArk Logo" width="200">
  <p>安全漏洞全生命周期管理系统，集成报告生成与数据分析功能</p>
</div>

[![License](https://img.shields.io/github/license/MochizukiSec/VulnArk.svg)](https://github.com/MochizukiSec/VulnArk/blob/main/LICENSE)

## 📖 概述

VulnArk是一个现代化的漏洞管理平台，帮助安全团队高效管理漏洞的全生命周期。系统支持漏洞的记录、跟踪、优先级管理、报告生成和数据分析，提供直观的用户界面和全面的API支持。

## ✨ 核心功能

- **漏洞管理**: 全面管理系统安全漏洞，基于CVSS评分评估风险
- **仪表盘分析**: 直观的安全数据可视化，帮助快速识别安全趋势
- **安全报告**: 支持多种报告类型和格式，包括PDF、Excel、Word和HTML
- **用户与权限**: 完善的用户管理和基于角色的权限控制系统
- **系统集成**: 灵活的API设计，便于与其他安全工具集成
- **安全与合规**: 内置的安全机制，确保系统合规和安全

## 🚀 快速开始

### 前提条件

- Go 1.19或更高版本
- Node.js 16或更高版本和npm
- MongoDB 5.0或更高版本
- 至少4GB内存和2核CPU
- 20GB可用磁盘空间

### 源码部署方式

#### 后端部署

1. 克隆仓库
   ```bash
   git clone https://github.com/MochizukiSec/VulnArk.git
   cd VulnArk
   ```

2. 配置环境变量
   创建`.env`文件在后端目录下（或复制`.env.example`）:
   ```bash
   cd backend
   cp .env.example .env
   # 编辑.env文件，配置MongoDB连接和其他参数
   ```

3. 安装依赖并构建
   ```bash
   go mod download
   go build -o vulnark-server
   ```

4. 启动后端服务
   ```bash
   ./vulnark-server
   ```

#### 前端部署

1. 进入前端目录
   ```bash
   cd frontend
   ```

2. 安装依赖
   ```bash
   npm install
   ```

3. 配置API地址
   编辑`.env.production`文件（如果不存在则创建）:
   ```
   VUE_APP_API_URL=http://your-backend-server:8000
   ```

4. 构建生产版本
   ```bash
   npm run build
   ```

5. 部署静态文件
   将`dist`目录下的文件部署到您的Web服务器（如Nginx、Apache等）

#### Nginx配置示例

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        root /path/to/frontend/dist;
        try_files $uri $uri/ /index.html;
    }

    location /api/ {
        proxy_pass http://localhost:8000/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## 🔧 配置选项

VulnArk提供多种配置选项，可通过环境变量进行设置。

### 环境变量

| 环境变量 | 说明 | 默认值 |
|----------|------|---------|
| PORT | API服务端口 | 8000 |
| MONGO_URI | MongoDB连接URI | mongodb://localhost:27017 |
| MONGO_DB_NAME | MongoDB数据库名称 | vulnark_db |
| JWT_SECRET | JWT密钥 | your-secret-key-for-production |
| ALLOWED_ORIGINS | CORS允许的源 | http://localhost:8080 |
| LOG_LEVEL | 日志级别 | info |

## 📊 系统架构

VulnArk采用现代化的微服务架构，分为前端、后端API和数据库三层：

![架构图](docs/images/architecture.svg)

- **前端**: Vue.js + Element Plus构建的SPA应用
- **后端API**: Go语言开发的RESTful API
- **数据库**: MongoDB提供的文档存储

## 🔒 系统安全

VulnArk内置了多层安全机制：

- JWT基于角色的身份认证
- API请求速率限制
- 数据输入验证和清洗
- 密码加密存储
- 审计日志记录

## 📋 系统要求

### 最低配置
- 2核CPU
- 4GB内存
- 20GB磁盘空间

### 推荐配置
- 4核CPU
- 8GB内存
- 50GB SSD磁盘空间

## 🛠 常见问题

### 1. 启动后端服务时MongoDB连接失败

**问题**: 启动后端服务时报错 "failed to connect to MongoDB"

**解决方案**:
- 确认MongoDB服务已正常运行
- 检查.env文件中的MONGO_URI配置是否正确
- 确认MongoDB用户名和密码正确（如果启用了身份验证）

### 2. 前端API请求失败

**问题**: 前端页面加载但无法获取数据，浏览器控制台显示API请求错误

**解决方案**:
- 确认后端服务正常运行
- 检查前端环境变量中的API地址配置是否正确
- 查看浏览器控制台错误信息，检查是否存在CORS问题
  - 如果存在CORS问题，请确保后端ALLOWED_ORIGINS环境变量包含前端域名
- 检查网络请求是否有身份验证错误，尝试重新登录

### 3. 用户注册/登录问题

**问题**: 无法注册新用户或登录失败

**解决方案**:
- 确认后端服务和数据库连接正常
- 检查日志中是否有详细错误信息
- 对于登录问题，可尝试重置密码

### 4. 如何更新到最新版本？

```bash
# 拉取最新代码
git pull origin main

# 后端更新
cd backend
go mod download
go build -o vulnark-server
# 重启后端服务

# 前端更新
cd frontend
npm install
npm run build
# 重新部署dist目录
```

### 5. 数据备份与恢复

**备份MongoDB数据**:

```bash
# 本地MongoDB
mongodump --db vulnark_db --out /backup/$(date +%Y%m%d)
```

**恢复MongoDB数据**:

```bash
# 本地MongoDB
mongorestore --db vulnark_db /backup/20230101/vulnark_db
```

### 6. 性能优化建议

如果系统运行缓慢，可以尝试以下优化措施：

- 为MongoDB创建适当的索引
  ```javascript
  db.vulnerabilities.createIndex({ "cvss_score": 1 })
  db.vulnerabilities.createIndex({ "status": 1 })
  db.vulnerabilities.createIndex({ "created_at": 1 })
  ```
- 增加后端服务的资源配置
- 实现API响应缓存
- 优化前端资源加载和渲染
- 考虑使用CDN加速静态资源
- 

## ☎️联系作者
![WechatIMG5](https://github.com/user-attachments/assets/97df7e28-5ada-47dc-a5be-e540122929cf)
![WechatIMG6](https://github.com/user-attachments/assets/da85b4c8-531e-4b7e-9234-b70bb1f14301)

## 🤝 贡献指南

我们欢迎社区贡献，请参阅[贡献指南](CONTRIBUTING.md)了解如何参与项目开发。

## 📜 许可证

本项目采用MIT许可证 - 详见[LICENSE](LICENSE)文件 
