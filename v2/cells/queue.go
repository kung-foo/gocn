// Tideland Go Cell Network - Cells - Queue
//
// Copyright (C) 2010-2014 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package cells

//--------------------
// IMPORTS
//--------------------

import (
	"github.com/tideland/goas/v2/loop"
)

//--------------------
// EVENT QUEUE
//--------------------

// EventQueueFactory describes a function returning individual
// implementations of the EventQueue interface. This way different
// types of event queues can be injected into environments.
type EventQueueFactory func(env Environment) (EventQueue, error)

// EventQueue describes the methods any queue implementation
// must provide.
type EventQueue interface {
	// Push appends an event to the end of the queue.
	Push(event Event) error

	// Events returns a channel delivering the event
	// fromt the beginning of the queue.
	Events() <-chan Event

	// Stop tells the queue to end working.
	Stop() error
}

//--------------------
// LOCAL EVENT QUEUE
//--------------------

// ringBuffer
type ringBuffer struct {
	events []Event
	start  int
	end    int
}

// newRingBuffer
func newRingBuffer(size int) *ringBuffer {
	return &ringBuffer{
		events: make([]Event, size),
		start:  0,
		end:    0,
	}
}

// peek returns the first event of the queue.
func (rb *ringBuffer) peek() Event {
	if rb.start == rb.end {
		return nil
	}
	return rb.events[rb.start]
}

// push adds an event to the end of the queue.
func (rb *ringBuffer) push(event Event) {
	rb.events[rb.end] = event
	rb.end++
	if rb.end == cap(rb.events) {
		rb.end = 0
	}
	if rb.end == rb.start {
		// Buffer full, so resize.
		tmp := make([]Event, cap(rb.events)*2)
		copy(tmp[0:], rb.events[rb.start:])
		if rb.end > 0 {
			copy(tmp[rb.end-1:], rb.events[:rb.end])
		}
		rb.start = 0
		rb.end = cap(rb.events)
		rb.events = tmp
	}
}

// pop removes the first event of the queue.
func (rb *ringBuffer) pop() {
	if rb.start == rb.end {
		return
	}
	rb.start++
	if rb.start == cap(rb.events) {
		rb.start = 0
	}
}

// localEventQueue implements a local in-memory event queue.
type localEventQueue struct {
	buffer *ringBuffer
	pushc  chan Event
	eventc chan Event
	loop   loop.Loop
}

// MakeLocalEventQueueFactory creates a factory for local
// event queues.
func MakeLocalEventQueueFactory(size int) EventQueueFactory {
	return func(env Environment) (EventQueue, error) {
		queue := &localEventQueue{
			buffer: newRingBuffer(size),
			pushc:  make(chan Event),
			eventc: make(chan Event),
		}
		queue.loop = loop.Go(queue.backendLoop)
		return queue, nil
	}
}

// Push appends an event to the end of the queue.
func (q *localEventQueue) Push(event Event) error {
	select {
	case q.pushc <- event:
	case <-q.loop.IsStopping():
		return newError(ErrStopping, "event queue")
	}
	return nil
}

// Events returns a channel delivering the event
// fromt the beginning of the queue.
func (q *localEventQueue) Events() <-chan Event {
	return q.eventc
}

// Stop tells the queue to end working.
func (l *localEventQueue) Stop() error {
	return l.loop.Stop()
}

// backendLoop realizes the backend of the queue.
func (q *localEventQueue) backendLoop(l loop.Loop) error {
	for {
		if q.buffer.peek() == nil {
			select {
			case <-q.loop.ShallStop():
				return nil
			case event := <-q.pushc:
				q.buffer.push(event)
			}
		}

		select {
		case <-q.loop.ShallStop():
			return nil
		case event := <-q.pushc:
			q.buffer.push(event)
		case q.eventc <- q.buffer.peek():
			q.buffer.pop()
		}
	}
}

// EOF
