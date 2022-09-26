package main

import (
	"chip-distribution-go/pkg/server"
	"fmt"

	"github.com/sirupsen/logrus"
)

var (
	Name       string = "chip-distribution"
	Version    string //版本
	CommitHash string //git 提交的 hash 值
	BuildTime  string //编译时间
)

func main() {
	showVersion()

	svr, err := server.NewChipServer(":8080")
	if err != nil {
		logrus.Fatalf("create chip server error: %v", err)
	}

	svr.Run()
}

func showVersion() {
	fmt.Printf(`
		Name: %s
		Version: %s
		CommitHash: %s
		BuildTime: %s

`, Name, Version, CommitHash, BuildTime)
}
