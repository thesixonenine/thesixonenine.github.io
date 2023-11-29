---
title: keep-your-fork-synced
date: 2021-02-20T11:20:21+0800
updated: 2021-02-22T09:56:21+0800
tags: ['Git']
categories: ['Git']
keywords: fork sync
description: fork项目后与原项目保持同步
---

[参考链接](https://docs.github.com/en/github/getting-started-with-github/fork-a-repo#keep-your-fork-synced)

> 将远程原始repo合并到远程fork repo的本地repo后推送到远程fork repo
>
> 1. clone远程fork repo
> 2. 添加远程原始repo为upstream
> 3. fetch&merge远程原始repo的分支到本地分支
> 4. push到远程fork repo分支

#### 具体操作

1. Set up Git
2. Create a local clone of your fork

```bash
git clone https://github.com/YOUR-USERNAME/REPO_NAME.git
```

3. Configure Git to sync your fork with the original repository

```bash
# check the current configured remote repository for your fork.
git remote -v
# add upstream
git remote add upstream https://github.com/ORIGINAL-USERNAME/REPO_NAME.git
# verify the new upstream repository
git remote -v
```

4. Fetch and merge original repository into your local repository

```bash
git fetch upstream
git merge upstream/master
```

5. push to your remote fork repo

```bash
git push origin master
```

