// Resource manager maps spacegame.Entity to a pixel.Sprite and ensures that the sprite is loaded
package spacegame

import (
	"fmt"

	"github.com/faiface/pixel"
)

type Resource struct {
	entity      Entity
	sprite      *pixel.Sprite
	scaleFactor float64
}

func (r Resource) Scale() float64 {
    return r.scaleFactor
}

type ResourceManager interface {
	Resource(renderable Entity) *Resource
}

type StandardResourceManager struct {
	basePath  string // for custom resource packs
	resources map[string]Resource
}

func NewStandardResourceManager(baseResourcePath string) *StandardResourceManager {
	return &StandardResourceManager{
		basePath:  baseResourcePath,
		resources: make(map[string]Resource),
	}
}

func (srm *StandardResourceManager) Resource(renderable Entity) *Resource {
	resource, ok := srm.resources[renderable.Name()]
	if !ok {
		return srm.load(renderable)
	}
	return &resource
}

func (srm *StandardResourceManager) load(renderable Entity) *Resource {
	imagePath := fmt.Sprintf("%s/%s", srm.basePath, renderable.ImagePath())

	pic, err := loadPicture(imagePath)

	if err != nil {
		// TODO: Better error handling :)
		panic(err)
		// Log the error and load some standard "missing" image
	}
	sprite := pixel.NewSprite(pic, pic.Bounds())

	resource := Resource{
		entity:     renderable,
		sprite:     sprite,
        scaleFactor: renderable.Bounds().Norm().H() / pic.Bounds().Norm().H(),
	}

	srm.resources[renderable.Name()] = resource

	return &resource
}
