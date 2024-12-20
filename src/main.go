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
	"time"
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

		_, err = os.Stat(fmt.Sprintf("%s/my-ships.json", rootDir))
		if err == nil {
			fmt.Println("Arquivo \"my-ships.json\" não deletado encontrado. Deseja ler navios do arquivo? (s/n)")
			read, _ := reader.ReadString('\n')
			read = strings.TrimSpace(read)

			if read == "s" || read == "S" {
				file, err := os.ReadFile("my-ships.json")
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

				opFile, err := os.ReadFile("op-ships.json")
				if err != nil {
					common.ClearScreen()
					fmt.Printf("falha ao ler arquivo: %s\nPor favor, insira novamente os navios...", err)
					break
				}

				if err = json.Unmarshal(opFile, &common.JsonShips); err != nil {
					common.ClearScreen()
					fmt.Printf("falha ao fazer Unmarshal do json: %s\nPor favor, insira novamente os navios...", err)
					break
				}

				common.IsReconnection = true
			}
		}

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
	common.SaveMyShips()

	if !common.IsServer {
		common.IsMyTurn = false
		defer common.Conn.Close()
	}

	for {
		common.PrintBattlefields()

		go func() {
			if defeat := actions.CheckDefeat(); defeat {
				fmt.Println("Jogo finalizado!\nDerrota...")
				time.Sleep(2 * time.Second)
				os.Exit(0)
			}
		}()

		if common.IsMyTurn {
			actions.Shoot(common.Conn, reader)
			common.IsMyTurn = false
		}

		fmt.Println("Aguardando oponente...")

		actions.HandleMessage(common.Conn)
	}
}
