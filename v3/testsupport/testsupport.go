// Tideland Go Cell Network - Test Support
//
// Copyright (C) 2010-2014 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package testsupport

//--------------------
// IMPORTS
//--------------------

import (
	"fmt"
	"time"

	"github.com/tideland/gocn/v3/cells"
)

//--------------------
// HELPERS
//--------------------

// PanicTopic lets the test behavior panic to check recovering.
const PanicTopic = "panic!"

// testBehavior implements a simple behavior used in the tests.
type testBehavior struct {
	ctx         cells.Context
	processed   []string
	recoverings int
}

// NewTestBehavior creates a behavior for testing. It collects and
// re-emits all events, returns them with the topic "processed" and
// delets all collected with the topic "reset".
func NewTestBehavior() cells.Behavior {
	return &testBehavior{nil, []string{}, 0}
}

func (t *testBehavior) Init(ctx cells.Context) error {
	t.ctx = ctx
	return nil
}

func (t *testBehavior) Terminate() error {
	return nil
}

func (t *testBehavior) ProcessEvent(event cells.Event) error {
	switch event.Topic() {
	case cells.ProcessedTopic:
		processed := make([]string, len(t.processed))
		copy(processed, t.processed)
		err := event.Respond(processed)
		if err != nil {
			return err
		}
	case cells.ResetTopic:
		t.processed = []string{}
	case cells.PingTopic:
		err := event.Respond(cells.PongResponse)
		if err != nil {
			return err
		}
	case PanicTopic:
		panic("Ouch!")
	default:
		t.processed = append(t.processed, fmt.Sprintf("%v", event))
		t.ctx.Emit(event)
	}
	return nil
}

func (t *testBehavior) Recover(r interface{}) error {
	t.recoverings++
	if t.recoverings > 5 {
		return cells.NewCannotRecoverError(t.ctx.ID(), r)
	}
	return nil
}

// LetItWork just consumes time.
func LetItWork() {
	time.Sleep(100 * time.Millisecond)
}

// EOF
