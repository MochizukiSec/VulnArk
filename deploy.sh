#!/bin/bash

# 设置颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== VulnArk 漏洞管理平台部署脚本 ===${NC}"
echo -e "${YELLOW}本脚本将帮助您部署VulnArk系统的所有组件${NC}"
echo

# 检查Docker是否已安装
echo -e "${YELLOW}正在检查Docker...${NC}"
if command -v docker &> /dev/null; then
    echo -e "${GREEN}Docker已安装√${NC}"
else
    echo -e "${RED}未检测到Docker，请先安装Docker和Docker Compose${NC}"
    echo "可参考: https://docs.docker.com/get-docker/"
    exit 1
fi

# 检查Docker Compose是否已安装
echo -e "${YELLOW}正在检查Docker Compose...${NC}"
if command -v docker-compose &> /dev/null || docker compose version &> /dev/null; then
    echo -e "${GREEN}Docker Compose已安装√${NC}"
else
    echo -e "${RED}未检测到Docker Compose，请先安装Docker Compose${NC}"
    echo "可参考: https://docs.docker.com/compose/install/"
    exit 1
fi

# 确保MySQL初始化目录存在
echo -e "${YELLOW}创建必要的目录...${NC}"
mkdir -p ./logs ./uploads ./mysql/init

# 确保文件权限正确
echo -e "${YELLOW}设置文件权限...${NC}"
chmod +x ./backend/Dockerfile
chmod +x ./frontend/Dockerfile
chmod +x ./deploy.sh
chmod +x -R ./mysql

# 复制Docker配置文件
echo -e "${YELLOW}配置后端...${NC}"
echo -e "${GREEN}后端配置已准备就绪${NC}"

# 询问是否需要自定义端口配置
echo -e "${YELLOW}是否需要自定义端口配置? (y/n)${NC}"
read -r customize_ports

if [[ "$customize_ports" == "y" ]]; then
    echo -e "${YELLOW}请输入前端Web端口 (默认: 80):${NC}"
    read -r web_port
    web_port=${web_port:-80}
    
    echo -e "${YELLOW}请输入后端API端口 (默认: 8080):${NC}"
    read -r api_port
    api_port=${api_port:-8080}
    
    echo -e "${YELLOW}请输入数据库端口 (默认: 3306):${NC}"
    read -r db_port
    db_port=${db_port:-3306}
    
    # 修改docker-compose.yml中的端口配置
    sed -i "s/- \"80:80\"/- \"$web_port:80\"/" docker-compose.yml
    sed -i "s/- \"8080:8080\"/- \"$api_port:8080\"/" docker-compose.yml
    sed -i "s/- \"3306:3306\"/- \"$db_port:3306\"/" docker-compose.yml
    
    echo -e "${GREEN}端口配置已更新${NC}"
else
    echo -e "${GREEN}将使用默认端口配置:${NC}"
    echo "  前端Web端口: 80"
    echo "  后端API端口: 8080"
    echo "  数据库端口: 3306"
fi

# 询问是否使用默认配置或自定义配置
echo -e "${YELLOW}是否要自定义数据库配置? (y/n)${NC}"
read -r customize_db

if [[ "$customize_db" == "y" ]]; then
    echo -e "${YELLOW}请输入MySQL root密码:${NC}"
    read -r db_root_password
    
    echo -e "${YELLOW}请输入VulnArk数据库用户名:${NC}"
    read -r db_user
    
    echo -e "${YELLOW}请输入VulnArk数据库密码:${NC}"
    read -r db_password
    
    echo -e "${YELLOW}请输入VulnArk数据库名:${NC}"
    read -r db_name
    
    # 修改docker-compose.yml中的数据库配置
    sed -i "s/MYSQL_ROOT_PASSWORD=root_password/MYSQL_ROOT_PASSWORD=$db_root_password/" docker-compose.yml
    sed -i "s/MYSQL_DATABASE=vulnark/MYSQL_DATABASE=$db_name/" docker-compose.yml
    sed -i "s/MYSQL_USER=vulnark/MYSQL_USER=$db_user/" docker-compose.yml
    sed -i "s/MYSQL_PASSWORD=vulnark_password/MYSQL_PASSWORD=$db_password/" docker-compose.yml
    
    # 同时更新后端服务的环境变量
    sed -i "s/DB_USER=vulnark/DB_USER=$db_user/" docker-compose.yml
    sed -i "s/DB_PASSWORD=vulnark_password/DB_PASSWORD=$db_password/" docker-compose.yml
    sed -i "s/DB_NAME=vulnark/DB_NAME=$db_name/" docker-compose.yml
    
    echo -e "${GREEN}数据库配置已更新${NC}"
else
    echo -e "${GREEN}将使用默认数据库配置:${NC}"
    echo "  Root密码: root_password"
    echo "  数据库名: vulnark"
    echo "  用户名: vulnark"
    echo "  密码: vulnark_password"
fi

# 开始构建和部署
echo -e "${YELLOW}开始构建和部署VulnArk...${NC}"
echo -e "${YELLOW}这可能需要一些时间，请耐心等待...${NC}"

# 使用docker-compose启动服务
docker-compose up -d --build

# 检查是否部署成功
if [ $? -eq 0 ]; then
    echo -e "${GREEN}VulnArk已成功部署!${NC}"
    
    # 获取本机IP地址
    if command -v ip &> /dev/null; then
        host_ip=$(ip route get 1 | sed -n 's/^.*src \([0-9.]*\) .*$/\1/p')
    elif command -v ifconfig &> /dev/null; then
        host_ip=$(ifconfig | grep -Eo 'inet (addr:)?([0-9]*\.){3}[0-9]*' | grep -Eo '([0-9]*\.){3}[0-9]*' | grep -v '127.0.0.1' | head -n 1)
    else
        host_ip="localhost"
    fi
    
    web_port=$(grep -o '"[0-9]*:80"' docker-compose.yml | grep -o '[0-9]*' || echo "80")
    api_port=$(grep -o '"[0-9]*:8080"' docker-compose.yml | grep -o '[0-9]*' || echo "8080")
    
    echo -e "${GREEN}前端访问地址: http://$host_ip:$web_port${NC}"
    echo -e "${GREEN}API访问地址: http://$host_ip:$web_port/api${NC}"
    echo -e "${YELLOW}默认管理员账号:${NC}"
    echo "  用户名: admin"
    echo "  密码: admin123"
    echo
    echo -e "${YELLOW}您可以使用以下命令来管理VulnArk:${NC}"
    echo "  查看容器状态: docker-compose ps"
    echo "  查看后端日志: docker-compose logs -f backend"
    echo "  查看前端日志: docker-compose logs -f frontend"
    echo "  重启服务: docker-compose restart"
    echo "  停止服务: docker-compose stop"
    echo "  启动服务: docker-compose start"
    echo "  完全删除: docker-compose down -v"
else
    echo -e "${RED}部署失败，请检查以上错误信息${NC}"
    echo -e "${YELLOW}您可以尝试手动运行: docker-compose up --build${NC}"
fi 