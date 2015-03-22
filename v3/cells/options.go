// Tideland Go Cell Network - Cells
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
	"github.com/tideland/goas/v2/identifier"
)

//--------------------
// OPTIONS
//--------------------

// Option allows to set an option of the environment.
type Option func(env Environment)

// Options is a set of options.
type Options []Option

// ID is the option to set the ID of the environment.
func ID(id string) Option {
	return func(env Environment) {
		e := env.(*environment)
		if id == "" {
			e.id = identifier.NewUUID().String()
		} else {
			e.id = id
		}
	}
}

// QueueFactory is the option to set the queue factory
// of the environment.
func QueueFactory(queueFactory EventQueueFactory) Option {
	return func(env Environment) {
		e := env.(*environment)
		if queueFactory == nil {
			queueFactory = makeLocalEventQueueFactory(10)
		}
		e.queueFactory = queueFactory
	}
}

// LocalQueueFactory is the option to set the queue factory of
// the environment to create local event queues.
func LocalQueueFactory(size int) Option {
	return func(env Environment) {
		if size < 5 {
			size = 5
		}
		QueueFactory(makeLocalEventQueueFactory(size))
	}
}

// EOF
