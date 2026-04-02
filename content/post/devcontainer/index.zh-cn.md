---
title: "devcontainer"
date: 2025-09-10T11:19:26
lastmod: 2026-04-02T12:00:21
categories: ['Docker']
keywords: devcontainer
description: Dev Container
---

**所有环境构建均在官方标准之上进行构建**

## 文档参考

- [官网](https://containers.dev)
- [模板列表](https://github.com/devcontainers/templates/tree/main/src)
- [chezmoi](https://chezmoi.io)
- [PAT](https://github.com/settings/personal-access-tokens)
- [ncat](https://nmap.org/ncat)

## 构建说明

**构建准备**: 创建 **dotfiles** 仓库及该仓库的只读 `PAT`

1. 基于官方基础镜像, 替换源和时区, 并指定代理 **HTTP_PROXY/HTTPS_PROXY**
2. 安装 `chezmoi` 来同步环境设置, 安装 `age` 来支持解密文件
3. **dotfiles** 仓库的只读 `PAT` 作为密码来解密用 **age** 加密的 **key.txt** 密钥文件
4. 初始化 `chezmoi` 时会使用  clone **dotfiles** 仓库需要走代理, 而 git 并不会使用 `HTTP_PROXY` 或 `HTTPS_PROXY`, 需要手动指定
5. 首次应用 `chezmoi` 时会先用 `PAT` 解密 **key.txt.age** 密钥文件, 然后再使用 **age** 来解密其他加密文件
6. 使用 **ssh** 来 clone 仓库, 如果还需要走代理, 则在第 2 步中安装 `ncat`

> 使用 `Docker Desktop` 4.63.0及以上版本则可以不用配置 **HTTP_PROXY/HTTPS_PROXY**, 直接在 `Docker Desktop` 中指定容器的代理即可. 参见 `Settings` -> `Resources` - `Proxies` - `Containers proxy`
> 如果需要配置 **HTTP_PROXY/HTTPS_PROXY**, 则安装并使用 `ncat` 来支持

> 由于 `age` 解密不支持命令行传入密码来完成自动解密, 所以安装并使用 `expect` 来支持

<details>
<summary>gpg aes-128-cbc 加解密示例</summary>

加密

```shell
openssl enc -aes-128-cbc -salt -pbkdf2 -iter 100000 -in ./test.txt -out ./test.txt.enc -pass env:PASSWORD
```

解密

```shell
openssl aes-128-cbc -d -salt -pbkdf2 -in ./test.txt.enc -out ./test.txt -pass env:PASSWORD
```

</details>


**构建参数**

用于构建 `Dev Container` 容器传入, 用来传递个性化信息或密钥信息, 避免后续环境中存在这些不该暴露的信息

- USERNAME: GitHub 的用户名, 用于拉取 **dotfiles** 仓库
- MY_PAT: GitHub dotfiles 仓库的 `PAT`, 用于拉取 **dotfiles** 仓库及后续解密 **key.txt.age** 密钥文件

**环境变量**

- HTTP_PROXY/HTTPS_PROXY: 代理, 用于安装和更新 `chezmoi`
- TZ=Asia/Shanghai: 上海时区
- DEBIAN_FRONTEND=noninteractive: 避免交互式提示

**构建命令**

1. 修改 `APT` 源为阿里云并更新
2. 安装 `chezmoi`, **ncat**(用于 **ssh** 走代理), age和expect(用于解密key.txt.age文件)
3. 切换到 vscode 用户并初始化 **dotfiles** 仓库

> 后续的软件安装及环境配置均由 **dotfiles** 中的脚本完成

> `chezmoi` 的 `run_once` 脚本中需要包含删除 `PAT` 的步骤, 以免泄漏 `PAT`

## 构建步骤

以 Java 为例, 其他镜像同理, 注意镜像默认使用的用户

**Dockerfile**

```Dockerfile
FROM mcr.microsoft.com/devcontainers/java:8-bookworm
LABEL authors="Simple"
ARG MY_PAT
ARG USERNAME=thesixonenine
ENV HTTP_PROXY=socks5://host.docker.internal:1080 HTTPS_PROXY=socks5://host.docker.internal:1080
ENV TZ=Asia/Shanghai DEBIAN_FRONTEND=noninteractive

RUN sed -i -e 's@deb.debian.org@mirrors.aliyun.com@;s@http:@https:@' /etc/apt/sources.list.d/debian.sources && \
    apt-get update > /dev/null && apt-get upgrade -y > /dev/null && \
    apt-get install -y ncat age expect > /dev/null && \
    sh -c "$(curl -fsLS get.chezmoi.io)" && \
    chsh -s /usr/bin/zsh vscode && \
    rm -rf /var/lib/apt/lists/*
USER vscode
RUN git config --global http.https://github.com.proxy $HTTP_PROXY && \
    chezmoi init https://$USERNAME:$MY_PAT@github.com/$USERNAME/dotfiles.git
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
        "USERNAME": "thesixonenine",
        "MY_PAT": "${localEnv:MY_PAT}"
    }
  },
  "postStartCommand": "chezmoi update",
  "customizations": { "vscode": { "settings": {
        "extensions.allowed": {
            "microsoft": false,
            "github": false
        }
  }}}
}
```

**Build**

需要在 GitHub 上申请 `PAT`, 权限只勾选 dotfiles 仓库的读取权限即可


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
docker build --no-cache --build-arg USERNAME=thesixonenine --build-arg MY_PAT=github_pat_* -t thesixonenine/dev-java:8-bookworm .
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
sed -i '10cRUN sed -i '\''s/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g'\'' /etc/apk/repositories' /tmp/vsch-simple/bootstrap-image/0.427.0/bootstrap.Dockerfile && \
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

## 更新基础镜像

**Powershell**

更新以 mcr 开头的镜像

```
docker images --format "{{.Repository}}:{{.Tag}}" | Where-Object { $_ -like "mcr*" -and $_ -notlike "*:<none>" } | ForEach-Object { docker pull $_ }
```

打印非 mcr 开头的镜像

```
docker images --format "{{.Repository}}:{{.Tag}}" | Where-Object { $_ -notlike "mcr*" } | ForEach-Object { Write-Host "$_" }
```


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


## 新环境初始化


很多时候需要临时在linux中试用一些软件或命令, 需要一个干净的环境, 这时可以直接使用 `mcr.microsoft.com/devcontainers/base:trixie`


```shell
docker run --rm -it mcr.microsoft.com/devcontainers/base:trixie /bin/zsh
```

```shell
sed -i -e 's@deb.debian.org@mirrors.aliyun.com@;s@http:@https:@' /etc/apt/sources.list.d/debian.sources && \
apt-get update > /dev/null && apt-get upgrade -y > /dev/null && \
apt-get install -y vim ncat age expect > /dev/null && \
sh -c "$(curl -fsLS get.chezmoi.io)" -- -b /bin
```

再从主机中复制 `chezmoi` 的初始化命令并执行


```shell
copy-chezmoi-init
```

最后刷新一下 shell 即可

```shell
source ~/.zshrc
```

### LazyVim

```shell
apt install -y git lazygit fd-find curl ripgrep neovim
```

```shell
ln -s $(which fdfind) /bin/fd
```

Version


```shell
git -v && \
lazygit -v && \
fdfind -V && \
curl -V && \
rg -V && \
nvim -v
```

Clone


```shell
git clone https://github.com/LazyVim/starter ~/.config/nvim
```

```shell
rm -rf ~/.config/nvim/.git
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
