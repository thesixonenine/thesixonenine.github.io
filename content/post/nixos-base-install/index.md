---
title: "nixos-base-install"
date: 2026-07-10T17:34:40+08:00
lastmod: 2026-07-10T17:34:40+08:00
categories: ['Linux']
tags: ['']
keywords: NixOS
description: NixOS 基础安装
image: 
---

### 从 ISO 文件安装

1. 从镜像源下载 ISO 文件
2. 烧录进U盘并插入电脑选择U盘启动/虚拟机挂载ISO并启动

可选[镜像源](https://mirrorz.org/os/NixOS): https://mirrorz.org/os/NixOS

### 从局域网连接 liveCD 环境

使用 `ip a` 查看 IP, 使用 `passwd` 设置密码, 然后就可以在局域网内通过 SSH 进行远程连接开始安装

### 切换到 root 用户

```bash
sudo -i
```

### 使用 disko 进行快速分区

```bash
vim disko.nix
```

`disko.nix`

```nix
{
  disko.devices = {
    disk.main = {
      type = "disk";
      device = "/dev/sda";
      content = {
        type = "gpt";
        partitions = {
          ESP = {
            priority = 1;
            size = "512M";
            type = "EF00";
            content = {
              type = "filesystem";
              format = "vfat";
              mountpoint = "/boot";
              mountOptions = [ "umask=0077" ];
            };
          };
          luks = {
            size = "100%";
            content = {
              type = "luks";
              name = "cryptroot";
              content = {
                type = "btrfs";
                extraArgs = [ "-f" ];
                subvolumes = {
                  "@root" = {
                    mountpoint = "/";
                    mountOptions = [ "compress=zstd:3" "noatime" "space_cache=v2" ];
                  };
                  "@home" = {
                    mountpoint = "/home";
                    mountOptions = [ "compress=zstd:3" "noatime" "space_cache=v2" ];
                  };
                  "@nix" = {
                    mountpoint = "/nix";
                    mountOptions = [ "compress-force=zstd:3" "noatime" ];
                  };
                  "@snapshots" = {
                    mountpoint = "/.snapshots";
                    mountOptions = [ "compress=zstd:3" "noatime" ];
                  };
                };
              };
            };
          };
        };
      };
    };
  };
}
```

执行分区

```bash
sudo NIX_CONFIG="access-tokens = github.com=github_pat_xxx" \
HTTP_PROXY="http://192.168.137.1:1080" HTTPS_PROXY="http://192.168.137.1:1080" \
nix --extra-experimental-features "nix-command flakes" run github:nix-community/disko/latest -- --mode disko ./disko.nix
```

### 生成初始化配置

```bash
nixos-generate-config --root /mnt
```

### 简单配置系统

```bash
vim /mnt/etc/nixos/configuration.nix
```

追加到 `configuration.nix`

```nix
networking.hostName = "nixos";
networking.proxy.default = "http://192.168.137.1:1080/";
networking.proxy.noProxy = "127.0.0.1,localhost,internal.domain";
time.timeZone = "Asia/Shanghai";
i18n.defaultLocale = "en_US.UTF-8";
i18n.extraLocaleSettings = {
  LC_ADDRESS = "zh_CN.UTF-8";
  LC_IDENTIFICATION = "zh_CN.UTF-8";
  LC_MEASUREMENT = "zh_CN.UTF-8";
  LC_MONETARY = "zh_CN.UTF-8";
  LC_NAME = "zh_CN.UTF-8";
  LC_NUMERIC = "zh_CN.UTF-8";
  LC_PAPER = "zh_CN.UTF-8";
  LC_TELEPHONE = "zh_CN.UTF-8";
  LC_TIME = "zh_CN.UTF-8";
};
users.users."simple" = {
  isNormalUser = true;
  shell = pkgs.zsh;
  description = "Simple";
  initialPassword = "1";
  extraGroups = [ "wheel" "networkmanager" ];
  openssh.authorizedKeys.keys = [ "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIEfY4AqFEB76gUXJKVifON936yf/MdsOKTsmioQ3HDKi" ];
};
services.openssh = {
  enable = true;
  ports = [ 22 ];
  settings = {
    PasswordAuthentication = true;
    PermitRootLogin = "no";
    AllowUsers = [ "simple" ];
  };
};
environment.systemPackages = with pkgs; [ vim git curl ];
programs.zsh.enable = true;
nix.settings.experimental-features = [ "nix-command" "flakes" ];
```

微调硬件配置: `在 subvol 中按disko.nix追加参数, 例如 "compress=zstd:3"`

```bash
vim /mnt/etc/nixos/hardware-configuration.nix
```

### 安装系统

```bash
sudo NIX_CONFIG="access-tokens = github.com=github_pat_xxx" \
HTTP_PROXY="http://192.168.137.1:1080" HTTPS_PROXY="http://192.168.137.1:1080" \
nixos-install
```

安装过程会要求输入`全盘加密密码`

### 关机

```bash
shutdown now
```

### 重新开机并进入系统

开机过程会要求输入`全盘加密密码`然后进入命令行, 接着输入用户名和密码进行登录并设置zsh
