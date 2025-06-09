package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

func ReadMultiLine(tips string) string {
	fmt.Println(tips)
	reader := bufio.NewReader(os.Stdin)
	var lines []string
	for {
		line, _ := reader.ReadString('\n')
		line = strings.TrimSuffix(line, "\n")
		lines = append(lines, line)
		if !strings.HasSuffix(strings.TrimRightFunc(line, unicode.IsSpace), "\\") {
			break
		}
	}
	return strings.Join(lines, "\n")
}
