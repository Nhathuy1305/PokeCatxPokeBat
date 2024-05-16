package network

import (
	"fmt"
	"net"
)

type Client struct {
	address string
}

func NewClient(address string) *Client {
	return &Client{address: address}
}

func (c *Client) Start() {
	conn, err := net.Dial("tcp", c.address)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to server:", c.address)
	// handle server communication
}
