package spacegame

import (
	"encoding/json"
	"log"
	"os"

	"github.com/faiface/pixel"
)

// Defines interaction between pilot (player or AI) and ship
type PilotableShip interface {
	Process(pilotAction) // process input events
	Serialize() SerializableShip
	ActivateSystem(system string, command pilotAction) // TODO: Maybe return something?
}

// TODO: Ship systems...
type ShipSystem interface {
	Name() string
    // Activate this system with a command
	Activate(command pilotAction)
    // Install this system on a ship
    Install(ship *Ship)
}

type Ship struct {
	name        string
	angle       float64
	coordinates pixel.Vec
	velocity    pixel.Vec
	bounds      pixel.Rect
	systems     map[string]ShipSystem
}

type ShipSystems map[string]ShipSystem

type SerializableShip struct {
	Name    string
	Length  float64
	Width   float64
	Systems ShipSystems
}

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

func DefaultShipEngine() *ShipEngine {
	return &ShipEngine{
		Acceleration: 3.0,
		MaxVel:       200.0,
		TurnSpeed:    3.0,
	}
}

// TODO: Put in some systems ?
func DefaultShipConfig(name string) SerializableShip {
	return SerializableShip{
		Name:    name,
		Length:  32,
		Width:   32,
		Systems: DefaultShipSystems(),
	}
}

func DefaultShipSystems() map[string]ShipSystem {
	sysmap := make(map[string]ShipSystem)

	sysmap["engine"] = DefaultShipEngine()

	return sysmap
}

func NewShip(name string) *Ship {
    ship := &Ship{
		name:        name,
		angle:       0.0,
		coordinates: pixel.V(0, 0),
		velocity:    pixel.V(0, 0),
		bounds:      pixel.R(0, 0, 32, 32),
		systems:     DefaultShipSystems(),
	}
    for _, system := range ship.systems {
        system.Install(ship)
    }
    return ship
}

func LoadShip(path string) (*Ship, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	var shipConfig SerializableShip
	shipConfig.Systems = ShipSystems{}
	decoder := json.NewDecoder(file)
	decoder.Decode(&shipConfig)

//	log.Println(shipConfig.Systems)

	return shipConfig.load(), nil
}

func (config SerializableShip) load() *Ship {
	ship := &Ship{
		name:        config.Name,
		angle:       0.0,
		coordinates: pixel.ZV,
		bounds:      pixel.R(0, 0, config.Width, config.Length),
		systems:     config.Systems,
	}

	for _, sys := range ship.systems {
		sys.Install(ship)
	}
	return ship
}

func (ss SerializableShip) SaveToFile(path string) error {
	// TODO: Check if file exists?

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	//    var b bytes.Buffer

//	log.Println(ss.Systems)

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	encoder.Encode(ss)
	return nil
}

func (s *Ship) SaveToFile(path string) error {
	// Create the Serializableship

	ss := s.Serialize()

	return ss.SaveToFile(path)
}

func (s *Ship) Serialize() SerializableShip {
//	log.Println(s.name, s.systems)
	ss := SerializableShip{
		Name:    s.name,
		Length:  s.bounds.Norm().H(),
		Width:   s.bounds.Norm().W(),
		Systems: s.systems,
	}

	return ss
}

func (s *Ship) Name() string {
	return s.name
}

func (s *Ship) Angle() float64 {
	return s.angle
}

func (s *Ship) Bounds() pixel.Rect {
	return s.bounds
}

func (s *Ship) Coordinates() pixel.Vec {
	return s.coordinates
}

func (s *Ship) Velocity() pixel.Vec {
	return s.velocity
}

func (s *Ship) Translate(by pixel.Vec) {
	s.coordinates = s.coordinates.Add(by)
}

/// TODO: Systems:
func (s *Ship) ActivateSystem(system string, command pilotAction) {
	// TODO:

	sys, ok := s.systems[system]
	if !ok {
		log.Printf("No system <%s> exists to process command <%s>\n", system, command.key)
		// No such system
		return
	}

	sys.Activate(command)
}

func (s *Ship) Process(a pilotAction) {
	switch a.key {
	case actionAccel:
        fallthrough
	case actionTurnLeft:
        fallthrough
	case actionTurnRight:
		s.ActivateSystem("engine", a)
	case actionTargetNext:
		s.ActivateSystem("scanner", a)
	}

}

// Custom marshaling

// func (ss SerializableShip) MarshalJson() ([]byte, error) {
//     return nil, nil
// }

// func (ss ShipSystems) MarshalJson() ([]byte, error) {
// 	buffer := bytes.NewBufferString("{\n\"Systems\"")
// 	length := len(ss)
// 	count := 0
// 	for key, value := range ss {
// 		jsonValue, err := json.Marshal(value)
// 		if err != nil {
// 			return nil, err
// 		}
// 		buffer.WriteString(fmt.Sprintf("\"%s\":%s", key, string(jsonValue)))
// 		count++
// 		if count < length {
// 			buffer.WriteString(",")
// 		}
// 	}
// 	buffer.WriteString("}\n}")
// 	return buffer.Bytes(), nil
// }
