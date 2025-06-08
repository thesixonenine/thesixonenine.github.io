package utils

import (
	"log"
	"os"
	"syscall"
)

// WriteToFile 写入文件
func WriteToFile(path string, marshal []byte) {
	err := os.WriteFile(path, marshal, syscall.O_RDWR|syscall.O_CREAT)
	if err != nil {
		log.Fatalf("写入文件异常[%s]\n", err.Error())
	}
}
