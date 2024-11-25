package main

import (
	"awesomeProject/src/server"
	"bufio"
	"flag"
	"fmt"
	"os"
	"time"
)

var (
	reader = bufio.NewReader(os.Stdin)
)

func main() {
	ip := flag.String("ip", "localhost", "IP do servidor")
	port := flag.Int("port", 8443, "Porta do servidor")
	flag.Parse()

	s := server.NewServer(*ip, *port)

	go func() {
		err := s.Start()
		if err != nil {
			fmt.Println("Erro ao iniciar servidor:", err)
			os.Exit(1)
		}
	}()

	time.Sleep(1 * time.Second)

	s.Handle()
}
