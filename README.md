# VulnArk - 漏洞管理平台

<div align="center">
  <img src="https://img.shields.io/badge/version-0.1.0-blue.svg" alt="Version">
  <img src="https://img.shields.io/badge/license-MIT-green.svg" alt="License">
  <img src="https://img.shields.io/badge/Go-1.18+-00ADD8.svg" alt="Go">
  <img src="https://img.shields.io/badge/Vue-3.x-4FC08D.svg" alt="Vue">
</div>

[English](#vulnark---vulnerability-management-platform) | [中文](#vulnark---漏洞管理平台-1)

## VulnArk - Vulnerability Management Platform

### Introduction

VulnArk is a modern vulnerability management platform designed to help security teams efficiently discover, track, and remediate security vulnerabilities across their organization. With powerful features like asset management, vulnerability tracking, knowledge base, and automated scanning, VulnArk provides a comprehensive solution for the entire vulnerability lifecycle management.

### Key Features

- **Dashboard**: Real-time overview of vulnerability statistics, trends, and recent activities
- **Vulnerability Management**: Create, update, track, and remediate vulnerabilities
- **Asset Management**: Manage and classify organizational assets with vulnerability mapping
- **Knowledge Base**: Document and share security best practices and remediation guides
- **Vulnerability Database**: Maintain a comprehensive database of known vulnerabilities
- **Scan Integration**: Schedule and manage automated vulnerability scans
- **User Management**: Role-based access control system
- **Notification System**: Customizable alerts for vulnerability events
- **AI-powered Analysis**: Risk assessment and prioritization with AI capabilities

### Technology Stack

- **Frontend**: Vue.js 3, Element Plus, ECharts
- **Backend**: Go (Gin framework)
- **Database**: MySQL 8.0
- **Deployment**: Docker & Docker Compose

### Quick Start

#### Docker Deployment (Recommended)

1. Clone the repository:
```bash
git clone https://github.com/yourusername/vulnark.git
cd vulnark
```

2. Run the deployment script:
```bash
chmod +x deploy.sh
./deploy.sh
```

3. Access the application:
   - Frontend: http://localhost
   - Default admin account: 
     - Username: `admin`
     - Password: `admin123`

#### Manual Deployment

For manual deployment instructions, please refer to the [Docker Deployment Guide](README.Docker.md).

### System Architecture

VulnArk follows a microservices architecture with three main components:

1. **Frontend Service**: Vue.js application served by Nginx
2. **Backend Service**: Go API server providing business logic and data access
3. **Database Service**: MySQL database for persistent storage

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│             │     │             │     │             │
│   Frontend  │────▶│   Backend   │────▶│   Database  │
│   (Nginx)   │     │   (Go API)  │     │   (MySQL)   │
│             │     │             │     │             │
└─────────────┘     └─────────────┘     └─────────────┘
```

### Screenshots

*Please add your application screenshots here*

### Configuration

VulnArk can be configured through several methods:

1. **Environment Variables**: Set in docker-compose.yml
2. **Config Files**: Modify backend/config/config.yaml
3. **Database Settings**: System settings stored in database

Key configuration options:

- Database connection settings
- JWT authentication settings
- Logging options
- Notification preferences
- AI service integration

For full configuration options, see the [Configuration Documentation](docs/configuration.md).

### Development Guide

#### Prerequisites

- Go 1.18+
- Node.js 16+
- MySQL 8.0+

#### Setup Development Environment

1. **Backend Development**:
```bash
cd backend
go mod download
go run main.go
```

2. **Frontend Development**:
```bash
cd frontend
npm install
npm run serve
```

### Contributing

We welcome contributions! Please follow these steps:

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to your branch
5. Create a Pull Request

Please make sure to follow our [Code of Conduct](CODE_OF_CONDUCT.md) and [Contribution Guidelines](CONTRIBUTING.md).

### License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## VulnArk - 漏洞管理平台

### 介绍

VulnArk 是一个现代化的漏洞管理平台，旨在帮助安全团队高效地发现、跟踪和修复组织内的安全漏洞。通过强大的功能，如资产管理、漏洞跟踪、知识库和自动扫描，VulnArk 为整个漏洞生命周期管理提供了全面的解决方案。

### 主要功能

- **仪表盘**：实时概览漏洞统计、趋势和最近活动
- **漏洞管理**：创建、更新、跟踪和修复漏洞
- **资产管理**：管理和分类组织资产，并映射相关漏洞
- **知识库**：记录和分享安全最佳实践和修复指南
- **漏洞库**：维护已知漏洞的综合数据库
- **扫描集成**：调度和管理自动化漏洞扫描
- **用户管理**：基于角色的访问控制系统
- **通知系统**：可自定义的漏洞事件告警
- **AI 驱动分析**：利用 AI 能力进行风险评估和优先级划分

### 技术栈

- **前端**：Vue.js 3、Element Plus、ECharts
- **后端**：Go（Gin 框架）
- **数据库**：MySQL 8.0
- **部署**：Docker 和 Docker Compose

### 快速开始

#### Docker 部署（推荐）

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

#### 手动部署

有关手动部署的说明，请参考 [Docker 部署指南](README.Docker.md)。

### 系统架构

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

### 截图展示

*请在此处添加应用截图*

### 配置说明

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

### 开发指南

#### 先决条件

- Go 1.18+
- Node.js 16+
- MySQL 8.0+

#### 设置开发环境

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

### 贡献指南

我们欢迎贡献！请遵循以下步骤：

1. 复刻（Fork）代码库
2. 创建功能分支
3. 提交您的更改
4. Push to your branch
5. Create a Pull Request

请确保遵循我们的 [行为准则](CODE_OF_CONDUCT.md) 和 [贡献指南](CONTRIBUTING.md).

### 许可证

该项目采用 MIT 许可证 - 详情请参见 [LICENSE](LICENSE) 文件。 