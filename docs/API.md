# VulnArk API 参考文档

本文档提供了VulnArk平台的API接口详细说明，便于系统集成和自动化。所有API均遵循RESTful设计规范，使用JWT进行认证。

## API 基础信息

- **基础URL**: `http://your-vulnark-server:8000/api`
- **API版本**: v1
- **内容类型**: `application/json`
- **认证方式**: JWT Token (Bearer Authentication)

## 目录

- [认证相关](#认证相关)
- [用户管理](#用户管理)
- [漏洞管理](#漏洞管理)
- [报告生成](#报告生成)
- [系统设置](#系统设置)
- [错误处理](#错误处理)

## 认证相关

### 用户登录

获取JWT认证令牌。

- **URL**: `/auth/login`
- **方法**: `POST`
- **认证要求**: 无

**请求参数**:

```json
{
  "email": "user@example.com",
  "password": "YourPassword123"
}
```

**成功响应** (200 OK):

```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "61fb1e7b9c6c3a5da2b0e624",
      "email": "user@example.com",
      "firstName": "John",
      "lastName": "Doe",
      "role": "admin",
      "status": "active",
      "department": "Security"
    }
  }
}
```

**错误响应** (401 Unauthorized):

```json
{
  "success": false,
  "error": "invalid_credentials",
  "message": "邮箱或密码不正确"
}
```

### 刷新令牌

刷新JWT令牌的有效期。

- **URL**: `/auth/refresh-token`
- **方法**: `POST`
- **认证要求**: 有效的JWT令牌

**成功响应** (200 OK):

```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

### 用户注销

使当前JWT令牌失效。

- **URL**: `/auth/logout`
- **方法**: `POST`
- **认证要求**: 有效的JWT令牌

**成功响应** (200 OK):

```json
{
  "success": true,
  "message": "成功注销"
}
```

## 用户管理

### 获取所有用户

获取所有用户的列表。

- **URL**: `/users`
- **方法**: `GET`
- **认证要求**: 有效的JWT令牌
- **权限要求**: 管理员角色

**查询参数**:

| 参数 | 类型 | 描述 |
|------|------|------|
| page | 整数 | 分页页码，默认为1 |
| limit | 整数 | 每页记录数，默认为20 |
| search | 字符串 | 搜索关键词 |
| role | 字符串 | 按角色筛选 |
| status | 字符串 | 按状态筛选 |

**成功响应** (200 OK):

```json
{
  "success": true,
  "data": {
    "users": [
      {
        "id": "61fb1e7b9c6c3a5da2b0e624",
        "email": "user@example.com",
        "firstName": "John",
        "lastName": "Doe",
        "role": "admin",
        "status": "active",
        "department": "Security",
        "createdAt": "2022-02-03T10:15:23Z",
        "updatedAt": "2022-02-03T10:15:23Z"
      }
      // 更多用户...
    ],
    "pagination": {
      "currentPage": 1,
      "totalPages": 5,
      "totalItems": 98,
      "itemsPerPage": 20
    }
  }
}
```

### 创建新用户

创建新用户账户。

- **URL**: `/users`
- **方法**: `POST`
- **认证要求**: 有效的JWT令牌
- **权限要求**: 管理员角色

**请求参数**:

```json
{
  "email": "newuser@example.com",
  "password": "SecurePassword123",
  "firstName": "Jane",
  "lastName": "Smith",
  "role": "regular",
  "status": "active",
  "department": "IT",
  "profilePic": "base64encodedimage" // 可选
}
```

**成功响应** (201 Created):

```json
{
  "success": true,
  "data": {
    "id": "61fb1e7b9c6c3a5da2b0e625",
    "email": "newuser@example.com"
  },
  "message": "用户创建成功"
}
```

### 获取单个用户

获取单个用户的详细信息。

- **URL**: `/users/{userId}`
- **方法**: `GET`
- **认证要求**: 有效的JWT令牌
- **权限要求**: 管理员角色或查询自己的信息

**成功响应** (200 OK):

```json
{
  "success": true,
  "data": {
    "id": "61fb1e7b9c6c3a5da2b0e624",
    "email": "user@example.com",
    "firstName": "John",
    "lastName": "Doe",
    "role": "admin",
    "status": "active",
    "department": "Security",
    "profilePic": "https://example.com/profiles/image.jpg",
    "createdAt": "2022-02-03T10:15:23Z",
    "updatedAt": "2022-02-03T10:15:23Z"
  }
}
```

### 更新用户

更新用户信息。

- **URL**: `/users/{userId}`
- **方法**: `PUT`
- **认证要求**: 有效的JWT令牌
- **权限要求**: 管理员角色或更新自己的信息

**请求参数**:

```json
{
  "firstName": "John",
  "lastName": "Updated",
  "department": "IT Security",
  "role": "regular", // 仅管理员可更改
  "status": "active", // 仅管理员可更改
  "newPassword": "NewSecurePassword123" // 可选
}
```

**成功响应** (200 OK):

```json
{
  "success": true,
  "message": "用户信息更新成功"
}
```

### 删除用户

删除用户账户。

- **URL**: `/users/{userId}`
- **方法**: `DELETE`
- **认证要求**: 有效的JWT令牌
- **权限要求**: 管理员角色

**成功响应** (200 OK):

```json
{
  "success": true,
  "message": "用户删除成功"
}
```

## 漏洞管理

### 获取漏洞列表

获取漏洞列表。

- **URL**: `/vulnerabilities`
- **方法**: `GET`
- **认证要求**: 有效的JWT令牌

**查询参数**:

| 参数 | 类型 | 描述 |
|------|------|------|
| page | 整数 | 分页页码，默认为1 |
| limit | 整数 | 每页记录数，默认为20 |
| search | 字符串 | 搜索关键词 |
| severity | 字符串 | 按严重程度筛选 |
| status | 字符串 | 按状态筛选 |
| startDate | 日期 | 开始日期筛选 |
| endDate | 日期 | 结束日期筛选 |

**成功响应** (200 OK):

```json
{
  "success": true,
  "data": {
    "vulnerabilities": [
      {
        "id": "61fb1e7b9c6c3a5da2b0e626",
        "title": "SQL注入漏洞",
        "description": "产品存在SQL注入漏洞，可能导致数据泄露...",
        "cvssScore": 8.5,
        "severity": "high",
        "status": "open",
        "discoveredDate": "2022-02-05T10:15:23Z",
        "system": "CRM系统",
        "assignedTo": "61fb1e7b9c6c3a5da2b0e624",
        "createdBy": "61fb1e7b9c6c3a5da2b0e625"
      }
      // 更多漏洞...
    ],
    "pagination": {
      "currentPage": 1,
      "totalPages": 10,
      "totalItems": 198,
      "itemsPerPage": 20
    }
  }
}
```

### 创建漏洞记录

创建新的漏洞记录。

- **URL**: `/vulnerabilities`
- **方法**: `POST`
- **认证要求**: 有效的JWT令牌
- **权限要求**: 管理员或普通用户角色

**请求参数**:

```json
{
  "title": "XSS漏洞",
  "description": "在用户输入表单中存在跨站脚本漏洞",
  "cvssScore": 6.5,
  "severity": "medium",
  "status": "open",
  "discoveredDate": "2022-02-10T15:30:00Z",
  "system": "Web应用",
  "assignedTo": "61fb1e7b9c6c3a5da2b0e624",
  "attachments": [
    {
      "name": "poc.txt",
      "content": "base64encodedcontent"
    }
  ]
}
```

**成功响应** (201 Created):

```json
{
  "success": true,
  "data": {
    "id": "61fb1e7b9c6c3a5da2b0e627"
  },
  "message": "漏洞记录创建成功"
}
```

### 获取漏洞详情

获取单个漏洞的详细信息。

- **URL**: `/vulnerabilities/{vulnerabilityId}`
- **方法**: `GET`
- **认证要求**: 有效的JWT令牌

**成功响应** (200 OK):

```json
{
  "success": true,
  "data": {
    "id": "61fb1e7b9c6c3a5da2b0e626",
    "title": "SQL注入漏洞",
    "description": "产品存在SQL注入漏洞，可能导致数据泄露...",
    "cvssScore": 8.5,
    "severity": "high",
    "status": "open",
    "discoveredDate": "2022-02-05T10:15:23Z",
    "system": "CRM系统",
    "assignedTo": {
      "id": "61fb1e7b9c6c3a5da2b0e624",
      "firstName": "John",
      "lastName": "Doe",
      "email": "john.doe@example.com"
    },
    "createdBy": {
      "id": "61fb1e7b9c6c3a5da2b0e625",
      "firstName": "Jane",
      "lastName": "Smith",
      "email": "jane.smith@example.com"
    },
    "attachments": [
      {
        "id": "61fb1e7b9c6c3a5da2b0e628",
        "name": "evidence.png",
        "url": "https://example.com/attachments/evidence.png",
        "uploadedAt": "2022-02-05T10:20:23Z"
      }
    ],
    "history": [
      {
        "action": "created",
        "timestamp": "2022-02-05T10:15:23Z",
        "user": "Jane Smith"
      },
      {
        "action": "status_changed",
        "oldValue": "new",
        "newValue": "open",
        "timestamp": "2022-02-05T10:30:23Z",
        "user": "John Doe"
      }
    ],
    "createdAt": "2022-02-05T10:15:23Z",
    "updatedAt": "2022-02-05T10:30:23Z"
  }
}
```

## 报告生成

### 生成报告

生成安全漏洞报告。

- **URL**: `/reports/generate`
- **方法**: `POST`
- **认证要求**: 有效的JWT令牌
- **权限要求**: 管理员或普通用户角色

**请求参数**:

```json
{
  "reportType": "summary", // summary, detailed, compliance, trend
  "format": "pdf", // pdf, excel, word, html
  "filters": {
    "startDate": "2022-01-01T00:00:00Z",
    "endDate": "2022-03-01T00:00:00Z",
    "severity": ["high", "critical"],
    "status": ["open", "in_progress"],
    "systems": ["Web应用", "CRM系统"]
  },
  "title": "第一季度安全漏洞报告",
  "includeCharts": true,
  "includeRemediation": true
}
```

**成功响应** (200 OK):

```json
{
  "success": true,
  "data": {
    "reportId": "61fb1e7b9c6c3a5da2b0e629",
    "downloadUrl": "https://example.com/reports/61fb1e7b9c6c3a5da2b0e629.pdf",
    "expiresAt": "2022-03-10T10:15:23Z"
  }
}
```

## 系统设置

### 获取系统设置

获取系统配置设置。

- **URL**: `/settings`
- **方法**: `GET`
- **认证要求**: 有效的JWT令牌
- **权限要求**: 管理员角色

**成功响应** (200 OK):

```json
{
  "success": true,
  "data": {
    "emailNotifications": true,
    "defaultReportFormat": "pdf",
    "severityLevels": ["info", "low", "medium", "high", "critical"],
    "statusLevels": ["new", "open", "in_progress", "resolved", "closed"],
    "autoAssignment": false,
    "retentionPeriod": 365
  }
}
```

### 更新系统设置

更新系统配置设置。

- **URL**: `/settings`
- **方法**: `PUT`
- **认证要求**: 有效的JWT令牌
- **权限要求**: 管理员角色

**请求参数**:

```json
{
  "emailNotifications": true,
  "defaultReportFormat": "pdf",
  "autoAssignment": true,
  "retentionPeriod": 730
}
```

**成功响应** (200 OK):

```json
{
  "success": true,
  "message": "系统设置更新成功"
}
```

## 错误处理

所有API错误响应遵循统一格式：

```json
{
  "success": false,
  "error": "error_code",
  "message": "人类可读的错误消息",
  "details": {} // 可选的详细错误信息
}
```

### 常见错误代码

| 错误代码 | HTTP状态码 | 描述 |
|----------|------------|------|
| invalid_credentials | 401 | 认证凭据无效 |
| token_expired | 401 | 令牌已过期 |
| token_invalid | 401 | 令牌无效 |
| permission_denied | 403 | 没有执行操作的权限 |
| resource_not_found | 404 | 请求的资源不存在 |
| validation_error | 422 | 请求数据验证失败 |
| rate_limit_exceeded | 429 | 超过API请求限制 |
| internal_error | 500 | 服务器内部错误 | 