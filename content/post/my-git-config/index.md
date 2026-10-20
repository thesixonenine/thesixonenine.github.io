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
