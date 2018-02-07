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
	Update(info SceneInformation)
}

type Ship struct {
	name        string
	angle       float64
	coordinates pixel.Vec
	velocity    pixel.Vec
	bounds      pixel.Rect
	systems     map[string]ShipSystem
}

type SerializableShip struct {
	Name    string
	Length  float64
	Width   float64
	Systems ShipSystems
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

// TODO: Refactor
func (s *Ship) Process(a pilotAction) {
	switch a.key {
	case actionAccel:
		fallthrough
	case actionTurnLeft:
		fallthrough
	case actionTurnRight:
		s.ActivateSystem("engine", a)

	case actionTargetCelestial:
		fallthrough
	case actionTargetNext:
		s.ActivateSystem("scanner", a)

	case actionAlign:
		var (
			scanner *ShipScanner
			engine  *ShipEngine
			target  Entity
			ok      bool
		)
		scanner, ok = s.systems["scanner"].(*ShipScanner)
		if !ok {
			return
		}
		engine, ok = s.systems["engine"].(*ShipEngine)
		if !ok {
			return
		}
		target = scanner.Target()
		log.Println("TARGET", target)
		// no target, face celestial
		if target == nil {
			target = scanner.Celestial()
			log.Println("TARGET", target, "target was nil")
		}
		if target == (*Celestial)(nil){
			log.Println("NO TARGET")
			return // nothing to align to // TODO: Maybe align to sun/origin?
		}
		log.Println("TARGET", target)
		engine.Align(target, a.dt)
		log.Println("Found target", (target).Name(), "at", (target).Coordinates())
	}
}

func (s *Ship) Update(info SceneInformation) {
	for _, sys := range s.systems {
		sys.Update(info)
	}
}
