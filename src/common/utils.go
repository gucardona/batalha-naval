package common

import (
	"awesomeProject/src/battlefield"
	"awesomeProject/src/ships"
	"bufio"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func PrintPlacementText(b *battlefield.Battlefield, i int) {
	fmt.Println("Seu campo de batalha:")
	fmt.Println()
	b.Display()
	fmt.Println()
	fmt.Printf("\nVocê tem %d navios para posicionar:\n", i)
	fmt.Println("  - (p) 1 porta-aviões (5 posições)")
	fmt.Println("  - (e) 1 encouraçado (4 posições)")
	fmt.Println("  - (c) 2 cruzadores (3 posições)")
	fmt.Println("  - (d) 2 destróieres (2 posições)")
	fmt.Println("\nDigite as coordenadas de cada navio no formato:\n'<tipo de navio> <linha> <coluna> <horizontal ou vertical>' (ex: p 1 1 h).")
	fmt.Println("  - tipos: p, e, c, d")
	fmt.Println("  - linha: 0-9")
	fmt.Println("  - coluna: 0-9")
	fmt.Println("  - horizontal ou vertical: h, v")
	fmt.Print("\nDigite o comando: ")
}

func ReadAndSplitCommand(reader *bufio.Reader, sep string) []string {
	command, _ := reader.ReadString('\n')
	command = strings.TrimSpace(command)
	commands := strings.Split(command, sep)
	return commands
}

func CreateShipByCommand(command string) (ships.Ship, error) {
	switch command {
	case "p", "P":
		return ships.NewPortaAvioes(), nil
	case "e", "E":
		return ships.NewEncouracado(), nil
	case "c", "C":
		return ships.NewCruzador(), nil
	case "d", "D":
		return ships.NewDestroier(), nil
	default:
		ClearScreen()
		fmt.Println("Tipo de navio inválido. Por favor, tente novamente.")
		return ships.Ship{}, errors.New("tipo de navio inválido")
	}
}

func PlacedAllShips(shipMap map[string]int, ship ships.Ship) bool {
	if shipMap[ship.Name] >= ship.Quantity {
		ClearScreen()
		fmt.Println("Você já posicionou todos os navios do tipo", ship.Name)
		fmt.Println("Por favor, posicione outro navio.")
		return true
	}
	return false
}

func ConvertCoordinates(commands []string) (int, int, error) {
	line, err := strconv.Atoi(commands[1])
	if err != nil {
		ClearScreen()
		fmt.Println("Erro ao converter coordenada X:", err)
		return 0, 0, err
	}

	column, err := strconv.Atoi(commands[2])
	if err != nil {
		ClearScreen()
		fmt.Println("Erro ao converter coordenada Y:", err)
		return 0, 0, err
	}

	return line, column, nil
}

func AddCoordinatesAndPlaceShip(b *battlefield.Battlefield, ship ships.Ship, line int, column int, commandHorizontal string) error {
	if err := ship.AddCoordinates(line, column, commandHorizontal == "h" || commandHorizontal == "H"); err != nil {
		ClearScreen()
		fmt.Println("Erro ao adicionar coordenada do navio:", err)
		fmt.Println("Por favor, posicione o navio novamente.")
		return err
	}

	if err := b.PlaceShip(ship); err != nil {
		ClearScreen()
		fmt.Println("Erro ao posicionar navio:", err)
		fmt.Println("Por favor, posicione o navio novamente.")
		return err
	}

	return nil
}

func IncrementShipMap(shipMap map[string]int, ship ships.Ship) {
	if _, exists := shipMap[ship.Name]; exists {
		shipMap[ship.Name]++
	} else {
		shipMap[ship.Name] = 1
	}
}

func ConvertPort(port string) (int, error) {
	convertedPort, err := strconv.Atoi(port)
	if err != nil {
		ClearScreen()
		fmt.Println("Erro ao converter porta:", err)
		fmt.Println("Por favor, informe o ip e a porta novamente.")
		return 0, err
	}
	return convertedPort, nil
}

func ClearScreen() {
	fmt.Print("\033[H\033[2J")
}
