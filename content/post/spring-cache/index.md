---
title: spring-cache
date: 2020-10-09T22:44:06+0800
lastmod: 2021-01-25T14:23:28+0800
tags:
- Java
- Spring
categories: ['Spring']
keywords: Spring Cache
description: Spring Cache相关
---

#### 参考文档

[spring-framework-reference](https://docs.spring.io/spring-framework/docs/current/spring-framework-reference/integration.html#cache)

#### 缓存抽象

- `org.springframework.cache.Cache`
- `org.springframework.cache.CacheManager`

**CacheManager**的作用

- 区分不同的缓存存储(`ConcurrentMapCacheManager`, `RedisCacheManager`)
- 管理同一缓存存储下的多个`Cache`

#### 缓存注解

- `@Cacheable` 缓存数据
- `@CacheEvict` 删除缓存
- `@CachePut` 更新缓存(不影响方法执行)
- `@Caching` 组合多个缓存操作
- `@CacheConfig` 在类级别共享一些与缓存相关的配置



