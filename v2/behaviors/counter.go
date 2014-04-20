// Tideland Go Cell Network - Behaviors - Counter
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
// COUNTER BEHAVIOR
//--------------------

// Counters is a set of named counters and their values.
type Counters map[string]int64

// CounterFunc is the signature of a function which analyzis
// an event and returns, which counters shall be incremented.
type CounterFunc func(id string, event cells.Event) []string

// counterBehavior counts events based on the counter function.
type counterBehavior struct {
	ctx         cells.CellContext
	counterFunc CounterFunc
	counters    Counters
}

// NewCounterBehavior creates a counter behavior based on the passed
// function. It increments and emits those counters named by the result
// of the counter function. The counters can be retrieved with the
// request "counters" and reset with "reset".
func NewCounterBehavior(cf CounterFunc) cells.Behavior {
	return &counterBehavior{nil, cf, make(Counters)}
}

// Init the behavior.
func (b *counterBehavior) Init(ctx cells.CellContext) error {
	b.ctx = ctx
	return nil
}

// Terminate the behavior.
func (b *counterBehavior) Terminate() error {
	return nil
}

// ProcessEvent counts the event for the return value of the counter func
// and emits this value.
func (b *counterBehavior) ProcessEvent(event cells.Event) error {
	if event.IsRequest() {
		switch event.Topic() {
		case "counters":
			return event.Respond(b.counters, nil)
		case "reset":
			b.counters = make(map[string]int64)
			return event.Respond(true, nil)
		}
		return NewIllegalRequestError(event)
	}
	cids := b.counterFunc(b.ctx.Id(), event)
	if cids != nil {
		for _, cid := range cids {
			v, ok := b.counters[cid]
			if ok {
				b.counters[cid] = v + 1
			} else {
				b.counters[cid] = 1
			}
			topic := "counter:" + cid
			b.ctx.RaiseAll(topic, b.counters[cid])
		}
	}
	return nil
}

// Recover from an error.
func (b *counterBehavior) Recover(err interface{}) error {
	return nil
}

// EOF
