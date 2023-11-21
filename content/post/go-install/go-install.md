---
title: go-install
date: 2022-02-23T16:49:28+0800
updated: 2022-03-01T10:54:28+0800
tags: 
- Go
- Linux
categories: 
- Go
- Linux
keywords: Go
description: CentOS 7 Go安装配置
url: '/p/go-install.html'
---

```bash
#!/bin/bash
echo '安装go1.17.7'
if ! command -v wget &> /dev/null
then
    echo "wget could not be found"
    exit
fi
# 下载
wget https://golang.google.cn/dl/go1.17.7.linux-amd64.tar.gz
# 删除旧版本, 解压新版本
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.17.7.linux-amd64.tar.gz

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
# 打印go版本
go version
```