// Tideland Go Cell Network - Cells
//
// Copyright (C) 2010-2014 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package cells

//--------------------
// IMPORTS
//--------------------

import (
	"time"

	"github.com/tideland/goas/v1/scene"
)

//--------------------
// ENVIRONMENT
//--------------------

// Environment is a set of networked cells.
type Environment interface {
	// ID returns the environment ID.
	ID() string

	// Options sets options of the environment.
	Options(options ...Option) Options

	// StartCell starts a new cell with a given id and its behavior.
	StartCell(id string, behavior Behavior) error

	// StopCell stops and removes the cell with the given id.
	StopCell(id string) error

	// HasCell returns true if the cell with the given id exists.
	HasCell(id string) bool

	// Subscribe assigns cells as receivers of the emitted
	// events of the first cell.
	Subscribe(emitterId string, subscriberIds ...string) error

	// Subscribers returns the subscribers of the passed id.
	Subscribers(id string) ([]string, error)

	// Unsubscribe removes the assignment of emitting und subscribed cells.
	Unsubscribe(emitterId string, unsubscriberIds ...string) error

	// Emit emits an event to the cell with a given id.
	Emit(id string, event Event) error

	// EmitNew creates an event and emits it to the cell with a given ID.
	EmitNew(id, topic string, payload interface{}, scn scene.Scene) error

	// Request creates and emits an event to the cell with the given ID.
	// It is intended as request which has to be responded to with
	// event.Respond().
	Request(id, topic string, payload interface{}, scn scene.Scene, timeout time.Duration) (interface{}, error)

	// Stop manages the proper finalization of an env.
	Stop() error
}

//--------------------
// BEHAVIOR
//--------------------

// Behavior is the interface that has to be implemented
// for the usage through the cells.
type Behavior interface {
	// Init the deployed behavior inside an environment.
	Init(ctx Context) error

	// Terminate the behavior.
	Terminate() error

	// ProcessEvent processes an event and can emit own events.
	ProcessEvent(event Event) error

	// Recover from an error.
	Recover(r interface{}) error
}

//--------------------
// CONTEXT
//--------------------

// Context gives a behavior access to its environment.
type Context interface {
	// Environment returns the environment of the cell.
	Environment() Environment

	// ID returns the ID of the cell.
	ID() string

	// Subscribers returns the ids of the subscriber cells
	Subscribers() []string

	// Emit emits an event to all subscribers of a cell.
	Emit(event Event) error

	// EmitNew creates an event and emits it to all subscribers of a cell.
	EmitNew(topic string, payload interface{}, scene scene.Scene) error
}

// EOF
