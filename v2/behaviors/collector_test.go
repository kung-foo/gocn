// Tideland Go Cell Network - Behaviors - Unit Tests - Collector
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

// TestCollectorBehavior tests the collector behavior.
func TestCollectorBehavior(t *testing.T) {
	assert := asserts.NewTestingAssertion(t, true)

	env := cells.NewEnvironment(cells.Id("collector-behavior"))
	defer env.Stop()

	env.StartCell("collector", behaviors.NewCollectorBehavior(10))

	for i := 0; i < 25; i++ {
		env.Raise("collector", "collect", i)
	}

	testsupport.LetItWork()

	var collected []behaviors.EventData
	err := env.Request("collector", "collected", nil, &collected)
	assert.Nil(err)
	assert.Length(collected, 10, "ten collected events")

	var ok bool
	err = env.Request("collector", "reset", nil, &ok)
	assert.Nil(err)
	assert.True(ok, "reset worked")

	err = env.Request("collector", "collected", nil, &collected)
	assert.Nil(err)
	assert.Empty(collected, "zero collected events")

	err = env.Request("collector", "someWeirdRequestThatDoesntExist", nil, nil)
	assert.NotNil(err)
}

// EOF
