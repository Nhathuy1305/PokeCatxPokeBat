package main

import (
	"log"
	"pokemon/config"
	"pokemon/internal/network"
)

func main() {
	config.InitConfig()
	server := network.NewServer(config.Conf.ServerAddress)
	log.Println("Starting server on", config.Conf.ServerAddress)
	server.Start()
}
