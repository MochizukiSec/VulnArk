# 使用輕量級nginx鏡像
FROM nginx:stable-alpine

# 設置工作目錄
WORKDIR /usr/share/nginx/html

# 創建一個簡單的靜態頁面來模擬前端服務
RUN echo '<!DOCTYPE html>\
<html>\
<head>\
    <meta charset="utf-8">\
    <title>VulnArk - 安全漏洞管理平台</title>\
    <style>\
        body {\
            font-family: Arial, sans-serif;\
            margin: 0;\
            padding: 0;\
            display: flex;\
            justify-content: center;\
            align-items: center;\
            min-height: 100vh;\
            background-color: #f5f5f5;\
        }\
        .container {\
            max-width: 800px;\
            padding: 30px;\
            background-color: white;\
            border-radius: 10px;\
            box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);\
            text-align: center;\
        }\
        h1 {\
            color: #2c3e50;\
        }\
        p {\
            color: #7f8c8d;\
            line-height: 1.6;\
        }\
    </style>\
</head>\
<body>\
    <div class="container">\
        <h1>VulnArk - 安全漏洞管理平台</h1>\
        <p>這是一個模擬的前端頁面。在實際部署中，這裡會顯示真正的VulnArk前端界面。</p>\
        <p>當前Docker配置僅供參考，不包含實際的功能代碼。</p>\
    </div>\
</body>\
</html>' > /usr/share/nginx/html/index.html

# 拷貝nginx配置文件（如果需要）
COPY frontend/nginx.conf /etc/nginx/conf.d/default.conf

# 暴露端口
EXPOSE 80

# 設置健康檢查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget -q -O - http://localhost/ || exit 1

# 啟動Nginx
CMD ["nginx", "-g", "daemon off;"] 