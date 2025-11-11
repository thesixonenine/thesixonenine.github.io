package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// YuanToFen 将 确定是单位是元的字符串 转成 单位是分的int
func YuanToFen(s string) (int, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty string")
	}

	// 检查是否带有小数点
	if strings.Contains(s, ".") {
		parts := strings.SplitN(s, ".", 2)
		intPart := parts[0]
		fracPart := parts[1]

		// 处理小数部分长度
		if len(fracPart) > 2 {
			return 0, fmt.Errorf("too many decimal places")
		}

		// 补齐到两位
		for len(fracPart) < 2 {
			fracPart += "0"
		}

		// 拼接整数部分和小数部分
		s = intPart + fracPart
	} else {
		// 没有小数点，就直接加上两位零
		s = s + "00"
	}
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("invalid number: %w", err)
	}

	return val, nil
}

// DivideHundred 将 数字 除以100 并 转成字符串
func DivideHundred(n int) string {
	sign := ""
	if n < 0 {
		sign = "-"
		n = -n
	} else if n == 0 {
		return "0"
	}
	s := strconv.Itoa(n)

	// 根据长度插入小数点
	l := len(s)
	if l == 1 {
		return sign + "0.0" + s
	} else if l == 2 {
		return sign + "0." + s
	} else {
		return sign + s[:l-2] + "." + s[l-2:]
	}
}
