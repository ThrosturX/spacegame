// Represents a ship's capability to target other ships
package spacegame

type Scanner interface {
    NextTarget() Entity // TargetableEntity interface TODO
    PrevTarget() Entity // TargetableEntity interface TODO
    Range() float64
    // TODO: Other scanner settings
}
