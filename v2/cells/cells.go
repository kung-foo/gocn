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

import ()

//--------------------
// ENVIRONMENT
//--------------------

// Environment is a set of networked cells.
type Environment interface {
	// Id returns the environment id.
	Id() string

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

	// Raise creates an event and emits it to the cell with a given id.
	Raise(id, topic string, payload interface{}) error

	// Request retrieves information from the cell with the given id.
	Request(id, topic string, payload, response interface{}) error

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
	Init(ctx CellContext) error

	// Terminate the behavior.
	Terminate() error

	// ProcessEvent processes an event and can emit own events.
	ProcessEvent(event Event) error

	// Recover from an error.
	Recover(r interface{}) error
}

//--------------------
// CELL CONTEXT
//--------------------

// CellContext gives a behavior access to its environment.
type CellContext interface {
	// Environment returns the environment of the cell.
	Environment() Environment

	// Id returns the id of the cell.
	Id() string

	// Subscribers returns the ids of the subscriber cells
	Subscribers() []string

	// Emit emits an event to one addressed subscriber of a cell.
	Emit(id string, event Event) error

	// EmitAll emits an event to all subscribers of a cell.
	EmitAll(event Event) error

	// Raise creates an event and emits it to one addressed subscriber of a cell.
	Raise(id, topic string, payload interface{}) error

	// RaiseAll creates an event and emits it to all subscribers of a cell.
	RaiseAll(topic string, payload interface{}) error
}

// EOF
