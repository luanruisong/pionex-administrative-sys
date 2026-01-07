package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"pionex-administrative-sys/server"
	"pionex-administrative-sys/utils/app"
	"pionex-administrative-sys/utils/app/daemon"
	"pionex-administrative-sys/utils/logger"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	defer func() {
		fmt.Println(logger.Sync())
	}()
	app.Parse()
	if app.Daemon() {
		d, err := daemon.Daemon()
		if err != nil {
			logger.Fatal("Failed to start daemon", zap.Error(err))
		}
		logger.Info("daemon started", zap.Int("pid", d.Pid), zap.Strings("args", d.Args))
		return
	}

	srv := server.New(app.Port())
	srv.Init()

	go func() {
		if err := srv.Run(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("run fatal", zap.Error(err))
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)

}
