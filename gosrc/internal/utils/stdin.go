package utils

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func ReadMultiLine(tips string) string {
	log.Println(tips)
	reader := bufio.NewReader(os.Stdin)
	var lines []string

	// 读取多行输入直到空行
	for {
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, " ")
}
