#!/bin/bash

echo "开始初始化数据库..."
echo "创建管理员账户..."

# 编译并运行初始化脚本
cd "$(dirname "$0")/.." && go run scripts/init_admin.go

echo "数据库初始化完成!" 