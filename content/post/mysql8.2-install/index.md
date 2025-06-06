---
title: "mysql8.2-install"
date: 2024-01-30T15:04:16+08:00
lastmod: 2025-01-17T14:51:32+08:00
tags: ['Linux']
categories: ['MySQL']
keywords: MySQL
description: CentOS 7 MySQL8.2.0安装配置
---

## MySQL8.2.0安装配置

```bash
#!/bin/bash

# 新增mysql用户组和用户
groupadd mysql
useradd -r -g mysql -s /bin/false mysql

# 安装需要的工具
yum -q -y install net-tools libaio numactl perl wget

# 创建目录
mkdir /usr/local/mysql8 && cd /usr/local/mysql8

# 下载MySQL相关的rpm文件
wget https://dev.mysql.com/get/Downloads/MySQL-8.2/mysql-community-common-8.2.0-1.el7.x86_64.rpm
wget https://dev.mysql.com/get/Downloads/MySQL-8.2/mysql-community-libs-8.2.0-1.el7.x86_64.rpm
wget https://dev.mysql.com/get/Downloads/MySQL-8.2/mysql-community-client-plugins-8.2.0-1.el7.x86_64.rpm
wget https://dev.mysql.com/get/Downloads/MySQL-8.2/mysql-community-client-8.2.0-1.el7.x86_64.rpm
wget https://dev.mysql.com/get/Downloads/MySQL-8.2/mysql-community-icu-data-files-8.2.0-1.el7.x86_64.rpm
wget https://dev.mysql.com/get/Downloads/MySQL-8.2/mysql-community-server-8.2.0-1.el7.x86_64.rpm

# 移除已有的mariadb
yum -q -y remove mariadb-*

# 安装MySQL
yum -y install mysql-community-{server,client,client-plugins,icu-data-files,common,libs}-*

# 初始化MySQL
mysqld --initialize --user=mysql
chown mysql:mysql /var/lib/mysql -R
# 启动MySQL
systemctl start mysqld
systemctl enable mysqld
echo ""
echo "下面是MySQL的临时密码, 请在登录MySQL后使用如下语句修改MySQL的密码"
grep 'temporary password' /var/log/mysqld.log
echo "修改root密码的SQL:"
echo "ALTER USER 'root'@'localhost' IDENTIFIED BY '123456';"
echo "刷新权限的SQL:"
echo "FLUSH PRIVILEGES;"
echo ""
mysql -uroot -p
```

创建库

**MySQL 8 开始默认字符集为 `utf8mb4`, 排序规则为 `utf8mb4_0900_ai_ci`**

> - 0900 表示基于 Unicode 9.0 标准
> - ai (Accent Insensitive) 表示区分重音符号不敏感
> - ci (Case Insensitive) 表示大小写不敏感
