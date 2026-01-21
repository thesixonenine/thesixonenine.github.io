---
title: centos7-config
date: 2022-01-16T15:24:10+0800
lastmod: 2026-01-21T11:10:37+0800
tags: ['Linux', 'Vim']
categories: ['Linux']
keywords:
description: CentOS7常用配置
---

## CA证书过期导致yum更新失败

由于`Let's Encrypt's`的CA证书过期, 而`yum`使用的`curl`依赖CA证书, 解决办法就是更新证书

```bash
# 如果执行过yum clean all, 则需要先关闭ssl验证, 等更新CA证书后再开
# sudo cat "sslverify=0" >> /etc/yum.conf
sudo yum install -y ca-certificates
sudo update-ca-trust extract
```

## VIM常用配置

**~/.vimrc**

```bash
:set fileencodings=utf-8
:set encoding=utf-8
:set nowrap
:set nu
:set softtabstop=4
:set shiftwidth=4
:set tabstop=4

" 设置主题配色
" colorscheme solarized
" colorscheme koehler
colorscheme desert
```

## 修改yum源

### 阿里云

```bash
sudo mv /etc/yum.repos.d/CentOS-Base.repo /etc/yum.repos.d/CentOS-Base.repo.backup
sudo curl -o /etc/yum.repos.d/CentOS-Base.repo https://mirrors.aliyun.com/repo/Centos-7.repo
sudo sed -i -e '/mirrors.cloud.aliyuncs.com/d' -e '/mirrors.aliyuncs.com/d' /etc/yum.repos.d/CentOS-Base.repo
sudo yum makecache
```

### ustc

```bash
sudo sed -e 's|^mirrorlist=|#mirrorlist=|g' \
         -e 's|^#baseurl=http://mirror.centos.org/centos|baseurl=https://mirrors.ustc.edu.cn/centos|g' \
         -i.bak \
         /etc/yum.repos.d/CentOS-Base.repo
sudo yum makecache
```

### tsinghua

```bash
sudo sed -e 's|^mirrorlist=|#mirrorlist=|g' \
         -e 's|^#baseurl=http://mirror.centos.org|baseurl=https://mirrors.tuna.tsinghua.edu.cn|g' \
         -i.bak \
         /etc/yum.repos.d/CentOS-*.repo
sudo yum makecache
```

## 升级 Git 版本

rpm 安装

```shell
sudo yum -y remove git
sudo yum -y remove git-*

sudo yum -y install https://packages.endpointdev.com/rhel/7/os/x86_64/endpoint-repo.x86_64.rpm
sudo yum -y install git

git --version
```

源码安装

```shell
# 移除旧版本
sudo yum remove git
sudo yum remove git-*
# 安装依赖
sudo yum install -y wget curl-devel expat-devel gettext-devel openssl-devel zlib-devel
sudo yum groupinstall -y "Development Tools"
# 下载并解压源码
wget --user-agent="Mozilla" https://mirrors.ustc.edu.cn/kernel.org/software/scm/git/git-2.51.0.tar.gz
tar -xvf git-2.51.0.tar.gz
cd git-2.51.0
# 编译并安装
make configure
sudo ./configure --prefix=/usr
sudo make
sudo make install
# 检查安装的 Git 版本
git --version
```

