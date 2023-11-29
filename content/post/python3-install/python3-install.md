---
title: python3-install
date: 2022-05-05T11:42:26+0800
updated: 2022-05-05T11:42:26+0800
tags: ['Python','Linux']
categories: ['Python','Linux']
keywords: Python
description: CentOS 7 Python3安装配置
---

## Python3.9.12 编译安装

以下命令均以 `root` 用户执行
```bash
#!/bin/bash
cd ~
# 安装wget用来下载Python源码
yum install -y wget
# 安装编译所需的工具
yum install -y gcc zlib-devel bzip2-devel openssl-devel ncurses-devel sqlite-devel readline-devel tk-devel gdbm-devel db4-devel libpcap-devel xz-devel libffi-devel
# 下载Python源码
wget https://www.python.org/ftp/python/3.9.12/Python-3.9.12.tar.xz
# 解压Python源码到/usr/local/
tar -xf Python-3.9.12.tar.xz -C /usr/local/
# 切换目录到源码目录准备开始编译安装
cd /usr/local/Python-3.9.12/
# 创建安装的目录
mkdir /usr/local/python3
# 配置参数
# 低版本的gcc版本中带有 --enable-optimizations 参数时会出现 Could not import runpy module 安装错误
# 1. 升级gcc至8.1.0
# 2. 去掉--enable-optimizations
./configure --with-ssl --prefix=/usr/local/python3 --with-ensurepip=install
# 编译
make
# 安装
make install
# 建立软链接
ln -s /usr/local/python3/bin/python3 /usr/bin/python3
ln -s /usr/local/python3/bin/pip3 /usr/bin/pip3
```
