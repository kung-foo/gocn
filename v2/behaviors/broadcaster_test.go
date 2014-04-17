// Tideland Go Cell Network - Behaviors - Unit Tests - Broadcaster
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

// TestBroadcasterBehavior tests the broadcast behavior.
func TestBroadcasterBehavior(t *testing.T) {
	assert := asserts.NewTestingAssertion(t, true)

	env := cells.NewEnvironment(cells.Id("broadcast-behavior"))
	defer env.Stop()

	env.StartCell("broadcast", behaviors.NewBroadcasterBehavior())
	env.StartCell("test-1", testsupport.NewTestBehavior())
	env.StartCell("test-2", testsupport.NewTestBehavior())
	env.Subscribe("broadcast", "test-1", "test-2")

	env.Raise("broadcast", "test", "a")
	env.Raise("broadcast", "test", "b")
	env.Raise("broadcast", "test", "c")

	testsupport.LetItWork()

	var processed []string
	err := env.Request("test-1", "processed", nil, &processed)
	assert.Nil(err)
	assert.Length(processed, 3)

	err = env.Request("test-2", "processed", nil, &processed)
	assert.Nil(err)
	assert.Length(processed, 3)
}

// EOF
