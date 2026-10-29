---
title: "my-git-config"
date: 2026-10-20T09:57:18+08:00
lastmod: 2026-10-20T09:57:52+0800
categories: ['Git']
keywords: git config
description: my git config
---

## Linux

### 全局配置

```bash
git config --global alias.co checkout && \
git config --global alias.ci commit && \
git config --global alias.br branch && \
git config --global alias.st status && \
git config --global alias.ps push && \
git config --global alias.pl pull && \
git config --global alias.ft fetch && \
git config --global alias.mg merge && \
git config --global alias.lg "log --color --graph --pretty=format:'%Cred%h%Creset -%C(yellow)%d%Creset %s %Cgreen(%cr) %C(bold blue)<%an>%Creset' --abbrev-commit"
```

### 个人配置

`放到全局配置中`

```bash
git config --global user.name Simple && \
git config --global user.email thesixonenine@outlook.com && \
git config --global user.signingkey ~/.ssh/id_ed25519.pub && \
git config --global commit.gpgsign true && \
git config --global tag.gpgsign true && \
git config --global gpg.format ssh && \
git config --global gpg.ssh.allowedSignersFile ~/.ssh/allowed_signers && \
git config --global pull.rebase true
```

`放到仓库配置中`

```bash
git config user.name Simple && \
git config user.email thesixonenine@outlook.com && \
git config user.signingkey ~/.ssh/id_ed25519.pub && \
git config commit.gpgsign true && \
git config tag.gpgsign true && \
git config gpg.format ssh && \
git config gpg.ssh.allowedSignersFile ~/.ssh/allowed_signers && \
git config pull.rebase true
```

## Windows

### 全局配置

```powershell
git config --global alias.co checkout && `
git config --global alias.ci commit && `
git config --global alias.br branch && `
git config --global alias.st status && `
git config --global alias.ps push && `
git config --global alias.pl pull && `
git config --global alias.ft fetch && `
git config --global alias.mg merge && `
git config --global alias.lg "log --color --graph --pretty=format:'%Cred%h%Creset -%C(yellow)%d%Creset %s %Cgreen(%cr) %C(bold blue)<%an>%Creset' --abbrev-commit"
```

### 个人配置

`放到全局配置中`

```powershell
git config --global user.name Simple && `
git config --global user.email thesixonenine@outlook.com && `
git config --global user.signingkey ~/.ssh/id_ed25519.pub && `
git config --global commit.gpgsign true && `
git config --global tag.gpgsign true && `
git config --global gpg.format ssh && `
git config --global gpg.ssh.allowedSignersFile ~/.ssh/allowed_signers && `
git config --global pull.rebase true
```

`放到仓库配置中`

```powershell
git config user.name Simple && `
git config user.email thesixonenine@outlook.com && `
git config user.signingkey ~/.ssh/id_ed25519.pub && `
git config commit.gpgsign true && `
git config tag.gpgsign true && `
git config gpg.format ssh && `
git config gpg.ssh.allowedSignersFile ~/.ssh/allowed_signers && `
git config pull.rebase true
```

## gitattributes

```txt
# reset gitattribute
# git rm --cached -r
# git reset --hard

* text=auto eol=lf

# lang
*.java text eol=lf
*.go text eol=lf
*.groovy text eol=lf
*.gradle text eol=lf
*.kt text eol=lf
*.kts text eol=lf
*.rs text eol=lf
*.html text eol=lf
*.css text eol=lf
*.js text eol=lf
*.jsx text eol=lf
*.ts text eol=lf
*.tsx text eol=lf
*.nix text eol=lf
*.svg text eol=lf

# config
*.xml text eol=lf
*.yml text eol=lf
*.yaml text eol=lf
*.toml text eol=lf
*.properties text eol=lf
/.gitattributes text eol=lf
/.gitignore text eol=lf
*.json text eol=lf
*.jsonc text eol=lf
*.pro text eol=lf

# framework
bun.lock text eol=lf
Cargo.lock text eol=lf

# document
*.md text eol=lf

# Linux start script should use lf
/gradlew text eol=lf

# These are Windows script files and should use crlf
*.bat text eol=crlf

# binary
*.jar binary
*.png binary
*.ico binary
*.icns binary
```
