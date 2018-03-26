package spacegame

import "image/color"

type Faction interface {

    // The name of the faction
	Name() string

    // Each faction should get a color. Ideas are that bigger factions get the medium to dark shades of colors while smaller factions get light shades and some grays.
	Color() color.Color 

	// Set relationships: (i) initialize (ii) increase, (iii) decrease
    // Modifies the relationship of another faction by delta
    ModifyRelationship(other Faction, delta float64)

	// Check relationships: (i) allies (ii) enemies (iii) neutral (iv) trust?
    // Returns the relationship of another faction
    Relationship(other Faction) float64
}
