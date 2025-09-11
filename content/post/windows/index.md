---
title: "windows"
date: 2023-07-21T10:14:26
lastmod: 2025-09-11T15:16:40
categories: ['Windows']
keywords: windows
description: Using Windows
# image: "https://msdesign.blob.core.windows.net/wallpapers/Microsoft_Nostalgic_Windows_Wallpaper_4k.jpg"
---


## Video References

<!-- {{< bilibili BV1dxT6zGESE >}} -->

<div class="video-wrapper">
<iframe src='https://player.bilibili.com/player.html?as_wide=1&high_quality=1&page=1&bvid=BV1dxT6zGESE&autoplay=0' scrolling='no' frameborder='no' framespacing='0' allowfullscreen='true'></iframe>
</div>

Using Windows from a Developer's Perspective

## Install Windows

### Reset Windows

- `win + i` 打开设置.
- 输入 `重置此电脑` 并按 `Enter` 确认.
- 点击 `开始`, 按提示操作即可.

### Install Windows From USB

[参考](https://mirrors.sdu.edu.cn/docs/guide/Windows-iso)

- 微软官方下载
  - 使用**任意设备**打开 [Windows 11](https://www.microsoft.com/zh-cn/software-download/windows11/)
  - 使用**除Windows**以外的设备打开 [Windows 10](https://www.microsoft.com/zh-cn/software-download/windows10ISO/)
- 使用 [Aria2](https://github.com/aria2/aria2) 下载 Windows 镜像, 推荐 `Windows 10/11 专业版`.
- 使用 [Ventoy](https://github.com/ventoy/Ventoy) 来制作启动盘.
- 将U盘插入电脑, 开机进入BIOS选择U盘启动, 按提示操作即可.

ISO文件Hash校验(SHA256)

```poweshell
if ((Get-FileHash -Algorithm SHA256 -Path "PATH/TO/FILE").Hash -ne "HASH") { Write-Host "文件校验失败" } else { Write-Host "文件校验成功" }
```

Windows11 跳过网络验证

- Shift + F10 打开 `CMD` 窗口并输入 regedit
- 定位到 `HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\OOBE`
- 新建DWORD(32位)值,名称为 `BypassNRO`, 值设置为1
- 命令行输入 `logoff` 或 `shutdown /r /t 0`

## System Setting

### Manual

- 桌面图标仅保留回收站

### Postpone Windows Updates

```powershell
reg add "HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\WindowsUpdate\UX\Settings" /v FlightSettingsMaxPauseDays /t reg_dword /d 10000 /f
# Type "Check for updates" in the search box and enter. Find "Pause updates," then click the dropdown menu to select a time.
```

### Install Script

在新安装的 Windows 系统上配置软件以便快速回到自己熟悉的开发环境, 特此记录以下脚本

软件安装统一使用 [Scoop](https://github.com/ScoopInstaller/Scoop) 来安装软件, 参考另一篇[博客](https://thesixonenine.site/p/scoop.html)

#### Install **Scoop**

打开自带的 **Windows PowerShell** 并执行以下命令

```
Set-ExecutionPolicy RemoteSigned -Scope CurrentUser
iwr -useb get.scoop.sh | iex

# 国内镜像
# iwr -useb https://gitee.com/glsnames/scoop-installer/raw/master/bin/install.ps1 | iex
# scoop config SCOOP_REPO 'https://gitee.com/glsnames/scoop-installer'
# https://gitee.com/scoop-bucket
# scoop bucket add extras https://gitee.com/scoop-bucket/extras.git
```

#### Install Software

```
scoop install git
scoop bucket add dorado https://github.com/chawyehsu/dorado.git
scoop install dorado/clash-for-windows
# 手动配置clash
scoop config proxy 127.0.0.1:10808
scoop config SCOOP_REPO "https://github.com/ScoopInstaller/Scoop"
scoop update
scoop install dorado/powershell

scoop bucket add extras
scoop bucket add java
scoop bucket add nerd-fonts https://github.com/matthewjberger/scoop-nerd-fonts
scoop bucket add thesixonenine-scoop-bucket https://github.com/thesixonenine/thesixonenine-scoop-bucket
scoop bucket add versions https://github.com/ScoopInstaller/Versions

if ([System.Environment]::OSVersion.Version.Major -eq 10) { 
    Write-Host "Windows 10, 需要安装windows-terminal" 
    scoop install extras/windows-terminal
} else if ([System.Environment]::OSVersion.Version.Major -eq 11) { 
    Write-Host "Windows 11, 不需要安装windows-terminal" 
}

scoop install oh-my-posh neofetch go gradle hugo-extended maven openssl proxychains python scrcpy
scoop install extras/powertoys extras/posh-git extras/git-aliases extras/scoop-completion extras/gpg4win extras/jetbrains-toolbox extras/filezilla extras/carnac extras/dismplusplus extras/everything extras/fiddler extras/geekuninstaller extras/jd-gui extras/openark extras/switchhosts extras/trafficmonitor extras/vscode extras/wireshark
Install-Module -Name DirColors -Proxy "127.0.0.1:10808"
```

Office 安装使用 `office-tool-plus` 或者 自定义部署配置文件部署安装

使用 `office-tool-plus` 部署
```
scoop install extras/windowsdesktop-runtime-lts extras/office-tool-plus
```

自定义部署配置文件部署
```
iwr "https://officecdn.microsoft.com/pr/wsus/setup.exe" -OutFile setup.exe
# move setup.exe into C:\
iwr "https://gist.githubusercontent.com/thesixonenine/173647918c69d9627eeb141a32d6ec57/raw/5ee850ca1fdacce442d94051fcb6f44598834093/Configuration.xml" -OutFile Configuration.xml
cd C:\
setup.exe /configure Configuration.xml
```

激活 Windows/Office 使用开源的[MAS](https://github.com/massgravel/Microsoft-Activation-Scripts)或者闭源的[HEU](https://github.com/zbezj/HEU_KMS_Activator)

**MAS激活**

```
irm https://get.activated.win | iex
```

配置同步


## Use WSL2

[安装参考](https://learn.microsoft.com/en-us/windows/wsl/install#install-wsl-command)

[环境配置](https://learn.microsoft.com/zh-cn/windows/wsl/setup/environment)

基本命令

```shell
# 安装
wsl --install Ubuntu
# 查看版本
wsl --version
# 删除
wsl --unregister Ubuntu
# 帮助
wsl --help
```

```toml
# C:\Users\simple\.wslconfig
[wsl2]
# default: same as Windows
processors=4
# default: 50% of available RAM
memory=8GB
# default: 25% of available RAM
swap=0
# default: true
localhostForwarding=true
```

```shell
# 检查系统中 PID 1 的主初始化进程是 init 还是 systemd
ps --no-headers -o comm 1
```

### Configure Mirror

```shell
# 查看 Ubuntu 版本
lsb_release -a
# 替换源
sudo sed -i 's@//.*archive.ubuntu.com@//mirrors.ustc.edu.cn@g' /etc/apt/sources.list.d/ubuntu.sources
# 替换 security 源
sudo sed -i 's/security.ubuntu.com/mirrors.ustc.edu.cn/g' /etc/apt/sources.list.d/ubuntu.sources
# 使用 HTTPS
sudo sed -i 's/http:/https:/g' /etc/apt/sources.list.d/ubuntu.sources
# 更新索引
sudo apt update
```

### Configure .shellrc

将shell配置追加到当前使用的shell

```shell
# load my shell config
echo "if [ -f ~/.shellrc ]; then source ~/.shellrc ; fi" >> ~/.bashrc
```

新建 **.shellrc** 并加入当前使用的shell

```shell
# func
function git-proxy-set {
proxy_ip=$(ip route show | grep -i default | awk '{ print $3}')
if [ -d ./.git ]; then
    git config http.proxy http://${proxy_ip}:10809
    git config https.proxy http://${proxy_ip}:10809
else
    git config --global http.proxy http://${proxy_ip}:10809
    git config --global https.proxy http://${proxy_ip}:10809
fi
}
function git-proxy-unset {
if [ -d ./.git ]; then
    git config --unset http.proxy
    git config --unset https.proxy
else
    git config --global --unset http.proxy
    git config --global --unset https.proxy
fi
}
function git-proxy-get {
if [ -d ./.git ]; then
    git config http.proxy
    git config https.proxy
else
    git config --global http.proxy
    git config --global https.proxy
fi
}
function git-proxy-command {
    echo "git-proxy-set git-proxy-unset git-proxy-get"
}

# enable passphrase prompt for gpg
export GPG_TTY=$(tty)

# env
JAVA_HOME=/home/simple/software/jdk8u432-b06
MAVEN_HOME=/home/simple/software/apache-maven-3.9.9
PATH=$JAVA_HOME/bin:$MAVEN_HOME/bin:$PATH

BOOKMARKS=/home/simple/gitee/Bookmarks

# homebrew
# ref: https://mirrors.ustc.edu.cn/help/brew.git.html#homebrew-linuxbrew
export HOMEBREW_BREW_GIT_REMOTE='https://mirrors.ustc.edu.cn/brew.git'
export HOMEBREW_CORE_GIT_REMOTE='https://mirrors.ustc.edu.cn/homebrew-core.git'
export HOMEBREW_BOTTLE_DOMAIN='https://mirrors.ustc.edu.cn/homebrew-bottles'
export HOMEBREW_API_DOMAIN='https://mirrors.ustc.edu.cn/homebrew-bottles/api'

# alias
alias ii="explorer.exe"
```

### Use Homebrew

```shell
# echo "export HOMEBREW_BREW_GIT_REMOTE='https://mirrors.ustc.edu.cn/brew.git'" >> ~/.shellrc
# echo "export HOMEBREW_CORE_GIT_REMOTE='https://mirrors.ustc.edu.cn/homebrew-core.git'" >> ~/.shellrc
# echo "export HOMEBREW_BOTTLE_DOMAIN='https://mirrors.ustc.edu.cn/homebrew-bottles'" >> ~/.shellrc
# echo "export HOMEBREW_API_DOMAIN='https://mirrors.ustc.edu.cn/homebrew-bottles/api'" >> ~/.shellrc
# source ~/.shellrc
/bin/bash -c "$(curl -fsSL https://mirrors.ustc.edu.cn/misc/brew-install.sh)"
```

### Use apt

```shell
sudo apt install -y build-essential curl git sudo wget file software-properties-common

# 安装并切换 zsh, 安装 oh-my-zsh
sudo apt install -y zsh && chsh -s /bin/zsh && \
    git clone --depth=1 https://mirrors.tuna.tsinghua.edu.cn/git/ohmyzsh.git ~/.oh-my-zsh && \
    cp ~/.oh-my-zsh/templates/zshrc.zsh-template ~/.zshrc && \
    echo "if [ -f ~/.shellrc ]; then source ~/.shellrc ; fi" >> ~/.zshrc

# 安装 oh-my-zsh 插件
git clone --depth=1 https://github.com/zsh-users/zsh-autosuggestions ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-autosuggestions && \
    git clone --depth=1 https://github.com/zsh-users/zsh-syntax-highlighting.git ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-syntax-highlighting && \
    sed -i 's/^plugins=(git)/plugins=(\ngit\n)/' ~/.zshrc && \
    sed -i 's/^plugins=(/&\nzsh-syntax-highlighting/' ~/.zshrc && \
    sed -i 's/^plugins=(/&\nzsh-autosuggestions/' ~/.zshrc && \
    source ~/.zshrc
```

### Install Software

#### neovim & lazyvim

```shell
brew install fzf neovim
git clone https://github.com/LazyVim/starter ~/.config/nvim
```

[配置字体](https://www.nerdfonts.com/)(Hack)并在 Windows Terminal 中选定

#### Configure Git

```shell
# 安装 Git
sudo apt install git

# 查看所有配置信息
git config --list
# 查看系统级(/etc/gitconfig)配置信息
git config --system --list
# 查看用户级(~/.gitconfig)配置信息
git config --global --list
# 查看仓库级(./.git/config)配置信息
git config --local --list

# 开始配置 Git
git config --global user.name Simple
# TODO 全局配置邮箱
# git config --global user.email xxx@xxx.com
git config --global alias.ci commit
git config --global alias.st status
git config --global alias.co checkout
git config --global alias.br branch
git config --global alias.ps push
git config --global alias.pl pull
git config --global alias.ft fetch
git config --global alias.mg merge
# 修改 ~/.gitconfig
# git config --global alias.lg log --color --graph --pretty=format:'%Cred%h%Creset -%C(yellow)%d%Creset %s %Cgreen(%cr) %C(bold blue)<%an>%Creset' --abbrev-commit
# 将帐密存储在 ~/.git-credentials
git config --global credential.helper store

# 配置 GPG
# 生成 GPG 密钥
gpg --full-generate-key
# 列出 GPG 密钥
gpg --list-secret-keys --keyid-format=long
# 导入 GPG 公钥
gpg --import public.asc
# 导入 GPG 私钥
gpg --import private.asc

# TODO 全局配置 GPG 签名
git config --global user.signingkey 8E61F4E8701DD140
git config --global commit.gpgsign true

# 支持 passphrase
# echo '# enable passphrase prompt for gpg' >> ~/.bashrc
# echo 'export GPG_TTY=$(tty)' >> ~/.bashrc

# 项目上单独配置
git config user.name Simple
# TODO 项目配置邮箱
# git config user.email xxx@xxx.com
# TODO 项目配置 GPG 签名
git config user.signingkey 8E61F4E8701DD140
git config commit.gpgsign true
```


#### docker

[官方文档](https://docs.docker.com/engine/install/ubuntu)

[源安装](https://mirrors.ustc.edu.cn/help/docker-ce.html)

```shell
# 下载安装脚本
curl --proxy 172.26.112.1:10809 -fsSL https://get.docker.com -o get-docker.sh
# 执行安装
sudo DOWNLOAD_URL=https://mirrors.ustc.edu.cn/docker-ce sh get-docker.sh

# https://docs.docker.com/engine/install/linux-postinstall
# 添加用户组
sudo groupadd docker
sudo usermod -aG docker $USER
# 退出重启
logout
# 查看版本
docker version
```

配置仓库镜像,阿里云个人镜像加速地址

```shell
sudo tee /etc/docker/daemon.json <<-'EOF'
{
  "registry-mirrors": ["https://dockerproxy.net"]
}
EOF
sudo systemctl daemon-reload
sudo systemctl restart docker
```

hello world

```shell
docker run hello-world
```

阿里云个人镜像操作

```shell
# 登录阿里云 Docker Registry
docker login --username=USERNAME registry.cn-chengdu.aliyuncs.com

# 从 Registry 中拉取镜像
docker pull registry.cn-chengdu.aliyuncs.com/NAMESPACE/REPO:[镜像版本号]

# 将镜像推送到 Registry
docker tag [ImageId] registry.cn-chengdu.aliyuncs.com/NAMESPACE/REPO:[镜像版本号]
docker push registry.cn-chengdu.aliyuncs.com/NAMESPACE/REPO:[镜像版本号]
```

**docker 代理源拉取**

将 `PROXY_DOMAIN` 替换成支持的代理源, 例如:

- docker.m.daocloud.io
- dockerproxy.net

```shell
docker pull PROXY_DOMAIN/library/ubuntu:latest
docker tag PROXY_DOMAIN/library/ubuntu:latest ubuntu:latest
docker rmi PROXY_DOMAIN/library/ubuntu:latest
```

## Configure Environment

### CMD

```bash
# 查看环境变量
echo %USERPROFILE%

# 设置会话级环境变量
set JAVA_HOME="\path\to\jdk"
# 设置用户级环境变量
setx JAVA_HOME "\path\to\jdk"
# 设置系统级环境变量
setx JAVA_HOME "\path\to\jdk" /M
```

### PowerShell

```powershell
# 查看环境变量
$env:JAVA_HOME
# 或者
[System.Environment]::GetEnvironmentVariable("JAVA_HOME")
[System.Environment]::GetEnvironmentVariable("JAVA_HOME", "Process")


# 设置会话级环境变量
$env:JAVA_HOME="\path\to\jdk"
# 或者
[System.Environment]::SetEnvironmentVariable("JAVA_HOME", "\path\to\jdk", "Process")
# 设置用户级环境变量
[System.Environment]::SetEnvironmentVariable("JAVA_HOME", "\path\to\jdk", "User")
# 设置会话级环境变量
[System.Environment]::SetEnvironmentVariable("JAVA_HOME", "\path\to\jdk", "Machine")
```
