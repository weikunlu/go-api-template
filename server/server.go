package server

import (
	"fmt"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/weikunlu/go-api-template/config"
	"io/ioutil"
	"os"
	"syscall"
)

var Router *gin.Engine

func init() {
	Router = SetupRouter()
}

func writePidFile(pidFile string) error {
	return ioutil.WriteFile(pidFile, []byte(fmt.Sprintf("%d", os.Getpid())), 0664)
}

func StartServer(graceful bool) {
	cfg := config.GetServerConfig()
	endPoint := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	if graceful {
		server := endless.NewServer(endPoint, Router)
		server.BeforeBegin = func(add string) {
			fmt.Printf("Actual pid is %d\n", syscall.Getpid())
			writePidFile("socket.pid")
		}

		err := server.ListenAndServe()
		if err != nil {
			fmt.Printf("Server err: %v\n", err)
		}
	} else {
		Router.Run(endPoint)
	}
}
