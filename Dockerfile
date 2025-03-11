# 使用阿里雲鏡像源的 Golang 鏡像作為構建環境
FROM registry.cn-hangzhou.aliyuncs.com/google_containers/golang:1.20-alpine AS builder

# 設置工作目錄
WORKDIR /app

# 拷貝go.mod和go.sum文件
COPY backend/go.mod backend/go.sum ./

# 下載依賴
RUN go mod download

# 拷貝源代碼
COPY backend/ ./

# 構建應用
RUN CGO_ENABLED=0 GOOS=linux go build -o vuln-management .

# 使用阿里雲鏡像源的輕量級alpine鏡像作為運行環境
FROM registry.cn-hangzhou.aliyuncs.com/google_containers/alpine:3.14

# 安裝ca-certificates，用於HTTPS請求
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk --no-cache add ca-certificates

WORKDIR /app

# 從構建階段拷貝編譯好的二進制文件
COPY --from=builder /app/vuln-management .

# 創建配置目錄並設置環境變量
COPY backend/config/ ./config/

# 設置時區
ENV TZ=Asia/Shanghai

# 暴露端口
EXPOSE 8000

# 設置健康檢查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget -q -O - http://localhost:8000/api/health || exit 1

# 運行應用
CMD ["./vuln-management"] 