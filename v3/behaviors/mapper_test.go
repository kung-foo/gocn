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

	"github.com/tideland/gocn/v3/behaviors"
	"github.com/tideland/gocn/v3/cells"
	"github.com/tideland/gocn/v3/testsupport"
	"github.com/tideland/gots/v3/asserts"
)

//--------------------
// TESTS
//--------------------

// TestMapperBehavior tests the mapping of events.
func TestMapperBehavior(t *testing.T) {
	assert := asserts.NewTestingAssertion(t, true)
	assertPayload := func(collected interface{}, index int, value string) {
		eventData, ok := collected.([]behaviors.EventData)
		assert.True(ok)
		payload, ok := eventData[index].Payload.(cells.Payload)
		assert.True(ok)
		upperText, ok := payload.Get("upper-text")
		assert.True(ok)
		assert.Equal(upperText, value)
	}
	env := cells.NewEnvironment(cells.ID("mapper-behavior"))
	defer env.Stop()

	mf := func(id string, event cells.Event) (cells.Event, error) {
		text, ok := event.Payload().Get(cells.DefaultPayload)
		if !ok {
			return event, nil
		}
		pv := cells.PayloadValues{
			"upper-text": strings.ToUpper(text.(string)),
		}
		payload := event.Payload().Apply(pv)
		return cells.NewEvent(event.Topic(), payload, event.Scene())
	}

	env.StartCell("mapper", behaviors.NewMapperBehavior(mf))
	env.StartCell("collector", behaviors.NewCollectorBehavior(10))
	env.Subscribe("mapper", "collector")

	env.EmitNew("mapper", "a", "abc", nil)
	env.EmitNew("mapper", "b", "def", nil)
	env.EmitNew("mapper", "c", "ghi", nil)

	testsupport.LetItWork()

	collected, err := env.Request("collector", cells.CollectedTopic, nil, nil, cells.DefaultTimeout)
	assert.Nil(err)
	assert.Length(collected, 3, "three mapped events")
	assertPayload(collected, 0, "ABC")
	assertPayload(collected, 1, "DEF")
	assertPayload(collected, 2, "GHI")
}

// EOF
