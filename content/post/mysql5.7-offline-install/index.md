---
title: "mysql5.7-offline-install"
date: 2024-10-10T16:33:43+08:00
lastmod: 2024-10-10T16:33:43+08:00
categories: ['MySQL']
tags: ['Linux']
keywords: MySQL
description: CentOS 7 MySQL5.7离线安装配置
---

## MySQL5.7.24离线安装配置


#### 卸载 CentOS 7 系统自带 mariadb

```bash
rpm -qa|grep mariadb
rpm -e --nodeps [item]
rm /etc/my.cnf
```

#### 检查用户和组

**不存在则创建**

```bash
cat /etc/group | grep mysql
cat /etc/passwd | grep mysql
# 新增mysql用户组和用户
groupadd mysql
useradd -r -g mysql -s /bin/false mysql
```

#### 将mysql上传并解压

```bash
cd /usr/local/
tar -zxf mysql-5.7.38-linux-glibc2.12-x86_64.tar.gz
mv mysql-5.7.38-linux-glibc2.12-x86_64 mysql
```

#### 更改所属的组和用户

```bash
cd /usr/local/
chown -R mysql mysql/
chgrp -R mysql mysql/
cd /usr/local/mysql/
mkdir data
chown -R mysql:mysql data
```

#### 创建my.cnf文件

```bash
vi /etc/my.cnf
```

**文件内容**

```toml
[mysql]
socket=/var/lib/mysql/mysql.sock
default-character-set=utf8mb4

[mysqld]
socket=/var/lib/mysql/mysql.sock
port=3306
basedir=/usr/local/mysql
datadir=/usr/local/mysql/data
max_connections=200
character-set-server=utf8mb4
default-storage-engine=INNODB
lower_case_table_names=1
max_allowed_packet=16M
explicit_defaults_for_timestamp=true

[mysql.server]
user=mysql
basedir=/usr/local/mysql
```

#### 安装mysql

```bash
/usr/local/mysql/bin/mysqld --initialize --user=mysql --basedir=/usr/local/mysql --datadir=/usr/local/mysql/data

echo 'export PATH=$PATH:/usr/local/mysql/bin' >> /etc/profile
source /etc/profile
```

#### 启动mysql

```bash
cp /usr/local/mysql/support-files/mysql.server  /etc/init.d/mysqld
# chmod 777 /usr/local/mysql/my.conf
chmod +x /etc/init.d/mysqld
mkdir /var/lib/mysql
chmod 777 /var/lib/mysql
# 启动
/etc/init.d/mysqld restart
# 重启
# /etc/init.d/mysqld restart
```

#### 修改mysql初始密码

```bash
# 获取初始密码
cat ~/.mysql_secret
# 登录
/usr/local/mysql/bin/mysql -uroot -p
# 修改密码
set PASSWORD = PASSWORD('123456');
flush privileges;
```
