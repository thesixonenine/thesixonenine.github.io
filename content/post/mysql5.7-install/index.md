---
title: mysql5.7-install
date: 2021-12-21T22:38:43+0800
lastmod: 2025-01-17T14:51:26+08:00
tags: ['Linux']
categories: ['MySQL']
keywords: MySQL
description: CentOS 7 MySQL5.7安装配置
---

## MySQL5.7.24安装配置

```bash
#!/bin/bash

# 新增mysql用户组和用户
groupadd mysql
useradd -r -g mysql -s /bin/false mysql

# 安装需要的工具
yum -q -y install net-tools libaio numactl perl wget

# 创建目录
mkdir /usr/local/mysql5.7 && cd /usr/local/mysql5.7

# 下载MySQL相关的rpm文件
wget https://dev.mysql.com/get/Downloads/MySQL-5.7/mysql-community-server-5.7.24-1.el7.x86_64.rpm
wget https://dev.mysql.com/get/Downloads/MySQL-5.7/mysql-community-common-5.7.24-1.el7.x86_64.rpm
wget https://dev.mysql.com/get/Downloads/MySQL-5.7/mysql-community-libs-5.7.24-1.el7.x86_64.rpm
wget https://dev.mysql.com/get/Downloads/MySQL-5.7/mysql-community-client-5.7.24-1.el7.x86_64.rpm

# 移除已有的mariadb
yum -q -y remove mariadb-*

# 安装MySQL
yum -y install mysql-community-{server,client,common,libs}-*

# 初始化MySQL
mysqld --initialize --user=mysql
chown mysql:mysql /var/lib/mysql -R
mysql_ssl_rsa_setup

# 启动MySQL
systemctl start mysqld
systemctl enable mysqld
echo ""
echo "下面是MySQL的临时密码, 请在登录MySQL后使用如下语句修改MySQL的密码"
grep 'temporary password' /var/log/mysqld.log
echo "修改root密码的SQL:"
echo "ALTER USER 'root'@'localhost' IDENTIFIED BY '123456';"
echo "修改root权限的SQL:"
echo "GRANT ALL PRIVILEGES ON *.* TO 'root'@'%' IDENTIFIED BY '123456';"
echo "刷新权限的SQL:"
echo "FLUSH PRIVILEGES;"
echo ""
mysql -uroot -p
```

## 创建库

```bash
CREATE DATABASE qing CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
```

## 创建用户

```bash
USE mysql;
CREATE USER 'dev'@'%' IDENTIFIED BY '123456';
# GRANT ALL PRIVILEGES ON *.* TO 'dev'@'%' WITH GRANT OPTION;
GRANT ALL PRIVILEGES ON qing.* TO 'dev'@'%' WITH GRANT OPTION;
```

## 创建只读用户
```bash
CREATE USER 'read_only'@'%' IDENTIFIED BY '123456';
GRANT SELECT ON product.* TO 'read_only'@'%' WITH GRANT OPTION;
FLUSH PRIVILEGES;
```

## 导出/入数据库
```
# 导出数据库所有(指定)表的结构, 加-d参数代表只导表结构
# mysqldump [选项] 数据库名 [表名] > 脚本名
mysqldump -uUSER -pPASSWORD [-d] dbname [table_name table_name...] > db_tables.sql
mysql -uUSER -pPASSWORD dbname < db_tables.sql
```
