// Resource manager maps spacegame.Entity to a pixel.Sprite and ensures that the sprite is loaded
// Future improvements include garbage collection and other stuff that resource managers generally do
package spacegame

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/faiface/pixel"
)

type Resource struct {
	entity      Entity
	sprite      *pixel.Sprite
	rect        pixel.Rect
	scaleFactor float64
	collection  string
}

func (r Resource) Entity() Entity {
	return r.entity
}

func (r Resource) Bounds() pixel.Rect {
	return r.rect
}

type ResourceManager interface {
	CreateResource(renderable Entity, path string)
	Find(search string) []Resource
	FindInCollection(collection string) []Resource
	ImportDefault() // TODO: Import(options GameOptions)
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

func (srm *StandardResourceManager) CreateResource(renderable Entity, path string) {
	imagePath := fmt.Sprintf("%s/%s", srm.basePath, path)

	pic, err := loadPicture(imagePath)
	if err != nil {
		// TODO: Better error handling :)
		panic(err)
		// Log the error and load some standard "missing" image
	}
	resource := srm.createResource(pic, renderable)
	//		scaleFactor: renderable.Bounds().Norm().H() / pic.Bounds().Norm().H(),

	srm.resources[renderable.Name()] = resource
}

func (srm *StandardResourceManager) Find(search string) []Resource {
	var matched []Resource
	for i, resource := range srm.resources {
		name := resource.entity.Name()
		match := strings.Index(name, search)
		if match != -1 {
			matched = append(matched, srm.resources[i])
		} else {
			log.Printf("%s didn't match %s\n", resource.entity.Name(), search)
		}
	}
	return matched
}
func (srm *StandardResourceManager) FindInCollection(collection string) []Resource {
	var matched []Resource
	for i, resource := range srm.resources {
		match := strings.Index(resource.collection, collection)
		if match != -1 {
			matched = append(matched, srm.resources[i])
		}
	}

	log.Println("Created this beautiful collection for", collection)
	log.Println(matched)
	return matched
}

// Imports everything we can find
func (srm *StandardResourceManager) ImportDefault() {
	// TODO: "makeWalkHandler"
	// import Ships
	// walk resources/entities/ships

	shipImporter := func(path string, info os.FileInfo, err error) error {
		// fail on error
		if err != nil {
			return err
		}

		// skip directories
		if info.IsDir() {
			return nil
		}

		// only import json files
		if filepath.Ext(path) != ".json" {
			return nil // TODO: log it?
		}
		ship, err := LoadShip(path)
		if err != nil {
			return err
		}
		collection, filename := filepath.Split(path)

		name := strings.Replace(filename, filepath.Ext(filename), "", 1)

		imagePath := fmt.Sprintf("%s/images/ships/%s.png", srm.basePath, strings.ToLower(name))

		pic, err := loadPicture(imagePath)
		if err != nil {
			// TODO: error handler
			panic(err)
		}

		// create the resource
		resource := srm.createResource(pic, ship)
		resource.collection = collection

		srm.resources[name] = resource

		log.Println("Imported", name)

		return nil
	}

	systemImporter := func(path string, info os.FileInfo, err error) error {
		// fail on error
		if err != nil {
			return err
		}

		// skip directories
		if info.IsDir() {
			return nil
		}

		// only import json files
		if filepath.Ext(path) != ".json" {
			return nil // TODO: log it?
		}
		sys, err := LoadSystem(path)
		if err != nil {
			return err
		}
		collection, _ := filepath.Split(path)

		for _, c := range sys.Celestials() {
			imagePath := fmt.Sprintf("%s/%s", srm.basePath, c.ImagePath())

			pic, err := loadPicture(imagePath)
			if err != nil {
				// TODO: error handler
				panic(err)
			}

			// create the resource
			resource := srm.createResource(pic, c)
			resource.collection = collection

			srm.resources[c.name] = resource
			log.Println("Imported", c.name)
		}

		log.Println("Imported", sys.name)

		return nil
	}

	entityImporter := func(path string, info os.FileInfo, err error) error {
		// fail on error
		if err != nil {
			return err
		}

		// skip directories
		if info.IsDir() {
			return nil
		}

		// prepare the resource
		collection, filename := filepath.Split(path)

		name := strings.Replace(filename, filepath.Ext(filename), "", 1)

		pic, err := loadPicture(path)
		if err != nil {
			// TODO: error handler
			panic(err)
		}

		// create the entity
		entity := NewBasicEntity(name, pic.Bounds())

		// create the resource
		resource := srm.createResource(pic, entity)
		resource.collection = collection

		srm.resources[name] = resource

		log.Println("Imported", name)

		return nil
	}

	shipPath := fmt.Sprintf("%s/entities/ships", srm.basePath)
	systemPath := fmt.Sprintf("%s/universe/systems", srm.basePath)
	starPath := fmt.Sprintf("%s/images/stars", srm.basePath)

	var err error

	// import ships
	err = filepath.Walk(shipPath, shipImporter)
	if err != nil {
		panic(err)
	}

	// import celestials
	err = filepath.Walk(systemPath, systemImporter)
	if err != nil {
		panic(err)
	}

	// import stars
	err = filepath.Walk(starPath, entityImporter)
	if err != nil {
		panic(err)
	}
}

func (srm *StandardResourceManager) Resource(renderable Entity) *Resource {
	resource, ok := srm.resources[renderable.Name()]
	if !ok {
		err := errors.New(fmt.Sprintf("Resource not found: %s", renderable.Name()))
		// TODO: Don't panic!
		panic(err)
	}
	return &resource
}

func (srm *StandardResourceManager) createResource(pic pixel.Picture, entity Entity) Resource {
	return Resource{
		entity: entity,
		sprite: pixel.NewSprite(pic, pic.Bounds()),
		rect:   pic.Bounds(),
	}
}
