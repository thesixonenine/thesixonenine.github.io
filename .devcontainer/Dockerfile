# Update the NODE_VERSION arg in docker-compose.yml to pick a Node version: 18, 16, 14
ARG NODE_VERSION=16
FROM mcr.microsoft.com/devcontainers/javascript-node:${NODE_VERSION}
LABEL authors="Simple"
# VARIANT can be either 'hugo' for the standard version or 'hugo_extended' for the extended version.
ARG VARIANT=hugo
# VERSION can be either 'latest' or a specific version number
ARG VERSION=latest
ARG GITHUB_PAT
ENV GITHUB_USERNAME=thesixonenine GITHUB_PAT=${GITHUB_PAT}
ENV HTTP_PROXY=socks5://host.docker.internal:1080 HTTPS_PROXY=socks5://host.docker.internal:1080 GOPROXY=https://goproxy.cn,direct
ENV DEBIAN_FRONTEND=noninteractive TZ=Asia/Shanghai SHELL=/bin/zsh

# Download Hugo
RUN sed -i -e 's@deb.debian.org@mirrors.aliyun.com@;s@http:@https:@' /etc/apt/sources.list.d/debian.sources && \
    apt-get update > /dev/null && apt-get upgrade -y > /dev/null && \
    apt-get install -y ncat ca-certificates openssl git curl && \
    sh -c "$(curl -fsLS get.chezmoi.io)" && \
    rm -rf /var/lib/apt/lists/* && \
    case ${VERSION} in \
    latest) \
    export VERSION=$(curl -s https://api.github.com/repos/gohugoio/hugo/releases/latest | grep "tag_name" | awk '{print substr($2, 3, length($2)-4)}') ;;\
    esac && \
    echo ${VERSION} && \
    case $(uname -m) in \
    aarch64) \
    export ARCH=ARM64 ;; \
    *) \
    export ARCH=64bit ;; \
    esac && \
    echo ${ARCH} && \
    wget -O ${VERSION}.tar.gz https://github.com/gohugoio/hugo/releases/download/v${VERSION}/${VARIANT}_${VERSION}_Linux-${ARCH}.tar.gz && \
    tar xf ${VERSION}.tar.gz && \
    mv hugo /usr/bin/hugo

# Hugo dev server port
EXPOSE 1313

# [Optional] Uncomment this section to install additional OS packages you may want.
#
# RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
#     && apt-get -y install --no-install-recommends <your-package-list-here>

# [Optional] Uncomment if you want to install more global node packages
# RUN sudo -u node npm install -g <your-package-list-here>

USER node
RUN chezmoi init https://$GITHUB_USERNAME:$GITHUB_PAT@github.com/$GITHUB_USERNAME/dotfiles.git
