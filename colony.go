package spacegame

// A colony inhabits celestials. Inhabitants can be expended to launch ships for defensive purposes or to go on missions. Population grows slowly over time, can be increased via some future planned "capture" mechanism and can be decreased via some future planned "siege" mechanism.
// Optionally, colonies can produce goods which may affect trade price of commodities.
type Colony interface {

	// The location of the colony
	Base() Celestial

	// The faction (clan) living in the colony
	Faction() Faction

	// Get population count
	Population() uint64

	// TODO: Launch ships, or react to something

}

// A clan inhabits a colony, so ClanColony implements its needs
type ClanColony struct {
	base       Celestial
	faction    Faction
	population uint64
}

// Add colonists (from a colonize mission)
func (cc *ClanColony) AddColonists(number uint64) {
	cc.population += number
}

func (cc *ClanColony) Base() Celestial {
	return cc.base
}

func (cc *ClanColony) Faction() Faction {
	return cc.faction
}

func (cc *ClanColony) Population() uint64 {
	return cc.population
}
