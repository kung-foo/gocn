// Tideland Go Cell Network - Cells - Environment
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
	"runtime"
	"time"

	"github.com/tideland/goas/v2/identifier"
	"github.com/tideland/goas/v2/logger"
	"github.com/tideland/goas/v3/errors"
)

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
		queueFactory: MakeLocalEventQueueFactory(10),
		cells:        newCluster(),
	}

	env.Options(options...)
	runtime.SetFinalizer(env, (*environment).Stop)
	logger.Infof("cells environment %q started", env.Id())

	return env
}

// Id returns the environment id.
func (env *environment) Id() string {
	return env.id
}

// Configuration returns the configuration of the environment.
func (env *environment) Options(options ...Option) Options {
	previous := make(Options, len(options))
	for index, option := range options {
		previous[index] = option(env)
	}
	return previous
}

// StartCell starts a new cell with a given id and its behavior.
func (env *environment) StartCell(id string, behavior Behavior) error {
	return env.cells.startCell(env, id, behavior)
}

// StopCell stops and removes the cell with the given id.
func (env *environment) StopCell(id string) error {
	return env.cells.stopCell(id)
}

// HasCell returns true if the cell with the given id exists.
func (env *environment) HasCell(id string) bool {
	_, err := env.cells.cell(id)
	return err == nil
}

// Subscribe assigns cells as receivers of the emitted
// events of the first cell.
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

// Subscribers returns the subscribers of the passed id.
func (env *environment) Subscribers(id string) ([]string, error) {
	return env.cells.subscribers(id)
}

// Unsubscribe removes the assignment of emitting und subscribed cells.
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

// Emit emits an event to the cell with a given id.
func (env *environment) Emit(id string, event Event) error {
	return env.cells.emit(id, event)
}

// Raise creates an event and emits it to the cell with a given id.
func (env *environment) Raise(id, topic string, payload interface{}) error {
	event, err := newEvent(topic, payload)
	if err != nil {
		return err
	}
	return env.Emit(id, event)
}

// Request retrieves information from the cell with the given id.
func (env *environment) Request(id, topic string, payload, response interface{}) error {
	request, err := newRequest(topic, payload)
	if err != nil {
		return err
	}
	err = env.Emit(id, request)
	if err != nil {
		return err
	}
	t := time.NewTicker(10 * time.Second)
	select {
	case reply := <-request.replyc:
		err = reply.Payload(response)
		if err != nil {
			return err
		}
		return reply.err
	case <-t.C:
		return errors.New(ErrInactive, errorMessages, id)
	}
}

// Stop manages the proper finalization of an env.
func (env *environment) Stop() error {
	env.cells.stop()
	runtime.SetFinalizer(env, nil)

	logger.Infof("cells environment %q terminated", env.Id())

	return nil
}

// EOF
