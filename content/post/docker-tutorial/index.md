---
title: docker-tutorial
date: 2022-03-22T18:27:00+0800
lastmod: 2025-10-21T14:42:02+0800
tags: ['Linux']
categories: ['Docker']
keywords: docker
description: 
---

## docker install

```bash
#!/bin/bash
sudo yum remove docker docker-client docker-client-latest \
	docker-common docker-latest docker-latest-logrotate \
	docker-logrotate docker-engine
sudo yum install -y yum-utils
sudo yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
sudo sed -i 's+download.docker.com+mirrors.ustc.edu.cn/docker-ce+' /etc/yum.repos.d/docker-ce.repo
sudo yum makecache fast
# 安装最新版并创建docker组(里面没有用户)
sudo yum install -y docker-ce docker-ce-cli containerd.io

# 安装指定版本
# 列出仓库可用的所有版本
# yum list docker-ce --showduplicates | sort -r
# 对于3:20.10.6-3.el8，它的版本号是:到-中间的部分，也就是: 20.10.6
# 将<VERSION_STRING>替换为20.10.6即可
# sudo yum install docker-ce-<VERSION_STRING> docker-ce-cli-<VERSION_STRING> containerd.io

# 将用户添加到docker组
sudo usermod -aG docker ${USER}
# 开机启动
sudo systemctl enable docker
# start docker
sudo systemctl start docker
sudo chmod a+rw /var/run/docker.sock
# update registry mirrors
sudo touch /etc/docker/daemon.json
sudo cat > /etc/docker/daemon.json <<EOF
{
    "registry-mirrors": [
        "https://dockerproxy.com",
		"https://docker.nju.edu.cn",
		"https://docker.mirrors.sjtug.sjtu.edu.cn",
		"https://hub-mirror.c.163.com",
        "https://mirror.baidubce.com",
        "https://cr.console.aliyun.com"
    ]
}
EOF
sudo systemctl restart docker
# hello-world
docker run hello-world
```

### mysql install

```shell
docker run -d --name mysql \
    --network host \
    --restart unless-stopped \
    -v /opt/mysql/data:/var/lib/mysql \
    -e MYSQL_ROOT_PASSWORD=123456 \
    mysql:8.2.0 \
    --innodb-dedicated-server=ON \
    --group-replication-consistency=AFTER \
    --transaction-isolation=READ-COMMITTED \
    --lower_case_table_names=1
```

### delete unused data

```shell
# delete unused none image
docker rmi $(docker images -f "dangling=true" -q)
# powershell
docker images -f "dangling=true" -q | % { docker rmi $_ }
# delete unused image container cache
docker system prune --all
# delete unused volume
docker volume prune
```

### stop all container

```shell
docker ps -q | xargs docker stop
# powershell
docker ps -q | % { docker stop $_ }
```

## docker uninstall

```shell
# docker uninstall
sudo yum remove docker-ce docker-ce-cli containerd.io
# delete all images, containers, and volumes
sudo rm -rf /var/lib/docker
sudo rm -rf /var/lib/containerd
```

## docker compose install

```shell
sudo yum install docker-compose-plugin
```

## portainer install

```shell
# portainer install
docker volume create portainer_data
docker run -d -p 8000:8000 -p 9000:9000 -p 9443:9443 \
	--name portainer --restart=always \
	-v /var/run/docker.sock:/var/run/docker.sock \
	-v portainer_data:/data \
	portainer/portainer-ce:2.11.1
```
