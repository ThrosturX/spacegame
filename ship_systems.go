package spacegame

import (
	"encoding/json"
	"log"
	"math"

	"github.com/faiface/pixel"
)

// TODO: Ship systems...
type ShipSystem interface {
	Name() string
	// Activate this system with a command
	Activate(command pilotAction)
	// Install this system on a ship
	Install(ship *Ship)
	Update(info SceneInformation)
}

type ShipSystems map[string]ShipSystem

func (s ShipSystems) MarshalJSON() ([]byte, error) {
	raw := map[string]struct {
		Type   string
		System interface{}
	}{}
	for name, sys := range s {
		typ := sys.Name()
		raw[name] = struct {
			Type   string
			System interface{}
		}{
			Type:   typ,
			System: sys,
		}
	}
	return json.Marshal(raw)
}

func (s ShipSystems) UnmarshalJSON(buf []byte) error {
	raw := map[string]struct {
		Type   string
		System json.RawMessage
	}{}
	err := json.Unmarshal(buf, &raw)
	if err != nil {
		return err
	}
	for name := range s {
		delete(s, name)
	}
	for name, rawsys := range raw {
		var sys ShipSystem
		switch rawsys.Type {
		case "engine":
			sys = &ShipEngine{}
		case "scanner":
			sys = &ShipScanner{}
		}
		err := json.Unmarshal(rawsys.System, sys)
		if err != nil {
			return err
		}
		s[name] = sys
	}
	return nil
}

type ShipEngine struct {
	Acceleration float64
	MaxVel       float64
	TurnSpeed    float64
	ship         *Ship
}

func (se ShipEngine) Name() string {
	return "engine"
}
func (se ShipEngine) Activate(command pilotAction) {
	log.Println("engine activate", command)
	switch command.key {
	case actionAccel:
		se.Accelerate(command.dt)
	case actionTurnLeft:
		se.Turn(command.dt)
	case actionTurnRight:
		se.Turn(-command.dt)
	}
}

func (se *ShipEngine) Install(ship *Ship) {
	se.ship = ship
	log.Println("engine installed on ship", ship.name)
}

func (se *ShipEngine) Update(info SceneInformation) {
	// No action necessary
}

func (se ShipEngine) Align(target Entity, dt float64) {
	if target == nil {
		return
	}
	log.Println("Want to align to", target.Name())
	location := se.ship.Coordinates()
	destination := target.Coordinates()
	direction := location.Add(destination)
	targetAngle := math.Atan2(direction.Y, direction.X)
	targetAngle = targetAngle*180/math.Pi + 270
	deltaAngle := targetAngle - se.ship.Angle()
	if deltaAngle > 1 {
		se.Turn(dt)
	} else if deltaAngle < -1 {
		se.Turn(-dt)
	}
	// Close enough, don't align more accurately
	// TODO: Make better engines align better
}

func (se ShipEngine) Turn(dt float64) {
	angle := se.TurnSpeed * dt
	se.ship.angle += angle
}

func (se ShipEngine) Accelerate(dt float64) {
	thrust := se.Acceleration * dt
	accelVec := pixel.V(0, thrust).Rotated(se.ship.angle)
	vel := se.ship.velocity.Add(accelVec)
	if vel.Len() > se.MaxVel {
		vel = vel.Unit().Scaled(se.MaxVel)
	}
	se.ship.velocity = vel
}

// TODO: quality (see colors? see sizes?), interference
// TODO: States (hostile, neutral, escort, fighter)
type ShipScanner struct {
	Accuracy       float64
	Range          float64
	targets        []Entity
	celestials     []*Celestial
	selectedTarget Entity
	selectedCelest *Celestial
}

func (sc ShipScanner) Name() string {
	return "scanner"
}
func (sc ShipScanner) Activate(command pilotAction) {
	log.Println("scanner activate", command)
	switch command.key {
	case actionTargetPrev:
		sc.PrevTarget()
	case actionTargetNext:
		sc.NextTarget()
	case actionTargetCelestial:
		sc.NextCelestial()
	}
}

func (sc ShipScanner) Install(ship *Ship) {
	log.Println("scanner installed on ship", ship.name)
}

func (sc ShipScanner) Update(info SceneInformation) {
	sc.targets = info.Entities
	sc.celestials = info.Celestials
}

func (sc ShipScanner) Celestial() Entity {
	return sc.selectedCelest
}

func (sc ShipScanner) NextCelestial() {
	if len(sc.celestials) == 0 {
		return
	}
	if sc.selectedCelest == nil {
		sc.selectedCelest = sc.celestials[0]
		return
	}
	for i, _ := range sc.celestials {
		// TODO: Use pointers?
		if sc.selectedCelest.Name() == sc.celestials[i].Name() {
			if i+1 < len(sc.celestials) {
				sc.selectedCelest = sc.celestials[i+1]
				return
			}
		}
	}
	// Cycle completed
	sc.selectedCelest = nil
}

func (sc ShipScanner) Target() Entity {
	return sc.selectedTarget
}

func (sc ShipScanner) nextTarget(delta int) Entity {
	// check if we have something targetted first
	if sc.selectedTarget == nil {
		// target the 'first' thing
		if len(sc.targets) > 0 {
			// decide direction
			if delta == 1 {
				return sc.targets[0]
			} else if delta == -1 {
				return sc.targets[len(sc.targets)-1]
			}
		}
		// no targets
		return nil
	}

	for i := range sc.targets {
		if sc.selectedTarget == sc.targets[i] {
			// check if there exists a next target
			if i+delta < len(sc.targets) {
				return sc.targets[i+delta]
			}
			// full circle
			return nil
		}
	}

	// nothing found
	return nil

}

func (sc ShipScanner) PrevTarget() {
	sc.selectedTarget = sc.nextTarget(-1)
}

func (sc ShipScanner) NextTarget() {
	sc.selectedTarget = sc.nextTarget(1)
}

func DefaultShipEngine() *ShipEngine {
	return &ShipEngine{
		Acceleration: 3.0,
		MaxVel:       200.0,
		TurnSpeed:    3.0,
	}
}

func DefaultShipScanner() *ShipScanner {
	return &ShipScanner{
		Accuracy: 100.0,
		Range:    300000.0,
	}
}

func DefaultShipSystems() map[string]ShipSystem {
	sysmap := make(map[string]ShipSystem)

	sysmap["engine"] = DefaultShipEngine()
	sysmap["scanner"] = DefaultShipScanner()

	return sysmap
}
