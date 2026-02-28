---
title: "termux-tutorial"
date: 2026-02-28T17:06:41+08:00
lastmod: 2026-02-28T17:06:41+08:00
categories: ['']
tags: ['']
keywords: 
description: 
image: 
draft: true
---

### 切换源

参考 [阿里云镜像](https://developer.aliyun.com/mirror/termux)

```shell
sed -i 's@^\(deb.*stable main\)$@#\1\ndeb https://mirrors.aliyun.com/termux/termux-packages-24 stable main@' $PREFIX/etc/apt/sources.list
```

```shell
pkg update && pkg install openssh iproute2 nmap expect chezmoi git vim
```

### 修改密码

```shell
passwd
```

### 查看ip

```shell
ip a
```

### 开启 sshd

```shell
sshd
```

### 切换到电脑上操作

初始化 chezmoi

```shell
echo "HTTP_PROXY=socks5://127.0.0.1:1080" >> ~/.bashrc
echo "HTTPS_PROXY=socks5://127.0.0.1:1080" >> ~/.bashrc
git config --global http.https://github.com.proxy $HTTP_PROXY
chezmoi init https://$USERNAME:$MY_PAT@github.com/$USERNAME/dotfiles.git
```
