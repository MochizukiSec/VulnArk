# 使用輕量級鏡像
FROM alpine:latest

# 設置工作目錄
WORKDIR /app

# 創建一個簡單的HTML文件
RUN echo '<!DOCTYPE html>\
<html>\
<head>\
    <meta charset="utf-8">\
    <title>VulnArk - 安全漏洞管理平台</title>\
    <style>\
        body { font-family: Arial, sans-serif; text-align: center; padding: 50px; }\
        h1 { color: #2c3e50; }\
        p { color: #7f8c8d; }\
    </style>\
</head>\
<body>\
    <h1>VulnArk - 安全漏洞管理平台</h1>\
    <p>這是一個模擬的前端頁面。在實際部署中，這裡會顯示真正的VulnArk前端界面。</p>\
    <p>當前Docker配置僅供參考，不包含實際的功能代碼。</p>\
</body>\
</html>' > /app/index.html

# 創建一個簡單的腳本來模擬web服務
RUN echo '#!/bin/sh\necho "啟動模擬的前端服務..."\necho "訪問 http://localhost:80 (實際環境中)"' > /app/start.sh && \
    chmod +x /app/start.sh

# 暴露端口
EXPOSE 80

# 啟動腳本
CMD ["/bin/sh", "-c", "/app/start.sh && tail -f /dev/null"] 