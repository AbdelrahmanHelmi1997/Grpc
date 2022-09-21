package main

import (
	"test/dataBase"
	"test/server"
)

func main() {

	dataBase.DB()
	server.GrpcServerConnection()

}
