package main

import (
	"fmt"
	"pokemon/internal/network"
)

func main() {
	client := network.NewClient("localhost:8080")
	fmt.Println("Connecting to server...")
	client.Start()
}
