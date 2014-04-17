// Tideland Go Cell Network - Cells - Cell
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

	"github.com/tideland/goas/v2/identifier"
	"github.com/tideland/goas/v2/logger"
	"github.com/tideland/goas/v2/loop"
	"github.com/tideland/goas/v2/monitoring"
)

//--------------------
// CELL
//--------------------

// cell for event processing.
type cell struct {
	env         *environment
	id          string
	behavior    Behavior
	subscribers *cluster
	queue       EventQueue
	loop        loop.Loop
	measuringId string
}

// newCell create a new cell around a behavior.
func newCell(env *environment, id string, behavior Behavior) (*cell, error) {
	// Init cell runtime.
	c := &cell{
		env:         env,
		id:          id,
		behavior:    behavior,
		subscribers: newCluster(),
		measuringId: identifier.Identifier("cells", env.id, "cell", identifier.TypeAsIdentifierPart(behavior)),
	}
	// Create queue.
	queue, err := env.queueFactory(env)
	if err != nil {
		return nil, annotateError(err, ErrCellInit, id)
	}
	c.queue = queue
	// Init behavior.
	if err := behavior.Init(c); err != nil {
		return nil, annotateError(err, ErrCellInit, id)
	}
	c.loop = loop.GoRecoverable(c.backendLoop, c.checkRecovering)

	logger.Infof("cell %q started", id)

	return c, nil
}

// Environment returns the environment the cell is running in.
func (c *cell) Environment() Environment {
	return c.env
}

// Id returns the id of the cell.
func (c *cell) Id() string {
	return c.id
}

// Subscribers returns the ids of the subscriber cells
func (c *cell) Subscribers() []string {
	return c.subscribers.ids()
}

// Emit emits an event to one addressed subscriber of a cell.
func (c *cell) Emit(id string, event Event) error {
	return c.subscribers.emit(id, event)
}

// EmitAll emits an event to all subscribers of a cell.
func (c *cell) EmitAll(event Event) error {
	return c.subscribers.emitAll(event)
}

// Raise creates an event and emits it to one addressed subscriber of a cell.
func (c *cell) Raise(id, topic string, payload interface{}) error {
	event, err := newEvent(topic, payload)
	if err != nil {
		return err
	}
	return c.Emit(id, event)
}

// RaiseAll creates an event and emits it to all subscribers of a cell.
func (c *cell) RaiseAll(topic string, payload interface{}) error {
	event, err := newEvent(topic, payload)
	if err != nil {
		return err
	}
	return c.EmitAll(event)
}

// processEvent tells the cell to process an event.
func (c *cell) processEvent(event Event) error {
	return c.queue.Push(event)
}

// subscribe adds the passed cells to the subscriptions.
func (c *cell) subscribe(sc *cluster) {
	c.subscribers.subscribe(sc)
}

// unsubscribe removes the passed cells from the subscriptions.
func (c *cell) unsubscribe(uc *cluster) {
	c.subscribers.unsubscribe(uc)
}

// stop terminates the cell.
func (c *cell) stop() error {
	defer func() {
		if err := c.queue.Stop(); err != nil {
			logger.Errorf("cannot stop queue of cell %q: %v", c.id, err)
		}
		logger.Infof("cell %q terminated", c.id)
	}()
	return c.loop.Stop()
}

// backendLoop is the backend for the processing of messages.
func (c *cell) backendLoop(l loop.Loop) error {
	monitoring.IncrVariable(c.measuringId)
	monitoring.IncrVariable(identifier.Identifier("cells", c.env.Id(), "total-cells"))
	defer monitoring.DecrVariable(identifier.Identifier("cells", c.env.Id(), "total-cells"))
	defer monitoring.DecrVariable(c.measuringId)

	for {
		select {
		case <-c.loop.ShallStop():
			return c.behavior.Terminate()
		case event := <-c.queue.Events():
			measuring := monitoring.BeginMeasuring(c.measuringId)
			err := c.behavior.ProcessEvent(event)
			if err != nil {
				c.loop.Kill(err)
				continue
			}
			measuring.EndMeasuring()
		}
	}
}

// checkRecovering checks if the cell may recover after a panic. It will
// signal an error and let the cell stop working if there have been 12 recoverings
// during the last minute or the behaviors Recover() signals, that it cannot
// handle the error.
func (c *cell) checkRecovering(rs loop.Recoverings) (loop.Recoverings, error) {
	logger.Errorf("recovering cell %q after error: %v", c.id, rs.Last().Reason)
	// Check frequency.
	if rs.Frequency(12, time.Minute) {
		return nil, newError(ErrRecoveredTooOften, rs.Last().Reason)
	}
	// Try to recover.
	if err := c.behavior.Recover(rs.Last().Reason); err != nil {
		return nil, annotateError(err, ErrEventRecovering, rs.Last().Reason)
	}
	return rs.Trim(12), nil
}

// EOF
