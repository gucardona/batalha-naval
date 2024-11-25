package client

import (
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

func (c *Client) Connect() (net.Conn, error) {
	serverAddr := fmt.Sprintf("%s:%d", c.ip, c.port)

	timeout := 10 * time.Second
	conn, err := net.DialTimeout("tcp", serverAddr, timeout)
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar ao servidor: %s", err)
	}
	fmt.Println("Conexão estabelecida com:", conn.RemoteAddr())
	return conn, nil
}

func (c *Client) SendMessage(conn net.Conn, message string) error {
	_, err := conn.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("falha ao enviar mensagem: %s", err)
	}

	return nil
}
