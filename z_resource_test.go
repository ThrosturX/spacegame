package spacegame

import (
	"encoding/json"
	"log"
	"reflect"
	"testing"
)

// Basic serialization test
func TestSerialization(t *testing.T) {
	// Create a SerializableShip
	ss1 := SerializableShip{
		Name:    "Starbridge",
		Length:  32,
		Width:   32,
		Systems: DefaultShipSystems(),
	}

	// Marshal it
	b, err := json.Marshal(ss1)
	if err != nil {
		panic(err)
	}

	var ss2 SerializableShip
	ss2.Systems = ShipSystems{}
	// Unmarshal it
	err = json.Unmarshal(b, &ss2)
	if err != nil {
		panic(err)
	}

	// Compare them
	if !reflect.DeepEqual(ss1, ss2) {
		log.Println(ss1.Systems["engine"])
		log.Println(ss2.Systems["engine"])
		t.Fail()
	}
}

// Tests entity <-> file relationship more thoroughly
func TestShipSaveToFile(t *testing.T) {
	// load the "Starbridge" ship
	ship1, err := LoadShip("resources/entities/ships/Starbridge.json")
	ship1 = NewShip("Starbridge")
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
		log.Println("not equal:")
		log.Println(ship1.systems["engine"])
		log.Println(ship2.systems["engine"])
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
