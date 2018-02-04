// TODO: Serialize this class to settings.json or something like that
package spacegame

import (
	"errors"
	"fmt"

	"github.com/faiface/pixel/pixelgl"
)

const (
	actionAccel      = "accel"
	actionTurnLeft   = "left"
	actionTurnRight  = "right"
	actionReverse    = "reverse"
	actionTargetNext = "targetNext"
)

type Controllable interface {
	CmdChan() chan pilotAction
}

type pilotAction struct {
	key string
	dt  float64
}

type Controller struct {
	// map keybind -> action key
	bindings map[int]string
	entity   Controllable
	window   *pixelgl.Window
}

func NewPlayerController(window *pixelgl.Window, entity Controllable) *Controller {
	c := &Controller{
		entity: entity,
		window: window,
	}
	c.ResetBindings()
	return c
}

func (c *Controller) relay(dt float64) {
	for key, cmd := range c.bindings {
		if c.window.Pressed(pixelgl.Button(key)) {
			c.entity.CmdChan() <- pilotAction{cmd, dt}
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
}
