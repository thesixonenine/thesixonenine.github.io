---
title: "git-repo-proxy"
date: 2023-12-04T15:44:18+08:00
lastmod: 2025-04-09T22:20:52+0800
categories: ['Git']
keywords: git repo proxy
description: git repo proxy
---

[参考](https://stackoverflow.com/a/67513102)

## Recommend

```shell
# cat $HOME/.ssh/config
# 代理服务器可能不允许22端口, 可以尝试在HTTPS端口使用SSH
# https://docs.github.com/zh/authentication/troubleshooting-ssh/using-ssh-over-the-https-port
Host github.com
  User git
  Port 443
  Hostname ssh.github.com
  IdentityFile "~/.ssh/id_rsa"
  TCPKeepAlive yes
  ProxyCommand "C:\Program Files\Git\mingw64\bin\connect" -H 127.0.0.1:1080 %h %p

# 正常填写代理
Host github.com
  User git
  Port 22
  Hostname github.com
  IdentityFile "~/.ssh/id_rsa"
  TCPKeepAlive yes
  ProxyCommand "C:\Program Files\Git\mingw64\bin\connect" -S 127.0.0.1:1080 %h %p
  # MacOS
  # ProxyCommand nc -v -x 127.0.0.1:1080 %h %p
Host ssh.github.com
  User git
  Port 443
  Hostname ssh.github.com
  IdentityFile "~/.ssh/id_rsa"
  TCPKeepAlive yes
  ProxyCommand "C:\Program Files\Git\mingw64\bin\connect" -S 127.0.0.1:1080 %h %p
```

## Method 1. git http + proxy http

```bash
git config --global http.proxy "http://127.0.0.1:1080"
git config --global https.proxy "http://127.0.0.1:1080"
```

## Method 2. git http + proxy shocks

```bash
git config --global http.proxy "socks5://127.0.0.1:1080"
git config --global https.proxy "socks5://127.0.0.1:1080"
```

## to unset

```bash
git config --global --unset http.proxy
git config --global --unset https.proxy
```

## Method 3. git ssh + proxy http

```bash
vim ~/.ssh/config
Host github.com
HostName github.com
User git
ProxyCommand socat - PROXY:127.0.0.1:%h:%p,proxyport=1087
```

# Method 4. git ssh + proxy socks

```bash
vim ~/.ssh/config
Host github.com
HostName github.com
User git
# Linux
ProxyCommand nc --proxy-type socks5 --proxy 127.0.0.1:1080 %h %p
# Windows
# fist, scoop install main/nmap
ProxyCommand ncat --proxy-type socks5 --proxy 127.0.0.1:1080 %h %p
```
