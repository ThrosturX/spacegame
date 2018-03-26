package spacegame

// A colony inhabits celestials. Inhabitants can be expended to launch ships for defensive purposes or to go on missions. Population grows slowly over time, can be increased via some future planned "capture" mechanism and can be decreased via some future planned "siege" mechanism.
// Optionally, colonies can produce goods which may affect trade price of commodities.
type Colony interface {

    // The location of the colony
    Base() Celestial

    // Get population count
    Population() uint64

    // TODO: Launch ships, or react to something

}

