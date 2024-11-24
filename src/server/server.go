package server

import (
	"fmt"
	"net"
)

type Server struct {
	ip   string
	port int
}

func NewServer(ip string, port int) *Server {
	return &Server{
		ip:   ip,
		port: port,
	}
}

func (s *Server) Start() {
	serverAddr := fmt.Sprintf("%s:%d", s.ip, s.port)

	ln, err := net.Listen("tcp", serverAddr)
	if err != nil {
		fmt.Println("Falha ao criar servidor:", err)
		return
	}
	defer ln.Close()

	fmt.Println("Servidor escutando em:", serverAddr)
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Falha ao aceitar conexão:", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Conexão estabelecida com:", conn.RemoteAddr())
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Erro ao ler dados do buffer:", err)
			return
		}
		fmt.Printf("Recebi: %s\n", buffer[:n])
	}
}
