---
title: message-read-design
date: 2022-06-15T11:25:01+0800
updated: 2022-06-15T16:21:01+0800
tags: 
- Redis 
- Bitmap
categories: ['Redis']
keywords: 公告通知
description: 简单的公告通知模块的设计
url: '/p/message-read-design.html'
---

## 需求背景

近日项目的后台系统中需要新增公告通知功能. 

## 需求分析

公告条数不会很多, 用户数量级在十万左右, 是典型的一对多的场景, 且有明显的布尔特征, 即 已读/未读.


## 方案设计

可以考虑使用 **Redis** 存储关联关系并使用 **bitmap** 来进行用户是否阅读公告的记录.

快速入门[Redis中bitmap的使用](https://zhuanlan.zhihu.com/p/401726844)

这里摘录一些本次需要注意的:

> Bitmap 不属于 Redis 的基本数据类型, 而是基于 String 类型进行的位操作.
> 
> Redis 中字符串的最大长度是 512M, 所以 bitmap 的 offset 值也是有上限的, 其最大值是 8 * 1024 * 1024 * 512  =  2^32

由于不同端(运营, 用户, 顾客)的用户数据可能在相同的表, 也可能在不同的表, 且公告需要支持只通知到某一些端, 故 Redis 的key设计如下

```
notice:通知端数字枚举:数据库中的消息id
```

### 新增或修改公告

根据公告设置的不同端, 调用 setBit 方法进行设置.

考虑后面查询脚本的编写方便及性能, 可以在公告的不同端中设置固定用户为已读

```java
public static final String notice_key =  "notice:%d:%d";

// platform: 通知端数字枚举
// msgId: 数据库中的消息id
// userId: 不同用户表(运营, 用户, 顾客)的主键id
// isRead: 是否阅读
redisTemplate.opsForValue().setBit(String.format(notice_key, platform, msgId), userId, isRead);
```

### 删除公告

删除公告直接组装 该公告需要删除的端的key 成List, 然后调用delete的重载方法进行删除

```java
redisTemplate.delete(K key)
redisTemplate.delete(Collection<K> keys)
```


### 设置已读/未读

直接组装key, 设置 isRead 即可

```java
redisTemplate.opsForValue().setBit(String.format(notice_key, productType, msgId), userId, isRead);
```

### 查询用户是否已读

用户查看公告列表(分页)时, 直接简单的循环查询公告是否已读 

```java
redisTemplate.opsForValue().getBit(String.format(notice_key, productType, msgId), userId);
```

### 未读/已读数量统计

需要使用 lua 脚本循环所有消息来进行统计.

用户登录后, 根据用户所在端可以拼出该端所有的消息的key: 
```text
"notice:" + productType + ":*"
```

然后调用如下 lua 脚本, 即可统计出未读和已读数量

```java
private final String READ_COUNT_SCRIPT =
// 第一步: 拿到所有的消息key
"local noticeKeys = redis.call('keys', KEYS[1]);" +
"local unread = 0;" +
"local read = 0;" +
// 第二步: 如果key不存在则直接返回
"if next(noticeKeys) == nil then return unread .. ':' .. read;end;" +
// 第三步: 拿到所有消息key的值
"local values = redis.call('mget', unpack(noticeKeys));" +
// 第四步: 循环消息key, 使用getbit命令判断是否已读, 分别对未读数量和已读数量加一
"for i = 1, #noticeKeys do " +
"  if(redis.call('getbit', noticeKeys[i], ARGV[1]) == 0) " +
"    then unread = unread + 1;" +
"  else read = read + 1;" +
"  end;" +
"end;" +
// 第五步: 返回未读数量和已读数量(lua中使用 .. 来连接字符串)
"return unread .. ':' .. read;";
List<String> keys = new ArrayList<>(1);
keys.add("notice:" + productType + ":*");
String unreadAndRead = redisTemplate.execute(new DefaultRedisScript<>(READ_COUNT_SCRIPT, String.class),
    new StringRedisSerializer(), new StringRedisSerializer(), keys, String.valueOf(userId));
```

## 后续踩坑记

原来的 lua 脚本第一行是
```lua
local noticeKeys = redis.call('keys', 'notice:' .. KEYS[1] .. ':*');
```

导致该功能上线后, 立即爆出异常
```text
-ERR bad lua script for redis cluster, all the keys that the script uses should be passed using the KEYS arrayrn
```

搜索相关报错后发现是[阿里云的redis产品限制](https://developer.aliyun.com/article/645851), 所有key都应该由 KEYS 数组来传递.
