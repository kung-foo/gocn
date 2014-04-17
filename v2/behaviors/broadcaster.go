// Tideland Go Cell Network - Behaviors - Broadcaster
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
	"git.tideland.biz/gocn/cells"
)

//--------------------
// BROADCASTER BEHAVIOR
//--------------------

// broadcasterBehavior is a simple repeater.
type broadcasterBehavior struct {
	ctx cells.CellContext
}

// NewBroadcasterBehavior creates a broadcasting behavior that just emits every
// received event. It's intended to work as an entry point for events, which
// shall be immediately processed by several subscribers.
func NewBroadcasterBehavior() cells.Behavior {
	return &broadcasterBehavior{}
}

// Init the behavior.
func (b *broadcasterBehavior) Init(ctx cells.CellContext) error {
	b.ctx = ctx
	return nil
}

// Terminate the behavior.
func (b *broadcasterBehavior) Terminate() error {
	return nil
}

// ProcessEvent emits the event to all subscribers.
func (b *broadcasterBehavior) ProcessEvent(event cells.Event) error {
	// Check for request.
	if event.IsRequest() {
		return NewIllegalRequestError(event)
	}
	// Process event.
	b.ctx.EmitAll(event)
	return nil
}

// Recover from an error.
func (b *broadcasterBehavior) Recover(err interface{}) error {
	return nil
}

// EOF
