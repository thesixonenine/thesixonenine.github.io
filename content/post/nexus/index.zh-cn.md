---
title: "nexus"
date: 2025-09-11T16:59:26
lastmod: 2025-09-17T14:35:26
categories: ['Java']
keywords: nexus
description: Using Nexus
---

## 安装

### Docker

```shell
docker run -d --name nexus -p 8081:8081 -u root --privileged=true sonatype/nexus3:latest
```

### 获取密码

```shell
docker exec nexus cat /nexus-data/admin.password
```

### 登录并改密码

[http://localhost:8081](http://localhost:8081)

## 仓库

仓库类别(type)分为三类

- group 聚合
- proxy 代理
- hosted 本地

## 上传本地仓库

### 新建的 hosted 类型仓库

1. 新建 maven2(hosted) 类型仓库
2. 填写 Name 唯一标识(例如 **example**)
3. **Version policy** 选择 **Mixed**
4. **Deployment Policy** 选择 **Allow redeploy**

### 将仓库加入聚合仓库

1. 选择默认的 **maven-public** 仓库
2. 在 **Member repositories** 中将新建的 **example** 仓库 从 **Available** 移到 **Members** 中并调整顺序

### 上传本地仓库的 jar 包到仓库中

**mavenupload.sh**

```shell
#!/bin/bash
# Get command line params
while getopts ":r:u:p:" opt; do
  case $opt in
    r) REPO_URL="$OPTARG"
    ;;
    u) USERNAME="$OPTARG"
    ;;
    p) PASSWORD="$OPTARG"
    ;;
esac
done

find . -type f -not -path './mavenupload\.sh*' -not -path '*/\.*' -not -path '*/\^archetype\-catalog\.xml*' -not -path '*/\^maven\-metadata\-local*\.xml' -not -path '*/\^maven\-metadata\-deployment*\.xml' | sed "s|^\./||" | xargs -I '{}' curl -s -u "$USERNAME:$PASSWORD" -X PUT -v -T {} ${REPO_URL}/{} ;
```

**上传本地仓库到 example 仓库中**

```shell
cd ~/.m2/repository
./mavenupload.sh -u admin -p 123456 -r http://localhost:8081/repository/example/
```

## 配置仓库代理

```xml
<settings>
  <mirrors>
    <mirror>
      <id>proxy</id>
      <mirrorOf>*</mirrorOf>
      <name>proxy</name>
      <url>http://localhost:8081/repository/public</url>
    </mirror>
  </mirrors>
</settings>
```

## 三方包

### 本地安装

必填参数包括 **file**, **groupId**, **artifactId**, **version**

选填参数 **maven.repo.local**

```shell
mvn install:install-file -Dmaven.repo.local= -DgroupId= -DartifactId= -Dversion= -Dfile= -Dpackaging=jar -DgeneratePom=true
```

eg:

```shell
mvn install:install-file "-Dmaven.repo.local=C:\Users\simple\.m2\repository" "-DgroupId=com.sankuai.sjst" "-DartifactId=MtOpJavaSDK" "-Dversion=1.0-SNAPSHOT" "-Dfile=C:\Users\simple\Desktop\MtOpJavaSDK-1.0-SNAPSHOT.jar" "-Dpackaging=jar" "-DgeneratePom=true"
```

### 部署到 Nexus

必填参数包括 **file**, **groupId**, **artifactId**, **version**, **repositoryId**, **url**

**repositoryId** 对应 settings.xml 中 server 节点的下的 id

**url** 指向 Hosted 仓库

```shell
mvn deploy:deploy-file -DgroupId= -DartifactId= -Dversion= -Dfile= -DrepositoryId= -Durl= -Dpackaging=jar -DgeneratePom=true
```

eg:

```shell
mvn deploy:deploy-file \
  -Dfile=C:\Users\simple\Desktop\MtOpJavaSDK-1.0-SNAPSHOT.jar \
  -DgroupId=com.sankuai.sjst \
  -DartifactId=MtOpJavaSDK \
  -Dversion=1.0-SNAPSHOT \
  -Dpackaging=jar \
  -DgeneratePom=true \
  -DrepositoryId=nexus \
  -Durl=http://localhost:8081/repository/my-hosted-repo/
```

自定义 Maven 的配置文件并在命令中追加参数 --settings /path/to/nexus-settings.xml

```xml
<settings xmlns="http://maven.apache.org/SETTINGS/1.0.0"
          xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
          xsi:schemaLocation="http://maven.apache.org/SETTINGS/1.0.0 https://maven.apache.org/xsd/settings-1.0.0.xsd">
  <servers>
    <server>
      <id>nexus</id>
      <username>Nexus 账号</username>
      <password>Nexus 密码</password>
    </server>
  </servers>
</settings>
```

eg:

```powershell
mvn deploy:deploy-file `
  "-Dfile=C:\Users\simple\Desktop\MtOpJavaSDK-1.0-SNAPSHOT.jar" `
  "-DgroupId=com.sankuai.sjst" `
  "-DartifactId=MtOpJavaSDK" `
  "-Dversion=1.0-SNAPSHOT" `
  "-Dpackaging=jar" `
  "-DgeneratePom=true" `
  "-DrepositoryId=nexus" `
  "-Durl=http://localhost:8081/repository/my-hosted-repo/" `
  --settings "C:\Users\simple\Desktop\nexus-settings.xml"
```
