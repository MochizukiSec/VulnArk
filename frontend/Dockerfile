FROM nginx:1.15.3-alpine as production-stage

MAINTAINER Vuln-Management

COPY ./nginx.conf /etc/nginx/nginx.conf

COPY ./dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]