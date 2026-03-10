---
title: "cloudflare-warp-in-docker"
date: 2026-03-10T15:05:50+08:00
lastmod: 2026-03-10T15:05:50+08:00
categories: ['Docker']
tags: ['Docker']
keywords: 
description: use the WARP client with Cloudflare Zero Trust in Docker
---


参考[repo](https://github.com/cmj2002/warp-docker)


## pull image


```powershell
docker pull swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/caomingjun/warp:latest `
docker tag swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/caomingjun/warp:latest docker.io/caomingjun/warp:latest `
docker rmi swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/caomingjun/warp:latest
```


```shell
docker pull swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/caomingjun/warp:latest \
docker tag swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/caomingjun/warp:latest docker.io/caomingjun/warp:latest \
docker rmi swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/caomingjun/warp:latest
```


## run image


```shell
docker run -d --name warp -p 1081:1080 -e WARP_SLEEP=2 -v warp:/var/lib/cloudflare-warp --cap-add MKNOD --cap-add AUDIT_WRITE --cap-add NET_ADMIN --device-cgroup-rule "c 10:200 rwm" --sysctl net.ipv6.conf.all.disable_ipv6=0 --sysctl net.ipv4.conf.all.src_valid_mark=1 caomingjun/warp
```


```powershell
docker run -d --name warp `
  -p 1081:1080 `
  -e WARP_SLEEP=2 `
  -v warp:/var/lib/cloudflare-warp `
  --cap-add MKNOD `
  --cap-add AUDIT_WRITE `
  --cap-add NET_ADMIN `
  --device-cgroup-rule "c 10:200 rwm" `
  --sysctl net.ipv6.conf.all.disable_ipv6=0 `
  --sysctl net.ipv4.conf.all.src_valid_mark=1 `
  caomingjun/warp
```


```shell
docker run -d --name warp \
  -p 1081:1080 \
  -e WARP_SLEEP=2 \
  -v warp:/var/lib/cloudflare-warp \
  --cap-add MKNOD \
  --cap-add AUDIT_WRITE \
  --cap-add NET_ADMIN \
  --device-cgroup-rule "c 10:200 rwm" \
  --sysctl net.ipv6.conf.all.disable_ipv6=0 \
  --sysctl net.ipv4.conf.all.src_valid_mark=1 \
  caomingjun/warp
```


## 登录


```shell
docker exec -it warp bash
```

```shell
warp-cli registration delete && warp-cli registration new thesixonenine
```

```shell
warp-cli registration token com.cloudflare.warp://thesixonenine.cloudflareaccess.com/auth?token=
```

```shell
warp-cli connect
```

```shell
warp-cli status
```

```shell
curl --socks5-hostname 127.0.0.1:1080 https://cloudflare.com/cdn-cgi/trace
```



以下部分摘抄自原repo的[文档](https://github.com/cmj2002/warp-docker/blob/main/docs/zero-trust.md)



If you want to use the WARP client with Cloudflare Zero Trust, just start the container without specifying license key, use `docker exec -it warp bash` to get into the container and follow these steps:

1. `warp-cli registration delete` to delete current registration
2. `warp-cli registration new <your-team-name>` to enroll the device
3. Open the link in the output in a browser and follow the instructions to complete the registration
4. On the success page, right-click and select **View Page Source**.
5. Find the HTML metadata tag that contains the token. For example, `<meta http-equiv="refresh" content"=0;url=com.cloudflare.warp://acmecorp.cloudflareaccess.com/auth?token=yeooilknmasdlfnlnsadfojDSFJndf_kjnasdf..." />`
6. Copy the URL field: `com.cloudflare.warp://<your-team-name>.cloudflareaccess.com/auth?token=<your-token>`
7. In the terminal, run the following command using the URL obtained in the previous step: `warp-cli registration token com.cloudflare.warp://<your-team-name>.cloudflareaccess.com/auth?token=<your-token>`. If you get an API error, then the token has expired. Generate a new one by refreshing the web page and quickly grab the new token from the page source.
8. `warp-cli connect` to reconnect using new registration.
9. Wait untill `warp-cli status` shows `Connected`.
10. Try `curl --socks5-hostname 127.0.0.1:1080 https://cloudflare.com/cdn-cgi/trace` to verify the connection.

This is only needed for the first time. After the device is enrolled, the registration information will be stored in the `./data` directory, if you don't delete them, the container will automatically use the registration information to connect to the WARP service after restart or recreate.
