package spacegame

import (
	"log"
	"math/rand"

	"github.com/faiface/pixel"
)

type Background interface {
	Displace(float64)
	Render()
}

type star struct {
	resource *Resource
	position pixel.Vec
	layer    float64
}

type Starscape struct {
	renderer  Renderer
	camera    Camera
	stars     []star
	resources []Resource
}

func NewStarscape(renderer Renderer, camera Camera) Starscape {
	const layerLow, layerHigh = 0.1, 1.35
	var (
		stars           []star
		resourceManager ResourceManager
	)
	resourceManager = renderer.ResourceManager()
	starResources := resourceManager.FindInCollection("star")
	if len(starResources) == 0 {
		log.Println("Could not create a starscape: No stars found")
		return Starscape{}
	}
	for i := 0; i < 100; i++ {
		var x, y float64
		index := rand.Intn(len(starResources))
		bounds := renderer.Bounds()
		x = float64(rand.Intn(int(bounds.W())))
		y = float64(rand.Intn(int(bounds.H())))

		star := star{
			resource: &starResources[index],
			position: pixel.V(x, y),
			layer:    layerLow + rand.Float64()*layerHigh,
		}
		stars = append(stars, star)
	}

	return Starscape{
		renderer:  renderer,
		camera:    camera,
		stars:     stars,
		resources: starResources,
	}
}

func (sc Starscape) Displace(dt float64) {
	vector := sc.camera.Motion().Scaled(dt)
	if vector.Len() == 0 {
		return
	}

	for i, star := range sc.stars {
		displacement := vector.Scaled(star.layer)
		sc.stars[i].position = star.position.Sub(displacement)

		// TODO: Reset the image
		// if position is out of bounds, reset the image and the coordinate component that was out of bounds
		// find the missing component
		if !sc.renderer.Bounds().Contains(sc.stars[i].position) {
			sc.recreate(i)
		}
	}
}

func (sc Starscape) Render() {
	// TODO: Sort by layer!
	for _, star := range sc.stars {
		sc.renderer.Render(star.resource.entity, star.position)
	}
}

func (sc Starscape) recreate(starIndex int) {
	candidatePosition := sc.stars[starIndex].position
	bounds := sc.renderer.Bounds()
	h := bounds.H()
	w := bounds.W()
	sc.stars[starIndex].resource = &sc.resources[rand.Intn(len(sc.resources))]

	if candidatePosition.X > w {
		candidatePosition.X = 0
		candidatePosition.Y = float64(rand.Intn(int(h)))
	} else if candidatePosition.X < 0 {
		candidatePosition.X = w
		candidatePosition.Y = float64(rand.Intn(int(h)))
	}
	if candidatePosition.Y > h {
		candidatePosition.X = float64(rand.Intn(int(w)))
		candidatePosition.Y = 0
	} else if candidatePosition.Y < 0 {
		candidatePosition.X = float64(rand.Intn(int(w)))
		candidatePosition.Y = h
	}

	sc.stars[starIndex].position = candidatePosition
}
