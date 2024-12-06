---
title: "hosts"
date: 2024-12-06T11:13:09+08:00
lastmod: 2024-12-06T11:13:09+08:00
categories: ['Script']
tags: ['hosts']
keywords: hosts
description: hosts
---

记录一些需要在hosts中修改的域名

## 被禁止的域名

```
# baidu统计
0.0.0.0                       hm.baidu.com
# EDGE浏览器禁用新闻
0.0.0.0                       ntp.msn.cn
0.0.0.0                       browser.events.data.msn.com
# 禁止firefox国内版
0.0.0.0                       www.firefox.com.cn
0.0.0.0                       firefox.com.cn
0.0.0.0                       download-ssl.firefox.com.cn
```

## 访问慢的域名

```json
[
    "cdn.jsdelivr.net",
    "outlook.live.com",
    "docs.live.net",
    "d.docs.live.net",
    "roaming.officeapps.live.com",
    "ocws.officeapps.live.com",
    "mobile.pipe.aria.microsoft.com",
    "onedrive.live.com",
    "api.onedrive.com",
    "skydrivesync.policies.live.net",
    "oneclient.sfx.ms",
    "storage.live.com",
    "skydrive.wns.windows.com",
    "contentsync.onenote.com",
    "www.onenote.com",
    "api.adoptopenjdk.net",
    "repo.maven.apache.org",
    "hub.docker.com",
    "services.gradle.org",
]
```
