#!/bin/bash
set -e

# Update the system clock
timedatectl set-ntp true
timedatectl set-timezone Asia/Shanghai
timedatectl

# Update the mirrors
[ "$(head -n1 /etc/pacman.d/mirrorlist)" = "## China" ] || sed -i '1i\
## China\
Server = https://mirrors.aliyun.com/archlinux/$repo/os/$arch\
Server = https://mirrors.tuna.tsinghua.edu.cn/archlinux/$repo/os/$arch\
Server = https://mirrors.ustc.edu.cn/archlinux/$repo/os/$arch\
' /etc/pacman.d/mirrorlist

head /etc/pacman.d/mirrorlist

# Partition the disks
DISK="/dev/sda"
EFI_SIZE="512M"
SWAP_SIZE="4G"

if ls "${DISK}"[0-9]* 1>/dev/null 2>&1; then
  echo "$DISK has Parted"
  fdisk -l
else

echo "Partition the disk $DISK ..."

fdisk "$DISK" <<EOF
g
n
1

+${EFI_SIZE}

t
1
n
2

+${SWAP_SIZE}

t
2
19
n
3


t
3
20
w
EOF

partprobe "$DISK" || true
sleep 2

echo "Format the partitions"

mkfs.fat -F32 "${DISK}1"
mkfs.ext4 -F "${DISK}3"
mkswap "${DISK}2"
swapon "${DISK}2"

lsblk "$DISK"
fi

echo 'Mount the file systems'
mount /dev/sda3 /mnt
mount --mkdir /dev/sda1 /mnt/boot
echo 'Install essential packages'
pacstrap /mnt base linux linux-firmware
echo 'Configure the system...'
echo 'Fstab'
genfstab -U /mnt >> /mnt/etc/fstab
echo 'Chroot'
arch-chroot /mnt
