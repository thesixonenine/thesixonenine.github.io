---
title: "rust-tutorial"
date: 2026-02-13T14:08:57+08:00
lastmod: 2026-02-13T14:08:57+08:00
categories: ['']
tags: ['']
keywords: 
description: 
image: 
---

## Install

### Windows

```powershell
winget install Microsoft.VisualStudio.2022.BuildTools --override "--quiet --wait --norestart --add Microsoft.VisualStudio.Workload.VCTools --includeRecommended"
```

```powershell
scoop install rustup-msvc
```

**check version**

```powershell
rustc --version
cargo --version
```

**hello world**

```powershell
cargo new hello
cd hello
cargo run
```

