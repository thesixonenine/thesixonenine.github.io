---
title: "git-repo-proxy"
date: 2023-12-04T15:44:18+08:00
lastmod: 2025-04-24T10:11:52+0800
categories: ['Git']
keywords: git repo proxy
description: git repo proxy
---

[参考](https://ericclose.github.io/git-proxy-config.html)

走 ssh 协议代理时, 代理服务器可能不允许22端口, 可以尝试[在443端口使用SSH](https://docs.github.com/zh/authentication/troubleshooting-ssh/using-ssh-over-the-https-port), 即在config文件中指定 Port 443

## Recommend

```shell
# cat $HOME/.ssh/config
Host github.com
  User git
  Port 443
  Hostname ssh.github.com
  IdentityFile "~/.ssh/id_rsa"
  ProxyCommand "C:\Program Files\Git\mingw64\bin\connect" -S 127.0.0.1:1080 %h %p
```

## http 协议代理

直接在git的配置文件中进行配置即可

- 用户级配置文件: ~/.gitconfig (git config --global)
- 仓库级配置文件: ./.git/config (git config)

```shell
以下的 http://127.0.0.1:1080 和 socks5://127.0.0.1:1080 可以互换

# 全局网站代理
git config --global http.proxy http://127.0.0.1:1080
# cat ~/.gitconfig
[http]
  proxy = http://127.0.0.1:1080

# 指定网站代理
git config --global http.https://github.com.proxy socks5://127.0.0.1:1080
# cat ~/.gitconfig
[http "https://github.com"]
  proxy = socks5://127.0.0.1:1080

# 取消代理
git config --global --unset http.proxy
```

## ssh 协议代理

```shell
# cat ~/.ssh/config
Host github.com
  User git
  # 使用 443 端口, 避免代理服务器不支持 22 端口
  Port 443
  Hostname ssh.github.com
  IdentityFile "~/.ssh/id_rsa"
  ProxyCommand "C:\Program Files\Git\mingw64\bin\connect" -H 127.0.0.1:1080 %h %p

# Windows 使用 Git for Windows 默认附带的 connect 程序
# -H 代表 http, -S 代表 socks
ProxyCommand "C:\Program Files\Git\mingw64\bin\connect" -H 127.0.0.1:1080 %h %p

# Linux 使用 nc 命令, 查看 nc 命令的参数来进行设置代理类型及IP:PORT
# socks5 (-X 5 是默认的, 可以不加)
ProxyCommand nc -X 5 -x 127.0.0.1:1080 %h %p
# http
ProxyCommand nc -X connect -x 127.0.0.1:1080 %h %p

# socks5
ProxyCommand nc --proxy-type socks5 --proxy 127.0.0.1:1080 %h %p
# http
ProxyCommand nc --proxy-type http --proxy 127.0.0.1:1080 %h %p
```

