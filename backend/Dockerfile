FROM golang:1.18-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装必要的依赖
RUN apk add --no-cache gcc musl-dev git

# 复制go.mod和go.sum文件并下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 编译应用
RUN CGO_ENABLED=1 GOOS=linux go build -a -o vulnark .

# 使用轻量级基础镜像
FROM alpine:3.16

# 安装必要的运行时依赖
RUN apk add --no-cache ca-certificates tzdata mysql-client

# 设置时区为亚洲/上海
ENV TZ=Asia/Shanghai

# 创建配置和数据目录
RUN mkdir -p /app/config /app/logs /app/uploads

# 从构建阶段复制编译好的二进制文件
COPY --from=builder /app/vulnark /app/

# 复制配置文件
COPY ./config/config.yaml /app/config/

# 设置工作目录
WORKDIR /app

# 暴露端口
EXPOSE 8080

# 启动应用
CMD ["./vulnark"] 