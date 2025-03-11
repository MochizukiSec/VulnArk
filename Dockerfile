# 使用輕量級鏡像
FROM alpine:latest

# 設置工作目錄
WORKDIR /app

# 創建一個簡單的文件來模擬後端服務
RUN echo 'echo "VulnArk後端服務模擬"' > /app/start.sh && \
    chmod +x /app/start.sh

# 設置時區（不使用apk安裝）
ENV TZ=Asia/Shanghai

# 暴露端口
EXPOSE 8000

# 啟動腳本（不需要健康檢查，直接保持容器運行）
CMD ["/bin/sh", "-c", "/app/start.sh && tail -f /dev/null"] 