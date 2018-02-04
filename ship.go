package spacegame

import (
	"encoding/json"
	"log"
	"os"

	"github.com/faiface/pixel"
)

// Defines interaction between pilot (player or AI) and ship
type PilotableShip interface {
	Process(pilotAction)                               // process input events
	Config() ShipConfig                                // TODO: directly accessible methods to what is needed...
	ActivateSystem(system string, command pilotAction) // TODO: Maybe return something?
	CmdChan() chan pilotAction
}

// TODO: Ship systems...
type ShipSystem interface {
	Name() string
	Activate(command pilotAction)
}

type Ship struct {
	angle       float64
	coordinates pixel.Vec
	velocity    pixel.Vec
	bounds      pixel.Rect
	cmdChan     chan pilotAction
	ShipConfig
	systems map[string]ShipSystem
}

type ShipConfig struct {
	Name        string
	Length      float64
	Width       float64
	ShipScanner Scanner
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

	switch command.key {
	case actionAccel:
		se.Accelerate(command.dt)
	case actionTurnLeft:
		se.Turn(command.dt)
	case actionTurnRight:
		se.Turn(-command.dt)
	}
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

// TODO: Put in some systems ?
func DefaultShipConfig(name string) ShipConfig {
	return ShipConfig{
		Name: name,
		////	Acceleration: 3.0,
		////	MaxVel:       200.0,
		////	TurnSpeed:    3.0,
		////	Length:       32,
		////	Width:        32,
	}
}

func NewShip(name string) *Ship {
	return &Ship{
		angle:       0.0,
		coordinates: pixel.V(0, 0),
		bounds:      pixel.R(0, 0, 32, 32),
		ShipConfig:  DefaultShipConfig(name),
	}
}

func LoadShip(path string) (*Ship, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	var shipConfig ShipConfig
	decoder := json.NewDecoder(file)
	decoder.Decode(&shipConfig)
	return shipConfig.load(), nil
}

func (config ShipConfig) load() *Ship {
	ship := &Ship{
		angle:       0.0,
		coordinates: pixel.ZV,
		bounds:      pixel.R(0, 0, config.Width, config.Length),
		ShipConfig:  config,
	}
	return ship
}

func (s *Ship) SaveToFile(path string) error {
	// TODO: Check if file exists?

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	encoder.Encode(s.ShipConfig)
	return nil
}

func (s *Ship) Config() ShipConfig {
	return s.ShipConfig
}

func (s *Ship) Name() string {
	return s.ShipConfig.Name
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
		// No such system
		return
	}

	sys.Activate(command)
}

func (s *Ship) thrusters(dt float64) {
	thrust := s.ShipConfig.Acceleration * dt
	accelVec := pixel.V(0, thrust).Rotated(s.angle)
	vel := s.velocity.Add(accelVec)
	if vel.Len() > s.ShipConfig.MaxVel {
		vel = vel.Unit().Scaled(s.ShipConfig.MaxVel)
	}
	s.velocity = vel
	log.Println(s.velocity)
}

func (s *Ship) turn(dt float64) {
	angle := s.ShipConfig.TurnSpeed * dt
	s.angle += angle
}

func (s *Ship) CmdChan() chan pilotAction {
	return s.cmdChan
}

func (s *Ship) Process(a pilotAction) {
	switch a {
	case actionAccel:
	case actionTurnLeft:
	case actionTurnRight:
		s.ActivateSystem("engine", a)
	case actionTargetNext:
		p.ship.ActivateSystem("scanner", a)
	}

}
