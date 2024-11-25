package server

import (
	"awesomeProject/src/battlefield"
	"awesomeProject/src/client"
	"awesomeProject/src/common"
	"awesomeProject/src/ships"
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
)

var (
	opponentJsonShips []common.ShipJSON
	targetIP          string
	targetPort        int
	b                 *battlefield.Battlefield
	c                 *client.Client
	conn              net.Conn
	reader            *bufio.Reader
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
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Falha ao aceitar conexão:", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	fmt.Println("Conexão estabelecida com:", conn.RemoteAddr())

	go func() {
		for {
			buffer := make([]byte, 1024)
			n, err := conn.Read(buffer)
			if err != nil {
				fmt.Println("Erro ao ler dados do buffer:", err)
				return
			}

			if err = json.Unmarshal(buffer[:n], &opponentJsonShips); err != nil {
				if !(string(buffer[:n]) == "\n") {
					s.handleShot(buffer[:n])
				}
				continue
			}

			fmt.Printf("Recebi: %s\n", string(buffer[:n]))
		}
	}()
}

func (s *Server) handleShot(message []byte) {
	coord := strings.TrimSpace(string(message))

	validCoordRegex := regexp.MustCompile(`^\d{2}$`)
	if !validCoordRegex.MatchString(coord) {
		return
	}

	x := int(coord[0] - '0')
	y := int(coord[1] - '0')

	common.ClearScreen()

	if b.Grid[x][y] == 1 {
		b.Grid[x][y] = 9
		fmt.Printf("Você foi atingido em: [%d, %d]\n", x, y)
	} else {
		fmt.Printf("Tiro do adversário na água em: [%d, %d]\n", x, y)
	}

	shoot(conn, reader)
}

func (s *Server) Handle() {
	reader = bufio.NewReader(os.Stdin)
	b = battlefield.NewBattlefield()

	for {
		if conn == nil {
			fmt.Println("Por favor, informe o IP e porta do servidor para se conectar (exemplo: 127.0.0.1:8080)")

			targetServer := common.ReadAndSplitCommand(reader, ":")

			var err error
			targetPort, err = common.ConvertPort(targetServer[1])
			if err != nil {
				continue
			}

			targetIP = targetServer[0]

			c = client.NewClient(targetIP, targetPort)
			conn, err = c.Connect()
			if err != nil {
				common.ClearScreen()
				fmt.Println("Erro ao conectar ao servidor:", err)
				fmt.Println("Por favor, tente novamente.")
				continue
			}

			break
		}
	}
	defer conn.Close()

	//ob := battlefield.NewBattlefield()

	var (
		shipList  []ships.Ship
		jsonShips []common.ShipJSON
	)

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

	var generatedJson string
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

		generatedJson, _ = common.GenerateJSON(shipList)
		_ = common.WriteJSONToFile(generatedJson)
	} else {
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
	}

	for {
		shoot(conn, reader)
	}
}

func shoot(conn net.Conn, reader *bufio.Reader) {
	fmt.Println("Seu campo de batalha:")
	b.Display()
	fmt.Println()
	fmt.Println("Envie a coordenada do seu tiro no formato 'linha coluna' (ex: 11).")
	chosenTargetCoord, _ := reader.ReadString('\n')
	chosenTargetCoord = strings.TrimSpace(chosenTargetCoord)
	validCoordRegex := regexp.MustCompile(`^\d{2}$`)

	if !validCoordRegex.MatchString(chosenTargetCoord) {
		common.ClearScreen()
		fmt.Println("Coordenada inválida. Por favor, tente novamente.")
	}

	if err := c.SendMessage(conn, chosenTargetCoord); err != nil {
		common.ClearScreen()
		fmt.Println("Erro ao enviar mensagem:", err)
		fmt.Println("Por favor, tente novamente.")
	}
}
