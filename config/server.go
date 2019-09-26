package config

import (
	"github.com/gin-gonic/gin"
	"github.com/weikunlu/go-api-template/core/util"
)

type ServerConfig struct {
	Host string
	Port string
}

var serverConfig *ServerConfig

func init() {
	cfg := GetAppConfig()

	// Switch to "release" mode in production
	if cfg.AppEnv != "local" ||
		(util.GetCommandOfExecution() == "event" || util.GetCommandOfExecution() == "migrate") {
		gin.SetMode(gin.ReleaseMode)
	}

	serverConfig = &ServerConfig{
		Host: cfg.ServerHost,
		Port: cfg.ServerPort,
	}
}

func GetServerConfig() *ServerConfig {
	return serverConfig
}
