---
title: "show-busy-java-threads"
date: 2024-12-06T14:18:32+08:00
lastmod: 2025-06-19T14:30:08+0800
categories: ['Java']
tags: ['Java', 'Linux']
keywords: ['Java', 'Linux', 'CPU']
description: 排查 Java 进程占用 CPU 高的堆栈
isCJKLanguage: true
---

> 更多实用脚本[参考](https://github.com/oldratlee/useful-scripts)

临时到服务器上排查 Java 进程占用 CPU 高的原因, 步骤如下:


1. 查询 CPU 占用高的进程 PID 为 X
2. 查询进程 X 中的线程资源占用情况
3. 找到占用量最高的线程 PID 为 Y
4. 计算 Y 对应的十六进制为 Z
5. 查询进程 X 中各线程的调用栈
6. 定位到 Z 所在的线程堆栈



```bash
# top 命令默认是按 CPU 的占用量从高到低排序
top
# 查询占用高的进程 X 中的线程资源占用情况
top -Hp X
# 占用量最高的线程 PID 为 Y
# 计算 Y 对应的十六进制为 Z
printf "%x\n" Y

# 查询进程中各线程的调用栈并定位到 Z
jstack X
```
