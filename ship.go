package spacegame

import (
	"encoding/json"
	"log"
	"os"

	"github.com/faiface/pixel"
)

type Ship struct {
	angle       float64
	coordinates pixel.Vec
	velocity    pixel.Vec
	bounds      pixel.Rect
	cmdChan     chan action
	ShipConfig
}

type ShipConfig struct {
	Name         string
	Acceleration float64
	MaxVel       float64
	TurnSpeed    float64
	Length       float64
	Width        float64
}

func DefaultShipConfig(name string) ShipConfig {
	return ShipConfig{
		Name:         name,
		Acceleration: 3.0,
		MaxVel:       200.0,
		TurnSpeed:    3.0,
		Length:       32,
		Width:        32,
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

func (s *Ship) CmdChan() chan action {
	return s.cmdChan
}

func (s *Ship) Process(a action) {
	switch a.key {
	case actionAccel:
		s.thrusters(a.dt)
	case actionTurnLeft:
		s.turn(a.dt)
	case actionTurnRight:
		s.turn(-a.dt)
	}
}
