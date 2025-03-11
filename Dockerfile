# 使用輕量級alpine鏡像作為運行環境
FROM alpine:latest

# 安裝ca-certificates，用於HTTPS請求
RUN apk --no-cache add ca-certificates

WORKDIR /app

# 設置時區
ENV TZ=Asia/Shanghai

# 暴露端口
EXPOSE 8000

# 設置健康檢查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget -q -O - http://localhost:8000/api/health || exit 1

# 指定運行命令
CMD echo "VulnArk後端服務" && \
    echo "請使用預構建鏡像運行此服務" && \
    echo "在一個真實環境中，這裡應該運行真正的後端服務" && \
    tail -f /dev/null 