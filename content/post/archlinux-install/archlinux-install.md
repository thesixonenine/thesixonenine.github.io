---
title: archlinux-install
date: 2020-10-08T13:37:34+0800
updated: 2020-10-09T14:00:00+0800
tags: ['Archlinux']
categories: ['Linux']
keywords: archlinux
description: archlinux 安装
---

## 安装

参考[官方Wiki](https://wiki.archlinux.org/index.php/Installation_guide)

### 设置字体

```bash
setfont /usr/share/kbd/consolefonts/LatGrkCyr-12×22.psfu.gz
```

### 连接网络

**静态IP**

```bash
cd /etc/systemd/network
touch eth0.network
echo '[Match]' >> eth0.network
echo 'Name=eth0' >> eth0.network
echo '[Network]' >> eth0.network
echo 'Address=192.168.137.12' >> eth0.network
echo 'Gateway=192.168.137.1' >> eth0.network
echo 'DNS=223.5.5.5' >> eth0.network
systemctl reenable systemd-networkd
```

### 更新系统时间

```bash
timedatectl set-ntp true
```

### 分区

```bash
# 查看分区
fdisk -l
```
**BIOS 与 MBR**

| 挂载点   | 分区        | 分区类型             | 建议大小 |
| -------- | ----------- | -------------------- | -------- |
| `/mnt`   | `/dev/sdX1` | Linux                | 剩余空间 |
| `[SWAP]` | `dev/sdX2`  | Linux swap(交换空间) | 大于512M |

**UEFI 与 GPT**

| 挂载点                  | 分区        | 分区类型               | 建议大小 |
| ----------------------- | ----------- | ---------------------- | -------- |
| `/mnt/boot`或`/mnt/efi` | `/dev/sdX1` | EFI系统分区            | 260-512M |
| `/mnt`                  | `/dev/sdX2` | Linux x86-64 根目录(/) | 剩余空间 |
| `[SWAP]`                | `/dev/sdX3` | Linux swap(交换空间)   | 大于512M |

**建立磁盘分区**

```bash
# 进入磁盘
fdisk /dev/sda
```

1. `g`创建一个空的gpt分区
2. `n`创建新分区(编号1), 大小为512M, 用作系统引导
3. `n`创建新分区(编号3), 与内存一样(4G), 用作SWAP
4. `n`创建新分区(编号2), 使用剩下所有的空间, 用作主分区(/)
5. `w`写入并退出

**格式化分区**

1. `mkfs.fat -F32 /dev/sda1` 格式化编号为1的引导分区(EFI系统分区)
2. `mkfs.ext4 /dev/sda2` 格式化编号为2的根分区(/)
3. `mkswap /dev/sda3` 格式化编号为3的swap分区
4. `swapon /dev/sda3` 打开swap

### 设置pacman的镜像源

```bash
vim /etc/pacman.d/mirrorlist
```

在第一行新增`Server = https://mirrors.tuna.tsinghua.edu.cn/archlinux/$repo/os/$arch`

### 挂载镜像

1. `mount /dev/sda2 /mnt` 将根目录挂载到`/mnt`
2. `mkdir /mnt/boot`创建EFI系统分区需要挂载的目录
3. `mount /dev/sda1 /mnt/boot` 将EFI系统分区挂载到`/mnt/boot`

### 安装必须的软件包

```bash
pacstrap /mnt base linux linux-firmware
```

### 配置系统

**生成fstab文件**

```bash
genfstab -U /mnt >> /mnt/etc/fstab
```

**切换到新安装的系统**

```bash
arch-chroot /mnt
```

**设置时区并同步系统时间**

```bash
ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
hwclock --systohc
```

**修改时区**

```bash
# 退出新安装的系统, 回到ISO中, 因为新安装的系统中连vim都没
exit
vim /etc/locale.gen
# 取消 en_US.UTF-8 UTF-8 这一行的注释
# 切换到新安装的系统
arch-chroot /mnt
# 生成 locale 信息
locale-gen
# 再次退出
exit
# 设置本地语言配置文件
echo 'LANG=en_US.UTF-8' >> /mnt/etc/locale.conf
```

**网络配置**

```bash
# 配置计算机名称
echo '12' >> /mnt/etc/hostname
echo '127.0.0.1	localhost' >> /mnt/etc/hosts
echo '::1       localhost' >> /mnt/etc/hosts
# 固定IP可以写入 echo '192.168.137.12 12' >> /mnt/etc/hosts
```

**修改root密码**

```bash
# 切换到新安装的系统
arch-chroot /mnt
passwd
```

**安装引导程序**

```bash
pacman -S grub efibootmgr intel-ucode os-prober
```

**配置GRUB**

```bash
mkdir /boot/grub
grub-mkconfig > /boot/grub/grub.cfg
grub-install --target=x86_64-efi --efi-directory=/boot --bootloader-id=GRUB
```

**安装程序**

```bash
pacman -S vim git zsh
```

**重启系统**

```bash
# 退出新安装的系统, 回到ISO中
exit
# 关机
shutdown -h now
```



## 使用

网络配置与远程登陆

### 设置网络

```bash
# 启用网络
ip link set eth0 up
# 设置静态IP
cd /etc/systemd/network
touch eth0.network
echo '[Match]' >> eth0.network
echo 'Name=eth0' >> eth0.network
echo '[Network]' >> eth0.network
echo 'Address=192.168.137.12' >> eth0.network
echo 'Gateway=192.168.137.1' >> eth0.network
echo 'DNS=223.5.5.5' >> eth0.network
systemctl restart systemd-resolved
systemctl restart systemd-networkd
systemctl enable systemd-resolved
systemctl enable systemd-networkd
```

### 安装SSH

```bash
pacman -S openssh
systemctl start sshd
systemctl enable sshd
```

> 默认是禁止root用户远程登录的
>
> 使用`echo 'PermitRootLogin yes' >> /etc/ssh/sshd_config`允许root用户远程登陆



### 安装oh-my-zsh

| Method    | Command                                                      |
| --------- | ------------------------------------------------------------ |
| **curl**  | `sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"` |
| **wget**  | `sh -c "$(wget -O- https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"` |
| **fetch** | `sh -c "$(fetch -o - https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"` |

修改主题

```bash
vim ~/.zshrc
ZSH_THEME="agnoster"
```





## 文件权限

### 类型区分

- `-` 普通文件
- `d` 目录文件
- `l` 链接文件
- `b` 块设备文件
- `c` 字符设备文件
- `p` 管道文件

### 操作理解

**文件**

| 权限      | 含义                         |
| --------- | ---------------------------- |
| 读(r-4)   | 读取文件内容                 |
| 写(w-2)   | 修改文件内容(新增/修改/删除) |
| 执行(x-1) | 执行脚本文件                 |

**目录**

| 权限      | 含义                           |
| --------- | ------------------------------ |
| 读(r-4)   | 读取目录内的文件列表           |
| 写(w-2)   | 可在目录内新增/删除/重命名文件 |
| 执行(x-1) | 进入目录                       |

