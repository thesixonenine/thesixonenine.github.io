---
title: scoop
date: 2022-08-08T11:11:21+0800
lastmod: 2022-08-31T17:19:21+0800
tags:
- Windows
- Scoop
categories:
- Windows
- Scoop
keywords: scoop
description: scoop使用
---

[`Scoop`](https://scoop.sh/) 是 `Windows` 下的一款软件包管理工具.

## 安装

设置 `PowerShell` 脚本执行策略, 然后下载安装脚本并执行

```powershell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
irm get.scoop.sh | iex
```

## 代理设置

```powershell
# 查看代理
scoop config proxy
# 设置代理
scoop config proxy 127.0.0.1:10809
# 取消代理
scoop config rm proxy
```

## `Bucket` 添加

```powershell
# 查看官方推荐仓库
scoop bucket known
# 添加bucket
scoop bucket add extras
scoop bucket add nerd-fonts
scoop bucket add java
# 添加自己的bucket
scoop bucket add thesixonenine-scoop-bucket https://github.com/thesixonenine/thesixonenine-scoop-bucket
```

## 推荐安装的软件

**其他安装**

- weasel [rime输入法](https://github.com/rime/weasel)
- Microsoft Edge [浏览器](https://www.microsoft.com/zh-cn/edge)
- QQ [聊天软件](https://im.qq.com/index)
- Intel Unison [电脑连接手机](https://www.microsoft.com/store/productId/9PP9GZM2GN26)

### 常用必备

- 7zip (压缩软件)
- aria2 (下载软件)
- bandizip (压缩软件)
- carnac (按键提示)
- colortool (配色工具)
- everything (文件搜索工具)
- googlechrome (浏览器)
- mpv (视频播放器)
- musikcube (音乐播放器)
- potplayer (视频播放器)
- qbittorrent-enhanced (BT下载软件)
- scrcpy (安卓手机投屏软件)
- sumatrapdf (PDF阅读工具)
- snipaste (截屏软件)
- trafficmonitor (电脑CPU,网速,内存监控软件)
- translucenttb (任务栏透明化工具)
- twinkle-tray (显示器亮度调整工具)

### 常用可选

- anki (卡片记忆)
- asciidocfx (asciidoc文档编写软件)
- asciidoctorj (asciidoc文档转换工具)
- dismplusplus (系统优化工具)
- ffmpeg (多媒体工具)
- figlet (基于ASCII字符组成的字符画)
- fscapture (截屏软件)
- msiafterburner (微星小飞机,显卡超频工具)
- obs-studio (视频录制和视频推流软件)
- pandoc (标记语言转换工具)
- youtube-dl (Youtube视频下载工具)

### 开发必备

- bind (服务器软件集, 包括host, dig)
- filezilla (FTP软件)
- git (版本控制软件)
- git-scripts (个人Git小脚本)
- gpg4win (基于GPG的非对称加密软件)
- gradle (构建工具)
- innounp (安装程序解包工具)
- iperf3 (带宽性能测量工具)
- jd-gui (Java反编译工具)
- jetbrains-toolbox (Jetbrains工具箱)
- maven (构建工具)
- nmap (端口扫描工具)
- notepadplusplus (文本编辑工具)
- oh-my-posh (命令行增强工具)
- openssl (ssl工具)
- powertoys (微软生产力工具集)
- proxifier (客户端代理软件, 让不支持代理的程序能使用代理)
- proxychains (命令行代理软件, 让程序的TCP连接请求使用代理)
- scoop-completion (scoop命令提示)
- sourcetree (Git可视化管理工具)
- switchhosts (hosts文件管理工具)
- v2rayn (基于v2Ray内核的Windows客户端)
- virtualbox (虚拟机软件)

### 开发可选

- arthas (Java诊断工具)
- eclipse-mat (Java堆内存分析器)
- fiddler (HTTP调试代理)
- frp (内网穿透和反向代理软件)
- jmeter (压力测试工具)
- redis (基于内存的KV数据库)
- tcping (类似ping, 但是使用的不是ICMP协议, 而是TCP协议)

### `Aria2` 配置参数

```powershell
# aria2 在 Scoop 中默认开启
scoop config aria2-enabled true
scoop config aria2-retry-wait 4
scoop config aria2-split 16
scoop config aria2-max-connection-per-server 16
scoop config aria2-min-split-size 4M
```

## 自制 `Bucket`

参考官方 [`WiKi`](https://github.com/ScoopInstaller/Scoop/wiki/Buckets#creating-your-own-bucket)

[My-Personal-Scoop-Bucket](https://github.com/thesixonenine/thesixonenine-scoop-bucket)

```powershell
scoop bucket add thesixonenine-scoop-bucket https://github.com/thesixonenine/thesixonenine-scoop-bucket
```
