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

	"github.com/tideland/gocn/v2/behaviors"
	"github.com/tideland/gocn/v2/cells"
	"github.com/tideland/gocn/v2/testsupport"
	"github.com/tideland/gots/v3/asserts"
)

//--------------------
// TESTS
//--------------------

// TestCounterBehavior tests the counting of events.
func TestCounterBehavior(t *testing.T) {
	assert := asserts.NewTestingAssertion(t, true)

	env := cells.NewEnvironment(cells.Id("counter-behavior"))
	defer env.Stop()

	cf := func(id string, event cells.Event) []string {
		var payload []string
		err := event.Payload(&payload)
		assert.Nil(err)
		return payload
	}
	env.StartCell("counter", behaviors.NewCounterBehavior(cf))

	env.Raise("counter", "count", []string{"a", "b"})
	env.Raise("counter", "count", []string{"a", "c", "d"})
	env.Raise("counter", "count", []string{"a", "d"})

	testsupport.LetItWork()

	var counters behaviors.Counters
	err := env.Request("counter", "counters", nil, &counters)
	assert.Nil(err)
	assert.Length(counters, 4, "four counted events")

	assert.Equal(counters["a"], int64(3))
	assert.Equal(counters["b"], int64(1))
	assert.Equal(counters["c"], int64(1))
	assert.Equal(counters["d"], int64(2))

	var ok bool
	err = env.Request("counter", "reset", nil, &ok)
	assert.Nil(err)
	assert.True(ok)

	counters = nil
	err = env.Request("counter", "counters", nil, &counters)
	assert.Nil(err)
	assert.Empty(counters, "zero counted events")
}

// EOF
