---
title: "devcontainer"
date: 2025-09-10T11:19:26
lastmod: 2025-09-23T17:01:26
categories: ['Docker']
keywords: devcontainer
description: Dev Container
---

**所有环境构建均在官方标准之上进行构建**

## 文档参考

[官网](https://containers.dev)

[模板列表](https://github.com/devcontainers/templates/tree/main/src)

## 构建步骤

**环境变量**

- GITHUB_USERNAME: GitHub 的用户名, 用于拉取 dotfiles 仓库
- GITHUB_PAT: GitHub dotfiles 仓库的 [PAT](https://github.com/settings/personal-access-tokens), 用于拉取 dotfiles 仓库. 后续解密并导入 **gpg** 密钥文件也需要
- HTTP_PROXY/HTTPS_PROXY: 代理, 用于安装和更新 [chezmoi](https://chezmoi.io)

**构建命令**

1. 修改 APT 源为阿里云并更新
2. 修改时区为上海
3. 安装 [chezmoi](https://chezmoi.io)
4. 切换到 vscode 用户并初始化 dotfiles 仓库

> 后续的软件安装及环境配置均由 dotfiles 中的脚本完成

## 构建步骤(以 Java 为例)

**Dockerfile**

```Dockerfile
FROM mcr.microsoft.com/devcontainers/java:8-bookworm
LABEL authors="Simple"
ARG GITHUB_USERNAME
ARG GITHUB_PAT
ENV GITHUB_USERNAME=${GITHUB_USERNAME} GITHUB_PAT=${GITHUB_PAT}
ENV HTTP_PROXY=socks5://host.docker.internal:1080 HTTPS_PROXY=socks5://host.docker.internal:1080
ENV DEBIAN_FRONTEND=noninteractive TZ=Asia/Shanghai

RUN sed -i "s/deb.debian.org/mirrors.aliyun.com/g" /etc/apt/sources.list.d/debian.sources && apt-get update > /dev/null && apt-get upgrade -y > /dev/null && \
    apt-get install -y tzdata > /dev/null && ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone && \
    sh -c "$(curl -fsLS get.chezmoi.io)"
USER vscode
RUN chezmoi init https://$GITHUB_USERNAME:$GITHUB_PAT@github.com/$GITHUB_USERNAME/dotfiles.git
ENTRYPOINT ["bash"]
```

**Build**

需要在 GitHub 上申请 [PAT](https://github.com/settings/personal-access-tokens), 权限只勾选 dotfiles 仓库的读取权限即可

```shell
docker build --no-cache --build-arg GITHUB_USERNAME=thesixonenine --build-arg GITHUB_PAT=github_pat_* -t thesixonenine/dev-java:8-bookworm .
```

```shell
docker run --rm -it --user vscode thesixonenine/dev-java:8-bookworm
# update
cd && chezmoi update
```

> 构建其他语言的镜像基本相同, 只需要替换基础镜像名称和构建的镜像名称即可.

## MySQL Client

**Dockerfile**

```Dockerfile
FROM alpine:latest
LABEL authors="Simple"
ENV TZ=Asia/Shanghai
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk add --no-cache tzdata mysql-client && \
    cp /usr/share/zoneinfo/${TZ} /etc/localtime && \
    echo "${TZ}" > /etc/timezone && \
    apk del tzdata && \
    echo -e "[client]\nskip-ssl" > /etc/my.cnf

WORKDIR /work
ENTRYPOINT ["mariadb"]
```

**Build**

```shell
docker build -t mysql-client:latest .
```

**Run SQL**

指定 **host** **port** **user** **password** **database** 及要执行的 **SQL**

```shell
docker run --rm mysql-client --host= --port= --user= --password= --database= -B -e "show tables;"
```
