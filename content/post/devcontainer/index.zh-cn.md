---
title: "devcontainer"
date: 2025-09-10T11:19:26
lastmod: 2025-10-20T10:23:02+0800
categories: ['Docker']
keywords: devcontainer
description: Dev Container
---

**所有环境构建均在官方标准之上进行构建**

## 文档参考

[官网](https://containers.dev)

[模板列表](https://github.com/devcontainers/templates/tree/main/src)

## 构建说明

1. 基于官方基础镜像, 替换源和时区, 并指定代理 **HTTP_PROXY/HTTPS_PROXY**
2. 安装 [chezmoi](https://chezmoi.io) 来同步环境设置
3. [chezmoi](https://chezmoi.io) 需要先创建 **dotfiles** 仓库及该仓库的只读 [PAT](https://github.com/settings/personal-access-tokens), [PAT](https://github.com/settings/personal-access-tokens) 同时也作为密码来解密用 **aes-128-cbc** 加密的 **gpg** 密钥文件
4. 初始化 [chezmoi](https://chezmoi.io) 时会 clone **dotfiles** 仓库需要走代理, 而 git 并不会使用 `HTTP_PROXY` 或 `HTTPS_PROXY`, 需要手动指定
5. 首次应用 [chezmoi](https://chezmoi.io) 时会先用 [PAT](https://github.com/settings/personal-access-tokens) 解密并导入 **gpg** 的密钥文件, 然后再使用 **gpg** 来解密 **ssh** 的公钥和私钥
6. 使用 **ssh** 来 clone 仓库, 如果还需要走代理, 则在第 2 步中安装 [ncat](https://nmap.org/ncat)

```shell
# 加密 gpg 密钥文件
openssl enc -aes-128-cbc -pbkdf2 -in ~/.gpg/SECRET.asc -out ~/.gpg/SECRET.asc.enc -pass env:GITHUB_PAT
openssl enc -aes-128-cbc -pbkdf2 -in ~/.gpg/public.asc -out ~/.gpg/public.asc.enc -pass env:GITHUB_PAT
# 解密 gpg 密钥文件
openssl aes-128-cbc -d -pbkdf2 -in ~/.gpg/SECRET.asc.enc -out ~/.gpg/SECRET.asc -pass env:GITHUB_PAT
openssl aes-128-cbc -d -pbkdf2 -in ~/.gpg/public.asc.enc -out ~/.gpg/public.asc -pass env:GITHUB_PAT
```


**环境变量**

- GITHUB_USERNAME: GitHub 的用户名, 用于拉取 **dotfiles** 仓库
- GITHUB_PAT: GitHub dotfiles 仓库的 [PAT](https://github.com/settings/personal-access-tokens), 用于拉取 **dotfiles** 仓库及后续解密并导入 **gpg** 密钥文件
- HTTP_PROXY/HTTPS_PROXY: 代理, 用于安装和更新 [chezmoi](https://chezmoi.io)
- TZ=Asia/Shanghai: 上海时区
- DEBIAN_FRONTEND=noninteractive: 避免交互式提示

**构建命令**

1. 修改 `APT` 源为阿里云并更新
2. 安装 [chezmoi](https://chezmoi.io), **ncat**(可选, 用于 **ssh** 走代理)
3. 切换到 vscode 用户并初始化 **dotfiles** 仓库

> 后续的软件安装及环境配置均由 **dotfiles** 中的脚本完成

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

RUN sed -i -e 's@deb.debian.org@mirrors.aliyun.com@;s@http:@https:@' /etc/apt/sources.list.d/debian.sources && \
    apt-get update > /dev/null && apt-get upgrade -y > /dev/null && \
    apt-get install -y ncat > /dev/null && \
    sh -c "$(curl -fsLS get.chezmoi.io)"
USER vscode
RUN git config --global http.https://github.com.proxy $HTTP_PROXY && \
    chezmoi init https://$GITHUB_USERNAME:$GITHUB_PAT@github.com/$GITHUB_USERNAME/dotfiles.git
ENTRYPOINT ["bash"]
```

如果基础镜像不支持通过环境变量 **TZ** 来修改时区, 则可以在更新软件包后安装 `tzdata` 来修改时区

```shell
apt-get install -y tzdata > /dev/null && ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
```

**devcontainer.json**

```json
{
  "name": "thesixonenine/dev-java:8-bookworm",
  "build": {
    "dockerfile": "Dockerfile",
    "args": {
        "GITHUB_USERNAME": "thesixonenine", 
        "GITHUB_PAT": "${localEnv:GITHUB_PAT}"
    }
  },
  "customizations": { "vscode": { "settings": {
        "extensions.allowed": {
            "microsoft": false,
            "github": false
        }
  }}}
}
```

**Build**

需要在 GitHub 上申请 [PAT](https://github.com/settings/personal-access-tokens), 权限只勾选 dotfiles 仓库的读取权限即可


**提前构建**

使用 [Dev Container CLI](https://github.com/devcontainers/cli) 构建镜像

```shell
npm install -g @devcontainers/cli --registry=https://registry.npmmirror.com
```

将以上 `Dockerfile` 和 `devcontainer.json` 置于当前目录下的 `.devcontainer` 目录中, 配置环境变量 `GITHUB_PAT` 后在当前目录下执行

```shell
devcontainer build --no-cache --workspace-folder . --image-name thesixonenine/dev-java:8-bookworm
```

**直接构建**

直接根据 `Dockerfile` 进行构建

```shell
docker build --no-cache --build-arg GITHUB_USERNAME=thesixonenine --build-arg GITHUB_PAT=github_pat_* -t thesixonenine/dev-java:8-bookworm .
```

```shell
docker run --rm -it --user vscode thesixonenine/dev-java:8-bookworm
# update
cd && chezmoi update
```

> 构建其他语言的镜像基本相同, 只需要替换基础镜像名称和构建的镜像名称即可.

## 使用构建

1. 打开 `VSCode` 并用 `Ctrl` + `Shift` + `p` 打开命令面板
2. 输入 `Dev Containers: Clone Repository in Named Container Volume...` 并选中
3. 输入仓库地址, 例如: `git@github.com:thesixonenine/dotfiles.git`
4. 选择新建命名的卷并输入名称, 与仓库名称相同即可
5. 按 `Enter` 确认使用默认的目录名称, 即与仓库名称相同

此时开始构建, 构建过程中需要从指定镜像开始, 例如 `/tmp/vsch-simple/bootstrap-image/0.427.0/bootstrap.Dockerfile`

> `0.427.0` 是指 `Dev Containers` 插件的版本, Identifier 是 `ms-vscode-remote.remote-containers`

在该 `Dockerfile` 中涉及 `alpine` 的镜像及软件安装, 可以先关闭 `VSCode`, 然后在 `WSL2` 中编辑该文件, 增加镜像源

直接命令修改

1. 修改软件镜像源
2. 修改npm镜像源
3. 指定git针对 `https://github.com` 进行代理, 以便在 `VSCode` 中使用 `https` 的方式进行 Clone 的时候加快速度

```shell
sed -i '10cRUN sed -i '\''s/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g'\'' /etc/apk/repositories' /tmp/vsch-simple/bootstrap-image/0.427.0/bootstrap.Dockerfile \
sed -i '23cRUN npm config set registry https://registry.npmmirror.com && git config --global http.https://github.com.proxy socks5://host.docker.internal:1080' /tmp/vsch-simple/bootstrap-image/0.427.0/bootstrap.Dockerfile
```

修改前的 `Dockerfile` 内容如下:

```Dockerfile
FROM mcr.microsoft.com/devcontainers/base:0-alpine-3.20

COPY host-ca-certificates.crt /tmp/host-ca-certificates.crt
RUN cat /tmp/host-ca-certificates.crt >> /etc/ssl/certs/ca-certificates.crt
RUN csplit -f /usr/local/share/ca-certificates/host-ca-certificate- -b '%02d.pem' -z -s /tmp/host-ca-certificates.crt '/-----BEGIN CERTIFICATE-----/' '{*}'
ENV NODE_EXTRA_CA_CERTS=/etc/ssl/certs/ca-certificates.crt

# Avoiding OpenSSH >8.8 for compatibility for now: https://github.com/microsoft/vscode-remote-release/issues/7482
RUN echo "@old https://dl-cdn.alpinelinux.org/alpine/v3.15/main" >> /etc/apk/repositories

RUN apk add --no-cache \
        git-lfs \
        nodejs \
        python3 \
        npm \
        make \
        g++ \
        docker-cli \
        docker-cli-buildx \
        docker-cli-compose \
        openssh-client-default@old \
        ;

RUN npm config set cafile /etc/ssl/certs/ca-certificates.crt && cd && npm i node-pty || echo "Continuing without node-pty."

COPY .vscode-remote-containers /root/.vscode-remote-containers
```

修改后的 `Dockerfile` 内容如下:

```Dockerfile
FROM mcr.microsoft.com/devcontainers/base:0-alpine-3.20

COPY host-ca-certificates.crt /tmp/host-ca-certificates.crt
RUN cat /tmp/host-ca-certificates.crt >> /etc/ssl/certs/ca-certificates.crt
RUN csplit -f /usr/local/share/ca-certificates/host-ca-certificate- -b '%02d.pem' -z -s /tmp/host-ca-certificates.crt '/-----BEGIN CERTIFICATE-----/' '{*}'
ENV NODE_EXTRA_CA_CERTS=/etc/ssl/certs/ca-certificates.crt

# Avoiding OpenSSH >8.8 for compatibility for now: https://github.com/microsoft/vscode-remote-release/issues/7482
RUN echo "@old https://dl-cdn.alpinelinux.org/alpine/v3.15/main" >> /etc/apk/repositories
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add --no-cache \
        git-lfs \
        nodejs \
        python3 \
        npm \
        make \
        g++ \
        docker-cli \
        docker-cli-buildx \
        docker-cli-compose \
        openssh-client-default@old \
        ;
RUN npm config set registry https://registry.npmmirror.com && git config --global http.https://github.com.proxy socks5://host.docker.internal:1080
RUN npm config set cafile /etc/ssl/certs/ca-certificates.crt && cd && npm i node-pty || echo "Continuing without node-pty."

COPY .vscode-remote-containers /root/.vscode-remote-containers
```

修改点如下:

**指定 apk 镜像源**

```Dockerfile
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
```

**指定 npm 镜像源**

```Dockerfile
RUN npm config set registry https://registry.npmmirror.com
```

另外, 还可以预先拉取该 `Dockerfile` 中的基础镜像 `mcr.microsoft.com/devcontainers/base:0-alpine-3.20`

## 镜像源替换

### debian

```shell
sudo sed -i -e 's@deb.debian.org@mirrors.aliyun.com@;s@http:@https:@' /etc/apt/sources.list.d/debian.sources
```

### ubuntu

```shell
sudo sed -i 's@//.*archive.ubuntu.com@//mirrors.aliyun.com@g' /etc/apt/sources.list.d/ubuntu.sources
```

### alpine

```shell
sudo sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
```

### npm

```shell
npm config set registry https://registry.npmmirror.com
```

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
