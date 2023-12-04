---
title: hyper-v-dynamic-port-tcp
date: 2022-03-28T09:56:26+0800
lastmod: 2022-03-28T09:56:26+0800
tags: 
- Hyper-V
- Windows
categories:
- Hyper-V
- Windows
keywords: Hyper-V
description: Hyper-V 占用端口导致软件无法启动的问题
---
## 背景和原因

平时使用 `Hyper-V` 作为虚拟机管理比较多, 而 `Windows` 会为其网络服务分配一些端口, 而如果其他软件也要使用这些端口就会出错, 比较典型的就是 `Intellij IDEA` 的启动, `Tomcat` 的8080端口等等.

## 问题查看

1. 查看当前的TCP动态端口范围

   ```powershell
   "当前的TCP动态端口范围"
   netsh interface ipv4 show dynamicport tcp
   ```

2. 查看已被使用的端口

   ```powershell
   "当前已被使用的端口, 如果其他软件需要这些端口, 则可能无法启动"
   netsh interface ipv4 show excludedportrange protocol=tcp
   ```

## 问题解决

解决办法就是重新设置TCP动态端口的范围, 让 `Hyper-V` 只在该范围内占用端口

```powershell
netsh int ipv4 set dynamicport tcp start=49152 num=16383
netsh int ipv4 set dynamicport udp start=49152 num=16383
netsh int ipv6 set dynamicport tcp start=49152 num=16383
netsh int ipv6 set dynamicport udp start=49152 num=16383
```
