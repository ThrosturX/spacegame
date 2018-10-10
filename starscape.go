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
	renderer    Renderer
	camera      Camera
	stars       []star
	resources   []Resource
	scaleFactor float64
}

func NewStarscape(renderer Renderer, camera Camera, density float64) Starscape {
	const layerLow, layerHigh = 0.1, 2.35
	var (
		stars           []star
		resourceManager ResourceManager
		bounds          = renderer.Bounds()
		w               = bounds.W()
		h               = bounds.H()
		numStars        = int(w*h*density) / 1000
		numDust         = numStars / 10
		scaleFactor     = 1.17 // TODO param like density
		extraW          = w*scaleFactor - w
		extraH          = w*scaleFactor - w
	)
    log.Println(extraW, extraH)
	resourceManager = renderer.ResourceManager()
	starResources := resourceManager.FindInCollection("star")
	dustResources := resourceManager.FindInCollection("dust")
	if len(starResources) == 0 {
		log.Println("Could not create a starscape: No stars found")
		return Starscape{}
	}
	if len(dustResources) == 0 {
		log.Println("Error creating a starscape: No dust found -- using stars instead")
		dustResources = starResources
	}
	for i := 0; i < numStars+numDust; i++ {
		var (
			x, y    float64
			topLow  float64 = 0.0
			topMult float64 = 1.00
		)
		resource := &starResources[rand.Intn(len(starResources))]
		x = -extraW/2 + float64(rand.Intn(int(w*scaleFactor)))
		y = -extraH/2 + float64(rand.Intn(int(h*scaleFactor)))

		// create dust
		if i > numStars {
			topLow = 3.7
			topMult = 16.76
			resource = &dustResources[rand.Intn(len(dustResources))]
		}

		star := star{
			resource: resource,
			position: pixel.V(x, y),
			layer:    topLow + layerLow + rand.Float64()*(layerHigh*topMult),
		}
		stars = append(stars, star)
	}

	return Starscape{
		renderer:    renderer,
		camera:      camera,
		stars:       stars,
		resources:   starResources,
		scaleFactor: scaleFactor,
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

		// if position is out of bounds, reset the image and the coordinate component that was out of bounds
		// find the missing component
		starBox := sc.renderer.Bounds()
		starBox = starBox.Resized(starBox.Center(), pixel.V(starBox.W()*sc.scaleFactor, starBox.H()*sc.scaleFactor))
		if !starBox.Contains(sc.stars[i].position) {
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
	extraH := h*sc.scaleFactor - h
	extraW := w*sc.scaleFactor - w
	sc.stars[starIndex].resource = &sc.resources[rand.Intn(len(sc.resources))]

	if candidatePosition.X > w+extraW {
		candidatePosition.X = 0-extraW
		candidatePosition.Y = float64(rand.Intn(int(h)))
	} else if candidatePosition.X < 0-extraW {
		candidatePosition.X = w+extraW
		candidatePosition.Y = float64(rand.Intn(int(h)))
	}
	if candidatePosition.Y > h+extraH {
		candidatePosition.X = float64(rand.Intn(int(w)))
		candidatePosition.Y = 0-extraH
	} else if candidatePosition.Y < 0-extraH {
		candidatePosition.X = float64(rand.Intn(int(w)))
		candidatePosition.Y = h+extraH
	}

	sc.stars[starIndex].position = candidatePosition
}
