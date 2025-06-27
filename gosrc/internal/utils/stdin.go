package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

// ReadContinuedLinesStdin reads continued lines from standard input.
//
// This is a convenience wrapper for ReadContinuedLines that uses os.Stdin.
// It prints the tips prompt and returns the concatenated input when
// a non-continued line is entered.
//
// Usage example:
//
//	input := ReadContinuedLinesStdin("Enter JSON (end with non-\\ line):")
//	// User enters:
//	// {\
//	//   "key": "value"\
//	// }
func ReadContinuedLinesStdin(tips string) string {
	return ReadContinuedLines(tips, os.Stdin)
}

// ReadContinuedLines reads multiple lines from an io.Reader until encountering
// a line that doesn't end with a continuation backslash.
//
// It prints the provided tips message before reading. Each line is processed as:
//   - Trailing newline is removed
//   - Trailing whitespace is trimmed for continuation check
//   - Final backslash (if present) is removed from continuation lines
//
// Lines are joined with newline separators in the returned string.
//
// Example termination condition:
//
//	"input\"   -> continues (backslash after trimming)
//	"input "   -> terminates (space after backslash: "input \")
//	"input"    -> terminates
func ReadContinuedLines(tips string, rd io.Reader) string {
	fmt.Println(tips)
	reader := bufio.NewReader(rd)
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

// KeepHeadAndAppend 保留开头几行并追加指定内容
func KeepHeadAndAppend(filePath string, keepHeadLine int, appendContent string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	var preservedLines []string
	scanner := bufio.NewScanner(file)
	for i := 0; i < keepHeadLine && scanner.Scan(); i++ {
		preservedLines = append(preservedLines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	file, err = os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	for _, line := range preservedLines {
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return err
		}
	}
	if _, err := writer.WriteString(appendContent); err != nil {
		return err
	}
	return writer.Flush()
}
