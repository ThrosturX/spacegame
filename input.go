// TODO: Serialize this class to settings.json or something like that
package spacegame

import (
	"errors"
	"fmt"

	"github.com/faiface/pixel/pixelgl"
)

const (
	actionAccel       = "accel"
	actionTurnLeft    = "left"
	actionTurnRight   = "right"
	actionReverse     = "reverse"
	actionAlign       = "align"
	actionTargetPrev  = "targetPrev"
	actionTargetNext  = "targetNext"
	actionLand        = "targetLand"
	actionClearTarget = "clearTarget"
)

type Controllable interface {
	Process(pilotAction)
}

type pilotAction struct {
	key string
	dt  float64
}

type Controller struct {
	// map keybind -> action key
	bindings  map[int]string
	entity    Controllable
	window    *pixelgl.Window
	repeaters map[string]bool
}

func NewPlayerController(window *pixelgl.Window, entity Controllable) *Controller {
	c := &Controller{
		entity: entity,
		window: window,
	}
	c.ResetBindings()
	c.repeaters = make(map[string]bool)
	c.repeaters[actionAccel] = true
	c.repeaters[actionAlign] = true
	c.repeaters[actionTurnLeft] = true
	c.repeaters[actionTurnRight] = true
	c.repeaters[actionReverse] = true
	c.repeaters[actionTargetPrev] = false
	c.repeaters[actionTargetNext] = false
	c.repeaters[actionLand] = false
    c.repeaters[actionClearTarget] = false

	return c
}

func (c *Controller) relay(dt float64) {
	for key, cmd := range c.bindings {
		var pressFunc func(pixelgl.Button) bool
		repeater, ok := c.repeaters[cmd]

		if !ok || !repeater {
			pressFunc = c.window.JustPressed
		} else {
			pressFunc = c.window.Pressed
		}

		if pressFunc(pixelgl.Button(key)) {
			c.entity.Process(pilotAction{cmd, dt})
		}
	}
}

func (c *Controller) SetKey(key pixelgl.Button, action string) error {

	if action, ok := c.bindings[int(key)]; ok {
		return errors.New(fmt.Sprintf("key <%s> is already bound to <%v>.", key, action))
	}
	c.bindings[int(key)] = action
	return nil
}

func (c *Controller) Unbind(key pixelgl.Button) {
	delete(c.bindings, int(key))
}

func (c *Controller) ResetBindings() {
	c.bindings = make(map[int]string)
	c.SetKey(pixelgl.KeyUp, actionAccel)
	c.SetKey(pixelgl.KeyLeft, actionTurnLeft)
	c.SetKey(pixelgl.KeyRight, actionTurnRight)
	c.SetKey(pixelgl.KeyDown, actionReverse)
	c.SetKey(pixelgl.KeyTab, actionTargetNext)
	//	c.SetKey(pixelgl._, actionTargetPrev)
	c.SetKey(pixelgl.KeyA, actionAlign)
	c.SetKey(pixelgl.KeyL, actionLand)
    c.SetKey(pixelgl.KeyC, actionClearTarget)
}
