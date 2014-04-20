// Tideland Go Cell Network - Behaviors - Unit Tests - Mapper
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

	"github.com/tideland/gocn/v2/behaviors"
	"github.com/tideland/gocn/v2/cells"
	"github.com/tideland/gocn/v2/testsupport"
	"github.com/tideland/gots/v3/asserts"
)

//--------------------
// TESTS
//--------------------

// TestMapperBehavior tests the mapping of events.
func TestMapperBehavior(t *testing.T) {
	assert := asserts.NewTestingAssertion(t, true)
	assertPayload := func(data behaviors.EventData, value string) {
		payload, ok := data.Payload.(string)
		assert.True(ok)
		assert.Equal(payload, value)
	}

	env := cells.NewEnvironment(cells.Id("mapper-behavior"))
	defer env.Stop()

	mf := func(id string, event cells.Event) (string, interface{}) {
		var payload string
		err := event.Payload(&payload)
		assert.Nil(err)
		return event.Topic(), strings.ToUpper(payload)
	}

	env.StartCell("map", behaviors.NewMapperBehavior(mf))
	env.StartCell("collect", behaviors.NewCollectorBehavior(10))
	env.Subscribe("map", "collect")

	env.Raise("map", "a", "abc")
	env.Raise("map", "b", "def")
	env.Raise("map", "c", "ghi")

	testsupport.LetItWork()

	var collected []behaviors.EventData
	err := env.Request("collect", "collected", nil, &collected)
	assert.Nil(err)
	assert.Length(collected, 3, "three mapped events")
	assertPayload(collected[0], "ABC")
	assertPayload(collected[1], "DEF")
	assertPayload(collected[2], "GHI")
}

// EOF
