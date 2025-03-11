# 構建階段
FROM registry.cn-hangzhou.aliyuncs.com/google_containers/node:16-alpine as build-stage

# 設置工作目錄
WORKDIR /app

# 拷貝package.json和package-lock.json (或 yarn.lock)
COPY frontend/package*.json ./

# 設置npm鏡像源為淘寶鏡像
RUN npm config set registry https://registry.npmmirror.com

# 安裝依賴
RUN npm install

# 拷貝前端源代碼
COPY frontend/ ./

# 構建生產版本
RUN npm run build

# 生產階段
FROM registry.cn-hangzhou.aliyuncs.com/google_containers/nginx:stable-alpine as production-stage

# 從構建階段拷貝構建結果到Nginx默認目錄
COPY --from=build-stage /app/dist /usr/share/nginx/html

# 拷貝自定義Nginx配置
COPY frontend/nginx.conf /etc/nginx/conf.d/default.conf

# 暴露端口
EXPOSE 80

# 設置健康檢查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget -q -O - http://localhost/ || exit 1

# 啟動Nginx
CMD ["nginx", "-g", "daemon off;"] 