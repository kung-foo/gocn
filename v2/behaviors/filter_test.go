// Tideland Go Cell Network - Behaviors - Unit Tests - Filter
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

// TestFilterBehavior tests the filter behavior.
func TestFilterBehavior(t *testing.T) {
	assert := asserts.NewTestingAssertion(t, true)

	env := cells.NewEnvironment(cells.Id("filter-behavior"))
	defer env.Stop()

	ff := func(id string, event cells.Event) bool {
		var payload string
		err := event.Payload(&payload)
		assert.Nil(err)
		return event.Topic() == payload
	}
	env.StartCell("filter", behaviors.NewFilterBehavior(ff))
	env.StartCell("test", testsupport.NewTestBehavior())
	env.Subscribe("filter", "test")

	env.Raise("filter", "a", "a")
	env.Raise("filter", "a", "b")
	env.Raise("filter", "b", "b")

	testsupport.LetItWork()

	var processed []string
	err := env.Request("test", "processed", nil, &processed)
	assert.Nil(err)
	assert.Length(processed, 2, "two filtered events")
}

// EOF
