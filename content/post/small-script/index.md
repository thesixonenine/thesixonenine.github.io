---
title: small-script
date: 2021-09-17T15:14:17
lastmod: 2025-07-01T16:26:53
categories: ['Script']
keywords: Script
description: 一些实用的脚本命令
---

## adb

### 查询手机的CPU架构

```shell
adb shell getprop ro.product.cpu.abi
```

## PowerShell

### 文件 MD5 生成与校验

#### 生成

```powershell
Get-FileHash -Algorithm MD5 text.txt | ForEach-Object { "$($_.Hash) text.txt" } > text.txt.md5
```

#### 校验

```powershell
if ((Get-FileHash -Algorithm MD5 text.txt).Hash -eq (Get-Content text.txt.md5).Split(" ")[0]) { "text.txt: OK" } else { "text.txt: FAILED" }
```

### 更新 lastmod

```powershell
(Get-Content /path/to/file) -replace '^lastmod: .*', "lastmod: $(Get-Date -Format 'yyyy-MM-ddTHH:mm:ss')" | Set-Content /path/to/file

function update-lastmod {
    param (
        [string]$FilePath
    )
    (Get-Content $FilePath) -replace '^lastmod: .*', "lastmod: $(Get-Date -Format 'yyyy-MM-ddTHH:mm:ss')" | Set-Content $FilePath
}
```

### winget 走代理

[winget替换源](https://mirrors.ustc.edu.cn/help/winget-source.html)

```powershell
# 管理员权限下开启
winget settings --enable ProxyCommandLineOptions
# 安装走代理
winget install JanDeDobbeleer.OhMyPosh -s winget --proxy http://127.0.0.1:10809
```

### 删除文件(目录)

```powershell
if (Test-Path .\dist.zip) { Remove-Item .\dist.zip; Write-Output "dist.zip 文件已删除" } else { Write-Output "dist.zip 文件不存在" }
if (Test-Path .\dist -PathType Container) { Remove-Item .\dist -Recurse; Write-Output "dist 目录已删除" } else { Write-Output "dist 目录不存在" }
```

### Redis 批量删除key

```powershell
redis-cli -a PASSWORD -p PORT -h HOST -n DB keys 'KEY*' |`
ForEach-Object {redis-cli -a PASSWORD -p PORT -h HOST -n DB del $_}
```

### 计算文件的哈希

支持的哈希函数: `MD5` `SHA1` `SHA256` `SHA384` `SHA512`

```powershell
Get-FileHash -Algorithm MD5 -Path .\filename.txt | Select-Object Hash
```

### 字符串编码解码

单/双引号 字符串的区别
1. 双引号中的`$`开头的变量或转义字符`` ` ``都会被处理, 单引号中的不会
2. 双引号中输出双引号, 则需要两个双引号, 即`""""`, 单引号输出单引号则需要两个单引号
3. 双引号中输出单引号或者单引号中输出双引号, 直接写即可

```powershell
# Base64
[System.Convert]::ToBase64String([System.Text.Encoding]::UTF8.GetBytes('Ab123!@#$%^&*()`'))
[System.Text.Encoding]::UTF8.GetString([System.Convert]::FromBase64String('QWIxMjMhQCMkJV4mKigpYA=='))
```

### 字符串对比

```powershell
Compare-Object -CaseSensitive "ABC" "abc"
```

### 类似 `tail -f` 监控文件变化并输出

```powershell
Get-Content -Path "filename.txt" -Wait
```

### 修改 `JAVA_HOME` 变量

```powershell
param($s)
# "Give permissions to HKLM\System\CurrentControlSet\Control\Session Manager\Environment to a desired user"
# [System.Environment]::SetEnvironmentvariable("JAVA_HOME", "C:\Program Files\Java\jdk-11.0.10", "Machine")
if ($s -eq 11) {
    [environment]::SetEnvironmentvariable('JAVA_HOME', 'C:\Program Files\Java\jdk-11.0.10', 'Machine')
}
if ($s -eq 8) {
    [environment]::SetEnvironmentvariable('JAVA_HOME', 'C:\Program Files\Java\jdk1.8.0_261', 'Machine')
}
[System.Environment]::GetEnvironmentvariable("JAVA_HOME", "Machine")
```

### 生成 `UUID`

```powershell
[System.Guid]::NewGuid().toString()
# [System.Guid]::NewGuid().toString("B")
# [System.Guid]::NewGuid().toString("B").toLower()
# [System.Guid]::NewGuid().toString("B").toUpper()
```

### 查询所有 `WIFI` 信息密码

```powershell
# 查询所有 WIFI 名称
netsh wlan show profiles | Where-Object {$_ -match 'All User Profile'}
# 查询指定 WIFI 的密码
netsh wlan show profile name="WIFI_NAME" key=clear |`
Where-Object {$_ -match 'Key Content'}
# 查询所有 WIFI 的名称及密码
netsh wlan show profiles |`
Where-Object {$_ -match 'All User Profile'} |`
foreach {$_.Substring(27) + ""} |`
foreach {Write-Host "WIFI Name: $_";netsh wlan show profile name="$_" key=clear |`
Where-Object {$_ -match 'Key Content'} |`
foreach {Write-Host "Password :" $_.Substring(29) "`n"}}
```

### 禁用 `Ctrl` + `Space` 切换输入法

将脚本改文件后缀为reg并双击导入

```
Windows Registry Editor Version 5.00

[HKEY_CURRENT_USER\Control Panel\Input Method\Hot Keys\00000010]
"Key Modifiers"=hex:00,c0,00,00
"Target IME"=hex:00,00,00,00
"Virtual Key"=hex:ff,00,00,00

[HKEY_CURRENT_USER\Control Panel\Input Method\Hot Keys\00000070]
"Key Modifiers"=hex:00,c0,00,00
"Target IME"=hex:00,00,00,00
"Virtual Key"=hex:ff,00,00,00

[HKEY_USERS\.DEFAULT\Control Panel\Input Method\Hot Keys\00000010]
"Key Modifiers"=hex:00,c0,00,00
"Target IME"=hex:00,00,00,00
"Virtual Key"=hex:ff,00,00,00

[HKEY_USERS\.DEFAULT\Control Panel\Input Method\Hot Keys\00000070]
"Key Modifiers"=hex:00,c0,00,00
"Target IME"=hex:00,00,00,00
"Virtual Key"=hex:ff,00,00,00
```

### `aria2.conf` 更新 `bt-tracker`

```powershell
$ConfigFile = "C:\aria2-1.35.0-win-64bit-build1\aria2_auto_update.conf"
$TrackersFile = "trackers_best.txt"
$DownloadLink = "https://raw.githubusercontent.com/ngosang/trackerslist/master/$TrackersFile"

Invoke-WebRequest -Uri $DownloadLink -OutFile $env:TEMP\$TrackersFile
$TrackersStream = (Get-Content $env:TEMP\$TrackersFile -Raw).Replace("`n`n", ",").Insert(0, "bt-tracker=")
$TrackersStream = $TrackersStream.Substring(0, $TrackersStream.Length - 1)
$ExcludeLineNum=(Select-String -Path $ConfigFile -SimpleMatch "bt-tracker=").LineNumber
$ConfigStream = Get-Content $ConfigFile -Encoding UTF8
$ConfigStream[$ExcludeLineNum-1]=$TrackersStream
Set-Content -Path $ConfigFile -Value $ConfigStream -Encoding UTF8
Remove-Item -Path $env:TEMP\trackers*
```

### 安装 `Firefox`

```powershell
# 下载最新版Firefox
Invoke-WebRequest -o ./ff-installer.exe 'https://download.mozilla.org/?product=firefox-latest&os=win64&lang=zh-CN'
# 安装
./ff-installer.exe
```

### 查询天气

```bash
curl wttr.in
```

### 查询 IP

```bash
echo "My public IP address is: $(curl -s https://myip.ipip.net)"
```

### 测试指定端口是否开放

```powershell
"HOST:PORT", "HOST:PORT" | % { $h, $p = $_.split(':'); $socket = New-Object Net.Sockets.TcpClient; $socket.ReceiveTimeout = 1500; try { $socket.Connect($h, $p); Write-Host "${h}:${p} is open" } catch { Write-Host "${h}:${p} is closed" }; $socket.Close() }
```

### 定时清除redis的指定前缀的key

```powershell
while ($true) { redis-cli -a PASSWORD -p PORT -h HOST -n 0 keys 'w*' | ForEach-Object {redis-cli -a PASSWORD -p PORT -h HOST -n 0 del $_} ; Start-Sleep -s 60 }
```

### 启动 Windows Terminal 时指定参数

```powershell
wt new-tab -p 'local' --title 'default' `; new-tab -p 'local' -d C:\Users\simple\Documents --title 'Documents' --tabColor '#07c160' `; new-tab -p 'local' -d C:\Users\simple\Desktop --title 'Desktop' --tabColor '#fa5151' `;focus-tab -t 0
```

### 统计当前目录下所有文件类型及数量

```powershell
Get-ChildItem -Recurse -File | Where-Object { $_.DirectoryName -notmatch '[\\/]\.(git|idea)$' -and $_.DirectoryName -notmatch '^\.(git|idea)[\\/]' } | ForEach-Object { $_.Extension } | Group-Object | Select-Object Name, Count | Sort-Object Name | ForEach-Object { if ($_.Name) { $_.Name.TrimStart('.') + " " + $_.Count } else { "<noext> " + $_.Count } }
```

## Bash

### 检查 SSH 公钥信息

```shell
ssh-keygen -l -f ~/.ssh/id_rsa.pub
ssh-keygen -E md5 -l -f ~/.ssh/id_rsa.pub
```

### 文件 MD5 生成与校验

#### 生成

```shell
md5sum text.txt > text.txt.md5
```

#### 校验

```shell
md5sum -c text.txt.md5
```

### 更新 lastmod

```shell
sed -i "/^lastmod: /s/^lastmod: .*/lastmod: $(date '+%Y-%m-%dT%H:%M:%S')/" /path/to/file

function update-lastmod {
    sed -i "/^lastmod: /s/^lastmod: .*/lastmod: $(date '+%Y-%m-%dT%H:%M:%S')/" "${1}"
}
```

### B站循环推流本地视频文件

```bash
# nasm, x264
url='rtmp://live-push.bilivideo.com/live-bvc/'
passwd='?streamname=&key=&schedule=rtmp&pflag=1'
file='/root/640_360.flv'
nohup /usr/local/bin/ffmpeg -re -stream_loop 3 -i ${file} -vcodec libx264 -acodec aac -f flv ${url}${passwd} 1>/dev/null 2>&1 &
```

### 添加 Tomcat 为 Systemd 服务

```bash
SERVICE_NAME=tomcat-service
TOMCAT_PATH=/path/to/tomcat
JAVA_PATH=/path/to/java/home
RUN_USER=root
RUN_GROUP=root

cat>/etc/systemd/system/${SERVICE_NAME}.service<<EOF
[Unit]
Description=${SERVICE_NAME}
After=network.target

[Service]
Type=forking
TimeoutSec=0
Environment=CATALINA_HOME=${TOMCAT_PATH}
Environment=CATALINA_BASE=${TOMCAT_PATH}
Environment=CATALINA_PID=${TOMCAT_PATH}/temp/tomcat.pid
Environment=JAVA_HOME=${JAVA_PATH}

ExecStart=${TOMCAT_PATH}/bin/startup.sh
ExecStop=${TOMCAT_PATH}/bin/shutdown.sh
ExecReload=/bin/kill -s HUP \$MAINPID

User=${RUN_USER}
Group=${RUN_GROUP}
[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable ${SERVICE_NAME}.service
systemctl status ${SERVICE_NAME}.service
systemctl list-unit-files | grep ${SERVICE_NAME}.service
systemctl list-units --type=service | grep tomcat
# systemctl start ${SERVICE_NAME}.service
```

### 添加 java -jar 为 Systemd 服务

```java
# 以下两种写法均可 https://docs.spring.io/spring-boot/reference/features/external-config.html#features.external-config.application-json
# system property
java -Xms128m -Dspring.profiles.active=prod -jar app.jar
# command line argument
java -Xms128m -jar app.jar --spring.profiles.active=prod
```

```bash
SERVICE_NAME=jar-file-manage
JAR_DIR=/path/to/jar/home/
JAR_NAME=file-manage.jar
SERVER_PORT=8210
JAR_ARGS="--spring.profiles.active=test --server.port=${SERVER_PORT}"
JVM_ARGS="-Xms256m -Xmx1024m"
JAVA_PATH=/path/to/java/home
RUN_USER=root
RUN_GROUP=root

cat>/etc/systemd/system/${SERVICE_NAME}.service<<EOF
[Unit]
Description=${SERVICE_NAME}
After=network.target

[Service]
Type=simple
User=${RUN_USER}
Group=${RUN_GROUP}
WorkingDirectory=${JAR_DIR}
ExecStart=${JAVA_PATH}/bin/java -jar ${JVM_ARGS} ${JAR_DIR}${JAR_NAME} ${JAR_ARGS}
ExecReload=/bin/kill -s HUP \$MAINPID
ExecStop=/bin/kill -s QUIT \$MAINPID
SuccessExitStatus=143
Restart=on-failure
[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable ${SERVICE_NAME}.service
systemctl status ${SERVICE_NAME}.service
systemctl list-unit-files | grep ${SERVICE_NAME}.service
systemctl list-units --type=service | grep jar
# systemctl start ${SERVICE_NAME}.service
```

### Redis 批量删除key

```bash
redis-cli -a PASSWORD -p PORT -h HOST -n DB keys 'KEY*' |\
xargs redis-cli -a PASSWORD -p PORT -h HOST -n DB del
```

### 查询目录和文件

#### `find` 命令

[参考](https://wangchujiang.com/linux-command/c/find.html)

##### 基本使用

```bash
# find 路径 参数
# "." 代表当前目录
# "-type f" 代表要查找的是普通文件
# "-name '*.out'" 代表按文件名称查找, 且是以".out"结尾的文件
find . -type f -name '*.out'
```

常用命令

| 参数     | 解释                                                         | 示例                                                         |
| -------- | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `-exec`  | 对匹配的文件执行该参数所给出的其他linux命令                  | 查找当前目录下所有.txt文件并拼接写入到all.txt文件<br> `find . -type f -name "*.txt" -exec cat {} \;> /all.txt` |
| `-ok`    | 同 `-exec`, 在执行命令前会确认                               |                                                              |
| `-type`  | 查找某一类型的文件                                           | `b` - 块设备文件; `d` - 目录; `c` - 字符设备文件; `p` - 管道文件;<br> `l` - 符号链接文件; `f` - 普通文件; `s` - socket文件 |
| `-name`  | 按文件名称查找                                               |                                                              |
| `-perm`  | 按文件权限查找                                               |                                                              |
| `-mtime` | -mtime -n +n 按**更改时间**查找<br>-n表示距现在n天以内<br>+n表示距现在n天以前 | **最后访问时间**: `-atime`(天) `-amin`(分)<br>**最后修改时间**: `-mtime`(天) `-mmin`(分)<br>**数据元(权限等)最后修改时间**: `-ctime`(天) -cmin(分) |

### 压缩/解压

#### `tar` 命令

##### 基本使用

```bash
# 压缩为*.tar.gz
tar -zcvf foldername.tar.gz ./foldername
# *.tar.gz解压缩
tar -zxvf foldername.tar.gz
```

##### 命令详细解释

| 参数                     | 解释                                  |
| ------------------------ | ------------------------------------- |
| `-z`                     | filter the archive through gzip       |
| `-c` `--create`          | create a new archive                  |
| `-x` `--extract` `--get` | extract files from an archive         |
| `-v` `--verbose`         | verbosely list files processed        |
| `-f` `--file`            | use archive file or device ARCHIVE    |
| `-r` `--append`          | append files to the end of an archive |
| `-t` `--list`            | list the contents of an archive       |

###  `aria2` 更新 `bt-tracker`

```bash
#!/bin/bash
list=`wget -qO- https://raw.githubusercontent.com/ngosang/trackerslist/master/trackers_all.txt|awk NF|sed ":a;N;s/\n/,/g;ta"`
if [ -z "`grep "bt-tracker" /conf/aria2.conf`" ]; then
    sed -i '$a bt-tracker='${list} /conf/aria2.conf
else
    sed -i "s@bt-tracker.*@bt-tracker=$list@g" /conf/aria2.conf
fi
```

### Windows 重启网络

```bash
net stop winnat
net start winnat
```

### 去除视频开头30秒

```bash
ffmpeg -ss 00:00:30 -i input.mp4 -c:v copy -c:a copy output.mp4
```

### ffmpeg循环视频推流

```bash
nohup ffmpeg -re -stream_loop 3 -i FileName.flv -vcodec libx264 -acodec aac -f flv rtmp://live-push.bilivideo.com/live-bvc/?streamname=&key=&schedule=rtmp&pflag=1 1>/dev/null 2>&1 &
```

### 多个文件中查询指定字符串并输出文件名

```bash
find . -type f -name "*.out" | xargs grep "220728115133646916" -l
```

### 替换多个目录中git仓库的地址

```bash
# 将 192.168.1.1 替换成 192.168.1.2
sed -i "s/192.168.1.1/192.168.1.2/g" `grep "192.168.1.1" -rl ./*/.git/config`
```

### adoc转docx

```cmd
@rem 使用asciidoctor转换adoc为docbook格式
asciidoctorj -b docbook README.adoc -o README.xml
@rem 使用pandoc转换为docx格式
pandoc -f docbook -t docx -o README.docx README.xml
```

### 测试指定端口是否开放

```bash
for hostport in HOST:PORT HOST:PORT; do (echo >/dev/tcp/${hostport/:/\/}) >/dev/null 2>&1 || echo "${hostport} is closed"; done
```
