---
title: mysql5.7-install
date: 2021-12-21T22:38:43+0800
updated: 2021-12-21T22:38:43+0800
tags: ['MySQL','Linux']
categories: ['MySQL','Linux']
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
yum install -y net-tools
yum install -y libaio
yum install -y numactl
yum install -y perl
yum install -y wget

cd /usr/local

# 下载MySQL相关的rpm文件
wget https://dev.mysql.com/get/Downloads/MySQL-5.7/mysql-community-server-5.7.24-1.el7.x86_64.rpm
wget https://dev.mysql.com/get/Downloads/MySQL-5.7/mysql-community-common-5.7.24-1.el7.x86_64.rpm
wget https://dev.mysql.com/get/Downloads/MySQL-5.7/mysql-community-libs-5.7.24-1.el7.x86_64.rpm
wget https://dev.mysql.com/get/Downloads/MySQL-5.7/mysql-community-client-5.7.24-1.el7.x86_64.rpm

# 移除已有的mariadb
yum remove -y mariadb-*

# 按顺序安装MySQL
rpm -i mysql-community-common-5.7.24-1.el7.x86_64.rpm
rpm -i mysql-community-libs-5.7.24-1.el7.x86_64.rpm
rpm -i mysql-community-client-5.7.24-1.el7.x86_64.rpm
rpm -i mysql-community-server-5.7.24-1.el7.x86_64.rpm

# 初始化MySQL
mysqld --initialize --user=mysql
chown mysql:mysql /var/lib/mysql -R
mysql_ssl_rsa_setup
systemctl start mysqld
systemctl enable mysqld
# mysqld_safe --user=mysql &
echo "下面试MySQL的临时密码, 请在登录MySQL登录后使用如下语句修改MySQL的密码"
grep 'temporary password' /var/log/mysqld.log
echo "修改root密码的SQL:"
echo "ALTER USER 'root'@'localhost' IDENTIFIED BY '123456';"
echo "修改root权限的SQL:"
echo "GRANT ALL PRIVILEGES ON *.* TO 'root'@'%' IDENTIFIED BY '123456';"
echo "刷新权限的SQL:"
echo "FLUSH PRIVILEGES;"
mysql -uroot -p
```

