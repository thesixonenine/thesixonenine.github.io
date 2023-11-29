---
title: rabbitmq-install-config
date: 2021-12-17T17:50:18+0800
updated: 2021-12-17T17:50:18+0800
tags: ['RabbitMQ','Linux']
categories: ['RabbitMQ','Linux']
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

