---
title: "redeem-code"
date: 2023-07-31T11:53:33+08:00
updated: 2023-07-31T15:40:33+08:00
tags:
  - Redis
categories: ['Redis']
keywords: redeem code
description: 兑换码设计
url: '/p/redeem-code.html'
---

兑换码的设计

## 需求背景

每个活动需要多次生成兑换码批次并导出兑换码,且记录兑换码的使用状态.每个兑换码只能使用一次.

## 需求分析

- 兑换码的生成要做到防暴力破解和防猜.
- 使用redis记录兑换码的使用情况 未使用/已使用
- 使用redis锁住该兑换码,避免兑换码被重复使用

## 方案设计

**使用数字+24位字母(排除I和O),共34个字符组成的13位字符串来作为兑换码.**

- 2位随机码
- 5位用做批次id:共支持45,435,424个批次(批次表的主键id)
- 4位用做库存顺序码共支持1,336,336
- 2位校验码,由前面的11位字符串进行CRC16编码得到int数,然后对34的2次方进行取余,确保生成的34进制字符最多两个,如果不足则补0
- 最后将以上13位字符按预定的打乱规则进行打乱,规则可以由13位中的某一位来指定,该位不进行打乱操作即可.

## 工具类

#### 34进制的定义

```java
// 34进制的字符存储
private static final char[] chars = new char[34];
static {
    int in = 0;
    // 0 - 9
    for (int i = 48; i <= 57; i++) {
        chars[in] = (char) i;
        in++;
    }
    // A - Z
    for (int i = 65; i <= 90; i++) {
        if (i == 73 || i == 79) {
            // I O
            continue;
        }
        chars[in] = (char) i;
        in++;
    }

}
```

#### 将十进制的数转换为34进制

```java
/**
 * 将十进制的数转换成34进制
 *
 * @param number 十进制数
 * @param length 34进制数的总长度, 
 *               如果转换后的不足该长度则补0, 
 *               超出该长度则删除末尾超出的部分
 * @return 34进制数
 */
private static String transRadix10To34(int number, int length) {
    int radix = 34;
    StringBuilder sb = new StringBuilder();
    while (number != 0) {
        sb.append(chars[number % radix]);
        number = number / radix;
    }
    if (length > 0) {
        while (sb.length() != length) {
            if (sb.length() < length) {
                sb.append(chars[0]);
            }
            if (sb.length() > length) {
                sb.deleteCharAt(sb.length() - 1);
            }
        }
    }
    return sb.reverse().toString();
}
```

#### 将34进制数转换成十进制

```java
/**
 * 将34进制数转换成十进制
 *
 * @param s 34进制数
 * @return 十进制数
 */
private static int transRadix34To10(String s) {
    int radix = 34;
    int sum = 0;

    int m = 1;
    int length = s.length();
    while (length > 0) {
        length--;
        char c = s.charAt(length);
        int index = findCharIndex(c);
        sum += index * m;
        m *= radix;
    }
    return sum;
}
/**
 * 根据字符查询在34进制中的索引位置
 */
private static int findCharIndex(char c) {
    int index = -1;
    for (int i = 0; i < chars.length; i++) {
        if (chars[i] == c) {
            index = i;
            break;
        }
    }
    return index;
}
```

#### 13位的兑换码的正则校验

```java
/**
 * 13位的兑换码的正则校验
 */
private static final Pattern CDKEY_PATTERN = Pattern.compile("^[0-9A-HJ-NP-Z]{13}$");
```