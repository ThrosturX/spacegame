package spacegame

import "fmt"

// TODO
type Controllable interface {
	CmdChan() chan action
	Vessel() Entity
	Process(action)
}

type Player struct {
	name    string
	ship    *Ship // change to Entity for fun?
	credits uint64
	cmdChan chan action
}

// Creates a new player with the specified name
// The player is assigned a Starbridge and 10000 credits
func NewPlayer(name string, resourceManager ResourceManager) *Player {
	ship := NewShip("Starbridge")

	player := &Player{
		name:    name,
		ship:    ship,
		credits: 10000,
		cmdChan: make(chan action),
	}

	return player
}

func (p *Player) Name() string {
	return p.name
}

func (p *Player) Ship() *Ship {
	return p.ship
}

func (p *Player) CmdChan() chan action {
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

func (p *Player) Process(a action) {
	fmt.Printf("Processing action %v\n", a)
	switch a.key {
	case actionAccel:
		p.ship.thrusters(a.dt)
	case actionTurnLeft:
		p.ship.turn(a.dt)
	case actionTurnRight:
		p.ship.turn(-a.dt)
	}
}
