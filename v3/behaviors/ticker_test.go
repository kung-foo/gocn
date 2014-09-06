// Tideland Go Cell Network - Behaviors - Unit Tests - Ticker
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
	"time"

	"github.com/tideland/gocn/v3/behaviors"
	"github.com/tideland/gocn/v3/cells"
	"github.com/tideland/gocn/v3/testsupport"
	"github.com/tideland/gots/v3/asserts"
)

//--------------------
// TESTS
//--------------------

// TestTickerBehavior tests the ticker behavior.
func TestTickerBehavior(t *testing.T) {
	assert := asserts.NewTestingAssertion(t, true)
	env := cells.NewEnvironment(cells.ID("ticker-behavior"))
	defer env.Stop()

	env.StartCell("ticker", behaviors.NewTickerBehavior(22*time.Millisecond))
	env.StartCell("test", testsupport.NewTestBehavior())
	env.Subscribe("ticker", "test")

	testsupport.LetItWork()

	processed, err := env.Request("test", cells.ProcessedTopic, nil, nil, cells.DefaultTimeout)
	assert.Nil(err)
	assert.Length(processed, 4)
}

// EOF
