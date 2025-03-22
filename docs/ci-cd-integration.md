# VulnArk CI/CD集成指南

本指南将帮助您将VulnArk漏洞管理平台与您的CI/CD流水线集成，实现自动化的安全扫描和漏洞管理。

## 目录

- [概述](#概述)
- [支持的CI/CD平台](#支持的cicd平台)
- [集成步骤](#集成步骤)
- [配置示例](#配置示例)
- [自定义数据格式](#自定义数据格式)
- [常见问题](#常见问题)

## 概述

VulnArk CI/CD集成允许您在持续集成和持续部署流程中自动执行安全扫描，并将发现的漏洞直接上报到VulnArk平台进行统一管理。

通过集成，您可以：

- 在开发早期发现并修复安全问题
- 将安全扫描结果集中到VulnArk中进行追踪和管理
- 对安全漏洞设置阻断规则，阻止有严重漏洞的代码进入生产环境
- 生成安全合规报告，满足合规要求

## 支持的CI/CD平台

VulnArk目前支持以下CI/CD平台的集成：

- Jenkins
- GitLab CI/CD
- GitHub Actions
- 自定义平台（通过API集成）

## 集成步骤

### 1. 创建集成配置

1. 登录VulnArk平台，进入【设置】>【集成管理】页面
2. 点击【添加集成】按钮
3. 选择集成类型（Jenkins、GitLab、GitHub或自定义）
4. 填写名称、描述等基本信息
5. 点击【创建集成】完成配置

### 2. 获取API密钥

1. 在集成列表中，找到您刚创建的集成配置
2. 点击【API密钥】按钮
3. 点击【重新生成】按钮生成API密钥
4. 复制并安全保存此API密钥（注意：密钥只会显示一次！）

### 3. 配置CI/CD平台

1. 在您的CI/CD平台中，添加VulnArk API密钥作为安全变量
   - Jenkins：在凭据管理中添加
   - GitLab：在项目设置的CI/CD变量中添加
   - GitHub：在仓库的Secrets中添加

2. 配置扫描工具，确保它们能够生成JSON格式的报告

3. 添加将扫描结果发送到VulnArk的步骤（参见[配置示例](#配置示例)）

## 配置示例

### Jenkins Pipeline

```groovy
pipeline {
    agent any

    environment {
        VULNARK_API_ENDPOINT = 'https://your-vulnark-instance.com'
        VULNARK_API_KEY = credentials('vulnark-api-key')
        VULNARK_INTEGRATION_TYPE = 'jenkins'
    }

    stages {
        stage('Security Scan') {
            steps {
                // 执行扫描并生成报告
                sh 'npm audit --json > npm-audit.json'
                
                // 发送结果到VulnArk
                sh '''
                    curl -X POST \
                        ${VULNARK_API_ENDPOINT}/api/v1/webhooks/${VULNARK_INTEGRATION_TYPE} \
                        -H "Content-Type: application/json" \
                        -H "X-API-Key: ${VULNARK_API_KEY}" \
                        -d @npm-audit.json
                '''
            }
        }
    }
}
```

### GitLab CI

```yaml
security_scan:
  stage: test
  script:
    # 执行扫描并生成报告
    - npm audit --json > npm-audit.json
    # 发送结果到VulnArk
    - |
      curl -X POST \
        ${VULNARK_API_ENDPOINT}/api/v1/webhooks/gitlab \
        -H "Content-Type: application/json" \
        -H "X-API-Key: ${VULNARK_API_KEY}" \
        -d @npm-audit.json
  variables:
    VULNARK_API_ENDPOINT: https://your-vulnark-instance.com
    # VULNARK_API_KEY在GitLab设置中配置为安全变量
```

### GitHub Actions

```yaml
- name: 安全扫描
  run: npm audit --json > npm-audit.json

- name: 发送结果到VulnArk
  run: |
    curl -X POST \
      ${{ secrets.VULNARK_API_ENDPOINT }}/api/v1/webhooks/github \
      -H "Content-Type: application/json" \
      -H "X-API-Key: ${{ secrets.VULNARK_API_KEY }}" \
      -d @npm-audit.json
```

## 自定义数据格式

VulnArk接受以下JSON格式的漏洞数据：

```json
{
  "findings": [
    {
      "title": "漏洞标题",
      "severity": "high",  // 严重程度: critical, high, medium, low, info
      "description": "漏洞描述",
      "cve_id": "CVE-2023-XXXX",  // 可选
      "references": "https://example.com/ref",  // 可选
      "location": "src/example.js:42"  // 可选
    }
  ]
}
```

如果您的扫描工具生成的格式与此不符，您需要编写转换脚本将其转换为VulnArk接受的格式。请参考示例中的转换脚本。

## 常见问题

### Q: 集成配置后无法接收扫描结果

**A:** 请检查以下几点：

1. API密钥是否正确配置
2. 网络连接是否正常（CI/CD服务器是否能访问VulnArk服务器）
3. 数据格式是否符合要求
4. 查看CI/CD日志中的错误信息

### Q: 支持哪些安全扫描工具？

**A:** VulnArk可以集成几乎所有能生成结构化输出（如JSON）的安全扫描工具，包括但不限于：

- NPM Audit
- OWASP Dependency Check
- SonarQube
- ESLint安全规则
- OWASP ZAP
- Snyk
- Trivy
- 等等...

### Q: 如何处理误报？

**A:** 将扫描结果导入VulnArk后，您可以在平台中将误报标记为"误报"状态。这些状态会被记录下来，将来相同的漏洞导入时会被自动标记。

### Q: 如何自定义结果处理逻辑？

**A:** 您可以创建自定义类型的集成，然后使用VulnArk API编写自己的处理逻辑。请参考[VulnArk API文档](https://docs.vulnark.example.com/api)了解更多信息。

---

如需更多帮助，请联系VulnArk支持团队。 