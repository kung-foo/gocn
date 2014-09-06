// Tideland Go Cell Network - Behaviors - Unit Tests - Router
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
	"strings"
	"testing"

	"github.com/tideland/gocn/v3/behaviors"
	"github.com/tideland/gocn/v3/cells"
	"github.com/tideland/gocn/v3/testsupport"
	"github.com/tideland/gots/v3/asserts"
)

//--------------------
// TESTS
//--------------------

// TestRouterBehavior tests the router behavior.
func TestRouterBehavior(t *testing.T) {
	assert := asserts.NewTestingAssertion(t, true)
	env := cells.NewEnvironment(cells.ID("router-behavior"))
	defer env.Stop()

	rf := func(id string, event cells.Event, subscribers []string) []string {
		return strings.Split(event.Topic(), ":")
	}
	env.StartCell("router", behaviors.NewRouterBehavior(rf))
	env.StartCell("test-1", testsupport.NewTestBehavior())
	env.StartCell("test-2", testsupport.NewTestBehavior())
	env.StartCell("test-3", testsupport.NewTestBehavior())
	env.StartCell("test-4", testsupport.NewTestBehavior())
	env.StartCell("test-5", testsupport.NewTestBehavior())
	env.Subscribe("router", "test-1", "test-2", "test-3", "test-4", "test-5")

	env.EmitNew("router", "test-1:test-2", "a", nil)
	env.EmitNew("router", "test-1:test-2:test-3", "b", nil)
	env.EmitNew("router", "test-3:test-4:test-5", "c", nil)

	testsupport.LetItWork()

	test := func(id string, length int) {
		processed, err := env.Request(id, cells.ProcessedTopic, nil, nil, cells.DefaultTimeout)
		assert.Nil(err)
		assert.Length(processed, length)
	}

	test("test-1", 2)
	test("test-2", 2)
	test("test-3", 2)
	test("test-4", 1)
	test("test-5", 1)
}

// EOF
