package main

import (
	"distributed-order-system-go/internal/server"
	"distributed-order-system-go/pkg/bootstrap"
	"distributed-order-system-go/pkg/global"
	"distributed-order-system-go/pkg/zookeeper"
)

func main() {

	bootstrap.InitializeConfig()

	global.App.DB = bootstrap.InitializeDB()

	global.App.Redis = bootstrap.InitializeRedis()

	global.App.DistributedLocker = zookeeper.NewZookeeperLock()

	// Start the server
	server.StartServer()
}
