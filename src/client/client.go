package client

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

type Client struct {
	ip   string
	port int
}

func NewClient(ip string, port int) *Client {
	return &Client{
		ip:   ip,
		port: port,
	}
}

func (c *Client) ConnectToServer() error {
	serverAddr := fmt.Sprintf("%s:%d", c.ip, c.port)

	timeout := 10 * time.Second
	conn, err := net.DialTimeout("tcp", serverAddr, timeout)
	if err != nil {
		return err
	}
	defer conn.Close()

	fmt.Println("Cliente conectado ao servidor:", serverAddr)
	return nil
}

func (c *Client) sendMessage(conn net.Conn, message string) {
	_, err := conn.Write([]byte(message + "\n"))
	if err != nil {
		fmt.Println("Failed to send message:", err)
		return
	}
	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Failed to receive response:", err)
		return
	}
	fmt.Println("Response received:", response)
}
