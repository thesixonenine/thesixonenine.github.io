---
title: go-install
date: 2022-02-23T16:49:28
lastmod: 2025-09-15T14:13:26
tags: ['Linux']
categories: ['Go']
keywords: Go
description: Go Intall on Linux
---

```bash
#!/bin/bash
filename=go1.25.1.linux-amd64.tar.gz
echo "Intall ${filename}"
if ! command -v wget > /dev/null;then
    echo "wget could not be found"
    exit
fi
# 下载
wget https://golang.google.cn/dl/${filename}
# 删除旧版本, 解压新版本
rm -rf /usr/local/go && tar -C /usr/local -xzf ${filename}

mkdir -p $HOME/go/src
mkdir -p $HOME/go/pkg
mkdir -p $HOME/go/bin

# 设置并生效环境变量
echo 'GOROOT=/usr/local/go' >> /etc/profile
echo 'PATH=$PATH:$GOROOT/bin:$HOME/go/bin' >> /etc/profile
source /etc/profile

# 设置代理
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct
# 还原代理
go env -w GOPROXY=https://proxy.golang.org,direct
# 打印go版本
go version
```