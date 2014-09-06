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
	"github.com/tideland/goas/v2/identifier"
)

//--------------------
// OPTIONS
//--------------------

// Option allows to set an option of the environment. It returns
// an option to reset it.
type Option func(env Environment) Option

// Options is a set of options.
type Options []Option

// ID is the option to set the ID of the environment.
func ID(id string) Option {
	return func(env Environment) Option {
		e := env.(*environment)
		previous := e.id
		if id == "" {
			e.id = identifier.NewUUID().String()
		} else {
			e.id = id
		}
		return ID(previous)
	}
}

// QueueFactory is the option to set the queue factory
// of the environment.
func QueueFactory(queueFactory EventQueueFactory) Option {
	return func(env Environment) Option {
		e := env.(*environment)
		previous := e.queueFactory
		if queueFactory == nil {
			queueFactory = MakeLocalEventQueueFactory(10)
		}
		e.queueFactory = queueFactory
		return QueueFactory(previous)
	}
}

// LocalQueueFactory is the option to set the queue factory of
// the environment to create local event queues.
func LocalQueueFactory(size int) Option {
	return func(env Environment) Option {
		if size < 5 {
			size = 5
		}
		return QueueFactory(MakeLocalEventQueueFactory(size))
	}
}

// EOF
