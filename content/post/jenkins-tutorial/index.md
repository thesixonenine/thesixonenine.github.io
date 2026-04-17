---
title: jenkins-tutorial
date: 2026-04-17T15:37:34+0800
lastmod: 2026-04-17T15:37:34+0800
tags: ['Jenkins', 'Docker']
categories: ['Jenkins']
keywords: jenkins
description: jenkins 安装
---

使用容器来搭建 `Jenkins` 环境

## 环境变量

有两个环境变量可以考虑提前设置

- `jenkins.install.runSetupWizard`: 为 `false` 代表不进行引导设置
- `hudson.model.DownloadService.noSignatureCheck`: 为 `true` 代表不检查插件签名

## 启动后再装插件

`jenkins.yml`

```yml
services:
  jenkins:
    image: jenkins/jenkins:lts-jdk21
    container_name: jenkins
    restart: always
    ports:
      - "8080:8080"
    environment:
      - TZ=Asia/Shanghai
      - JAVA_OPTS=-Xms1024m -Xmx1024m -Duser.timezone=Asia/Shanghai -Dfile.encoding=UTF-8 -Djenkins.install.runSetupWizard=false -Dhudson.model.DownloadService.noSignatureCheck=true
    # 宿主机docker组GID: stat -c '%g' /var/run/docker.sock
    group_add:
      - "989"
    volumes:
      - jenkins_home:/var/jenkins_home
      - /usr/bin/docker:/bin/docker
      - /var/run/docker.sock:/var/run/docker.sock
volumes:
  jenkins_home:
```

## 自定义镜像装好插件再启动

`plugins.txt`

```text
git:latest
localization-zh-cn:latest
```

`Dockerfile`

```Dockerfile
FROM jenkins/jenkins:lts-jdk21
USER root
COPY plugins.txt /usr/share/jenkins/ref/plugins.txt
RUN jenkins-plugin-cli --plugin-file /usr/share/jenkins/ref/plugins.txt
USER jenkins
```

`jenkins.yml`

```yml
services:
  jenkins:
    build:
      context: .
      dockerfile: Dockerfile
    image: thesixonenine/jenkins:lts-jdk21
    container_name: jenkins
    restart: always
    ports:
      - "8080:8080"
    environment:
      - TZ=Asia/Shanghai
      - JAVA_OPTS=-Xms1024m -Xmx1024m -Duser.timezone=Asia/Shanghai -Dfile.encoding=UTF-8 -Djenkins.install.runSetupWizard=false -Dhudson.model.DownloadService.noSignatureCheck=true
    # 宿主机docker组GID: stat -c '%g' /var/run/docker.sock
    group_add:
      - "989"
    volumes:
      - jenkins_home:/var/jenkins_home
      - /usr/bin/docker:/bin/docker
      - /var/run/docker.sock:/var/run/docker.sock
volumes:
  jenkins_home:
```

如果是启动并进行引导设置, 即没有设置 `jenkins.install.runSetupWizard=false`, 则需要查看默认admin的密码

```shell
docker exec jenkins cat /var/jenkins_home/secrets/initialAdminPassword
```


## 替换插件源

`查看`

```shell
docker exec jenkins cat /var/jenkins_home/hudson.model.UpdateCenter.xml
```

`替换`, 使用华为源(`https://mirrors.huaweicloud.com/jenkins/updates/update-center.json`)

```shell
docker exec jenkins sed -i 's/https:\/\/updates.jenkins.io/https:\/\/mirrors.huaweicloud.com\/jenkins\/updates/g' /var/jenkins_home/hudson.model.UpdateCenter.xml
```
