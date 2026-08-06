---
title: "nixos-flake-install"
date: 2026-08-06T11:06:19+0800
lastmod: 2026-08-06T11:55:19+0800
categories: ['Linux']
tags: ['']
keywords: NixOS
description: 使用 Flake 管理 NixOS
image: 
---

> 之前的基础安装已经启用了 Flake 特性, 所以现在直接切换到 Flake 管理

新增并编辑 `flake.nix`

```bash
sudo vim /etc/nixos/flake.nix
```

写入如下内容

- 使用 `unstable` 版本是因为有 `flake.lock` 文件来锁定依赖
- 引入 `disko` 后可以去掉 `configuration.nix` 中的 `disko`

```nix
{
  description = "My NixOS Flake";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    disko.url = "github:nix-community/disko";
  };
  outputs = inputs@{ self, nixpkgs, disko, ... }: {
    # Define a system called "nixos"
    nixosConfigurations.nixos = nixpkgs.lib.nixosSystem {
      system = "x86_64-linux";
      specialArgs = { inherit inputs; };
      modules = [ 
        ./configuration.nix
        disko.nixosModules.disko
        ./disko.nix
      ];
    };
    # You can define many systems in one Flake file.
    # NixOS will choose one based on your hostname.
    #
    # nixosConfigurations."nixos2" = nixpkgs.lib.nixosSystem {
    #   system = "x86_64-linux";
    #   modules = [
    #     ./configuration2.nix
    #   ];
    # };
  };
}
```

再去掉 `configuration.nix` 中的 `disko`

```nix
imports = [ ./hardware-configuration.nix ];
```

最后执行如下命令来生成新世代以便应用配置


更新 lock 文件

```bash
cd /etc/nixos
```

```bash
sudo NIX_CONFIG="access-tokens = github.com=github_pat_xxx" \
HTTP_PROXY="http://192.168.137.1:1080" HTTPS_PROXY="http://192.168.137.1:1080" \
nix flake update
```

生成新世代

```bash
sudo NIX_CONFIG="access-tokens = github.com=github_pat_xxx" \
HTTP_PROXY="http://192.168.137.1:1080" HTTPS_PROXY="http://192.168.137.1:1080" \
nixos-rebuild switch --flake /etc/nixos#nixos
```

关机再开机/重启, 选择新的世代进入系统

```bash
sudo shutdown now
```

```bash
sudo reboot
```

此时相比之前的三个文件又多了两个文件

- `flake.nix`: Flake 的入口
- `flake.lock`: 版本锁文件, 确保系统可复现
- `disko.nix`: 分区信息
- `configuration.nix`: 基础配置
- `hardware-configuration.nix`: 硬件信息

