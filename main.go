package main

import (
	"fmt"
	"github.com/smarterwallet/demand-abstraction-serv/global"
	"github.com/smarterwallet/demand-abstraction-serv/route"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/smarterwallet/demand-abstraction-serv/config"
)

func main() {
	var env = os.Getenv("GO_ENV")
	if env != "production" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	global.InitLogger()

	cfg := &config.Config{}
	if err := config.LoadConfig(cfg); err != nil {
		panic(err)
	}
	server := route.NewHTTPServer(cfg)
	server.Start()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("Shutting down server...")

	server.Stop()
	fmt.Println("Server exiting")
}
