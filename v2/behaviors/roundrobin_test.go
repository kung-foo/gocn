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

	"git.tideland.biz/gocn/behaviors"
	"git.tideland.biz/gocn/cells"
	"git.tideland.biz/gocn/testsupport"
	"git.tideland.biz/gots/asserts"
)

//--------------------
// TESTS
//--------------------

// TestRoundRobinBehavior tests the round robin behavior.
func TestRoundRobinBehavior(t *testing.T) {
	assert := asserts.NewTestingAssertion(t, true)

	env := cells.NewEnvironment(cells.Id("round-robin-behavior"))
	defer env.Stop()

	env.StartCell("round-robin", behaviors.NewRoundRobinBehavior())
	env.StartCell("test-1", testsupport.NewTestBehavior())
	env.StartCell("test-2", testsupport.NewTestBehavior())
	env.StartCell("test-3", testsupport.NewTestBehavior())
	env.StartCell("test-4", testsupport.NewTestBehavior())
	env.StartCell("test-5", testsupport.NewTestBehavior())
	env.Subscribe("round-robin", "test-1", "test-2", "test-3", "test-4", "test-5")

	// Just 23 to let 'test-4' and 'test-5' receive less events.
	for i := 0; i < 23; i++ {
		err := env.Raise("round-robin", "round", i)
		assert.Nil(err)
	}

	testsupport.LetItWork()

	test := func(id string, length int) {
		var processed []string
		err := env.Request(id, "processed", nil, &processed)
		assert.Nil(err)
		assert.Length(processed, length)
	}

	test("test-1", 5)
	test("test-2", 5)
	test("test-3", 5)
	test("test-4", 4)
	test("test-5", 4)
}

// EOF
