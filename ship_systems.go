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
	case actionReverse:
		se.reverse(command.dt)
	}
}

func (se *ShipEngine) Install(ship *Ship) {
	se.ship = ship
	log.Println("engine installed on ship", ship.name)
}

func (se *ShipEngine) Update(info SceneInformation) {
	// No action necessary
}

func (se *ShipEngine) Align(target Entity, dt float64) {
	if target == nil {
		return
	}
	location := se.ship.Coordinates()
	destination := target.Coordinates()
	log.Printf("Want to align to %s at %v from %v", target.Name(), target.Coordinates(), location)
	direction := location.Sub(destination)
	targetAngle := math.Atan2(direction.Y, direction.X)
	se.align(targetAngle, dt)
}

func (se *ShipEngine) align(angle float64, dt float64) {
	// For now: Close enough, don't align more accurately
	const deadZone = 0.05
	// rotate so zero is upwards
	angle = angle + math.Pi/2
	log.Println("trying to align towards", angle, "; current:", se.ship.Angle())
	// should be easy, make the ship's angle go towards param angle
	// get difference
	deltaAngle := normalizeAngle(se.ship.Angle() - angle)
	// TODO: Make better engines align better
	if deltaAngle < -deadZone {
		se.Turn(dt)
	} else if deltaAngle > deadZone {
		se.Turn(-dt)
	}
}

func (se *ShipEngine) reverse(dt float64) {
	// naive reverse, turn around or do nothing

	// ship's angle should go towards the inverse of its velocity
	targetAngle := normalizeAngle(se.ship.Velocity().Angle() + 2*math.Pi)
	se.align(targetAngle, dt)
}

func (se ShipEngine) Turn(dt float64) {
	angle := se.ship.angle + se.TurnSpeed*dt

	se.ship.angle = normalizeAngle(angle)
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

// TODO: Interface!!
// TODO: quality (see colors? see sizes?), interference, range
// TODO: States (hostile, neutral, escort, fighter)
type ShipScanner struct {
	Accuracy       float64
	Range          float64
	targets        []Entity
	celestials     []Celestial
	selectedTarget Entity
	selectedCelest Celestial
}

func (sc ShipScanner) Name() string {
	return "scanner"
}
func (sc *ShipScanner) Activate(command pilotAction) {
	log.Println("scanner activate", command)
	switch command.key {
	case actionTargetPrev:
		sc.PrevTarget()
	case actionTargetNext:
		sc.NextTarget()

	case actionClearTarget:
		sc.ClearTarget()

	case actionLand:
		// The ship wants to land. A good scanner would target the nearest celestial...
		sc.NextCelestial()
	}
}

func (sc ShipScanner) Install(ship *Ship) {
	log.Println("scanner installed on ship", ship.name)
}

func (sc *ShipScanner) Update(info SceneInformation) {
	sc.targets = info.Entities
	sc.celestials = info.Celestials
}

func (sc *ShipScanner) Celestial() Celestial {
	if sc.selectedCelest == nil {
		return nil
	}
	return sc.selectedCelest
}

// TODO: Cycle UNLESS we have no target, then land at the nearest.
// "Clear Target" action should be required to have no target unless the pilot just entered the system.
func (sc *ShipScanner) NextCelestial() {
	if len(sc.celestials) == 0 {
		log.Println("no celestials")
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

func (sc *ShipScanner) Target() Entity {
	return sc.selectedTarget
}

func (sc *ShipScanner) ClearTarget() {
	// clear ship
	if sc.selectedTarget != nil {
		sc.selectedTarget = nil
	} else {
		// no ship? clear celestial
		sc.selectedCelest = nil
	}
}

func (sc *ShipScanner) nextTarget(delta int) Entity {
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

func (sc *ShipScanner) PrevTarget() {
	sc.selectedTarget = sc.nextTarget(-1)
}

func (sc *ShipScanner) NextTarget() {
	sc.selectedTarget = sc.nextTarget(1)
}

func DefaultShipEngine() *ShipEngine {
	return &ShipEngine{
		Acceleration: 3.0,
		MaxVel:       1.0,
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
