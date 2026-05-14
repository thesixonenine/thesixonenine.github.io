#!/bin/bash
set -e

echo 'Time'
ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
hwclock --systohc
echo 'Localization'
sed -i 's/^#en_US\.UTF-8/en_US.UTF-8/' /etc/locale.gen
locale-gen
echo 'LANG=en_US.UTF-8' > /etc/locale.conf
echo 'Arch' > /etc/hostname
echo '127.0.0.1	localhost' >> /etc/hosts
echo '::1       localhost' >> /etc/hosts


echo 'Root password'
echo "root:123456" | chpasswd

echo 'Boot loader'
pacman --noconfirm -S grub efibootmgr intel-ucode os-prober
mkdir /boot/grub
grub-mkconfig > /boot/grub/grub.cfg
grub-install --target=x86_64-efi --efi-directory=/boot --bootloader-id=GRUB

pacman --noconfirm -S vim git zsh networkmanager openssh
systemctl enable NetworkManager
systemctl enable sshd

mkdir -p ~/.ssh
echo 'ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIEfY4AqFEB76gUXJKVifON936yf/MdsOKTsmioQ3HDKi' >> ~/.ssh/authorized_keys
chmod 700 ~/.ssh
chmod 600 ~/.ssh/authorized_keys
echo 'PermitRootLogin yes' >> /etc/ssh/sshd_config
