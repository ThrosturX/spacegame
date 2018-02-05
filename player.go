package spacegame

import "fmt"

type Player struct {
	name    string
	ship    PilotableShip
	credits uint64
	cmdChan chan pilotAction
}

// Creates a new player with the specified name
// The player is assigned a Starbridge and 10000 credits
func NewPlayer(name string, resourceManager ResourceManager) *Player {
	ship := NewShip("Starbridge")

	player := &Player{
		name:    name,
		ship:    ship,
		credits: 10000,
		cmdChan: make(chan pilotAction),
	}

	return player
}

func (p *Player) Name() string {
	return p.name
}

func (p *Player) Ship() *Ship {
    var (
        ship *Ship
        ok bool
    )
    if ship, ok = p.ship.(*Ship); !ok {
        panic("Wrong type for player ship")
        return nil
    }
    return ship
}

func (p *Player) CmdChan() chan pilotAction {
	return p.cmdChan
}

func (p *Player) Vessel() Entity {
	return p.Ship()
}

func (p *Player) tick() {
	for cmd := range p.cmdChan {
		p.Process(cmd)
	}
}

func (p *Player) Process(a pilotAction) {
	switch a.key {
	case actionAccel:
        fallthrough
	case actionTurnLeft:
        fallthrough
	case actionTurnRight:
        fallthrough
    case actionTargetNext:
        fmt.Printf("Processing action %v\n", a)
        p.ship.Process(a)
	}
}
