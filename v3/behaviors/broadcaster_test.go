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

	"github.com/tideland/gocn/v3/behaviors"
	"github.com/tideland/gocn/v3/cells"
	"github.com/tideland/gocn/v3/testsupport"
	"github.com/tideland/gots/v3/asserts"
)

//--------------------
// TESTS
//--------------------

// TestBroadcasterBehavior tests the broadcast behavior.
func TestBroadcasterBehavior(t *testing.T) {
	assert := asserts.NewTestingAssertion(t, true)
	env := cells.NewEnvironment(cells.ID("broadcast-behavior"))
	defer env.Stop()

	env.StartCell("broadcast", behaviors.NewBroadcasterBehavior())
	env.StartCell("test-1", testsupport.NewTestBehavior())
	env.StartCell("test-2", testsupport.NewTestBehavior())
	env.Subscribe("broadcast", "test-1", "test-2")

	env.EmitNew("broadcast", "test", "a", nil)
	env.EmitNew("broadcast", "test", "b", nil)
	env.EmitNew("broadcast", "test", "c", nil)

	testsupport.LetItWork()

	processed, err := env.Request("test-1", cells.ProcessedTopic, nil, nil, cells.DefaultTimeout)
	assert.Nil(err)
	assert.Length(processed, 3)

	processed, err = env.Request("test-2", cells.ProcessedTopic, nil, nil, cells.DefaultTimeout)
	assert.Nil(err)
	assert.Length(processed, 3)
}

// EOF
