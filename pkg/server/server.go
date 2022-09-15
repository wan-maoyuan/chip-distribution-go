package server

import (
	"chip-distribution-go/pkg/service"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ChipServer struct {
	g          *gin.Engine
	server     *http.Server
	service    *service.ChipService
	singalStop chan os.Signal
}

func NewChipServer(addr string) (*ChipServer, error) {
	var server = &ChipServer{}

	service, err := service.NewChipService()
	if err != nil {
		return server, err
	}

	server.service = service
	server.singalStop = make(chan os.Signal)
	server.g = gin.Default()
	server.server = &http.Server{
		Addr:    addr,
		Handler: server.g,
	}

	return server, nil
}

func (server *ChipServer) Run() {
	server.registerApi()

	go func() {
		if err := server.server.ListenAndServe(); err != nil {
			logrus.Errorf("http Server Run error: %v", err)
		}
	}()

	signal.Notify(server.singalStop, syscall.SIGINT, syscall.SIGTERM)
	<-server.singalStop
	logrus.Info("chip server shutdown ......")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	server.clean(ctx)
}

func (server *ChipServer) registerApi() {
	server.g.POST("/upload_excel", server.service.UploadExcel)
	server.g.GET("/shutdown", func(c *gin.Context) {
		go server.shutdown()
		c.String(http.StatusOK, "OK")
	})
	logrus.Info("register chip api route success")
}

func (server *ChipServer) shutdown() {
	time.Sleep(time.Second)
	server.singalStop <- syscall.SIGINT
}

func (server *ChipServer) clean(ctx context.Context) {
	if err := server.server.Shutdown(ctx); err != nil {
		logrus.Errorf("http server shutdown error: %v", err)
	}

	server.service.Close()
}
