FROM nginx:1.15.6-alpine

RUN ["rm", "/etc/nginx/conf.d/default.conf"]

COPY updatehub.conf /etc/nginx/conf.d/
