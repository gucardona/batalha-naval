package main

import (
	"awesomeProject/src/battlefield"
	"awesomeProject/src/client"
	"awesomeProject/src/common"
	"awesomeProject/src/server"
	"awesomeProject/src/ships"
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	ip := flag.String("ip", "localhost", "IP do servidor")
	port := flag.Int("port", 8443, "Porta do servidor")
	flag.Parse()

	s := server.NewServer(*ip, *port)

	go s.Start()
	time.Sleep(1 * time.Second)

	b := battlefield.NewBattlefield()
	var (
		shipList  []ships.Ship
		jsonShips []common.ShipJSON
	)

	fmt.Println("\nHora de colocar seus navios no campo de batalha!")
	reader := bufio.NewReader(os.Stdin)

	for {
		rootDir, err := os.Getwd()
		if err != nil {
			common.ClearScreen()
			fmt.Println("Erro ao obter o diretório root do projeto:", err)
			fmt.Println("Por favor, insira novamente os navios...")
			break
		}

		_, err = os.Stat(fmt.Sprintf("%s/ships.json", rootDir))
		fmt.Println(err)
		if err == nil {
			fmt.Println("Arquivo \"ships.json\" não deletado encontrado. Deseja ler navios do arquivo? (s/n)")
			read, _ := reader.ReadString('\n')
			read = strings.TrimSpace(read)

			if read == "s" || read == "S" {
				file, err := os.ReadFile("ships.json")
				if err != nil {
					common.ClearScreen()
					fmt.Printf("falha ao ler arquivo: %s\nPor favor, insira novamente os navios...", err)
					break
				}

				if err = json.Unmarshal(file, &jsonShips); err != nil {
					common.ClearScreen()
					fmt.Printf("falha ao fazer Unmarshal do json: %s\nPor favor, insira novamente os navios...", err)
					break
				}
			}
		}
		common.ClearScreen()
		break
	}

	if len(jsonShips) <= 0 {
		for {
			fmt.Println("Deseja colocar os navios manualmente ou automaticamente? (m/a)")
			manual, _ := reader.ReadString('\n')
			manual = strings.TrimSpace(manual)

			switch manual {
			case "m", "M":
				shipList = common.ManuallyPlaceShips(reader, b)
				break
			case "a", "A":
				shipList = common.AutomaticallyPlaceShips(b)
				break
			default:
				common.ClearScreen()
				fmt.Println("Opção inválida. Por favor, tente novamente.")
				continue
			}
			break
		}

		generatedJson, _ := common.GenerateJSON(shipList)
		_ = common.WriteJSONToFile(generatedJson)
	}

	var err error
	shipList, err = common.ConvertJsonToShip(jsonShips)
	if err != nil {
		fmt.Println("Erro ao converter JSON para navios:", err)
	}
	for _, ship := range shipList {
		if err = b.PlaceShip(ship); err != nil {
			fmt.Println("Erro ao colocar navio no campo de batalha:", err)
		}
	}

	common.ClearScreen()
	fmt.Println("Seu campo de batalha:")
	b.Display()
	fmt.Println()

	var (
		targetIP   string
		targetPort int
	)
	for {
		fmt.Println("Por favor, informe o IP e porta do servidor para se conectar (exemplo: 127.0.0.1:8080)")

		targetServer := common.ReadAndSplitCommand(reader, ":")

		targetPort, err = common.ConvertPort(targetServer[1])
		if err != nil {
			continue
		}

		targetIP = targetServer[0]

		c := client.NewClient(targetIP, targetPort)
		if err = c.ConnectToServer(); err != nil {
			common.ClearScreen()
			fmt.Printf("Erro ao conectar ao servidor %s:%d: %s\n", targetIP, targetPort, err)
			continue
		}

		break
	}
}
