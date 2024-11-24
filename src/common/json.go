package common

import (
	"awesomeProject/src/ships"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type ShipJSON struct {
	Tipo     string   `json:"tipo"`
	Posicoes [][2]int `json:"posicoes"`
}

func GenerateJSON(ships []ships.Ship) (string, error) {
	var jsonShips []ShipJSON

	for _, ship := range ships {
		var positions [][2]int
		for _, coord := range ship.Coordinates {
			positions = append(positions, [2]int{coord.Y, coord.X})
		}

		jsonShip := ShipJSON{
			Tipo:     ship.Name,
			Posicoes: positions,
		}
		jsonShips = append(jsonShips, jsonShip)
	}

	jsonData, err := json.Marshal(jsonShips)
	if err != nil {
		return "", fmt.Errorf("failed to generate JSON: %w", err)
	}

	return string(jsonData), nil
}

func WriteJSONToFile(data string) error {
	fileName := "ships.json"
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = file.Write([]byte(data))
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}

func ConvertJsonToShip(jsonShips []ShipJSON) ([]ships.Ship, error) {
	var shipList []ships.Ship

	for _, jsonShip := range jsonShips {
		var ship ships.Ship

		switch strings.ToLower(jsonShip.Tipo) {
		case "porta-avioes":
			ship = ships.NewPortaAvioes()
		case "encouracado":
			ship = ships.NewEncouracado()
		case "cruzador":
			ship = ships.NewCruzador()
		case "destroier":
			ship = ships.NewDestroier()
		default:
			return nil, fmt.Errorf("tipo de navio inválido: %s", jsonShip.Tipo)
		}

		previousCoord := [2]int{999, 999}
		for _, coord := range jsonShip.Posicoes {
			if previousCoord[0] == 999 && previousCoord[1] == 999 {
				previousCoord = coord
			} else {
				if previousCoord[0] == coord[0] {
					ship.IsHorizontal = true
				}
			}

			ship.Coordinates = append(ship.Coordinates, ships.Coordinate{X: coord[0], Y: coord[1]})
		}

		shipList = append(shipList, ship)
	}

	return shipList, nil
}
