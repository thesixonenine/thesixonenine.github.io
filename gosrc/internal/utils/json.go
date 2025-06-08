package utils

import (
	"bytes"
	"encoding/json"
	"log"
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
