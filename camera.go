package spacegame

import (
	"github.com/faiface/pixel"
)

type Camera interface {
	Render(Renderer, Entity)
	Motion() pixel.Vec
}

type ChaseCamera struct {
	target Entity
}

func NewChaseCamera(target Entity) *ChaseCamera {
	return &ChaseCamera{
		target: target,
	}
}

func (c *ChaseCamera) Render(renderer Renderer, entity Entity) {
	// we have to transfer the entity from game space to screen space

	// vector from origin to position
	offset := c.target.Coordinates().To(pixel.ZV)

	position := entity.Coordinates().Add(offset)

	// we want the player in the center
	position = position.Add(renderer.Center())

	renderer.Render(entity, position)
}

func (c *ChaseCamera) Motion() pixel.Vec {
	return c.target.Velocity()
}
