// Tideland Go Cell Network - Behaviors - Finite State Machine
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
// FSM BEHAVIOR
//--------------------

// State is the signature of a function or method which processes
// an event and returns the following state or an error.
type State func(ctx cells.CellContext, event cells.Event) (State, error)

// fsmBehavior runs the finite state machine.
type fsmBehavior struct {
	ctx   cells.CellContext
	state State
	done  bool
	err   error
}

// NewFSMBehavior creates a finite state machine behavior based on the
// passed initial state function. The function is called with the event
// has to return the next state, which can be the same one. In case of
// nil the stae will be transfered into a generic end state, if an error
// is returned the state is a generic error state.
func NewFSMBehavior(state State) cells.Behavior {
	return &fsmBehavior{nil, state, false, nil}
}

// Init the behavior.
func (b *fsmBehavior) Init(ctx cells.CellContext) error {
	b.ctx = ctx
	return nil
}

// Terminate the behavior.
func (b *fsmBehavior) Terminate() error {
	return nil
}

// ProcessEvent executes the state function and stores
// the returned new state.
func (b *fsmBehavior) ProcessEvent(event cells.Event) error {
	if event.IsRequest() {
		switch event.Topic() {
		case "isOk?":
			return event.Respond(b.done, nil)
		case "isError?":
			return event.Respond(b.err != nil, b.err)
		}
	}
	if b.done {
		return nil
	}
	state, err := b.state(b.ctx, event)
	if err != nil {
		b.done = true
		b.err = err
	} else if state == nil {
		b.done = true
	}
	b.state = state
	return nil
}

// Recover from an error.
func (b *fsmBehavior) Recover(err interface{}) error {
	b.done = true
	b.err = NewCannotRecoverError(err)
	return nil
}

// IsOkEndState checks if the finite state machine with the given
// id is in a positive end state.
func IsOkEndState(env cells.Environment, id string) (bool, error) {
	var okState bool
	err := env.Request(id, "isOk?", nil, &okState)
	return okState, err
}

// IsErrorEndState checks if the finite state machine with the given
// id is in an faulty end state.
func IsErrorEndState(env cells.Environment, id string) (bool, error) {
	var errState bool
	err := env.Request(id, "isError?", nil, &errState)
	return errState, err
}

// EOF
