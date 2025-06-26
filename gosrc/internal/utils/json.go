package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"
)

// JSONIndent 进行 JSON 格式化
func JSONIndent(marshal []byte) []byte {
	var out bytes.Buffer
	err := json.Indent(&out, marshal, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	return out.Bytes()
}

func ReadJSONFile[T any](path string) T {
	file, err := os.Open(path)
	if err != nil {
		// 文件不存在时创建空文件
		if os.IsNotExist(err) {
			file, err = os.Create(path)
			if err != nil {
				log.Fatalf("创建文件[%s]失败: %s", path, err)
			}
		} else {
			log.Fatalf("打开文件[%s]失败: %s", path, err)
		}
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatalf("关闭文件[%s]异常,err[%s]\n", path, err)
		}
	}()
	content, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("读取文件[%s]异常,err[%s]\n", path, err)
	}
	var m T
	if len(content) == 0 {
		return m
	}
	err = json.Unmarshal(content, &m)
	if err != nil {
		log.Fatalf("解析JSON失败[%s]:err[%s]", path, err)
	}
	return m
}
