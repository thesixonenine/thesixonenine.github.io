---
title: seal-the-vpn
date: 2025-12-10T23:11:21+0800
lastmod: 2025-12-10T23:35:21+0800
tags: ['Docker']
categories: ['Docker']
keywords: vpn
description: 封印vpn
---

使用 `Docker` 来运行常见的 VPN 软件以便与主机隔离

- qianxin(奇安信)
- easyconnect
- atrust

## install docker

### install docker on Mac

```shell
https://desktop.docker.com/mac/main/arm64/Docker.dmg
```

### install docker on Windows

```shell
https://desktop.docker.com/win/main/amd64/Docker%20Desktop%20Installer.exe
```

## qianxin

### pull image

```shell
docker pull lukbinx/qianxin-client:1.2.1.463
```

**china mirror**

```shell
docker pull swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/lukbinx/qianxin-client:1.2.1.463
docker tag swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/lukbinx/qianxin-client:1.2.1.463 docker.io/lukbinx/qianxin-client:1.2.1.463
docker rmi swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/lukbinx/qianxin-client:1.2.1.463
```

### run image

```shell
docker run -d --name qianxin-lukbinx --restart=always -e TZ=Asia/Shanghai -e KASM_VNC_SSL=0 -e KASM_NO_AUTH=1 -p 26901:6901 -p 21080:1080 --cap-add=NET_ADMIN --device=/dev/net/tun:/dev/net/tun --shm-size=512m --ulimit nofile=1048576:1048576 lukbinx/qianxin-client:1.2.1.463
```

### add hosts

```shell
docker exec -it -u root qianxin-lukbinx /bin/bash
```

```shell
cat >> /etc/hosts << 'EOF'
# add hosts
EOF
```

### login vpn

http://localhost:26901

### open browser use socks5 vpn

#### Mac

```shell
open -a /Applications/Google\ Chrome.app/ --args --proxy-server=127.0.0.1:21080
```

#### Windows

**edge**

```shell
& "C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe" --proxy-server="socks5://127.0.0.1:21080"
```

**chrome**

```shell
& "C:\Program Files (x86)\Google\Chrome\Application\chrome.exe" --proxy-server="socks5://127.0.0.1:21080"
```
