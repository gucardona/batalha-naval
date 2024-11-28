package main

import (
	"awesomeProject/src/actions"
	"awesomeProject/src/common"
	"awesomeProject/src/mode"
	"awesomeProject/src/ships"
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	common.ShipList = []*ships.Ship{}

	common.ClearScreen()
	common.ShowIntro()
	ip := flag.String("ip", "localhost", "IP do servidor")
	port := flag.Int("port", 8080, "Porta do servidor")
	flag.Parse()

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\nHora de colocar seus navios no campo de batalha!")

	for {
		rootDir, err := os.Getwd()
		if err != nil {
			common.ClearScreen()
			fmt.Println("Erro ao obter o diretório root do projeto:", err)
			fmt.Println("Por favor, insira novamente os navios...")
			break
		}

		_, err = os.Stat(fmt.Sprintf("%s/ships.json", rootDir))
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

				if err = json.Unmarshal(file, &common.JsonShips); err != nil {
					common.ClearScreen()
					fmt.Printf("falha ao fazer Unmarshal do json: %s\nPor favor, insira novamente os navios...", err)
					break
				}
			}
		}
		common.ClearScreen()
		break
	}

	if len(common.JsonShips) <= 0 {
		for {
			fmt.Println("Deseja colocar os navios manualmente ou automaticamente? (m/a)")
			manual, _ := reader.ReadString('\n')
			manual = strings.TrimSpace(manual)

			switch manual {
			case "m", "M":
				common.ShipList = common.ManuallyPlaceShips(reader, common.B)
				break
			case "a", "A":
				common.ShipList = common.AutomaticallyPlaceShips(common.B)
				break
			default:
				common.ClearScreen()
				fmt.Println("Opção inválida. Por favor, tente novamente.")
				continue
			}
			break
		}
	} else {
		var err error
		common.ShipList, err = common.ConvertJsonToShip(common.JsonShips)
		if err != nil {
			fmt.Println("Erro ao converter JSON para navios:", err)
		}
		for _, ship := range common.ShipList {
			if err = common.B.PlaceShip(ship); err != nil {
				fmt.Println("Erro ao colocar navio no campo de batalha:", err)
			}
		}
	}
	jsonShip, _ := common.GenerateJSON(common.ShipList)
	_ = common.WriteJSONToFile(jsonShip)

	if common.Conn == nil {
		for {
			ok := false

			fmt.Println("Deseja atuar como servidor ou cliente? (s/c)")
			serverOrClient, _ := reader.ReadString('\n')
			serverOrClient = strings.TrimSpace(serverOrClient)

			switch serverOrClient {
			case "c", "C":
				mode.ClientMode(reader)
				ok = true
				break
			case "s", "S":
				mode.ServerMode(ip, port)
				ok = true
				break
			default:
				common.ClearScreen()
				fmt.Println("Opção inválida. Por favor, tente novamente.")
			}

			if ok {
				break
			}
		}
	}

	if !common.IsServer {
		common.IsMyTurn = false
		defer common.Conn.Close()
	}

	for {
		fmt.Println("Seu campo de batalha:")
		fmt.Println()
		common.B.Display()
		fmt.Println()
		fmt.Println("Campo de batalha do adversário:")
		fmt.Println()
		common.OB.Display()
		fmt.Println()

		if defeat := actions.CheckGameOver(common.ShipList); defeat {
			fmt.Println("Jogo finalizado!\nDerrota...")
			break
		}

		if common.IsMyTurn {
			actions.Shoot(common.Conn, reader)
			common.IsMyTurn = false
		}

		fmt.Println("Aguardando oponente...")
		actions.HandleMessage(common.Conn)
	}
}
