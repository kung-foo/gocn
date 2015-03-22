// Tideland Go Cell Network - Cells - Queue
//
// Copyright (C) 2010-2015 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package cells

//--------------------
// IMPORTS
//--------------------

import (
	"github.com/tideland/goas/v2/loop"
	"github.com/tideland/goas/v3/errors"
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

// eventLink stores an event and a link to the next event.
type eventLink struct {
	event Event
	next  *eventLink
}

// ringBuffer stores events in a linked ring.
type ringBuffer struct {
	start   *eventLink
	end     *eventLink
	current *eventLink
}

// newRingBuffer creates a new ring buffer.
func newRingBuffer(size int) *ringBuffer {
	rb := &ringBuffer{}
	rb.start = &eventLink{}
	rb.end = rb.start
	for i := 0; i < size-1; i++ {
		link := &eventLink{}
		rb.end.next = link
		rb.end = link
	}
	rb.end.next = rb.start
	return rb
}

// Len returns the number of events in the buffer.
func (rb *ringBuffer) Len() int {
	l := 0
	current := rb.start
	for current.event != nil {
		l++
		current = current.next
		if current == rb.start {
			break
		}
	}
	return l
}

// Cap returns the capacity of the buffer.
func (rb *ringBuffer) Cap() int {
	c := 1
	current := rb.start
	for current.next != rb.start {
		c++
		current = current.next
	}
	return c
}

// Peek returns the first event of the queue.
func (rb *ringBuffer) Peek() Event {
	return rb.start.event
}

// Push adds an event to the end of the queue.
func (rb *ringBuffer) Push(event Event) {
	if rb.end.next.event == nil {
		rb.end.next.event = event
		rb.end = rb.end.next
		return
	}
	link := &eventLink{
		event: event,
		next:  rb.start,
	}
	rb.end.next = link
	rb.end = rb.end.next
}

// Pop removes the first event of the queue.
func (rb *ringBuffer) Pop() {
	if rb.start.event == nil {
		return
	}
	rb.start.event = nil
	rb.start = rb.start.next
}

// localEventQueue implements a local in-memory event queue.
type localEventQueue struct {
	buffer *ringBuffer
	Pushc  chan Event
	eventc chan Event
	loop   loop.Loop
}

// makeLocalEventQueueFactory creates a factory for local
// event queues.
func makeLocalEventQueueFactory(size int) EventQueueFactory {
	return func(env Environment) (EventQueue, error) {
		queue := &localEventQueue{
			buffer: newRingBuffer(size),
			Pushc:  make(chan Event),
			eventc: make(chan Event),
		}
		queue.loop = loop.Go(queue.backendLoop)
		return queue, nil
	}
}

// Push appends an event to the end of the queue.
func (q *localEventQueue) Push(event Event) error {
	select {
	case q.Pushc <- event:
	case <-q.loop.IsStopping():
		return errors.New(ErrStopping, errorMessages, "event queue")
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
		if q.buffer.Peek() == nil {
			// Empty buffer.
			select {
			case <-l.ShallStop():
				return nil
			case event := <-q.Pushc:
				q.buffer.Push(event)
			}
		} else {
			// At least one event in buffer.
			select {
			case <-l.ShallStop():
				return nil
			case event := <-q.Pushc:
				q.buffer.Push(event)
			case q.eventc <- q.buffer.Peek():
				q.buffer.Pop()
			}
		}
	}
}

// EOF
