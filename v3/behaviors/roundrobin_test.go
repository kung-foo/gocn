// Tideland Go Cell Network - Behaviors - Unit Tests
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

// TestRoundRobinBehavior tests the round robin behavior.
func TestRoundRobinBehavior(t *testing.T) {
	assert := asserts.NewTestingAssertion(t, true)
	env := cells.NewEnvironment(cells.ID("round-robin-behavior"))
	defer env.Stop()

	env.StartCell("round-robin", behaviors.NewRoundRobinBehavior())
	env.StartCell("round-robin-1", testsupport.NewTestBehavior())
	env.StartCell("round-robin-2", testsupport.NewTestBehavior())
	env.StartCell("round-robin-3", testsupport.NewTestBehavior())
	env.StartCell("round-robin-4", testsupport.NewTestBehavior())
	env.StartCell("round-robin-5", testsupport.NewTestBehavior())
	env.Subscribe("round-robin", "round-robin-1", "round-robin-2", "round-robin-3", "round-robin-4", "round-robin-5")

	// Just 23 to let 'round-robin-4' and 'round-robin-5' receive less events.
	for i := 0; i < 23; i++ {
		err := env.EmitNew("round-robin", "round", i, nil)
		assert.Nil(err)
	}

	testsupport.LetItWork()
	testsupport.LetItWork()

	test := func(id string, length int) {
		processed, err := env.Request(id, cells.ProcessedTopic, nil, nil, cells.DefaultTimeout)
		assert.Nil(err)
		assert.Length(processed, length)
	}

	test("round-robin-1", 5)
	test("round-robin-2", 5)
	test("round-robin-3", 5)
	test("round-robin-4", 4)
	test("round-robin-5", 4)
}

// EOF
