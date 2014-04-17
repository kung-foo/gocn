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

	"github.com/tideland/gocn/v2/cells"
)

//--------------------
// HELPERS
//--------------------

// testBehavior implements a simple behavior used in the tests.
type testBehavior struct {
	ctx       cells.CellContext
	processed []string
}

// NewTestBehavior creates a behavior for testing. It collects and
// re-emits all events, returns them with the topic "processed" and
// delets all collected with the topic "reset".
func NewTestBehavior() cells.Behavior {
	return &testBehavior{nil, []string{}}
}

func (t *testBehavior) Init(ctx cells.CellContext) error {
	t.ctx = ctx
	return nil
}

func (t *testBehavior) Terminate() error {
	return nil
}

func (t *testBehavior) ProcessEvent(event cells.Event) error {
	if event.IsRequest() {
		switch event.Topic() {
		case "processed":
			return event.Respond(t.processed, nil)
		case "reset":
			return event.Respond(true, nil)
		}
		return fmt.Errorf("illegal request: %q", event.Topic())
	}
	t.processed = append(t.processed, event.String())
	t.ctx.EmitAll(event)
	return nil
}

func (t *testBehavior) Recover(r interface{}) error {
	return nil
}

// LetItWork just consumes time.
func LetItWork() {
	time.Sleep(100 * time.Millisecond)
}

// EOF
