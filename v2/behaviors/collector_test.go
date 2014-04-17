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

	"git.tideland.biz/gocn/behaviors"
	"git.tideland.biz/gocn/cells"
	"git.tideland.biz/gots/asserts"
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
	t.Logf("response: %v", err)
}

// EOF
