// Tideland Go Cell Network - Behaviors - Mapper
//
// Copyright (C) 2010-2014 Frank Mueller / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package behaviors

//--------------------
// IMPORTS
//--------------------

import (
	"github.com/tideland/gocn/v2/cells"
)

//--------------------
// MAPPER BEHAVIOR
//--------------------

// MapFunc is a function type mapping an event to another one.
type MapFunc func(id string, event cells.Event) (topic string, payload interface{})

// mapperBehavior maps the received event to a new event.
type mapperBehavior struct {
	ctx     cells.CellContext
	mapFunc MapFunc
}

// NewMapperBehavior creates a map behavior based on the passed function.
// It emits the mapped events.
func NewMapperBehavior(mf MapFunc) cells.Behavior {
	return &mapperBehavior{nil, mf}
}

// Init the behavior.
func (b *mapperBehavior) Init(ctx cells.CellContext) error {
	b.ctx = ctx
	return nil
}

// Terminate the behavior.
func (b *mapperBehavior) Terminate() error {
	return nil
}

// ProcessEvent executes the simple action func.
func (b *mapperBehavior) ProcessEvent(event cells.Event) error {
	if event.IsRequest() {
		return NewIllegalRequestError(event)
	}
	topic, payload := b.mapFunc(b.ctx.Id(), event)
	b.ctx.RaiseAll(topic, payload)
	return nil
}

// Recover from an error.
func (b *mapperBehavior) Recover(err interface{}) error {
	return nil
}

// EOF
