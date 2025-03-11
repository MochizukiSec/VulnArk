# VulnArk 入门指南

本文档将帮助您快速上手 VulnArk 平台，包括安装部署、基本配置和初始使用。

## 目录

- [系统要求](#系统要求)
- [安装方法](#安装方法)
  - [使用 Docker Compose](#使用-docker-compose)
  - [使用 Kubernetes](#使用-kubernetes)
- [初始配置](#初始配置)
- [首次登录](#首次登录)
- [基本操作指南](#基本操作指南)
- [常见问题](#常见问题)

## 系统要求

在开始部署 VulnArk 前，请确保您的系统满足以下要求：

### 最低配置
- 2核CPU
- 4GB内存
- 20GB磁盘空间
- Docker 20.10.x 或更高版本
- Docker Compose 2.x 或更高版本

### 推荐配置
- 4核CPU
- 8GB内存
- 50GB SSD磁盘空间
- 现代化的Linux发行版（如Ubuntu 20.04/22.04, CentOS 8/9）

## 安装方法

### 使用 Docker Compose

Docker Compose 是最简单的部署方式，适合单机部署或小规模团队使用。

1. 克隆仓库
   ```bash
   git clone https://github.com/username/vulnark.git
   cd vulnark
   ```

2. 配置环境变量（可选）
   ```bash
   # 复制示例配置文件
   cp .env.example .env
   
   # 根据需要编辑配置
   nano .env
   ```

3. 启动服务
   ```bash
   docker-compose up -d
   ```

4. 检查服务状态
   ```bash
   docker-compose ps
   ```

成功部署后，您可以通过 http://localhost:8080 访问VulnArk前端界面，通过 http://localhost:8000 访问API服务。

### 使用 Kubernetes

对于大规模部署或生产环境，推荐使用Kubernetes进行部署。

1. 获取Kubernetes配置文件
   ```bash
   git clone https://github.com/username/vulnark.git
   cd vulnark/kubernetes
   ```

2. 创建Kubernetes命名空间
   ```bash
   kubectl create namespace vulnark
   ```

3. 创建配置
   ```bash
   # 创建配置
   kubectl apply -f config.yaml -n vulnark
   
   # 创建密钥
   kubectl apply -f secrets.yaml -n vulnark
   ```

4. 部署应用
   ```bash
   kubectl apply -f frontend.yaml -n vulnark
   kubectl apply -f backend.yaml -n vulnark
   kubectl apply -f mongodb.yaml -n vulnark
   ```

5. 检查部署状态
   ```bash
   kubectl get pods -n vulnark
   ```

6. 设置Ingress或负载均衡器以暴露服务（根据您的K8s环境配置）

## 初始配置

### 数据库初始化

首次启动时，系统会自动初始化数据库并创建管理员账户：

```
用户名: admin@vulnark.com
密码: Admin@123
```

强烈建议首次登录后立即修改默认密码。

### 自定义配置

VulnArk提供多种配置选项，可通过环境变量进行设置：

1. 创建docker-compose.override.yml文件
   ```yaml
   version: '3.8'
   
   services:
     backend:
       environment:
         - JWT_SECRET=your-custom-secret-key
         - LOG_LEVEL=debug
     
     frontend:
       environment:
         - CUSTOM_TITLE=我的安全平台
   ```

2. 重启服务应用更改
   ```bash
   docker-compose down
   docker-compose up -d
   ```

## 首次登录

1. 访问 http://localhost:8080
2. 使用默认管理员账户登录
3. 完成以下初始设置：
   - 修改默认管理员密码
   - 设置系统基本参数
   - 创建部门和用户账户
   - 配置邮件通知（可选）

## 基本操作指南

### 用户管理

1. 导航至"用户管理"模块
2. 创建新用户，设置用户角色（管理员、普通用户或只读用户）
3. 分配部门和权限

### 漏洞管理

1. 导航至"漏洞管理"模块
2. 创建新漏洞记录或导入漏洞数据
3. 设置漏洞严重程度和修复期限
4. 分配责任人和跟踪修复进度

### 报告生成

1. 导航至"报告中心"
2. 选择报告类型（摘要、详细、合规或趋势分析）
3. 设置筛选条件和时间范围
4. 生成并下载报告

## 常见问题

### 数据库连接失败

**问题**: 系统无法连接到MongoDB数据库
**解决方案**:
1. 检查MongoDB服务是否正常运行:
   ```bash
   docker-compose ps mongodb
   ```
2. 检查数据库连接配置:
   ```bash
   docker-compose logs backend | grep "MongoDB"
   ```

### 无法登录系统

**问题**: 使用默认管理员账户无法登录
**解决方案**:
1. 重置管理员密码:
   ```bash
   docker-compose exec backend ./vulnark-api reset-admin
   ```
2. 检查日志获取更多信息:
   ```bash
   docker-compose logs frontend
   docker-compose logs backend
   ```

### 系统性能问题

**问题**: 系统响应缓慢
**解决方案**:
1. 检查资源使用情况:
   ```bash
   docker stats
   ```
2. 考虑增加资源配置:
   ```yaml
   # docker-compose.override.yml
   services:
     backend:
       deploy:
         resources:
           limits:
             cpus: '1'
             memory: 1G
   ```

如需更多帮助，请参阅[故障排除指南](TROUBLESHOOTING.md)或[提交问题](https://github.com/username/vulnark/issues)。 