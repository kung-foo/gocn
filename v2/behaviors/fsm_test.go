// Tideland Go Cell Network - Behaviors - Unit Tests - Finite State Machine
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
	"fmt"
	"testing"

	"git.tideland.biz/gocn/behaviors"
	"git.tideland.biz/gocn/cells"
	"git.tideland.biz/gocn/testsupport"
	"git.tideland.biz/gots/asserts"
)

//--------------------
// TESTS
//--------------------

// TestFSMBehavior tests the finite state machine behavior.
func TestFSMBehavior(t *testing.T) {
	assert := asserts.NewTestingAssertion(t, true)

	env := cells.NewEnvironment(cells.Id("fsm-behavior"))
	defer env.Stop()

	checkCents := func(id string) int {
		var cents int
		err := env.Request(id, "cents", nil, &cents)
		assert.Nil(err)
		return cents
	}
	info := func(id string) string {
		var info string
		err := env.Request(id, "info", nil, &info)
		assert.Nil(err)
		return info
	}
	grabCents := func() int {
		var cents int
		err := env.Request("restorer", "grab", nil, &cents)
		assert.Nil(err)
		return cents
	}

	lockA := lockMachine{}
	lockB := lockMachine{}

	env.StartCell("lock-a", behaviors.NewFSMBehavior(lockA.locked))
	env.StartCell("lock-b", behaviors.NewFSMBehavior(lockB.locked))
	env.StartCell("restorer", newRestorerBehavior())

	env.Subscribe("lock-a", "restorer")
	env.Subscribe("lock-b", "restorer")

	// 1st run: emit not enough and press button.
	env.Raise("lock-a", "coin", 20)
	env.Raise("lock-a", "coin", 20)
	env.Raise("lock-a", "coin", 20)
	env.Raise("lock-a", "button-press", nil)

	testsupport.LetItWork()

	assert.Equal(checkCents("lock-a"), 0)
	assert.Equal(grabCents(), 60)

	// 2nd run: unlock the lock and lock it again.
	env.Raise("lock-a", "coin", 50)
	env.Raise("lock-a", "coin", 20)
	env.Raise("lock-a", "coin", 50)

	testsupport.LetItWork()

	assert.Equal(info("lock-a"), "state 'unlocked' with 20 cents")

	env.Raise("lock-a", "button-press", nil)

	testsupport.LetItWork()

	assert.Equal(checkCents("lock-a"), 00)
	assert.Equal(info("lock-a"), "state 'locked' with 0 cents")
	assert.Equal(grabCents(), 20)

	// 3rd run: put a screwdriwer in the lock.
	env.Raise("lock-a", "screwdriver", nil)

	testsupport.LetItWork()

	ok, err := behaviors.IsOkEndState(env, "lock-a")
	assert.Nil(err)
	assert.True(ok, "ok maybe the wrong term ;)")

	// 4th run: try an illegal action.
	env.Raise("lock-b", "chewing-gum", nil)

	testsupport.LetItWork()

	ok, err = behaviors.IsErrorEndState(env, "lock-b")
	assert.True(ok)
	assert.ErrorMatch(err, "illegal topic in state 'locked': chewing-gum")
}

//--------------------
// HELPERS
//--------------------

// cents retrieves the cents out of the payload of an event.
func payloadCents(event cells.Event) int {
	var cents int
	event.Payload(&cents)
	return cents
}

// lockMachine will be unlocked if enough money is inserted. After
// that it can be locked again.
type lockMachine struct {
	cents int
}

// locked represents the locked state receiving coins.
func (m *lockMachine) locked(ctx cells.CellContext, event cells.Event) (behaviors.State, error) {
	if event.IsRequest() {
		switch event.Topic() {
		case "cents":
			return m.locked, event.Respond(m.cents, nil)
		case "info":
			info := fmt.Sprintf("state 'locked' with %d cents", m.cents)
			return m.locked, event.Respond(info, nil)
		}
		err := fmt.Errorf("illegal request in state 'locked': %v", event)
		return m.locked, event.Respond(nil, err)
	}
	switch event.Topic() {
	case "coin":
		cents := payloadCents(event)
		if cents < 1 {
			return nil, fmt.Errorf("do not insert buttons")
		}
		m.cents += cents
		if m.cents > 100 {
			m.cents -= 100
			return m.unlocked, nil
		}
		return m.locked, nil
	case "button-press":
		if m.cents > 0 {
			ctx.Raise("restorer", "drop", m.cents)
			m.cents = 0
		}
		return m.locked, nil
	case "screwdriver":
		// Allow a screwdriver to bring the lock into an undefined state.
		return nil, nil
	}
	return m.locked, fmt.Errorf("illegal topic in state 'locked': %s", event.Topic())
}

// unlocked represents the unlocked state receiving coins.
func (m *lockMachine) unlocked(ctx cells.CellContext, event cells.Event) (behaviors.State, error) {
	if event.IsRequest() {
		switch event.Topic() {
		case "cents":
			return m.locked, event.Respond(m.cents, nil)
		case "info":
			info := fmt.Sprintf("state 'unlocked' with %d cents", m.cents)
			return m.unlocked, event.Respond(info, nil)
		}
		err := fmt.Errorf("illegal request in state 'unlocked': %v", event)
		return m.unlocked, event.Respond(nil, err)
	}
	switch event.Topic() {
	case "coin":
		cents := payloadCents(event)
		ctx.RaiseAll("return", cents)
		return m.unlocked, nil
	case "button-press":
		ctx.Raise("restorer", "drop", m.cents)
		m.cents = 0
		return m.locked, nil
	}
	return m.unlocked, fmt.Errorf("illegal topic in state 'unlocked': %s", event.Topic())
}

type restorerBehavior struct {
	ctx   cells.CellContext
	cents int
}

func newRestorerBehavior() cells.Behavior {
	return &restorerBehavior{nil, 0}
}

func (b *restorerBehavior) Init(ctx cells.CellContext) error {
	b.ctx = ctx
	return nil
}

func (b *restorerBehavior) Terminate() error {
	return nil
}

func (b *restorerBehavior) ProcessEvent(event cells.Event) error {
	if event.IsRequest() {
		switch event.Topic() {
		case "grab":
			cents := b.cents
			b.cents = 0
			return event.Respond(cents, nil)
		}
		return nil
	}
	switch event.Topic() {
	case "drop":
		b.cents += payloadCents(event)
	}
	return nil
}

func (b *restorerBehavior) Recover(err interface{}) error {
	return nil
}

// EOF
