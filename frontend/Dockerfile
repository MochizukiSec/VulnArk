# 构建阶段
FROM node:16-alpine as build-stage

# 设置工作目录
WORKDIR /app

# 拷贝package.json和package-lock.json (或 yarn.lock)
COPY frontend/package*.json ./

# 安装依赖
RUN npm install

# 拷贝前端源代码
COPY frontend/ ./

# 构建生产版本
RUN npm run build

# 生产阶段
FROM nginx:stable-alpine as production-stage

# 从构建阶段拷贝构建结果到Nginx默认目录
COPY --from=build-stage /app/dist /usr/share/nginx/html

# 拷贝自定义Nginx配置
COPY frontend/nginx.conf /etc/nginx/conf.d/default.conf

# 暴露端口
EXPOSE 80

# 设置健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost/ || exit 1

# 启动Nginx
CMD ["nginx", "-g", "daemon off;"] 