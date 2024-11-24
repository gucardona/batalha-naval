package ships

type Ship struct {
	Name         string
	Size         int
	Quantity     int
	Coordinates  []Coordinate
	IsHorizontal bool
}

type Coordinate struct {
	X int
	Y int
}

func NewShip(shipType ShipType) Ship {
	return Ship{
		Name:     shipType.Name,
		Size:     shipType.Size,
		Quantity: shipType.Quatity,
	}
}

func (s *Ship) AddCoordinates(x, y int, isHorizontal bool) error {
	s.IsHorizontal = isHorizontal
	s.Coordinates = make([]Coordinate, s.Size)

	coordinate := Coordinate{X: x, Y: y}

	for i := 0; i < s.Size; i++ {
		s.Coordinates[i] = coordinate
		if isHorizontal {
			coordinate.Y++
		} else {
			coordinate.X++
		}
	}

	return nil
}

func NewPortaAvioes() Ship {
	ship := NewShip(PortaAvioes)
	return ship
}

func NewEncouracado() Ship {
	ship := NewShip(Encouracado)
	return ship
}

func NewCruzador() Ship {
	ship := NewShip(Cruzador)
	return ship
}

func NewDestroier() Ship {
	ship := NewShip(Destroier)
	return ship
}
