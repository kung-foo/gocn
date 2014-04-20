// Tideland Go Cell Network - Behaviors - Filter
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
// FILTER BEHAVIOR
//--------------------

// FilterFunc is a function type checking if an event shall be filtered.
type FilterFunc func(id string, event cells.Event) bool

// filterBehavior is a simple repeater using the filter
// function to check if an event shall be emitted.
type filterBehavior struct {
	ctx        cells.CellContext
	filterFunc FilterFunc
}

// NewFilterBehavior creates a filter behavior based on the passed function.
// It emits every received event for which the filter function returns true.
func NewFilterBehavior(ff FilterFunc) cells.Behavior {
	return &filterBehavior{nil, ff}
}

// Init the behavior.
func (b *filterBehavior) Init(ctx cells.CellContext) error {
	b.ctx = ctx
	return nil
}

// Terminate the behavior.
func (b *filterBehavior) Terminate() error {
	return nil
}

// ProcessEvent emits the event when the filter func returns true.
func (b *filterBehavior) ProcessEvent(event cells.Event) error {
	if event.IsRequest() {
		return NewIllegalRequestError(event)
	}
	if b.filterFunc(b.ctx.Id(), event) {
		b.ctx.EmitAll(event)
	}
	return nil
}

// Recover from an error.
func (b *filterBehavior) Recover(err interface{}) error {
	return nil
}

// EOF
