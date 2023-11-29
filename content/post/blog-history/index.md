---
title: "blog-history"
date: 2023-11-29T17:13:10+08:00
lastUpdated: 2023-11-29T17:14:10+08:00
categories: ['Blog']
keywords: blog
description: blog history
---

今天将blog的markdown文章从`文章名.md`改成了`Page Bundles`模式.

hugo官方对`Page Bundles`模式的介绍在[这里](https://gohugo.io/content-management/page-bundles)

```bash
#/bin/bash
# 从 文章名.md 改成 Page Bundles模式
find . -type f -name "*.md" -exec sh -c 'mv "$1" "$(dirname "$1")/index.md"' _ {} \;
# 从 Page Bundles模式 改成 文章名.md
for file in */*.md; do mv "$file" "$(dirname "$file")/$(basename "$(dirname "$file").md")"; done
```
