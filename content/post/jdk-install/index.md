---
title: jdk-install
date: 2021-12-21T22:35:21+0800
lastmod: 2021-12-21T22:35:21+0800
tags: ['JDK']
categories: ['Java']
keywords: jdk
description: jdk安装配置
---

## JDK17安装

```bash
#!/bin/bash

# 安装必须工具
yum -q -y install wget

# jdk17下载页: https://www.oracle.com/java/technologies/javase/jdk17-archive-downloads.html
wget https://download.oracle.com/java/17/archive/jdk-17.0.10_linux-x64_bin.rpm
# 安装
rpm -i jdk-17.0.10_linux-x64_bin.rpm

# 设置并生效环境变量
echo 'JAVA_HOME=/usr/lib/jvm/jdk-17-oracle-x64' >> /etc/profile
echo 'PATH=$PATH:$JAVA_HOME/bin' >> /etc/profile
source /etc/profile

# 打印JDK安装版本信息
echo '################################################'
echo '### JAVA_HOME=/usr/lib/jvm/jdk-17-oracle-x64 ###'
echo '################################################'
java -version
```

## JDK8安装

```bash
#!/bin/bash

# 手动上传rpm包
# jdk8下载页: https://www.oracle.com/java/technologies/javase/javase8u211-later-archive-downloads.html

# 安装
rpm -i jdk-8u301-linux-x64.rpm

# 设置并生效环境变量
echo 'JAVA_HOME=/usr/java/jdk1.8.0_301-amd64' >> /etc/profile
echo 'PATH=$PATH:$JAVA_HOME/bin' >> /etc/profile
source /etc/profile

# 打印JDK安装版本信息
echo '##############################################'
echo '### JAVA_HOME=/usr/java/jdk1.8.0_301-amd64 ###'
echo '##############################################'
java -version
```

