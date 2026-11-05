---
title: gradle-tutorial
date: 2026-11-04T14:13:28
lastmod: 2026-11-04T16:26:18
tags: ['Java']
categories: ['Java']
keywords: Gradle
description: gradle tutorial
---

`只推荐使用 Gradle Wrapper`


- `GRADLE_USER_HOME`: 默认为 `~/.gradle`


## Proxy

[官方文档](https://docs.gradle.org/current/userguide/networking.html#networking_gradle)

推荐在 `GRADLE_USER_HOME` 下增加配置来设置代理

`~/.gradle/gradle.properties`

```properties
systemProp.http.proxyHost=127.0.0.1
systemProp.http.proxyPort=1080
systemProp.https.proxyHost=127.0.0.1
systemProp.https.proxyPort=1080
systemProp.socksProxyHost=127.0.0.1
systemProp.socksProxyPort=1080

org.gradle.jvmargs=-Xmx4g -Dfile.encoding=UTF-8
org.gradle.parallel=true
org.gradle.caching=true
org.gradle.warning.mode=all
org.gradle.configuration-cache=true

systemProp.file.encoding=UTF-8
```

或者设置 `GRADLE_OPTS` 环境变量, 这里可以解析 `HTTP_PROXY` 环境得到 `Host` 和 `Port` 并追加到 `~/.zshrc` 中

```bash
export GRADLE_OPTS="$GRADLE_OPTS -Dhttp.proxyHost=127.0.0.1 -Dhttp.proxyPort=1080 -Dhttps.proxyHost=127.0.0.1 -Dhttps.proxyPort=1080"
```

## Mirror

推荐在 `GRADLE_USER_HOME` 下增加配置来提前设置镜像源

`~/.gradle/init.d/mirrors.init.gradle.kts`

```kotlin
import org.gradle.api.artifacts.repositories.MavenArtifactRepository
import org.gradle.api.artifacts.repositories.IvyArtifactRepository
import org.gradle.api.artifacts.repositories.FlatDirectoryArtifactRepository

gradle.beforeSettings {
    val pluginMirrors = listOf(
        "https://maven.aliyun.com/repository/public",
        "https://maven.aliyun.com/repository/google",
        "https://maven.aliyun.com/repository/spring",
        "https://maven.aliyun.com/repository/gradle-plugin",
        "https://mirrors.tencent.com/nexus/repository/gradle-plugins"
    )
    val dependencyMirrors = listOf(
        "https://maven.aliyun.com/repository/public",
        "https://maven.aliyun.com/repository/google",
        "https://maven.aliyun.com/repository/spring",
        "https://maven.aliyun.com/repository/gradle-plugin",
        "https://mirrors.tencent.com/nexus/repository/gradle-plugins"
    )
    pluginMirrors.forEachIndexed { index, mirror ->
        pluginManagement.repositories.maven {
            name = "PluginMirror$index"
            url = uri(mirror)
        }
    }
    dependencyMirrors.forEachIndexed { index, mirror ->
        dependencyResolutionManagement.repositories.maven {
            name = "DependencyMirror$index"
            url = uri(mirror)
        }
    }
    gradle.settingsEvaluated {
        fun printRepositories(title: String, repositories: Iterable<*>) {
            println("\n[$title]")
            repositories.forEachIndexed { index, repo ->
                val location = when (repo) {
                    is MavenArtifactRepository -> repo.url.toString()
                    is IvyArtifactRepository -> repo.url.toString()
                    is FlatDirectoryArtifactRepository ->
                        repo.dirs.joinToString()
                    else -> "Unknown repository type"
                }
                println("${index + 1}. ${repo.let { (it as org.gradle.api.artifacts.repositories.ArtifactRepository).name }} -> $location")
            }
        }
        printRepositories(
            "Plugin Management",
            pluginManagement.repositories
        )
        printRepositories(
            "Dependency Resolution",
            dependencyResolutionManagement.repositories
        )
        println()
    }
}
println("[Gradle Mirror] Loaded")
```

## Command

### check mirrors

```bash
./gradlew --no-build-cache --no-configuration-cache
```

or

```bash
./gradlew build --refresh-dependencies --info 2>&1 | grep -E 'aliyun|mavenCentral|Downloading'
```

### use mirror to download gradle wrapper

```bash
find . -name 'gradle-wrapper.properties' -exec sed -i 's/services\.gradle\.org\/distributions/mirrors.cloud.tencent.com\/gradle/g' {} +
```

revert

```bash
find . -name 'gradle-wrapper.properties' -exec sed -i 's/mirrors\.cloud\.tencent\.com\/gradle/services.gradle.org\/distributions/g' {} +
```
