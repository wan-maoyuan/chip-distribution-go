package main

import (
	"chip-distribution-go/pkg/server"

	"github.com/sirupsen/logrus"
)

func main() {
	svr, err := server.NewChipServer(":8080")
	if err != nil {
		logrus.Fatalf("create chip server error: %v", err)
	}

	svr.Run()
}
