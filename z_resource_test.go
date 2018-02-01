package spacegame

import (
	"reflect"
	"testing"
)

func TestShipSaveToFile(t *testing.T) {
	// load the "Starbridge" ship
	ship1, err := LoadShip("resources/entities/ships/Starbridge.json")
	if err != nil {
		panic(err)
	}

	// save it to file
	tmpPath := "/tmp/ship1.json"
	ship1.SaveToFile(tmpPath)

	// load from saved file
	ship2, err := LoadShip(tmpPath)
	if err != nil {
		panic(err)
	}

	if !reflect.DeepEqual(ship1, ship2) {
		t.Fail()
	}
}

func TestSystemSaveToFile(t *testing.T) {
	// load the "Vera" system
	sys1, err := LoadSystem("resources/universe/systems/Vera.json")
	if err != nil {
		panic(err)
	}

	// save it to file
	tmpPath := "/tmp/sys1.json"
	sys1.SaveToFile(tmpPath)

	// load from saved file
	sys2, err := LoadSystem(tmpPath)
	if err != nil {
		panic(err)
	}

	if !reflect.DeepEqual(sys1, sys2) {
		t.Fail()
	}
}
