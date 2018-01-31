package spacegame

// TODO
type Controllable interface {
	Vessel() Entity
}

type Player struct {
	name    string
	ship    *Ship // change to Entity for fun?
	credits uint64
}

// Creates a new player with the specified name
// The player is assigned a Starbridge and 10000 credits
func NewPlayer(name string) *Player {
	return &Player{
		name:    name,
		ship:    NewShip("Starbridge", "ships/starbridge.png"),
		credits: 10000,
	}
}

func (p Player) Name() string {
	return p.name
}

func (p Player) Ship() *Ship {
	return p.ship
}
