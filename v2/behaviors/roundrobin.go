// Tideland Go Cell Network - Behaviors - Round-Robin
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
	"sort"

	"github.com/tideland/gocn/v2/cells"
)

//--------------------
// ROUND ROBIN BEHAVIOR
//--------------------

// roundRobinBehavior emit the received events round robin to its
// subscribers in a very simple way.
type roundRobinBehavior struct {
	ctx     cells.CellContext
	current int
}

// NewRoundRobinBehavior creates a behavior emitting the received events to
// its subscribers in a very simple way. Subscriptions or unsubscriptions
// during runtime may influence the order.
func NewRoundRobinBehavior() cells.Behavior {
	return &roundRobinBehavior{nil, 0}
}

// Init the behavior.
func (b *roundRobinBehavior) Init(ctx cells.CellContext) error {
	b.ctx = ctx
	return nil
}

// Terminate the behavior.
func (b *roundRobinBehavior) Terminate() error {
	return nil
}

// ProcessEvent emits the event round robin to the subscribers.
func (b *roundRobinBehavior) ProcessEvent(event cells.Event) error {
	if event.IsRequest() {
		return NewIllegalRequestError(event)
	}
	subscribers := b.ctx.Subscribers()
	if len(subscribers) > 0 {
		var id string
		sort.Strings(subscribers)
		if b.current < len(subscribers) {
			id = subscribers[b.current]
			b.current++
		} else {
			id = subscribers[0]
			b.current = 1
		}
		return b.ctx.Emit(id, event)
	}
	return nil
}

// Recover from an error.
func (b *roundRobinBehavior) Recover(err interface{}) error {
	return nil
}

// EOF
