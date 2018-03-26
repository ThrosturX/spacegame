package spacegame

import (
	"log"
	"math"
)

func normalizeAngle(rads float64) float64 {
	var out float64 = rads
	for out > math.Pi {
		out = out - 2*math.Pi
	}
	for out <= -math.Pi {
		out = out + 2*math.Pi
	}
	if out != rads {
		log.Printf("Normalized %f to %f\n", rads, out)
	}
	return out
}
