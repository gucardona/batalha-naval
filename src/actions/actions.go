package actions

import (
	"awesomeProject/src/common"
	"awesomeProject/src/ships"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"regexp"
	"strings"
)

func Shoot(conn net.Conn, reader *bufio.Reader) {
	fmt.Println("Envie a coordenada do seu tiro no formato 'linha coluna' (ex: 11).")
	chosenTargetCoord, _ := reader.ReadString('\n')
	chosenTargetCoord = strings.TrimSpace(chosenTargetCoord)
	validCoordRegex := regexp.MustCompile(`^\d{2}$`)

	if !validCoordRegex.MatchString(chosenTargetCoord) {
		common.ClearScreen()
		fmt.Println("Coordenada inválida. Por favor, tente novamente.")
	}

	if _, err := conn.Write([]byte(chosenTargetCoord)); err != nil {
		common.ClearScreen()
		fmt.Println("Erro ao enviar mensagem:", err)
		fmt.Println("Por favor, tente novamente.")
	}
}

func HandleShot(message []byte) {
	coord := strings.TrimSpace(string(message))

	validCoordRegex := regexp.MustCompile(`^\d{2}$`)
	if !validCoordRegex.MatchString(coord) {
		return
	}

	x := int(coord[0] - '0')
	y := int(coord[1] - '0')

	common.ClearScreen()

	if common.B.Grid[x][y] == 1 {
		fmt.Printf("Você foi atingido em: [%d, %d]\n", x, y)
		common.B.Grid[x][y] = 9
		for _, ship := range common.ShipList {
			if isHit(ship, x, y) {
				if isDestroyed(ship) {
					fmt.Printf("O navio %s foi destruído!\n", ship.Name)
				}
				fmt.Println(*ship)
				generatedJson, _ := common.GenerateJSON(common.ShipList)
				_ = common.WriteJSONToFile(generatedJson)
			}
		}
	} else {
		fmt.Printf("Tiro do adversário na água em: [%d, %d]\n", x, y)
	}

	common.IsMyTurn = true
}

func CheckGameOver(list []*ships.Ship) bool {
	for _, ship := range list {
		if !isDestroyed(ship) {
			return false
		}
	}
	return true
}

func isDestroyed(ship *ships.Ship) bool {
	for _, coord := range ship.Coordinates {
		if common.B.Grid[coord.X][coord.Y] != 9 {
			return false
		}
	}

	return true
}

func isHit(ship *ships.Ship, x int, y int) bool {
	for i := range ship.Coordinates {
		if ship.Coordinates[i].X == x && ship.Coordinates[i].Y == y {
			ship.Coordinates[i].X = -x
			if ship.Coordinates[i].X == 0 {
				ship.Coordinates[i].X = -9090
			}
			ship.Coordinates[i].Y = -y
			if ship.Coordinates[i].Y == 0 {
				ship.Coordinates[i].Y = -9090
			}
			fmt.Println(ship.Coordinates[i].X)
			fmt.Println(ship.Coordinates[i].Y)
			return true
		}
	}

	return false
}

func HandleMessage(conn net.Conn) {
	n, err := conn.Read(common.Buffer)
	if err != nil {
		if err == io.EOF {
			fmt.Println("Connection closed by EOF")
			os.Exit(0)
		}
		fmt.Println("Erro ao ler dados do buffer:", err)
		return
	}

	if err = json.Unmarshal(common.Buffer[:n], &common.OpponentJsonShips); err != nil {
		if !(string(common.Buffer[:n]) == "\n") {
			HandleShot(common.Buffer[:n])
		}
	} else {
		opponentShips, err := common.ConvertJsonToShip(common.OpponentJsonShips)
		if err != nil {
			fmt.Println("Erro ao converter JSON de navios do oponente:", err)
		}

		for _, ship := range opponentShips {
			if err = common.OB.PlaceShip(ship); err != nil {
				fmt.Println("Erro ao colocar navio no campo de batalha do oponente:", err)
			}
		}
	}
}

func SendJSONShips() {
	generatedJson, _ := common.GenerateJSON(common.ShipList)

	if _, err := common.Conn.Write([]byte(generatedJson + "\n")); err != nil {
		fmt.Println("failed to send json with ships:", err)
	}
}

func ReceiveOpponentShips() {
	n, err := common.Conn.Read(common.Buffer)
	if err != nil {
		fmt.Println("Erro ao ler dados do buffer:", err)
		return
	}

	if err = json.Unmarshal(common.Buffer[:n], &common.OpponentJsonShips); err != nil {
		fmt.Println("Erro ao fazer Unmarshal do JSON:", err)
		return
	}

	opponentShips, err := common.ConvertJsonToShip(common.OpponentJsonShips)
	if err != nil {
		fmt.Println("Erro ao converter JSON de navios do oponente:", err)
	}

	for _, ship := range opponentShips {
		if err = common.OB.PlaceShip(ship); err != nil {
			fmt.Println("Erro ao colocar navio no campo de batalha do oponente:", err)
		}
	}
}

func SendAndReceiveShips() {
	SendJSONShips()
	ReceiveOpponentShips()
}

func ReceiveAndSendShips() {
	ReceiveOpponentShips()
	SendJSONShips()
}
