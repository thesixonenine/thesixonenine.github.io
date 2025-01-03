---
title: "windows"
date: 2023-07-21T10:14:26+08:00
lastmod: 2025-01-03T10:59:26+08:00
categories: ['Windows']
keywords: windows
description: Windows相关
image: "https://msdesign.blob.core.windows.net/wallpapers/Microsoft_Nostalgic_Windows_Wallpaper_4k.jpg"
---


从一个开发者的角度来使用 `Windows`

## 系统安装

### 重置 Windows

- `win + i` 打开设置.
- 输入 `重置此电脑` 并按 `Enter` 确认.
- 点击 `开始`, 按提示操作即可.

### U盘安装

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

## 系统设置

### 手动设置

- 桌面图标仅保留回收站

### 安装脚本

在新安装的 Windows 系统上配置软件以便快速回到自己熟悉的开发环境, 特此记录以下脚本

软件安装统一使用 [Scoop](https://github.com/ScoopInstaller/Scoop) 来安装软件, 参考另一篇[博客](https://thesixonenine.site/p/scoop.html)

#### 安装 **Scoop**

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

#### 安装软件

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
irm https://massgrave.dev/get | iex
```

配置同步


## 使用 WSL2

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

### 配置源

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

### 配置 .shellrc

将shell配置追加到当前使用的shell

```shell
# load my shell config
echo "if [ -f ~/.shellrc ]; then source ~/.shellrc ; fi" >> ~/.bashrc
```

新建 **.shellrc** 并加入当前使用的shell

```shell
# func
function git_proxy_set {
proxy_ip=$(ip route show | grep -i default | awk '{ print $3}')
if [ -d ./.git ]; then
    git config http.proxy http://${proxy_ip}:10809
    git config https.proxy http://${proxy_ip}:10809
else
    git config --global http.proxy http://${proxy_ip}:10809
    git config --global https.proxy http://${proxy_ip}:10809
fi
}
function git_proxy_unset {
if [ -d ./.git ]; then
    git config --unset http.proxy
    git config --unset https.proxy
else
    git config --global --unset http.proxy
    git config --global --unset https.proxy
fi
}
function git_proxy_get {
if [ -d ./.git ]; then
    git config http.proxy
    git config https.proxy
else
    git config --global http.proxy
    git config --global https.proxy
fi
}
function git_proxy_command {
    echo "git_proxy_set git_proxy_unset git_proxy_get"
}

# enable passphrase prompt for gpg
export GPG_TTY=$(tty)

# env
JAVA_HOME=/home/simple/software/jdk8u432-b06
MAVEN_HOME=/home/simple/software/apache-maven-3.9.9
PATH=$JAVA_HOME/bin:$MAVEN_HOME/bin:$PATH

# homebrew
# ref: https://mirrors.ustc.edu.cn/help/brew.git.html#homebrew-linuxbrew
export HOMEBREW_BREW_GIT_REMOTE='https://mirrors.ustc.edu.cn/brew.git'
export HOMEBREW_CORE_GIT_REMOTE='https://mirrors.ustc.edu.cn/homebrew-core.git'
export HOMEBREW_BOTTLE_DOMAIN='https://mirrors.ustc.edu.cn/homebrew-bottles'
export HOMEBREW_API_DOMAIN='https://mirrors.ustc.edu.cn/homebrew-bottles/api'

# alias
alias ii="explorer.exe"
```

### 安装包管理器 Homebrew

```shell
# echo "export HOMEBREW_BREW_GIT_REMOTE='https://mirrors.ustc.edu.cn/brew.git'" >> ~/.shellrc
# echo "export HOMEBREW_CORE_GIT_REMOTE='https://mirrors.ustc.edu.cn/homebrew-core.git'" >> ~/.shellrc
# echo "export HOMEBREW_BOTTLE_DOMAIN='https://mirrors.ustc.edu.cn/homebrew-bottles'" >> ~/.shellrc
# echo "export HOMEBREW_API_DOMAIN='https://mirrors.ustc.edu.cn/homebrew-bottles/api'" >> ~/.shellrc
# source ~/.shellrc
/bin/bash -c "$(curl -fsSL https://mirrors.ustc.edu.cn/misc/brew-install.sh)"
```

### 安装软件

#### neovim & lazyvim

```shell
brew install fzf neovim
git clone https://github.com/LazyVim/starter ~/.config/nvim
```

[配置字体](https://www.nerdfonts.com/)(Hack)并在 Windows Terminal 中选定

#### Git配置

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
