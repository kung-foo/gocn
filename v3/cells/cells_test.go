// Tideland Go Cell Network - Cells - Unit Tests
//
// Copyright (C) 2010-2014 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package cells_test

//--------------------
// IMPORTS
//--------------------

import (
	"fmt"
	"testing"
	"time"

	"github.com/tideland/goas/v3/errors"
	"github.com/tideland/gocn/v3/cells"
	"github.com/tideland/gocn/v3/testsupport"
	"github.com/tideland/gots/v3/asserts"
)

//--------------------
// TESTS
//--------------------

// TestEvent tests the event construction.
func TestEvent(t *testing.T) {
	assert := asserts.NewTestingAssertion(t, true)

	event, err := cells.NewEvent("foo", "bar", nil)
	assert.Nil(err)
	assert.Equal(event.Topic(), "foo")
	assert.Equal(event.String(), "<event: \"foo\" / payload: <\"default\": bar>>")

	bar, ok := event.Payload().Get(cells.DefaultPayload)
	assert.True(ok)
	assert.Equal(bar, "bar")

	_, err = cells.NewEvent("", nil, nil)
	assert.True(cells.IsNoTopicError(err))

	_, err = cells.NewEvent("yadda", nil, nil)
	assert.Nil(err)
}

// TestRingBuffer tests the buffer used by the local event queue.
func TestRingBuffer(t *testing.T) {
	assert := asserts.NewTestingAssertion(t, true)
	buffer := cells.NewRingBuffer(5)

	// Initial empty buffer.
	assert.Nil(buffer.Peek())
	assert.Equal(buffer.Cap(), 5)

	// Push first event.
	event, err := cells.NewEvent("buffer-test-1", "foo", nil)
	assert.Nil(err)
	buffer.Push(event)
	assert.Equal(buffer.Peek().Topic(), "buffer-test-1")
	assert.Length(buffer, 1)

	// Fill buffer.
	for i := 2; i < 6; i++ {
		topic := fmt.Sprintf("buffer-test-%d", i)
		event, err := cells.NewEvent(topic, "foo", nil)
		assert.Nil(err)
		buffer.Push(event)
	}
	assert.Length(buffer, 5)
	assert.Equal(buffer.Cap(), 5)

	// Add one to increase the buffer.
	event, err = cells.NewEvent("buffer-test-6", "foo", nil)
	assert.Nil(err)
	buffer.Push(event)
	assert.Length(buffer, 6)
	assert.Equal(buffer.Cap(), 6)

	// Pop the first event.
	buffer.Pop()
	assert.Length(buffer, 5)
	assert.Equal(buffer.Cap(), 6)

	// Add one w/o increasing the buffer.
	event, err = cells.NewEvent("buffer-test-7", "foo", nil)
	assert.Nil(err)
	buffer.Push(event)
	assert.Length(buffer, 6)
	assert.Equal(buffer.Cap(), 6)

	// Add one to increase the buffer.
	event, err = cells.NewEvent("buffer-test-8", "foo", nil)
	assert.Nil(err)
	buffer.Push(event)
	assert.Length(buffer, 7)
	assert.Equal(buffer.Cap(), 7)

	// Peek again.
	assert.Equal(buffer.Peek().Topic(), "buffer-test-2")

	// Pop it almost empty.
	for i := 0; i < 6; i++ {
		buffer.Pop()
	}
	assert.Length(buffer, 1)
	assert.Equal(buffer.Cap(), 7)

	// Now pop it empty.
	buffer.Pop()
	assert.Length(buffer, 0)
	assert.Equal(buffer.Cap(), 7)
	assert.Nil(buffer.Peek())

	// One final push to empty buffer.
	event, err = cells.NewEvent("buffer-test-9", "foo", nil)
	assert.Nil(err)
	buffer.Push(event)
	assert.Length(buffer, 1)
	assert.Equal(buffer.Cap(), 7)
}

// TestLocalEventQueue tests the local event queue.
func TestLocalEventQueue(t *testing.T) {
	count := 10000
	assert := asserts.NewTestingAssertion(t, true)
	factory := cells.MakeLocalEventQueueFactory(10)
	queue, err := factory(nil)
	assert.Nil(err)

	for i := 0; i < count; i++ {
		event, err := cells.NewEvent("queue-test", i, nil)
		assert.Nil(err)

		assert.Nil(queue.Push(event))
	}

	for i := 0; i < count; i++ {
		_, ok := <-queue.Events()
		assert.True(ok)
	}

	select {
	case <-queue.Events():
		assert.Fail("didn't expected any queued event")
	case <-time.After(100 * time.Millisecond):
		assert.True(true)
	}

	assert.Nil(queue.Stop())
}

// BenchmarkLocalEventQueue tests the performance of the local event queue.
func BenchmarkLocalEventQueue(b *testing.B) {
	factory := cells.MakeLocalEventQueueFactory(10)
	queue, err := factory(nil)
	if err != nil {
		b.Fatalf("cannot create queue: %v", err)
	}

	for i := 0; i < b.N; i++ {
		event, _ := cells.NewEvent("queue-test", i, nil)
		queue.Push(event)
	}

	for i := 0; i < b.N; i++ {
		<-queue.Events()
	}

	queue.Stop()
}

// TestEnvironmentAddRemove tests adding, checking and
// removing of cells.
func TestEnvironmentStartStop(t *testing.T) {
	assert := asserts.NewTestingAssertion(t, true)

	env := cells.NewEnvironment()
	defer env.Stop()

	err := env.StartCell("foo", testsupport.NewTestBehavior())
	assert.Nil(err)

	hasFoo := env.HasCell("foo")
	assert.True(hasFoo)

	env.StopCell("foo")
	hasFoo = env.HasCell("foo")
	assert.False(hasFoo)

	hasBar := env.HasCell("bar")
	assert.False(hasBar)
	env.StopCell("bar")
	hasBar = env.HasCell("bar")
	assert.False(hasBar)
}

// TestEnvironmentSubscribeUnsubscribe tests subscribing,
// checking and unsubscribing of cells.
func TestEnvironmentSubscribeUnsubscribe(t *testing.T) {
	assert := asserts.NewTestingAssertion(t, true)

	env := cells.NewEnvironment()
	defer env.Stop()

	err := env.StartCell("foo", testsupport.NewTestBehavior())
	assert.Nil(err)
	err = env.StartCell("bar", testsupport.NewTestBehavior())
	assert.Nil(err)
	err = env.StartCell("baz", testsupport.NewTestBehavior())
	assert.Nil(err)
	err = env.StartCell("yadda", testsupport.NewTestBehavior())
	assert.Nil(err)

	err = env.Subscribe("humpf", "foo")
	assert.True(errors.IsError(err, cells.ErrInvalidID))
	err = env.Subscribe("foo", "humpf")
	assert.True(errors.IsError(err, cells.ErrInvalidID))

	err = env.Subscribe("foo", "bar", "baz")
	assert.Nil(err)
	subs, err := env.Subscribers("foo")
	assert.Nil(err)
	assert.Contents("bar", subs)
	assert.Contents("baz", subs)

	err = env.Unsubscribe("foo", "bar")
	assert.Nil(err)
	subs, err = env.Subscribers("foo")
	assert.Nil(err)
	assert.Contents("baz", subs)

	err = env.Unsubscribe("foo", "baz")
	assert.Nil(err)
	subs, err = env.Subscribers("foo")
	assert.Nil(err)
	assert.Empty(subs)
}

// TestEnvironmentScenario tests creating and using the
// environment in a simple way.
func TestEnvironmentScenario(t *testing.T) {
	assert := asserts.NewTestingAssertion(t, true)
	env := cells.NewEnvironment(cells.ID("scenario"))
	defer env.Stop()

	err := env.StartCell("foo", testsupport.NewTestBehavior())
	assert.Nil(err)
	err = env.StartCell("bar", testsupport.NewTestBehavior())
	assert.Nil(err)
	err = env.StartCell("collector", testsupport.NewTestBehavior())
	assert.Nil(err)

	err = env.Subscribe("foo", "bar")
	assert.Nil(err)
	err = env.Subscribe("bar", "collector")
	assert.Nil(err)

	err = env.EmitNew("foo", "lorem", 4711, nil)
	assert.Nil(err)
	err = env.EmitNew("foo", "ipsum", 1234, nil)
	assert.Nil(err)
	response, err := env.Request("foo", cells.PingTopic, nil, nil, cells.DefaultTimeout)
	assert.Nil(err)
	assert.Equal(response, cells.PongResponse)

	collected, err := env.Request("collector", cells.ProcessedTopic, nil, nil, cells.DefaultTimeout)
	assert.Nil(err)
	assert.Length(collected, 2, "two collected events")
	assert.Contents(`<event: "lorem" / payload: <"default": 4711>>`, collected)
	assert.Contents(`<event: "ipsum" / payload: <"default": 1234>>`, collected)
}

// EOF
