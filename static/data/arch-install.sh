#!/bin/bash
timedatectl set-ntp true
timedatectl set-timezone Asia/Shanghai
timedatectl

awk '
/^## China/ { in_china = 1; print; next }
in_china && /^##/ { in_china = 0 }
in_china { china = china $0 "\n"; next }
{ other = other $0 "\n" }
END { printf "%s%s", china, other }
' /etc/pacman.d/mirrorlist
