package hub

import (
	"context"
	"embed"
	"log"
	"os"
	"os/signal"
	"sync"
	"xbsrebuild/hub/api"
)

var (
	//go:embed templates/*
	f embed.FS
)

func Run(server, port string) {
	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	controllerServer := api.NewServer(server, port, ctx)
	controllerServer.InitRouters(f)
	// 运行 ResultSrver 调度服务
	go func() {
		controllerServer.Run(&wg)
	}()

	// 等待中断信号以优雅地关闭服务器
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	cancel()
	wg.Wait()
	log.Println("Shutdown Server ...")
}
