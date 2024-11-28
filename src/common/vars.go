package common

import (
	"awesomeProject/src/battlefield"
	"awesomeProject/src/ships"
	"net"
)

var (
	ShipList []*ships.Ship

	JsonShips         = make([]*ShipJSON, 0)
	OpponentJsonShips = make([]*ShipJSON, 0)

	B  = battlefield.NewBattlefield()
	OB = battlefield.NewBattlefield()

	Conn   net.Conn
	Buffer = make([]byte, 8192)

	IsMyTurn = true
	IsServer = false
)
