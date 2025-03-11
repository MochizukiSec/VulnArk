# 使用官方Go镜像作为构建环境
FROM golang:1.20-alpine AS builder

# 设置工作目录
WORKDIR /app

# 拷贝go.mod和go.sum文件
COPY backend/go.mod backend/go.sum ./

# 下载依赖
RUN go mod download

# 拷贝源代码
COPY backend/ ./

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -o vuln-management .

# 使用轻量级alpine镜像作为运行环境
FROM alpine:latest

# 安装ca-certificates，用于HTTPS请求
RUN apk --no-cache add ca-certificates

WORKDIR /app

# 从构建阶段拷贝编译好的二进制文件
COPY --from=builder /app/vuln-management .

# 创建配置目录并设置环境变量
COPY backend/config/ ./config/

# 设置时区
ENV TZ=Asia/Shanghai

# 暴露端口
EXPOSE 8000

# 设置健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget -q -O - http://localhost:8000/api/health || exit 1

# 运行应用
CMD ["./vuln-management"] 