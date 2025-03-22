#!/bin/bash

# 清理测试数据脚本

echo "开始清理测试数据..."

# 验证环境变量或设置默认值
DB_USER=${DB_USER:-"root"}
DB_PASS=${DB_PASS:-"password"}
DB_HOST=${DB_HOST:-"localhost"}
DB_PORT=${DB_PORT:-"3306"}
DB_NAME=${DB_NAME:-"vulnark"}

# 执行SQL脚本
mysql -u"$DB_USER" -p"$DB_PASS" -h"$DB_HOST" -P"$DB_PORT" "$DB_NAME" < clean_test_data.sql

if [ $? -eq 0 ]; then
    echo "测试数据清理成功！"
else
    echo "测试数据清理失败，请检查数据库连接参数和权限。"
    exit 1
fi

echo "数据清理完成。" 