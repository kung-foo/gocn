// Tideland Go Cell Network - Behaviors - Unit Tests - Counter
//
// Copyright (C) 2010-2014 Frank Mueller / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package behaviors_test

//--------------------
// IMPORTS
//--------------------

import (
	"testing"

	"github.com/tideland/gocn/v3/behaviors"
	"github.com/tideland/gocn/v3/cells"
	"github.com/tideland/gocn/v3/testsupport"
	"github.com/tideland/gots/v3/asserts"
)

//--------------------
// TESTS
//--------------------

// TestCounterBehavior tests the counting of events.
func TestCounterBehavior(t *testing.T) {
	assert := asserts.NewTestingAssertion(t, true)
	env := cells.NewEnvironment(cells.ID("counter-behavior"))
	defer env.Stop()

	cf := func(id string, event cells.Event) []string {
		payload, ok := event.Payload().Get(cells.DefaultPayload)
		if !ok {
			return []string{}
		}
		return payload.([]string)
	}
	env.StartCell("counter", behaviors.NewCounterBehavior(cf))

	env.EmitNew("counter", "count", []string{"a", "b"}, nil)
	env.EmitNew("counter", "count", []string{"a", "c", "d"}, nil)
	env.EmitNew("counter", "count", []string{"a", "d"}, nil)

	testsupport.LetItWork()

	counters, err := env.Request("counter", cells.CountersTopic, nil, nil, cells.DefaultTimeout)
	assert.Nil(err)
	assert.Length(counters, 4, "four counted events")

	c := counters.(behaviors.Counters)

	assert.Equal(c["a"], int64(3))
	assert.Equal(c["b"], int64(1))
	assert.Equal(c["c"], int64(1))
	assert.Equal(c["d"], int64(2))

	err = env.EmitNew("counter", cells.ResetTopic, nil, nil)
	assert.Nil(err)

	counters, err = env.Request("counter", cells.CountersTopic, nil, nil, cells.DefaultTimeout)
	assert.Nil(err)
	assert.Empty(counters, "zero counted events")
}

// EOF
