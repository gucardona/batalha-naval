package common

import (
	"awesomeProject/src/battlefield"
	"awesomeProject/src/ships"
	"bufio"
	"math/rand"
	"time"
)

func ManuallyPlaceShips(reader *bufio.Reader, b *battlefield.Battlefield) []ships.Ship {
	shipMap := make(map[string]int)
	var shipList []ships.Ship

	i := 6
	for i > 0 {
		PrintPlacementText(b, i)
		commands := ReadAndSplitCommand(reader, " ")

		ship, err := CreateShipByCommand(commands[0])
		if err != nil {
			continue
		}

		if PlacedAllShips(shipMap, ship) {
			continue
		}

		line, column, err := ConvertCoordinates(commands)
		if err != nil {
			continue
		}

		if err = AddCoordinatesAndPlaceShip(b, ship, line, column, commands[3]); err != nil {
			continue
		}

		IncrementShipMap(shipMap, ship)

		i--
		shipList = append(shipList, ship)
		ClearScreen()
	}

	return shipList
}

func AutomaticallyPlaceShips(b *battlefield.Battlefield) []ships.Ship {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	var shipList []ships.Ship

	for _, shipType := range ships.ShipTypes {
		for i := 0; i < shipType.Quatity; i++ {
			placed := false
			for !placed {
				x := random.Intn(battlefield.Size - 1)
				y := random.Intn(battlefield.Size - 1)
				isHorizontal := random.Intn(2) == 0

				ship := ships.NewShip(shipType)

				if err := ship.AddCoordinates(x, y, isHorizontal); err != nil {
					continue
				}

				if err := b.PlaceShip(ship); err != nil {
					continue
				}

				shipList = append(shipList, ship)
				placed = true
			}
		}
	}

	return shipList
}
