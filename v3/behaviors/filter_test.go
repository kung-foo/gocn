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

	"github.com/tideland/gocn/v3/behaviors"
	"github.com/tideland/gocn/v3/cells"
	"github.com/tideland/gocn/v3/testsupport"
	"github.com/tideland/gots/v3/asserts"
)

//--------------------
// TESTS
//--------------------

// TestFilterBehavior tests the filter behavior.
func TestFilterBehavior(t *testing.T) {
	assert := asserts.NewTestingAssertion(t, true)
	env := cells.NewEnvironment(cells.ID("filter-behavior"))
	defer env.Stop()

	ff := func(id string, event cells.Event) bool {
		dp, ok := event.Payload().Get(cells.DefaultPayload)
		if !ok {
			return false
		}
		payload := dp.(string)
		return event.Topic() == payload
	}
	env.StartCell("filter", behaviors.NewFilterBehavior(ff))
	env.StartCell("collector", behaviors.NewCollectorBehavior(10))
	env.Subscribe("filter", "collector")

	env.EmitNew("filter", "a", "a", nil)
	env.EmitNew("filter", "a", "b", nil)
	env.EmitNew("filter", "b", "b", nil)

	testsupport.LetItWork()

	collected, err := env.Request("collector", cells.CollectedTopic, nil, nil, cells.DefaultTimeout)
	assert.Nil(err)
	assert.Length(collected, 2, "two collected events")
}

// EOF
