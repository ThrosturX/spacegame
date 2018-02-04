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

	offset := c.target.Coordinates().To(pixel.ZV)

	position := entity.Coordinates().Add(offset)

	renderer.Render(entity, position)
}

func (c *ChaseCamera) Motion() pixel.Vec {
	return c.target.Velocity()
}
