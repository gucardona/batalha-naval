package ships

type ShipType struct {
	Name    string
	Size    int
	Quatity int
}

var (
	PortaAvioes = ShipType{Name: "porta-avioes", Size: 5, Quatity: 1}
	Encouracado = ShipType{Name: "encouracado", Size: 4, Quatity: 1}
	Cruzador    = ShipType{Name: "cruzador", Size: 3, Quatity: 2}
	Destroier   = ShipType{Name: "destroier", Size: 2, Quatity: 2}
)

var ShipTypes = []ShipType{
	PortaAvioes,
	Encouracado,
	Cruzador,
	Destroier,
}
