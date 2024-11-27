---
title: rabbitmq-install-config
date: 2021-12-17T17:50:18+0800
lastmod: 2024-11-27T14:53:18+0800
tags: ['Linux']
categories: ['RabbitMQ']
keywords: rabbitmq
description: rabbitmq安装配置
---

## Yum安装

```bash
#!/bin/bash

# 安装rabbitmq-server
yum install epel-release -y
yum install rabbitmq-server -y
# 查看防火墙是否开启
# systemctl status firewalld.service
# systemctl status iptables.service
# 开放端口
# firewall-cmd --zone=public --add-port=15672/tcp --permanent
# firewall-cmd --reload

# 启用网页插件
rabbitmq-plugins enable rabbitmq_management
# 设置开机启动rabbitmq
systemctl enable rabbitmq-server
# 启动rabbitmq
systemctl start rabbitmq-server

# 状态查看
rabbitmqctl status

# 版本查看
rabbitmqctl status | grep RabbitMQ

# 修改默认的账户guest的密码
rabbitmqctl change_password guest guest123
# 新增用户
rabbitmqctl add_user mytest mytest123
# 删除用户
# rabbitmqctl delete_user mytest

# 默认角色
# administrator
# monitoring
# policymaker
# management

# 设置用户的角色
rabbitmqctl set_user_tags mytest administrator

# 查询所有用户的权限
# rabbitmqctl list_permissions
# 查看virtual host为/的所有用户权限
# rabbitmqctl list_permissions -p /
# 查询指定用户的权限
# rabbitmqctl list_user_permissions mytest
# 清除用户权限
# rabbitmqctl clear_permissions mytest

# 设置用户可以访问的virtual host
rabbitmqctl set_permissions -p / mytest ".*" ".*" ".*"
```

## Docker 安装

### Dockerfile

```bash
FROM rabbitmq:3.13.2-management
WORKDIR /plugins
# 在release页面查看对应插件版本并下载
# https://github.com/rabbitmq/rabbitmq-delayed-message-exchange/releases
COPY rabbitmq_delayed_message_exchange-3.13.0.ez ./
RUN rabbitmq-plugins enable rabbitmq_delayed_message_exchange

CMD ["rabbitmq-server"]
```

### Build Image

```bash
docker build -t thesixonenine/rabbitmq-delayed:3.13.2-management .
```

### Run Container

```bash
docker run -d --name=rabbitmq-delayed -p 15672:15672 -p 5672:5672 -e RABBITMQ_DEFAULT_USER=admin -e RABBITMQ_DEFAULT_PASS=admin123 thesixonenine/rabbitmq-delayed:3.13.2-management
```

## 离线安装特定版本

[版本对应](https://v3-12.rabbitmq.com/which-erlang.html)

```bash
# 下载 rabbitmq-server
# https://github.com/rabbitmq/rabbitmq-server/releases/download/v3.7.14/rabbitmq-server-generic-unix-3.7.14.tar.xz
# 下载对应版本的 erlang
# https://github.com/rabbitmq/erlang-rpm/releases/download/v21.3.8.21/erlang-21.3.8.21-1.el7.x86_64.rpm
rpm -ivh erlang-21.3.8.21-1.el7.x86_64.rpm
# 查看 erlang 版本
erl -version
xz -d rabbitmq-server-generic-unix-3.7.14.tar.xz
tar -xvf rabbitmq-server-generic-unix-3.7.14.tar
# 修改 hosts 文件, 追加主机名称与本机IP的关系, 主机名称使用 hostname 命令查看
# vi /etc/hosts
# 127.0.0.1 c12
# 启动 rabbitmq
./rabbitmq_server-3.7.14/sbin/rabbitmq-server -detached
# 查看状态
./rabbitmq_server-3.7.14/sbin/rabbitmqctl status
```
