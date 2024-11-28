---
title: "node-docker-build"
date: 2024-11-28T11:19:43+08:00
lastmod: 2024-11-28T11:19:43+08:00
categories: ['Docker']
tags: ['']
keywords: node docker build
description: 使用 Docker 来构建前端项目
---

## Dockerfile

```dockerfile
FROM node:16-buster-slim AS modules
COPY package.json package-lock.json /app/
WORKDIR /app
ENV LC_ALL=en_US.UTF-8
# 将老的镜像地址替换为新的镜像地址并安装指定版本的依赖
RUN sed -i 's|http://registry.npm.taobao.org|https://registry.npmmirror.com|g' package-lock.json && \
    sed -i 's|https://registry.npm.taobao.org|https://registry.npmmirror.com|g' package-lock.json && \
    sed -i 's|https://registry.npmjs.org|https://registry.npmmirror.com|g' package-lock.json && \
    sed -i 's|https://registry.nlark.com|https://registry.npmmirror.com|g' package-lock.json && \
    npm config set registry http://registry.npmmirror.com && \
    npm install --production && \
    npm install --silent chalk@1.1.3 && \
    npm install --silent vue-template-compiler@2.7.16

FROM node:16-buster-slim AS bulider
WORKDIR /app
COPY . .
ENV LC_ALL=en_US.UTF-8

COPY --from=modules /app/node_modules /app/node_modules
RUN sed -i 's|deb.debian.org|mirrors.ustc.edu.cn|g' /etc/apt/sources.list && \
    apt update && apt install -y zip --no-install-recommends && \
    npm run test && \
    zip -r dist.zip ./dist/

FROM alpine
COPY --from=bulider /app/dist.zip /dist.zip
CMD ["echo", "use 'docker cp container:/dist.zip ./dist.zip'"]
```

## 使用

```shell
# 构建
docker build -t vue-project .
# 运行
docker run --name vue-project-container vue-project
# 取出 dist.zip 到当前目录
docker cp vue-project-container:/dist.zip ./dist.zip
# 移除容器
docker rm vue-project-container
# 移除镜像
docker rmi vue-project
```
