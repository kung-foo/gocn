// Tideland Go Cell Network - Cells - Event
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
	"encoding/json"
	"fmt"

	"github.com/tideland/goas/v3/errors"
)

//--------------------
// ENCODE / DECODE
//--------------------

// encode encodes a payload into a slice of bytes.
func encode(payload interface{}) ([]byte, error) {
	encoded, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.Annotate(err, ErrEncoding, errorMessages)
	}
	return encoded, nil
}

// decode decodes encoded data into a payload.
func decode(encoded []byte, payload interface{}) error {
	if err := json.Unmarshal(encoded, payload); err != nil {
		return errors.Annotate(err, ErrDecoding, errorMessages)
	}
	return nil
}

//--------------------
// REPLY
//--------------------

// reply contains the response of a request.
type reply struct {
	encoded []byte
	err     error
}

// newReply creates a new response with the given payload and error.
func newReply(payload interface{}, rerr error) (*reply, error) {
	encoded, err := encode(payload)
	if err != nil {
		return nil, err
	}
	return &reply{encoded, rerr}, nil
}

// Payload returns the payload of the reply.
func (r *reply) Payload(payload interface{}) error {
	return decode(r.encoded, payload)
}

//--------------------
// EVENT
//--------------------

// Event transports what to process.
type Event interface {
	// Topic returns the topic of the event.
	Topic() string

	// Payload returns the payload of the event into
	// the passed interface.
	Payload(payload interface{}) error

	// IsRequest tells if the event is a request.
	IsRequest() bool

	// Respond allows the cell to respond with a payload
	// and/or an error.
	Respond(payload interface{}, err error) error

	// String returns a string representation of the event.
	String() string
}

// event implements the Event interface.
type event struct {
	topic   string
	encoded []byte
	replyc  chan *reply
}

// newEvent creates a new event with the given topic and payload.
func newEvent(topic string, payload interface{}) (*event, error) {
	if topic == "" {
		return nil, errors.New(ErrNoTopic, errorMessages)
	}
	encoded, err := encode(payload)
	if err != nil {
		return nil, err
	}
	return &event{topic, encoded, nil}, nil
}

// newEvent creates a new request event with the given topic and payload.
func newRequest(topic string, payload interface{}) (*event, error) {
	if topic == "" {
		return nil, errors.New(ErrNoTopic, errorMessages)
	}
	encoded, err := encode(payload)
	if err != nil {
		return nil, err
	}
	return &event{topic, encoded, make(chan *reply, 1)}, nil
}

// Topic returns the topic of the event.
func (e *event) Topic() string {
	return e.topic
}

// Payload returns the payload of the event.
func (e *event) Payload(payload interface{}) error {
	return decode(e.encoded, payload)
}

// IsRequest tells if the event is a request.
func (e *event) IsRequest() bool {
	return e.replyc != nil
}

// Respond allows the cell to respond with a payload and/or an error.
func (e *event) Respond(payload interface{}, err error) error {
	if e.replyc == nil {
		return errors.New(ErrNoRequest, errorMessages)
	}
	reply, err := newReply(payload, err)
	if err != nil {
		return err
	}
	e.replyc <- reply
	return nil
}

// String returns a string representation of the event.
func (e *event) String() string {
	var decoded interface{}
	if err := decode(e.encoded, &decoded); err != nil {
		decoded = err
	}
	if e.IsRequest() {
		return fmt.Sprintf("<request: %q / %v>", e.topic, decoded)
	}
	return fmt.Sprintf("<event: %q / %v>", e.topic, decoded)
}

// EOF
