package server

import (
	"awesomeProject/src/common"
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

func (s *Server) Start() error {
	serverAddr := fmt.Sprintf("%s:%d", s.ip, s.port)

	ln, err := net.Listen("tcp", serverAddr)
	if err != nil {
		return fmt.Errorf("falha ao iniciar servidor: %s", err)
	}
	defer ln.Close()

	fmt.Println("Servidor escutando em:", serverAddr)
	for {
		common.Conn, err = ln.Accept()
		if err != nil {
			fmt.Println("Falha ao aceitar conexão:", err)
			continue
		}
	}
}
