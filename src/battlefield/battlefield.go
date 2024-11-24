package battlefield

import (
	"awesomeProject/src/ships"
	"fmt"
)

const Size = 10

type Battlefield struct {
	Grid [Size][Size]int
}

func NewBattlefield() *Battlefield {
	return &Battlefield{}
}

func (b *Battlefield) Display() {
	fmt.Print("     ")
	for y := 0; y < Size; y++ {
		fmt.Printf("%d ", y)
	}

	fmt.Print("\n   +")
	for i := 0; i < 21; i++ {
		fmt.Print("-")
	}
	fmt.Print("+\n")

	for x := 0; x < Size; x++ {
		fmt.Printf("%2d | ", x)
		for y := 0; y < Size; y++ {
			fmt.Printf("%d ", b.Grid[x][y])
		}
		fmt.Println("|")
	}

	fmt.Print("   +")
	for i := 0; i < 21; i++ {
		fmt.Print("-")
	}
	fmt.Print("+")
}

func (b *Battlefield) PlaceShip(ship ships.Ship) error {
	for _, coord := range ship.Coordinates {
		if coord.X >= Size || coord.Y >= Size {
			return fmt.Errorf("coordenada %v fora da grade", coord)
		}

		if b.Grid[coord.X][coord.Y] != 0 {
			return fmt.Errorf("coordenada %v já está ocupada", coord)
		}
	}

	for _, coord := range ship.Coordinates {
		b.Grid[coord.X][coord.Y] = 1
	}

	return nil
}
