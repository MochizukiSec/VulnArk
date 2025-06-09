# VulnArk - 漏洞管理平台

<div align="center">
  <img src="https://img.shields.io/badge/version-0.1.0-blue.svg" alt="Version">
  <img src="https://img.shields.io/badge/license-MIT-green.svg" alt="License">
  <img src="https://img.shields.io/badge/Go-1.18+-00ADD8.svg" alt="Go">
  <img src="https://img.shields.io/badge/Vue-3.x-4FC08D.svg" alt="Vue">
</div>

## 作者最近正在重构此项目，将采用Rust作为后端语言，优化前后端代码，新增核心功能。将会在2025年7月前发布。

## 📋 更新日志

### v0.1.1 (2025-06-05)
**🔧 数据库连接修复更新**

本次更新主要修复了Docker部署环境下MySQL连接失败的问题：

#### 🚀 新增功能
- 添加数据库连接重试机制（5次重试，每次间隔5秒）
- 增强数据库连接错误处理和日志输出
- 支持环境变量优先级配置读取

#### 🐛 问题修复
- **修复配置文件路径不匹配问题**：Docker容器现在正确使用`config.docker.yaml`配置文件
- **修复数据库配置读取错误**：修正配置路径从`database.mysql.*`到`database.*`
- **修复环境变量处理问题**：添加完整的环境变量处理逻辑，优先读取Docker环境变量
- **优化部署脚本**：简化配置文件复制流程，避免配置冲突

#### 📝 技术改进
- 重构数据库连接初始化逻辑
- 添加`getConfigString()`和`getConfigInt()`辅助函数
- 增强连接状态日志，便于问题排查
- 统一Docker和开发环境的配置文件结构

#### 📖 文档更新
- 删除README.md中的英文内容，只保留中文说明
- 优化文档结构和标题层级

---

## 介绍

VulnArk 是一个现代化的漏洞管理平台，旨在帮助安全团队高效地发现、跟踪和修复组织内的安全漏洞。通过强大的功能，如资产管理、漏洞跟踪、知识库和自动扫描，VulnArk 为整个漏洞生命周期管理提供了全面的解决方案。

## 主要功能

- **仪表盘**：实时概览漏洞统计、趋势和最近活动
- **漏洞管理**：创建、更新、跟踪和修复漏洞
- **资产管理**：管理和分类组织资产，并映射相关漏洞
- **知识库**：记录和分享安全最佳实践和修复指南
- **漏洞库**：维护已知漏洞的综合数据库
- **扫描集成**：调度和管理自动化漏洞扫描
- **用户管理**：基于角色的访问控制系统
- **通知系统**：可自定义的漏洞事件告警
- **AI 驱动分析**：利用 AI 能力进行风险评估和优先级划分

## 技术栈

- **前端**：Vue.js 3、Element Plus、ECharts
- **后端**：Go（Gin 框架）
- **数据库**：MySQL 8.0
- **部署**：Docker 和 Docker Compose

## 快速开始

### Docker 部署（推荐）

1. 克隆代码库：
```bash
git clone https://github.com/yourusername/vulnark.git
cd vulnark
```

2. 运行部署脚本：
```bash
chmod +x deploy.sh
./deploy.sh
```

3. 访问应用：
   - 前端：http://localhost
   - 默认管理员账号：
     - 用户名：`admin`
     - 密码：`admin123`

### 手动部署

有关手动部署的说明，请参考 [Docker 部署指南](README.Docker.md)。

## 系统架构

VulnArk 采用微服务架构，主要包含三个组件：

1. **前端服务**：由 Nginx 提供服务的 Vue.js 应用
2. **后端服务**：提供业务逻辑和数据访问的 Go API 服务器
3. **数据库服务**：用于持久化存储的 MySQL 数据库

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│             │     │             │     │             │
│    前端     │────▶│    后端     │────▶│   数据库    │
│   (Nginx)   │     │  (Go API)   │     │   (MySQL)   │
│             │     │             │     │             │
└─────────────┘     └─────────────┘     └─────────────┘
```

## 截图展示

<img width="1511" alt="iShot_2025-03-23_08 29 13" src="https://github.com/user-attachments/assets/423d2888-6100-4b6f-bd12-0c165e8b3cc7" />
<img width="1441" alt="iShot_2025-03-23_08 29 25" src="https://github.com/user-attachments/assets/9d6268ee-1b50-4072-abd3-42772e72136b" />
<img width="1511" alt="iShot_2025-03-23_08 27 45" src="https://github.com/user-attachments/assets/06ab81aa-577a-48c8-8ba8-f25d51ee634a" />
<img width="1474" alt="iShot_2025-03-23_08 29 36" src="https://github.com/user-attachments/assets/bc43f23a-23ee-4023-834d-bc0251698c7b" />
<img width="1507" alt="iShot_2025-03-23_08 27 54" src="https://github.com/user-attachments/assets/a0f8ac6c-27f2-44e7-b274-375fcabf7bd6" />
<img width="1508" alt="iShot_2025-03-23_08 28 01" src="https://github.com/user-attachments/assets/6309135e-56bf-49e7-b0ca-882bb86e9663" />
<img width="1503" alt="iShot_2025-03-23_08 28 22" src="https://github.com/user-attachments/assets/ef1fb8ff-63ec-4c42-97a8-e0e39b21f611" />
<img width="1497" alt="iShot_2025-03-23_08 28 29" src="https://github.com/user-attachments/assets/77062651-3a04-46fd-b313-ac3acd0b6f17" />
<img width="1509" alt="iShot_2025-03-23_08 28 37" src="https://github.com/user-attachments/assets/9f2954b2-1ca3-4e64-9755-eaed5bb5edeb" />
<img width="1505" alt="iShot_2025-03-23_08 28 51" src="https://github.com/user-attachments/assets/bbc39ac5-5a06-4576-ae5f-3eb63f0336a5" />
<img width="1509" alt="iShot_2025-03-23_08 29 05" src="https://github.com/user-attachments/assets/db939164-aee1-46a4-aeb0-7905605dac2a" />

## 配置说明

VulnArk 可以通过多种方式进行配置：

1. **环境变量**：在 docker-compose.yml 中设置
2. **配置文件**：修改 backend/config/config.yaml
3. **数据库设置**：存储在数据库中的系统设置

主要配置选项：

- 数据库连接设置
- JWT 认证设置
- 日志选项
- 通知首选项
- AI 服务集成

完整配置选项，请参见 [配置文档](docs/configuration.md)。

## 开发指南

### 先决条件

- Go 1.18+
- Node.js 16+
- MySQL 8.0+

### 设置开发环境

1. **后端开发**：
```bash
cd backend
go mod download
go run main.go
```

2. **前端开发**：
```bash
cd frontend
npm install
npm run serve
```

## 贡献指南

我们欢迎贡献！请遵循以下步骤：

1. 复刻（Fork）代码库
2. 创建功能分支
3. 提交您的更改
4. 推送到您的分支
5. 创建拉取请求（Pull Request）

请确保遵循我们的 [行为准则](CODE_OF_CONDUCT.md) 和 [贡献指南](CONTRIBUTING.md)。

## 许可证

该项目采用 MIT 许可证 - 详情请参见 [LICENSE](LICENSE) 文件。 
