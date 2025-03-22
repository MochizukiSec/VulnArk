# VulnArk Docker 部署指南

VulnArk是一个现代化的漏洞管理平台，本文档提供了使用Docker进行部署的详细指南。

## 系统要求

- Docker Engine 19.03.0+
- Docker Compose 1.29.0+
- 内存: 最小2GB, 推荐4GB+
- 存储: 最小10GB, 推荐20GB+

## 快速开始

1. 克隆仓库：

```bash
git clone https://github.com/yourusername/vulnark.git
cd vulnark
```

2. 执行部署脚本：

```bash
chmod +x deploy.sh
./deploy.sh
```

3. 打开浏览器访问：http://localhost

4. 使用默认管理员账号登录：
   - 用户名: admin
   - 密码: admin123

## 手动部署

如果您想手动部署，可以按照以下步骤操作：

1. 创建必要的目录：

```bash
mkdir -p logs uploads mysql/init
```

2. 编辑配置文件：

```bash
# 复制Docker环境配置文件
cp backend/config/config.docker.yaml backend/config/config.yaml
```

3. 启动服务：

```bash
docker-compose up -d
```

## 配置说明

### 环境变量

您可以在`docker-compose.yml`中修改以下环境变量：

- `DB_HOST`: 数据库主机名 (默认: mysql)
- `DB_PORT`: 数据库端口 (默认: 3306)
- `DB_USER`: 数据库用户名 (默认: vulnark)
- `DB_PASSWORD`: 数据库密码 (默认: vulnark_password)
- `DB_NAME`: 数据库名称 (默认: vulnark)

### 端口映射

- 前端服务: `80` -> 主机端口 `80`
- 数据库服务: `3306` -> 主机端口 `3306` (可选，用于直接连接数据库)

## 数据持久化

所有数据将通过以下方式保持持久化：

- MySQL数据存储在名为`mysql-data`的Docker卷中
- 日志文件映射到主机的`./logs`目录
- 上传文件映射到主机的`./uploads`目录
- 配置文件映射到主机的`./backend/config`目录

## 常用操作

- 查看容器状态：`docker-compose ps`
- 查看服务日志：`docker-compose logs -f [服务名]`
- 重启服务：`docker-compose restart [服务名]`
- 停止所有服务：`docker-compose stop`
- 启动所有服务：`docker-compose start`
- 重建服务：`docker-compose up -d --build [服务名]`
- 完全删除所有服务及数据：`docker-compose down -v`

## 自定义配置

### 自定义MySQL配置

修改`./mysql/my.cnf`文件后重启MySQL服务：

```bash
docker-compose restart mysql
```

### 自定义后端配置

修改`./backend/config/config.yaml`文件后重启后端服务：

```bash
docker-compose restart backend
```

## 故障排除

1. 如果无法访问前端页面，检查容器状态和日志：

```bash
docker-compose ps
docker-compose logs -f frontend
```

2. 如果API请求失败，检查后端日志：

```bash
docker-compose logs -f backend
```

3. 如果数据库连接失败，检查MySQL日志和配置：

```bash
docker-compose logs -f mysql
```

## 安全建议

1. 修改默认的用户名和密码
2. 修改`config.yaml`中的JWT密钥
3. 在生产环境中使用HTTPS
4. 限制MySQL端口仅供内部访问

## 备份数据

备份数据库：

```bash
docker exec vulnark-mysql mysqldump -u root -p vulnark > backup_$(date +%F).sql
``` 