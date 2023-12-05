---
title: nginx-install
date: 2022-05-05T14:26:34+0800
lastmod: 2022-05-05T14:26:34+0800
tags: ['Linux']
categories: ['Nginx']
keywords: Nginx
description: CentOS 7 Nginx安装配置
---

## Nginx 编译安装

以下命令均以 `root` 用户执行
```bash
#!/bin/bash
cd ~
# 安装wget用来下载Nginx源码
yum install -y wget
# 安装编译所需的工具
yum install -y gcc openssl openssl-devel pcre pcre-devel zlib zlib-devel
# 下载Nginx源码
wget https://nginx.org/download/nginx-1.21.6.tar.gz
# 解压Nginx源码到/usr/local/
tar -xf nginx-1.21.6.tar.gz -C /usr/local/
# 切换目录到源码目录准备开始编译安装
cd /usr/local/nginx-1.21.6
# 创建安装的目录
mkdir /usr/local/nginx
# 配置参数
./configure --prefix=/usr/local/nginx --with-http_ssl_module
# 编译
make
# 安装
make install
# 建立软链接
ln -s /usr/local/nginx/sbin/nginx /usr/bin/nginx


```

## 设置Nginx开机自启

以下命令均以 `root` 用户执行
```bash
#!/bin/bash
# 写入自启动文件
touch /lib/systemd/system/nginx.service
cat > /lib/systemd/system/nginx.service <<EOF
[Unit]
Description=nginx service
After=network.target 
   
[Service] 
Type=forking 
ExecStart=/usr/local/nginx/sbin/nginx
ExecReload=/usr/local/nginx/sbin/nginx -s reload
ExecStop=/usr/local/nginx/sbin/nginx -s quit
PrivateTmp=true 
   
[Install] 
WantedBy=multi-user.target
EOF

# 查看nginx是否在自启动列表中
systemctl list-unit-files | grep nginx
# 设置nginx开机自启
systemctl enable nginx
# 启动nginx
systemctl start nginx.service
```

## 查看开机启动项

```bash
# systemctl list-unit-files
# systemctl list-unit-files | grep enabled
```

## Nginx Yum安装

[参考官方文档](http://nginx.org/en/linux_packages.html#RHEL-CentOS)

以下命令均以 `root` 用户执行
```bash
#!/bin/bash
yum install -y yum-utils
touch /etc/yum.repos.d/nginx.repo
cat > /etc/yum.repos.d/nginx.repo << EOF
[nginx-stable]
name=nginx stable repo
baseurl=http://nginx.org/packages/centos/$releasever/$basearch/
gpgcheck=1
enabled=1
gpgkey=https://nginx.org/keys/nginx_signing.key
module_hotfixes=true
EOF

yum install -y nginx
```

## Nginx 配置目录浏览

以下命令均以 `root` 用户执行
```bash
# 如果需要帐号密码访问, 则需要生成密码文件并在配置中指定文件
yum install -y httpd
htpasswd -bc /usr/local/nginx/nginx_passwd simple 123456

# 指定目录开启
mkdir /usr/local/nginx/html/flow
# 修改server块中的配置
location /flow {
	autoindex on; # 打开目录浏览功能
	autoindex_exact_size off;
	autoindex_localtime on; # 以服务器的文件时间作为显示的时间
	charset utf-8,gbk; # 展示中文文件名
	auth_basic "need password";
	auth_basic_user_file /usr/local/nginx/nginx_passwd;
	root   /usr/local/nginx/html;
}
```
