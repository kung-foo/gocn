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

	"github.com/tideland/gocn/v3/behaviors"
	"github.com/tideland/gocn/v3/cells"
	"github.com/tideland/gocn/v3/testsupport"
	"github.com/tideland/gots/v3/asserts"
)

//--------------------
// TESTS
//--------------------

// TestCollectorBehavior tests the collector behavior.
func TestCollectorBehavior(t *testing.T) {
	assert := asserts.NewTestingAssertion(t, true)
	env := cells.NewEnvironment(cells.ID("collector-behavior"))
	defer env.Stop()

	env.StartCell("collector", behaviors.NewCollectorBehavior(10))

	for i := 0; i < 25; i++ {
		env.EmitNew("collector", "collect", i, nil)
	}

	testsupport.LetItWork()

	collected, err := env.Request("collector", cells.CollectedTopic, nil, nil, cells.DefaultTimeout)
	assert.Nil(err)
	assert.Length(collected, 10)

	err = env.EmitNew("collector", cells.ResetTopic, nil, nil)
	assert.Nil(err)

	collected, err = env.Request("collector", cells.CollectedTopic, nil, nil, cells.DefaultTimeout)
	assert.Nil(err)
	assert.Length(collected, 0)
}

// EOF
