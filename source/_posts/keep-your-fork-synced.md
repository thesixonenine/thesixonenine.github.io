---
title: keep-your-fork-synced
date: 2021-02-20 11:20:21
updated: 2021-02-20 11:20:21
tags:
categories: Git
keywords: fork sync
description: fork项目后与原项目保持同步 
---

[参考链接](https://docs.github.com/en/github/getting-started-with-github/fork-a-repo#keep-your-fork-synced)

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

