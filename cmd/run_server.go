package cmd

import (
	"github.com/weikunlu/go-api-template/server"
)

func RunServer() error {

	//err := database.NewDatabase()
	//if err != nil {
	//	fmt.Printf("init error %v", err.Error())
	//	return nil
	//}
	//defer database.GetDb().Close()

	server.StartServer(false)
	return nil
}
