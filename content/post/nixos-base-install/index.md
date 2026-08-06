---
title: "nixos-base-install"
date: 2026-07-10T17:34:40+08:00
lastmod: 2026-08-06T10:32:19+0800
categories: ['Linux']
tags: ['']
keywords: NixOS
description: 基于 disko 的 NixOS 基础安装(无 Flake)
image: 
---

参考 [disko](https://github.com/nix-community/disko/blob/master/docs/quickstart.md) 和 [NixOS Manual](https://nixos.org/manual/nixos/stable/#sec-installation-manual)

## 从 ISO 文件安装

1. 从镜像源下载 ISO 文件
2. 烧录进U盘并插入电脑选择U盘启动/虚拟机挂载ISO并启动(都需要关闭安全启动)

可选[镜像源](https://mirrorz.org/os/NixOS): https://mirrorz.org/os/NixOS

## 从局域网连接 liveCD 环境

使用 `ip a` 查看 IP, 使用 `passwd` 设置密码, 然后就可以在局域网内通过 `ssh nixos@192.168.137.124` 进行远程连接开始安装

## 切换到 root 用户

```bash
sudo -i
```

## 使用 disko 进行快速分区

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

`disko-without-luks.nix`

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
          root = {
            size = "100%";
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
}
```

执行分区

```bash
sudo NIX_CONFIG="access-tokens = github.com=github_pat_xxx" \
HTTP_PROXY="http://192.168.137.1:1080" HTTPS_PROXY="http://192.168.137.1:1080" \
nix --extra-experimental-features "nix-command flakes" run github:nix-community/disko/latest -- --mode disko ./disko.nix
```

## 生成初始化配置

```bash
nixos-generate-config --no-filesystems --root /mnt
```

## 简单配置系统

只需要编辑 `/mnt/etc/nixos/configuration.nix`, 其默认配置如下

```nix
{ config, lib, pkgs, ... }: {
  imports = [ ./hardware-configuration.nix ];
  boot.loader.systemd-boot.enable = true;
  boot.loader.efi.canTouchEfiVariables = true;
  boot.kernelPackages = pkgs.linuxPackages_latest;
  networking.networkmanager.enable = true;
  system.stateVersion = "26.05";
}
```

配置流程如下

- 将 `disko.nix` 文件移动到 `/mnt/etc/nixos/` 下
- 在 `/mnt/etc/nixos/configuration.nix` 中引入 `disko` NixOS module 和 `disko.nix` (后续使用 Flake 管理后将恢复为只引入硬件配置)
- 在 `/mnt/etc/nixos/configuration.nix` 中追加其他必要配置


**移动 `disko.nix`**

```bash
cp disko.nix /mnt/etc/nixos/
```

开始编辑 `/mnt/etc/nixos/configuration.nix`

```bash
vim /mnt/etc/nixos/configuration.nix
```

**引入 disko 配置**, 后续使用 Flake 管理后将恢复为只引入 `hardware-configuration.nix`

```nix
imports =
 [
   ./hardware-configuration.nix
   "${builtins.fetchTarball "https://github.com/nix-community/disko/archive/master.tar.gz"}/module.nix"
   ./disko.nix
 ];
```

**追加其他必要配置到 `configuration.nix`**

- 网络代理
- 时区与语言
- 新建用户及其密码, 密钥, 默认shell等等
- 必要软件如vim, git, curl等
- 设置 GitHub PAT(公共仓库只读即可), 避免被 GitHub 限制以及后续命令的简洁
- 启用 Flake 特性(后续用 Flake 进行管理)

所有配置参考[官方手册](https://nixos.org/manual/nixos/stable/options)

```nix
networking.hostName = "nixos";
networking.proxy.default = "http://192.168.137.1:1080";

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
  description = "Simple";
  initialPassword = "1";
  extraGroups = [ "wheel" "networkmanager" ];
  openssh.authorizedKeys.keys = [ "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIEfY4AqFEB76gUXJKVifON936yf/MdsOKTsmioQ3HDKi" ];
};
services.openssh = {
  enable = true;
  settings = {
    PermitRootLogin = "no";
    AllowUsers = [ "simple" ];
  };
};
environment.systemPackages = with pkgs; [ vim git curl ];
environment.variables.EDITOR = "vim";
nix.settings.experimental-features = [ "nix-command" "flakes" ];
# 注意修改为自己的 GitHub PAT
nix.settings.access-tokens = [ "github.com=github_pat_xxx" ];
```

## 安装系统

> 仅在执行安装命令时需要手动指定 GitHub PAT 和 代理, 后续命令不再需要指定这两个参数, 因为在上一步简单配置中已经指定了 `nix.settings.access-tokens` 和 `networking.proxy.default`

```bash
sudo NIX_CONFIG="access-tokens = github.com=github_pat_xxx" \
HTTP_PROXY="http://192.168.137.1:1080" HTTPS_PROXY="http://192.168.137.1:1080" \
nixos-install
```

`for flake`

```bash
sudo NIX_CONFIG="access-tokens = github.com=github_pat_xxx" \
HTTP_PROXY="http://192.168.137.1:1080" HTTPS_PROXY="http://192.168.137.1:1080" \
nixos-install --flake '/mnt/etc/nixos#nixos'
```

如果使用了 `luks` 则安装过程会要求设置 `磁盘加密密码`

## 重启并进入系统

关机

```bash
sudo shutdown now
```

电脑拔掉U盘/虚拟机卸载ISO

开机并设置从 `systemd-bootx64.efi` 启动 (如果使用了 `luks` 则会要求输入 `磁盘加密密码` )进入命令行, 接着输入用户名和密码进行登录并设置 `zsh`

简单安装就到这里, 后续会继续优化, 包括: 切换到 `Flake`, 使用 `home-manager`, 安装 `hyprland` / `niri`, `noctalia` 等等, 以此构建可正常使用的 NixOS

## 静态IP配置(可选)

- 编辑静态 IP 的配置
- 如果是安装完成进入系统后编辑则需要生成新世代
- 在安装时编辑后第一次进入系统/生成新世代后进入系统
- 停止已有连接,启动配置连接,删除已停止连接

编辑 `configuration.nix`, 可以是在安装时编辑(`/mnt/etc/nixos/configuration.nix`), 也可以在安装完成进入系统后编辑(`/etc/nixos/configuration.nix`)

```bash
sudo vim /etc/nixos/configuration.nix
```

追加如下内容

```nix
  networking.networkmanager.ensureProfiles.profiles = {
    "eth0" = {
      connection = {
        id = "eth0";
        uuid = "9f6f3a52-1b88-4b0d-a2d0-8b7e3b4c9a01";
        type = "ethernet";
        interface-name = "eth0";
        autoconnect = true;
      };
      ipv4 = {
        method = "manual";
        address1 = "192.168.137.10/24,192.168.137.1";
        dns = "223.5.5.5;223.6.6.6;";
      };
      ipv6.method = "disabled";
      ethernet = {};
    };
  };
```

如果是安装时编辑则继续安装并重启进入系统, 否则还需要执行如下命令生成新世代并重启进入系统

```bash
sudo nixos-rebuild switch
```

> 之前已经配置了 `users.users."simple".extraGroups = [ "wheel" "networkmanager" ];`, 所以下面的 `nmcli` 命令均不需要 `sudo`

查看现有连接

```bash
nmcli -f NAME,UUID,FILENAME connection show
```

停止已有连接

```bash
nmcli connection down "Wired connection 1"
```

启动配置连接

```bash
nmcli connection up eth0
```

删除已停止连接

```bash
nmcli connection delete "Wired connection 1"
```

**后续日常使用目标**

- 开机并输入加密密码进入 `tty` 环境
- 用账户密码登录
- 执行 `start-hyprland` 启动 hyprland 和 noctalia

<details>
<summary>附: 安装好后系统的配置</summary>

查看配置文件

```bash
ls -hl /etc/nixos
```

可以看到三个文件 

- `disko.nix`: 与之前安装步骤中的内容相同
- `configuration.nix`: 内容即是之前安装步骤修改完成后的
- `hardware-configuration.nix`: 根据载入 ISO 时选项不同与安装环境(实体机/虚拟机)有略微不同

`configuration.nix`

```nix
{ config, lib, pkgs, ... }: {
  imports =
  [
    ./hardware-configuration.nix
    "${builtins.fetchTarball "https://github.com/nix-community/disko/archive/master.tar.gz"}/module.nix"
    ./disko.nix
  ];
  boot.loader.systemd-boot.enable = true;
  boot.loader.efi.canTouchEfiVariables = true;
  boot.kernelPackages = pkgs.linuxPackages_latest;
  networking.networkmanager.enable = true;
  system.stateVersion = "26.05";
  networking.hostName = "nixos";
  networking.proxy.default = "http://192.168.137.1:1080";
  networking.networkmanager.ensureProfiles.profiles = {
    "eth0" = {
      connection = {
        id = "eth0";
        uuid = "9f6f3a52-1b88-4b0d-a2d0-8b7e3b4c9a01";
        type = "ethernet";
        interface-name = "eth0";
        autoconnect = true;
      };
      ipv4 = {
        method = "manual";
        address1 = "192.168.137.10/24,192.168.137.1";
        dns = "223.5.5.5;223.6.6.6;";
      };
      ipv6.method = "disabled";
      ethernet = {};
    };
  };
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
    settings = {
      PermitRootLogin = "no";
      AllowUsers = [ "simple" ];
    };
  };
  environment.systemPackages = with pkgs; [ vim git curl ];
  environment.variables.EDITOR = "vim";
  programs.zsh.enable = true;
  nix.settings.experimental-features = [ "nix-command" "flakes" ];
}
```

`hardware-configuration.nix`

```nix
# Do not modify this file!  It was generated by ‘nixos-generate-config’
# and may be overwritten by future invocations.  Please make changes
# to /etc/nixos/configuration.nix instead.
{ config, lib, pkgs, modulesPath, ... }:

{
  imports = [ ];

  boot.initrd.availableKernelModules = [ "sd_mod" "sr_mod" ];
  boot.initrd.kernelModules = [ ];
  boot.kernelModules = [ ];
  boot.extraModulePackages = [ ];

  nixpkgs.hostPlatform = lib.mkDefault "x86_64-linux";
  virtualisation.hypervGuest.enable = true;
}
```

</details>
