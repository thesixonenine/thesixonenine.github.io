---
title: "git-repo-proxy"
description: 
date: 2023-12-04T15:44:18+08:00
lastmod: 2023-11-29T17:14:10+08:00
categories: ['Git']
keywords: git repo proxy
description: git repo proxy
---

[参考](https://stackoverflow.com/a/67513102)

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
