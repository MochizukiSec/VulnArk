# VulnArk 管理员账户创建指南

如果您需要重新创建管理员账户，可以选择以下几种方法之一。默认创建的管理员账户信息如下：

- **用户名**: `admin`
- **密码**: `admin123`
- **角色**: `admin`（管理员）
- **邮箱**: `admin@vulnark.com`

## 方法一：使用脚本自动创建

这是最简单的方法，直接执行脚本即可：

```bash
# 进入后端目录
cd backend

# 赋予脚本执行权限
chmod +x create_admin.sh

# 运行脚本
./create_admin.sh
```

## 方法二：直接执行SQL

如果您有权限直接访问MySQL数据库，可以执行以下命令：

```bash
# 使用MySQL客户端
mysql -uroot -proot123456 vulnark < create_admin.sql

# 或者进入MySQL后执行
mysql -uroot -proot123456 vulnark
source create_admin.sql
```

## 方法三：使用Go程序创建

如果您的环境中有Go编译器，可以使用Go程序创建管理员账户：

```bash
# 进入后端目录
cd backend

# 编译并运行
go run cmd/create_admin/main.go
```

## 方法四：使用API创建（需要现有管理员账户）

如果您已经有一个可用的管理员账户，可以使用API创建新的管理员账户：

```bash
# 获取token
TOKEN=$(curl -s -X POST -H "Content-Type: application/json" -d '{"username":"已有管理员账户","password":"密码"}' http://localhost:8080/api/v1/auth/login | grep -o '"token":"[^"]*' | cut -d'"' -f4)

# 创建新管理员
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN" -d '{
  "username": "admin",
  "password": "admin123",
  "email": "admin@vulnark.com",
  "real_name": "系统管理员",
  "role": "admin",
  "active": true
}' http://localhost:8080/api/v1/admin/users
```

## 特殊账号：测试管理员

如果以上方法都不可行，系统内置了一个测试账号，可以用来紧急访问：

- **用户名**: `testadmin999`
- **密码**: `testpass123`

这个账号有特殊处理逻辑，即使密码验证机制出现问题，也可以登录。登录后，您可以使用管理员功能创建新的管理员账户。 