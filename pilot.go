package spacegame

type Pilot interface {
    Name() string
    Faction() Faction
    Ship() PilotableShip
    //    Fleet() *Fleet // TODO: Fleets: Escorts and coordinated groups of AI
    Update(dt float64) // Tick. For players this is input handling, for AI this is thinking
}
