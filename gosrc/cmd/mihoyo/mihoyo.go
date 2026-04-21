package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	runCommand("go", "run", "./cmd/genshin-impact")
	runCommand("go", "run", "./cmd/star-rail")
	runCommand("go", "run", "./cmd/zzz")
	runCommand("go", "run", "./cmd/pre")
}
func runCommand(name string, args ...string) {
	fmt.Printf("正在运行: %s %v\n", name, args)
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// Run 会阻塞直到命令执行完成
	err := cmd.Run()
	if err != nil {
		log.Printf("命令执行失败: %v\n", err)
	}
}
