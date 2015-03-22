// Tideland Go Cell Network - Cells - Environment
//
// Copyright (C) 2010-2015 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package cells

//--------------------
// IMPORTS
//--------------------

import (
	"fmt"
	"runtime"
	"time"

	"github.com/tideland/goas/v1/scene"
	"github.com/tideland/goas/v2/identifier"
	"github.com/tideland/goas/v2/logger"
	"github.com/tideland/goas/v3/errors"
)

//--------------------
// CONST
//--------------------

const DefaultTimeout = 5 * time.Second

//--------------------
// ENVIRONMENT
//--------------------

// Environment implements the Environment interface.
type environment struct {
	id           string
	queueFactory EventQueueFactory
	cells        *cluster
}

// NewEnvironment creates a new environment.
func NewEnvironment(options ...Option) Environment {
	env := &environment{
		id:           identifier.NewUUID().String(),
		queueFactory: makeLocalEventQueueFactory(10),
		cells:        newCluster(),
	}
	for _, option := range options {
		option(env)
	}
	runtime.SetFinalizer(env, (*environment).Stop)
	logger.Infof("cells environment %q started", env.ID())
	return env
}

// ID is specified on the Environment interface.
func (env *environment) ID() string {
	return env.id
}

// StartCell is specified on the Environment interface.
func (env *environment) StartCell(id string, behavior Behavior) error {
	return env.cells.startCell(env, id, behavior)
}

// StopCell is specified on the Environment interface.
func (env *environment) StopCell(id string) error {
	return env.cells.stopCell(id)
}

// HasCell is specified on the Environment interface.
func (env *environment) HasCell(id string) bool {
	_, err := env.cells.cell(id)
	return err == nil
}

// Subscribe is specified on the Environment interface.
func (env *environment) Subscribe(emitterId string, subscriberIds ...string) error {
	cell, err := env.cells.cell(emitterId)
	if err != nil {
		return err
	}
	scm, err := env.cells.subset(subscriberIds...)
	if err != nil {
		return err
	}
	cell.subscribers.subscribe(scm)
	return nil
}

// Subscribers is specified on the Environment interface.
func (env *environment) Subscribers(id string) ([]string, error) {
	return env.cells.subscribers(id)
}

// Unsubscribe is specified on the Environment interface.
func (env *environment) Unsubscribe(emitterId string, unsubscriberIds ...string) error {
	cell, err := env.cells.cell(emitterId)
	if err != nil {
		return err
	}
	uscm, err := env.cells.subset(unsubscriberIds...)
	if err != nil {
		return err
	}
	cell.subscribers.unsubscribe(uscm)
	return nil
}

// Emit is specified on the Environment interface.
func (env *environment) Emit(id string, event Event) error {
	return env.cells.emitDirect(id, event)
}

// EmitNew is specified on the Environment interface.
func (env *environment) EmitNew(id, topic string, payload interface{}, scene scene.Scene) error {
	event, err := NewEvent(topic, payload, scene)
	if err != nil {
		return err
	}
	return env.Emit(id, event)
}

// Request is specified on the Environment interface.
func (env *environment) Request(
	id, topic string,
	payload interface{},
	scn scene.Scene,
	timeout time.Duration,
) (interface{}, error) {
	responseChan := make(chan interface{}, 1)
	p := NewPayload(payload).Apply(PayloadValues{ResponseChanPayload: responseChan})
	err := env.EmitNew(id, topic, p, scn)
	if err != nil {
		return nil, err
	}
	select {
	case response := <-responseChan:
		if err, ok := response.(error); ok {
			return nil, err
		}
		return response, nil
	case <-time.After(timeout):
		op := fmt.Sprintf("request %q to %q", topic, id)
		return nil, errors.New(ErrTimeout, errorMessages, op)
	}
}

// Stop manages the proper finalization of an env.
func (env *environment) Stop() error {
	env.cells.stop()
	runtime.SetFinalizer(env, nil)
	logger.Infof("cells environment %q terminated", env.ID())
	return nil
}

// EOF
