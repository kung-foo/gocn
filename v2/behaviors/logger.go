// Tideland Go Cell Network - Behaviors - Logger
//
// Copyright (C) 2010-2013 Frank Mueller / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package behaviors

//--------------------
// IMPORTS
//--------------------

import (
	"github.com/tideland/goas/v2/logger"
	"github.com/tideland/gocn/v2/cells"
)

//--------------------
// LOGGER BEHAVIOR
//--------------------

// loggerBehavior is a behaior for the logging of events.
type loggerBehavior struct {
	ctx cells.CellContext
}

// NewLoggerBehavior creates a logging behavior. It logs emitted
// events with info level.
func NewLoggerBehavior() cells.Behavior {
	return &loggerBehavior{}
}

// Init the behavior.
func (b *loggerBehavior) Init(ctx cells.CellContext) error {
	b.ctx = ctx
	return nil
}

// Terminate the behavior.
func (b *loggerBehavior) Terminate() error {
	return nil
}

// ProcessEvent logs the event at info level.
func (b *loggerBehavior) ProcessEvent(event cells.Event) error {
	logger.Infof("(%s) processing: %q", b.ctx.Id(), event)
	return nil
}

// Recover from an error. Can't even log, it's a logging problem.
func (b *loggerBehavior) Recover(err interface{}) error {
	return nil
}

// EOF
