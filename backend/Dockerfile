FROM alpine:latest

ENV TZ=Asia/Shanghai

# 安装必要的依赖
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY ./vulnark .

EXPOSE 8000

CMD ["./vulnark"]