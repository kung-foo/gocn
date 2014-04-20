// Tideland Go Cell Network - Behaviors - Collector
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
// COLLECTOR BEHAVIOR
//--------------------

// EventData represents the pure collected event data.
type EventData struct {
	Topic   string
	Payload interface{}
}

// newEventData returns the passed event as event data to collect.
func newEventData(event cells.Event) EventData {
	data := EventData{}
	data.Topic = event.Topic()
	event.Payload(&data.Payload)
	return data
}

// collectorBehavior collects events for debugging.
type collectorBehavior struct {
	ctx       cells.CellContext
	max       int
	collected []EventData
}

// NewCollectorBehaviorFactory creates a collector behavior. It collects
// a configured maximum number events emitted directly or by subscription.
// The event is passed through. The collected events can be requested with
// the request "collected" and reset with "reset".
func NewCollectorBehavior(max int) cells.Behavior {
	return &collectorBehavior{nil, max, []EventData{}}
}

// Init the behavior.
func (b *collectorBehavior) Init(ctx cells.CellContext) error {
	b.ctx = ctx
	return nil
}

// Terminate the behavior.
func (b *collectorBehavior) Terminate() error {
	return nil
}

// ProcessEvent collects and re-emits events.
func (b *collectorBehavior) ProcessEvent(event cells.Event) error {
	// Check for request.
	if event.IsRequest() {
		switch event.Topic() {
		case "collected":
			return event.Respond(b.collected, nil)
		case "reset":
			b.collected = []EventData{}
			return event.Respond(true, nil)
		}
		return NewIllegalRequestError(event)
	}
	// Process event.
	if len(b.collected) == b.max {
		b.collected = b.collected[1:]
	}
	b.collected = append(b.collected, newEventData(event))
	b.ctx.EmitAll(event)
	return nil
}

// Recover from an error.
func (b *collectorBehavior) Recover(err interface{}) error {
	return nil
}

// EOF
