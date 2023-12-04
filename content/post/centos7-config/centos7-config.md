---
title: centos7-config
date: 2022-01-16T15:24:10+0800
lastmod: 2022-02-09T14:20:10+0800
tags: 
- CentOS
- Linux
categories: 
- Linux
keywords:
description: CentOS7常用配置
---

## CA证书过期导致yum更新失败

由于`Let's Encrypt's`的CA证书过期, 而`yum`使用的`curl`依赖CA证书, 解决办法就是更新证书

```bash
# 如果执行过yum clean all, 则需要先关闭ssl验证, 等更新CA证书后再开
# sudo cat "sslverify=0" >> /etc/yum.conf
sudo yum install -y ca-certificates
sudo update-ca-trust extract
```

## VIM常用配置

**~/.vimrc**

```bash
:set fileencodings=utf-8
:set encoding=utf-8
:set nowrap
:set nu
:set softtabstop=4
:set shiftwidth=4
:set tabstop=4

" 设置主题配色
" colorscheme solarized
" colorscheme koehler
colorscheme desert
```

## PowerShell常用配置

**Microsoft.PowerShell_profile.ps1**

```powershell
Import-Module posh-git
Import-Module oh-my-posh
# Set-PoshPrompt -Theme Paradox
Set-PoshPrompt -Theme nu4a

# 这个主题会改tab名称
# Set-PoshPrompt -Theme pure
# Set-PoshPrompt -Theme ys

# 设置预测文本来源为历史记录
Set-PSReadLineOption -PredictionSource History
# 设置 Tab 键补全
Set-PSReadLineKeyHandler -Key "Tab" -Function MenuComplete

cls
if (Test-Path ~\Desktop\VPN.lnk) {
    Remove-Item -Recurse ~\Desktop\VPN.lnk
    Write-Host "VPN.lnk 删除成功"
}
# 传递指定公钥到服务器上
function ssh-copy-id([string]$userAtMachine, $args){   
    $publicKey = "$ENV:USERPROFILE" + "/.ssh/id_rsa.pub"
    if (!(Test-Path "$publicKey")){
        Write-Error "ERROR: failed to open ID file '$publicKey': No such file"            
    }
    else {
        & cat "$publicKey" | ssh $args $userAtMachine "umask 077; test -d .ssh || mkdir .ssh ; cat >> .ssh/authorized_keys || exit 1"      
    }
}
# 代理相关


################
### HTTP 代理 ###
################
function proxy_set([string]$protocol, [int]$port, $args){
    if ([String]::IsNullOrEmpty($protocol) -or 'http','socks5' -cnotcontains $protocol) {
        # 区分大小写, 且左边不包含右边
        "Invalid protocol[${protocol}], The default value of socks5 will be used"
        $protocol = "socks5"
    }
    if ($port -le 0 -or $port -gt 65535) {
        "Invalid port[${port}], The default value of 10808 will be used"
        $port = 10808
    }
    Set-Item Env:http_proxy "${protocol}://127.0.0.1:${port}"
    Set-Item Env:https_proxy "${protocol}://127.0.0.1:${port}"
}
function proxy_unset {
    Remove-Item Env:http_proxy
    Remove-Item Env:https_proxy
}
function proxy_get {
    "http_proxy  = ${env:http_proxy}"
    "https_proxy = ${env:https_proxy}"
}
function proxy_test {
    curl -v http://www.google.com
}
function proxy_commad {
    "设置代理: proxy_set"
    "重置代理: proxy_unset"
    "查看代理: proxy_get"
    "测试代理: proxy_test"
}

###############
### Git 代理 ###
###############
function git_proxy_set([string]$protocol, [int]$port, $args){
    if ([String]::IsNullOrEmpty($protocol) -or 'http','socks5' -cnotcontains $protocol) {
        # 区分大小写, 且左边不包含右边
        "Invalid protocol[${protocol}], The default value of socks5 will be used"
        $protocol = "socks5"
    }
    if ($port -le 0 -or $port -gt 65535) {
        "Invalid port[${port}], The default value of 10808 will be used"
        $port = 10808
    }
    git config http.proxy "${protocol}://127.0.0.1:${port}"
    git config https.proxy "${protocol}://127.0.0.1:${port}"
}
function git_proxy_unset {
    git config --unset http.proxy
    git config --unset https.proxy
}
function git_proxy_get {
    git config http.proxy
    git config https.proxy
}
function git_proxy_commad {
    "设置代理: git_proxy_set"
    "重置代理: git_proxy_unset"
    "查看代理: git_proxy_get"
}
```

## 修改yum源

### 阿里云

```bash
sudo mv /etc/yum.repos.d/CentOS-Base.repo /etc/yum.repos.d/CentOS-Base.repo.backup
sudo curl -o /etc/yum.repos.d/CentOS-Base.repo https://mirrors.aliyun.com/repo/Centos-7.repo
sudo sed -i -e '/mirrors.cloud.aliyuncs.com/d' -e '/mirrors.aliyuncs.com/d' /etc/yum.repos.d/CentOS-Base.repo
sudo yum makecache
```

### ustc

```bash
sudo sed -e 's|^mirrorlist=|#mirrorlist=|g' \
         -e 's|^#baseurl=http://mirror.centos.org/centos|baseurl=https://mirrors.ustc.edu.cn/centos|g' \
         -i.bak \
         /etc/yum.repos.d/CentOS-Base.repo
sudo yum makecache
```

### tsinghua

```bash
sudo sed -e 's|^mirrorlist=|#mirrorlist=|g' \
         -e 's|^#baseurl=http://mirror.centos.org|baseurl=https://mirrors.tuna.tsinghua.edu.cn|g' \
         -i.bak \
         /etc/yum.repos.d/CentOS-*.repo
sudo yum makecache
```



