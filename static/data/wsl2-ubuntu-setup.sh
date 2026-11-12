#!/usr/bin/env bash

set -Eeuo pipefail

# ==========================================
# WSL2 Ubuntu 初始化脚本
# ==========================================

trap 'echo "[ERROR] 执行失败，行号: $LINENO" >&2' ERR

# ==========================================
# 1. 配置 /etc/wsl.conf
# ==========================================

WSL_CONF="/etc/wsl.conf"

echo "[INFO] 配置 WSL interop..."

sudo touch "$WSL_CONF"

if sudo grep -qE '^[[:space:]]*\[interop\][[:space:]]*$' "$WSL_CONF"; then
    if sudo sed -n '/^[[:space:]]*\[interop\][[:space:]]*$/,/^[[:space:]]*\[/p' "$WSL_CONF" \
        | grep -qE '^[[:space:]]*appendWindowsPath[[:space:]]*='; then

        sudo sed -i '/^[[:space:]]*\[interop\][[:space:]]*$/,/^[[:space:]]*\[/ s/^[[:space:]]*appendWindowsPath[[:space:]]*=.*/appendWindowsPath=false/' "$WSL_CONF"
    else
        sudo sed -i '/^[[:space:]]*\[interop\][[:space:]]*$/a appendWindowsPath=false' "$WSL_CONF"
    fi
else
    sudo tee -a "$WSL_CONF" > /dev/null <<'EOF'

[interop]
appendWindowsPath=false
EOF
fi

echo "[INFO] /etc/wsl.conf 配置完成。"

# ==========================================
# 2. 配置 Ubuntu APT 镜像源
# ==========================================

echo "[INFO] 检查 Ubuntu 版本..."

lsb_release -a

APT_SOURCE="/etc/apt/sources.list.d/ubuntu.sources"

if [[ -f "$APT_SOURCE" ]]; then
    echo "[INFO] 替换 Ubuntu APT 镜像源为 USTC..."

    sudo sed -i 's@//.*archive.ubuntu.com@//mirrors.ustc.edu.cn@g' "$APT_SOURCE"
    sudo sed -i 's/security.ubuntu.com/mirrors.ustc.edu.cn/g' "$APT_SOURCE"
    sudo sed -i 's/http:/https:/g' "$APT_SOURCE"

    echo "[INFO] 更新 APT 软件包索引..."

    sudo apt update
else
    echo "[ERROR] 未找到 APT 源文件: $APT_SOURCE" >&2
    exit 1
fi

echo "[INFO] Ubuntu APT 镜像源配置完成。"

# ==========================================
# 3. 配置 Bash 代理管理功能
# ==========================================

BASHRC="$HOME/.bashrc"
PROXY_MARKER="# >>> WSL proxy configuration >>>"

echo "[INFO] 配置 Bash 代理管理功能..."

touch "$BASHRC"

if ! grep -qF "$PROXY_MARKER" "$BASHRC"; then
    cat >> "$BASHRC" <<'EOF'

# >>> WSL proxy configuration >>>

function proxy_on {
    local host
    local port="${1:-1080}"
    if ! [[ "$port" =~ ^[0-9]+$ ]] || (( port < 1 || port > 65535 )); then
        echo "端口号必须是 1-65535 之间的整数"
        return 1
    fi
    host=$(ip route show default | awk '/default/ {print $3; exit}')
    if [ -z "$host" ]; then
        echo "无法获取 Windows 主机 IP"
        return 1
    fi

    export HTTP_PROXY="http://${host}:${port}"
    export HTTPS_PROXY="http://${host}:${port}"
    export ALL_PROXY="http://${host}:${port}"
    export NO_PROXY="localhost,127.0.0.1"

    export http_proxy="$HTTP_PROXY"
    export https_proxy="$HTTPS_PROXY"
    export all_proxy="$ALL_PROXY"

    echo "代理已设置 ${host}:${port}"
}
# default proxy_on
proxy_on
function proxy_off {
    unset HTTP_PROXY HTTPS_PROXY ALL_PROXY NO_PROXY
    unset http_proxy https_proxy all_proxy
}

# <<< WSL proxy configuration <<<
EOF
    echo "[INFO] Bash 代理配置追加完成。"
else
    echo "[INFO] Bash 代理配置已经存在，跳过追加。"
fi

# ==========================================
# 4. 启动新的 Bash Shell
# ==========================================

echo "[INFO] 启动新的 Bash Shell，重新加载 ~/.bashrc..."

bash

echo "[INFO] 新 Bash Shell 已退出，继续执行初始化脚本..."

# ==========================================
# 5. 安装基础开发工具和配置 Zsh
# ==========================================

echo "[INFO] 安装基础开发工具..."

sudo apt install -y \
    build-essential \
    curl \
    git \
    sudo \
    wget \
    file \
    unzip \
    software-properties-common \
    zsh

# ------------------------------------------
# 设置默认 Shell 为 Zsh
# ------------------------------------------

echo "[INFO] 设置默认 Shell 为 Zsh..."

if [[ "$(getent passwd "$USER" | cut -d: -f7)" != "/bin/zsh" ]]; then
    chsh -s /bin/zsh
else
    echo "[INFO] 默认 Shell 已经是 Zsh，跳过修改。"
fi

# ------------------------------------------
# 安装 Oh My Zsh
# ------------------------------------------

OH_MY_ZSH_DIR="$HOME/.oh-my-zsh"

if [[ ! -d "$OH_MY_ZSH_DIR" ]]; then
    echo "[INFO] 从南京大学镜像站克隆 Oh My Zsh..."

    git clone --depth=1 \
        https://mirror.nju.edu.cn/git/ohmyzsh.git \
        "$OH_MY_ZSH_DIR"
else
    echo "[INFO] Oh My Zsh 已经存在，跳过克隆。"
fi

# ------------------------------------------
# 配置 .zshrc
# ------------------------------------------

if [[ ! -f "$HOME/.zshrc" ]]; then
    echo "[INFO] 初始化 .zshrc..."

    cp "$OH_MY_ZSH_DIR/templates/zshrc.zsh-template" "$HOME/.zshrc"
else
    echo "[INFO] .zshrc 已经存在，保留现有配置。"
fi

echo "[INFO] Zsh 和 Oh My Zsh 配置完成。"

# ==========================================
# 6. 安装 Zsh 插件
# ==========================================

ZSH_CUSTOM="${ZSH_CUSTOM:-$HOME/.oh-my-zsh/custom}"
ZSH_PLUGIN_DIR="$ZSH_CUSTOM/plugins"

echo "[INFO] 安装 Zsh 插件..."

mkdir -p "$ZSH_PLUGIN_DIR"

# ------------------------------------------
# 安装 zsh-autosuggestions
# ------------------------------------------

AUTOSUGGESTIONS_DIR="$ZSH_PLUGIN_DIR/zsh-autosuggestions"

if [[ ! -d "$AUTOSUGGESTIONS_DIR" ]]; then
    echo "[INFO] 安装 zsh-autosuggestions..."

    git clone --depth=1 \
        https://github.com/zsh-users/zsh-autosuggestions \
        "$AUTOSUGGESTIONS_DIR"
else
    echo "[INFO] zsh-autosuggestions 已经存在，跳过安装。"
fi

# ------------------------------------------
# 安装 zsh-syntax-highlighting
# ------------------------------------------

SYNTAX_HIGHLIGHTING_DIR="$ZSH_PLUGIN_DIR/zsh-syntax-highlighting"

if [[ ! -d "$SYNTAX_HIGHLIGHTING_DIR" ]]; then
    echo "[INFO] 安装 zsh-syntax-highlighting..."

    git clone --depth=1 \
        https://github.com/zsh-users/zsh-syntax-highlighting.git \
        "$SYNTAX_HIGHLIGHTING_DIR"
else
    echo "[INFO] zsh-syntax-highlighting 已经存在，跳过安装。"
fi

# ------------------------------------------
# 配置 .zshrc 插件
# ------------------------------------------

ZSHRC="$HOME/.zshrc"

if grep -qE '^plugins=\(git\)[[:space:]]*$' "$ZSHRC"; then
    sed -i '/^plugins=(git)[[:space:]]*$/c\
plugins=(\
    git\
    zsh-autosuggestions\
    zsh-syntax-highlighting\
)' "$ZSHRC"
else
    # 仅在插件尚未配置时进行追加
    if ! grep -qE '^[[:space:]]*zsh-autosuggestions[[:space:]]*$' "$ZSHRC"; then
        sed -i '/^plugins=(/a\    zsh-autosuggestions' "$ZSHRC"
    fi

    if ! grep -qE '^[[:space:]]*zsh-syntax-highlighting[[:space:]]*$' "$ZSHRC"; then
        sed -i '/^plugins=(/a\    zsh-syntax-highlighting' "$ZSHRC"
    fi
fi

echo "[INFO] Zsh 插件配置完成。"

# ------------------------------------------
# 重新加载 Zsh 配置
# ------------------------------------------

echo "[INFO] 检查 Zsh 配置..."

zsh -c 'source "$HOME/.zshrc"'

echo "[INFO] Zsh 配置加载成功。"

# ==========================================
# 7. 启动新的 Zsh Shell
# ==========================================

echo "[INFO] 启动新的 Zsh Shell，加载 ~/.zshrc..."

zsh

echo "[INFO] 新 Zsh Shell 已退出，继续执行初始化脚本..."

# ==========================================
# 后续功能将在这里继续追加
# ==========================================

echo "[INFO] WSL2 初始化脚本执行完成。"

