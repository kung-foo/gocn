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

	"git.tideland.biz/gocn/behaviors"
	"git.tideland.biz/gocn/cells"
	"git.tideland.biz/gocn/testsupport"
	"git.tideland.biz/gots/asserts"
)

//--------------------
// TESTS
//--------------------

// TestTickerBehavior tests the ticker behavior.
func TestTickerBehavior(t *testing.T) {
	assert := asserts.NewTestingAssertion(t, true)

	env := cells.NewEnvironment(cells.Id("ticker-behavior"))
	defer env.Stop()

	env.StartCell("ticker", behaviors.NewTickerBehavior(25*time.Millisecond))
	env.StartCell("test", testsupport.NewTestBehavior())
	env.Subscribe("ticker", "test")

	time.Sleep(110 * time.Millisecond)

	var processed []string
	err := env.Request("test", "processed", nil, &processed)
	assert.Nil(err)
	assert.Length(processed, 4)
}

// EOF
