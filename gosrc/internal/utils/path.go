package utils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

func LastModVerPath(p string) (string, error) {
	entries, err := os.ReadDir(p)
	if err != nil {
		return "", err
	}
	type dirInfo struct {
		name     string
		modTime  time.Time
		filePath string
	}
	var validDirs []dirInfo

	for _, entry := range entries {
		if entry.IsDir() {
			dirName := entry.Name()
			if isValidDirName(dirName) {
				// 获取完整的目录路径
				fullPath := filepath.Join(p, dirName)

				// 获取目录信息
				info, err := os.Stat(fullPath)
				if err != nil {
					fmt.Printf("无法获取目录信息: %s - %v\n", dirName, err)
					continue
				}
				validDirs = append(validDirs, dirInfo{
					name:     dirName,
					modTime:  info.ModTime(),
					filePath: fullPath,
				})
			}
		}
	}

	if len(validDirs) == 0 {
		return "", errors.New("未找到符合条件的目录")
	}

	// 按修改时间排序
	sort.Slice(validDirs, func(i, j int) bool {
		return validDirs[i].modTime.After(validDirs[j].modTime)
	})

	// 获取最近修改的目录
	recentDir := validDirs[0]

	//fmt.Printf("最近修改的目录: %s\n", recentDir.name)
	//fmt.Printf("修改时间: %s\n", recentDir.modTime.Format("2006-01-02 15:04:05"))
	//fmt.Printf("完整路径: %s\n", recentDir.filePath)
	return recentDir.name, nil
}
func isValidDirName(name string) bool {
	if name == "" {
		return false
	}

	hasDigit := false
	hasDot := false
	for _, ch := range name {
		switch {
		case ch >= '0' && ch <= '9':
			hasDigit = true
		case ch == '.':
			hasDot = true
		default:
			return false
		}
	}
	return hasDigit && hasDot
}
