package mode

import (
	"awesomeProject/src/actions"
	"awesomeProject/src/client"
	"awesomeProject/src/common"
	"awesomeProject/src/server"
	"bufio"
	"fmt"
	"os"
	"time"
)

func ClientMode(reader *bufio.Reader) {
	for {
		fmt.Println("Por favor, informe o IP e porta do servidor para se conectar (exemplo: 127.0.0.1:8080)")

		targetServer := common.ReadAndSplitCommand(reader, ":")

		var err error
		targetPort, err := common.ConvertPort(targetServer[1])
		if err != nil {
			common.ClearScreen()
			fmt.Println("Erro ao converter porta:", err)
			fmt.Println("Por favor, tente novamente.")
			continue
		}

		targetIP := targetServer[0]

		c := client.NewClient(targetIP, targetPort)
		common.Conn, err = c.Connect()
		if err != nil {
			common.ClearScreen()
			fmt.Println("Erro ao conectar ao servidor:", err)
			fmt.Println("Por favor, tente novamente.")
			continue
		}

		break
	}

	actions.ReceiveAndSendShips()
	common.SaveOpponentShips()
}

func ServerMode(ip *string, port *int) {
	common.IsServer = true
	s := server.NewServer(*ip, *port)
	go func() {
		err := s.Start()
		if err != nil {
			fmt.Println("Erro ao iniciar servidor:", err)
			os.Exit(1)
		}
	}()
	time.Sleep(1 * time.Second)
	fmt.Println("Esperando conexão de cliente...")

	for {
		if common.Conn != nil {
			break
		}
	}

	actions.SendAndReceiveShips()
	common.SaveOpponentShips()
}
