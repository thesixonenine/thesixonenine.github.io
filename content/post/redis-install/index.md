---
title: redis-install
date: 2021-12-21T22:40:56+0800
lastmod: 2021-12-21T22:40:56+0800
tags: ['Linux']
categories: ['Redis']
keywords: Redis
description: CentOS 7 Redis安装配置
---

## Redis安装配置

```bash
#!/bin/bash

# 安装必须工具
yum -q -y install gcc wget tar

cd /usr/local/

# 下载并解压redis, 其他版本替换版本号即可 6.2.14
wget http://download.redis.io/releases/redis-5.0.4.tar.gz
tar -zxf redis-5.0.4.tar.gz
cd redis-5.0.4

# 编译和安装redis
make
make install PREFIX=/usr/local/redis
cp ./redis.conf ../redis/bin/
cd /usr/local/redis/bin/

# 创建软链接
ln -s /usr/local/redis/bin/redis-cli /usr/bin/redis-cli
ln -s /usr/local/redis/bin/redis-server /usr/bin/redis-server

# 生成随机密码
PD=`date | md5sum | cut -b 1-8`


### 修改配置文件

# 后台运行
sed -i 's/^daemonize\ no$/daemonize\ yes/g' redis.conf
# 关闭保护模式
sed -i 's/^protected-mode\ yes$/protected-mode\ no/g' redis.conf
# 修改密码
sed -i "s/^#\ requirepass\ foobared$/requirepass\ $PD/g" redis.conf


### 设置开机启动

touch /etc/systemd/system/redis.service
cat>/etc/systemd/system/redis.service<<EOF
[Unit]
Description=redis-server
After=network.target
[Service]
Type=forking
ExecStart=/usr/local/redis/bin/redis-server /usr/local/redis/bin/redis.conf
PrivateTmp=true
[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl start redis.service
systemctl enable redis.service

### 打印随机密码
echo '### PASSWORD ###'
echo $PD
echo '################'
echo ''
echo '############################################'
echo '### Redis Installed In /usr/local/redis/ ###'
echo '############################################'

redis-server -v
echo '################'

# 查询redis进程
ps -ef | grep -v grep | grep redis-server
```

